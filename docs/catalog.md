# Component catalog

The demo executable now contains a server-rendered catalog at `/catalog` with
all 35 current public templ components. Search by name, category or usage notes;
filter by category; use ordinary pagination links. Search text is bounded to 120
Unicode characters, category choices are validated, and page numbers are clamped.
No search index service, browser runtime or external request is involved.

Each component has:

- A short description, props type name and composition guidance.
- A live preview, plus a standalone preview route.
- The complete templ recipe as selectable source text.
- A `.templ` download and a link to a fuller working demo.

Layout and AppShell previews open separately to avoid nesting document/main
landmarks inside the catalog. Other previews render inline. Controls retain
native behavior; individual input previews do not imply stored values or form
submission handlers. The Form/Button recipes have working native search actions,
ThemeToggle uses the existing demo handler, and file selection remains local.

## One recipe source

`examples/basic/catalog_samples.templ` is compiled with the demo and embedded as
text. Marker-delimited blocks supply the source shown and downloaded by the
catalog. This keeps preview markup and copied code together. Each download has
`package examples`, imports the public components package, and contains its own
recipe function plus any necessary local helper.

Change the package declaration to match your application. Adapt demo destinations
and IDs; do not assume `/catalog`, `/dashboard` or the theme endpoint exist in
another app. Consumers using the library's generated Go do not need the templ
CLI. Applications editing templ source do need a matching Go templ generator.
No Node/npm step is introduced.

`/catalog/{component}/example.templ` returns a text attachment, not executable
HTML. Only registered component identifiers are accepted; requests never map
user input to filesystem paths. GET and HEAD are supported; mutations are rejected.
The existing script-blocking CSP remains in effect. Catalog-only CSS is embedded
and served by the example app at `/demo.css`, outside the library asset payload.

## Maintaining coverage

For a new public templ component:

1. Add a self-contained marked `ExampleName()` recipe to `catalog_samples.templ`.
2. Register its name, category, description, usage notes, demo and function in
   `catalog.go`. Do not mutate the registry while requests are running.
3. Format and generate templ Go from the repository root.
4. Run the catalog and full project checks. `TestCatalogCoversPublicTemplComponents`
   compares exported templ declarations with the registry and fails on omissions.

HTTP tests cover every detail, preview and download route, parallel rendering,
HEAD semantics, document landmarks, script policy, escaping, safe attachments,
search/pagination boundaries and registered theme-return paths. Source coverage
is specifically for public templ components, not every exported Go type or helper.

The developer recipe check generates and compiles all 35 downloads as separate Go
packages using a temporary module with a local replacement of GothUI. That proves
the copied recipes compile independently; it is not a published-module install
check. The temporary fixture packages are removed after validation, with concise
compile logs retained. Browser layout, keyboard and assistive-technology checks
remain separate release gates outside the restricted sandbox.
