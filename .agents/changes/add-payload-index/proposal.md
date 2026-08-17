# add-payload-index

Tracker: none

## The ask, verbatim

> Algunos directorios (skills/, agents/) contienen muchas subcarpetas/recursos; un índice
> o docs más enlazados al código ayudaría a navegar el proyecto.

## Reading

`skills/` and `agents/` are the payload — the point of the project — and there is no one
page that lists what ships. `docs/FLOW.md` describes eight phases in prose; mapping a
phase to the directory that implements it is left to the reader.

The trap to avoid is the one `AGENTS.md` was bitten by twice: a hand-written list of
directories drifts, and the paragraph that said "ten" over eleven directories is why
`docs/SPEC.md` is now the only place the spec list lives. So an index here should be
**generated or gated**, not typed — `scripts/check-payload` already walks the payload and
reads frontmatter, so it is the natural place to emit or verify the list.

Note the constraint: `scripts/` and `docs/` are not installed, so no skill may reference
the index. It is for people reading the repository, not for the payload.
