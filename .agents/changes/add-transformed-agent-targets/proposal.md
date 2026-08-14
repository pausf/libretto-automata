# add-transformed-agent-targets

Tracker: none
Queued: 2026-08-14

## The ask, verbatim

> Fase 3 — agents transformados: md→TOML para Codex, remapeo de frontmatter
> para OpenCode, y el modelo de propiedad para archivos generados. Esta fase es
> la cara y la que merece spec propia.

> Los custom prompts de Codex los dejamos fuera a propósito: apostar por algo
> que su propio vendor deprecó es deuda el día uno.

> y tambien ten e cuenta que nuestra skills mucha veces dice, pregunta como
> claude, pero se tiene que adaptar s los 3

## Reading

Install the agents kind into Codex (`~/.codex/agents/*.toml` — a real md→TOML
transformation) and OpenCode (`.opencode/agent/*.md` — frontmatter remap:
`mode: subagent`, `model: provider/model-id`, permission map). Generated files
cannot be symlinks, so this phase must design a new ownership marker that keeps
`prune` and `uninstall` safe — that is the expensive part and why it deserves
its own spec. Codex custom prompts are deliberately out of scope: the vendor
deprecated them in favour of skills.
