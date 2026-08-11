# Where this stands

Handoff note, last updated 2026-08-03. Read this first, then
[SPEC.md](SPEC.md) · [DESIGN.md](DESIGN.md) · [FLOW.md](FLOW.md).

Live task state lives in `.agents/specs/*/spec.md`, at the bottom of each spec.
**[PLAN.md](PLAN.md) is superseded and stale** — it still shows `git init`,
`internal/repo/git.go` and the README rewrite as pending, and all three landed.
Two ledgers means the one nobody reads is the one that lies. It goes, or it becomes
a pointer.

## The one thing to understand

**Both halves exist now, and neither has been exercised.**

```
skills/     20 items   (phase skills + evidence + review lenses + vendored delegates)
agents/      7 items   (spec-writer, work-reviewer, 5 review lenses)
commands/    5 items   (libretto-flow, libretto-status, libretto-review,
                        libretto-queue, libretto-next)
```

The CLI, the panel and the symlink logic are *delivery*, and they are green:
**263 test functions pass**, `gofmt` clean, `go vet` clean. The payload — the
author's own flow as skills and commands — is written and statically checked:
`scripts/check-payload` passes, `spec-drift --self-test` passes, and
`spec-drift --anchors` resolves **223 citations, file and test name**.

What none of that proves is behaviour. A skill is a prompt and a prompt is checked
by running it. **The flow has now been run end to end several times** — the payload
spec records what each run observed, and which claims are still only claims. The
queue is the newest of those: written, statically checked, never run.

**These counts go stale silently.** Nothing checks them, and they sat at 10/1/2 and
161/156 while the real numbers doubled — noticed only because a reader asked whether
the README had been updated. Re-measure before trusting them; the commands that
produce them are in `AGENTS.md` under Gates.

## The open question — answered

**What does this flow have that gentle-ai's does not?** The answer, from the
sessions that built this repo rather than from theory:

**The design came before the spec.** The order of work was explore → name → *draw
the panel* → spec → plan → implement. gentle-ai runs proposal → spec → design →
tasks. Seeing the panel is what produced the requirements: the first palette
**satisfied the spec and was unreadable** at 1.4:1 border contrast. No spec written
before seeing a terminal catches that; a WCAG measurement after seeing it caught it
immediately.

The second answer is smaller and turned out to matter more: **work is found, not
fetched.** Phase 1 asks home first — unchecked boxes in `.agents/changes/*/plan.md`
— then a tracker, then what the user said. Every change in this repository arrived
by the third route.

## What is done

| | | Proof |
|---|---|---|
| 0 | Groundwork, `git init`, 16 commits, version stamped from git | `make build` |
| 1 | `internal/target` — `Target`, Claude Code, global/project scope | `internal/target` tests |
| 2 | `internal/link` — ownership, scan, five-state classification | `internal/link` tests |
| 3 | `plan.go`, `apply.go`, prune — the code that writes to disk | `internal/link` tests |
| 4 | `internal/repo/git.go` and the update flow, wired into the CLI | `internal/repo/git_test.go` |
| 5 | subcommands: `status` `preview` `install` `uninstall` `update` `doctor` `prune` `version` `help`, TTY detection | `cmd/libretto` tests |
| 6 | logo, theme, fluid panel, model, destination switch, in-place actions and confirmation | `internal/ui` + `panelrun_test.go` |
| 7 | `Makefile`, README rewritten, licence | — |
| — | the payload: 7 phase skills, `evidence`, both commands, `spec-writer`, `spec-drift`, `check-payload` | `scripts/check-payload` |

`libretto install` installs into the global config or into the project, per the
active destination in the strip. `uninstall` undoes it. `prune` asks in place and
removes only on the second press, scoped to the destination it was planned for.

### The five states

`libretto status` reports, per item per destination:

| State | Meaning | Remedy |
|---|---|---|
| `linked` | owned symlink, right destination | none |
| `missing` | in the repo, absent from the target | `install` |
| `wrong target` | owned symlink, wrong destination | `install` repoints |
| `conflict` | something foreign in the way | **none** — reported, never touched |
| `stale` | owned symlink with no item behind it | `prune` |

`stale` deliberately subsumes "broken link": an owned link whose source was
deleted has no item behind it. Same diagnosis, same remedy, one concept.

### Ownership — the file that matters

`internal/link/own.go` decides whether a symlink belongs to this repo. Get it
wrong one way and `install` overwrites gentle-ai's work; wrong the other way and
`prune` deletes it. Three non-obvious things it handles, each with a dedicated
test:

1. **Relative links** resolve against the link's own directory, never the process
   working directory.
2. **Path normalisation** resolves symlinks before comparing. On macOS `/tmp` →
   `/private/tmp` and `t.TempDir()` sits under a symlinked `/var`, so one
   directory has two spellings and a naive compare calls our own links foreign.
   Works on paths that do not exist, because a broken link is still ours.
3. **Containment by path segment**, not string prefix. `strings.HasPrefix` would
   place `/repo-backup` inside `/repo` — that is how a tool deletes from the
   wrong tree.

## What is pending

Eleven open boxes, in the specs that own them. In the order they are worth doing:

| Where | What | Why now |
|---|---|---|
| payload | **run the flow end to end against a real task** | the only thing that checks a prompt. Start with the failure paths, which are written and never run: an unconfigured tracker, a board URL where a key was expected, a trivial task that should skip phase 2, a sub-agent that hits a question it must not answer |
| cli | `install`: present the plan before applying, as `prune` already does | the one asymmetry left in the destructive-action story — `prune` shows what it will do, `install` just does it |
| cli | `doctor`: target directory missing or unwritable; repo dirty or behind | `doctor` currently reports less than `status` can already see |
| cli | `7.3` delete `install.sh` | goes only once `libretto install` is verified against a real `~/.claude`. Still referenced by AGENTS.md and the cli spec |
| repo-sync | a fake `Git` and the flow's own tests | the interface exists precisely to make this testable, and nothing uses it that way yet |
| panel | `6.6` target-strip golden files | behaviour is tested (`TestStripRowsReportTheirOwnState`); the rendering is not pinned. No `testdata/` exists |
| panel | `6.7` `teatest` end-to-end flow | `panelrun_test.go` drives a real `tea.Program` by hand and covers the glue. `teatest` is not a dependency. **Decide whether this box still means anything** before adding one |
| cli | `5.2` `--json` for `status` and `doctor` | nothing consumes it yet. Lowest value on this list |
| cli | `5.3` the panel path under a real TTY | needs a pty; `COLUMNS` already makes layout checkable in a pipe |
| payload | an independent verifier, never run by whoever wrote the code | worth having after the flow has been run once, not before |

**`panel/6.5` was closed in this pass.** The confirmation exists, in the model
rather than as a Huh form, with four tests behind it
(`TestFooterOffersTheAnswersWhileAsking`, `TestPanelPruneConfirmsInPlace`,
`TestPanelPruneOnOnePressRemovesNothing`, `TestPanelUninstallNeedsTwoPresses`).
The box had gone stale against its own spec body, which already describes the
behaviour in prose.

## Decisions not to relitigate

- **Name and symbol.** Libretto Automata; directory `libretto-automata`; package
  `cmd/libretto`; binary `bin/libretto`.

  **The command is `libretto`.** `make link` also links `libretto-automata`, so
  somebody who remembers the project rather than the command still finds it. The
  binary prints whatever name it was invoked as, so every link describes itself
  correctly and help never names a command that does not exist.

  `lib` was the earlier choice and is **rejected**: it reads as a system directory,
  not a tool. The tempting short alternatives are worse — `maestro`, `presto`,
  `score` and `aria` all belong to real projects, and `opus` collides with a Claude
  model name in a repository full of AI tooling. `libretto` is the project's own
  noun and nobody else claims it.

  `𝄞` is README-only and must never reach a terminal (SMP plane, renders as tofu).
  `♩♪♫♬` are banned too (East Asian Ambiguous Width tears the layout).
- **Stack.** Go + Bubbletea + Lipgloss. **Huh was planned and never needed** — the
  confirmation is two states in the model and one footer line, which is less code
  than wiring a form library into a panel that already owns its own layout. opencode
  was measured, not assumed: its binary contains `opentui`, `yoga`, `solid-js`,
  `kitty`, `sixel` and no Go runtime, so it is TypeScript on its own framework.
  Higher ceiling, irrelevant for a symlink installer.
- **Symlink per item, never per directory.** `~/.claude/skills/` is shared with
  gentle-ai's sync.
- **Contrast floor 4.5:1 for text, 3:1 for borders**, enforced by a WCAG test.
  Recession is achromatic, not faded.
- **Fluid panel, 58–98 content columns**, centred on both axes when there is room.
- **One menu row, one colour.** Selected row gold end to end. Disabled rows are
  *not* dimmed — colour carries selection and nothing else.
- **A skill may only invoke what gets installed.** `install` links `skills/`,
  `agents/` and `commands/` and nothing else, so a tool a skill needs ships inside
  the skill's own directory — which is why `spec-drift` lives under
  `skills/record-work/` and not in `scripts/`.
- **Landed changes are deleted, not archived.** Git history is the archive; a
  `changes/archive/` directory nobody reopens is growth.
- **Module path is `github.com/pausf/libretto-automata`**, confirmed against the
  real repository. It had been `pausanchezv`, a guess that would have broken
  `go install` and every import for anyone but us.
- **Out of scope:** `CLAUDE.md` and `settings.json` (other tooling rewrites them);
  targets other than Claude Code; requiring `ponytail` or `caveman` (called when
  present, reported by `doctor`, never required); migrating the author's existing
  skills out of `~/.claude` (only after `libretto install` is proven).

## Spec drift — closed, and how it stays closed

`skills/record-work/spec-drift` makes the question mechanical in both directions,
and the anchor earned it immediately: the first pass cited a test
(`TestClaudeHonoursClaudeHome`) that **did not exist**, and a file-level check had
passed it. The tool verifies the test name, not just the file.

The lesson worth keeping: the drift was not created by carelessness about
documents. It was created by writing `install`, `doctor` and `prune` in one sitting
and reconciling the spec in the next — precisely the gap the same-commit rule
exists to close, and it happened anyway to the people who wrote the rule. This
handoff note is the same failure at a longer wavelength: it described an empty
`skills/` directory for three days after ten skills landed in it.
