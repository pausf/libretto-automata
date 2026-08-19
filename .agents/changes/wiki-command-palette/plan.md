# Plan — wiki-command-palette

Durable decisions: the ones in Prior decisions of spec.md

## Summary

Build the index in Go (`searchIndex(caps)` → struct → `encoding/json` with
default HTML escaping), embed it as `<script type="application/json"
id="index">`, add the palette markup + CSS + inline script (open on the
shortcut or the search box click, substring filter over the parsed index,
grouped render, arrows/Enter/Escape), and pin structure with three tests.

## Technical context

Branch `feat/add-wiki-html-output` (PR #67). Blast radius: `cmd/libretto/wiki.go`,
`cmd/libretto/wiki_test.go`. Gates: the six.

## The approach

1. `type searchEntry struct{ Cap, Text string }`; `searchIndex` fills three
   slices (criteria texts, decisions, capability names) in caps order.
2. `json.Marshal` (default escaping turns `<` into `<` — the `</script>`
   fence) into the block after `<main>`.
3. Markup: `<div id="palette" hidden>` with input + `<div id="palette-results">`;
   CSS overlay reusing tokens; script: keydown (metaKey||ctrlKey with 'k', and
   '/' outside inputs) toggles `hidden`, input filters, arrows move
   `aria-selected`, Enter sets `location.hash` to the row's capability and
   hides, Escape hides. Classes and `hidden`, no style juggling.
4. The search box in the home gains the shortcut hint and opens the palette on
   click — the card filter beneath it stays untouched.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| Fuzzy ranking | A knob nobody asked for; substring is the recorded assumption. |
| data-* attributes for the index | Three groups with owners outgrow attributes; one JSON block is one contract. |
| Web Component / template element | The page has one palette; a div with `hidden` is the whole need. |

## Risks

| Risk | What catches it |
|---|---|
| A criterion containing `</script>` terminates the block | encoding/json's default escaping; `TestWikiSearchIndexEscapesContent` feeds exactly that fixture. Force-red: `SetEscapeHTML(false)`, watch it fail, restore. |
| The palette steals the editor's keys in other contexts | It binds only ⌘/Ctrl-K and bare `/` outside editable elements; browser-land beyond structure, untested by decision like the router. |

## Validation and rollback

Six gates; force red as above. Render-and-look: open the real page, press the
key, search, land. Rollback: one revert.

## Complexity deliberately kept

None — one block, one overlay, one listener set.
