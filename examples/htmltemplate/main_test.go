package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Tend-Stack/TendUI/v2026"
)

func TestPageAndEmbeddedStylesheet(t *testing.T) {
	h := handler()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK || w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("page status/type: %d %q", w.Code, w.Header().Get("Content-Type"))
	}
	for _, want := range []string{"Hello, Go developer.", `href="` + tendui.StylesheetURL + `"`, `action="/" method="get"`, `for="name"`} {
		if !strings.Contains(w.Body.String(), want) {
			t.Errorf("missing page content: %s", want)
		}
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tendui.StylesheetURL, nil))
	if w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/css") || !strings.Contains(w.Body.String(), "--tendui-primary") {
		t.Fatalf("stylesheet unavailable: status %d, type %q", w.Code, w.Header().Get("Content-Type"))
	}
}

func TestNameIsEscapedInTextAndAttributeContexts(t *testing.T) {
	name := `"><script>alert('x')</script><input autofocus onfocus="alert(1)">`
	w := httptest.NewRecorder()
	handler().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/?name="+url.QueryEscape(name), nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	body := w.Body.String()
	if strings.Contains(body, "<script") || strings.Count(body, "<input") != 1 || strings.Contains(body, `onfocus="`) {
		t.Fatal("user input escaped its text or attribute context")
	}
	for _, want := range []string{`<h1>Hello, &#34;&gt;&lt;script&gt;`, `value="&#34;&gt;&lt;script&gt;`, `&lt;/script&gt;`} {
		if !strings.Contains(body, want) {
			t.Errorf("missing escaped user text: %s", want)
		}
	}
}

func TestRouteAndMethodBoundaries(t *testing.T) {
	h := handler()
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{http.MethodGet, "/missing", http.StatusNotFound},
		{http.MethodPost, "/", http.StatusMethodNotAllowed},
		{http.MethodGet, tendui.AssetPath + "missing.css", http.StatusNotFound},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.status {
			t.Errorf("%s %s: got %d, want %d", tc.method, tc.path, w.Code, tc.status)
		}
	}
}
