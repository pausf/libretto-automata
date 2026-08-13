# ask-release-label-at-attacca-end

Tracker: none
Queued: 2026-08-13

## The ask, verbatim

> una mejora seria que el attacca al finalizar la pr etcc te pida poner lo de la
> realease para que no vuelva a petar la PR sabes? que o haga como pregunta de claude
> con 3 opciones

## Reading

After attacca opens the PR and gives its report, ask the bump once with
`AskUserQuestion` — release:patch / release:minor / release:major, the run's own
reading recommended — and apply the answer with `gh pr edit`. This does not bend
attacca's "never label": the reading stays the user's, the agent only types it
(AGENTS.md: "The reading is yours. The typing is not."). Two constraints carried from
the lesson in `.agents/lessons.md` (2026-08-13 · improve-metrics-report): the question
comes after everything else, never as a mid-run stop; and unanswered or headless
(`libretto loop`) runs end exactly as today — unlabeled, with the closing report
predicting the red `require-release-label` check in its own line. `release:major`
keeps its ask-and-wait rule with no exception while 0.x.
