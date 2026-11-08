// A standard-library HTML application with TendUI's CSS embedded in its binary.
package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Tend-Stack/TendUI/v2026"
)

//go:embed page.html
var pageHTML string

var page = template.Must(template.New("page").Parse(pageHTML))

type pageData struct {
	Name          string
	StylesheetURL string
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(tendui.AssetPath, tendui.AssetHandler())
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "Go developer"
		}
		// Render before writing headers so an error cannot send a partial page.
		var body bytes.Buffer
		if err := page.Execute(&body, pageData{Name: name, StylesheetURL: tendui.StylesheetURL}); err != nil {
			http.Error(w, "Unable to render page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(body.Bytes())
	})
	return mux
}

func main() {
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
