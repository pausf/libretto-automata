# reduce-flow-stops

Tracker: none

## The ask, in the words it was asked in

> cosas raras que veo, no para de preguntarme al pasar de fases y eso no me gusta, solo
> me deberia pregunta en la fase de la spec por si quiero cambiar algo y en la del plan
> por si quiero cambiar algo, ya no deberia preguntarme mas hasta el momento del push y
> commit, cuando el reviewer encuentra algo en la mr rapidamente tiene que solucionar
> todo lo mal que encuentre sin preguntar

## Reading

Two things, and they are separable.

**1 — Too many stops.** The flow currently waits at four points: after phase 1
(`find-work` reports what it found), after phases 2–3 (the spec), after phase 5 (the
plan), and after phase 7 (`present-work`). The ask keeps exactly two of them — the spec
and the plan, the two places where the user might want to change something — and moves
every other question to the end, at commit and push.

That means:

- phase 1 reports and continues into phase 2 in the same turn
- phase 7 presents and continues into phase 8 in the same turn
- the push/commit question at the end of phase 8 stays, and is the only remaining stop

**2 — The reviewer fixes rather than asks.** When the reviewer returns findings, they get
fixed immediately and without a question. Today `review-work` "reports, it never blocks",
and acting on a finding is described as a new pass through phase 6 that the user
authorises.

## Ambiguity for phase 4

The ask says *"cuando el reviewer encuentra algo en la mr"*. Two reviewers exist and only
one of them is ours:

- `review-work` — the phase 6→7 seam, reviews **our own** finished work against its spec
- `review-project` / `/libretto-review` — reviews **somebody else's** PR/MR in a cloned
  repository, and has no mandate to change their code

"Fix everything it finds without asking" makes sense for the first and cannot apply to the
second. Confirm before phase 2 writes it.
