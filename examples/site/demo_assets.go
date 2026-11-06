package site

import (
	"bytes"
	_ "embed"
	"net/http"
	"time"
)

// Example-only presentation stays out of the library's CSS payload.
//
//go:embed demo.css
var demoCSS []byte

func serveDemoCSS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")
	http.ServeContent(w, r, "demo.css", time.Time{}, bytes.NewReader(demoCSS))
}
