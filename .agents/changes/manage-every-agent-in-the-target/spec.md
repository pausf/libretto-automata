# Agent Models — delta

Targets: agent-models
Governs: internal/agentmodel/**

The package stops assuming the repository and starts working on a directory it is
handed.

## Outcomes

`Agents` and `Apply` take **the directory holding the agent files**, not a repository
root they join `agents/` onto.

```go
Agents(dir string) ([]Agent, error)
Apply(dir string, names []string, model string) error
```

- Every `*.md` in that directory is an agent, whether libretto created it or not.
- **Nothing else changes.** The frontmatter rules, the byte-for-byte promise, the
  validate-the-whole-set-first guarantee and the catalogue all hold exactly as they
  did — they were never about where the file came from.
- A directory that does not exist reports no agents, not an error. A target with no
  `agents/` yet is a state.

## Scope boundaries

**In:** the parameter change, and dropping the `agents/` join.

**Out:**

- **Knowing what a target is.** The package takes a path. `internal/target` stays
  unimported here — that is what keeps the layering true while the reach widens.
- **Deciding whether a file may be written.** Ownership is a `cli` and `panel`
  concern; this package writes where it is pointed.
- **Following a symlink to decide anything.** A symlinked agent file is written
  through, which edits its destination — that is ordinary file behaviour and the
  callers explain it.

## Constraints

- Still standard library only. Still no YAML.
- `SetModel` already refuses a file with no frontmatter, and that refusal is now the
  guard that keeps a stray `README.md` in an agents directory from being mangled.
  It was written for the repository's own tidy files; it now meets whatever is
  actually on disk.

## Prior decisions

- **The clause it reverses, and why.** The capability spec says the package "never
  touches an install target … Where the targets come in is the CLI's problem". Half
  of that survives: the CLI still decides *which* directory. What falls is the
  assumption that the answer is always this repository's `agents/`.

  The reason it fell is measured, not theoretical: the user has 22 agents, none of
  them libretto's, and the feature edited seven files that reach nothing. **The
  contract was written from the same misreading as the code, which is why no gate
  caught it.**

- **A path, not a target.** Handing the package a directory rather than a
  `target.Target` was chosen over the obvious alternative of importing
  `internal/target` here. It delivers the same reach and keeps the package testable
  against a `t.TempDir()` with no target in sight.

## Task breakdown

1. `internal/agentmodel`: `Agents` and `Apply` take the agents directory; drop the
   `dir` constant and the join. Update the tests to build a bare directory rather
   than a repo shape.

## Verification criteria

- a directory of agent files is listed whatever its path
  Proof: internal/agentmodel/apply_test.go TestAgentsListsEveryAgentSorted
- current models are read per agent
  Proof: internal/agentmodel/apply_test.go TestAgentsReportsEachCurrentModel
- **a directory that does not exist reports no agents rather than an error**
  Proof: internal/agentmodel/apply_test.go TestAgentsOnAMissingDirectoryIsEmptyNotAnError
- a file with no frontmatter in the directory is refused, and the set is untouched
  Proof: internal/agentmodel/apply_test.go TestApplyWritesNothingWhenAnyAgentIsUnwritable
- one model reaches every agent in the set
  Proof: internal/agentmodel/apply_test.go TestApplyReachesEveryAgentInTheSet
- **writing through a symlinked agent file edits its destination**
  Proof: internal/agentmodel/apply_test.go TestApplyThroughASymlinkWritesTheDestination
