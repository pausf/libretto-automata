# Panel — delta

Targets: panel
Governs: internal/ui/**

A second screen. Everything else in the panel spec stands, including the rule that
this package never reads the filesystem.

## Outcomes

The menu gains one entry. Choosing it replaces the menu with a **model selector** and
`esc` brings the menu back.

The selector is a list of every agent, one row each: its name, the model it runs on
today, and a mark.

```
  ❯ [x] review-design        haiku
    [x] review-tests         haiku
    [ ] review-security      opus
    [ ] spec-writer          (session)

    space mark · a all · m model · esc back
```

- **Marking is multi.** `space` marks the row under the cursor, `a` marks every row —
  and `a` again clears them, because a key that only ever adds leaves no way back but
  pressing space eleven times.
- **One model, applied to the marked set, in one act.** `m` opens the catalogue; the
  chosen model goes to every marked row at once. That is the ordinary case, not the
  advanced one: making the prose lenses cheap is one gesture, not four.
- **Nothing marked means nothing happens**, and the panel says so. It does not fall
  back to "the row under the cursor" — a selector with a marking mechanism that
  sometimes ignores it teaches the user not to trust the marks.
- **The mark is legible without colour.** `[x]` and `[ ]`, *and* the theme's emphasis.
  Both signals, the same rule the destination strip already follows — colour alone
  fails a non-colour terminal.
- After applying, the rows show their new models. A screen that requires a reopen to
  tell the truth is a screen that lies for as long as it is open.
- The selector obeys the width rules the panel already has: flush frame at every
  width, no tearing when narrow.

## Scope boundaries

**Out:** filtering or searching the list — four to eight agents fit on a screen, and a
filter is machinery for a scale this does not have; per-row model choice in the same
gesture as marking; any confirmation prompt, because writing a frontmatter key is
reversible in one keystroke and `y/n` is for the destructive actions.

## Constraints

- **`internal/ui` stays free of the filesystem** (`internal/ui/model.go:11`). The
  selector gets its rows and applies its choice through callbacks supplied by the
  caller, exactly as `Refresh` and `Runner` already do. No new dependency on
  `internal/target`, none on `internal/agentmodel`.
- Screen state lives in the existing `Model`; `Update` stays free of I/O so transitions
  are testable by calling it directly.
- The confirmation machinery (`pending`, `pendingScope`) is for destructive actions and
  is not reused here.
- No `𝄞`, no `♩♪♫♬` — `AGENTS.md` bans them from anything a terminal renders.
- ASCII-safe mode is already a panel concept and the selector honours it.

## Task breakdown

6. `internal/ui`: the selector screen — rows, marks, `a`, the model chooser, `esc`,
   and the two callbacks that keep the package off the filesystem.
7. `cmd/libretto`: wire the callbacks to `internal/agentmodel`, and add the menu entry.

## Verification criteria

- the menu entry opens the selector and `esc` returns to the menu
  Proof: internal/ui/models_test.go TestSelectorOpensFromTheMenuAndEscapeReturns
- `space` marks and unmarks the row under the cursor
  Proof: internal/ui/models_test.go TestSpaceMarksAndUnmarksTheCurrentRow
- `a` marks every row, and again clears every row
  Proof: internal/ui/models_test.go TestMarkAllTogglesEveryRow
- a chosen model reaches every marked row and no unmarked one
  Proof: internal/ui/models_test.go TestChosenModelReachesOnlyTheMarkedRows
- applying with nothing marked changes nothing and says so
  Proof: internal/ui/models_test.go TestApplyingWithNothingMarkedSaysSo
- the rows show the new model without reopening the screen
  Proof: internal/ui/models_test.go TestRowsShowTheNewModelAfterApplying
- an agent with no declared model renders as running the session's
  Proof: internal/ui/models_test.go TestUndeclaredAgentRendersAsSession
- the mark is visible with colour stripped
  Proof: internal/ui/models_test.go TestMarkIsLegibleWithoutColour
- the selector frame is flush at every width
  Proof: internal/ui/models_test.go TestSelectorFrameIsFlushAtEveryWidth
- a failing apply reports the error and leaves the screen usable
  Proof: internal/ui/models_test.go TestFailedApplyIsReportedAndTheScreenSurvives
