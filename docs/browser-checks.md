# Browser verification

`Browser verification` runs Chromium, Firefox, and WebKit against an isolated Go
consumer fixture. It exercises native forms with JavaScript disabled, keyboard
popover/dialog behavior, catalog search and downloads, profile validation,
theme persistence, mobile layouts, reduced motion, print, and forced colors
where the browser supports emulation. Axe checks all 37 component detail pages
in light and dark modes. Test failures are not retried or silently excluded.

A separate test-only consumer loads pinned HTMX 2.0.11 with default
response-status handling. It checks fragment editing, validation forms returned
with 200, and the theme redirect contract. Native invalid submissions use 422.
The fixture is under `tests/browser/fixture`; the normal example handler never
loads this runtime. The browser suite's Node packages are development tooling,
not a requirement to build, consume, or deploy TendUI.

## Run locally

From the repository root, build the fixture:

```sh
go build -o tests/results/browser/fixture ./tests/browser/fixture
```

Then, from `tests/browser`:

```sh
npm ci --ignore-scripts
npx playwright install --with-deps chromium firefox webkit
npm test
```

The test runner starts and stops its own loopback server. It does not reuse an
unidentified process on the port. Reports, failure traces, screenshots, and a
Chromium print PDF stay under `tests/results/browser`; CI retains them seven days.
Do not check binaries, browser installations, or `node_modules` into Git.

## Scope and remaining manual review

Browser engine tests do not represent physical mobile devices or all assistive
technologies. Axe is an automated rules check, not a conformance certification.
Screen-reader announcements, visual PDF review, and target-device review remain
manual release evidence. WebKit's unsupported forced-colors emulation is an
explicit skip; supporting engines still exercise the controls in that mode.

Inspect the run for the exact proposed release commit. A workflow definition
or an earlier green run does not establish that current code passed.

## Reusing examples

`examples/site.Handler()` exposes the catalog and native demo workflows for the
documentation website. Mount it at the origin root because links and forms use
absolute paths. This example package is not part of the stable component API.
`go run ./examples/basic` remains the local executable entry point. Demo state
uses bounded browser cookies and never constitutes production authentication.
