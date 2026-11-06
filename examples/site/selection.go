package site

import "net/http"

// Review is deliberately non-mutating. Real bulk actions need authorization per ID.
func reviewSelection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid or oversized form", http.StatusBadRequest)
		return
	}
	names := r.PostForm["project"]
	if len(names) > len(demoProjects) {
		http.Error(w, "too many selections", http.StatusBadRequest)
		return
	}
	valid := map[string]bool{}
	for _, p := range demoProjects {
		valid[p.Name] = true
	}
	selected := map[string]bool{}
	for _, name := range names {
		if !valid[name] {
			http.Error(w, "unknown project", http.StatusBadRequest)
			return
		}
		selected[name] = true
	}
	var projects []project
	for _, p := range demoProjects {
		if selected[p.Name] {
			projects = append(projects, p)
		}
	}
	status, message := http.StatusOK, ""
	if len(projects) == 0 {
		status, message = http.StatusUnprocessableEntity, "Select at least one project on the dashboard, then try again."
	}
	renderPage(w, r, "Selection review · GothUI", selectionReceipt(projects, message), status)
}
