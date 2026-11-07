# TendUI preview verification

The supported preview is `v2026.2.0-preview.2`. Install it explicitly:

```sh
go get github.com/Tend-Stack/TendUI/v2026@v2026.2.0-preview.2
```

## Published preview

The following GitHub Actions runs passed for the published preview. These links
identify the tested revision directly; later documentation updates do not extend
those results to unrelated code changes.

| Gate | Evidence | Result |
| --- | --- | --- |
| Native Go compatibility and security | [Verify TendUI](https://github.com/Tend-Stack/TendUI/actions/runs/42313429148) | Linux, macOS and Windows on Go 1.26.0 and 1.27.3; independent vulnerability check passed. |
| Browser behavior | [Browser verification](https://github.com/Tend-Stack/TendUI/actions/runs/42313429136) | 59 checks passed; one unsupported WebKit forced-color emulation check skipped. |
| Installation from the published tag | [Release verification](https://github.com/Tend-Stack/TendUI/actions/runs/42313876893) | Published-tag consumer and quickstart checks passed. |
| Release notes | [TendUI preview](https://github.com/Tend-Stack/TendUI/releases/tag/v2026.2.0-preview.2) | Prerelease; APIs may change before stable release. |

On 2026-11-07, anonymous GitHub source access and the Go proxy module download
were checked again. The downloaded module checksum matched the website's pinned
`go.sum`. Both repositories are public; installation does not require a GitHub
token or private-module configuration.

## Website and browser coverage

The website consumes the pinned preview. Its [publication-polish browser run](https://github.com/Tend-Stack/TendUI-web/actions/runs/42337550182)
passed 20 checks on desktop Chromium and mobile WebKit. These include native
navigation, theme changes, onboarding downloads, containment, and light/dark
hover and keyboard-focus treatments. The [website build checks](https://github.com/Tend-Stack/TendUI-web/actions/runs/42337550151)
also passed for that tested revision. CI success is separate from deployment;
verify the served content after a website change is merged.

The library suite exercises Chromium, Firefox and WebKit. Native flows use
scripts-disabled contexts. Optional HTMX integration tests use a local,
test-only runtime; TendUI does not bundle that runtime with the library.
Automated checks cover accessibility rules, keyboard behavior, native overlays,
reduced motion, motion opt-out, responsive containment and print styles.

## Limits and repeatability

These checks do not certify WCAG conformance, every screen reader, every physical
device, every zoom level, or every consumer customization. Manual assistive
technology and physical-device reviews remain necessary. WebKit forced-color
emulation is unsupported in the test environment and is reported as skipped.

Workflows and fixtures are committed for repeatability. CI artifacts expire after
seven days; maintainers retain concise release evidence separately. Screenshots
and PDFs are review evidence, not pixel-by-pixel regression baselines. See
[development](development.md), [CI](ci.md), and [accessibility](accessibility.md)
for the checks to run when changing components or applications.
