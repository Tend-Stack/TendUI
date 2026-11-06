# Icons that travel with your Go binary

GothUI includes 24 original geometric SVG drawings. The paths compile into the
component package and render inline, so they need no image server, icon font,
sprite download or JavaScript. They use a shared 24-unit grid, two-unit rounded
strokes and `currentColor` to follow surrounding text, button and theme colors.

Open `/catalog/icon` in the demo to explore the set, sizes, action examples and
downloadable recipe. `IconNames()` returns a fresh, sorted list of names for
application galleries; callers cannot mutate the library registry.

## Use words with actions

```templ
@components.Button(components.ButtonProps{
    Text: "Search projects", Type: components.ButtonTypeSubmit,
    IconBefore: components.IconSearch,
})
@components.Link(components.LinkProps{
    Text: "Download report", Href: "/report.csv",
    IconAfter: components.IconDownload, Outline: true,
})
```

`IconBefore` and `IconAfter` are optional on Button and Link. Their icons are
decorative: the existing control text supplies its accessible name. Native form
behavior, disabled state, overlay commands and optional HTMX attributes keep
their existing contracts. No new element is rendered when both icon fields are
zero. Use keyed props literals, as in all examples, so optional additions do not
require changing positional struct initializers.

An icon-only action must name the control explicitly:

```templ
@components.Button(components.ButtonProps{
    IconBefore: components.IconClose, AriaLabel: "Close project preview",
    CommandFor: "project-preview", Command: components.DialogClose,
})
```

The button remains the focus target; its SVG is hidden from assistive technology
and is not focusable. An icon-only action with no nonblank `AriaLabel` returns a
render error before writing markup. Visible words are usually clearer, especially
for unfamiliar actions. An aria label is not a visible tooltip. If overriding a
control that already has visible text, include those words in the accessible name
so speech-control users can refer to the visible label.

## Decorative and meaningful images

```templ
<p>
    @components.Icon(components.IconProps{Name: components.IconCheck})
    All changes saved.
</p>
@components.Icon(components.IconProps{
    Name: components.IconCheck, Label: "All changes saved",
    Size: components.IconSizeLarge,
})
```

Icons are decorative by default (`aria-hidden="true"`). A nonblank `Label` on
Icon gives a standalone SVG `role="img"` and an escaped `aria-label`. Avoid naming
it twice when adjacent text or a containing control already conveys the meaning.
A static icon does not announce a live status change: the application owns live
regions and asynchronous feedback. Symbols and color should not be the only way
to explain important state.

Sizes are `IconSizeSmall` (1em), `IconSizeMedium` (1.25em, the default) and
`IconSizeLarge` (1.75em). They scale with text settings; without CSS the SVG has
intrinsic 24-by-24 dimensions. Icons do not shrink inside flex controls. Unknown
sizes fall back to medium. There is no icon animation or hover transformation.

`IconNone` emits nothing. An unknown nonempty name returns an error instead of
silently showing an unrelated symbol. Values cannot become remote URLs or SVG
markup. As with other component errors, use buffered rendering when sending
HTTP responses; GothUI's `render.HTML` handles errors without sending partial
pages. Directional icons are literal: an application with RTL navigation chooses
the corresponding left/right drawing for its own meaning.

## Available drawings

| Purpose | IconName constants (prefix each with `components.`) |
| --- | --- |
| Direction | `IconArrowLeft`, `IconArrowRight`, `IconArrowUp`, `IconArrowDown`, `IconChevronDown`, `IconChevronRight` |
| Actions | `IconCheck`, `IconClose`, `IconPlus`, `IconMinus`, `IconSearch`, `IconMenu`, `IconDownload`, `IconUpload`, `IconEdit`, `IconExternalLink` |
| Status and theme | `IconInfo`, `IconWarning`, `IconSun`, `IconMoon` |
| Content | `IconUser`, `IconFolder`, `IconFile`, `IconCalendar` |

These drawings are shipped under GothUI's project license. They are authored in
`components/icons.go`; no third-party icon collection was copied or added.

## Evidence and remaining review

Tests parse the rendered SVG namespace, dimensions, path, color and accessibility
attributes; exercise all names, escaping and invalid values; and check icon actions
retain native/HTMX attributes. The catalog recipe compiles independently.

A light/dark contact sheet at 24 and approximately 38 pixels was rendered and
visually inspected using the installed librsvg/Cairo document renderer. It is
retained under `tests/results/icons/`. That verifies the drawings, not CSS layout
or a browser accessibility tree. Browser review remains required for keyboard
focus, zoom, long labels, forced colors, print, RTL usage and screen-reader names.
