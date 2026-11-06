# Owner-applied repository settings

This directory is a reviewable configuration, not proof that GitHub settings
are active. An owner with repository administration access should apply the
main ruleset once the contribution workflow is merged and its check names exist.

```sh
gh api --method POST repos/Tend-Stack/TendUI/rulesets --input .github/settings/main-ruleset.json
```

Read existing rulesets first. If one already has this name, update it by its ID
with PUT rather than creating duplicates. The ruleset requires DCO and all CI,
one approving code-owner review, dismissal of stale approvals, approval of the
latest push, resolved conversations and an up-to-date branch. It blocks force
pushes and deletion. Repository admins have an explicit PR-only bypass for
maintainer exceptions; record the reason in the PR and preserve contributor DCO.

Also enable GitHub private vulnerability reporting and web commit sign-off,
disable squash/rebase merges to preserve commit trailers, and allow merge commits.
Require approval for fork workflows according to the organization's trust policy.
Never run untrusted fork code on private self-hosted runners or with deployment
secrets. Workflows use hosted runners, read-only PR permissions and no
`pull_request_target` code execution. Tag release credentials run only on maintainer tags.

Use existing `bug`, `enhancement`, `documentation`, `good first issue` and
`help wanted` labels. Curate bounded starter tasks; do not label complex API or
security changes as beginner work merely to attract contributions.
