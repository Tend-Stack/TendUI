// A complete TendUI application. The generated page and CSS ship in its binary.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Tend-Stack/TendUI/v2026"
	"github.com/a-h/templ"
)

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(tendui.AssetPath, tendui.AssetHandler())
	mux.Handle("GET /{$}", templ.Handler(page()))
	return mux
}

func main() {
	server := &http.Server{
		Addr: "127.0.0.1:8080", Handler: handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
