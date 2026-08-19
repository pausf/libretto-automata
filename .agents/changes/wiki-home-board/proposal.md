# wiki-home-board

Tracker: none

## The ask

From the design canvas the user approved 2026-08-19 ("ok implementalo"): the wiki's
home becomes a project board. This change is the home surface; three siblings follow
(spec page v2, command palette, flow board). Run under /libretto-attacca.

## Reading

The home gains: a "Recently changed" rail (top specs by git date, with the last
commit's subject); an "In flight" strip (open changes and their boxes, read from
.agents/changes/, plus the queued-ideas count); a segmented global health bar
(EARS % / pre-EARS / unproven) with a per-card health dot; a governed-tree footer
(% of tracked files a Governs: glob claims, orphan count); and staggered motion
behind prefers-reduced-motion. Every datum from sources the repo already reads:
git, the spec files themselves, .agents/changes.

Assumed (user not present): vertical cut — extraction lands with its consumer, four
changes not five. Absolute dates over "hace N días" — relative dates need a clock
and the determinism criterion forbids one.
