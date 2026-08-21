# add-loop-budget-and-run-log

Tracker: none
Queued: 2026-08-21

## The ask, verbatim

> Mejoras a `loop`: presupuesto de coste por ronda (alimentando `metrics`), un log
> persistido por ronda (hoy imprime y se pierde). Los rieles de seguridad (spend limits,
> checkpoint por iteración) son exactamente lo que la comunidad de Ralph-loops añadió en
> 2025-26.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Give `libretto loop` a per-round cost budget and a persisted per-round run log, so a run
is inspectable after the fact and `metrics` can consume it.
