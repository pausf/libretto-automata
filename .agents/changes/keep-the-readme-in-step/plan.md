# keep-the-readme-in-step — plan

Spec: `.agents/changes/keep-the-readme-in-step/spec.md` (Targets: readme)
Branch: `docs/keep-the-readme-in-step`

**Goal:** the README carries every command, and a command that arrives without a row fails
the suite.

**Order is the whole plan.** The guard is written first and watched failing against today's
README, because a guard written after the fix has never been observed to catch anything.

## Global constraints

- Six gates before every commit: `gofmt -l .`, `go vet ./...`, `go test ./... -count=1`,
  `scripts/check-payload`, `spec-drift --self-test`, `spec-drift --anchors`.
- `cmd/libretto/readme_test.go`, `package main`, reaching files through the existing
  `repoFile` helper. No new helper unless `flat()` and `section()` genuinely do not fit.
- Conventional commits, no AI attribution, one commit per task.
- **No version bump.** `release:patch` at most on the request.

## Task 1 — the guard, red

**Files:** modify `cmd/libretto/readme_test.go`

Depends on: nothing.

- [ ] `TestEveryCommandIsInTheReadme` — read `commands/` from disk with `os.ReadDir`, take
      each `*.md` basename without its extension, and assert it appears in
      `section(t, readme, "## Commands")`. The failure names the command that is missing.
- [ ] **run it and watch it fail**, in the foreground, output read. Expected: one failure
      naming `libretto-attacca`. A guard that has never been red has never been shown to
      work.
- [ ] commit the failing test on its own, so the red is in history and not only in a
      transcript.

**Closes:** *a command file with no mention in the README's Commands section fails the
suite* — Proof: `cmd/libretto/readme_test.go TestEveryCommandIsInTheReadme`. Closed on the
green run in task 2, not here.

## Task 2 — the row, green

**Files:** modify `README.md`

Depends on: task 1.

- [ ] the `/libretto-attacca` row in the Commands table, one line, matching the shape of
      the five beside it.
- [ ] one line in **Your first run**, where the stops are walked — a command whose whole
      subject is those stops belongs where they are described.
- [ ] run the suite and watch it pass. Same command, opposite result, both observed.
- [ ] gates, then commit.

**Closes:** the same criterion, now green — plus *the sections stay in reading order*
(Proof: `TestReadmeSectionsAreInReadingOrder`), *the first-run walk still names the flow and
its stops* (Proof: `TestReadmeWalksAFirstRun`), and *every relative link resolves* (Proof:
`TestReadmeLinksResolve`).

## Task 3 — land it

Depends on: tasks 1 and 2.

- [ ] apply the delta onto `.agents/specs/readme/spec.md`: the guard as an outcome, its two
      prior decisions with their ceiling, and the `commands/**` path if `Governs:` needs it.
- [ ] delete `.agents/changes/keep-the-readme-in-step/`.
- [ ] `spec-drift --anchors`, then the six gates.
- [ ] one commit: the applied delta and the deleted folder together.

**Closes:** *every `Proof:` citation resolves, file and test name* — Proof:
`skills/record-work/spec-drift --anchors`.

## What can start now

Task 1. Task 2 waits on it by design — the point is the red — and task 3 waits on both.
