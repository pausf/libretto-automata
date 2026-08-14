# Delta: the destination strip shows four rows

Targets: panel

## Outcomes

- The strip lists all four destinations — global, project, codex, opencode —
  each row reporting its own state from its own scan, exactly as the two
  existing rows do. `tab` cycles through all of them.
- An unconfigured codex or opencode renders as `○ … not configured` — the
  row DESIGN.md has shown as the intended future since the interface was
  drawn. No error, no missing row.
- A skills-only destination's status summary counts skills alone; it never
  mentions agents or commands it does not accept, and the `models` menu row
  is absent when the active destination takes no agents.
- The active-row contract holds across four rows: `❯` and gold end to end,
  both signals, one active row always.

## Scope boundaries

**In:** the rows `cmd/libretto` feeds the panel, golden files for the
four-row strip, footer unchanged in shape.

**Out, named:**
- Any second ordering of the rows — by name, by configured-first, by tool.
  The panel spec forbids it and this change adds no exception.
- Per-tool icons or branding in the strip. Rows are text; the name is the
  identity.
- UI-side knowledge of targets. `internal/ui` still never imports
  `internal/target`; four rows arrive through the same `Refresh`/`Runner`
  seam as two did.

## Constraints

- `nextScope()` is already index-based and dimension-agnostic; the change
  feeds it more rows and touches none of its logic.
- The strip renders arbitrary-length row lists today (`demoPanel` proves it
  with a codex row); this change makes that fixture real, not decorative.
- Panel actions keep receiving the destination index, never a captured
  target.

## Prior decisions

- Footer keeps saying `tab scope`. Renaming the word to "destination" is a
  wording change with layout consequences (footer width is a measured
  constraint) and is not what this change is for. Ceiling: if user confusion
  shows up, rename in its own change with new goldens.

## Task breakdown

1. Extend the row loop in `panelData` to the four-destination order.
2. Golden files for the four-row strip, configured and not.
3. Panel-action isolation test against a new destination.

## Verification criteria

- the strip renders four rows, each with its own state, active marked on
  exactly one
  Proof: cmd/libretto/panelrun_test.go TestStripShowsAllFourDestinations
- an unconfigured new destination renders as not configured and stays
  selectable
  Proof: cmd/libretto/panelrun_test.go TestUnconfiguredDestinationRow
- panel actions act on the active destination only, including the new ones
  Proof: cmd/libretto/scope_test.go TestPanelPruneActsOnTheActiveDestinationOnly
- the four-row strip matches its goldens in colour and mono
  Proof: internal/ui/panel_test.go TestFourDestinationStripGolden
- the `models` menu row is absent when the active destination accepts no
  agents
  Proof: cmd/libretto/panelrun_test.go TestModelsRowAbsentForSkillsOnlyDestination
