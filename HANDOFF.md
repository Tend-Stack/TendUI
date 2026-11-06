# GothUI development handoff

GothUI is an original semantic CSS and typed templ component library for Go.
The public collaboration source is https://github.com/Tend-Stack/GothUI.
The standalone website lives at https://github.com/Tend-Stack/GothUI-Web.

## Current implementation

- Embedded CSS; no Node/npm, CDN, external fonts, or bundled JavaScript.
- 36 templ components, 24 inline SVG icons, and 13 demo areas.
- Light/dark/system appearance, color presets, density, typography, widths,
  elevation, and bounded CSS motion with reduced-motion support.
- Native form workflows and optional HTMX attributes; applications supply
  routing, storage, request protection, and any HTMX runtime themselves.
- Go 1.26 minimum; templ v0.3.1100 is the only runtime dependency.
- YEAR.RELEASE.FIX versions, with module path `/v2026` for the 2026 line.

## Development entry points

See [README.md](README.md) for setup and component usage,
[docs/development.md](docs/development.md) for generation and validation,
[docs/ci.md](docs/ci.md) for the native verification matrix, and
[ROADMAP.md](ROADMAP.md) for remaining work.

Run `go run ./cmd/gothui-dev check` after code changes. Template edits also need
`generated` and `recipes` checks. Commit generated Go alongside templates.
Use `go run ./cmd/gothui-dev build` to refresh the single review executable;
run `./tests/results/demos/gothui-demo -open` on a host with browser access.

## Release boundary

The initial GitHub native matrix and vulnerability job passed. Require green
checks at the exact proposed release SHA. Still outstanding: real-browser and
assistive-technology review, optional HTMX browser verification, a release tag,
and installation from that published version. Do not infer these from Go tests.

Keep focused local test evidence under ignored `tests/results/`. Reuse one
review binary, remove completed disposable build environments, and preserve
active work. Fetch and integrate GitHub `main` before publishing changes.
