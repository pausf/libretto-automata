# Decisions

### Session 2026-08-21

- Q: does `libretto land` verify the landing or perform it? → A: verify-only — a
  read-only gate over the staged index. Applying a delta is a semantic merge a binary
  cannot do, and the observed failures were *incomplete* landings, which a verifier
  catches. If wrong: the command grows a write mode later; verify-only is
  forward-compatible. (assumed)
- Q: what does it check? → A: parts 2 and 4 of the landing — every file of the change
  folder deleted (no partial deletion), and each delta's `Targets:` capability spec
  modified in the same staged diff. Part 3 stays owned by `spec-drift --retired`;
  duplicating it in the binary is two sources of truth. If wrong: the check list grows,
  the ownership split does not. (assumed)
- Q: staged index or an existing commit? → A: the staged index, pre-commit — before the
  mistake is history. Change name optional: with none given, infer from staged change
  folder deletions. If wrong: a `--commit <sha>` mode can be added without moving the
  default. (assumed)
- Q: does a stale wiki block? → A: no — warn only, same rule as record-work's "a
  missing convenience never blocks a landing". (assumed)
- Q: what is a "stale wiki", operationally? → A: a marker-owned view (`README.md` /
  `wiki.html`, per `wiki.go`'s markers) exists in the specs directory, a capability
  spec is in the staged diff, and the view is not — derived from record-work's own
  clause. If wrong: one paragraph and two criteria in the cli delta move. (assumed)
- Q: is deleting a change folder with no delta a broken landing? → A: no — it is
  abandoning a queued idea, which the flow promises costs nothing; part 2 passes
  vacuously and the report says so. If wrong: one criterion flips. (assumed)
- Q: approach — self-contained land.go, an internal/land package, or payload bash?
  → A: self-contained `cmd/libretto/land.go` with an exec'd-git seam like metrics —
  one file per command, stdlib only. internal/ package lost (one consumer is a
  premature abstraction); payload bash lost (the command must run on installed-binary
  machines, and payload never implements delivery). If wrong: the logic lifts into
  internal/ the day a second consumer exists. (assumed)
- Q: does a staged rename OUT of the change folder satisfy part 4 for that file?
  → A: yes — the contract is "nothing under the folder survives the commit", and a
  rename-out leaves the folder empty; rename sources count as removals. If wrong:
  the parser reclassifies R sources and one test row flips. (assumed)
- Q: what base does "the staged diff" use? → A: ordinary `diff --cached` against
  HEAD; in-progress merges are out of scope. If wrong: a MERGE_HEAD path is added
  without moving the default. (assumed)
- Q: distinct exit codes for "nothing is landing" vs "parts missing"? → A: no —
  the spec distinguishes only zero from non-zero and main exits 1 on any error;
  distinct codes wait for a caller that needs them. (assumed)
- Q: (cutter finding) a `libretto` on PATH too old to know `land` exits non-zero
  with *unknown command*, naming no part — does that wedge the landing? → A: no —
  it falls on the absent side: a binary that does not know the command IS absent
  for this purpose; say the landing is unverified and continue. If wrong: one
  sentence in the clause and one clause in the criterion move. (assumed)
- Q: does the payload learn about the command? → A: yes, minimally — a second delta on
  the payload capability: record-work gains "where `libretto` is on PATH, run
  `libretto land` before the landing commit", same shape as the existing wiki clause,
  absent-binary path unchanged. A verifier nothing invokes verifies nothing. If wrong:
  the clause is one sentence to remove. (assumed)
