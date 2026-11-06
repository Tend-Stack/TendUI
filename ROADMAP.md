# TendUI product roadmap

The objective is an elegant Go-native UI library that ships in a single Go binary,
with the breadth, polish and developer experience expected of modern UI systems.
“First in class” is an aspiration, not a claim established by a test suite
or component count.

## Product contract

- Embedded, locally served CSS; no Node/npm, external fonts, CDN or asset build.
- No bundled client JavaScript. Native HTML/CSS first; optional HTMX contracts
  for consumers who already choose to supply that runtime.
- Generated templ Go included. The component runtime is the one Go dependency;
  assets and HTTP rendering helpers stay standard-library-only.
- Semantic markup, keyboard usability, responsive layouts, light/dark/system
  themes, explicit reduced-motion support, and progressive enhancement.
- Application-owned routing, authentication, persistence and request protection.
- Stable, typed APIs with safe defaults, composable slots, clear examples and
  measured size/rendering cost. No false equivalence between server rendering
  and browser-side behaviors that need a runtime.

## Current foundation

Implemented and covered by Go tests: embedded asset HTTP behavior; forms and
validation feedback; typed appearance settings/exports; application shells;
breadcrumbs and pagination; native dialogs/popovers; full-page/fragment helpers;
inline editing with native fallback; short CSS motion; card surfaces; loading,
empty and metric states; reusable DataTable with native selection review; page
tabs, step indicators and grouped disclosures; composable fields/input addons;
a native three-step form with final revalidation; and a searchable catalog
covering all 37 public templ components with compiled/downloadable recipes.
The Icon component includes 24 original inline SVG drawings and optional
button/link adornments with explicit accessible naming for icon-only actions.
A bounded CSV upload preview exercises native multipart forms and Go validation
without persisting files. Stack and Cluster provide typed layout composition
with a shared density-aware spacing scale. Prose adds reading measures, long-form
HTML typography and widget exclusions. Thirteen demo areas exercise these APIs.

The Go-only developer command now checks runtime dependency boundaries, race/vet,
generated-source reproducibility and independent catalog recipe compilation. It
records rendering/HTTP benchmarks and enforces CSS/demo size budgets, with Git
version stamping and checked atomic replacement of the review executable.
See docs/development.md and docs/performance.md for current evidence and limits.
A compiler-resolved public API baseline now guards 246 declarations/method-set
entries across the three public packages. It flags additions and changes for
review without claiming semantic compatibility; docs/api-compatibility.md defines
the update and migration process.
A six-job GitHub verification matrix is configured for Go 1.26.0/1.27.3 on Linux,
macOS and Windows. Local native Linux checks pass on both Go versions, including
the minimum compiler's API, consumer installation and executable smoke checks.
Foreign-target builds pass; GitHub runs the native matrix for pushed source.
Require its success at the release SHA; docs/ci.md distinguishes compilation
evidence from native results. A separate pinned
vulnerability job uses text-mode failure semantics. A full official database
snapshot scan completed locally with no known findings on Go 1.27.3/Linux amd64;
docs/ci.md records its dated scope and the separate remote security job.
Allocation profiling removed redundant table-cell wrappers and attribute
conversion: the local 100-row benchmark uses 618 allocations instead of 1,821,
with unchanged output across the captured parity fixtures. See docs/performance.md.

Expanded source-color checks cover text, links and control borders across all
built-in surfaces. A consolidated print palette takes priority over dark themes
and retains readable controls without requiring background graphics. See
docs/accessibility.md for measured scope and the consumer review matrix.

Expanded asset HTTP coverage verifies content-derived ETags, stale validators,
HEAD, single/multipart ranges and concurrent readers. docs/assets.md records the
stable-URL cache contract and distinguishes in-process checks from live delivery.

The current preview has a complete documentation website, shared
component catalog, examples and playground, plus passing native, security and
browser verification. See [release verification](docs/release-verification.md)
for historical pre-rename evidence. Current preview details are in
[preview release notes](docs/releases/2026.2.0-preview.1.md).
The development tracks below describe continuing product work beyond this baseline.

## Next development tracks

1. **Composition and information density.** Table contracts, bidirectional sort
   navigation, loading/empty states, native selection review, page tabs, steps,
   disclosures, and Stack/Cluster layout recipes are implemented. Next: composable
   application actions with explicit server ownership; prove responsive and
   keyboard behavior in browsers. Avoid decorative controls that suggest
   unsupported behavior.
2. **Forms for real applications.** Field label/help/error bindings, input addons,
   native multi-step form validation and a bounded CSV upload preview are
   implemented. Next: application-owned mutation/idempotency examples where
   useful. Validate keyboard and assistive-technology behavior. Keep validation,
   storage and request protection separate from the library.
3. **Visual language.** Inline SVG icons with typed names and accessible action
   composition, a shared spacing scale, Prose typography tokens and four shared
   elevation levels with per-card overrides are implemented. Next: additional surface recipes,
   responsive page patterns and careful interaction states. Motion must remain optional and bounded. Verify icons and control
   layouts in browsers as well as the retained static drawing contact sheet.
4. **Developer experience.** Native catalog search, category filters, preview
   pages, copyable/downloadable recipes and public templ coverage checks are
   implemented. A conservative public API change guard and migration process are
   also in place. Next: release-specific migration guidance and broader consumer
   compatibility evidence. Keep demos runnable from one binary.
5. **Evidence and release.** Local benchmarks, asset/binary size budgets and
   generated-source consistency checks are implemented. The browser matrix is implemented. Next: stable visual regression baselines
   and manual screen-reader/device checks;
   passing GitHub CI at the release SHA; versioned release documentation
   and installation proof. Table allocation profiling and
   the first optimization are complete; profile representative application work
   before further rendering changes, retaining semantic/escaping contracts.

## Acceptance gates before calling the broad goal complete

- Every public component has documented semantics, composition rules, state
  coverage and a working demo. Test malformed/hostile values where relevant.
- Native no-JavaScript flows work end to end. Optional HTMX behavior is verified
  in a separate consumer fixture and is never implied by HTTP-only tests.
- Browser evidence covers normal/reduced motion, print, forced colors, mobile,
  keyboard interaction and modern native feature fallbacks. Record exceptions.
- The full suite, race detector and vet pass; templ output is reproducible;
  installation and a consumer build succeed from the intended published version.
- Performance/size results are measured against recorded budgets. All assets
  are local and embedded; no runtime package installation or external fetches.
- GitHub contains the verified release, and the release uses
  YEAR.RELEASE.FIX versioning.

Keep bounded test evidence and one review executable. Do not erase older history
or publish claims of browser compatibility, accessibility conformance, feature
parity or superiority without the matching evidence.
