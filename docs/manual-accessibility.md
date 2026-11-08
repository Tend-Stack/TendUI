# Manual accessibility session

Use this checklist with the actual website and assistive technology. Automated
axe checks, screenshots and browser emulation do not establish screen-reader or
physical-device results. This document is a protocol, **not a completed audit**.
Track observations in [issue #4](https://github.com/Tend-Stack/TendUI/issues/4).

## Before starting

Record the date, tester, device model, OS/version, browser/version, screen reader
and version, input method, and settings such as zoom or reduced motion. Record
the website build and library revision from [the version endpoint](https://tendui.com/version).
For a local build, record the Git commit and library version. If the deployed
revision changes during the session, split the results rather than combining them.

Aim for a desktop keyboard/screen-reader pass and a physical phone pass with
VoiceOver or TalkBack. State the combinations actually tested. Use sample data
such as `Alex` and `alex@example.com`; do not record personal information in
public evidence. Agree on the combination in the issue first to avoid duplicate work.

## Short walkthrough

| Check | Steps | Observe and record |
| --- | --- | --- |
| Shared navigation | On [home](https://tendui.com/), navigate by landmarks and headings, follow Docs, open Examples and choose a demo. On desktop, use Tab and the skip link. | Meaningful names and order; main content reachable; menu state conveyed; focus visible and not obscured by the header. |
| Catalog | On [catalog](https://tendui.com/catalog), find the search field, enter `button`, submit, open the result and return. | Field label, search action and results understandable; component page heading discoverable after navigation; no lost keyboard access. A full page load need not announce as a live region. |
| Forms and errors | On [form lab](https://tendui.com/forms), submit empty fields normally, then use **Test server validation**. Follow an error-summary link, fix fields and submit valid sample data. | Required state, native error and server error text available; summary links reach the right control; descriptions and existing values preserved; success discoverable after navigation. Record what is spoken, not only what is visible. |
| Dialog | On [interactions](https://tendui.com/interactions), activate **Review demo request**. Navigate inside, close, reopen, then dismiss using the platform's available dismiss action. | Dialog name/content announced, modal focus stays inside, closing returns focus to the opener. Record unsupported native opener behavior separately and test **Use the ordinary-page confirmation**. |
| Popover | On the same page, activate **Show project details**, reach its contents, then close it. Repeat with **Show pinned note** and its explicit close control. | Content reachable in reading/focus order; a nonmodal popover does not trap focus; dismissal leaves a sensible focus position. Record inline fallback behavior if unsupported. |
| Theme | Toggle the header theme, revisit Docs and Catalog, and repeat a control check in both modes. | Toggle has a meaningful name; preference persists; text, focus and control state stay distinguishable. Do not infer contrast compliance from visual inspection alone. |
| Physical phone and zoom | On a real phone, repeat navigation, search and form recovery in portrait/landscape with larger text. On desktop, inspect Catalog and Playground at 200% and 400% zoom. | No clipped controls, unreachable content, overlapping header or unnecessary two-direction scrolling; native gestures and touch targets usable. Record actual device and zoom separately from emulation. |
| Reduced motion | Enable the OS reduced-motion setting, reload [motion](https://tendui.com/motion), then compare its explicit Off option. | All content remains present and usable; decorative entrances, lift and skeleton pulses stop. |

If a check fails, note the exact route, settings, steps, expected behavior and
actual observation. Include a short spoken-output transcript when useful. Link a
focused library or website bug instead of hiding the failure behind a broad pass.
Native browser support limitations should have a tested fallback; mark them
**Unsupported** rather than **Pass** if the primary behavior could not be tested.

## Copyable result record

Use **Not run**, **Pass**, **Fail** or **Unsupported** for each row. Keep the initial
Not run values until a person performs that check. An unsupported or failed check
needs notes and a fallback result; omit claims about combinations not tested.

```text
Date/time (UTC):
Tester:
Website build / exact source revision:
Library version / revision:
Device model and physical or emulated:
OS and version:
Browser and version:
Screen reader and version:
Input method / zoom / larger-text / motion settings:

Navigation: Not run
Catalog: Not run
Forms and error recovery: Not run
Dialog and focus return: Not run
Popover and fallback: Not run
Light/dark theme: Not run
Physical phone / zoom / reflow: Not run
Reduced motion: Not run

Observed spoken output and focus behavior:
Limitations / unsupported features:
Bug links and minimal reproduction steps:
Evidence links (no personal data):
Revision unchanged at end of session:
```

Keep the issue open until the agreed manual coverage has actual results and
failures have been triaged. This scoped review does not certify WCAG conformance
or accessibility of applications that consume TendUI.
