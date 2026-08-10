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

## Scope boundaries

**In:** one workflow file, the six gates, per-package coverage as output, and `rg` in
the runner.

**Out:**

- **A coverage threshold.** See prior decisions — measurement, not a seventh gate.
- **Building or releasing anything.** No artefacts, no `make build` beyond what
  `go test` compiles, no tags, no publishing. The tag is a human act by `AGENTS.md`
  and stays one.
- **A matrix.** One Go version — the one `go.mod` pins. This tool runs on the
  developer's machine, not on a support matrix.
- **Caching.** Six gates on a repository this size do not need it, and a cache that
  serves a stale module is a gate that passed for the wrong reason.
- **`.gitlab-ci.yml`.** Settled in the proposal: the remote is GitHub.
- **Turning branch protection on.** A repository setting, not a file. Named as owed.
- **Anything that writes.** No pushes, no comments, no auto-formatting commits. A
  workflow with write access is a workflow that can be made to write.

## Constraints

- **Go is pinned by `go.mod`: 1.26.5.** The workflow reads it from there rather than
  restating it, so the two cannot drift.
- `check-payload` and `spec-drift` are `#!/usr/bin/env bash` and both call `rg`.
  `spec-drift` also calls `git`.
- **`spec-drift` runs from the repository, not from `~/.claude`.** The installed copy
  may be a different version; CI checks the tree it was given.
- The workflow needs **read-only** permissions. Nothing here writes to the repository.
- No secrets. Nothing in this workflow authenticates to anything.
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

## Verification criteria

The workflow's own correctness is checked by running it, and until it has run on a real
request that is a claim rather than a fact. What can be checked here is checked here:

- **the `gofmt` step fails on an unformatted file**, rather than passing because
  `gofmt -l` exits zero
  Proof: cmd/libretto/gates_test.go TestGofmtStepFailsOnUnformattedOutput
- **every gate named in `AGENTS.md` appears in the workflow**, so the two cannot drift
  apart silently
  Proof: cmd/libretto/gates_test.go TestWorkflowRunsEveryGateAgentsNames
- the workflow asks for read-only permissions
  Proof: cmd/libretto/gates_test.go TestWorkflowIsReadOnly
- `make gates` runs the same six commands as the workflow
  Proof: cmd/libretto/gates_test.go TestMakeGatesMatchesTheWorkflow
- **owed, and not provable here:** one real request with the workflow attached, green;
  and branch protection turned on so the check is required. Until both, "cannot merge
  until green" is a sentence in a file rather than a rule anybody is held to.
