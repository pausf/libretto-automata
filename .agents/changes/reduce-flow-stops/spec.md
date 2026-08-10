# reduce-flow-stops — delta

Targets: payload

Two stops, not four. And a reviewer that fixes instead of reporting.

## Outcomes

**The flow stops in exactly two places, plus the last question.**

| Phase | Today | After |
|---|---|---|
| 1 · find-work | reports, waits | reports, **carries straight into phase 2** |
| 2–3 · write-spec | waits | **waits** — the user may want to change the contract |
| 5 · write-plan | waits | **waits** — the user may want to change the order |
| 6 · build-and-check | no wait | no wait |
| 6→7 · review-work | reports, no wait | **fixes**, no wait |
| 7 · present-work | reports, waits | reports, **carries straight into phase 8** |
| 8 · record-work | commits, then asks push + request | unchanged — **the last question** |

Three stops in total, and every one of them is a place where the user changes
something: the contract, the order, or whether the world sees it.

**The review seam acts on what it finds, without asking.** `review-work` relays the
reviewer's findings and then fixes every one of them — the same pass, no question, no
new round trip — re-runs the proofs the fix touches, and hands phase 7 both the
findings and what was done about each.

**Phase 7 still says everything it says today.** What disappears is the *stop* after it,
never a sentence of the report. The report and the commits land in the same turn.

## Scope boundaries

**In:** `docs/FLOW.md`, `commands/libretto-flow.md`, and the four skills whose stop or
standing changes — `find-work`, `review-work`, `present-work`, plus the payload spec
that promises four waits. One line in the `review-project` spec that justifies its own
no-fix rule by analogy to `review-work`, which no longer holds.

**Out:**

- **the reviewer becoming a writer.** `work-reviewer` stays read-only — `Read`, `Grep`,
  `Glob`, `Bash` and nothing else. Its independence is the whole point of the seam, and
  an agent that fixes what it finds is grading its own repair. The fixing happens in
  `review-work`, after the verdict is in and attributed.
- **`review-project` fixing anything.** It reviews somebody else's repository. Reporting
  is not a limitation there, it is the contract. Settled with the user 2026-08-10.
- **an unbounded fix loop.** One fix pass. A finding whose fix fails twice is a stopped
  item under `skills/evidence/`, reported to phase 7 as unfixed — never retried until it
  goes green.
- **asking about a finding.** Including the ones that look like product decisions. A
  finding cites a pillar or a proof by contract; if it cannot be fixed without a
  decision that is not ours, it is *reported* to phase 7 with the rest, not turned back
  into a question. The user reads it at phase 8, which is where they already were.
- **removing phase 1's artifact.** `proposal.md` and the branch still exist before phase
  1 reports. What goes is the wait after it, not the file.
- **a flag, a setting or a "quiet mode".** The stops are what they are for everyone. A
  toggle is a second contract to keep in sync, and the payload already rules it out for
  change size.

## Constraints

Every phase is still **invoked**, including the ones that decline — the phase owns its
own judgment and an orchestrator that pre-empts it has hidden a decision. Removing a
wait is not removing a phase, and the two are easy to conflate in prose.

The trivial lane already collapses to one turn and one question. This change makes the
spec'd lane *closer* to it without merging them: the spec'd lane keeps two stops the
trivial lane has none of.

`skills/evidence/` is unchanged and binds the fix pass exactly as it binds phase 6:
nothing reported that was not observed, no proof weakened to make a fix look done.

Wording lives in four places that must agree — the command, the flow doc, each skill,
and the payload spec. A stop removed in three of them and left in the fourth is worse
than not removing it, because whichever one the session happens to read wins.

## Prior decisions

- **The reviewer is `review-work`, not `review-project`.** Asked and answered by the
  user, 2026-08-10, before this spec was written. "Fix everything without asking" is
  coherent only for our own code.
- **The fix pass lives in `review-work`, not in a new phase and not in `build-and-check`.**
  The seam already has the findings in hand; routing them back through phase 6 was the
  round trip being removed. Not renumbering anything — the seam stayed unnumbered for
  the same reason.
- **Findings are still relayed attributed and unedited**, before the fixes. What was
  wrong is part of the record even when it is no longer true; a report that shows only
  the repaired state hides that the builder got it wrong.
- **Phase 8's question survives, and stays one question** — push and the pull request
  together, as already decided.
- **`work-reviewer` keeps its read-only tool grant.** Named here so the next session
  does not "helpfully" add `Edit` to close the loop inside the agent.
- **Ceiling named:** one fix pass, no re-review. The fixes are not themselves reviewed by
  a fresh context, so a fix that introduces a new defect is caught by the proofs or not
  at all. The replacement, the day that bites, is a second reviewer pass over the fix
  diff only — not a loop, a single bounded second look.

## Task breakdown

- [ ] 1 · `review-work` — add the fix pass: relay, then fix, then re-run the touched
      proofs; the two-failures stop; findings that cannot be fixed reported not asked
- [ ] 2 · `find-work` — the phase ends at the artifact and the reading, with no wait
- [ ] 3 · `present-work` — report, then carry into phase 8 in the same turn
- [ ] 4 · `commands/libretto-flow.md` — the two remaining waits, and the review seam's
      new standing
- [ ] 5 · `docs/FLOW.md` — the reasoning behind three stops instead of four-plus-one
- [ ] 6 · `.agents/specs/review-project/spec.md` — the analogy that no longer holds
- [ ] 7 · land the delta onto `.agents/specs/payload/spec.md` and delete the change

## Verification criteria

- every referenced skill still exists and every frontmatter still parses after the
  rewording
  Proof: scripts/check-payload
- no skill invokes a path that does not get installed — `review-work` gains a fix pass,
  not a tool outside the payload
  Proof: scripts/check-payload
- every `Proof:` citation in every spec still resolves, file and test name, after the
  payload and `review-project` specs move
  Proof: skills/record-work/spec-drift --anchors
- the Go suite is untouched by a payload-only change and still passes
  Proof: internal/link/apply_test.go TestApplyIsIdempotent

**What none of that verifies is the behaviour**, which is the whole change — a skill is
a prompt, and a prompt is checked by running it. Stated as what the next real run must
observe, not as criteria citing tests that cannot exist:

- phase 1 reports and phase 2 begins in the same turn, with `proposal.md` on disk before
  the report
- the spec stop and the plan stop both still happen
- a reviewer finding is fixed in the seam, with no question asked, and phase 7 carries
  both the original finding and its repair
- phase 7 reports and phase 8 commits in the same turn
- exactly one question is asked in the whole run after the plan: push and open the
  request
