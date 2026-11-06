// Package tendui serves embedded, semantic CSS without external assets or build tools.
package tendui

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"fmt"
	"net/http"
	"time"
)

const (
	// AssetPath is the default mount point for AssetHandler.
	AssetPath = "/tendui/static/"
	// StylesheetURL is the stylesheet used by components.Layout.
	StylesheetURL = AssetPath + "tendui.css"
	// Version identifies the current release line; development binaries should use git describe.
	Version = "2026.2.0-preview.1"
)

//go:embed static
var assets embed.FS

// AssetHandler serves the embedded stylesheet at StylesheetURL. Mount directly:
// mux.Handle(tendui.AssetPath, tendui.AssetHandler()). No StripPrefix is needed.
// Each request has its own reader, making a shared handler safe for concurrent use.
// Unknown assets and directory listings return 404; only GET and HEAD are allowed.
func AssetHandler() http.Handler {
	css, err := assets.ReadFile("static/tendui.css")
	if err != nil {
		panic("tendui: embedded stylesheet missing: " + err.Error())
	}
	etag := fmt.Sprintf(`"%x"`, sha256.Sum256(css))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != StylesheetURL {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
		w.Header().Set("ETag", etag)
		http.ServeContent(w, r, "tendui.css", time.Time{}, bytes.NewReader(css))
	})
}
