# raise-the-update-notice

Targets: panel

Governs: internal/ui/**

Move the update notice from its band between the menu and the destination strip to a
captioned banner above the logo, so that news reads as news rather than as one more row.

## Outcomes

With `UpdateNotice` set, the panel opens with a banner and the logo comes after it.

The banner is **its own box inside the frame** — a bordered rule carrying the caption
`NEW VERSION` in its top edge, and the notice text on the line inside it:

```
╭────────────────────────────────────────────────────╮
│  ╭─ NEW VERSION ──────────────────────────────╮    │
│  │ v0.6.0 → v1.0.2 available · choose update  │    │
│  ╰────────────────────────────────────────────╯    │
│  ░▒▓█ ══════════════════════════════════ █▓▒░      │
│                                                    │
│      ▄▀▀▄   …the logo…                             │
```

Everything the current row promises survives the move, unchanged:

- it names **both versions and the action** — `v0.6.0 → v1.0.2 available · choose update`
- **silence is the default**: nothing shows when the versions match, when the check has
  not answered, or when it could not
- it is **elided to the content area**, so a long notice cannot tear either border
- the **narrow layout keeps it**, without the box, because `renderNarrow` drops borders
  and not content
- the check still **never blocks a paint or a keypress** — `Init` returns it as a command

And the band between the menu and the destination strip is **gone**, not duplicated. One
fact gets one place on the screen.

## Scope boundaries

**In:** where the notice is drawn, and what it is drawn as. `Theme.Render`,
`Theme.renderNarrow`, and the criteria in `.agents/specs/panel/spec.md` that name the old
position.

**Out, named so it cannot creep in:**

- **the tombstone defect.** `git ls-remote` cannot see a `retract` block, so a checkout is
  currently offered `v1.0.2` — the retraction tombstone — while an installed copy correctly
  gets `v0.7.0` from the module proxy. Real, separate, and not this change.
- **`subcommandNotice`.** The plain uncoloured line after `status` and friends stays plain
  and uncoloured. It is not inside a frame and has no border to inherit.
- **dismissing the banner.** Already a recorded non-goal on the capability spec; a taller
  notice does not reopen it.
- **a new palette pair.** See the assumption below.
- **`preview` wiring the check.** It renders once and statically; that silence is a
  criterion, not an oversight.

## Constraints

- **The 4.5:1 contrast floor holds in both palettes.** `darkTheme` and `lightTheme` were
  built to it after a version shipped at 1.4:1, and a banner is the loudest thing on the
  panel to get it wrong on.
- **The frame must stay flush at every width.** The banner is two nested borders, so it is
  measured and padded like any other row and elided to `cw-4` before it is boxed.
- **Below `MinPanelWidth` there is no box at all** — `renderNarrow` owns that width.
- **The panel holds no notion of what a version is.** The callback returns the finished
  row; the semver comparison stays in `repo-sync`.
- Three extra rows at the top of the frame. The panel already grows by a separator and a
  row today, so the delta is two rows, and only when there is news.

## Prior decisions

- **`UpdateNotice`, not `Notice`.** Action feedback is overwritten by the first `install`;
  news that one keypress deletes is news nobody finishes reading. Unchanged here.
- **Not the footer.** The footer's last resort is dropping the version, so anything beside
  it inherits that fate. Unchanged, and the move is *away* from the footer, not toward it.
- **Gold means attention, not error.** Being one release behind is news, not a fault.
- **The row this replaces was placed deliberately**, and its reason — "not the footer, not
  `Notice`" — argued against two specific alternatives. Neither was "above the logo", so
  this delta reverses a placement rather than relitigating a settled argument.

### Assumed under `/libretto-attacca`, not asked

Three questions the flow would have put to the user. Each carries what changes if the
assumption is wrong.

1. **A bordered box, not a filled tag.** The reference image is a *filled* blue label, and
   the terminal analogue would be a background colour with ink on it. That needs a new
   theme pair verified at 4.5:1 in both palettes, and `Fg` is the only helper the theme
   has. A bordered box reuses `Gold`, which already clears the floor as a foreground in
   both palettes, so it adds no contrast risk.
   **If wrong:** the user wants the filled sticker. That is a `Tag(ink, bg)` helper and two
   new palette values, each measured against its terminal background before it ships.
2. **The caption is `NEW VERSION`, verbatim from the image.** English, like every other
   string in the panel.
   **If wrong:** a one-word string change and its assertion.
3. **The banner goes above the logo, inside the frame** — not above the frame, and not
   below the logo. The user said "sobre el logo", and outside the frame there is nothing to
   align to and nothing measuring the width.
   **If wrong:** the logo stays first and the banner moves under it, which is the same rows
   in a different order.

## Task breakdown

1. Draw the banner: a `banner` helper on `Theme` producing the boxed, captioned rows for a
   given content width, and `Render` emitting it above the logo instead of below the menu.
2. Keep the narrow layout telling the truth: `renderNarrow` shows the notice above the
   mark, still without a border.
3. Land the delta on `.agents/specs/panel/spec.md`, replacing the position criterion and
   its `Proof:` rather than leaving two placements described.

## Verification criteria

- **the banner renders above the logo**, and no longer between the menu and the strip
  Proof: internal/ui/notice_test.go TestUpdateNoticeRendersAboveTheLogo
- **it is boxed and captioned** — its own border inside the frame, carrying `NEW VERSION`
  Proof: internal/ui/notice_test.go TestUpdateNoticeBannerIsBoxedAndCaptioned
- an empty notice changes nothing about the panel, banner included
  Proof: internal/ui/notice_test.go TestPanelOmitsUpdateNoticeWhenEmpty
- **the narrow layout keeps it**, above the mark and without the box
  Proof: internal/ui/notice_test.go TestNarrowLayoutKeepsUpdateNotice
- **the frame does not tear with a notice set**, at every width — both borders flush
  Proof: internal/ui/fluid_test.go TestFrameHoldsWithUpdateNotice
- news and action feedback both render, because they are two fields
  Proof: internal/ui/notice_test.go TestUpdateNoticeAndActionFeedbackCoexist
- running an action does not clear the news, and neither does moving the cursor
  Proof: internal/ui/notice_test.go TestActionFeedbackDoesNotOverwriteUpdateNotice
  Proof: internal/ui/notice_test.go TestNavigationDoesNotClearUpdateNotice
- `Init` returns the check as a command, so the first paint never waits on it
  Proof: internal/ui/notice_test.go TestInitReturnsReleaseCheckCommand
- **no check configured means no command and no banner**, which is `preview`
  Proof: internal/ui/notice_test.go TestNoReleaseCheckMeansNoCommandAndNoNotice
