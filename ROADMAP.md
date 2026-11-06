# GothUI product roadmap

The objective is an elegant Go-native UI library that ships in a single Go binary,
with the breadth, polish and developer experience expected of modern UI systems.
Pico was the baseline, not the feature ceiling. “First in class” is an aspiration,
not a claim established by a test suite or component count.

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
covering all 32 public templ components with compiled/downloadable recipes.
Twelve demo areas exercise these APIs.

Current evidence is local. Browser review is incomplete and remote publication
is pending. The active broad goal remains unfinished.

## Next development tracks

1. **Composition and information density.** Table contracts, bidirectional sort
   navigation, loading/empty states, native selection review, page tabs, steps
   and disclosures are implemented. Next: richer layout recipes and composable
   application actions with explicit server ownership; prove responsive and
   keyboard behavior in browsers. Avoid decorative controls that suggest
   unsupported behavior.
2. **Forms for real applications.** Field label/help/error bindings, input addons
   and native multi-step form validation are implemented. Next: constrained file
   handling recipes and application-owned mutation/idempotency examples where
   useful. Validate keyboard and assistive-technology behavior. Keep validation,
   storage and request protection separate from the library.
3. **Visual language.** Coherent space/type/elevation tokens, additional surface
   recipes, an icon strategy without external requests, responsive page patterns,
   and careful interaction states. Motion must remain optional and bounded.
4. **Developer experience.** Native catalog search, category filters, preview
   pages, copyable/downloadable recipes and public templ coverage checks are
   implemented. Next: API compatibility checks, package documentation and
   versioned migration guidance. Keep demos runnable from one binary.
5. **Evidence and release.** Browser matrix and visual regression baselines;
   keyboard/screen-reader/contrast checks; benchmarks and asset/binary size
   budgets; Go version/build CI; generated-source consistency; exact Gitea/GitHub
   history reconciliation; versioned release documentation and installation proof.

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
- Gitea history is integrated safely, both intended remotes contain the verified
  release, and the release uses the operator's YEAR.RELEASE.FIX versioning.

Keep bounded test evidence and one review executable. Do not erase older history
or publish claims of browser compatibility, accessibility conformance, feature
parity or superiority without the matching evidence.
