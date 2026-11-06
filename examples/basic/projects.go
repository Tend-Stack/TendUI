package main

import (
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
	d := dashboardOptions{Search: q.Get("q"), Status: choice(q.Get("status"), "all", "all", "active", "paused", "complete"), Sort: choice(q.Get("sort"), "name", "name", "progress", "due"), PageSize: 2, query: cloneQuery(q)}
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
		switch d.Sort {
		case "progress":
			a, _ := strconv.Atoi(filtered[i].Progress)
			b, _ := strconv.Atoi(filtered[j].Progress)
			if a != b {
				return a > b
			}
		case "due":
			// Sample dates all belong to November and have zero-padded days.
			if filtered[i].Due != filtered[j].Due {
				return filtered[i].Due < filtered[j].Due
			}
		}
		return filtered[i].Name < filtered[j].Name
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
