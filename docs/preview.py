#!/usr/bin/env python3
"""Truecolor preview of the Libretto Automata panel.
Throwaway: proves the palette before it becomes internal/ui/theme.go.
"""

PANEL = [
    "╭──────────────────────────────────────────────────────────╮",
    "│  ░▒▓█ ════════════════════════════════════════ █▓▒░      │",
    "│                                                          │",
    "│    ▄▀▀▄                                                  │",
    "│   ▐▌  ▐▌    █    ▀█▀  █▀▄  █▀▄  █▀▀  ▀█▀  ▀█▀  ▄▀▄       │",
    "│ ──█▄▄▄▀──   █     █   █▀▄  █▀▄  █▀    █    █   █ █  ──   │",
    "│ ──█▀▀▀▄──   █▄▄  ▄█▄  █▄▀  █ ▀  █▄▄   █    █   ▀▄▀  ──   │",
    "│   ▐▌  ▐▌                                                 │",
    "│   ▐▙▄▄▟▘  ▏ A U T O M A T A                              │",
    "│    ▐▌     ▏ the libretto is written first ·              │",
    "│   ▄▀      ▏ the automaton performs it                    │",
    "│           ▏ b y   p a u s f                              │",
    "│                                                          │",
    "│  ░▒▓█ ════════════════════════════════════════ █▓▒░      │",
    "├──────────────────────────────────────────────────────────┤",
    "│                                                          │",
    "│  ❯ ▸ install     link the score into ~/.claude           │",
    "│    ▸ update      git pull · relink · report              │",
    "│    ▸ status      12 linked · 0 broken                    │",
    "│    ▸ doctor      diagnose the orchestra                  │",
    "│                                                          │",
    "├──────────────────────────────────────────────────────────┤",
    "│  ● claude   12 skills · 8 agents · 4 commands            │",
    "│  ○ codex    not configured                               │",
    "╰──────────────────────────────────────────────────────────╯",
    "      v0.1.0                    ↑↓ · ⏎ select · q quit",
]

# palette
PARCHMENT = "F5E6C8"   # gradient start: LIBRETTO left edge
GOLD      = "E8B44A"   # gradient end, clef, cursor
STEEL     = "8A8AA0"   # AUTOMATA, labels — the machine
MUTED     = "6A6A78"   # tagline, descriptions
DIM       = "3A3A42"   # staff lines, borders
OFF       = "4A4A55"   # inactive bullet, footer
GREEN     = "7FD1A0"   # configured
RAMP      = {"░": "4A3B1E", "▒": "6E5424", "▓": "A8822F"}  # shading rail

BRIGHT    = PARCHMENT  # selected menu label

WORDMARK_ROWS = (4, 5, 6)
WORDMARK_COLS = range(14, 52)
CLEF = set("▄▀█▐▌▙▟▘")
RESET = "\033[0m"


def fg(hexcolor: str) -> str:
    r, g, b = (int(hexcolor[i:i + 2], 16) for i in (0, 2, 4))
    return f"\033[38;2;{r};{g};{b}m"


def lerp(a: str, b: str, t: float) -> str:
    """Interpolate two hex colours. t clamped to [0,1]."""
    t = max(0.0, min(1.0, t))
    out = ""
    for i in (0, 2, 4):
        ca, cb = int(a[i:i + 2], 16), int(b[i:i + 2], 16)
        out += f"{round(ca + (cb - ca) * t):02x}"
    return out


def colour_of(row: int, col: int, ch: str) -> str | None:
    """The single source of truth for the panel's colouring."""
    if ch == " ":
        return None

    # 1. wordmark gradient wins inside its box
    if row in WORDMARK_ROWS and col in WORDMARK_COLS:
        span = len(WORDMARK_COLS) - 1
        return lerp(PARCHMENT, GOLD, (col - WORDMARK_COLS.start) / span)

    # 2. shading rail ramp
    if ch in RAMP:
        return RAMP[ch]

    # 3. clef and block glyphs
    if ch in CLEF:
        return GOLD

    # 4. structure: borders, staff lines, thin rules
    if ch in "╭╮╰╯├┤│─═▏":
        return DIM

    # 5. row-specific text
    if row == 8:                      # A U T O M A T A
        return STEEL
    if row in (9, 10):                # tagline
        return MUTED
    if row == 11:                     # signature, italic
        return MUTED
    if row == 16:                     # selected menu item
        if ch == "❯" or col == 5:
            return GOLD
        return BRIGHT if col < 19 else MUTED
    if row in (17, 18, 19):           # unselected menu items
        if col == 5:
            return OFF
        return STEEL if col < 19 else MUTED
    if row == 22:
        if ch == "●":
            return GREEN
        return STEEL if col < 11 else MUTED
    if row == 23:                     # an unconfigured target is grey throughout
        return OFF
    if row == 25:                     # footer
        return OFF
    return None


def render() -> str:
    out = []
    for r, line in enumerate(PANEL):
        buf, cur = "\033[3m" if r == 11 else "", None  # the signature is italic
        for c, ch in enumerate(line):
            want = colour_of(r, c, ch)
            if want != cur:
                buf += RESET if want is None else fg(want)
                cur = want
            buf += ch
        out.append(buf + RESET)
    return "\n".join(out)


def demo() -> None:
    """One runnable check: the colouring must never break the geometry."""
    import re
    strip = re.compile(r"\033\[[0-9;]*m")
    for r, coloured in enumerate(render().split("\n")):
        plain = strip.sub("", coloured)
        assert plain == PANEL[r], f"row {r} altered by colouring"
    widths = {len(l) for l in PANEL[:25]}
    assert widths == {60}, f"panel rows must all be 60 cols, got {widths}"
    assert lerp(PARCHMENT, GOLD, 0.0) == PARCHMENT.lower()
    assert lerp(PARCHMENT, GOLD, 1.0) == GOLD.lower()
    print("checks ok")


if __name__ == "__main__":
    print()
    print(render())
    print()
    demo()
