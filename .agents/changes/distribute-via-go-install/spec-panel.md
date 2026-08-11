# Panel — delta

Targets: panel

One row: a newer release exists, and `update` is how you get it.

## Outcomes

- **The notice is its own row inside the frame**, between the menu and the destination
  strip. Not the footer: the footer already drops `p.Version` when the legend and the
  version cannot both fit, and a notice that disappears at 96 columns is a notice that
  was never read. Not appended to a target row either — the strip is about where links
  go.
- **It names both versions and the action.** `v0.2.0 → v0.3.0 available · choose update`.
  A row saying "an update is available" with no version is a row that cannot be checked
  against `libretto version`, and one with no action is a notification with nowhere to go.
- **It uses the attention colour, not the error colour.** Being one release behind is
  news, not a fault.
- **Nothing shows when the versions match, when the check has not answered yet, or when
  the running version is unidentifiable.** Absence is the default, and it is silent.
- **It is not `Panel.Notice`.** That field is action feedback, and the first `install`
  overwrites it — the same overwrite that once ate the selector's key legend, recorded in
  `footer`'s own comment. A separate field means one press of ⏎ does not delete the news.
- **It survives the narrow layout.** `renderNarrow` drops the border, not the content. A
  degraded layout that also degrades what it is telling you is two losses for one
  terminal width.
- **The check never blocks a paint or a keypress.** It arrives as a message; the panel
  renders complete without it and re-renders when it lands. A panel that waits on
  `git ls-remote` before its first frame hangs on bad DNS, and the user's only recourse
  is ⌃C on a tool that looked broken.
- **This package still performs no I/O.** The check is a callback supplied by the caller,
  in the shape `WithRefresh` and `WithRunner` already established, and `Init()` returns
  the command that runs it — the first thing `Init` has ever had to do.

## Scope boundaries

**In:** the notice row, its placement in both layouts, the field it lives in, the
callback seam and the message that delivers the answer.

**Out:**

- **performing the check.** `internal/repo` answers, `cmd/libretto` wires. Unchanged rule.
- **a dedicated key.** `update` is already in the menu with a cursor on it; a `u`
  accelerator is a second way to reach a row that is one keystroke away.
- **dismissing it.** State to persist, a file to write, and the notice is one row that
  goes away when you update.
- **a changelog, release notes, or what changed.** That is the remote's job and a browser's.
- **animation, spinner, or "checking…".** A row that appears when there is news beats a
  row that narrates a network call the user did not ask for.

## Constraints

- No filesystem, no network, no `os/exec` in `internal/ui`. The callback is the seam.
- The fluid frame's width arithmetic must account for the row like every other, so the
  border cannot tear — `fluid_test.go` is the existing guard and the new row goes through
  it.
- `Update` stays free of I/O and testable by direct call.

## Prior decisions

- **Cached, not live.** The panel reads whatever the check cache holds; `doctor` is the
  command that pays for a live call. Recorded in the `cli` delta.
- **The action is `update`.** Bootstrapper shape means `update` already lands the newer
  tag: pull, rebuild, relink. No new menu entry.

## Task breakdown

1. `Panel.UpdateNotice string` — the notice, distinct from `Notice`. Not `Panel.Update`:
   a field named `Update` sitting beside a Bubbletea `Update` method is a name that has
   to be disambiguated by type every time it is read.
2. Render it as a row between the menu and the strip, in the attention colour, and in
   `renderNarrow`.
3. `WithReleaseCheck(func() (string, error))` and an `Init()` that returns it as a
   `tea.Cmd`.
4. A message type carrying the finished notice text, handled in `Update`.

**The comparison happens on the caller's side, not here.** The callback returns the row to
show or the empty string, so semver ordering lives in exactly one place — `internal/repo`
— rather than in a package whose whole promise is that it holds no logic about the world.
A second implementation of "is this newer" is a second implementation that can disagree.

## Verification criteria

```
Proof: internal/ui/notice_test.go TestPanelRendersUpdateNoticeBetweenMenuAndStrip
Proof: internal/ui/notice_test.go TestPanelOmitsUpdateNoticeWhenEmpty
Proof: internal/ui/notice_test.go TestNarrowLayoutKeepsUpdateNotice
Proof: internal/ui/notice_test.go TestUpdateNoticeAndActionFeedbackCoexist
Proof: internal/ui/fluid_test.go TestFrameHoldsWithUpdateNotice
Proof: internal/ui/notice_test.go TestInitReturnsReleaseCheckCommand
Proof: internal/ui/notice_test.go TestUpdateNoticeSetFromMessage
Proof: internal/ui/notice_test.go TestActionFeedbackDoesNotOverwriteUpdateNotice
Proof: internal/ui/notice_test.go TestNoReleaseCheckMeansNoCommandAndNoNotice
```
