# Reading content with Prose

`Prose` gives documentation, help text and release notes a bounded reading width
and consistent text spacing. It styles native HTML supplied by your application.
It emits a plain div, adds no landmark or heading, and introduces no Markdown
parser, sanitizer, client script, font download or asset build.

```templ
templ helpArticle(title, summary string) {
    <section aria-labelledby="help-title">
        @components.Prose(components.ProseProps{
            ID: "help-body", Measure: components.MeasureStandard,
        }) {
            <h2 id="help-title">{ title }</h2>
            <p>{ summary }</p>
            <h3>Before you begin</h3>
            <ul><li>Choose a project.</li><li>Review its settings.</li></ul>
            <div data-prose="off">
                @components.Alert(components.AlertProps{
                    Title: "Need a hand?",
                    Message: "Contact your project administrator.",
                    Tone: components.ToneInfo,
                })
            </div>
        }
    </section>
}
```

Use ordinary templ interpolation for text. Prose preserves child components as
rendered; it does not make untrusted `templ.Raw` or application-generated HTML
safe. If an application accepts Markdown or rich HTML, its parser and sanitizing
policy remain application concerns before content reaches this wrapper.

## Measure and typography

`ProseProps.Measure` accepts `MeasureNarrow` (55ch), `MeasureStandard` (65ch) and
`MeasureWide` (80ch). Empty or unknown values normalize to Standard. The wrapper
fills available inline space up to that maximum; it never sets a minimum reading
width. `ch` depends on the chosen font's character metrics and does not promise
an exact number of characters per line. Text direction and language inherit
normally. The optional ID is escaped; applications own uniqueness.

The appearance font and root text-size settings continue to apply. Prose uses a
1.75 line height for paragraphs, list items, definitions and quotes, with extra
space before headings. Smaller heading levels remain at least the root text
size within Prose. Nested lists, definition pairs and captions get consistent
spacing; inline code can wrap. Preformatted code keeps its whitespace and scrolls
using the existing native CSS overflow behavior. Code inside a pre does not
compound the outer font-size reduction.

Five general typography tokens now describe the existing global defaults,
without changing their values:

| Token | Default |
| --- | --- |
| `--gothui-leading-body` | `1.6` |
| `--gothui-leading-heading` | `1.2` |
| `--gothui-heading-1` | `clamp(2.25rem, 5vw, 4rem)` |
| `--gothui-heading-2` | `clamp(1.65rem, 3vw, 2.4rem)` |
| `--gothui-heading-3` | `1.35rem` |

Reading-specific tokens are `--gothui-prose-narrow`, `--gothui-prose-standard`,
`--gothui-prose-wide`, `--gothui-prose-leading`, `--gothui-prose-space` (1.25rem)
and `--gothui-prose-heading-gap` (2.25rem). Override them in a stylesheet loaded
after GothUI. Reading rhythm intentionally does not shrink with compact control
density; use the dedicated tokens if your document needs a different rhythm.
Prose adds no color palette and follows the existing light/dark/print tokens.

## Embedded components and content structure

Wrap a widget in `data-prose="off"` to exclude it and its subtree from the prose
descendant typography rules. The rest of GothUI's component CSS still applies,
and the wrapper remains within the reading column. It does not create a full
CSS reset or let content escape that width. Put a new Prose outside an opted-out
subtree if you want its reading styles again. First/last direct children still
receive outer margin trimming at the reading-container boundary.

Use the appropriate heading level for your page; size is not a reason to skip
levels. Keep lists and tables in their valid HTML structure. Prose does not
manufacture an article/card: place it inside a semantic section, main area or
application-owned article as appropriate. Use Stack with center alignment if
you want a narrower reading column centered in a wider area; Prose itself stays
at the normal inline start.

For a wide table, use DataTable or a `.table-scroll` wrapper. For a horizontally
scrolling code example, provide a useful accessible name and native focusability
so keyboard users can reach it, as the demo does. Neither text clipping nor line
clamping is used. Print rules, motion preferences and ordinary link navigation
remain those of the surrounding application.

## Demo and evidence

Open `/patterns?section=reading`. Its native GET form changes reading measure;
the Appearance page changes fonts, text size and theme. The example covers
headings, inline emphasis, nested lists, a quote, definitions, a labeled scrollable
code block, a caption and an opted-out alert. `/catalog/prose` provides a compiled,
downloadable standalone recipe.

Go tests cover normalized measures, shipped style mappings, ID escaping, unchanged
child markup, render errors/cancellation and widget-exclusion selector guards.
HTTP tests cover the real form route, selected measures, malformed query fallback,
semantic content and HEAD behavior. These do not measure typography in a browser.
Review each measure with sans/serif/mono, large text, compact density, dark/light,
RTL, narrow viewports, zoom and print before claiming visual compatibility.
