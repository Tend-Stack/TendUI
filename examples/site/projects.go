package site

import (
	"cmp"
	"github.com/Tend-Stack/GothUI/v2026/components"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type dashboardOptions struct {
	Search, Status, Sort                          string
	Projects                                      []project
	Page, PageSize, Total, TotalPages, Start, End int
	query                                         url.Values
}

func cloneQuery(q url.Values) url.Values {
	result := make(url.Values, len(q))
	for key, values := range q {
		result[key] = append([]string(nil), values...)
	}
	return result
}

func (d dashboardOptions) PageURL(page int) string {
	q := cloneQuery(d.query)
	q.Set("q", d.Search)
	q.Set("status", d.Status)
	q.Set("sort", d.Sort)
	q.Set("per_page", strconv.Itoa(d.PageSize))
	q.Set("page", strconv.Itoa(page))
	return "/dashboard?" + q.Encode()
}

func readDashboard(q url.Values) dashboardOptions {
	d := dashboardOptions{Search: q.Get("q"), Status: choice(q.Get("status"), "all", "all", "active", "paused", "complete"), Sort: choice(q.Get("sort"), "name", "name", "name-desc", "progress", "progress-asc", "due", "due-desc"), PageSize: 2, query: cloneQuery(q)}
	if q.Get("per_page") == "5" {
		d.PageSize = 5
	}
	var filtered []project
	for _, p := range demoProjects {
		if (d.Status == "all" || p.Status == d.Status) && strings.Contains(strings.ToLower(p.Name+" "+p.Category), strings.ToLower(d.Search)) {
			filtered = append(filtered, p)
		}
	}
	// Sort the request's private slice; never mutate the shared sample data.
	sort.SliceStable(filtered, func(i, j int) bool {
		a, b := filtered[i], filtered[j]
		result := 0
		switch d.Sort {
		case "progress", "progress-asc":
			av, _ := strconv.Atoi(a.Progress)
			bv, _ := strconv.Atoi(b.Progress)
			result = cmp.Compare(av, bv)
		case "due", "due-desc":
			// These fixed demo dates all belong to November and use zero-padded days.
			result = strings.Compare(a.Due, b.Due)
		default:
			result = strings.Compare(a.Name, b.Name)
		}
		if d.Sort == "progress" || strings.HasSuffix(d.Sort, "-desc") {
			result = -result
		}
		if result == 0 {
			return a.Name < b.Name
		}
		return result < 0
	})
	d.Total = len(filtered)
	d.TotalPages = max(1, (d.Total+d.PageSize-1)/d.PageSize)
	requested, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		requested = 1
	}
	d.Page = max(1, min(requested, d.TotalPages))
	start := (d.Page - 1) * d.PageSize
	d.End = min(start+d.PageSize, d.Total)
	d.Projects = filtered[start:d.End]
	if d.Total > 0 {
		d.Start = start + 1
	}
	return d
}

func (d dashboardOptions) SortURL(column string) string {
	q := cloneQuery(d.query)
	next := column
	if d.Sort == column {
		next = column + "-desc"
		if column == "progress" {
			next = "progress-asc"
		}
	}
	q.Set("sort", next)
	q.Set("page", "1")
	return "/dashboard?" + q.Encode()
}

func (d dashboardOptions) StatusURL(status string) string {
	q := cloneQuery(d.query)
	q.Set("status", status)
	q.Set("page", "1")
	return "/dashboard?" + q.Encode()
}

func (d dashboardOptions) StatusTabs() []components.TabItem {
	var items []components.TabItem
	for _, item := range []struct{ value, label string }{{"all", "All projects"}, {"active", "Active"}, {"complete", "Completed"}, {"paused", "On hold"}} {
		items = append(items, components.TabItem{Label: item.label, Href: d.StatusURL(item.value), Current: d.Status == item.value})
	}
	return items
}

func (d dashboardOptions) sortColumn(label, key string, rowHeader bool) components.TableColumn {
	direction := components.SortNone
	if d.Sort == key || strings.HasPrefix(d.Sort, key+"-") {
		direction = components.SortAscending
		if d.Sort == "progress" || strings.HasSuffix(d.Sort, "-desc") {
			direction = components.SortDescending
		}
	}
	next, _ := url.Parse(d.SortURL(key))
	order := "ascending"
	if next.Query().Get("sort") == "progress" || strings.HasSuffix(next.Query().Get("sort"), "-desc") {
		order = "descending"
	}
	return components.TableColumn{Label: label, SortURL: d.SortURL(key), SortLabel: "Sort " + label + " " + order, Sort: direction, RowHeader: rowHeader}
}

func (d dashboardOptions) Table() components.DataTableProps {
	p := components.DataTableProps{Caption: "Sample project portfolio", Striped: true, Empty: emptyProjects(), Columns: []components.TableColumn{
		{Label: "Select"}, d.sortColumn("Project", "name", true), {Label: "Discipline"}, {Label: "Status"}, d.sortColumn("Due", "due", false), d.sortColumn("Progress", "progress", false),
	}}
	for _, project := range d.Projects {
		p.Rows = append(p.Rows, components.TableRow{Cells: []components.TableCell{
			{Content: selectProject(project)}, {Text: project.Name}, {Text: project.Category}, {Text: project.Status}, {Text: project.Due}, {Content: projectProgress(project)},
		}})
	}
	return p
}
