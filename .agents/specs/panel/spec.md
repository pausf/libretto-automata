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

### The footer's legend belongs to the screen it is under

Each screen lists its own keys: the menu's, the selector's, the catalogue's, and — over
any of them — the answers while a confirmation is open.

A legend that lists `⏎ select` over a selector where `⏎` opens a catalogue and `space`
is what marks is worse than no legend: it is read once, believed, and never re-read.
The selector's real keys used to live only in the opening notice, which the first apply
overwrote and which was then gone for the rest of the session — which is how `a`, a key
this spec had promised since the selector shipped, stayed invisible to anyone who did
not read the one frame before their first keystroke.

**The footer never outgrows the frame, at any width and for a version of any length.**
A tag name is arbitrary, so no fixed pair of legends is short enough on its own, and a
footer wider than the frame drags the whole centred block off the terminal — the same
tearing the fluid frame exists to prevent, arriving from underneath it. Three steps:
the full legend; a tight second one when it will not fit, which loses the verbs and no
key; and then the version, because the legend is what the panel is operated with and
`libretto version` still prints the version.

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

### A newer release gets one row, where it cannot vanish

The row sits **inside the frame, between the menu and the destination strip**, in the
attention colour — news, not error, because being one release behind is not a fault.

**Not the footer.** The footer already drops the version when the legend and the version
cannot both fit, and a tag name has no length limit, so a notice there would disappear at
96 columns. A notice that can vanish is a notice that was never read.

**Not `Panel.Notice` either.** That field is action feedback, and the first `install`
overwrites it — the same overwrite that once ate the selector's key legend. News deleted by
one keypress is news nobody finishes reading. It gets its own field, and moving the cursor
does not clear it.

- **It names both versions and the action** — `v0.2.0 → v0.3.0 available · choose update`. A
  row saying only "an update is available" cannot be checked against `libretto version`, and
  one with no action is a notification with nowhere to go.
- **Silence is the default.** Nothing shows when the versions match, when the check has not
  answered, or when it could not.
- **It survives the narrow layout.** `renderNarrow` drops the border, not the content —
  degrading what the panel is telling you as well as how it is drawn is two losses for one
  terminal width.
- **It is elided to the content area** like a report line, so a long notice cannot tear the
  frame.
- **The check never blocks a paint or a keypress.** `Init` returns it as a command; the panel
  renders complete without it and re-renders when it lands. Asking inline would put a
  subprocess and a network round trip in front of the first frame, and on bad DNS the user's
  only recourse is ⌃C on a tool that looks broken.
- **The callback returns the finished row, not a version.** So the semver comparison stays in
  `repo-sync` — one implementation of "is this newer", not two that can disagree — and this
  package keeps holding no notion of what a version is.
- **No check configured means no command and no notice.** That is `preview` and every test
  that does not care.

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

### The model selector — the panel's second screen

One menu entry opens it and `esc` brings the menu back. It is a **replacement, not an
overlay**: two lists of rows on one frame, each with its own cursor, is two things to
misread at a glance.

The rows are **the active destination's agents** — the same set `models` lists under
the same flag — and `tab` changes them.

```
  ❯ [ ] all
    ──────────────────────────────────────
    [x] review-lens-design   haiku        (session)   shared
    [x] sdd-apply            sonnet       low
    ──────────────────────────────────────
    [ ] review-risk          opus         xhigh
    ──────────────────────────────────────
    [ ] jd-judge-a           (session)    (session)
```

- **The rows are grouped by model**, in the catalogue's order — cheapest first, the
  session default after them, and a model this build does not know about last. That is
  the ranking the menu tally already uses, shared rather than restated: the tally row
  and the list below it are one screen, and two orderings of the same models would be
  read as a bug. Names sort inside a group, because two agents on one model in an
  order nobody chose is a list that reorders itself between sessions.
- **A rule divides one group from the next**, and a rule divides the `all` row from
  the list. One per boundary and none anywhere else — never above the `all` row, never
  after the last group, and a screen whose agents all share one model shows the `all`
  row's rule alone. It is an indented dim rule, deliberately **not** the frame's
  `├───┤` junction: that glyph means "a new section of the panel", and three of them
  inside one list would read as three panels. A blank line was the other candidate and
  it is worse — inside a bordered frame it reads as the end of the list.
- **Marking is multi.** `space` marks the row under the cursor; `a` marks every row
  and `a` again clears them — a key that only ever adds leaves no way back but
  pressing space once per row.
- **`all` is a row, and it is the discoverable half of `a`.** It sits above the first
  group, `space` on it marks every agent and `space` again clears them, and its box
  reads `[x]` only when every agent is marked — a master checkbox that stays ticked
  after one row is cleared is a checkbox that lies. `a` keeps working: taking the key
  away to justify the row would cost whoever already found it. The row is a marking
  control, not an agent — it never reaches `ApplyModel` and never appears in
  `MarkedAgents`, and `m` on it opens the catalogue for whatever is marked, like any
  other row. No agents means no `all` row.
- **The cost of the row is that screen position stops being the agent index.** Agent
  `i` sits at cursor `i+1`. A wrong offset marks the neighbouring agent and says
  nothing, which is the one failure here that is silent, so it carries a proof of its
  own.
- **One model, applied to the marked set, in one act.** `m` opens the catalogue and
  the choice reaches every marked row at once. That is the ordinary case, not the
  advanced one: making the prose lenses cheap is one gesture, not four.
- **Nothing marked means nothing happens**, and the panel says so. It never falls
  back to the row under the cursor — a selector whose marking mechanism is sometimes
  ignored teaches the user not to trust the marks.
- **The mark is legible without colour**: `[x]` / `[ ]` *and* the theme's emphasis.
  The same both-signals rule the destination strip follows.
- Rows show their new models straight after applying. A screen that needs a reopen to
  tell the truth lies for as long as it is open.
- **`esc` leaves the selector, never the program.** Sharing one key between "go back"
  and "exit" is how somebody loses the panel trying to back out of a screen.
- **`shared` is a warning, not decoration.** A row whose file this repository owns is
  reached from more than one destination: applying to it changes every project on the
  machine, applying to an unmarked one changes this destination only. The word, not a
  colour — the strip already shipped that mistake once.
- **`tab` reloads the rows, and abandons the switch if they will not load.** Moving
  the destination first and loading second leaves the strip naming one place while the
  rows below belong to another — the exact divergence the strip exists to prevent,
  produced by the code meant to honour it. Found in review, after a test that passed
  with the destination index hardcoded.
- **Effort is a second value column, and one more key.** Every row shows the level it
  declares, or the session word when it declares none — the same rendering the model
  column uses for the same state. `e` opens the effort catalogue over the marked rows,
  exactly as `m` opens the model one: same cursor, same escape, same nothing-marked
  notice, same one-act apply, same refusal to write a value every marked row already
  has. **`m` and `enter` keep meaning the model** — `enter` is the one gesture nobody
  reads the legend for, and rebinding it to pick up a second key would be a silent
  change to a reflex.
- **The effort catalogue offers only what the marked set can actually run**, and refuses
  to open when that is nothing.

  Each row carries the levels its model supports, as data through the existing seam. The
  offer is the **intersection** across the marked rows, never the union: the apply is
  all-or-nothing, so a level offered because *one* marked row can take it is a level
  guaranteed to be refused. An empty intersection is a notice naming the model that has
  none and how to proceed — *"haiku has no effort levels — unmark those rows, or move
  them off it with m"* — not a menu.

  **This reverses a decision made when the column landed**, and the reversal is the
  point: that version said which levels a model supports is not the panel's knowledge,
  and it renders the refusal it is handed. The refusal arrives at apply time. So `e` on
  two Haiku rows opened five levels, the user picked one, and the screen said no to a
  choice it had just offered — reported as a bug, correctly. **A screen that offers yes
  and then says no is worse than one that says no first**, and the layering it was
  protecting is intact: the rule still lives in `agent-models`, it just travels as data
  instead of as an error.

  The session default survives every narrowing, because removing a key is legal on any
  agent. A list whose only surviving entry is *remove* does not open either — a
  one-choice menu reads as the feature being broken rather than inapplicable.

  **Applying to the capable subset was the alternative and it loses on an existing
  promise:** nothing marked means nothing happens, and a selector whose marks are
  sometimes ignored teaches the user not to trust them. Silently skipping the Haiku rows
  is that, with extra steps.
- **The list is windowed to the terminal, and the window follows the cursor.**

  29 agents drew 29 rows, and `PlaceVertical` centring a block taller than the screen
  pushed the wordmark, the strip and the first rows off the top — so the rows you could
  not reach took the ones you could with them. Reported at exactly that scale.

  A bound alone would be the wrong fix. The action report is bounded because a report is
  *read*; this list is **operated**, and a row you cannot reach is a row you cannot mark.
  So every agent stays reachable and the window moves the least it can to keep the cursor
  inside it.

  - **It says how many rows are out of view**, at each end — the same rule as the
    truncated report, and it matters more here because the hidden rows are reachable.
  - **`all` is never scrolled away.** It is the master checkbox for the list under it, and
    a checkbox that scrolls away from what it governs is a checkbox you cannot find.
  - **An open catalogue takes its rows from the window, not from the frame.** The choices
    are what the user is looking at while it is open.
  - **An unknown height bounds nothing.** `preview`, a pipe and most tests never send a
    size, and hiding rows against a height nobody gave us would be inventing one.
  - **A group rule divides visible rows only.** A rule drawn because the row above it is
    off screen separates nothing from nothing.
  - **Below 36 lines with a catalogue open, the panel overflows deliberately.** 25 lines
    of chrome, measured, plus a three-row floor, plus the widest catalogue. A screen that
    is chrome with a peephole in it is not worth degrading into, and that terminal
    overflowed before this change with no window at all.
- **Rows group by model, and by nothing else.** Effort is a column, not a second
  grouping: grouping by the pair turns four groups into twenty on a screen whose whole
  argument is that a reader sees the shape at a glance. Applying a level moves no row.
- **The row is built against the frame, not summed and hoped.** The name column is
  measured from the longest name — not borrowed from the main menu's constant — and then
  measured *against what the frame can spare*. `pad` never truncates by design, so a row
  wider than the interior does not clip: it shifts everything after it and tears the
  border off every row on the screen. Two value columns, the `shared` warning and the
  longest name the payload ships came to 64 columns against a 58-column interior, and
  that is exactly what happened. So the name yields to the budget, **visibly, with an
  ellipsis** — a name that says it was cut is a smaller lie than a frame that came
  apart, and it is the one column whose content the two beside it cannot be inferred
  from. The value columns and the warning never yield: dropping one to save a name hides
  state rather than shortening it.

  Found by review, and the gate that should have caught it passed through the tear —
  its fixture's longest name was 15 runes and its border filter skipped agent rows
  entirely. That is a finding about the gate as much as about the row, which is why
  three criteria below name the row's width directly rather than only the frame's.
- **Applying the model every marked row already has says so and writes nothing.**
  `SetModel` will not rewrite a file that already declares the model, deliberately —
  but from outside, "nothing happened because nothing needed to" and "nothing happened
  because it is broken" are the same picture. Twice in one session the first was read
  as the second.
- An empty or absent agents directory shows a plain line saying so, not an empty box.
- **The menu row reports rather than describing itself** — `2 on haiku · 3 on session`,
  next to `status`. Every other row carries live state, and the tally is the question
  the screen was opened to answer. The session default sorts last: it is not a price,
  and a row opening with the one entry that cannot answer "how much is still
  expensive?" is not answering it.

## Scope boundaries

**In:** the wordmark, palettes, contrast, layout, the fluid frame, the menu, the target
strip, the release-notice row, the Bubbletea model and its navigation.

**Out:**

- **checking for a release.** `repo-sync` answers and `cli` decides; this package renders a
  string it was handed. The callback is the seam, in the shape `WithRefresh` and `WithRunner`
  already established.
- **a key for the update.** `update` is already a menu row with a cursor on it; an
  accelerator is a second way to reach something one keystroke away.
- **dismissing the notice.** State to persist and a file to write, for a row that goes away
  when you update.
- **a changelog, release notes, or what changed.** That is the remote's job and a browser's.
- **a spinner or a "checking…" row.** A row that appears when there is news beats one that
  narrates a network call nobody asked for.
- **performing an action.** The model reports the choice; `cli` runs it. That holds
  for the selector too: its rows and its catalogue arrive through callbacks, so this
  package still knows nothing about what an agent file is — or what a symlink is — and
  imports neither `internal/target` nor `internal/agentmodel`. The destination index is
  passed to those callbacks, never captured when the panel opened: the same rule
  `Runner` follows, and the same failure if it were not.
- **filtering or searching the agent list.** Four to eight rows fit on a screen; a
  filter is machinery for a scale this does not have.
- **any second ordering of the rows** — by name, by `shared`, by destination. One
  order, and it is the one that answers the question the screen was opened with.
- **a header naming each group's model, and collapsing a group.** Every row already
  carries its model in the model column, so a header would print it twice and cost a
  row per group. Seven rows fit. The same holds for the effort column.
- **grouping the rows by effort, or by the pair.** Named above: twenty groups is a
  different screen.
- **a third screen for the effort.** The catalogue is a mode over the same rows, like the
  model's. A screen that navigates to another screen to change one line is a screen with
  a hallway in it.
- **confirming an effort change.** `y/n` is for the destructive actions. Writing a line
  into a file the user can read is not one.
- **`ultracode` as a choice.** Not an effort level, and not something a file can declare.
- **confirming a model change.** `y/n` is for the destructive actions. Writing a
  frontmatter key is reversible in one keystroke.
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
what it contains. **The footer sits outside the frame and must not be wider than it**,
for a version of any length — otherwise it, and not the frame, decides how wide the
centred block is.

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
- **The selector sorts its rows rather than leaving them alphabetical and ruling only
  where the model changes.** Asked and answered by the user: the rows move. Grouping
  that leaves the order alone answers "where does this agent sit?"; grouping that
  reorders answers "what is still expensive?", which is what the screen is opened for.
- **One ranking function for the tally and the selector.** The delta that introduced
  grouping first invented a second order — unknown models before the session default,
  where the tally puts them after. Reading `Tally` is what caught it. Two orders on one
  screen is a bug people report as a bug.
- **A pinned version string is not a worst case for the footer.** The first footer-width
  criterion sampled one 22-character version and passed while 25 overflowed by a
  column. There is no longest tag name, so the sweep is over lengths and the last
  resort is dropping the version, not choosing a shorter legend.
- **The release notice is a frame row, not a footer item — and that follows from the line
  above.** The footer's last resort is dropping the version, so anything put beside it
  inherits the same fate. News that disappears at 96 columns is news that was never read.
- **Cached in the panel, live in `doctor`.** A panel that waits on the network before
  painting hangs on bad DNS; a diagnostic the user typed can afford five seconds.
- **`UpdateNotice`, not `Update`.** A field named `Update` beside a Bubbletea `Update`
  method has to be disambiguated by type on every read.
- **One `update` row, and its description names the outcome rather than the mechanism** —
  `bring this installation up to date`, the same string on every machine.

  Two drafts got this wrong in opposite directions. The first had **two rows**, one per command,
  with the inapplicable one disabled. The second had one row whose description **changed per
  machine** — `pull this checkout · rebuild · relink` against `install the newest version ·
  relink` — on the reasoning that a row should be honest about what pressing it does.

  What is honest about that is **the command's own output**. A menu row is a label, and a label
  naming an implementation detail is the thing this whole change came from: `git pull` in front of
  somebody who only wanted to use the tool. `pull this checkout` is the same mistake in quieter
  words.

  The consequence worth having: the description does not vary, so nothing computes it and nothing
  needs to know which mode the panel is in.

## Task breakdown

- [x] wordmark, ASCII fallback, geometry
- [x] palettes, ramp, contrast enforcement
- [x] fluid panel, centring, degradation
- [x] menu rendering and selection
- [x] the Bubbletea model and navigation
- [x] the active destination, visible and switchable
- [x] 6.5 confirmation for destructive actions — in the model, not a Huh form
- [x] the release-notice row, its field, and the `Init` seam that fills it
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
- **a destructive action asks first**, with the plan on screen and the answers offered
  Proof: internal/ui/panel_test.go TestDestructiveActionAsksBeforeActing
- answering no changes nothing
  Proof: internal/ui/panel_test.go TestAnsweringNoChangesNothing
- **no key but `y` carries it out** — not enter, not navigation, not any other rune
  Proof: internal/ui/panel_test.go TestOnlyYesCarriesADestructiveActionOut
- navigation dismisses the question rather than leaving it stale
  Proof: internal/ui/panel_test.go TestNavigationDismissesTheQuestion
- the question names the destination it was asked for
  Proof: internal/ui/panel_test.go TestTheQuestionNamesItsDestination
- an empty plan asks nothing
  Proof: internal/ui/panel_test.go TestNothingToRemoveAsksNothing
- a non-destructive action runs at once
  Proof: internal/ui/panel_test.go TestNonDestructiveActionRunsAtOnce
- **the footer offers only the answers while asking**
  Proof: internal/ui/panel_test.go TestFooterOffersTheAnswersWhileAsking
- **the second press really removes, and only from the active destination**
  Proof: cmd/libretto/panelrun_test.go TestPanelPruneConfirmsInPlace
- **showing the plan removes nothing**
  Proof: cmd/libretto/panelrun_test.go TestPanelPruneOnOnePressRemovesNothing
- **uninstall is offered, enabled, marked destructive, and names its destination**
  Proof: cmd/libretto/uninstall_test.go TestPanelOffersUninstallAsDestructive
- **an unconfirmed uninstall from the panel removes nothing; a confirmed one does**
  Proof: cmd/libretto/uninstall_test.go TestPanelUninstallNeedsTwoPresses
- report lines keep their head, so the verb survives the elision
  Proof: internal/ui/panel_test.go TestReportLinesKeepTheirHead
- **the rows report their own state and can differ**
  Proof: cmd/libretto/scope_test.go TestStripRowsReportTheirOwnState
- the status row follows the active destination rather than summing both
  Proof: cmd/libretto/scope_test.go TestStatusRowFollowsTheActiveScope

The model selector:

- the menu entry opens it and `esc` returns to the menu
  Proof: internal/ui/models_test.go TestSelectorOpensFromTheMenuAndEscapeReturns
- **`esc` in the selector does not quit the panel**
  Proof: internal/ui/models_test.go TestEscapeInTheSelectorDoesNotQuitThePanel
- `space` marks and unmarks the row under the cursor
  Proof: internal/ui/models_test.go TestSpaceMarksAndUnmarksTheCurrentRow
- `a` marks every row, and again clears every row
  Proof: internal/ui/models_test.go TestMarkAllTogglesEveryRow
- a chosen model reaches every marked row and no unmarked one
  Proof: internal/ui/models_test.go TestChosenModelReachesOnlyTheMarkedRows
- **applying with nothing marked changes nothing and says so**
  Proof: internal/ui/models_test.go TestApplyingWithNothingMarkedSaysSo
- the rows show the new model without reopening the screen
  Proof: internal/ui/models_test.go TestRowsShowTheNewModelAfterApplying
- an agent with no declared model renders as running the session's
  Proof: internal/ui/models_test.go TestUndeclaredAgentRendersAsSession
- **the mark is visible with colour stripped**
  Proof: internal/ui/models_test.go TestMarkIsLegibleWithoutColour
- the selector frame is flush at every width
  Proof: internal/ui/models_test.go TestSelectorFrameIsFlushAtEveryWidth
- **no row ever outgrows the frame**, at any width, with the longest name the payload
  ships and the `shared` warning on it
  Proof: internal/ui/models_test.go TestTheSelectorRowNeverOutgrowsTheFrameAtAnyWidth
- a name the frame cannot afford is elided visibly, and both value columns survive
  Proof: internal/ui/models_test.go TestALongNameIsElidedRatherThanTearingTheFrame
- nothing is elided when there is room
  Proof: internal/ui/models_test.go TestNamesAreNotElidedWhenThereIsRoom
- a row shows its declared effort, and the session word when it declares none
  Proof: internal/ui/models_test.go TestRowsShowTheirEffort
- `e` opens the effort catalogue, and escape returns to the rows without quitting
  Proof: internal/ui/models_test.go TestEOpensTheEffortCatalogueAndEscapeReturns
- **`m` and `enter` still open the model catalogue**
  Proof: internal/ui/models_test.go TestEnterStillOpensTheModelCatalogue
- pressing `e` with nothing marked says so and opens nothing
  Proof: internal/ui/models_test.go TestChoosingEffortWithNothingMarkedSaysSo
- a chosen level reaches every marked row and no unmarked one, and never the model apply
  Proof: internal/ui/models_test.go TestChosenEffortReachesOnlyTheMarkedRows
- the rows show the new effort without reopening the screen
  Proof: internal/ui/models_test.go TestRowsShowTheNewEffortAfterApplying
- **a refused apply leaves every row as it was, keeps the marks, and shows the reason**
  Proof: internal/ui/models_test.go TestARefusedEffortApplyChangesNoRow
- applying the level every marked row already has says so and calls nothing
  Proof: internal/ui/models_test.go TestApplyingTheEffortTheyAlreadyHaveSaysNothingChanged
- the rows still group by model alone
  Proof: internal/ui/models_test.go TestRowsStillGroupByModelAlone

The effort catalogue is narrowed to what the marked set can run:

- **`e` over rows whose model has no levels opens nothing and names the model**
  Proof: internal/ui/models_test.go TestEOnAModelWithNoEffortLevelsSaysSoAndOpensNothing
- **a mixed marked set refuses rather than applying to the capable subset**
  Proof: internal/ui/models_test.go TestAMixedMarkedSetRefusesRatherThanApplyingToTheSubset
- unmarking the incapable row opens the catalogue — the refusal is a property of the set,
  not a state the screen sticks in
  Proof: internal/ui/models_test.go TestUnmarkingTheIncapableRowOpensTheCatalogue
- **the offer is the intersection**, proven on a model with four of the five levels
  Proof: internal/ui/models_test.go TestTheCatalogueOffersOnlyTheLevelsTheMarkedSetCanRun
- **the cursor applies the entry under it in a narrowed list**, not the entry at that index
  of the full catalogue
  Proof: internal/ui/models_test.go TestTheCursorAppliesTheEntryUnderItInANarrowedCatalogue

The window:

- **the panel fits the terminal height**, at every height it can, with either catalogue
  open and 29 agents
  Proof: internal/ui/models_test.go TestTheSelectorFitsTheTerminalHeight
- **every agent is reachable and markable below the fold**, and the row it marked is on
  screen
  Proof: internal/ui/models_test.go TestEveryAgentIsReachableAndMarkableBelowTheFold
- it says how many rows are out of view, at each end, and stops claiming rows that are not
  Proof: internal/ui/models_test.go TestTheSelectorSaysHowManyRowsAreOutOfView
- `all` is never scrolled away
  Proof: internal/ui/models_test.go TestTheAllRowIsNeverScrolledAway
- an unknown height shows every row and claims nothing is hidden
  Proof: internal/ui/models_test.go TestAnUnknownHeightShowsEveryRow
- **the rows give way to an open catalogue, not the frame**
  Proof: internal/ui/models_test.go TestTheWindowShrinksWhileACatalogueIsOpen
- a failing apply reports the error and leaves the screen usable
  Proof: internal/ui/models_test.go TestFailedApplyIsReportedAndTheScreenSurvives
- the menu row reports a tally of agents by model, not a description of itself
  Proof: internal/ui/models_test.go TestMenuRowReportsTheModelTally
- the session default sorts last in that tally
  Proof: internal/ui/models_test.go TestTallyPutsTheSessionDefaultLast
- the tally refreshes after applying
  Proof: internal/ui/models_test.go TestMenuRowTallyRefreshesAfterApplying

- a shared row is marked and a local one is not
  Proof: internal/ui/models_test.go TestSharedAgentsAreMarked
- **the marker survives colour being stripped**
  Proof: internal/ui/models_test.go TestSharedMarkerIsLegibleWithoutColour
- `tab` reloads the rows for the new destination, and asks for that destination
  Proof: internal/ui/models_test.go TestTabReloadsTheSelectorForTheNewDestination
- **a failed reload abandons the whole switch** — rows and destination both stay
  Proof: internal/ui/models_test.go TestAFailedReloadKeepsTheRowsAndSaysSo
- an empty agents directory renders a plain statement, not an empty box
  Proof: internal/ui/models_test.go TestAnEmptyAgentSetSaysSo
- the tally counts the active destination's agents
  Proof: cmd/libretto/models_test.go TestMenuTallyCountsTheActiveTargetsAgents
- the model column starts in the same place whatever the names are
  Proof: internal/ui/models_test.go TestTheModelColumnLinesUpWhateverTheNamesAre
- **applying a model every marked row already has says nothing changed**
  Proof: internal/ui/models_test.go TestApplyingTheModelTheyAlreadyHaveSaysNothingChanged

Grouping, the `all` row, and the legend:

- **rows are grouped by model** — every row of one model contiguous, groups in
  catalogue order with the session default last, names sorted inside a group
  Proof: internal/ui/models_test.go TestRowsAreGroupedByModel
- a model the catalogue does not know still renders, in its own group, last — the same
  position the tally gives it, from the same shared ranking
  Proof: internal/ui/models_test.go TestAnUnknownModelGetsItsOwnGroup
- **one rule per boundary and none anywhere else** — one under `all`, one wherever the
  model changes, none above `all`, none after the last group, and a single-model screen
  shows the `all` row's rule alone
  Proof: internal/ui/models_test.go TestGroupRuleSitsOnlyBetweenGroups
- `space` on the `all` row marks every agent, and again clears them
  Proof: internal/ui/models_test.go TestSpaceOnTheAllRowMarksEveryAgent
- the `all` row's box is `[x]` only while every agent is marked
  Proof: internal/ui/models_test.go TestTheAllRowBoxFollowsEveryAgent
- **`all` never reaches the apply callback and never appears in `MarkedAgents`**
  Proof: internal/ui/models_test.go TestTheAllRowIsNeverAppliedAsAnAgent
- **the cursor marks the row it points at** — the `all` row does not shift the mark by
  one, which is the only silent failure this screen has
  Proof: internal/ui/models_test.go TestTheCursorMarksTheRowItPointsAt
- **the `all` row is legible with colour stripped**
  Proof: internal/ui/models_test.go TestTheAllRowIsLegibleWithoutColour
- **the footer lists the keys of the screen it is under** — the menu's unchanged, the
  selector's naming `space`, `a`, `m` and `tab`, the catalogue's only `⏎` and `esc`,
  and a confirmation still winning over all of them
  Proof: internal/ui/panel_test.go TestTheFooterFollowsTheScreen
- **no footer is wider than its frame**, at every width down to the floor and at every
  version length — swept, not sampled
  Proof: internal/ui/panel_test.go TestTheFooterNeverOutgrowsTheFrame
- when neither legend fits beside the version, the version goes and the legend stays
  Proof: internal/ui/panel_test.go TestTheLegendOutlivesTheVersion

The release notice:

- **it renders between the menu and the destination strip**, not in the footer where the
  version already gets dropped
  Proof: internal/ui/notice_test.go TestPanelRendersUpdateNoticeBetweenMenuAndStrip
- an empty notice changes nothing about the panel
  Proof: internal/ui/notice_test.go TestPanelOmitsUpdateNoticeWhenEmpty
- **the narrow layout keeps it** — it drops the border, not the content
  Proof: internal/ui/notice_test.go TestNarrowLayoutKeepsUpdateNotice
- news and action feedback both render, because they are two fields
  Proof: internal/ui/notice_test.go TestUpdateNoticeAndActionFeedbackCoexist
- **the frame does not tear with a notice set**, at every width — the row is measured and
  padded like any other
  Proof: internal/ui/fluid_test.go TestFrameHoldsWithUpdateNotice
- `Init` returns the check as a command, so the first paint never waits on it
  Proof: internal/ui/notice_test.go TestInitReturnsReleaseCheckCommand
- **no check configured means no command and no notice**, which is `preview`
  Proof: internal/ui/notice_test.go TestNoReleaseCheckMeansNoCommandAndNoNotice
- the message sets the row, and an empty answer sets nothing
  Proof: internal/ui/notice_test.go TestUpdateNoticeSetFromMessage
- **running an action does not clear the news** — the overwrite that once ate the
  selector's key legend
  Proof: internal/ui/notice_test.go TestActionFeedbackDoesNotOverwriteUpdateNotice
- and neither does moving the cursor, which does clear action feedback
  Proof: internal/ui/notice_test.go TestNavigationDoesNotClearUpdateNotice
