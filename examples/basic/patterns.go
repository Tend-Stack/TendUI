package main

import (
	"github.com/Tend-Stack/GothUI/v2026/components"
	"net/url"
	"strconv"
)

type patternOptions struct {
	Section string
	Step    int
	Loading bool
}

func readPatterns(q url.Values) patternOptions {
	step, err := strconv.Atoi(q.Get("step"))
	if err != nil || step < 1 || step > 3 {
		step = 1
	}
	return patternOptions{Section: choice(q.Get("section"), "navigation", "navigation", "tables", "disclosures"), Step: step, Loading: q.Get("state") == "loading"}
}
func (p patternOptions) Tabs() []components.TabItem {
	return []components.TabItem{
		{Label: "Navigation", Href: "/patterns?section=navigation", Current: p.Section == "navigation"},
		{Label: "Table states", Href: "/patterns?section=tables", Current: p.Section == "tables"},
		{Label: "Disclosures", Href: "/patterns?section=disclosures", Current: p.Section == "disclosures"},
	}
}
func (p patternOptions) Steps() []components.StepItem {
	var items []components.StepItem
	for i, label := range []string{"Choose a pattern", "Explore its states", "Use it in Go"} {
		href := ""
		if i+1 < p.Step {
			href = "/patterns?section=navigation&step=" + strconv.Itoa(i+1)
		}
		items = append(items, components.StepItem{Label: label, Href: href, Current: i+1 == p.Step, Completed: i+1 < p.Step})
	}
	return items
}
