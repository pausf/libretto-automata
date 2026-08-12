# add-agent-effort-level

Tracker: none

## What was asked

> me han pedido una nueva feaute que es pner el nivel de esfuerzo del modelo que se
> elija, es decir en el cli que puedas eleir por ejemplo OPUS esfuerzo zhight, hight ,
> medium etcc,, bueno investifa cuales son lo smodelos que se pueden selecionar

Verbatim, in the words it was asked in. Read as: the model selector already picks *which*
model an agent runs on; it should also pick *how hard that model thinks*.

## What the investigation found

The ask contained a research question, and it is answered before any spec is written
because the answer decides the shape of the feature.

Claude Code exposes effort in exactly the place this repository already writes to — a
subagent file's frontmatter, one key beside `model:`:

```
effort: xhigh
```

Source: <https://code.claude.com/docs/en/sub-agents> · *Supported frontmatter fields*,
and <https://code.claude.com/docs/en/model-config> · *Adjust effort level*, read
2026-08-12.

### The levels, and which model has which

| Model | Levels |
|---|---|
| Fable 5 | `low` `medium` `high` `xhigh` `max` |
| Opus 5, Sonnet 5, Opus 4.8, Opus 4.7 | `low` `medium` `high` `xhigh` `max` |
| Opus 4.6, Sonnet 4.6 | `low` `medium` `high` `max` — no `xhigh` |
| anything else, **Haiku included** | no effort at all |

Three facts that matter more than the table:

- **`high` is the default** on every model that supports effort. Absence of the key is
  not "no effort", it is `high` — which is the same relationship `model:` already has
  with the session default, and a different one from what "unset means minimum" would
  suggest.
- **An unsupported level does not fail.** Claude Code falls back to the highest supported
  level at or below the one asked for, so `xhigh` on Opus 4.6 runs as `high`. A wrong
  choice degrades; it does not break.
- **`max` is session-only in settings**, but frontmatter is not settings. The
  `effortLevel` settings key rejects `max`; the subagent `effort` field lists it.
- **`ultracode` is not an effort level.** It is a Claude Code session mode. It has no
  place in a per-agent catalogue.

### Where the repository stands today

`internal/agentmodel/` owns the frontmatter line and the catalogue. Its catalogue is
`default · haiku · sonnet · opus`, and **Haiku supports no effort** — so this feature has
a hole in it by construction on one of the four entries, and that hole is the design
question phase 2 has to answer rather than paper over.

Three surfaces render the model today and would each grow an effort dimension:
`libretto models`, `libretto models set`, and the panel's selector screen
(`internal/ui/models.go`).
