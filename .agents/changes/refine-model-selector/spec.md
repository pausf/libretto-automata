# refine-model-selector — delta

Targets: panel
Governs: internal/ui/**

The model selector reads as a flat alphabetical list, its persistent legend belongs
to the other screen, and the key that marks every row is discoverable only from a
notice that the next action overwrites. Three changes, one screen.

## Outcomes

### The rows are grouped by model, with a rule between groups

Rows arrive in whatever order the caller supplies. The selector **sorts them by the
agent's current model** and draws one separator line between groups, so a glance
answers "what is still on the expensive model?" without reading every row.

```
  ❯ [ ] all
    ──────────────────────────────────────────
    [ ] review-lens-design       haiku    shared
    [ ] review-lens-intent       haiku    shared
    [ ] review-lens-tests        haiku    shared
    ──────────────────────────────────────────
    [ ] review-lens-reliability  (session)  shared
    [ ] review-lens-security     (session)  shared
    [ ] spec-writer              (session)  shared
    [ ] work-reviewer            (session)  shared
```

- **Group order is the catalogue's order, session last.** The same rule `Tally`
  already follows and for the same reason: the catalogue is cheapest-first, and the
  session default is not a price — a screen opened to answer "how much of this is
  still expensive?" does not lead with the one group that cannot answer it.
- **Within a group, rows sort by name.** Two agents on one model in an order nobody
  chose is a list that reorders itself between sessions.
- **A model absent from the catalogue still gets a group, last** — after session.
  That is where `Tally` already puts it, and the two must not disagree: the tally row
  and the selector below it are one screen. An agent declaring a model this build does
  not know about is a state, not an error, and it must not vanish or land in the wrong
  pile.
- **The separator is a rule, not a blank line.** A blank line inside a bordered frame
  reads as the end of the list. It is an indented dim rule, **not** the frame's
  `├───┤` junction: that glyph means "new section of the panel", and three of them
  inside one list would read as three panels.
- **One group means no separator.** A rule with nothing on one side of it is a
  division of nothing.

### `all` is a row, and space on it marks everything

A row labelled `all` sits **above the first group**, separated from it by the same
rule. `space` on it marks every agent; `space` again clears them.

- **Its box reflects the truth**: `[x]` only when every agent row is marked. A
  master checkbox that stays `[x]` after one row is unmarked is a checkbox that lies.
- **`a` keeps working**, unchanged. The row is the discoverable path, the key is the
  fast one; removing the key to justify the row would take something away from
  whoever already found it.
- **`all` is not an agent.** It never reaches `ApplyModel` as a name, and it does not
  appear in `MarkedAgents`. The applied set is the marked agent rows, exactly as
  today.
- **`m`/`enter` on the `all` row opens the catalogue** for whatever is marked, like
  any other row. The row is a marking control, not a mode.
- **No agents means no `all` row.** The empty-destination line stands alone.

### The footer legend tells the truth about the screen it is on

The footer currently reads `↑↓ · ⏎ select · tab scope · q quit` on both screens. In
the selector, `⏎` opens the model catalogue, `space` marks, and `a` marks everything
— none of which it says.

- **In the selector the footer reads** `↑↓ · space mark · a all · m model · tab scope
  · esc back`.
- **It is the footer, not the notice.** The opening notice is transient: the first
  apply, the first error, the first anything replaces it, and the legend is gone for
  the rest of the session. A legend that disappears the moment the screen is used is
  not a legend.
- **While the catalogue is open it reads** `↑↓ · ⏎ apply · esc back`, because those
  are the only keys that do anything. The same rule the confirm prompt already
  follows.
- **The menu screen's footer is unchanged.**

## Scope boundaries

**In:** row order, the group rule, the `all` row, and the footer's per-screen hints.
All of it inside `internal/ui/`.

**Out:**

- **Sorting by anything else** — name, shared, destination. One order, and it is the
  one that was asked for.
- **A user-visible toggle between grouped and flat.** A setting for a list of seven
  rows is machinery for a scale this does not have.
- **A group header naming the model.** Every row in the group already carries its
  model in the model column; a header would print it twice and cost a row per group.
- **Collapsing a group.** Seven rows fit.
- **Changing what `ApplyModel` receives.** The selection semantics are untouched;
  only how the rows are ordered and how the set gets marked.
- **`filtering`, searching, mouse** — already out in the capability spec and still out.

## Constraints

- **`internal/ui` stays ignorant of the filesystem.** Rows still arrive through
  `ListAgents` and the catalogue through `WithAgents`. The sort happens on rows the
  package was handed, using the catalogue it was handed, and imports nothing new.
- **The rule obeys the fluid frame.** It is drawn at the content width the frame is
  using, like `separator` in `Render`, and does not tear at 58 or at 98 columns.
- **Legible without colour.** The `all` row's `[x]`/`[ ]` and the group rule are
  glyphs, not shades — the same both-signals rule the mark and the `shared` warning
  already follow.
- **No astral-plane runes, no ambiguous-width glyphs** in the rule.
- **The footer must still fit** at `MinContentWidth`; the selector's hint string is
  longer than the menu's and the gap calculation already clamps at 1.

## Prior decisions

- **The user chose sorted, not flat-with-a-rule.** Asked in this session: grouping
  either reorders the rows or leaves them alphabetical and rules only where the model
  changes. The answer was *ordenar por modelo* — the rows move.
- **Catalogue order, session last, unknown after that** is `Tally`'s existing rule
  (`internal/ui/models.go`), settled when the menu tally was written. The selector does
  not restate it — it shares the ranking function, so the tally row and the list below
  it cannot drift. This spec first said unknown models sort *before* session; reading
  `Tally` showed that was a second order invented for one screen, and it was wrong.
- **`a` marks and unmarks with one key**, decided in the capability spec: a key that
  only ever adds leaves no way back but pressing space once per row. The `all` row
  inherits it.
- **The mark is `[x]`/`[ ]` *and* the theme's emphasis, never emphasis alone** — the
  capability spec's both-signals rule, and the reason the `all` row gets a box.
- **The cursor indexes the rendered list, not the agent slice.** The `all` row means
  screen position and agent index are no longer the same number; whichever way the
  implementation resolves that, `MarkedAgents` must keep returning agent names only.
  This is the one place where a wrong-by-one is silent, so it carries its own proof.

## Task breakdown

1. **Sort the rows by model, name within model**, catalogue order with session last
   and unknown models between the two. Selector-side, on the rows as handed in.
2. **Draw the group rule** between model groups at the frame's content width, and
   never above the first group or below the last.
3. **Add the `all` row**: rendered above the first group with its own rule below,
   cursor-reachable, `space` marks and clears every agent, box `[x]` only when all
   are marked, excluded from `MarkedAgents` and from `ApplyModel`.
4. **Make the footer per-screen**: menu, selector, and catalogue-open hints.
5. **Land the delta on `.agents/specs/panel/spec.md`** in the same commit as the code
   that taught it, and delete the change folder.

## Verification criteria

1. Rows come back grouped: every row of one model is contiguous, groups run in
   catalogue order with the session default last, names sort inside a group.
   Proof: internal/ui/models_test.go TestRowsAreGroupedByModel

2. A model the catalogue does not know still renders, in its own group, last — the
   same position `Tally` gives it, from the same shared ranking.
   Proof: internal/ui/models_test.go TestAnUnknownModelGetsItsOwnGroup

3. A rule is drawn between groups and nowhere else — none above the first group, none
   after the last, and none at all when every agent shares one model.
   Proof: internal/ui/models_test.go TestGroupRuleSitsOnlyBetweenGroups

4. `space` on the `all` row marks every agent; `space` again clears them.
   Proof: internal/ui/models_test.go TestSpaceOnTheAllRowMarksEveryAgent

5. The `all` row's box is `[x]` only when every agent is marked, and drops back to
   `[ ]` as soon as one is unmarked.
   Proof: internal/ui/models_test.go TestTheAllRowBoxFollowsEveryAgent

6. `all` never reaches the apply callback and never appears in `MarkedAgents`; with
   the cursor on it and rows marked, the applied set is the agent names alone.
   Proof: internal/ui/models_test.go TestTheAllRowIsNeverAppliedAsAnAgent

7. With the `all` row present, `space` on an agent row marks that agent and no other
   — the cursor offset does not shift the mark by one.
   Proof: internal/ui/models_test.go TestTheCursorMarksTheRowItPointsAt

8. The footer in the selector lists `space`, `a` and `m`; the menu's footer is
   unchanged; the catalogue-open footer lists `⏎ apply · esc back`.
   Proof: internal/ui/panel_test.go TestTheFooterFollowsTheScreen

9. The frame stays flush at every width with the `all` row and the group rules
   present.
   Proof: internal/ui/models_test.go TestSelectorFrameIsFlushAtEveryWidth

10. The `all` row and the group rule are legible with colour stripped.
    Proof: internal/ui/models_test.go TestTheAllRowIsLegibleWithoutColour
