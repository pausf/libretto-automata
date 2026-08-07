# Plan — add-review-work-skill

**Goal:** the work is looked at by someone who did not write it — a fresh
`work-reviewer` subagent in the 6→7 seam, reporting and never blocking.

Derived from `spec.md` in this folder (`Targets: payload`). One writer: the
orchestrator marks the boxes; sub-agents report.

Every task closes against the spec's verification criteria: the two static ones
(`scripts/check-payload`, `spec-drift --anchors`) plus the observation it must
produce, named per task.

## Tasks

- [x] **1 · `agents/work-reviewer.md`** — done, commit 07f3969, six gates green. — the fresh reviewer. Inputs it is handed
      (spec delta path, capability spec paths, diff range), its own rules carried
      explicitly (evidence first, re-run every named `Proof:` in the foreground,
      no commit, no push, no writes anywhere — a shell it keeps read-only by rule),
      and what it returns: findings one line each citing a pillar or a proof, or an
      explicit "nothing found". Waits on: nothing.
      Closes: check-payload passes on the new agent; frontmatter `name:` matches
      filename.

- [x] **2 · `skills/review-work/SKILL.md`** — done, commit 122540c, six gates green. — the seam. Step 0 applicability (no
      spec → decline in one line, no wait), assemble the reviewer's inputs, launch
      exactly one `work-reviewer`, relay its findings to phase 7 attributed and
      unedited. Frontmatter contract like the other seven skills (`license`,
      `author`, `version`). Waits on: 1 (it names the agent it launches).
      Closes: check-payload passes; every skill it references exists.

- [x] **3 · wire `commands/libretto-flow.md`** — done, commit 384a24b, six gates green. — invoke `review-work` between
      phases 6 and 7, with the same invoked-even-when-empty rule as every phase.
      Waits on: 2.
      Closes: check-payload — the command's references all resolve.

- [x] **4 · document the seam in `docs/FLOW.md`** — done, commit 5f456ed, six gates green. — the reviewer between build and
      present, report-never-block, and the note that the open "where does the
      artifact get looked at" question stays open (this reviews contracts, not
      pixels). Waits on: 2 (documents what exists, not what is hoped).
      Closes: spec-drift --anchors still resolves 169+ citations.

- [ ] **5 · run it once for real** — the spec's three observations: a reviewer that
      saw no session context returns findings into phase 7's report; a planted
      failing `Proof:` is caught by running it; a specless change gets a one-line
      decline. Waits on: 3, 4.
      Closes: the three observations recorded in the spec delta, as observations.

- [ ] **6 · landing** — phase 8 applies the delta onto
      `.agents/specs/payload/spec.md` (open verifier box checked, decisions carried
      in) and deletes this folder, same commit as the final code. Waits on: 5.
      Owner: phase 8, not this plan's executor.

## Can start now

Task 1. Tasks 3–4 are independent of each other once 2 lands.
