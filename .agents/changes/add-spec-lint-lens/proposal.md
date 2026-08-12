# add-spec-lint-lens

Tracker: none
Queued: 2026-08-12

## The ask, verbatim

> Spec linter / lens de specs. Un review lens que lee el spec antes de la fase
> 5: ambigüedades, criterios no testeables, "Governs" que no existen. Las
> cinco lenses actuales revisan código; nadie revisa el contrato.

## Reading

A review lens for the contract itself, running between write-spec and
write-plan: ambiguous wording, untestable criteria, `Governs:` paths that do
not exist, criteria that duplicate or contradict another capability. The five
existing lenses all read code; this one reads the spec before the plan is
built on it.
