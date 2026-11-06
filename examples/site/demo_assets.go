package site

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"net/http"
	"time"
)

// Example-only presentation stays out of the library's CSS payload.
//
//go:embed demo.css
var demoCSS []byte

var demoStylePath = fmt.Sprintf("/demo-assets/%x.css", sha256.Sum256(demoCSS))

//go:embed demo-gopher.svg
var demoGopher []byte

func serveDemoGopher(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
	http.ServeContent(w, r, "demo-gopher.svg", time.Time{}, bytes.NewReader(demoGopher))
}

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
