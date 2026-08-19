# Tasks — wiki-command-palette

Execution: `build-and-check` (phase 6). One writer owns this file.

- [x] 1. Index, palette, wiring — and all three proofs, one commit.
  searchIndex three groups in caps order; json with default HTML escaping into
  <script type="application/json" id="index"> after <main>; palette hidden by
  default, inline script (shortcut toggles, substring filter, grouped render,
  arrows/Enter/Escape via classes and hidden), search box opens it, card filter
  untouched. Closes when the three cited tests pass (escape test force-red via
  SetEscapeHTML(false), observed, restored), existing suite green, six gates.
  Not splittable: anchors per commit.
  - Waits on: nothing.
  - Evidence: six gates green on commit "feat(cli): the palette searches the
    whole contract"; SetEscapeHTML(false) force-red observed on the escape
    test's own message, restored. One test assertion was corrected mid-build:
    the decisions-block test asserted over the whole page and the index now
    legitimately carries every decision — scoped to the block its criterion
    promises.

- [x] 2. The look — overlay CSS on existing tokens, observed on the real page
  (open, press the key, search, land). Builder's judgment is the definition
  (recorded assumption). Existing theme/self-contained/determinism proofs green.
  - Waits on: 1.
  - Evidence: opened for real via --open; the index parses on the generated
    page — 565 criteria, 143 decisions, 14 capabilities; palette dormant markup
    present.
