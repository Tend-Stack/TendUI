# Embedded stylesheet delivery

GothUI reads `static/gothui.css` from its embedded filesystem when an AssetHandler
is created. The returned handler serves that immutable byte slice with a new
reader for each request. Consumers need no runtime file, CDN, npm package or
stylesheet build.

```go
mux.Handle(gothui.AssetPath, gothui.AssetHandler())
```

Mount it directly at `/gothui/static/`. Do not add `http.StripPrefix`; the handler
expects `/gothui/static/gothui.css`, the URL used by Layout. It serves only that
file, rejects methods other than GET/HEAD with 405 and an Allow header, and does
not expose directory listings or the rest of the embedded filesystem.

## Revalidation across builds

The stylesheet URL is stable. Responses use a quoted SHA-256 ETag of the exact
served bytes and `Cache-Control: public, max-age=0, must-revalidate`. There is no
invented Last-Modified timestamp for an embedded asset. Independent handlers from
the same binary identify the same stylesheet.

A cache may keep the CSS but must revalidate before reuse. If-None-Match matching
the current ETag returns 304 with no body; weak validators, validator lists and
`*` follow Go's standard ServeContent handling. A normal GET with a stale validator receives 200
with the current full stylesheet. A changed CSS payload therefore has a changed
validator even when the development release-line constant is unchanged. HEAD
returns metadata without a body; a normal HEAD includes the full content length.

Preserve these headers through a reverse proxy. Adding long freshness or
`immutable` caching to this stable URL changes the contract and can keep old CSS
alive across an upgrade. Applications serving their own fingerprinted asset URLs
must own those URLs, cache policy and corresponding head links.

The example's additional `/demo.css` uses `no-store`; example HTML also uses
`no-store`. Appearance exports are application-owned response paths. Asset
revalidation is not an HTML/HTMX cache policy: applications still decide page and
fragment cacheability, representation validators and Vary headers.

## Byte ranges and preconditions

The handler delegates HTTP conditional and range semantics to `http.ServeContent`.
Bounded, open-ended and suffix byte ranges return 206 and the relevant
Content-Range. A range extending beyond EOF is clamped; one wholly beyond EOF
returns 416 with the current file length. Multiple ranges use multipart/byteranges
with an individual content type/range and the exact bytes in each part.

A matching strong If-Range validator permits the requested partial response.
A stale or weak If-Range validator falls back to the complete 200 representation.
A matching If-None-Match takes precedence over Range and returns 304; an unmatched
If-Match fails with 412 before range processing. No shared reader offset is used,
so simultaneous range requests cannot change each other's returned section.

## Verification scope

`asset_http_test.go` checks the actual exported handler, including the served-byte
hash, revalidation headers, HEAD lengths, empty conditional responses, stale
validators, If-Range/precondition precedence, single/multipart ranges and 32
parallel readers. `gothui_test.go` covers mounting, method restrictions, unknown
assets and directories. These tests use in-process HTTP recorders and the race
detector; they do not claim a live browser, proxy or network test.

Run the focused checks with:

```sh
go test -race .
go vet .
```

No cache-policy or handler implementation change was needed for this coverage.
If a deployed page still uses old styles, inspect the actual response's ETag and
Cache-Control at the browser/proxy boundary and confirm which binary is serving
it. A local source checkout or freshly built executable alone does not prove
that the running process has changed.
