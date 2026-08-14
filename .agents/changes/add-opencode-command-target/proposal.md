# add-opencode-command-target

Tracker: none

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

## Reading, corrected at pickup — 2026-08-14

**Two of the three guesses above are wrong, and both were checked against
`sst/opencode` rather than the docs page.** `packages/opencode/src/config/command.ts`
loads markdown commands with:

```ts
Glob.scan("{command,commands}/**/*.md", { cwd: dir, absolute: true, dot: true, symlink: true })
```

- **The directory is `commands` as well as `command`.** Both names are globbed, so
  the plural — which is what `dirUnderRoot` already produces for every kind — needs
  no special case. The docs page names only the plural; the queued reading named only
  the singular; the source accepts both.
- **`symlink: true`.** Symlinks are followed, so the existing install mechanism
  works unchanged. This is the fact the whole change rests on.
- **No transform is needed.** `ConfigCommandV1.Info` is an Effect `Schema.Struct`,
  which ignores excess properties — proven by the loader passing `name` into a struct
  that does not declare it. An unknown key like `allowed-tools` is dropped silently,
  not rejected. And this payload's commands carry `description:` and nothing else, so
  there is no key to map even if mapping were possible.

What survives from the queued reading: OpenCode does not read `.claude/commands`, so
a real install step is needed, and it depends on the target plumbing that landed in
`add-multi-tool-skill-targets`.

**The second half of the ask is not in this change.** Adapting the payload's
Claude-specific instruction wording to three hosts is a payload-wide content change
that touches every skill, not a target plumbing change. Captured separately — see
`.agents/changes/adapt-payload-wording-to-three-hosts/`.
