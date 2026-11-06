package components

import (
	"math"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Tend-Stack/GothUI/v2026"
	"github.com/a-h/templ"
)

func TestAppearanceSafetyAndLayout(t *testing.T) {
	bad := Appearance{Preset: "</style>", Density: "url(evil)", Corners: "invalid", Typography: "invalid", TextSize: "invalid", Width: "invalid", Elevation: "invalid"}
	if bad.Normalize() != (Appearance{}).Normalize() {
		t.Fatal("unknown settings did not fall back to defaults")
	}
	if strings.Contains(bad.CSS(), "evil") || strings.Contains(bad.CSS(), "</style>") {
		t.Fatal("untrusted CSS emitted")
	}
	a := Appearance{Preset: PresetForest, Density: DensityCompact, Corners: CornersSquare, Typography: TypographySerif, TextSize: TextLarge, Width: WidthWide, Elevation: ElevationFlat}
	attrs := a.Attributes()
	attrs["data-preset"] = "changed"
	if a.Attributes()["data-preset"] != "forest" {
		t.Fatal("shared attribute map mutated")
	}
	html := render(t, Layout(LayoutProps{Title: "Appearance", Appearance: a, Head: templ.Raw(`<link rel="stylesheet" href="/custom.css">`)}, nil))
	for _, want := range []string{`data-preset="forest"`, `data-density="compact"`, `data-corners="square"`, `data-typography="serif"`, `data-text-size="large"`, `data-width="wide"`, `data-elevation="flat"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Index(html, "/custom.css") < strings.Index(html, gothui.StylesheetURL) {
		t.Fatal("custom stylesheet must follow base CSS")
	}
}

func stylesheet(t *testing.T) string {
	t.Helper()
	w := httptest.NewRecorder()
	gothui.AssetHandler().ServeHTTP(w, httptest.NewRequest("GET", gothui.StylesheetURL, nil))
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	return w.Body.String()
}

func TestAppearanceExportsUseExistingTokens(t *testing.T) {
	css := stylesheet(t)
	ref := regexp.MustCompile(`var\((--gothui-[a-z-]+)`)
	for _, a := range []Appearance{
		{}, {Preset: PresetNeutral, Density: DensityCompact, Corners: CornersSquare, Typography: TypographySerif, TextSize: TextLarge, Width: WidthNarrow, Elevation: ElevationFlat},
		{Preset: PresetForest, Corners: CornersSoft, Typography: TypographyMono, Width: WidthWide}, {Preset: PresetRose},
	} {
		for _, m := range ref.FindAllStringSubmatch(a.CSS(), -1) {
			if !strings.Contains(css, m[1]+":") {
				t.Errorf("export references missing token %s", m[1])
			}
		}
		for attr, value := range a.Attributes() {
			selector := `:root:where([` + attr + `="` + value.(string) + `"])`
			if !strings.Contains(css, selector) {
				t.Errorf("no CSS selector for %s", selector)
			}
		}
	}
}

func luminance(hex string) float64 {
	var rgb [3]float64
	for i := range rgb {
		v, _ := strconv.ParseUint(hex[1+i*2:3+i*2], 16, 8)
		c := float64(v) / 255
		if c <= 0.04045 {
			rgb[i] = c / 12.92
		} else {
			rgb[i] = math.Pow((c+0.055)/1.055, 2.4)
		}
	}
	return .2126*rgb[0] + .7152*rgb[1] + .0722*rgb[2]
}

func contrast(a, b string) float64 {
	x, y := luminance(a), luminance(b)
	return (math.Max(x, y) + .05) / (math.Min(x, y) + .05)
}

func TestPresetTextContrast(t *testing.T) {
	css := stylesheet(t)
	colors := map[string]string{}
	for _, m := range regexp.MustCompile(`--gothui-preset-([a-z-]+):\s*(#[0-9a-f]{6})`).FindAllStringSubmatch(css, -1) {
		colors[m[1]] = m[2]
	}
	for _, preset := range []string{"indigo", "neutral", "forest", "rose"} {
		for _, state := range []string{"light", "hover-light", "dark", "hover-dark"} {
			color := colors[preset+"-"+state]
			if color == "" {
				t.Fatalf("missing color %s-%s", preset, state)
			}
			text, background := "#ffffff", "#ffffff"
			if strings.HasSuffix(state, "dark") {
				text = colors[preset+"-on-dark"]
				background = "#1b2131"
			}
			if text == "" {
				t.Fatal("missing dark button text")
			}
			if ratio := contrast(color, text); ratio < 4.5 {
				t.Errorf("%s %s button contrast %.2f", preset, state, ratio)
			}
			if ratio := contrast(color, background); ratio < 4.5 {
				t.Errorf("%s %s link contrast %.2f", preset, state, ratio)
			}
		}
	}
}
