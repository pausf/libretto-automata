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
- Q: does the payload learn about the command? → A: yes, minimally — a second delta on
  the payload capability: record-work gains "where `libretto` is on PATH, run
  `libretto land` before the landing commit", same shape as the existing wiki clause,
  absent-binary path unchanged. A verifier nothing invokes verifies nothing. If wrong:
  the clause is one sentence to remove. (assumed)
