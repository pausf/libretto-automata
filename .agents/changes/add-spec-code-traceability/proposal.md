# add-spec-code-traceability

Tracker: none
Queued: 2026-08-12

## The ask, verbatim

> Trazabilidad spec→código completa (libretto trace o un skill). Hoy spec-drift
> solo avisa sobre código staged. Esto sería el mapa completo: qué código no
> gobierna ningún spec (código huérfano), qué Proof: apunta a tests débiles,
> qué capability tiene criterios sin proof. Es LA capacidad que convierte los
> specs de documentación en contrato auditable.

## Reading

A whole-repo traceability report over `.agents/specs/`: orphan code no spec
governs, criteria without a `Proof:`, proofs pointing at weak tests. Extends
what `spec-drift` and `--anchors` already check from staged-only to the full
tree. Open decision: CLI subcommand (`libretto trace`) or a payload skill.
