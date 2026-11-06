package components

import (
	"bytes"
	"context"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestIconDrawingsAndAccessibleNames(t *testing.T) {
	names := IconNames()
	if len(names) != 24 {
		t.Fatalf("expected the documented 24 drawings, got %d", len(names))
	}
	for i, name := range names {
		if i > 0 && names[i-1] >= name {
			t.Fatal("names must be unique and sorted")
		}
		t.Run(string(name), func(t *testing.T) {
			html := render(t, Icon(IconProps{Name: name}))
			var svg struct {
				XMLName   xml.Name
				Width     string `xml:"width,attr"`
				Height    string `xml:"height,attr"`
				Hidden    string `xml:"aria-hidden,attr"`
				Focusable string `xml:"focusable,attr"`
				Stroke    string `xml:"stroke,attr"`
				Path      []struct {
					D string `xml:"d,attr"`
				} `xml:"path"`
			}
			if err := xml.Unmarshal([]byte(html), &svg); err != nil {
				t.Fatal(err)
			}
			if svg.XMLName.Space != "http://www.w3.org/2000/svg" || svg.Width != "24" || svg.Height != "24" || svg.Hidden != "true" || svg.Focusable != "false" || svg.Stroke != "currentColor" || len(svg.Path) != 1 || svg.Path[0].D == "" {
				t.Fatalf("broken SVG contract: %s", html)
			}
			for _, forbidden := range []string{"<script", "<image", "<use", "href=", "style=", "role=", "aria-label="} {
				if strings.Contains(html, forbidden) {
					t.Fatalf("unexpected active/external/named markup: %s", html)
				}
			}
		})
	}
	names[0] = "changed"
	if IconNames()[0] == "changed" {
		t.Fatal("caller mutated shared icon registry")
	}
	label := `Complete <script> " onload="bad`
	html := render(t, Icon(IconProps{Name: IconCheck, Label: label, Size: "unrecognized"}))
	var svg struct {
		Role  string `xml:"role,attr"`
		Label string `xml:"aria-label,attr"`
		Size  string `xml:"data-size,attr"`
	}
	if err := xml.Unmarshal([]byte(html), &svg); err != nil {
		t.Fatal(err)
	}
	if svg.Role != "img" || svg.Label != label || svg.Size != "medium" || strings.Contains(html, "aria-hidden=") || strings.Contains(html, "<script>") || strings.Contains(html, ` onload="`) {
		t.Fatalf("accessible name escaped incorrectly: %s", html)
	}
	if got := render(t, Icon(IconProps{Name: IconNone})); got != "" {
		t.Fatalf("zero-value icon emitted markup: %q", got)
	}
	if got := render(t, Icon(IconProps{Name: IconInfo, Label: " \t\n"})); !strings.Contains(got, `aria-hidden="true"`) {
		t.Fatal("blank label made a named image")
	}
}

func TestIconActionsPreserveNamesAndNativeBehavior(t *testing.T) {
	html := render(t, Button(ButtonProps{Text: "Search", Type: ButtonTypeSubmit, IconBefore: IconSearch, IconAfter: IconArrowRight, Name: "action", Value: "search", HxPost: "/search"}))
	if strings.Count(html, "<svg") != 2 || strings.Contains(html, `role="img"`) || !strings.Contains(html, `type="submit"`) || !strings.Contains(html, `name="action"`) || !strings.Contains(html, `hx-post="/search"`) || !strings.Contains(html, "Search") {
		t.Fatal(html)
	}
	for _, action := range []templ.Component{
		Button(ButtonProps{IconBefore: IconClose, AriaLabel: "Close preview", CommandFor: "preview", Command: DialogClose}),
		Link(LinkProps{Href: "/catalog", IconBefore: IconMenu, AriaLabel: "Browse components"}),
	} {
		html := render(t, action)
		if !strings.Contains(html, "aria-label=") || !strings.Contains(html, `aria-hidden="true"`) || strings.Contains(html, `role="img"`) {
			t.Fatal(html)
		}
	}
	// No-icon controls retain their existing behavior, without SVG wrappers.
	if html := render(t, Button(ButtonProps{Text: "Continue"})); strings.Contains(html, "svg") || strings.Contains(html, "aria-label=") {
		t.Fatal(html)
	}
	for _, c := range []templ.Component{
		Icon(IconProps{Name: `"><script>`}),
		Button(ButtonProps{Text: "Find", IconBefore: "typo"}),
		Link(LinkProps{Text: "Find", Href: "/", IconAfter: "typo"}),
		Button(ButtonProps{IconBefore: IconSearch, Text: " \t", AriaLabel: "\n"}),
		Link(LinkProps{Href: "/", IconAfter: IconSearch}),
	} {
		var out bytes.Buffer
		if err := c.Render(context.Background(), &out); err == nil || out.Len() != 0 {
			t.Fatalf("invalid icon action rendered partial content: %v %s", err, out.String())
		}
	}
}
