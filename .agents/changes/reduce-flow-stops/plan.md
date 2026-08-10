# reduce-flow-stops — plan

Spec: `.agents/changes/reduce-flow-stops/spec.md` (Targets: payload)
Branch: `feat/reduce-flow-stops`

Seven tasks. Tasks 1–6 are independent of one another — they are six separate files
saying one thing, and the risk is not ordering, it is one of them being left saying the
old thing. Task 7 waits on all six.

Phase 6 does not fan out. Serial, one writer, one commit per task.

## Tasks

- [ ] **1 · `skills/review-work/SKILL.md` — the fix pass**
      Step 3 stays (relay, attributed and unedited). A new step 4 fixes every finding in
      the same pass, re-runs the proofs the fix touches, and adds no question. The
      two-failures stop from `skills/evidence/`. A finding that cannot be fixed without a
      decision that is not ours is reported to phase 7, never turned back into a
      question. The closing "Then stop" keeps its meaning — presenting is still phase 7's.
      Frontmatter `description:` still ends "it reports, it never blocks" — that sentence
      is now false and is part of this task.
      From: Outcomes · "the review seam acts on what it finds"
      Closes: `scripts/check-payload`
      Waits on: nothing

- [ ] **2 · `agents/work-reviewer.md` — confirm, do not change**
      Read it and verify the read-only grant and the no-edit rule are intact and still
      correct beside task 1. If it already says the right thing, this task closes with
      that observation and no diff. Named as a task because "the agent stays read-only"
      is a decision in the spec, and a decision nobody checked is a decision nobody made.
      From: Scope boundaries · "the reviewer becoming a writer" is out
      Closes: `scripts/check-payload`
      Waits on: nothing

- [ ] **3 · `skills/find-work/SKILL.md` — no wait after the artifact**
      Line 39, "Then confirm the reading, and stop." The artifact and the confirmation
      both survive; the stop goes. Line 278's "Then stop. Nothing else happens here."
      becomes the phase's boundary, not a wait. Line 172's stop stays untouched — that is
      the unconfigured-tracker stop, which is a real blocker and not ceremony.
      From: Outcomes · phase 1 carries into phase 2
      Closes: `scripts/check-payload`
      Waits on: nothing

- [ ] **4 · `skills/present-work/SKILL.md` — report, then phase 8**
      Line 12 and lines 127–129. Every sentence of the report survives; the wait does
      not. The trivial-lane exception at line 127 stops being an exception, because both
      lanes now carry straight on.
      From: Outcomes · phase 7 carries into phase 8
      Closes: `scripts/check-payload`
      Waits on: nothing

- [ ] **5 · `commands/libretto-flow.md` — two waits, not four**
      Lines 57, 76, 101, 139. Remove the phase 1 and phase 7 waits, keep the spec and
      plan waits, and restate the review seam's new standing at lines 120–127.
      From: Outcomes · the stop table
      Closes: `scripts/check-payload`
      Waits on: nothing

- [ ] **6 · `docs/FLOW.md` — the reasoning**
      Lines 153–176 and 225. Why three stops and not four: each remaining stop is a place
      the user changes something. Why the seam fixes rather than reports, and why the
      reviewer itself stays read-only. Not installed, so nothing depends on it — which is
      exactly why it is the one that gets left behind.
      From: Constraints · four places must agree
      Closes: observation only — `docs/` is outside every `Governs:`
      Waits on: nothing

- [ ] **7 · land the delta**
      Apply the delta onto `.agents/specs/payload/spec.md` — the four-waits promise in
      Outcomes, the "a review that blocks, fixes or opines" boundary, the ceremony
      section, and the two prior decisions about the reviewer's standing. Fix
      `.agents/specs/review-project/spec.md` line 112, whose analogy no longer holds.
      Delete `.agents/changes/reduce-flow-stops/`. One commit with the final code.
      From: Task breakdown · 7
      Closes: `skills/record-work/spec-drift --anchors`
      Waits on: 1, 2, 3, 4, 5, 6

## Can start now

Tasks 1 through 6, in that order.
