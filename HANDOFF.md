# GothUI implementation handoff

## Native platform preparation and CI milestone (2026-11-06)

Added .github/workflows/verify.yml with six native jobs: Go 1.26.0 and 1.27.3 on
ubuntu-24.04, macos-15 and windows-2025. GOTOOLCHAIN=local prevents silent compiler
upgrades. Each job downloads/verifies modules, runs the standard check (including
the public API baseline, race suite and vet), checks pinned templ generation,
compiles recipes and executes the native demo's version smoke check. Jobs have
20-minute timeouts, read-only contents permission, no persistent checkout
credentials and seven-day report retention. No binaries, releases, tags or mirror
writes are published by the workflow. Official Actions release commit pins and
runner labels were verified through GitHub sources; docs/ci.md links provenance.
The workflow is prepared locally and has not run remotely. Local PyYAML parsing
and structural checks are not an actionlint pass or GitHub execution result.

The developer builder now uses .exe on Windows for the Go tool, candidate and
retained demo. It checks the candidate before replacement and records its actual
filename. Windows rename can fail when the old file is in use; it is never
pre-deleted. The existing Unix rename behavior remains. An integration test builds
and executes a temporary module under a path with spaces, then changes its source
and verifies a dirty-version replacement without candidate-file leftovers. The
test is included in the native CI matrix and has passed locally on Linux only.

The demo's -open now dispatches to Linux xdg-open, macOS open or the Windows URL
handler via rundll32.exe. The generated local URL remains a process argument,
never shell-interpolated. Desktop launches remain unverified. .gitattributes
normalizes text checkout line endings for native generated-source comparison.
No public API, CSS, templates or Go dependencies changed.

Validation: the full seven-package race suite and vet passed. After adding the
native build integration fixture, the affected command package passed race tests
and vet again. All 212 public API entries still match. The demo cross-compiled
with Go 1.27.3/CGO disabled for Linux arm64 (8,192,160 bytes), macOS arm64
(8,407,570 bytes), and Windows amd64 (9,249,792 bytes). Go build metadata and
ELF/Mach-O/PE architecture headers were verified. These are compilation/format
results, not execution on those foreign systems. The foreign artifacts were
removed immediately, and concise results remain under tests/results/platform.
The native demo is rebuilt once after the local commit to stamp its final source
revision; its version smoke result and size/hash remain in validation/build.json.

Delivery recheck: shell Git access to both remotes fails DNS resolution, and the
unarchived repository; its commits endpoint returns 409 "Git Repository is empty."

Housekeeping: test-owned module/Git fixtures and cross-build directories clean up
automatically. Remove the task-owned /tmp/gothui-platform-cache after final native
build. Retain one review binary, the icon contact sheet, current validation and
platform reports, and prior browser/TCP failure evidence. Superseded successful
API-command-only logs are removed after replacement checks. The continuing
developer owns retained evidence until review or replacement. No active user
process, reference source, shared cache, worktree or database was changed.

The broad goal remains active: minimum-version and non-Linux native CI, browser/
assistive-technology evidence, optional HTMX consumer execution, repository
publication and a clean published install still need proof.

## Public API baseline and migration milestone (2026-11-06)

Added the Go-only `gothui-dev api` command and included it in `check`. It records
compiler-resolved declarations for the root, components and render packages in
api/public.json: 212 declarations/method-set entries at this initial development
baseline. This is not a comparison with an earlier published or archived API.
No public library declarations, CSS, templates, runtime dependencies or demo
behavior changed in this milestone.

The internal/apiguard package captures exact constant values, function and
variable types, type definitions/aliases, generic constraints, struct fields and
tags, reachable local private types, and value/pointer method sets including
promoted methods. Generic instances do not overwrite their declaration's method
signatures. Deterministic comparison flags every difference, including additions.
Package discovery rejects new public import packages until the explicit package
list is extended; internal, cmd, examples and tests are intentionally excluded.
Imported type internals, behavior, CSS/HTML semantics and cross-platform results
are outside this declaration guard's proof.

The command writes api-current.json, api-changes.json and (only on a match)
api.json under tests/results/validation, with revision/dirty/toolchain/target
provenance on reports. It clears stale API evidence before each attempt and never
rewrites the checked-in baseline. docs/api-compatibility.md documents intentional
baseline review, consumer migration and annual module versioning. README,
development guide and roadmap distinguish declaration agreement from semantic
compatibility and retain browser/consumer/release gates.

Validation: the full seven-package race suite and vet passed. After the final
public-package lookup refactor, the affected command package passed race tests
and vet again. Fixtures cover removals/additions, fields/order/tags, constants and
iota, aliases, variadics, interfaces, generic constraints/instances, pointer and
promoted methods, inferred variables, recursive/private result types and package
changes. Body/comment-only changes and equivalent constant expressions compare
equal. A temporary real Go module/Git repository proves compiler-export loading,
baseline matching, rejection of an int-for-string field change, stale-success
removal, baseline preservation and rejection of an unreviewed public package.
Fixture repositories are automatically removed by testing.T.TempDir.

The existing f0f1f42 icon-demo binary is retained because this milestone changes
developer tooling only; its version remains truthful to the UI it contains.
No unnecessary demo rebuild, template generation or recipe recompilation was
performed. Their prior reports still identify their actual source revision.
The final standalone API check records the committed tooling revision. Remove
initial bootstrap/focused logs superseded by the completed suite and the task's
/tmp/gothui-api-cache after verification. Keep bounded validation evidence, one
review binary, the icon contact sheet and prior browser/TCP failure logs; the
continuing developer owns their replacement or removal after review. No user
process, shared cache, reference source, remote branch or release tag changed.

The broad goal remains active. Outstanding evidence includes browser and
assistive-technology behavior, optional HTMX consumer execution, supported Go/
platform CI, repository reconciliation and clean published installation. The
API gate is progress toward a stable library, not a completed release claim.

## Inline icons and action composition milestone (2026-11-05)

Added Icon as the 33rd public templ component, with 24 original geometric SVG
drawings on a consistent 24-unit grid. Typed IconName and IconSize values select
paths and text-relative dimensions. IconNames returns a fresh sorted list.
The paths are compiled into Go and rendered inline; no fonts, sprites, image
requests, runtime dependencies or JavaScript were added. CSS grows by 347 bytes
to 32,607 bytes. Icons inherit currentColor, remain non-focusable and have no
motion of their own. Directional drawings are literal; RTL applications choose
the appropriate counterpart.

Button and Link now accept IconBefore, IconAfter and AriaLabel. Icons inside
controls are decorative; visible control text supplies the name. Icon-only
controls require a nonblank AriaLabel. Unknown names and missing icon-only names
fail before component markup is emitted. IconNone emits nothing. Standalone
Icon.Label produces an escaped aria-label and role=img; the default is hidden
from assistive technology. Existing form types, commands, optional HX attributes
and disabled behavior are preserved. Use keyed props literals for additive API
changes; do not imply positional struct literals remain source-compatible.

The new `/catalog/icon` entry includes the full set, sizes, decorative/meaningful
image examples, labeled icon actions and a disabled upload example. Existing
Button and Link recipes also demonstrate icons. The catalog category Visuals,
downloadable recipe, README and docs/icons.md are updated. All recipes stay
self-contained with the public components import and embedded library CSS.

Validation passed: all six packages with the race detector and vet, all 31
generated templates reproducible, and all 33 recipes independently compiled.
Tests parse SVG namespace/dimensions/paths and naming attributes across every
drawing; check escaped hostile labels, invalid names, empty-label rejection,
safe size fallback, registry isolation and native action behavior. Existing
catalog HTTP tests exercise all new detail/preview/download endpoints.

A contact sheet was rendered through installed librsvg/Cairo document APIs and
visually inspected at 24 and approximately 38 pixels on light and dark surfaces.
That review caught and corrected a pencil end-cap alignment issue; the final
sheet reflects the corrected paths. The initial Python GI/Cairo bridge lacked
its foreign-context binding, and its pixbuf fallback required blocked D-Bus;
direct native document rendering succeeded without a display service. This is
static SVG drawing evidence, not browser layout/accessibility evidence. No
browser-control tool is available in this session; earlier sandbox launch
failure evidence remains retained. Browser/keyboard/forced-color/print review,
API compatibility automation and remote publication remain open.

The single demo executable is rebuilt with the final clean local Git version;
restart it outside the sandbox to explore `/catalog/icon`. Current validation
reports retain their actual revision and dirty-state provenance. Keep one SVG
and PNG contact sheet plus its provenance note under tests/results/icons until
review or replacement, owned by the continuing developer. Remove the task's
temporary proof generator, superseded focused-test log and task-owned Go cache
after final build. No user process, reference source, worktree, database or
shared cache is changed. The broad goal remains active and incomplete.

## Repeatable validation and performance milestone (2026-11-05)

Added Go-only `cmd/gothui-dev` commands: check, generated, recipes, bench and build.
The checks enforce Go formatting, stdlib-only root/render dependencies, templ-only
component runtime dependencies, a 64 KiB CSS ceiling, race tests and vet. Generated
source is compared byte-for-byte after regeneration in a temporary directory;
orphan output is rejected. Recipe extraction is shared between catalog downloads
and developer checks through internal/recipes; malformed/duplicate boundaries
are rejected. All 32 recipes compile in separate temporary packages.

Reports under tests/results/validation include revision, dirty state, toolchain
and UTC time. Each command removes its old success report before starting, so a
failure cannot leave stale success evidence. Logs are replaced, not accumulated.
Temporary generation/recipe directories are removed automatically on return.
The `build` command uses Git describe (or a pre-tag development fallback), includes
untracked work in dirty detection, rejects malformed version identifiers and
checks CGO-disabled build metadata, size, Linux ELF interpreter absence and the
new demo `-version` flag before replacing the one retained executable. It creates
no tags, commits, remote writes or running servers.

Validation passed: race tests for all six packages, vet, all 30 generated template
outputs reproduced, and 32 independent recipe compilations. New tests cover
recipe boundaries, version provenance including untracked state, malformed tags,
report metadata and exclusion of generated evidence from source walks. The
version smoke check executes successfully without a listener. First stripped
Linux amd64 build: 8,949,920 bytes, below the 16 MiB review ceiling. CSS remains
32,260 bytes. A final post-commit build records the delivered local revision in
build.json and the binary; other reports honestly retain the dirty pre-commit
source provenance under which checks ran.

Three-sample benchmarks with Go 1.27.3, Linux amd64 and GOMAXPROCS=2 provide an
initial baseline: median 45.2 microseconds for a preconstructed 100-row table
render, 34.1 microseconds for the dashboard handler and 4.81 microseconds for CSS
GET. HTTP cases include ResponseRecorder allocation; component cases discard
output. These do not measure network, browser, concurrency or application data
work. docs/performance.md retains measurements and variation, including the
table's 1,821 allocations as a candidate for later profiling. No timing threshold
or comparison with frontend frameworks is claimed.

Added components package documentation, docs/development.md, performance evidence
and README commands. Public API compatibility automation, supported Go/platform
CI, browser/assistive-technology evidence, optional consumer HTMX validation and
published installation remain outstanding. The broad goal stays active.

Housekeeping: superseded successful catalog logs are removed after the replacement
suite passes. Keep the bounded validation folder, one demo executable and prior
browser/TCP failure logs. The continuing developer owns them until replacement or
review completion. The task-owned /tmp/gothui-release-cache is cleaned after the
final build; shared caches, user processes, source/reference folders and repository
history are untouched. No new worktree, database or backup was created.

## Searchable component catalog milestone (2026-11-05)

Added the twelfth demo area `/catalog`, covering all 32 current public templ
components with native search, category filtering, bounded pagination, usage
notes, inline/standalone previews and downloadable self-contained templ recipes.
Layout and AppShell use standalone previews to avoid nested document/main
landmarks. The overview CTA/directory and workspace navigation link to the catalog.

Recipes are compiled from and embedded as `catalog_samples.templ`; marked blocks
supply both displayed and downloadable source. Downloads use a package declaration
and the public components import, with any local helpers contained in the recipe.
The registry supplies metadata and demo links. Only registered component paths
are served; no user input becomes a filesystem path. Search text is capped at
120 Unicode characters, categories validated and pages clamped. Theme-return
allowlisting recognizes registered catalog details/previews but excludes downloads.

Catalog CSS (1,366 bytes) is embedded in the example app and served at `/demo.css`,
leaving the library stylesheet unchanged at 32,260 bytes. Recipe source is 9,943
bytes. Script-blocking CSP remains unchanged; source downloads are plain-text
attachments with nosniff, and GET/HEAD semantics are covered. No frontend tools,
client runtime, external assets or Go dependencies were added.

Validation: all four packages passed `go test -race ./...`; `go vet ./...` passed.
After the final overview-link changes, the affected example package's race suite
passed again. Catalog tests cover every detail/preview/download (parallel subtests),
HEAD, document landmarks, escaping, search/pagination, attachment headers and
registered paths. A coverage test compares exported component declarations in
components/*.templ against the registry, failing if a public templ component is
missing. This does not claim coverage of every exported Go helper/type.

All 32 downloadable recipes were generated and compiled in separate packages in
a temporary module with a local GothUI replacement. The generator initially noted
that templ was not a direct requirement in that temporary module; Go resolved the
existing pinned transitive runtime and all 32 packages compiled successfully.
This proves standalone recipe compilation, not published-module installation.
The task-owned temporary recipe directory is removed after validation; concise
generation/build logs remain as provenance. Browser review remains unverified.

The statically linked review executable is rebuilt and atomically replaces
`tests/results/demos/gothui-demo`; existing user processes are untouched. Restart
from a normal terminal, then visit `/catalog`. Current catalog validation/recipe/
build logs replace successful forms logs; retain prior browser/TCP failure logs.
The continuing developer owns the single binary and bounded logs until review or
replacement. The task-specific Go cache is cleaned after building. No new worktree,
database, backup, server process or release tag was created.

The broad goal remains active. ROADMAP.md still requires browser evidence, API/
compatibility documentation, measured build/render/size gates, and safe exact
remote publication. The next useful milestone is reproducible validation/build
commands and performance evidence, rather than treating the catalog as a release.

## Composable fields and native wizard milestone (2026-11-05)

Added Field with typed label/help/error/required metadata and fresh control
attribute maps, plus InputGroup for escaped decorative prefixes/suffixes. Field
requires an explicit non-whitespace ID and label; the application owns uniqueness
and spreads ControlAttributes onto one raw native control. Existing convenience
Input/Select/Textarea components remain unchanged. Button.FormNoValidate emits
only for submitters, enabling Back actions without blocking native constraints.

The eleventh demo `/wizard` provides a three-step project-brief flow: details,
preferences, review, then a validated preview. Draft values travel through native
visible/hidden form controls; Back preserves incomplete edits and choices. Every
forward transition rechecks prior details; Finish validates all fields again.
Hidden input is not treated as trustworthy or proof of an authorized workflow.
No account, cookie, database, project creation, publication or external sending is
involved. GET always starts fresh; browser history can retain submitted values.

The POST handler bounds bodies to 8 KB, rejects duplicate known fields, invalid
UTF-8, malformed step/action combinations and unknown preferences, and shares the
server's cross-origin protection. Field errors return 422 on the relevant step;
malformed requests return 400, oversized bodies 413 and wrong methods 405. The
forward submitter precedes Back in DOM order to make implicit Enter advance.
The final preview is intentionally non-mutating; production ordered approvals,
idempotency, persistent state and authorization remain application responsibilities.

Validation: templ formatting/generation completed; all four packages passed
`go test -race ./...`; `go vet ./...` passed. Tests use actual generated hidden
inputs to exercise details -> preferences -> Back -> edited details -> review ->
finish, and cover escaping, field association/attribute isolation, native
submitter flags, hidden-state tampering, final revalidation, malformed/oversized
bodies, cross-origin rejection, method handling and fresh-request independence.
The formatter normalized an existing Disclosure block without behavior changes.
Embedded CSS is 32,260 bytes. Browser keyboard/focus/screen-reader/layout testing
remains pending; HTTP tests do not establish those results.

The one review binary is rebuilt statically and replaced atomically at
`tests/results/demos/gothui-demo`. Restart from a normal terminal, then visit
`/wizard`. Current forms check/build logs supersede successful data-suite logs;
keep prior motion-browser.log and demo-launch.log for unresolved sandbox limits.
The continuing developer owns the binary/current logs until review or replacement.
The task's temporary Go cache is cleaned after building; no extra worktree,
database, backup or running server was introduced.

The broad goal remains active and incomplete. ROADMAP.md keeps its release gates;
remote publication and browser proof remain outstanding. Next useful work is a
searchable native component catalog with copyable recipes, then measured build/
performance and compatibility gates. The form guide is at docs/forms.md.

## Data and navigation composition milestone (2026-11-05)

Added DataTable with required captions/column labels, row shape validation,
escaped text or templ component cells, a single active sort header and optional
row header column, alignment, compact/striped styles, visual selected rows,
loading status outside the busy table, and an empty-state slot. It retains native
table semantics inside a labeled keyboard-reachable scroll region. Sorting and
selection remain application-owned; no client data model or runtime was added.

Added Tabs (ordinary page links/current/disabled state), Steps (ordered progress,
current step, completed labels and application-owned destinations), and Disclosure
(native details/summary with optional group name). Contradictory navigation state
fails rendering. The new tenth demo `/patterns` provides page sections, a stateless
three-step reading guide, table loading/empty previews and grouped disclosures.
`docs/data-navigation.md` documents APIs, accessibility contracts and boundaries.

The dashboard now composes DataTable and status Tabs. Header links toggle both
sort directions and preserve query/filter settings while returning to page one.
Native checkboxes submit the visible-page selection to `/dashboard/selection`.
That bounded, cross-origin-protected POST validates known sample names,
deduplicates and produces a read-only receipt. No data, cookie or project changes
occur; empty selection returns 422, invalid/excessive input 400 and wrong methods
405. A production action would need authorization and persistent record IDs.

Validation: all four packages pass `go test -race ./...`; `go vet ./...` passes.
Tests cover column/row invariants, safe cells/URLs, navigation semantics/conflicts,
loading/empty markup, sorting directions, query preservation, native selection
validation/no mutation, cross-origin rejection and new demo request methods.
An old filter test assumed empty results removed the whole table; it now checks
that empty results have no selectable rows, matching the intentional retained
caption/header/empty-cell structure. A CSS review fixed inherited semantic stripes
so DataTable.Striped=false really removes striping. CSS is 30,549 bytes.

The review binary is rebuilt statically with CGO_ENABLED=0 and replaced atomically
at `tests/results/demos/gothui-demo`. Existing user processes are left alone;
restart from a normal terminal and visit `/dashboard` and `/patterns`. Current
data-suite/build logs supersede successful motion logs. Retain the unresolved
motion-browser.log and TCP demo-launch.log evidence. The continuing developer
owns the one binary/current logs until review or replacement; this task's Go cache
is removed after building. No new worktree, database, backup or daemon was created.

This is verified local progress toward the active broad goal, not completion.
Browser/assistive-technology validation and remote publication remain outstanding.
Prior sandbox browser/socket failures are unchanged and have not been relabeled
as passing checks. Next roadmap work includes composable fields and a real native
multi-step form, followed by a searchable component catalog and measured release
gates. Keep ROADMAP.md's full acceptance criteria intact.

## Motion and expressive surfaces milestone (2026-11-05)

The broader user goal is now active: build an elegant, comprehensive Go-native
UI library in a single binary, extending beyond Pico. `ROADMAP.md` preserves the
full objective, development tracks and release gates. This milestone advances
that goal; it does not establish feature parity, superiority, accessibility
conformance or release completion.

Added typed card surfaces (tinted, outlined, gradient), one-shot fade/rise/scale
entrances with bounded delays, a Reveal wrapper, two-pixel hover/focus lift,
short control transitions and native overlay opening effects. Text remains on
opaque surfaces. MotionSystem respects OS reduced motion; MotionOff zeroes all
motion durations/delays and hover displacement, including CSS/Go exports and
saved preferences. Reduced-motion and print styles disable animations. No JS,
external assets or new dependencies were introduced.

Added Skeleton (bounded decorative shapes, visible status, three pulse cycles),
EmptyState (recovery action slot) and Stat (description-list metric). Dashboard
empty results use the new component. The ninth demo, `/motion`, provides replay,
motion/theme/preset controls, all new surfaces/states, and native overlays.
`docs/motion.md` documents semantics, duration tokens, usage and manual checks.
The operator's PureCSS-Animations repository was reviewed as visual reference;
the authored token-based styles do not bundle its source or its JS carousel.

Verification: templ formatting/generation completed, `go test -race ./...` passed
for all four packages and `go vet ./...` passed. Tests cover safe props, bounded
placeholder output, semantic labels/escaping, motion preference save/preview/
export, native route behavior and CSS policy guards. Embedded stylesheet size is
27,791 bytes at this milestone. CSS-string tests do not prove browser behavior.
An actual installed headless Chromium launch failed before rendering with
`sandbox_host_linux.cc: shutdown: Operation not permitted`; browser compatibility,
visual, keyboard and assistive-technology gates remain pending.

The review binary is rebuilt with CGO_ENABLED=0 and replaced atomically at
`tests/results/demos/gothui-demo`, preserving any user-owned running process.
Restart it from a normal terminal and visit `/motion`. Retain the current motion
validation/build logs and browser-launch failure, plus the earlier TCP failure;
remove superseded editing logs and the task-specific Go cache. The continuing
developer owns these artifacts until review or replacement. No extra worktree,
backup, test database, server process or release tag was created. Remote delivery
remains pending under the documented network restrictions. Next work should
follow ROADMAP.md, especially reusable data/navigation patterns and browser gates.

## Pages, fragments, and inline editing milestone (2026-11-05)

Added the standard-library-only `render` package. `render.HTML` accepts templ
components through a minimal interface, selects full pages or optional fragments,
merges Vary, respects HEAD, buffers rendering, and sends generic 500 errors
without exposing partial markup or internal errors. `render.IsFragment` excludes
history restoration and boosted navigation. The helper owns no authentication,
cache policy, cookies, CSP or validation status policy. Existing demo rendering
uses this helper; the root asset package stays standard-library-only.

The eighth page, `/editing`, demonstrates native edit/save/cancel/reset flows
and an optional HTMX fragment contract using the existing component attributes.
Fictional profile values are bounded, validated, escaped and stored only in an
unsigned HttpOnly browser session cookie, with SameSite=Strict and Secure over
TLS. This is cosmetic sample storage, not an authentication/session design.
Restoring the sample clears the cookie. Invalid submissions preserve entered
values and show linked errors without writing a cookie. Native errors return
422; HTMX form fragments explicitly return 200 to work with default response
handling. Successful native POSTs redirect; fragment POSTs return the panel.
Requests retain the example's cross-origin protection and a 4 KB body limit.
No scripts or dependencies were added. Native selects and appearance are unchanged.
The formatter also normalized an existing dashboard table header.

Verification: `go test -race ./...` passed for all four packages and
`go vet ./...` passed. Tests cover full/partial/history responses, merged Vary,
HEAD, error buffering, canceled contexts and transport errors, edit form hooks,
save/reset/cancel destinations, isolation, escaping, cookie limits and TLS flags,
validation, method/body limits, and cross-origin rejection. Generated templ Go
is included. These are component and HTTP contract tests, not HTMX browser tests.
Browser review remains manual because the session cannot bind TCP sockets.

The one review executable at `tests/results/demos/gothui-demo` is rebuilt and
replaced atomically; no existing user server is terminated. Restart it from a
normal terminal, then open `/editing`. Keep only the four current editing check
logs and the earlier socket-launch failure alongside this binary. Earlier passing
navigation logs are superseded by this full suite; their results remain below.
The continuing developer owns these review artifacts and can replace/remove them
after review. The task-specific temporary Go cache is cleaned after building.
Remote publication remains pending under the previously documented network
restriction. No tag or release is created. Richer table abstractions and testing
with a consumer-supplied HTMX runtime remain future work.

## Native interactions and navigation milestone (2026-11-05)

Added Dialog, Popover, Breadcrumbs, and Pagination with generated Go. Button
accepts typed native dialog commands and popover actions; native targets always
render as type=button. Invalid overlay IDs fail rendering. Dialogs include a
native close form and labels; popovers have automatic/manual modes with explicit
close controls. No scripts or dependencies were added, and native selects remain
unchanged as accepted by the operator.

The new `/interactions` page demonstrates popovers, a modal, and a bounded,
cross-origin-protected confirmation POST that saves/mutates nothing. An ordinary
page provides the same confirmation flow without newer dialog command support.
The example now has seven pages. The dashboard adds sorting, page size, and real
pagination while retaining search/status/appearance in page URLs. Sorting uses
a request-owned slice and does not mutate shared sample data.

Component/HTTP tests cover native attributes, safe labels/URLs, render errors,
pagination boundaries including maximum integers, bounded output, filters/sort,
confirmation methods/body limits, fallback rendering, and cross-origin rejection.
Race-enabled tests and vet are part of this milestone's verification. Native
focus, keyboard behavior, and browser layout/compatibility remain pending manual
review; this restricted session cannot start the local HTTP server. Earlier
operator approval establishes the appearance/native-select baseline only.

The current demo executable is replaced atomically at
`tests/results/demos/gothui-demo`; existing demo processes remain untouched and
must be restarted to load this build. One executable and concise test logs are
retained for operator review, owned by the developer continuing this task and
replaceable/removable after review. The task's temporary Go cache is removed.
Earlier successful-run logs are removed after the new full validation; their
results remain summarized below. The latest four navigation validation logs and
the earlier socket-launch failure log remain as bounded review evidence.
Remote delivery remains pending; advanced HTMX workflows and richer table
abstractions remain later roadmap work.

## Appearance and application components milestone (2026-11-05)

Implemented independent typed appearance settings: four accent presets
(Indigo/Neutral/Forest/Rose), density, corners, system font family, text size,
content width, and card shadows. Defaults preserve the original appearance.
CSS tokens are shared by attribute styling and downloadable CSS overrides.
Layout accepts Appearance and an optional trusted Head slot.

New components: Badge, Alert, ErrorSummary, Textarea, RadioGroup, Switch,
FileInput, and AppShell. Checkbox gains error/description attributes; Form
supports multipart encoding. The form lab now uses the reusable textarea,
checkbox validation attributes, and linked error summary. No JS or new Go
dependencies were introduced.

The demo now has six pages. `/settings` previews/saves/resets preferences and
exports CSS/Go settings; `/workspace` demonstrates the app shell. The playground
exposes appearance settings and all new controls, badges, and alerts. Preferences
are normalized before use; there is no library-owned state or storage. Settings
POSTs share the example's cross-origin protection. File selection is a local
preview only, not an upload service.

Race-enabled tests pass for rendering, escaping, control semantics, preference
cookies, normalization, overrides, exports, and cross-origin protection. Preset
button/link color pairs are checked for at least 4.5:1 contrast in normal/hover
light/dark states; this is not a full accessibility audit. Generated source and
CSS token references are checked. Browser/layout review remains outstanding
because this session cannot bind sockets. The operator's earlier positive
browser review covered the prior demos, not this milestone.

The single retained review binary at `tests/results/demos/gothui-demo` is
replaced with this milestone's build. An already-running older binary is not
stopped; the operator must restart it to load new embedded CSS and routes.
The continuing developer owns removal/replacement of this binary after review.
Temporary compilation caches are removed at handoff; concise local evidence
remains under `tests/results/` for the next review/run.

Repository publication remains pending under the earlier delivery notes. This
milestone covers appearance, feedback, controls, and shell foundations. Advanced
HTMX interactions, dialogs/popovers, richer table abstractions, and other items
from the broader roadmap are not included in this milestone.

## Demo expansion (2026-11-05)

The basic example now serves four connected pages: overview (`/`), component
playground (`/playground`), form lab (`/forms`), and sample project dashboard
(`/dashboard`). Playground options and dashboard filters use native GET forms.
The form lab exercises browser validation, server errors, and success without
saving data. Theme toggles return to a validated same-origin demo Referer and
preserve filters while removing explicit theme overrides.

The `-open` flag launches the default Linux browser after binding the listener;
`-addr` selects a port. A launch attempt still failed with `socket: operation
not permitted`, so no browser session or server is running. The pages are
implemented and HTTP-tested but still need desktop browser review.

`go test -race ./...` passed after the expansion. Tests cover routes, rendered
control states, themes, input escaping, successful/invalid form submissions,
oversized bodies, cross-origin rejection, project filters, and safe redirects.
See `examples/basic/README.md` for launch and manual review steps.

One executable is retained under `tests/results/demos/gothui-demo` for operator
review. Launch it with `-open` from a normal desktop terminal. The next demo
build should replace it; the developer completing review owns its removal
after review is complete. The demo build cache is temporary and removed at
handoff. Concise demo logs remain under `tests/results/` until the next run.

## Implemented

- Embedded, concurrent-safe CSS handler with ETag revalidation, GET/HEAD,
  explicit mount path, and no directory listing.
- Original semantic CSS: light/dark/system themes, responsive grids, buttons,
  cards, native forms/disclosures, tables, visible focus, and reduced-motion rules.
- Typed templ sources and generated Go for layout, theme toggle, buttons, links,
  cards, forms, input/select/checkbox controls, and grids.
- Native form theme switching with validated cookies, cross-origin protection,
  and an HTTP-level HTMX refresh contract. No bundled JavaScript.
- Runnable showcase, compiled quick-start snippet, MIT licensing and Pico attribution.

The operator explicitly accepted templ as the only direct Go dependency and
no Node/npm/bundled JavaScript. The pinned runtime requires Go 1.26.0 or newer.
The module is `github.com/Tend-Stack/GothUI/v2026` so the requested annual
`v2026.1.0` tag will be compatible with Go semantic import versioning.

## Verification on 2026-11-05

- `templ fmt .`: passed, six source files.
- `templ generate`: passed; generated Go is included.
- `go test -race ./...`: passed for all three packages.
- `go vet ./...`: passed.
- README's complete main.go snippet: compiled using `go test` in a temporary
  package; temporary package removed afterward.
- Runtime dependency inspection: only the standard library and templ packages.

Tests exercise asset mounting, caching, HEAD/404/405 responses, concurrent
requests under the race detector, escaping, unsafe link rejection, optional
HTMX attributes, semantic slots, labeled controls, selected/checked states,
cookie theme round trips, the HTMX response header, request methods, and
cross-origin POST rejection.

Live browser verification is outstanding. The sandbox refused the example's
TCP listener (`socket: operation not permitted`). The example compiled and its
real HTTP handler was exercised in-process via httptest; this is not browser
or visual validation. No example process remains running.

The system's `/snap/bin/go` wrapper also refused to run in this sandbox. Checks
used the existing `/snap/go/current/bin/go` executable (Go 1.27.3) and, for templ
formatting, placed that directory first in PATH. No system installation changed.

## Repository delivery remains outstanding

The nested project has its own Git repository and these remotes:

- github: https://github.com/Tend-Stack/GothUI.git

restricted shell. The GitHub connector reported an empty, private, unarchived
content or repository settings were changed, and no release tag was created.

Next owner: the developer continuing this task with repository network access.
this independent initial history over it. Unarchive only if necessary for the
requested conversion. Rerun affected checks after integration, validate the
showcase in a browser with scripts disabled at desktop/mobile sizes and both
themes, and publish matching commits to both remotes. Then create the first
annual release tag. Private module installation also requires authentication;
repository visibility has deliberately not been changed.

## Housekeeping

The parent `pico-main/`, `pico-main.zip`, and `converzation` are operator inputs
and remain untouched. This nested checkout is the deliverable, not a disposable
worktree. No backups, extra worktrees, databases, or dependency installs were
created. The task-owned temporary Go build cache and quick-start compilation
fixture were removed after validation.

Small logs under `tests/results/` are ignored by Git and retained as local
verification evidence. The next task owner may replace them on the next run or
remove them after delivery; this file preserves the useful result summary.
