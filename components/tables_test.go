package components

import (
	"bytes"
	"context"
	"github.com/a-h/templ"
	"strings"
	"testing"
)

func TestDataTableSemanticsAndEscaping(t *testing.T) {
	p := DataTableProps{Caption: "Projects <today>", Striped: true, Compact: true,
		Columns: []TableColumn{{Label: "Name", RowHeader: true, Sort: SortAscending, SortURL: "/projects?sort=name&dir=desc", SortLabel: "Sort Name descending"}, {Label: "Total", Alignment: AlignEnd}},
		Rows:    []TableRow{{Cells: []TableCell{{Text: "<unsafe>"}, {Content: Badge(BadgeProps{Text: "12"})}}, Selected: true}},
	}
	html := render(t, DataTable(p))
	for _, want := range []string{`<caption>Projects &lt;today&gt;</caption>`, `scope="row"`, `scope="col"`, `aria-sort="ascending"`, `data-align="end"`, `data-selected="true"`, `tabindex="0"`, `role="region"`, `aria-label="Sort Name descending"`, `&lt;unsafe&gt;`, `class="badge"`, "table-striped", "table-compact"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Count(html, "aria-sort=") != 1 || strings.Contains(html, `role="grid"`) {
		t.Fatal("incorrect table accessibility contract")
	}
	p.Columns[0].SortURL = "javascript:alert(1)"
	html = render(t, DataTable(p))
	if strings.Contains(html, "javascript:") {
		t.Fatal("unsafe URL escaped templ filtering")
	}
}

func TestDataTableEmptyLoadingAndInvalidShapes(t *testing.T) {
	p := DataTableProps{Caption: "Sample", Columns: []TableColumn{{Label: "Name"}, {Label: "Status"}}, Empty: EmptyState(EmptyStateProps{Title: "Empty sample"})}
	html := render(t, DataTable(p))
	if !strings.Contains(html, `colspan="2"`) || !strings.Contains(html, "Empty sample") {
		t.Fatal("empty table contract")
	}
	p.Loading, p.LoadingLabel = true, "Waiting <sample>"
	p.Rows = []TableRow{{Cells: []TableCell{{Text: "stale content"}}}}
	html = render(t, DataTable(p))
	if !strings.Contains(html, `aria-busy="true"`) || !strings.Contains(html, `role="status"`) || !strings.Contains(html, "Waiting &lt;sample&gt;") || strings.Contains(html, "stale content") {
		t.Fatal("loading must replace stale rows and expose a status")
	}
	for _, invalid := range []DataTableProps{
		{}, {Caption: "Missing columns"}, {Columns: []TableColumn{{Label: "Name"}}},
		{Caption: "Blank label", Columns: []TableColumn{{Label: " "}}},
		{Caption: "Short row", Columns: []TableColumn{{Label: "Name"}}, Rows: []TableRow{{}}},
		{Caption: "Long row", Columns: []TableColumn{{Label: "Name"}}, Rows: []TableRow{{Cells: []TableCell{{}, {}}}}},
		{Caption: "Bad sort", Columns: []TableColumn{{Label: "Name", Sort: "bad"}}},
		{Caption: "Two sorts", Columns: []TableColumn{{Label: "Name", Sort: SortAscending}, {Label: "Status", Sort: SortDescending}}},
		{Caption: "Two row headers", Columns: []TableColumn{{Label: "Name", RowHeader: true}, {Label: "Status", RowHeader: true}}},
	} {
		var b bytes.Buffer
		if err := DataTable(invalid).Render(context.Background(), &b); err == nil {
			t.Errorf("accepted invalid table: %+v", invalid)
		}
	}
}

func TestWayfindingUsesNativeNavigation(t *testing.T) {
	html := render(t, Tabs(TabsProps{Items: []TabItem{{Label: "Current <page>", Current: true}, {Label: "Next", Href: "/next"}, {Label: "Later", Disabled: true}}}))
	for _, want := range []string{`aria-current="page"`, `href="/next"`, `aria-disabled="true"`, "Current &lt;page&gt;"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, "tablist") || strings.Contains(html, `role="tab"`) || strings.Contains(html, "aria-selected") {
		t.Fatal("page links must not advertise scripted tabs")
	}
	html = render(t, Steps(StepsProps{Items: []StepItem{{Label: "Start", Href: "/start", Completed: true}, {Label: "Review", Current: true}, {Label: "Finish"}}}))
	if strings.Count(html, `aria-current="step"`) != 1 || !strings.Contains(html, "Completed") || strings.Count(html, "<a ") != 1 {
		t.Fatal("steps lost progress or unavailable destination semantics")
	}
	html = render(t, Disclosure(DisclosureProps{Summary: "Show <details>", Name: "faq", Open: true}))
	for _, want := range []string{"<details", "<summary>Show &lt;details&gt;</summary>", `name="faq"`, " open"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	if strings.Contains(html, "aria-expanded") {
		t.Fatal("native state must not be hard-coded")
	}
}

func TestWayfindingRejectsContradictoryState(t *testing.T) {
	for _, c := range []templ.Component{
		Tabs(TabsProps{Items: []TabItem{{Label: "A", Current: true}, {Label: "B", Current: true}}}),
		Tabs(TabsProps{Items: []TabItem{{Label: "A", Current: true, Disabled: true}}}),
		Tabs(TabsProps{Items: []TabItem{{Label: "A"}}}),
		Tabs(TabsProps{Items: []TabItem{{Href: "/"}}}),
		Steps(StepsProps{Items: []StepItem{{Label: "A", Current: true, Completed: true}}}),
		Steps(StepsProps{Items: []StepItem{{Label: "A", Current: true}, {Label: "B", Current: true}}}),
	} {
		var b bytes.Buffer
		if err := c.Render(context.Background(), &b); err == nil {
			t.Fatal("accepted contradictory state")
		}
	}
	if render(t, Tabs(TabsProps{})) != "" || render(t, Steps(StepsProps{})) != "" {
		t.Fatal("empty navigation should render nothing")
	}
}
