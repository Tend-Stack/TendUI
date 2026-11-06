// Package render serves complete HTML pages or optional HTMX fragments.
// It uses only the standard library; templ components satisfy Component directly.
package render

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
)

// Component is implemented by templ.Component and other streaming HTML renderers.
type Component interface {
	Render(context.Context, io.Writer) error
}

// IsFragment reports whether this request asks for a local HTMX update.
// Boosted navigation and history restoration receive a complete page.
// These headers control presentation only; they are not an authorization boundary.
func IsFragment(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true" &&
		r.Header.Get("HX-History-Restore-Request") != "true" &&
		r.Header.Get("HX-Boosted") != "true"
}

// HTML renders fragment for IsFragment requests, or page otherwise. A nil
// fragment falls back to page. It merges Vary without removing existing values,
// buffers before committing the status, and suppresses the body for HEAD.
//
// status must permit an HTML body (200..599, excluding 204, 205 and 304).
// On a render/configuration error it sends a generic 500 and returns the error
// for logging. A transport error is returned after the response was committed;
// callers must not write another response. Cache policy, CSP, authentication,
// CSRF protection and validation status policy belong to the application.
func HTML(w http.ResponseWriter, r *http.Request, status int, page, fragment Component) error {
	for _, name := range []string{"HX-Request", "HX-History-Restore-Request", "HX-Boosted"} {
		vary(w.Header(), name)
	}
	selected := page
	if fragment != nil && IsFragment(r) {
		selected = fragment
	}
	var buf bytes.Buffer
	var err error
	switch {
	case status < 200 || status > 599 || status == 204 || status == 205 || status == 304:
		err = errors.New("render: status must permit an HTML body")
	case selected == nil:
		err = errors.New("render: missing component")
	default:
		err = selected.Render(r.Context(), &buf)
	}
	if err != nil {
		w.Header().Del("Content-Length")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		if r.Method != http.MethodHead {
			_, _ = io.WriteString(w, "unable to render page\n")
		}
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if r.Method == http.MethodHead {
		return nil
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func vary(h http.Header, name string) {
	for _, line := range h.Values("Vary") {
		for _, field := range strings.Split(line, ",") {
			field = strings.TrimSpace(field)
			if field == "*" || strings.EqualFold(field, name) {
				return
			}
		}
	}
	h.Add("Vary", name)
}
