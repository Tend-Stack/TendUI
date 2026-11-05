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

Run the complete local showcase, including a working form and theme endpoint:

```sh
go run ./examples/basic
```

Open `http://127.0.0.1:8080`. The example deliberately blocks scripts with its
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
  *_templ.go              generated Go, included for consumers
examples/basic/           runnable showcase and cookie handler
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
303 to the homepage. Cookies are HttpOnly, SameSite=Strict, and Secure over TLS;
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

- Collaboration source: https://j.santanafam.com/Rubirosa/GothUI
- Go import repository: https://github.com/Tend-Stack/GothUI

Before publishing this initial implementation, fetch and inspect the Gitea
history, integrate any existing work without overwriting it, and publish the
same verified history to GitHub. See `HANDOFF.md` for current verification and
delivery status.

## License

MIT. The CSS is a new implementation inspired by Pico's principles, not a
copy of its generated stylesheet. Pico's MIT notice is preserved separately
in `THIRD_PARTY_NOTICES.md` as attribution for the supplied reference.
