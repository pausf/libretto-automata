# Delta — `libretto wiki`

Targets: cli

A new read-mostly subcommand in the `loop`/`metrics` family: it runs inside any
project, reads that project's consolidated specification, and writes one markdown
wiki index next to it. The wiki is generated output, never payload — nothing here
is symlinked.

## Outcomes

Running `libretto wiki` in a project that holds a consolidated specs directory
leaves a `README.md` inside that directory: an index table (capability · governs ·
what it promises) followed by one section per capability carrying its opening
paragraph, its `Governs:` line, its criteria, and a link to the full spec. The
forge renders it when anyone browses the directory. Running it again after a spec
moves refreshes the same file; running it twice with nothing changed produces
identical bytes.

## Scope boundaries

In:

- `libretto wiki` — discover the specs dir, parse `*/spec.md`, write `README.md`
  into the specs dir
- a generation marker in the output, and a refusal to touch a `README.md` that
  does not carry it
- deterministic output, capabilities in lexical order

Out, named:

- **HTML, `--serve`, any web output.** Markdown on the forge is the wiki. HTML
  comes back only if someone actually needs to browse specs off-forge.
- **install-time generation and any node script.** Wrong moment (specs move after
  install, install happens once) and a dependency outside the stack.
- **the forge's own wiki tabs** (GitHub Wiki, GitLab Wiki). Separate repos with
  separate auth; the file in-tree versions with the code.
- **multi-file wiki output.** One README until one file measurably fails —
  a project would need hundreds of criteria before it does.
- **watching / auto-regeneration by the binary.** Regeneration rides the flow
  (see the payload delta), not a daemon.

## Constraints

- Go stdlib only — no markdown parser dependency. The extraction is line-shaped
  (`Governs:` lines, `##` headings, criterion bullets), which is `bufio.Scanner`
  territory, not a parser's.
- **A criterion, for extraction, is a bullet with a `Proof:` line beneath it** —
  the same definition `spec-drift` already counts. Two tools with two definitions
  of "criterion" would disagree about what a capability promises.
- A spec with no prose paragraph before its first heading gets a section carrying
  its `Governs:` line and criteria alone — absent prose is absent, never an error.
- Specs dir discovery matches `spec-drift`'s order exactly: `.agents/specs`,
  `specs`, `openspec`, `docs/specs`, `spec` — first hit wins. Two tools with two
  discovery orders would disagree about which specification a project has.
- The ownership rule from AGENTS.md applies to the output: never overwrite
  anything the tool did not create. The marker is how the tool recognises its own.
- Tests never touch a real project: fixture specs under `t.TempDir()`.

## Prior decisions

- **Generator is the Go binary, not a skill-shipped script, not install-time node.**
  Asked 2026-08-18, user chose `libretto wiki`. Testable with `go test`, zero new
  dependencies, third member of the family (`loop`, `metrics`) that already reads a
  project's `.agents/` from the cwd. If a host without the binary ever needs the
  wiki, that is the condition that brings a skill-shipped generator back.
- **Output is markdown in the repo, not HTML, not served.** Asked 2026-08-18, user
  chose markdown. The forge already renders it; it versions and diffs with the code.
- **The wiki lands at `<specs-dir>/README.md`** so the forge shows it on directory
  browse with zero clicks. Assumed (not asked): follows directly from "markdown the
  forge renders"; what changes if wrong is one constant.
- **A generated index is not a second source of truth.** This repo's `docs/SPEC.md`
  rule ("the list lives in one place") forbids a second *hand-maintained* copy,
  because the copy nobody edits is the one that drifts. The wiki is a marked,
  regenerated view of the specs themselves — a cache, not a source. The drift risk
  is real all the same, and it is exactly why regeneration rides the landing step.
- **Regeneration rides `record-work`'s landing step** — see `payload-spec.md` in
  this change. Assumed (not asked): a generated file nobody regenerates is drift,
  and this repo's own history says a number kept in two places rots. If the user
  prefers manual-only, the payload delta is the part to drop.

## Task breakdown

1. Spec-dir discovery and spec parsing (capability name, `Governs:`, first
   paragraph, criterion bullets) — pure functions over fixture input.
2. Rendering: index table + per-capability sections + marker header, deterministic.
3. The subcommand: wiring, refusal path, exit codes, help text.
4. The payload delta: the regeneration line in `record-work`.

## Verification criteria

- The `wiki` subcommand shall discover the consolidated specs directory in the
  order `.agents/specs`, `specs`, `openspec`, `docs/specs`, `spec`, first hit wins.
  Proof: cmd/libretto/wiki_test.go TestWikiDiscoversSpecsDirInDriftOrder
- When run in a project whose specs directory holds `*/spec.md` files, `wiki`
  shall write `README.md` inside that directory with an index row and a section
  per capability, each section carrying the capability's first paragraph, its
  `Governs:` line, its criteria, and a relative link to its spec.
  Proof: cmd/libretto/wiki_test.go TestWikiWritesIndexAndSections
- The generated file shall open with a marker comment naming `libretto wiki` as
  its generator and the command that refreshes it.
  Proof: cmd/libretto/wiki_test.go TestGeneratedReadmeCarriesTheMarker
- If the target `README.md` exists and does not carry the marker, then `wiki`
  shall refuse, report the conflict, exit non-zero, and leave the file untouched.
  Proof: cmd/libretto/wiki_test.go TestWikiNeverOverwritesAHandWrittenReadme
- If no specs directory is found, or the directory found holds no `*/spec.md`,
  then `wiki` shall report which of the two it is in one line, exit non-zero, and
  write nothing — an empty specification has nothing to index, and a marker-only
  README would be noise pretending to be a wiki. (Settled 2026-08-18 from the
  cutter's finding; the builder does not decide this silently.)
  Proof: cmd/libretto/wiki_test.go TestWikiReportsNoSpecsAndExitsNonZero
- The `wiki` subcommand shall produce byte-identical output for unchanged input.
  Proof: cmd/libretto/wiki_test.go TestWikiOutputIsDeterministic
- The `wiki` subcommand shall write nothing but the one `README.md`.
  Proof: cmd/libretto/wiki_test.go TestWikiWritesNothingButTheReadme
