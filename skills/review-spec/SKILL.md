---
name: review-spec
description: "Trigger: a change document's checklist has just been written and code is about to start; reviewing a capability spec or a change's contract for ambiguity, untestable criteria, dead Governs: paths or a criterion that contradicts another capability. Launches one fresh spec-reviewer, relays it unedited, then asks the flow's one stop before code — the only review in the flow that reads the contract instead of the code."
license: MIT
metadata:
  author: pausf
  version: "2.1"
---

## What this does

**Every other review in this flow reads code. This one reads the promise the code will
be measured against**, once the change document is whole — contract, approach and
checklist — and before the first line of code.

That is where a contract is cheapest to fix and most expensive to leave. A vague
criterion costs one edit here. Past phase 6 it has become code; and at phase 7 the
reviewer measures the work against a sentence nobody could have failed, so the work
passes and the promise was never real.

The five code lenses all run **after** the work exists. By then the contract is a fact.

**Only on the contract lane.** A change step 0 sent down the trivial lane has no contract
to read, and running this there is ceremony on a one-line change.

## Launch exactly one fresh reviewer

The launch, the relay and the stop below are **the session's conduct** — a
`spec-reviewer` invoking this skill for the reading starts at *What it does not do* and
follows nothing above that heading.

**Launch exactly one `spec-reviewer`**, foreground, given by path the change document, every
`spec-<capability>.md` beside it, and the capability spec(s) their `Targets:` lines name —
and none of the conversation. Not a summary of what was decided, not which findings would
be welcome: telling the reviewer what to expect is priming the witness. A change begun in
the old layout is given its delta, `plan.md` and `tasks.md` in their place.

The session that wrote the document cannot read it cold. It remembers what it *considered*
as vividly as what it *chose*, and reads every gap as filled. A reviewer with no memory
can only read what is on the page — which makes it the only check this flow has that the
document says enough to be built from, box by box.

One, not a panel — the same answer `review-work` gives. **Silence is not a clean review**:
a reviewer that died or returned nothing did not read the document. Launch it again, once.

## Relay it, unedited

Its findings go to the user in its words, attributed, before the question — never
summarised, never filtered for the ones the session agrees with. A finding the session
thinks is wrong is relayed anyway, with the session's reading beside it; the user decides
between them at the stop.

## The stop — the one before code

Then stop. **This is the flow's only stop before code**, and it carries everything phase 2,
phase 5 and the 5→6 seam produced: the contract, the approach and the checklist, with the
reviewer's findings on the table while the agreement is still being made.

**Ask it with `AskUserQuestion` — the host's native prompt, or in conversation where there
is none — never as a sentence at the end of the report.** Carry on to phase 6 —
recommended, naming the first box — change the document first, or stop here. **And room
to answer differently**, because the reason to reject a document is rarely one of three
things somebody guessed in advance.

Same rule as every stop. The reason lives once, in `skills/record-work/`.

Under `/libretto-attacca` the stop is answered by the invocation, and the run that wrote
the document acts on the findings before phase 6 — the command owns that rule.

## What it does not do

**It never asks whether the spec is a good idea.** That was phase 2's question and it was
answered. This reads only whether what is written can be built against and measured.

**It never rewrites the spec.** It reports; the session that wrote the document decides. A
reviewer that edits the contract it is reviewing is the author twice over.

## Run the mechanical half first

Three of the questions below are lookups, not judgment, and a human reading for them
misses them. `record-work` ships the tool:

```
<skill-base>/../record-work/spec-drift --trace
<skill-base>/../record-work/spec-drift --ears
```

`--trace` reports dead `Governs:` globs, criteria with no `Proof:` beneath them, code
no spec claims, and joined criteria — two or more `shall` resting on one `Proof:` with no
ceiling named. A joined criterion is a prompt to split the criterion or name its ceiling,
never a finding on its own: whether the one test holds both promises is a reading.

`--ears` reports criteria that cannot be failed because they were written in prose rather
than in one of the five EARS patterns — hard on a change delta, a warning on a capability
spec that predates the syntax.

Read both before reading the spec — findings they already name do not need to be found
again, and the ones they cannot name are the ones this skill is for.

**`--ears` checks for `shall` and nothing else, which leaves this skill the whole
question.** A criterion can carry the keyword and still be unfailable: *"When the user
asks, the system shall behave correctly"* passes the gate and means nothing. The gate
buys the shape; whether the response is concrete is the reading below.

**Every finding below is a labelled judgment call.** A spec is a document a person wrote
for other people; a rule enforced as law here produces a spec written to satisfy the
checker.

## Ambiguity — can two readers build different things?

The test is not "is this unclear". It is: **could two competent people implement this
sentence and disagree about who was right?** If yes, name both readings — a finding that
says "ambiguous" without showing the fork is a finding the author cannot act on.

The words that carry it, and what they hide:

| Word | The question it is not answering |
|---|---|
| *fast*, *quick*, *responsive* | how fast, measured with what |
| *properly*, *correctly*, *gracefully* | correct by whose definition |
| *should*, *ideally*, *where possible* | is this a promise or a preference |
| *handle*, *support*, *manage* | do what, on which input |
| *etc.*, *and so on*, *among others* | the list is the contract; a trailing *etc.* is a blank cheque |

A criterion containing **and** is often two criteria that will be half-met and reported
as done. Split it or say why it is atomic.

## Testability — could this fail?

**The strongest question in this skill, and the cheapest to ask: what would a failing
run of this criterion look like?** If no concrete input produces a red result, the
criterion cannot be verified — only asserted — and it will be marked done by whoever
reads it most generously.

Three shapes that cannot fail:

- **describing structure, not behaviour** — "the parser is modular" is unfalsifiable;
  "the parser rejects a trailing comma" is one test.
- **naming an internal instead of an effect** — "the cache is used" is satisfied by a
  call that does nothing; "a second identical request issues no network call" is not.
- **a promise about the future** — "will scale", "is extensible". Nothing today can fail it.

A criterion whose `Proof:` names a test that only asserts the code ran — no input, no
expected value, an error swallowed — is untestable with a citation attached, which is
worse than none: it reads as proven. **`--trace` cannot see this and will not tell you.**
Read the test the citation names.

## `Governs:` — does the boundary describe reality?

`--trace` catches a glob matching nothing. It cannot catch these:

- **two capabilities claiming one path.** Both specs now own it, so both are the place
  a change gets recorded, so it gets recorded in neither. Say which one owns it.
- **a glob wider than the capability.** `internal/**` on a spec about one package makes
  every unrelated change look like drift, and a drift check that cries wolf is muted.
- **the capability's real code sitting outside its globs.** The spec is intact and
  governs nothing that matters.

## Contradiction across capabilities

The one finding that cannot be seen from inside a single spec, and the reason this is
read against the others rather than alone: **a criterion that another capability's
criterion forbids.** Both are true in their own document, the conflict surfaces in the
code that has to satisfy both, and by then it is a bug with two specs defending it.

Same for silent duplication — the same promise stated in two capabilities drifts, and
the copy nobody edited is the one that reads as authoritative.

## The checklist — could each box be built from this?

The question a cutter used to ask, now asked by the one cold reader left: **for each box,
what do the contract and the approach together fail to answer for building it?** A task
implied with no mechanism described, a criterion with no proof saying what closes it, two
statements that cannot both hold, a boundary that does not say which side something falls
on.

Name the gap and the part of the document that should have answered it, never the answer.
A gap goes back to the contract or the approach; it never becomes a box that says "figure
out X".

## Proportion

**This skill can generate infinite findings and that is its failure mode, not its
strength.** A spec is prose; every sentence in it can be read as slightly imprecise, and
a review that returns thirty findings gets skimmed and dismissed entire.

The filter: **would this finding change what somebody builds?** If the ambiguity resolves
the same way under every plausible reading, it is not a finding. Report what forks the
work, and say plainly when the answer is nothing.

## Output

Group by spec file. Per finding, in one or two lines:

- the criterion, quoted
- which of the five it is — ambiguity, testability, governs, contradiction, or a box the
  document cannot build
- **for ambiguity, both readings.** For testability, the failing run that does not exist.
- what would settle it — never the settled version. That is the author's to write.

End with what was read and found clean. A review that lists only problems reads as a
review that found everything wrong with everything, and the document's author cannot tell
what was checked from what was skipped.
