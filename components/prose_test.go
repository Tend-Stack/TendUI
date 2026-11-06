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

func TestProseDefaultsAndMeasureContract(t *testing.T) {
	css := stylesheet(t)
	for _, tc := range []struct {
		measure TextMeasure
		want    string
	}{
		{"", "standard"}, {MeasureNarrow, "narrow"}, {MeasureStandard, "standard"}, {MeasureWide, "wide"}, {`bad" onclick="alert(1)`, "standard"},
	} {
		html := render(t, Prose(ProseProps{Measure: tc.measure}))
		if html != `<div class="prose" data-measure="`+tc.want+`"></div>` {
			t.Errorf("unsafe or unexpected measure wrapper: %s", html)
		}
		if !strings.Contains(css, `.prose[data-measure="`+tc.want+`"]`) || !strings.Contains(css, "--gothui-prose-"+tc.want+":") {
			t.Errorf("measure %s has no shipped style", tc.want)
		}
	}
}

func TestProsePreservesEmbeddedComponentAndEscapesID(t *testing.T) {
	child := Alert(AlertProps{Title: "A <note>", Message: "A & B", Tone: ToneInfo})
	expectedChild := render(t, child)
	ctx := templ.WithChildren(context.Background(), child)
	var out bytes.Buffer
	if err := Prose(ProseProps{ID: `reading"><script>`, Measure: MeasureNarrow}).Render(ctx, &out); err != nil {
		t.Fatal(err)
	}
	html := out.String()
	if !strings.Contains(html, `id="reading&#34;&gt;&lt;script&gt;"`) || strings.Contains(html, "<script>") {
		t.Fatal("unsafe reading ID")
	}
	if strings.Count(html, expectedChild) != 1 || !strings.HasSuffix(html, expectedChild+"</div>") {
		t.Fatal("Prose rewrote or duplicated child markup")
	}
	for _, unwanted := range []string{`style=`, `role="article"`, `tabindex=`, `<h1>`} {
		if strings.Contains(html, unwanted) {
			t.Errorf("Prose invented document semantics: %s", unwanted)
		}
	}
}

func TestProseChildErrorsAndCancellation(t *testing.T) {
	sentinel := errors.New("content failed")
	child := templ.ComponentFunc(func(context.Context, io.Writer) error { return sentinel })
	ctx := templ.WithChildren(context.Background(), child)
	if err := Prose(ProseProps{}).Render(ctx, io.Discard); !errors.Is(err, sentinel) {
		t.Errorf("lost content error: %v", err)
	}
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	if err := Prose(ProseProps{}).Render(ctx, io.Discard); !errors.Is(err, context.Canceled) {
		t.Errorf("lost cancellation: %v", err)
	}
}

func TestProseDescendantRulesKeepWidgetOptOut(t *testing.T) {
	css := stylesheet(t)
	// Guard the scoping contract: every descendant typography rule excludes an
	// opted-out widget and its subtree. This is not computed-browser verification.
	count := 0
	for _, line := range strings.Split(css, "\n") {
		if strings.HasPrefix(line, ".prose :where(") {
			count++
			if !strings.Contains(line, `:not(:where([data-prose="off"], [data-prose="off"] *))`) {
				t.Errorf("unscoped prose rule: %s", line)
			}
		}
	}
	if count == 0 {
		t.Fatal("no prose typography rules found")
	}
}
