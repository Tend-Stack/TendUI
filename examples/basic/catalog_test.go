package main

import (
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestCatalogCoversPublicTemplComponents(t *testing.T) {
	files, err := filepath.Glob("../../components/*.templ")
	if err != nil || len(files) == 0 {
		t.Fatalf("component sources unavailable: %v", err)
	}
	exported := map[string]bool{}
	declaration := regexp.MustCompile(`(?m)^templ ([A-Z][A-Za-z0-9_]*)\(`)
	for _, file := range files {
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range declaration.FindAllStringSubmatch(string(source), -1) {
			exported[m[1]] = true
		}
	}
	seen := map[string]bool{}
	for _, e := range catalogEntries {
		if !exported[e.Name] || seen[e.Name] || e.Example == nil {
			t.Fatalf("invalid/duplicate catalog entry %s", e.Name)
		}
		seen[e.Name] = true
		if e.Summary == "" || e.Notes == "" {
			t.Errorf("missing guidance for %s", e.Name)
		}
		source, err := catalogSource(e)
		if err != nil || !strings.Contains(source, "templ Example"+e.Name+"()") || !strings.HasPrefix(source, "package examples\n") {
			t.Fatalf("missing reusable source for %s: %v", e.Name, err)
		}
		demo, err := url.Parse(e.Demo)
		if err != nil || !demoPath(demo.Path) {
			t.Errorf("unknown demo for %s", e.Name)
		}
	}
	for name := range exported {
		if !seen[name] {
			t.Errorf("public component %s is absent from catalog", name)
		}
	}
}

func TestCatalogEveryDetailPreviewAndDownload(t *testing.T) {
	h := newHandler()
	for _, entry := range catalogEntries {
		t.Run(entry.Name, func(t *testing.T) {
			t.Parallel()
			for _, suffix := range []string{"", "/preview", "/example.templ"} {
				path := entry.URL() + suffix
				w := httptest.NewRecorder()
				h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
				if w.Code != 200 {
					t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
				}
				if suffix == "/example.templ" {
					want, _ := catalogSource(entry)
					if w.Body.String() != want || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" || w.Header().Get("Content-Disposition") != `attachment; filename="`+entry.ID()+`_example.templ"` {
						t.Fatal("download drift or unsafe response")
					}
				} else {
					body := w.Body.String()
					if strings.Count(body, "<html ") != 1 || strings.Count(body, "<main") != 1 || strings.Contains(body, "<script") || !strings.Contains(w.Header().Get("Content-Security-Policy"), "default-src 'none'") {
						t.Fatalf("document/runtime contract failed for %s", path)
					}
					if suffix == "" && (!strings.Contains(body, "package examples") || !strings.Contains(body, "Templ example")) {
						t.Fatal("missing source on detail page")
					}
				}
				head := httptest.NewRecorder()
				h.ServeHTTP(head, httptest.NewRequest("HEAD", path, nil))
				if head.Code != 200 || head.Body.Len() != 0 || head.Header().Get("Content-Type") != w.Header().Get("Content-Type") {
					t.Fatalf("HEAD contract %s", path)
				}
			}
		})
	}
}

func TestCatalogSearchPaginationAndQuerySafety(t *testing.T) {
	all := readCatalog(url.Values{})
	if all.Total != len(catalogEntries) || len(all.Entries) > 12 || all.Page != 1 {
		t.Fatal("catalog totals")
	}
	seen := map[string]bool{}
	for page := 1; page <= all.Pages; page++ {
		u, _ := url.Parse(all.PageURL(page))
		d := readCatalog(u.Query())
		for _, e := range d.Entries {
			if seen[e.Name] {
				t.Fatal("duplicate across pages")
			}
			seen[e.Name] = true
		}
	}
	if len(seen) != len(catalogEntries) {
		t.Fatal("pagination omitted a component")
	}
	filtered := readCatalog(url.Values{"category": {"Forms"}, "q": {"FIELD"}, "theme": {"dark"}})
	if filtered.Total == 0 {
		t.Fatal("case-insensitive search failed")
	}
	for _, e := range filtered.Entries {
		if e.Category != "Forms" {
			t.Fatal("category mismatch")
		}
	}
	u, _ := url.Parse(filtered.PageURL(2))
	if u.Query().Get("q") != "FIELD" || u.Query().Get("category") != "Forms" || u.Query().Get("theme") != "dark" {
		t.Fatal("pagination lost filters")
	}
	empty := readCatalog(url.Values{"q": {"not-a-component-at-all"}, "page": {"999"}})
	if empty.Page != 1 || empty.Total != 0 || len(empty.Entries) != 0 {
		t.Fatal("empty pagination bounds")
	}
	long := readCatalog(url.Values{"q": {strings.Repeat("語", 300)}})
	if utf8.RuneCountInString(long.Query) != 120 {
		t.Fatal("search not bounded")
	}
	h := newHandler()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/catalog?q="+url.QueryEscape(`<script>alert(1)</script>`), nil))
	if w.Code != 200 || strings.Contains(w.Body.String(), "<script>") || !strings.Contains(w.Body.String(), "&lt;script&gt;") {
		t.Fatal("search escaping failed")
	}
}

func TestCatalogRoutesAndDemoAssets(t *testing.T) {
	h := newHandler()
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"GET", "/catalog", 200}, {"GET", "/catalog/unknown", 404}, {"GET", "/catalog/button/unknown", 404},
		{"POST", "/catalog", 405}, {"POST", "/catalog/button/example.templ", 405},
		{"GET", "/demo.css", 200}, {"HEAD", "/demo.css", 200}, {"POST", "/demo.css", 405},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.status {
			t.Errorf("%s %s: %d", tc.method, tc.path, w.Code)
		}
		if tc.method == "HEAD" && w.Body.Len() != 0 {
			t.Fatal("HEAD body")
		}
		if tc.path == "/demo.css" && tc.status == 200 && w.Header().Get("Content-Type") != "text/css; charset=utf-8" {
			t.Fatal("CSS MIME")
		}
	}
	if !catalogPagePath("/catalog/button") || !catalogPagePath("/catalog/layout/preview") || catalogPagePath("/catalog/button/example.templ") || catalogPagePath("/catalog/not-real") {
		t.Fatal("catalog referer allowlist")
	}
	if strings.Contains(string(demoCSS), "@import") || strings.Contains(string(demoCSS), "url(") {
		t.Fatal("external example CSS dependency")
	}
}
