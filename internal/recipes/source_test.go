package recipes

import (
	"strings"
	"testing"
)

func TestRecipeBoundaries(t *testing.T) {
	source := "// catalog:Button\ntempl ExampleButton() {}\n// catalog:end\n"
	names := Names(source)
	if len(names) != 1 || names[0] != "Button" {
		t.Fatal(names)
	}
	got, err := Source(source, "Button")
	if err != nil || !strings.Contains(got, "package examples") || !strings.Contains(got, "templ ExampleButton()") || strings.Contains(got, "catalog:end") {
		t.Fatal(got, err)
	}
	for _, tc := range []struct{ source, name string }{{source, "../Button"}, {source, "Missing"}, {source + source, "Button"}, {"// catalog:Button\n// catalog:Card\n// catalog:end", "Button"}, {"// catalog:Button\nmissing end", "Button"}} {
		if _, err := Source(tc.source, tc.name); err == nil {
			t.Fatalf("accepted malformed recipe %q", tc.name)
		}
	}
}
