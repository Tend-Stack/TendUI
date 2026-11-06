# Data and navigation patterns

TendUI's table and navigation components preserve native HTML semantics. The Go
application chooses data, sort order, destinations and progress. No browser-side
model, reflection layer, JavaScript runtime or external asset is required.

## DataTable

```go
props := components.DataTableProps{
    Caption: "Project portfolio",
    Striped: true,
    Columns: []components.TableColumn{
        {Label: "Project", RowHeader: true,
            Sort: components.SortAscending,
            SortURL: "/projects?sort=name&direction=desc",
            SortLabel: "Sort Project descending"},
        {Label: "Owner"},
        {Label: "Hours", Alignment: components.AlignEnd},
    },
    Rows: []components.TableRow{
        {Cells: []components.TableCell{
            {Text: "Studio website"},
            {Text: "Alex"},
            {Text: "24"},
        }},
    },
}
// Render components.DataTable(props) inside your templ page.
```

Caption and every column label are required. Each non-loading row must have
exactly as many cells as columns. Invalid shapes fail rendering before table
markup is emitted, rather than silently misaligning data. At most one column
can declare a sort direction and one can act as the row header. Supported
sort directions are none, ascending and descending.

Text cells are escaped. A cell's optional `Content templ.Component` takes
precedence over Text and can hold a badge, progress indicator, link or checkbox.
Use application-authored components; do not wrap untrusted input in `templ.Raw`.
Alignment accepts start, center or end, with invalid values normalized to start.
Only supply the current page of rows; paginate large datasets in your handler.

`SortURL` creates a real navigation link with a visible direction marker. The
handler must validate the sort key, sort records and set `Sort` to describe the
returned order. Set `SortLabel` to describe the next action when toggling order.
The component does not reorder data. Only the active column receives `aria-sort`.
The scroll container is labeled and keyboard reachable; native table headers and
row scopes remain intact at narrow widths. This is a table, not an ARIA grid.

`Striped` and `Compact` control per-table styling. `Loading` suppresses stale
rows, marks the table busy and places a status label outside that busy region.
`LoadingLabel` can be localized. An empty table retains its caption and columns;
pass `Empty` to supply recovery content inside a full-width cell.

Selection remains a native form concern. Put a labeled checkbox in each row,
wrap the table with a Form, and validate every submitted ID on the server.
`TableRow.Selected` is only a visual highlight; it does not set a checkbox or
provide selection behavior. There is no script-powered select-all control or
implicit cross-page selection. Do not place forms inside a table already wrapped
by another form.

The dashboard demonstrates both sort directions for name, progress and due date,
status navigation, filtering and pagination. Header/status URLs preserve the
existing query and reset to page one. Its review action accepts only known sample
names, deduplicates them, limits request size/count, rejects cross-origin unsafe
requests, and renders a non-mutating receipt. It is not a production bulk-action
service: real actions require authorization for every selected record and proper
application storage. No selection survives navigation in this demo.

## Tabs, Steps and Disclosure

```templ
@components.Tabs(components.TabsProps{
    Label: "Project sections",
    Items: []components.TabItem{
        {Label: "Overview", Href: "/projects", Current: true},
        {Label: "Activity", Href: "/projects/activity"},
        {Label: "Reports", Disabled: true},
    },
})
@components.Steps(components.StepsProps{
    Label: "Setup progress",
    Items: []components.StepItem{
        {Label: "Details", Href: "/setup/details", Completed: true},
        {Label: "Review", Current: true},
        {Label: "Finish"},
    },
})
@components.Disclosure(components.DisclosureProps{
    Summary: "What happens next?", Name: "setup-help",
}) {
    <p>Your Go handler decides how to validate and save the form.</p>
}
```

Tabs is a navigation landmark with ordinary links and `aria-current="page"`.
It intentionally does not use tablist/tab roles or promise arrow-key activation.
Use native Tab/Enter navigation. Disabled items have no link. A current item may
omit its URL; other available items require one. Empty labels, multiple current
items, or a current-and-disabled item fail rendering.

Steps uses an ordered list, numbered markers, `aria-current="step"` and completed
status text for assistive technology. `CompletedLabel` permits localization.
Only provide a URL for a step the application permits visiting. Multiple current
steps or a current-and-completed step fail rendering. The component does not
persist form values or authorize skipping validation; that belongs to the handler.
The `/patterns` walkthrough is a stateless reading guide, not a setup wizard.

Disclosure is native details/summary with children in its body. `Open` is an
initial boolean state. A shared `Name` groups disclosures in supporting browsers;
use unique names for unrelated groups. Older browsers may keep multiple panels
open. The summary is text, keeping controls and links in the disclosure body.
Do not manually set aria-expanded: the browser owns the open state.

## Validation boundaries

Go tests cover row/header contracts, escaping, safe URLs, current-state conflicts,
loading and empty markup, sorting/filter preservation, selection validation and
native request paths. Full race tests and vet are required for delivery. Browser
layout, keyboard/screen-reader behavior and native grouped-disclosure compatibility
still require manual verification outside the current restricted sandbox.

Reference patterns: [WAI table guidance](https://www.w3.org/WAI/ARIA/apg/patterns/table/),
[WAI step indicators](https://www.w3.org/WAI/tutorials/forms/multi-page/), and
[MDN native disclosures](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/details).
