package components

import (
	"fmt"
	"sort"
	"strings"
)

// IconName selects an inline drawing authored for GothUI. No URL or SVG markup
// supplied by an application is interpreted as an icon.
type IconName string

const (
	IconNone         IconName = ""
	IconArrowLeft    IconName = "arrow-left"
	IconArrowRight   IconName = "arrow-right"
	IconArrowUp      IconName = "arrow-up"
	IconArrowDown    IconName = "arrow-down"
	IconCheck        IconName = "check"
	IconClose        IconName = "close"
	IconPlus         IconName = "plus"
	IconMinus        IconName = "minus"
	IconSearch       IconName = "search"
	IconMenu         IconName = "menu"
	IconChevronDown  IconName = "chevron-down"
	IconChevronRight IconName = "chevron-right"
	IconDownload     IconName = "download"
	IconUpload       IconName = "upload"
	IconSun          IconName = "sun"
	IconMoon         IconName = "moon"
	IconInfo         IconName = "info"
	IconWarning      IconName = "warning"
	IconExternalLink IconName = "external-link"
	IconEdit         IconName = "edit"
	IconUser         IconName = "user"
	IconFolder       IconName = "folder"
	IconFile         IconName = "file"
	IconCalendar     IconName = "calendar"
)

type IconSize string

const (
	IconSizeSmall  IconSize = "small"
	IconSizeMedium IconSize = "medium"
	IconSizeLarge  IconSize = "large"
)

// IconProps defaults to a decorative, non-focusable icon. A nonblank Label
// makes a standalone icon an image with an accessible name. Inside a named
// link/button, keep the icon decorative and name the control instead.
type IconProps struct {
	Name  IconName
	Label string
	Size  IconSize
}

// All paths use a 24-unit viewBox and the same rounded two-unit stroke.
// These geometric drawings are original library assets, not a remote sprite.
var iconPaths = map[IconName]string{
	IconArrowLeft:    "M20 12H4m7-7-7 7 7 7",
	IconArrowRight:   "M4 12h16m-7-7 7 7-7 7",
	IconArrowUp:      "M12 20V4m-7 7 7-7 7 7",
	IconArrowDown:    "M12 4v16m-7-7 7 7 7-7",
	IconCheck:        "m4 12 5 5L20 6",
	IconClose:        "m6 6 12 12M6 18 18 6",
	IconPlus:         "M4 12h16M12 4v16",
	IconMinus:        "M4 12h16",
	IconSearch:       "M17 10a7 7 0 1 1-14 0 7 7 0 0 1 14 0Zm-2 5 6 6",
	IconMenu:         "M4 6h16M4 12h16M4 18h16",
	IconChevronDown:  "m5 9 7 7 7-7",
	IconChevronRight: "m9 5 7 7-7 7",
	IconDownload:     "M12 3v12m-5-5 5 5 5-5M4 16v5h16v-5",
	IconUpload:       "M12 15V3m-5 5 5-5 5 5M4 16v5h16v-5",
	IconSun:          "M16 12a4 4 0 1 1-8 0 4 4 0 0 1 8 0ZM12 2v2m0 16v2M2 12h2m16 0h2M5 5l1.5 1.5m11 11L19 19M5 19l1.5-1.5m11-11L19 5",
	IconMoon:         "M20 15A9 9 0 0 1 9 3a9 9 0 1 0 11 12Z",
	IconInfo:         "M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM12 11v6m0-10v.1",
	IconWarning:      "M12 3 2 21h20L12 3Zm0 6v5m0 3v.1",
	IconExternalLink: "M14 3h7v7m0-7L10 14M10 5H3v16h16v-7",
	IconEdit:         "m14 6 4 4M4 20l5-1L21 7l-4-4L5 15l-1 5Z",
	IconUser:         "M16 7a4 4 0 1 1-8 0 4 4 0 0 1 8 0ZM4 21v-2a8 6 0 0 1 16 0v2",
	IconFolder:       "M3 5h7l2 3h9v12H3V5Z",
	IconFile:         "M5 3h9l5 5v13H5V3Zm9 0v5h5M8 13h8m-8 4h6",
	IconCalendar:     "M4 5h16v16H4V5Zm4-3v6m8-6v6M4 11h16m-12 4h1m6 0h1m-8 3h1m6 0h1",
}

// IconNames returns a fresh, sorted list of all drawings, excluding IconNone.
func IconNames() []IconName {
	names := make([]IconName, 0, len(iconPaths))
	for name := range iconPaths {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })
	return names
}

func validateIcon(name IconName) error {
	if name != IconNone {
		if _, ok := iconPaths[name]; !ok {
			return fmt.Errorf("gothui: unknown icon %q", name)
		}
	}
	return nil
}

func validateActionIcons(before, after IconName, text, label string) error {
	for _, name := range []IconName{before, after} {
		if err := validateIcon(name); err != nil {
			return err
		}
	}
	if (before != IconNone || after != IconNone) && strings.TrimSpace(text) == "" && strings.TrimSpace(label) == "" {
		return fmt.Errorf("gothui: an icon-only action requires an AriaLabel")
	}
	return nil
}
