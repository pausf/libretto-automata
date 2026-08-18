# Delta — the wiki row in the panel

Targets: cli

The panel gains a `wiki` row, and only where pressing it can succeed: project
scope, in a project that has a consolidated specs directory. Global scope never
shows it — `~/.claude` has no specification, and the user drew that line
explicitly.

## Outcomes

Opening the panel in a project with specs, with the scope on `project`, shows a
`wiki` row after `status`; pressing it runs the plain `wiki` command against the
working directory and reports its output in place, like every other row. Flipping
the scope to `global` makes the row disappear. A project with no specs directory
never shows it in either scope.

## Scope boundaries

In: the conditional row in `panelData`, the `wiki` case in `dispatch` (which needs
the project directory threaded through `runCaptured` and `dispatch`), and the
tests beside their siblings.

Out, named:

- **a `--html` variant row.** The plain run already refreshes a marked
  `wiki.html`; a second row for one flag doubles the menu for no new capability.
  Comes back if creating the viewer from the panel is actually missed.
- **any change to the CLI surface or the wiki command itself.** This is wiring.

## Constraints

- The panel's labels are the subcommand names — one list of actions, not two.
  The row is `wiki` and dispatches to the same function the CLI runs.
- `dispatch` currently receives no project directory; it gains one parameter
  rather than calling `os.Getwd` itself — `run` resolves the project once and
  threads it, and a second resolver is the two-answers bug the file already
  records.
- Visibility reuses the wiki's own discovery order — not a second list.

## Prior decisions

- **Project scope only, by instruction.** The user drew the line 2026-08-18: the
  row appears under project scope and never under global.
- **No specs, no row** — the `models` precedent: an entry that opens an empty
  screen is a promise the panel cannot keep. Assumed rather than asked; if a
  user wants the row to appear and explain what is missing, that is the change.
- **The row runs the plain command**, not `--html`, because the panel's contract
  is label = subcommand, and the plain run refreshes every marked view anyway.

## Task breakdown

1. Thread the project directory to `dispatch`, add the `wiki` case and the
   conditional row, with the two tests.

## Verification criteria

- Where the panel scope is project and the project holds a consolidated specs
  directory, `panelData` shall include an enabled `wiki` row; where the scope is
  global, or the project holds no specs directory, it shall not include one.
  Proof: cmd/libretto/wiki_test.go TestPanelOffersWikiOnlyInAProjectWithSpecs
- When the `wiki` row is dispatched, the panel shall run the plain `wiki`
  command against the project directory, leaving the generated `README.md` in
  the project's specs directory.
  Proof: cmd/libretto/wiki_test.go TestDispatchRunsWiki
