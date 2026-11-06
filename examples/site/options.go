package site

import (
	"context"
	"net/http"

	"github.com/Tend-Stack/TendUI/v2026/components"
	"github.com/a-h/templ"
)

// Options customizes the example site when another site mounts it.
type Options struct {
	// Header renders the header of every demo page. active is the demo path,
	// such as "/catalog" or "/forms". Nil uses the built-in SiteHeader.
	Header func(active string) templ.Component
	// DefaultTheme applies when neither a theme query nor the tendui_theme
	// cookie chooses one. Empty keeps the light default.
	DefaultTheme components.Theme
	// HostPage reports whether the host serves path, so the theme toggle may
	// return there. Nil allows only demo pages.
	HostPage func(path string) bool
	// Head adds trusted elements to every demo page head, such as the host's
	// header stylesheet. It renders after the demo stylesheet.
	Head templ.Component
}

type optionsKey struct{}

// HandlerWithOptions is Handler with a host-supplied header and default theme.
func HandlerWithOptions(options Options) http.Handler {
	handler := newHandler()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), optionsKey{}, options)))
	})
}

func siteOptions(ctx context.Context) Options {
	options, _ := ctx.Value(optionsKey{}).(Options)
	return options
}

var exampleGroups = []struct {
	Title string
	Pages []components.NavItem
}{
	{"Forms & workflows", []components.NavItem{{Label: "Form lab", Href: "/forms"}, {Label: "Inline editing", Href: "/editing"}, {Label: "Form wizard", Href: "/wizard"}, {Label: "CSV preview", Href: "/uploads"}}},
	{"Application layouts", []components.NavItem{{Label: "Dashboard", Href: "/dashboard"}, {Label: "Workspace", Href: "/workspace"}, {Label: "Composition patterns", Href: "/patterns"}}},
	{"Style & behavior", []components.NavItem{{Label: "Appearance studio", Href: "/settings"}, {Label: "Native interactions", Href: "/interactions"}, {Label: "Motion & surfaces", Href: "/motion"}}},
}

// ExampleNavigation returns the grouped application examples, marking active.
func ExampleNavigation(active string) []components.NavGroup {
	groups := make([]components.NavGroup, 0, len(exampleGroups))
	for _, group := range exampleGroups {
		items := make([]components.NavItem, len(group.Pages))
		for i, page := range group.Pages {
			page.Current = page.Href == active
			items[i] = page
		}
		groups = append(groups, components.NavGroup{Title: group.Title, Items: items})
	}
	return groups
}

// IsExamplePath reports whether path is one of the grouped application examples.
func IsExamplePath(path string) bool {
	for _, group := range exampleGroups {
		for _, page := range group.Pages {
			if page.Href == path {
				return true
			}
		}
	}
	return false
}

func demoNavigation(active string) []components.SiteHeaderItem {
	return []components.SiteHeaderItem{
		{Label: "Overview", Href: "/", Current: active == "/"},
		{Label: "Components", Href: "/catalog", Current: active == "/catalog"},
		{Label: "Examples", Current: IsExamplePath(active), Groups: ExampleNavigation(active), Note: "Runnable pages with native controls and Go handlers."},
		{Label: "Playground", Href: "/playground", Current: active == "/playground"},
	}
}
