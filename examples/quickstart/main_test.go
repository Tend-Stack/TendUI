package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Tend-Stack/TendUI/v2026"
)

func TestQuickstartServesPageAndEmbeddedCSS(t *testing.T) {
	h := handler()
	for _, tc := range []struct{ path, contentType, text string }{
		{"/", "text/html", "Hello, TendUI."},
		{tendui.StylesheetURL, "text/css", "--tendui-primary"},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != 200 || !strings.HasPrefix(w.Header().Get("Content-Type"), tc.contentType) || !strings.Contains(w.Body.String(), tc.text) {
			t.Fatalf("%s: unexpected quickstart response: %d", tc.path, w.Code)
		}
		if strings.Contains(w.Body.String(), "<script") {
			t.Fatal("quickstart must work without scripts")
		}
		if tc.path == "/" && !strings.Contains(w.Body.String(), tendui.StylesheetURL) {
			t.Fatal("page must link the embedded stylesheet")
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/missing", nil))
	if w.Code != 404 {
		t.Fatal("unknown paths must not render the home page")
	}
}
