# add-unattended-flow-mode

Tracker: none
Queued: 2026-08-11

## The ask, verbatim

> hacer un modo libretto-flow automatico, es decir que no haga preguntas de nada, y que
> haga la pr y push sin parar es decir que pase todos los pasos automaticamente

## Reading

An unattended mode for `/libretto-flow`: run phases 1–8 end to end without stopping at any
question, and take the branch all the way to a pushed PR.

The thing it will have to settle is which of the flow's stops are *questions* and which are
*gates*. `AGENTS.md` puts push and the release bump behind "ask first" on purpose, and a mode
that answers those on the user's behalf is the failure that published `v1.0.0`. Skipping the
six gates is a different and worse thing than skipping a clarifying question.
