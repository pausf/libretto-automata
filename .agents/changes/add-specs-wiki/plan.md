# Plan — add-specs-wiki

Durable decisions: the ones in Prior decisions of spec.md

## Summary

Add `libretto wiki` as a sibling of `loop` and `metrics`: a subcommand that runs in
any project, discovers its consolidated specs directory, extracts each capability's
shape with line-oriented scanning (no markdown parser), and renders one
deterministic `README.md` into that directory behind a generation marker. Then one
paragraph in `record-work`'s landing step so the file refreshes whenever the
specification actually moves. All new Go code lives in `cmd/libretto/wiki.go` plus
its test file, mirroring how `loop` and `metrics` are laid out.

## Technical context

- Go 1.26.5, stdlib only for this change — `os`, `path/filepath`, `bufio`,
  `strings`, `bytes`. No new dependency, so no ask-first conversation.
- Commands dispatch from a `case` switch in `cmd/libretto/main.go`; `loop.go` and
  `metrics.go` are the precedents for a project-reading subcommand and for test
  style (`go test`, fixtures under `t.TempDir()`, no framework).
- Gates: the six in AGENTS.md. `--anchors` currently reports seven red citations
  for `cmd/libretto/wiki_test.go` — they go green in this change, by writing the
  tests the spec named.
- Blast radius: `cmd/libretto/wiki.go` (new), `cmd/libretto/wiki_test.go` (new),
  `cmd/libretto/main.go` (dispatch + help text), `AGENTS.md` (one line in the
  commands block), `skills/record-work/SKILL.md` (the regeneration paragraph).
  Five files, two of them new.

## The approach

One file, three layers, each a pure function until the last:

1. **Discover**: walk the fixed order `.agents/specs`, `specs`, `openspec`,
   `docs/specs`, `spec` relative to the working directory; first existing
   directory wins. Same list as `spec-drift`, stated in the spec as a constraint.
2. **Extract**: for each `<dir>/<capability>/spec.md`, scan lines once. Capability
   name is the directory name. `Governs:` is the first line with that prefix. The
   intro is the first run of prose lines before the first `##` heading, skipping
   the `#` title, `Governs:`/`Targets:` lines and blanks. A criterion is a `- `
   bullet block that carries a `Proof:` line inside it — `spec-drift`'s own
   definition — captured as the bullet's text without the proof line.
3. **Render**: marker comment first, then a title, then the index table
   (capability · governs · intro's first sentence), then one `##` section per
   capability with intro, `Governs:`, criteria bullets and a relative link.
   Capabilities in lexical order; `os.ReadDir` already returns them sorted, and
   nothing passes through a map, so determinism is structural rather than patched.

Writing: if `README.md` exists and its first line is not the marker, print the
conflict and exit 1 touching nothing. Otherwise write the one file and print one
line naming it and the capability count.

The payload half is prose: `record-work`'s landing step gains a short paragraph —
where `libretto` is on PATH and a specs dir exists, run `libretto wiki` before the
landing commit so the refreshed index rides it; where absent, say the wiki may be
stale and continue.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| Skill-shipped bash generator (the `spec-drift` precedent) | Works without the binary, but structured markdown generation in bash is fragile and effectively untestable next to `go test`; the user chose the binary. Comes back if a host without the binary ever needs the wiki. |
| Node script at install time (the original idea) | Wrong moment — install happens once, specs move on every landing — and node is a new dependency in a repo whose ladder ends at "ask first". |
| Hand-maintained index updated by the flow | Zero code, but it is discipline pretending to be a mechanism; this repo's own history (the "ten specs" over eleven directories) shows how such an index rots. |
| HTML output or `wiki --serve` | The forge already renders markdown in-tree; a renderer or a server would be over-engineering rung zero. Comes back only if someone must browse specs off-forge. |
| A markdown parser dependency (goldmark) | The extraction is line-shaped — prefixes, headings, bullets. `bufio.Scanner` covers it; a parser is a dependency bought to avoid twenty lines. |

## Risks

| Risk | What catches it |
|---|---|
| Legacy specs (prose criteria, no `Proof:` bullets) render empty criteria sections | By design — absent is absent, per the spec's constraint. A fixture spec with no proofs and no intro proves the fallback renders instead of erroring: `TestWikiWritesIndexAndSections` carries that fixture. |
| In this repo, a generated `.agents/specs/README.md` reads as a rival to `docs/SPEC.md` | Settled in the spec's Prior decisions: a marked, regenerated view is a cache, not a source. The marker names its generator. |
| Regeneration instruction ignored by sessions (prose, not code) | Accepted, and named in payload-spec.md — no citation can test obedience. The mitigation is the marker telling any reader how to refresh. |
| Output accidentally nondeterministic on other platforms | No map iteration anywhere in the render path, and `TestWikiOutputIsDeterministic` compares two runs byte for byte. |

## Validation and rollback

`go test ./... -count=1` carries the seven new tests the spec already names; all
six gates before each commit, and `--anchors` flips from seven broken citations to
green when `wiki_test.go` lands. Per `skills/evidence/`, the test to force red on
purpose before believing it is `TestWikiNeverOverwritesAHandWrittenReadme` — an
overwrite bug would pass every other test silently, and this is the one that
guards the repo's oldest rule.

Rollback: one revert, nothing migrates, no generated file is load-bearing.

## Complexity deliberately kept

The refusal path (marker check before writing) is more code than an unconditional
write. It stays: "never overwrite anything the tool did not create" is the
project's most defended boundary, and the wiki writes into directories the user
also writes into.
