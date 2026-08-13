# ship-frugal-model-defaults — showing it on the screen

Targets: panel

The request names the surface: *"podriamos añadir en libretto recomendaciones por agente
de cual seria recomendable poner aqui"* — `aqui` being the `libretto models` selector.

**And the frame cannot afford a third column.** That is measured, not feared: at the
narrowest interior the panel supports, 58 columns, the row already spends 2 on the indent,
6 on the cursor and box, 12 on the model, 10 on the effort and 9 on the `shared` warning,
leaving **19 for the name against a floor of 12**. A `recommended` column wide enough to
hold `(session)` is nine of those, which puts the name under its floor — and `pad` does
not truncate, so the row does not clip, it tears the border off every row on the screen.
The panel's own criterion already settles what happens next: *the value columns and the
warning never yield.*

So the recommendation appears **where the choice is made rather than at rest**: inside the
model and effort catalogues, at the moment the user opens them.

## Outcomes

**1 · The catalogue marks what is recommended for the agents that are marked.**

Pressing `m` opens the model catalogue for the marked set; the recommended entry carries a
visible mark. `e` does the same for effort, inside the narrowing that screen already
applies.

**A mark and nothing else. The reason is not rendered here**, and that is forced rather
than chosen: the catalogue line is `cursor + 5 + name + label` inside an interior that
narrows to 58 columns, and the reasons run to seventy runes. A reason on this screen tears
the frame — the same arithmetic that kept a third column off the row, arriving one screen
over. The reasons live in `libretto models`, which has no width to defend.

**2 · What "disagree" means is decided per catalogue, and an agent with no opinion does
not vote.**

Three rules, because the word covers three different sets and each was read both ways:

| The marked set | The model catalogue | The effort catalogue |
|---|---|---|
| all recommend the same model, differing efforts | marks that model | marks nothing |
| recommend different models | marks nothing, and says the marked agents differ | marks nothing |
| some recommended, some not known | the known ones decide; unknown agents do not vote | same |
| exactly one agent marked | marks its recommendation | marks its recommendation |

**Each catalogue compares only its own field.** A set agreeing on `sonnet` and disagreeing
on effort is not "a set that disagrees" when the question on screen is which model.

**An agent with no recommendation abstains rather than blocking.** Marking every agent on a
machine that has the payload plus twenty-two of the user's own would otherwise mark nothing
for ever, which is the common gesture producing the useless answer.

**3 · The effort half is judged against the model the agent declares, not the one it is
recommended.**

The effort catalogue already narrows to what the marked rows can run *now*, and that is the
set the user can actually pick from. So a recommended level outside that offer is simply not
marked.

**Named consequence:** `review-lens-design` declares `sonnet` and is recommended `haiku`,
which has no effort levels at all — so its effort catalogue marks nothing, and that is
indistinguishable on screen from "no recommendation exists". The listing is where those two
are told apart, which is the second thing pushing the reasons there.

**4 · An agent with no recommendation changes nothing on the screen.**

No mark, no reason, no empty row where one would be. A user's own agent is not something
this repository has an opinion about, and a blank marker column would read as a verdict of
"none recommended" rather than "not known".

**5 · The screen still never applies anything.**

The mark is a mark. The user moves the cursor and presses enter exactly as before, onto
the recommended entry or past it. **Nothing preselects the recommendation** — putting the
cursor on it would be the tool typing the answer with extra steps.

## Scope boundaries

**In:**

- `internal/ui/models.go` — the mark inside both catalogues, and the differ notice
- `internal/ui/models_test.go`
- `.agents/specs/panel/spec.md` — the criteria

**Out, and named so it cannot be quietly added:**

- **No third value column, at any width.** Measured above. *Brings it back:* nothing that
  is a change to this screen — it would be a change to what the frame is.
- **No reason text on this screen**, in the catalogue or anywhere else. The frame cannot
  hold seventy runes at 58 columns, and eliding a reason to fit produces a sentence that
  stops mid-argument. *Brings it back:* a wider `MinContentWidth`, which is a change to the
  panel rather than to this feature.
- **No extra line under the catalogue.** The open-catalogue row reservation is part of the
  window arithmetic that keeps the screen inside the terminal height, and a line not in
  that count is a scroll bug on a short terminal. The differ notice reuses the footer the
  screen already has.
- **No divergence glyph on the resting row.** A single character would fit the budget, and
  it was rejected: a glyph needs a legend, the legend row is already full at
  `space mark · a all · m model · e effort · esc back`, and the panel spec forbids the
  colour-only affordance that would let it fit. *Brings it back:* a measured complaint
  that the recommendation is invisible until three keystrokes in.
- **No preselection, no reordering, no filtering by recommendation.** The catalogue's
  render order is contracted cheapest-first and this does not touch it; a second ordering
  is already on the panel's own out-list.
- **No new key binding.** `m` and `e` already open the two catalogues, and `r` staying
  free is worth more than saving a keystroke.
- **The panel still learns nothing about agent files.** The recommendation arrives on
  `AgentRow` through the same adapter every other field does, so `internal/ui` never
  imports `internal/agentmodel`.

## Constraints

- **`TestTheSelectorRowNeverOutgrowsTheFrameAtAnyWidth` must still pass unchanged.** It is
  a boundary this change respects, **not a criterion**: it renders resting rows with no
  catalogue open, so no implementation of this change makes it red. Citing it as proof
  would be a criterion green before the work starts, which is the defect the target spec
  already caught on itself once.
- **`TestRowsShowTheirEffort` asserts `(session)` appears exactly once** on a row that
  declares no effort, and the CLI's sibling asserts exactly twice on its line. Both are
  counting assertions that a recommendation rendered in the wrong place would break. They
  are the tripwire, and they stay.
- **Legible without colour.** The mark is a character, not a hue — the same rule the row
  mark and the `shared` warning already keep.
- All six gates pass before any commit.

## Prior decisions

**Settled by measurement, this change:**

- **The row budget has no room.** 58 − 2 − 6 − 12 − 10 − 9 = 19 against a floor of 12.
- **`r` is free** — the selector's switch consumes `esc q tab s up k down j space a m
  enter e` and nothing else. It stays free.

**Assumed under `/libretto-attacca`, because nobody was there to answer.** Each names what
changes if it is wrong:

- **B1 · The recommendation is shown at the moment of choosing rather than at rest.** The
  request said `aqui` — this screen — and the frame's arithmetic decided which part of it.
  *If wrong:* the alternative is the resting-row divergence glyph named in the boundaries,
  which needs a legend slot this screen does not have spare, or a wider `MinContentWidth`,
  which is a change to the panel rather than to this feature.
- **B2 · A mixed marked set marks nothing rather than marking each agent's own.** The
  catalogue is one list of models, not one per agent, so per-agent marks would need a
  column inside the catalogue for the same reason the row cannot have one. *If wrong:* the
  catalogue grows a per-agent line and the disagreement case disappears.

- **B3 · An unknown agent abstains rather than blocking**, and the effort half is judged
  against the declared model rather than the recommended one. Both were genuine forks that
  review found unpinned. *If wrong:* each is one condition, and both fail towards showing
  no mark rather than a wrong one.

## Task breakdown

1. `AgentRow` carries the recommendation. **The adapter that fills it is `cli`'s** — see
   the sibling delta; this one only reads the field.
2. The model catalogue marks the recommended entry for a set that agrees, unknown agents
   abstaining.
3. A set disagreeing on the model marks nothing and says so in the footer.
4. An unrecommended agent renders the catalogue exactly as today.
5. The effort catalogue, same rules, judged against the declared model inside its existing
   narrowing.
6. The catalogue holds at 58 columns with the longest label the table ships.
7. Six gates, then apply this delta onto `.agents/specs/panel/spec.md`.

## Verification criteria

- **the open catalogue never outgrows the frame**, at the narrowest interior, with the
  longest entry the recommendation makes possible. **This is the criterion the row test
  cannot be** — the existing width test renders resting rows with no catalogue open, so
  the one line this change can lengthen is measured by nothing until here.
  Proof: internal/ui/models_test.go TestTheOpenCatalogueHoldsTheFrameAtItsNarrowest

- **the model catalogue marks the recommended entry.** One marked agent with a
  recommendation; the mark is on the right row.
  Proof: internal/ui/models_test.go TestTheModelCatalogueMarksTheRecommendation

- **a set disagreeing on the model is marked nowhere, and the screen says so.** Two agents
  recommended onto different models; no entry carries the mark.
  Proof: internal/ui/models_test.go TestAMixedMarkedSetIsNotGivenOneRecommendation

- **an agent with no recommendation abstains instead of blocking.** A known agent marked
  together with an unknown one still marks the known one's recommendation — otherwise
  marking everything on a real machine answers nothing, for ever.
  Proof: internal/ui/models_test.go TestAnUnknownAgentDoesNotBlockTheOthersRecommendation

- **a set agreeing on the model but differing on effort still marks the model.** Each
  catalogue compares its own field; a disagreement about depth is not a disagreement about
  tier.
  Proof: internal/ui/models_test.go TestDisagreementIsJudgedPerCatalogue

- **an agent with no recommendation leaves the catalogue exactly as it was.** No mark, no
  reason, no blank row standing in for one.
  Proof: internal/ui/models_test.go TestAnUnrecommendedAgentAddsNothingToTheCatalogue

- **the effort catalogue marks against the declared model, inside its existing
  narrowing**, and marks nothing when the recommended level is not in the offer.
  Proof: internal/ui/models_test.go TestTheEffortCatalogueMarksTheRecommendation

- **the mark is legible without colour.** A character, like the row mark and the `shared`
  warning, because colour alone is an affordance this panel has already refused once.
  Proof: internal/ui/models_test.go TestTheRecommendationMarkIsLegibleWithoutColour

- **nothing is preselected.** Opening the catalogue leaves the cursor where it has always
  opened; the recommendation is marked, never chosen.
  Proof: internal/ui/models_test.go TestTheRecommendationIsNeverPreselected
