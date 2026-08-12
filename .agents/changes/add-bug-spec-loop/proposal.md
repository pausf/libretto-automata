# add-bug-spec-loop

Tracker: none
Queued: 2026-08-12

## The ask, verbatim

> Bug→spec loop. Cuando aparece un bug, un skill que obliga a escribir el
> criterio que faltaba en el spec ANTES de arreglar el código. Cada bug es un
> agujero del contrato; hoy el flow no tiene fase para eso — entra por
> find-work como tarea genérica y el spec se entera después, si se entera.

## Reading

A skill (or a find-work branch) for bug intake: every bug is a missing or
wrong criterion, so the fix starts by amending the owning capability spec with
the criterion that would have caught it — plus its `Proof:` — before any code
moves. Makes the contract grow from failures instead of drifting past them.
