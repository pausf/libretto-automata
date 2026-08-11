# Panel — delta

Targets: panel

The notice row is unchanged. One menu entry and one string are not.

## Outcomes

- **`upgrade` joins the menu**, and `update` stays. Two entries, because they are two
  commands — see the `cli` delta. Hiding one based on mode was the alternative and it loses:
  a menu that changes shape between machines is a menu whose screenshots and instructions
  are wrong somewhere.
- **The one that does not apply is shown disabled**, not hidden. That is already this
  spec's rule — *actions that cannot run are shown disabled rather than hidden, because the
  panel does not promise what it cannot do and does not hide what is coming* — and it
  applies here unchanged.
- **The notice names the command that will actually work.** `v0.3.0 → v0.4.0 available ·
  choose upgrade` on an installed copy; `choose update` in a checkout. A row naming a
  command that refuses is worse than no row, and the string is still handed in whole from
  `cli` — this package continues to know nothing about which mode it is in.

Everything else about the row stands: its position between the menu and the strip, its own
field, the attention colour, silence by default, surviving the narrow layout, and arriving
as a command so the first paint never waits.

## Scope boundaries

**In:** the second menu entry, and the notice text arriving already correct.

**Out:**

- **deciding which mode the panel is in.** `cli` decides and hands over a string; this
  package renders it. Unchanged, and the reason is unchanged: a package that renders and
  also detects is a package that cannot be tested without a filesystem.
- **a confirmation before upgrading.** `y/n` is for the destructive actions. An upgrade
  keeps the previous version and rolls back with a symlink swap.
- **progress while it downloads.** The action already reports inside the frame when it
  finishes, and a progress bar is machinery for a tarball of markdown.
  *Ceiling:* a payload large enough that the panel looks hung.

## Prior decisions

- **Two entries, both always present.** Recorded because "just hide the one that does not
  apply" is the obvious instinct and it makes the menu machine-dependent.

## Task breakdown

1. The `upgrade` entry, enabled by the same mode check `cli` already computes.
2. Nothing in the renderer changes; the notice string arrives correct.

## Verification criteria

```
Proof: internal/ui/menu_test.go TestBothUpgradeAndUpdateAreOffered
Proof: internal/ui/menu_test.go TestTheInapplicableActionIsDisabledNotHidden
Proof: cmd/libretto/version_test.go TestReleaseNoticeNamesTheCommandForTheMode
```

The third is cited by the `cli` delta too, deliberately: it is one behaviour with two
owners, and a criterion each spec states in its own terms is how both stay honest about
depending on it.
