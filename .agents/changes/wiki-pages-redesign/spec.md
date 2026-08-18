# Delta — the wiki becomes home and pages

Targets: cli

The infinite scroll dies. The viewer opens on a home of capability cards — name,
criteria count with a mini-bar, last-changed date from git — and each card opens
that capability's own page. Minimalist and visual at once: the information moves
into the cards, and the big hero chart and the scroll-spy go with the long list
that justified them.

## Outcomes

Opening a generated wiki shows the home: the project's name, two quiet totals,
a search box, and a responsive grid of cards, one per capability. A card carries
the capability's name, its criteria count with a small bar proportional to the
project's maximum, and the date its spec last changed. Clicking a card opens
that capability's page — back link, name, date, `Governs:`, intro, criteria —
alone on screen, addressable as `wiki.html#<capability>`. Back returns home.
With JS off, everything renders stacked and every anchor still navigates.

## Scope boundaries

In: the template's information architecture (home + per-capability pages via a
hash router), the card grid with mini-bars, the git last-changed date behind an
injectable seam, the home search filtering cards, and the supersession of the
hero chart and the scroll-spy.

Out, named:

- **real files per page.** One self-contained `wiki.html` stays the contract;
  `#capability` is the page address. Multi-file returns if someone needs URLs a
  hash cannot serve (crawlers, per-page HTTP).
- **commit-count / churn per spec.** Offered, not chosen. One `git log -1` per
  capability is the whole git surface.
- **history timelines, authors, diffs.** The wiki shows what is promised and
  when it last moved; git itself answers the rest.

## Constraints

- Everything not superseded stays green unchanged: marker, ownership, escaping,
  self-contained, byte-determinism for unchanged input *and history*, themes
  (token-complete, dark-first), motion behind reduced-motion, `--open`, the
  panel row, the plain-run refresh, the markdown side.
- The git date comes through `var wikiGitDate(projectDir, specPath) string` —
  exec `git log -1 --format=%as -- <spec>` by default, injectable in tests; any
  failure (no git, no repo, no history) yields the empty string and the date is
  simply absent, never an error and never a broken layout.
- The router is enhancement: visibility is toggled by JS on `hashchange`; the
  document order is home first, then every capability article, so no-JS renders
  the complete wiki.
- Cards are plain `<a href="#name">`; the mini-bar is the same integer
  arithmetic as the chart it replaces.

## Prior decisions

- **Single file, pages by hash.** Asked 2026-08-18; multi-file lost on N+1
  generated files to own, mark and clean, and it is the recorded return path.
- **Key info is last-changed date and criteria count.** Asked 2026-08-18;
  commit-count was offered and not chosen. Governs stays on the page, not the
  card — the card decides the click, the page carries the contract.
- **This supersedes two of modernize-wiki-design's criteria, hours old and
  unreleased**: the hero bar chart (its information moves into the cards'
  mini-bars) and the scroll-spy (there is no long scroll left to spy on). Their
  criteria and tests are removed at the same commit that removes the code —
  `--anchors` runs per commit, so a dead citation cannot wait for the landing.

## Task breakdown

1. The template's new architecture, the git seam, the router, the card grid —
   with the new tests and the superseded ones removed, one commit.

## Verification criteria

- The generated page shall carry a home section with one card per capability —
  each card a plain `#`-anchor showing the capability's name, its criteria
  count, and a mini-bar whose width is proportional to the project's maximum
  count — followed by one article per capability carrying a home link, its
  `Governs:` line, intro and criteria; zero-criteria capabilities get a card
  and a page like every other.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLIsHomeAndPages
- Where the git seam yields a date for a spec, the card and its article shall
  both carry it; where it yields nothing, neither shall render a date and the
  run shall succeed unchanged.
  Proof: cmd/libretto/wiki_test.go TestWikiDatesComeFromGitAndDegrade
- The page's routing shall ship as an inline `hashchange` script toggling
  visibility only — document order home-first, every article present in the
  markup, every navigation a plain anchor — so a JS-less render shows the whole
  wiki and every link still lands.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLRouterIsAnEnhancement
- The home search shall filter cards client-side by capability name and
  criteria text, inline and offline, replacing the old criteria filter.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLHomeSearchIsInline
- The redesign shall keep every non-superseded `wiki.html` proof green
  unchanged, and shall remove the superseded chart and scroll-spy criteria,
  tests and code in the same commit.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLWritesTheViewer

The proof anchors the keep-green half; the removal half cannot cite a test that
no longer exists — it is enforced by `--anchors` failing any surviving citation
of the removed tests, the same commit. Ceiling named, not covered over.
