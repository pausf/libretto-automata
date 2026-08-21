# Plan — add-land-command

Durable decisions: the ones in the delta's Prior decisions (`.agents/changes/add-land-command/spec.md`, carried verbatim from `decisions.md`, session 2026-08-21 — verify-only, parts 2+4, staged index, warn-only wiki, minimal payload clause).

Contract: `.agents/changes/add-land-command/spec.md` (`Targets: cli`) and `.agents/changes/add-land-command/spec-payload.md` (`Targets: payload`). This document is the how; the what stays there.

## Summary

A new self-contained subcommand, `cmd/libretto/land.go`, that reads the staged index through an exec'd-git seam (the `gitRunner` pattern `metrics.go` already owns) and verifies the two mechanically checkable parts of the landing contract — the change folder fully deleted (part 4) and each delta's `Targets:` capability spec moved in the same staged diff (part 2) — exiting zero only when both hold, naming every miss in one run, attributing part 3 to `spec-drift --anchors`, and warning on stderr about a stale marker-owned wiki view without touching the exit code. It performs nothing and writes nothing. The sibling payload edit adds one guarded clause to `skills/record-work/SKILL.md` so something actually invokes the verifier.

## Technical context

- **Go 1.26.5, no new dependencies.** Stdlib plus `exec`'d git only; the module's five direct dependencies are untouched. `go.mod` does not change.
- **Tests live beside the command:** `cmd/libretto/land_test.go`, standard library only, table-driven. Anything that needs a real index builds a real temporary git repository and is gated behind `testing.Short()` (`make test-short` skips it), matching the suite's convention. Pure parsing helpers are additionally testable against fixed `gitRunner` output, as `metrics` does.
- **The six gates** (same list AGENTS.md names): `gofmt -l` empty, `go vet ./...`, `go test ./... -count=1`, `scripts/check-payload`, `spec-drift --self-test`, `spec-drift --anchors`. The last one is also what forces every spec criterion's `Proof:` to resolve to a real test — the new criteria in the delta cannot land unproven.
- **Generated artifacts:** none change. The wiki views (`README.md` / `wiki.html`) are only *read* for the staleness warning; `land` never regenerates them.
- **Blast radius — six files, two new:**
  - `cmd/libretto/land.go` (new — the whole command)
  - `cmd/libretto/land_test.go` (new)
  - `cmd/libretto/main.go` (one `case "land":` in `run`'s switch, one `usage()` block)
  - `cmd/libretto/root.go` (`land` joins the `needsPayload` exemption list beside `models`, `update`, `loop`, `metrics`)
  - `cmd/libretto/main_test.go` (`TestHelpNamesLand`)
  - `skills/record-work/SKILL.md` (the sibling delta's one clause + `version:` 1.3 → 1.4)
  - At landing time the deltas are applied onto `.agents/specs/cli/spec.md` and `.agents/specs/payload/spec.md` — the landing contract, not a task here.

## The approach

**Entry and seam.** `run()` gains `case "land": return land(os.Stdout, os.Stderr, args[1:], execGit(projectDir))`. `execGit` and the `gitRunner` type already exist in `metrics.go` (`cmd/libretto/metrics.go:72-81`) and are reused as-is — no new seam, no shared package. Outside a repository, the first git call fails and `land` returns `fmt.Errorf("not a git repository, or git is unavailable: %w", err)`, the same shape as `metrics`. Argument parsing mirrors `metrics`: at most one positional `<change>`, any `-flag` refused with an `invokedAs()`-bearing remedy.

**Anchoring.** Everything is asked of the repository root: `root := git("rev-parse", "--show-toplevel")`, trimmed — `metrics` records the run-from-a-subdirectory bug this avoids. All git pathspecs below are repo-root-relative; filesystem checks join onto `root`.

**Reading the staged index.** One call: `git diff --cached --name-status -z`. `-z` is not optional — it is the mechanism for paths with spaces and it disables rename-path quoting. Parsing walks NUL-separated tokens: a status token (`A`, `M`, `D`, `T`, or `R<score>`/`C<score>`) followed by one path, except `R`/`C` which are followed by two (source, then destination). The parser produces two sets: `removed` (paths gone from the index after the commit: `D` entries plus rename sources) and `touched` (paths added or modified: `A`, `M`, `T`, rename/copy destinations). Everything downstream reads these two sets, so the parsing is one function, unit-testable against fixed strings without a repository. A rename out of the change folder counts its source as removed — the contract is that nothing under the folder survives the commit, and a rename-out leaves the folder empty (assumed, logged). The base is ordinary `HEAD`; in-progress merges are out of scope (assumed, logged).

**Discovery, mirrored on both axes.** Change roots are a package-level `landChangeRoots = []string{".agents/changes", "changes", "openspec/changes"}` — spec-drift's `CHANGE_ROOTS` (`skills/record-work/spec-drift:82`) verbatim, with a comment naming the script as the authority. The specs directory reuses `findSpecsDir(root)` from `wiki.go:159` directly — it already encodes the order (`.agents/specs`, `specs`, `openspec`, `docs/specs`, `spec`), first hit wins; no fourth list is written.

**Inferring the landing change(s).** With no argument: every path in `removed` that lies under `<changeRoot>/<name>/` contributes `<name>`; the resulting set (sorted, for a deterministic report) is what gets verified — one folder or several, each on its own contract. With a name given: the folder is located as `<changeRoot>/<name>` in `HEAD` (`git ls-tree` non-empty) — HEAD, not disk, because a fully staged landing may have already removed it from the working tree. Either way, a change with zero staged deletions under its folder is refused non-zero ("nothing is landing"), and an empty inferred set is refused the same way — a verifier that verified nothing must not show green.

**Part 4 — folder deleted.** For each change folder: `git ls-tree -r --name-only -z HEAD -- <folder>/` lists every tracked file the commit starts from; each one not in `removed` is a survivor, named in the report. Then `git ls-files --others --exclude-standard -z -- <folder>/` lists untracked files still on disk under the folder; each is named too — a commit cannot delete a folder it leaves a file in. (`--exclude-standard` keeps ignored files out; an ignored scratch file does not survive into anyone's commit.)

**Part 2 — delta applied.** The change's delta files are the `.md` files `ls-tree` returned for the folder. Each is read from HEAD with `git show HEAD:<path>` — never from disk or the index, because the landing has already deleted them there. Fence-stripping mirrors spec-drift's `defenced` (`spec-drift:119`): a line matching an opening/closing fence toggles an in-fence flag and is itself dropped; only out-of-fence lines are scanned. On the surviving lines, `^\s*Targets:\s*(.+)` yields capability names (whitespace-split — one name per delta today, but the grammar tolerates more). For each capability `c`, the check is: `<specsDir>/c/spec.md` present in `touched`. A folder whose deltas carry no `Targets:` at all passes part 2 vacuously and the report says so in words ("no delta — an abandoned proposal, not a landing").

**The report and exit codes.** Stdout, grouped per change: a header line naming the change, then one line per finding — `part 4: <file> survives the commit`, `part 2: <capability> spec did not move` — all findings in one run, never stopping at the first. Every run, pass or fail, prints the attribution line: `part 3 (decisions retired) is owned by spec-drift --anchors — not checked here`. Passing runs say what passed, per part. Exit: `land` returns `nil` on a clean contract and an error naming how many parts are missing otherwise — `main` already prints that to stderr and exits 1. "Nothing is landing" is likewise an error return with an `invokedAs()` remedy. Distinct exit codes for the two failure kinds wait for a caller that needs them (assumed, logged).

**The wiki staleness warning.** Reuses `ownsFile` and the two markers from `wiki.go`. Mechanics: if `findSpecsDir` found a directory, and a path matching `<specsDir>/*/spec.md` is in `touched`, then for each view (`README.md` under the wiki marker, `wiki.html` under the HTML marker): the view exists on disk *and* opens with its marker *and* its repo-relative path is not in `touched` → one warning line to **stderr**. Exit code untouched, always. A view without the marker is foreign and silently ignored — `wiki`'s own precedent. No marked view, no spec staged, or the view riding the same diff: silence.

**Read-only, proven.** `land` calls only `rev-parse`, `diff --cached`, `ls-tree`, `ls-files --others`, `show` — all reads. The test `TestLandChangesNothing` snapshots the repository (working tree bytes, `git status --porcelain -z`, `.git/index` bytes and all refs) before and after runs on the passing path and each failing path, and asserts identity.

**Wiring.** `needsPayload` in `root.go` gains `"land"` in its exemption case with a one-line reason (the machine record-work invokes it on is exactly the machine that installed the binary and nothing else). `usage()` in `main.go` gains the `land [<change>]` line, and `main_test.go` gains `TestHelpNamesLand`.

**The payload edit.** One paragraph in the "Landing a change consolidates it" section of `skills/record-work/SKILL.md`, beside and in the grammar of the wiki clause (~lines 104–108): bold guarded instruction ("**Where `libretto` is on PATH, run `libretto land` before the landing commit**"), fix-and-re-run on non-zero, plain-prose absent path ("say the landing is unverified and continue"). Sequenced after the wiki clause so the refreshed index is staged when `land` runs. Frontmatter `version: 1.3` → `1.4`. `scripts/check-payload` proves the reference is legal.

**Build order:** (1) diff parsing + discovery + inference, unit-tested on fixed runner output; (2) part 4; (3) part 2 with HEAD-reading and fence-stripping; (4) report/exit codes/attribution; (5) wiki warning; (6) wiring (`run`, `needsPayload`, `usage`); (7) integration tests over real temp repos; (8) the SKILL.md clause. The spec's task breakdown maps onto this one-to-one.

## The alternatives it beat

| Alternative | Why it lost |
|---|---|
| `internal/land` package | One consumer; a package boundary for a single command is a premature abstraction. Lifts into `internal/` the day a second consumer exists (decision recorded 2026-08-21). |
| Payload bash (a script beside spec-drift) | The command must run on machines that installed the binary and nothing else — exactly the machine record-work invokes it on — and payload never implements delivery. |
| go-git library instead of exec'd git | No new dependency (the ladder is stdlib first); exec'd git is the existing seam (`metrics`, `wiki`), gives the real index semantics for free, and go-git's index/status behavior is a second implementation of git to trust. |
| `git status --porcelain -z` as the single read | One call instead of three, but it conflates index and worktree state and its rename encoding differs by version/config; `diff --cached` + `ls-tree` + `ls-files --others` each answer exactly one question. |

## Risks

| Risk | Mechanism that catches it |
|---|---|
| Rename (`R100`) misparsed — source not counted as removed, or dest missed as touched; a landing that renames a file out of the folder wrongly fails part 4 | Unit table over the `-z` parser with `R100\0old\0new\0` cases; integration case with a staged rename inside `TestLandFailsAPartialFolderDeletion`'s table (rename out = not a survivor) |
| Paths with spaces or non-ASCII torn by line-splitting | `-z` everywhere plus a table row whose change folder contains `a file with spaces.md`, asserted in `TestLandFailsAnUntrackedLeftover` / parser unit tests |
| `HEAD` does not exist (unborn branch, first commit) — `ls-tree`/`show` fail | The failure surfaces as the git-naming error path, never a panic or an empty green report; a dedicated table case (temp repo with no commits) asserts non-zero and an error naming git |
| Discovery lists diverge from spec-drift (`CHANGE_ROOTS`, `findSpecsDir` order) | `TestLandDiscoversEveryChangeRoot` exercises all three roots; `findSpecsDir` is *shared* with `wiki`, not copied; the roots constant carries the ceiling comment (see below) |
| Fence-stripping disagrees with `defenced` (e.g. indented fence, `Targets:` inside a fence) | `TestLandReadsTargetsFromHeadAndSkipsFences` mirrors spec-drift's own self-test fixture byte-for-byte |
| Untracked leftover detection misses an edge or reports ignored scratch files | `--exclude-standard` + a table row with a `.gitignore`d file under the folder asserting no failure, and an untracked one asserting failure, in `TestLandFailsAnUntrackedLeftover` |
| The command mutates state it promised not to (e.g. `ls-files` refreshing the index stat cache) | `TestLandChangesNothing` — byte comparison of `.git/index`, refs and worktree before/after every path |
| Wiki warning changes the exit code or fires on a foreign view | `TestLandWarnsOnAStaleWikiWithoutBlocking` asserts exit zero with the warning; `TestLandStaysSilentWhenTheWikiIsCurrentOrForeign` covers the unmarked view |

## Validation and rollback

The six gates carry it: `gofmt`, `go vet`, `go test ./... -count=1` (all the new tests the delta's criteria name, `TestLandPassesACompleteLanding` through `TestHelpNamesLand`), `scripts/check-payload` (the SKILL.md clause's legality), `spec-drift --self-test`, and `spec-drift --anchors` (every new criterion's `Proof:` must resolve — an unwritten test fails the gate, not just the intention). `make test-short` must stay green with the repo-building tests skipped.

**Forced red before believed:** `TestLandFailsAPartialFolderDeletion` — invert the part-4 membership check and watch it fail, then restore; and `TestLandChangesNothing` — temporarily make `land` touch a file and watch the snapshot comparison catch it, then remove. Both tests exist to catch a silent wrong answer, which is exactly the kind of test that can pass green while asserting nothing; a first-run green proves nothing until each has been seen red for the right reason.

Rollback: revert the commit — two new files, four small edits, no data, no migration, no state.

## Complexity deliberately kept

- **`landChangeRoots` duplicates spec-drift's `CHANGE_ROOTS`** — three strings in Go mirroring three in bash. A shared source would mean the binary parsing the script or the script reading Go; both cost more than three strings. Ceiling: unify (or add a cross-checking test that greps the script) the day the two lists diverge — the comment on the constant names spec-drift as the authority.
- **Three git calls where a clever one might do** (`diff --cached`, `ls-tree`, `ls-files --others`) — each answers one question with stable, versioned output; the merged alternative was rejected as a risk (row 4 of the alternatives table).
- **The fence-stripper is a third copy of the toggle** (spec-drift's `defenced`, wiki's implicit handling, now land's). It is five lines; extracting it buys nothing until a fourth reader appears, and the fixture-shared test is what keeps the copies honest.
