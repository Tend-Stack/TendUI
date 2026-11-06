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
5. Builds a CGO-disabled native demo, checks its size and embedded build version,
   and executes `-version`. Linux also rejects an ELF interpreter.
6. Retains logs and JSON reports for seven days, even when an earlier step fails.

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

These are CI orchestration dependencies, not GothUI runtime or frontend assets.
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

The workflow parses locally and its matrix, commit pins, read permission and
retention settings were inspected. It has **not run on GitHub**. Local native
race/vet checks pass with Go 1.27.3. The build-helper test executes and replaces a
native Linux executable successfully. The demo also cross-compiles for Linux
arm64, macOS arm64 and Windows amd64, with target build metadata and ELF/Mach-O/PE
architecture headers checked. Those disposable foreign executables were removed.
Evidence is retained under `tests/results/validation/` and `tests/results/platform/`.
Cross-compilation cannot prove foreign execution, race behavior or browser UI.
No minimum-version or non-Linux runtime pass is claimed yet.

Shell Git access currently fails DNS resolution for both configured remotes.
The GitHub connector reports the import repository as private, unarchived and
empty; Gitea history/archive status remains unverified. Reconcile Gitea history
before publishing verified commits to both intended repositories. No Gitea runner
or publication automation has been configured by this workflow.

After publication, inspect all six job results and their exact source SHA before
calling the matrix verified. Release installation, optional HTMX execution and
browser/accessibility checks remain independent acceptance gates.
