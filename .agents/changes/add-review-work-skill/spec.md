# Independent review — the work is looked at by someone who did not write it

Targets: payload

The flow's own rule is that nothing is true until observed — and phase 7 is a
self-report by the same agent that wrote the code. This delta adds the missing
observer: between build (6) and present (7), a **fresh subagent with none of the
session's context** reads the contract and the diff, re-runs the proofs, and returns
findings. Phase 7 then presents the work *including* that verdict.

It closes the open box in the payload spec — *"an independent verifier: check the
implementation against the spec's criteria, never run by whoever wrote the code."*

## Outcomes

- After phase 6 finishes and before phase 7 reports, `review-work` runs: it assembles
  what the reviewer needs — the change's spec delta, the capability specs it targets,
  and the diff of the change's branch against its base — and launches one
  `work-reviewer` subagent.
- The reviewer starts from zero. It sees the spec and the diff, never the
  conversation. What the session believes about the work is exactly what it must not
  inherit.
- The reviewer **re-runs every `Proof:` named by the criteria the change touches**,
  in the foreground, and reads the result. A proof is *verified* only when the
  reviewer observed it pass; phase 6 having reported it green is not evidence here —
  trusting the builder's report is the failure this skill exists to remove.
- The reviewer checks the diff against the spec's own terms: each outcome has code
  behind it, each scope boundary held, no criterion is satisfied only on paper.
- Findings come back as the subagent's return value — one line each, with where and
  what — and phase 7 presents them **attributed to the reviewer**, next to the
  builder's own report.
- **Findings never block.** Phase 8 remains the user's decision, taken with the
  verdict in front of them.
- A change with no spec gets no review: nothing to review against. The skill declines
  in one line, like any phase with nothing to do, and the trivial lane stays
  collapsed.

## Scope boundaries

**In:** the `review-work` skill, the `work-reviewer` agent, their wiring into
`commands/libretto-flow.md`, the seam documented in `docs/FLOW.md`, and one
paragraph in `skills/present-work/SKILL.md` so phase 7 carries the verdict when
reached by its own trigger, outside `libretto-flow`. That last file was missing from
this list; the reviewer's first real run found the gap — outcome promised, phase-7
skill unchanged — and the scope learned it here.

**Out:**

- **a blocking gate.** Decided in phase 4 and consistent with the payload's standing
  rule — drift detection warns and never blocks. A reviewer that can stop the flow in
  someone else's project is a reviewer that gets deleted.
- **a fix loop.** The reviewer reports; it edits nothing, and no correction round
  runs inside this skill. Acting on findings is the user's call, and the fix is a new
  pass through phase 6.
- **general code review.** Style, performance, taste — out. The reviewer checks the
  work against *its contract*. A finding cites a pillar or a proof, not an opinion.
- **re-reviewing after a fix, automatically.** Running review again is the
  orchestrator invoking the skill again, not machinery inside it.
- **renumbering the flow.** The phases stay 1–8; this runs in the seam between 6
  and 7. Independence comes from the fresh context, not from a number.
- FLOW.md's open question — where the *rendered artifact* gets looked at — stays
  open. This reviewer reads specs, diffs and test output, not pixels.

## Constraints

- Everything the payload spec already imposes: frontmatter `name:` equals the
  directory or filename, the skill invokes only what gets installed, no path from
  this repository named as though the user's project had it, one author per file.
- **The reviewer's independence is structural, not requested.** It is a separate
  agent file launched fresh — like `spec-writer` — not a promise inside the
  orchestrator's own context.
- **Sub-agents start with no rules loaded**, so `work-reviewer` carries its own:
  evidence first, nothing reported that was not observed, no commit, no push, no
  writes anywhere. Unlike `spec-writer` it gets a shell — running the proofs is its
  job — which makes "write nothing" a rule it keeps rather than one that is enforced.
  The agent file says so, the way `spec-writer` says it about its write scope.
- The reviewer finds the proofs the way `spec-drift --anchors` resolves them: the
  spec's citation names the file and the test. It runs what is named; it does not
  invent a broader suite to run.
- Both new files pass `scripts/check-payload`; the skill carries the same
  frontmatter contract as the other seven (`license`, `author`, `version`).

## Prior decisions

- **Report, never block** — user, 2026-08-07, phase 4 of this change. The stops in
  this flow exist so the *user* can say no; a machine saying it introduces the first
  automatic gate the payload has deliberately never had.
- **A new skill in the 6→7 seam, no renumbering** — user, 2026-08-07, same round.
  Confirmed explicitly: the skill launches a subagent that starts from zero, outside
  the session's context window.
- **The reviewer re-runs proofs rather than trusting phase 6's report.** From
  `skills/evidence/`: an observation made by someone else, relayed, is a claim. The
  reviewer's entire value is that its observations are its own.
- **No spec, no review.** The trivial lane collapsed the ceremony because there was
  no contract to disagree with; a reviewer with no contract has nothing to check.
  Declining is said in one line, per "a phase is invoked even when it has nothing to
  do".
- **The orchestrator relays, the reviewer never writes.** Same single-author rule as
  the fan-out: findings travel in the return value, and phase 7 carries them into the
  report.

## Task breakdown

- [ ] `agents/work-reviewer.md` — the fresh reviewer: inputs (spec paths, diff
      range), its rules, what it returns
- [ ] `skills/review-work/SKILL.md` — the seam: decide applicability, assemble
      inputs, launch, relay
- [ ] `commands/libretto-flow.md` — invoke `review-work` between phases 6 and 7
- [ ] `docs/FLOW.md` — document the seam and mark the open verifier box answered
- [ ] delta applied onto `.agents/specs/payload/spec.md` at landing — the open task
      box checked, these decisions carried in

## Verification criteria

- frontmatter parses, `name:` matches directory and filename, and every skill
  `libretto-flow` references — now including `review-work` — exists
  Proof: scripts/check-payload
- neither new file invokes a path that does not get installed
  Proof: scripts/check-payload
- the payload's existing anchors still resolve after the delta lands
  Proof: skills/record-work/spec-drift --anchors

**What the static checks cannot prove**, stated as observations this change must
produce rather than as citations to tests that cannot exist — a skill is a prompt,
and a prompt is checked by running it:

- a run of the flow where `review-work` launches a reviewer that saw no session
  context and its findings appear in phase 7's report, attributed
- a run against a change whose `Proof:` names a test that does not pass, where the
  reviewer catches it by running it, not by reading phase 6's claim
- a change with no spec, where the skill declines in one line and adds no wait
