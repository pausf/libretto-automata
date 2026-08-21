# add-queue-command

Tracker: none
Queued: 2026-08-21

## The ask, verbatim

> `libretto queue` en el binario. El CLI es deliberadamente ciego a la cola, pero `wiki`
> YA parsea las líneas `Queued:` — el dato está adentro, solo falta el comando de
> lectura. Barato y mejora el flujo diario.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

A read-only `libretto queue` listing queued ideas oldest-first, reusing the `Queued:`
parsing the wiki flow board already does. Reverses the queue-blind non-goal in
`.agents/specs/payload/spec.md`.
