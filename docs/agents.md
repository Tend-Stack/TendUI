# Using TendUI with a coding agent

Use this guide when building an application with TendUI. The repository-root
[AGENTS.md](../AGENTS.md) is for contributors changing the library itself.
The website version is at [tendui.com/docs/agents](https://tendui.com/docs/agents).

## Give the agent the right context

Tell it whether your application uses plain HTML, Go `html/template`, or Go with
templ. HTMX is an optional enhancement, not a requirement. Share your installed
TendUI version, relevant route and component, and the behavior you want to test.
Use the [quickstart](quickstart.md) and inspect the component source at the
installed version; the online catalog may describe a newer version.

Installing a Go module does not install project instructions into your app.
Merge the following section into your own root `AGENTS.md`; preserve your existing
instructions and adapt commands to your project.

```markdown
## TendUI in this application

- Read go.mod (or the pinned CSS source) to identify the installed TendUI version.
  Do not upgrade it or guess component props/classes to complete a task.
- Prefer semantic HTML and native browser controls. Do not add Node/npm, a CSS
  build step, or client JavaScript solely to use TendUI.
- With Go, import github.com/Tend-Stack/TendUI/v2026 as package tendui.
  Mount tendui.AssetHandler() at tendui.AssetPath without StripPrefix and link
  tendui.StylesheetURL. CSS is embedded and ships inside the executable.
- With templ, use typed components from the /components package. Inspect the
  installed API before writing props. Edit .templ files, regenerate with the
  generator matching the application's templ dependency, and commit generated Go.
- For plain HTML or html/template, serve the pinned stylesheet as a local asset;
  keep existing escaping and server-side rendering. templ is optional.
- Applications own routes, validation, authentication, CSRF protection, sessions
  and persistence. Never copy demo handlers as production security logic.
- HTMX attributes alone do not execute requests. Add the optional runtime only
  when requested; keep native form/link fallbacks and correct fragment responses.
- Use documented --tendui-* tokens and data-theme="light" or "dark"; omit the
  theme attribute for system preference. Preserve labels, focus, keyboard access,
  responsive layout, error messages and reduced-motion behavior.
- Run the application's relevant tests and build. Check the page and stylesheet
  routes and the requested user flow; report checks performed and limitations.
```

## Make the instructions discoverable

[Codex reads AGENTS.md](https://learn.chatgpt.com/docs/agent-configuration/agents-md).
For Claude Code, add this import to a neighboring `CLAUDE.md`, outside a code
fence in that file, preserving any existing instructions:

```markdown
@AGENTS.md
```

This uses [Claude Code's documented imports](https://code.claude.com/docs/en/memory#import-additional-files)
and keeps one shared policy. Other agents may require you to select or explicitly
ask them to read the file. Start a fresh session after setup and ask the agent to
summarize the applicable instructions and installed version before editing.
Instructions guide behavior; tests, reviews and tool permissions still matter.

## A useful first request

> Read our project instructions and installed TendUI API. Add a labeled settings
> form using our existing stack and server routes. Preserve native submission,
> light/dark themes and keyboard behavior. Show the changed files, run relevant
> tests, and explain anything you could not verify. Do not add dependencies.

For upstream contributions, read [CONTRIBUTING.md](../CONTRIBUTING.md). AI-assisted
changes still require human review, testing and an authorized DCO sign-off.
