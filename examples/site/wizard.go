package site

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Tend-Stack/GothUI/v2026/components"
)

type wizardData struct {
	Name, Slug, Summary, Visibility, Accent string
	Updates                                 bool
	Step                                    int
	Done                                    bool
	Errors                                  []components.FieldError
}

var slugPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{1,39}$`)

func newWizard() wizardData {
	return wizardData{Step: 1, Visibility: "private", Accent: "indigo"}
}

func detailsErrors(d wizardData) []components.FieldError {
	var result []components.FieldError
	if n := utf8.RuneCountInString(d.Name); n < 2 || n > 60 {
		result = append(result, components.FieldError{ID: "brief-name", Message: "Use a project name between 2 and 60 characters."})
	}
	if !slugPattern.MatchString(d.Slug) {
		result = append(result, components.FieldError{ID: "brief-slug", Message: "Use a 2–40 character slug starting with a lowercase letter, followed by lowercase letters, numbers or hyphens."})
	}
	if utf8.RuneCountInString(d.Summary) > 400 {
		result = append(result, components.FieldError{ID: "brief-summary", Message: "Keep the summary to 400 characters or fewer."})
	}
	return result
}

func preferencesErrors(d wizardData) []components.FieldError {
	var result []components.FieldError
	if d.Visibility != "private" && d.Visibility != "public" {
		result = append(result, components.FieldError{ID: "brief-visibility", Message: "Choose private or public visibility."})
	}
	switch d.Accent {
	case "indigo", "forest", "rose", "neutral":
	default:
		result = append(result, components.FieldError{ID: "brief-accent", Message: "Choose one of the available accent colors."})
	}
	return result
}

// All draft state travels in the form; nothing is stored or trusted as validated.
func projectWizard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	d := newWizard()
	status := http.StatusOK
	switch r.Method {
	case http.MethodGet, http.MethodHead:
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 8192)
		if err := r.ParseForm(); err != nil {
			code := http.StatusBadRequest
			var sizeError *http.MaxBytesError
			if errors.As(err, &sizeError) {
				code = http.StatusRequestEntityTooLarge
			}
			http.Error(w, "invalid or oversized form", code)
			return
		}
		// Reject ambiguous or malformed draft state instead of choosing one duplicate.
		for _, key := range []string{"stage", "action", "name", "slug", "summary", "visibility", "accent", "updates"} {
			values := r.PostForm[key]
			if len(values) > 1 || (len(values) == 1 && !utf8.ValidString(values[0])) {
				http.Error(w, "invalid draft state", http.StatusBadRequest)
				return
			}
		}
		stage, err := strconv.Atoi(r.PostForm.Get("stage"))
		action := r.PostForm.Get("action")
		if err != nil || stage < 1 || stage > 3 ||
			!(action == "next" && stage < 3 || action == "back" && stage > 1 || action == "finish" && stage == 3) {
			http.Error(w, "invalid step transition", http.StatusBadRequest)
			return
		}
		updates := r.PostForm.Get("updates")
		if updates != "" && updates != "on" {
			http.Error(w, "invalid update preference", http.StatusBadRequest)
			return
		}
		d = wizardData{Name: strings.TrimSpace(r.PostForm.Get("name")), Slug: strings.TrimSpace(r.PostForm.Get("slug")), Summary: strings.TrimSpace(r.PostForm.Get("summary")), Visibility: r.PostForm.Get("visibility"), Accent: r.PostForm.Get("accent"), Updates: updates == "on", Step: stage}
		if action == "back" {
			d.Step-- // Preserve even incomplete edits; validation resumes on forward movement.
		} else if errs := detailsErrors(d); len(errs) > 0 {
			d.Step, d.Errors, status = 1, errs, http.StatusUnprocessableEntity
		} else if stage >= 2 {
			if errs := preferencesErrors(d); len(errs) > 0 {
				d.Step, d.Errors, status = 2, errs, http.StatusUnprocessableEntity
			} else if action == "finish" {
				d.Done = true
			} else {
				d.Step = 3
			}
		} else {
			d.Step = 2
		}
	default:
		w.Header().Set("Allow", "GET, HEAD, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	title := "Project brief · GothUI"
	if !d.Done {
		title = "Step " + strconv.Itoa(d.Step) + " of 3 · " + title
	}
	renderPage(w, r, title, wizardPage(d), status)
}

func (d wizardData) Steps() []components.StepItem {
	var steps []components.StepItem
	for i, label := range []string{"Project details", "Preferences", "Review"} {
		steps = append(steps, components.StepItem{Label: label, Current: !d.Done && i+1 == d.Step, Completed: d.Done || i+1 < d.Step})
	}
	return steps
}

func (d wizardData) Field(key string) components.FieldProps {
	p := components.FieldProps{ID: "brief-" + key}
	switch key {
	case "name":
		p.Label, p.Description, p.Required = "Project name", "Required. Use 2–60 characters and fictional details for this demo.", true
	case "slug":
		p.Label, p.Description, p.Required = "Project URL slug", "Required. 2–40 lowercase letters, numbers or hyphens, starting with a letter. The /projects/ prefix is supplied for you.", true
	case "summary":
		p.Label, p.Description = "Short summary", "Optional. Up to 400 characters."
	case "visibility":
		p.Label, p.Description, p.Required = "Visibility", "Required. This is a preview preference; no project will be published.", true
	case "accent":
		p.Label, p.Description, p.Required = "Accent color", "Required. Choose the proposed visual accent for the brief.", true
	}
	for _, e := range d.Errors {
		if e.ID == p.ID {
			p.Error = e.Message
		}
	}
	return p
}
