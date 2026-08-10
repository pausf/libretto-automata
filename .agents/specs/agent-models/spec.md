# Agent Models

Governs: internal/agentmodel/**

Choose which model each payload agent runs on, so the work that is pattern-matching
over prose stops being billed at the rate of the work that is not.

The lever is **not fewer agents and not fewer tokens**. It is cheaper tokens for the
lenses that read prose and answer in prose.

## Outcomes

An agent's model is declared in its own frontmatter, and the user sets it without
opening a file.

```
---
name: review-design
description: …
tools: Read, Grep, Glob, Skill
model: haiku
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

| Value | Means |
|---|---|
| `opus` | the most capable. Max plans; metered on Pro |
| `sonnet` | the default working model |
| `haiku` | the cheap one — what this capability exists to make reachable |
| *default* | no `model:` key at all: the agent runs on whatever the session runs on |

Choosing *default* **removes** the key rather than writing a word meaning "no choice".
An absent key is already the language's way of saying it, and two spellings of one
state is a difference somebody will eventually treat as meaningful.

## Scope boundaries

**In:** reading and writing the `model:` key of the agent files in `agents/`, the
catalogue of legal values, and refusing everything else.

**Out, named so it cannot arrive quietly:**

- **A different model per install target.** Settled below under prior decisions: it
  cannot be done without breaking the symlink invariant, and it is not worth that.
- **Models for skills or commands.** Only agents carry a `model:`; a skill has no such
  key and inventing one would be a key no host reads.
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
  key** (`scripts/check-payload:62-67`). `model:` therefore needs no gate change.
  Verified, not assumed.
- **The writer must refuse a file with no frontmatter**, rather than inventing one. A
  file that does not open with `---` on line 1 is not an agent file by this repo's own
  check, and writing a block into it would manufacture an agent out of a document.
- This package **never touches an install target**. It reads and writes files in the
  repository. Where the targets come in is the CLI's problem, and the answer there is
  that they resolve to these same files.

## Prior decisions

- **One source of truth: the repository's own `agents/*.md`.** Asked whether the choice
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

- **The default value is an absent key, not the word `inherit`.** See outcomes.

## Task breakdown

1. `internal/agentmodel`: read the declared model from one agent file, and report a
   file that declares none as *default* rather than as an error.
2. `internal/agentmodel`: write or remove the key, leaving every other byte alone, and
   refuse a file with no frontmatter.
3. `internal/agentmodel`: the catalogue — the legal values, their labels, and the
   validation that refuses everything else.
4. `internal/agentmodel`: apply one model to a set of agents as one act, validating the
   whole set before writing any of it.

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
- an unknown model name is refused
  Proof: internal/agentmodel/catalogue_test.go TestUnknownModelIsRefused

Applying to a set:

- one model reaches every agent in the set
  Proof: internal/agentmodel/apply_test.go TestApplyReachesEveryAgentInTheSet
- **a bad name in the set means nothing is written at all** — not even the agents that
  came before it
  Proof: internal/agentmodel/apply_test.go TestApplyWritesNothingWhenAnyAgentIsUnwritable
