# Brief — add-land-command

## Conventions in use

- The CLI lives in `cmd/libretto/` (one file per command, e.g. `loop.go`, `wiki.go`,
  `metrics.go`), reusable logic in `internal/`. Dispatch is in `cmd/libretto/main.go`.
- Read-only commands exit 0 and report; verifying commands exit non-zero naming what
  failed. `prune`/`uninstall` are plan-then-apply with `--yes`; a verifier has no apply
  half.
- Tests: standard library only, table-driven, real temporary git repositories for
  anything touching git (see `internal/repo/`, `spec-drift --self-test`'s `land`
  helper). `CLAUDE_HOME` guards anything touching a target — not relevant here, `land`
  never touches `~/.claude`.
- Spec criteria are EARS with `Proof:` citations naming file AND test. The cli
  capability spec is `.agents/specs/cli/spec.md` (`Governs: cmd/libretto/** install.sh`).
- Payload skills never reference `scripts/` or `docs/`; a skill may reference the
  `libretto` binary only guarded by "where `libretto` is on PATH", with an explicit
  absent-binary path (precedent: record-work's `libretto wiki` clause,
  `skills/record-work/SKILL.md` around lines 104–108).

## Constraints everyone inherits

- The landing contract is owned by `skills/record-work/SKILL.md` ("Landing a change
  consolidates it"): one commit carrying (1) final code, (2) delta applied onto each
  `Targets:` capability spec, (3) durable decisions retired into *Prior decisions*,
  (4) change folder deleted. All four or none.
- Part 3 is already gated by `spec-drift --retired` (inside `--anchors`). `libretto
  land` must NOT duplicate it — ownership stays with spec-drift.
- Layout discovery: changes live under `.agents/changes/` in this repo, but spec-drift
  discovers `changes`, `openspec/changes` too; the command should discover the same way
  or the delta must say why not.
- Go 1.26.5, no new dependencies (five direct ones today; the ladder is stdlib first).
- The binary never carries a hardcoded version; nothing here changes versioning.

## Decisions already settled

See `decisions.md` in this folder — five assumed decisions, each naming what changes if
wrong. Verify-only; checks parts 2+4 on the staged index; change name optional;
stale wiki warns, never blocks; record-work gains a guarded invocation clause.

## Vocabulary

- **landing** — the final commit of a change, the four-part contract above.
- **change folder** — `.agents/changes/<change>/` (proposal, delta spec(s), plan,
  tasks, brief, decisions).
- **delta** — a change's `spec.md` (or several) carrying `Targets: <capability>`.
- **capability spec** — `.agents/specs/<capability>/spec.md`, the consolidated truth.
- **staged index** — what `git diff --cached` shows; the commit about to exist.
- **partial landing** — a landing commit missing one of the four parts; the failure
  this command exists to catch.

## Six-pillar structure

Each delta fills: Outcomes · Scope boundaries · Constraints · Prior decisions · Task
breakdown · Verification criteria (EARS, each with `Proof:`). Delta files carry
`Targets: <capability>` on top. Two deltas in this change:

- `spec.md` → `Targets: cli` — the `libretto land` command itself.
- `spec-payload.md` → `Targets: payload` — record-work's guarded invocation clause.
  Small by design; do not grow it a task breakdown beyond the one edit.
