# support-codex-opencode-in-loop

Tracker: none
Queued: 2026-08-21

## The ask, verbatim

> soporte codex/opencode — instalás a tres tools pero el loop solo maneja `claude`.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

`libretto loop` hands a fixed prompt to `claude` (`cmd/libretto/loop.go`) while install
already targets codex and opencode. Let the loop drive the other two tools it installs
to.
