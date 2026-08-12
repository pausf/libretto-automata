# ask-before-planning-and-sync-main

Tracker: none

## The ask, verbatim

> estaria bien que cuando hagas el push y commit en la fase 8 vuelvas a main y actualices
> la rama, aparte de esto me gustaria que en el momento de crear la spec o el plan lo que
> veas mas sentido hagas mas preguntas del estilo claude, para que el plan se cree entre
> los 2 y no solo tu, puedes hacer una media de 3 pregunta de puntos importantes o dudas
> que tengas con la tarea

## Reading

Two changes to the flow, and they are unrelated except that both arrived in one sentence.

1. **Phase 8, after push.** Once the work is pushed and the request is open, return to the
   base branch and bring it up to date, so the next session does not start from a stale
   `main`. This one is not hypothetical: phase 1 of *this very run* reported a finished
   branch as work in flight, because local `main` was seven commits behind and the branch
   had already been merged and tagged.

2. **Phases 2 and 5, before writing.** Ask around three questions about the points that
   genuinely matter or that the flow is unsure of, so the spec and the plan are built with
   the user rather than handed to them. Whichever of the two phases the questions fit
   better is where they go.

The thing to settle is what counts as a question worth asking. `AGENTS.md` already says
*do not ask what the code can tell you*, and the flow's own rule is that a stop whose only
answer is "yes, carry on" is a round trip charged for a rubber stamp. Three questions per
phase, asked badly, is exactly that failure three times over.
