package site

import "net/http"

func confirmDemo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	if err := r.ParseForm(); err != nil || r.PostForm.Get("decision") != "approve" {
		http.Error(w, "Choose the demo approval action", http.StatusBadRequest)
		return
	}
	// This demonstration performs no mutation, storage, or external requests.
	http.Redirect(w, r, "/interactions?confirmed=1", http.StatusSeeOther)
}
