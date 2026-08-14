# split-readme-into-sections

Tracker: none
Queued: 2026-08-14

## The ask, verbatim

> El README es muy largo y a veces prescriptivo; estructurarlo en secciones más cortas
> (quickstart, usage, contrib) haría la adopción más rápida.

## Reading

380 lines, 18 headings, and the first thing a newcomer meets before `## Install` is a
badge row and a rationale. The structure the ask names is roughly already there —
*What you get*, *Install*, *Your first run*, *Commands* — so the change is about length
and register, not about inventing sections.

The prescriptive part is the harder half and worth deciding explicitly: a README that
argues its own design decisions is this project's house style everywhere else
(`AGENTS.md`, `docs/FLOW.md`), and moving the arguments out to `docs/` is the obvious cut
— but it is also how a README becomes a page that says nothing and links elsewhere.
