# Verification across Go versions and operating systems

`.github/workflows/verify.yml` prepares native verification on three GitHub-hosted
runner images and two explicitly selected Go versions:

| Runner | Go versions | Purpose |
| --- | --- | --- |
| `ubuntu-24.04` | 1.26.0, 1.27.3 | Linux native checks and ELF interpreter verification |
| `macos-15` | 1.26.0, 1.27.3 | macOS native checks |
| `windows-2025` | 1.26.0, 1.27.3 | Windows native checks and `.exe` execution |

Go 1.26.0 is the declared minimum compatibility target; Go 1.27.3 is the measured
local toolchain. `GOTOOLCHAIN=local` prevents automatic upgrades from concealing a
failure with the selected compiler. Review the matrix when updating the supported
Go policy. Runner labels were checked against the official
[runner image inventory](https://github.com/actions/runner-images).

## What each job does

1. Checks out the triggering source with complete version history and without
   persisting repository credentials.
2. Selects Go and verifies downloaded module checksums.
3. Runs `gothui-dev check`: formatting, dependency boundaries, CSS size, public
   API comparison, race tests and vet.
4. Installs the pinned templ generator, compares generated source, and compiles
   every catalog download independently.
5. Installs the clean committed source through a local versioned module proxy
   into an isolated consumer, tests component rendering and embedded CSS without
   `replace` or templ generation, and verifies module checksums. This checks
   packaging, not installation from a published repository tag.
6. Builds a CGO-disabled native demo, checks its size and embedded build version,
   and executes `-version`. Linux also rejects an ELF interpreter.
7. Retains logs and JSON reports for seven days, even when an earlier step fails.

Jobs have a 20-minute timeout and independent matrix results. A newer run for the
same workflow/ref cancels its predecessor. Triggers are pull requests, pushes to
main, and manual dispatch. Repository permissions are limited to reading contents;
the workflow creates no releases, tags, commits, deployments or mirror pushes.
It uploads reports only, not six duplicate demo binaries. Module and build
caching uses setup-go's normal cache support.

The official Actions are pinned to commit IDs, verified on November 6, 2026:

| Action release | Pinned source |
| --- | --- |
| checkout v7.0.1 | [3d3c42e5aac5](https://github.com/actions/checkout/commit/3d3c42e5aac5ba805825da76410c181273ba90b1) |
| setup-go v7.0.0 | [b7ad1dad31e0](https://github.com/actions/setup-go/commit/b7ad1dad31e06c5925ef5d2fc7ad053ef454303e) |
| upload-artifact v7.0.2 | [cf430e030ddb](https://github.com/actions/upload-artifact/commit/cf430e030ddbb5b0abf93d22962f4752f3646cd9) |

## Known-vulnerability job

A separate `security` job runs on Ubuntu with Go 1.27.3 and pinned
`golang.org/x/vuln/cmd/govulncheck@v1.9.0`. It scans reachable source code in all
packages, including the example and developer tool, against the live Go
vulnerability database. It records the exact source revision, compiler/target
and scanner output. Scan errors and reachable findings fail the job; evidence
is retained for seven days through an `always()` artifact step. No suppressions,
continue-on-error or fabricated successful report is used.

The enforcing scan uses **text mode**. The scanner's
[documented JSON/SARIF/OpenVEX modes](https://pkg.go.dev/golang.org/x/vuln@v1.9.0/cmd/govulncheck#hdr-Exit_codes)
can return success even when findings exist. A local synthetic conversion fixture
with one symbol-level finding verified exit 3 for text and exit 0 for JSON using
the pinned scanner. That validates exit handling only, not the absence of real
vulnerabilities. `go run` may wrap a nonzero tool exit as exit 1; either fails CI.

The initial direct scan failed while fetching the database. A later download of
the complete official `vulndb.zip` snapshot succeeded through the available HTTPS
download path. The pinned scanner completed locally against that snapshot using
its supported file protocol: no known vulnerabilities found in all seven root
packages, templ v0.3.1100 and the Go 1.27.3 standard library on Linux amd64.
The snapshot's database modification time is 2026-11-05T23:49:57Z. Its SHA-256 is
`0a54b7eb56bf031dc214fc4f29ca6db71073ab324595c3de9485f97a334e9e22`.
All 4,789 indexed OSV reports were present; no subset or suppressions were used.
This is a dated snapshot result, not an evergreen security guarantee. The
remote security job independently reads the live database. Inspect all seven
results at the exact published SHA. Source analysis also has
limits around reflection, unsafe calls and other build configurations; see the
[Go vulnerability-management guide](https://go.dev/doc/security/vuln/).

Go 1.26.0 remains a minimum-language/toolchain compatibility test, not a deployment
recommendation. Use patched compilers for deployed binaries. As verified on
November 6, 2026, the official
[release history](https://go.dev/doc/devel/release) lists Go 1.26.10 and 1.27.3 as
the patched releases on those lines. The local demo uses Go 1.27.3. Recheck this
policy when preparing a release; it is not an evergreen claim about the latest
version.

These are CI orchestration dependencies and a Go development-time scanner,
not GothUI runtime or frontend assets.
The application remains a Go binary with embedded CSS and no bundled JavaScript.

## Native build details

On Windows, the toolchain path, candidate executable and retained demo use `.exe`.
The builder verifies the candidate before replacing the retained file. It does
not delete a working executable first. Renaming on Unix provides atomic replacement;
Windows replacement depends on filesystem/handle behavior and can fail if another
process holds the executable. Stop an old demo before replacing it on Windows.

The build-helper integration test creates an isolated module in a path containing
spaces, builds and executes a clean candidate, then replaces it with a changed
source build. It verifies the version change and absence of leftover candidate
files. Running that test on Linux does not prove its Windows result.

`.gitattributes` keeps text checkout line endings consistent for generated-source
comparison. The demo's `-open` dispatches to `xdg-open` on Linux, `open` on macOS,
or the Windows URL handler through `rundll32.exe`. URLs are process arguments,
not interpolated into a shell. CI's version smoke check does not open a browser;
desktop browser launch remains a separate manual check.

## Current evidence and publication boundary

The workflow's matrix, commit pins, read permission and retention settings were
inspected locally. Remote results are available in
[Verify GothUI](https://github.com/Tend-Stack/GothUI/actions/workflows/verify.yml);
require all seven jobs to pass for the intended release SHA. Local native
race/vet checks pass with Go 1.26.0 and Go 1.27.3. The minimum compiler was
downloaded temporarily from Go's official distribution and its published SHA-256
and archive size verified before execution. The API baseline, versioned local
consumer installation, native demo build and version smoke also pass on 1.26.0.
Its read-only cache warning exposed a developer-tool bug: parsed command output
included stderr. The helper now forwards successful diagnostics separately from
stdout and preserves failure diagnostics. Regression tests pass on both compilers.
The build-helper test executes and replaces a
native Linux executable successfully. The demo also cross-compiles for Linux
arm64, macOS arm64 and Windows amd64, with target build metadata and ELF/Mach-O/PE
architecture headers checked. Those disposable foreign executables were removed.
Evidence is retained under `tests/results/validation/`, `tests/results/platform/`,
`tests/results/compatibility-go126/` and `tests/results/security/`.
Cross-compilation cannot prove foreign execution, race behavior or browser UI.
Non-Linux native execution remains unverified. The temporary minimum compiler
and its disposable demo are removed after validation; the retained review binary
continues to use Go 1.27.3. A Firefox headless attempt
with an isolated temporary profile also failed under the sandbox namespace
restriction; its profile was removed. Chromium/browser results remain unverified.

history, GitHub rejected its private author email with GH007. With operator
approval, 26 commits' author/committer emails were corrected to the GitHub no-reply
identity. Names, dates, messages and every file tree were preserved; the descendant
commit already using no-reply only changed its parent reference. The original
operator and retained until an explicit cleanup decision after verified delivery.
history was accepted by GitHub. GitHub remains private; its email-privacy setting
was not disabled. This checkout uses no-reply for future commits. No release tag,

After publication, inspect all six compatibility results plus the security job
at the exact source SHA before calling the matrix verified. Release installation, optional HTMX execution and
browser/accessibility checks remain independent acceptance gates.
