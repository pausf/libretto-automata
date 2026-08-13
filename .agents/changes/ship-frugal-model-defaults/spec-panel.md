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

Pressing `m` opens the model catalogue for the marked set. The entry the repository
recommends carries a visible mark and its per-agent reason. `e` does the same for effort,
inside the narrowing that screen already applies.

**2 · Marked agents that disagree are said to disagree, never averaged.**

Marking `review-lens-security` and `review-lens-design` together is marking two agents with
different recommendations. The catalogue says so and marks nothing, because a single mark
would be a recommendation for a set that has none. **This is the same rule as
`Nothing marked means nothing happens`** — an affordance that is sometimes right teaches
the user to ignore it.

**3 · An agent with no recommendation changes nothing on the screen.**

No mark, no reason, no empty row where one would be. A user's own agent is not something
this repository has an opinion about, and a blank marker column would read as a verdict of
"none recommended" rather than "not known".

**4 · The screen still never applies anything.**

The mark is a mark. The user moves the cursor and presses enter exactly as before, onto
the recommended entry or past it. **Nothing preselects the recommendation** — putting the
cursor on it would be the tool typing the answer with extra steps.

## Scope boundaries

**In:**

- `internal/ui/models.go` — the mark and the reason inside both catalogues
- `internal/ui/models_test.go`
- `.agents/specs/panel/spec.md` — the criteria

**Out, and named so it cannot be quietly added:**

- **No third value column, at any width.** Measured above. *Brings it back:* nothing that
  is a change to this screen — it would be a change to what the frame is.
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

- **`TestTheSelectorRowNeverOutgrowsTheFrameAtAnyWidth` must still pass unchanged**, with
  the longest name the payload ships and the `shared` warning on it. This change adds
  nothing to a row, so it should — and if it does not, the change is wrong rather than the
  test.
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

## Task breakdown

1. `AgentRow` carries the recommendation, filled by the adapter — no new import in
   `internal/ui`.
2. The model catalogue marks the recommended entry for a homogeneous marked set, with the
   reason.
3. A mixed set marks nothing and says why.
4. An unrecommended agent renders the catalogue exactly as today.
5. The effort catalogue, same three rules, inside its existing narrowing.
6. Six gates, then apply this delta onto `.agents/specs/panel/spec.md`.

## Verification criteria

- **the row is untouched, and the frame still holds at every width.** This change adds
  nothing to a resting row; the existing width test passes with no edit, which is the
  evidence that the third column was avoided rather than squeezed in.
  Proof: internal/ui/models_test.go TestTheSelectorRowNeverOutgrowsTheFrameAtAnyWidth

- **the model catalogue marks the recommended entry, and says why.** One marked agent with
  a recommendation; the mark is on the right row and the reason is on screen.
  Proof: internal/ui/models_test.go TestTheModelCatalogueMarksTheRecommendation

- **a marked set that disagrees is marked nowhere.** Two agents recommended onto different
  models; no entry carries the mark, and the screen says the marked agents differ.
  Proof: internal/ui/models_test.go TestAMixedMarkedSetIsNotGivenOneRecommendation

- **an agent with no recommendation leaves the catalogue exactly as it was.** No mark, no
  reason, no blank row standing in for one.
  Proof: internal/ui/models_test.go TestAnUnrecommendedAgentAddsNothingToTheCatalogue

- **the effort catalogue keeps the same three rules**, inside the narrowing it already
  applies to the marked set's model.
  Proof: internal/ui/models_test.go TestTheEffortCatalogueMarksTheRecommendation

- **the mark is legible without colour.** A character, like the row mark and the `shared`
  warning, because colour alone is an affordance this panel has already refused once.
  Proof: internal/ui/models_test.go TestTheRecommendationMarkIsLegibleWithoutColour

- **nothing is preselected.** Opening the catalogue leaves the cursor where it has always
  opened; the recommendation is marked, never chosen.
  Proof: internal/ui/models_test.go TestTheRecommendationIsNeverPreselected
