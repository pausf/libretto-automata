# check-clause-coverage-in-anchors

Tracker: none
Queued: 2026-08-21

## The ask, verbatim

> Cobertura cláusula→test en `--anchors`. El patrón "un criterio cita un test que cubre
> media cláusula" apareció seis veces en tu ledger de lecciones — es tu bug recurrente
> número uno. Hoy `--anchors` prueba que la cita resuelve, no que el test ejerza la
> condición del `When`/`While`. Aunque sea heurístico, se paga solo.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Extend `spec-drift --anchors` so a cited test is checked — even heuristically — for
exercising the trigger condition named in the criterion's EARS clause. Prior art:
spec-kit `/speckit.analyze` cross-artifact coverage.
