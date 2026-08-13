# add-token-cost-to-metrics

Tracker: none

## The ask, verbatim

> que mas podemos aplicar en este proyecto, como podemos mejorar el sdd , hacerlo mas
> rapido utilicemos un poco de brainstorming — me gustaria mas lo de reducir token pero
> sin perder funcionalidad

Queued as "5" off the brainstorm's numbered list: the one candidate that is Go code.

## Reading

`libretto metrics` derives everything from git and measures no tokens, so any context
diet is guesswork — which already bit once: the five review lenses cost ~307k and the
lens count turned out not to be the lever.

Extend `libretto metrics` to read the Claude Code session transcripts
(`~/.claude/projects/*.jsonl`, whose entries already carry usage fields, including the
cache-read/cache-write split) and attribute token cost per change and per phase.
Cache-aware: input, output, cache creation and cache reads reported separately, because
cache reads bill at 0.1x and lumping them in makes cheap context look expensive.

This is the measurement the other queued change's follow-up (any skill-prose diet)
waits on. Research trail: session memory `token-reduction-is-not-cost-reduction`.
