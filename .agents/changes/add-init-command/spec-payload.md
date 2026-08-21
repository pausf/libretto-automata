Targets: payload

# Delta — the bootstrap non-goal gets an owner, and the refusal sites point at it

The payload capability forbids creating a `changes/` directory in a project that has
none (`.agents/specs/payload/spec.md:155`). That prohibition was written when nothing
sanctioned existed to do it; `libretto init` (the sibling delta, `Targets: cli`) now
does. This delta rewrites the non-goal so implicit creation by skills stays forbidden
while explicit, user-invoked `libretto init` becomes the named owner of bootstrap, and
gives the two refusal sites one guarded sentence each pointing at it.

## Outcomes

- The payload spec's Scope boundaries non-goal at line 155 no longer reads as a blanket
  ban. Replacement bullet, landing onto the capability spec:

  > - **implicit creation of a `changes/` directory in a project that has none.** A
  >   staging area nobody empties is a second source of truth. Bootstrap has exactly
  >   one owner — explicit, user-invoked `libretto init` — and a skill or command that
  >   finds the directory missing still refuses, pointing at the owner where `libretto`
  >   is on PATH and refusing identically where it is absent.

- `skills/write-spec/SKILL.md` (the "Never create a `changes/` directory" paragraph,
  line ~176) gains one sentence, guarded by the record-work grammar: where `libretto`
  is on PATH, the sanctioned bootstrap is `libretto init`; where it is absent, the
  refusal stands unchanged.
- `commands/libretto-queue.md` (the "Never create the directory" section, line ~70)
  gains the same guarded sentence, matching the same grammar.
- `skills/write-spec/SKILL.md` frontmatter `version:` bumps `"1.2"` → `"1.3"` — its
  instructions changed meaningfully. `commands/*.md` carry no `version:` frontmatter
  (verified: only `description:`), so `libretto-queue.md` bumps nothing.

## Scope boundaries

**In:** the one non-goal bullet in `.agents/specs/payload/spec.md`; one sentence in
each of the two refusal sites; the write-spec frontmatter version.

**Out — non-goals:**

- **anything about what `init` does.** The command's behaviour — what it creates, when
  it refuses, idempotence — is the sibling delta (`spec.md`, `Targets: cli`). This
  delta names the command in one clause and never its internals.
- **weakening either refusal.** Both sites refuse exactly as before, PATH or no PATH.
  The pointer only tells the user where the sanctioned bootstrap lives; it never runs
  it, suggests running it unprompted, or turns the refusal into a fallback.
- **a third refusal site.** No other payload document creates or declines to create
  `.agents/changes/`; none gains a pointer.
- **restating the rule.** The two sites keep pointing at each other for the reasoning
  ("the same rule `write-spec` holds") — the pointer is one sentence, not a copy of
  the non-goal.

## Constraints

- The guarded-binary-reference grammar is `skills/record-work/SKILL.md:104-108`:
  "Where `libretto` is on PATH … where the binary is absent, say … and continue" — an
  explicit absent path, never a bare "run libretto init". Both pointers match it.
- Skills never reference `scripts/` or `docs/`; the pointer names only the binary,
  which is legal under the existing record-work precedent (`libretto wiki`,
  `libretto land` clauses).
- `scripts/check-payload` must stay green: frontmatter parses, references resolve.
- The capability spec edit is a delta — it lands onto `.agents/specs/payload/spec.md`
  in the landing commit, same as any other delta. Pre-2026-08-17 neighbouring criteria
  stay untouched.

## Prior decisions

From `decisions.md`, Session 2026-08-21, pau.sanchez:

- Q: what happens to the payload non-goal ("no creating a changes/ directory in a
  project that has none")? → A: it is rewritten, not deleted: implicit creation by
  skills stays forbidden; `libretto init` becomes the named, explicit, user-invoked
  owner of bootstrap. The two refusal sites (write-spec, libretto-queue) gain a
  guarded one-line pointer ("where `libretto` is on PATH: `libretto init`"), absent
  path unchanged. If wrong: the pointers are one sentence each to remove. (assumed)

Standing, not relitigated here: skills must be self-sufficient once installed
(payload spec, Outcomes); the guarded-binary grammar with an explicit absent path
(record-work precedent, landed 2026-08); commands carry no version frontmatter.

## Task breakdown

1. Rewrite the non-goal bullet in `.agents/specs/payload/spec.md` Scope boundaries
   (line 155) to the replacement text above. Lands with this delta's landing commit.
2. Add the guarded pointer sentence to `skills/write-spec/SKILL.md`'s refusal
   paragraph and bump frontmatter `version:` to `"1.3"`.
3. Add the guarded pointer sentence to `commands/libretto-queue.md`'s refusal section.
   No version to bump — the file has no version frontmatter.

Each unit is one file; any order; all three land together with the sibling delta.

## Verification criteria

- **When** `write-spec` or `libretto-queue` finds no `.agents/changes/` in the
  project, it **shall** still refuse to create the directory; the payload spec
  **shall** name explicit, user-invoked `libretto init` as the sole owner of
  bootstrap and **shall** keep implicit creation by skills a non-goal.
  **Ceiling named:** the anchor keeps the sentences findable in their files, never
  proves a session obeyed prose — the same limit every skill criterion here lives with.
  Proof: .agents/specs/payload/spec.md
- **Where** `libretto` is on PATH, `skills/write-spec/SKILL.md`'s refusal paragraph
  **shall** point at `libretto init` as the sanctioned bootstrap; **where** the binary
  is absent, the paragraph **shall** leave the refusal unchanged — the pointer only
  says where bootstrap lives.
  Proof: skills/write-spec/SKILL.md
- **Where** `libretto` is on PATH, `commands/libretto-queue.md`'s refusal section
  **shall** point at `libretto init` with the same guarded grammar; **where** the
  binary is absent, it **shall** say so and stop, exactly as it does today.
  Proof: commands/libretto-queue.md
- **When** this delta lands, `skills/write-spec/SKILL.md` frontmatter **shall** carry
  `version: "1.3"`, and its frontmatter **shall** still parse with `name:` matching
  the directory.
  Proof: scripts/check-payload
