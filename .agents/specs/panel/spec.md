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

## Scope boundaries

**In:** the wordmark, palettes, contrast, layout, the fluid frame, the menu, the target
strip, the Bubbletea model and its navigation.

**Out:**

- performing any action. The model dispatches; `cli` and `linking` do the work.
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
