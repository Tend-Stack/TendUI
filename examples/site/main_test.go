package site

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestThemeRoundTrip(t *testing.T) {
	h := newHandler()
	var cookie *http.Cookie
	for _, want := range []string{"dark", "light"} {
		r := httptest.NewRequest("POST", "/api/theme/toggle", nil)
		if cookie != nil {
			r.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 303 || w.Header().Get("Location") != "/" {
			t.Fatalf("redirect: %d", w.Code)
		}
		cookies := w.Result().Cookies()
		if len(cookies) != 1 {
			t.Fatal("missing cookie")
		}
		cookie = cookies[0]
		if cookie.Value != want || !cookie.HttpOnly || cookie.SameSite != http.SameSiteStrictMode {
			t.Fatalf("cookie: %+v", cookie)
		}
		r = httptest.NewRequest("GET", "/", nil)
		r.AddCookie(cookie)
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 200 || !strings.Contains(w.Body.String(), `data-theme="`+want+`"`) {
			t.Fatal("theme not applied")
		}
		if strings.Contains(w.Body.String(), "<script") {
			t.Fatal("unexpected JavaScript")
		}
	}
}

func TestThemeHTMXAndRequestBoundaries(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct {
		name, method, path, hx, origin string
		code                           int
	}{
		{"htmx", "POST", "/api/theme/toggle", "true", "", 200},
		{"method", "GET", "/api/theme/toggle", "", "", 405},
		{"cross-origin", "POST", "/api/theme/toggle", "", "https://evil.example", 403},
		{"missing", "GET", "/missing", "", "", 404},
		{"page-method", "POST", "/", "", "", 405},
		{"head", "HEAD", "/", "", "", 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "https://example.com"+tc.path, nil)
			r.Header.Set("HX-Request", tc.hx)
			r.Header.Set("Origin", tc.origin)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.code {
				t.Fatalf("got %d, want %d", w.Code, tc.code)
			}
			if tc.hx != "" && (w.Header().Get("HX-Redirect") != "/" || !w.Result().Cookies()[0].Secure) {
				t.Fatal("HTMX navigation or secure cookie missing")
			}
			if tc.method == "HEAD" && w.Body.Len() != 0 {
				t.Fatal("HEAD body emitted")
			}
		})
	}
}
