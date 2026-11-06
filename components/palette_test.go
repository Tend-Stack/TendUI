package components

import (
	"regexp"
	"strings"
	"testing"
)

// These tests inspect literal source colors, not a browser's computed cascade.
// Keep selectors explicit so missing or reformatted theme blocks fail visibly.
func paletteBody(t *testing.T, css, selector string) string {
	t.Helper()
	_, tail, ok := strings.Cut(css, selector+" {")
	if !ok {
		t.Fatalf("missing palette selector %s", selector)
	}
	body, _, ok := strings.Cut(tail, "}")
	if !ok {
		t.Fatalf("unclosed palette selector %s", selector)
	}
	return body
}

func literalPalette(t *testing.T, body string) map[string]string {
	t.Helper()
	result := map[string]string{}
	for _, match := range regexp.MustCompile(`--gothui-([a-z-]+):\s*(#[0-9a-f]{6})(?:\s*!important)?\s*;`).FindAllStringSubmatch(body, -1) {
		result[match[1]] = match[2]
	}
	return result
}

func paletteColor(t *testing.T, colors map[string]string, name string) string {
	t.Helper()
	color := colors[name]
	if color == "" {
		t.Fatalf("missing literal color %s", name)
	}
	return color
}

func TestScreenPaletteContrastAcrossSurfaces(t *testing.T) {
	css := stylesheet(t)
	presets := literalPalette(t, css)
	for _, theme := range []struct{ name, selector string }{{"light", ":root"}, {"dark", `:root[data-theme="dark"]`}} {
		colors := literalPalette(t, paletteBody(t, css, theme.selector))
		for _, bg := range []string{"background", "surface", "subtle"} {
			background := paletteColor(t, colors, bg)
			for _, fg := range []string{"text", "muted", "danger", "info", "success", "warning", "border"} {
				minimum := 4.5
				if fg == "border" {
					minimum = 3
				}
				ratio := contrast(paletteColor(t, colors, fg), background)
				if ratio < minimum {
					t.Errorf("%s %s on %s: %.4f below %.1f", theme.name, fg, bg, ratio, minimum)
				}
			}
			for _, preset := range []string{"indigo", "neutral", "forest", "rose"} {
				for _, state := range []string{"", "hover-"} {
					key := "preset-" + preset + "-" + state + theme.name
					ratio := contrast(paletteColor(t, presets, key), background)
					if ratio < 4.5 {
						t.Errorf("%s link on %s: %.4f below 4.5", key, bg, ratio)
					}
				}
			}
		}
	}
}

func TestSystemDarkPaletteMatchesExplicitDark(t *testing.T) {
	css := stylesheet(t)
	explicit := paletteBody(t, css, `:root[data-theme="dark"]`)
	system := paletteBody(t, css, `:root:not([data-theme="light"]):not([data-theme="dark"])`)
	// Include variable references and elevation, not just literal hex values.
	declarations := regexp.MustCompile(`(?:--gothui-[a-z-]+|color-scheme):\s*[^;]+;`)
	want, got := declarations.FindAllString(explicit, -1), declarations.FindAllString(system, -1)
	if strings.Join(want, "\n") != strings.Join(got, "\n") {
		t.Fatal("system dark and explicit dark theme declarations differ")
	}
}

func TestPrintPaletteWinsOverScreenThemes(t *testing.T) {
	css := stylesheet(t)
	_, printCSS, ok := strings.Cut(css, "@media print {")
	if !ok {
		t.Fatal("missing print styles")
	}
	body := paletteBody(t, printCSS, ":root")
	colors := literalPalette(t, body)
	for _, name := range []string{"background", "surface", "subtle", "text", "muted", "border", "primary", "primary-hover", "on-primary", "danger", "info", "success", "warning", "focus"} {
		color := paletteColor(t, colors, name)
		if !strings.Contains(body, "--gothui-"+name+": "+color+" !important;") {
			t.Errorf("print %s needs priority over explicit/system dark selectors and appearance exports", name)
		}
	}
	if !strings.Contains(body, "color-scheme: light !important;") {
		t.Error("native controls must use light print scheme")
	}
	for _, bg := range []string{"background", "surface", "subtle"} {
		if colors[bg] != "#ffffff" {
			t.Errorf("print %s must be white", bg)
		}
		for _, fg := range []string{"text", "muted", "primary", "primary-hover", "danger", "info", "success", "warning"} {
			if ratio := contrast(colors[fg], colors[bg]); ratio < 4.5 {
				t.Errorf("print %s on %s: %.4f", fg, bg, ratio)
			}
		}
		if ratio := contrast(colors["border"], colors[bg]); ratio < 3 {
			t.Errorf("print border on %s: %.4f", bg, ratio)
		}
	}
	// Printing with backgrounds disabled must not leave white button labels.
	buttons := paletteBody(t, printCSS, `button, .button, input:is([type="submit"], [type="button"], [type="reset"])`)
	for _, declaration := range []string{"background: transparent !important;", "color: var(--gothui-text) !important;", "border-color: var(--gothui-border) !important;"} {
		if !strings.Contains(buttons, declaration) {
			t.Errorf("print buttons missing %s", declaration)
		}
	}
	step := paletteBody(t, printCSS, `.steps [aria-current="step"] .step-number`)
	for _, declaration := range []string{"background: transparent;", "color: var(--gothui-text);", "border-width: 2px;"} {
		if !strings.Contains(step, declaration) {
			t.Errorf("print current-step marker missing %s", declaration)
		}
	}
	// Preserve a native checked-state mark without depending on a colored pseudo-element.
	switches := paletteBody(t, printCSS, `input[type="checkbox"][role="switch"]`)
	if !strings.Contains(switches, "appearance: auto;") || !strings.Contains(paletteBody(t, printCSS, `input[role="switch"]::before`), "content: none;") {
		t.Error("print switch must use its native checkbox state")
	}
}
