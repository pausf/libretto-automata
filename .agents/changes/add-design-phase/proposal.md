# Split the plan from the task list

Tracker: none

## The request, in the words it was asked in

> mi Lead ha probado la APP y dice que echa de menos que las specs sean spec y que los
> planes esten bien montado y no solamente una lista de tareas

And, once the shape was on the table:

> pero habra un plan aparte no?
>
> las task por cierto se deberia tirar en un subagente limpio

## What is actually wrong

The flow has two artifacts where the industry has three, and the one it is missing is
the middle one.

| Artifact | Answers | In this flow, today |
|---|---|---|
| spec | what, and why | `spec.md` — exists |
| **plan** | **how** | **does not exist** |
| tasks | in what order, what is done | `plan.md` — exists, wearing the other one's name |

`skills/write-plan/SKILL.md:12` says it out loud: *"turn the specs' task breakdown into
the single file that says what is done and what is next"*. That is a task list. It is a
task list by construction, not by neglect — the task breakdown is already the sixth
pillar of the spec, so phase 5 was only ever transcribing it and adding checkboxes.

Nowhere in eight phases does anything write down **how** the change will be built: the
approach chosen, the approaches rejected and why, what could go wrong, how it gets
validated. That reasoning happens — it just evaporates with the session that had it.

A Lead opening `plan.md` and finding a checklist is the correct reaction to a file
named `plan` that is not one.

## What is proposed

Phase 5 keeps its number and changes its output. A new unnumbered seam cuts the tasks.

- **5 · The plan** — `plan.md` becomes the *how*: technical context, the decision and
  the alternatives it beat, risks, how it gets validated, and what was deliberately
  left complex.
- **Between 5 and 6 · The tasks** — `tasks.md`, the checklist that was `plan.md`.
  Cut by one fresh `task-cutter` subagent that reads the spec and the plan and nothing
  else.

The naming now matches spec-kit, which is where anybody arriving at this flow has seen
these three words before: `spec` → `plan` → `tasks`.

### Why the task cutter is a fresh subagent

Same argument the 6→7 review seam already makes, applied one phase earlier. The session
that argued its way to a design remembers what it *considered* as vividly as what it
*chose*, and cuts boxes against the argument rather than against the document. A cutter
with no memory of the discussion can only cut what was actually written down — which
doubles as a check that the plan says enough to be built from.

### Why unnumbered, and not phase 6 with everything after it shifted

`docs/FLOW.md:251` already settled this for the review seam: *"independence comes from
the fresh context, not from renumbering everything that says eight."* Renumbering would
touch README, AGENTS.md, four docs and every skill that names a phase, and buy nothing.

## Scope

The rename is a contract change: `libretto loop` and `libretto metrics` both read
`plan.md` by path. Pre-1.0, so it is a minor, and the commit carries `!`.

`metrics` works retroactively over every change ever landed — all of which have
`plan.md` in their history. It falls back rather than losing that history.

## What is deliberately not in this change

**The spec's `Task breakdown` pillar stays.** Removing it ripples into `spec-writer`,
the shared brief's six-pillar structure and `review-spec`, and it is not what was
reported. It is retargeted in one line — it feeds phase 5, and the cutter re-cuts
against the plan rather than transcribing it. Bring it back if the pillar turns out to
be what keeps producing transcribed checklists.

**EARS syntax for verification criteria** was the other half of the research and is not
here. It is independent of this change and cheap on its own.
