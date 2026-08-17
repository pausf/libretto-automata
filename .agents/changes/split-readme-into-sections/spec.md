# split-readme-into-sections — delta

Targets: readme

The ask: the README is long and sometimes prescriptive; shorter sections — quickstart, usage,
contrib — would speed adoption.

**The sections it names already exist.** *What you get*, *Install*, *Your first run*, *Commands*,
then reference detail, then *Learn more* — that is quickstart, usage and the rest, and `readme`
outcome 1 already holds them in that order with a test. So this is **not a restructure.** It is
length and register, which is what the ask's own second half said.

## What reading it found, and neither was in the ask

1. **`/libretto-attacca` is explained twice inside *Your first run*** — once as *"answers all
   three in advance"* immediately after step 7, then again two paragraphs later as *"answers
   those stops in advance"* with a fuller version. A reader who noticed would wonder which one
   was current.

2. **76 lines of model-selection detail carry arguments**, which `readme` outcome 4 already says
   belong in `docs/`: why per-agent models exist at all, and the fifteen-line explanation of how
   an alias resolves per provider with its `os.Getenv` justification. Those are arguments about
   *why it is that way*, and the capability's own prior decision draws exactly this line — what
   a command **does** stays, why it is that way goes.

**And the capability named this change's mechanism before it was needed.** `readme` outcome 6's
ceiling says: *the replacement, the day that matters, is a criterion about what a row must
contain.* Its `contributing` sibling says the same about file length. A README with no length
criterion grows back, and every argument moved out is one paragraph of room for the next one.

## Outcomes

1. **`README.md` is at most 380 lines**, from 389 — a ratchet against growth rather than a
   dramatic cut. The reduction comes from one duplicate removed and **five** arguments
   relocated, each of which lands in `docs/` and is checked to have landed. No fact is deleted
   to reach it.

   **This number was 340 in the first draft, and 340 was wrong.** It was chosen before anything
   was measured, on an estimate that the relocations came to ~60 lines. They came to 13: a
   relocation swaps prose for a pointer, so it saves far less than it reads like it should.
   Grinding to 340 would have meant cutting the commands table, the model and effort listings,
   the five states or the environment table — **reference a reader wants**, and cutting it would
   have produced a worse README that passes its own guard, which is the exact failure this
   capability was created to prevent.

   Corrected after measurement rather than met. The record of the wrong number is kept because a
   ceiling nobody can see the derivation of is a ceiling the next person raises.

2. **`/libretto-attacca` is described once.** The fuller of the two paragraphs survives; the
   earlier one-liner goes, because a step list is not where a mode gets explained.

3. **Five more arguments live in `docs/` rather than `README.md`**, joining the eight outcome 4
   already tracks:

   | Argument | Lands in |
   |---|---|
   | the spec stop is the cheap place to disagree, because a wrong sentence costs a line and the same mistake costs a day as code | `docs/FLOW.md` |
   | phase 1's source order is the point — starting new work while a change sits half-finished is how the half-finished thing gets abandoned | `docs/FLOW.md` |
   | `attacca` will not answer a gate, and never merges, tags or releases | `docs/FLOW.md` |
   | an alias resolves from the environment per provider, and that resolution is `os.Getenv` with no request and no credential | `docs/DESIGN.md` |
   | why other people's skills are vendored rather than depended on, and only the parts the flow calls by name | `docs/DESIGN.md` |

   Each is **absent from `README.md` and present in `docs/DESIGN.md` or `docs/FLOW.md`** — the
   same both-ends check outcome 4 already makes, because in a single-file diff *moved* and
   *deleted* are indistinguishable.

   **Present in *either*, not in the named one.** The table above says where each argument went;
   the check does not pin it there, because the existing eight anchors are checked against the two
   documents concatenated and splitting the behaviour would give nine anchors two rules. Both
   files own reasoning by design, so pinning an argument to one of them buys precision nobody
   needs and costs a false failure the day an argument legitimately moves between them.

4. **`CONTRIBUTING.md` is linked from *Learn more*.** This is the ask's *contrib* half, and it is
   this capability's job rather than `contributing`'s: that capability's scope boundary says a new
   front-door link is `readme`'s change, not its own.

## Scope boundaries

**In:** `README.md`, `docs/FLOW.md`, `docs/DESIGN.md`, and `cmd/libretto/readme_test.go`.

**Out**, named:

- **the section order does not change, and no section is added or removed.** `readme` outcome 1
  holds five load-bearing headings in order and its test would catch a reorder. The ask asked for
  sections that already exist
- **no facts are deleted to hit the line count.** Every command, flag, state and environment
  variable stays. A line count met by removing reference content is a worse README that passes
  its own guard, which is the failure the whole capability was created to prevent
- **the README is not split into multiple files.** The ask's wording allows that reading; it is
  refused, because a front door made of links to other front doors is the failure `readme`'s own
  prior decisions warn about — *how a README becomes a page that says nothing and links elsewhere*
- **the two Mermaid diagrams are untouched.** Outcome 7 holds them and they are inside *What you
  get*
- **the badge row is untouched.** PR 2 of this stack just made the tests badge honest
- **`AGENTS.md` is untouched.** Different reader, and the `contributing` capability owns that
  boundary now
- **no rewording for its own sake.** A prose polish pass is invisible in review and is how a
  documentation change becomes unreviewable

## Constraints

- **`readme`'s existing seven outcomes all keep passing.** Section order, install-steps-only,
  the first-run walk, the eight already-moved arguments, link resolution, every command named,
  and the two diagrams. This change is measured by the whole of `TestReadme*` staying green, not
  only by its own new case
- **whitespace is normalised before matching**, per `flat()`. That file has shipped two guards
  unable to fire because Markdown wrapped, and the third caught a phrase across a line break
- **an argument survives a paraphrase**, which `readme` already records as the honest ceiling of
  a substring check: the README kept an argument while swapping two words and the guard stayed
  green. So each anchor is chosen as the phrase a paraphrase would have to lose
- the line count is `wc -l` on the file, so a criterion about it is exact rather than a judgment
- all six gates green

## Prior decisions

- **380 lines, corrected from 340 after measuring** — *assumed, nobody was asked.*
  `/libretto-attacca` answered the contract stop. The first number came from an estimate and the
  estimate was out by a factor of four. **The honest choice at that point was between cutting
  reference content to hit a number and moving the number**, and the scope boundary already
  forbade the first. What changes if 380 is wrong: one line of test.

  **Ceiling, and it is real:** a line count says nothing about density. A README rewritten into
  380 very long lines passes. It is a ratchet against growth, not a measure of readability — and
  `readme` already records that readability is the one thing no test here holds.

- **The ask's premise was half right, and the half that was wrong is worth stating.** *Long and
  sometimes prescriptive* — the prescriptive part was real: one mode explained twice, and five
  arguments sitting in a front door that `docs/` already owned. The *long* part mostly is not
  fixable by moving prose, because what remains is the commands table, the model and effort
  listings, the five states and the environment table. That is a reference section, and a
  reference that stops listing things is not shorter, it is incomplete.

  This is the second ask in this batch whose premise did not survive reading the code, after the
  badge one. Both are recorded rather than quietly built around.

- **The arguments move rather than being deleted** — *assumed.* Deleting them would be faster and
  would lose the reasoning permanently; `readme`'s prior decisions record that this exact choice
  was put to the user once and answered *moved, not cut*. Following that answer rather than
  re-asking it.

- **`docs/FLOW.md` for the three flow arguments and `docs/DESIGN.md` for the alias one** —
  *assumed.* Those files already own those subjects, and `readme` records that a third file would
  split one subject three ways.

- **The fuller `attacca` paragraph survives, not the earlier one-liner** — *assumed.* The
  one-liner sits inside a numbered step list, where a mode does not get explained.

- **The line count is a criterion rather than a note** — *assumed.* Without it the file grows back
  and the next reader inherits the same complaint. This is the mechanism two of this repository's
  own capability ceilings already proposed for themselves.

## Task breakdown

1. **The new guards first, red**: the line count against 389, the five new anchors against a
   README that still carries them, and the `CONTRIBUTING.md` link against a *Learn more* that
   does not.
2. **Then the edits**: remove the duplicate, move the five arguments into `docs/`, add the link.
3. Land the delta onto `readme` and delete the change folder.

**Steps 1 and 2 are one box.** The guards alone leave the suite red; the edits alone move
arguments with nothing checking they arrived, which is the exact difference between *moved* and
*deleted* that outcome 4 exists for.

## Verification criteria

- **`README.md` is at most 380 lines.**
  Proof: cmd/libretto/readme_test.go TestReadmeStaysShort

- **each of the five relocated arguments is absent from `README.md` and present in
  `docs/DESIGN.md` or `docs/FLOW.md`**, matched over the flattened document, so a deletion cannot
  pass as a move.
  Proof: cmd/libretto/readme_test.go TestMovedReasoningLandedInDocs

- **`README.md` links to `CONTRIBUTING.md`, and that link resolves.** Both halves in the test that
  owns README links: a resolving-links test passes a README with no such link at all, so its
  presence is asserted there rather than inferred.
  Proof: cmd/libretto/readme_test.go TestReadmeLinksResolve

  **The presence assertion does not live in the length test.** It was written there first, and a
  future reader who deleted the link would have been told the README was too long — a failing run
  that exists but points at the wrong thing is a worse guard than a missing one, because it sends
  the fix somewhere else.

- **every existing `readme` criterion still passes** — section order, install steps, the first-run
  walk, the original eight moved arguments, link resolution, every command named, the two
  diagrams.
  Proof: cmd/libretto/readme_test.go TestReadmeSectionsAreInReadingOrder

**What nothing here tests:** whether the result reads better to a stranger. `readme` already
records that as the one thing no test holds, and that pretending otherwise produced the README
this capability replaced — every fact in that version was correct.
