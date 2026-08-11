# notify-users-of-new-updates

Tracker: none
Queued: 2026-08-11

## The ask, verbatim

> alguna veces me olvido de install desde libreto de nuevos command etc.. que vinieron
> en la nueva update es decir me falta fichero por linkar por que son nuevos, que
> podemos hacer para que lo usuario sepan que tiene una actualizacion nueva ?

And later, when the question came back after `distribution` shipped:

> cuando hay una nueva version del cli como se entera el usuario?

## What already landed, and what it did not answer

`ship-payload-as-release-asset` closed **half** of this, and it is worth being exact about
which half:

- **The mechanical half is done.** `update` relinks after installing, so a version that adds
  a skill no longer leaves it unlinked. That was the first sentence of the original ask.
- **The notifying half is not.** There are exactly two places a user finds out a newer
  version exists: the panel's notice row, and `libretto doctor`.

**Both require running `libretto`, and after installing, most people never do.** The payload
is used *inside Claude Code* — you install once and then work with the skills for months
without typing the command again. The panel also needs a TTY, so `libretto status` in a
script says nothing about a new version at all.

So the notice lives where the user is not. That is the whole of what is left.

**This was briefly marked as answered and it was wrong.** Recorded because "the relink covers
it" is the mistake a reader would make again: the relink covers the *consequence* of a new
version, not the *knowing*.

## The scope, as discussed

Two candidates, and the second crosses a boundary this project guards.

| | What it does | Cost |
|---|---|---|
| **Notify on any subcommand** | `status`, `install`, `doctor` print the notice row when the cache knows of a newer version | small — the cache already exists, and it is one line at the end |
| **Let the payload say it** | a skill, or the `libretto-status` command, mentions the version when the flow starts | medium, **and it puts the CLI inside the payload** — a boundary the project keeps |

The first is where to start. It is honest about the real limit: **if the user never runs the
CLI, no amount of CLI notices reaches them.** What can be guaranteed is that *whenever they
run it, for any reason*, they find out — and today that only happens if they chose the panel.

The second is the only thing that reaches somebody who genuinely never runs the command, and
its cost is the one to argue about: `payload`'s spec keeps skills free of this repository's
own tooling, precisely so an installed skill works in a project that has never heard of
`libretto`. A skill that announced CLI versions would break that.

## What this will have to settle

- **Where the line goes**, and whether it is silent when the cache has no answer — the panel's
  rule is silence, and a subcommand printing "could not check" on every run would be noise.
- **Whether a script sees it.** Piped output carries no escape codes by an existing promise; a
  notice on `status` would be a new line in output somebody may be parsing.
- **The TTL question.** A day is right for a panel launch. Somebody running `install` right
  after a release still waits up to a day to be told, and whether that matters is a decision.
- **Whether `payload` should be reopened at all**, or whether the answer is that the CLI cannot
  reach a user who does not run it, stated as a limit rather than solved.
