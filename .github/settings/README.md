# Owner-applied repository settings

The **Reviewed main contributions** ruleset was applied and read back as active
on 2026-11-07 UTC ([view ruleset](https://github.com/Tend-Stack/TendUI/rules/26745168)).
This records the verified setup date; GitHub remains the source of truth for
current enforcement. Review the live rules before changing them.

To update this existing ruleset after reviewing the JSON:

```sh
gh api --method PUT repos/Tend-Stack/TendUI/rulesets/26745168 --input .github/settings/main-ruleset.json
```

Do not POST another copy. The ruleset requires DCO and all CI,
one approving code-owner review, dismissal of stale approvals, approval of the
latest push, resolved conversations and an up-to-date branch. It blocks force
pushes and deletion. Repository admins have an explicit PR-only bypass for
maintainer exceptions; record the reason in the PR and preserve contributor DCO.

Separate settings still require verification: enable GitHub private vulnerability reporting and web commit sign-off,
disable squash/rebase merges to preserve commit trailers, and allow merge commits.
Require approval for fork workflows according to the organization's trust policy.
Never run untrusted fork code on private self-hosted runners or with deployment
secrets. Workflows use hosted runners, read-only PR permissions and no
`pull_request_target` code execution. Tag release credentials run only on maintainer tags.

Use existing `bug`, `enhancement`, `documentation`, `good first issue` and
`help wanted` labels. Curate bounded starter tasks; do not label complex API or
security changes as beginner work merely to attract contributions.
