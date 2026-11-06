package site

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Tend-Stack/TendUI/v2026/components"
	"github.com/a-h/templ"
)

func TestHandlerWithOptionsDefaultThemeAndHeader(t *testing.T) {
	var seen []string
	h := HandlerWithOptions(Options{
		DefaultTheme: components.ThemeDark,
		Head:         templ.Raw(`<link rel="stylesheet" href="/host.css">`),
		Header: func(active string) templ.Component {
			return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				seen = append(seen, active)
				_, err := io.WriteString(w, `<header class="host-header">Host</header>`)
				return err
			})
		},
	})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/catalog", nil))
	body := w.Body.String()
	if w.Code != 200 || !strings.Contains(body, `data-theme="dark"`) || !strings.Contains(body, "host-header") {
		t.Fatalf("host options not applied: %d", w.Code)
	}
	if strings.Contains(body, `class="site-header"`) || len(seen) != 1 || seen[0] != "/catalog" {
		t.Fatalf("built-in header rendered or wrong active path: %v", seen)
	}
	if !strings.Contains(body, `href="/host.css"`) {
		t.Fatal("host head elements missing")
	}
	if !strings.Contains(body, `class="skip-link"`) {
		t.Fatal("skip link must remain with a host header")
	}

	r := httptest.NewRequest("GET", "/catalog", nil)
	r.AddCookie(&http.Cookie{Name: "tendui_theme", Value: "light"})
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), `data-theme="light"`) {
		t.Fatal("a saved light theme must override the default")
	}

	// With no cookie the default is dark, so the first toggle saves light.
	r = httptest.NewRequest("POST", "/api/theme/toggle", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if cookies := w.Result().Cookies(); len(cookies) != 1 || cookies[0].Value != "light" {
		t.Fatal("first toggle should switch away from the default")
	}
}

func TestHandlerWithOptionsHostReturnPath(t *testing.T) {
	h := HandlerWithOptions(Options{HostPage: func(path string) bool { return strings.HasPrefix(path, "/docs") }})
	for _, tc := range []struct{ ref, location string }{
		{"https://example.com/docs/install", "/docs/install"},
		{"https://example.com/private", "/"},
		{"https://evil.example/docs", "/"},
	} {
		r := httptest.NewRequest("POST", "https://example.com/api/theme/toggle", nil)
		r.Header.Set("Referer", tc.ref)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 303 || w.Header().Get("Location") != tc.location {
			t.Errorf("%s: got %q", tc.ref, w.Header().Get("Location"))
		}
	}
}

func TestBuiltInDemoHeader(t *testing.T) {
	w := httptest.NewRecorder()
	newHandler().ServeHTTP(w, httptest.NewRequest("GET", "/forms", nil))
	body := w.Body.String()
	for _, want := range []string{`class="site-header"`, `<summary data-current="true">`, `href="/forms" aria-current="page"`, `data-icon-only="true"`, `aria-label="Toggle theme"`} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(body, "◐") || strings.Contains(body, "<script") {
		t.Fatal("unexpected glyph or script")
	}
	if !IsExamplePath("/forms") || IsExamplePath("/catalog") {
		t.Fatal("example path classification")
	}
}
