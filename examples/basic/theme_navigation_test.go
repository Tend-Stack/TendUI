package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestThemeNavigationClearsPreviewAndAppliesCookie(t *testing.T) {
	h := newHandler()
	for _, hx := range []bool{false, true} {
		name := "native"
		if hx {
			name = "htmx"
		}
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "https://example.com/api/theme/toggle", nil)
			r.Header.Set("Referer", "https://example.com/playground?theme=dark&theme=light&state=disabled")
			if hx {
				r.Header.Set("HX-Request", "true")
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			location, status := w.Header().Get("Location"), 303
			if hx {
				location, status = w.Header().Get("HX-Redirect"), 200
				if w.Header().Get("HX-Refresh") != "" || w.Header().Get("Location") != "" || w.Body.Len() != 0 {
					t.Fatal("HTMX response must request only a navigation to the cleaned URL")
				}
			}
			if w.Code != status || location != "/playground?state=disabled" {
				t.Fatalf("preview still overrides cookie: status=%d location=%q", w.Code, location)
			}
			cookies := w.Result().Cookies()
			if len(cookies) != 1 || cookies[0].Value != "light" || !cookies[0].HttpOnly || !cookies[0].Secure || cookies[0].SameSite != http.SameSiteStrictMode {
				t.Fatal("incorrect saved theme or cookie policy", cookies)
			}
			// Follow the destination with the saved cookie. This verifies the HTTP
			// contract, not that an HTMX runtime performs the browser navigation.
			r = httptest.NewRequest("GET", "https://example.com"+location, nil)
			r.AddCookie(cookies[0])
			w = httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != 200 || !strings.Contains(w.Body.String(), `data-theme="light"`) || !strings.Contains(w.Body.String(), `disabled`) {
				t.Fatal("destination did not apply saved theme and retain other preview settings")
			}
		})
	}
}

func TestHTMXThemeValidatesCurrentPageURL(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct{ name, hx, current, ref, location, theme string }{
		{"current wins", "true", "https://example.com/motion?theme=dark&motion=off", "https://example.com/forms?theme=light", "/motion?motion=off", "light"},
		{"current without referer", "true", "https://example.com/dashboard?theme=dark&q=a%26b", "", "/dashboard?q=a%26b", "light"},
		{"external current", "true", "https://evil.example/motion?theme=dark", "https://example.com/forms", "/forms", "dark"},
		{"wrong scheme", "true", "http://example.com/motion?theme=dark", "", "/", "dark"},
		{"wrong port", "true", "https://example.com:8443/motion?theme=dark", "", "/", "dark"},
		{"userinfo", "true", "https://user@example.com/motion?theme=dark", "", "/", "dark"},
		{"relative", "true", "/motion?theme=dark", "", "/", "dark"},
		{"unknown page", "true", "https://example.com/api/theme/toggle?theme=dark", "", "/", "dark"},
		{"malformed", "true", "https://example.com/%zz", "", "/", "dark"},
		{"native ignores hx current", "", "https://example.com/motion?theme=dark", "https://example.com/forms", "/forms", "dark"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "https://example.com/api/theme/toggle", nil)
			r.Header.Set("HX-Request", tc.hx)
			r.Header.Set("HX-Current-URL", tc.current)
			r.Header.Set("Referer", tc.ref)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			location, status := w.Header().Get("Location"), 303
			if tc.hx == "true" {
				location, status = w.Header().Get("HX-Redirect"), 200
			}
			cookies := w.Result().Cookies()
			if w.Code != status || location != tc.location || len(cookies) != 1 || cookies[0].Value != tc.theme {
				t.Fatalf("unexpected response: status=%d location=%q cookies=%v", w.Code, location, cookies)
			}
		})
	}
	r := httptest.NewRequest("POST", "https://example.com/api/theme/toggle", nil)
	r.Header.Set("HX-Request", "true")
	r.Header.Set("HX-Current-URL", "https://example.com/motion")
	r.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden || len(w.Result().Cookies()) != 0 || w.Header().Get("HX-Redirect") != "" {
		t.Fatal("presentation headers bypassed cross-origin protection")
	}
}
