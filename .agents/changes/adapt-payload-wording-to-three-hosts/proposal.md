# adapt-payload-wording-to-three-hosts

Tracker: none
Queued: 2026-08-14

## The ask, verbatim

> y tambien ten e cuenta que nuestra skills mucha veces dice, pregunta como
> claude, pero se tiene que adaptar s los 3

## Where it came from

Captured out of `add-opencode-command-target` on 2026-08-14, which carried this
sentence in its proposal and could not honour it. That change is target plumbing —
two lines saying OpenCode accepts the commands kind. This is a content pass over
every skill and command in the payload. Landing them together would mix a target
change with a payload rewrite and neither would be reviewable.

## Reading

The payload names Claude Code's vocabulary throughout, and now installs into three
hosts. Known divergences, read off `sst/opencode` on 2026-08-14:

- **`Skill(skill="find-work")`.** OpenCode has a skill tool, and its parameter is
  `name`, not `skill`. The instruction is prose a model reads rather than a literal
  call, so it degrades rather than breaking — but it is Claude's spelling.
- **`AskUserQuestion`.** OpenCode has `question.ts`. Named differently; the skills
  already carry the escape hatch "or in conversation where the native prompt does not
  exist", which is the pattern the rest should follow.
- **Codex** accepts skills only, so its exposure is the `SKILL.md` bodies rather than
  the commands.

## The question this has to answer first

Whether to name each host's tool, or to describe the *capability* and let the host
bind it — "load the `find-work` skill", "ask the user with the native prompt if there
is one". The second is one text for three hosts and is probably right, but it trades
away the precision that makes a `Skill(...)` line unambiguous today.

Not decided. That is the decision this change exists to make, and it is the user's.
