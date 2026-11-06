package components

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestFieldAssociationsAndEscaping(t *testing.T) {
	p := FieldProps{ID: "name", Label: "Name <required>", Description: "Use & keep a short name", Error: "Invalid <value>", Required: true}
	attrs := p.ControlAttributes()
	if attrs["id"] != "name" || attrs["required"] != true || attrs["aria-describedby"] != "name-description name-error" || attrs["aria-invalid"] != "true" {
		t.Fatalf("bad associations: %#v", attrs)
	}
	attrs["id"] = "changed"
	if p.ControlAttributes()["id"] != "name" {
		t.Fatal("shared attribute map")
	}
	html := render(t, Field(p))
	for _, want := range []string{`for="name"`, `id="name-description"`, `id="name-error"`, "Name &lt;required&gt;", "Invalid &lt;value&gt;", "Use &amp; keep"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, `role="alert"`) {
		t.Fatal("field errors must not duplicate summary live announcements")
	}
	optional := FieldProps{ID: "optional", Label: "Optional"}.ControlAttributes()
	for _, key := range []string{"required", "aria-invalid", "aria-describedby"} {
		if _, exists := optional[key]; exists {
			t.Errorf("unexpected %s", key)
		}
	}
	html = render(t, InputGroup(InputGroupProps{Prefix: "<unit>", Suffix: "USD"}))
	if strings.Count(html, `aria-hidden="true"`) != 2 || !strings.Contains(html, "&lt;unit&gt;") {
		t.Fatal("addons must be escaped decoration")
	}
}

func TestFieldsRejectBrokenLabelTargets(t *testing.T) {
	for _, p := range []FieldProps{{}, {ID: "two words", Label: "Name"}, {ID: "tab\tname", Label: "Name"}, {ID: "valid", Label: " "}} {
		var b bytes.Buffer
		if Field(p).Render(context.Background(), &b) == nil {
			t.Fatalf("accepted invalid label target: %+v", p)
		}
	}
}

func TestFormNoValidateOnlyOnSubmitButtons(t *testing.T) {
	for _, p := range []ButtonProps{
		{Text: "Back", Type: ButtonTypeSubmit, FormNoValidate: true},
		{Text: "Normal", Type: ButtonTypeSubmit},
		{Text: "Open", Type: ButtonTypeSubmit, FormNoValidate: true, CommandFor: "dialog", Command: DialogShowModal},
		{Text: "Reset", Type: ButtonTypeReset, FormNoValidate: true},
	} {
		html := render(t, Button(p))
		want := p.Text == "Back"
		if strings.Contains(html, "formnovalidate") != want {
			t.Fatalf("invalid submitter semantics: %s", html)
		}
	}
}
