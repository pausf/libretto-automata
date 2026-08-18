# Tasks — add-wiki-panel-entry

Execution: build-and-check (phase 6), one fresh session per open box.

- [x] **1. The wiki row, wired end to end, with its proof** — one commit.
  - Extract `findSpecsDir(projectDir) (string, bool)` from `wiki` in
    `cmd/libretto/wiki.go`, so visibility and the command share one discovery.
  - Add `projectDir` as one parameter to `dispatch` and `runCaptured` in
    `cmd/libretto/main.go`, threaded from `run`'s existing single resolution
    through `panelUI`/`panelRefresh`; never a second `os.Getwd`. Update every
    existing call site, including tests (the compiler enumerates them).
  - Add the `wiki` case to `dispatch`: `wiki(os.Stdout, projectDir, nil)`.
  - In `panelData`, after the `status` row: if `scope == target.ProjectScope`
    and `findSpecsDir` hits, append
    `{Label: "wiki", Desc: "render this project's specs into <dir>", Enabled: true}`.
  - Write both tests in `cmd/libretto/wiki_test.go`:
    `TestPanelOffersWikiOnlyInAProjectWithSpecs` (row present under project
    scope with specs; absent under global scope on the same fixture; absent
    in a project with no specs dir) and `TestDispatchRunsWiki` (dispatching
    the row runs plain `wiki` against the project dir and leaves the
    generated `README.md` in the specs directory).
  - Force the global-scope arm of `TestPanelOffersWikiOnlyInAProjectWithSpecs`
    red on purpose (break the scope condition), watch it fail, restore.
  - Traces to: both delta verification criteria
    (`.agents/changes/add-wiki-panel-entry/spec.md`).
  - Closes when: both named tests pass and all six gates are green on the
    commit — including `spec-drift --anchors`, which resolves both `Proof:`
    citations only once the tests exist.
  - Waits on: nothing.
  - Evidence: all six gates green on commit "feat(cli): the wiki row in the
    panel"; the global-scope arm forced red with the condition broken, observed
    failing with its own message, restored.

The capability spec delta lands separately in the final commit; it is not a box.
