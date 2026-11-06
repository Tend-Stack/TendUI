package main

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestDashboardHeaderSortingAndTabs(t *testing.T) {
	for _, tc := range []struct{ sort, first string }{
		{"name", "Billing integration"}, {"name-desc", "Team handbook"}, {"progress", "Brand guidelines"}, {"progress-asc", "Billing integration"}, {"due", "Team handbook"}, {"due-desc", "Billing integration"},
	} {
		d := readDashboard(url.Values{"sort": {tc.sort}, "per_page": {"5"}})
		if d.Projects[0].Name != tc.first {
			t.Errorf("%s: first=%s", tc.sort, d.Projects[0].Name)
		}
	}
	q := url.Values{"q": {"design"}, "sort": {"name"}, "page": {"2"}, "per_page": {"5"}, "theme": {"dark"}, "motion": {"off"}}
	d := readDashboard(q)
	sortURL, _ := url.Parse(d.SortURL("name"))
	if sortURL.Query().Get("sort") != "name-desc" || sortURL.Query().Get("page") != "1" {
		t.Fatal("sort must toggle direction and reset page")
	}
	for _, key := range []string{"q", "per_page", "theme", "motion"} {
		if sortURL.Query().Get(key) != q.Get(key) {
			t.Errorf("lost %s", key)
		}
	}
	statusURL, _ := url.Parse(d.StatusURL("complete"))
	if statusURL.Query().Get("status") != "complete" || statusURL.Query().Get("q") != "design" || q.Get("page") != "2" {
		t.Fatal("tab state mutated or lost")
	}
	w := httptest.NewRecorder()
	newHandler().ServeHTTP(w, httptest.NewRequest("GET", "/dashboard?sort=name-desc", nil))
	for _, want := range []string{`aria-sort="descending"`, `aria-label="Project status"`, `action="/dashboard/selection"`, `name="project"`, `aria-label="Sort Project ascending"`} {
		if !strings.Contains(w.Body.String(), want) {
			t.Errorf("missing %s", want)
		}
	}
}

func TestSelectionReviewValidationAndNoMutation(t *testing.T) {
	original := append([]project(nil), demoProjects...)
	h := newHandler()
	for _, tc := range []struct {
		method, body string
		status       int
		want         string
	}{
		{"POST", "project=Studio+website&project=Studio+website", 200, "1 sample projects selected"},
		{"POST", "project=Studio+website&project=Customer+portal", 200, "2 sample projects selected"},
		{"POST", "", 422, "Select at least one"}, {"POST", "project=unknown", 400, "unknown project"},
		{"POST", strings.Repeat("project=Studio+website&", 6), 400, "too many"},
		{"POST", "project=" + strings.Repeat("x", 5000), 400, "oversized"}, {"POST", "project=%ZZ", 400, "invalid"},
		{"GET", "", 405, "method not allowed"},
	} {
		r := httptest.NewRequest(tc.method, "/dashboard/selection", strings.NewReader(tc.body))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.want) || len(w.Result().Cookies()) != 0 {
			t.Fatalf("review %q: %d %s", tc.body[:min(len(tc.body), 60)], w.Code, w.Body.String())
		}
	}
	r := httptest.NewRequest("POST", "/dashboard/selection", strings.NewReader("project=Studio+website"))
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("cross-origin action allowed")
	}
	if !reflect.DeepEqual(original, demoProjects) {
		t.Fatal("review mutated shared records")
	}
}

func TestPatternsNavigationAndStates(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct{ path, want string }{
		{"/patterns?section=composition", `class="layout-stack"`},
		{"/patterns", `aria-current="step"`}, {"/patterns?step=3", "Compose it with typed props"},
		{"/patterns?section=tables", `colspan="2"`}, {"/patterns?section=tables&state=loading", `aria-busy="true"`},
		{"/patterns?section=disclosures", `name="pattern-faq"`}, {"/patterns?section=invalid&step=999", "Start with the interaction"},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != 200 || !strings.Contains(w.Body.String(), tc.want) || strings.Contains(w.Body.String(), "<script") {
			t.Fatalf("bad pattern %s", tc.path)
		}
	}
	for _, method := range []string{"HEAD", "POST"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, "/patterns", nil))
		if method == "HEAD" && (w.Code != 200 || w.Body.Len() != 0) {
			t.Fatal("bad HEAD")
		}
		if method == "POST" && w.Code != 405 {
			t.Fatal("unexpected mutation method")
		}
	}
}

func TestCompositionPatternNativeActions(t *testing.T) {
	h := newHandler()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/patterns?section=composition&density=compact", nil))
	for _, want := range []string{`data-density="compact"`, `id="composition-title"`, `data-layout-justify="between"`, `data-layout-justify="end"`, `action="/catalog" method="get"`, `name="q"`, `type="reset"`, `type="submit"`, `/catalog/stack`, `/catalog/cluster`} {
		if w.Code != 200 || !strings.Contains(w.Body.String(), want) {
			t.Errorf("missing composition contract %s", want)
		}
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/catalog?q=Stack", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `href="/catalog/stack"`) {
		t.Fatal("composition search destination")
	}
}
