# Motion and expressive surfaces

TendUI uses native CSS for short, optional motion. There is no animation runtime,
DOM observer, client state manager, or asset build step. These styles ship in the
same embedded CSS as the rest of the library.

## Choose a motion preference

```go
appearance := components.Appearance{Motion: components.MotionSystem}
// Or components.MotionOff to disable TendUI's movement explicitly.
```

Pass that appearance to `LayoutProps`. System mode enables movement only inside
`prefers-reduced-motion: no-preference`; it never overrides a user's reduced-motion
setting. Off sets durations, stagger delays and hover displacement to zero.
The appearance editor saves this setting and includes it in CSS/Go exports.
An OS preference cannot be detected by the Go handler; the CSS media query owns it.

## Entrances and surfaces

```templ
@components.Card(components.CardProps{
    Surface: components.SurfaceGradient,
    Entry: components.EntryRise,
    Delay: components.DelayShort,
    HoverLift: true,
}) {
    <h2>Your next project</h2>
    <p>An opaque surface with a theme-aware gradient border.</p>
    <a href="/projects/new">Create a project</a>
}
```

Surfaces: default, tinted, outlined and gradient. The gradient decorates the border;
text stays on a solid background. A tinted surface uses the existing subtle token.
HoverLift is a two-pixel decoration for pointer hover and keyboard focus within.
It does not make the article interactive: provide a real link/button for actions.
Focus outlines are not animated or removed. On touch-only devices, pointer hover
motion is not applied.

### Elevation

The shared `Elevation` type has four choices: `ElevationFlat`, `ElevationSubtle`,
`ElevationRaised` and `ElevationFloating`. `Appearance.Elevation` sets the page's
card and overlay shadow. Raised remains the default. The Appearance editor and
playground expose all four choices and include them in saved preferences and
CSS/Go downloads. `/motion` compares them and lets you change the page setting.

```templ
@components.Card(components.CardProps{
    Surface: components.SurfaceTinted,
    Elevation: components.ElevationSubtle,
}) {
    <h2>A quieter surface</h2>
    <p>A small shadow separates this card from its background.</p>
}
```

Card's optional `Elevation` overrides the page setting, including a flat page.
Empty or unknown values inherit the page setting. An override affects only that
card; a nested card still follows the page unless it has its own override.
Outlined cards always remain shadow-free. Shadows do not change stacking order,
create an interaction or communicate state by themselves. All surfaces retain
borders, and print and forced-color modes suppress shadows.

CSS consumers can use `data-card-elevation="subtle"` on an article. The shadow
tokens are `--tendui-elevation-flat`, `--tendui-elevation-subtle`,
`--tendui-elevation-raised` and `--tendui-elevation-floating`; light, explicit dark
and system dark provide consistent variants. `--tendui-shadow` is the page alias
used by default cards and overlays. Custom styles should preserve the print and
forced-color overrides. Review depth, clipping and readability in a browser;
source-level tests do not measure the rendered effect.

### Entrance behavior

Entrance effects: `EntryFade`, `EntryRise` (half a rem), and `EntryScale` (98% to
100%). The zero value applies none. Invalid enum values fall back to none.
Entrances run once when inserted, including an HTMX fragment inserted by a
consumer's runtime. They are not scroll-reveal observers. `Reveal(RevealProps)`
provides the same effect around children in a div; use Card's fields when an
extra wrapper would break the layout or HTML content model.

Default entrances last 320 ms. Optional delays are 80, 160 and 240 ms. No final
transform or opacity is retained after the animation; content remains visible
without CSS or when animation is disabled. Native dialogs/popovers animate on
opening. Closing is immediate, preserving native dismissal/focus behavior and
avoiding reliance on discrete exit-transition support.

The token API is `--tendui-duration-fast`, `--tendui-duration-enter`,
`--tendui-duration-loading`, `--tendui-stagger-step`, `--tendui-lift-distance` and
`--tendui-ease-out`. Custom durations should remain short. Custom CSS must continue
to respect reduced motion, including pseudo-elements and any new transforms.

## Loading, empty and metric states

```templ
@components.Skeleton(components.SkeletonProps{
    Label: "Loading projects…", Lines: 3, Avatar: true,
})
@components.EmptyState(components.EmptyStateProps{
    Title: "No projects yet", Message: "Start with a small idea.",
}) {
    <a href="/projects/new">Create a project</a>
}
@components.Stat(components.StatProps{
    Label: "Projects completed", Value: "24", Detail: "This month",
})
```

Skeleton has a visible status label; shapes are decorative and hidden from
assistive technology. It defaults to three lines and clamps the count to eight.
Only the shapes pulse, for three cycles (3.6 seconds by default), then settle.
There is no endless animation. The application must replace the placeholder
when data is available and display a recovery state if loading fails. Never
pretend a request is in progress merely to demonstrate an animation.

EmptyState includes the heading and message; its children supply recovery actions.
It is static, without an automatic alert announcement. Stat uses a description
list to keep a metric label and value associated even without styling.

## Verification and reference

The `/motion` demo includes replay, motion-off, light/dark/preset controls,
interactive overlays, all surfaces and state components. The Appearance page
provides persistence. Automated Go tests cover safe props, bounded output,
preference persistence/exports, HTML semantics, routes and CSS contract guards.
Those guards do not prove browser animation, layout or assistive technology behavior.

Manual release checks remain required: desktop/mobile Chromium, Firefox and
WebKit; reduced motion; forced colors; keyboard focus; print; JavaScript disabled;
closed overlays remaining hidden; no content lost after animation completion.
The installed headless Chromium cannot launch inside the current sandbox
(`sandbox_host_linux.cc: shutdown: Operation not permitted`), so browser results
are not claimed.

The operator supplied [PureCSS-Animations](https://github.com/sera619/PureCSS-Animations)
as a visual reference, especially its gradient-border and hover examples. This
implementation is independently authored against TendUI's tokens and accessibility
constraints; no source files from that collection are bundled. The repository
also labels its image carousel as not CSS-only, so that example is outside the
current no-JavaScript contract.

Motion policy follows the [W3C reduced-motion technique](https://www.w3.org/WAI/WCAG22/Techniques/css/C39)
and [MDN's media-query reference](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/At-rules/@media/prefers-reduced-motion).
