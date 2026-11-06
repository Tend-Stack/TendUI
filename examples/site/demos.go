package site

import (
	"errors"
	"github.com/Tend-Stack/GothUI/v2026/components"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
)

type playgroundOptions struct {
	Appearance                   components.Appearance
	Label, Variant, State, Theme string
	Header, Footer, Checked      bool
	Clicked                      bool
}

func choice(value, fallback string, allowed ...string) string {
	for _, option := range allowed {
		if value == option {
			return value
		}
	}
	return fallback
}

func readPlayground(q url.Values) playgroundOptions {
	label := q.Get("label")
	if label == "" {
		label = "Try this button"
	}
	if len(label) > 200 {
		label = string([]rune(label)[:min(200, len([]rune(label)))])
	}
	return playgroundOptions{
		Label: label, Variant: choice(q.Get("variant"), "primary", "primary", "outline"),
		State:   choice(q.Get("state"), "normal", "normal", "disabled", "invalid", "readonly"),
		Theme:   choice(q.Get("theme"), "saved", "saved", "light", "dark", "system"),
		Header:  !q.Has("apply") || q.Get("header") == "on",
		Footer:  !q.Has("apply") || q.Get("footer") == "on",
		Checked: q.Get("checked") == "on", Clicked: q.Get("clicked") == "yes",
	}
}

type formData struct {
	Name, Email, Role, Message  string
	Consent, Submitted, Success bool
	Errors                      map[string]string
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	if err := r.ParseForm(); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "Demo form is limited to 16 KB", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
		}
		return
	}
	d := formData{
		Name: strings.TrimSpace(r.PostForm.Get("name")), Email: strings.TrimSpace(r.PostForm.Get("email")),
		Role: r.PostForm.Get("role"), Message: strings.TrimSpace(r.PostForm.Get("message")),
		Consent: r.PostForm.Get("consent") == "on", Submitted: true, Errors: map[string]string{},
	}
	if d.Name == "" || len([]rune(d.Name)) > 80 {
		d.Errors["name"] = "Enter a name between 1 and 80 characters."
	}
	address, err := mail.ParseAddress(d.Email)
	if err != nil || address.Address != d.Email {
		d.Errors["email"] = "Enter an email address, such as alex@example.com."
	}
	if choice(d.Role, "", "developer", "designer", "exploring") == "" {
		d.Errors["role"] = "Choose what best describes you."
	}
	if len([]rune(d.Message)) > 1000 {
		d.Errors["message"] = "Keep your message under 1,000 characters."
	}
	if !d.Consent {
		d.Errors["consent"] = "Check the demo consent box to continue."
	}
	status := http.StatusUnprocessableEntity
	if len(d.Errors) == 0 {
		d.Success = true
		status = http.StatusOK
	}
	renderPage(w, r, "Form lab · GothUI", formLab(d), status)
}

type project struct {
	Name, Category, Status, Due string
	Progress                    string
}

var demoProjects = []project{
	{"Studio website", "Design", "active", "Nov 12", "72"},
	{"Customer portal", "Product", "active", "Nov 18", "45"},
	{"Brand guidelines", "Design", "complete", "Nov 04", "100"},
	{"Billing integration", "Engineering", "paused", "Nov 24", "28"},
	{"Team handbook", "Operations", "complete", "Nov 02", "100"},
}

func fieldHint(message, fallback string) string {
	if message != "" {
		return message
	}
	return fallback
}

func formErrors(d formData) []components.FieldError {
	var result []components.FieldError
	for _, id := range []string{"name", "email", "role", "message", "consent"} {
		if message := d.Errors[id]; message != "" {
			result = append(result, components.FieldError{ID: id, Message: message})
		}
	}
	return result
}
