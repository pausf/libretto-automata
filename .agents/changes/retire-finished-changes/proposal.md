# Finish landing the two changes that merged without landing

Tracker: none

## The request

`/libretto-flow`, with no argument, immediately after this was reported:

> **Las dos carpetas de change siguen en `main`.** Yo dije explícitamente que las dejaba
> hasta el merge para que el que revisara el PR pudiera leer el proposal y el plan — y
> después el merge pasó sin borrarlas.

Phase 1 read it the same way from the tree, independently: eight boxes, none open, two
change folders present, no open request, nothing queued.

## What is wrong

`add-design-phase` and `retire-plan-decisions` merged to `main` as `v0.30.0` with their
work finished and their landing incomplete. `record-work` names four things the last
commit of a change does together:

1. the final code — **done**
2. the delta applied onto the capability spec — **done**
3. the plan's durable decisions retired into *Prior decisions* — **not done**
4. the change folder deleted — **not done**

> All four or none. A delta applied without deleting the change leaves two documents
> describing the same capability, and the next reader has no way to tell which one is
> current.

That is the state on `main` right now: `payload/spec.md` carries the outcomes, and two
folders beside it describe the same capability in the past tense.

## Why step 3 is the expensive half

Four decisions currently exist only inside plans that have to be deleted:

- **`plan.md` was reused for the technical approach rather than adding a `design.md`.**
  A new file beside an unchanged `plan.md` leaves the reported complaint exactly where
  it was.
- **The task cut runs in a subagent, not a numbered phase.** Independence comes from the
  fresh context; renumbering buys nothing.
- **The retirement gate compares the *section*, never the file.** Requiring any edit to a
  capability spec is green on every landing by definition and measures nothing.
- **Its escape is a declaration in the plan, not a flag.** A flag is typed by whoever
  wants the commit through; a line in the plan is written by the person who knew.

Every one of those will be re-litigated by somebody who cannot see it. That is the
failure the previous change built a gate against.

## And it is that gate's first real test

`spec-drift --retired` has never run against a real landing, because the deletion that
fires it has never happened. This change is the deletion.

**It should fail before it passes.** Deleting the folders without moving the decisions is
exactly the commit the gate exists to refuse — so that is the order: stage the deletion,
watch it refuse, then migrate and watch it accept. A gate whose first real encounter is a
pass has proved nothing about itself.

## Scope

`.agents/specs/payload/spec.md` gains the decisions. Both change folders go. Nothing else
moves — no code, no skill, no CLI behaviour.

## What is deliberately not in this change

**The two capability specs are not migrated to EARS.** They are amended here, and the
migration rule says a capability moves to EARS when a delta lands on it — but this delta
adds no criteria, only prior decisions, and rewriting 545 legacy criteria is the change
that was already declined once on its own merits.

**No decision about `add-design-phase`'s sequencing is kept.** How *that* change was cut
into boxes dies with it, correctly. What survives is what would otherwise be decided
again, differently.
