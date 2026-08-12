# ask-before-planning-and-sync-main — plan

Spec: `spec.md` in this folder. Targets `payload`.

Five tasks. Tasks 1 and 2 are the two halves of the ask and are independent of each
other — either can start now. Tasks 3 and 4 describe what 1 and 2 did, so they wait on
both. Task 5 lands the delta and closes the change.

```
1 ─┐
   ├─▶ 3 ─▶ 4 ─▶ 5
2 ─┘
```

## Tasks

- [ ] **1 · Phase 8 returns to the base branch**
  `skills/record-work/SKILL.md` — after the push and the request are both confirmed,
  derive the base branch, `git checkout` it and `git pull --ff-only`. Yes path only.
  Dirty tree or a refused fast-forward is reported, never forced. The feature branch is
  left alone.
  *From:* spec · outcome 1 · **Waits on:** nothing
  *Closes:* the run answering yes ends on the base branch, current, with the branch intact

- [ ] **2 · Phase 2 asks up to three**
  `skills/write-spec/SKILL.md` — step 4 becomes: up to three questions the code cannot
  answer, one `AskUserQuestion` call, before the spec file is written, zero when there is
  nothing real to ask. Answers land under prior decisions.
  *From:* spec · outcome 2 · **Waits on:** nothing
  *Closes:* a phase 2 with open points asks at most three in one call; one with none says so

- [ ] **3 · The command mirrors both**
  `commands/libretto-flow.md` — phase 8's line gains the return to base; phase 4's line
  says up to three, phase 2 only, zero allowed. No new stop in the table.
  *From:* spec · scope boundaries · **Waits on:** 1, 2
  *Closes:* `scripts/check-payload`

- [ ] **4 · The reasoning goes in `docs/FLOW.md`**
  Why the base branch is checked out rather than fetched, and why the questions sit at
  phase 2 alone with no minimum. Reasoning belongs here, not in the skills.
  *From:* spec · prior decisions · **Waits on:** 3
  *Closes:* read back; no gate covers prose

- [ ] **5 · Land the delta**
  Apply both outcomes onto `.agents/specs/payload/spec.md`, delete this change folder,
  same commit as the last code. Run all six gates.
  *From:* spec · task breakdown · **Waits on:** 4
  *Closes:* `skills/record-work/spec-drift --anchors`, `scripts/check-payload`, `go test ./...`
