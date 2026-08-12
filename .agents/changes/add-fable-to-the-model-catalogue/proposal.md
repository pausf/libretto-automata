# add-fable-to-the-model-catalogue

Tracker: none
Source: what the user said, under `/libretto-attacca`

## The ask, in the words it was asked in

> en lo modelos y esfuerzo añade Fable de claude

## The reading

`libretto models` and `libretto models effort` offer a catalogue of four values — the
session default, `haiku`, `sonnet`, `opus`. Claude's `fable` alias is missing from it,
so no agent can be put on Fable through this tool and the selector cannot offer it.

The plumbing for Fable is already half here and was built without the catalogue entry:
`internal/agentmodel/provider.go` already reads `Fable 5` out of a pinned model id
(`pinPattern` matches `fable`) and already knows Fable 5 runs all five effort levels
(`effortByVersion`). What is missing is the one thing that makes the alias reachable —
a `catalogue` row — and the alias in the provider tables that resolve it.

So this is one entry in the model catalogue, not a new key and not a new capability.
The effort half of the ask needs no new level: Fable 5 runs the same five, and it will
be offered them the moment the alias resolves.
