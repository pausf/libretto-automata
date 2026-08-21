Targets: cli

# `libretto init` — bootstrap the `.agents/` layout

One subcommand that creates the minimal `.agents/` layout in the current project —
and nothing else. Today the layout appears as a side effect of the first change, and
git drops the empty `changes/` directory on every fresh checkout, which is how a
queued idea was lost. `init` is the explicit, user-invoked owner of bootstrap; the
skills keep refusing to do it implicitly (the payload half is the sibling delta,
`spec-payload.md`).

> The brief names `cmd/libretto/land.go` as the precedent to follow. That file does
> not exist in the tree; the shape below follows the actual freshest precedents:
> `cmd/libretto/metrics.go` (the `gitRunner` seam, the not-a-git-repository refusal,
> repository-root resolution via `rev-parse --show-toplevel`) and `wiki.go`.

## Outcomes

`libretto init`, run inside a git repository, ends with these three artifacts present
at the **repository root** — never the cwd; `metrics.go` records the cwd-relative bug
this avoids:

| Artifact | Content |
|---|---|
| `.agents/specs/.gitkeep` | empty — a keep-file so git carries the directory |
| `.agents/changes/.gitkeep` | empty — same, and the fix for the vanishing queue |
| `.agents/lessons.md` | the standard append-only header, exactly as this repository's own `.agents/lessons.md` opens: a `# Lessons` heading and the "Append-only. Entries are spent by `libretto-retro` and gain a `Resolved:` line there; they are never edited otherwise." paragraph |

Only what is missing is created. Everything already present is left byte-identical.
The report names the root it acted on before acting (the standing CLI rule), then one
line per artifact saying `created` or `exists`. A second run is a no-op that reports
all three as present and exits 0.

Refusals, each an error return (stderr via `main`, exit 1):

- **outside a git repository** — the `metrics` shape: an error naming git, never an
  empty report or a silent bootstrap into a directory nothing tracks. The whole point
  of the keep-files is surviving git; a non-repo has nothing to survive.
- **a foreign layout** — any other spec-drift consolidation root present as a
  directory: specs roots `specs`, `openspec`, `docs/specs`, `spec`; change roots
  `changes`, `openspec/changes` (the lists in
  `skills/record-work/spec-drift`, lines 80–82, are the layout authority). The
  refusal names what it found. The project already has a convention, and inventing a
  second one is two sources of truth.

`init` never commits, never writes `.agents/specs/README.md` (that file is
`libretto wiki`'s, marker-owned), and never touches any install target — no
`CLAUDE_HOME`, no payload tree.

Wiring, all three in the same change:

- a `case "init"` in `run`'s switch (`cmd/libretto/main.go`), dispatching to the new
  command with the project directory and an exec'd `gitRunner`
- a usage line in `usage()`
- an `init` arm in `needsPayload` (`cmd/libretto/root.go`) — `init` reads the
  project, not the payload, so it must work on a machine with nothing installed,
  exactly like `metrics` and `loop`

## Scope boundaries

**In:** the `init` subcommand — repository-root resolution, the foreign-layout probe,
the create-what-is-missing plan, the report, the three wiring points, its tests, and
dogfooding: `.agents/changes/.gitkeep` is created and committed **in this
repository** as part of this change.

**Out (non-goals):**

- the payload side — the rewritten non-goal in `.agents/specs/payload/spec.md` and
  the guarded pointers at the two refusal sites. The sibling delta
  (`spec-payload.md`, `Targets: payload`) owns all of it.
- a sample spec, `docs/SPEC.md`, `CLAUDE.md`, `AGENTS.md`, or any content beyond the
  three artifacts. The list grows one file at a time, if ever.
- adopting or migrating a foreign layout, and any flag to do so. Refusal is the whole
  answer today; a `--layout` flag is a later change if ever asked for.
- committing, staging, or any git write. `init` runs exactly one git read
  (`rev-parse --show-toplevel`).
- any `--force`. There is no overwrite path to force.
- a panel row. `init` is a once-per-project command typed on purpose; a menu row for
  it earns its place when somebody misses it.

## Constraints

- Go 1.26.5, stdlib only — no new dependencies.
- **Never overwrite anything, at all.** The repository rule is "never overwrite what
  the tool did not create"; for `init` that tightens to never overwriting anything,
  because `init` cannot tell its own `.gitkeep` from anyone else's and does not need
  to: present is present.
- spec-drift's discovery lists are the layout authority. The foreign-root lists in
  the Go code mirror `SPEC_ROOTS`/`CHANGE_ROOTS` in
  `skills/record-work/spec-drift` (minus the `.agents/` pair, which is the layout
  being bootstrapped); a divergence is a bug in the Go copy.
- The repository root comes from `git rev-parse --show-toplevel` through the
  `gitRunner` seam (`metrics.go` line 72), injected so tests drive parsing against
  fixed output. The refusal message follows `metrics.go` line 338: an error naming
  git ("not a git repository, or git is unavailable"), wrapped so `main` prints it
  on stderr with exit 1.
- Output conventions: the root named before acting; plain text, no escape codes; the
  release notice rides `defer noticeAfter` like every other subcommand and changes
  nothing here.
- Tests: stdlib, table-driven, in `cmd/libretto/init_test.go`; anything that builds
  a real temp git repository is gated behind `testing.Short()`; everything else runs
  against a fake `gitRunner` and `t.TempDir()`. `CLAUDE_HOME` is irrelevant — `init`
  never touches a target — but no test may write outside `t.TempDir()`.
- Scope flags (`--global`/`--project`, tool flags) are stripped by `scopeFlags`
  before dispatch and have no meaning for `init`; it neither reads nor refuses them,
  the same as `metrics` and `loop`.

## Prior decisions

From `decisions.md`, session 2026-08-21 — all assumed, carried as they stand:

- **What `init` creates** — the `.agents/` layout only, minimally:
  `.agents/specs/.gitkeep`, `.agents/changes/.gitkeep`, and `.agents/lessons.md`
  with the standard append-only header. No sample spec, no `docs/SPEC.md`, no
  CLAUDE.md (out of scope by standing decision). If wrong: the list grows one file
  at a time. (assumed, 2026-08-21)
- **Foreign layout refuses** — when the project already consolidates elsewhere
  (`specs/`, `openspec/`, `docs/specs/`, `spec/`), refuse and name what was found;
  inventing a second layout is two sources of truth. A partial `.agents/` gets its
  missing pieces filled. If wrong: a flag to adopt another layout is a later change.
  (assumed, 2026-08-21)
- **Idempotent** — a second run creates nothing, overwrites nothing, reports
  everything as already present, exits 0. Never overwrite anything, ever. If wrong:
  nothing; this one has no other defensible reading. (assumed, 2026-08-21)
- **Outside a git repository: refuse**, same shape as `metrics`/`land` — the whole
  point of the keep-files is surviving git, and a non-repo has nothing to survive.
  If wrong: a warn-and-continue mode is one branch. (assumed, 2026-08-21)
- **Never commits** — it creates files and reports; committing the user's repository
  is the user's act. (assumed, 2026-08-21)
- **This repository dogfoods the fix** — `.agents/changes/.gitkeep` is created and
  committed as part of this change, closing the vanishing-queue hole that lost
  ask-release-label-at-attacca-end. `.agents/specs/` already survives. If wrong: one
  file to delete. (assumed, 2026-08-21)
- **The payload non-goal is rewritten, not deleted** — implicit creation by skills
  stays forbidden; `libretto init` becomes the named explicit owner. Owned by the
  sibling delta; recorded here only so this delta is read against it. (assumed,
  2026-08-21)

Derived in this delta from the second decision's own rationale, not a new decision:
when a foreign root and a partial `.agents/` are **both** present, the refusal wins
— filling `.agents/` beside another convention is exactly the two-sources-of-truth
the refusal exists to prevent.

## Task breakdown

Each unit independently assignable; 1 before 2; 4 any time after 1.

1. **`cmd/libretto/init.go`** — the command: resolve the repository root through a
   `gitRunner` (refuse with the git-naming error on failure); probe the foreign
   roots and refuse naming what was found; stat the three artifacts; create only the
   missing ones (`MkdirAll` + empty keep-files + the lessons header as a `const`);
   print the root, then `created`/`exists` per artifact.
2. **Wiring** — `case "init"` in `run`'s switch (`main.go`), the usage line
   (`usage()`), the `init` arm in `needsPayload` (`root.go`).
3. **Tests** — `cmd/libretto/init_test.go` (fake `gitRunner` + `t.TempDir()`;
   real-git cases behind `testing.Short()`), plus the help line in
   `cmd/libretto/main_test.go`.
4. **Dogfood** — create `.agents/changes/.gitkeep` by hand, in exactly the shape
   `init` would write it, and commit it with this change. Never a fresh build — the
   session rule stands, and the suite is what proves the command. (Amended 2026-08-21
   after a cutter finding: this line originally said "run the built binary", which
   the plan countermands.)

## Verification criteria

- **When run in a git repository with no `.agents/` layout, `init` shall create
  `.agents/specs/.gitkeep`, `.agents/changes/.gitkeep` and `.agents/lessons.md` at
  the repository root, report each as created, and exit zero.**
  Proof: cmd/libretto/init_test.go TestInitBootstrapsAFreshProject
- **When run a second time over its own output, `init` shall create nothing, modify
  nothing, report all three artifacts as already present, and exit zero.**
  Proof: cmd/libretto/init_test.go TestInitIsIdempotent
- **If a partial `.agents/` layout exists, then `init` shall create only the missing
  artifacts and shall leave every existing file byte-identical** — an existing
  `lessons.md` with entries in it is the case that matters.
  Proof: cmd/libretto/init_test.go TestInitFillsAPartialLayoutWithoutTouchingWhatExists
- **If a foreign consolidation root is present (`specs`, `openspec`, `docs/specs`,
  `spec`, `changes`, or `openspec/changes` as a directory), then `init` shall refuse,
  name what it found, create nothing, and exit non-zero.**
  Proof: cmd/libretto/init_test.go TestInitRefusesAForeignLayoutAndNamesIt
- **If a foreign root and a partial `.agents/` are both present, then `init` shall
  still refuse and create nothing** — the refusal wins over the fill.
  Proof: cmd/libretto/init_test.go TestInitRefusalWinsOverAPartialLayout
- **If the directory is not inside a git repository or git is unavailable, then
  `init` shall refuse with an error naming git and shall create nothing.**
  Proof: cmd/libretto/init_test.go TestInitRefusesOutsideAGitRepository
- **When run from a subdirectory of a repository, `init` shall act on the repository
  root, never the cwd.**
  Proof: cmd/libretto/init_test.go TestInitActsOnTheRepositoryRootNotTheCwd
- **The command shall run exactly one git invocation, `rev-parse --show-toplevel`,
  and shall never commit** — asserted through the fake `gitRunner`'s call log.
  Proof: cmd/libretto/init_test.go TestInitNeverCommitsAndRunsOnlyRevParse
- **A created `lessons.md` shall carry the standard append-only header** — the
  `# Lessons` heading and the append-only paragraph naming `libretto-retro` and
  `Resolved:`.
  Proof: cmd/libretto/init_test.go TestInitLessonsHeaderIsTheStandardOne
- **`init` shall write nothing under `.agents/specs/` beyond the keep-file** — in
  particular no `README.md`, which `libretto wiki` owns by marker.
  Proof: cmd/libretto/init_test.go TestInitWritesNoReadmeIntoSpecs
- **`init` shall work with no payload installed** — the `needsPayload` exemption,
  proven the way `models` proves its own.
  Proof: cmd/libretto/init_test.go TestInitWorksWithNoPayload
- **The report shall name the repository root it acted on before the per-artifact
  lines, and shall carry no escape codes when piped.**
  Proof: cmd/libretto/init_test.go TestInitNamesTheRootAndCarriesNoEscapeCodes
- **`help` shall list `init` with a line saying what it bootstraps.**
  Proof: cmd/libretto/main_test.go TestHelpNamesInit
