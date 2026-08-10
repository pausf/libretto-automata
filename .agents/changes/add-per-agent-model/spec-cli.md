# CLI — delta

Targets: cli
Governs: cmd/libretto/**

One new subcommand. Everything else in the CLI spec stands.

## Outcomes

| Command | Does |
|---|---|
| `models` | list every agent with the model it runs on, read-only |
| `models set <model> <agent>…` | write that model into each named agent |
| `models set <model> --all` | write it into every agent |

The name is `models`, not `save-token`. Nothing here is a token and nothing is saved
in the sense that word implies — the subcommand reads and writes a field, and it is
named after the field.

- **`models` alone changes nothing.** Same shape as `status`: read, report, exit zero.
- **`set` with no agents and no `--all` is an error**, not "all of them". A destructive
  default that fires on a forgotten argument is how every agent on the machine silently
  becomes `haiku`.
- **`--all` is explicit for exactly that reason**, and it is the bulk case the panel
  makes ordinary.
- **The whole set is validated before anything is written.** An unknown model, or an
  agent name the repo does not have, means nothing changed. Partial application would
  leave the user with no way to know how far it got.
- Output is plain text with no escape codes, like every other non-panel command.

### Scope, and the honest part

`--global` and `--project` keep their existing meaning: they select the target whose
installed agents are listed. They do **not** produce two different model settings.

Both targets' agent entries are symlinks to the same file in this repository
(`internal/link/own.go:23`), and a write through a symlink writes its destination. So
`models set haiku spec-writer` under either scope changes the same file and takes
effect everywhere.

**The command says so, once, when it writes.** A user who passed `--project`
reasonably expects a project-local effect; leaving that expectation unchallenged is
how a shared setting gets discovered by accident later.

An agent listed in one scope and not the other is a real difference and is shown as
such — that is the scope flag earning its place here. The listing marks every agent
the repository has but the target does not, and **only `linked` counts**: a conflict
is somebody else's file in our slot, and an agent whose slot is occupied does not
reach that target however much the repository wishes it did.

This clause had no criterion in the first draft of this delta, and the code did not
implement it — the listing was byte-identical under both flags and the green suite
had no way to notice. That is what an outcome with no `Proof:` behind it costs, and
it is why the two criteria below exist.

## Scope boundaries

**Out:** any flag that writes into a target rather than the repository; a `--force`;
reading the model of an agent this repo does not own.

## Constraints

- Dispatch follows the existing `switch` in `run` (`cmd/libretto/main.go:137`); a new
  case, not a new dispatch mechanism.
- The scope flags are already parsed once, before dispatch (`main.go:71-76`). This
  subcommand consumes the result; it does not parse them again.
- Exit non-zero when anything was refused, consistent with `install`.

## Task breakdown

5. `cmd/libretto`: the `models` subcommand — list, set, `--all`, validation, exit codes,
   and the one line that tells the truth about scope.

## Verification criteria

- `models` lists every agent with its current model and writes nothing
  Proof: cmd/libretto/models_test.go TestModelsListsEveryAgentAndChangesNothing
- an agent with no declared model is listed as running the session's
  Proof: cmd/libretto/models_test.go TestModelsShowsDefaultForAnUndeclaredAgent
- `set` applies one model to several named agents
  Proof: cmd/libretto/models_test.go TestModelsSetAppliesToEveryNamedAgent
- `set --all` reaches every agent
  Proof: cmd/libretto/models_test.go TestModelsSetAllReachesEveryAgent
- `set` with no agents and no `--all` is an error and writes nothing
  Proof: cmd/libretto/models_test.go TestModelsSetWithoutAgentsIsAnError
- an unknown model exits non-zero and writes nothing
  Proof: cmd/libretto/models_test.go TestModelsSetRejectsAnUnknownModel
- an unknown agent name exits non-zero and leaves the valid ones untouched
  Proof: cmd/libretto/models_test.go TestModelsSetRejectsAnUnknownAgentAndWritesNothing
- an agent the repository has but this target does not is marked as such
  Proof: cmd/libretto/models_test.go TestModelsMarksAgentsThatDoNotReachThisScope
- **the two scopes do not produce the same listing** — the flag changes something
  Proof: cmd/libretto/models_test.go TestModelsListingDiffersBetweenScopes
- writing under `--project` says the effect is not project-local
  Proof: cmd/libretto/models_test.go TestModelsSetUnderProjectScopeSaysTheEffectIsShared
- the subcommand is reachable from dispatch
  Proof: cmd/libretto/main_test.go TestRunDispatch
- output carries no escape codes
  Proof: cmd/libretto/models_test.go TestModelsOutputHasNoEscapeCodes
