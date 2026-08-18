# Plan — add-wiki-panel-entry

Durable decisions: the ones in Prior decisions of spec.md

## Summary

Thread the project directory into the panel's action path, add a `wiki` case to
`dispatch`, and append a conditional `wiki` row in `panelData` after `status` —
present only when the scope is project and the wiki's own discovery finds a specs
directory. Two tests beside the existing panel tests.

## Technical context

- Branch `feat/add-wiki-html-output` (PR #67), on top of the landed wiki work.
- Blast radius: `cmd/libretto/main.go` (`panelData`, `dispatch`, `runCaptured`,
  their call sites in `panelUI`/`panelRefresh`), `cmd/libretto/wiki.go` (expose
  discovery as a small helper), `cmd/libretto/wiki_test.go` (two tests), existing
  test call sites of `runCaptured`/`dispatch` gain the new parameter. No new files.
- Gates: the six; the two new citations are red until the tests land in the same
  commit — single-box cut, same as the previous two changes.

## The approach

1. Extract `findSpecsDir(projectDir) (string, bool)` from `wiki` so visibility and
   the command share one discovery — not a second list.
2. `dispatch(action, root, projectDir, tg, confirm)` and
   `runCaptured(action, root, projectDir, tg, confirm)`: one added parameter,
   resolved once in `run` as it already is; the `wiki` case calls
   `wiki(os.Stdout, projectDir, nil)` (the redirected stdout inside `runCaptured`
   is what captures it, as with every other action).
3. In `panelData`, after the `status` row: if `scope == target.ProjectScope` and
   discovery hits, append `{Label: "wiki", Desc: "render this project's specs into
   <dir>", Enabled: true}`.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| `dispatch` calling `os.Getwd()` itself | `run` resolves the project once and threads it; a second resolver is the two-answers bug `main.go` already records happening with the project root. |
| Showing the row disabled under global with an explanation | More visible, but the panel's precedent (`models`) is absence, and the user's instruction was "only project". Comes back if absence confuses. |
| A `--html` row too | One flag, two rows, no new capability — the plain run refreshes a marked viewer already. |

## Risks

| Risk | What catches it |
|---|---|
| The parameter threading misses a `runCaptured` call site | The compiler — the signature change makes every stale call a build error. |
| The row appears in global scope through the toggle path | `TestPanelOffersWikiOnlyInAProjectWithSpecs` asserts absence under `GlobalScope` on the same fixture. |

## Validation and rollback

`go test ./cmd/libretto/`; all six gates per commit. The test to force red on
purpose is `TestPanelOffersWikiOnlyInAProjectWithSpecs`'s global-scope arm —
visibility bugs pass silently when only the happy path is asserted; break the
scope condition, watch it red, restore. Rollback: one revert.

## Complexity deliberately kept

None — this is wiring, and the one seam it adds (`findSpecsDir`) removes a
duplicated list rather than adding structure.
