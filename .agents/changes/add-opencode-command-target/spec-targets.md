# Delta: OpenCode accepts commands, not skills alone

Targets: targets

## Outcomes

- **`OpencodeTool` accepts two kinds: `skills` and `commands`.** In both scopes —
  `~/.config/opencode/commands` globally, `<dir>/.opencode/commands` in a project —
  because tool and scope are orthogonal axes and a kind added to a tool arrives in
  every scope that tool has.
- **The commands directory is `<root>/commands`, produced by the same
  `dirUnderRoot` every other kind uses.** No per-tool directory name, no special
  case. OpenCode globs `{command,commands}/**/*.md`, so the plural this repository
  already produces is one of the two names it looks for.
- **Installation stays a symlink.** No copy, no rewrite, no frontmatter transform.
  The loader passes `symlink: true`, so a linked file is read exactly as a real one.
- **`Agents` remains rejected** by OpenCode, so nothing that keys off
  `Accepts(Agents)` changes — the panel's `models` row stays absent there, as it is
  for Codex.
- **`CodexTool` is untouched.** It stays skills-only.

## Scope boundaries

**In:** `Opencode.Kinds()`, the criteria that assert it, and the prose in this spec
that says "skills only" about a target where it is no longer true.

**Out, named:**

- **Any content transform.** Dropping or mapping frontmatter keys, rewriting bodies,
  generating a second copy of a command. It is out because it is not needed — see
  *Prior decisions* for the evidence — and it stays out until a key exists that
  OpenCode rejects rather than ignores. What brings it back: a command in this
  payload that needs a frontmatter key OpenCode's `Info` struct refuses.
- **The singular `command/` directory.** OpenCode accepts it; so does the plural, and
  the plural is what the shared code already emits. Supporting both would mean a
  per-target directory-name override for zero behavioural difference.
- **Codex commands.** Codex's equivalent is custom prompts, which the vendor has
  deprecated. Recorded as a prior decision in this spec already; unchanged here.
- **`agents` for OpenCode.** Still queued as `add-transformed-agent-targets`, and
  still a transform rather than a link, which is the reason it is a separate change.
- **Adapting the payload's Claude-specific instruction wording** for three hosts.
  Captured as `adapt-payload-wording-to-three-hosts`; see *Prior decisions*.

## Constraints

- `dirUnderRoot` is shared by every target. A change to it would move Claude's and
  Codex's directories too, so the commands directory has to be right *because* it is
  `<root>/commands`, not in spite of it.
- Nothing outside `internal/target` may learn that OpenCode now takes commands.
  Every caller already iterates `t.Kinds()` and asks `t.Dir(k)` — verified across
  `internal/link/state.go`, `internal/link/scan.go` and `cmd/libretto/models.go`,
  which is the whole reason this change is two lines of production code.
- `OPENCODE_HOME` still overrides the root, and every test still points it at a
  temporary directory.

## Prior decisions

- **Every tool has both scopes** — user correction, 2026-08-14. A kind added to a
  tool therefore arrives in both scopes with no new rows anywhere.
- **Per-tool targets with disjoint roots** — user decision, 2026-08-14. Unchanged.
- **The directory name, the symlink support and the absence of a transform were read
  off `sst/opencode`, not the docs page** — 2026-08-14. `packages/opencode/src/config/command.ts`
  globs `{command,commands}/**/*.md` with `symlink: true`, and
  `packages/core/src/v1/config/command.ts` declares `Info` as an Effect
  `Schema.Struct` — which ignores excess properties, proven by the loader passing a
  `name` field the struct does not declare. **The docs page disagrees with the queued
  proposal about the directory name and the source settles it**: docs say plural,
  the proposal said singular, the loader takes both.
- **ASSUMED — the payload's commands install into OpenCode as they are, naming
  Claude's tool spellings.** No reading of the code settles this; under
  `/libretto-attacca` it is answered rather than asked. Our commands say things like
  `Skill(skill="find-work")`, and OpenCode's skill tool takes `name` rather than
  `skill`. The instruction is prose a model reads, not a literal call, so it
  degrades — an OpenCode session reading it invokes its own skill tool — but the
  wording is Claude's.
  **What changes if this is wrong:** if the wording has to be host-neutral *before*
  the commands may be installed, then this change blocks on
  `adapt-payload-wording-to-three-hosts` and ships after it instead of before. The
  mechanism here would not change; only the order would.
  **Why this way round:** the plumbing is two lines and verifiable; the wording is a
  content pass over the whole payload. Shipping them together would mix a target
  change with a payload rewrite, and neither would be reviewable.
- **Ceiling named:** this makes OpenCode a real commands destination with Claude's
  vocabulary inside the commands. The upgrade path is the captured change above, not
  a transform in this one.

## Task breakdown

1. `Opencode.Kinds()` returns `{Skills, Commands}`, and the type comment stops
   saying it accepts only skills.
2. Amend the criteria below and the "skills only" prose in this spec.
3. Amend the install-isolation test that asserts OpenCode creates no `commands`
   directory — the contract it guards has moved.

## Verification criteria

- **opencode serves skills at `<root>/skills` and commands at `<root>/commands`,
  and rejects agents** — the plural directory is the one the shared helper emits and
  one of the two OpenCode globs
  Proof: internal/target/opencode_test.go TestOpencodeAcceptsSkillsAndCommands
- the opencode root still honours `OPENCODE_HOME` and falls back to
  `~/.config/opencode`
  Proof: internal/target/opencode_test.go TestOpencodeRootResolution
- **every tool resolves in both scopes onto its own root — skills only for codex,
  skills and commands for opencode**
  Proof: internal/target/scope_test.go TestResolveToolScopeMatrix
- **an opencode install links a command as a symlink into `<root>/commands`**, links
  the skill too, creates no `agents` directory, and leaves every other destination
  alone
  Proof: cmd/libretto/scope_test.go TestInstallOpencodeLeavesOthersAlone

  **The symlink half is the load-bearing word**, and it was missing from the first
  draft of this criterion: "links commands" is satisfied by a copy, and a copy is the
  one outcome this change promises not to produce. `isSymlinkTo` on the command is
  what makes the promise falsifiable.

  **Atomic on purpose, not by accident:** these are four properties of one install
  run, and splitting them into four criteria would mean four tests setting up the
  same fixture to assert one thing each. The joint form is only safe because the
  proof is a single test that asserts all four — a criterion joined by *and* whose
  proof checks one half is the failure mode this shape usually hides.
