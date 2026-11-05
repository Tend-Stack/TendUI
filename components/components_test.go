package components

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b bytes.Buffer
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestLayoutAndEscaping(t *testing.T) {
	html := render(t, Layout(LayoutProps{Title: "<script>alert(1)</script>", Theme: ThemeDark}, Button(ButtonProps{Text: "<b>safe</b>"})))
	for _, want := range []string{`lang="en"`, `data-theme="dark"`, `/gothui/static/gothui.css`, `&lt;script&gt;`, `&lt;b&gt;safe&lt;/b&gt;`, `type="button"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, "<script") {
		t.Fatal("script emitted")
	}
	if strings.Contains(render(t, Layout(LayoutProps{Theme: "invalid"}, nil)), "data-theme") {
		t.Fatal("invalid theme emitted")
	}
}

func TestButtonHTMXAndSafeLink(t *testing.T) {
	html := render(t, Button(ButtonProps{Text: "Save", Type: ButtonTypeSubmit, Primary: true, Outline: true, Disabled: true, HxPost: "/save", HxTarget: "#result", HxSwap: "outerHTML", HxIndicator: "#busy"}))
	for _, want := range []string{`type="submit"`, `class="primary outline"`, `disabled`, `hx-post="/save"`, `hx-target="#result"`, `hx-swap="outerHTML"`, `hx-indicator="#busy"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in %s", want, html)
		}
	}
	if strings.Contains(render(t, Button(ButtonProps{})), "hx-") {
		t.Fatal("empty HTMX attributes emitted")
	}
	if strings.Contains(render(t, Link(LinkProps{Href: "javascript:alert(1)"})), "javascript:") {
		t.Fatal("unsafe URL emitted")
	}
}

func TestCardSlotsAndControls(t *testing.T) {
	html := render(t, Card(CardProps{Header: templ.Raw("<h2>Title</h2>"), Body: templ.Raw("<p>Body</p>"), Footer: templ.Raw("End")}))
	for _, want := range []string{"<article", "<header><h2>Title</h2></header>", "<p>Body</p>", "<footer>End</footer>"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(render(t, Card(CardProps{})), "<header>") {
		t.Fatal("empty header emitted")
	}
	html = render(t, Input(InputProps{Label: "Email", Name: "email", Type: InputEmail, Required: true, Invalid: true, DescribedBy: "hint"}))
	for _, want := range []string{"<label>", `type="email"`, `required`, `aria-invalid="true"`, `aria-describedby="hint"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	html = render(t, Select(SelectProps{Label: "Choice", Value: "b", Options: []Option{{Value: "a", Label: "A"}, {Value: "b", Label: "B"}}}))
	if !strings.Contains(html, `value="b" selected`) {
		t.Fatal(html)
	}
	if !strings.Contains(render(t, Checkbox(CheckboxProps{Checked: true})), `value="on" checked`) {
		t.Fatal("checkbox defaults missing")
	}
	html = render(t, Form(FormProps{Method: MethodPost, Action: "/save", HxPost: "/save", HxGet: "/preview", HxTarget: "#result", HxSwap: "innerHTML", HxIndicator: "#busy"}))
	for _, want := range []string{`method="post"`, `action="/save"`, `hx-post="/save"`, `hx-get="/preview"`, `hx-target="#result"`, `hx-swap="innerHTML"`, `hx-indicator="#busy"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
}

func TestThemeToggle(t *testing.T) {
	html := render(t, ThemeToggle())
	for _, want := range []string{`action="/api/theme/toggle"`, `method="post"`, `type="submit"`, `hx-post="/api/theme/toggle"`, `hx-target="html"`, `hx-swap="outerHTML"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	html = render(t, ThemeToggle(ThemeToggleProps{Label: "Tema", CSRFToken: `"<>`}))
	if !strings.Contains(html, "Tema") || !strings.Contains(html, `name="csrf_token" value="&#34;&lt;&gt;"`) {
		t.Fatal(html)
	}
}
