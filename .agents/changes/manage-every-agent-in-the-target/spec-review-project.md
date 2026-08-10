# Review Project — delta

Targets: review-project
Governs: skills/review-project/** skills/review-security/** skills/review-design/** skills/review-reliability/** skills/review-tests/** agents/review-lens-security.md agents/review-lens-design.md agents/review-lens-reliability.md agents/review-lens-tests.md agents/review-lens-intent.md commands/libretto-review.md

The five lens agents are namespaced `review-lens-*`. Nothing else about the review
changes.

## Outcomes

```
agents/review-lens-security.md      agents/review-lens-design.md
agents/review-lens-reliability.md   agents/review-lens-tests.md
agents/review-lens-intent.md
```

- **The agents are renamed; the skills they apply are not.** `review-lens-design`
  invokes `Skill(skill="review-design")`. They live in different directories, only
  the agents collided, and renaming a skill nobody asked about would be scope
  arriving unannounced.
- Every reference in the payload and the code moves with them: the launch table in
  `skills/review-project/SKILL.md`, and the `Governs:` line above. The capability
  spec keeps its old names until phase 8 folds this delta onto it — it describes
  `main`, and `main` still has the old files.
- The four skill names stay written out in the `Skill(skill="…")` form
  `scripts/check-payload` scans. That requirement has now moved house twice and has
  not changed.

## What this fixes

`install --global` on the machine that reported it:

```
conflict     agents/review-reliability.md
```

The user already had an agent of that name — a different agent, their own R3
reliability reviewer. `install` reported it and refused to touch it, which is correct
and is also permanent: **a name we do not own is a name we can never install under.**
The lens was unreachable in the global target and would have stayed that way.

## Scope boundaries

**In:** the five agent filenames, their `name:` frontmatter, and every reference.

**In, and worth naming because a reviewer could not tie it to an outcome
otherwise:** `agents/review-lens-intent.md` carries `model: haiku`. That line is the
user's own edit, made from the panel before this change opened, carried onto the
branch by `2767491` rather than discarded. It is a setting, not a behaviour change
this work chose — but it arrived inside a rename commit and the commit message was
the only place that said so, which is not a contract.

**Out:**

- **Renaming the skills.** They did not collide. `install --global` reports exactly
  one conflict on this machine and it is the agent.
- **Renaming `spec-writer` or `work-reviewer`.** Checked, both free, and neither is a
  lens. Namespacing every payload agent on principle is a change nobody asked for.
- **Any behavioural change to a lens** — what it looks for, what it returns, the order
  they run in, the tool grants. This is five filenames.

## Constraints

- `scripts/check-payload` requires each agent's `name:` to equal its filename. Five
  files, five frontmatter lines, and the gate is what proves they agree.
- A dangling agent name fails as a lens that silently never ran, so the reference
  hunt is the work, not the rename.

## Prior decisions

- **All five, not just the one that broke.** Only `review-reliability` collided today.
  But five of the reporting machine's own agents already sit under a bare `review-`
  prefix, so `review-security`, `review-design` and `review-tests` were free by luck
  rather than by design — and `review-intent` is the same family launched by the same
  skill. A family named one way in four files and another way in the fifth is an
  inconsistency somebody straightens later with less context than we have now.

- **`review-lens` is the prefix because it is the one name that provably did not
  collide.** It was the single agent before the split, and nothing in the reporting
  machine's `~/.claude/agents` begins with it.

- **The collision was introduced by the split**, in this session, and would not exist
  otherwise. Recording that rather than presenting the rename as a tidy-up: the
  previous design had one agent with an unusual name and no way to hit this.

## Task breakdown

5. Rename the five lens agents to `review-lens-*`, update the launch table and the
   `Governs:` line, and confirm no reference survives.

## Verification criteria

- frontmatter parses and every agent's `name:` matches its new filename
  Proof: scripts/check-payload
- the four lens skills are still named in the form the static check can see
  Proof: scripts/check-payload
- every `Proof:` citation in the amended specs resolves
  Proof: skills/record-work/spec-drift --anchors
- **no reference to an old lens agent name survives in the payload or the code.**
  Checked with `rg --hidden`, because `.agents` starts with a dot and the obvious
  spelling of this check reports a clean repo while the specs still name the deleted
  agent — which is exactly what happened the last time a lens was renamed:

  ```
  rg -n --hidden --glob '!.git/*' --glob '!.agents/**' \
    'agents/review-(security|design|reliability|tests|intent)\.md'
  ```

  **The first draft of this criterion claimed the whole repository and was false when
  written.** `.agents/specs/review-project/spec.md` still names the five old files in
  its `Governs:` line and in prose, because a capability spec describes what is on
  `main` until phase 8 folds this delta onto it. Claiming otherwise made a criterion
  that failed the moment a reviewer ran it — and the reviewer did.

  No `Proof:` — no test asserts this. `check-payload` verifies that each agent that
  *exists* is coherent; it cannot see a name that should have stopped existing. An
  observation, until a check earns the citation.

- **the capability spec's own references move when the delta lands**, not before —
  `Governs:` and the three prose mentions. Phase 8's consolidation is where that is
  checked, and `--anchors` is the check.

- behaviour is a prompt and is checked by running it. **Owed:** one real review
  launched with the renamed agents, five reports relayed. Until then the rename is
  believed rather than observed.
