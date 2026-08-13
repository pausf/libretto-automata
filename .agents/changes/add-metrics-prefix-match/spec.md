# add-metrics-prefix-match — delta

Targets: cli

## Outcomes

`libretto metrics <name>` accepts an unambiguous prefix: `metrics drain` finds
`queue/drain-six` when nothing else starts that way. An ambiguous prefix is refused
naming the candidates; an unknown name keeps its existing refusal.

## Scope boundaries

In: `filterName`'s matching. Out: substring or fuzzy matching, interactive selection,
any change to `loop`'s change lookup.

## Constraints

An exact match always wins, even when it is also a prefix of another name — otherwise a
change whose full name prefixes a sibling becomes unreachable.

## Prior decisions

- **Assumed (attacca, 2026-08-13): prefix match, not substring.** A prefix is what a
  human types from memory of the start of a name; substring matches surprise (`metrics
  add` would match three of today's four changes). If wrong, the comparison swaps to
  `strings.Contains` with the same ambiguity rule.

## Task breakdown

1. `filterName`: exact, else unique prefix, else error listing candidates.
2. Tests for all three outcomes.

## Verification criteria

- a unique prefix selects; an ambiguous one is refused with the candidates named; exact
  still wins over a longer sibling.
  Proof: cmd/libretto/metrics_test.go TestAPrefixSelectsAChangeUnlessAmbiguous
