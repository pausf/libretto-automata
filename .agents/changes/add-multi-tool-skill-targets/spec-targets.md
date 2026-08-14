# Delta: two skills-only targets — codex and opencode

Targets: targets

## Outcomes

- `target.Codex` exists: root `~/.agents`, overridable by `AGENTS_HOME` —
  libretto's own test-safety override, exactly `CLAUDE_HOME`'s role; Codex
  itself does not read the variable. Accepts
  exactly one kind, `Skills`, at `~/.agents/skills`. Codex CLI discovers skills
  there, and OpenCode reads the same path — one link serves both tools.
- `target.Opencode` exists: root `~/.config/opencode`, overridable by
  `OPENCODE_HOME` — the same libretto-only override, not a variable OpenCode
  reads; accepts exactly one kind, `Skills`, at `~/.config/opencode/skills`.
- `Scope` gains two values, `codex` and `opencode`, and `Resolve` returns the
  matching target for each. An unrecognised scope still resolves to global —
  that promise does not move.
- Both new targets report `Exists()` from their root directory, exactly as
  Claude does: absent means unconfigured, never an error.
- An unresolvable home yields empty directories, not a panic — same promise as
  the existing targets, extended to the new ones.

## Scope boundaries

**In:** the two implementations, their env overrides, their `Resolve` wiring,
their tests.

**Out, named:**
- Project-scoped variants of the new targets (Codex's project `.agents/skills`,
  OpenCode's project `.opencode/skill`). Deferred: OpenCode already reads a
  project `.claude/skills`, which the existing project scope serves today.
  What brings it back: a user asking for per-project Codex skills.
- Kinds other than skills for these targets. Commands and agents are the two
  queued follow-up changes (`add-opencode-command-target`,
  `add-transformed-agent-targets`), not this one.
- Codex custom prompts (`~/.codex/prompts/`). Vendor-deprecated; betting on
  them is debt on day one.
- Any registry, config file or plugin mechanism. Adding a target stays "one
  implementation of one interface".

## Constraints

- No code outside `internal/target` may derive a path from `~/.agents` or
  `~/.config/opencode` — the same rule that already protects `~/.claude`.
  Documentation strings (help's env table, the README) name the defaults,
  exactly as they name `~/.claude` today.
- A target that accepts only skills must cause no error about agents or
  commands anywhere downstream. `link.Counts` already omits rejected kinds and
  `agentsDir` already returns "" — the constraint is to keep that true.
- The env overrides exist for test safety: no test may write to a real
  `~/.agents` or `~/.config/opencode`, exactly as `CLAUDE_HOME` protects
  `~/.claude`.
- `dirUnderRoot` and `accepts` are the shared helpers; new targets use them
  rather than reimplementing.

## Prior decisions

- **Two per-tool targets with disjoint roots**, not one shared `agents` target
  — user decision, 2026-08-14. Per-tool rows keep doctor/uninstall honest;
  the shared-path fact (OpenCode also reads `~/.agents/skills`) is documented,
  not modelled.
- **A command acts on exactly one chosen destination, never all** — user
  decision, 2026-08-14 ("que el usuario pueda elegir cuál quiere, no que
  install sea para todos"). The new values extend the choice; nothing iterates.
- **New targets are single-root in this phase** (no project variant) — assumed
  from the same answer; what changes if wrong: a second `Resolve` dimension and
  more strip rows, additive later.
- **STATE.md's "Out of scope: targets other than Claude Code" is reversed** —
  user authorisation, 2026-08-14. The entry is rewritten to record the
  reversal and cite the feasibility research (memory note
  `multi-tool-target-viability`), not deleted.
- Roots verified against vendor docs and source on 2026-08-14: Codex reads
  `~/.agents/skills` (docs), OpenCode reads `~/.config/opencode/skill(s)`,
  `~/.claude/skills` and `~/.agents/skills` (sst/opencode source).

## Task breakdown

1. `internal/target/codex.go` — `Codex` struct, `NewCodex()`, `AGENTS_HOME`.
2. `internal/target/opencode.go` — `Opencode` struct, `NewOpencode()`,
   `OPENCODE_HOME`.
3. `Scope` values `CodexScope`, `OpencodeScope` and their `Resolve` arms.
4. Tests mirroring the Claude set for each target.

## Verification criteria

- Codex resolves `AGENTS_HOME` first, `~/.agents` as fallback
  Proof: internal/target/codex_test.go TestCodexRootResolution
- Codex serves skills at `<root>/skills` and rejects agents and commands
  Proof: internal/target/codex_test.go TestCodexAcceptsOnlySkills
- Opencode resolves `OPENCODE_HOME` first, `~/.config/opencode` as fallback
  Proof: internal/target/opencode_test.go TestOpencodeRootResolution
- Opencode serves skills at `<root>/skills` and rejects agents and commands
  Proof: internal/target/opencode_test.go TestOpencodeAcceptsOnlySkills
- `Resolve` returns the codex and opencode targets for their scopes, and an
  unknown scope still falls back to global
  Proof: internal/target/scope_test.go TestResolveNewDestinations
- No two destinations share a root
  Proof: internal/target/scope_test.go TestScopesNeverShareARoot
- An unresolvable home yields empty dirs for the new targets
  Proof: internal/target/target_test.go TestUnresolvableRootYieldsEmptyDirs
- both new targets report `Exists()` from their root, absent meaning
  unconfigured
  Proof: internal/target/codex_test.go TestCodexExists
