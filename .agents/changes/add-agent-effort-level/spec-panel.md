# The selector's effort column — delta

Targets: panel

Governs: internal/ui/models.go

The model selector marks agents and applies one model to the marked set. It gains one
column and one key, and nothing else about the screen moves.

## Outcomes

- **Every row shows its effort** after its model, and shows the session word when it
  declares none — the same rendering the row already uses for an undeclared model, because
  it is the same state.
- **`e` opens the effort catalogue** over the marked rows, exactly as `m` opens the model
  catalogue: same cursor, same escape, same "nothing marked" notice, same one-act apply.
- **`m` and `enter` are unchanged.** The model is what the screen is named after and what
  `enter` has always meant; rebinding it to a chooser would be a silent change to a
  reflex.
- **Rows show their new effort straight after applying**, without reopening the screen —
  the promise already made for the model.
- **Applying the level every marked row already has says so and writes nothing.** The
  writer refuses to rewrite an identical value deliberately, and a screen that reported
  success for a no-op would be reporting the wrong thing.
- **A marked row whose model cannot run the level is refused, named, and nothing is
  applied to any row.** The all-or-nothing act is the panel's promise too, and a partial
  apply on a screen showing twelve rows is a state nobody can read back.

## Scope boundaries

**In:** the effort column, the `e` chooser, its apply path, and the refusal above.

**Out:**

- **Grouping by effort.** Rows group by model, cheapest first, and that stays the only
  grouping. Grouping by the pair turns four groups into twenty on a screen whose whole
  argument is that a reader can see the shape at a glance.
- **A third screen.** The chooser is a mode over the same rows, like the model chooser. A
  screen that navigates to another screen to change one line is a screen with a hallway in
  it.
- **A header naming each group's effort.** Every row carries its effort in the effort
  column; a header would print it twice. The same reasoning already rejected a model
  header.
- **Confirming an effort change.** `y/n` is for the destructive actions. Writing a line
  into a file the user can read is not one.
- **`ultracode` as a choice.** Not a level, and not something a file can declare.

## Constraints

- **`internal/ui` imports neither `internal/target` nor `internal/agentmodel`**, and this
  does not change it. The levels arrive as a slice of choices and the apply as a function
  value, exactly as `ModelChoice` and `ApplyModel` already do. The panel reports the
  choice; `cli` runs it.
- **Which levels a model supports is not the panel's knowledge.** It renders the refusal it
  is handed. A `haiku` check inside the UI would be the third copy of one rule.
- **The row is built against the frame, not summed and hoped.** Two value columns, the
  `shared` warning and the longest name the payload ships came to 64 columns against a
  58-column interior, and `pad` refuses to truncate — so the overflow did not clip, it
  tore six columns of border off every row on the screen. The name column is therefore
  computed from what the frame can spare and the name yields to it, **visibly, with an
  ellipsis**: a name that says it was cut is a smaller lie than a frame that came apart,
  and it is the one column whose content the two beside it cannot be inferred from.

  The value columns and the warning never yield. Dropping one to save a name would hide
  state rather than shorten it.

  The layout tests are the gate, and the one that existed was not enough: its fixture's
  longest name was 15 runes and its border filter skipped agent rows, so it passed
  through the tear. That is a finding about the gate, not only about the row, and the
  criteria below name the tests that close it.

## Prior decisions

- **`e`, not a rebinding of `enter`.** *(assumed under `/libretto-attacca` — nobody was
  asked)* `e` is unbound on this screen today, verified against every `case` in
  `internal/ui/models.go`, `panel.go` and `model.go`. **What changes if this is wrong:**
  one string in one switch.

- **A column, not a grouping.** *(assumed — same run)* Stated above with its reason. The
  panel spec already records that the user chose "rows move when the model changes"; moving
  them on two axes was never what was agreed, and twenty groups is a different screen.
  **What changes if this is wrong:** `sortRowsByModel` gains a secondary key and
  `groupRule` a second level. The column stays either way.

## Task breakdown

1. `AgentRow` carries `Effort`; the row renders it, session word included.
2. `EffortChoice` and `ApplyEffort` join `ModelChoice` and `ApplyModel` on the model's
   `WithAgents` seam.
3. `e` opens the chooser; escape, cursor and the nothing-marked notice mirror `m`.
4. The apply path refreshes the rows and reports the refusal it is handed.

## Verification criteria

- a row shows its declared effort, and the session word when it declares none
  Proof: internal/ui/models_test.go TestRowsShowTheirEffort
- `e` opens the effort catalogue, and escape returns to the rows without quitting
  Proof: internal/ui/models_test.go TestEOpensTheEffortCatalogueAndEscapeReturns
- **`m` and `enter` still open the model catalogue**
  Proof: internal/ui/models_test.go TestEnterStillOpensTheModelCatalogue
- pressing `e` with nothing marked says so and opens nothing
  Proof: internal/ui/models_test.go TestChoosingEffortWithNothingMarkedSaysSo
- a chosen level reaches every marked row and no unmarked one
  Proof: internal/ui/models_test.go TestChosenEffortReachesOnlyTheMarkedRows
- the rows show the new effort without reopening the screen
  Proof: internal/ui/models_test.go TestRowsShowTheNewEffortAfterApplying
- **a refused apply leaves every row as it was and shows the reason**
  Proof: internal/ui/models_test.go TestARefusedEffortApplyChangesNoRow
- applying the level every marked row already has says so and calls nothing
  Proof: internal/ui/models_test.go TestApplyingTheEffortTheyAlreadyHaveSaysNothingChanged
- the rows still group by model only
  Proof: internal/ui/models_test.go TestRowsStillGroupByModelAlone
- **no row ever outgrows the frame**, at any width, with the longest name the payload
  ships and the `shared` warning on it
  Proof: internal/ui/models_test.go TestTheSelectorRowNeverOutgrowsTheFrameAtAnyWidth
- a name the frame cannot afford is elided visibly, and the value columns survive
  Proof: internal/ui/models_test.go TestALongNameIsElidedRatherThanTearingTheFrame
- nothing is elided when there is room
  Proof: internal/ui/models_test.go TestNamesAreNotElidedWhenThereIsRoom
