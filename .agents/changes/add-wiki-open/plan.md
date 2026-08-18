# Plan — add-wiki-open

Durable decisions: the ones in Prior decisions of spec.md

## Summary

Add `--open` to the wiki's flag loop: it implies the HTML view, and after a
successful write hands the path to an opener seam — a package `var` wrapping
`exec.Command(openerArgv(runtime.GOOS, path)...)`. Switch the panel row's
dispatch to pass `--open` and update its description. Rewrite the one test the
contract replaces (`TestDispatchRunsWiki` → `TestPanelWikiRowOpensTheViewer`),
add three new ones. An opener error after a successful write surfaces non-zero,
naming the written path — the cutter asked, the spec now answers.

## Technical context

- Branch `feat/add-wiki-html-output` (PR #67), on top of the panel-row landing.
- Blast radius: `cmd/libretto/wiki.go` (flag, seam, `openerArgv`),
  `cmd/libretto/main.go` (the dispatch case gains `[]string{"--open"}`, the row
  description), `cmd/libretto/wiki_test.go` (three new tests, one renamed and
  rewritten). No new files.
- Gates: the six; the four citations go green in the one box, as before.

## The approach

1. `openerArgv(goos, path) []string` — pure: darwin → `open path`, everything
   else → `xdg-open path`.
2. `var openViewer = func(path string) error` — default builds the argv for
   `runtime.GOOS` and runs it detached (`Start`, not `Run` — the panel must not
   block on the browser process). Tests swap the var and record the path.
3. In `wiki`: `--open` sets htmlMode and a flag; after the `--html` write and
   report, call `openViewer(htmlPath)` and report `opened <path>` on success.
4. `dispatch`'s wiki case: `wiki(os.Stdout, projectDir, []string{"--open"})`;
   the row description becomes "render this project's specs and open the viewer".

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| net/http on localhost + open the URL | A server inside a panel action is a lifetime nobody manages (who stops it? when?); `file://` is the same tab. The spec records the bring-back condition. |
| Opening via `exec` inline, no seam | Untestable without launching browsers in CI; the seam is one var. |
| `Run` (wait for the opener) | `open`/`xdg-open` return fast, but a misconfigured handler could block the panel mid-action; `Start` detaches and the report stays honest ("opened" = handed to the OS). |

## Risks

| Risk | What catches it |
|---|---|
| Opener invoked after a refused write | `TestWikiOpenDoesNotOpenOnFailure` — the foreign-file fixture with a recording seam asserting zero calls. Force-red target: move the open call above the write guard, watch it fail, restore. |
| The row silently keeps the plain run | `TestPanelWikiRowOpensTheViewer` drives `dispatch("wiki", …)` and asserts the seam fired. |

## Validation and rollback

`go test ./cmd/libretto/`; six gates per commit. Force-red:
`TestWikiOpenDoesNotOpenOnFailure` by calling the opener before the guard.
The real browser launch is untested by decision (CI cannot open one); the
builder runs `wiki --open` once by hand and reports what appeared. Rollback:
one revert.

## Complexity deliberately kept

The seam (`openViewer` var + `openerArgv` func) is two names where inline exec
would be none — kept because it is the only thing standing between the suite
and a browser window per test run.
