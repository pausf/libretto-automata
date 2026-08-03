# Targets

Governs: internal/target/**

What an installable destination is. One interface, one implementation, and a hard rule
that nothing outside this package may assume either.

## Outcomes

A target declares three things and nothing else is asked of it:

- its **name**, for reporting
- its **root** path
- which **kinds** it accepts, and the directory each kind lives in

A kind knows its own shape: whether its items are directories or files, and which file
extension counts. Callers ask the kind rather than deciding for themselves.

### Two scopes, never both at once

A command acts on **one** destination, chosen rather than assumed:

| Scope | Root |
|---|---|
| `GlobalScope` | `~/.claude`, honouring `CLAUDE_HOME` |
| `ProjectScope` | `<dir>/.claude`, where `dir` is the working directory |

`Resolve(scope, dir)` returns the target for a scope. An unrecognised scope resolves
to **global**, not to nothing — a typo must not silently produce a rootless target
that writes nowhere and reports success.

The project target reports **configured** only once its root exists. A project with
no `.claude/` has not opted in, which is a state and not an error; install creates
the directory rather than demanding it.

Both scopes accept the same three kinds. Accepting different sets would mean an item
that installs in one scope and silently vanishes in the other, which is
indistinguishable from a bug.

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
- targets other than Claude Code. The interface exists for them; the implementations
  do not, and inventing one with no user is speculation.
- reading or writing anything inside a target. Targets describe locations; `link-state`
  and `linking` act on them.

## Constraints

**Nothing outside this package may name `~/.claude`, and nothing may assume a target
has all three kinds.** Both assumptions would compile and both would be wrong the
moment a second target exists. A target that accepts only skills must cause no error
about agents.

**`CLAUDE_HOME` overrides the root.** This is not a convenience feature — **it is what
makes the entire test suite safe.** Every test points it at a temporary directory, so
no test can touch a real configuration. Remove it and the suite becomes something you
cannot run twice.

**An unresolvable root yields empty directories, not an error and not a guess.** A
target that cannot say where it lives is a target nothing should be written to.

## Prior decisions

- Three kinds: `skills`, `agents`, `commands`. Skills are directories; agents and
  commands are `.md` files. This mirrors what Claude Code actually loads.
- Adding a target is adding one implementation of one interface — no registry file, no
  plugin mechanism, no configuration.
- `Exists()` is discovered by interface assertion rather than required of every target,
  so a target that cannot answer is assumed present instead of blocking.
- **Project scope is the working directory, never the repo root.** `repoRoot()` finds
  *libretto's own* repository — the source of the items. Conflating the two would
  install the payload into libretto's own `.claude/` no matter where it was run.

## Task breakdown

Complete. Shipped in phase 1; scopes added later.

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
- **the two scopes never share a root**, or isolation is a claim with nothing behind it
  Proof: internal/target/scope_test.go TestScopesNeverShareARoot
- an unknown scope falls back to global rather than to nothing
  Proof: internal/target/scope_test.go TestUnknownScopeFallsBackToGlobal
