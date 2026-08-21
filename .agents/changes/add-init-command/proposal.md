# add-init-command

Tracker: none
Queued: 2026-08-21

## The ask, verbatim

> `libretto init` — bootstrap de `.agents/` en un proyecto virgen. Hoy es un non-goal
> explícito, y el costo es real: git no trackea directorios vacíos, así que
> `.agents/changes/` desaparece en un checkout fresco — de hecho perdiste una idea
> encolada. `init` + durabilidad de la cola matan dos pájaros.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Bootstrap `.agents/specs/`, `.agents/changes/` and `lessons.md` in a fresh project, and
make the changes directory survive git (a tracked keep-file or equivalent). Reverses the
non-goal in `.agents/specs/payload/spec.md`. Prior art: OpenSpec `/opsx:onboard`.
