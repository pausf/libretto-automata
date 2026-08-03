# Where this stands

Handoff note, last updated 2026-07-31. Read this first, then
[SPEC.md](SPEC.md) · [DESIGN.md](DESIGN.md) · [PLAN.md](PLAN.md).

## The one thing to understand

**The mechanism is built. The product is empty.**

```
skills/     empty
agents/     empty
commands/   empty
```

The CLI, the panel, the symlink logic and 155 tests are all *delivery*. What
Libretto Automata is actually for is the payload — the author's own SDD flow, as
skills and commands. None of that is written yet.

Writing the payload is **not blocked** by any remaining CLI phase. A skill is a
markdown file with frontmatter. `install.sh` still exists and still symlinks, so
skills can be written and installed today while the Go CLI is finished in
parallel.

## The open question

**What does this flow have that gentle-ai's does not?**

If the answer turns out to be `sdd-explore`, `sdd-propose`, `sdd-spec`,
`sdd-design`, `sdd-tasks`, `sdd-apply`, `sdd-verify`, `sdd-archive` under new
names, this is a fork with a new logo and is not worth building.

One observation from the session that built this repo, offered as evidence rather
than theory:

**The design came before the spec.** The actual order of work was: explore →
name → *draw the panel* → write the spec → write the plan → implement. gentle-ai
runs proposal → spec → design → tasks. This project did not, and it was not
sloppiness: seeing the panel is what produced the requirements. Requirements
changed *after* each render — the palette, the centring, the fluid width, the
single-colour menu.

The proof that this ordering was right: the first palette **satisfied the spec
and was unreadable**. Contrast of 1.4:1 on borders. No spec written before
seeing a terminal would have caught that; a WCAG measurement after seeing it
did, immediately.

If that generalises, this flow has a phase gentle-ai's lacks, and it belongs
early. Next session starts by naming the phases.

## What is done

| Phase | | Notes |
|---|---|---|
| 0 | Groundwork | Go 1.26.5, `go mod init`, Makefile. **No `git init` yet.** |
| 1 | `internal/target` | `Target` interface + Claude Code. `CLAUDE_HOME` override is what makes the whole suite safe. |
| 2 | `internal/link` | `own.go` ownership predicate, `scan.go` enumeration, `state.go` five-state classification. Read-only. |
| 6 | `internal/ui` | Logo, theme, fluid panel, Bubbletea model. |

**155 tests green.** `gofmt` clean, `go vet` clean.

Working commands: `libretto` (TUI), `libretto status`, `libretto preview`, `libretto version`,
`libretto help`. `install`, `update`, `doctor` and `prune` are present in the menu and
refuse to run — deliberately, until Phase 3 has tests.

### The five states

`libretto status` reports, per item per target:

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

| | |
|---|---|
| **0.2** | `git init` + first commit. Nothing is versioned. This also makes the Spec-Anchored rule ("doc and code in the same commit") enforceable rather than aspirational. |
| **3.1–3.3** | `plan.go`, `apply.go`, prune. The code that writes to disk. Menu buttons stay disabled until its tests are green. |
| **4.1–4.2** | `internal/repo/git.go`, the update flow. |
| **5.1–5.3** | Remaining subcommands, `--json`, TTY-detection tests. |
| **6.5–6.7** | Huh confirmation form, target-strip goldens, `teatest` flow. |
| **7.2–7.3** | README rewrite, delete `install.sh` (only after `libretto install` is proven). |
| — | **The payload: the SDD skills and commands.** Not in PLAN.md at all — it needs its own plan. |

### Spec drift — closed

All of it. SPEC.md now describes what is true, carries a `Governs:` anchor and
`Proof:` citations, and lists what is not built yet in its own section rather than
mixing it in with what is.

Nine divergences were reconciled: the three that had accumulated since Phase 2
(`preview` missing from R8, five environment variables documented nowhere, R5
splitting broken from stale) and six created by Phase 3 (`prune --yes`, doctor's
prerequisites section, and four R5/R9 promises that turned out to be future work).

`skills/record-work/spec-drift` now makes the question mechanical in both directions, and the
anchor immediately earned it: the first pass cited a test
(`TestClaudeHonoursClaudeHome`) that **did not exist**. A file-level check had
passed it. The tool now verifies the test name, not just the file.

The lesson worth keeping: the drift was not created by carelessness about
documents. It was created by writing `install`, `doctor` and `prune` in one sitting
and reconciling the spec in the next — which is precisely the gap the same-commit
rule exists to close, and it happened anyway to the people who wrote the rule.

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
- **Stack.** Go + Bubbletea + Lipgloss + Huh. opencode was measured, not assumed:
  its binary contains `opentui`, `yoga`, `solid-js`, `kitty`, `sixel` and no Go
  runtime, so it is TypeScript on its own framework. Higher ceiling, irrelevant
  for a symlink installer.
- **Symlink per item, never per directory.** `~/.claude/skills/` is shared with
  gentle-ai's sync.
- **Contrast floor 4.5:1 for text, 3:1 for borders**, enforced by a WCAG test.
  Recession is achromatic, not faded.
- **Fluid panel, 58–98 content columns**, centred on both axes when there is room.
- **One menu row, one colour.** Selected row gold end to end. Disabled rows are
  *not* dimmed — colour carries selection and nothing else.
- **Out of scope:** `CLAUDE.md` and `settings.json` (other tooling rewrites them);
  targets other than Claude Code; migrating the author's existing skills out of
  `~/.claude` (that happens only after `libretto install` is proven).

## Loose end — closed

The module path is `github.com/pausf/libretto-automata`, confirmed against the real
repository. It had been `pausanchezv`, which was a guess, and publishing with it would
have broken `go install` and every import for anyone but us.
