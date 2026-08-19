# Delta — the capability page carries its history and its proof

Targets: cli

Each capability's page stops being a flat list: an activity sparkline beside the
title, the proof behind every criterion as a chip, the Prior decisions section
surfaced above the criteria, and a Related row linking the capabilities this
spec's own text names.

## Outcomes

A capability page shows: a sparkline of commit activity over the spec (monthly
buckets from the git history, anchored to the last commit's month — no clock);
each criterion followed by a chip naming its cited proof, green when it
resolves, amber when it does not or the criterion predates EARS; a Prior
decisions block, when the spec has that section, rendered above the criteria
with its first bullets; and a Related row linking every other capability whose
name the spec's text mentions. Every piece absent when its data is: no history,
no sparkline; no section, no block; no mentions, no row.

## Scope boundaries

In: a commit-dates seam, monthly bucketing as a pure function, Prior decisions
extraction in `parseSpec`, the mention scan, the page template additions, chips
from the criterion data wiki-home-board already carries.

Out, named: the palette and the flow board (siblings); authors, diffs, or
anything beyond dates and counts; cross-links resolved semantically (a mention
is a word-boundary name match in the spec's text, nothing smarter — assumed
under attacca; if it over-links, the fix is scoping the scan to prose).

## Constraints

- The sparkline buckets anchor to the most recent commit's month, last eight
  buckets — a function of history alone, so byte-determinism holds.
- The dates seam follows the existing contract: any failure yields nothing and
  the sparkline is absent.
- Prior decisions render the section's first three bullets with a count; the
  full section lives on in the spec file itself, one link away.
- Existing wiki proofs stay green unchanged.

## Prior decisions

- **Buckets anchor to the last commit, not to today.** A "now"-anchored histogram
  changes bytes at every month boundary with unchanged history.
- **A mention is a word-boundary name match.** Assumed 2026-08-19; smarter
  linking waits for a real false-positive complaint.
- **Chips reuse wiki-home-board's proof machinery** — one definition of
  "resolves", not two.

## Task breakdown

1. Seam + parsing + page render with the proofs, one commit; then the look.

## Verification criteria

- Where the dates seam yields commit dates for a spec, the page shall carry a
  sparkline of monthly commit counts — eight buckets ending at the most recent
  commit's month, bar heights proportional to counts — and shall omit it when
  the seam yields nothing.
  Proof: cmd/libretto/wiki_test.go TestWikiPageCarriesTheActivitySparkline
- Each criterion on a page shall carry a chip naming its cited proof file and,
  when one is named, its test, marked green where the criterion carries `shall` and its proof
  resolves, amber otherwise.
  Proof: cmd/libretto/wiki_test.go TestWikiCriteriaCarryProofChips
- Where a spec holds a `Prior decisions` section, the page shall render up to
  its first three bullets with the section's total count, above the criteria;
  where it holds none, the block shall be absent.
  Proof: cmd/libretto/wiki_test.go TestWikiPageSurfacesPriorDecisions
- Where a spec file's raw text mentions another capability's name on a word
  boundary,
  the page shall carry a Related row linking each mentioned capability once,
  never linking itself; with no mentions the row shall be absent.
  Proof: cmd/libretto/wiki_test.go TestWikiPageLinksRelatedCapabilities
- The additions shall keep every existing wiki proof green unchanged.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLIsHomeAndPages
