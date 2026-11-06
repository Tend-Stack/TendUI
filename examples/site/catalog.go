package site

import (
	_ "embed"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/Tend-Stack/GothUI/v2026/internal/recipes"
	"github.com/Tend-Stack/GothUI/v2026/render"
	"github.com/a-h/templ"
)

//go:embed catalog_samples.templ
var catalogSamples string

type catalogEntry struct {
	Name, Category, Summary, Notes, Demo string
	Example                              func() templ.Component
}

func (e catalogEntry) ID() string       { return strings.ToLower(e.Name) }
func (e catalogEntry) URL() string      { return "/catalog/" + e.ID() }
func (e catalogEntry) Standalone() bool { return e.Name == "Layout" || e.Name == "AppShell" }

var catalogEntries = []catalogEntry{
	{"Layout", "Layout", "A complete HTML document with embedded CSS and theme attributes.", "Use once per document. Pass the body component and optional trusted Head metadata. Appearance controls tokens; the application owns its preferences.", "/settings", ExampleLayout},
	{"AppShell", "Layout", "An application header, sidebar and main landmark.", "Place inside Layout and do not nest shells or main landmarks. Navigation uses a native disclosure; Actions and Footer accept components.", "/workspace", ExampleAppShell},
	{"Grid", "Layout", "Responsive columns that preserve document order.", "Place block-level content in its children. Columns stack as available space narrows; no JavaScript measures the viewport.", "/playground", ExampleGrid},
	{"Stack", "Layout", "Vertical rhythm with a density-aware gap.", "Owns direct children's block margins. A plain div, not a landmark or list; keep semantic sections and labels in the children. Align defaults to stretch.", "/patterns?section=composition", ExampleStack},
	{"Cluster", "Layout", "Wrapping action rows with typed spacing and alignment.", "Preserves document and focus order. Each line wraps naturally; JustifyBetween distributes spare space within each line. Use native links and submit buttons for actions.", "/patterns?section=composition", ExampleCluster},
	{"Prose", "Layout", "Readable long-form HTML with a bounded line length.", "Use native headings, paragraphs and lists. It does not parse Markdown or sanitize raw HTML. Wrap embedded widgets in data-prose=off to retain their own typography.", "/patterns?section=reading", ExampleProse},
	{"Icon", "Visuals", "A cohesive set of 24 inline SVG drawings with typed names.", "Icons are decorative by default. Label standalone images; name icon-only actions on the button or link. Unknown names return a render error. No fonts, sprites, remote requests or JavaScript are used.", "/catalog/icon/preview", ExampleIcon},
	{"Card", "Layout", "A semantic article with surface and content slots.", "Header, Body and Footer accept components. Surface, Entry, Delay and HoverLift are optional; a card is not automatically interactive.", "/motion", ExampleCard},
	{"Button", "Actions", "Typed submit, reset and native command buttons.", "Native HTMX fields are optional. Set Type explicitly for forms; FormNoValidate supports Back actions. Native overlay targets force type=button.", "/wizard", ExampleButton},
	{"Link", "Actions", "A real link with button styling.", "Use Href for navigation. Outline and Primary affect presentation. Optional HTMX attributes do not load a runtime.", "/dashboard", ExampleLink},
	{"ThemeToggle", "Actions", "A native POST form for changing the color mode.", "Place outside other forms. Supply your application's endpoint, cookies and CSRF policy; this demo already has the handler. The optional HTMX response redirects to a cleaned same-origin URL.", "/settings", ExampleThemeToggle},
	{"Form", "Forms", "Native GET/POST submission with optional HTMX hooks.", "Action and Method remain functional without scripts. Multipart enables encoding only; the application must bound and validate file uploads.", "/forms", ExampleForm},
	{"Input", "Forms", "A labeled input with typed input kinds and native states.", "Use Required, Disabled, ReadOnly, Invalid and DescribedBy as needed. It already includes a label; do not put it inside Field.", "/forms", ExampleInput},
	{"Select", "Forms", "A labeled native select with typed options.", "Value chooses the selected option. The browser owns the picker and keyboard behavior. Validate submitted values in Go.", "/playground", ExampleSelect},
	{"Checkbox", "Forms", "A native checkbox with label and validation metadata.", "An unchecked checkbox is normally absent from submitted form data. Checked is server-rendered state, not shared library state.", "/forms", ExampleCheckbox},
	{"Textarea", "Forms", "A labeled multiline control for longer input.", "Rows and MaxLength are browser hints. Validate content and length on the server, preserving invalid input for correction.", "/forms", ExampleTextarea},
	{"RadioGroup", "Forms", "A fieldset and legend containing same-name radio options.", "The browser supplies radio keyboard behavior. Options, Value, Required and Disabled remain native form state.", "/playground", ExampleRadioGroup},
	{"Switch", "Forms", "A checkbox presented as a labeled switch.", "It submits like a checkbox. Use for a boolean choice; this preview is not wired to a saved preference.", "/playground", ExampleSwitch},
	{"FileInput", "Forms", "A labeled native file chooser.", "Accept filters the picker as a hint. It is not MIME validation or an upload service. The component preview performs no upload; the CSV demo shows bounded Go validation without saving files.", "/uploads", ExampleFileInput},
	{"Field", "Forms", "A composable label with linked help and error text.", "Spread ControlAttributes on one raw control. Use a unique ID. Required and aria-invalid are connected to actual rendered descriptions.", "/wizard", ExampleField},
	{"InputGroup", "Forms", "Text prefixes and suffixes around a raw control.", "Decorations are hidden from assistive technology. Put meaningful currency, units or URL context in the Field label or description.", "/wizard", ExampleInputGroup},
	{"Badge", "Feedback", "A compact text status in one of five tones.", "Use meaningful text so color is not the only status cue. Badge itself is not a live announcement or action.", "/playground", ExampleBadge},
	{"Alert", "Feedback", "A message with an optional heading and content slot.", "Live is opt-in for dynamically inserted feedback. Static notices should remain static to avoid unnecessary announcements.", "/forms", ExampleAlert},
	{"ErrorSummary", "Feedback", "A linked summary of fields that need correction.", "Point each error ID at an actual control. Pair it with Field or other per-control error descriptions and preserve submitted values.", "/wizard", ExampleErrorSummary},
	{"Skeleton", "Feedback", "A labeled loading placeholder with bounded decorative motion.", "Only the decorative shapes pulse, for three cycles. Replace it when data arrives; show a recovery state on failure. This preview is not a running request.", "/motion", ExampleSkeleton},
	{"EmptyState", "Feedback", "A clear explanation and recovery action for missing content.", "The children slot supplies your action. Static empty states do not automatically announce an alert.", "/dashboard?q=no-matching-project", ExampleEmptyState},
	{"Stat", "Feedback", "An associated metric label, value and detail.", "Uses a description list. Values are strings so the application controls formatting and localization.", "/motion", ExampleStat},
	{"DataTable", "Data", "Semantic rows, headers, sorting links and state slots.", "Caption and column labels are required. Row shapes are validated. Go owns sorting, pagination and selection; cells accept escaped text or components.", "/dashboard", ExampleDataTable},
	{"Breadcrumbs", "Navigation", "A linked trail to the current page.", "The final item is always the current page. Supply meaningful labels and application-controlled destinations.", "/patterns", ExampleBreadcrumbs},
	{"Pagination", "Navigation", "Bounded page links that keep large collections navigable.", "URLForPage preserves your filters. Boundary controls are non-interactive. The page window remains bounded even for very large totals.", "/dashboard", ExamplePagination},
	{"Tabs", "Navigation", "Page navigation styled as tabs.", "These are ordinary links, not an ARIA tablist. Use Tab/Enter, with Current for the active page and Disabled for unavailable destinations.", "/patterns", ExampleTabs},
	{"Steps", "Navigation", "A numbered description of application-owned progress.", "Current uses aria-current=step. Only give a URL to a step the application permits visiting. The component does not validate workflow permissions.", "/wizard", ExampleSteps},
	{"Disclosure", "Navigation", "A native details/summary disclosure.", "Name groups panels in supporting browsers; older browsers may allow several to stay open. Keep controls in the body, outside Summary.", "/patterns?section=disclosures", ExampleDisclosure},
	{"Dialog", "Overlays", "A native modal with a labeled title and close form.", "Provide a unique ID and a native opener. Place outside other forms because it contains a method=dialog form. Keep an ordinary-page fallback for unsupported commands.", "/interactions", ExampleDialog},
	{"Popover", "Overlays", "A non-modal native panel with automatic or manual dismissal.", "Provide a unique ID and an opener. It is not an ARIA menu. Older browsers may show its content inline.", "/interactions", ExamplePopover},
	{"Reveal", "Motion", "A short, one-shot entrance around content.", "Use fade, rise or scale. It runs on insertion, not scrolling. Respects reduced motion and MotionOff; the div wrapper must fit its HTML parent.", "/motion", ExampleReveal},
}

var catalogCategories = []string{"Layout", "Visuals", "Actions", "Forms", "Feedback", "Data", "Navigation", "Overlays", "Motion"}

func catalogSource(e catalogEntry) (string, error) {
	return recipes.Source(catalogSamples, e.Name)
}

func catalogRoute(path string) (catalogEntry, string, bool) {
	if !strings.HasPrefix(path, "/catalog/") {
		return catalogEntry{}, "", false
	}
	id, action, _ := strings.Cut(strings.TrimPrefix(path, "/catalog/"), "/")
	if action != "" && action != "preview" && action != "example.templ" {
		return catalogEntry{}, "", false
	}
	for _, e := range catalogEntries {
		if id == e.ID() {
			return e, action, true
		}
	}
	return catalogEntry{}, "", false
}

func catalogPagePath(path string) bool {
	_, action, ok := catalogRoute(path)
	return ok && action != "example.templ"
}

type catalogOptions struct {
	Query, Category    string
	Entries            []catalogEntry
	Page, Pages, Total int
	query              url.Values
}

func readCatalog(q url.Values) catalogOptions {
	query := []rune(strings.TrimSpace(q.Get("q")))
	if len(query) > 120 {
		query = query[:120]
	}
	d := catalogOptions{Query: string(query), Category: "all", query: cloneQuery(q)}
	for _, category := range catalogCategories {
		if q.Get("category") == category {
			d.Category = category
		}
	}
	for _, e := range catalogEntries {
		haystack := strings.ToLower(e.Name + " " + e.Category + " " + e.Summary + " " + e.Notes)
		if (d.Category == "all" || e.Category == d.Category) && strings.Contains(haystack, strings.ToLower(d.Query)) {
			d.Entries = append(d.Entries, e)
		}
	}
	sort.Slice(d.Entries, func(i, j int) bool { return d.Entries[i].Name < d.Entries[j].Name })
	d.Total = len(d.Entries)
	d.Pages = max(1, (d.Total+11)/12)
	page, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		page = 1
	}
	d.Page = max(1, min(page, d.Pages))
	start := (d.Page - 1) * 12
	d.Entries = d.Entries[start:min(start+12, d.Total)]
	return d
}
func (d catalogOptions) PageURL(page int) string {
	q := cloneQuery(d.query)
	q.Set("q", d.Query)
	q.Set("category", d.Category)
	q.Set("page", strconv.Itoa(page))
	return "/catalog?" + q.Encode()
}

func serveCatalog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/catalog" {
		renderPage(w, r, "Component catalog · GothUI", catalogPage(readCatalog(r.URL.Query())), http.StatusOK)
		return
	}
	e, action, ok := catalogRoute(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if action == "preview" {
		if e.Name == "Layout" {
			htmlHeaders(w)
			if err := render.HTML(w, r, http.StatusOK, e.Example(), nil); err != nil {
				log.Printf("catalog preview: %v", err)
			}
		} else if e.Name == "AppShell" {
			renderPage(w, r, "AppShell preview · GothUI", e.Example(), http.StatusOK)
		} else {
			renderPage(w, r, e.Name+" preview · GothUI", catalogPreviewPage(e), http.StatusOK)
		}
		return
	}
	source, err := catalogSource(e)
	if err != nil {
		log.Printf("catalog: %v", err)
		http.Error(w, "example unavailable", http.StatusInternalServerError)
		return
	}
	if action == "example.templ" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Content-Disposition", `attachment; filename="`+e.ID()+`_example.templ"`)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Length", strconv.Itoa(len(source)))
		if r.Method != http.MethodHead {
			_, _ = w.Write([]byte(source))
		}
		return
	}
	renderPage(w, r, e.Name+" · GothUI catalog", catalogDetail(e, source), http.StatusOK)
}
