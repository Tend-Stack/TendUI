# Developing TendUI

Consumers build with Go alone; the checked-in generated components use templ as
their only runtime dependency. Contributors need Git, Go 1.26 or newer, and the
generator matching the pinned runtime when checking or changing templates:

```sh
go install github.com/a-h/templ/cmd/templ@v0.3.1100
```

Run the following from the repository root. The developer command is Go code,
uses the toolchain that compiled it, and requires no frontend tools.

| Command | Evidence |
| --- | --- |
| `go run ./cmd/tendui-dev check` | Go formatting, production dependency boundaries, CSS size, public API baseline, all race tests and vet |
| `go run ./cmd/tendui-dev api` | Compiler-resolved public declarations, candidate snapshot and reviewed-baseline comparison |
| `go run ./cmd/tendui-dev generated` | Matching templ version and byte-identical generated Go for every template; rejects orphan output |
| `go run ./cmd/tendui-dev recipes` | Each downloadable catalog recipe generated and compiled in an independent package |
| `go run ./cmd/tendui-dev consumer` | Clean committed module installed from a local versioned proxy, without `replace` or a templ generator |
| `go run ./cmd/tendui-dev bench` | Three samples per server-side benchmark, allocation counts and environment metadata |
| `go run ./cmd/tendui-dev build` | Versioned, CGO-disabled demo; size, embedded version and Linux ELF interpreter checks |

`TEMPL=/absolute/path/to/templ` overrides generator discovery. Go module/cache
environment settings are inherited. `check` requires a platform/toolchain that
supports Go's race detector, including its C compiler prerequisites. No compiler
or templ executable is needed to run the completed demo binary.

Run commands sequentially against a stable working tree. Logs and JSON reports
replace the previous run under `tests/results/validation/`. Reports include the
Git revision, whether tracked or untracked source work is present, toolchain and
UTC time. A failed command removes its previous success report before starting;
its log may contain partial results. A dirty report identifies a working-tree
measurement, not an exact committed release. These commands are local checks,
not CI execution or publication evidence.

## Templates and recipe validation

After changing templates, run `templ fmt .` and `templ generate`, and commit both
the `.templ` and `_templ.go` files. The `generated` check itself does not overwrite
source: it copies templates and module metadata to a temporary directory,
regenerates there, compares outputs, and removes that directory on return.

The `recipes` command uses the same extractor as catalog downloads. Each marked
recipe is compiled without its neighbors in a temporary module with a local
TendUI replacement. This catches examples that accidentally depend on another
recipe's imports or helpers. It does **not** prove a published `go get` works.
The temporary module is removed on success and ordinary errors. An abrupt kill
can leave a `tendui-generated-*` or `tendui-recipes-*` directory in the OS temporary
directory; confirm its process has stopped before removing it.

## Versioned consumer packaging

Run `go mod download` first to cache dependencies, then run `tendui-dev consumer`
from a clean commit. It archives that commit into a temporary file-based module
proxy under the declared release version, without creating a tag or publishing.
A separate consumer uses `go get` for the root, components and render packages,
then runs a rendering/embedded-CSS test and `go mod verify`. The consumer has no
local `replace` and never invokes templ. The check verifies escaped content,
the layout stylesheet link, asset delivery and the absence of a script element.

Dependency archives come from the existing Go download cache and are verified
against the committed checksums. Installation uses an isolated temporary module
cache, so the unpublished fixture cannot enter the shared cache. Both the proxy
and cache are removed on success and ordinary failures; an abrupt kill can leave
a `tendui-consumer-*` directory. The invoking developer owns cleanup after
confirming the process has stopped. The report records archive hash, module sum,
source revision and cleanup status. Dirty sources are rejected without changes.
The archive preserves committed file bytes and omits Git's explicit directory
entries, matching [Go's module archive behavior](https://go.dev/ref/mod#zip-files).
Go validates the resulting archive during installation; `go mod verify` checks
both its download and extracted contents.

This proves versioned packaging and consumption locally. It does **not** prove
repository access, a published tag, public checksum-server acceptance or a live
published installation. Those remain separate release gates.

## Review executable and versioning

```sh
go run ./cmd/tendui-dev build
./tests/results/demos/tendui-demo -version
./tests/results/demos/tendui-demo -open
```

On Windows the retained executable is `tests/results/demos/tendui-demo.exe`.

Builds use `CGO_ENABLED=0`, `-trimpath` and stripped debug symbols. The command
checks the candidate executable before replacing the single retained review
binary (atomic rename on Unix). Windows replacement can fail while an executable
is in use; the tool does not delete the old binary first. Linux builds must have
no ELF interpreter; all builds must record
CGO disabled and print the expected version without starting an HTTP listener.
The report records byte size and SHA-256. An existing demo process is left alone;
restart it to see new code. Native builds are intended here: the version smoke
check executes the binary locally, so this command is not a cross-build helper.

Versions use year-based tags, such as `v2026.1.0`. With an appropriate tag, the
build embeds `git describe --tags --long --dirty`. Before the first release tag,
it uses `<Version>-dev-g<commit>` (for example, `2026.2.0-preview.1-dev-g<commit>`). Untracked files also add `-dirty`. Malformed
release identifiers are rejected. The command creates no tags or commits and
publishes nothing.

## Budgets and remaining release gates

The current review ceilings are 64 KiB for uncompressed library CSS and 16 MiB
for the stripped demo. These ceilings catch accidental growth; they are not
browser performance targets. Budget changes require an explained tradeoff and
fresh measurements. See [performance evidence](performance.md).

Before a release, validate the supported Go versions and operating systems,
browser/keyboard/accessibility behavior, optional HTMX integration in a consumer
fixture, clean published installation and matching repository histories. Passing
local HTTP tests does not establish any of those results. Public API changes
also need migration guidance. The [API review guard](api-compatibility.md) detects
declaration drift; it does not prove behavioral or visual compatibility.

The [CI guide](ci.md) describes the six-job native verification matrix, local
cross-compilation evidence and how to check remote results at the release SHA.

Keep only current validation reports, the review executable and evidence for
unresolved failures. The continuing developer owns their replacement or removal
after review. Temporary fixtures are owned by the invoking command; shared Go
caches remain user-owned and are not purged by this tool.


## Known-vulnerability scan

From a network-enabled development environment, run the pinned Go tool:

```sh
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 -show verbose ./...
```

This development command does not add a TendUI runtime dependency or alter its
module requirements. It needs the Go module proxy and vulnerability database.
Use text mode for a failing gate: machine-readable output modes can return zero
with findings. Treat a database/network error as an incomplete check, never a
clean result. The separate prepared CI security job preserves logs and exact
source/toolchain metadata; [the CI guide](ci.md) records its scope and evidence.
The ordinary `tendui-dev check` remains usable without a vulnerability database.

When the scanner cannot reach HTTPS directly, Go also publishes a
[complete database snapshot and file-protocol instructions](https://go.dev/doc/security/vuln/database).
Download the official archive, validate its paths before extraction, and pass the
extracted database URL with `govulncheck -db file:///absolute/path -show verbose ./...`.
Record the source URL, archive SHA-256, database modification time, exact project
revision, compiler/target, scanner version and exit status. A dated snapshot scan
does not establish later database results or other build targets. Retain one
compressed snapshot and concise results for review; remove its extracted files
and task build cache. The local evidence and limitations are in [the CI guide](ci.md).
