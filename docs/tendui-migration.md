# Moving to TendUI

TendUI is the new name of this library. The renamed module is available on
`main`; historical `v2026.1.0` tags retain their original module path and have
not been rewritten. Pin the resolved development version for reproducible builds.

Update Go imports to `github.com/Tend-Stack/TendUI/v2026` and the root package
identifier to `tendui`. Regenerate your own templ files after changing imports.
The component APIs remain available, with the additive `SiteHeader` component
and host-page options for the example handler.

The asset endpoint is now `/tendui/static/tendui.css`. Prefer `tendui.AssetPath`
and `tendui.StylesheetURL` in Go code. CSS custom properties and prefixed utility
classes now use `tendui` rather than the previous prefix. Update application
overrides, downloaded recipes and appearance exports together with the stylesheet.

The example site's preference cookies now use the `tendui_` prefix. Existing
demo preferences are not migrated; visitors can save their choices again.
The standalone examples retain their default appearance, while tendui.com
defaults to dark and shares its theme toggle across Home, Docs and the demos.

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
