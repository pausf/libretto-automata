# Delta: `--opencode` is no longer a skills-only destination

Targets: cli

## Outcomes

- **`help` describes `--opencode` as taking skills and commands**, and `--codex` as
  skills only. The two rows stop being identical, because the destinations stopped
  being identical.
- **The skills-only rule survives unchanged and applies to Codex alone.** A
  destination installs the kinds it accepts and the rest are absent from the run and
  the summary, never errors — that is still exactly right, and it is what makes this
  whole change two lines of production code. What was wrong was the list of
  destinations it applied to, which is a correction to this spec's prose rather than
  a promise about the tool, so it lives in the task breakdown and not here.
- **`install --opencode` links skills and commands** and creates no `agents`
  directory. Its summary says so, through the same counting path every destination
  uses.
- The environment table, the flag list and the dispatch/help/README agreement gates
  are unchanged. No flag is added, removed or renamed.

## Scope boundaries

**In:** the two help lines, the skills-only paragraph, and the criteria below.

**The README row is not this spec's**, and the first draft claimed it. `Governs:` here
is `cmd/libretto/** install.sh`; `README.md` belongs to the **readme** capability, and
a path two specs claim is a path recorded in neither. The row still has to change —
`--opencode` is described there as skills-only and that becomes false — but it is a
factual correction that moves no `readme` criterion, so it lands as a documentation
edit under that capability's ownership with no delta of its own. Named here so the
edit is not invisible, not so this spec owns it.

**Out, named:**

- **A new flag.** `--opencode` already exists and already resolves; only what it
  accepts has moved.
- **`doctor` growing a per-tool prerequisite check** — is OpenCode installed, does it
  have a commands directory. Named out when the flag was added, 2026-08-14, and
  unchanged: it reports the destination it was handed. What brings it back is still
  users confused by linking into a tool they do not have.
- **A summary that distinguishes a command from a skill by kind.** The existing
  per-kind counting already does this; nothing new is needed and nothing new is
  added.

## Constraints

- Configuration stays environment-only. No new variable — `OPENCODE_HOME` already
  covers the root and the commands directory hangs off it.
- The help text is asserted by test against the dispatch table, so a row that
  describes a destination wrongly is a documentation bug the gates cannot catch. It
  is caught here by reading, which is the honest statement of the ceiling.

## Prior decisions

- **One command, one destination** — user decision, 2026-08-14. Unchanged: this adds
  a kind to a destination, not a second destination to a run.
- The wording discipline for the destination table is the existing five rows'. A row
  says what the destination is and what it takes, in the same shape as its neighbours.

## Task breakdown

1. `cmd/libretto/main.go` help: `--opencode` reads *skills and commands*.
2. The skills-only paragraph in this spec names Codex.
3. `README.md` destination table: the `--opencode` row loses "— skills only". Owned by
   the **readme** capability, listed here because this change is what makes it false.

## Verification criteria

- **an opencode install links skills and commands and leaves every other destination
  alone**
  Proof: cmd/libretto/scope_test.go TestInstallOpencodeLeavesOthersAlone
- a codex install still links skills only and leaves every other destination alone
  Proof: cmd/libretto/scope_test.go TestInstallCodexLeavesOthersAlone
- `help` still names every destination flag and every destination environment
  variable
  Proof: cmd/libretto/main_test.go TestHelpNamesEveryDestination

  **Ceiling named:** it proves the flag appears, never that the row beside it says
  something true — `--opencode` described as skills-only would satisfy it. The row's
  accuracy is held by reading, and the replacement the day that bites is a criterion
  about what a row must contain.
