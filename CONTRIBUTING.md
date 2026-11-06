# Contributing to TendUI

Thank you for helping improve TendUI. Start with a bug, documentation fix,
accessibility improvement or an issue labeled `good first issue` or `help wanted`.
Discuss new components, public API changes, dependencies and large refactors in an
issue before implementation. A proposal should explain the use case, native HTML
semantics, keyboard behavior, narrow-screen behavior and alternatives considered.

## Contribution terms: MIT + DCO

Contributions are submitted under this repository's [MIT license](LICENSE).
You keep your copyright; no copyright assignment or separate CLA is required.
Every non-merge commit must include a sign-off certifying the [Developer
Certificate of Origin 1.1](DCO). Read it before signing. Obtain any employer or
third-party permission needed, retain attribution, and do not submit incompatible
licensed material. The sign-off is a certification, not a cryptographic signature.

```sh
git commit -s -m "Describe the change"
```

Use your own name and the same email as the commit author. A GitHub-provided
noreply address is accepted. Names, commit emails and sign-offs become public
and remain in history. Never add another person's sign-off without authorization.
For your own last commit, use `git commit --amend --signoff`; coordinate with
reviewers before rewriting an already-pushed contribution branch. Maintainers do
not add missing contributor certifications on someone else's behalf.

AI-assisted contributions follow exactly the same rules: a human contributor
must understand, review and test the result and certify the right to submit it.
Identify significant generated material or reused sources in the PR description.

## Workflow

1. Fork the repository and create a focused branch from current `main`.
2. Make the smallest complete change; keep unrelated cleanup separate.
3. Add proportionate tests and update documentation for changed behavior.
4. Regenerate templ output when templates change; commit both source and generated Go.
5. Open a draft PR early when feedback would help. Link the issue and describe validation.
6. Respond to review; maintainers merge after required checks and approval.

Go 1.26 or newer is required. Template authors use templ v0.3.1100. Node/npm are
only used for optional contributor browser tooling, never the shipped UI.

```sh
go run ./cmd/tendui-dev check
go run ./cmd/tendui-dev generated
go run ./cmd/tendui-dev recipes
# After committing, verify isolated consumption:
go run ./cmd/tendui-dev consumer
```

CI supplies native-platform and browser coverage. Do not claim a check ran when
it did not; explain unavailable checks in the PR. Keep generated test evidence
under ignored `tests/results/` and remove owned disposable artifacts after use.

## Design and compatibility

- Semantic HTML and native browser behavior first; no bundled client JavaScript,
  remote fonts or runtime asset downloads.
- Keep Go configuration typed. Applications own routing, persistence, sessions,
  request protection and any optional HTMX runtime.
- Preserve keyboard access, accessible names, theme behavior and reduced motion.
- Public API changes need an explicit compatibility review and migration notes,
  even during preview. Never silently replace old tags.
- Include screenshots for visual changes and narrow-screen checks where relevant.
- Never commit credentials, private infrastructure details, personal datasets or logs.

## Reviews and releases

CODEOWNERS identifies the maintainer. External contributions require maintainer
approval, DCO and applicable CI. An administrator may make a documented exception
through a PR; never use an exception to accept an unsigned contributor submission.
Prefer merge commits to preserve contributor sign-offs. Release tags are created
by maintainers after exact-source checks; previews use `YEAR.RELEASE.FIX-preview.N`.

Report security issues privately using [SECURITY.md](SECURITY.md). Community
participation follows [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
