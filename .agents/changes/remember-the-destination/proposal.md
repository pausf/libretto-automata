# Remember the destination

Tracker: none

## What was asked

Verbatim, in the words it arrived in — paraphrasing a request loses the part that is not
understood yet:

> "Vamos a crear una nueva feature ahora mismo el selector de project o global siempre
> esta por default como global, estaria bien guardar estas de como se queda para no tener
> que ir siempre haciendo tab"

## The reading, to be confirmed

The panel always opens with the destination on **global**. Working inside a project means
pressing `tab` every single time it opens. The destination should be remembered as it was
left, so the gesture is not repeated.

## Why it is worth doing

The strip exists because *where did that just install?* is a question best answered before
it is asked — see `.agents/specs/panel/spec.md`. A destination that resets to `global` on
every launch answers it correctly and then makes the user re-state their answer, which is
the same friction one level down.

## What this proposal does not settle

Three things, named here so they are not quietly decided by whoever types first. Two are
product calls and the third touches behaviour that already works:

- **Where the preference is stored.** Beside the project, or in the user's own config.
  These are different claims: *this repository installs into the project* is not the same
  as *I prefer the project*.
- **Per project, or one global preference.** Remembering `project` while sitting in a
  different repository could install somewhere unintended — destructive and silent, the
  worst pair.
- **Whether the subcommands honour it, or only the panel.** `libretto install` with no
  flag goes to global today.

## What is already settled

**An explicit `--project` or `--global` always wins**, whatever is remembered. A
preference that can override a hand-typed flag is a preference that removes from the
wrong place, and `internal/link/own.go` exists precisely because that class of mistake is
unrecoverable.

## Non-goals

- **prompting at startup.** The panel spec already rejects it: an answer given at the top
  of a session is invisible by the time a key is pressed.
- remembering anything else about the session — cursor row, last action, window size.
  One preference, because one is what was asked for.
