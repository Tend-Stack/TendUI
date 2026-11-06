// Run with: go run ./examples/basic
package main

import (
	"bytes"
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"time"

	"github.com/Tend-Stack/GothUI/v2026"
	"github.com/Tend-Stack/GothUI/v2026/components"
	"github.com/a-h/templ"
)

// Set with -ldflags "-X main.version=$(git describe --tags --long --dirty)".
var version = gothui.Version + "-dev"

func theme(r *http.Request) components.Theme {
	switch r.URL.Query().Get("theme") {
	case "light":
		return components.ThemeLight
	case "dark":
		return components.ThemeDark
	case "system":
		return components.ThemeSystem
	}
	cookie, err := r.Cookie("gothui_theme")
	if err == nil && cookie.Value == "system" {
		return components.ThemeSystem
	}
	if err == nil && cookie.Value == string(components.ThemeDark) {
		return components.ThemeDark
	}
	return components.ThemeLight
}

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(gothui.AssetPath, gothui.AssetHandler())
	mux.HandleFunc("/appearance.css", exportAppearance)
	mux.HandleFunc("/appearance.go", exportAppearance)
	mux.HandleFunc("/api/theme/toggle", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		current := theme(r)
		if ref := demoReferer(r); ref != nil {
			pageRequest := r.Clone(r.Context())
			pageRequest.URL = ref
			current = theme(pageRequest)
		}
		if current == components.ThemeSystem {
			current = theme(r)
		}
		next := components.ThemeDark
		if current == components.ThemeDark {
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
		http.Redirect(w, r, themeReturnPath(r), http.StatusSeeOther)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !demoPath(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		if r.URL.Path == "/settings" && r.Method == http.MethodPost {
			saveAppearance(w, r)
			return
		}
		if r.URL.Path == "/forms" && r.Method == http.MethodPost {
			submitForm(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			allowed := "GET, HEAD"
			if r.URL.Path == "/forms" || r.URL.Path == "/settings" {
				allowed += ", POST"
			}
			w.Header().Set("Allow", allowed)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body templ.Component
		title := "GothUI · HTML, beautifully enough"
		switch r.URL.Path {
		case "/playground":
			title = "Component playground · GothUI"
			options := readPlayground(r.URL.Query())
			options.Appearance = appearance(r)
			body = playground(options)
		case "/settings":
			title = "Appearance · GothUI"
			body = appearancePage(appearance(r), themeName(theme(r)), r.URL.Query().Get("saved") == "1")
		case "/workspace":
			title = "Workspace · GothUI"
			body = workspace()
		case "/forms":
			title = "Form lab · GothUI"
			body = formLab(formData{Errors: map[string]string{}})
		case "/dashboard":
			title = "Project dashboard · GothUI"
			body = dashboard(readDashboard(r.URL.Query()))
		default:
			body = showcase(r.URL.Query().Get("name"), version)
		}
		renderPage(w, r, title, body, http.StatusOK)
	})

	// Reject browser cross-origin unsafe requests. The strict cookie is additional protection.
	return http.NewCrossOriginProtection().Handler(mux)
}

func demoPath(path string) bool {
	switch path {
	case "/", "/playground", "/forms", "/dashboard", "/settings", "/workspace":
		return true
	}
	return false
}

// Only return to a known same-origin demo; never redirect to a supplied external host.
func themeReturnPath(r *http.Request) string {
	u := demoReferer(r)
	if u == nil {
		return "/"
	}
	q := u.Query()
	q.Del("theme") // The freshly saved cookie takes precedence after toggling.
	u.RawQuery = q.Encode()
	return u.RequestURI()
}

func demoReferer(r *http.Request) *url.URL {
	u, err := url.Parse(r.Referer())
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if err != nil || u.Host != r.Host || u.Scheme != scheme || !demoPath(u.Path) {
		return nil
	}
	return u
}

func renderPage(w http.ResponseWriter, r *http.Request, title string, body templ.Component, status int) {
	page := components.Layout(components.LayoutProps{
		Title: title, Language: "en", Theme: theme(r), Appearance: appearance(r),
		Description: "Interactive GothUI demos: semantic components, native forms, and layouts.",
	}, body)
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
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		_, _ = w.Write(buf.Bytes())
	}
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "local listen address")
	open := flag.Bool("open", false, "open the component playground in the default browser")
	flag.Parse()
	server := &http.Server{
		Addr: *addr, Handler: newHandler(),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second,
	}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}
	address := "http://" + listener.Addr().String()
	log.Printf("GothUI %s at %s/playground", version, address)
	if *open {
		go func() {
			if err := exec.Command("xdg-open", address+"/playground").Run(); err != nil {
				log.Printf("Open %s/playground in your browser (%v)", address, err)
			}
		}()
	}
	log.Fatal(server.Serve(listener))
}
