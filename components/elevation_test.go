package components

import (
	"strings"
	"testing"
)

func TestCardElevationAndFallbacks(t *testing.T) {
	for _, elevation := range []Elevation{ElevationFlat, ElevationSubtle, ElevationRaised, ElevationFloating} {
		t.Run(string(elevation), func(t *testing.T) {
			a := Appearance{Elevation: elevation}
			if a.Normalize().Elevation != elevation || a.Attributes()["data-elevation"] != string(elevation) {
				t.Fatal("page preference lost")
			}
			if !strings.Contains(a.CSS(), "--tendui-shadow: var(--tendui-elevation-"+string(elevation)+")") {
				t.Fatal("export differs from page preference")
			}
			html := render(t, Card(CardProps{Elevation: elevation, ID: "outer", Body: Card(CardProps{ID: "inner"})}))
			if strings.Count(html, "data-card-elevation=") != 1 || !strings.Contains(html, `data-card-elevation="`+string(elevation)+`"`) {
				t.Fatal("card override missing or repeated on nested card", html)
			}
			if !strings.Contains(html, `<article id="inner">`) || strings.Contains(html, "style=") || strings.Contains(html, "tabindex") {
				t.Fatal("card semantics or default nested markup changed", html)
			}
		})
	}
	for _, value := range []Elevation{"", "invalid", `floating" onclick="bad`} {
		html := render(t, Card(CardProps{Elevation: value}))
		if html != render(t, Card(CardProps{})) {
			t.Fatal("unknown elevation must inherit page default", html)
		}
	}
	if (Appearance{}).Normalize().Elevation != ElevationRaised {
		t.Fatal("existing page default changed")
	}
}

// Source guards check selector scope/order; actual cascade and visual depth still
// require browser review. Low-specificity card selectors let later print and
// forced-color rules win and keep the outlined surface shadow-free.
func TestElevationStyleContracts(t *testing.T) {
	css := stylesheet(t)
	for _, level := range []string{"flat", "subtle", "raised", "floating"} {
		selector := `article:where([data-card-elevation="` + level + `"])`
		body := paletteBody(t, css, selector)
		if !strings.Contains(body, "box-shadow: var(--tendui-elevation-"+level+");") || strings.Contains(body, "--tendui-shadow:") {
			t.Fatal("card override must not replace an inherited token", body)
		}
		if !strings.Contains(css, "--tendui-elevation-"+level+":") {
			t.Fatal("missing elevation token", level)
		}
	}
	outline := paletteBody(t, css, `article[data-surface="outlined"]`)
	if !strings.Contains(outline, "box-shadow: none;") {
		t.Fatal("outlined surface must remain shadow-free")
	}
	_, printCSS, _ := strings.Cut(css, "@media print {")
	if !strings.Contains(paletteBody(t, printCSS, "article"), "box-shadow: none;") {
		t.Fatal("print must suppress card overrides")
	}
	// This rule is inside forced-colors and after all card override selectors.
	forced := "@media (forced-colors: active) {\n  article[data-surface=\"gradient\"]"
	_, forcedCSS, found := strings.Cut(css, forced)
	if !found || !strings.Contains(paletteBody(t, forcedCSS, "article, .tendui-dialog, .tendui-popover"), "box-shadow: none;") {
		t.Fatal("forced colors must suppress card and overlay shadows")
	}
}
