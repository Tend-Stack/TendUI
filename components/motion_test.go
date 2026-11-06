package components

import (
	"strings"
	"testing"
)

func TestMotionPreferencesAndSafeAttributes(t *testing.T) {
	if (Appearance{Motion: "unexpected"}).Normalize().Motion != MotionSystem {
		t.Fatal("unknown motion must follow the system")
	}
	a := Appearance{Motion: MotionOff}
	html := render(t, Layout(LayoutProps{Appearance: a}, Card(CardProps{Entry: EntryRise, Delay: DelayMedium, Surface: SurfaceGradient, HoverLift: true})))
	for _, want := range []string{`data-motion="off"`, `data-enter="rise"`, `data-enter-delay="medium"`, `data-surface="gradient"`, `data-hover="lift"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	for _, prop := range []string{"duration-fast", "duration-enter", "duration-loading", "stagger-step"} {
		if !strings.Contains(a.CSS(), "--tendui-"+prop+": 0ms") {
			t.Errorf("motion export did not disable %s", prop)
		}
	}
	if !strings.Contains(a.CSS(), "--tendui-lift-distance: 0px") {
		t.Fatal("off export still moves hover cards")
	}
	html = render(t, Card(CardProps{Entry: "bad", Delay: "bad", Surface: "bad"}))
	if strings.Contains(html, "data-enter") || strings.Contains(html, "data-surface") || strings.Contains(html, "style=") {
		t.Fatal("unvalidated decoration emitted")
	}
	attrs := entryAttrs(EntryFade, DelayShort)
	attrs["data-enter"] = "changed"
	if entryAttrs(EntryFade, DelayShort)["data-enter"] != "fade" {
		t.Fatal("shared mutable attributes")
	}
	if len(entryAttrs(EntryNone, DelayLong)) != 0 {
		t.Fatal("delay without entrance")
	}
}

func TestLoadingAndEmptyStateSemantics(t *testing.T) {
	for _, tc := range []struct{ lines, want int }{{0, 3}, {-1, 3}, {1, 1}, {8, 8}, {999999, 8}} {
		html := render(t, Skeleton(SkeletonProps{Label: "Loading <projects>", Lines: tc.lines, Avatar: true}))
		if strings.Count(html, `class="skeleton-line"`) != tc.want || !strings.Contains(html, `aria-hidden="true"`) || !strings.Contains(html, `role="status"`) || !strings.Contains(html, "Loading &lt;projects&gt;") {
			t.Fatalf("incorrect skeleton %d: %s", tc.lines, html)
		}
	}
	html := render(t, EmptyState(EmptyStateProps{Title: "No <results>", Message: "Try another filter"}))
	if !strings.Contains(html, "No &lt;results&gt;") || strings.Contains(html, `role="alert"`) {
		t.Fatal("static empty state must escape content without announcing an alert")
	}
	html = render(t, Stat(StatProps{Label: "Requests <today>", Value: "1,234", Detail: "Sample only"}))
	if !strings.Contains(html, "<dl") || !strings.Contains(html, "<dt>Requests &lt;today&gt;</dt>") || !strings.Contains(html, "1,234") {
		t.Fatal("metric lost label/value semantics")
	}
}

func TestMotionStylesHaveBoundedDecorativeLoops(t *testing.T) {
	css := stylesheet(t)
	for _, required := range []string{"@media (prefers-reduced-motion: no-preference)", "@media (prefers-reduced-motion: reduce)", "animation: none !important", "tendui-loading var(--tendui-duration-loading) ease-in-out 3", "animation-iteration-count: 1", "@media (forced-colors: active)"} {
		if !strings.Contains(css, required) {
			t.Errorf("missing motion contract %q", required)
		}
	}
	if strings.Contains(css, "infinite") || strings.Contains(css, "@import") {
		t.Fatal("unbounded motion or external CSS dependency")
	}
}
