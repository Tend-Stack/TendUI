package components

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestNativeButtonContracts(t *testing.T) {
	html := render(t, Button(ButtonProps{Text: "Open", Type: ButtonTypeSubmit, CommandFor: "confirm", Command: DialogShowModal}))
	for _, want := range []string{`type="button"`, `commandfor="confirm"`, `command="show-modal"`} {
		if !strings.Contains(html, want) {
			t.Error(want)
		}
	}
	html = render(t, Button(ButtonProps{Text: "Toggle", Type: ButtonTypeSubmit, PopoverTarget: "help"}))
	for _, want := range []string{`type="button"`, `popovertarget="help"`, `popovertargetaction="toggle"`} {
		if !strings.Contains(html, want) {
			t.Error(want)
		}
	}
	html = render(t, Button(ButtonProps{Type: ButtonTypeSubmit, CommandFor: "dialog", Command: "unsupported", PopoverTarget: "ignored"}))
	if !strings.Contains(html, `type="button"`) || strings.Contains(html, "commandfor=") || strings.Contains(html, "popovertarget=") {
		t.Fatal("invalid native action should be inert")
	}
}

func TestOverlayLabelsAndCloseActions(t *testing.T) {
	html := render(t, Dialog(DialogProps{ID: "confirm", Title: "<Review>", Description: "Check & confirm", Footer: templ.Raw("<p>Footer</p>")}))
	for _, want := range []string{`<dialog id="confirm"`, `aria-labelledby="confirm-title"`, `aria-describedby="confirm-description"`, `id="confirm-title"`, `&lt;Review&gt;`, `Check &amp; confirm`, `method="dialog"`, `autofocus`, `<p>Footer</p>`} {
		if !strings.Contains(html, want) {
			t.Error(want)
		}
	}
	if strings.Contains(html, ` open`) || strings.Contains(html, "tabindex=") || strings.Contains(html, "<script") {
		t.Fatal("dialog should be closed and native")
	}
	for _, mode := range []PopoverMode{PopoverAuto, PopoverManual} {
		html = render(t, Popover(PopoverProps{ID: "help", Title: "Help", Mode: mode}))
		for _, want := range []string{`popover="` + string(mode) + `"`, `aria-labelledby="help-title"`, `popovertarget="help"`, `popovertargetaction="hide"`} {
			if !strings.Contains(html, want) {
				t.Error(want)
			}
		}
		if strings.Contains(html, `role="menu"`) {
			t.Fatal("generic content must not claim menu semantics")
		}
	}
}

func TestOverlayIDsAreRequired(t *testing.T) {
	for _, id := range []string{"", "two ids", "tab\tid", "line\nid"} {
		for _, c := range []templ.Component{Dialog(DialogProps{ID: id}), Popover(PopoverProps{ID: id})} {
			var buf bytes.Buffer
			if err := c.Render(context.Background(), &buf); err == nil {
				t.Fatalf("invalid ID %q accepted", id)
			}
		}
	}
}
