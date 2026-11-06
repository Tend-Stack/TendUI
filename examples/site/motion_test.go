package site

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Tend-Stack/TendUI/v2026/components"
)

func TestMotionStudioRoutesAndPreferences(t *testing.T) {
	h := newHandler()
	for _, path := range []string{"/motion", "/motion?motion=off&theme=dark&preset=rose", "/motion?motion=invalid"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		body := w.Body.String()
		if w.Code != 200 {
			t.Fatalf("%s: %d", path, w.Code)
		}
		for _, want := range []string{`data-surface="gradient"`, `data-enter="rise"`, `class="skeleton"`, `class="empty-state"`, `class="stat"`, `name="motion"`} {
			if !strings.Contains(body, want) {
				t.Errorf("missing %s", want)
			}
		}
		motion := "system"
		if strings.Contains(path, "motion=off") {
			motion = "off"
		}
		if !strings.Contains(body, `data-motion="`+motion+`"`) || strings.Contains(body, "<script") || strings.Contains(body, "style=") {
			t.Fatal("motion preference or no-script contract failed")
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("HEAD", "/motion", nil))
	if w.Code != 200 || w.Body.Len() != 0 {
		t.Fatal("HEAD must not have a body")
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("POST", "/motion", nil))
	if w.Code != 405 {
		t.Fatal("preview must not accept POST")
	}
}

func TestMotionPreferenceSaveAndExports(t *testing.T) {
	h := newHandler()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/settings", strings.NewReader("motion=off"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(w, r)
	if w.Code != 303 {
		t.Fatal(w.Code)
	}
	page := httptest.NewRequest("GET", "/motion", nil)
	for _, cookie := range w.Result().Cookies() {
		page.AddCookie(cookie)
	}
	if appearance(page).Motion != components.MotionOff {
		t.Fatal("motion was not saved")
	}
	q := url.Values{"motion": {"system"}}
	if appearanceFromValues(appearance(page), q).Motion != components.MotionSystem {
		t.Fatal("preview did not override cookie")
	}
	if appearanceValues(appearance(page)).Get("motion") != "off" {
		t.Fatal("motion lost from export URL")
	}
	for _, tc := range []struct{ path, want string }{
		{"/appearance.css?motion=off", "--tendui-duration-enter: 0ms"},
		{"/appearance.go?motion=off", `Motion: "off"`},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != 200 || !strings.Contains(w.Body.String(), tc.want) {
			t.Fatalf("bad export %s: %s", tc.path, w.Body.String())
		}
	}
}
