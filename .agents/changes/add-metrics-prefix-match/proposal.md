# add-metrics-prefix-match

Tracker: none
Queued: 2026-08-13

## The ask, verbatim

> 3. Menor: filterName es match exacto (metrics.go:305). libretto metrics drain no
> encuentra queue/drain-six y te obliga a escribir el nombre completo de memoria. Un
> prefix match con error si es ambiguo es más amable y son cinco líneas.

## Reading

`libretto metrics <change>` should accept an unambiguous prefix (or substring) and fail
with the candidate list when the match is ambiguous. Small ergonomic fix, ~5 lines in
`filterName`.
