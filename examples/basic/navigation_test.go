package main

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestProjectSortingPaginationAndState(t *testing.T) {
	original := append([]project(nil), demoProjects...)
	q := url.Values{"q": {"design"}, "status": {"all"}, "sort": {"progress"}, "per_page": {"2"}, "page": {"999"}, "theme": {"dark"}, "preset": {"forest"}}
	d := readDashboard(q)
	if d.Total != 2 || d.Page != 1 || d.Projects[0].Name != "Brand guidelines" {
		t.Fatalf("bad filtered page: %+v", d)
	}
	u, err := url.Parse(d.PageURL(2))
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"q", "status", "sort", "per_page", "theme", "preset"} {
		if u.Query().Get(key) != q.Get(key) {
			t.Errorf("lost query %s", key)
		}
	}
	if q.Get("page") != "999" {
		t.Fatal("URL generation mutated input")
	}
	d = readDashboard(url.Values{"page": {"3"}})
	if d.Total != 5 || d.TotalPages != 3 || d.Page != 3 || d.Start != 5 || d.End != 5 || len(d.Projects) != 1 {
		t.Fatalf("bad final page: %+v", d)
	}
	d = readDashboard(url.Values{"sort": {"due"}, "per_page": {"5"}})
	if len(d.Projects) != 5 || d.Projects[0].Name != "Team handbook" {
		t.Fatal("due date sort failed")
	}
	for _, value := range []string{"-1", "not-a-number", "99999999999999999999999999"} {
		d = readDashboard(url.Values{"page": {value}, "q": {"no-match"}})
		if d.Page != 1 || len(d.Projects) != 0 || d.Start != 0 {
			t.Fatal("invalid page not normalized")
		}
	}
	if !reflect.DeepEqual(original, demoProjects) {
		t.Fatal("shared projects were reordered")
	}
}

func TestDashboardNavigationHTTP(t *testing.T) {
	h := newHandler()
	for _, path := range []string{"/dashboard", "/dashboard?page=3", "/dashboard?sort=progress&per_page=2&theme=dark"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != 200 || !strings.Contains(w.Body.String(), `aria-label="Pagination"`) || !strings.Contains(w.Body.String(), `aria-label="Breadcrumb"`) {
			t.Errorf("navigation missing on %s", path)
		}
	}
}

func TestInteractionsAndConfirmation(t *testing.T) {
	h := newHandler()
	for _, path := range []string{"/interactions", "/interactions?view=confirm", "/interactions?confirmed=1"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		html := w.Body.String()
		if w.Code != 200 || !strings.Contains(html, `command="show-modal"`) || !strings.Contains(html, `popover="manual"`) || strings.Contains(html, "<script") {
			t.Fatal("invalid interaction page")
		}
		if strings.Contains(path, "view=confirm") && !strings.Contains(html, `id="inline-confirm"`) {
			t.Fatal("fallback missing")
		}
	}
	for _, tc := range []struct {
		method, body, origin string
		code                 int
	}{
		{"POST", "decision=approve", "", 303}, {"POST", "decision=reject", "", 400}, {"POST", strings.Repeat("x", 3000), "", 400}, {"GET", "", "", 405}, {"POST", "decision=approve", "https://evil.example", 403},
	} {
		r := httptest.NewRequest(tc.method, "https://example.com/interactions/confirm", strings.NewReader(tc.body))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("Origin", tc.origin)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != tc.code {
			t.Fatalf("confirmation: %d, want %d", w.Code, tc.code)
		}
		if tc.code == 303 && w.Header().Get("Location") != "/interactions?confirmed=1" {
			t.Fatal("incorrect confirmation redirect")
		}
	}
}
