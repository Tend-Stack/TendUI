# Accessibility and print review

TendUI supplies semantic HTML, native controls, visible keyboard focus and
appearance tokens. An accessible application also depends on its content,
headings, labels, navigation, validation flow and any custom CSS or optional
HTMX behavior. Passing the checks below is not a WCAG conformance claim.

## Palette checks

The Go tests read the embedded stylesheet's literal colors and calculate sRGB
relative luminance. The screen matrix covers light/dark page, card and subtle
surfaces, including tinted cards and striped rows:

- Body text, muted text and the four feedback tones: at least 4.5:1.
- All four presets' links, including hover colors: at least 4.5:1.
- Control-border color against all three surfaces: at least 3:1.
- Filled-button labels against normal/hover backgrounds: at least 4.5:1.
- System dark and explicit dark declarations must match, including token aliases.

The chosen thresholds follow W3C's explanations of
[text contrast](https://www.w3.org/WAI/WCAG22/Understanding/contrast-minimum.html)
and [non-text contrast](https://www.w3.org/WAI/WCAG22/Understanding/non-text-contrast.html).
The tests use unrounded values. Inactive controls are excluded from these
thresholds; the current disabled opacity remains unchanged. Focus rings use the
primary accent, whose contrast against these surfaces is covered by the link
checks. Their visibility, clipping, area and placement still need browser review.

For the current default tokens, the lowest tested screen text ratio is 5.42:1
(light muted text on subtle background); the lowest tested border ratio is
3.21:1 (light border on subtle background). Dark minima are 6.36:1 for text and
3.40:1 for borders. These are source-color calculations, not pixel measurements.
The checks do not model arbitrary consumer overrides, inherited opacity,
background images, transparency, antialiasing or browser forced-color replacement.

Run the palette checks with:

```sh
go test ./components -run 'Test(ScreenPalette|SystemDark|PrintPalette|PresetText)' -v
```

They also run through `go run ./cmd/tendui-dev check` with the full race suite.

## Print behavior

Print uses a white canvas, black body text, dark neutral labels and visible gray
borders regardless of the screen theme. Feedback retains its text label and
border structure while using a neutral print color. The root print palette uses
`!important` deliberately: the explicit and system-dark selectors have higher
specificity than a plain `:root`, and consumer appearance exports load later.
A normal-priority print root would leave parts of a dark palette active.

Buttons and button-like links print with dark text, a transparent background and
a border. The current-step marker also retains dark text and a stronger outline
without a fill. Switches use their native checkbox appearance so checked state does
not rely on a colored CSS thumb. This avoids requiring background graphics just
to read a label or identify a checked control. Browsers and users may adjust
printed colors; see [MDN's print-color-adjust reference](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/Properties/print-color-adjust).
TendUI does not force `print-color-adjust: exact` or override print-dialog choices.

The print rules are consolidated at the end of the stylesheet. Animation and
hover displacement stop, shadows disappear, the application shell becomes a
normal block without a viewport-height minimum, and tables leave their scroll
containers. Sidebars, skip links, theme toggles, native overlays and pagination
remain omitted as before. Main content, ordinary navigation links, table captions
and form controls remain available. Printing does not fetch hidden pages,
expand all disclosures or include closed overlays; applications should offer a
separate complete document when that is needed.

Custom print styling can override TendUI's important print palette with matching
important declarations in a later print stylesheet. Review the result with both
background-graphics settings. The tests guard source rules and palette ratios;
the v2026.1.0 dashboard PDF was also reviewed for pagination, readable controls
and background omission. This does not establish every page/theme combination.

## Consumer review matrix

Use the same built demo, record its `-version`, browser/version and operating
system, and retain only the screenshots/PDFs needed to explain a finding.

| Review | Suggested demo and expected observation |
| --- | --- |
| Keyboard | `/forms` and `/patterns?section=composition`: visible focus, source-order navigation, labels and native submit/reset controls |
| Error recovery | `/wizard` and `/uploads`: linked summaries lead to the relevant field; validation text remains available |
| Zoom and narrow views | `/catalog`, `/dashboard`, `/patterns?section=composition`: content remains usable at narrow widths and 200%/400% zoom; long rows stay reachable |
| RTL | A consumer page with `dir="rtl"` on its content wrapper: start/end placement follows direction without reversing DOM order |
| Forced colors | `/playground`: native checkbox/switch state, control boundaries, focus and inline icons remain distinguishable |
| Reduced motion | `/motion`: content remains present; entrances, lift and skeleton animation stop under the OS preference and explicit Off |
| Print | `/playground`, `/dashboard`, `/workspace`, `/motion`: repeat light/dark/system with background graphics on/off; inspect labels, switches, surfaces, clipped tables and extra pages |
| Assistive technology | Check names, landmarks, status/error announcements, table captions/headers and native overlay focus with the target screen reader |

Modern native dialog, popover and grouped-details support varies by browser;
retain the documented ordinary-page/inline fallbacks. A successful HTTP response
or render test proves neither keyboard behavior nor screen-reader output. The
local sandbox cannot launch browsers or bind a TCP socket, so the release suite
runs real Chromium, Firefox and WebKit engines in GitHub Actions. It passed 53
checks, including optional HTMX execution; one unsupported WebKit forced-color
emulation check is skipped. Automated axe rules cover all 37 component pages in
both themes. Keyboard, native overlays, narrow layouts, reduced motion and print
have focused coverage. See [release verification](release-verification.md).
The broader consumer matrix above, especially screen-reader output, zoom and
physical devices, still needs review for each assembled application.
