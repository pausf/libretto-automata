# Delta — the wiki home becomes the project board

Targets: cli

The home stops being only an index: it says what moved, what is being built, and
how healthy the contract is — from git, the spec files, and `.agents/changes/`,
read at generation time.

## Outcomes

The generated home shows, above the cards: a Recently-changed rail (the three
most recently changed capabilities, each with its date and the subject of the
last commit that touched its spec); an In-flight strip (each change under the
changes directory with open boxes, its closed/total count, plus the queued-ideas
count); and a segmented health bar — EARS criteria green, pre-EARS amber — with
the unproven count beside it. Each card carries a health dot (green: every cited
proof resolves and every criterion carries `shall`; amber otherwise). A footer
states how much of the tracked tree some `Governs:` glob claims and how many
files none does. Cards and sections enter with staggered motion, bars grow — all
behind `prefers-reduced-motion`. A project without git, without changes, without
a queue renders without those pieces: absent is absent, never an error.

## Scope boundaries

In: the home template additions and their motion; extraction for commit subject,
box counts, queue count, EARS/proof/governed checks — each behind a seam or as a
pure function; the tests.

Out, named: the spec page, the palette, the flow board (sibling changes);
relative dates ("2 days ago" — needs a clock, and determinism forbids one;
absolute dates only); any new dependency; blocking on any advisory number.

## Constraints

- **spec-drift stays the authority; the wiki is an advisory mirror.** The Go
  checks reimplement the minimal definitions the wiki displays — a criterion is
  a bullet with `Proof:` beneath it; EARS is the presence of `shall` (emphasis
  stripped); a proof resolves when the cited file exists and, for a Go test
  citation, declares the named test — and a divergence from spec-drift is a bug
  in the wiki, never a second truth. Assumed under attacca; what changes if
  wrong: the checks move behind an interface or get deleted in favour of running
  the script.
- Git access stays behind seams (`wikiGitDate` exists; a subject seam joins it);
  box counting mirrors the loop's own definition of a box.
- Determinism for unchanged input and history holds — no clock anywhere.
- Existing wiki proofs stay green unchanged.

## Prior decisions

- **Vertical cut: extraction lands with its consumer.** Assumed 2026-08-19 — a
  data layer with no surface is the horizontal box that cannot merge alone.
- **Absolute dates, no clock.** Assumed: the canvas showed "hace 1 día"; a
  relative date is a function of now and breaks byte-determinism.
- **Bar arithmetic: green is integer division, amber is the remainder.** Assumed
  2026-08-19 from the cutter's finding — the two widths always sum to 100, so the
  bar closes; if the rounding should bias the other way, it is one line.
- **A queue alone does not summon the strip.** Assumed from the cutter's finding:
  with queued proposals but zero open changes the strip is absent, queue count and
  all — the strip reports work in flight, and the flow board (sibling change) is
  where the queue lives on its own. If wrong, the fix is rendering the count solo.
- **The health dot is two states, green and amber.** No red: the wiki reports,
  gates block; a red dot would read as a gate.

## Task breakdown

1. Extraction + home surface + motion, with the proofs, one commit.

## Verification criteria

- The home shall carry a recently-changed rail of up to three of the most recently
  changed capabilities — most recent first, each with its date and the last
  commit subject from the git seam — omitting the rail entirely when no
  capability has a date.
  Proof: cmd/libretto/wiki_test.go TestWikiHomeCarriesTheRecentRail
- Where the project's changes directory holds a checklist (`tasks.md`, or the legacy `plan.md` the loop also reads) with open boxes, the
  home shall carry an in-flight strip naming each such change with its
  closed/total box count, and the count of queued proposals (a `proposal.md` carrying a `Queued:` line); where there are
  none, the strip shall be absent.
  Proof: cmd/libretto/wiki_test.go TestWikiHomeCarriesTheInFlightStrip
- The home shall carry a segmented health bar whose green and amber widths are
  the integer percentages of criteria with and without `shall` (emphasis
  stripped), and the count of criteria whose cited proof does not resolve — a
  proof resolving meaning the cited file exists and, when a test is named, the
  file declares it.
  Proof: cmd/libretto/wiki_test.go TestWikiHomeMeasuresContractHealth
- Each capability card shall carry a health marker: green where every criterion
  carries `shall` and every cited proof resolves, amber otherwise.
  Proof: cmd/libretto/wiki_test.go TestWikiCardsCarryTheHealthDot
- The home footer shall state the count of tracked files claimed by some
  `Governs:` glob and the count claimed by none, computed over the files git
  tracks; where git is unavailable the footer shall be absent.
  Proof: cmd/libretto/wiki_test.go TestWikiFooterMeasuresGovernedTree
- The redesign shall keep every existing wiki proof green unchanged, and its
  motion shall stay inside the reduced-motion guard.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLMotionRespectsReducedMotion
