# add-contributing-guide

Tracker: none
Queued: 2026-08-14

## The ask, verbatim

> Falta (o no es obvio) un archivo CONTRIBUTING.md o guía de estilo/PR que explique cómo
> colaborar o el flujo de trabajo esperado para contribuciones externas.

## Reading

Confirmed: there is no `CONTRIBUTING.md` at the root. GitHub links one automatically from
the PR and issue forms when it exists, so today an external contributor gets nothing at
the moment they need it.

What it would say already exists and is written well — `AGENTS.md` carries the gates, the
commit convention, the `release:` label rule and the boundaries. So the risk is the one
`CLAUDE.md` already names about itself: two files in sync is two sources of truth, and the
one that wins is the one nobody edited. The likely shape is a short `CONTRIBUTING.md` that
points at `AGENTS.md` and adds only what is contributor-specific — the six gates as a
checklist, and that a PR merges with a `release:` label or the run refuses.
