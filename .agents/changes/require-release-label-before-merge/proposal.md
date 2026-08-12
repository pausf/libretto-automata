# require-release-label-before-merge

Tracker: none

## The ask, verbatim

> 1. Chequeo del label ANTES del merge — la más importante. Hoy el release.yml refuse si
> falta el label release:*, pero se entera después del merge: main queda adelantado y sin
> tag, exactamente el estado que AGENTS.md llama "un release que nadie puede instalar".
> Un job chiquito en gates.yml (o aparte) que corra en pull_request y falle si el PR no
> tiene exactamente un label release:* convierte ese fallo post-merge en un check rojo
> pre-merge. Combinado con branch protection que lo exija, el agujero se cierra. Son ~15
> líneas.

## Reading

A new CI job that runs on `pull_request` and fails unless the PR carries exactly one
`release:patch` / `release:minor` / `release:major` label. Today `release.yml` refuses
only after the merge, leaving `main` ahead of its last tag — the exact state AGENTS.md
names as "a release nobody can install". This moves that refusal to a red check before
the merge. Branch protection making the check required is the second half, done by the
user on GitHub, not by this change.
