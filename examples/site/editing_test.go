package site

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func editingRequest(method, path, body string, hx bool, cookie *http.Cookie) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if hx {
		r.Header.Set("HX-Request", "true")
	}
	if cookie != nil {
		r.AddCookie(cookie)
	}
	return r
}

func TestEditingRepresentations(t *testing.T) {
	h := newHandler()
	for _, edit := range []string{"", "?edit=1"} {
		for _, hx := range []bool{false, true} {
			for _, restore := range []bool{false, true} {
				r := editingRequest("GET", "/editing"+edit, "", hx, nil)
				if restore {
					r.Header.Set("HX-History-Restore-Request", "true")
				}
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				body := w.Body.String()
				if w.Code != 200 || strings.Contains(body, "<html") != (!hx || restore) || strings.Count(body, `id="profile-panel"`) != 1 {
					t.Fatalf("bad representation hx=%v restore=%v: %d", hx, restore, w.Code)
				}
				if strings.Contains(body, "<script") || !strings.Contains(w.Header().Get("Content-Security-Policy"), "default-src 'none'") {
					t.Fatal("script policy changed")
				}
				if edit != "" && (!strings.Contains(body, `hx-post="/editing"`) || !strings.Contains(body, `action="/editing"`) || !strings.Contains(body, `href="/editing"`)) {
					t.Fatal("progressive form hooks missing")
				}
			}
		}
	}
}

func TestEditingSaveCancelResetAndIsolation(t *testing.T) {
	h := newHandler()
	for _, hx := range []bool{false, true} {
		values := url.Values{"intent": {"save"}, "name": {"<Alex>"}, "bio": {"Hello & welcome"}}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, editingRequest("POST", "/editing", values.Encode(), hx, nil))
		want := 303
		if hx {
			want = 200
		}
		if w.Code != want {
			t.Fatalf("save status %d", w.Code)
		}
		cookies := w.Result().Cookies()
		if len(cookies) != 1 || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteStrictMode || cookies[0].MaxAge != 0 {
			t.Fatal("cookie policy")
		}
		if !hx && w.Header().Get("Location") != "/editing?notice=saved" {
			t.Fatal("native POST must redirect")
		}
		for _, path := range []string{"/editing", "/editing?edit=1"} {
			view := httptest.NewRecorder()
			h.ServeHTTP(view, editingRequest("GET", path, "", hx, cookies[0]))
			if !strings.Contains(view.Body.String(), "&lt;Alex&gt;") || strings.Contains(view.Body.String(), "<Alex>") {
				t.Fatal("saved values missing or unescaped")
			}
		}
		fresh := httptest.NewRecorder()
		h.ServeHTTP(fresh, editingRequest("GET", "/editing", "", false, nil))
		if !strings.Contains(fresh.Body.String(), "Alex Morgan") {
			t.Fatal("profile leaked into another session")
		}
		reset := httptest.NewRecorder()
		h.ServeHTTP(reset, editingRequest("POST", "/editing", "intent=reset", hx, cookies[0]))
		if reset.Code != want || len(reset.Result().Cookies()) != 1 || reset.Result().Cookies()[0].MaxAge != -1 {
			t.Fatal("reset did not clear cookie")
		}
		if hx && !strings.Contains(reset.Body.String(), "Alex Morgan") {
			t.Fatal("reset fragment stale")
		}
	}
}

func TestEditingValidationAndRequestLimits(t *testing.T) {
	h := newHandler()
	for _, hx := range []bool{false, true} {
		w := httptest.NewRecorder()
		values := url.Values{"intent": {"save"}, "name": {"  "}, "bio": {strings.Repeat("b", 241)}}
		h.ServeHTTP(w, editingRequest("POST", "/editing", values.Encode(), hx, nil))
		want := 422
		if hx {
			want = 200
		}
		if w.Code != want || len(w.Result().Cookies()) != 0 || strings.Count(w.Body.String(), `aria-invalid="true"`) != 2 || !strings.Contains(w.Body.String(), `href="#profile-name"`) {
			t.Fatalf("validation lost: %d", w.Code)
		}
	}
	for _, tc := range []struct {
		method, body string
		status       int
	}{
		{"PUT", "", 405}, {"POST", "intent=unknown", 400}, {"POST", "intent=save&name=" + strings.Repeat("x", 4096), 400}, {"POST", "intent=save&name=%ZZ", 400},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, editingRequest(tc.method, "/editing", tc.body, false, nil))
		if w.Code != tc.status || len(w.Result().Cookies()) != 0 {
			t.Fatalf("bad rejection: %d", w.Code)
		}
	}
	r := editingRequest("POST", "/editing", "intent=reset", true, nil)
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatalf("cross-origin POST accepted: %d", w.Code)
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, editingRequest("HEAD", "/editing?edit=1", "", true, nil))
	if w.Code != 200 || w.Body.Len() != 0 {
		t.Fatal("HEAD semantics")
	}
}

func TestProfileCookieBoundaries(t *testing.T) {
	for _, value := range []string{"bad!", strings.Repeat("x", 3073), base64.RawURLEncoding.EncodeToString([]byte(`{"name":""}`))} {
		r := editingRequest("GET", "/editing", "", false, &http.Cookie{Name: profileCookie, Value: value})
		if readProfile(r) != defaultProfile() {
			t.Fatal("invalid cookie accepted")
		}
	}
	// HTML-sensitive values expand during JSON encoding; all valid sizes must round-trip.
	p := demoProfile{Name: strings.Repeat("<", 60), Bio: strings.Repeat("&", 240)}
	b, _ := json.Marshal(p)
	r := editingRequest("GET", "/editing", "", false, &http.Cookie{Name: profileCookie, Value: base64.RawURLEncoding.EncodeToString(b)})
	if readProfile(r) != p {
		t.Fatal("valid maximum-size profile lost")
	}
	r = editingRequest("POST", "https://example.com/editing", "intent=save&name=Alex", true, nil)
	w := httptest.NewRecorder()
	newHandler().ServeHTTP(w, r)
	if len(w.Result().Cookies()) != 1 || !w.Result().Cookies()[0].Secure {
		t.Fatal("TLS cookie must be secure")
	}
}
