# Go html/template quickstart

Use TendUI with Go's standard-library `html/template` and ordinary semantic HTML.
This example imports only the asset package: no templ components, generator,
Node.js, npm or client JavaScript are needed. The HTML template and CSS both ship
inside the executable through `go:embed`.

From a checkout containing this example:

```sh
go run ./examples/htmltemplate
```

Open <http://127.0.0.1:8080>. Submit a sample name and try angle brackets or quotes:
`html/template` escapes it in the heading and input value. Keep user data as
ordinary strings; do not convert it to `template.HTML` or build markup by string
concatenation. The GET form places the name in the URL; use sample data only.
The demo follows system light/dark preference. Set `data-theme="light"` or
`data-theme="dark"` on `<html>` for an explicit theme.

To use it in your own directory, copy `main.go` and `page.html`, then initialize
and pin the library:

```sh
go mod init example.com/my-tendui-app
go get github.com/Tend-Stack/TendUI/v2026@v2026.2.0-preview.2
go run .
```

The recipe is new on main; it uses the asset API already present in that preview.
Use Go 1.26 or newer. `AssetHandler()` is mounted at `tendui.AssetPath` without
`StripPrefix`, and the template links `tendui.StylesheetURL` so the route and
stylesheet stay aligned. The handler renders to a buffer before sending headers
and serves only the root page and embedded assets.

Run its focused HTTP checks with `go test ./examples/htmltemplate`. To ship it,
build with `go build -o tendui-demo ./examples/htmltemplate`; the resulting binary
does not need `page.html` or CSS files beside it. Your application supplies
persistence, authentication and request protection when those are needed.

Choose the [Go/templ quickstart](../../docs/quickstart.md) for typed component
composition, or [plain HTML/CSS](../html) for a page without Go. All three use the
same semantic stylesheet.
