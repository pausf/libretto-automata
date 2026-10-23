# trim-skill-reading-cost

Tracker: none
Queued: 2026-10-23

## The ask, verbatim

> Pendiente aparte: la "dieta" de los skills (se leyeron ~3.000 líneas de skills para 15 líneas de código). ¿La encolo como idea?

> si

## Reading

A trial run of `/libretto-attacca` on a throwaway project read about 3,000 lines of skill
text for a change of about 15 lines. Goal: the reading cost of a small change scales with
the change, not with the payload.

Measured on that run:

- record-work is ~70% push, forge and bump, all irrelevant without a remote.
- write-spec is mostly fan-out, which a single capability never uses.
- find-work is ~60% Jira and queue.
- write-tasks delegates to writing-plans, whose format contradicts the flow: a
  docs/superpowers path, per-step boxes and an execution question.

Small contradictions found on the same run:

- present-work's "show then commit" against commit-as-you-go.
- record-work announces "three things" and lists four.
- write-spec's EARS example is a joined criterion that `--anchors` warns on.
- The reviewers run async, although the skills say foreground.
