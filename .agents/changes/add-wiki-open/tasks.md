# Tasks — add-wiki-open

Execution: build-and-check, one fresh session per open box. Boxes are cut so each
merges alone — code and proof in the same commit, tree green, all six gates.

The capability spec delta is not a box: it lands on `.agents/specs/cli/spec.md` in
the final commit, when the change folder is deleted.

- [ ] **1. `--open` end to end: flag, seam, argv function, panel row — with all four proofs**
  - Traces: spec.md, all four verification criteria; plan.md "The approach" 1–4.
  - What: in `cmd/libretto/wiki.go`, accept `--open` (implies the HTML view; `--open`
    and `--html --open` identical), add `openerArgv(goos, path) []string` (darwin →
    `open path`, everything else → `xdg-open path`) and
    `var openViewer = func(path string) error` defaulting to a detached
    `exec.Command(...).Start()`; after a successful HTML write, call
    `openViewer(htmlPath)` and report `opened <path>` on success — never before or
    after a refused/failed generation; an opener error surfaces non-zero naming the
    written path. In `cmd/libretto/main.go`, the panel wiki dispatch case becomes
    `wiki(os.Stdout, projectDir, []string{"--open"})` and the row description says
    the viewer opens. In `cmd/libretto/wiki_test.go`: add
    `TestWikiOpenGeneratesAndOpensTheViewer` (including the erroring-seam subcase),
    `TestWikiOpenDoesNotOpenOnFailure` (foreign-file fixture, recording seam, zero
    calls — force it red once by moving the open call above the write guard, watch
    it fail, restore), `TestOpenerArgvPerPlatform`; rename and rewrite
    `TestDispatchRunsWiki` as `TestPanelWikiRowOpensTheViewer` (drives
    `dispatch("wiki", …)`, asserts the seam fired). Tests swap the seam var — no
    test ever launches a browser. Update the usage text alongside the flag. No new
    files.
  - Closes when: all four Proof citations pass in `go test ./cmd/libretto/`; the six
    gates green (`spec-drift --anchors` resolves the four citations against
    `cmd/libretto/wiki_test.go`); and the one hand check named by the plan is done —
    the builder runs `libretto wiki --open` once for real and reports what appeared.
  - Waits on: nothing.
  - Evidence: —
