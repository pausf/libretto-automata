# keep-the-readme-in-step — plan

Specs: `spec.md` (Targets: readme) · `spec-stops.md` (Targets: payload)
Branch: `docs/keep-the-readme-in-step`

**Goal:** the README carries every command, a command that arrives without a row fails the
suite, and every stop in the flow is a native question rather than a paragraph.

Two deltas, one change. They share nothing but the branch — tasks 1–3 are the README, tasks
4–6 are the stops, and neither waits on the other.

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

- [x] `TestEveryCommandIsInTheReadme` — read `commands/` from disk with `os.ReadDir`, take
      each `*.md` basename without its extension, and assert it appears in the README. The
      failure names the command that is missing.
      **Corrected mid-task:** written first against `section(t, readme, "## Commands")`,
      which is the *binary's* subcommand table — the slash commands live in the first-run
      door list. It failed all six. Now the whole file, and the spec's outcome 1 moved with
      it, in this commit.
- [x] **run it and watch it fail**, in the foreground, output read. One failure naming
      `libretto-attacca`. A guard that has never been red has never been shown to work.
- [x] commit the failing test on its own, so the red is in history and not only in a
      transcript.

**Closes:** *a command file with no mention in the README fails the suite* — Proof:
`cmd/libretto/readme_test.go TestEveryCommandIsInTheReadme`. Closed on the green run in
task 2, not here.

## Task 2 — the row, green

**Files:** modify `README.md`

Depends on: task 1.

- [x] the `/libretto-attacca` line in the first-run door list, matching the shape of the
      five beside it.
- [x] one line in **Your first run**, where the stops are walked — a command whose whole
      subject is those stops belongs where they are described.
- [x] run the suite and watch it pass. Same command, opposite result, both observed.
- [x] gates, then commit.

**Closes:** the same criterion, now green — plus *the sections stay in reading order*
(Proof: `TestReadmeSectionsAreInReadingOrder`), *the first-run walk still names the flow and
its stops* (Proof: `TestReadmeWalksAFirstRun`), and *every relative link resolves* (Proof:
`TestReadmeLinksResolve`).

## Task 3 — the README delta, applied

Depends on: tasks 1 and 2.

- [x] apply `spec.md` onto `.agents/specs/readme/spec.md`: the guard as an outcome, its two
      prior decisions with their ceiling, and the `commands/**` path if `Governs:` needs it.
- [x] gates, then commit.

**Closes:** *every `Proof:` citation resolves, file and test name* — Proof:
`skills/record-work/spec-drift --anchors`.

## Task 4 — the two stops that wait

**Files:** modify `skills/write-spec/SKILL.md`, `skills/write-plan/SKILL.md`

Depends on: nothing.

- [x] `write-spec` — its stop asked with `AskUserQuestion`: carry on to the plan
      (recommended, saying what runs next), or change the contract first.
- [x] `write-plan` — the same shape: start the work, or change the order first.
- [x] each says *that* it asks natively and what its options mean, and neither restates
      `record-work`'s argument for why. That argument lives once.
- [x] gates, then commit.

**Closes:** *every referenced skill exists and frontmatter parses* — Proof:
`scripts/check-payload`.

## Task 5 — the choice phase 1 already makes

**Files:** modify `skills/find-work/SKILL.md`

Depends on: nothing. Independent of task 4.

- [ ] the in-flight choice asked with `AskUserQuestion`, the same way the which-task
      question three sections below it already is. It already says *never choose*; what it
      never said is how to ask.
- [ ] gates, then commit.

**Closes:** same proof as task 4.

## Task 6 — said once where the stops are argued for, then land

Depends on: tasks 3, 4 and 5.

- [ ] `commands/libretto-flow.md` and `docs/FLOW.md` — the stops are native questions, one
      line each, beside the table that already lists them.
- [ ] apply `spec-stops.md` onto `.agents/specs/payload/spec.md`, including the ceiling: no
      check can tell a native prompt from a paragraph, and the guard that would is named
      and deliberately not built.
- [ ] delete `.agents/changes/keep-the-readme-in-step/`.
- [ ] `spec-drift --anchors`, then the six gates.
- [ ] one commit: the final prose, both applied deltas' remainder, and the deleted folder.

**Closes:** *every `Proof:` citation resolves* — Proof:
`skills/record-work/spec-drift --anchors`.

## What can start now

**Tasks 1 and 4 and 5**, independently. Task 2 waits on 1 by design — the point is the red
— task 3 waits on both, and task 6 waits on everything because it deletes the folder.
