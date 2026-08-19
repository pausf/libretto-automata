# Plan — wiki-flow-board

Durable decisions: the ones in Prior decisions of spec.md

## Summary

Parse the lessons ledger's entry headers into phase counts, list the queue with
dates (readInFlight already finds them — it grows a names+dates return), render
one `#flow` article (class `cap`, so the router pages it free), link it from the
home when present.

## Technical context

Branch `feat/add-wiki-html-output` (PR #67). Blast radius: `cmd/libretto/wiki.go`,
`cmd/libretto/wiki_test.go`. Gates: the six.

## The approach

1. `readLessons(projectDir) []phaseCount` — scan `.agents/lessons.md` for
   `## ` headers, split on ` · `, count field 3; ordered by first appearance.
2. `readInFlight` return grows `queue []wikiQueued{name, date}` (it already
   opens every proposal.md); sorted by date then name.
3. Render `<article class="cap" id="flow">` after the capability articles:
   bars reuse the `.bar` styles; queue rows monospace name + date.
4. Home: a small `#flow` link beside the governed footer when the article
   exists.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| Rendering the metrics cost table | The CLI answers it; the wiki version reuses the metrics walker as its own change — the recorded condition. |
| A separate non-cap article class | The router already pages anything with class `cap`; a second class is a second router rule. |

## Risks

| Risk | What catches it |
|---|---|
| A ledger header with extra `·` in the change name miscounts | Split from the right: the phase is the last field. Force-red target: split from the left, watch the fixture with a dotted change name fail. |

## Validation and rollback

Six gates; force red as above; render-and-look (copy out of the worktree before
opening — the recorded lesson). Rollback: one revert.

## Complexity deliberately kept

None.
