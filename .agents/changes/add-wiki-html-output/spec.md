# Delta — `libretto wiki --html`

Targets: cli

The wiki gains a second output: one self-contained HTML viewer, generated from the
same extraction the markdown index already uses, into the same specs directory. A
browser opens it locally; nothing serves it. This is the bring-back condition
add-specs-wiki wrote down ("HTML comes back if someone needs to browse specs
off-forge"), arriving with the user pointing at a rendered demo of this repository's
own specs.

## Outcomes

Running `libretto wiki --html` in a project with a consolidated specs directory
leaves a `wiki.html` beside the generated `README.md`: a sidebar index of
capabilities with criteria counts, a section per capability carrying its intro,
`Governs:` line and criteria, and a client-side filter over the criteria — usable
offline, from `file://`, with no build step and no server. Running plain
`libretto wiki` afterwards refreshes both marked views in one pass, so the landing
regeneration that record-work already mandates keeps the HTML fresh without the
skill growing a second instruction.

## Scope boundaries

In:

- `--html` on the `wiki` subcommand: render the extracted capabilities through an
  embedded HTML template into `<specs-dir>/wiki.html`
- the same ownership model as the markdown output: a marker line, a refusal to
  touch a `wiki.html` that does not carry it
- HTML-escaping of all spec-sourced text — a spec is user content, and user content
  never becomes markup
- plain `wiki` refreshing a marked `wiki.html` when one exists
- light and dark themes via `prefers-color-scheme`, and a criteria filter in
  vanilla inline JS

Out, named:

- **`--serve`.** A static file the browser opens covers the ask; a daemon comes
  back only if someone needs a URL other people open.
- **publishing anywhere** (artifacts, pages, a forge). The file is local output;
  where it goes afterwards is the user's.
- **markdown rendering inside criteria beyond bold and inline code.** The specs use
  little else in criterion bullets; a full renderer is a dependency this change
  does not need.
- **a second search scope.** The filter covers criteria text and capability names,
  as the demo did; filtering intros and Governs globs comes back if someone misses
  a hit they expected.
- **configurable styling or templates.** One embedded template; a knob nobody
  asked for is a knob to maintain.

Untested by decision, and named as such (the cutter asked): the inline filter's
behaviour and the light/dark theming run in a browser Go tests cannot drive. The
builder opens the generated page and looks, per build-and-check's render-and-look
rule, and the evidence names what was looked at. The bold/backtick conversion is
*not* on that list — it is provable in Go and criterion 1's proof asserts it.

## Constraints

- Stdlib only, again: `html` for escaping, the template as an embedded string. No
  goldmark, no bundler, no JS toolchain — the page's JS is a filter in a few lines.
- Discovery, the empty-directory refusal and the missing-directory refusal are the
  existing code paths, shared, not reimplemented — `--html` changes the renderer
  and the target filename, nothing upstream of them.
- The page must hold with no network: styles inline, JS inline. Remote fonts are
  the one allowed external reference (they degrade to declared fallback stacks
  offline).
- The HTML marker is the first line, as an HTML comment, so the ownership check
  stays "read one line", same as the markdown side.

## Prior decisions

- **A static file, not a server.** The user's example was a page reached by a
  link; `file://` reaches it. The markdown wiki's no-daemon reasoning holds
  unchanged. 2026-08-18.
- **One `wiki` run refreshes every marked view present.** The alternative was
  amending record-work again to name `--html` at landings; teaching the binary to
  refresh what it owns keeps the payload instruction at one sentence and cannot
  drift from the set of outputs. Decided here, 2026-08-18.
- **The template ships inside the binary.** A template file on disk would be
  payload the linker has to place and the ownership scan has to classify; a const
  string is none of that and the template has exactly one consumer.
- **Escaping is a criterion, not a detail.** The specs being rendered are
  arbitrary project content; `<script>` in a criterion must render as text. This
  is the trust-boundary rule write-spec never lets be scoped out.

## Task breakdown

1. The HTML renderer over the existing extraction: template, escaping, marker.
2. Flag parsing, the shared-path wiring, and the both-views refresh.
3. The payload nothing: record-work's instruction already covers regeneration —
   verified by reading, not changed.

## Verification criteria

- When invoked with `--html`, the `wiki` subcommand shall write `wiki.html` into
  the discovered specs directory — a single self-contained page carrying a
  navigation entry and a section per capability, each section with the
  capability's intro, its `Governs:` line and its criteria — and shall not write
  or modify `README.md` in that run. A `**bold**` and a backticked span in a
  fixture criterion shall arrive as `<strong>` and `<code>` markup.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLWritesTheViewer
- The generated `wiki.html` shall open with a first-line HTML comment marker
  naming `libretto wiki` as its generator and the command that refreshes it.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLCarriesTheMarker
- If the target `wiki.html` exists and does not carry the marker, then `wiki
  --html` shall refuse, report the conflict, exit non-zero, and leave the file
  untouched.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLNeverOverwritesAForeignFile
- The `wiki --html` output shall carry every spec-sourced string HTML-escaped, so
  that markup or script inside a spec renders as text and never executes.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLEscapesSpecContent
- The `wiki --html` output shall reference no external resource other than font
  stylesheets from `fonts.googleapis.com` (and the `fonts.gstatic.com` files they
  pull) — no external scripts, no other stylesheets, no remote images — so the
  page works from `file://` with no network, fonts degrading to their declared
  fallback stacks.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLIsSelfContained
- The `--html` run shall produce byte-identical output for unchanged input.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLIsDeterministic
- When run without `--html` and the specs directory holds a `wiki.html` carrying
  the marker, the `wiki` subcommand shall refresh that file in the same run, so
  one landing-step invocation keeps every marked view current; where the
  `wiki.html` present carries no marker, the plain run shall leave it alone and
  not fail — the tool refreshes only what it owns, and a foreign file provokes a
  refusal only when `--html` targets it explicitly. A plain run that errored on a
  foreign `wiki.html` would block every landing regeneration in that project.
  Proof: cmd/libretto/wiki_test.go TestPlainWikiRefreshesAMarkedHTMLView
- The `wiki --html` run shall write nothing but the one `wiki.html`.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLWritesNothingButTheOneFile
- If `wiki` is given an argument other than `--html`, then it shall report the
  unknown argument and exit non-zero, writing nothing.
  Proof: cmd/libretto/wiki_test.go TestWikiRejectsAnUnknownFlag
