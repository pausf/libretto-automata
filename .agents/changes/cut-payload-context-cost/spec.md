# cut-payload-context-cost

Targets: payload

Three rules into the payload, all of them about what a phase *loads* rather than what it
does. The flow's behaviour does not change: the same phases run, the same artifacts come
out, and the same questions get asked. What changes is how much context each phase pays
to produce them.

**Read this delta as a cost change, never as a capability change.** A rule here that
alters what the flow decides has escaped its scope, and the scope boundaries below name
that explicitly.

**All three land in one file — `skills/write-spec/SKILL.md` — plus one uncontracted
paragraph in `docs/FLOW.md`.** That is narrower than the proposal drew it, and every
narrowing is recorded under prior decisions with what would widen it again.

## Outcomes

**1 · A phase opens the spec that governs its work, never the corpus.**

`write-spec` step 2 states how to find the ground truth in the specification, and the
instruction is selection: identify the capability the change targets — from the change's
own `Targets:`, from a `Governs:` line, or from the project's index if it has one — and
open that spec. Reading every file under the spec root is named as the thing not to do.

The rule is stated **without naming this repository's layout**, because skills install
into projects that have none. That property is a review-time reading and **not
machine-checked**: no `check_wiring` row can prove the absence of an unbounded set, and
the one concrete token available — `docs/SPEC.md` — legitimately appears in the same
paragraph as part of `spec-drift`'s documented search order. The criterion claims only
what a row can hold up.

**2 · A spec-writer reads the brief sections that touch it, plus the vocabulary.**

`brief.md` carries a **fixed, enumerated** set of section headings — the five the skill
already lists as the brief's contents: conventions, constraints, settled decisions,
vocabulary, and the six-pillar structure. Fixed rather than per-change, so a prompt naming
sections names them from a set both sides know.

Each spec-writer's prompt names the sections its subtask touches. **The vocabulary section
is named for every writer, always**, because that section is the entire mechanism stopping
two writers from giving one concept two names.

**The return contract does not change, and gets no criterion.** `agents/spec-writer.md`
already promises "what you found that the brief did not know … a question only the user
can settle" and already forbids the empty return. Adding "never restatements" would be
prose restating a promise that exists, and its row would be green the moment it was
written — which is the failure the payload spec records twice.

**3 · The flow states that model and effort do not move mid-phase.**

`skills/write-spec/SKILL.md` carries it in step 2b, beside the fan-out that pays for it:
N spec-writers is N contexts, and a switch part-way through rebills every one of them.
The statement names the mechanism — a switch invalidates the cached prefix, so the whole
context is rebilled at full input price instead of the cached fraction.

`docs/FLOW.md` carries the same reasoning under *Delegation*, **as documentation and not
as contract**: no capability governs `docs/`, so no criterion here cites it. That is
stated rather than left to be noticed.

## Scope boundaries

**In:**

- `skills/write-spec/SKILL.md` — the selection rule in step 2; the brief's fixed heading
  set and the sub-agent prompt contract in step 2b; the cache-stability rule in step 2b
- `docs/FLOW.md` — the cache-stability reasoning under *Delegation*, uncontracted
- `scripts/check-payload` — three `check_wiring` rows, all three over
  `skills/write-spec/SKILL.md`
- `.agents/specs/payload/spec.md` — three criteria, applied at phase 8

**Out, and named so it cannot be quietly added:**

- **`skills/find-work/SKILL.md` gets no selection rule.** Phase 1 reads
  `.agents/changes/`, never `.agents/specs/` — verified by quotation, not assumed. A rule
  forbidding a read that does not happen is prose that can only ever be green.
- **`skills/review-work/SKILL.md` gets no selection rule either.** Step 1 already assembles
  "the capability spec(s) its `Targets:` names" — selection is what it already does.
- **`agents/spec-writer.md` is not edited.** Rule 2's return half was dropped for
  duplicating what that file already promises; its reading half is carried by the prompt
  the orchestrator writes, which is `write-spec`'s text.
- **`skills/review-project/SKILL.md` is not edited.** It was rule 3's second home until
  review found that `review-project` owns `skills/review-*/**` more specifically than
  `payload` owns `skills/**` — so a payload criterion over that file would be a second
  spec claiming a path, which is the condition where a change gets recorded in neither.
- **No new skill, no new agent, no new command.**
- **No criterion cites `docs/FLOW.md`.** Nothing governs `docs/`; a `Proof:` over an
  unowned file is a citation no drift check can anchor.
- **No measurement, no instrumentation, no token accounting.** That is
  `add-token-cost-to-metrics`, queued separately.
- **The fan-out still hands over a path, never excerpted text.** Minimisation is at read
  time — which sections the writer opens — not at write time. Excerpting the brief into N
  prompts is N copies in N contexts, which reverses the decision the brief implements.

## Constraints

- **Markdown only. No Go changes**, except the one shell script that checks the markdown.
- **A payload rule about prompt behaviour cannot be machine-checked.** `check-payload`
  proves wiring — that the decisive words are still in the file that owns them — and each
  criterion says so in its own body, matching the `**Proved as wiring only**` precedent
  already in the payload spec.
- **Every `check_wiring` row is written and observed FAILING before the prose it
  describes exists.** A row added to match prose already in the file has never proved it
  could fail. The payload spec records this happening twice.
- **The pattern is the skill's own wording**, not this spec's phrasing. A row matching a
  phrase the file never contained is green on its first run and proves nothing.
- **One row names one file** — `check_wiring` is `rg` against a single path. A rule with
  two homes needs two rows or it is half-proved, which is why rule 3 has one contracted
  home and says so.
- All six gates pass before any commit.

## Prior decisions

**Settled by reading the code, this change:**

- **`review-work` selects by `Targets:` and `write-spec` will select the same way.** Two
  routes to "which capability" would be two answers; the index is the fallback for a
  change that has no delta yet, not the primary.
- **The brief stays one re-runnable file with one author.** Section-scoped *reading* keeps
  the regenerate-from-a-corrected-brief property; section-scoped *writing* would destroy
  it.
- **The brief's heading set is fixed and enumerated, not per-change.** The five bullets
  the skill already lists become the five headings. Per-change headings would leave
  "names the sections its subtask touches" unverifiable against anything.

**Assumed under `/libretto-attacca`, because nobody was there to answer.** Each names what
changes if it is wrong:

- **A1 · Rule 1 is stated convention-agnostically and lands in `write-spec` only.** The
  proposal named three skills; two of them already satisfy it or never read specs at all.
  *If wrong:* the trimmed two get the same sentence and two more rows — additive, no
  rework.
- **A2 · The brief gains a fixed heading set, and the prompt names sections rather than
  carrying excerpts.** Excerpting reverses `write-spec`'s own "one path is one line".
  *If wrong:* the prompt contract changes shape and `agents/spec-writer.md` changes with
  it; the headings stay useful either way.
- **A3 · Rule 3's contracted home is `skills/write-spec/SKILL.md`, with `docs/FLOW.md`
  carrying the reasoning uncontracted.** The mid-phase-mutable dial is the *session's*
  model and effort — agent frontmatter is static and written by `libretto models` — so the
  payload cannot enforce this, only state it where a reader deciding a fan-out will meet
  it. *If wrong:* it moves to `skills/evidence/`, the only skill every phase invokes —
  rejected here because `evidence`'s one rule is about observation and a cost rule does
  not follow from it.
- **A4 · `write-spec`'s "This repository is one." is corrected in the same change.** It
  claims this repo is the single-file-spec case; it has thirteen capability directories
  under `.agents/specs/`, which `spec-drift` resolves first. It is a false sentence in the
  exact paragraph rule 1 rewrites. *If wrong:* revert one sentence.

**Amended after `review-spec`, before phase 5 read this.** Five findings, all acted on,
all of them narrowing:

- rule 3's second home removed — `review-project` owns that file more specifically
- rule 3's `docs/FLOW.md` half declared uncontracted — nothing governs `docs/`
- rule 2's return-contract half dropped entirely — it duplicated `agents/spec-writer.md`
  and its row would have been green on first run
- rule 1's "names no repository layout" clause moved out of the criterion — an absence
  claim over an unbounded set cannot go red
- the brief's headings fixed and enumerated — the ambiguity forked the prompt contract

## Task breakdown

1. Add three `check_wiring` rows to `scripts/check-payload`, all over
   `skills/write-spec/SKILL.md`, run it, **observe all three red**, record the output.
2. Rule 1 — `write-spec` step 2: the selection instruction; and the convention-agnostic
   index paragraph, correcting A4's false line.
3. Rule 2 — `write-spec` step 2b: the fixed heading set, and the prompt naming sections
   plus vocabulary.
4. Rule 3 — `write-spec` step 2b: the cache-stability statement; and `docs/FLOW.md`
   *Delegation*.
5. Run `scripts/check-payload`, observe the three rows green.
6. Apply the three criteria onto `.agents/specs/payload/spec.md` (phase 8, with the code).

## Verification criteria

- **a phase opens the spec that governs its work, never the corpus.** `write-spec` step 2
  instructs selection by `Targets:`, `Governs:` or the project's index, and names reading
  the whole spec root as the thing not to do.
  **Proved as wiring only** — that the instruction is still in the file that owns it, not
  that a session obeyed it. That it names no repository's layout is a review-time reading
  and is deliberately not cited here: no row can prove an absence over an unbounded set.
  Proof: scripts/check-payload

- **a spec-writer is given named brief sections, and always the vocabulary.** The brief's
  headings are a fixed enumerated set; the prompt names the ones the subtask touches and
  names the vocabulary regardless. The return contract is unchanged and carries no
  criterion, because `agents/spec-writer.md` already promises it.
  **Proved as wiring only.**
  Proof: scripts/check-payload

- **model and effort do not move mid-phase.** `write-spec` step 2b states that a switch
  invalidates the cached prefix and rebills the whole context at full input price, beside
  the fan-out where N contexts pay it.
  **Proved as wiring only** — and this one cannot be more than wiring even in principle,
  because the dial it constrains is the session's and the payload cannot read it.
  `docs/FLOW.md` carries the same reasoning and is **not** cited: no capability governs
  `docs/`.
  Proof: scripts/check-payload

Each row above is added to `scripts/check-payload` **before** the prose it matches, and
observed failing. A row that was green on its first run is removed and rewritten, not
kept.
