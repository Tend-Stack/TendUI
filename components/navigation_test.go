package components

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"testing"
)

func TestBreadcrumbs(t *testing.T) {
	if render(t, Breadcrumbs(BreadcrumbsProps{})) != "" {
		t.Fatal("empty trail should not render a landmark")
	}
	html := render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{{Label: "Home", Href: "/"}, {Label: "<Current>", Href: "/ignored"}}}))
	for _, want := range []string{`aria-label="Breadcrumb"`, `<ol>`, `href="/"`, `aria-current="page"`, `&lt;Current&gt;`, `aria-hidden="true"`} {
		if !strings.Contains(html, want) {
			t.Error(want)
		}
	}
	if strings.Contains(html, `href="/ignored"`) {
		t.Fatal("current page should not link to itself")
	}
	html = render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{{Label: "Bad", Href: "javascript:bad()"}, {Label: "Here"}}}))
	if strings.Contains(html, "javascript:") {
		t.Fatal("unsafe link")
	}
}

func TestPaginationBoundariesAndEscaping(t *testing.T) {
	for _, tc := range []struct{ page, total int }{{1, 3}, {3, 3}, {-10, 3}, {999, 3}} {
		called := []int{}
		html := render(t, Pagination(PaginationProps{Page: tc.page, TotalPages: tc.total, URLForPage: func(page int) string {
			called = append(called, page)
			return "/items?status=active&page=" + strconv.Itoa(page)
		}}))
		if strings.Count(html, `aria-current="page"`) != 1 || !strings.Contains(html, "&amp;page=") {
			t.Fatal(html)
		}
		for _, page := range called {
			if page < 1 || page > tc.total {
				t.Fatalf("out-of-range URL call %d", page)
			}
		}
		if tc.page <= 1 && strings.Contains(html, `rel="prev"`) {
			t.Fatal("previous link on first page")
		}
		if tc.page >= tc.total && strings.Contains(html, `rel="next"`) {
			t.Fatal("next link on last page")
		}
	}
	if render(t, Pagination(PaginationProps{TotalPages: 1})) != "" {
		t.Fatal("unnecessary pagination")
	}
	var buf bytes.Buffer
	if err := Pagination(PaginationProps{TotalPages: 3}).Render(context.Background(), &buf); err == nil {
		t.Fatal("missing URL callback was accepted")
	}
	html := render(t, Pagination(PaginationProps{Page: 2, TotalPages: 3, URLForPage: func(int) string { return "javascript:bad()" }}))
	if strings.Contains(html, "javascript:") {
		t.Fatal("unsafe URL")
	}
}

func TestPaginationWindowIsBounded(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	for _, total := range []int{0, 1, 2, 3, 5, 20, maxInt} {
		for _, page := range []int{-1, 1, 2, 10, total, maxInt} {
			window := paginationWindow(page, total)
			if len(window) > 9 {
				t.Fatalf("unbounded window: %v", window)
			}
			last := 0
			found := false
			for _, n := range window {
				if n == 0 {
					continue
				}
				if n < 1 || n > total || n <= last {
					t.Fatalf("invalid window %v for %d/%d", window, page, total)
				}
				last = n
				if n == currentPage(page, total) {
					found = true
				}
			}
			if total > 0 && (!found || window[0] != 1 || window[len(window)-1] != total) {
				t.Fatalf("missing boundary/current page: %v", window)
			}
		}
	}
}
