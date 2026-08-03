# Design

Visual and technical decisions for Libretto Automata. This is the reference
for anything that renders. **Do not redraw the art from memory — copy it from
this file.**

## Name and metaphor

An 18th-century automaton was a machine that played music by reading a score. A
human wrote the notes; the machine executed them. It never improvised.

Written first, performed second. That is the whole product.

- Project: **Libretto Automata**
- Directory: `libretto-automata` — no symbol, no spaces, so `cd`, symlinks and
  scripts never need quoting
- Binary: `bin/libretto`; command `libretto`

## The `𝄞` rule

`𝄞` (U+1D11E MUSICAL SYMBOL G CLEF) lives in Unicode's SMP plane. SF Mono does
not have it. Most Nerd Fonts do not have it. On the author's machine it renders
via font fallback; on a user's machine it renders as `□`.

**Use `𝄞` only in `README.md`**, which is rendered by GitHub as web content.
**Never emit `𝄞` to a terminal.** The terminal clef is drawn by hand from
block glyphs (below).

The same rule kills `♩ ♪ ♫ ♬`: they are BMP, but classified East Asian
Ambiguous Width. Terminals disagree on whether they occupy one column or two,
and a two-column render tears the box layout apart. Do not use them.

## Main panel

```
╭──────────────────────────────────────────────────────────╮
│  ░▒▓█ ════════════════════════════════════════ █▓▒░      │
│                                                          │
│    ▄▀▀▄                                                  │
│   ▐▌  ▐▌    █    ▀█▀  █▀▄  █▀▄  █▀▀  ▀█▀  ▀█▀  ▄▀▄       │
│ ──█▄▄▄▀──   █     █   █▀▄  █▀▄  █▀    █    █   █ █  ──   │
│ ──█▀▀▀▄──   █▄▄  ▄█▄  █▄▀  █ ▀  █▄▄   █    █   ▀▄▀  ──   │
│   ▐▌  ▐▌                                                 │
│   ▐▙▄▄▟▘  ▏ A U T O M A T A                              │
│    ▐▌     ▏ the libretto is written first ·              │
│   ▄▀      ▏ the automaton performs it                    │
│           ▏ b y   p a u s f                              │
│                                                          │
│  ░▒▓█ ════════════════════════════════════════ █▓▒░      │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  ❯ ▸ install     link the score into ~/.claude           │
│    ▸ update      git pull · relink · report              │
│    ▸ status      12 linked · 0 broken                    │
│    ▸ doctor      diagnose the orchestra                  │
│                                                          │
├──────────────────────────────────────────────────────────┤
│  ● claude   12 skills · 8 agents · 4 commands            │
│  ○ codex    not configured                               │
╰──────────────────────────────────────────────────────────╯
      v0.1.0                    ↑↓ · ⏎ select · q quit
```

Composition, top to bottom:

1. **Shading rails** — `░▒▓█ ═══ █▓▒░`. The ramp fades the staff into the
   border instead of stopping dead.
2. **Clef** — drawn at left, crossing two staff lines, exactly as on a real
   stave.
3. **Wordmark** — `LIBRETTO` in a three-row half-block face. Half blocks give
   double vertical resolution, which is what makes it read as a typeface rather
   than as ASCII art.
4. **Subtitle block** — `AUTOMATA`, the tagline, and the signature, each line
   prefixed `▏` as a thin rule.
5. **Menu** — one row per action, `❯` marking the cursor, `▸` as a static
   bullet. The selected row is gold end to end.
6. **Target strip** — one row per agent target. `●` configured, `○` not. This
   is where `codex` and anything else will appear.
7. **Footer**, outside the border — version left, key hints right.

## Small mark

For the shell prompt prefix, spinners, and one-line output:

```
 ▄▀
▄█▀
▀▄
```

## Palette — warm is human, cold is machine

The colour scheme encodes the metaphor. The part a person writes is warm ink on
parchment. The part a machine executes is cold steel. Warmth is scarce and
deliberate: it marks the clef, the wordmark, and the one row the cursor is on.
Everything else recedes.

### The contrast floor

**Every text colour must clear 4.5:1 against the terminal background, and every
border or rule 3:1.** This is not a nicety. The first palette shipped borders at
`#3A3A42` (1.4:1 on a dark terminal) and descriptions at `#6A6A78` (3.4:1); it
looked refined in a mockup and read as "nothing is visible" in a real terminal.

Recession is expressed by **staying achromatic**, not by fading toward the
background. A grey that is still clearly legible recedes plenty next to gold.

The floor is enforced by a WCAG luminance test over both palettes against four
backgrounds — see `internal/ui/contrast_test.go`. Changing a colour without
running it is how the panel goes invisible again.

### Named colours — dark terminal

| Name | Hex | Used for |
|---|---|---|
| `parchment` | `#F7EAD2` | wordmark gradient start |
| `gold` | `#F0BE52` | gradient end; clef; signature; the selected menu row, end to end |
| `steel` | `#E4E4EC` | ordinary text: `AUTOMATA`, unselected menu rows, target name |
| `muted` | `#AEAEBE` | tagline; target counts; the notice line |
| `dim` | `#6A6A78` | borders, staff lines, `═` rails, `▏` rules |
| `off` | `#8A8A98` | `○`, unconfigured target row, footer |
| `green` | `#8FE3B0` | `●` configured |
| `error` | `#FF8878` | failures, conflicts |

### Named colours — light terminal

Same hues, lightness axis inverted. `steel` becomes near-black for the same
reason it is near-white on dark: it is the panel's ordinary text colour.

| Name | Hex |
|---|---|
| `parchment` | `#5A4A24` |
| `gold` | `#8A6A1E` |
| `steel` | `#1C1C26` |
| `muted` | `#4E4E5A` |
| `dim` | `#85858F` |
| `off` | `#6E6E7C` |
| `green` | `#1B7A47` |
| `error` | `#A82A1A` |

### The shading rail

`░▒▓█` is a four-step ramp into `gold`, so the staff fades out of the border
rather than stopping dead. Consecutive steps must stay at least 1.2:1 apart or
the ramp reads as a flat block.

| Glyph | Dark | Light |
|---|---|---|
| `░` | `#8A6E36` | `#C4A868` |
| `▒` | `#B08C42` | `#A88A44` |
| `▓` | `#D2A84A` | `#97781F` |
| `█` | `#F0BE52` | `#8A6A1E` |

### Overriding detection

Terminals misreport their background often enough that `LIBRETTO_THEME=dark|light`
forces a palette. An unrecognised value falls back to detection rather than
erroring.

### Colouring rules

Resolved in this order, first match wins. Row and column indices are into the
panel art as written above, zero-based.

1. A space is never coloured.
2. **Rows 4–6, columns 14–51** — the wordmark box. Colour is
   `lerp(parchment, gold, (col-14)/37)`, interpolated **per column, not per
   letter**, so the gradient runs smooth across the whole word instead of
   stepping letter by letter. This rule wins over every glyph rule below, which
   is why the wordmark's `█` reads gradient while the clef's `█` reads gold.
3. `░ ▒ ▓` — the ramp table above.
4. `▄ ▀ █ ▐ ▌ ▙ ▟ ▘` — `gold`. The clef.
5. `╭ ╮ ╰ ╯ ├ ┤ │ ─ ═ ▏` — `dim`. All structure, one colour. Staff lines and
   borders share it deliberately, so a staff line running into a border is one
   continuous stroke.
6. Row-specific text in the art block, by row index into `artRows`:
   - `AUTOMATA` → `steel`, the ordinary text colour.
   - The two tagline rows → `muted`.
   - The signature row → `gold`. It is letter-spaced like `AUTOMATA` and reads
     warm, so it carries more weight than the muted tagline above it. This is the
     one place warmth is spent on something other than the clef, the wordmark and
     the selected row.
7. Menu rows. **One row, one colour.** Labels begin at column 7 and
   descriptions at column 19 in every row, which is what lets one rule cover
   them all.
   - Unselected → `steel`, the ordinary text colour.
   - Selected → `gold`, **end to end**: cursor, bullet, label and description
     together.

   An earlier version split each row across a cursor colour, a label colour and
   a description colour. That made the cursor compete with the label for
   attention instead of leading the eye to it. A single sweep of gold is
   unmistakable and needs no legend.

   Disabled rows are **not** dimmed. Colour carries selection and nothing else;
   an action that cannot run refuses behaviourally, with a notice. Encoding
   availability in the palette would put two meanings on one channel.
8. Target rows. `●` `green`, name `steel`, counts `muted`. An unconfigured
   target (row 22) is `off` end to end — the whole row recedes, not just its
   bullet.
9. Row 24, the footer → `off`.

Wrap every colour in `lipgloss.AdaptiveColor` so the panel survives a light
terminal.

**Colouring must never change geometry.** Stripping every ANSI escape from a
rendered row must reproduce the source row byte for byte. That is a test, not a
hope — see PLAN.md 6.1.

A working truecolor implementation of all of the above lives in
`docs/preview.py`. It is a throwaway reference for porting to
`internal/ui/theme.go`, not shipped code.

## Layout rules

**Never hardcode a width or count spaces.** The art above is a mockup, not a
layout. Lipgloss does layout:

- `Style.Width()`, `Padding()`, `Align()`
- `Style.Border(lipgloss.RoundedBorder())`
- `lipgloss.JoinHorizontal` for clef-beside-wordmark
- `lipgloss.JoinVertical` for stacking sections
- `lipgloss.Place` for the footer

Read `tea.WindowSizeMsg` and degrade cleanly. Below the width the logo needs,
drop to the small mark plus a plain title. A torn box is worse than no box.

### Fluid width

The panel tracks the terminal width between two bounds:

| | Columns | Why |
|---|---|---|
| Minimum content | 58 (`ArtWidth`) | narrower than the drawing, nothing to draw |
| Maximum content | 98 | past this the layout gets worse, not better |

The ceiling is not arbitrary. At 140 columns a label sits at column 4 and its
description ends forty columns away with a desert between them; the eye has to
travel. Newspapers use narrow columns for the same reason. Beyond the ceiling the
panel stops growing and centres instead.

What stretches and what does not:

- **Stretch:** the border, the section rules `├───┤`, and the shading rails. The
  rails are symmetric — `░▒▓█ ═══ █▓▒░` with equal margins. The fixed-width
  version had six leftover columns on the right; a fluid rail has no excuse.
- **Never stretches:** the clef and wordmark. They are a drawing. They are
  centred within the content width and nothing more.
- **Stays left:** the menu and the target strip. A centred list is hard to read
  because the eye loses the edge it returns to.

The footer keeps a six-column margin on both sides, mirroring the art's own
indent, so it reads as balanced under the box rather than running to its corners.

### Centring

The panel is centred in the terminal, each axis only when that axis has room to
spare:

- **Horizontally** when the terminal is wider than `PanelWidth`. A 60-column box
  pinned to the left edge of a 200-column terminal looks unfinished.
- **Vertically** when the terminal is taller than the rendered block. This
  applies in the TUI, which owns the whole alternate screen. `preview` prints
  inline and passes no height, so it never pads vertically.

With no known size — a pipe — the block stays flush left and unpadded, so piped
output carries no alignment whitespace. `COLUMNS` is honoured as a width
fallback, which is the only way to check centring without a terminal.

The notice line is composed **inside** the centred block, not appended after it.
Appending would shift the panel up by a row the moment feedback appeared.

## Font fallback

The art uses three glyph classes, all from Unicode's Block Elements range
(U+2580–U+259F). Verified by extracting them from the panel above:

| Class | Glyphs | Support |
|---|---|---|
| Half blocks and full block | `▀ ▄ █ ▌ ▐` | universal |
| Shading | `░ ▒ ▓` | universal |
| Eighth block (the `▏` rule prefix) | `▏` | universal in monospace fonts |
| **Quadrants** | `▘ ▙ ▟` | very high, **not** universal |

Only the quadrants are a risk, and there are exactly three of them: `▘` closes
the clef's lower loop, `▙` and `▟` form its bowl.

Ship a second clef variant that replaces `▘▙▟` with `▀█▀`, selected by
`LIBRETTO_ASCII=safe`. It looks slightly squarer and never breaks.

**No font autodetection.** There is no reliable way to query terminal font
coverage. An env var the user sets once is honest; a guess is not.

## Stack

Go, with **Bubbletea** for state, **Lipgloss** for style, **Huh** for forms.

Chosen because it ships a single binary with no runtime to install,
cross-compiles, starts in about 2ms, and Charm is the strongest TUI toolkit
available. It is also gentle-ai's stack, so the two can interoperate later.

**opencode was measured, not assumed.** Inspecting the installed v1.17.15
binary found `opentui`, `yoga`, `solid-js`, `truecolor`, `kitty` and `sixel`
strings, and no Go runtime. opencode is TypeScript on its own TUI framework with
flexbox layout and real pixel graphics via the Kitty protocol and Sixel. That
ceiling is higher than Lipgloss's, but this tool is a symlink installer with a
four-item menu — it would use none of it, and would drag in a Bun runtime to
get there.

## Why symlinks, per item

Copying means every edit needs a reinstall. Symlinking means the repo *is* the
installation: edit a skill, it is live in the next session.

**Per item, never per directory.** `~/.claude/skills/` also holds skills
installed by gentle-ai's sync. Symlinking the whole directory would erase them.
Symlinking each item lets both live in the same folder, and `git pull` only ever
touches what this repo owns.

Precedent: `~/.claude/skills/jira-cli` is already a symlink
(`-> ../../.agents/skills/jira-cli`) and Claude Code loads it without
complaint. The mechanism is verified, not hoped for.

## Why the target interface exists

`internal/target` defines a `Target` interface with one implementation
(`claude.go`).

Normally a single-implementation interface is speculative and should be cut.
Here a second target is committed, not imagined — Codex, and others after it.
The interface is what makes that a new file instead of a refactor.

If that commitment ever disappears, collapse the interface.
