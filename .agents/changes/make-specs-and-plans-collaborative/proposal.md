# make-specs-and-plans-collaborative

Tracker: none

## The ask, verbatim

> Analiza el proyecto y busca en internet manera de interaccion con el usuario cuando se
> hace spec y planes, noto que el mio es muy creado por la IA y lo que sea, me falta como
> mas interaccion por parte del usuario, que sea más un copañero, otra cosa que me
> gustarua es que crear una spec y crear el plan lo haga un subagente aparte con 0
> contexto y no en el main contexto

After research (Spec Kit, Kiro, OpenSpec, Harper Reed's interview pattern, superpowers
brainstorming) a design was proposed and the user answered **"dale"** to:

1. **`decisions.md`** in the change folder — a verbatim, dated Q→A log. It is the input
   contract for the writer subagents and the source of *Prior decisions* at landing.
2. **Phase 2 interviews before drafting**: questions one at a time with `AskUserQuestion`
   (options + a reasoned recommendation first, ~5 max, "no more questions" always an
   option), each answer recorded verbatim in the log. Replaces the single batched call
   made after the pillars are already drafted.
3. **The spec is drafted by `spec-writer`**, generalized from the multi-spec fan-out to
   the single-spec case, reading brief + `decisions.md`. Where the log does not reach it
   writes `[NEEDS CLARIFICATION: question]`, never guesses; the orchestrator asks the
   remaining markers, logs the answers, and patches inline.
4. **Phase 5 presents 2–3 approaches with tradeoffs via `AskUserQuestion`** — the user
   chooses; chosen and rejected (with why) go to the log. This retires the 2026-08-12
   decision that phase 5 asks nothing, and says so.
5. **New `plan-writer` agent** modeled on task-cutter: zero context, no Write, returns
   markdown; drafts the plan from spec + `decisions.md`. The rejected-alternatives
   pillar becomes reconstructible from the log instead of living only in conversation.

Deliberately out: section-by-section validation, extra approval gates. The stop count
stays at three — the new questions ride the phases that already have them.

Not a bug: nothing already built behaved wrong. This is a change of how phases 2 and 5
interact and where their documents get written.
