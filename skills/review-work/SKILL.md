---
name: review-work
description: "Trigger: phase 6 finished and the work has a spec; before present-work reports; asking whether the work matches its contract. The seam between build and present. Launches one fresh work-reviewer subagent that re-runs the proofs and returns findings — it reports, it never blocks."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

The seam between phases 6 and 7 of the Libretto flow: hand the finished work to
someone who did not write it.

This skill does not review anything itself — it cannot. Whoever is reading this
carried the session that produced the code, and what that session believes about the
work is precisely the claim under review. The review runs in a **fresh `work-reviewer`
subagent** that starts with none of this context. This skill decides whether a review
applies, assembles what the reviewer needs, launches exactly one, and relays what
comes back.

`skills/evidence/` governs here too: the reviewer's verdict is relayed as observed,
never summarised into something kinder.

## Step 0 — Is there anything to review against?

**No spec, no review.** A change phase 2 sent down the trivial lane has no contract,
and a reviewer with no contract has nothing to check — what would a finding cite?

Decline in one line — *review: no spec, nothing to review against* — and add no
wait. The trivial lane stays collapsed; this skill never becomes the ceremony that
phase 2 just waived.

## Step 1 — Assemble the reviewer's world

The reviewer gets three things, all paths and ranges, no narration:

- **the contract**: the change's spec delta (`.agents/changes/<change>/spec.md` or
  the project's own layout), and the capability spec(s) its `Targets:` names
- **the work**: the diff range, `<base>...<branch>` — read the base off the
  repository rather than assuming it
- **the folder**: the change directory, for the proposal and the plan if the
  reviewer wants the stated intent

What it never gets: this conversation, the session's summary of what was built, or
any hint of which findings would be welcome. Telling the reviewer what to expect is
priming the witness.

## Step 2 — Launch exactly one

One `work-reviewer`, foreground, and its return value read in full. The agent file
carries its own rules — evidence first, re-run every proof the change touches, no
edits, no commit, no push.

One, not a panel. A second reviewer costs a second full pass and buys correlation
dressed as corroboration — two runs of the same model over the same diff are not two
opinions. The day that changes is the day the payload decides reviews need lenses,
and that is a spec change, not a quiet doubling here.

**Silence is not a clean review.** A subagent that died or returned nothing did not
review the work. Say so and launch it again — once. A reviewer that cannot finish
twice is a stopped item, per `skills/evidence/`.

## Step 3 — Relay, attributed and unedited

The findings go into phase 7's report in their own section, attributed to the
reviewer, next to the builder's own account:

- each proof the reviewer ran, pass or fail, one line each
- each finding as returned — never softened, never merged into the builder's prose,
  never filtered by "that one is minor"
- or the reviewer's explicit **"nothing found"**, which is a statement, not an
  absence

**Findings never block.** Phase 8 remains the user's decision, taken with the
verdict in front of them. What a finding changes is what the user knows, not what
the user may do. Acting on one is a new pass through phase 6, not a loop inside
this skill.

## Output

One line when declining. Otherwise: the reviewer's proof verdicts and findings,
relayed as returned, ready for phase 7 to carry.

Then stop. Presenting is phase 7's, and the decision the findings inform is the
user's.
