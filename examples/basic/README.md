# GothUI demos

Four connected pages, served by one Go executable. All pages use the embedded
GothUI stylesheet and generated templ components. No frontend tooling or
JavaScript is required; the response CSP blocks scripts.

From the repository root:

```sh
go run ./examples/basic -open
```

`-open` opens the playground in the default Linux browser using `xdg-open`,
after the listener starts. Without that flag, open the links below manually.
Choose a different port with `-addr 127.0.0.1:8081`. Stop with Ctrl+C.

| Page | Local address | Try it |
| --- | --- | --- |
| Overview | http://127.0.0.1:8080/ | Component introduction and links to each demo |
| Playground | http://127.0.0.1:8080/playground | Themes, button labels/variants, disabled/invalid/read-only states, card header/footer, checked controls |
| Form lab | http://127.0.0.1:8080/forms | Required fields, invalid email, server validation, reset, success receipt |
| Dashboard | http://127.0.0.1:8080/dashboard | Search, status filters, empty results, responsive cards and scrollable tables |

The playground uses GET controls: click **Apply settings** to refresh the preview.
Its URL preserves settings so you can bookmark a test case. The preview button
submits a native form and confirms the click. Its fields demonstrate visual and
native interaction states; the form lab tests actual submission validation.

Theme selection in the playground previews light, dark, or system preferences.
The navigation toggle saves a cookie, removes the preview override, and returns
to the current demo with other filters preserved when the browser sends a
same-origin Referer. Without Referer it returns to the overview. When using
system preference, the server cannot detect the OS theme; the first toggle
uses its saved/default theme as the starting point.

The form lab accepts POSTs but saves no data and sends no messages. **Submit
demo** uses browser validation. **Test server validation** deliberately bypasses
browser checks to exercise HTTP 422 field errors. The server validates the name,
email, role, consent, and message length and limits request bodies to 16 KB.
**Reset edits** restores values from the current document; **Start a fresh test**
clears the form. Successful submissions show an escaped, in-memory receipt.

Dashboard data is fixed sample content. Filters work on the server, and portfolio
totals intentionally remain independent of the filtered table.

## Manual review

1. Test at desktop width and at 375 px. Cards should stack; only the table should
   need horizontal scrolling.
2. Test light/dark themes and system preference. A saved light theme must remain
   light even when the OS is dark.
3. Disable JavaScript. Apply playground settings, submit the form, filter the
   dashboard, expand disclosures, and toggle the theme.
4. Use Tab, Shift+Tab, Enter, and Space. Verify visible focus, labels, disabled
   controls, the skip link, and the scrollable table region.
5. Test blank fields, malformed email, long names, a successful submission,
   and angle brackets in name/message text.
6. Search for `portal`, filter `paused`, and search for a nonexistent project.
7. Confirm no external assets are requested and no script is loaded.

Automated HTTP tests cover page routes, playground states, escaping, form
validation/body limits, filters, cross-origin rejection, and safe theme redirects:

```sh
go test -race ./examples/basic
```

The current sandbox cannot bind local TCP sockets.
The demos must be launched from a normal local terminal for browser review.
A single review binary may be retained at `tests/results/demos/gothui-demo`;
run it with `-open` to avoid recompilation. The next demo build can replace it,
and it can be removed after review.
