package site

import (
	"encoding/base64"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAppearanceCookieRoundTrip(t *testing.T) {
	h := newHandler()
	values := url.Values{"preset": {"forest"}, "density": {"compact"}, "corners": {"soft"}, "typography": {"serif"}, "text_size": {"large"}, "width": {"wide"}, "elevation": {"flat"}, "theme": {"system"}}
	r := httptest.NewRequest("POST", "https://example.com/settings", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 303 || w.Header().Get("Location") != "/settings?saved=1" {
		t.Fatal("settings not redirected")
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 2 {
		t.Fatal("missing preferences")
	}
	for _, cookie := range cookies {
		if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteStrictMode {
			t.Fatal("cookie attributes")
		}
	}
	for _, path := range []string{"/settings", "/workspace", "/playground", "/forms", "/dashboard"} {
		r = httptest.NewRequest("GET", path, nil)
		for _, c := range cookies {
			r.AddCookie(c)
		}
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		html := w.Body.String()
		for _, want := range []string{`data-preset="forest"`, `data-density="compact"`, `data-corners="soft"`, `data-typography="serif"`, `data-text-size="large"`, `data-width="wide"`, `data-elevation="flat"`} {
			if !strings.Contains(html, want) {
				t.Errorf("%s missing %s", path, want)
			}
		}
		if w.Code != 200 || strings.Contains(html, `data-theme=`) || strings.Contains(html, "<script") {
			t.Errorf("incorrect system theme or script on %s", path)
		}
	}
	r = httptest.NewRequest("GET", "/workspace?preset=rose&theme=dark", nil)
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), `data-preset="rose"`) || !strings.Contains(w.Body.String(), `data-theme="dark"`) {
		t.Fatal("preview should override cookie")
	}
	r = httptest.NewRequest("POST", "/settings", strings.NewReader("action=reset"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	for _, c := range w.Result().Cookies() {
		if c.MaxAge != -1 {
			t.Fatal("reset did not clear cookie")
		}
	}
}

func TestAppearanceExportsAndUntrustedPreferences(t *testing.T) {
	h := newHandler()
	for _, ext := range []string{"css", "go"} {
		r := httptest.NewRequest("GET", "/appearance."+ext+"?preset=rose&density=compact&theme=dark", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 200 || !strings.Contains(w.Header().Get("Content-Disposition"), "attachment;") {
			t.Fatal("missing download")
		}
		if ext == "css" {
			if !strings.Contains(w.Body.String(), "--tendui-preset-rose-dark") {
				t.Fatal("wrong preset export")
			}
		} else {
			if _, err := parser.ParseFile(token.NewFileSet(), "theme.go", w.Body.Bytes(), parser.AllErrors); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(w.Body.String(), `components.Theme("dark")`) {
				t.Fatal("mode missing")
			}
		}
		w = httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("HEAD", "/appearance."+ext, nil))
		if w.Code != 200 || w.Body.Len() != 0 {
			t.Fatal("HEAD returned a body")
		}
		w = httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("POST", "/appearance."+ext, nil))
		if w.Code != 405 {
			t.Fatal("download accepts POST")
		}
	}
	for _, value := range []string{"invalid-base64", base64.RawURLEncoding.EncodeToString([]byte(`{"preset":"</style><script>bad()</script>","width":"invalid"}`))} {
		r := httptest.NewRequest("GET", "/settings", nil)
		r.AddCookie(&http.Cookie{Name: appearanceCookie, Value: value})
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 200 || !strings.Contains(w.Body.String(), `data-preset="indigo"`) || strings.Contains(w.Body.String(), "<script>") {
			t.Fatal("untrusted cookie was not normalized")
		}
	}
	r := httptest.NewRequest("POST", "https://example.com/settings", nil)
	r.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("cross-origin settings were accepted")
	}
}
