# Integrating TendUI

Install the pinned preview:

```sh
go get github.com/Tend-Stack/TendUI/v2026@v2026.2.0-preview.2
```

Use `github.com/Tend-Stack/TendUI/v2026` for embedded assets and
`github.com/Tend-Stack/TendUI/v2026/components` for typed templ components.
The root Go package is `tendui`. Authoring your own templ files requires the
generator version listed in `go.mod`; consuming the shipped Go components does
not require generation.

The asset endpoint is `/tendui/static/tendui.css`. Prefer `tendui.AssetPath`
and `tendui.StylesheetURL` in Go code. CSS custom properties and prefixed utility
classes use `tendui`; keep application overrides and appearance exports aligned
with the stylesheet version you install.

The demo site's preference cookies use the `tendui_` prefix. Consumer applications
own their preference storage and routing. The standalone examples retain their
default appearance; tendui.com defaults to dark and shares its theme toggle across
Home, Docs and the demos.

## CSS with any HTML stack

Copy `static/tendui.css` into your application's static assets and link it:

```html
<link rel="stylesheet" href="/assets/tendui.css">
```

Use semantic HTML such as `button`, `input`, `article` and `details`. There is
no CSS compilation step or required browser JavaScript. In Go, the asset handler
serves that same stylesheet from the compiled binary through `//go:embed`.
The templ runtime is required only for the Go component wrappers; HTMX is an
optional enhancement supplied by the application. Authoring new `.templ` files
requires templ generation, and building a Go app still requires the Go compiler.
