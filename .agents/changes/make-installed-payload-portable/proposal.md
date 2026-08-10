# make-installed-payload-portable

Tracker: none

## The ask, verbatim

> 1. Path hardcodeado de spec-drift — 5 referencias a ~/.claude/skills/record-work/spec-drift
> en record-work y write-spec. Rompe install --project: el skill instalado en un proyecto
> apunta a un binario que puede no existir. Es el bug más real porque le pega a un usuario
> externo.
> 2. doctor no chequea rg ni jq — y spec-drift sin rg sale 0 en silencio: un gate que
> miente en verde. Dos onPath() en main.go:977 y un guard en el script. Barato y cierra
> un falso negativo.

Items 1 and 2 of a five-item improvement list; the user grouped them as one
portability change ("1 y 2 juntos — los dos son portabilidad del payload instalado").

## Reading

Both bugs bite only the *installed* payload, outside this repository:

1. `record-work` and `write-spec` name `spec-drift` by the absolute path
   `~/.claude/skills/record-work/spec-drift`. Under `install --project` the payload
   lives in `<project>/.claude`, so the path points at a binary that may not exist.
2. `libretto doctor` checks `jira` and `gh`/`glab` but not `rg`/`jq`
   (`cmd/libretto/main.go:977-993`), and `spec-drift` without `rg` exits 0 silently —
   a green gate that proved nothing.

## Not in this change

Items 3–5 of the same list (libretto-status "reporting mode", vendored superpowers
references, stale AGENTS.md/PLAN.md) — queued as their own changes, per the same
grouping.
