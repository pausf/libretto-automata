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

## The three open questions, answered

Asked with the alternatives and the risk of each stated. These are the user's answers, and
they are recorded here rather than left in the conversation — an answer that lives only in
a transcript gets asked again next session, and does not come out the same.

- **Where the preference is stored: in the global `.claude`.** `~/.claude`, honouring
  `CLAUDE_HOME`. One consequence worth naming because it was not the reason for the
  choice but pays for it anyway: `CLAUDE_HOME` is what already makes the suite safe to run
  twice, so the preference is test-isolated without a second mechanism.
- **One preference, not per project.** The alternative — keyed by working directory —
  was offered together with the case against a single value: the panel can open on
  `project` while sitting in a repository that was never meant to have it. Accepted
  knowingly, and the next line is why the cost is small.
- **Only the panel honours it.** `libretto install` with no flag still goes to global.
  A command typed into a terminal does not change meaning because of state left by an
  earlier session. This also bounds the risk above: the remembered value decides which
  side the panel *opens* on, and the destination strip is visible before any key that
  acts — installing stays an explicit gesture against a destination you can see.

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
