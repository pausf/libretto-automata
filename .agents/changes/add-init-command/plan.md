# Plan — add-init-command

Durable decisions: the ones in the delta's Prior decisions.

## Summary

Add `libretto init` as a self-contained subcommand in a new `cmd/libretto/init.go`: it resolves the repository root with exactly one `git rev-parse --show-toplevel` through the existing `gitRunner` seam, refuses outside a git repository or when a foreign consolidation root exists, then stats the three `.agents/` artifacts and creates only what is missing — never overwriting anything — reporting the root first and one `created`/`exists` line per artifact. Wiring is three touches (the `run` switch, `usage()`, the `needsPayload` exemption), tests live in `cmd/libretto/init_test.go` against a fake `gitRunner` and `t.TempDir()`, and the payload side is two one-sentence guarded pointers (write-spec, libretto-queue) plus a version bump, with the capability-spec non-goal rewrite deferred to the landing commit. This repository dogfoods the fix by committing `.agents/changes/.gitkeep` with the change.

## Technical context

- **Language/deps:** Go 1.26.5, stdlib only, no new dependencies. The command follows the mould of `cmd/libretto/metrics.go` — `gitRunner func(args ...string) (string, error)` at metrics.go:72, `execGit(dir)` at metrics.go:74, and the refusal wording at metrics.go:338 (`"not a git repository, or git is unavailable: %w"`).
- **Tests:** stdlib, table-driven, in `cmd/libretto/init_test.go`; anything building a real temp git repository gated behind `testing.Short()`; everything else drives a fake `gitRunner` returning fixed output plus `t.TempDir()`. One help-line test lands in `cmd/libretto/main_test.go`. No test writes outside `t.TempDir()`.
- **Six gates**: gofmt, go vet, go test, check-payload, spec-drift --self-test, spec-drift --anchors. `scripts/check-payload` must stay green after the SKILL.md frontmatter bump.
- **Blast radius: 8 files.** New: `cmd/libretto/init.go`, `cmd/libretto/init_test.go`. Modified: `cmd/libretto/main.go`, `cmd/libretto/root.go`, `cmd/libretto/main_test.go`. Payload (code-time): `skills/write-spec/SKILL.md` (one sentence + `version: "1.2"` → `"1.3"`), `commands/libretto-queue.md` (one sentence, no version frontmatter to bump). Dogfood: `.agents/changes/.gitkeep` (new, empty, committed here). The `.agents/specs/payload/spec.md` non-goal rewrite is a ninth touch but lands at landing time, not with this code.

## The approach

**Root resolution and the git refusal.** `initCmd(w io.Writer, git gitRunner) error` (name avoids the `init` keyword) calls `git("rev-parse", "--show-toplevel")` exactly once. On error, return `fmt.Errorf("not a git repository, or git is unavailable: %w", err)` — the metrics.go:338 shape; `main` prints it on stderr and exits 1. Trim the output with `strings.TrimSpace` to get the root. Everything after this line is `filepath.Join(root, ...)`, never cwd-relative — metrics.go records the cwd bug this avoids, and the subdirectory test proves it.

**The foreign-root probe.** Two package-level slices mirroring `skills/record-work/spec-drift` lines 80-82, minus the `.agents/` pair (the layout being bootstrapped):

- spec roots: `specs`, `openspec`, `docs/specs`, `spec`
- change roots: `changes`, `openspec/changes`

Each probed as `os.Stat(filepath.Join(root, p))` and `IsDir()` — a *file* named `specs` is not a consolidation root and does not refuse. Collect every hit, and if any exist, return an error naming all of them (e.g. `this project already consolidates elsewhere: specs/, changes/ — init bootstraps only the .agents/ layout`). The probe runs **before** any stat-and-create, so the refusal wins over a partial `.agents/` by ordering alone — no flag, no branch. A comment above the slices names spec-drift as the authority and states that divergence is a bug in this copy.

**The three-artifact stat-and-create.** A fixed ordered table of three artifacts: `.agents/specs/.gitkeep` (empty), `.agents/changes/.gitkeep` (empty), `.agents/lessons.md` (header). For each: `os.Stat`; if present in any form, report `exists` and touch nothing (never overwrite anything, at all — present is present); if missing, `os.MkdirAll(filepath.Dir(path), 0o755)` then `os.WriteFile(path, content, 0o644)`. Keep-files are zero-byte. The lessons header is a `const`, byte-for-byte what this repository's own `.agents/lessons.md` lines 1-4 open with:

```
# Lessons

Append-only. Entries are spent by `libretto-retro` and gain a `Resolved:` line there;
they are never edited otherwise.
```

(heading, blank line, the two-line paragraph, trailing newline). Nothing else is written — in particular no `.agents/specs/README.md`, which `libretto wiki` owns by marker.

**Report and exit semantics.** First line names the root (the standing rule is the root before acting), then one line per artifact: `created  .agents/specs/.gitkeep` / `exists   .agents/changes/.gitkeep`. Plain text, no escape codes. All-present second run prints three `exists` lines and returns nil (exit 0). Refusals return an error: stderr + exit 1 via `main`. The release notice rides the existing `defer noticeAfter` in `run` untouched.

**Wiring, three touches.**
- `main.go` `run` switch: `case "init": return initCmd(os.Stdout, execGit(projectDir))` — beside `metrics`, which already gets `execGit(projectDir)`.
- `usage()`: one line, e.g. `%[2]s init          bootstrap the .agents/ layout at the repository root`.
- `root.go` `needsPayload`: add `"init"` to the exemption list, with one comment line — `init` reads the project, not the payload, and must work on a machine with nothing installed.

Scope flags are stripped by `scopeFlags` before dispatch and mean nothing to `init` — no code needed, same as `metrics` and `loop`.

**The payload edits — what is code-time and what is landing-time.**
- *Code-time, in this change:* `skills/write-spec/SKILL.md` — one sentence appended to the "Never create a `changes/` directory" paragraph (line ~176), in the record-work guarded grammar: where `libretto` is on PATH, the sanctioned bootstrap is `libretto init`; where it is absent, the refusal stands unchanged. Frontmatter `version: "1.2"` → `"1.3"`. And `commands/libretto-queue.md` — the same guarded sentence in the "Never create the directory in a project that has none" section (line ~70); no version frontmatter exists there, so nothing bumps. Neither sentence weakens the refusal — it only points.
- *Landing-time:* the replacement non-goal bullet lands onto `.agents/specs/payload/spec.md:155` in the landing commit, like any other delta application. Not touched by this change's code commits.

**The dogfood step.** Do **not** run the built binary here — never build after changes is a repo session rule. Instead create `.agents/changes/.gitkeep` by hand in exactly the shape `init` would: an empty file at the repository root's `.agents/changes/`, committed with this change. That closes the vanishing-queue hole that lost ask-release-label-at-attacca-end (`.agents/specs/` already survives — it has content). The binary's behaviour is proven by the test suite, not by a manual run.

## The alternatives it beat

| Alternative | Why it lost |
|---|---|
| `internal/bootstrap` package | One consumer. The logic lifts into a package the day a second consumer exists; until then a package boundary is ceremony (decisions.md, session 2026-08-21). |
| Guided onboarding content (sample spec, `docs/SPEC.md`, CLAUDE.md) | Ponytail — the list grows one file at a time when someone asks. Three artifacts are the whole contract today. |
| Tracking `.agents/changes/` with a committed placeholder `README.md` instead of `.gitkeep` | A README carries content that must say something, be maintained, and — worst — collides with the marker-ownership convention `libretto wiki` uses for `.agents/specs/README.md`, inviting the same ambiguity into `changes/`. `.gitkeep` is zero bytes, universally understood as "carry this directory", and can never drift. |

## Risks

- **Foreign-root false positive on a generic `changes/` directory.** A project with an unrelated top-level `changes/` dir refuses. Mechanism: acceptable by design — spec-drift would consolidate against that same directory, so refusing is the consistent answer; the refusal *names* what it found so the user can tell, and `TestInitRefusesAForeignLayoutAndNamesIt` pins the naming.
- **Clobbering an existing `lessons.md` in a partial layout.** Mechanism: the stat-before-write per artifact; `TestInitFillsAPartialLayoutWithoutTouchingWhatExists` seeds a `lessons.md` with entries and asserts byte-identity after the run.
- **Creating at cwd instead of the repository root.** Mechanism: every path is joined onto the trimmed rev-parse output; `TestInitActsOnTheRepositoryRootNotTheCwd` runs from a subdirectory (real repo case behind `testing.Short()`) and asserts the artifacts land at the root.
- **`needsPayload` wiring missed.** If `init` is not exempted, a fresh `go install` machine — the exact machine `init` is for — gets "no payload at …" instead of a bootstrap. Mechanism: `TestInitWorksWithNoPayload` drives `run` with an empty payload root the way `models` proves its own exemption.
- **Refusal ordering inverted (fill happens beside a foreign root).** Mechanism: the probe runs before the stat-and-create in a straight line of code; `TestInitRefusalWinsOverAPartialLayout` seeds both a foreign root and a partial `.agents/` and asserts refusal plus zero writes.
- **A second git call sneaks in later.** Mechanism: the fake `gitRunner` keeps a call log; `TestInitNeverCommitsAndRunsOnlyRevParse` asserts the log is exactly `[rev-parse --show-toplevel]`.

## Validation and rollback

All six gates green: gofmt, go vet, go test (including the new `init_test.go` table and `TestHelpNamesInit` in main_test.go), check-payload (frontmatter of the bumped SKILL.md still parses, references resolve), spec-drift --self-test, spec-drift --anchors.

New tests, per the delta's Proof lines: TestInitBootstrapsAFreshProject, TestInitIsIdempotent, TestInitFillsAPartialLayoutWithoutTouchingWhatExists, TestInitRefusesAForeignLayoutAndNamesIt, TestInitRefusalWinsOverAPartialLayout, TestInitRefusesOutsideAGitRepository, TestInitActsOnTheRepositoryRootNotTheCwd, TestInitNeverCommitsAndRunsOnlyRevParse, TestInitLessonsHeaderIsTheStandardOne, TestInitWritesNoReadmeIntoSpecs, TestInitWorksWithNoPayload, TestInitNamesTheRootAndCarriesNoEscapeCodes, TestHelpNamesInit.

**Forced red before believed** — the two that could pass vacuously:
1. `TestInitFillsAPartialLayoutWithoutTouchingWhatExists` — temporarily make `initCmd` overwrite unconditionally (skip the stat) and confirm the byte-identity assertion goes red. A test whose seeded file happens to equal the written content would pass either way; forcing it proves the seed has entries the header lacks.
2. `TestInitRefusalWinsOverAPartialLayout` — temporarily invert the precedence (stat-and-create before the probe) and confirm the zero-writes assertion goes red. Otherwise the test could pass because the fixture accidentally had no foreign root.

Rollback: revert the change's commits — `init.go`/`init_test.go` disappear, the three wiring lines and two payload sentences revert, `.agents/changes/.gitkeep` is one file to delete.

## Complexity deliberately kept

The mirrored foreign-root lists in `init.go` duplicate `SPEC_ROOTS`/`CHANGE_ROOTS` from `skills/record-work/spec-drift:80-82`. Kept deliberately, with a ceiling comment naming spec-drift as the single authority and declaring any divergence a bug in the Go copy — a shared source (parsing the shell script, or a generated file) would be more machinery than two four-and-two element lists that change roughly never.
