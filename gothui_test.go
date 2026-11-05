package gothui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAssets(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(AssetPath, AssetHandler())
	request := func(method, path, etag string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, path, nil)
		r.Header.Set("If-None-Match", etag)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}
	w := request("GET", StylesheetURL, "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), ":root") || w.Header().Get("Content-Type") != "text/css; charset=utf-8" {
		t.Fatalf("CSS response: %d %v", w.Code, w.Header())
	}
	for _, tc := range []struct {
		method, path, etag string
		code               int
	}{
		{"HEAD", StylesheetURL, "", 200},
		{"GET", StylesheetURL, w.Header().Get("ETag"), 304},
		{"POST", StylesheetURL, "", 405},
		{"GET", AssetPath, "", 404},
		{"GET", AssetPath + "missing.css", "", 404},
		{"GET", AssetPath + "static/gothui.css", "", 404},
	} {
		t.Run(tc.method+tc.path+tc.etag, func(t *testing.T) {
			t.Parallel()
			got := request(tc.method, tc.path, tc.etag)
			if got.Code != tc.code {
				t.Fatalf("status %d, want %d", got.Code, tc.code)
			}
			if (tc.method == "HEAD" || tc.code == 304) && got.Body.Len() != 0 {
				t.Fatal("unexpected response body")
			}
		})
	}
}
