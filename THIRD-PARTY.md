# Third-party items

This repository ships skills written by other people, so that installing Libretto
Automata gives you a working flow without requiring you to install anything else
first. They are vendored deliberately, not by accident.

## Vendored

| Item | Origin | Version | Licence |
|---|---|---|---|
| `skills/writing-plans/` | [obra/superpowers](https://github.com/obra/superpowers) | 6.2.0 | MIT · Copyright (c) 2025 Jesse Vincent |
| `skills/test-driven-development/` | [obra/superpowers](https://github.com/obra/superpowers) | 6.2.0 | MIT · Copyright (c) 2025 Jesse Vincent |
| `skills/using-git-worktrees/` | [obra/superpowers](https://github.com/obra/superpowers) | 6.2.0 | MIT · Copyright (c) 2025 Jesse Vincent |

Full licence text: [`LICENSE-superpowers`](LICENSE-superpowers).

Copied unmodified. If one of them needs changing for this flow, the change goes in
the Libretto skill that calls it — never into the vendored copy, so the copy stays
comparable with upstream.

## Why vendored and not required

The flow's own skills are thin on purpose: `write-plan` delegates plan structure,
`build-and-check` delegates test discipline and worktree setup. A thin skill whose
delegate is missing is not thin, it is broken.

Being installed on the author's machine is not the same as being installed on
yours. These three ship with the repository so the flow works on a machine that has
nothing else.

## Keeping them current

Pinned at 6.2.0. They will drift from upstream, and that is the accepted cost of
working out of the box.

To update: replace the directory from upstream, bump the version in the table above,
and re-read the Libretto skill that calls it in case what it delegates has moved.

## Not vendored

Referenced when present, never required. The flow works without them.

| Item | Origin | Role |
|---|---|---|
| [ponytail](https://github.com/DietrichGebert/ponytail) | DietrichGebert | decides how much gets built — the ladder, and the list of things never trimmed |
| caveman | juliusbrussee | decides how much gets said |
| `chained-pr` | gentle-ai | chains dependent branches instead of racing them at the trunk |

These install themselves as plugins and namespace their own skills, so vendoring
them would create a second copy of something the user may already have chosen a
version of. `libretto doctor` reports whether they are present.

## Naming

Vendored items keep their upstream names. Plugin-installed skills are namespaced by
their plugin — `superpowers:test-driven-development` — so a vendored
`test-driven-development` in `~/.claude/skills/` does not collide with a plugin copy.
Both can exist. If both do, they are the same guidance at two versions, and the
plugin's is likely newer.
