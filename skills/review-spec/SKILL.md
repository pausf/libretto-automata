---
name: review-spec
description: "Trigger: a spec has just been written and a plan is about to be built on it; reviewing a capability spec or a change delta for ambiguity, untestable criteria, dead Governs: paths or a criterion that contradicts another capability. The seam between phases 3 and 5 — the only review in the flow that reads the contract instead of the code."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

**Every other review in this flow reads code. This one reads the promise the code will
be measured against**, in the seam between phase 3 writing the spec and phase 5 building
a plan on it.

That seam is where a spec is cheapest to fix and most expensive to leave. A vague
criterion costs one edit here. Past phase 5 it has become tasks; past phase 6 it has
become code; and at phase 7 the reviewer measures the work against a sentence nobody
could have failed, so the work passes and the promise was never real.

The five code lenses all run **after** the work exists. By then the contract is a fact.

## What it does not do

**It never asks whether the spec is a good idea.** That was phase 2's question and it was
answered. This reads only whether what is written can be built against and measured.

**It never rewrites the spec.** It reports; the phase-3 author decides. A reviewer that
edits the contract it is reviewing is the author twice over.

## Run the mechanical half first

Three of the questions below are lookups, not judgment, and a human reading for them
misses them. `record-work` ships the tool:

```
<skill-base>/../record-work/spec-drift --trace
```

It reports dead `Governs:` globs, criteria with no `Proof:` beneath them, and code no
spec claims. Read its output before reading the spec — findings it already names do not
need to be found again, and the ones it cannot name are the ones this skill is for.

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
- which of the four it is — ambiguity, testability, governs, contradiction
- **for ambiguity, both readings.** For testability, the failing run that does not exist.
- what would settle it — never the settled version. That is the author's to write.

End with what was read and found clean. A review that lists only problems reads as a
review that found everything wrong with everything, and the phase-3 author cannot tell
what was checked from what was skipped.
