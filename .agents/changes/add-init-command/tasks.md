# Tasks — add-init-command

Executor: build-and-check (phase 6). One fresh session per open box; each box leaves
the tree green and mergeable on its own. Boxes are marked by the orchestrator only,
in the same commit as the work that closed them.

- [ ] **1. The `init` command, wired and proven end to end** — new
  `cmd/libretto/init.go` (`initCmd(w io.Writer, git gitRunner) error`): one
  `rev-parse --show-toplevel` through the `gitRunner` seam, trimmed root, foreign-root
  probe (spec roots `specs`/`openspec`/`docs/specs`/`spec`, change roots
  `changes`/`openspec/changes`, `IsDir()` only, probe before any write, refusal names
  every hit; ceiling comment naming spec-drift lines 80–82 as the authority), then the
  three-artifact stat-and-create (never overwrite; lessons header as the exact
  four-line `const`), report root-first then `created`/`exists` per artifact. Wiring:
  `case "init"` in `run`'s switch (`main.go`), the `usage()` line, the `init` arm in
  `needsPayload` (`root.go`). Proof in the same commit: `cmd/libretto/init_test.go`
  (fake `gitRunner` + `t.TempDir()`; real-git cases behind `testing.Short()`) and
  `TestHelpNamesInit` in `cmd/libretto/main_test.go` — all thirteen tests the plan
  names, including the two forced-red checks
  (`TestInitFillsAPartialLayoutWithoutTouchingWhatExists` with the stat skipped,
  `TestInitRefusalWinsOverAPartialLayout` with the precedence inverted; observe red,
  revert).
  Closes: all thirteen cli-delta criteria — bootstrap, idempotence, partial fill,
  foreign refusal, refusal-over-fill, non-repo refusal, root-not-cwd, single git
  invocation, lessons header, no README under specs, works with no payload,
  root-first escape-free report, help line. Each Proof: test must exist under its
  cited name.
  Waits on: nothing.

- [ ] **2. The two guarded payload pointers** — one sentence appended to
  `skills/write-spec/SKILL.md`'s "Never create a `changes/` directory" paragraph
  (~line 176) and one to `commands/libretto-queue.md`'s refusal section (~line 70),
  both in the record-work guarded grammar (`skills/record-work/SKILL.md:104-108`):
  where `libretto` is on PATH the sanctioned bootstrap is `libretto init`; where
  absent, the refusal stands unchanged. Bump write-spec frontmatter
  `version: "1.2"` → `"1.3"`; libretto-queue has no version frontmatter, bump nothing.
  Neither sentence weakens a refusal.
  Closes: the payload delta's three pointer/version criteria — proof is the sentences
  present in the two cited files and `scripts/check-payload` green (frontmatter
  parses, `name:` matches the directory).
  Waits on: nothing. Independent of box 1.

- [ ] **3. Dogfood the fix in this repository** — create
  `.agents/changes/.gitkeep` by hand (empty file at the repository root), exactly the
  shape `init` would write, and commit it with this change. Do **not** build and run
  the binary — never-build stands; the command's behaviour is box 1's test suite's to
  prove.
  Closes: the "this repository dogfoods the fix" prior decision — verified by the
  file existing empty at the root and appearing in `git ls-files`. No EARS criterion
  names this; the delta's scope and the plan's dogfood step are its source.
  Waits on: nothing. Independent of boxes 1 and 2.

Not a box: the `.agents/specs/payload/spec.md` non-goal rewrite and the cli delta's
application land in the final commit with the change-folder deletion, per the landing
contract. All six gates green before every commit.
