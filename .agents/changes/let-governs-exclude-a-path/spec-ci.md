Targets: ci

# let-governs-exclude-a-path — the ci half

## Outcomes

1. ci's `Governs:` claims the three tests that prove its workflows and Makefile:
   `cmd/libretto/release_workflow_test.go`, `cmd/libretto/label_workflow_test.go` and
   `cmd/libretto/gates_test.go`, which cli now excludes.

## Scope boundaries

**In:** ci's `Governs:` line. **Out:** any ci promise; the tests themselves.

## Constraints

- The three tests stay where they are; only ownership moves.

## Prior decisions

- **ci owns the tests of its workflows.** Decided 2026-10-23: the three files read
  `.github/workflows/*.yml` and `Makefile`, which ci governs, so a change to them is a ci
  change; under cli they made drift warnings point at the wrong spec.

## Task breakdown

- ci's `Governs:` line.

## Verification criteria

- **When** `cmd/libretto/gates_test.go` is staged alone, spec-drift's default mode **shall**
  name ci's spec and not cli's. Ceiling named: the self-test proves the rule on a fixture;
  this repository's case is observed once, at the landing box, with the worktree's spec-drift.
  Proof: skills/record-work/spec-drift --self-test
