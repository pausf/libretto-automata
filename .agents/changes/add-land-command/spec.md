Targets: cli

# `libretto land` — a read-only verifier of the landing commit

A landing is one commit doing four things — final code, delta applied onto each
`Targets:` capability spec, durable decisions retired, change folder deleted — and it
failed twice by half-finishing silently. `land` reads the **staged index**, before the
mistake is history, and says which parts of that contract the commit about to exist is
missing. It verifies; it never performs. Applying a delta is a semantic merge a binary
cannot do, and the observed failures were incomplete landings, which a verifier catches.

## Outcomes

`libretto land [<change>]` inspects `git diff --cached` and exits zero only when the two
parts a binary can check both hold:

| Part | Check |
|---|---|
| 4 — folder deleted | every tracked file under the change folder is a staged deletion, and nothing under the folder survives the commit — a tracked file not staged for deletion and an untracked file still on disk both fail, by name. A commit cannot delete a folder it leaves something in. |
| 2 — delta applied | for each capability named on a `Targets:` line in any of the change's delta files, that capability's `spec.md` under the discovered specs directory is **added or modified in the same staged diff**. |

**Part 3 is deliberately not checked.** `spec-drift --retired`, inside `--anchors`, owns
it, and a second implementation is two sources of truth about what "retired" means. The
failure report says so — it attributes part 3 to `spec-drift --anchors` rather than
staying silent about a part it does not verify, so a green `land` is never read as the
whole contract passing.

**The deltas are read from `HEAD`, not from disk or the index.** A staged deletion
removes the file from the index and the landing may already have removed it from the
working tree, so `HEAD` is the one place the `Targets:` lines are guaranteed to still
be. `Targets:` lines inside fenced code blocks do not count, matching `spec-drift`'s
own fence-stripping — the two tools must not disagree about what a delta targets.

**The change name is optional.** With none given, the change is inferred from staged
deletions under a change root; with more than one folder landing, each is verified —
each has its own contract, and refusing an unusual-but-legal commit helps nobody. With
none inferable, or with a named change that has no staged deletions, the exit is
non-zero saying nothing is landing: a verifier that exits zero having verified nothing
is the silent half-landing wearing a green light.

**A folder whose deletion carries no delta passes part 2 vacuously, and the report says
so.** Deleting a queued proposal is abandoning an idea, which the flow says costs
nothing — it is not a landing and must not be failed as a broken one.

**Discovery mirrors what already exists, on both axes.** Change roots are the three
`spec-drift` reads — `.agents/changes`, `changes`, `openspec/changes` — and the specs
directory is `findSpecsDir`'s order (`.agents/specs`, `specs`, `openspec`, `docs/specs`,
`spec`). A fourth list in the binary would disagree with the others exactly when it
matters.

**A stale wiki warns and never blocks.** Stale, operationally: a marked view
(`README.md` under the wiki marker, or `wiki.html` under the HTML marker) exists in the
specs directory, a capability `spec.md` is in the staged diff, and that view is not.
That is the mechanical reading of record-work's own clause — the refreshed index rides
the same commit as the delta that changed it. The warning goes to stderr and the exit
code is untouched; a foreign view (no marker) is ignored silently, matching `wiki`'s
own precedent. No marked view present means nothing to say.

**Every failure names its part.** The non-zero exit carries one line per missing part —
which file survived, which capability spec did not move — so the fix is legible from
the output. All failures are reported in one run; stopping at the first would make the
repair iterative for no reason.

## Scope boundaries

**In:** the `land` subcommand — dispatch entry in `run`, a `usage()` line, discovery of
change roots and the specs directory, the two checks, the wiki warning, exit codes, and
the report.

**Out:**

- **performing the landing.** No write mode, no staging, no `git` mutation of any kind.
  Verify-only is forward-compatible: a write mode is a later change, not a hedge here.
- **part 1** (the final code — unverifiable mechanically) and **part 3** (owned by
  `spec-drift --retired`; `land` names the owner and checks nothing).
- **running `spec-drift`, or any gate.** `land` is one more voice, not an orchestrator.
- **a `--commit <sha>` mode.** The staged index is the default and the decision log
  records that a commit mode can be added later without moving it.
- **the payload's skill text.** record-work's guarded invocation clause is the sibling
  delta (`Targets: payload`); nothing here specifies a word of it.
- **verifying delta content** — EARS, Proof: citations, pillar structure. `--anchors`
  owns those.

## Constraints

- **Read-only, absolutely.** No file written, no index touched, no ref moved. The
  repository's bytes before and after a run are identical, and that is tested rather
  than assumed — this command runs immediately before the most destructive commit in
  the flow.
- **Stdlib plus exec'd git only.** No new dependency; git through a `gitRunner`-style
  seam as `metrics` does, so tests drive real parsing and the integration tests build
  real temporary repositories (`make test-short` skips those, per the suite's
  convention).
- **No payload required.** `land` reads the project being landed, not `~/.claude` and
  not the payload tree, so it joins the `needsPayload` exemption list (`models`,
  `update`, `loop`, `metrics`) — a machine that installed the binary and nothing else
  is exactly the machine record-work invokes it on.
- **Outside a git repository it refuses** with the same shape as `metrics`: an error,
  not an empty report.
- **Warnings on stderr, report on stdout.** The exit code means the contract, nothing
  else; the wiki warning never changes it.
- **One file per command:** `cmd/libretto/land.go`, tests beside it. Reusable logic
  only if a second consumer exists; today none does.
- **`invokedAs()` in every remedy line**, as everywhere else — the binary is linked
  under more than one name.

## Prior decisions

Carried from `decisions.md`, session 2026-08-21, as they stand:

- **Verify-only — a read-only gate over the staged index.** Applying a delta is a
  semantic merge a binary cannot do, and the observed failures were *incomplete*
  landings, which a verifier catches. If wrong: the command grows a write mode later;
  verify-only is forward-compatible. (assumed)
- **It checks parts 2 and 4 of the landing** — every file of the change folder deleted
  (no partial deletion), and each delta's `Targets:` capability spec modified in the
  same staged diff. Part 3 stays owned by `spec-drift --retired`; duplicating it in the
  binary is two sources of truth. If wrong: the check list grows, the ownership split
  does not. (assumed)
- **The staged index, pre-commit — before the mistake is history.** Change name
  optional: with none given, infer from staged change folder deletions. If wrong: a
  `--commit <sha>` mode can be added without moving the default. (assumed)
- **A stale wiki does not block — warn only**, same rule as record-work's "a missing
  convenience never blocks a landing". (assumed)
- **The payload learns about the command, minimally** — a second delta on the payload
  capability: record-work gains "where `libretto` is on PATH, run `libretto land`
  before the landing commit", same shape as the existing wiki clause, absent-binary
  path unchanged. A verifier nothing invokes verifies nothing. If wrong: the clause is
  one sentence to remove. (assumed — owned by the sibling delta; recorded here only so
  this spec is not read as the whole change.)

## Task breakdown

- [ ] `cmd/libretto/land.go`: change-root and specs-dir discovery, change inference
      from staged deletions, the named-change path
- [ ] part 4: tracked files under the folder vs staged deletions; untracked leftovers
- [ ] part 2: `Targets:` extraction from `HEAD` with fences stripped; capability spec
      presence in the staged diff; the vacuous-delta path
- [ ] the report: one line per missing part, part 3 attributed to `spec-drift`,
      exit codes
- [ ] the wiki warning: marker-owned views, stderr, exit code untouched
- [ ] wiring: `run`'s switch, `needsPayload` exemption, the `usage()` line
- [ ] tests: table-driven over real temporary git repositories, gated behind
      `testing.Short()` where they build one

## Verification criteria

These run against real temporary git repositories built per case — an index is a git
fact, and a faked one proves nothing about the command that reads the real thing.

- **When the staged index deletes every file of the change folder and modifies each
  delta's `Targets:` capability spec, `land` shall exit zero.**
  Proof: cmd/libretto/land_test.go TestLandPassesACompleteLanding
- **If a tracked file under the change folder is not a staged deletion, then `land`
  shall exit non-zero and name that file under part 4.**
  Proof: cmd/libretto/land_test.go TestLandFailsAPartialFolderDeletion
- **If an untracked file remains on disk under the change folder, then `land` shall
  exit non-zero and name it** — the commit cannot delete a folder it leaves a file in.
  Proof: cmd/libretto/land_test.go TestLandFailsAnUntrackedLeftover
- **If a delta's `Targets:` capability spec is neither added nor modified in the staged
  diff, then `land` shall exit non-zero and name the capability under part 2.**
  Proof: cmd/libretto/land_test.go TestLandFailsWhenTheCapabilitySpecDidNotMove
- **When both checked parts are missing, `land` shall name both in one run**, not stop
  at the first.
  Proof: cmd/libretto/land_test.go TestLandNamesEveryMissingPart
- **`land` shall read `Targets:` lines from `HEAD`**, so a delta already deleted from
  the working tree is still verified, **and shall ignore a `Targets:` inside a fenced
  code block**, matching spec-drift's fence-stripping.
  Proof: cmd/libretto/land_test.go TestLandReadsTargetsFromHeadAndSkipsFences
- **When several delta files in one folder carry `Targets:` lines, `land` shall check
  every named capability**, failing on any one whose spec did not move.
  Proof: cmd/libretto/land_test.go TestLandChecksEveryDeltasTarget
- **When no change name is given, `land` shall infer the change from staged deletions
  under a change root**, and shall verify each folder when more than one is landing.
  Proof: cmd/libretto/land_test.go TestLandInfersTheChangeFromStagedDeletions
  Proof: cmd/libretto/land_test.go TestLandVerifiesEveryStagedChangeFolder
- **If no staged deletion lies under any change root and no name was given, then `land`
  shall exit non-zero saying nothing is landing**; a named change with no staged
  deletions shall be refused the same way.
  Proof: cmd/libretto/land_test.go TestLandWithNothingStagedRefuses
  Proof: cmd/libretto/land_test.go TestLandRefusesANamedChangeWithNothingStaged
- **When a change folder's deletion carries no delta, `land` shall pass part 2
  vacuously and say so** — an abandoned queued proposal is not a broken landing.
  Proof: cmd/libretto/land_test.go TestLandAllowsADeltalessFolderDeletion
- **`land` shall discover change folders under `.agents/changes`, `changes` and
  `openspec/changes`** — the same three roots spec-drift reads.
  Proof: cmd/libretto/land_test.go TestLandDiscoversEveryChangeRoot
- **`land` shall not check part 3, and its report shall attribute it to
  `spec-drift --anchors`** — a landing that would fail `--retired` still passes `land`.
  Proof: cmd/libretto/land_test.go TestLandLeavesPartThreeToSpecDrift
- **When a marked wiki view exists and a capability spec is staged while that view is
  not, `land` shall warn on stderr and shall not change the exit code.**
  Proof: cmd/libretto/land_test.go TestLandWarnsOnAStaleWikiWithoutBlocking
- **When no marked view exists, or the view rides the same staged diff, `land` shall
  print no warning**, and a foreign unmarked view shall be ignored silently.
  Proof: cmd/libretto/land_test.go TestLandStaysSilentWhenTheWikiIsCurrentOrForeign
- **`land` shall write nothing** — the repository's files, index and refs are
  byte-identical before and after a run, on the passing path and on every failing one.
  Proof: cmd/libretto/land_test.go TestLandChangesNothing
- **`land` shall run with no payload installed** — it reads the project, not the tree.
  Proof: cmd/libretto/land_test.go TestLandWorksWithNoPayload
- **If the working directory is not inside a git repository, then `land` shall exit
  non-zero with an error naming git**, never an empty report.
  Proof: cmd/libretto/land_test.go TestLandOutsideARepositoryFails
- **`help` shall name `land`**, or the feature was not delivered.
  Proof: cmd/libretto/main_test.go TestHelpNamesLand
