# Tasks — add-specs-wiki

Executing phase: `build-and-check` (phase 6). One fresh session per open box; each box
leaves the tree green, all six gates passing, and mergeable on its own. The capability
spec delta is not a box — it lands once, in the final landing commit.

- [x] 1. `libretto wiki`, end to end: `cmd/libretto/wiki.go` (discovery in the fixed
  order `.agents/specs`, `specs`, `openspec`, `docs/specs`, `spec`; line-oriented
  extraction of capability name, `Governs:`, intro paragraph, `Proof:`-backed
  criterion bullets; deterministic render — marker first, index table, one section
  per capability with a relative link; refusal when an existing `README.md` lacks
  the marker; one-line report, non-zero exit and no write when no specs dir is
  found or the one found holds no `*/spec.md`, naming which of the two it is),
  dispatch + help text in `cmd/libretto/main.go`, and the one-line entry in the
  AGENTS.md commands block. All fixtures under `t.TempDir()`.
  Traces to: all seven criteria in `spec.md` (delta, Targets: cli).
  Closes when: the seven named tests in `cmd/libretto/wiki_test.go` exist and pass —
  `TestWikiDiscoversSpecsDirInDriftOrder`, `TestWikiWritesIndexAndSections` (carrying
  the no-proofs/no-intro legacy fixture), `TestGeneratedReadmeCarriesTheMarker`,
  `TestWikiNeverOverwritesAHandWrittenReadme` (forced red on purpose first, then
  reverted, per the plan's validation section), `TestWikiReportsNoSpecsAndExitsNonZero`
  (covering both the missing directory and the empty one),
  `TestWikiOutputIsDeterministic`, `TestWikiWritesNothingButTheReadme` — and
  `spec-drift --anchors` flips from seven red citations to green.
  Waits on: nothing. Independent.
  Evidence: all six gates green on commit "feat(cli): libretto wiki";
  TestWikiNeverOverwritesAHandWrittenReadme forced red with the guard disabled,
  observed failing, guard restored.

- [ ] 2. The regeneration paragraph in `skills/record-work/SKILL.md`'s landing step:
  where `libretto` is on PATH and a consolidated specs dir exists, run `libretto wiki`
  before the landing commit so the refreshed index rides it; where the binary is
  absent, say the wiki may be stale and continue. Bump the skill's frontmatter
  `version:` (its instructions change meaningfully).
  Traces to: the single criterion in `payload-spec.md` (delta, Targets: payload).
  Closes when: the paragraph reads as stated, `scripts/check-payload` passes, and the
  `skills/record-work/SKILL.md` anchor resolves under `spec-drift --anchors`.
  Waits on: box 1 — a skill instructing sessions to run a subcommand that does not
  exist yet is an instruction that fails everywhere it is read.
  Evidence: —
