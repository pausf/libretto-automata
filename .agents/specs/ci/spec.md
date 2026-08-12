# Continuous Integration

Governs: .github/** Makefile

The six gates, run by a machine that is not the author's, before a request can merge.

## Outcomes

One workflow, `.github/workflows/gates.yml`, on every push and every pull request. It
runs **the same six gates `AGENTS.md` names**, in the same order, with the same
commands:

```
gofmt -l .                          must print nothing
go vet ./...
go test ./... -count=1
scripts/check-payload
skills/record-work/spec-drift --self-test
skills/record-work/spec-drift --anchors
```

- **The same six, not a CI-flavoured subset.** A pipeline that checks four of them
  teaches people the other two are optional, and the two it would drop — the payload
  ones — guard the reason the project exists.
- **Each gate is its own step**, so a red build names which gate failed on the summary
  page rather than in the middle of a log.
- **`gofmt` fails when it prints anything.** `gofmt -l` exits zero with a list of
  unformatted files, so a naive step passes on a repository that is not formatted. The
  step checks the output, not the status.
- **`rg` is installed in the runner.** `check-payload` and `spec-drift` are bash and
  both shell out to it; without it they fail for a reason that has nothing to do with
  the code.
- **The workflow reports coverage per package and fails on none of it.**
- **A request cannot merge until the workflow is green** — that last part is a branch
  protection setting in the repository, not a file, and this change cannot make it
  true. It is named here so it is not mistaken for done.

### `--anchors` is a merge gate, and that is the point

`spec-drift --anchors` exits non-zero when a `Proof:` citation names a test that does
not exist. Mid-change that is **normal and expected**: `write-spec` says to cite tests
before writing them, so a branch part-way through a change has broken citations by
design.

It is still the right gate, because of when a request gets opened. Phase 8 folds the
delta onto the capability spec and deletes the change folder in the same commit — so a
request that has been through the flow has no forward citations left and `--anchors` is
green. Every request opened this session was.

**The consequence, stated so it is a decision rather than a surprise:** a request
opened part-way through a change fails this gate until the change is consolidated.
That is the gate working. A specification citing tests that do not exist is exactly
what should not merge.

### The label is required before the merge, not discovered after it

`release.yml` refuses an unlabelled merge — after the merge, when `main` is already
ahead of its last tag and nobody can install what landed. A third workflow,
`.github/workflows/require-release-label.yml`, makes the same refusal a red check on
the request itself:

- **it triggers on `labeled` and `unlabeled`**, which are not in the default
  `pull_request` set — without them, adding the label leaves the check red until the
  next push, and a check only a push can satisfy is a check people learn to work
  around. This is why it is its own file: giving those types to `gates.yml` would
  re-run all six gates every time a label moves.
- it fails on zero `release:patch|minor|major` labels, naming the three and
  `AGENTS.md#Versioning`; fails on two or more; passes on exactly one.
- **matching is by whole label name.** Names may contain spaces, so a space-joined
  loop counts a label merely *containing* `release:patch` as the real thing. The
  labels arrive as `toJSON` output and `jq` compares whole names.
- **it checks nothing out and holds no permissions** (`permissions: {}`). The labels
  are in the event payload, so the job never runs contributor code.
- `release.yml`'s own refusal stays, as the second fence for when branch protection
  is off. It still parses space-joined — the same word-split hole this check closed —
  and hardening it is a follow-up, not part of this promise.
- it never judges whether the label is the *correct* bump. That is a reading of
  `.agents/specs/` only a human makes — the same boundary `release.yml` keeps.

### Releasing happens on merge, and `make release` stays as the fallback

**This section used to say the opposite**, and the reversal is deliberate rather than
forgotten. It read *"releasing is a `make` target, not a workflow"*, and it named its own
ceiling: *the day releases are frequent enough that doing it by hand gets skipped, a workflow
is the answer.* That day arrived and was measurable — `main` sat three commits past `v0.5.0`
untagged, and the panel reported `v0.5.0-3-g376fc46` to a user who read it as a bug. The
procedure was correct and nobody ran it, which is the failure mode a manual step has.

What stopped being true is *"the tag is a human act"*. What is still true, and is the reason
this is not simply automation, is that **the bump is a human judgment** — so the human still
makes the decision, on the request, as a label, and the workflow only executes it.

```
merge a request carrying release:patch|minor|major
  └─ .github/workflows/release.yml
       gates → derive the bump → tag main → push the tag → create the Release
```

`make release` is kept, unchanged, as the hand-run path for a tag created outside the flow.
It is no longer the documented route.

**A workflow that decided the bump itself would make both halves false**, and that is the line
this design does not cross: with no label the run refuses, names the three labels, and points
at `AGENTS.md#Versioning`. It never picks one.

- **The payload is not an asset, because it does not need to be.** It ships inside the Go module,
  so `go install <module>/cmd/libretto@latest` brings it down with the binary and checks it
  against the checksum database. The module proxy resolves `@latest` from **tags**, which is why
  installing works against a repository with no Releases at all.
- **So this target is for humans:** a Release page carrying the tag's own notes, so somebody can
  read what changed without reading the log. Skipping it costs nothing mechanical.
- **The job checks out `main`, never the request's own tip.** It is the one job holding
  `contents: write` and it is fired by a pull-request event, so running the contributor's tree
  would put the token in reach of whoever opened the request. The merged work is already on
  `main`, and that is what gets tagged.
- **A closed request is not a merged one.** Without the `merged == true` guard, closing a request
  without merging would tag `main` at whatever it happens to be.
- **Releases are serialised.** Two merges landing together would read the same last tag and the
  second would try to create one that exists. Never cancel-in-progress: a half-finished release
  is worse than a queued one.
- **Untrusted text reaches the scripts through the environment.** A label, a title and a body are
  text somebody else wrote, and `${{ }}` expanded inside a shell line is that text becoming
  executable.
- **The push comes before the target, and the order is load-bearing** — in `make release`, which
  is still there. `gh release create`
  creates the tag itself when it is not on the remote, at the default branch's HEAD — a second
  tag with your name on it pointing somewhere you did not choose. `--verify-tag` refuses instead.
- **A tag is not a Release.** A tag is a git ref; a Release is a GitHub object. `git push origin
  v0.5.0` creates the first and not the second, and this repository sat at four tags and zero
  Releases — verified, not assumed.
- **It refuses a dirty tree and a HEAD that is not at a tag**, and runs the six gates first. A
  release is the one moment where "the gates were green earlier" is not good enough.
- **Re-running is a no-op** when the Release already exists, rather than an error or a second one.

## Scope boundaries

**In:** three workflow files, the six gates, per-package coverage as output, `rg` in the runner,
the `release` target with its four preconditions, **one job that tags `main` and publishes
the Release when a request merges**, and **one check that refuses an unlabelled request
before it can merge**.

**Out:**

- **A coverage threshold.** See prior decisions — measurement, not a seventh gate.
- **Building artefacts in CI.** No tarballs, no cross-compilation, no `make build` beyond what
  `go test` compiles. Tagging and publishing are now in; building things to attach is not, for
  the reason the next bullet gives.
- **Deriving the bump from the commit log.** `AGENTS.md` makes a new promise in an existing spec
  a minor, and that is a reading of `.agents/specs/` rather than of `type:` prefixes. A workflow
  reading commit types picks the floor and is wrong precisely when a contract moves. The bump
  comes from a label on the request, and its absence stops the run.
- **`release:major`, while this is `0.x`.** The label exists and the workflow honours it, but
  `AGENTS.md` reserves it: in `0.x` a reversed promise is a minor, and `1.0.0` is a declaration
  about the tool's contract that nobody has made. **This spec's own reversal is what taught
  that** — it was labelled `release:major`, published `v1.0.0` and `v1.0.1`, and the tool's
  contract had not moved by a line. The tags were withdrawn; the versions are cached in
  `proxy.golang.org` permanently and neither number can ever name different content.
- **Defaulting to patch when the label is missing.** That is the silently-wrong bump wearing a
  different hat. The run refuses and says which label to add.
- **release assets.** Nothing is attached: the payload ships inside the Go module and the proxy
  resolves `@latest` from tags. A draft specified a tarball, checksums and four cross-compiled
  binaries; `distribution` records why none is needed.
- **A matrix.** One Go version — the one `go.mod` pins. This tool runs on the
  developer's machine, not on a support matrix.
- **Caching.** Six gates on a repository this size do not need it, and a cache that
  serves a stale module is a gate that passed for the wrong reason.
- **`.gitlab-ci.yml`.** Settled in the proposal: the remote is GitHub.
- **Branch protection as a file.** It is a repository setting and no file in this tree can
  express it. **It is now on** — see the criteria below for what it is set to and how that was
  read back. It stays named here because a reader looking for it in the tree will not find it,
  and "absent from the repository" must not read as "not done".
- **Write access anywhere except the release job.** `gates.yml` stays `contents: read`, and
  nothing anywhere comments, auto-formats or commits. `release.yml` gets `contents: write`
  because a tag and a Release cannot be created without it — and it is the reason that job
  checks out `main` rather than the request's own tip. **A workflow with write access is a
  workflow that can be made to write, so the one that has it never runs contributor code.**

## Constraints

- **Go is pinned by `go.mod`: 1.26.5.** The workflow reads it from there rather than
  restating it, so the two cannot drift.
- `check-payload` and `spec-drift` are `#!/usr/bin/env bash` and both call `rg`.
  `spec-drift` also calls `git`.
- **`spec-drift` runs from the repository, not from `~/.claude`.** The installed copy
  may be a different version; CI checks the tree it was given.
- **`gates.yml` needs read-only permissions and no secrets.** Nothing in it writes to the
  repository or authenticates to anything.
- **`release.yml` needs `contents: write` and `GITHUB_TOKEN`, and nothing else.** No secret is
  added to the repository for it: the token Actions already issues is enough to push a tag and
  create a Release.
- **A tag pushed with `GITHUB_TOKEN` does not trigger workflows.** GitHub suppresses it to stop
  a workflow retriggering itself, so the Release must be created in the same job that pushes the
  tag. A second workflow listening on `push: tags` would never run, and the only evidence would
  be a tag with no Release.
- **`fetch-depth: 0` on the release checkout.** The default is depth 1 with no tags, and
  `git describe --tags` against that computes the next version from zero.
- **The workflow's own tests live in `cmd/libretto/gates_test.go`.** A Go test needs a
  package, and a directory holding one test file and no code to justify it is the
  over-engineering this repository argues with elsewhere. `cmd` is already where
  repository-level concerns land. The `cli` spec governs that path; this one owns the
  workflow, and a `Proof:` may cite a test wherever the test can honestly live.

## Prior decisions

- **GitHub Actions, not `.gitlab-ci.yml`.** From the proposal: the remote is GitHub, so
  a GitLab file would sit in the tree looking like a gate while every request merged
  unchecked. Both files was rejected too — two definitions of one set of gates diverge,
  and the one nobody watches is the one that lies.

- **Coverage is measured and printed, never enforced.** A threshold is a seventh gate,
  and this repository has been deliberate about having six. It also changes what tests
  get written: a number to defend produces tests that move the number, which is the
  opposite of the discipline the rest of the flow asks for. **The ceiling:** if
  coverage falls somewhere that matters and nobody notices for a release, a threshold
  on that package alone is the answer — not a repository-wide floor.

- **The six gates are copied, not re-invented.** The workflow runs the exact commands
  `AGENTS.md` lists. If a gate changes, it changes in one place and the workflow
  follows; a CI file that paraphrases its gates is a second source of truth for what
  "checked" means.

## Task breakdown

1. `.github/workflows/gates.yml`: the six gates, one step each, `rg` installed, Go from
   `go.mod`, coverage printed, read-only permissions.
2. `Makefile`: a `gates` target that runs the same six locally, so the workflow and the
   hand-run are the same list rather than two lists that agree today.
3. `Makefile`: a `release` target — four preconditions, the gates, then `gh release create
   --verify-tag --notes-from-tag`, idempotent against an existing Release.
4. `.github/workflows/release.yml`: on a merged request into `main` — `main` checked out at
   full depth, the gates, the bump read off a `release:` label or the run refused, an
   annotated tag carrying the request's title and body, the tag pushed, and the Release
   created in the same job.

## Verification criteria

The workflow's own correctness is checked by running it, and until it has run on a real
request that is a claim rather than a fact. What can be checked here is checked here:

- **the `gofmt` step fails on an unformatted file**, rather than passing because
  `gofmt -l` exits zero
  Proof: cmd/libretto/gates_test.go TestGofmtStepFailsOnUnformattedOutput
- **every gate named in `AGENTS.md` appears in the workflow**, so the two cannot drift
  apart silently
  Proof: cmd/libretto/gates_test.go TestWorkflowRunsEveryGateAgentsNames
- the gates workflow asks for read-only permissions
  Proof: cmd/libretto/gates_test.go TestWorkflowIsReadOnly
- `make gates` runs the same six commands as the workflow
  Proof: cmd/libretto/gates_test.go TestMakeGatesMatchesTheWorkflow
- the release target runs the gates before it publishes anything
  Proof: cmd/libretto/gates_test.go TestReleaseTargetRunsTheGatesFirst
- **it verifies the tag rather than letting `gh` invent one** at the default branch's HEAD
  Proof: cmd/libretto/gates_test.go TestReleaseVerifiesTheTagRatherThanCreatingIt
The release workflow's own criteria, each one a failure that is silent or dangerous rather
than loud:

- **it checks out `main` and never the request's own tip** — the job holds `contents: write`
  Proof: cmd/libretto/release_workflow_test.go TestReleaseWorkflowChecksOutMainAndNotTheRequestHead
- **it fetches deep enough to see the tags**, so the next version is not computed from zero
  Proof: cmd/libretto/release_workflow_test.go TestReleaseWorkflowFetchesEnoughHistoryToSeeTheTags
- **it creates the Release in the same run that pushes the tag**, because a `GITHUB_TOKEN` tag
  push triggers no workflow
  Proof: cmd/libretto/release_workflow_test.go TestReleaseWorkflowPublishesInTheSameRunThatPushesTheTag
- **it refuses when no `release:` label is present** rather than defaulting to patch
  Proof: cmd/libretto/release_workflow_test.go TestReleaseWorkflowRefusesWithoutABumpLabel
- it runs the gates before it tags
  Proof: cmd/libretto/release_workflow_test.go TestReleaseWorkflowRunsTheGatesBeforeItTags
- concurrent merges are serialised and never cancelled
  Proof: cmd/libretto/release_workflow_test.go TestReleaseWorkflowSerialisesConcurrentMerges
- a request that was closed without merging does not tag anything
  Proof: cmd/libretto/release_workflow_test.go TestReleaseWorkflowOnlyRunsOnAMergedRequest
- **no `run:` script expands a label, a title or a body** — untrusted text arrives through `env:`
  Proof: cmd/libretto/release_workflow_test.go TestReleaseWorkflowNeverExpandsUntrustedTextInsideAScript

The label check's criteria — the script ones run the workflow's own extracted script
under bash, because "refuses zero and two" is behaviour, not shape:

- the workflow listens for label changes, so the check re-evaluates without a push
  Proof: cmd/libretto/label_workflow_test.go TestLabelWorkflowRunsWhenLabelsChange
- the script refuses zero release labels
  Proof: cmd/libretto/label_workflow_test.go TestLabelScriptRefusesZeroLabels
- the script refuses two release labels
  Proof: cmd/libretto/label_workflow_test.go TestLabelScriptRefusesTwoLabels
- the script accepts exactly one
  Proof: cmd/libretto/label_workflow_test.go TestLabelScriptAcceptsExactlyOneLabel
- matching is by whole label name, never by substring or word
  Proof: cmd/libretto/label_workflow_test.go TestLabelScriptMatchesLabelNamesExactly
- the workflow asks for no permissions at all
  Proof: cmd/libretto/label_workflow_test.go TestLabelWorkflowIsReadOnly
- the workflow never checks out contributor code
  Proof: cmd/libretto/label_workflow_test.go TestLabelWorkflowNeedsNoCheckout
- no `run:` script expands a label inline
  Proof: cmd/libretto/label_workflow_test.go TestLabelWorkflowNeverExpandsUntrustedTextInsideAScript

**Paid, by observation:** the check ran green on its own request (#40, the first it
could run on), and it now sits beside `gates` in the required status checks — set via
`gh api` after that first run, then read back from the forge rather than trusted:

```
checks: ["gates", "label"]   pr_required: true   strict: false   admins_enforced: false
```

- **paid, by observation:** the gates workflow has run green on real requests, and **two real
  merges have tagged and published** — requests #25 and #26, runs `31485394056` and
  `31485661677`, every step `success`, the tag on the remote pointing at `main`'s tip, a Release
  carrying the request's title, and `libretto version` reporting it from a fresh build.

  **The mechanism was proven; the numbers it produced were withdrawn.** Both runs derived their
  bump correctly from the label they were given — the labels were wrong, not the workflow. See
  the `release:major` boundary above.
- **branch protection is on**, so "a request cannot merge until the workflow is green" is a rule
  and no longer a sentence. Read back from the forge rather than trusted to the call that set
  it:

  ```
  checks: ["gates"]   pr_required: true   approvals: 0   strict: false   admins_enforced: false
  ```

  Four of those five are decisions worth stating:

  - **`approvals: 0`.** A required approval on a single-maintainer repository blocks every
    request permanently, because nobody can approve their own. Requiring the *request* is what
    matters here — it is what routes every change past the gates.
  - **`admins_enforced: false`** leaves an escape hatch. A protection that can lock the only
    maintainer out of their own `main` gets switched off entirely the first time it does, and
    then there is no protection at all.
  - **`strict: false`** — a request does not have to be rebased onto a moved `main` before
    merging. The gates run against the merge result regardless, and requiring a rebase per merge
    buys friction rather than safety at this size.
  - **Tags are not covered.** Branch protection applies to `refs/heads`, and there are no tag
    rulesets — verified, because the release workflow pushes a tag and nothing else. Protecting
    `main` would otherwise have silently broken every release.

  **A request introducing a change to this workflow triggers the changed version on its own
  merge.** Observed, not reasoned: request #25 added `release.yml` and its own merge ran it,
  deriving the bump from its own label and publishing the Release. `pull_request` events run the
  workflow file as it exists in the request rather than in the base.

  **That is a property to be careful with, not a convenience.** A request that breaks this
  workflow breaks the release of the merge that lands it, and the evidence arrives after the
  merge. So a change here carries its `release:` label like any other, and the run is read
  afterwards rather than assumed.
