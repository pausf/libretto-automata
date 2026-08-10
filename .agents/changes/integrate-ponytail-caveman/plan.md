# Plan: integrate-ponytail-caveman

Spec: `.agents/changes/integrate-ponytail-caveman/spec.md` (targets `payload`).
One writer: the orchestrator. Sub-agents report; boxes are marked here, when a task
is verified, not at the end.

Every commit passes all six gates — task 6 is the final whole-payload run, not the
first time the gates see the work.

## Tasks

- [x] 1. vendor `skills/ponytail/` and `skills/ponytail-debt/` from
      DietrichGebert/ponytail at the current upstream commit, unmodified, plus
      `LICENSE-ponytail`. Verify self-sufficiency: no reference outside each skill's
      directory — if `ponytail-debt` reaches outside itself, it is dropped, never
      patched.
      Waits on: nothing.
      Closes: spec criterion "frontmatter parses and name equals directory"
      (scripts/check-payload) for the ponytail pair.

- [x] 2. vendor `skills/caveman/` and `skills/caveman-commit/` from
      JuliusBrussee/caveman at the current upstream commit, unmodified, plus
      `LICENSE-caveman`. Same self-sufficiency check, same drop-never-patch rule.
      Waits on: nothing. Independent of task 1.
      Closes: the same criterion for the caveman pair.

- [x] 3. THIRD-PARTY.md: move both from *Not vendored* to *Vendored* with pinned
      version/commit and licence lines; rewrite the *Not vendored* rationale and
      *Naming* to state the reversal and why (target user starts from zero;
      namespacing already answers the collision).
      Waits on: 1 and 2 (the pinned commits are known only after the copies exist).

- [x] 4. calling skills: `write-spec` (`ponytail:ponytail` → `ponytail`, ledger
      reference), `build-and-check`, `present-work`, `record-work`
      (`caveman-commit`), `commands/libretto-flow.md` — conditional "if installed"
      prose reconciled to shipped-by-default, adaptations only in callers, never in
      the copies.
      Waits on: 1 and 2 (check-payload fails on a reference to a skill not yet on
      disk).
      Closes: spec criterion "every skill the flow references exists"
      (scripts/check-payload).

- [x] 5. docs: `docs/FLOW.md` "present rather than vendored" paragraph and README's
      companions section state the new decision.
      Waits on: 3 (the docs cite the decision THIRD-PARTY.md records).

- [x] 6. all six gates green over the enlarged payload, including
      `spec-drift --anchors` after the prose edits.
      Waits on: 1–5.
      Closes: spec criteria "no vendored skill invokes a path that does not get
      installed" and "the payload spec's own anchors still resolve".

## Can start now

Tasks 1 and 2, independently. Everything else is downstream.
