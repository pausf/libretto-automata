# Right-size the flow

Tracker: none

## What was asked

> "después de todo lo que hemos hecho ahora que mejorarías de mi libretto para que no
> vuelvas a preguntar tanto? por que no veo que hayas lanzado todos estos pasos al
> mandarte una tarea […] puede que no hayamos olvidado lo de que tiene que crear una
> branch nueva y que cuando haga el commit también cree la PR?"

Then, after the five improvements were proposed and accepted:

> "vale pues escribe esta tarea como un changes para que cuando lance el flow sea la
> siguiente tarea"

## Where it came from

A session that asked the user four times to update two documentation files. The
complaint was about the asking; the investigation found the asking was correct and
something else was wrong.

**The four stops are all mandated by the payload:**

| Stop | The instruction |
|---|---|
| after phase 1 | `commands/libretto-flow.md` — "Report what was found, then **wait**" |
| after phase 7 | "Then **wait**. Phase 8 begins when the user says it does" |
| the push | `skills/record-work/SKILL.md` — "ask once: push, yes or no" |
| the PR | `skills/record-work/SKILL.md` — "Opening a pull request is a **separate** question" |

So this is not an agent being timid. **The flow has one gear**, and it charges a
two-file documentation edit the same ceremony as a six-task feature.

## The two real defects underneath

**1 · The branch is decided one phase too late.** `build-and-check` owns the branch
(`skills/build-and-check/SKILL.md`, "Branch, and commit as you go"), but the only rule
that makes it mandatory — "never commit directly to the base branch" — lives in
`record-work`, phase 8. In the session that produced this proposal, two files were
edited while on `main` and the branch was created at commit time. It worked because
`git checkout -b` carries uncommitted changes. With `main` moved, or a conflicting
file in the way, it does not.

**2 · A skipped phase and a forgotten phase look identical.** `write-spec` Step 0 says
skipping the spec is "a legitimate outcome of it" — of the *phase*, not of the
orchestrator pre-empting it. That session never invoked `write-spec` or
`build-and-check` at all. The judgment was the same; the record that the judgment
happened was missing, and that is exactly what the user noticed.

## Non-goals

- **Removing the push question.** It stays. Pushing is the user's decision and always
  was.
- **Making `spec-drift` a gate.** It warns and exits 0, deliberately.
- **Letting sub-agents commit or push.** Unchanged.
- **A new configuration surface.** No flag, no env var, no `--fast`. The size of the
  change is already known from the change itself; asking the user to declare it is a
  second source of truth about the same fact.

## Cost of not doing it

The flow's own author routes around it. A process that costs four round trips for a
typo gets abandoned for typos first, then for small features, and what remains is a
ceremony reserved for work important enough to deserve it — which is the opposite of
a habit.
