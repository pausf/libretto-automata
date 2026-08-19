# Plan — wiki-home-board

Durable decisions: the ones in Prior decisions of spec.md

## Summary

Grow the extraction into a small analysis pass over data already in hand, then
render the four new home blocks. All git behind seams; all checks pure functions
over `wikiCapability` plus a criteria structure that now keeps the `Proof:` line
it used to drop.

## Technical context

- Branch `feat/add-wiki-html-output` (PR #67). Blast radius: `cmd/libretto/wiki.go`,
  `cmd/libretto/wiki_test.go`. Two files. Gates: the six, citations red until the
  one box lands, as every change on this branch.

## The approach

1. **Extraction**: `wikiCriterion{text, proofFile, proofTest string}` replaces the
   plain string (render side unchanged — it prints `.text`). `parseSpec` keeps the
   first `Proof:` line per bullet, split into file (first token) and test (second,
   when present).
2. **Seams**: `wikiGitDate` stays; add `wikiGitSubject(projectDir, specPath)`
   (`git log -1 --format=%s -- <spec>`, failure → "") and
   `wikiGitTracked(projectDir)` (`git ls-files -z`, failure → nil).
3. **Pure checks** over the parsed set, one small func each: `isEARS(text)`
   (strip `*` and backticks, word-boundary `shall`); `proofResolves(projectDir,
   c)` (file exists; `.go` test citation → file contains `func <Test>(`);
   `capHealthy(caps[i])`; `governedSplit(tracked, allGoverns)` (glob match via
   `path.Match` per segment — the same matching family spec-drift's globs use).
4. **Project state**: `readInFlight(projectDir)` walks `.agents/changes/*/`
   counting `- [ ]`/`- [x]` in `tasks.md` (fallback `plan.md`, the loop's rule)
   and `Queued:` proposals. Changes with zero open boxes are skipped.
5. **Render**: rail (up to 3 by date desc, name+date+subject), strip, segmented
   bar + unproven count, per-card dot, footer. Motion joins the existing
   reduced-motion block: staggered `rise` delays and `fill` bars, sentinel
   `/*end-motion*/` kept.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| Shelling out to spec-drift for health | The script lives in the payload install, at a path the binary cannot assume in an arbitrary project; and parsing its human output is a second contract. The minimal Go mirror is ~40 lines with spec-drift kept authoritative by decision. |
| Relative dates | A function of the clock; determinism criterion forbids it. |
| One `git log` for all subjects | Subjects join dates per-spec through the same seam shape; batching is an optimization for a page generated at landings. |

## Risks

| Risk | What catches it |
|---|---|
| Go health mirror drifts from spec-drift | Named prior decision: advisory mirror, spec-drift authoritative; the criterion pins the exact definitions mirrored. |
| Glob matching diverges from spec-drift's | `governedSplit` uses full-path segments with `**` handling mirroring the script's `matches_any`; the footer test carries a `**` fixture. Force-red target: break the `**` arm, watch the footer test fail, restore. |
| Strip renders for landed changes | Only open boxes count; zero-open changes skipped, tested. |

## Validation and rollback

`go test ./cmd/libretto/`; the six gates. Force red: the `**` arm of
`governedSplit`. Render-and-look at the end of the change against this repo via a
throwaway worktree. Rollback: one revert.

## Complexity deliberately kept

`wikiCriterion` as a struct where a string sufficed yesterday — the proof data is
what three of the six criteria consume, and the spec page (next change) needs it
too.
