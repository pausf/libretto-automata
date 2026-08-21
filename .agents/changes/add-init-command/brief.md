# Brief — add-init-command

## Conventions in use

- CLI: one file per command in `cmd/libretto/` (`land.go` is the freshest precedent —
  small, self-contained, exec'd git via the `gitRunner` seam from `metrics.go`).
  Dispatch in `main.go`'s `run` switch, a `usage()` line, `needsPayload` exemptions in
  `root.go` for commands that read the project rather than the payload tree.
- Refusal shape outside a git repository: an error naming git, never an empty report
  (`metrics`, `land`).
- `invokedAs()` in remedy lines. Errors to stderr via `main`; exit 1 on error return.
- Tests: stdlib, table-driven, real temp git repos gated behind `testing.Short()`;
  `CLAUDE_HOME` irrelevant here (init never touches a target).
- Spec criteria in EARS with `Proof:` naming file AND test.
- Payload skills reference the binary only guarded by "where `libretto` is on PATH"
  with an explicit absent path (record-work's wiki and land clauses are precedent).
- The payload capability spec is `.agents/specs/payload/spec.md`; the cli one is
  `.agents/specs/cli/spec.md` (`Governs: cmd/libretto/** install.sh`).

## Constraints everyone inherits

- **Never overwrite anything the tool did not create** — for `init` this tightens to:
  never overwrite anything at all; only create what is missing.
- The existing non-goal at `.agents/specs/payload/spec.md:155` ("creating a `changes/`
  directory in a project that has none — a staging area nobody empties…") is being
  rewritten by this change, not violated silently: implicit creation by skills stays
  forbidden; explicit `libretto init` becomes the owner.
- The refusal sites keep refusing: `skills/write-spec/SKILL.md:176` and
  `commands/libretto-queue.md:70`. They gain only a guarded pointer to `libretto init`.
- spec-drift's discovery lists are the layout authority: specs roots `.agents/specs
  specs openspec docs/specs spec`, change roots `.agents/changes changes
  openspec/changes`. `init` must not invent a fourth layout — it bootstraps `.agents/`
  and refuses when another consolidation convention already exists.
- Go 1.26.5, no new dependencies. The binary never carries a hardcoded version.
- `libretto wiki` owns `.agents/specs/README.md` generation — `init` must NOT write a
  README there (the marker-ownership check would treat it as foreign).

## Decisions already settled

See `decisions.md` — seven assumed decisions: minimal three-artifact bootstrap
(two `.gitkeep` + `lessons.md` header), refuse on foreign layout, idempotent and
never-overwrite, refuse outside git, never commit, dogfood `.agents/changes/.gitkeep`
in this repo, non-goal rewritten with guarded pointers at the two refusal sites.

## Vocabulary

- **bootstrap** — creating the `.agents/` layout in a project that lacks it.
- **keep-file** — `.gitkeep`, an empty tracked file whose only job is making git
  carry an otherwise-empty directory.
- **foreign layout** — a consolidation convention other than `.agents/` (any other
  spec-drift root present).
- **partial layout** — `.agents/` present with some of the three artifacts missing.
- **refusal site** — a payload document that declines to create `.agents/changes/`
  implicitly (write-spec, libretto-queue).

## Six-pillar structure

Each delta fills: Outcomes · Scope boundaries · Constraints · Prior decisions · Task
breakdown · Verification criteria (EARS + `Proof:`). Two deltas:

- `spec.md` → `Targets: cli` — the `libretto init` command.
- `spec-payload.md` → `Targets: payload` — the non-goal rewritten and the two guarded
  pointers. Small by design.
