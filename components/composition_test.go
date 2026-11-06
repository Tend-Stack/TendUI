package components

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestCompositionDefaultsAndClosedOptions(t *testing.T) {
	for _, tc := range []struct {
		name      string
		component templ.Component
		want      []string
	}{
		{"stack default", Stack(StackProps{}), []string{`class="layout-stack"`, `data-gap="md"`, `data-layout-align="stretch"`}},
		{"cluster default", Cluster(ClusterProps{}), []string{`class="layout-cluster"`, `data-gap="md"`, `data-layout-align="center"`, `data-layout-justify="start"`}},
		{"stack unknown", Stack(StackProps{Gap: `bad" onclick="alert(1)`, Align: "reverse"}), []string{`data-gap="md"`, `data-layout-align="stretch"`}},
		{"stack baseline", Stack(StackProps{Align: LayoutAlignBaseline}), []string{`data-layout-align="stretch"`}},
		{"cluster unknown", Cluster(ClusterProps{Gap: "bad", Align: "reverse", Justify: "unsafe"}), []string{`data-gap="md"`, `data-layout-align="center"`, `data-layout-justify="start"`}},
		{"cluster aligned", Cluster(ClusterProps{Gap: SpaceSmall, Align: LayoutAlignBaseline, Justify: LayoutJustifyBetween}), []string{`data-gap="sm"`, `data-layout-align="baseline"`, `data-layout-justify="between"`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			html := render(t, tc.component)
			for _, want := range tc.want {
				if !strings.Contains(html, want) {
					t.Errorf("missing %s in %s", want, html)
				}
			}
			for _, unwanted := range []string{`style=`, `onclick=`, `role=`, `tabindex=`, `id=`, "reverse", "unsafe"} {
				if strings.Contains(html, unwanted) {
					t.Errorf("unexpected %s in %s", unwanted, html)
				}
			}
		})
	}
}

func TestCompositionNestedChildrenAndEscaping(t *testing.T) {
	var body bytes.Buffer
	inner := Cluster(ClusterProps{ID: `actions"><script>`, Gap: SpaceSmall})
	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return inner.Render(templ.WithChildren(ctx, templ.Raw(`<a href="/first">First</a><button type="submit">Second</button>`)), w)
	})
	if err := Stack(StackProps{Gap: SpaceExtraLarge}).Render(templ.WithChildren(context.Background(), children), &body); err != nil {
		t.Fatal(err)
	}
	html := body.String()
	for _, want := range []string{`class="layout-stack" data-gap="xl"`, `class="layout-cluster" data-gap="sm"`, `id="actions&#34;&gt;&lt;script&gt;"`, `First</a><button type="submit">Second`} {
		if !strings.Contains(html, want) {
			t.Errorf("lost nested contract %q in %s", want, html)
		}
	}
	if strings.Count(html, "First</a>") != 1 || strings.Contains(html, "<script>") {
		t.Fatal("children duplicated or ID unescaped")
	}
	if !strings.HasSuffix(html, "</div></div>") {
		t.Fatal("layout wrappers did not close")
	}
}

func TestCompositionPropagatesRenderErrorsAndCancellation(t *testing.T) {
	sentinel := errors.New("child failed")
	child := templ.ComponentFunc(func(context.Context, io.Writer) error { return sentinel })
	for _, component := range []templ.Component{Stack(StackProps{}), Cluster(ClusterProps{})} {
		ctx := templ.WithChildren(context.Background(), child)
		if err := component.Render(ctx, io.Discard); !errors.Is(err, sentinel) {
			t.Fatalf("lost child error: %v", err)
		}
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		if err := component.Render(ctx, io.Discard); !errors.Is(err, context.Canceled) {
			t.Fatalf("lost cancellation: %v", err)
		}
	}
}

// Cross-check public options against the shipped stylesheet, not a separate CSS fixture.
func TestCompositionOptionsHaveStyles(t *testing.T) {
	css := stylesheet(t)
	for _, space := range []Space{SpaceNone, SpaceExtraSmall, SpaceSmall, SpaceMedium, SpaceLarge, SpaceExtraLarge} {
		html := render(t, Stack(StackProps{Gap: space}))
		selector := `data-gap="` + string(space) + `"`
		if !strings.Contains(html, selector) || !strings.Contains(css, "["+selector+"]") {
			t.Errorf("unbacked gap %s", space)
		}
	}
	for _, align := range []LayoutAlign{LayoutAlignStart, LayoutAlignCenter, LayoutAlignEnd, LayoutAlignStretch, LayoutAlignBaseline} {
		selector := `data-layout-align="` + string(align) + `"`
		if !strings.Contains(render(t, Cluster(ClusterProps{Align: align})), selector) || !strings.Contains(css, "["+selector+"]") {
			t.Errorf("unbacked alignment %s", align)
		}
	}
	for _, justify := range []LayoutJustify{LayoutJustifyStart, LayoutJustifyCenter, LayoutJustifyEnd, LayoutJustifyBetween} {
		selector := `data-layout-justify="` + string(justify) + `"`
		if !strings.Contains(render(t, Cluster(ClusterProps{Justify: justify})), selector) || !strings.Contains(css, "["+selector+"]") {
			t.Errorf("unbacked justification %s", justify)
		}
	}
}
