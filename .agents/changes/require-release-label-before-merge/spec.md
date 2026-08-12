# Require the release label before the merge

Targets: ci

Governs: .github/workflows/require-release-label.yml cmd/libretto/label_workflow_test.go

`release.yml` refuses a merge with no `release:` label — after the merge, when `main` is
already ahead of its last tag. That is the exact state AGENTS.md names as "a release
nobody can install", discovered by whoever reads the failed run. This change moves the
refusal to a red check on the request itself.

## Outcomes

A new workflow, `.github/workflows/require-release-label.yml`, that:

- runs on `pull_request` against `main`, **including the `labeled` and `unlabeled`
  event types** — adding the label must turn the check green without a new push, and
  removing it must turn the check red
- fails when the request carries zero `release:patch|minor|major` labels, and says
  which three are accepted and where the rule is written (`AGENTS.md#Versioning`)
- fails when it carries two or more, with the same message shape `release.yml` uses
- passes when it carries exactly one
- holds read-only permissions and **checks nothing out** — the labels arrive in the
  event payload, so the job never runs contributor code and never needs the tree

With branch protection requiring this check, a request cannot merge unlabelled, and
`release.yml`'s own refusal becomes the second fence instead of the only one.

## Scope boundaries

**In:** the workflow file, its tests, and — at the end of the flow, once the check has
run on a real request — adding it to `main`'s required status checks via `gh api`, read
back rather than trusted.

**Out:**

- **Touching `gates.yml`.** The check needs `labeled`/`unlabeled` triggers, and giving
  those to `gates.yml` would re-run all six gates every time a label moves. Named so it
  is not "tidied" into one file later.
- **Weakening or removing `release.yml`'s own refusal.** Branch protection is a
  setting, and a setting can be off. The post-merge refusal stays as defence in depth.
- **Validating the label's *correctness*** (patch vs minor vs major). That is a reading
  of `.agents/specs/` only a human makes — same boundary as `release.yml`.
- **Exempting any request.** AGENTS.md: every merge to `main` gets a tag, so every
  request needs its label. No docs-only escape hatch.

## Constraints

- Untrusted text (label names) reaches the script through `env:`, never `${{ }}`
  inside `run:` — the standing rule from `release.yml`, and its test applies here too.
- Workflow tests live in `cmd/libretto`, as static reads of the YAML plus the extracted
  script run under `bash` — the convention `gates_test.go` and
  `release_workflow_test.go` set.
- The check must exist and have run once before branch protection can require it, so
  the protection step happens after the request is open, not before.

## Prior decisions

- **A separate file, not a job in `gates.yml`** — settled by the trigger types, 2026-08-12.
- **Exactly one label, not at least one** — mirrors `release.yml`, which errors on two.
- **The user chose (2026-08-12): the branch-protection setting is done here**, via
  `gh api` at the end of phase 8, read back like the `gates` check was.

## Task breakdown

1. Write `.github/workflows/require-release-label.yml`.
2. Write `cmd/libretto/label_workflow_test.go` proving the criteria below.
3. After the request is open and the check has run: add it to required checks, read the
   protection back, and record the result in this spec's criteria.

## Verification criteria

- the workflow listens for label changes, so the check re-evaluates without a push
  Proof: cmd/libretto/label_workflow_test.go TestLabelWorkflowRunsWhenLabelsChange
- the script refuses zero release labels
  Proof: cmd/libretto/label_workflow_test.go TestLabelScriptRefusesZeroLabels
- the script refuses two release labels
  Proof: cmd/libretto/label_workflow_test.go TestLabelScriptRefusesTwoLabels
- the script accepts exactly one
  Proof: cmd/libretto/label_workflow_test.go TestLabelScriptAcceptsExactlyOneLabel
- the workflow asks for read-only permissions
  Proof: cmd/libretto/label_workflow_test.go TestLabelWorkflowIsReadOnly
- the workflow never checks out contributor code
  Proof: cmd/libretto/label_workflow_test.go TestLabelWorkflowNeedsNoCheckout
- no `run:` script expands a label inline — untrusted text arrives through `env:`
  Proof: cmd/libretto/label_workflow_test.go TestLabelWorkflowNeverExpandsUntrustedTextInsideAScript
