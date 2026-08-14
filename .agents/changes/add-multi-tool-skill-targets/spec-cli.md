# Delta: the CLI offers four destinations

Targets: cli

## Outcomes

- `--codex` and `--opencode` join `--global` and `--project`. Any two
  destination flags at once is an error, same as `--global --project` today.
- Every scope-taking subcommand (`status`, `preview`, `install`, `doctor`,
  `prune`, `uninstall`, `update`, `models`) works against the new destinations
  through the same single-`Resolve` path. No subcommand iterates destinations.
- `install --codex` links only the skills kind — agents and commands are
  simply absent from the run, not errors, and the summary says what was linked.
- The remembered panel destination recognises `codex` and `opencode`; any
  other stored word still falls back to global.
- The environment table grows by `AGENTS_HOME` and `OPENCODE_HOME`, in `help`
  and in the README — the same wording discipline as the existing five.

## Scope boundaries

**In:** the flags, the destination order the panel iterates, the env table,
help text, README rows for the new flags.

**Out, named:**
- Short flags for the new destinations. `-g`/`-p` predate them; a `-c` that
  might one day mean something else is not worth squatting now.
- Installing to several destinations in one run. One command, one destination
  — user decision, 2026-08-14.
- `doctor` growing per-tool prerequisite checks (is Codex installed, is
  OpenCode installed). It reports the destination it was handed, as today.
  What brings it back: users confused by linking into a tool they don't have.

## Constraints

- `scopeFlags` keeps returning `chosen` so the panel can tell "no flag" from
  an explicit one; the new flags flow through the same mechanism.
- The dispatch/help/README agreement gates hold: any flag the dispatch
  accepts appears in `help`, and every command stays in the README.
- Configuration remains environment-only (R11): two new variables, both with
  working defaults, nothing else exposed.

## Prior decisions

- The destination list the panel iterates (`scopeOrder`) is the one place the
  four are ordered: global, project, codex, opencode. Claude's two rows first
  because they are configured on every machine this tool has users on today.
- The env-table sentence "five environment variables" becomes "seven" — the
  count lives in `cli/spec.md`, help and README, and all three move in this
  change or the number drifts.

## Task breakdown

1. `scopeFlags` accepts the two new flags, mutual exclusion preserved.
2. `scopeOrder` gains the two destinations; `destination()` bounds follow.
3. `rememberedScope` recognises the new words.
4. Env docs: help literal, README, `cli/spec.md` table.
5. Scope-isolation tests for the new destinations.

## Verification criteria

- `--codex` and `--opencode` resolve their targets, and any two destination
  flags together is an error
  Proof: cmd/libretto/scope_test.go TestDestinationFlags
- `install --codex` links skills only and leaves every other destination alone
  Proof: cmd/libretto/scope_test.go TestInstallCodexLeavesOthersAlone
- `install --opencode` links skills only and leaves every other destination alone
  Proof: cmd/libretto/scope_test.go TestInstallOpencodeLeavesOthersAlone
- no flag still means global, for every subcommand
  Proof: cmd/libretto/scope_test.go TestDefaultScopeIsGlobal
- the remembered destination round-trips codex and opencode and falls back to
  global on anything else
  Proof: cmd/libretto/remembered_test.go TestRememberedDestinationRecognisesNewTargets
- help names the new flags and both env variables
  Proof: cmd/libretto/main_test.go TestHelpNamesEveryDestination
