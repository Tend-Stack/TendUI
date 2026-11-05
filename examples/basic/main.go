// Run with: go run ./examples/basic
package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/Tend-Stack/GothUI/v2026"
	"github.com/Tend-Stack/GothUI/v2026/components"
)

// Set with -ldflags "-X main.version=$(git describe --tags --long --dirty)".
var version = gothui.Version + "-dev"

func theme(r *http.Request) components.Theme {
	cookie, err := r.Cookie("gothui_theme")
	if err == nil && cookie.Value == string(components.ThemeDark) {
		return components.ThemeDark
	}
	return components.ThemeLight
}

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(gothui.AssetPath, gothui.AssetHandler())
	mux.HandleFunc("/api/theme/toggle", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next := components.ThemeDark
		if theme(r) == components.ThemeDark {
			next = components.ThemeLight
		}
		http.SetCookie(w, &http.Cookie{
			Name: "gothui_theme", Value: string(next), Path: "/",
			MaxAge: 60 * 60 * 24 * 365, HttpOnly: true,
			Secure: r.TLS != nil, SameSite: http.SameSiteStrictMode,
		})
		if r.Header.Get("HX-Request") == "true" {
			// A real navigation updates the root attribute and head reliably.
			// Use a non-redirect status: HTMX does not process headers on 3xx.
			w.Header().Set("HX-Refresh", "true")
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		page := components.Layout(components.LayoutProps{
			Title: "GothUI · HTML, beautifully enough", Language: "en", Theme: theme(r),
			Description: "Semantic CSS and typed Go components. No frontend build pipeline.",
		}, showcase(r.URL.Query().Get("name"), version))
		var buf bytes.Buffer
		if err := page.Render(r.Context(), &buf); err != nil {
			log.Printf("render: %v", err)
			http.Error(w, "unable to render page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'")
		if r.Method != http.MethodHead {
			_, _ = w.Write(buf.Bytes())
		}
	})
	// Reject browser cross-origin unsafe requests. The strict cookie is additional protection.
	return http.NewCrossOriginProtection().Handler(mux)
}

func main() {
	server := &http.Server{
		Addr: "127.0.0.1:8080", Handler: newHandler(),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second,
	}
	log.Printf("GothUI %s at http://%s", version, server.Addr)
	log.Fatal(server.ListenAndServe())
}
