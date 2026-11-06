# Start with TendUI

Choose Go/templ for typed components and embedded assets, or plain HTML for a
stylesheet that works with any stack. Neither path requires Node.js or npm.

## Go and templ

Use Go 1.26 or newer.

```sh
mkdir my-tendui-app
cd my-tendui-app
go mod init example.com/my-tendui-app
go get github.com/Tend-Stack/TendUI/v2026@v2026.2.0-preview.2
go install github.com/a-h/templ/cmd/templ@v0.3.1100
```

Save `page.templ`:

```templ
package main

import "github.com/Tend-Stack/TendUI/v2026/components"

templ page() {
	@components.Layout(components.LayoutProps{
		Title:    "My TendUI app",
		Language: "en",
		Theme:    components.ThemeSystem,
	}, body())
}

templ body() {
	<main>
		<h1>Hello, TendUI.</h1>
		@components.Card(components.CardProps{}) {
			<p>This page and its stylesheet ship in one Go binary.</p>
			<details>
				<summary>What does the browser need?</summary>
				<p>Just HTML and CSS. This disclosure works without JavaScript.</p>
			</details>
		}
	</main>
}
```

Save `main.go`:

```go
// A complete TendUI application. The generated page and CSS ship in its binary.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Tend-Stack/TendUI/v2026"
	"github.com/a-h/templ"
)

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(tendui.AssetPath, tendui.AssetHandler())
	mux.Handle("GET /{$}", templ.Handler(page()))
	return mux
}

func main() {
	server := &http.Server{
		Addr: "127.0.0.1:8080", Handler: handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
```

```sh
templ generate
go run .
```

Open <http://127.0.0.1:8080>. If `templ` is not on your PATH, run the executable
in `$(go env GOPATH)/bin/templ` instead. The generator is only needed when
writing or changing your own `.templ` files. The library ships its generated Go.

Build with `go build -o my-tendui-app .`. The binary contains your rendered-page
code and TendUI's stylesheet through `//go:embed`. Deploy the executable; no
runtime asset directory is needed. Add application-specific routing, storage,
request protection and production server configuration as required.

From this repository, `go run ./examples/quickstart` runs the same checked-in
example without installing a generator. Its HTTP test verifies the page and
embedded stylesheet; the consumer packaging check builds it in a separate module.

## Plain HTML and CSS

Use [examples/html](../examples/html). Copy its `index.html` and
`static/tendui.css` into one folder, then open `index.html`. This works without
Go, templ, JavaScript, or a build tool. It also works on any static web host.

Change `data-theme="light"` to `dark`, or omit it to follow system preferences.
Keep any optional HTMX runtime and application behavior under your own control.
