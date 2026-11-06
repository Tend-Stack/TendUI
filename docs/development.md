# Developing GothUI

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
| `go run ./cmd/gothui-dev check` | Go formatting, production dependency boundaries, CSS size, all race tests and vet |
| `go run ./cmd/gothui-dev generated` | Matching templ version and byte-identical generated Go for every template; rejects orphan output |
| `go run ./cmd/gothui-dev recipes` | Each downloadable catalog recipe generated and compiled in an independent package |
| `go run ./cmd/gothui-dev bench` | Three samples per server-side benchmark, allocation counts and environment metadata |
| `go run ./cmd/gothui-dev build` | Versioned, CGO-disabled demo; size, embedded version and Linux ELF interpreter checks |

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
GothUI replacement. This catches examples that accidentally depend on another
recipe's imports or helpers. It does **not** prove a published `go get` works.
The temporary module is removed on success and ordinary errors. An abrupt kill
can leave a `gothui-generated-*` or `gothui-recipes-*` directory in the OS temporary
directory; confirm its process has stopped before removing it.

## Review executable and versioning

```sh
go run ./cmd/gothui-dev build
./tests/results/demos/gothui-demo -version
./tests/results/demos/gothui-demo -open
```

Builds use `CGO_ENABLED=0`, `-trimpath` and stripped debug symbols. The command
checks the candidate executable before atomically replacing the single retained
review binary. Linux builds must have no ELF interpreter; all builds must record
CGO disabled and print the expected version without starting an HTTP listener.
The report records byte size and SHA-256. An existing demo process is left alone;
restart it to see new code. Native builds are intended here: the version smoke
check executes the binary locally, so this command is not a cross-build helper.

Versions use year-based tags, such as `v2026.1.0`. With an appropriate tag, the
build embeds `git describe --tags --long --dirty`. Before the first release tag,
it uses `2026.1.0-dev-g<commit>`. Untracked files also add `-dirty`. Malformed
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
also need migration guidance; automated API compatibility comparison remains a
separate roadmap item.

Keep only current validation reports, the review executable and evidence for
unresolved failures. The continuing developer owns their replacement or removal
after review. Temporary fixtures are owned by the invoking command; shared Go
caches remain user-owned and are not purged by this tool.
