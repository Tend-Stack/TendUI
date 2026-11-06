# GothUI 2026.1.0 verification

Verified on 2026-11-06. The immutable `v2026.1.0` tag points to
`2f622c95fbf71df942822d6697093b15fc8270db`. Later documentation-only commits
record these results without moving the release tag.

| Gate | Evidence | Result |
| --- | --- | --- |
| Native Go compatibility and security | [Verify GothUI](https://github.com/Tend-Stack/GothUI/actions/runs/42211995225) | Six native Linux/macOS/Windows jobs on Go 1.26.0/1.27.3, plus the vulnerability job, passed at the tagged commit. |
| Real browser behavior | [Browser verification](https://github.com/Tend-Stack/GothUI/actions/runs/42211995320) | 53 passed, zero failed/flaky; one unsupported WebKit forced-color emulation check skipped. |
| Installation from the actual tag | [Release GothUI](https://github.com/Tend-Stack/GothUI/actions/runs/42213916698) | Fresh consumer downloaded `v2026.1.0`, passed `TestPublishedModule` and `go mod verify`, with no replacement or templ generation. |
| Release publication | [v2026.1.0](https://github.com/Tend-Stack/GothUI/releases/tag/v2026.1.0) | Published as a regular release after the tag-installation check passed. |
| Documentation website | [Website browser checks](https://github.com/Tend-Stack/GothUI-Web/actions/runs/42212696656), [native/container checks](https://github.com/Tend-Stack/GothUI-Web/actions/runs/42212696776) | Eight browser checks passed on desktop Chromium and mobile WebKit at website commit `14874b6286e1d5cd77fc2b1628814c60a5212ae6`; native tests and the actual container smoke test passed. |

## Browser scope

The pinned Playwright 1.65.0 suite runs Chromium, Firefox and WebKit on the
Ubuntu CI runner. Native flows use scripts-disabled contexts. The separate
HTMX fixture loads a local, test-only HTMX 2.0.11 runtime and verifies fragment
editing, validation recovery and theme redirects. Neither project bundles this
runtime into its shipped application.

Axe checks WCAG 2 A/AA and 2.1 AA rules on all 36 component detail pages in both
themes. Application containment is checked at 320, 768 and 1920 pixels. Focused
checks cover keyboard skip/focus, popovers, native modal focus and background
inertness, reduced motion, explicit motion off, forced colors and print tokens.
The Chromium dashboard PDF was reviewed: readable monochrome controls and table,
no horizontal clipping, and headings stay with their following content.

The website suite checks nine documentation/home/example routes, native theme
navigation, catalog discovery, keyboard skip navigation, reduced motion, mobile
containment, and the 1920-pixel maximum at a 2560-pixel viewport. Code examples
are keyboard-focusable for horizontal scrolling. Screenshots and PDFs are
review evidence, not pixel-by-pixel regression baselines.

## Limits and repeatability

These results do not certify WCAG conformance, every screen reader, every
physical device, every zoom level, or every consumer customization. WebKit
forced-color emulation is explicitly unsupported. Assembled applications still
need testing with their target assistive technology and supported devices.

The repositories currently require GitHub authorization. Tag installation was
verified with authorized credentials, not anonymous public access. Consumers
need normal Git credentials and `GOPRIVATE=github.com/Tend-Stack/GothUI` while
that visibility remains in effect.

Workflows and fixtures are committed for repeatability. CI artifacts expire
after seven days; focused local reports and representative render artifacts live
under ignored `tests/results/`. No credentials or private infrastructure details
are included in the repository or retained test reports.
