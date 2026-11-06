// Package site provides the runnable TendUI catalog and application examples.
// Its routes and demo workflows are examples, not a stable component API.
package site

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Tend-Stack/TendUI/v2026"
	"github.com/Tend-Stack/TendUI/v2026/components"
	"github.com/Tend-Stack/TendUI/v2026/render"
	"github.com/a-h/templ"
)

const version = tendui.Version

// Handler returns an independent HTTP handler for the catalog and native demos.
// Mount at the origin root; demo links and form actions use absolute paths.
func Handler() http.Handler { return newHandler() }

func theme(r *http.Request) components.Theme {
	switch r.URL.Query().Get("theme") {
	case "light":
		return components.ThemeLight
	case "dark":
		return components.ThemeDark
	case "system":
		return components.ThemeSystem
	}
	cookie, err := r.Cookie("tendui_theme")
	if err == nil && cookie.Value == "system" {
		return components.ThemeSystem
	}
	if err == nil && cookie.Value == string(components.ThemeDark) {
		return components.ThemeDark
	}
	if err == nil && cookie.Value == string(components.ThemeLight) {
		return components.ThemeLight
	}
	if fallback := siteOptions(r.Context()).DefaultTheme; fallback != "" {
		return fallback
	}
	return components.ThemeLight
}

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(tendui.AssetPath, tendui.AssetHandler())
	mux.HandleFunc("/appearance.css", exportAppearance)
	mux.HandleFunc("/appearance.go", exportAppearance)
	mux.HandleFunc("/interactions/confirm", confirmDemo)
	mux.HandleFunc("/editing", editingDemo)
	mux.HandleFunc("/wizard", projectWizard)
	mux.HandleFunc("/uploads", uploadDemo)
	mux.HandleFunc("/uploads/sample.csv", downloadProjectSample)
	mux.HandleFunc("/demo.css", serveDemoCSS)
	mux.HandleFunc(demoStylePath, serveDemoCSS)
	mux.HandleFunc("/demo-gopher.svg", serveDemoGopher)
	mux.HandleFunc("/catalog", serveCatalog)
	mux.HandleFunc("/catalog/", serveCatalog)
	mux.HandleFunc("/dashboard/selection", reviewSelection)
	mux.HandleFunc("/api/theme/toggle", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		current := theme(r)
		if ref := themePageURL(r); ref != nil {
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
			Name: "tendui_theme", Value: string(next), Path: "/",
			MaxAge: 60 * 60 * 24 * 365, HttpOnly: true,
			Secure: r.TLS != nil, SameSite: http.SameSiteStrictMode,
		})
		if r.Header.Get("HX-Request") == "true" {
			// A real navigation updates the root attribute and head reliably.
			// Use a non-redirect status: HTMX does not process headers on 3xx.
			// Refreshing the current URL would keep a theme preview that overrides
			// the new cookie. Navigate to the same cleaned destination as native POST.
			w.Header().Set("HX-Redirect", themeReturnPath(r))
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
		title := "TendUI · HTML, beautifully enough"
		switch r.URL.Path {
		case "/playground":
			title = "Component playground · TendUI"
			options := readPlayground(r.URL.Query())
			options.Appearance = appearance(r)
			body = playground(options)
		case "/patterns":
			title = "Application patterns · TendUI"
			body = patterns(readPatterns(r.URL.Query()))
		case "/motion":
			title = "Motion & surfaces · TendUI"
			body = motionStudio(appearance(r), themeName(theme(r)))
		case "/interactions":
			title = "Native interactions · TendUI"
			body = interactions(r.URL.Query().Get("confirmed") == "1", r.URL.Query().Get("view"))
		case "/settings":
			title = "Appearance · TendUI"
			body = appearancePage(appearance(r), themeName(theme(r)), r.URL.Query().Get("saved") == "1")
		case "/workspace":
			title = "Workspace · TendUI"
			body = workspace()
		case "/forms":
			title = "Form lab · TendUI"
			body = formLab(formData{Errors: map[string]string{}})
		case "/dashboard":
			title = "Project dashboard · TendUI"
			body = dashboard(readDashboard(r.URL.Query()))
		default:
			body = showcase(r.URL.Query().Get("name"), version)
		}
		renderPage(w, r, title, body, http.StatusOK)
	})

	// Reject browser cross-origin unsafe requests. The strict cookie is additional protection.
	return http.NewCrossOriginProtection().Handler(mux)
}

func returnPath(r *http.Request, path string) bool {
	if demoPath(path) {
		return true
	}
	host := siteOptions(r.Context()).HostPage
	return host != nil && host(path)
}

func demoPath(path string) bool {
	switch path {
	case "/", "/playground", "/forms", "/dashboard", "/settings", "/workspace", "/interactions", "/editing", "/motion", "/patterns", "/wizard", "/catalog", "/uploads":
		return true
	}
	return catalogPagePath(path)
}

// Only return to a known same-origin demo; never redirect to a supplied external host.
func themeReturnPath(r *http.Request) string {
	u := themePageURL(r)
	if u == nil {
		return "/"
	}
	q := u.Query()
	q.Del("theme") // The freshly saved cookie takes precedence after toggling.
	u.RawQuery = q.Encode()
	return u.RequestURI()
}

// HTMX's current URL follows client-side history changes. Treat it as untrusted
// presentation input, validate it like Referer, and never use it for authorization.
func themePageURL(r *http.Request) *url.URL {
	if r.Header.Get("HX-Request") == "true" {
		if u := sameOriginDemoURL(r, r.Header.Get("HX-Current-URL")); u != nil {
			return u
		}
	}
	return sameOriginDemoURL(r, r.Referer())
}

func sameOriginDemoURL(r *http.Request, value string) *url.URL {
	u, err := url.Parse(value)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if err != nil || u.User != nil || u.Host != r.Host || u.Scheme != scheme || !returnPath(r, u.Path) {
		return nil
	}
	return u
}

func renderPage(w http.ResponseWriter, r *http.Request, title string, body templ.Component, status int) {
	renderResponse(w, r, title, body, nil, status)
}

func renderResponse(w http.ResponseWriter, r *http.Request, title string, body, fragment templ.Component, status int) {
	page := components.Layout(components.LayoutProps{
		Title: title, Language: "en", Theme: theme(r), Appearance: appearance(r),
		Description: "Interactive TendUI demos: semantic components, native forms, and layouts.",
		Head:        demoHead(),
	}, body)
	htmlHeaders(w)
	if err := render.HTML(w, r, status, page, fragment); err != nil {
		log.Printf("render: %v", err)
	}
}

func htmlHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; img-src 'self'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'")
}
