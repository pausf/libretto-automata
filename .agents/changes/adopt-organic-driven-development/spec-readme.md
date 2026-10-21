Targets: readme

# readme — the front door describes one stop before code

The README describes the flow in four places: the prose and the flow diagram in
`## What you get`, the numbered walk in `## Your first run`, the phase table below it, and
the `/libretto-attacca` paragraph. Today all four describe the shape this change retires:
a spec stop, a plan, tasks cut by a fresh agent, a second stop, then the push. This delta
brings them to the new shape: **one change document, one stop before code, a fresh reviewer
that reads the document before that stop, then the push question**. The delivery diagram,
the section order, the commands, and every other section stay as they are.

## Outcomes

1. **The flow diagram has two decisions, not three.** It still sits inside `## What you get`
   as the second of two Mermaid fences. Phases 1 to 5 write **one change document** (contract,
   approach, checklist). A **fresh reviewer** reads it. Then come **one stop** on the whole
   document, phases 6, 7 and 8, and **the push** as the last decision. No node draws a separate
   plan stop, a task cut, or a fresh agent that cuts tasks. Every label is ASCII, as this
   capability's *Constraints* already require. The current label `5→6` does not meet that
   rule.

2. **The `What you get` prose matches the diagram.** Before any code you get one document
   (what "done" means, what is out, how each promise is proven, how it gets built, and the
   checklist). A reviewer that saw none of the session reads it, you agree it once, and the
   last question is whether to push. The fresh reviewer that checks the *result* (phase 7)
   stays. The two readers are described as two, not merged into one.

3. **The first-run walk describes one stop before code.** Step 2 and step 4 of today's walk
   ("stops at the spec", "stops once the tasks are cut") collapse into one stop on the change
   document, which carries the pre-code reviewer's findings, including what the document
   fails to answer for each box. The closing sentence counts the places the flow waits as
   two: the stop and the push. The walk still names `/libretto-flow`, and the push stays
   "asked, never assumed".

4. **The phase table and the attacca paragraph stop describing retired pieces.** The `5→6`
   row no longer says the tasks are "cut by a fresh agent". The `/libretto-attacca` paragraph
   no longer says "same plan and tasks" as separate documents. The skill names in the table
   follow the payload. `write-spec`, `write-plan` and `write-tasks` remain skills per the
   `payload` delta, so they stay. A `review-spec` row is added only if the payload makes it
   the stop's owner, which `spec-payload.md` does.

5. **No acronym, no credit.** `README.md` never contains the acronym quoted in the
   2026-10-21 naming answer, as a whole word, and never contains the method's spelled-out
   name. The credit to gentle-ai lives only in `docs/FLOW.md`. This capability has never
   credited influences. *Standing on other people's work* lists vendored skills, and
   gentle-ai is not vendored, so the README does not gain a credit line.

6. **The front door does not grow.** `TestReadmeStaysShort`'s ceiling of 380 stays where it
   is, and the rewrite is net non-positive in lines. Read shows `README.md` at 380 lines
   today, so the ceiling has no slack. This delta does not raise it.

## Scope boundaries

**In:** the four places in `README.md` named above, the matching assertions in
`cmd/libretto/readme_test.go`, and the delta onto `.agents/specs/readme/spec.md` (Outcomes 3
and 7 and their criteria).

**Out:**

- **`docs/FLOW.md`, including the gentle-ai credit line.** That belongs to the `payload`
  delta, unit 7. This delta only *depends* on one anchor in it (see *Constraints*).
- **Every skill, agent and command file.** Those belong to `payload`.
- **Any Go behaviour outside `readme_test.go`.** That belongs to `cli`.
- **The delivery diagram, the section order, the install section, the commands table, the
  states, the environment table, `Learn more`, and the footer.** None of these describes the
  flow's stops.
- **Wording.** As the capability already says, no criterion pins a sentence. The guards
  below assert shape (decision-node count, ASCII labels) and two vocabulary terms the brief
  fixes (`change document`, `push`). They never assert a sentence.
- **Renumbering the phases.** The README keeps "an eight-phase … flow", because the `payload`
  delta keeps eight.

## Constraints

- **The proof reads files as text**, with no Markdown or Mermaid parser. A decision node is
  counted by its Mermaid rhombus syntax, `<id>{"`, inside the flow fence.
- **Every substring assertion runs over `flat()`**, per this capability's prior decision. Two
  guards here have already shipped unable to fire.
- **Diagram labels are ASCII; styling stays Mermaid defaults.** These are unchanged. This
  delta turns the label rule into a check, because the rule was in prose and `5→6` broke it
  unnoticed.
- **Line ceiling 380, unchanged** (Outcome 6).
- **`TestMovedReasoningLandedInDocs` depends on `docs/FLOW.md`, which the `payload` delta
  rewrites.** The anchor `more expensive than this one` (subject *"the spec stop is the cheap
  place to disagree"*) sits in FLOW.md's paragraph on the spec stop (FLOW.md:416-420). That
  paragraph is exactly what payload unit 7 rewrites for one stop. If the rewrite keeps the
  phrase, the subject label in the test is renamed to the stop and nothing else moves. If it
  drops the phrase, the anchor is chosen again under this capability's two anchor rules:
  counted against the *base* of both destination documents, and chosen from the reasoning,
  never from the fact. It must never be deleted to get green. Either way, the README unit
  lands after FLOW.md has its final wording.
- **No new dependency.** Use the standard library, as the rest of `readme_test.go` does.

## Prior decisions

### Answered by the user, 2026-10-21 (verbatim from the change's decision log)

- Q: for substantial work, what replaces the proposal → brief → spec → plan → tasks chain? → A: Un doc de cambio + specs vivos (Recommended) — one document per change (intent, EARS criteria, approach, checklist); capability specs with Governs:/Proof: and the landing stay
- Q: with one change document, where does the flow stop for the go-ahead before the push? → A: Un stop: el doc entero (Recommended) — contract, approach and checklist reviewed together, once, before the first line of code; then only the push question
- Q: who writes the change document, and where does the fresh eye stay? → A: Sesión escribe, revisor fresco lee (Recommended) — the session writes the document inline; one fresh reviewer (review-spec) reads it before the stop; spec-writer, plan-writer and task-cutter launches go from three to one review
- Q: what is this called in Libretto? → A: Sin sigla, vocabulario propio (Recommended) — no "ODD" in the payload; the flow stays "the Libretto flow", lanes named in its own vocabulary; FLOW.md credits gentle-ai as inspiration in one line
- Q: without task-cutter, who checks that the document is enough to build every box? → A: El revisor fresco lo hereda (Recommended) — the one fresh reviewer reads contract and checklist and reports, per box, what the document fails to answer

### Derived (not user answers)

- **The acronym rule reaches the README**, by the orchestrator's instruction for this
  subtask. The brief's correction scopes the *payload* rule to `skills/`, `agents/` and
  `commands/`, and the README is none of those. The README is still the public description
  of the flow, so it follows the same vocabulary.

### Reversed by this delta

- **Outcome 7, "the flow diagram, the eight phases in order with the three stops (spec, plan,
  push) marked as decisions"** is reversed. The flow diagram now has two decisions: the stop
  on the change document, and the push.
- **Outcome 7's criterion, "one of them names the three stops — spec, plan and push"**, is
  reversed. It is replaced by the decision-node count below.
- **Outcome 3's criterion, "names `/libretto-flow`, the spec stop, the plan stop and the
  push"**, is reversed. It now names `/libretto-flow`, the change document and the push.

### Unchanged

Mermaid, not committed images. Two diagrams, not one. Both stay inside `What you get`. The
diagrams are judged by looking as well as by the guard. Reasoning lives in `docs/`, and the
wording is not pinned.

## Task breakdown

1. **Guards first, watched red.** Amend `TestWhatYouGetCarriesTheDiagrams` and
   `TestReadmeWalksAFirstRun`, and add `TestReadmeNamesNoMethodAcronym`, per *Verification
   criteria*. Run them against today's README. The decision count (3) and the ASCII label
   (`5→6`) must fail, and so must `change document` in the walk. The acronym guard has
   nothing to fail on today, so it is watched red on a hand-broken copy and then reverted.
   Waits on nothing.
2. **Rewrite the four places in `README.md`** (Outcomes 1-4) so the guards go green and the
   line count does not rise. Render the flow diagram with mermaid-cli and look at it before
   the review seam, per the capability's *judged by looking* decision. Waits on 1, and on the
   `payload` delta settling the skill and agent names the phase table cites.
3. **Re-check the moved-reasoning anchor against the final `docs/FLOW.md`** (see
   *Constraints*). Rename the subject, or choose the anchor again under the two anchor rules.
   Waits on payload unit 7.

## Verification criteria

- **When** `README.md`'s `## What you get` section is read, it **shall** contain exactly two
  ```` ```mermaid ```` fences, and the fence that contains decision nodes **shall** contain
  exactly two of them, written `<id>{"`. The last of those two **shall** name `push`.
  Watched red on today's README, which has three.
  Proof: cmd/libretto/readme_test.go TestWhatYouGetCarriesTheDiagrams

- **If** any character outside ASCII appears inside either Mermaid fence in `## What you
  get`, **then** the guard **shall** fail, naming the character. Watched red on today's
  `5→6` label.
  Proof: cmd/libretto/readme_test.go TestWhatYouGetCarriesTheDiagrams

- **When** `## Your first run` is read through `flat()`, it **shall** contain
  `/libretto-flow`, `change document` and `push`. Watched red on today's README, which has no
  `change document`.
  Proof: cmd/libretto/readme_test.go TestReadmeWalksAFirstRun

  **Ceiling named:** this proves the vocabulary arrived. It does not prove the walk describes
  *one* stop. A walk that names the change document and still draws two stops stays green.
  The count of stops in prose is judged by reading, and the decision-node count above holds
  the diagram's count.

- **If** `README.md`, read through `flat()`, contains the acronym quoted in the 2026-10-21
  naming answer as a whole word (case-sensitive), or the method's spelled-out name
  (case-insensitive), **then** the guard **shall** fail, naming which one it found. Watched
  red on a hand-broken copy, because today's README carries neither.
  Proof: cmd/libretto/readme_test.go TestReadmeNamesNoMethodAcronym

- **While** this delta is applied, `README.md` **shall** stay at or under 380 lines, with
  the ceiling constant unedited.
  Proof: cmd/libretto/readme_test.go TestReadmeStaysShort

- **When** `docs/FLOW.md` has its final wording from the `payload` delta, every subject in
  `TestMovedReasoningLandedInDocs` **shall** still be absent from `README.md` and present in
  `docs/DESIGN.md` or `docs/FLOW.md`, with no subject removed from the list.
  Proof: cmd/libretto/readme_test.go TestMovedReasoningLandedInDocs
