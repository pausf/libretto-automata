# 𝄞 Libretto Automata

> The libretto is written first. The automaton performs it.

An 18th-century automaton was a machine that played music by reading a score. A
human wrote the notes; the machine executed them. It never improvised — it did
exactly what the paper said. If the performance was wrong, the paper was wrong.

That is how this works. You write the spec. The agent performs it. When the
output is bad, you fix the libretto, not the automaton.

## What this is

My own Claude Code setup — skills, agents, commands — kept in git instead of
living loose inside `~/.claude` where a sync can overwrite it.

It installs by **symlink**, not by copy. Edit here, it's live immediately.
`git pull` updates everything.

## Install

```bash
git clone <this repo> ~/gitrepos/libretto-automata
cd ~/gitrepos/libretto-automata
./install.sh
```

## Commands

| | |
|---|---|
| `./install.sh` | symlink every item into `~/.claude` |
| `./install.sh status` | what's linked, what's missing |
| `./install.sh doctor` | find broken links |

Links are made **per item**, never per directory — so this coexists with
anything else installed into the same folders. Anything already there that
isn't a symlink is left untouched.

## Layout

```
skills/      one directory per skill
agents/      one .md per agent
commands/    one .md per slash command
```

Drop something in, run `./install.sh`, it's live.

## Not managed here

`CLAUDE.md` and `settings.json` — those get rewritten by other tooling, so
they stay hand-managed.
