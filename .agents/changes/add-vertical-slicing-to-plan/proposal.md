# add-vertical-slicing-to-plan

Tracker: none

## The ask

> El enfoque RPI (Research, Plan, Implement) no debe tomarse como un dogma inmutable, sino
> como una herramienta para reducir la frustración y el coste de programar con agentes de IA.
>
> [...]
>
> La regla de oro: rechazar por completo los planes que generan por defecto las herramientas
> nativas, porque tienden a planificar en "capas horizontales" (primero backend, luego
> frontend, luego tests). Exigen aplicar Vertical Slicing (Particionado Vertical). El plan
> debe estructurarse para crear pequeños flujos funcionales de extremo a extremo que puedan
> probarse y subirse a producción inmediatamente, reduciendo el riesgo.
>
> [...] que podria aportar en nuestro proyecto?

**The attribution was in the original ask and has been removed at the user's explicit
request** — the idea is recorded, the source is not named anywhere in this change. This
note exists because the ask is normally kept verbatim, so a reader comparing this to the
convention should know the edit was asked for rather than accidental.

## Reading

Of RPI's three phases, two are already covered better here — research is phase 3's brief
and the specs it produces, continuous verification is `skills/evidence/`. The gap is the
third: **vertical slicing appears nowhere in the payload.** Grepped `write-plan`,
`build-and-check` and the vendored `writing-plans` for *vertical / slice / horizontal* —
one hit, and it is about end-to-end tests, not about how a plan is cut.

The gap is structural rather than an oversight. Phase 3 is "one spec per subtask" and
specs are per capability, so the plan inherits a cut along components — which is exactly
the horizontal layering the ask rejects. `write-plan` requires each task to trace to a spec
and to name the criterion that closes it, and never requires the task to be shippable on
its own.

What the change would add is one line of contract in `write-plan`: a box closes something
that can be verified and merged alone; two boxes where the first only makes sense once the
second lands are one badly cut box.

It also buys something `libretto loop` does not have today — the loop runs one session per
open box with no guarantee that a single box leaves the tree green and mergeable.
