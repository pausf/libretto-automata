# ship-frugal-model-defaults — plan

Live state. The orchestrator is the only writer; a box is marked the moment its task is
verified, and it ships in the commit that closed it.

Execution is `build-and-check`, phase 6.

Two deltas, one per capability:
`.agents/changes/ship-frugal-model-defaults/spec-cli.md` ·
`.agents/changes/ship-frugal-model-defaults/spec-panel.md`

**Red-first per task.** Every criterion here cites a Go test that does not exist yet, and
each can be written and watched fail before its production code. **Three of them cannot go
red against a do-nothing stub** — the table-integrity walks are true of an empty table —
so those are forced red deliberately, by breaking the table on purpose, and that is
recorded rather than assumed.

## The table — nothing renders until this exists

- [x] **1 · `recommend` and its table** — seven entries, `false` for anything else.
      `cmd/libretto/recommend.go`: a `recommendation` carrying model, effort and reason;
      `recommend(name)` returning it with `false` for an unknown agent. Seven entries, the
      values and reasons from the `cli` delta's A1 table.
      **In `cmd/libretto`, never `internal/agentmodel`** — that package settled that it
      does not know this repository, and a map of payload agent names is the payload's
      agent list.
      *From:* cli outcomes 1 · *Closes:* criteria 2, 3.
      *Blocks:* everything. Nothing blocks it.
      Proof: cmd/libretto/recommend_test.go TestAnUnknownAgentGetsNoRecommendation
      Proof: cmd/libretto/recommend_test.go TestEveryRecommendationCarriesAReason

- [x] **2 · The three integrity walks** — all four guards forced red on purpose, observed, reverted:
      `sonnet-5` → *which the catalogue does not know*; `xhigh` on `haiku` → *which has no
      levels*; `spec-writer` onto `fable` → *nothing is steered onto the priciest tiers*;
      a reason emptied → *a verdict without one is an instruction*.
      Every entry passes `agentmodel.Valid`; every recommended effort passes
      `ValidEffort` **and** `SupportsEffort` for its recommended model; nothing names
      `opus` or `fable`.
      **These are true of an empty table**, so each is forced red by breaking the table on
      purpose — an entry naming `sonnet-5`, a level on `haiku`, an entry on `fable` —
      observed failing, and reverted. A guard that has never fired is a guard nobody has
      tested.
      *From:* cli constraints · A2 · *Closes:* criteria 1, 4. *Waits on:* 1.
      Proof: cmd/libretto/recommend_test.go TestEveryRecommendationIsTypeable
      Proof: cmd/libretto/recommend_test.go TestNothingIsRecommendedOntoThePriciestTiers

## The listing — where the reasons live

- [ ] **3 · The recommendation column and its reason**
      `libretto models` prints both. An agent with no recommendation prints nothing — no
      blank cell standing in for an opinion that does not exist.
      **`TestModelsListsEffortBesideTheModel` asserts `(session)` exactly twice on the
      inheriting line and must stay green unedited.** It is the tripwire for a column
      rendered into the wrong place.
      *From:* cli outcomes 2 · *Closes:* criterion 5. *Waits on:* 1.
      Proof: cmd/libretto/models_test.go TestModelsListsTheRecommendationAndItsReason

- [ ] **4 · Saying when an agent runs against its recommendation**
      A recommendation nobody can tell they are ignoring changes nothing.
      `review-lens-design` is that case on the day this ships — it declares `sonnet` and is
      recommended `haiku`, deliberately, per A3.
      *From:* cli outcomes 2 · A3 · *Closes:* criterion 6. *Waits on:* 3.
      Proof: cmd/libretto/models_test.go TestModelsMarksAgentsRunningAgainstTheRecommendation

- [ ] **5 · One field through the existing adapter**
      `agentRows` in `cmd/libretto/main.go` gains the recommendation exactly as it gained
      `Efforts`. **`internal/ui` imports no `agentmodel` after this, and that is the
      assertion** — the seam is the point, not the field.
      *From:* cli outcomes 3 · *Closes:* criterion 7. *Waits on:* 1.
      Proof: cmd/libretto/models_test.go TestAgentRowsCarryTheRecommendation

## The screen — a mark, and never a reason

- [ ] **6 · The catalogue holds the frame at its narrowest**
      58 columns, catalogue open, longest entry the table makes possible.
      **This is the criterion the existing row-width test cannot be**: that test renders
      resting rows with no catalogue open, so the one line this change can lengthen is
      measured by nothing until here. Written first for that reason.
      *From:* panel outcomes 1 · constraints · *Closes:* criterion 8. *Waits on:* 5.
      Proof: internal/ui/models_test.go TestTheOpenCatalogueHoldsTheFrameAtItsNarrowest

- [ ] **7 · The mark, for a set that agrees**
      One marked agent, then several agreeing. A character, not a hue.
      *From:* panel outcomes 1 · *Closes:* criteria 9, 13. *Waits on:* 6.
      Proof: internal/ui/models_test.go TestTheModelCatalogueMarksTheRecommendation
      Proof: internal/ui/models_test.go TestTheRecommendationMarkIsLegibleWithoutColour

- [ ] **8 · The three ways a set can fail to agree**
      Different models → nothing marked, footer says so. An unknown agent **abstains**
      rather than blocking — marking everything on a real machine must not answer nothing
      for ever. A set agreeing on model and differing on effort still marks the model,
      because each catalogue compares its own field.
      *From:* panel outcomes 2 · B3 · *Closes:* criteria 10, 11, 12. *Waits on:* 7.
      Proof: internal/ui/models_test.go TestAMixedMarkedSetIsNotGivenOneRecommendation
      Proof: internal/ui/models_test.go TestAnUnknownAgentDoesNotBlockTheOthersRecommendation
      Proof: internal/ui/models_test.go TestDisagreementIsJudgedPerCatalogue

- [ ] **9 · The effort catalogue, judged against the declared model**
      Inside the narrowing it already applies. A recommended level outside that offer marks
      nothing — including `review-lens-design`, whose `haiku` recommendation has no levels
      at all.
      *From:* panel outcomes 3 · *Closes:* criterion 14. *Waits on:* 8.
      Proof: internal/ui/models_test.go TestTheEffortCatalogueMarksTheRecommendation

- [ ] **10 · Nothing is preselected, and an unknown agent changes nothing**
      The cursor opens where it has always opened. An agent with no recommendation renders
      the catalogue exactly as today — no mark, no blank row standing in for one.
      *From:* panel outcomes 4, 5 · *Closes:* criteria 15, 16. *Waits on:* 7.
      Proof: internal/ui/models_test.go TestTheRecommendationIsNeverPreselected
      Proof: internal/ui/models_test.go TestAnUnrecommendedAgentAddsNothingToTheCatalogue

## Closing

- [ ] **11 · Six gates, then apply both deltas**
      `gofmt -l .`, `go vet ./...`, `go test ./... -count=1`, `scripts/check-payload`,
      `spec-drift --self-test`, `spec-drift --anchors`.
      Then fold the `cli` delta onto `.agents/specs/cli/spec.md` and the `panel` delta onto
      `.agents/specs/panel/spec.md`, **both in the same commit**, and delete the change
      folder with them. Half-consolidated is the one state with no honest description.
      *From:* both task breakdowns · *Waits on:* 10.

      `.agents/specs/cli/spec.md` is also touched by the open PR #47. Both are appends to
      different sections; if git disagrees, the merge is the merger's to resolve and
      neither side is rewritten here.

## What can start now

**Task 1, alone.** Nothing renders, marks or abstains until the table exists.
