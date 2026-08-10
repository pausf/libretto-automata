# add-per-agent-model — plan

**Goal:** choose which model each payload agent runs on — from `libretto models` and
from the panel — so the lenses that pattern-match over prose stop being billed at the
rate of the ones that reason.

**Architecture:** one new package, `internal/agentmodel`, that reads and writes the
`model:` key in `agents/*.md` and owns the catalogue of legal values. The CLI and the
panel are two doors onto it; `internal/ui` stays off the filesystem and reaches it
through callbacks, as it already does for `Refresh` and `Runner`. The payload side is
independent: `review-lens` splits into four agent files so design and tests can be
cheap without security following them down.

Specs: `spec.md` (agent-models) · `spec-cli.md` (cli) · `spec-panel.md` (panel) ·
`spec-review-project.md` (review-project), all in this folder. One writer: the
orchestrator marks these boxes; sub-agents report.

## Global constraints

- six gates before every commit: `gofmt -l .` silent, `go vet ./...`,
  `go test ./... -count=1`, `scripts/check-payload`, `spec-drift --self-test`,
  `spec-drift --anchors`
- **Go 1.26.5, standard library only.** No YAML dependency — frontmatter here is a
  fenced block of `key: value` lines and must stay that
- `CLAUDE_HOME` points at a temp dir in anything that touches a target; never a real
  `~/.claude`
- the writer changes one line and leaves every other byte alone
- no `𝄞`, no `♩♪♫♬` anywhere a terminal renders
- **the repo's own `agents/*.md` are the files under edit.** Tests write to copies in
  `t.TempDir()`, never to the payload itself
- never pipe a gate into `head` when the exit code matters
- no push from any step; the release tag is phase 8's call, not a task here

---

### Task 1: read the declared model

Depends on: nothing. **Start here.**
From spec: `spec.md` breakdown 1. Closes: the three reading criteria.

**Files:** create `internal/agentmodel/frontmatter.go`,
`internal/agentmodel/frontmatter_test.go`

- the frontmatter is the block between a `---` on line 1 and the next `---`. A file
  that does not open with `---` has none — that is the same rule
  `scripts/check-payload:58` already applies, and it must not diverge from it
- `ReadModel(path) (string, error)` returns the declared value, or the empty string
  meaning *default* when the key is absent. **Absent is not an error**
- a `model:` line after the closing `---` is body text and is not the model
- [ ] **1.1** `TestReadModelReturnsTheDeclaredModel`,
  `TestReadModelReportsDefaultWhenTheKeyIsAbsent`, `TestReadModelIgnoresTheBody` —
  written and failing first
- [ ] **1.2** implement until they pass
- [ ] **1.3** six gates, exit codes read from files
- [ ] **1.4** committed

### Task 2: write and remove the key

Depends on: Task 1 (shares the frontmatter boundary logic).
From spec: `spec.md` breakdown 2. Closes: the five writing criteria.

**Files:** modify `internal/agentmodel/frontmatter.go`,
`internal/agentmodel/frontmatter_test.go`

- `SetModel(path, model string) error`. The empty string removes the key
- **byte-for-byte** everywhere but the one line. The test compares the whole file
  against the original with that line substituted — not a spot check on the key
- a new key goes inside the frontmatter block; position within it is not contracted,
  so pick one and let the test pin it
- a file with no frontmatter is refused and left **unchanged on disk**. The test
  asserts both halves; an error return with a mangled file is the failure worth
  catching
- setting the model an agent already has writes nothing — same mtime, same bytes
- [ ] **2.1** `TestSetModelInsertsWithoutDisturbingTheFile`,
  `TestSetModelReplacesInPlace`, `TestSetModelDefaultRemovesTheKey`,
  `TestSetModelRefusesAFileWithoutFrontmatter`, `TestSetModelIsIdempotent` — failing
  first
- [ ] **2.2** implement until they pass
- [ ] **2.3** six gates
- [ ] **2.4** committed

### Task 3: the catalogue

Depends on: nothing. **Independent — may run alongside Task 1.**
From spec: `spec.md` breakdown 3. Closes: the two catalogue criteria.

**Files:** create `internal/agentmodel/catalogue.go`,
`internal/agentmodel/catalogue_test.go`

- exactly four entries: `opus`, `sonnet`, `haiku`, and default (the empty string)
- each carries the label the CLI and panel print, including which plan it needs.
  **The plan text is a claim about somebody's subscription** — it says what the plan
  tier is called, never what this user has, because the binary cannot know
- `Valid(model string) bool`; everything not in the catalogue is refused
- the catalogue is the single list. The CLI and the panel both read it; neither
  keeps its own copy of the model names
- [ ] **3.1** `TestCatalogueListsTheSubscriptionModels`, `TestUnknownModelIsRefused` —
  failing first
- [ ] **3.2** implement until they pass
- [ ] **3.3** six gates
- [ ] **3.4** committed

### Task 4: apply one model to a set

Depends on: Tasks 2 and 3.
From spec: `spec.md` breakdown 4. Closes: the two applying criteria.

**Files:** create `internal/agentmodel/apply.go`, `internal/agentmodel/apply_test.go`

- `Agents(repoRoot) ([]Agent, error)` — every `agents/*.md`, name and current model,
  sorted, so the CLI and the panel list the same things in the same order
- `Apply(repoRoot string, names []string, model string) error`
- **validate the whole set before writing any of it.** Model in the catalogue, every
  name a real agent, every file readable with frontmatter. Then write
- the test for that is the important one: a set whose last member is unwritable
  leaves the *first* member unchanged. A partial write is the failure mode this
  design exists to prevent, so it gets the test, not a comment
- [ ] **4.1** `TestApplyReachesEveryAgentInTheSet`,
  `TestApplyWritesNothingWhenAnyAgentIsUnwritable` — failing first
- [ ] **4.2** implement until they pass
- [ ] **4.3** six gates
- [ ] **4.4** committed

### Task 5: `libretto models`

Depends on: Task 4.
From spec: `spec-cli.md` breakdown 5. Closes: all ten cli criteria.

**Files:** create `cmd/libretto/models.go`, `cmd/libretto/models_test.go`; modify
`cmd/libretto/main.go` (the `switch` at :137, and `usage()` at :1001)

- `models` lists every agent and its model, changes nothing, exits zero
- `models set <model> <agent>…` and `models set <model> --all`
- **`set` with neither agents nor `--all` is an error.** Not "all of them"
- unknown model or unknown agent → exit non-zero, nothing written
- `--global` / `--project` are already parsed before dispatch (`main.go:71-76`); this
  consumes the result. They pick **which target's installed agents are listed**
- **on write, one line saying the effect is not scoped.** Both targets symlink to the
  same repo file, so the change is shared. Say it once, at the moment it matters
- plain text, no escape codes, like every other non-panel command
- [ ] **5.1** the ten tests named in `spec-cli.md`, failing first
- [ ] **5.2** implement until they pass
- [ ] **5.3** six gates — `TestRunDispatch` covers the new case
- [ ] **5.4** committed

### Task 6: the selector screen

Depends on: nothing in Go — `internal/ui` never touches the filesystem, so this is
built against callbacks and fake rows. **Independent of Tasks 1–5.**
From spec: `spec-panel.md` breakdown 6. Closes: eight of the twelve panel criteria.

**Files:** create `internal/ui/models.go`, `internal/ui/models_test.go`; modify
`internal/ui/model.go` (screen state), `internal/ui/panel.go` (rendering)

- one row per agent: name, current model, mark. `space` marks, `a` marks all and `a`
  again clears all, `m` opens the catalogue, `esc` returns to the menu
- **nothing marked → nothing happens, and the panel says so.** No falling back to the
  cursor row
- the mark is `[x]` / `[ ]` **and** the theme's emphasis. Both signals, never one —
  the same rule the destination strip follows
- rows show their new model straight after applying
- `Update` stays free of I/O; the two callbacks come in the way `WithRefresh` and
  `WithRunner` already do
- the `pending` / `pendingScope` confirmation machinery is **not** reused; this is not
  a destructive action
- [ ] **6.1** the eight selector tests named in `spec-panel.md`, failing first
- [ ] **6.2** implement until they pass
- [ ] **6.3** six gates
- [ ] **6.4** committed

### Task 7: wire the panel to the package

Depends on: Tasks 4, 5 and 6.
From spec: `spec-panel.md` breakdown 7. Closes: the menu-row tally criteria.

**Files:** modify `cmd/libretto/main.go` (menu construction and callbacks),
`internal/ui/models_test.go`

- the menu gains `models`, and its row **reports rather than describes itself** —
  `2 on haiku · 3 on session`, cheapest first, undeclared counted as `session`
- the tally refreshes after applying, through the existing `Refresh`
- the callbacks are supplied here, in `cmd`, so `internal/ui` gains no dependency on
  `internal/agentmodel` and none on `internal/target`
- [ ] **7.1** `TestMenuRowReportsTheModelTally`,
  `TestMenuRowTallyRefreshesAfterApplying` — failing first
- [ ] **7.2** implement until they pass
- [ ] **7.3** six gates, plus `make preview` looked at — the panel is the one thing a
  test cannot fully judge
- [ ] **7.4** committed

### Task 8: split `review-lens` into four

Depends on: nothing. Payload only, no Go. **Independent — may start immediately.**
From spec: `spec-review-project.md` breakdown 8. Closes: its check-payload and
`--anchors` criteria.

**Files:** create `agents/review-security.md`, `agents/review-design.md`,
`agents/review-reliability.md`, `agents/review-tests.md`; delete
`agents/review-lens.md`; modify `skills/review-project/SKILL.md` (the launch section
at :255-285), `docs/FLOW.md`, and every remaining reference

- each new agent: `name:` matching its filename, the same `Read, Grep, Glob, Skill`
  grant `review-lens` declared, and the body naming its one skill
- `review-intent` keeps `Bash`; it is untouched apart from being allowed a `model:`
- the four `Skill(skill="…")` names stay written out in the form the static reference
  check can see. That requirement moves house, it does not disappear
- **hunt every reference to `review-lens`.** `rg -n 'review-lens'` over the whole repo
  must come back empty before this task closes; a dangling agent name fails as a lens
  that silently never ran
- [ ] **8.1** write the four agents, delete the old one
- [ ] **8.2** update the launch site, `docs/FLOW.md`, and the remaining references
- [ ] **8.3** `rg -n 'review-lens' .` returns nothing
- [ ] **8.4** six gates
- [ ] **8.5** committed

---

## What can start now

**Tasks 1, 3, 6 and 8** — none of them waits on anything.

Task 8 is the only one that touches no Go, so it is also the only one that can run
truly beside the others without two agents landing in the same package.

Then: 2 (after 1) → 4 (after 2, 3) → 5 (after 4) → 7 (after 4, 5, 6).

## Owed before this is done, not by a test

- **one real review run with the four split agents**, five reports relayed, and the
  token cost measured against the 307k baseline with design and tests on a cheaper
  model. Until that run the saving is a prediction, which is what the whole change is
  for. Same standing debt `review-project` already carries for its tool grants
- `make preview` looked at by a human at more than one width
