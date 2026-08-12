# Plan — add-agent-effort-level

One writer: the orchestrator. A box is marked the moment its criterion is observed, and
it ships in the same commit as the task that closed it.

Three deltas, and the dependency is a straight line out of `internal/agentmodel`: nothing
in the CLI or the panel can be written before the package knows what an effort level is,
and both of them can be written at once after that.

## 1 · The package — the foundation everything waits on

- [x] **1.1 Parametrise the frontmatter reader and writer by key.**
      `const key = "model:"` becomes an argument; `ReadModel`/`SetModel` keep their names
      and their behaviour byte for byte.
      Spec: agent-models · task 1 · Closes: TestSetModelInsertsWithoutDisturbingTheFile
      and every existing frontmatter criterion, unchanged and still green.
      Waits on: nothing.

- [x] **1.2 The effort catalogue.** Five levels weakest first, their labels, `ValidEffort`,
      and which catalogue model supports effort at all — `haiku` none, the rest all five.
      Spec: agent-models · task 3 · Closes: TestEffortCatalogueListsTheFiveLevels,
      TestUnknownEffortIsRefused, TestWhichModelsSupportEffort.
      Waits on: nothing. Independent of 1.1 — different file, no shared symbol.

- [x] **1.3 `ReadEffort` and `SetEffort`** over 1.1's shared implementation.
      Spec: agent-models · task 2 · Closes: TestReadEffortReturnsTheDeclaredEffort,
      TestReadEffortReportsDefaultWhenTheKeyIsAbsent, TestReadEffortIgnoresTheBody,
      TestSetEffortInsertsWithoutDisturbingTheFile, TestSetEffortReplacesInPlace,
      TestSetEffortDefaultRemovesTheKey, TestSetEffortIsIdempotent,
      TestTheTwoKeysDoNotDisturbEachOther.
      Waits on: 1.1.

- [x] **1.4 `Agent` carries `Effort`.** One field, read in `Agents()` beside the model.
      Spec: agent-models · task 5 · Closes: TestAgentsReportsEachCurrentEffort.
      Waits on: 1.3.

- [x] **1.5 `ApplyEffort`** — one level onto a set as one act, the level and every target
      agent's model validated before the first file is opened.
      Spec: agent-models · task 4 · Closes: TestApplyEffortReachesEveryAgentInTheSet,
      TestApplyEffortWritesNothingWhenAnyAgentCannotRunIt,
      TestApplyEffortAllowsAnAgentOnTheSessionModel.
      Waits on: 1.2, 1.4.

- [x] **1.6 `Apply` clears a stale effort** when the model it writes supports none.
      Spec: agent-models · task 5 · Closes:
      TestApplyModelClearsEffortWhenTheModelSupportsNone.
      Waits on: 1.2, 1.3.

## 2 · The CLI — after the package, independent of the panel

- [x] **2.1 The effort column and the levels trailer in `models`.**
      Spec: cli · task 1 · Closes: TestModelsListsEffortBesideTheModel,
      TestModelsListsTheEffortCatalogue.
      Waits on: 1.4, 1.2.

- [x] **2.2 `models effort <level> <agents…|--all>`**, reusing `set`'s argument handling
      rather than reimplementing the `--all` refusal.
      Spec: cli · task 2 · Closes: TestModelsEffortWritesOnlyTheNamedAgents,
      TestModelsEffortAllReachesEveryAgent, TestModelsEffortRefusesWithNothingNamed,
      TestModelsEffortRefusesAnUnknownLevel,
      TestModelsEffortRefusesAnAgentThatCannotRunTheLevel,
      TestModelsEffortDefaultRemovesTheKey.
      Waits on: 1.5.

- [x] **2.3 `models set` reports a cleared effort** on the rows it cleared.
      Spec: cli · task 3 · Closes: TestModelsSetReportsAClearedEffort.
      Waits on: 1.6, 2.1.

## 3 · The panel — after the package, independent of the CLI

- [x] **3.1 `AgentRow` carries `Effort`, and the row renders it**, session word included,
      inside the existing width budget.
      Spec: panel · task 1 · Closes: TestRowsShowTheirEffort,
      TestRowsStillGroupByModelAlone.
      Waits on: nothing in the package — `internal/ui` imports no `agentmodel`. Rows
      arrive as data. **Can start immediately, in parallel with 1.x.**

- [x] **3.2 `EffortChoice` and `ApplyEffort` on the `WithAgents` seam.**
      Spec: panel · task 2 · Closes: covered by 3.3's and 3.4's criteria; this task has no
      criterion of its own and exists because the seam is the thing that must not grow a
      second shape.
      Waits on: 3.1.

- [x] **3.3 `e` opens the chooser**; escape, cursor and the nothing-marked notice mirror
      `m`, and `m`/`enter` keep meaning the model.
      Spec: panel · task 3 · Closes: TestEOpensTheEffortCatalogueAndEscapeReturns,
      TestEnterStillOpensTheModelCatalogue, TestChoosingEffortWithNothingMarkedSaysSo.
      Waits on: 3.2.

- [x] **3.4 The apply path** refreshes the rows and reports the refusal it is handed.
      Spec: panel · task 4 · Closes: TestChosenEffortReachesOnlyTheMarkedRows,
      TestRowsShowTheNewEffortAfterApplying, TestARefusedEffortApplyChangesNoRow.
      Waits on: 3.3.

- [x] **3.5 Wire the panel's seam to the real package** in `cmd/libretto`, so the `e` key
      reaches `ApplyEffort` rather than a test double.
      Spec: panel · constraints — the panel reports the choice, `cli` runs it.
      Closes: no unit criterion; proven by the panel run test that already covers the
      model path (`cmd/libretto/panelrun_test.go`).
      Waits on: 2.2, 3.4.

## 4 · The front door

- [x] **4.1 `README.md` and `AGENTS.md`** gain the verb in the tables that already list
      `models set`, and the README's sample output gains the effort column.
      Spec: cli · task 4 · Closes: nothing enforces this — the readme gate checks payload
      slash commands, not subcommands. It is here because that is the only thing that
      will make it happen.
      Waits on: 2.2.

## 5 · Landing

- [x] **5.1 All six gates.** `gofmt -l .`, `go vet ./...`, `go test ./... -count=1`,
      `scripts/check-payload`, `spec-drift --self-test`, `spec-drift --anchors`.
      Waits on: everything above.

- [ ] **5.2 Fold the three deltas onto `agent-models`, `cli` and `panel`**, delete the
      change folder, in the same commit as the last code.
      Waits on: 5.1.

## What can start now

**1.1, 1.2 and 3.1** — three independent starting points. 1.1 and 1.2 touch different
files in the same package; 3.1 touches a package that imports neither.

Everything else waits on one of them.
