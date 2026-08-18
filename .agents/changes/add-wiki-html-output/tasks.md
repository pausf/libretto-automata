# Tasks — add-wiki-html-output

Executing phase: build-and-check
Spec: .agents/changes/add-wiki-html-output/spec.md
Plan: .agents/changes/add-wiki-html-output/plan.md

One note that shaped the cut: the delta's `Proof:` citations are committed on this
branch, and `spec-drift --anchors` requires each cited test to exist in
`cmd/libretto/wiki_test.go`. No commit passes the six gates until all nine tests
exist and pass — so the feature and its proofs are one box, not three layers.

- [ ] 1. `wiki --html` end to end: the seam, the renderer, the refresh rule, all
      nine tests, one commit.
      Refactor `wiki(w, projectDir, args)` to parse `--html` (any other arg is an
      error) and extract a shared `writeMarked(path, marker, render)` ownership
      seam; add `renderWikiHTML` over the existing `[]wikiCapability` — embedded
      const template, first-line HTML-comment marker naming `libretto wiki --html`
      as the refresh command, inline token CSS with a `prefers-color-scheme` dark
      block, sidebar nav with criteria counts, per-capability sections (intro,
      `Governs:`, criteria), inline vanilla-JS filter; every spec-sourced string
      through `html.EscapeString`, bold/backtick conversion applied *after*
      escaping and asserted to produce `<strong>`/`<code>`; plain run additionally
      re-renders a `wiki.html` that carries the HTML marker, silently skips one
      that does not, does nothing when absent; `--html` touches only `wiki.html`.
      Wire `main.go` to pass the remaining args and extend the usage text. Write
      the nine tests under the exact names the delta cites; force
      `TestWikiHTMLEscapesSpecContent` red first (fixture criterion carrying a
      script tag), then make it green — the plan's named force-red test. Open the
      generated page and look (filter, light/dark) — untested-by-decision, the
      evidence names what was looked at.
      Traces to: all nine delta verification criteria; plan approach steps 1–3.
      Closes when: `go test ./cmd/libretto/ -count=1` passes with all nine new
      tests present, and all six gates are green — in particular
      `spec-drift --anchors` resolves every citation for the first time on this
      branch.
      Waits on: nothing.
      Evidence: —

- [ ] 2. Documentation surface and the payload nothing.
      Add the one AGENTS.md line covering `wiki --html` beside the existing wiki
      command documentation, and verify by reading — not by editing — that
      record-work's landing-regeneration instruction already covers both views
      through the refresh-what-is-marked rule. Record the verification outcome
      here.
      Traces to: plan blast radius ("one AGENTS.md line"); spec task breakdown
      item 3.
      Closes when: the AGENTS.md line is committed, all six gates green, and this
      box carries a one-line note confirming record-work needed no change.
      Waits on: box 1.
      Evidence: —

The capability spec delta is not a box. It lands once, in the final commit, when
the change folder is deleted — same commit as the last code, per the landing rule.
