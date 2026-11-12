# CI/CD Operations

This document provides guidance for AI agents working with CI/CD operational tasks in the Optimism monorepo.

For Docker image build failures — especially flaky `apt`/`apk`/`curl` downloads from package registries and CDNs — see [docker.md](docker.md).

## Watching CI after a push

Watch every push to a terminal state — `AGENTS.md` requires it. Most jobs run on
CircleCI and report as commit statuses, so `gh` sees them alongside the GitHub Actions
and Wiz checks:

```bash
gh pr checks <pr> --watch --fail-fast   # blocks until done, exits on the first failure
gh pr checks <pr> --required            # only the merge-gating checks
gh pr checks <pr> --json name,bucket,link --jq '.[]|select(.bucket=="fail")'
```

`--watch` blocks until every reported check settles — the `main` workflow alone runs
~25 minutes, past most agent command timeouts — so background it or re-invoke it
instead of running it as a blocking call.

Notes that matter in practice:

- **Wait for the gates the ruleset requires,** not the individual jobs: the four
  CircleCI fan-in gates (`ci-gate`, `required-contracts-ci`, `required-rust-ci`,
  `required-rust-e2e`) plus the `dependency-review` GitHub Actions check. A gate reports
  last, so it can still be pending while every job you were watching is green. On skip
  paths the gate is produced by an `always-succeed` companion, so a gate that never
  reports at all is a config bug (see [ci-config-review.md](ci-config-review.md) item 2),
  not something to wait out.
- **A first push is not the only push to watch.** Rebases, review fixups, and
  merge-queue rebases each start a new pipeline against a different merge base.
- **Triage before rerunning.** Rule out an inherited failure (next section), then check
  whether the test is a known flake. The `generate-flaky-tests-report` job publishes a
  `flaky-test-reports` artifact, but it covers `mc-acceptance-tests` only, is scoped to
  the pipeline's own branch (use a `develop` pipeline's copy, not your PR's) and does not
  run on fast paths; for every other suite, look for an open flake issue instead. A rerun
  that hides a real regression costs more than the minutes it saved, and a confirmed
  flake needs an issue, not a silent retry. Reruns through the CircleCI v2 API need a
  personal API token in `CIRCLE_TOKEN` — not the `CIRCLECI_CLI_TOKEN` the CLI reads for
  [ci-config-review.md](ci-config-review.md)'s `circleci config validate --org-slug`, and
  not the `CIRCLE_API_TOKEN` in `.circleci/`, which is the in-job context token. For
  flakes in `mc-acceptance-tests/`/`mc-devstack/`,
  [flake-prevention.md](flake-prevention.md) catalogues the recurring causes.

## Diagnosing a CI failure on a feature branch

Before assuming a red check is caused by your change, rule out a failure the branch
**inherited** — especially a long-lived branch that has drifted from `develop`:

1. **Check the diff scope:** `git diff origin/develop...HEAD --name-only` (three dots).
   If the failing test exercises code your branch never touches, it almost certainly
   isn't your regression.
2. **Check `develop`:** the failure may be a known flake, or a real bug a *later*
   `develop` commit already fixed — your branch just predates the fix. Look for an open
   flake issue or a recent fix PR touching the failing area.
3. **Rebase onto latest `develop`** before deeper debugging — stale branches miss
   upstream fixes, and a rebase often clears failures that were never about your change.

Worked example (#21356): `go-tests-short` failed on `mc-deployer` integration tests on a
branch that only added a new `mc-core/types` package. The diff touched nothing in
`mc-deployer`; the failures were a data-dependent flake already fixed on `develop` by
#21396. Rebasing made CI green with no change to the branch's own work.

## TODO Checker Failures

The repo runs a scheduled CircleCI job every 4 hours that validates TODO comments don't reference closed GitHub issues. When this job fails, issues need to be reopened.

### Quick Instructions

1. Find the failed TODO checker job in CircleCI (scheduled workflow named `scheduled-todo-issues`)
2. Identify which issues were closed but still have active TODOs in the codebase
3. For each issue:
   - Determine who closed it (using GitHub timeline API)
   - Read the actual TODO comment from the code
   - Reopen with proper attribution and context
   - Include file location and CircleCI job link

### Detailed Workflow

For complete step-by-step instructions with all commands and error handling, see:
**[.claude/skills/fix-todo/SKILL.md](../../.claude/skills/fix-todo/SKILL.md)**

The skill includes:
- Detailed commands for querying CircleCI and GitHub APIs
- How to find who closed an issue
- Comment template for reopening
- Error handling for edge cases
- Output format and requirements
