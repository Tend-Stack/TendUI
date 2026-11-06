# Composing application layouts

Use `Stack` for vertical groups, `Cluster` for wrapping rows, and `Grid` for
responsive columns. They compose with Card, Form, Field, native headings and
ordinary HTML; no client runtime measures or rearranges the page. Explore
`/patterns?section=composition` and the Stack/Cluster catalog recipes.

```templ
templ projectActions() {
    @components.Stack(components.StackProps{Gap: components.SpaceLarge}) {
        <h2>Project settings</h2>
        <p>Review the appearance and component examples.</p>
        @components.Cluster(components.ClusterProps{
            Gap: components.SpaceSmall,
            Justify: components.LayoutJustifyEnd,
        }) {
            <a href="/catalog">Browse components</a>
            @components.Link(components.LinkProps{
                Text: "Appearance", Href: "/settings",
            })
        }
    }
}
```

## Spacing and density

Both props accept `Gap components.Space`. Empty and unknown values normalize to
`SpaceMedium`; choose `SpaceNone` explicitly for no gap. The scale derives from
`--gothui-space`, so the existing root density setting and exported appearance
stylesheet apply automatically.

| Go option | CSS token | Multiple of base space |
| --- | --- | --- |
| `SpaceNone` | Literal zero | 0 |
| `SpaceExtraSmall` | `--gothui-space-xs` | 0.25 |
| `SpaceSmall` | `--gothui-space-sm` | 0.5 |
| `SpaceMedium` | `--gothui-space-md` | 1 |
| `SpaceLarge` | `--gothui-space-lg` | 1.5 |
| `SpaceExtraLarge` | `--gothui-space-xl` | 2 |

Comfortable density uses a 1rem base; compact uses 0.75rem. Customize the scale
at `:root` in your stylesheet loaded after GothUI. Root custom properties derive
at the root: changing only `--gothui-space` on a descendant does not recalculate
root-derived tokens there. For a local adjustment, set the specific scale token
on that container instead. Prefer the typed Gap prop for ordinary layout choices.

Each layout resets **only direct children's block margins** and supplies a gap
between them. Descendant margins and padding remain the child's own concern.
Nested flows select their own gap; a large outer stack does not force large gaps
inside an inner cluster. The legacy `.stack` and `.cluster` CSS utilities retain
their original behavior. The new components use `.layout-stack` and
`.layout-cluster`; existing pages do not change until they opt in.

## Alignment and wrapping

`StackProps.Align` defaults to `LayoutAlignStretch`; it also accepts Start,
Center and End. Those choices position children along the inline axis. An
explicit-width child can still keep its width; stretch does not override native
sizing constraints. Baseline is not a Stack mode and falls back to stretch.

`ClusterProps.Align` defaults to `LayoutAlignCenter`; Start, End, Stretch and
Baseline are also available. It positions children across each row.
`ClusterProps.Justify` defaults to `LayoutJustifyStart` and also accepts Center,
End and Between. Between distributes remaining space **within each line**; it
is not a multi-column grid. An item on its own line remains at the start.

Clusters always wrap, and both components allow items to shrink within their
available inline width. Long ordinary text can wrap. Content with its own
non-wrapping or fixed-size rules, such as a preformatted code block, still needs
an appropriate overflow container. Start/end follow the document's writing
direction. Neither component reverses or sorts children, so source, reading and
keyboard focus order stay the same. No inline styles are emitted.

## Semantics and forms

These components emit a plain `div` with optional ID. They add no landmark, ARIA
role, focus stop or behavior. Use a real section and heading around a meaningful
region, or a form around controls. Put actual `li` children in a `ul`/`ol`;
a layout div cannot replace a list item, table row or table cell. Keep the
application responsible for unique IDs, routing and validation.

Stacks work well for label-wrapped Input components or Field groups. A Cluster
can hold native submit/reset buttons and ordinary navigation links. Put it inside
the Form if the buttons submit that form. It does not make a card or text block
interactive. Use Grid when columns, rather than a row of content-sized items,
should share available space.

## Evidence and review

Tests cover safe normalization, independent nested props, escaped IDs, child
render errors/cancellation, source order and option-to-stylesheet mappings. The
composition example's GET form reaches the searchable catalog; both downloadable
recipes are compiled independently. These checks do not measure browser layout.
Review narrow/wide views, compact/comfortable density, RTL, text zoom, keyboard,
forced colors and print in a real browser before claiming visual compatibility.
