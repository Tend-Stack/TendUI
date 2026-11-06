# Public API review and migration

GothUI's public Go surface is recorded in `api/public.json`. The initial
development baseline contained 212 declarations and method-set entries across
the root asset package, `components` and `render`. This is a baseline of the
current implementation, not a comparison with a published release or the archived
project. No release has been declared by creating this file.

## Run the guard

```sh
go run ./cmd/gothui-dev api
```

The standard `go run ./cmd/gothui-dev check` includes the same API guard before
the race suite and vet. Ordinary `go test ./...` tests the guard but does not
compare the library with the checked-in baseline; use the developer command for
that comparison.

The guard uses Go compiler export data and the standard library's type importer.
It adds no dependency. It records:

- Exported constants with resolved types and exact values, including `iota`.
- Function signatures, variable types, aliases and generic constraints.
- Type definitions, including struct field order, embedding and tags.
- Value and pointer method sets, including promoted exported methods.
- Local private types reachable through public declarations.

Every difference fails the check, including additions. This deliberately calls
for review instead of labeling an arbitrary change backward compatible. Comments
and function bodies do not alter the snapshot; equivalent constant expressions
resolve to the same value. Parameter-name changes or toolchain formatting changes
can still produce a harmless difference that needs review.

The command writes a candidate to `tests/results/validation/api-current.json`
and differences to `api-changes.json`. A matching run writes `api.json`, with
revision, dirty state, Go toolchain, target OS and architecture. The previous API
reports are removed before a new run, so failure cannot leave an older success
report looking current. The guard never rewrites `api/public.json`.

## Review an intentional change

1. Read the difference and identify affected consumer code. Check zero-value
   behavior, interface satisfaction, serialization and generated templ examples.
2. Preserve existing behavior where possible. Add opt-in names or fields rather
   than silently changing existing semantics. Update documentation and tests for
   the actual behavior, not just the declaration.
3. Write migration guidance with before/after usage when consumers must change.
   Evaluate CSS, markup, accessibility and native/optional HTMX behavior separately.
4. After reviewing the candidate, replace the baseline explicitly and inspect the
   Git diff:

   ```sh
   cp tests/results/validation/api-current.json api/public.json
   git diff -- api/public.json
   go run ./cmd/gothui-dev check
   ```

5. Commit the reviewed baseline together with the implementation, tests and
   migration note. Merely replacing it does not establish compatibility.

The baseline is intentionally small in scope: public library packages only.
The example server, internal helpers, tests and developer command are not
supported consumer APIs. Package discovery rejects new packages outside those
directories until they are added to the explicit public package list in
`cmd/gothui-dev/api.go`. Document their contracts and update the baseline too.

## What a matching result does not prove

This is a conservative declaration-change guard, not a complete Go compatibility
solver. Imported types are referenced by package-qualified name, so changes to
templ or the standard library must be reviewed separately. The comparison runs
for the selected Go build target and build tags; one native run does not verify
all supported platforms or toolchains. Compiler export formatting can differ
between toolchains. The initial run is Linux amd64 with Go 1.27.3; minimum-version
and platform CI remains a release gate.

Matching declarations do not prove identical defaults, errors, validation,
concurrency behavior, HTML structure, CSS tokens, appearance, focus behavior or
HTTP responses. Existing semantic tests, native-flow tests, browser review,
consumer compilation and release installation checks remain necessary. A new
struct field can break positional literals and affect comparability even when
keyed literals still compile; do not describe all field additions as universally
source compatible.

## Consumer practices and version lines

Use keyed props and option literals, as the catalog examples do:

```go
components.Button(components.ButtonProps{
    Text: "Save", Type: components.ButtonTypeSubmit,
})
```

Keep routing, authorization, persistence, CSRF protection and request validation
in the application. Pin module versions and review dependency updates. Recheck
the browser behavior of custom CSS and optional HTMX integrations on upgrades.

Versions follow the operator's `YEAR.RELEASE.FIX` scheme. The Go module's semantic
major is the year: `/v2026` corresponds to tags such as `v2026.1.0`. Fix releases
and later feature releases within the same major should preserve supported
consumer contracts. When a redesign would break them, preserve the old API or
introduce opt-in new names; a feature-number increase is not a Go major bump.

An annual module line such as `/v2027` requires explicit import updates and a
migration guide. Do not rewrite old tags. The first published version still needs
repository reconciliation, release documentation and a clean consumer install;
this local baseline does not substitute for those gates.


## Development addition: Stack and Cluster

The reviewed composition change adds 22 declarations (234 baseline entries in
total): Stack/StackProps, Cluster/ClusterProps, Space with six values,
LayoutAlign with five values, and LayoutJustify with four values. No existing
signature, type, struct field, constant value or method set changes. Existing
consumers do not need a source migration.

The new CSS classes are opt-in. Existing `.stack`, `.cluster`, Grid and Card
behavior remains unchanged. Five new root spacing tokens derive from the
existing density base. When adopting a Stack or Cluster, remove redundant
per-child block margins: the container owns that spacing. Do not mechanically
wrap list/table children in divs; preserve their HTML content model. See
[the layout guide](layout.md) for defaults, nested gaps and semantics. Adding
these names is a reviewed development change, not evidence of a published release
or browser compatibility.
