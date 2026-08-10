# Review Project — delta

Targets: review-project
Governs: skills/review-project/** agents/review-security.md agents/review-design.md agents/review-reliability.md agents/review-tests.md agents/review-intent.md

`review-lens` becomes four agents. Everything else in the review-project spec stands —
five lenses, one frozen diff, reports never blocks, briefs are the diversity.

## Outcomes

`agents/review-lens.md` is replaced by four files, one per lens:

```
agents/review-security.md     agents/review-design.md
agents/review-reliability.md  agents/review-tests.md
```

- **Each carries its own `model:`**, which is the whole reason for the split. The four
  are no longer one thing with a parameter: design and tests can run cheap while
  security does not.
- Each declares the same tool grant `review-lens` declared — `Read, Grep, Glob, Skill`
  — and applies exactly one skill, named in its own body rather than passed in a brief.
- **`review-intent` is unchanged** except that it may now declare a `model:` too. It
  was already its own file and already the only lens carrying `Bash`.
- The launch site in `skills/review-project/SKILL.md` names the four agents instead of
  launching one agent four times. The four skill names stay written out in the
  `Skill(skill="…")` form the static reference check can see — that requirement does
  not move, it just moves house.
- Five fresh subagents in parallel, one brief each, none carrying the conversation.
  Unchanged.

## Scope boundaries

**Out:** changing what any lens looks for, what it returns, or the order they run in.
This delta changes how many files the lenses live in and nothing about the review.

## Constraints

- `scripts/check-payload:62-67` requires each agent's `name:` to match its filename.
  Four new names, four new filenames, and the old one gone.
- Any reference to `review-lens` anywhere in the payload dies with the file. A dangling
  agent name fails the reachability check rather than failing silently at runtime as a
  lens that never ran, which is the failure mode the check exists for.
- `docs/FLOW.md` and the review-project spec's own `Governs:` name the agent; both move.

## Prior decisions

- **This reverses a recorded decision, and on purpose.** "One `review-lens` agent, not
  four" (`skills/review-project/SKILL.md:277`) rested on the four lenses differing in
  exactly one thing. A per-lens model is a second thing. The premise is what changed;
  the reasoning was never wrong.
- **The cost is accepted, not hidden:** four near-identical bodies that can drift apart.
  The ceiling — if they do drift, or if a fifth lens arrives, the answer is to generate
  them from one source at build time, not to hand-maintain five copies.

## Task breakdown

8. Split `agents/review-lens.md` into the four lens agents; update the launch site in
   `skills/review-project/SKILL.md`, `docs/FLOW.md`, and every remaining reference.

## Verification criteria

- frontmatter parses, each new agent's `name:` matches its filename, every referenced
  skill resolves, and no reference to the deleted `review-lens` survives
  Proof: scripts/check-payload
- the four lens skills are still named in the form the static check can see
  Proof: scripts/check-payload
- every `Proof:` citation in the amended specs resolves
  Proof: skills/record-work/spec-drift --anchors
- behaviour is a prompt and is checked by running it. **Owed before this is fact:** one
  real review launched with the four split agents, five reports relayed unmerged, and
  the token cost measured against the 307k baseline with design and tests on a cheaper
  model. Until that run, the saving is a prediction — which is the same standing debt
  the review-project spec already carries for its tool grants.
