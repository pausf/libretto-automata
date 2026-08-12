# add-ralph-loop-runner

Tracker: none
Queued: 2026-08-12

## The ask, verbatim

> ok metelo en queue el add-ralph-loop-runner para el libretto-attacca,
> entiendo que tendras que crear un capreta llamada loops en el rpyecto y que
> tambien haga symlink como command , skills etc ?

## Reading

A Ralph-Wiggum-style loop over the existing flow: relaunch a fresh-context
session per task (`claude -p`), each iteration reading state from
`.agents/changes/<change>/plan.md`, taking one unchecked box, building,
gating, committing, exiting. Attacca already has the guardrails (assume-and-
record, stop after two failed gates, verifiable done); what is missing is only
the fresh-context engine between tasks.

Open decision the user raised: where it lives. A new `loops/` payload category
symlinked like skills/commands — or no new category at all: the loop must run
OUTSIDE the session it relaunches, so a skill/command cannot BE the loop,
which points at a `libretto loop <change>` Go subcommand (ships in the binary,
nothing to symlink). Decide in phase 2, not here.

Researched (2026-08-12): the canonical Ralph playbook (paddo.dev, endorsed by
Huntley) does create a file structure — loop.sh, PROMPT_plan.md,
PROMPT_build.md, AGENTS.md, IMPLEMENTATION_PLAN.md, specs/*.md. Every one of
those except loop.sh already has a richer Libretto equivalent (the write-plan
and build-and-check skills, AGENTS.md gates, plan.md checkboxes,
.agents/specs/). Copying the structure would duplicate state; the only
genuinely missing artifact is the loop engine itself.

Guardrails from the sources, to keep: machine-verifiable done (boxes + gates),
iteration cap, stuck detection (same box unchecked two rounds → stop and
report), never push/merge/tag beyond what attacca already answers.
