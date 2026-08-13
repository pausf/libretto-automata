# cut-payload-context-cost

Targets: payload

Three rules into the payload, all of them about what a phase *loads* rather than what it
does. The flow's behaviour does not change: the same phases run, the same artifacts come
out, and the same questions get asked. What changes is how much context each phase pays
to produce them.

**Read this delta as a cost change, never as a capability change.** A rule here that
alters what the flow decides has escaped its scope, and the scope boundaries below name
that explicitly.

## Outcomes

**1 · A phase opens the spec that governs its work, never the corpus.**

`write-spec` step 2 states how to find the ground truth in the specification, and the
instruction is selection: identify the capability the change targets — from the change's
own `Targets:`, from a `Governs:` line, or from the project's index if it has one — and
open that spec. Reading every file under the spec root is named as the thing not to do.

The rule is stated **without naming this repository's layout**. `docs/SPEC.md` is a local
convention and skills install into projects that have none; a skill that names it works
only here, which is the failure `AGENTS.md` forbids in as many words.

**2 · A spec-writer reads the brief sections that touch it, plus the vocabulary, and
returns deltas and objections.**

`brief.md` carries stable section headings. Each spec-writer's prompt names the sections
its subtask touches; **the vocabulary section is named for every writer, always**, because
that section is the entire mechanism stopping two writers from giving one concept two
names.

The return contract tightens to deltas and objections — what the brief did not know, what
it got wrong, what only the user can settle — and never a restatement of what the brief
already said. **Returning nothing found stays a legal, required answer**: "never
restatements" must not read as licence for an empty return, because silence-is-not-success
is load-bearing in two files.

**3 · The flow states that model and effort do not move mid-phase.**

`docs/FLOW.md` carries it under *Delegation*, where per-agent model choice already lives,
and `skills/review-project/SKILL.md` carries it where it is dearest — five lenses, five
contexts, one switch rebilling all of them. The statement names the mechanism: a switch
invalidates the cached prefix, so the whole context is rebilled at full input price
instead of the cached fraction.

## Scope boundaries

**In:**

- `skills/write-spec/SKILL.md` — the selection rule in step 2; the brief's section
  structure and the sub-agent prompt contract in step 2b; the return contract
- `agents/spec-writer.md` — reading named sections, and the tightened return
- `docs/FLOW.md` — the cache-stability rule under *Delegation*
- `skills/review-project/SKILL.md` — the same rule where the lens fan-out pays it
- `scripts/check-payload` — one `check_wiring` row per rule
- `.agents/specs/payload/spec.md` — three criteria, applied at phase 8

**Out, and named so it cannot be quietly added:**

- **`skills/find-work/SKILL.md` gets no selection rule.** Phase 1 reads
  `.agents/changes/`, never `.agents/specs/` — verified, not assumed. A rule forbidding
  a read that does not happen is prose that can only ever be green, and the payload spec
  already carries the incident where exactly that kind of row proved nothing.
- **`skills/review-work/SKILL.md` gets no selection rule either.** Step 1 already assembles
  "the capability spec(s) its `Targets:` names" — selection is what it already does. The
  rule generalises review-work's precedent to write-spec rather than restating it back at
  the file it came from.
- **No new skill, no new agent, no new command.** Three prose rules do not need a home of
  their own.
- **No broadcast of rule 3 to every skill.** Two homes, chosen because the reasoning lives
  in one and the cost is concrete in the other. A standing rule copied into eight files is
  eight copies to keep in sync, and the one nobody edits is the one that reads as
  authoritative.
- **No measurement, no instrumentation, no token accounting.** That is
  `add-token-cost-to-metrics`, queued separately.
- **The `write-spec` fan-out still hands over a path, never excerpted text.** Minimisation
  is at read time — which sections the writer opens — not at write time. Excerpting the
  brief into N prompts is N copies in N contexts, which reverses the decision the brief
  exists to implement.

## Constraints

- **Markdown only. No Go changes**, except the one shell script that checks the markdown.
- **A payload rule about prompt behaviour cannot be machine-checked.** `check-payload`
  proves wiring — that the decisive words are still in the file that owns them — and each
  criterion says so in its own body, matching the `**Proved as wiring only**` precedent
  already in the payload spec.
- **Every `check_wiring` row is written and observed FAILING before the prose it
  describes exists.** A row added to match prose already in the file has never proved it
  could fail. The payload spec records this happening twice; the rule is not new here, only
  applied.
- **The pattern is the skill's own wording**, not this spec's phrasing. A row matching a
  phrase the file never contained is green on its first run and proves nothing.
- **`check-payload` fails on a `docs/…` path referenced from a skill unless the path
  exists in this repo** or sits under the `docs/(specs|superpowers)` exclusion. Rule 1
  must therefore state the index generically and add no new `docs/` reference.
- All six gates pass before any commit.

## Prior decisions

**Settled by reading the code, this change:**

- **`review-work` selects by `Targets:` and `write-spec` will select the same way.** Two
  routes to "which capability" would be two answers, and the index is the fallback for a
  change that has no delta yet — not the primary.
- **The brief stays one re-runnable file with one author.** `write-spec` already turns on
  being able to correct the brief and regenerate affected specs; section-scoped *reading*
  keeps that, section-scoped *writing* would destroy it.

**Assumed under `/libretto-attacca`, because nobody was there to answer.** Each names what
changes if it is wrong:

- **A1 · Rule 1 is stated convention-agnostically and lands in `write-spec` only.** The
  proposal named three skills; two of them already satisfy it or never read specs at all,
  which the recon established by quotation. *If wrong:* the trimmed two get the same
  sentence and two more `check_wiring` rows — additive, no rework of what lands here.
- **A2 · The brief gains stable section headings, and the prompt names sections rather
  than carrying excerpts.** The alternative — excerpting into each prompt — measurably
  reverses `write-spec`'s own "one path is one line". *If wrong:* the prompt contract
  changes shape and `agents/spec-writer.md` changes with it; the brief's headings stay
  useful either way.
- **A3 · Rule 3 lands in `docs/FLOW.md` and `skills/review-project/SKILL.md`, and nowhere
  else.** The mid-phase-mutable dial is the *session's* model and effort, not agent
  frontmatter — frontmatter is static and written by `libretto models`. A rule the payload
  cannot enforce is stated where a reader deciding a fan-out will meet it. *If wrong:* it
  moves to `skills/evidence/`, the only skill every phase invokes — rejected here because
  `evidence`'s one rule is about observation and a cost rule does not follow from it.
- **A4 · `write-spec`'s "This repository is one." is corrected in the same change.** It
  claims this repo is the single-file-spec case; it has thirteen capability directories
  under `.agents/specs/`, which `spec-drift` resolves first. It is a false sentence in the
  exact paragraph rule 1 rewrites, and leaving a known-false line in a file being edited
  is how it survives another four tags. *If wrong:* revert one sentence.

## Task breakdown

1. Add three `check_wiring` rows to `scripts/check-payload`, run it, **observe all three
   red**, and record the output.
2. Rule 1 — `skills/write-spec/SKILL.md` step 2: the selection instruction; and the
   convention-agnostic index sentence, correcting A4's false line.
3. Rule 2 — `skills/write-spec/SKILL.md` step 2b: brief section headings, the prompt
   naming sections plus vocabulary, the tightened return; and `agents/spec-writer.md` to
   match.
4. Rule 3 — `docs/FLOW.md` *Delegation*, and `skills/review-project/SKILL.md`.
5. Run `scripts/check-payload`, observe the three rows green.
6. Apply the three criteria onto `.agents/specs/payload/spec.md` (phase 8, with the code).

## Verification criteria

- **a phase opens the spec that governs its work, never the corpus.** `write-spec` step 2
  instructs selection by `Targets:`, `Governs:` or the project's index, and names reading
  the whole spec root as the thing not to do. Stated without naming any repository's
  layout, so it survives installation into a project with a different convention.
  **Proved as wiring only** — that the instruction is still in the file that owns it, not
  that a session obeyed it.
  Proof: scripts/check-payload

- **a spec-writer is given sections and returns deltas.** The brief carries stable
  headings; the prompt names the sections the subtask touches and always names the
  vocabulary; the return is what the brief did not know, never a restatement of what it
  did. The empty return stays illegal — silence is not success, unchanged.
  **Proved as wiring only.**
  Proof: scripts/check-payload

- **model and effort do not move mid-phase.** The flow states that a switch invalidates
  the cached prefix and rebills the whole context at full input price, and states it where
  a fan-out is being decided rather than in every skill.
  **Proved as wiring only** — and this one cannot be more than wiring even in principle,
  because the dial it constrains is the session's and the payload cannot read it.
  Proof: scripts/check-payload

Each row above is added to `scripts/check-payload` **before** the prose it matches, and
observed failing. A row that was green on its first run is removed and rewritten, not
kept.
