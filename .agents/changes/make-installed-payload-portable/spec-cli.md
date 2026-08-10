# Delta: doctor reports the gate tools

Targets: cli

## Outcomes

- `libretto doctor`'s prerequisite report includes `rg` and `jq`, each attributed to
  the skill that needs it: `rg` to `record-work` (spec-drift runs on it) and to
  `find-work`'s plan scan; `jq` to `find-work` (parsing `jira --raw`). Each note
  names the install command, matching the existing `jira` row's style.
- A machine missing either shows `—` for that row and still exits 0 when links are
  clean — the prerequisite report stays informational.

## Scope boundaries

In: two `Prereq` rows in `prerequisites()` (`cmd/libretto/main.go`), one test.

Out, named:

- No doctor check for `fd`, `bat`, `eza` or any other nicety — nothing in the
  payload's gates depends on them. A row per fashionable tool is noise.
- No exit-code change. That is a standing spec clause, not an option.

## Constraints

- **The prerequisite report never affects the exit code** — already specified
  (cli spec, "prerequisites" clause) and proven by
  `TestPrerequisitesDoNotAffectTheExitCode`. The new rows inherit it.
- `onPath()` exists (`main.go:1011`); no new plumbing.

## Prior decisions

- `rg` and `jq` are hard dependencies of the payload's own gates (AGENTS.md names
  them under shell tooling); doctor's job is "what the payload expects on this
  machine", so omitting them was a gap, not a choice.

## Task breakdown

1. Add the `rg` and `jq` rows to `prerequisites()`; extend the doctor test to see
   them in the report.

## Verification criteria

- The report lists `rg` and `jq` with skill attribution and an install hint.
  Proof: cmd/libretto/main_test.go TestPrerequisitesIncludeTheGateTools
- Missing tools still do not affect the exit code.
  Proof: cmd/libretto/main_test.go TestPrerequisitesDoNotAffectTheExitCode
