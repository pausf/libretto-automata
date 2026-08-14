# Targets

Governs: internal/target/**

What an installable destination is. One interface, four implementations, and a hard
rule that nothing outside this package may assume where any of them lives.

## Outcomes

A target declares three things and nothing else is asked of it:

- its **name**, for reporting
- its **root** path
- which **kinds** it accepts, and the directory each kind lives in

A kind knows its own shape: whether its items are directories or files, and which file
extension counts. Callers ask the kind rather than deciding for themselves.

### One (tool, scope) pair at a time, never several

A command acts on **one** destination, chosen rather than assumed. The two axes
are deliberately never fused into one list — a flat list grows as tools × scopes,
two axes as tools + scopes:

| Tool | Global root | Project root | Kinds |
|---|---|---|---|
| `ClaudeTool` | `~/.claude`, honouring `CLAUDE_HOME` | `<dir>/.claude` | skills, agents, commands |
| `CodexTool` | `~/.agents`, honouring `AGENTS_HOME` | `<dir>/.agents` | skills only |
| `OpencodeTool` | `~/.config/opencode`, honouring `OPENCODE_HOME` | `<dir>/.opencode` | skills, commands |

The two overrides on the new rows are libretto's own, for test safety — exactly
`CLAUDE_HOME`'s role. Codex and OpenCode do not read those variables. Codex CLI
discovers Claude-compatible `SKILL.md` directories under `~/.agents/skills`, and
OpenCode reads that same path plus `~/.claude/skills` and its own
`~/.config/opencode/skills` — so one codex link also serves OpenCode, a fact the
docs carry and the model does not.

**OpenCode's commands directory is `<root>/commands`, which is what `dirUnderRoot`
already produces**, so the kind needed no per-target directory name. OpenCode globs
`{command,commands}/**/*.md` with `symlink: true`, so the plural is one of the two
names it looks for and a linked file is read as a real one. Its frontmatter schema is
an Effect `Schema.Struct`, which ignores keys it does not declare — so no transform,
no key mapping, and nothing to drop. Read off `sst/opencode` on 2026-08-14, not the
docs page, which names only the plural and would have left the singular looking
required.

`Resolve(scope, dir)` returns the target for a scope. An unrecognised scope resolves
to **global**, not to nothing — a typo must not silently produce a rootless target
that writes nowhere and reports success.

The project target reports **configured** only once its root exists. A project with
no `.claude/` has not opted in, which is a state and not an error; install creates
the directory rather than demanding it.

The two Claude scopes accept the same three kinds — an item that installs in one
and silently vanishes in the other would be indistinguishable from a bug. **Codex accepts skills alone; OpenCode accepts skills and
commands.** A kind a target does not accept is absent from that destination, never an
error. Agents are the one kind neither takes, because both load a different agent
format and that needs a transform rather than a link — still queued as
`add-transformed-agent-targets`.

**A kind added to a tool arrives in both of that tool's scopes.** That is not a
separate decision, it is what orthogonal axes mean: the accepted set belongs to the
tool, and the scope only moves the root.

**There is no `All()`.** It listed every known target so callers could iterate, and
every caller did — which with two destinations would write to both on every run,
the exact surprise scopes exist to remove. The extensibility seam is `Resolve`.

## Scope boundaries

**In:** the interface, the kind vocabulary, the Claude Code implementation, root
resolution.

**Out:**

- **`CLAUDE.md` and `settings.json`.** Other tooling rewrites regions of these files.
  Linking them would start a fight this tool cannot win, and losing it means
  corrupting the user's configuration.
- project-scoped variants of the codex and opencode targets. OpenCode already reads
  a project's `.claude/skills`, which the project scope serves; per-project Codex
  skills come back the day a user asks for them.
- **the agents kind, for codex and opencode.** Both load a different agent format, so
  it needs a frontmatter transform rather than a symlink — which is a different
  mechanism, not a bigger `Kinds()`. `add-transformed-agent-targets` owns it.
- **any content transform for OpenCode commands.** Out because it is not needed:
  unknown frontmatter keys are ignored, not rejected, and this payload's commands carry
  `description:` and nothing else. What brings it back is a command needing a key
  OpenCode's schema refuses.
- **the singular `command/` directory.** OpenCode accepts it and the plural equally, and
  the plural is what the shared helper already emits. Supporting both would mean a
  per-target directory-name override for no behavioural difference.
- **commands for codex.** Its equivalent is custom prompts, which are
  vendor-deprecated; betting on them is debt on day one.
- reading or writing anything inside a target. Targets describe locations; `link-state`
  and `linking` act on them.

## Constraints

**No code outside this package may derive a path from `~/.claude`, `~/.agents` or
`~/.config/opencode`, and nothing may assume a target has all three kinds.** Both
assumptions would compile and both are wrong now that skills-only targets exist. A
target that accepts only skills must cause no error about agents. Documentation
strings — help's env table, the README — name the defaults, exactly as they always
named `~/.claude`.

**`CLAUDE_HOME` overrides the root, and `AGENTS_HOME`/`OPENCODE_HOME` do the same
for the new targets.** This is not a convenience feature — **it is what
makes the entire test suite safe.** Every test points it at a temporary directory, so
no test can touch a real configuration. Remove it and the suite becomes something you
cannot run twice.

**An unresolvable root yields empty directories, not an error and not a guess.** A
target that cannot say where it lives is a target nothing should be written to.

## Prior decisions

- Three kinds: `skills`, `agents`, `commands`. Skills are directories; agents and
  commands are `.md` files. This mirrors what Claude Code actually loads.
- Adding a target is adding one implementation of one interface — no registry file, no
  plugin mechanism, no configuration. Proven twice on 2026-08-14: codex and opencode
  each arrived as one file and one `Resolve` arm.
- **Per-tool targets with disjoint roots**, not one shared `agents` target — user
  decision, 2026-08-14. Per-tool rows keep doctor and uninstall honest; the
  shared-path fact is documented, not modelled.
- **Every tool has both scopes** — user correction, 2026-08-14, reversing the
  same day's single-root reading: "en project solo se puede instalar en claude
  no tiene sentido". Tool and scope are separate axes precisely so the matrix
  costs no new rows anywhere.
- **OpenCode's commands support was read off the source, and the docs page was
  wrong-adjacent** — 2026-08-14. `packages/opencode/src/config/command.ts` globs
  `{command,commands}/**/*.md` with `symlink: true`;
  `packages/core/src/v1/config/command.ts` declares the frontmatter as an Effect
  `Schema.Struct`, which ignores excess properties — proven by the loader passing a
  `name` field the struct does not declare. The queued proposal said `command/`
  singular and predicted a transform; the docs page says `commands/` plural only. The
  source settles both: either name works, and no transform is needed. **A vendor's
  docs page is a summary of its source, and this is the second time on this capability
  that the source disagreed with it.**
- **ASSUMED, unattended — the payload's commands install into OpenCode carrying
  Claude's tool spellings.** `/libretto-attacca`, 2026-08-14, answered rather than
  asked. Commands say things like `Skill(skill="find-work")`; OpenCode's skill tool
  takes `name`, not `skill`. It degrades rather than breaking — the line is prose a
  model reads, not a literal call — but the vocabulary is Claude's.
  **What changes if this is wrong:** the install has to wait on
  `adapt-payload-wording-to-three-hosts`, which is captured. The mechanism would not
  change; only the order would.
  **Ceiling named:** OpenCode is a real commands destination with Claude's vocabulary
  inside the commands. The upgrade path is that captured change, never a transform
  here.
- Roots verified against vendor docs and the sst/opencode source on 2026-08-14; the
  reversal of "targets other than Claude Code" is recorded, dated, in
  `docs/STATE.md`.
- `Exists()` is discovered by interface assertion rather than required of every target,
  so a target that cannot answer is assumed present instead of blocking.
- **Project scope is the working directory, never the repo root.** `repoRoot()` finds
  *libretto's own* repository — the source of the items. Conflating the two would
  install the payload into libretto's own `.claude/` no matter where it was run.

## Task breakdown

Complete. Shipped in phase 1; scopes added later; codex and opencode added 2026-08-14.

## Verification criteria

- a kind reports its own item shape, so no caller has to know it
  Proof: internal/target/target_test.go TestKindItemShape
- **the root honours `CLAUDE_HOME` and falls back to `~/.claude`**
  Proof: internal/target/target_test.go TestClaudeRootResolution
- each accepted kind resolves to its directory
  Proof: internal/target/target_test.go TestClaudeDirs
- an unknown kind is rejected rather than invented
  Proof: internal/target/target_test.go TestClaudeRejectsUnknownKind
- presence is reported honestly
  Proof: internal/target/target_test.go TestClaudeExists
- an unresolvable root yields empty directories rather than a bad path
  Proof: internal/target/target_test.go TestUnresolvableRootYieldsEmptyDirs
- the project target roots at `<dir>/.claude` and accepts the same three kinds
  Proof: internal/target/scope_test.go TestProjectRoot
- it reports configured only once its root exists
  Proof: internal/target/scope_test.go TestProjectExists
- an unrooted project is inert rather than resolving to `/.claude`
  Proof: internal/target/scope_test.go TestProjectWithNoDirectoryIsInert
- each scope resolves to its own root
  Proof: internal/target/scope_test.go TestResolveScope
- **no two destinations share a root**, or isolation is a claim with nothing behind it
  Proof: internal/target/scope_test.go TestScopesNeverShareARoot
- an unknown scope falls back to global rather than to nothing
  Proof: internal/target/scope_test.go TestUnknownScopeFallsBackToGlobal
- the codex root honours `AGENTS_HOME` and falls back to `~/.agents`
  Proof: internal/target/codex_test.go TestCodexRootResolution
- codex serves skills at `<root>/skills` and rejects agents and commands
  Proof: internal/target/codex_test.go TestCodexAcceptsOnlySkills
- codex reports presence honestly
  Proof: internal/target/codex_test.go TestCodexExists
- the opencode root honours `OPENCODE_HOME` and falls back to `~/.config/opencode`
  Proof: internal/target/opencode_test.go TestOpencodeRootResolution
- **opencode serves skills at `<root>/skills`, commands at `<root>/commands` — the
  plural the shared helper already emits, and one of the two names OpenCode globs —
  and rejects agents**
  Proof: internal/target/opencode_test.go TestOpencodeAcceptsSkillsAndCommands
- **every tool resolves in both scopes onto its own root, with the same accepted kinds
  in both scopes: skills for codex, skills and commands for opencode**
  Proof: internal/target/scope_test.go TestResolveToolScopeMatrix

**That a command actually gets linked — as a symlink, into `<root>/commands` — is
`cli`'s criterion, not this one.** This spec says where a kind lives and which target
takes it; `linking` and `cli` own what gets written there. Stating it in both would be
one promise in two documents, and the copy nobody edits is the one that reads as
authoritative.
