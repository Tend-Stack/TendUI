package components

import (
	"strings"
	"testing"
)

func TestExtendedControlSemantics(t *testing.T) {
	html := render(t, Textarea(TextareaProps{Label: "Notes", Name: "notes", Value: "</textarea><script>bad()</script>", Required: true, ReadOnly: true, MaxLength: 100, Rows: 4, Invalid: true, DescribedBy: "notes-hint"}))
	for _, want := range []string{`maxlength="100"`, `rows="4"`, `required`, `readonly`, `aria-invalid="true"`, `aria-describedby="notes-hint"`, `&lt;/textarea&gt;`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, "<script>") {
		t.Fatal("textarea escaped its element")
	}
	html = render(t, RadioGroup(RadioGroupProps{Name: "plan", Legend: "Choose a plan", Value: "b", Required: true, Options: []Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}, {Label: "C", Value: "c", Disabled: true}}}))
	if strings.Count(html, " checked") != 1 || !strings.Contains(html, `value="b" checked`) || !strings.Contains(html, "<legend>Choose a plan</legend>") {
		t.Fatal(html)
	}
	html = render(t, Switch(SwitchProps{Label: "Updates", Name: "updates", Checked: true, Disabled: true}))
	for _, want := range []string{`type="checkbox"`, `role="switch"`, `checked`, `disabled`} {
		if !strings.Contains(html, want) {
			t.Error(want)
		}
	}
	html = render(t, FileInput(FileInputProps{Label: "Files", Name: "files", Accept: "image/*", Multiple: true, Required: true}))
	for _, want := range []string{`type="file"`, `accept="image/*"`, `multiple`, `required`} {
		if !strings.Contains(html, want) {
			t.Error(want)
		}
	}
	if !strings.Contains(render(t, Form(FormProps{Method: MethodPost, Multipart: true})), `enctype="multipart/form-data"`) {
		t.Fatal("missing multipart encoding")
	}
}

func TestFeedbackAndShellSemantics(t *testing.T) {
	html := render(t, Alert(AlertProps{Title: "Notice", Message: "<unsafe>", Tone: ToneInfo}))
	if strings.Contains(html, "role=") || !strings.Contains(html, "&lt;unsafe&gt;") {
		t.Fatal(html)
	}
	if !strings.Contains(render(t, Alert(AlertProps{Tone: ToneDanger, Live: true})), `role="alert"`) {
		t.Fatal("missing urgent announcement")
	}
	if !strings.Contains(render(t, Alert(AlertProps{Tone: ToneSuccess, Live: true})), `role="status"`) {
		t.Fatal("missing polite announcement")
	}
	if !strings.Contains(render(t, Badge(BadgeProps{Text: "Ready", Tone: "unknown"})), `data-tone="neutral"`) {
		t.Fatal("invalid badge tone")
	}
	html = render(t, ErrorSummary(ErrorSummaryProps{Errors: []FieldError{{ID: "email", Message: "Check <email>"}}}))
	if !strings.Contains(html, `href="#email"`) || !strings.Contains(html, "Check &lt;email&gt;") {
		t.Fatal(html)
	}
	if render(t, ErrorSummary(ErrorSummaryProps{})) != "" {
		t.Fatal("empty summary rendered")
	}
	html = render(t, AppShell(AppShellProps{Brand: "Studio", Title: "Workspace", Navigation: []NavItem{{Label: "Home", Href: "/", Current: true}, {Label: "Bad", Href: "javascript:bad()"}}}))
	for _, want := range []string{`<main id="main"`, `href="#main"`, `aria-label="Application"`, `aria-current="page"`, `<details open>`} {
		if !strings.Contains(html, want) {
			t.Error(want)
		}
	}
	if strings.Count(html, "<main ") != 1 || strings.Contains(html, "javascript:") {
		t.Fatal("invalid shell landmark or unsafe URL")
	}
}
