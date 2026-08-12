# Plan — require-release-label-before-merge

From `spec.md` (Targets: ci). One commit for tasks 1+2 together: AGENTS.md wants the
test in the same commit as the logic it proves.

- [x] 1. `.github/workflows/require-release-label.yml` — the check itself.
      Spec: Outcomes. Criteria: all seven, written here, proven by task 2.
      Waits on: nothing. Evidence: commit b3597fd, six gates green.
- [x] 2. `cmd/libretto/label_workflow_test.go` — the seven proofs, static YAML reads
      plus the extracted script run under bash with 0/1/2 labels.
      Spec: Verification criteria. Waits on: task 1 (same commit).
      Evidence: go test ./cmd/libretto -run TestLabel, 7/7 PASS.
- [ ] 3. Add the check to `main`'s required status checks via `gh api`, read the
      protection back, record the result in the spec's criteria.
      Spec: Prior decisions (user chose 2026-08-12). Waits on: the request being open
      and the check having run once — phase 8, after the push question.
