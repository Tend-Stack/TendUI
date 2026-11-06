package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestElevationPreferencesAcrossDemosAndExports(t *testing.T) {
	h := newHandler()
	for _, level := range []string{"subtle", "floating"} {
		r := httptest.NewRequest("POST", "/settings", strings.NewReader("elevation="+level))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 303 {
			t.Fatal("cannot save elevation", w.Code)
		}
		cookies := w.Result().Cookies()
		for _, path := range []string{"/playground", "/settings", "/motion", "/appearance.css", "/appearance.go"} {
			r = httptest.NewRequest("GET", path, nil)
			for _, cookie := range cookies {
				r.AddCookie(cookie)
			}
			w = httptest.NewRecorder()
			h.ServeHTTP(w, r)
			want := `data-elevation="` + level + `"`
			if path == "/appearance.css" {
				want = "--gothui-shadow: var(--gothui-elevation-" + level + ")"
			} else if path == "/appearance.go" {
				want = `Elevation: "` + level + `"`
			}
			if w.Code != 200 || !strings.Contains(w.Body.String(), want) {
				t.Fatalf("%s lost saved %s elevation", path, level)
			}
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/motion?elevation=flat&theme=dark&motion=off", nil))
	body := w.Body.String()
	for _, want := range []string{`data-elevation="flat"`, `data-theme="dark"`, `data-motion="off"`, `name="elevation"`, `data-card-elevation="flat"`, `data-card-elevation="subtle"`, `data-card-elevation="raised"`, `data-card-elevation="floating"`} {
		if !strings.Contains(body, want) {
			t.Errorf("comparison demo missing %s", want)
		}
	}
	if strings.Contains(body, "<script") || strings.Contains(body, "style=") {
		t.Fatal("comparison demo must work without scripts or inline styles")
	}
}
