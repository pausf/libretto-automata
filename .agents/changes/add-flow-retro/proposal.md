# add-flow-retro

Tracker: none

## The ask, verbatim

> creo que deberiamos hacer un nuevo comando como que sea para mejorar los flujos, me
> explico, cada vez hacemos un flow con una tarea estaria guay que todos los errores que
> se han cometido en ese flow se apuntente en algun lado para luego cuando lanzamos X
> comando podamos decirle oye soluciona todo esto para que no vuelva a pasar, es decir
> como autoaprendizaje del flow , por ejemplo si en un flow el usuario te dice muchas
> veces oye pero esto lo has hehco mal el tag tal no va a asi o cualquier cosa apuntanrlo
> para solucionarlo, y cuando se lance el comando que lo arregle. supongo que estaria
> bien en la metricas cuantos errores comete la IA, me eniendes lo que quiero decir?

Chosen in conversation: **option A — ledger + explicit retro**, over an automatic
retro inside phase 8 (rejected for now: editing the payload as a side effect of
closing a task is too much power without eyes on it).

## The reading

Three pieces, agreed in the brainstorm before this flow started:

1. **Capture.** The flow's phase skills gain an instruction: when the user corrects
   the work mid-flow, append one line to a lessons ledger (`.agents/lessons.md`) —
   what happened, what the user said, which phase/skill was active. Append-only,
   cheap, no friction during the flow.
2. **Retro command.** A new payload command reads the ledger and classifies each
   lesson: *project knowledge* (the fix belongs in the project where the flow ran —
   its specs / AGENTS.md) versus *flow defect* (the fix belongs in the libretto
   payload skill, proposed as a diff, never self-applied silently).
3. **Metrics.** `libretto metrics` gains a corrections count per change — how many
   times the AI got corrected.

Open question deferred to the spec: whether the retro also has a cross-project mode
that gathers lessons from several projects before touching the payload.

This is new capability, not a bug.
