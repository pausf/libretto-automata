# Delta — the wiki viewer's modern face

Targets: cli

The generated `wiki.html` gets 2026's visual language where it carries
information: a bento hero with the specification's numbers and an inline SVG bar
chart of criteria per capability, a glass sidebar with scroll-spy, a pure-CSS
reading progress bar, and scroll-driven section reveals. Behaviour is untouched —
extraction, marker, ownership, escaping, self-containment and determinism keep
their existing criteria and their existing tests, which double as the redesign's
regression net.

## Outcomes

Opening a generated wiki shows a hero bento: stat cards (capabilities, criteria)
and a horizontal bar chart, one bar per capability, width proportional to its
criteria count, drawn as inline SVG from the same extraction — no chart library.
A thin reading-progress bar tracks scroll at the top, sections reveal as they
enter the viewport, and the sidebar highlights the section in view — progress and
reveals in CSS scroll-driven animations, the spy in inline IntersectionObserver.
Dark is the first-class theme with a complete light palette beside it, and every
motion sits behind `prefers-reduced-motion: no-preference`.

## Scope boundaries

In: the embedded template's CSS, the hero markup, the SVG chart rendering, the
scroll-spy script, and the tests that pin the new structure.

Out, named:

- **any behavioural change** — flags, discovery, refresh, ownership, escaping.
- **chart libraries or any external asset.** The bars are rectangles with
  computed widths; a library for that is the dependency ladder's last rung
  bought at its first.
- **more chart types.** A donut of proofs-vs-prose or a Governs treemap comes
  back when someone asks a question those answer; the bar chart answers "where
  does the contract live", which is the wiki's own question.
- **view transitions.** One page, no navigations to animate between.
- **kinetic typography.** The trend reports themselves confine it to marketing
  heroes; a reference page is read, not performed.

## Constraints

- Every existing `wiki.html` criterion and test stays green unchanged — the
  redesign fails if any anchor the tests assert (nav hrefs, section ids, escaped
  content, marker first line, font hosts, byte determinism) moves.
- Bar widths are computed from integers with integer arithmetic, so determinism
  stays structural.
- Motion: only inside `@media (prefers-reduced-motion: no-preference)`; the page
  is complete with zero animation.
- Both themes complete at the token level; contrast of body text on ground is
  measured in the plan, not eyeballed — the 1.4:1 palette is why.
- The scroll-spy is enhancement: with JS off, the sidebar still navigates.

## Prior decisions

- **Bento + data over editorial-minimal and glass-total.** User chose 2026-08-18
  after web research; glass stays confined to the sidebar, which is what the
  2026 retrospectives themselves recommend.
- **The chart is criteria-per-capability bars.** It answers the wiki's own
  question and derives from data already extracted; richer charts wait for a
  question they answer.
- **Dark-first, light complete.** The research is unambiguous that dark mode
  earns its place on content-heavy reference pages; the light theme remains a
  full palette, not an inversion.

## Task breakdown

1. The template rewrite with the chart, the motion, the spy — and the three new
   structure tests, existing suite green throughout.

## Verification criteria

- The generated page shall carry a hero with the capability count and the
  criteria total, and an inline `<svg>` bar chart with one labelled bar per
  capability whose width is proportional to its criteria count — zero-criteria
  capabilities drawn at zero width, never omitted.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLCarriesTheBentoHeroAndChart
- The page shall declare its reading-progress and reveal animations only inside
  `prefers-reduced-motion: no-preference`, and shall reference
  `animation-timeline` rather than scroll listeners for both.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLMotionRespectsReducedMotion
- The stylesheet shall define the complete palette — glass surfaces included —
  as tokens on `:root` with a full dark redefinition in a
  `prefers-color-scheme: dark` block, and no hex colour literal outside the two
  token blocks; components take colour only through `var()`.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLThemesAreTokenComplete
- The sidebar's entries shall remain plain `#`-anchor links — navigation with JS
  off — and the scroll-spy shall ship as an inline script referencing
  `IntersectionObserver` and `aria-current`, enhancement only. (Its runtime
  behaviour is browser-land, declared untested-by-decision with the filter and
  the theming; the structure is what Go pins.)
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLScrollSpyIsAnEnhancement
- The redesign shall keep every existing `wiki.html` proof green, unchanged.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLWritesTheViewer
