# TendUI demos

Twelve connected demo areas, including a component catalog, served by one Go executable. All pages use the embedded
TendUI stylesheet and generated templ components. No frontend tooling or
JavaScript is required; the response CSP blocks scripts.

From the repository root:

```sh
go run ./examples/basic -open
```

To build a self-contained executable with the stylesheet embedded:

```sh
go run ./cmd/tendui-dev build
./tests/results/demos/tendui-demo -version
./tests/results/demos/tendui-demo -open
```

`-open` opens the playground after the listener starts, using `xdg-open` on Linux,
`open` on macOS or the Windows URL handler. Without that flag, open the links
below manually. On Windows, the built executable is `tendui-demo.exe`.
Choose a different port with `-addr 127.0.0.1:8081`. Stop with Ctrl+C.
The build command checks size and version before replacing the review binary.
`-version` prints the embedded Git identifier and exits without listening.
See [the development guide](../../docs/development.md) for reproducible checks.

| Page | Local address | Try it |
| --- | --- | --- |
| Overview | http://127.0.0.1:8080/ | Component introduction and links to each demo |
| Catalog | http://127.0.0.1:8080/catalog | Search 33 components, filter categories, open previews and download templ recipes |
| Playground | http://127.0.0.1:8080/playground | Themes, button labels/variants, disabled/invalid/read-only states, card header/footer, checked controls |
| Form lab | http://127.0.0.1:8080/forms | Required fields, invalid email, server validation, reset, success receipt |
| Dashboard | http://127.0.0.1:8080/dashboard | Search, status filters, sorting, page size, pagination, empty results and scrollable tables |
| Appearance | http://127.0.0.1:8080/settings | Presets, spacing, corners, typography, text size, width, shadows, saved preferences, CSS/Go downloads |
| Workspace | http://127.0.0.1:8080/workspace | Reusable app header/sidebar/main layout, native collapsible navigation, feedback and cards |
| Interactions | http://127.0.0.1:8080/interactions | Automatic/manual popovers, native modal, confirmation POST, ordinary-page fallback, breadcrumbs |
| Inline editing | http://127.0.0.1:8080/editing | Edit, save, cancel, restore, linked validation errors and optional fragment responses |
| Motion & surfaces | http://127.0.0.1:8080/motion | Replay entrances, compare motion off, card surfaces, skeleton/empty/metric states and overlay entrances |
| Patterns | http://127.0.0.1:8080/patterns | Page tabs, step walkthrough, table loading/empty states and grouped native disclosures |
| Form wizard | http://127.0.0.1:8080/wizard | Details, preferences, review, Back edit preservation, linked errors and final validation |

The playground uses GET controls: click **Apply settings** to refresh the preview.
Its URL preserves settings so you can bookmark a test case. The preview button
submits a native form and confirms the click. Its fields demonstrate visual and
native interaction states; the form lab tests actual submission validation.

Expand **Appearance settings** to combine visual options without saving them.
The added component section shows textareas, radio groups, switches, file inputs,
all five badge tones, and alerts. Choosing a file remains local to your browser.

The Appearance page has **Preview changes**, **Save preferences**, and **Restore
defaults**. Preview affects only the URL; saving stores preferences in HttpOnly,
SameSite=Strict cookies, with Secure enabled over TLS. The library itself stores
no settings. Downloads reflect the currently displayed preview, so apply or save
edits before downloading. CSS overrides require the base TendUI stylesheet; the
Go export also includes light/dark/system mode. The Workspace page uses the
reusable `components.AppShell` and the same saved preferences as other demos.

Theme selection in the playground previews light, dark, or system preferences.
The navigation toggle saves a cookie, removes the preview override, and returns
to the current demo with other filters preserved when the browser sends a
same-origin Referer. Without Referer it returns to the overview. When using
system preference, the server cannot detect the OS theme; the first toggle
uses its saved/default theme as the starting point.

The form lab accepts POSTs but saves no data and sends no messages. **Submit
demo** uses browser validation. **Test server validation** deliberately bypasses
browser checks to exercise HTTP 422 field errors. The server validates the name,
email, role, consent, and message length and limits request bodies to 16 KB.
**Reset edits** restores values from the current document; **Start a fresh test**
clears the form. Successful submissions show an escaped, in-memory receipt.

Dashboard data is fixed sample content. Filters, sorting, and pagination work on
the server, and portfolio totals intentionally remain independent of the filtered
table. The default page size is two projects; choose five to see all records.
Page links preserve search, status, sorting, size, and appearance query settings.
Submitting filters starts again at the first page. Sorting never mutates shared
sample records. Due-date sorting is specific to this sample's November dates.

Interactions use native HTML with no scripts. Automatic popovers support
light dismissal; manual popovers require their Close button or another toggle.
The dialog opens using native invoker commands and includes a non-destructive
close control. The approval form POSTs to `/interactions/confirm` and redirects
to a success message without changing or saving records. The ordinary-page
confirmation link exercises the same handler in browsers without dialog command
support. Native browser behavior still needs manual compatibility testing.

Inline editing stores a fictional display name/introduction in a bounded session
cookie (HttpOnly, SameSite=Strict, Secure over TLS). **Cancel** reloads the last
saved sample; **Restore sample** deletes that cookie. No account, database or
shared server-side profile exists. Session cookies may survive browser restart
when the browser restores its session, so use Restore to clear the sample.
POST bodies are limited to 4 KB; names allow 1–60 Unicode characters and optional
introductions allow up to 240. Unsaved invalid values remain in the form.

The native flow performs full-page navigation. The same handler returns only
`#profile-panel` when `HX-Request: true` is supplied; boosted navigation and
history restoration receive full pages. HTMX attributes remain inert because
this demo does not load HTMX and its CSP blocks scripts. Verify the server
contract independently with:

```sh
curl -i -H 'HX-Request: true' 'http://127.0.0.1:8080/editing?edit=1'
curl -i -H 'HX-Request: true' -H 'HX-History-Restore-Request: true' 'http://127.0.0.1:8080/editing?edit=1'
```

Native validation errors use HTTP 422; HTMX validation fragments use 200 so a
consumer's default HTMX response handling can display them. Neither writes a
cookie. Successful native saves redirect; successful HTMX saves return the panel
and a cookie. The library helper does not choose this application status policy.

Motion & surfaces uses one-shot CSS entrances and short hover/focus feedback.
Use **Apply & replay entrances** to reload with a chosen mode and theme. The
skeleton is clearly labeled as a preview with no request running; decoration
stops after three cycles. Motion follows OS reduced motion unless explicitly
turned off, and the Appearance page saves the setting across demos. Surface
variants preserve opaque text backgrounds. See [the motion guide](../../docs/motion.md).

The Dashboard now uses the reusable DataTable. Click a sortable column header to
change direction; the links preserve query settings and restart at page one.
Status tabs offer another way to filter. Select rows on the current page and
choose **Review selected projects** for a read-only receipt. Selection is not
saved across pages. Unknown names and excessive submissions are rejected; a
blank selection produces a useful error. No project is changed or sent anywhere.

Patterns contains link-based tabs, a three-step reading guide, loading/empty
table previews and native grouped disclosures. Loading is labeled as a preview;
no request is running. The walkthrough URL chooses its step and saves no data.
See [the composition guide](../../docs/data-navigation.md) for the public APIs.

The Form wizard builds a fictional project brief over three native POST steps.
Use Back to keep edits, Continue to validate the current and prior details, and
Finish preview to recheck everything. Completion creates or publishes nothing.
Draft values travel with the form, not in cookies or a database. A fresh GET
starts over. Browser history may retain submissions; use fictional information.
See [the form guide](../../docs/forms.md) for validation and state ownership.

The Catalog covers all public templ components. Search is case-insensitive across
names, categories and notes, with twelve results per page. Open an entry to try
its preview and select or download the example source. The same source is compiled
for the preview; adapt its package, IDs and demo destinations in another app.
Layout and AppShell open standalone to preserve document landmarks. All recipe
text and the example-only `/demo.css` stylesheet are embedded in this binary.
See [the catalog guide](../../docs/catalog.md).

## Manual review

1. Test at desktop width and at 375 px. Cards should stack; only the table should
   need horizontal scrolling.
2. Test light/dark themes and system preference. A saved light theme must remain
   light even when the OS is dark.
3. Disable JavaScript. Apply playground settings, submit the form, filter the
   dashboard, expand disclosures, and toggle the theme.
4. Use Tab, Shift+Tab, Enter, and Space. Verify visible focus, labels, disabled
   controls, the skip link, and the scrollable table region.
5. Test blank fields, malformed email, long names, a successful submission,
   and angle brackets in name/message text.
6. Search for `portal`, filter `paused`, and search for a nonexistent project.
7. Confirm no external assets are requested and no script is loaded.
8. On Appearance, try Forest/compact, Rose/soft corners, Neutral/monospace,
   large text, narrow/wide content, and flat/raised cards in both light and dark.
   Save preferences, navigate to Workspace and Form lab, then restore defaults.
9. Download CSS and Go settings after previewing changes. Load CSS after the
   base stylesheet, or use the exported Appearance and Mode with LayoutProps.
10. Test the workspace skip link and collapse/expand sidebar navigation at
    desktop and phone widths. Test radio arrow keys, switch Space toggling,
    textarea editing, and disabled file selection in the playground.
11. On Interactions, use Enter/Space to open each panel, Tab through its controls,
    Escape to dismiss the automatic popover/modal, and the Close button for the
    manual popover. Verify modal focus containment and return to its opener.
12. Approve a demo request inside the dialog and through the ordinary-page
    confirmation. Both should show the success message. Canceling/closing should
    perform no POST. Repeat with JavaScript disabled.
13. On Dashboard, change sort order, use page 2/3, change the page size, and filter
    the results. Verify link URLs retain filters and pagination stays in bounds.

14. On Inline editing, save a fictional name, refresh, reopen the editor, change
    it and cancel. Check that the saved value remains. Clear the name, use Test
    server validation, follow the error link, then correct and save. Restore the
    sample, then check another browser session has independent values.

15. On Motion & surfaces, compare Follow system and Off, then enable OS reduced
    motion. Watch entrances finish, keyboard-focus the gradient card link, and
    open/close overlays. Check that skeleton decoration stops after three cycles.
    Repeat in dark mode, forced colors, print preview and a narrow viewport.

16. On Dashboard, toggle each sortable header twice. Verify the arrow, order and
    active sort agree. Use status tabs with a search query and move between pages.
    Select two sample projects and review; repeat with no selection.
17. On Patterns, navigate through all three sections and walkthrough steps with
    the keyboard. Check current markers, unavailable future steps, grouped details,
    table loading/empty content, narrow layouts and screen-reader labels.

18. On Form wizard, submit blank/invalid details through Test server validation,
    follow the error links, and correct the fields. Continue, change preferences,
    go Back, edit the details, and continue again. Verify preferences survive.
19. On each wizard step, use Enter in a text control to check forward submission;
    use Back with incomplete values to confirm native constraints do not block it.
    Complete the review, inspect escaped punctuation, then start a fresh brief.

20. On Catalog, search `field`, filter Forms, page through all results and try an
    empty search result. Open a detail, test the native control, then compare its
    source with the downloaded file. Open both layout previews separately.
21. Test catalog search, code-block scrolling, breadcrumbs, pagination and return
    links at phone width and by keyboard. Toggle the theme from a detail page and
    verify the same entry remains open. Check that source is selectable without JS.
22. Open `/catalog/icon`. Compare all 24 drawings in light/dark modes and at larger
    text sizes. Tab through action examples: only controls should take focus,
    and the icon-only link should announce "Browse all components". Check forced
    colors, print, zoom and long control labels; decorative SVGs should not be
    announced separately. The static contact sheet is not browser evidence.

Automated HTTP tests cover page routes, playground states, escaping, form
validation/body limits, filters, cross-origin rejection, and safe theme redirects:

```sh
go test -race ./examples/site
```

The current sandbox cannot bind local TCP sockets.
The demos must be launched from a normal local terminal for browser review.
A single review binary may be retained at `tests/results/demos/tendui-demo`;
run it with `-open` to avoid recompilation. The next demo build can replace it,
and it can be removed after review.

After a rebuild, stop your running demo with Ctrl+C and start that executable
again. Refreshing the browser alone cannot update CSS embedded in an older
running binary. The build does not terminate your existing process.
