# TendUI

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/brand/logo-dark.svg">
  <img src="docs/brand/logo-light.svg" alt="Tend UI" width="260">
</picture>

**Beautiful interfaces. No build pipeline.**

Pure CSS. Semantic HTML. No Node.js, no npm, no bundled JavaScript.
First-class Go and templ support.

Use the stylesheet with any HTML stack. In Go, `//go:embed` bundles the CSS
inside your application: build and ship one executable, with no runtime asset
downloads. Typed templ components and optional HTMX attributes are ready when
you need them.

TendUI pairs original semantic CSS with composable Go components. Build
interfaces with native HTML, typed configuration, and application-owned behavior.

See [migration notes](docs/tendui-migration.md) for the renamed Go module, CSS
namespace, and use with other HTML stacks.

## Install and run

Requires **Go 1.26+**, matching the pinned `github.com/a-h/templ v0.3.1100`
runtime. The asset package imports only the standard library; the components
have one direct Go dependency, templ. It is not a literally dependency-free
module, and this version does not support Go 1.21.

Install the pinned TendUI preview (API changes remain possible before a stable release):

```sh
go get github.com/Tend-Stack/TendUI/v2026@v2026.2.0-preview.2
```

Generated `*_templ.go` files ship with the module, so consumers do not need
the templ CLI to use these components from Go.

A complete `main.go`:

```go
package main

import (
    "log"
    "net/http"

    "github.com/Tend-Stack/TendUI/v2026"
    "github.com/Tend-Stack/TendUI/v2026/components"
    "github.com/a-h/templ"
)

func main() {
    mux := http.NewServeMux()
    mux.Handle(tendui.AssetPath, tendui.AssetHandler())
    page := components.Layout(components.LayoutProps{
        Title: "My Go app",
        Language: "en",
        Theme: components.ThemeLight,
    }, components.Card(components.CardProps{
        Body: components.Button(components.ButtonProps{
            Text: "Hello from Go", Primary: true,
        }),
    }))
    mux.Handle("/", templ.Handler(page))
    log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
}
```

For a complete, tested Go/templ application, see [the quickstart](docs/quickstart.md)
and [`examples/quickstart`](examples/quickstart). To use CSS with any stack, see
[the plain HTML example](examples/html); Go and templ are optional for that path.

Run the local demos, including a component playground, form lab, dashboard,
appearance editor, workspace, native interactions, inline editing, motion and surfaces, and theme endpoint:

```sh
go run ./examples/basic -open
```

The `-open` flag opens `http://127.0.0.1:8080/playground` in your default Linux
browser. See [the demo guide](examples/site/README.md) for all pages and test
scenarios. The example deliberately blocks scripts with its
Content Security Policy. Theme switching and the form work with JavaScript
disabled.

## Packages and assets

```text
tendui.go                 embedded CSS handler
static/tendui.css         readable semantic stylesheet
components/
  layout.templ            document, title, language, theme, stylesheet
  theme.templ             native POST theme form with optional HTMX hooks
  buttons.templ           buttons and button-styled links
  cards.templ             article with optional header, body, footer
  forms.templ             forms, inputs, selects, checkboxes, grids
  fields.templ            composable labels, help/error bindings and input addons
  controls.templ          textareas, radio groups, switches, file inputs
  feedback.templ          alerts, badges, linked error summaries
  shell.templ             application header, sidebar, main content
  overlays.templ          native dialog and popover panels
  navigation.templ        breadcrumbs and bounded pagination
  tables.templ            data tables with column/row contracts and state slots
  wayfinding.templ        page tabs, progress steps and native disclosures
  motion.templ            typed one-shot entrance effects
  states.templ            loading placeholders, empty states and metrics
  appearance.go           typed settings, attributes, and CSS exports
  *_templ.go              generated Go, included for consumers
render/                   buffered page/fragment HTTP responses
examples/site/            thirteen demo pages, searchable catalog, preference cookies and downloads
```

Mount `AssetHandler()` directly at `/tendui/static/`. Do **not** apply
`http.StripPrefix`: the handler already recognizes the complete path. It serves
only `tendui.css`, supports GET/HEAD, conditional ETag requests and byte ranges,
and returns 404 for directories and unknown assets. Revalidation prevents a
browser from keeping stale CSS across binary upgrades. See the
[asset delivery guide](docs/assets.md) for cache/proxy behavior, range requests
and verification scope.

## Components

All text is escaped by templ. Link and form URLs use templ's safe URL handling.
Keep URL destinations and HTMX endpoints application-controlled. Labels wrap
controls, so an ID is optional and repeated components do not invent duplicate
IDs. Supply meaningful labels for every input.

```templ
@components.Form(components.FormProps{
    Action: "/profile", Method: components.MethodPost,
    HxPost: "/profile", HxTarget: "#result", HxSwap: "innerHTML",
    HxIndicator: "#saving",
}) {
    @components.Grid(components.GridProps{}) {
        @components.Input(components.InputProps{
            Name: "email", Label: "Email", Type: components.InputEmail,
            Required: true, Autocomplete: "email",
        })
        @components.Select(components.SelectProps{
            Name: "locale", Label: "Language", Value: "en",
            Options: []components.Option{
                {Label: "English", Value: "en"},
                {Label: "Español", Value: "es"},
            },
        })
    }
    @components.Checkbox(components.CheckboxProps{
        Name: "updates", Label: "Receive updates",
    })
    @components.Button(components.ButtonProps{
        Text: "Save", Type: components.ButtonTypeSubmit, Primary: true,
    })
    <span id="saving" class="htmx-indicator" role="status">Saving…</span>
}
<div id="result" aria-live="polite"></div>
```

`ButtonProps`, `LinkProps`, and `FormProps` expose `HxGet`, `HxPost`, `HxTarget`,
`HxSwap`, and `HxIndicator` directly. Set either HxGet or HxPost for a request,
not both. Empty attributes are omitted. `Button` defaults to `type="button"`
to avoid accidental submission; `Form` defaults to GET. Primary and Outline
can be combined. Raw buttons use the primary color by default.

`CardProps` accepts `Header`, `Body`, and `Footer` components. Cards also accept
templ children after Body. Forms and grids use templ children. Pass ordinary
components rather than `templ.Raw` for untrusted content.
Cards can optionally override the page's shadow with `ElevationFlat`,
`ElevationSubtle`, `ElevationRaised`, or `ElevationFloating`. The zero value
inherits the page setting. See [surface elevation](docs/motion.md#elevation)
for theme tokens, nesting, outlined cards and accessibility behavior.

## Appearance and application components

Appearance choices are independent. The zero value preserves the original
Indigo, comfortable, rounded defaults:

```go
appearance := components.Appearance{
    Preset:     components.PresetForest,
    Density:    components.DensityCompact,
    Corners:    components.CornersSoft,
    Typography: components.TypographySans,
    TextSize:   components.TextStandard,
    Width:      components.WidthWide,
    Elevation:  components.ElevationFlat,
}
page := components.Layout(components.LayoutProps{
    Title: "My workspace",
    Theme: components.ThemeDark,
    Appearance: appearance,
}, body)
```

Available presets: Indigo, Neutral, Forest, Rose. Choose comfortable/compact
spacing, rounded/nearly-square/soft corners, system sans/serif/monospace fonts,
standard/large text, narrow/standard/wide content, and flat/subtle/raised/floating shadows.
No font downloads are needed. Light/dark/system mode remains an independent
`LayoutProps.Theme` setting.

`Appearance.Normalize()` validates preferences and falls back to defaults.
`Attributes()` produces native `data-*` attributes for applications that own
their document shell. `CSS()` exports overrides referencing TendUI's canonical
tokens, so it **requires tendui.css** and must be loaded after it. CSS exports
configure appearance; keep light/dark/system mode on the root HTML element.
`LayoutProps.Head` accepts trusted application components after the base
stylesheet, for custom stylesheet links or metadata.

Open `/settings` in the demo to preview, save, reset, and download appearance
settings. `/appearance.css` and `/appearance.go` export the selected settings.
The library owns no persistence; the example uses bounded, validated preference
cookies. Query previews take precedence over saved preferences without changing
them. Cookie data never controls authorization or application behavior.

New components:

| Component | Contract |
| --- | --- |
| `Badge(BadgeProps)` | Text plus neutral/info/success/warning/danger tone |
| `Alert(AlertProps)` | Title, message, optional children; `Live` opts into announcements |
| `ErrorSummary(ErrorSummaryProps)` | Error messages linked to application-supplied field IDs |
| `Textarea(TextareaProps)` | Label, value, rows, maximum length, required/disabled/read-only/error states |
| `RadioGroup(RadioGroupProps)` | Native fieldset/legend and labeled options sharing one name |
| `Switch(SwitchProps)` | Native checkbox with switch semantics and keyboard interaction |
| `FileInput(FileInputProps)` | Native file chooser with accept/multiple/required states |
| `AppShell(AppShellProps)` | Brand, actions, sidebar links, optional footer, and one main landmark |

`CheckboxProps` now also supports `Invalid` and `DescribedBy`. Use
`FormProps{Method: components.MethodPost, Multipart: true}` for actual uploads.
Applications must validate file types, bound request sizes, and handle storage;
the `accept` attribute is only a browser hint. The playground file picker does
not upload data.

Static alerts do not create live regions. Set `Live: true` for inserted feedback:
danger uses `role="alert"`; other tones use `role="status"`. These attributes
do not guarantee announcements on a full page navigation, so visible error
summaries and linked field errors remain available.

Use `AppShell` as the body inside `Layout`, with page content as templ children.
It owns the main landmark; do not nest another main or application shell inside
it. `MainID` defaults to `main` and connects the skip link. Sidebar navigation
uses a native disclosure that works on desktop and mobile without JavaScript.
See `/workspace` and `examples/site/workspace.templ` for a complete example.

## Native interactions and navigation

`Dialog` and `Popover` use native HTML. Supply a nonempty ID with no whitespace;
invalid IDs return a render error. IDs must be unique within the document, and
their `-title` and `-description` suffixes are reserved for accessible labels.

```templ
@components.Button(components.ButtonProps{
    Text: "Review request",
    CommandFor: "review-request",
    Command: components.DialogShowModal,
})
@components.Dialog(components.DialogProps{
    ID: "review-request", Title: "Review request",
    Description: "Check the details before submitting.",
}) {
    @components.Form(components.FormProps{
        Action: "/requests/approve", Method: components.MethodPost,
    }) {
        @components.Button(components.ButtonProps{
            Text: "Approve", Type: components.ButtonTypeSubmit,
        })
    }
}
```

Keep dialogs outside other forms. The dialog includes a native `method="dialog"`
close form and initially focuses that non-destructive close control. Closing a
dialog does not submit an application request; use a separate normal form for
an action that must reach Go. Native opener/closer fields on Button always force
`type="button"`, even if Type was set to submit. Unsupported command values
produce an inert button. Do not combine native commands and HTMX requests on
the same control.

Modal opening uses the newer `commandfor`/`command="show-modal"` attributes.
Provide an ordinary-page alternative for unsupported browsers, as `/interactions`
does. No JavaScript fallback is injected. See the browser's
[native dialog commands](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dialog#modal_dialogs_using_invoker_commands).

For popovers, use `ButtonProps{PopoverTarget: "help"}` and
`Popover(PopoverProps{ID: "help", Title: "Help"})` with templ children. The default
automatic mode supports dismissal outside the panel; `Mode: PopoverManual`
requires an explicit close/toggle action. Both variants include a close button.
`PopoverAction` supports Toggle, Show, and Hide. Panels are non-modal, centered
by default, and do not claim ARIA menu semantics. Unsupported browsers may show
the content inline. See the [native popover attribute](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/popover).

`Breadcrumbs(BreadcrumbsProps{Items: ...})` renders ancestor links and a final,
non-linked current page. Separators are hidden from assistive technology.

`Pagination(PaginationProps{Page: page, TotalPages: total, URLForPage: pageURL})`
renders normal links, current/disabled states, and a bounded page window with
ellipses. Zero/one-page results omit the navigation. Current page numbers are
clamped to range; a missing URL callback for multiple pages returns a render
error. The application supplies totals, page data, and URLs that retain its
filters. No database, router, or query convention is imposed by the component.
Labels can be localized through the props.

The dashboard example demonstrates filter → sort → paginate on a request-owned
copy of sample records. Page URLs preserve search, status, sorting, page size,
and appearance query settings. Invalid or excessive page numbers are normalized;
empty results keep their recovery link.

## Themes and the zero-JavaScript contract

Set `LayoutProps.Theme` to `ThemeLight` or `ThemeDark` to emit the root
`data-theme` attribute. Leave it at `ThemeSystem` (the empty default) to follow
`prefers-color-scheme`. Override `--tendui-*` properties in your own CSS loaded
after TendUI. Dark tokens have no effect on explicitly light pages.

`ThemeToggle()` renders a small native form that POSTs to `/api/theme/toggle`.
Place it **outside other forms**. You can pass `ThemeToggleProps` for a localized
label and your application's CSRF token. The consumer owns cookies, request
protection, routing, and page rendering; the library has no global theme state.

The included endpoint alternates a validated `light`/`dark` cookie and returns
303 to the current demo when its same-origin Referer is available, otherwise to
the homepage. Cookies are HttpOnly, SameSite=Strict, and Secure over TLS;
the example uses Go's cross-origin request protection. Behind a TLS-terminating
proxy, configure your application's Secure-cookie policy explicitly. The first
example visit uses light so the server can deterministically toggle it; the
library's system theme remains available for applications that prefer it.

The button preserves the requested `hx-post="/api/theme/toggle"`,
`hx-target="html"`, and `hx-swap="outerHTML"` attributes. These are inert without
HTMX. **HTMX itself is JavaScript** and is never included or loaded by TendUI.
Consumers that already use it can use the hooks without writing custom JS.

A document root swap is not a reliable way to replace a page's root attributes
and head. The example detects `HX-Request: true` and sends **200 plus
`HX-Redirect` to the cleaned demo URL** after setting the cookie. This requests a
full navigation. Both native and HTMX paths remove the `theme` query parameter
so a previous preview cannot override the new cookie, while preserving other
settings. For HTMX, a validated same-origin `HX-Current-URL` takes precedence over
Referer; invalid values fall back to a validated Referer, then the homepage.
Only known demo pages are accepted. These headers control presentation, not access.
HTMX does not process its response headers on 3xx responses. See the upstream
[redirect header documentation](https://htmx.org/headers/hx-redirect/) and
[swap reference](https://htmx.org/attributes/hx-swap/).
The isolated browser consumer verifies this navigation with HTMX 2.0.11 in
Chromium, Firefox and WebKit. The runtime is a test dependency and is not bundled.

## Full pages and optional HTMX fragments

`render` is a standard-library-only package. Its small `Component` interface is
implemented by templ components; no adapter or additional dependency is needed.
Use the same route for ordinary navigation and partial updates:

```go
import (
    "log"
    "net/http"

    "github.com/Tend-Stack/TendUI/v2026/components"
    "github.com/Tend-Stack/TendUI/v2026/render"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {
    panel := components.Alert(components.AlertProps{Title: "Profile", Message: "Your profile is ready."})
    page := components.Layout(components.LayoutProps{Title: "Profile"}, panel)
    w.Header().Set("Cache-Control", "private, no-store")
    if err := render.HTML(w, r, http.StatusOK, page, panel); err != nil {
        log.Printf("render profile: %v", err)
    }
}
```

`HTML` chooses the fragment for `HX-Request: true`, but serves the full page for
history restoration or boosted navigation. A nil fragment always serves the page.
It preserves existing `Vary` values and adds the three headers that affect the
representation. Rendering is buffered so a component error produces a generic
500 instead of partially written HTML; errors are returned for logging. HEAD
has no body, including on failures. After calling it, do not write another error
response: a render failure has already sent 500, and a transport failure may
have committed the response. Only the selected component is rendered.

The helper leaves authentication, CSRF protection, cache policy, CSP and status
selection to the application. HTMX headers are presentation hints, never proof
of identity or permission. For cacheable pages, use representation-specific ETags.
Buffering is intended for bounded UI pages/fragments, not streaming responses.

The [inline editing example](examples/site/editing.go) pairs native form actions
and link destinations with `HxPost`/`HxGet`, `HxTarget` and `HxSwap`. It returns
303 after a successful native POST and 200 with the updated panel for an HTMX
request. Invalid native submissions return 422. Invalid HTMX submissions return
200 **with field errors and no save**, because default HTMX does not swap 422
responses. This policy is explicit in the example, not imposed by `render.HTML`.
See the upstream [response handling](https://htmx.org/docs/#response-handling),
[history](https://htmx.org/docs/#history) and
[caching](https://htmx.org/docs/#caching) documentation.

The included demo remains JavaScript-free and blocks scripts with CSP. Its native
flow and HTTP fragment contract are testable separately. Applications enabling
HTMX must supply their own runtime and appropriate CSP; this project neither
loads nor bundles it. The example cookie is unsigned cosmetic sample data,
validated and escaped on every read. It is not an account/session implementation.

## Motion, surfaces and useful states

Cards now accept typed `Surface`, `Entry`, `Delay` and `HoverLift` options.
Choose tinted, outlined or gradient borders; use fade, rise or scale entrances.
`Reveal` wraps other content with the same entrance system. `Skeleton`,
`EmptyState` and `Stat` cover loading, recovery and metric presentation.

`Appearance.Motion` accepts `MotionSystem` (default) or `MotionOff`. System mode
respects OS reduced motion. Off also disables hover displacement; CSS/Go exports
include the choice. Entrances run once, and loading decoration stops after three
cycles. All styling stays in the embedded CSS, with no runtime dependency.

See [the motion guide](docs/motion.md), try `/motion`, and track the broader
[product roadmap](ROADMAP.md). See the [recorded browser results](docs/release-verification.md).

## Data tables and navigation patterns

`DataTable` provides semantic headers/row scopes, application-owned sort URLs,
cell component slots, striped/compact styles, and loading/empty states. The
Dashboard uses it for bidirectional sorting and native checkbox selection review.
`Tabs` supplies ordinary page navigation; `Steps` labels application-owned
progress; `Disclosure` uses native details/summary with optional grouping.

See the [data and navigation guide](docs/data-navigation.md) and `/patterns`
for composition rules, working examples and validation boundaries.

## Composable fields and multi-step forms

`Field` pairs a raw control with its label, description and error through typed
`FieldProps.ControlAttributes()`. `InputGroup` adds text prefixes/suffixes while
preserving native control behavior. `ButtonProps.FormNoValidate` supports Back
actions that submit edits without native constraints blocking navigation.

The `/wizard` project-brief example preserves draft values across Back/Continue,
checks all submitted fields on completion, and shows a validated preview without
creating or storing a project. See [the form guide](docs/forms.md) for composition,
validation and workflow ownership rules.

The `/uploads` example connects the native file chooser to a bounded Go CSV
preview: download a sample, upload up to 64 KiB/50 rows, and inspect validated
results or line-specific errors. Uploaded content stays in request memory;
nothing is saved or imported. See the form guide for limits and HTTP behavior.

## Accessibility and print

Palette checks cover all built-in presets across page, card and subtle surfaces.
Print uses a neutral palette independent of the screen theme, with readable
outlined actions and native checkbox rendering for switches. See the
[accessibility and print guide](docs/accessibility.md) for the measured scope,
customization rules and browser checks still required.

## Layout composition

`Stack` provides vertical rhythm; `Cluster` wraps actions and other related items
into rows. Typed gap and alignment props use a density-aware spacing scale,
preserve document order, and compose with Grid, Card and native forms. Existing
CSS utilities retain their behavior. Try `/patterns?section=composition` or read
[the layout guide](docs/layout.md) for spacing tokens, semantics and examples.

## Reading content

`Prose` styles application-owned HTML for help pages, guides and release notes.
Choose a typed reading measure, use native headings/lists/code, and opt embedded
widgets out of prose typography. Try `/patterns?section=reading` and see the
[reading-content guide](docs/prose.md) for tokens, semantics and examples.

## Inline icons

The Icon component provides 24 original SVG drawings, compiled into Go and
rendered inline. They inherit text color and use no external image or font files.
Button and Link accept optional `IconBefore`, `IconAfter` and `AriaLabel` fields.
Icon-only actions require an explicit accessible name; decorative icons stay
hidden from assistive technology. Explore `/catalog/icon` and the
[icon guide](docs/icons.md) for sizes, composition and verification limits.

## Searchable component catalog

Open `/catalog` to browse all 37 public templ components by name, category or
usage notes. Each entry provides a preview, guidance and a downloadable `.templ`
recipe drawn from the same source compiled for that preview. Layouts open as
standalone pages to preserve document landmarks. Search and pagination use Go
requests; no client runtime or external assets are added.

Catalog-only CSS and recipe source are embedded in the example executable, not
the library stylesheet. See [the catalog guide](docs/catalog.md) for recipe
adaptation, coverage checks and validation boundaries.

## Maintaining the library

Only contributors changing `.templ` files need the generator:

```sh
go install github.com/a-h/templ/cmd/templ@v0.3.1100
templ fmt .
templ generate
go test -race ./...
go vet ./...
```

Commit generated Go alongside its source. A normal Go build consumes CSS
directly through `go:embed`; CSS editing needs no asset build step. Test output
belongs under `tests/results/` and is ignored except for its placeholder.

The Go-only developer tool also checks generated-source consistency, compiles
all catalog downloads independently, records benchmarks and builds the versioned
review executable. See [the development guide](docs/development.md) and
[measured performance](docs/performance.md):

```sh
go run ./cmd/tendui-dev check
go run ./cmd/tendui-dev api
go run ./cmd/tendui-dev generated
go run ./cmd/tendui-dev recipes
go run ./cmd/tendui-dev bench
go run ./cmd/tendui-dev build
./tests/results/demos/tendui-demo -version
```

`check` already includes `api`; the separate command is useful during API review.
The checked-in declaration baseline catches changed props, constants, signatures
and method sets. See [API review and migration](docs/api-compatibility.md) for
intentional updates and the limits of this check.

The [CI matrix](docs/ci.md) checks native Linux, macOS and Windows with Go 1.26.0
and 1.27.3, plus a vulnerability scan. Browser checks cover Chromium, Firefox and
WebKit. These jobs must pass on the preview source; see
[preview release notes](docs/releases/2026.2.0-preview.2.md). The older
[release verification](docs/release-verification.md) documents the pre-rename
release only and is not evidence for the current preview.

Versions follow `YEAR.RELEASE.FIX`, with SemVer prerelease suffixes for previews:
`v2026.2.0-preview.2`. Install the explicit tag; a preview is not a stable API promise.
Development builds use `git describe --tags --long --dirty`.

Go treats 2026 as the semantic major version, so the module path ends in
`/v2026`. A new annual line such as 2027 will need `/v2027` and consumer import
updates. Fixes and previews on the 2026 line retain `/v2026`.

## Project links

- Source and issues: https://github.com/Tend-Stack/TendUI
- Website source: https://github.com/Tend-Stack/TendUI-web

Fetch and integrate GitHub `main` before contributing. See `HANDOFF.md` for the
current development entry points and remaining release checks.

## License

MIT. See [LICENSE](LICENSE).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for setup, review and MIT + DCO sign-off.
Report vulnerabilities privately using [SECURITY.md](SECURITY.md); community
participation follows [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
