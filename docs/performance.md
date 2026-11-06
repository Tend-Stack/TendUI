# Performance evidence

GothUI benchmarks separate rendering from HTTP response construction. They do not
measure browser paint, CSS animation smoothness, network latency, database work,
concurrent request throughput or application startup.

Reproduce the local suite with:

```sh
go run ./cmd/gothui-dev bench
```

The command runs three 200 ms samples per case with `GOMAXPROCS=2` and allocation
reporting. Logs and revision/dirty-state metadata are in
`tests/results/validation/benchmarks.{log,json}`. Use a quiet machine and compare
equivalent toolchains and inputs before drawing conclusions about changes.
This short baseline is descriptive; there is no timing-based pass threshold.

## Initial baseline

Measured on November 5, 2026 (America/New_York), Linux amd64, AMD Ryzen AI Max+ 395,
Go 1.27.3. Source: the developer-tooling working tree based on commit `72c60f9`;
the report correctly records it as dirty. Values below are medians of three
samples, rounded for readability. Allocation bytes are rounded where samples
differed slightly.

| Case | Time per operation | Bytes allocated | Allocations |
| --- | ---: | ---: | ---: |
| CSS GET | 4.81 µs | 66,717 | 18 |
| CSS conditional GET (304) | 0.574 µs | 1,072 | 12 |
| Card render | 0.161 µs | 216 | 5 |
| Table with 10 rows | 5.29 µs | 5,233 | 201 |
| Table with 100 rows | 45.2 µs | 46,046 | 1,821 |
| Full document with 100-row table | 47.4 µs | 47,393 | 1,865 |
| Dashboard GET | 34.1 µs | 47,589 | 704 |
| Catalog GET | 37.4 µs | 50,731 | 494 |
| Field catalog detail GET | 13.8 µs | 16,135 | 214 |
| Wizard GET | 16.1 µs | 29,623 | 311 |
| Editing GET | 12.9 µs | 16,769 | 234 |

Component cases render preconstructed components to `io.Discard`; they include
templ rendering and validation but exclude fixture construction. Handler cases
include routing, headers, page composition and an `httptest.ResponseRecorder`,
which allocates memory to capture the response. Requests are constructed outside
the measured loop. CSS GET allocation includes captured body buffers, so it is
not the resident cost of embedding the stylesheet. The suite is sequential.

There is visible timing noise: the editing samples ranged from 11.5 to 27.0 µs.
Do not interpret small differences as improvements. Future optimization should
profile representative workloads, preserve escaping and component semantics,
and compare multiple runs. Table allocations grow with row count; that is a
measured candidate for profiling, not justification for bypassing safety checks.

## Size checks

The embedded library CSS is currently 32,260 bytes, below its 65,536-byte review
ceiling. Catalog-only styles remain in the example app. The stripped demo has a
16,777,216-byte ceiling; `go run ./cmd/gothui-dev build` records its exact size,
hash, version and toolchain in `tests/results/validation/build.json`.
The first measured Linux amd64 build was 8,949,920 bytes (8.54 MiB), with CGO
disabled and no ELF interpreter. Size varies with platform and Go toolchain.

These measurements are not comparisons with Tailwind or JavaScript frameworks.
GothUI's single-binary delivery and native interactions are architectural
properties; competitive performance and accessibility claims need relevant
browser and application evidence that this suite does not supply.
