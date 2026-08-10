---
name: review-work
description: "Trigger: phase 6 finished and the work has a spec; before present-work reports; asking whether the work matches its contract. The seam between build and present. Launches one fresh work-reviewer subagent that re-runs the proofs, relays what it found attributed and unedited, then fixes every finding without asking."
license: MIT
metadata:
  author: pausf
  version: "1.1"
---

## What this does

The seam between phases 6 and 7 of the Libretto flow: hand the finished work to
someone who did not write it.

This skill does not review anything itself — it cannot. Whoever is reading this
carried the session that produced the code, and what that session believes about the
work is precisely the claim under review. The review runs in a **fresh `work-reviewer`
subagent** that starts with none of this context. This skill decides whether a review
applies, assembles what the reviewer needs, launches exactly one, relays what comes
back — and then fixes it.

**The reviewer reads; this skill writes.** The split is the value of the seam: an agent
that repairs what it found is grading its own repair, and a verdict is only independent
while the hands that earned it are not the ones being judged.

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

**Relay before repairing, always.** By the time phase 7 reports, most of these will
already be fixed — and a report showing only the repaired state hides that the builder
got it wrong. What was wrong is part of the record even when it is no longer true.

## Step 4 — Fix every one of them, and ask nothing

The findings do not go to the user as a question. They get fixed, here, in this pass.

Sending a defect back to the user to authorise its own repair is a round trip that buys
nothing: a finding cites a pillar or a proof by contract, so it is a defect against a
contract the user already agreed to. There is no version of "yes, leave it broken" worth
a stop.

Per finding, in the order they came back:

1. fix it — `skills/build-and-check/` governs how much, the same as any phase 6 work
2. re-run **the proofs that finding touched**, in the foreground, and read the output.
   Not the whole suite unless the fix reaches that far
3. record what changed, one line, for phase 7 to carry beside the finding

**Two failed attempts on one finding stops that finding.** Not a third try. It goes to
phase 7 as *found, not fixed*, with what failed and what is still unknown — the stopped-item
rule from `skills/evidence/`, unchanged. One stopped finding does not stop the others.

**A finding that cannot be fixed without a decision that is not ours is reported, not
asked.** A product tradeoff, a boundary the spec never drew — it goes to phase 7 with the
rest and the user meets it at phase 8, which is where they already were. Turning it into a
question here reintroduces the stop this seam exists without.

**One fix pass. No re-review.** The fixes are not themselves seen by a fresh context, so a
fix that introduces a new defect is caught by the proofs or not at all. Named as a ceiling
rather than hidden: the replacement, the day it bites, is one bounded second look at the
fix diff — never a loop.

Nothing here weakens a proof to make a fix look done. `skills/evidence/` binds this pass
exactly as it binds phase 6.

## Output

One line when declining. Otherwise, in this order:

- the reviewer's proof verdicts and findings, relayed as returned
- what was done about each: fixed with its re-run, or stopped with what failed, or
  reported-not-fixed because the decision was not ours

Then carry straight into phase 7. Presenting is phase 7's job, but it is not another
turn — nothing here is waiting on the user.
