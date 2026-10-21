# Plan — adopt-organic-driven-development

Durable decisions: yes — the Prior decisions of all three deltas (the reversals, the derived
discriminator and box definition, the two cli bullets `spec-cli.md` names)

Contract: [spec-payload.md](spec-payload.md) · [spec-cli.md](spec-cli.md) ·
[spec-readme.md](spec-readme.md) · [proposal.md](proposal.md) · [decisions.md](decisions.md)

## Summary

The contract lane moves from five documents, three fresh writers and two stops to one change
document — `proposal.md`, grown in place, written inline, read once by a fresh `spec-reviewer`,
agreed at one stop — with both layouts readable everywhere. One branch, one PR, one commit per
unit, **readers before writers**: `spec-drift` and the Go readers learn the new layout without
losing the old; then the skills, agents and commands change what they write; then `docs/FLOW.md`
and the README describe the result; then the landing applies three deltas and deletes this
folder. One merge, one tag, `release:minor`, `!` on the landing.

## Technical context

- **Stack.** Go 1.26.5, stdlib, five direct dependencies, none added. Tooling: bash, `rg`, `git`,
  `awk` (`spec-drift`, `check-payload`).
- **Gates.** The six in `AGENTS.md`. `--retired` stays inside `--anchors`; the count stays six.
- **Base.** PR #71 (`feat/add-land-command`) merges first. It overlaps this change in
  `cmd/libretto/main.go`, `main_test.go`, `skills/record-work/SKILL.md`, and the `cli` and
  `payload` capability specs — every line number the deltas cite there is pre-#71 and will shift.
- **`land` needs no code.** `headTargets` (`feat/add-land-command:cmd/libretto/land.go` ~296-328)
  reads `Targets:` from every `.md` the folder holds at `HEAD`, fences stripped. cli unit 5 is
  tests N1-N3 only.
- **The main checkout is contaminated** by `add-init-command`: untracked `cmd/libretto/init.go`,
  `init_test.go`, modified `main.go`, `main_test.go`. They break `go vet`/`go test` there. Also
  an untracked `.agents/specs/wiki.html` of unknown owner.
- **The payload is live.** `~/.claude/skills/*` and `agents/*` are symlinks into the main
  checkout, so building in that checkout rewrites the skills the building session runs on.
- **Blast radius, ~35 files plus the landing.**

| Area | Files |
|---|---|
| Skills (9 files) | `write-spec`, `write-plan`, `write-tasks`, `review-spec`, `find-work`, `record-work` (SKILL.md and `spec-drift`, 765 lines), `review-work`, `build-and-check` |
| Agents (5) | `spec-reviewer.md` new · `plan-writer.md`, `task-cutter.md` deleted · `spec-writer.md`, `work-reviewer.md` edited |
| Commands (4) | `libretto-flow`, `libretto-attacca`, `libretto-next`, `libretto-status` |
| Tooling (2) | `scripts/check-payload` (456 lines), `docs/PAYLOAD.md` regenerated |
| Docs, no criterion (3) | `docs/FLOW.md`, `docs/STATE.md`, `AGENTS.md` |
| CLI (10) | `loop.go`, `root.go`, `metrics.go`, `wiki.go`, `recommend.go` and their tests, plus `land_test.go` |
| README (2) | `README.md` (380 lines, at its ceiling), `cmd/libretto/readme_test.go` |
| Landing | three capability specs, regenerated wiki README, this folder deleted |

## The approach

**Where it is built.** A separate git worktree of this branch. The main checkout returns to `main`
with `add-init-command`'s files untouched, so the gates run in a clean tree and `~/.claude` keeps
running `main`'s payload while the branch rewrites it.

**Layer 1 — readers learn both layouts; nothing writes the new one yet.**
- *payload unit 1, `spec-drift`.* `--retired` also covers a deleted `proposal.md` whose `HEAD`
  copy carries `Targets:` outside fences; a change deleting both it and `plan.md` is reported
  once. Default mode ignores a staged `proposal.md` whose only difference is inside
  `## Checklist`. Self-test cases use real temporary repos, as the existing `land()` helper does.
- *cli units 1-2.* One section-scoped counter beside `countBoxes`, `planLine` the only box
  definition. `checklistPath` picks `proposal.md` only when its `## Checklist` holds a box, then
  `tasks.md`, then `plan.md`. The queued refusal comes before the PATH check.
- *cli unit 3, `metrics`.* A third `git log` for `proposal.md`, first, with whole-file context; a
  query counts as found only when it yields an in-scope box; `churn` tracks the section in pre-
  and post-image. `fakeGit` refuses a `proposal.md` query without `--unified=`.
- *cli unit 4, `wiki`.* `readInFlight` goes through `checklistPath` and `planLine`.
- *payload unit 2, `find-work` and `/libretto-status`.* A reader that is a skill: after the Go
  readers, before any writer.
- *cli unit 5, `land`.* Tests N1-N3 only, after #71 merges and the branch rebases.

**Layer 2 — writers move.**
- *payload unit 3.* `write-spec`, `write-plan`, `write-tasks` write their parts into `proposal.md`;
  `record-work`, `review-work`, `work-reviewer`, `build-and-check` find the contract in either
  layout.
- *payload unit 4.* Nothing launches `plan-writer` or `task-cutter`; both files deleted. Safe on
  this branch because this plan and its `tasks.md` already exist. Narrows the fan-out,
  regenerates `docs/PAYLOAD.md`, retires or rewords the listed `check-payload` rows.
- *payload unit 5 — one indivisible commit.* `agents/spec-reviewer.md`, `review-spec`'s launch,
  relay and stop, the go-ahead removed from `write-spec` and `write-tasks`, the stop tables in
  `libretto-flow`/`-attacca`/`-next`, the stop rows in `check-payload` — **and cli unit 6,
  `recommend`**, whose R2 is red until the agent file exists.
- *payload unit 8.* `spec-writer` reads `proposal.md` by path. After units 3 and 4.
- *payload unit 6.* The acronym `check-payload` row. After units 3-5.

**Layer 3 — prose.** *payload unit 7*: `docs/FLOW.md` with its one gentle-ai credit line,
`docs/STATE.md`, `AGENTS.md`. Then the readme units: guards amended and watched red, `README.md`
rewritten at net ≤ 0 lines, the moved-reasoning anchor re-checked against the final FLOW.md —
`TestMovedReasoningLandedInDocs` anchors on `more expensive than this one` in the paragraph unit 7
rewrites; if the phrase dies, the anchor is re-chosen under the readme's two anchor rules, never
deleted.

**Layer 4 — landing.** One commit: three deltas applied, their durable Prior decisions retired, the
two cli bullets moved, the wiki regenerated, this folder deleted, `!` in the subject. This folder
is **old layout** — `plan.md` and `tasks.md`, three `spec-*.md` with `Targets:`, a `proposal.md`
with none — so its own landing is the first real exercise of old-layout support through
`--retired`'s `plan.md` arm, `land`'s delta path, and the readers' `tasks.md` fallback.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| Chained PRs per capability | every merge tags a release, so main would publish readers for a layout no skill produces yet, or a README describing a flow that does not exist |
| One PR, skills first | while the branch is open, --retired and loop misread the new-layout changes, this one included |
| Building in the main checkout | its tree carries another change's broken files, and `~/.claude` links into it — the session would run the half-rewritten payload it is building |
| New `Targets:` discovery code in `land` | `headTargets` already reads every `.md` at `HEAD`; three tests, no code |
| `## Verification criteria` as a second contract arm | it let `spec-drift` and `land` give opposite verdicts on one commit |
| `proposal.md` existing selects it as the checklist | every change in both layouts has one; existence would hide every old change's `tasks.md`/`plan.md` |
| `metrics` querying `proposal.md` with a 3-line diff | a hunk does not carry its section heading, so section scope cannot be decided |

The first two are the fork's rejected rows, verbatim. "This one included" was written before the
plan observed that this change is itself old layout and read correctly by both readers today; the
row stands as the user chose it, and the risk it names is real for every new-layout change opened
while such a branch is in flight.

## Risks

| Risk | What catches it |
|---|---|
| Gates pass or fail because of `add-init-command`'s files | the worktree has none of them; `git diff --cached --name-only` is read before every commit |
| The rebase onto #71 conflicts in `record-work/SKILL.md` and both capability specs | the rebase is done in the worktree, which carries no dirty files; `check-payload` rows match phrases, not lines |
| A tree with no pre-code stop, or two | payload unit 5 is one commit; its rows — one present, two `check_absent` — prove it |
| A recommendation key names an agent that does not ship | R2 reads `../../agents`; key and agent file land in one commit |
| `metrics` shadows every old change's churn | M2, forced red under the naive order first |
| `TestMovedReasoningLandedInDocs` loses its anchor | README after FLOW.md; anchor re-chosen, never deleted |
| README past 380 lines | `TestReadmeStaysShort`, ceiling unedited |
| Dangling `~/.claude` links, unlinked `spec-reviewer` | **accepted**: `libretto install` and `libretto prune` after merge, named in the report |
| The acronym leaks | the `check-payload` row and `TestReadmeNamesNoMethodAcronym`, each watched red on a hand-broken copy |
| A contract-lane `proposal.md` without `Targets:` is invisible to `--retired` and `land` | **accepted, named ceiling** in `spec-payload.md` Prior decisions |
| `metrics` moves a landed change's numbers on real history | `go run ./cmd/libretto metrics` before and after box 3, compared per landed change (no binary built) |
| Writer commits leave the payload incoherent in between | payload units 3, 4, 5, 8 and cli unit 6 ship as one box, one commit (decided at the 5→6 seam) |

## Validation and rollback

The six gates, in the worktree, before every commit; every test in the commit of the code it proves.

**Red today, no forcing needed:** the readme guards (decision count, the non-ASCII `5→6` label,
`change document` in the walk); W3; L1-L5 and C1-C2; M1, M3, M4; the new `spec-drift` self-test
cases; the new `check-payload` rows — each observed red on today's tree before it is trusted.

**Green for the wrong reason unless forced red:**
1. N1-N3 — `land` already reads every `.md`. Exclude `proposal.md`/`spec-*.md` from
   `headTargets`, watch each fail, revert.
2. R1 — no retired key exists today. Add a `plan-writer` key, watch it fail, revert. R2 goes red
   by adding the `spec-reviewer` key before the agent file.
3. The `--retired` discriminator — removing `Targets:` from the fixture must turn exit 1 into 0.
4. M2 — must fail with `proposal.md` queried first under today's non-empty-diff test.
5. The acronym rows — zero matches today; each watched red on a hand-broken copy.
6. The old-layout proofs `spec-cli.md` lists stay green **unedited**; editing one is a reversed
   promise and the commit stops.

**Rollback is not one line, and that is the finding.** Before merge: drop the commits. After
merge: the tag cannot be recalled, so a revert is a new merge with its own tag; a change begun in
the new layout in between becomes invisible to the old readers; and users relink with
`libretto install` and `libretto prune`.

## Complexity deliberately kept

- **Two layouts in every reader.** Nothing in flight is touched. The `loop`/`wiki` fallbacks carry
  a `ponytail:` naming when they go; `metrics`' old-layout queries are permanent, said in a plain
  comment.
- **Whole-file context and per-image section tracking in `metrics`.** Named limit: a heading rename
  that moves untouched boxes is not counted.
- **One `git log` per file name, now three.** An existing Prior decision.
- **`Targets:` as the one contract discriminator**, shared by `spec-drift` and `land`.
- **`spec-writer` survives for the fan-out only.**
- **One PR, not a chain** — the larger review, because every merge tags.
- **Payload unit 5 as one commit across five or more files** — a partial stop change leaves none.
