# CLI

Governs: cmd/libretto/** install.sh

The command surface. Every action is reachable without a terminal, so scripts and CI
can use it, and the same binary shows a panel when a human is driving.

## Outcomes

`libretto` with a TTY and no subcommand shows the panel. Otherwise it behaves as a
plain command.

| Command | Does |
|---|---|
| `install` | link every item into every target |
| `update` | pull, relink, rebuild when Go changed, report |
| `status` | every item's state, read-only |
| `doctor` | what needs attention, plus what the payload expects on this machine |
| `prune` | show removable links; `--yes` removes them |
| `preview` | print the panel once, no TUI |
| `version`, `help` | say so |

**Exit codes carry meaning.** Non-zero whenever something was not linked or needs
attention, so a script can gate on it. A conflict counts: an item that did not get
linked is an incomplete install, whatever the reason.

**Without a TTY and without a subcommand it prints usage and exits non-zero.** It never
opens a TUI it cannot drive.

**Help names the command it was invoked as.** The binary can be linked under any name,
so help that hardcodes one lies to everyone using another.

## Scope boundaries

**In:** argument parsing, dispatch, output formatting, exit codes, environment
variables, repo-root discovery, the prerequisite report.

**Out:**

- classification and mutation — `link-state` and `linking` own those
- git — `repo-sync` owns it
- rendering — `panel` owns it
- **requiring any optional companion.** The prerequisite report says what is present.
  It never blocks.
- configuration files. There are none, deliberately.

## Constraints

**`prune` is dry by default.** With no flag it prints what it would remove and changes
nothing; `--yes` carries it out. This is stricter than confirming only when
interactive — a pipe is no reason to be less careful, and a destructive command that
acts before being asked twice eventually deletes the wrong thing.

**The prerequisite report never affects the exit code.** None of those tools is
required, and failing on an optional absence trains people to ignore the section.

**A companion counts as present wherever it legitimately lives** — the plugin tree,
`skills/`, or `commands/`. Checking only one place reports an installed tool as missing,
which is worse than not checking, because it sends people to install what they have.

**Five environment variables, each with a working default:**

| Variable | Effect |
|---|---|
| `CLAUDE_HOME` | Claude Code's root instead of `~/.claude`. What makes the test suite safe. |
| `LIBRETTO_ROOT` | the repo location, instead of deriving it from the binary |
| `LIBRETTO_ASCII` | `safe` swaps quadrant glyphs for half blocks |
| `LIBRETTO_THEME` | `dark` or `light`, instead of detecting |
| `COLUMNS` | layout width when stdout is not a terminal |

A value that never varies is not configuration, so nothing else is exposed.

**The version is stamped from git at build time, never held in a source constant.**
`make build` passes `-ldflags "-X main.version=$(git describe --tags --always --dirty)"`,
so the binary reports `v0.2.0`, or `v0.2.0-3-gabc123` past a tag, or `v0.2.0-dirty` over
uncommitted changes.

A constant drifts from the tag the moment someone forgets to bump it, and it drifts
**silently** — the binary keeps announcing a version nobody released. Asking git means
the answer cannot be wrong.

A build without those flags reports `dev`, not a version number. **A binary that cannot
prove its version says so rather than claiming one**, which is the same rule as
everywhere else here: nothing is asserted that was not observed.

**`install.sh` still exists** as the bootstrap for a machine without Go. It is a
prototype and it lacks the ownership re-check, so it must not be presented as
equivalent. It goes when `libretto install` has been verified against a real
`~/.claude`.

## Prior decisions

- The command is `libretto`; `libretto-automata` is linked too. `lib` was rejected: it
  reads as a system directory.
- `make link` symlinks rather than copies, so `make build` updates the installed
  command and no stale binary can pretend to be current. It refuses to overwrite a
  link that is not ours — the tool does not get an exception from its own rule.
- The repo root comes from the binary's compile-time location, with `LIBRETTO_ROOT` as
  the override. A tool that guesses its own repo from the working directory installs
  the wrong thing from the wrong place.
- After a rebuild the process states that it is still the old binary, rather than
  implying the upgrade took effect mid-run.

## Task breakdown

- [x] dispatch, `status`, `preview`, `version`, `help`
- [x] `install`, `doctor`, `prune --yes`
- [x] `update` composed over `repo-sync`
- [ ] 5.2 `--json` for `status` and `doctor`
- [ ] 5.3 TTY-detection tests
- [ ] `doctor`: target directory missing or unwritable
- [ ] `doctor`: repo state — uncommitted changes, behind the remote
- [ ] `install`: present the plan before applying, as `prune` already does
- [ ] 7.3 delete `install.sh`, once `libretto install` is verified against a real
      `~/.claude` with a throwaway item

## Verification criteria

**This capability has no tests of its own.** `cmd/libretto` is 662 lines and the suite
reports `no test files`. The logic it composes is covered — `linking`, `link-state`,
`targets` and `panel` carry 79 tests between them — but the composition is not: exit
codes, dispatch, flag handling and the prerequisite report are held up by nothing but
manual runs.

The gap is stated rather than hidden, and it is the largest one in the project.

Criteria owed a test:

- a conflict makes `install` exit non-zero
- `prune` without `--yes` changes nothing
- `prune --yes` removes only what the plan named
- no TTY and no subcommand prints usage and exits non-zero
- `status` output is plain, with no escape codes, when piped
- help names the invoked command rather than a fixed string
- the prerequisite report never changes the exit code
- a companion installed as a command, not a plugin, is found

The one thing verified mechanically today is the payload the CLI installs:

  Proof: scripts/check-payload
