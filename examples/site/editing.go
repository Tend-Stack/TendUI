package site

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Tend-Stack/TendUI/v2026/components"
	"github.com/Tend-Stack/TendUI/v2026/render"
)

const profileCookie = "tendui_demo_profile"

type demoProfile struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type editingData struct {
	Profile demoProfile
	Editing bool
	Notice  string
	Errors  []components.FieldError
}

func defaultProfile() demoProfile {
	return demoProfile{Name: "Alex Morgan", Bio: "Building thoughtful things with Go."}
}

func profileErrors(p demoProfile) []components.FieldError {
	var errs []components.FieldError
	if !utf8.ValidString(p.Name) || strings.TrimSpace(p.Name) == "" || utf8.RuneCountInString(p.Name) > 60 {
		errs = append(errs, components.FieldError{ID: "profile-name", Message: "Use a display name between 1 and 60 characters."})
	}
	if !utf8.ValidString(p.Bio) || utf8.RuneCountInString(p.Bio) > 240 {
		errs = append(errs, components.FieldError{ID: "profile-bio", Message: "Keep the introduction to 240 characters or fewer."})
	}
	return errs
}

func readProfile(r *http.Request) demoProfile {
	c, err := r.Cookie(profileCookie)
	if err != nil || len(c.Value) > 3072 {
		return defaultProfile()
	}
	b, err := base64.RawURLEncoding.DecodeString(c.Value)
	var p demoProfile
	if err != nil || json.Unmarshal(b, &p) != nil || len(profileErrors(p)) != 0 {
		return defaultProfile()
	}
	return p
}

func editingDemo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	d := editingData{Profile: readProfile(r), Editing: r.URL.Query().Get("edit") == "1"}
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		switch r.URL.Query().Get("notice") {
		case "saved":
			d.Notice = "Sample profile saved in this browser session."
		case "reset":
			d.Notice = "The sample profile has been restored."
		}
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid or oversized form", http.StatusBadRequest)
			return
		}
		intent := r.PostForm.Get("intent")
		if intent != "save" && intent != "reset" {
			http.Error(w, "unknown action", http.StatusBadRequest)
			return
		}
		cookie := &http.Cookie{Name: profileCookie, Path: "/", HttpOnly: true,
			SameSite: http.SameSiteStrictMode, Secure: r.TLS != nil}
		notice := "saved"
		if intent == "reset" {
			cookie.MaxAge = -1
			d.Profile = defaultProfile()
			d.Notice = "The sample profile has been restored."
			notice = "reset"
		} else {
			d.Profile = demoProfile{Name: strings.TrimSpace(r.PostForm.Get("name")), Bio: strings.TrimSpace(r.PostForm.Get("bio"))}
			d.Errors = profileErrors(d.Profile)
			if len(d.Errors) != 0 {
				d.Editing = true
				status := http.StatusUnprocessableEntity
				if render.IsFragment(r) {
					// Default HTMX swaps 2xx responses, not 422. The returned form
					// contains errors; no cookie is written on invalid submissions.
					status = http.StatusOK
				}
				renderEditing(w, r, d, status)
				return
			}
			b, _ := json.Marshal(d.Profile) // Only validated strings; cannot fail.
			cookie.Value = base64.RawURLEncoding.EncodeToString(b)
			d.Notice = "Sample profile saved in this browser session."
		}
		http.SetCookie(w, cookie)
		d.Editing = false
		if !render.IsFragment(r) {
			http.Redirect(w, r, "/editing?notice="+notice, http.StatusSeeOther)
			return
		}
	default:
		w.Header().Set("Allow", "GET, HEAD, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	renderEditing(w, r, d, http.StatusOK)
}

func renderEditing(w http.ResponseWriter, r *http.Request, d editingData, status int) {
	renderResponse(w, r, "Inline editing · TendUI", editingPage(d), profilePanel(d), status)
}

func profileError(d editingData, id string) string {
	for _, e := range d.Errors {
		if e.ID == id {
			return e.Message
		}
	}
	return ""
}
