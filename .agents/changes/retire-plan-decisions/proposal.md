# The plan's durable decisions must survive the plan

Tracker: none

## The request, in the words it was asked in

> ahora hagamos esta El gate del plan que se muere — el segundo de los dos que te
> expliqué. plan.md se borra al aterrizar y con él la tabla de alternativas; escribí la
> regla de que la decisión durable se muda a Prior decisions, y nada la mide. Sigue
> siendo una promesa escrita sin mecanismo, que es el patrón exacto que este repo
> documenta haber roto tres veces.

## The hole, and it was opened yesterday

`plan.md` is as temporary as the change that holds it. When the change lands, phase 8
applies the delta onto the capability spec and deletes the folder — plan included.

Which destroys the alternatives table: the section the previous change argued is the one
that pays, on the grounds that **a diff shows what was built and nothing shows what was
not built and why.** It survives in git history. Nobody reads deleted files.

`skills/write-plan/SKILL.md` states the answer:

> **The decision worth keeping outlives the file.** A choice that will still constrain
> work after this change lands belongs in the capability spec's *Prior decisions*, and
> phase 8 is where it moves.

That is correct, and **nothing measures it.** A promise written without a mechanism is
the pattern `AGENTS.md` documents having paid for three times: the versioning section
that sat unread for four tags, the paragraph that said "ten" over eleven directories,
the 0/24 boxes that lived only in a working tree. Each of those was written down by
somebody who meant it.

## What is proposed

A gate on the landing commit, not on the plan.

**When a commit deletes `.agents/changes/*/plan.md`, some capability spec's *Prior
decisions* section must change in the same commit.** Nothing else — not the whole spec,
that section.

**And an escape that is a declaration rather than a flag.** A change genuinely may hold
no decision worth keeping: a rename, a typo fix that grew a plan, a change whose only
alternative was "do nothing". So the plan may carry, on one line:

```
Durable decisions: none
```

Which the gate reads out of the version being deleted and accepts. It costs one line and
it is written *while the plan is being written*, by the person who knows — not typed by
whoever is trying to get a commit through.

## Scope

`spec-drift` gains the check, inside `--anchors`, on the same argument the EARS half
already used: the count "six gates" is written in ten places and a number kept in ten
places drifts.

## What is deliberately not in this change

**Nothing reads whether the migrated decision is the *right* one, or even related.** The
gate proves the section moved in the same commit. A line added to *Prior decisions* that
has nothing to do with the plan passes it. That is judgment, `review-work` has the diff
and the contract, and a script that guessed would be wrong exactly where it matters.

**No retrofit.** Changes already landed have no plan to migrate from.
