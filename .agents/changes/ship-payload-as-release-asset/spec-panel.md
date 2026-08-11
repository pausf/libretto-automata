# Panel — delta

Targets: panel

The notice row is unchanged. One row's description is not.

## Outcomes

- **`update` stays one row.** An earlier draft added a second for `upgrade` and disabled whichever
  did not apply; collapsing the two commands removed the row along with the command.
- **The row names the mechanism this machine will use** — `pull this checkout · rebuild · relink`
  in a checkout, `install the newest version · relink` otherwise. One row, honest about what
  pressing it does, without becoming two.
- **And it no longer says `git pull`.** That string in front of somebody who only wanted to use
  the tool is the complaint this whole change came from.
- **The notice names `update`**, which is the only command there is. It used to have to choose
  between two names per machine.

Everything else about the notice row stands: its position between the menu and the strip, its own
field, the attention colour, silence by default, surviving the narrow layout, and arriving as a
command so the first paint never waits.

## Scope boundaries

**In:** the `update` row's description.

**Out:**

- **deciding which mode the panel is in.** `cli` decides and hands over strings; this package
  renders them and continues to know nothing about a filesystem.
- **a second row, or hiding one.** Both were drafted and both are gone with the second command.
- **a confirmation before updating.** `y/n` is for the destructive actions.
- **progress while it installs.** The action reports inside the frame when it finishes.
  *Ceiling:* a payload large enough that the panel looks hung.

## Prior decisions

- **One row, because there is one command.** Recorded because the draft that had two argued for
  always showing both and disabling the inapplicable one — which was right, for a design that no
  longer exists.

## Task breakdown

1. The row's description, per mode. Nothing in the renderer changes.

## Verification criteria

```
Proof: cmd/libretto/update_release_test.go TestTheUpdateRowNamesTheMechanism
Proof: internal/ui/notice_test.go TestPanelRendersUpdateNoticeBetweenMenuAndStrip
Proof: internal/ui/notice_test.go TestPanelOmitsUpdateNoticeWhenEmpty
Proof: internal/ui/notice_test.go TestNarrowLayoutKeepsUpdateNotice
Proof: internal/ui/notice_test.go TestUpdateNoticeAndActionFeedbackCoexist
Proof: internal/ui/notice_test.go TestInitReturnsReleaseCheckCommand
Proof: internal/ui/notice_test.go TestNoReleaseCheckMeansNoCommandAndNoNotice
Proof: internal/ui/notice_test.go TestUpdateNoticeSetFromMessage
Proof: internal/ui/notice_test.go TestActionFeedbackDoesNotOverwriteUpdateNotice
Proof: internal/ui/notice_test.go TestNavigationDoesNotClearUpdateNotice
Proof: internal/ui/fluid_test.go TestFrameHoldsWithUpdateNotice
```

One of the two menu tests the draft added is gone with the row it covered. The other outlived it:
`TestTheInapplicableActionIsDisabledNotHidden` was exercising a standing rule of this spec — the
panel does not hide what it cannot do — and that rule did not depend on there being two rows.

```
Proof: internal/ui/menu_test.go TestTheInapplicableActionIsDisabledNotHidden
```

`internal/ui` is otherwise untouched by this change, which is the sign the seam between it and
`cli` held: the panel never learned what a module cache is.
