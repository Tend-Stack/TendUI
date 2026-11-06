package site

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestDemoPages(t *testing.T) {
	h := newHandler()
	for _, path := range []string{"/", "/playground", "/forms", "/dashboard"} {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
			if w.Code != 200 {
				t.Fatalf("status %d", w.Code)
			}
			html := w.Body.String()
			for _, want := range []string{`id="main"`, `href="/playground"`, `href="/forms"`, `href="/dashboard"`, `/tendui/static/tendui.css`} {
				if !strings.Contains(html, want) {
					t.Errorf("missing %s", want)
				}
			}
			if strings.Contains(html, "<script") || strings.Contains(html, "onclick=") {
				t.Fatal("unexpected JavaScript")
			}
			if !strings.Contains(w.Header().Get("Content-Security-Policy"), "default-src 'none'") {
				t.Fatal("missing restrictive CSP")
			}
		})
	}
}

func TestPlaygroundStates(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct{ query, want string }{
		{"state=disabled", `name="preview" value="A little less setup." aria-describedby="preview-hint" disabled`},
		{"state=readonly", `readonly`},
		{"state=invalid", `aria-invalid="true"`},
		{"theme=dark", `data-theme="dark"`},
		{"theme=light", `data-theme="light"`},
		{"variant=outline&clicked=yes", `Button submitted successfully.`},
		{"label=%3Cscript%3E", `&lt;script&gt;`},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/playground?"+tc.query, nil))
		if w.Code != 200 || !strings.Contains(w.Body.String(), tc.want) {
			t.Errorf("query %s: missing %s", tc.query, tc.want)
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/playground?apply=1&theme=system", nil))
	for _, unwanted := range []string{`data-theme=`, `Your preview card`, `Optional footer · Ready for your content`} {
		if strings.Contains(w.Body.String(), unwanted) {
			t.Errorf("unexpected %s", unwanted)
		}
	}
}

func TestFormValidationAndEscaping(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct {
		name   string
		values url.Values
		code   int
		want   string
	}{
		{"empty", url.Values{}, 422, "Please check the highlighted fields."},
		{"bad-email", url.Values{"name": {"Alex"}, "email": {"not-an-email"}, "role": {"developer"}, "consent": {"on"}}, 422, "Enter an email address"},
		{"success", url.Values{"name": {"<script>Alex</script>"}, "email": {"alex@example.com"}, "role": {"developer"}, "consent": {"on"}, "message": {"<b>Hello</b>"}}, 200, "&lt;script&gt;Alex&lt;/script&gt;"},
		{"invalid-role", url.Values{"name": {"Alex"}, "email": {"alex@example.com"}, "role": {"admin"}, "consent": {"on"}}, 422, "Choose what best describes you."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/forms", strings.NewReader(tc.values.Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.code || !strings.Contains(w.Body.String(), tc.want) {
				t.Fatalf("status %d, expected %d containing %s", w.Code, tc.code, tc.want)
			}
			if strings.Contains(w.Body.String(), "<script>") || strings.Contains(w.Body.String(), "<b>Hello</b>") {
				t.Fatal("unescaped input")
			}
		})
	}
	r := httptest.NewRequest("POST", "/forms", strings.NewReader("message="+strings.Repeat("x", 17000)))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized body: %d", w.Code)
	}
	r = httptest.NewRequest("POST", "https://example.com/forms", nil)
	r.Header.Set("Origin", "https://other.example")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatalf("cross-origin form: %d", w.Code)
	}
}

func TestDashboardFilters(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct{ query, want, absent string }{
		{"?q=portal&status=active", "Customer portal", "Studio website"},
		{"?status=paused", "Billing integration", "Customer portal"},
		{"?q=does-not-exist", "No matching projects.", `name="project"`},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/dashboard"+tc.query, nil))
		html := w.Body.String()
		if w.Code != 200 || !strings.Contains(html, tc.want) || strings.Contains(html, tc.absent) {
			t.Errorf("unexpected filter result: %s", tc.query)
		}
	}
}

func TestThemeReturnsToDemo(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct{ ref, location, cookie string }{
		{"https://example.com/playground?theme=dark&state=disabled", "/playground?state=disabled", "light"},
		{"https://example.com/dashboard?status=active", "/dashboard?status=active", "dark"},
		{"https://evil.example/forms", "/", "dark"},
		{"https://example.com/unknown", "/", "dark"},
	} {
		r := httptest.NewRequest("POST", "https://example.com/api/theme/toggle", nil)
		r.Header.Set("Referer", tc.ref)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 303 || w.Header().Get("Location") != tc.location || w.Result().Cookies()[0].Value != tc.cookie {
			t.Errorf("unexpected theme response for %s", tc.ref)
		}
	}
	r := httptest.NewRequest("POST", "https://example.com/api/theme/toggle", nil)
	r.Header.Set("Referer", "https://example.com/playground?theme=system")
	r.AddCookie(&http.Cookie{Name: "tendui_theme", Value: "dark"})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 303 || w.Result().Cookies()[0].Value != "light" {
		t.Fatal("system preview should toggle from the saved theme")
	}
}
