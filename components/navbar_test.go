package components

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestSiteHeader(t *testing.T) {
	html := render(t, SiteHeader(SiteHeaderProps{
		Brand:      templ.Raw("<span>Brand</span>"),
		BrandLabel: "Brand home",
		Items: []SiteHeaderItem{
			{Label: "Docs", Href: "/docs", Current: true},
			{Label: "Examples", Href: "/examples", Groups: []NavGroup{{Title: "Forms", Items: []NavItem{{Label: "<Form lab>", Href: "/forms"}}}}},
		},
		Section: []NavItem{{Label: "Install", Href: "/docs/install", Current: true}},
		Actions: ThemeToggle(ThemeToggleProps{IconOnly: true}),
		Extras:  templ.Raw(`<a href="https://example.com/source">Source</a>`),
	}))
	for _, want := range []string{
		`<header class="site-header">`,
		`aria-label="Brand home"`,
		`<nav class="site-nav" aria-label="Primary">`,
		`<a href="/docs" aria-current="page">Docs</a>`,
		`<details class="site-menu">`,
		`<p class="site-menu-title">Forms</p>`,
		`&lt;Form lab&gt;`,
		`<a href="/examples">View all`,
		`<summary aria-label="Menu">`,
		`<div class="site-sheet-section">`,
		`<a href="/docs/install" aria-current="page">Install</a>`,
		`<div class="site-sheet-footer">`,
		`data-icon-only="true"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Count(html, "<header") != 1 || strings.Contains(html, "<script") || strings.Contains(html, "<Form lab>") {
		t.Fatal("expected one escaped, script-free header")
	}
	if strings.Count(html, `href="https://example.com/source"`) != 2 {
		t.Fatal("extras render in the bar and the compact menu")
	}
}

func TestSiteHeaderMinimal(t *testing.T) {
	html := render(t, SiteHeader(SiteHeaderProps{Items: []SiteHeaderItem{{Label: "Home", Href: "/"}}}))
	if strings.Contains(html, "site-actions") || strings.Contains(html, "aria-current") || strings.Contains(html, "site-sheet-section") {
		t.Fatal("optional parts should be omitted")
	}
	if !strings.Contains(html, `<a class="site-brand" href="/">`) {
		t.Fatal("brand defaults to the root")
	}
}

func TestAppShellHeaderSlot(t *testing.T) {
	html := render(t, AppShell(AppShellProps{Brand: "Old", Header: templ.Raw(`<header class="custom"></header>`)}))
	if strings.Contains(html, "app-header") || !strings.Contains(html, `class="custom"`) {
		t.Fatal("Header should replace the built-in header")
	}
}
