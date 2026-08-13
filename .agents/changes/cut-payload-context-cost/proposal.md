# cut-payload-context-cost

Tracker: none
Queued: 2026-08-13

## The ask, verbatim

> que mas podemos aplicar en este proyecto, como podemos mejorar el sdd , hacerlo mas
> rapido utilicemos un poco de brainstorming — me gustaria mas lo de reducir token pero
> sin perder funcionalidad

Queued as "2+3+4" off the brainstorm's numbered list: the three payload-rule changes,
bundled because all three are markdown edits to the same skills and none touches Go.

## Reading

Three rules into the payload skills, from the 2026-08-13 token-cost research sweep
(three parallel research agents; findings in the session memory
`token-reduction-is-not-cost-reduction`):

1. **Select, don't load** — `find-work`, `write-spec` and `review-work` read
   `docs/SPEC.md` as the index and open only the governing capability spec, never all
   of `.agents/specs/`. Measured >10x gap between loading a corpus and loading what is
   needed.
2. **Minimal fan-out brief** — in phase 3, each spec-writer receives its subtask plus
   only the brief sections that touch it, and returns deltas and objections, never
   restatements. Every parallel worker re-pays the full brief today.
3. **Cache stability rule** — FLOW.md and the skills state: never switch model or
   effort mid-phase. Each switch invalidates the cached prefix and rebills the whole
   context at full input price instead of 0.1x. Caching is ~87% of agent billing
   (arXiv 2607.12161).
