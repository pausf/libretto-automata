# CLI — delta

Targets: cli
Governs: cmd/libretto/**

`models` acts on the target the scope flag names, and the flag finally means what it
says.

## Outcomes

`models` and `models set` read and write **the agents directory of the selected
target** — `~/.claude/agents` under `--global`, `<cwd>/.claude/agents` under
`--project` — not this repository's `agents/`.

Every `*.md` there is listed, whether libretto created it or not. On the machine that
prompted this change that is 22 agents rather than 7, and the seven were the wrong
seven: none of them were installed in the target being displayed.

### Where a write lands, and saying so

An agent file in a target is one of two things, and they behave differently:

| The file is | Writing it affects | Marked |
|---|---|---|
| a symlink into this repository | **every project on the machine** — one file, many targets | `shared` |
| a real file in that target | that target only | *(nothing)* |

**The listing marks the shared ones**, and `set` says which kind it just wrote. This
replaces the blanket "in effect for every project" line, which is now true of only
some rows and would be a lie about the rest.

That reverses the cli spec's current claim that `models` is the one command whose
scope does not decide where it writes. It does now.

### Writing files libretto did not create

**`set` edits real files this tool did not create, and does so without asking.**
Chosen deliberately over a per-write confirmation, because the user's own agents are
the point of the feature and a prompt on every one of 22 rows is a prompt people
learn to dismiss.

This is not `--force` returning by the back door, and the distinction is the whole
of it: `AGENTS.md` forbids **overwriting** something the tool did not create —
replacing a foreign item with our symlink, destroying what was there. This replaces
one line of frontmatter, in a file the user named, at their request. Nothing is
clobbered, nothing changes hands, and `install` still refuses a conflict exactly as
before.

- A target with no `agents/` directory says so and exits zero.
- A `*.md` there with no frontmatter is refused by name, and **nothing in the set is
  written** — the existing all-or-nothing guarantee now protects foreign files too.

## Scope boundaries

**In:** the target's agents directory as the subject of `models`, the shared marker,
and the honest per-write message.

**Out:**

- **Editing anything but `model:`.** One key, in files we do not own.
- **Creating an agent** in a target, or deleting one.
- **`--all` reaching outside the selected target.** One scope per invocation, as
  every other command already promises.
- **Reconciling the two scopes.** `models --global` and `models --project` are two
  questions with two answers; nothing merges them.

## Constraints

- Scope flags are parsed once before dispatch (`main.go:71-76`); this consumes the
  result, as it already does.
- `link.Owned(repoRoot, path)` already answers "is this ours" and is the only thing
  the shared marker needs. No second definition of ownership.
- A target that does not accept agents at all lists nothing rather than failing.
- Plain text, no escape codes, unchanged.

## Prior decisions

- **The scope warning was right about the danger and wrong about the scope.** The
  shipped message told every user their write was machine-wide. That was true when
  the subject was always the repository's own file; it is now true only of the
  symlinked rows, and telling somebody their local edit is global is the same class
  of error as the silence it replaced.

## Task breakdown

2. `cmd/libretto`: point `models` at the target's agents directory, mark shared rows,
   and make the post-write message name which kind was written.

## Criteria this delta retires

Three criteria in the cli capability spec describe behaviour that no longer exists.
Their tests are gone, so `--anchors` reports them broken until phase 8 folds this delta
in — that is the delta doing its job, not rot, and it is named here so the two are
distinguishable:

| Retired criterion | Why |
|---|---|
| `TestModelsListsEveryAgentAndChangesNothing` | replaced by `TestModelsListsTheTargetsAgents` — "every agent" now means the target's |
| `TestModelsMarksAgentsThatDoNotReachThisScope` | the `· not linked here` marker has nothing left to mark: every listed agent is in the target by construction |
| `TestModelsSetUnderProjectScopeSaysTheEffectIsShared` | the message was unconditional and is now per row; `TestModelsSetSaysWhenAWriteIsShared` and `TestModelsSetDoesNotOverclaimALocalWrite` replace it as a pair |

## Verification criteria

- the listing is the target's agents, not the repository's
  Proof: cmd/libretto/models_test.go TestModelsListsTheTargetsAgents
- **an agent the target has and the repository does not is listed and editable**
  Proof: cmd/libretto/models_test.go TestModelsEditsAnAgentTheRepositoryDoesNotShip
- a symlink into the repository is marked shared; a real file is not
  Proof: cmd/libretto/models_test.go TestModelsMarksSharedAgents
- writing a shared agent says the effect reaches every project
  Proof: cmd/libretto/models_test.go TestModelsSetSaysWhenAWriteIsShared
- **writing a target-local agent does not claim to be machine-wide**
  Proof: cmd/libretto/models_test.go TestModelsSetDoesNotOverclaimALocalWrite
- the two scopes list different agents
  Proof: cmd/libretto/models_test.go TestModelsListingDiffersBetweenScopes
- a target with no agents directory says so and exits zero
  Proof: cmd/libretto/models_test.go TestModelsWithNoAgentsSaysSo
- a stray file with no frontmatter is refused and the whole set is left alone
  Proof: cmd/libretto/models_test.go TestModelsSetRefusesAStrayFileAndWritesNothing
- `set` with no agents and no `--all` is still an error
  Proof: cmd/libretto/models_test.go TestModelsSetWithoutAgentsIsAnError
- output carries no escape codes
  Proof: cmd/libretto/models_test.go TestModelsOutputHasNoEscapeCodes
