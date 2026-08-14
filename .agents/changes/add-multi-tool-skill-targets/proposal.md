# add-multi-tool-skill-targets

Tracker: none

## The ask, verbatim

> Fase 1 — skills multi-target: targets codex y opencode que enlazan skills a
> ~/.agents/skills (y ~/.config/opencode como opción). Puro symlink, cero
> transformación, el selector de destino en el panel que pediste.

> y tambien ten e cuenta que nuestra skills mucha veces dice, pregunta como
> claude, pero se tiene que adaptar s los 3

## Reading

New `target.Target` implementations for Codex CLI and OpenCode that symlink the
skills kind only — both tools read Claude-compatible `SKILL.md` directories, so
no transformation is needed. Codex reads `~/.agents/skills/`; OpenCode reads
that same path plus `~/.claude/skills/` and `~/.config/opencode/skill(s)/`, so
one link into `~/.agents/skills` serves both. The panel gains a target selector
alongside the existing global/project scope choice. Feasibility research with
paths and citations: memory note `multi-tool-target-viability` (2026-08-14).

The trailing note applies here too: payload skills whose prose says "Claude"
where it means "the agent" must read correctly in all three tools once they are
installed beyond `~/.claude`.
