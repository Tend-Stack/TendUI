# TendUI development handoff

TendUI is an original semantic CSS and typed templ component library for Go.
The public collaboration source is https://github.com/Tend-Stack/TendUI.
The standalone website lives at https://github.com/Tend-Stack/TendUI-web.

## Current implementation

- Embedded CSS; no Node/npm, CDN, external fonts, or bundled JavaScript.
- 37 templ components, 24 inline SVG icons, and 13 demo areas.
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

Run `go run ./cmd/tendui-dev check` after code changes. Template edits also need
`generated` and `recipes` checks. Commit generated Go alongside templates.
Use `go run ./cmd/tendui-dev build` to refresh the single review executable;
run `./tests/results/demos/tendui-demo -open` on a host with browser access.

## Release boundary

The current renamed preview is `v2026.2.0-preview.2`. Install that explicit tag;
see [release notes](docs/releases/2026.2.0-preview.2.md) and
[quickstarts](docs/quickstart.md).
Historical pre-rename tags and evidence remain immutable and are clearly labeled.

Require passing native Go/OS, security and browser CI on the preview source.
The release workflow must also prove installation from the published Git tag,
separately from the local file-proxy packaging check. Automated accessibility
checks do not constitute a screen-reader or physical-device conformance review.

Keep focused evidence under ignored `tests/results/`; fetch and integrate
GitHub `main` before publishing. Private infrastructure references remain in
older library history, so review publication separately before changing visibility.
