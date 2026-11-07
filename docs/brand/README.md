# TendUI brand assets

The compact dark/light SVG lockups and demo mascot come from the operator-provided
Tend UI Gopher Brand Kit. These are genuine vector geometry and editable wordmark
text, with system font fallbacks and no external fonts, scripts or linked images.
The demo uses the transparent vector mascot; no illustration PNGs are embedded
in the library stylesheet or component packages.

The website repository prepares the favicon, social previews and downloadable
web assets. See `https://tendui.com/brand/github-social-preview.png` for the
1280 × 640 GitHub social preview image. Uploading a social preview in repository
settings is separate from displaying the README logo.

## README banners

The light and dark README banners use the owner-supplied Tend UI Gopher Brand Kit
as their source, with an AI-assisted layout correction. The mascot and wordmark
form a compact group; **MODERN UI FOR GO** aligns beneath the wordmark, and the
complete mascot has padding above and below. Light and dark versions retain the
library's illustrated identity, distinct from the website repository's
browser-window banner.

[readme-banners.json](readme-banners.json) records their origin, processing,
dimensions, byte sizes and SHA-256 hashes. These documentation-only images are
not embedded in the runtime stylesheet or component packages. The README's
`picture` element selects the matching color scheme.
