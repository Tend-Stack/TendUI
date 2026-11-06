package components

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

type Preset string
type Density string
type Corners string
type Typography string
type TextSize string
type ContentWidth string
type Elevation string
type MotionPreference string

const (
	PresetIndigo       Preset           = "indigo"
	PresetNeutral      Preset           = "neutral"
	PresetForest       Preset           = "forest"
	PresetRose         Preset           = "rose"
	DensityComfortable Density          = "comfortable"
	DensityCompact     Density          = "compact"
	CornersRounded     Corners          = "rounded"
	CornersSquare      Corners          = "square"
	CornersSoft        Corners          = "soft"
	TypographySans     Typography       = "sans"
	TypographySerif    Typography       = "serif"
	TypographyMono     Typography       = "mono"
	TextStandard       TextSize         = "standard"
	TextLarge          TextSize         = "large"
	WidthStandard      ContentWidth     = "standard"
	WidthNarrow        ContentWidth     = "narrow"
	WidthWide          ContentWidth     = "wide"
	ElevationRaised    Elevation        = "raised"
	ElevationFlat      Elevation        = "flat"
	ElevationSubtle    Elevation        = "subtle"
	ElevationFloating  Elevation        = "floating"
	MotionSystem       MotionPreference = "system"
	MotionOff          MotionPreference = "off"
)

// Appearance configures independent CSS tokens. Its zero value preserves GothUI's defaults.
// Normalize untrusted preferences before persistence; Layout always normalizes them too.
type Appearance struct {
	Preset     Preset           `json:"preset"`
	Density    Density          `json:"density"`
	Corners    Corners          `json:"corners"`
	Typography Typography       `json:"typography"`
	TextSize   TextSize         `json:"text_size"`
	Width      ContentWidth     `json:"width"`
	Elevation  Elevation        `json:"elevation"`
	Motion     MotionPreference `json:"motion"`
}

func allowed[T ~string](value, defaultValue T, choices ...T) T {
	for _, choice := range choices {
		if value == choice {
			return value
		}
	}
	return defaultValue
}

func (a Appearance) Normalize() Appearance {
	a.Preset = allowed(a.Preset, PresetIndigo, PresetIndigo, PresetNeutral, PresetForest, PresetRose)
	a.Density = allowed(a.Density, DensityComfortable, DensityComfortable, DensityCompact)
	a.Corners = allowed(a.Corners, CornersRounded, CornersRounded, CornersSquare, CornersSoft)
	a.Typography = allowed(a.Typography, TypographySans, TypographySans, TypographySerif, TypographyMono)
	a.TextSize = allowed(a.TextSize, TextStandard, TextStandard, TextLarge)
	a.Width = allowed(a.Width, WidthStandard, WidthStandard, WidthNarrow, WidthWide)
	a.Elevation = allowed(a.Elevation, ElevationRaised, ElevationRaised, ElevationFlat, ElevationSubtle, ElevationFloating)
	a.Motion = allowed(a.Motion, MotionSystem, MotionSystem, MotionOff)
	return a
}

// Attributes returns a fresh map, safe to use concurrently and spread on your own html element.
func (a Appearance) Attributes() templ.Attributes {
	a = a.Normalize()
	return templ.Attributes{
		"data-preset": string(a.Preset), "data-density": string(a.Density),
		"data-corners": string(a.Corners), "data-typography": string(a.Typography),
		"data-text-size": string(a.TextSize), "data-width": string(a.Width),
		"data-elevation": string(a.Elevation),
		"data-motion":    string(a.Motion),
	}
}

// CSS exports overrides referencing the canonical tokens in gothui.css; no colors are duplicated.
// Load after gothui.css. Root attributes and these exports use the same specificity.
func (a Appearance) CSS() string {
	a = a.Normalize()
	var b strings.Builder
	b.WriteString("/* GothUI appearance. Load after gothui.css. */\n:root {\n")
	for _, part := range []string{"light", "hover-light", "dark", "hover-dark", "on-dark"} {
		fmt.Fprintf(&b, "  --gothui-accent-%s: var(--gothui-preset-%s-%s);\n", part, a.Preset, part)
	}
	for _, part := range []string{"space", "control-y", "control-x", "control-height", "card-padding", "grid-gap"} {
		fmt.Fprintf(&b, "  --gothui-%s: var(--gothui-density-%s-%s);\n", part, a.Density, part)
	}
	fmt.Fprintf(&b, "  --gothui-radius: var(--gothui-corners-%s);\n  --gothui-font: var(--gothui-font-%s);\n  --gothui-text-size: var(--gothui-size-%s);\n  --gothui-content-width: var(--gothui-width-%s);\n  --gothui-shadow: var(--gothui-elevation-%s);\n}\n", a.Corners, a.Typography, a.TextSize, a.Width, a.Elevation)
	if a.Motion == MotionOff {
		b.WriteString(":root { --gothui-duration-fast: 0ms; --gothui-duration-enter: 0ms; --gothui-duration-loading: 0ms; --gothui-stagger-step: 0ms; --gothui-lift-distance: 0px; }\n")
	}
	return b.String()
}
