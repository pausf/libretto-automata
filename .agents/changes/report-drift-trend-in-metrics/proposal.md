# report-drift-trend-in-metrics

Tracker: none
Queued: 2026-08-21

## The ask, verbatim

> Drift trend en `metrics`: cuántas veces el spec de cada capability se movió DESPUÉS
> del código. Convierte `spec-drift` (warn-only) en una señal de salud medible. Ninguna
> tool de la industria tiene esto todavía.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Derive per-capability drift events from git history — spec edits landing after the code
they govern — and report the trend in `libretto metrics`.
