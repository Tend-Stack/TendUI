# GothUI

Semantic CSS and typed templ components, embedded in your Go binary.
No Node.js, npm, Sass, CDN, web fonts, bundled JavaScript, or asset compilation.

GothUI reimplements Pico's semantic styling principles in original, plain CSS.
It is independently named and does not use Pico's logo. The supplied Pico
source and conversation remain outside this project as reference material.

## Install and run

Requires **Go 1.26+**, matching the pinned `github.com/a-h/templ v0.3.1100`
runtime. The asset package imports only the standard library; the components
have one direct Go dependency, templ. It is not a literally dependency-free
module, and this version does not support Go 1.21.

After repository publication, consumers can use:

```sh
go get github.com/Tend-Stack/GothUI/v2026
```

The GitHub repository is currently private. Until that changes, Go downloads
require repository credentials and `GOPRIVATE=github.com/Tend-Stack/GothUI`.
Generated `*_templ.go` files ship with the module, so consumers do not need
the templ CLI to use these components from Go.

A complete `main.go`:

```go
package main

import (
    "log"
    "net/http"

    "github.com/Tend-Stack/GothUI/v2026"
    "github.com/Tend-Stack/GothUI/v2026/components"
    "github.com/a-h/templ"
)

func main() {
    mux := http.NewServeMux()
    mux.Handle(gothui.AssetPath, gothui.AssetHandler())
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

Run the local demos, including a component playground, form lab, dashboard,
appearance editor, workspace, native interactions, and theme endpoint:

```sh
go run ./examples/basic -open
```

The `-open` flag opens `http://127.0.0.1:8080/playground` in your default Linux
browser. See [the demo guide](examples/basic/README.md) for all pages and test
scenarios. The example deliberately blocks scripts with its
Content Security Policy. Theme switching and the form work with JavaScript
disabled.

## Packages and assets

```text
gothui.go                 embedded CSS handler
static/gothui.css         readable semantic stylesheet
components/
  layout.templ            document, title, language, theme, stylesheet
  theme.templ             native POST theme form with optional HTMX hooks
  buttons.templ           buttons and button-styled links
  cards.templ             article with optional header, body, footer
  forms.templ             forms, inputs, selects, checkboxes, grids
  controls.templ          textareas, radio groups, switches, file inputs
  feedback.templ          alerts, badges, linked error summaries
  shell.templ             application header, sidebar, main content
  overlays.templ          native dialog and popover panels
  navigation.templ        breadcrumbs and bounded pagination
  appearance.go           typed settings, attributes, and CSS exports
  *_templ.go              generated Go, included for consumers
examples/basic/           seven demo pages, preference cookies, and downloads
```

Mount `AssetHandler()` directly at `/gothui/static/`. Do **not** apply
`http.StripPrefix`: the handler already recognizes the complete path. It serves
only `gothui.css`, supports GET/HEAD, conditional ETag requests and byte ranges,
and returns 404 for directories and unknown assets. Revalidation prevents a
browser from keeping stale CSS across binary upgrades.

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
standard/large text, narrow/standard/wide content, and raised/flat cards.
No font downloads are needed. Light/dark/system mode remains an independent
`LayoutProps.Theme` setting.

`Appearance.Normalize()` validates preferences and falls back to defaults.
`Attributes()` produces native `data-*` attributes for applications that own
their document shell. `CSS()` exports overrides referencing GothUI's canonical
tokens, so it **requires gothui.css** and must be loaded after it. CSS exports
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
See `/workspace` and `examples/basic/workspace.templ` for a complete example.

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
`prefers-color-scheme`. Override `--gothui-*` properties in your own CSS loaded
after GothUI. Dark tokens have no effect on explicitly light pages.

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
HTMX. **HTMX itself is JavaScript** and is never included or loaded by GothUI.
Consumers that already use it can use the hooks without writing custom JS.

A document root swap is not a reliable way to replace a page's root attributes
and head. The example detects `HX-Request: true` and sends **200 plus
`HX-Refresh: true`** after setting the cookie. That performs a full navigation,
not an instant partial swap; the new document then has the correct theme.
HTMX does not process its response headers on 3xx responses. See the upstream
[response header documentation](https://htmx.org/docs/#response-headers) and
[swap reference](https://htmx.org/attributes/hx-swap/).

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

Versions follow `YEAR.RELEASE.FIX` (initial line: `2026.1.0`). No release tag is
created until publication is verified. To embed an exact development identifier
after tagging:

```sh
go build -ldflags "-X main.version=$(git describe --tags --long --dirty)" ./examples/basic
```

Go treats 2026 as the semantic major version, so the module path ends in
`/v2026`. This is required for a future `v2026.1.0` tag to resolve correctly.
A new annual line such as 2027 will need `/v2027` and consumer import updates;
fixes on the 2026 line keep `/v2026`. See Go's
[major version rules](https://go.dev/doc/modules/major-version).

Configured repositories:

- Go import repository: https://github.com/Tend-Stack/GothUI

history, integrate any existing work without overwriting it, and publish the
same verified history to GitHub. See `HANDOFF.md` for current verification and
delivery status.

## License

MIT. The CSS is a new implementation inspired by Pico's principles, not a
copy of its generated stylesheet. Pico's MIT notice is preserved separately
in `THIRD_PARTY_NOTICES.md` as attribution for the supplied reference.
