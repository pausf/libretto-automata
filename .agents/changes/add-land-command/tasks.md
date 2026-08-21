# Checklist — add-land-command

Executor: `build-and-check` (phase 6). One fresh session per open box; each box merges
alone, tree green, code and proof in the same commit. The capability spec application
(`.agents/specs/cli/spec.md`, `.agents/specs/payload/spec.md`) is **not** a box — it
lands with the final commit per the landing contract, and `Durable decisions` are
carried in the plan.

- [x] **1. The `land` verifier, end to end** — done 2026-08-21, commit 4b2b6c2;
  evidence: 18 named tests + parser table green (`go test ./... -count=1`, 55s),
  forced-red observed on `TestLandFailsAPartialFolderDeletion` and
  `TestLandChangesNothing`, gofmt/vet clean. — `cmd/libretto/land.go` +
  `cmd/libretto/land_test.go`, plus wiring in the same commit: `case "land":` in `run`
  (`cmd/libretto/main.go`), the `usage()` line, `land` in `root.go`'s `needsPayload`
  exemption, `TestHelpNamesLand` in `cmd/libretto/main_test.go`. Inside the one file:
  the `-z` diff parsing into `removed`/`touched` (unit-tested on fixed `gitRunner`
  output, incl. `R100` and spaced paths), root anchoring via `rev-parse`, discovery
  (`landChangeRoots` mirroring spec-drift with the authority comment; `findSpecsDir`
  reused from `wiki.go`), change inference and the named-change path, both refusals
  ("nothing is landing"), part 4 (`ls-tree` survivors + `ls-files --others
  --exclude-standard` leftovers), part 2 (`git show HEAD:` reads, fence-stripping
  mirroring `defenced`'s fixture, every `Targets:` capability, the vacuous-delta path),
  the grouped report with the part-3 attribution line on every run, error-return exit
  semantics, `invokedAs()` in remedies, the unborn-`HEAD` error path, and the read-only
  snapshot proof. This is one box, not several: a merged verifier that exits zero
  having checked only some of its contract is the silent half-landing the spec exists
  to catch, so no intermediate cut merges honestly.
  Closes: `TestLandPassesACompleteLanding`, `TestLandFailsAPartialFolderDeletion`
  (**seen red first**, per the plan), `TestLandFailsAnUntrackedLeftover`,
  `TestLandFailsWhenTheCapabilitySpecDidNotMove`, `TestLandNamesEveryMissingPart`,
  `TestLandReadsTargetsFromHeadAndSkipsFences`, `TestLandChecksEveryDeltasTarget`,
  `TestLandInfersTheChangeFromStagedDeletions`, `TestLandVerifiesEveryStagedChangeFolder`,
  `TestLandWithNothingStagedRefuses`, `TestLandRefusesANamedChangeWithNothingStaged`,
  `TestLandAllowsADeltalessFolderDeletion`, `TestLandDiscoversEveryChangeRoot`,
  `TestLandLeavesPartThreeToSpecDrift`, `TestLandChangesNothing` (**seen red first**),
  `TestLandWorksWithNoPayload`, `TestLandOutsideARepositoryFails`, `TestHelpNamesLand`.
  Repo-building tests gated behind `testing.Short()`; `make test-short` stays green.
  Waits on: nothing.

- [x] **2. The stale-wiki warning** — done 2026-08-21, commit e759ae1; evidence: both
  wiki tests green, forced-red observed (warning made to flip the exit code, test
  caught it), gofmt/vet/test/-short all green. — in `land.go`, reusing `ownsFile` and both markers
  from `wiki.go`: marked view exists on disk, a `<specsDir>/*/spec.md` is in `touched`,
  the view is not → one line to stderr; exit code untouched on every path; unmarked
  (foreign) view ignored silently; silence when no marked view exists or the view rides
  the same diff.
  Closes: `TestLandWarnsOnAStaleWikiWithoutBlocking`,
  `TestLandStaysSilentWhenTheWikiIsCurrentOrForeign`.
  Waits on: box 1.

- [x] **3. The record-work clause** — done 2026-08-21; evidence: clause beside the
  wiki clause with the too-old-binary case on the absent side, version 1.3 → 1.4,
  `scripts/check-payload` exit 0 read to completion. — one paragraph in the "Landing a change
  consolidates it" section of `skills/record-work/SKILL.md`, beside and in the grammar
  of the wiki clause: bold guarded instruction ("**Where `libretto` is on PATH, run
  `libretto land` before the landing commit**"), fix-and-re-run on non-zero and never
  commit past it, the too-old-binary case falling on the absent side (unverified and
  continue), plain-prose absent path, sequenced after the wiki clause; frontmatter
  `version: 1.3` → `1.4`.
  Closes: the two payload-delta criteria — the clause present in
  `skills/record-work/SKILL.md` in the shape the criterion states, and
  `scripts/check-payload` run and read green in the same commit.
  Waits on: box 1 (the clause names a command; it must exist in the tree the clause
  ships in). Independent of box 2.
