# `models effort` — delta

Targets: cli

Governs: cmd/libretto/models.go

`libretto models` reads the model of every agent and `libretto models set` writes it. Both
grow a second dimension, and the writing half gets its own verb rather than a flag.

## Outcomes

```
libretto models                                    # model and effort, per agent
libretto models set opus review-lens-design        # unchanged
libretto models effort xhigh review-lens-design    # new
libretto models effort xhigh --all
libretto models effort default review-lens-tests   # removes the key
```

- **`models` prints effort in its own column**, beside the model, and prints `(session)`
  for an agent that declares none — the same word the model column already prints for the
  same state.
- **The listing's trailer names the levels**, the way it already names the models: the
  five, weakest first, each with the one line saying when to reach for it, under the same
  `Resolved` date.
- **`models effort` is a sibling verb of `models set`**, taking a level, then agent names
  or `--all`. Not a flag on `set`: the two keys are independent, and `set opus --effort
  xhigh` would make changing one without touching the other the awkward path.
- **Naming no agents and passing no `--all` is refused**, identically to `set`. A
  destructive default that fires on a forgotten argument is how every agent on the machine
  silently becomes the same thing.
- **A level the catalogue does not know is refused**, naming the five.
- **An agent whose model cannot run the level is refused by name**, and nothing is
  written — the whole set, or none of it, as `set` already promises.
- **`models set` reports a cleared effort on the row it cleared.** Moving an agent to a
  model with no effort support drops its `effort:` key, and a drop the user cannot see is
  a silent edit to a prompt file.
- **The `shared` mark and its trailer are unchanged**, and apply to both verbs: a row
  whose file this repository owns reaches every destination linked to it.

## Scope boundaries

**In:** the `models effort` verb, the effort column in `models`, the levels trailer, and
the refusals above.

**Out:**

- **`--effort` on `models set`.** Named above. Two keys, two verbs.
- **A third verb that sets both at once.** Two commands is not a burden worth a
  combinatorial flag surface, and the day it is, it is one line of dispatch.
- **`models effort` on a target that takes no agents.** Already handled: `agentsDir`
  returns empty and the listing says so rather than failing. Unchanged.
- **Changing the exit-code contract.** A refusal is a non-zero exit and an error on
  stderr, exactly as `set` refuses today.

## Constraints

- **All validation stays in `internal/agentmodel`.** This file dispatches, resolves the
  directory from the target, and renders. It does not learn which model supports which
  level — that would put one rule in two places and let the panel reach the weaker copy.
- **`cmd/libretto` has tests now** (`models_test.go`), so this arrives with them rather
  than being held up by manual runs.
- **The word `default` is the CLI's spelling of the absent key**, for effort as it already
  is for the model. Both verbs map it to the one `agentmodel.Default` — there is no
  separate `EffortDefault`, because two names for the empty string is the same "two
  spellings of one state" the capability already refused when it declined to write the
  word `inherit`.

## Prior decisions

- **A verb, not a flag.** *(assumed under `/libretto-attacca` — nobody was asked)* The
  alternative was `models set opus --effort xhigh`, which reads well for a fresh agent and
  badly for every later edit: it forces the model to be restated to change the effort, and
  a restated model is a write nobody asked for. **What changes if this is wrong:** the
  dispatch gains one branch and the flag delegates to the same `ApplyEffort`. No format,
  no validation, no test of the writer moves.

- **Effort is a column, not a second grouping.** See the panel delta — the decision is one
  decision and it is recorded there, because that is where it is visible.

## Task breakdown

1. `models` renders the effort column and the levels trailer.
2. `models effort <level> <agents…|--all>` dispatches, with `set`'s argument handling
   reused rather than reimplemented.
3. `models set` reports a cleared effort on the affected rows.
4. `README.md` and `AGENTS.md` gain the verb in the tables that already list `models set`.
   **No delta targets `readme`:** that spec promises the front door says what to type, and
   a new row in an existing table keeps that promise rather than moving it. The
   `readme` gate (`cmd/libretto/readme_test.go`) checks payload slash commands, not the
   binary's subcommands — verified, not assumed — so nothing enforces this row and it is
   listed here because that is the only thing that will make it happen.

## Verification criteria

- the listing shows each agent's effort, and `(session)` when it declares none
  Proof: cmd/libretto/models_test.go TestModelsListsEffortBesideTheModel
- the trailer names the five levels
  Proof: cmd/libretto/models_test.go TestModelsListsTheEffortCatalogue
- `models effort` writes the level to the named agents and no others
  Proof: cmd/libretto/models_test.go TestModelsEffortWritesOnlyTheNamedAgents
- `--all` reaches every agent in the destination
  Proof: cmd/libretto/models_test.go TestModelsEffortAllReachesEveryAgent
- **no agents named and no `--all` is refused, and nothing is written**
  Proof: cmd/libretto/models_test.go TestModelsEffortRefusesWithNothingNamed
- an unknown level is refused and names the five
  Proof: cmd/libretto/models_test.go TestModelsEffortRefusesAnUnknownLevel
- **an agent on Haiku is refused by name and the whole set is left alone**
  Proof: cmd/libretto/models_test.go TestModelsEffortRefusesAnAgentThatCannotRunTheLevel
- `default` removes the key rather than writing a word
  Proof: cmd/libretto/models_test.go TestModelsEffortDefaultRemovesTheKey
- **`models set haiku` on an agent declaring an effort reports the clearing**
  Proof: cmd/libretto/models_test.go TestModelsSetReportsAClearedEffort
