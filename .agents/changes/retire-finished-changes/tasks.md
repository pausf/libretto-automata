# Tasks — a finished change that never landed is reported, not omitted

Spec: [spec.md](spec.md) · Plan: [plan.md](plan.md)

Phase 6 (`build-and-check`) drives this file. One box, one commit, one mergeable tree.
The orchestrator owns this file; sub-agents report, they never mark.

**Cut by a fresh agent with none of the design conversation.** The `task-cutter` type was
not registered in the session that ran this seam — the payload item was added and linked
mid-run, and the host loads its agent registry at startup — so the cut ran under a
general-purpose agent carrying `agents/task-cutter.md` as its contract. The independence
comes from the fresh context, not from the label, and that part held.

---

- [x] **Phase 1 names the finished-and-not-landed case**

  `skills/find-work/SKILL.md`, source 1. The two `rg -c` scans already there stay
  exactly as they are; what is added is what their *difference* means, written as the
  plan's table of three outcomes so the empty open scan is a case with a name rather
  than an absence to notice:
  in the open scan → work in flight, the existing report; not in the open scan but in
  the closed one → finished and not landed, reported by name, every box closed, folder
  still present; in neither → a captured idea or a change not cut yet, say nothing.
  Prose only — the skill installs into projects with no `scripts/` and no `docs/`, so
  it may reference neither.

  Plus the `check_wiring` row in `scripts/check-payload` that holds the mandate up, in
  the block at the end of the wiring rows. Before wiring it, `rg` the chosen phrase
  across `skills/`, `commands/` and `agents/` and confirm it appears nowhere else —
  a row matching a phrase that also lives in an example is a row that proves nothing.

  Closes when: `scripts/check-payload` passes with the new row, **and** the row was
  first forced red on purpose by deleting the line it matches and re-running the gate,
  then restored. `check_wiring` is a line-scoped literal match and goes green the
  moment the phrase exists anywhere in the file — a green first run is not evidence
  here. (Plan, *Validation and rollback*.)

  Waits on: nothing. Can start now.

- [x] **`/libretto-status` carries the same report**

  `commands/libretto-status.md` gains the same line, as a delegation to the scan the
  skill owns — not a second description of it. The command already says that a status
  command walking the directories its own way is a second answer to "what work exists";
  the line added here says what the skill reports, and does not restate how it is
  detected.

  Plus the second `check_wiring` row, on the same terms as the first: phrase confirmed
  unique across the payload before it is wired.

  Closes when: `scripts/check-payload` passes with the second row, **and** that row was
  forced red on purpose by deleting the line it matches, then restored.

  Waits on: box 1 — the command delegates to the skill's scan, so the wording it points
  at has to exist before it can point at it.

- [ ] **Land the two changes that were finished and never landed**

  The decisions in `.agents/changes/add-design-phase/` and
  `.agents/changes/retire-plan-decisions/` retired into `.agents/specs/payload/spec.md`
  *Prior decisions*, and **both folders deleted in that same commit** — the constraint
  is the spec's, and `spec-drift --retired` refuses the commit that gets it wrong.

  Order inside the box, and it is load-bearing: stage the two deletions **first**, with
  `payload`'s *Prior decisions* untouched, and run `spec-drift --anchors`. Read the
  refusal. Only then write the migration. A gate whose first real encounter is a pass
  has proved nothing about itself, and this is that gate's first real exercise.

  **Six decisions, not four**, and the sixth is not `payload`'s: five land on
  `payload`'s *Prior decisions* and the `metrics`-fallback one lands on `cli`'s, because
  `cmd/libretto/**` is what `cli` governs. The list is each change's `spec.md` *Prior
  decisions* section — never the plan's `Durable decisions:` line, which is a claim about
  whether the list is empty and not the list itself. All four gaps below are answered in
  the spec's *Scope boundaries* and *Prior decisions*; this box is no longer blocked.

  Closes when: all six gates pass on the staged landing —
  `gofmt -l .` silent, `go vet ./...`, `go test ./... -count=1`,
  `scripts/check-payload`, `spec-drift --self-test`, `spec-drift --anchors` — and the
  earlier refusal of `--anchors` on the deletions-only stage was observed and recorded.

  Waits on: the four gaps being answered. Touches no file boxes 1 and 2 touch, so it is
  mechanically independent of both.

---

## What the cutter said the documents failed to answer

Verbatim, in its words. Every factual claim in them was checked against the tree and
every one held.

- **Which decisions get retired, and how many.** The spec's *Task breakdown* and the
  plan's *Summary* both say **four**; the two change folders carry **six** *Prior
  decisions* between them (three in `add-design-phase/spec.md`, three in
  `retire-plan-decisions/spec.md`), and neither document names which four. The spec's
  *Scope boundaries* should have said.
- **Where a retired decision that is not `payload`'s goes.** `add-design-phase`'s third
  decision — the `metrics` fallback being a second `git log` rather than a two-path
  pathspec — is about `cmd/libretto/`, which `payload`'s `Governs:` does not cover, and
  the spec puts only `.agents/specs/payload/spec.md` in scope. The spec should have said
  whether it is dropped, reworded, or lands on `cli`.
- **Which document is the source of the decisions being retired** — the change's
  `plan.md` `Durable decisions:` line or its `spec.md` *Prior decisions* section. They
  disagree already: `retire-plan-decisions/plan.md` declares "the two in *Prior
  decisions* below" over a section holding three, and `add-design-phase/plan.md` carries
  no `Durable decisions:` line at all. The spec should have said which one the
  retirement reads.
- **The wording each `check_wiring` row matches.** The plan names the *property* the
  phrase must have (unique across the payload, sitting in the mandate rather than in an
  example) but not the phrase, so the row and the sentence it guards are both chosen at
  build time, and the gate's whole contract is that literal string. The plan's
  *Validation and rollback* should have named both.
