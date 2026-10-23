# The fan-out — extra capabilities, when context cannot hold them

Read this when the change amends more than one capability. It covers writing their
`spec-<capability>.md` files and, only when this session's context cannot hold the research,
launching one `spec-writer` per extra capability.

**A change amending more than one capability** carries the primary capability's contract
in `proposal.md` and one `spec-<capability>.md` beside it per extra capability, each with
its own `Targets:` line and its own six pillars. The decision log, the approach, the
`Durable decisions:` line and the checklist exist once, in `proposal.md`. With room in
context, the session writes those files inline as well.

**The context is the only reason to fan out.** Launch a writer for an extra capability
only when this session's context cannot hold its research — ten subtasks researched in one
thread will exhaust it, which is the exact problem fanning out exists to solve. Fanning out
because the answer is not yet obvious turns a reflex into a research project — if two
readings of the code already settle it, settle it and move on.

The risk a fan-out creates is a set of specs that each make sense alone and contradict
each other — different names for one concept, incompatible assumptions, the same
decision made twice two ways. Consistency has to come from somewhere other than a
single author.

**It comes from `proposal.md`, written first.** Its decision log, the primary contract and
**the vocabulary** for every concept the specs will share are the ground every writer
reads, so two writers cannot give one thing two names. No separate brief and no separate
log are written beside it: the change document already is both, and a second file about
the same change is a second source of truth.

## Do not change the model or the effort part-way through a fan-out

A phase's context is billed cheaply while its prefix stays identical between calls —
the run is charged a fraction of the input price for everything the previous call
already established. **Switching the model or the effort level invalidates that prefix
and rebills the whole context at full input price**, and a fan-out is where that hurts
most: N writers is N contexts, and a switch part-way through pays for all of them
twice.

So pick the tier before the fan-out starts and leave it alone until the phase ends.
Between phases is free — that is a new prefix either way.

**This is a statement, not a check.** The dial is the session's, and no skill can read
it or stop a hand from moving it; what a skill can do is put the cost where somebody
about to launch eight writers will see it.

## The bridge runs one way

```
session  ──writes──▶ proposal.md ──read-only──▶ N writers
writer   ──writes──▶ its own spec-<capability>.md   (one author per file)
writer   ──returns──▶ findings in its reply           (never to a file)
```

**Writers never write to `proposal.md`.** N writers on one file is the lost-update race
that `skills/write-tasks/` exists to avoid — two reads of the same content, two writes,
the second silently discarding the first. Findings come back in the return value.

Every file in this phase has exactly one author: the session writes the change document,
and each `spec-<capability>.md` is written by the session or by one writer.

## What each writer is given

Each one is launched as the `spec-writer` agent — that contract exists for exactly this
seat, and a generic sub-agent would start without it. In its prompt:

- **the path to `proposal.md`** — its log, its contract, and always the vocabulary. The
  log carries the user's answers verbatim, which the writer carries into its *Prior
  decisions* pillar as they stand. The vocabulary is never optional, because it is the
  entire mechanism stopping two writers from giving one concept two names
- **the path, never the text.** Excerpting is what puts N copies in N contexts
- its own capability — which one, and what of the request it covers
- the boundary between its spec and its siblings', stated explicitly
- the one path it may write, and nothing else

And the standing rules, because **a sub-agent starts with none of them**:

- invoke `evidence` before anything: nothing reported that was not observed, no gate
  in the background, no test weakened to go green
- **do not commit.** The session commits, because only it sees the whole task
- **do not push.** Ever. Pushing is a question for the user, and a sub-agent cannot
  ask it
- **do not touch `proposal.md`** or a sibling's spec
- read-only everywhere except the one file assigned

## When a writer gets stuck

There is no channel while it runs. It cannot ask the user — it does not have the
conversation — and it cannot ask the session, because the return value happens once, at
the end.

So the rule has to be explicit, or the default takes over: **a model with no
instruction guesses plausibly**, which is the exact failure this flow exists to
prevent.

When a writer hits something it cannot legitimately decide:

1. **Write everything it can.** A spec with five solid pillars and one open question
   is useful. Nothing is not.
2. **Mark the gap in the file itself** with `[NEEDS CLARIFICATION: question]`, where the
   answer belongs. Not a placeholder that reads like an answer.
3. **Return the question**, phrased so the session can ask it without having to
   reconstruct the context.
4. **Never invent past it.** Not a default, not the more common of two precedents,
   not "the obvious choice". A guess that looks like a decision is worse than a hole,
   because a hole gets filled and a guess gets built on.

Questions funnel through the session, always. If writers could ask directly, several
would interleave questions from work the user cannot see, and each answer would arrive
without the context that made it a question.

## Resolving the markers

Read every returned spec. Where a writer left `[NEEDS CLARIFICATION: question]`, its
inputs did not say enough, and the resolution is the interview's tail:

1. Ask each marker's question — same shape as step 4, one at a time, a reasoned
   recommendation first.
2. Log the answer under `## Decision log`, verbatim, like every other answer.
3. **Replace the bracket expression — the brackets and everything inside them, nothing
   outside them — with the logged answer.** That is the one declared exception to
   one-author-per-file, and it is scoped to exactly that expression: an answer that
   changes anything beyond it is a relaunch of the writer with the updated log, not a
   patch.

**Under `/libretto-attacca` a marker takes the recommended answer**, logged with the
`(assumed)` suffix and naming what changes if it is wrong — the same rule as the
questions that were never asked.

## Writers do not talk to each other

By design — no shared mutable state, no negotiation, no ordering to get wrong.

The consequence is that a boundary drawn badly between two specs is not fixed by the
two of them noticing. It surfaces as two returns describing overlapping ground, and
**the session resolves it after the fact**: one decision, applied to both files,
recorded once.

Both writers were right about their own scope. `proposal.md` was wrong about the line
between them, and that is a finding, not a failure.

## What comes back

The path it wrote, and a short list of what it found that `proposal.md` did not settle —
a convention it missed, a constraint that contradicts it, a question only the user can
settle.

**Read every one of those returns before accepting the set.** They are the interesting
part: a writer that found the document wrong found something real. When two returns
disagree about the same thing, that disagreement is a decision, and it is resolved here —
once, then applied to both files. Never left as two specs each quietly assuming it won.

**Silence is not success.** A writer that died, returned nothing, or came back empty did
not finish. Do not accept its spec because it failed to complain. A returned spec with no
markers and no findings is still read before it is accepted.

If `proposal.md` turns out to be wrong on something structural, correct it and re-run the
affected specs. Patching the files individually reproduces the drift a shared document
was there to prevent.
