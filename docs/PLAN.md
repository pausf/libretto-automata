# Implementation plan

Ordered work items derived from [SPEC.md](SPEC.md). Each has a done criterion.
Requirement tags (`R1`…`R10`) point back at the spec.

**Logic before logo.** The linker and its tests come before anything that
renders. A beautiful panel wrapped around code that clobbers a file is worse
than no panel.

## Phase 0 — Groundwork

- [x] **0.1 · Docs** — `SPEC.md`, `DESIGN.md`, `PLAN.md` written.
      *Done:* every requirement has at least one scenario and appears here.
- [ ] **0.2 · `git init` and first commit** — docs and the empty item
      directories, no code yet.
      *Done:* `git log` shows one commit; `git status` clean.
- [x] **0.3 · Go toolchain** — `brew install go`. Not currently installed.
      *Done:* `go version` prints.
- [x] **0.4 · `go mod init`** — module path matching the eventual repo URL.
      *Done:* `go.mod` exists, `go build ./...` succeeds on an empty tree.

## Phase 1 — Targets (R10)

- [x] **1.1 · `internal/target/target.go`** — the `Target` interface: name, root
      path, accepted kinds, and the directory for a kind.
      *Done:* interface compiles; nothing outside this package mentions
      `~/.claude`.
- [x] **1.2 · `internal/target/claude.go`** — Claude Code. Root from
      `$CLAUDE_HOME`, falling back to `~/.claude`. Kinds: `skills`, `agents`,
      `commands`.
      *Done:* the env override is a table-driven test, so every later test can
      point at `t.TempDir()`.

The env override is not a feature for users — it is what makes the whole suite
testable without touching the real `~/.claude`.

## Phase 2 — Scanning (R4)

- [x] **2.1 · `internal/link/scan.go`** — enumerate repo items per kind; classify
      each destination as `linked`, `missing`, `wrong target`, `conflict`, and
      find `stale` owned links.
      *Done:* one table-driven test per state, all under `t.TempDir()`.
- [x] **2.2 · Ownership test** — a link resolving inside the repo is owned; one
      resolving outside is foreign. Cover relative links and a symlinked repo
      path.
      *Done:* both directions asserted. This predicate is what protects R2, so
      it gets its own test.

## Phase 3 — Linking (R1, R2, R6, R9)

- [x] **3.1 · `internal/link/plan.go`** — turn a scan into a list of intended
      actions. Pure function, no filesystem writes.
      *Done:* an already-correct tree produces an empty plan (R9).
- [x] **3.2 · `internal/link/apply.go`** — execute a plan: create links, repoint
      owned links, create missing directories, skip conflicts.
      *Done:* the spec's R1 and R2 scenarios pass as tests, including the
      foreign-directory case asserting the directory still exists afterwards.
- [x] **3.3 · Prune** — delete stale owned links only.
      *Done:* the R6 scenario passes; a test asserts a foreign entry beside two
      stale links survives.

Phase 3 is the dangerous code. It does not ship without the tests above green.

## Phase 4 — Repo state (R3, R5)

- [ ] **4.1 · `internal/repo/git.go`** — dirty check, remote presence, ahead or
      behind, pull. Shell out to `git`; keep it behind a small interface so tests
      can fake it.
      *Done:* dirty and no-remote paths tested against a fake; one integration
      test against a real temp repo, skipped under `testing.Short()`.
- [ ] **4.2 · Update flow** — pull only when clean; relink always; report added,
      removed, repaired, conflicting.
      *Done:* all three R3 scenarios pass.

## Phase 5 — Plain output (R8)

- [ ] **5.1 · `cmd/libretto` subcommands** — `install`, `update`, `status`, `doctor`,
      `prune`. Stdlib `flag`; no CLI framework for five verbs.
      *Done:* each runs headless and exits with the code the spec requires.
- [ ] **5.2 · `--json`** — for `status` and `doctor`. Stdout is JSON and nothing
      else.
      *Done:* golden-file test; `libretto status --json | jq` parses.
- [ ] **5.3 · TTY detection** — no TTY means plain text, no ANSI, no panel. No
      subcommand and no TTY means usage plus non-zero exit.
      *Done:* piped output contains no escape codes, asserted by test.

At the end of Phase 5 the tool is fully usable. Everything after this is looks.

## Phase 6 — The panel (R7)

- [x] **6.1 · `internal/ui/theme.go`** — port the palette and the nine ordered
      colouring rules from DESIGN.md, plus the per-column gradient helper.
      `docs/preview.py` is the working reference; translate it, do not reinvent
      it.
      *Done:* gradient interpolation unit-tested at both endpoints and the
      midpoint; **and** a test asserting that stripping ANSI escapes from every
      rendered row reproduces the source art byte for byte, so colour can never
      shift the geometry.
- [x] **6.2 · `internal/ui/logo.go`** — clef, wordmark, small mark, and the
      `LIBRETTO_ASCII=safe` variant. Art copied verbatim from DESIGN.md.
      *Done:* golden files for both variants; a test asserts no rune above
      U+FFFF is ever emitted, which is the `𝄞` rule made executable.
- [x] **6.3 · `internal/ui/panel.go`** — compose with `JoinHorizontal`,
      `JoinVertical`, `Place`. No hardcoded widths, no counted spaces.
      *Done:* golden files at 100, 80 and 40 columns; the 40-column case shows
      the degraded layout with intact borders.
- [x] **6.4 · `internal/ui/model.go`** — Bubbletea model: menu navigation,
      running an action, showing its result, `tea.WindowSizeMsg`.
      *Done:* `Model.Update()` tested directly for each transition, per
      `~/.claude/skills/go-testing/SKILL.md`.
- [ ] **6.5 · Confirmation form** — Huh, for `prune` and for applying a plan.
      *Done:* declining changes nothing, asserted by test.
- [ ] **6.6 · Target strip** — one row per target, `●` or `○`, with counts.
      *Done:* in the golden files.
- [ ] **6.7 · `teatest` flow** — launch, arrow to `status`, select, read output,
      quit.
      *Done:* one interactive test, skipped under `testing.Short()`.

## Phase 7 — Close out

- [x] **7.1 · `Makefile`** — `build`, `test`, `install`, `fmt`.
      *Done:* `make test` runs the suite; `make build` produces `./bin/libretto`.
- [ ] **7.2 · `README.md`** — rewrite for the real CLI. `𝄞` allowed here and
      nowhere else.
      *Done:* install and usage instructions match shipped behaviour.
- [ ] **7.3 · Delete `install.sh`** — the bash prototype.
      *Done:* deleted, only after `libretto install` is verified against a real
      `~/.claude` with a throwaway skill.

## Not in this plan

- Moving the author's existing skills out of `~/.claude` into this repo. It is a
  one-time migration, it needs its own decision, and it happens after `install`
  is proven.
- Any target other than Claude Code.
- Homebrew formula and release automation.

## Verification

```bash
cd ~/gitrepos/libretto-automata
go test ./...                    # every filesystem test under t.TempDir()
go test -short ./...             # excludes real-git and teatest flows
make build

./bin/libretto                        # the panel
LIBRETTO_ASCII=safe ./bin/libretto    # the no-quadrant fallback
./bin/libretto status --json | jq
./bin/libretto status | cat           # no TTY: plain, no colour
```

The suite must never read or write the real `~/.claude`. `CLAUDE_HOME` pointed
at `t.TempDir()` is how that is guaranteed, and it is why Phase 1 comes first.

Final manual check, once: put a throwaway skill in `skills/`, run
`./bin/libretto install` for real, confirm Claude Code loads it, then remove it and
run `./bin/libretto prune`.
