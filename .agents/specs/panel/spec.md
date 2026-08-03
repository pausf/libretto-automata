# Panel

Governs: internal/ui/**

The interactive surface: logo, palette, layout, menu. The part of this project where a
requirement written before seeing a terminal is a requirement that turns out to be
wrong.

## Outcomes

With a TTY and no subcommand, one panel: the wordmark, a menu of actions with live
state, and a strip listing every target.

It is **readable and intact at every terminal width**. Not "usually" — the frame is
flush at every width tested, and below the logo's minimum it degrades to a small mark
and a plain title rather than tearing.

Selection is unmistakable: the selected row is gold end to end, description included.

**The destination is always visible.** The strip lists every place items can be
installed and marks the active one the way the menu marks its selection: a `❯` cursor
in the same column, **and** gold end to end. `tab` moves it.

**Both signals, not either.** Colour alone fails a non-colour terminal and a
colour-blind reader — strip it and the rows are identical. And encoding selection in
the *bullet* instead is worse: `●` and `○` already mean configured and not-yet, and a
green `●` reads as "on" more strongly than a gold `◉` ring, which reads as an unticked
radio button. That version shipped, made the inactive destination look selected, and
got correct behaviour reported as a bug. One channel, one meaning.

**And gold is the only colour that may mean selected.** The inactive destination's
bullet was green; green says "on" loudly enough that it kept reading as the chosen row
even with gold on the active one. Two colours arguing about selection is one colour too
many, so the strip is achromatic apart from the active row. The glyph's *shape* carries
configured-ness — which is the recession rule the palette already follows.

A prompt asked once at startup would be worse: an answer given at the top of a session
is invisible by the time you press a key, and *where did that just install?* is the
question this strip exists to answer before it is asked.

**Each row reports its own state**, not the repository's contents. `2 missing` beside
`1 linked · 1 missing` tells you which destination has what; the item counts of the
repo, filtered by the kinds a target accepts, are identical for every destination by
construction and so say nothing at all. **Two rows that cannot differ are two rows that
mislead**, and they looked authoritative while doing it.

**`tab` is listed in the footer and the switch announces itself.** Both are required,
not polish. The key changes the *destination*, not the cursor, so the `❯` stays exactly
where it was — anybody watching it sees nothing move and reports the key as broken.
That is not hypothetical; it is what happened the first time this shipped. A bullet
changing from `○` to `◉` and a path changing inside a description are evidence you have
to go looking for, and a key whose effect must be hunted for is a key that does not
work.

### Actions run in place and report inside the frame

`enter` carries the action out and shows what it did, in a section below the strip.
The panel does not close.

Staying is the point. Quitting to print the report closes the loop on the terminal and
reopens it in the reader's head — install, relaunch, tab, install again. Keeping the
panel up puts the destination, its state and the last report on screen together, which
is the only arrangement where *did that go where I meant?* has an answer you can see
rather than remember.

The report is **the command's own output**, captured rather than re-rendered. A second
rendering of the same facts is a second thing that can disagree with the first.

**The destination is passed to the action, never captured when the panel opened.** A
closure over the starting scope would send `prune` at the old destination after a tab —
the strip reading `project` while links vanish from the global config. Destructive and
silent, the worst pair.

**Where the project is, is decided once and threaded down.** Resolving it separately in
the strip and in the action gave two answers that agreed only by accident, and they did
disagree: the strip read one root while the action wrote to another.

After a successful action the figures are asked for again. They described the state
*before* it ran, and a panel whose strip contradicts its own report is worse than one
that shows neither.

### A destructive action is asked twice, in place

The first press runs it dry and shows exactly what it would remove. The second, on the
same row and the same destination, carries it out.

**Moving the cursor or switching destination disarms it.** A confirmation can only ever
apply to the plan just read, for the destination that was on screen — otherwise it could
be spent on a row walked to afterwards, or on the other config entirely. That last one is
the case that deletes from the wrong place, so it has its own test.

Nothing to remove is not something to confirm: an empty plan says so and arms nothing.

`install` and `uninstall` sit together at the top so the pair reads as a pair, and
`uninstall`'s row names the destination the same way `install`'s does.

**Only actions marked destructive ask twice.** Asking for everything teaches people to
press twice for everything, which is how a confirmation stops being read.

The earlier behaviour told the user to leave the panel and type `prune --yes` — a
confirmation step that throws away the plan it was confirming.

A refused action runs nothing and leaves no report, and with no runner wired every action
refuses — which is honest.

An earlier version set the notice to `running install…` and ran nothing, with a test
asserting that notice. **A panel that says it is working and is not is worse than a
panel with no actions at all**, and a test that pins a placeholder in place is worse
than no test: it makes the lie look checked.

The report is bounded at `MaxResultRows` and says how many lines it dropped. A
truncated list that does not admit it is a list that lies.

## Scope boundaries

**In:** the wordmark, palettes, contrast, layout, the fluid frame, the menu, the target
strip, the Bubbletea model and its navigation.

**Out:**

- **performing an action.** The model reports the choice; `cli` runs it.
- **dimming disabled rows.** Colour carries selection and nothing else, so a disabled
  row keeps full contrast and reads as available-but-inert rather than as noise.
- images, sixel, mouse. A symlink installer does not need them.
- animation. Nothing moves.

## Constraints

**Contrast floor: 4.5:1 for text, 3:1 for borders, enforced by a test rather than by
taste.** This is the constraint that earned its place: the first palette **satisfied
the spec and was unreadable** at 1.4:1 on borders. Nothing written before seeing a
terminal caught it; measuring the render caught it immediately.

Recession is achromatic, not faded — a dimmer version of a colour loses contrast, a
grey of the same lightness does not.

**No astral-plane runes in the art.** `𝄞` lives in the README and must never reach a
terminal: it is outside the BMP and renders as tofu in most fonts. `♩♪♫♬` are banned
too — East Asian Ambiguous Width makes their advance unpredictable and tears the
layout. `LIBRETTO_ASCII=safe` swaps quadrant glyphs for half blocks for fonts that
lack them.

**Fluid frame, 58–98 content columns**, centred on both axes when there is room.
Padding never truncates: a layout that silently cuts a word is a layout that lies about
what it contains.

**Every colour in every theme covers every layer.** A theme with a hole renders a
default nobody chose.

## Prior decisions

- One menu row, one colour. The gold spans the whole selected row rather than
  highlighting a fragment.
- The status row carries the live tally, so the panel states the truth rather than a
  label.
- Actions that cannot run are shown disabled rather than hidden — the panel does not
  promise what it cannot do, and it does not hide what is coming.
- `COLUMNS` is honoured when stdout is not a terminal, so centring is checkable in a
  pipe. Without it, layout would be untestable outside an interactive session.

## Task breakdown

- [x] wordmark, ASCII fallback, geometry
- [x] palettes, ramp, contrast enforcement
- [x] fluid panel, centring, degradation
- [x] menu rendering and selection
- [x] the Bubbletea model and navigation
- [x] the active destination, visible and switchable
- [ ] 6.5 confirmation form for destructive actions
- [ ] 6.6 target-strip golden files
- [ ] 6.7 `teatest` end-to-end flow

## Verification criteria

Layout and geometry:

- the frame is flush at every width
  Proof: internal/ui/fluid_test.go TestFrameIsFlushAtEveryWidth
- content width clamps to the 58–98 range
  Proof: internal/ui/fluid_test.go TestContentWidthClamping
- the panel centres when there is room and stops growing at the ceiling
  Proof: internal/ui/panel_test.go TestPanelIsCentredWhenThereIsRoom
- vertical centring holds too
  Proof: internal/ui/panel_test.go TestPanelIsCentredVertically
- a narrow terminal degrades without tearing
  Proof: internal/ui/panel_test.go TestNarrowTerminalDegradesWithoutBorders
- padding never truncates
  Proof: internal/ui/panel_test.go TestPadNeverTruncates
- the rail stays fluid and symmetric, even absurdly narrow
  Proof: internal/ui/fluid_test.go TestRailSurvivesAbsurdlyNarrowWidths

Colour, and the constraint that earned its test:

- **every palette meets the contrast floor**
  Proof: internal/ui/contrast_test.go TestPalettesAreReadable
- every theme covers every layer
  Proof: internal/ui/logo_test.go TestThemesCoverEveryLayer
- ramp steps are visually distinct
  Proof: internal/ui/contrast_test.go TestRampStepsAreDistinct
- colour is actually emitted, not merely configured
  Proof: internal/ui/logo_test.go TestColouringActuallyEmitsColour
- colouring preserves the art rather than reflowing it
  Proof: internal/ui/logo_test.go TestColouringPreservesArt
- `LIBRETTO_THEME` overrides detection
  Proof: internal/ui/contrast_test.go TestThemeEnvOverride

The wordmark:

- **no astral runes reach the art**
  Proof: internal/ui/logo_test.go TestNoAstralRunesInArt
- the ASCII-safe mode removes quadrants
  Proof: internal/ui/logo_test.go TestASCIISafeRemovesQuadrants
- the art is centred, not stretched
  Proof: internal/ui/fluid_test.go TestArtIsCentredNotStretched
- the gradient runs across the word
  Proof: internal/ui/logo_test.go TestWordmarkGradientRunsAcrossTheWord

The menu:

- the selected row is gold end to end
  Proof: internal/ui/menu_test.go TestSelectedRowIsGoldEndToEnd
- **disabled rows are not dimmed**
  Proof: internal/ui/menu_test.go TestDisabledRowsAreNotDimmed
- each row uses a single colour
  Proof: internal/ui/menu_test.go TestEachMenuRowUsesASingleColour

The model:

- navigation wraps at both ends
  Proof: internal/ui/panel_test.go TestModelNavigationWraps
- selecting a disabled action refuses rather than pretending
  Proof: internal/ui/panel_test.go TestModelSelectingADisabledActionRefuses
- a resize reaches the panel
  Proof: internal/ui/panel_test.go TestModelWindowSizeReachesThePanel
- **the active destination renders differently from the others**
  Proof: internal/ui/panel_test.go TestPanelShowsTheActiveScope
- **`tab` moves it, wraps, and asks for figures that match**
  Proof: internal/ui/panel_test.go TestModelSwitchesScope
- **a failed refresh leaves the panel as it was** rather than showing one
  destination's counts under another's name
  Proof: internal/ui/panel_test.go TestModelKeepsStateWhenRefreshFails
- with no refresh wired the key is inert
  Proof: internal/ui/panel_test.go TestModelScopeKeyIsInertWithoutRefresh
- **the footer lists `tab`**, because a key that moves nothing visible reads as broken
  Proof: internal/ui/panel_test.go TestFooterListsTheScopeKey
- **switching names the destination it now acts on**
  Proof: internal/ui/panel_test.go TestSwitchingScopeSaysSo
- **the active destination is gold end to end**, and no inactive row carries gold
  Proof: internal/ui/panel_test.go TestActiveDestinationIsGoldEndToEnd
- **with colour removed it is still marked** — the cursor survives
  Proof: internal/ui/panel_test.go TestActiveDestinationIsMarkedWithoutColour
- **no second colour competes with gold for meaning**
  Proof: internal/ui/panel_test.go TestNoSecondColourCompetesWithSelection
- **selecting an action runs it and stays**, with the report on screen
  Proof: internal/ui/panel_test.go TestSelectingAnEnabledActionRunsAndStays
- **the action is told which destination**, after a tab
  Proof: internal/ui/panel_test.go TestTheRunnerIsToldWhichDestination
- a refused action runs nothing and leaves no report
  Proof: internal/ui/panel_test.go TestSelectingADisabledActionRunsNothing
- with no runner wired every action refuses
  Proof: internal/ui/panel_test.go TestNoRunnerMeansEveryActionRefuses
- **every enabled menu label has a dispatch case**
  Proof: cmd/libretto/scope_test.go TestEveryMenuLabelDispatches
- prune chosen from the menu is still dry
  Proof: cmd/libretto/scope_test.go TestDispatchedPruneIsDry
- **a real program run installs and reports in place** — driven through `tea.WithInput`
  Proof: cmd/libretto/panelrun_test.go TestPanelRunsInstallAndReportsInPlace
- the capture returns lines and puts stdout back
  Proof: cmd/libretto/panelrun_test.go TestRunCapturedReturnsLinesAndRestoresStdout
- **a failing action still reports** what half-happened
  Proof: cmd/libretto/panelrun_test.go TestRunCapturedKeepsOutputOnFailure
- **the strip and the action agree on where the project is**
  Proof: cmd/libretto/scope_test.go TestStripAndRunnerAgreeOnTheProjectRoot
- **prune from the panel touches only the active destination**
  Proof: cmd/libretto/scope_test.go TestPanelPruneActsOnTheActiveDestinationOnly
- **a destructive action needs a second press**, and the plan stays on screen for it
  Proof: internal/ui/panel_test.go TestDestructiveActionNeedsASecondPress
- moving the cursor disarms it
  Proof: internal/ui/panel_test.go TestMovingTheCursorDisarms
- **switching destination disarms it** — the case that would delete from the wrong place
  Proof: internal/ui/panel_test.go TestSwitchingDestinationDisarms
- an empty plan arms nothing
  Proof: internal/ui/panel_test.go TestNothingToRemoveIsNotArmed
- a non-destructive action runs on the first press
  Proof: internal/ui/panel_test.go TestNonDestructiveActionRunsAtOnce
- **the second press really removes, and only from the active destination**
  Proof: cmd/libretto/panelrun_test.go TestPanelPruneConfirmsInPlace
- **one press removes nothing**
  Proof: cmd/libretto/panelrun_test.go TestPanelPruneOnOnePressRemovesNothing
- **uninstall is offered, enabled, marked destructive, and names its destination**
  Proof: cmd/libretto/uninstall_test.go TestPanelOffersUninstallAsDestructive
- **one press of uninstall removes nothing; the second does**
  Proof: cmd/libretto/uninstall_test.go TestPanelUninstallNeedsTwoPresses
- report lines keep their head, so the verb survives the elision
  Proof: internal/ui/panel_test.go TestReportLinesKeepTheirHead
- **the rows report their own state and can differ**
  Proof: cmd/libretto/scope_test.go TestStripRowsReportTheirOwnState
- the status row follows the active destination rather than summing both
  Proof: cmd/libretto/scope_test.go TestStatusRowFollowsTheActiveScope
