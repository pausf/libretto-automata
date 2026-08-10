# Panel — delta

Targets: panel
Governs: internal/ui/**

The selector shows the active destination's agents, and says which rows are shared.

## Outcomes

The rows are whatever the active target holds — the same set `models` lists under the
same flag — and **`tab` changes them**. Switching destination while the selector is
open reloads it, because a screen that keeps showing one target's agents under
another's name is the lie the destination strip exists to prevent.

A row that is a symlink into this repository is marked `shared`:

```
  ❯ [x] review-design        haiku      shared
    [x] sdd-apply            sonnet
    [ ] review-risk          opus
    [ ] jd-judge-a           (session)

    space mark · a all · m model · esc back
```

- **`shared` is a warning, not decoration.** Applying to a shared row changes every
  project on the machine; applying to an unmarked one changes this target only.
- Legible with colour stripped, like every other signal in this panel: the word, not
  a colour.
- The menu row's tally counts the target's agents, so it answers "how much of *this
  destination* is still expensive" rather than a number about a directory the user is
  not looking at.
- An empty or absent agents directory shows a plain line saying so, not an empty box.

## Scope boundaries

**In:** the rows following the active destination, the shared marker, `tab` reloading.

**Out:**

- **Confirming a write to a foreign file.** Settled in `spec-cli.md`; the panel does
  not reintroduce a prompt the CLI decided against.
- **Showing both destinations at once.** One list, one destination, `tab` between —
  the panel's existing grammar.
- **`internal/ui` learning what a target or a symlink is.** Rows arrive carrying a
  `Shared` flag; the caller decides what set them.

## Constraints

- The package still touches no filesystem (`internal/ui/model.go:11`) and imports
  neither `internal/target` nor `internal/agentmodel`. `AgentRow` gains one bool.
- The reload on `tab` goes through the callbacks already there — no new plumbing, and
  a failed reload leaves the previous rows and says so, the way `nextScope` already
  handles a failed refresh.
- Frame flush at every width, no `𝄞`, ASCII-safe honoured — unchanged.

## Prior decisions

- **The marker is a word, not a colour.** The panel already learned this the
  expensive way: the destination strip encoded selection in a bullet, shipped, and
  had correct behaviour reported as a bug. Colour is the signal that vanishes first.

## Task breakdown

3. `internal/ui`: `AgentRow.Shared`, render the marker, reload the rows on `tab`.
4. `cmd/libretto`: build rows from the active target and set `Shared` from
   `link.Owned`.

## Verification criteria

- a shared row is marked and a local one is not
  Proof: internal/ui/models_test.go TestSharedAgentsAreMarked
- **the marker survives colour being stripped**
  Proof: internal/ui/models_test.go TestSharedMarkerIsLegibleWithoutColour
- `tab` in the selector reloads the rows for the new destination
  Proof: internal/ui/models_test.go TestTabReloadsTheSelectorForTheNewDestination
- a failed reload keeps the rows it had and says so
  Proof: internal/ui/models_test.go TestAFailedReloadKeepsTheRowsAndSaysSo
- an empty agents directory renders a plain statement, not an empty box
  Proof: internal/ui/models_test.go TestAnEmptyAgentSetSaysSo
- the tally counts the active destination's agents
  Proof: cmd/libretto/models_test.go TestMenuTallyCountsTheActiveTargetsAgents
- marking, `a`, the catalogue and apply all behave as before
  Proof: internal/ui/models_test.go TestChosenModelReachesOnlyTheMarkedRows
