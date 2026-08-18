# Plan — modernize-wiki-design

Durable decisions: the ones in Prior decisions of spec.md

## Summary

Rewrite `wikiCSS` and the body markup in `renderWikiHTML`, add a `renderChart`
helper emitting the SVG bars from `[]wikiCapability`, add the scroll-spy script
beside the filter, and pin the new structure with three tests while the existing
fifteen stay green untouched.

## Technical context

- Branch `feat/add-wiki-html-output` (PR #67). Blast radius:
  `cmd/libretto/wiki.go` (template consts + `renderChart`),
  `cmd/libretto/wiki_test.go` (three tests). Two files.
- The dataviz discipline loads before any chart markup is written — phase 6
  invokes the `dataviz` skill first; its rules govern the bars' form.
- Gates: the six; three new citations red until the box lands, as always.

## The approach

1. **Palette, measured** (WCAG ratios computed, not eyeballed):

   | Token | Dark | Light | Contrast on ground |
   |---|---|---|---|
   | ink | `#E7E5DE` on `#101114` | `#1F201C` on `#FAF9F6` | 14.98 / 15.56 |
   | ink-soft | `#A3A49B` | `#5A5B54` | 7.50 / 6.52 |
   | accent | `#5CB8A6` | `#1E6E60` | 7.97 / 5.77 |

   All ≥ 4.5:1. Glass surface and chart track are alpha tokens defined beside
   the hex tokens; components use `var()` only. Dark is the bare `:root`
   (dark-first per contract); light rides `prefers-color-scheme: light`.
2. **Hero bento**: CSS grid of cards — capability count, criteria total, and one
   wide card holding the chart.
3. **Chart**: `renderChart(caps)` returns `<svg role="img">` with one
   `<rect>` + `<text>` label per capability, width = `count * 100 / max` viewBox
   units (integer math), zero-count bars at width 0 but present with label.
   Sorted as the page is (lexical), bars take `fill="var(--accent)"`.
4. **Motion**: `#progress` fixed top bar and `.cap` reveal keyframes, both
   declared only inside `@media (prefers-reduced-motion: no-preference)` using
   `animation-timeline: scroll()` / `view()`. No scroll listeners.
5. **Scroll-spy**: IntersectionObserver toggling `aria-current` on nav items;
   pure enhancement — anchors work with JS off.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| A JS chart library (Chart.js et al.) | External script breaks the self-contained criterion outright; the chart is rectangles. |
| Canvas chart | Needs JS to draw, invisible with JS off; SVG is markup, printable, deterministic. |
| JS scroll listeners for progress/reveals | The research's own headline is that scroll-driven CSS replaced them; listeners cost jank and a test surface. |
| Light-first with a dark block | The contract says dark-first (content-heavy reference page); light remains a complete palette either way. |

## Risks

| Risk | What catches it |
|---|---|
| The redesign moves an anchor the old tests assert | The existing fifteen run in the same gate; the spec names them as a criterion. |
| A colour literal sneaks outside the token blocks | `TestWikiHTMLThemesAreTokenComplete` scans for hex outside the two blocks. Force-red target: drop a hex into a component rule, watch it fail, restore. |
| Zero-criteria capabilities vanish from the chart | The fixture's `checkout` has zero; the chart test asserts its label is present. |

## Validation and rollback

`go test ./cmd/libretto/` — eighteen wiki tests after this. Force red on purpose:
`TestWikiHTMLThemesAreTokenComplete` via a stray hex. Render-and-look: regenerate
against this repo's specs and open it — the look is the deliverable, so the look
is the check; contrast is already measured above. Rollback: one revert.

## Complexity deliberately kept

The chart helper is a second render function where inline string building would
do — kept because the SVG has real arithmetic in it and a separate pure function
is what makes the proportionality assertable.
