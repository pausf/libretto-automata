# add-opencode-command-target

Tracker: none
Queued: 2026-08-14

## The ask, verbatim

> Fase 2 — commands para OpenCode: conversión casi trivial de md a
> .opencode/command/.

> y tambien ten e cuenta que nuestra skills mucha veces dice, pregunta como
> claude, pero se tiene que adaptar s los 3

## Reading

Install the commands kind into OpenCode's `.opencode/command/` (project) or
`~/.config/opencode/command/` (global). The format is near-identical to Claude
commands — `$ARGUMENTS` and `$1..$n` work as-is — but OpenCode does NOT read
`.claude/commands` directly (verified in sst/opencode source), so this needs a
real install step. Frontmatter keys that OpenCode lacks (`allowed-tools`) need
dropping or mapping, which may make this a light transform rather than a pure
symlink. Depends on the target plumbing from add-multi-tool-skill-targets.
