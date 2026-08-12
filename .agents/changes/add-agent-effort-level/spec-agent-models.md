# Agent effort level — delta

Targets: agent-models

Governs: internal/agentmodel/**

An agent already declares *which* model it runs on. This adds *how hard that model
thinks*, in the same file, on its own line:

```
---
name: review-lens-design
model: opus
effort: xhigh
---
```

The lever is the same one the capability was built for, pointed the other way. `model:`
buys cheaper tokens for work that does not need capability; `effort:` buys fewer tokens
from a model you are keeping, and buys more of them for the one place it pays.

## Outcomes

- **Reading:** for every agent in the directory, the user sees its effort alongside its
  model — a declared level, or the session's, when the file declares none.
- **Writing:** the user marks agents, picks one level, and that level is written into
  every marked agent's frontmatter in one act. Same all-or-nothing act as the model.
- **The rest of the file is untouched**, byte for byte. The same promise, for the same
  reason: an agent file is a prompt.
- **`effort:` and `model:` are independent keys.** Either can be set without disturbing
  the other. Two keys, two acts, because an agent staying on Opus while its effort drops
  is the whole point of the feature.
- **A level the catalogue does not know is refused**, and nothing is written.
- **A level the agent's model cannot run is refused**, and nothing is written. This is
  the criterion that makes the capability honest rather than decorative — see below.

### The catalogue of levels

| Value | Means |
|---|---|
| `low` | short, scoped, latency-sensitive work that is not intelligence-sensitive |
| `medium` | cost-sensitive work that can trade off some intelligence |
| `high` | the balance point, and the default on every model that supports effort |
| `xhigh` | deeper reasoning at higher spend |
| `max` | the deepest; prone to overthinking, worth measuring before adopting |
| *default* | no `effort:` key at all: the agent runs at whatever the session runs at |

Read off <https://code.claude.com/docs/en/model-config> · *Adjust effort level* and
<https://code.claude.com/docs/en/sub-agents> · *Supported frontmatter fields*, 2026-08-12.
The same `Resolved` date the model catalogue already publishes covers this table, because
it decays for the same reason and there is no second date worth maintaining.

`ultracode` is **not** in it. It is a Claude Code session mode that sends `xhigh` and
turns on workflow orchestration; it is not a level and no frontmatter accepts it.

### Which model has which

| Catalogue model | Levels it supports |
|---|---|
| `opus` → Opus 5 | all five |
| `sonnet` → Sonnet 5 | all five |
| `haiku` → Haiku 4.5 | **none. Haiku does not support effort at all** |
| *default* — the session's model | **all five, because the binary does not know what the session runs** |

Two of those four rows are decisions rather than facts, and both are stated under prior
decisions below.

## Scope boundaries

**In:** reading and writing the `effort:` key of agent files in a directory, the
catalogue of levels, which catalogue model supports which, and refusing everything else.

**Out, named so it cannot arrive quietly:**

- **Fusing effort into the model catalogue.** `opus/xhigh` as a single selectable value
  would turn four entries into twenty and make "keep the model, drop the effort"
  impossible to express. Two keys stay two keys.
- **`ultracode`.** Not a level. See above.
- **Effort for skills.** Claude Code does accept `effort:` in skill frontmatter, and this
  capability's existing boundary already excludes skills for the model. It stays
  excluded, for the same reason, in the same words: only one kind of item is the subject
  here. A skill's effort is a different feature and would need its own catalogue of
  what a skill even is to this tool.
- **Asking Anthropic which levels this account may use.** Organisation effort caps exist
  and are invisible from here. No API call, no token — the same standing refusal.
- **A different effort per install target.** Falls to the same symlink invariant that
  settled it for the model. One effort per agent, per machine.
- **Rewriting existing agent files.** Nothing gains an `effort:` key until asked.

## Constraints

- **Go 1.26.5, standard library only.** The frontmatter reader already exists and already
  refuses general YAML. Adding a second key must reuse it rather than grow a parser.
- **`scripts/check-payload` requires `name:` and `description:` and ignores every other
  key** — verified at `scripts/check-payload:62-67` when `model:` landed. `effort:` needs
  no gate change for the same reason, and the same verification applies.
- **The frontmatter functions are currently written around one hardcoded key**
  (`internal/agentmodel/frontmatter.go`, `const key = "model:"`). A second key means that
  constant becomes a parameter. `ReadModel` and `SetModel` keep their names and their
  signatures — they are the contract three surfaces call — and gain `ReadEffort` and
  `SetEffort` as siblings over one shared keyed implementation. Two copies of the
  byte-for-byte logic is two places for the promise to break.
- **This package still never learns what an install target is.** No `internal/target`
  import, testable against a bare `t.TempDir()`.
- **`Apply`'s all-or-nothing guarantee extends unchanged.** The whole set is validated
  before the first file is opened, and the pairing check is part of that validation, not
  a per-file decision taken while writing.

## Prior decisions

- **The default value is an absent key**, exactly as it is for the model, and it means
  *the session's effort* — not `high`. The docs say `high` is the model's default, and
  that is a fact about Claude Code, not about this file: an agent with no `effort:`
  inherits the session, and the session may have been launched at `low`. Printing
  `(session)` for it reuses the word the model column already prints, because it is the
  same state and a second spelling would read as a different one.

- **Effort on Haiku is refused, not written and ignored.** *(assumed under
  `/libretto-attacca` — nobody was asked)* Haiku appears in no row of the docs' effort
  table, so the key would be inert. Writing it anyway leaves a line in a prompt file
  that claims a setting the model has no concept of, and the catalogue's own comment
  already names that failure mode: *"the kind of confident wrong answer that gets
  believed"*. **What changes if this is wrong:** the check becomes a warning instead of a
  refusal — one branch, in `ApplyEffort`, and no file format moves.

- **Effort on an agent with no `model:` is allowed.** *(assumed — same run)* The session
  model is unknowable from here, so refusing would be a guess in the other direction, and
  the capability's standing rule is that this binary never claims to know what the
  account or the session is running. It says what it wrote; the host resolves it. **What
  changes if this is wrong:** nothing in the format — only the one condition that decides
  whether `Default` counts as effort-capable.

- **Moving an agent to a model with no effort clears its `effort:` key, and says so.**
  *(assumed — same run)* `models set haiku` on an agent declaring `effort: xhigh` has
  three possible answers: refuse the model change, keep the dead key, or drop it. Refusing
  makes the cheap model the awkward one to reach, which is backwards for a capability
  built to make Haiku reachable. Keeping it leaves the lie the decision above exists to
  prevent. So it drops, in the same act, reported on the row. **What changes if this is
  wrong:** `Apply` stops clearing and starts refusing — the strictness moves, the format
  does not.

- **An unsupported level is not silently downgraded here.** Claude Code itself falls back
  to the highest supported level at or below the one asked for, so `xhigh` on Opus 4.6
  runs as `high`. This tool does not reproduce that: it writes aliases, not versions, and
  every alias in its catalogue either supports all five or none. Reimplementing the
  host's fallback would mean pinning versions, which the capability already refused.

## Task breakdown

1. `internal/agentmodel`: parametrise the frontmatter reader and writer by key, keeping
   `ReadModel`/`SetModel` byte-identical in behaviour.
2. `internal/agentmodel`: `ReadEffort` and `SetEffort` over that shared implementation.
3. `internal/agentmodel`: the effort catalogue — the five levels, their labels, the
   validation, and which catalogue model supports effort at all.
4. `internal/agentmodel`: `ApplyEffort` — one level onto a set as one act, validating the
   level and every target agent's model before writing any of it.
5. `internal/agentmodel`: `Agent` carries `Effort`, and `Apply` clears a stale effort key
   when the model it writes supports none.

## Verification criteria

Reading:

- a declared effort is read back
  Proof: internal/agentmodel/effort_test.go TestReadEffortReturnsTheDeclaredEffort
- a file with no `effort:` reads as default, not as an error
  Proof: internal/agentmodel/effort_test.go TestReadEffortReportsDefaultWhenTheKeyIsAbsent
- an `effort:` in the body rather than the frontmatter is not read as the effort
  Proof: internal/agentmodel/effort_test.go TestReadEffortIgnoresTheBody

Writing, and the byte-for-byte promise:

- inserting the key leaves every other line identical, `model:` included
  Proof: internal/agentmodel/effort_test.go TestSetEffortInsertsWithoutDisturbingTheFile
- replacing an existing level changes that line and no other
  Proof: internal/agentmodel/effort_test.go TestSetEffortReplacesInPlace
- choosing default removes the key rather than writing a word
  Proof: internal/agentmodel/effort_test.go TestSetEffortDefaultRemovesTheKey
- setting the level an agent already has rewrites nothing
  Proof: internal/agentmodel/effort_test.go TestSetEffortIsIdempotent
- **the model survives an effort write and the effort survives a model write** — two
  independent keys, proven independent
  Proof: internal/agentmodel/effort_test.go TestTheTwoKeysDoNotDisturbEachOther

The catalogue:

- the five levels are listed, weakest first, and `ultracode` is not among them
  Proof: internal/agentmodel/effort_test.go TestEffortCatalogueListsTheFiveLevels
- an unknown level is refused
  Proof: internal/agentmodel/effort_test.go TestUnknownEffortIsRefused
- **`haiku` supports no effort; `opus`, `sonnet` and the session default support all five**
  Proof: internal/agentmodel/effort_test.go TestWhichModelsSupportEffort

Applying to a set:

- one level reaches every agent in the set
  Proof: internal/agentmodel/effort_test.go TestApplyEffortReachesEveryAgentInTheSet
- **an agent on Haiku means nothing is written at all** — not even the agents before it
  Proof: internal/agentmodel/effort_test.go TestApplyEffortWritesNothingWhenAnyAgentCannotRunIt
- an agent with no declared model accepts a level
  Proof: internal/agentmodel/effort_test.go TestApplyEffortAllowsAnAgentOnTheSessionModel
- an unknown level is refused before any file is opened
  Proof: internal/agentmodel/effort_test.go TestApplyEffortRefusesAnUnknownLevel
- listing reports each agent's effort beside its model
  Proof: internal/agentmodel/effort_test.go TestAgentsReportsEachCurrentEffort
- **moving an agent to Haiku clears a declared effort rather than leaving it dead**
  Proof: internal/agentmodel/effort_test.go TestApplyModelClearsEffortWhenTheModelSupportsNone
- **and the clearing is narrow: a model that does support effort leaves the level alone**
  Proof: internal/agentmodel/effort_test.go TestApplyModelKeepsEffortWhenTheModelSupportsIt
