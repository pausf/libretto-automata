# add-multi-tool-skill-targets — Implementation Plan

> **For agentic workers:** this plan is live state, read and marked by phase 6
> (`build-and-check`). One writer: the orchestrator marks boxes, in the same
> commit as the work that closed them. Steps use checkbox (`- [ ]`) syntax.

**Goal:** install the skills kind into Codex CLI (`~/.agents/skills`) and
OpenCode (`~/.config/opencode/skills`) as two new symlink targets, selectable
from the CLI flags and the panel strip, with the payload's prose reading
correctly in all three tools.

**Architecture:** two new implementations of the existing `target.Target`
interface, two new `Scope` values through the existing `Resolve` seam, the
panel's destination list (`scopeOrder`) grown from two to four. No new
mechanism anywhere — the interface was built for exactly this.

**Tech Stack:** Go 1.26.5 stdlib only. No new dependencies.

## Global Constraints

- No code outside `internal/target` derives a path from `~/.agents`,
  `~/.config/opencode` or `~/.claude` (spec-targets; doc strings excepted).
- No test writes to a real home directory: every fixture sets `CLAUDE_HOME`,
  `AGENTS_HOME` and `OPENCODE_HOME` to temp dirs (AGENTS.md · Never).
- A command acts on exactly one destination; nothing iterates targets to act.
- All six gates pass before every commit; conventional commits, no AI
  attribution.
- `ponytail:` comments in English.

Spec traceability: each task names its delta. Criteria quoted are the delta's
`Proof:` lines — a task closes when its named tests pass.

---

### Task 1: `target.Codex` — waits on nothing

**Delta:** spec-targets · **Proof:** `TestCodexRootResolution`,
`TestCodexAcceptsOnlySkills`, `TestCodexExists`

**Files:** create `internal/target/codex.go`, `internal/target/codex_test.go`

**Interfaces — produces:** `target.Codex` (implements `Target` + `Exists()`),
`target.NewCodex()`, `target.EnvAgentsHome = "AGENTS_HOME"`.

- [x] Write the failing tests, mirroring `target_test.go`'s Claude set:
  root resolution (`AGENTS_HOME` wins via `t.Setenv`, `~/.agents` fallback),
  `Kinds()` exactly `[Skills]`, `Dir(Skills) == root/skills`,
  `Accepts(Agents)==false`, `Accepts(Commands)==false`, `Exists()` false on a
  missing dir and true after `os.MkdirAll`.
- [x] Run `go test ./internal/target/ -run TestCodex -v` — expect compile
  failure (`Codex` undefined).
- [x] Implement `codex.go` as `claude.go`'s 1:1 sibling (~40 lines):

```go
// EnvAgentsHome overrides the Codex target's root. Libretto-only, for test
// safety — Codex itself does not read it. Same role as EnvClaudeHome.
const EnvAgentsHome = "AGENTS_HOME"

// Codex is OpenAI Codex CLI, rooted at ~/.agents. OpenCode reads the same
// skills directory, so one link here serves both tools.
type Codex struct{ root string }

func NewCodex() Codex {
	if r := os.Getenv(EnvAgentsHome); r != "" {
		return Codex{root: r}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return Codex{}
	}
	return Codex{root: filepath.Join(home, ".agents")}
}

func (c Codex) Name() string      { return "codex" }
func (c Codex) Root() string      { return c.root }
func (c Codex) Kinds() []Kind     { return []Kind{Skills} }
func (c Codex) Accepts(k Kind) bool { return accepts(c.Kinds(), k) }
func (c Codex) Dir(k Kind) string {
	if c.root == "" {
		return ""
	}
	return dirUnderRoot(c.root, k)
}
func (c Codex) Exists() bool {
	if c.root == "" {
		return false
	}
	fi, err := os.Stat(c.root)
	return err == nil && fi.IsDir()
}
```

- [x] `go test ./internal/target/ -count=1` green; commit
  `feat(target): codex target, skills only, rooted at ~/.agents`.

### Task 2: `target.Opencode` — waits on nothing (parallel with 1)

**Delta:** spec-targets · **Proof:** `TestOpencodeRootResolution`,
`TestOpencodeAcceptsOnlySkills` (+ `Exists` case inside)

**Files:** create `internal/target/opencode.go`,
`internal/target/opencode_test.go`

**Interfaces — produces:** `target.Opencode`, `target.NewOpencode()`,
`target.EnvOpencodeHome = "OPENCODE_HOME"`.

- [x] Same test set as Task 1 with root fallback
  `filepath.Join(home, ".config", "opencode")`; run, expect compile failure.
- [x] Implement `opencode.go` exactly as Task 1's shape: `Name() "opencode"`,
  env const `OPENCODE_HOME`, comment noting it is libretto-only.
- [x] Suite green; commit
  `feat(target): opencode target, skills only, rooted at ~/.config/opencode`.

### Task 3: the two new scopes through `Resolve` — waits on 1, 2

**Delta:** spec-targets · **Proof:** `TestResolveNewDestinations`,
`TestScopesNeverShareARoot`, `TestUnresolvableRootYieldsEmptyDirs`

**Files:** modify `internal/target/scope.go` (consts + `Resolve`),
`internal/target/scope_test.go`, `internal/target/target_test.go`

**Interfaces — produces:** `target.CodexScope Scope = "codex"`,
`target.OpencodeScope Scope = "opencode"`; `Resolve` returns them.

- [ ] Failing tests: `TestResolveNewDestinations` asserts
  `Resolve(CodexScope, "")` names `codex`, `Resolve(OpencodeScope, "")` names
  `opencode`, and `Resolve("nonsense", "")` still names the global target.
  Extend `TestScopesNeverShareARoot` to all four under the three env vars set
  to distinct temp dirs. Extend `TestUnresolvableRootYieldsEmptyDirs` to the
  two new targets (empty `HOME` and env unset ⇒ `Dir(Skills)==""`).
- [ ] Implement: two consts beside `GlobalScope`/`ProjectScope`; `Resolve`
  gains a `switch` with the two arms before the existing global fallback.
- [ ] `go test ./internal/target/ -count=1` green; commit
  `feat(target): resolve codex and opencode destinations`.

### Task 4: fixture safety — the three env vars everywhere — waits on 3

**Delta:** spec-targets (constraints) · closes the AGENTS.md "never write a
real home" rule for the new roots. **This lands before any cmd test can
touch the new scopes.**

**Files:** modify `cmd/libretto/helpers_test.go` (`newFixture`),
`internal/link/state_test.go` (`sandbox()`)

- [ ] `newFixture` and `sandbox()` additionally `t.Setenv(target.EnvAgentsHome,
  t.TempDir())` and `t.Setenv(target.EnvOpencodeHome, t.TempDir())`, exposing
  both dirs on the fixture as `f.Codex`, `f.Opencode` for later tasks.
- [ ] Full suite green (no behaviour change expected); commit
  `test: sandbox the codex and opencode roots in every fixture`.

### Task 5: CLI flags and destination order — waits on 4

**Delta:** spec-cli · **Proof:** `TestDestinationFlags`,
`TestDefaultScopeIsGlobal`,
`TestRememberedDestinationRecognisesNewTargets`

**Files:** modify `cmd/libretto/main.go` (`scopeFlags`, `scopeOrder`,
help literal), `cmd/libretto/remembered.go` (`rememberedScope`),
`cmd/libretto/scope_test.go`, `cmd/libretto/remembered_test.go`

**Interfaces — produces:** flags `--codex`, `--opencode`;
`scopeOrder = [global, project, codex, opencode]` (order is the strip order).

- [ ] Failing tests: `TestDestinationFlags` table-drives `scopeFlags` over
  `--codex`, `--opencode`, and every two-flag pair (`-g --codex`,
  `--codex --opencode`, …) expecting one destination or an error; remembered
  test round-trips `codex`/`opencode` and falls back to global on `garbage`.
- [ ] Implement: `scopeFlags` tracks one `chosen` string and errors on a
  second destination flag with the existing two-flag error shape;
  `scopeOrder` grows to four; `rememberedScope` recognises the two new words;
  help's flag/env literal gains `--codex`, `--opencode`, `AGENTS_HOME`,
  `OPENCODE_HOME` (defaults named, same wording as `CLAUDE_HOME`'s row).
- [ ] Suite green; commit
  `feat(cli): codex and opencode destination flags`.

### Task 6: scope isolation + skills-only behaviour end to end — waits on 5

**Delta:** spec-cli · **Proof:** `TestInstallCodexLeavesOthersAlone`,
`TestInstallOpencodeLeavesOthersAlone`, `TestHelpNamesEveryDestination`

**Files:** modify `cmd/libretto/scope_test.go`, `cmd/libretto/main_test.go`

- [ ] Failing tests, mirroring `TestInstallProjectScopeLeavesGlobalAlone`:
  `install --codex` creates links under `f.Codex/skills` only — nothing under
  the global, project or opencode roots, and **no `agents/` or `commands/`
  dir appears under `f.Codex`**; symmetrical test for opencode; help output
  contains the two flags and the two env vars.
- [ ] Run: red for the right reason (feature exists, isolation asserted) —
  then green without code changes, or fix what leaks. Commit
  `test(cli): destination isolation for codex and opencode`.

### Task 7: the four-row panel strip — waits on 5

**Delta:** spec-panel · **Proof:** `TestStripShowsAllFourDestinations`,
`TestUnconfiguredDestinationRow`,
`TestPanelPruneActsOnTheActiveDestinationOnly`,
`TestModelsRowAbsentForSkillsOnlyDestination`,
`TestFourDestinationStripGolden`

**Files:** modify `cmd/libretto/panelrun_test.go`,
`cmd/libretto/scope_test.go`, `internal/ui/panel_test.go` (+ new goldens
under `internal/ui/testdata/`)

**Interfaces — consumes:** `scopeOrder` from Task 5. `panelData` already
loops it; `nextScope` already wraps `len(Targets)`. Expect **no production
change** in `internal/ui` — this task proves it.

- [ ] Failing tests: strip renders four rows each with own state and exactly
  one active; an unconfigured codex renders `○ … not configured` and stays
  tab-selectable; panel prune with codex active touches only `f.Codex`;
  `models` menu row absent when codex is the active destination; golden
  files for the four-row strip, colour and mono
  (`go test ./internal/ui/ -update` convention if present, else write
  goldens by hand from `make preview`-style forced-colour output).
- [ ] Green; goldens reviewed by eye before committing. Commit
  `feat(panel): four destinations on the strip`.

### Task 8: README + docs — waits on 5 (parallel with 6, 7)

**Delta:** spec-cli (env table), spec-targets (STATE.md reversal, user-
authorised 2026-08-14)

**Files:** modify `README.md` (flags + env vars), `docs/STATE.md` (the
"Out of scope: targets other than Claude Code" entry), `docs/DESIGN.md`
only if its codex row example now misleads

- [ ] README: `--codex`/`--opencode` beside `--global`/`--project`; env table
  gains the two vars with defaults. `docs/STATE.md`: rewrite the out-of-scope
  entry to record the reversal, dated, citing the 2026-08-14 feasibility
  research (Codex reads `~/.agents/skills`; OpenCode reads it too plus
  `~/.claude/skills` and `~/.config/opencode`). Never delete the entry —
  a reversed decision that vanishes gets relitigated.
- [ ] `go test ./cmd/libretto/ -run TestEveryCommandIsInTheReadme -count=1`
  green (no new command, must stay green); commit
  `docs: codex and opencode destinations in README and STATE`.

### Task 9: payload prose reads in all three tools — waits on nothing (parallel)

**Delta:** spec-payload · **Proof:** `scripts/check-payload`

**Files:** modify `scripts/check-payload`, plus every `skills/**/*.md` and
`commands/**/*.md` the audit flags

- [ ] Audit: `rg -n '\bClaude\b' skills commands --glob '*.md'` — classify
  each hit: **fact** (a real path `~/.claude`, `CLAUDE.md`, `CLAUDE_HOME`,
  or the product name `Claude Code` in a sentence about Claude Code
  specifically) or **addressee** (the text means *the agent running this*).
- [ ] Rewrite the addressee hits: "the agent" / imperative voice; where a
  Claude-only mechanism is invoked (`AskUserQuestion`, `Skill` tool), keep it
  and name the generic fallback in place ("or ask in conversation when the
  native prompt does not exist"). Frontmatter `version:` bumps only where an
  instruction changed meaning, not for pronoun swaps.
- [ ] Add the check to `scripts/check-payload`, allowlist-only, no judgment:

```bash
# prose must not address Claude where it means the running agent.
# every factual use matches the allowlist; anything else fails.
prose_claude() {
  rg -n '\bClaude\b' skills commands --glob '*.md' 2>/dev/null |
    rg -v 'CLAUDE_HOME|CLAUDE\.md|~/\.claude|\.claude/|Claude Code' || true
}
```
  wired into the script's existing fail-collection pattern (match its style;
  the allowlist is the entire classification).
- [ ] Force it red once on purpose (plant `ask Claude` in a scratch skill,
  watch the gate fail, revert) — green on first run is not evidence. Then
  `scripts/check-payload` green; commit
  `feat(payload): tool-agnostic prose, enforced by check-payload`.

### Task 10: gates, boxes, close — waits on all

- [ ] All six gates: `gofmt -l .` prints nothing, `go vet ./...`,
  `go test ./... -count=1`, `scripts/check-payload`,
  `spec-drift --self-test`, `spec-drift --anchors` (every `Proof:` above now
  names a test that exists).
- [ ] Every box above marked in the commits that closed it; plan never
  batched.
