# Agent Models

Governs: internal/agentmodel/**

Choose which model each payload agent runs on, and how hard that model thinks, so the
work that is pattern-matching over prose stops being billed at the rate of the work
that is not.

The lever is **not fewer agents and not fewer tokens**. It is cheaper tokens for the
lenses that read prose and answer in prose.

Two keys, and they are independent. `model:` chooses the tier; `effort:` chooses the
depth inside it. Keeping an agent on Opus while its effort drops is the case the second
key exists for, so neither is expressible only through the other.

## Outcomes

An agent's model and effort are declared in its own frontmatter, and the user sets each
without opening a file.

```
---
name: review-design
description: …
tools: Read, Grep, Glob, Skill
model: haiku
---
```

```
---
name: review-security
model: opus
effort: xhigh
---
```

- **Reading:** for every agent the repo ships, the user can see which model it runs on
  today — a declared one, or the session's, when the file declares none.
- **Writing:** the user marks agents — one, several, or all — picks one model, and that
  model is written into every marked agent's frontmatter in one act.
- **The rest of the file is untouched.** Byte for byte, apart from the one line that
  changed. An agent file is a prompt; a writer that reflows it changes the prompt.
- **A model the catalogue does not know is refused**, and nothing is written. Not one
  agent, not the ones before the bad name in the list.
- **The catalogue names what a Claude subscription offers**, and says which plan each
  one needs, because the binary cannot ask Anthropic what the user is paying for.

| Value | Resolves to | Means |
|---|---|---|
| `opus` | Opus 5 | the most capable. Max plans; metered on Pro |
| `sonnet` | Sonnet 5 | the default working model |
| `haiku` | Haiku 4.5 | the cheap one — what this capability exists to make reachable |
| *default* | — | no `model:` key at all: the agent runs on whatever the session runs on |

**The value written is the alias; the version is what it means today.** `opus` keeps
meaning "the Opus tier" after the model behind it is replaced, which is why the
frontmatter carries the alias — an agent file written today should not need editing
the day a new Opus ships.

But `opus` alone does not answer the question somebody opens the selector to ask, so
the catalogue carries the version too, **with the date it was last checked**. That
column decays. It is stated rather than hidden because this repository has already
shipped the other kind: a test-count badge that read 117 against an actual 221,
because nobody recomputed it and nothing said when it was written.

Choosing *default* **removes** the key rather than writing a word meaning "no choice".
An absent key is already the language's way of saying it, and two spellings of one
state is a difference somebody will eventually treat as meaningful. The same holds for
`effort:`, and there is one `Default` constant serving both — a second name for the
empty string would be the same difference in the code.

### Effort — the depth inside the tier

| Value | Means |
|---|---|
| `low` | short, scoped work that is not intelligence-sensitive |
| `medium` | cost-sensitive work that can trade off some intelligence |
| `high` | the balance point, and the host's own default |
| `xhigh` | deeper reasoning at higher spend |
| `max` | the deepest; prone to overthinking. Measure before adopting |
| *default* | no `effort:` key at all: the agent runs at whatever the session runs at |

**An absent `effort:` means the session's effort, not `high`.** `high` is the host's
default for a model; the file's silence is inheritance, and the session may have been
launched at `low`. The listing prints `(session)` for it — the same word the model
column prints, because it is the same state.

`ultracode` is **not** a level. It is a Claude Code session mode that sends `xhigh` and
turns on workflow orchestration, and no frontmatter accepts it.

**Which model has which:**

| Value | Effort |
|---|---|
| `opus`, `sonnet` | all five |
| `haiku` | **none. Haiku supports no effort at all** |
| *default* | all five, because the session's model is unknowable from here |

- **A level on a model that has none is refused**, and nothing is written. Writing it
  would leave a line in a prompt file claiming a setting the model has no concept of,
  which is the confident wrong answer this catalogue's whole posture is against.
- **Moving an agent onto a model with no effort clears its `effort:` key**, and the
  callers say so on the row it happened to. Refusing the model change instead would
  make the cheap model the awkward one to reach, which is backwards for a capability
  built to make Haiku reachable; leaving the key would be the lie above. The clearing
  is narrow: a model that does support effort leaves the level alone, because changing
  the tier is not a request to change the depth.
- **The levels are answerable, not only refusable.** `EffortsFor` names what a model can
  run so a caller can offer a choice instead of discovering the refusal after one was
  made. The panel had only the apply-time error to go on, and the cost was a menu of five
  levels over two Haiku rows. It returns a slice rather than a bool because the host's own
  table already lists a model with four of the five, so a bool would have to become this
  the day one enters the catalogue.
- **An unsupported level is never silently downgraded here.** The host itself falls
  back to the highest supported level at or below the one asked for — `xhigh` runs as
  `high` on Opus 4.6 — and this tool does not reproduce that. It writes aliases, not
  versions, and every alias it knows either supports all five or none. Reproducing the
  fallback would mean pinning versions, which is refused above.

The `Resolved` date covers this table too. It decays for the same reason, and a second
date is a second thing to forget.

## Scope boundaries

**In:** reading and writing the `model:` and `effort:` keys of the agent files in a
directory, the catalogue of legal values for each, which model supports effort at all,
and refusing everything else.

**Out, named so it cannot arrive quietly:**

- **Fusing effort into the model catalogue.** `opus/xhigh` as one selectable value
  would turn four entries into twenty and make "keep the model, drop the effort"
  inexpressible. Two keys stay two keys.
- **`ultracode`.** Not a level, and not something a file can declare.
- **Organisation effort caps.** They exist and are invisible from here. Same standing
  refusal as the plan tiers: no API call, no token.
- **A different model per install target.** Settled below under prior decisions: it
  cannot be done without breaking the symlink invariant, and it is not worth that.
- **Models or effort for skills or commands.** Only agents are the subject here. A host
  does accept `effort:` in skill frontmatter, and it stays out for the same reason the
  model does: a skill's settings are a different feature and would need this tool to
  learn what a skill is to it.
- **Asking Anthropic which models this account may use.** No API call, no token, no
  network. The catalogue is static. It comes back the day the CLI has an authenticated
  way to ask that does not involve this tool touching a credential — which
  `AGENTS.md` forbids outright.
- **Per-invocation model choice.** The frontmatter is the whole mechanism. A caller
  that wants a different model for one run is a different feature.
- **Migrating existing agent files.** Nothing is rewritten until the user asks.

## Constraints

- **Go 1.26.5, standard library only.** Frontmatter here is a fenced block of
  `key: value` lines at the top of the file. It is not general YAML and must not become
  a reason to add a YAML dependency — `AGENTS.md` puts a new dependency behind an ask,
  and this does not clear the bar.
- **`scripts/check-payload` requires `name:` and `description:` and ignores every other
  key** (`scripts/check-payload:62-67`). `model:` and `effort:` therefore need no gate
  change. Verified, not assumed.
- **One reader and one writer serve both keys**, parametrised by key rather than copied.
  The byte-for-byte promise below is the whole value of that file, and two copies of it
  is two places for it to break — only one of which anybody would think to re-read.
  `ReadModel`/`SetModel` and `ReadEffort`/`SetEffort` are four names over one
  implementation.
- **`Apply` and `ApplyEffort` share the resolution and the all-or-nothing refusal.** The
  guarantee is the same guarantee in both, and a second copy is a copy a caller can
  reach the weaker version of.
- **The writer must refuse a file with no frontmatter**, rather than inventing one. A
  file that does not open with `---` on line 1 is not an agent file by this repo's own
  check, and writing a block into it would manufacture an agent out of a document.
- This package **never learns what an install target is**. It is handed a directory
  and works on every `*.md` in it, whoever created them. Which directory is the CLI's
  problem, and that is what keeps the layering true while the reach is wide: no
  `internal/target` import, and the whole package testable against a bare `t.TempDir()`.
- **A stale link is survivable; a stray file is not.** An entry that cannot be opened
  at all — a symlink whose destination was renamed or deleted — is skipped and named
  in a second return value. Renaming one agent leaves a dangling link in every target
  that had the old name, and taking eleven readable agents down because a twelfth is
  broken is a listing that punishes the ordinary case. A file that *is* present and is
  not an agent stays an error: `Apply`'s all-or-nothing guarantee rests on it.
- **A directory that does not exist reports no agents rather than an error.** A target
  that has never had one installed is a state, and making every caller special-case
  `os.IsNotExist` to render an empty list is how that state becomes a crash.
- **A symlinked agent file is written through to its destination.** Ordinary file
  behaviour rather than anything this package does — but the callers promise it to the
  user, so it is pinned by a test rather than left to be true by accident.

## Prior decisions

- **The subject is a directory, not the repository.** The first version of this
  capability read `<repo>/agents` and nothing else. On the machine that reported it
  that meant editing seven files installed nowhere while the user's own 22 agents were
  invisible — *"the agents available right now"* had been read as *"the agents this
  payload ships"*. **The contract was written from the same misreading as the code,
  which is why no gate caught it.**

  What survives of the original decision is below: the CLI still decides which
  directory. What fell is the assumption that the answer is always this repository's.

- **Superseded — kept because the reasoning still holds for the shared case.** Asked
  whether the choice
  should be global or per-project, the answer was both — and the code says both resolve
  to the same file. `~/.claude/agents/x.md` and `<cwd>/.claude/agents/x.md` are symlinks
  to the same repo file (`internal/link/own.go:23`), and a write through a symlink
  writes its destination. A genuinely per-target model needs a **real file** in the
  target, which `internal/link/state.go:110` classifies as `Conflict` — never touched,
  always reported, `install` exits non-zero. The ceiling: **one model per agent, for
  every project on the machine.** What would lift it is teaching `Owned()` a third
  concept, "a real file we generated" — and that predicate is the one `own.go:9` calls
  the most consequential in the program. Not for this.

- **`review-lens` is split into four agent files.** Its one-file design rested on a
  stated premise — the four lenses "differ in exactly one thing"
  (`skills/review-project/SKILL.md:277`). A per-lens model is a second thing, so the
  premise no longer holds. Chosen over writing models into the skill's launch table,
  because that would make the selector write two formats and parse a markdown table
  that any reformat breaks. The cost, accepted knowingly: four near-identical bodies
  that can drift. The reasoning is recorded in the review-project spec.

- **The default value is an absent key, not the word `inherit`.** See outcomes. It holds
  for both keys, and there is one constant behind it.

- **Effort on Haiku is refused rather than written and ignored**, and effort on an agent
  declaring no model is allowed. *(Both assumed under `/libretto-attacca`, 2026-08-12 —
  nobody was asked.)* The first because an inert key is the confident wrong answer this
  capability is built against; the second because the session's model is unknowable from
  here and refusing would be the guess in the other direction. **What changes if either
  is wrong:** one condition each — the refusal becomes a warning, or `Default` stops
  counting as effort-capable. No file format moves either way.

- **Moving to a model with no effort clears the level rather than refusing the move.**
  *(Assumed, same run.)* Reasoning under outcomes. **What changes if this is wrong:** the
  strictness moves, the format does not.

## Task breakdown

1. `internal/agentmodel`: read the declared model from one agent file, and report a
   file that declares none as *default* rather than as an error.
2. `internal/agentmodel`: write or remove the key, leaving every other byte alone, and
   refuse a file with no frontmatter.
3. `internal/agentmodel`: the catalogue — the legal values, their labels, and the
   validation that refuses everything else.
4. `internal/agentmodel`: apply one model to a set of agents as one act, validating the
   whole set before writing any of it.
5. `internal/agentmodel`: the same four for `effort:` — read, write, its catalogue and
   which model supports it, and apply one level to a set as one act with the pairing
   check inside that validation.
6. `internal/agentmodel`: `Agent` carries both keys, and applying a model clears an
   effort the new model cannot run.

## Verification criteria

Reading:

- a declared model is read back
  Proof: internal/agentmodel/frontmatter_test.go TestReadModelReturnsTheDeclaredModel
- a file with no `model:` reads as default, not as an error
  Proof: internal/agentmodel/frontmatter_test.go TestReadModelReportsDefaultWhenTheKeyIsAbsent
- a `model:` appearing in the body rather than the frontmatter is not read as the model
  Proof: internal/agentmodel/frontmatter_test.go TestReadModelIgnoresTheBody

Writing, and the byte-for-byte promise:

- inserting a key leaves every other line identical
  Proof: internal/agentmodel/frontmatter_test.go TestSetModelInsertsWithoutDisturbingTheFile
- replacing an existing key changes that line and no other
  Proof: internal/agentmodel/frontmatter_test.go TestSetModelReplacesInPlace
- choosing default removes the key rather than writing a word
  Proof: internal/agentmodel/frontmatter_test.go TestSetModelDefaultRemovesTheKey
- a file with no frontmatter is refused and left unchanged
  Proof: internal/agentmodel/frontmatter_test.go TestSetModelRefusesAFileWithoutFrontmatter
- setting the model an agent already has rewrites nothing
  Proof: internal/agentmodel/frontmatter_test.go TestSetModelIsIdempotent

The catalogue, and refusing what is not in it:

- the catalogue lists exactly the subscription models plus default
  Proof: internal/agentmodel/catalogue_test.go TestCatalogueListsTheSubscriptionModels
- **every real model names the version its alias resolves to** — an alias alone does
  not say what it means
  Proof: internal/agentmodel/catalogue_test.go TestEveryRealModelNamesItsVersion
- the date those versions were checked is stated, so the staleness is visible
  Proof: internal/agentmodel/catalogue_test.go TestTheResolvedDateIsStated
- an unknown model name is refused
  Proof: internal/agentmodel/catalogue_test.go TestUnknownModelIsRefused

Applying to a set:

- a directory of agent files is listed whatever its path
  Proof: internal/agentmodel/apply_test.go TestAgentsListsEveryAgentSorted
- current models are read per agent
  Proof: internal/agentmodel/apply_test.go TestAgentsReportsEachCurrentModel
- **a directory that does not exist reports no agents rather than an error**
  Proof: internal/agentmodel/apply_test.go TestAgentsOnAMissingDirectoryIsEmptyNotAnError
- **a stale link is skipped and named, not fatal to the whole listing**
  Proof: internal/agentmodel/apply_test.go TestAgentsSkipsAStaleLinkAndNamesIt
- a present file that is not an agent is still an error
  Proof: internal/agentmodel/apply_test.go TestAgentsStillFailsOnAPresentFileWithNoFrontmatter
- **writing through a symlinked agent file edits its destination**
  Proof: internal/agentmodel/apply_test.go TestApplyThroughASymlinkWritesTheDestination
- one model reaches every agent in the set
  Proof: internal/agentmodel/apply_test.go TestApplyReachesEveryAgentInTheSet
- **a bad name in the set means nothing is written at all** — not even the agents that
  came before it
  Proof: internal/agentmodel/apply_test.go TestApplyWritesNothingWhenAnyAgentIsUnwritable

Effort, reading:

- a declared effort is read back
  Proof: internal/agentmodel/effort_test.go TestReadEffortReturnsTheDeclaredEffort
- a file with no `effort:` reads as default, not as an error
  Proof: internal/agentmodel/effort_test.go TestReadEffortReportsDefaultWhenTheKeyIsAbsent
- an `effort:` in the body rather than the frontmatter is not read as the effort
  Proof: internal/agentmodel/effort_test.go TestReadEffortIgnoresTheBody
- listing reports each agent's effort beside its model
  Proof: internal/agentmodel/effort_test.go TestAgentsReportsEachCurrentEffort

Effort, writing, and the byte-for-byte promise:

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

The effort catalogue:

- the five levels are listed, weakest first, and `ultracode` is not among them
  Proof: internal/agentmodel/effort_test.go TestEffortCatalogueListsTheFiveLevels
- an unknown level is refused
  Proof: internal/agentmodel/effort_test.go TestUnknownEffortIsRefused
- **`haiku` supports no effort; `opus`, `sonnet` and the session default support all five**
  Proof: internal/agentmodel/effort_test.go TestWhichModelsSupportEffort
- **the levels a given model can run are answerable before a choice is offered**, weakest
  first, and nothing at all for a model that has none
  Proof: internal/agentmodel/effort_test.go TestEffortsForNamesWhatAModelCanRun

Applying a level to a set:

- one level reaches every agent in the set
  Proof: internal/agentmodel/effort_test.go TestApplyEffortReachesEveryAgentInTheSet
- **an agent on Haiku means nothing is written at all** — not even the agents before it
  Proof: internal/agentmodel/effort_test.go TestApplyEffortWritesNothingWhenAnyAgentCannotRunIt
- an agent with no declared model accepts a level
  Proof: internal/agentmodel/effort_test.go TestApplyEffortAllowsAnAgentOnTheSessionModel
- an unknown level is refused before any file is opened
  Proof: internal/agentmodel/effort_test.go TestApplyEffortRefusesAnUnknownLevel
- **moving an agent to Haiku clears a declared effort rather than leaving it dead**
  Proof: internal/agentmodel/effort_test.go TestApplyModelClearsEffortWhenTheModelSupportsNone
- **and the clearing is narrow: a model that does support effort leaves the level alone**
  Proof: internal/agentmodel/effort_test.go TestApplyModelKeepsEffortWhenTheModelSupportsIt
