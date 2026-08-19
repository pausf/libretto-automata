# Delta — the palette searches the whole contract

Targets: cli

One keystroke opens a palette that searches everything the wiki knows — criteria,
Prior decisions, capability names — grouped, keyboard-driven, landing on the
owning page. For a contract of hundreds of criteria this is the difference
between a reference and a wall.

## Outcomes

Pressing the shortcut (or clicking the search affordance) opens an overlay:
typing filters a pre-built index of every criterion, every Prior decisions
bullet, and every capability name; results group under those three headings with
the owning capability named on each row; arrow keys move the selection, Enter
navigates to `#<capability>` and closes, Escape closes. The index rides the page
as one escaped JSON block built at render time — offline, no network, no
dependency. With JS off the palette simply never opens and the page loses
nothing it had.

## Scope boundaries

In: the index block, the palette markup/CSS/script, the shortcut wiring, the
tests that pin structure and data.

Out, named: fuzzy ranking (substring match only — assumed under attacca; ranking
returns when substring demonstrably fails a real search); searching Governs
globs or intro prose (the three groups are the contract's own units); replacing
the home card filter (it stays — different gesture, different scope).

## Constraints

- The index is JSON in a `<script type="application/json">` block, built with
  `encoding/json` and HTML-escaped for `</script>` safety by encoder defaults —
  spec content is data here exactly as everywhere else, never markup.
- Palette behaviour is browser-land, untested by decision like the router; Go
  pins the structure: the block parses as JSON, carries the three groups, and
  the script references the block and `keydown`.
- Byte-determinism holds — the index is a pure function of the extraction.
- Existing wiki proofs stay green unchanged.

## Prior decisions

- **Substring match, no ranking.** Assumed 2026-08-19; a ranking nobody asked
  for is a knob to maintain.
- **The index rides as JSON, not as data-attributes.** Three groups with owner
  fields outgrow attribute encoding; one parsed block is one contract.
- **The look's definition is the builder's judgment under the existing
  tokens.** Assumed 2026-08-19 (the cutter asked): the design canvas is a
  conversation artifact, not a repository one, so the render-and-look closes on
  observation against the shipped theme.
- **The home card filter stays.** Two gestures, two scopes: the filter narrows
  cards in view, the palette searches text across pages.

## Task breakdown

1. Index + palette + wiring with the proofs, one commit; then the look.

## Verification criteria

- The generated page shall carry one `<script type="application/json">` index
  block that parses as JSON and holds three groups — criteria, decisions,
  capabilities — each entry naming its owning capability and its text, built
  from the same extraction the pages render.
  Proof: cmd/libretto/wiki_test.go TestWikiCarriesTheSearchIndex
- The page shall carry the palette overlay markup hidden by default and an
  inline script that reads the index block, listens on `keydown`, and navigates
  by setting the location hash — no external script, no network, classes over
  inline style juggling.
  Proof: cmd/libretto/wiki_test.go TestWikiPaletteIsInlineAndDormant
- Spec content in the index shall arrive JSON-encoded such that a criterion
  containing `</script>` or markup cannot terminate the block or execute.
  Proof: cmd/libretto/wiki_test.go TestWikiSearchIndexEscapesContent
- The additions shall keep every existing wiki proof green unchanged.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLIsDeterministic
