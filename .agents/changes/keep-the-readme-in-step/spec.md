# keep-the-readme-in-step

Targets: readme

The README lists the commands a reader can type. `/libretto-attacca` shipped and the list
did not move.

**The fix is one row. The change is the guard**, because the rule this broke was already
written: `readme/spec.md` says the Commands section is "the reference table, one line per
command", and it did not stop anything. A rule that asks a person to remember has already
failed once here, and the second time it fails it will be a different command.

## Outcomes

1. **Every file in `commands/` has a row in the README's Commands table**, and a command
   that arrives without one **fails `go test ./...`**. Not a warning, not a note in a
   review — the same gate that already refuses a broken README link.

2. `/libretto-attacca` is in that table, and in the first-run walk beside the stops it
   answers. The walk names where the flow pauses; a command whose entire subject is those
   pauses belongs where they are described, not only in a reference list.

3. The guard names the missing command in its failure, so the fix is obvious from the
   output without opening two files.

## Scope boundaries

**In:** the README's Commands table and first-run walk, and one test in
`cmd/libretto/readme_test.go`.

**Out:**

- **skills and agents.** The README has no table of them and does not want one — it is the
  front door, and forty skills is reference material. A guard over `commands/` covers what
  a reader actually types.
- **wording.** The guard asserts the command's name appears, never how the row reads.
  `readme/spec.md` already scopes prose out for the reason that a test pinning a sentence
  gets deleted by whoever improves the sentence.
- **ordering inside the table.** A reader scanning six rows does not need them sorted, and
  a sort assertion is a second thing to fight when a row is added.
- **`scripts/check-payload`.** It stays about the payload's own integrity — frontmatter,
  references, reachability. Whether the README mentions something is a README fact.
- **a version bump.** The user's explicit instruction. Patch at most, and this change
  carries no promise to a user that was not already made.

## Constraints

- The test lives in `cmd/libretto/readme_test.go`, in `package main`, and reaches files
  through the existing `repoFile` helper — the same relative-path convention every test
  beside it uses.
- It reads `commands/` from disk rather than from a list in the test. **A hardcoded list is
  the same failure one level down**: somebody adds a command, forgets the list, and the
  guard passes.
- `flat()` and `section()` already exist in that file and are the right tools. No new
  helper unless one of them genuinely does not fit.
- `𝄞` and `♩♪♫♬` never enter the README's neighbours; the README itself is the one place
  `𝄞` is allowed, and this change does not add one.

## Prior decisions

- **The guard is a Go test, not a `check-payload` check.** Both are repo-only gates and
  either would run, so the tie is broken by ownership: `readme/spec.md` governs `README.md`
  and `cmd/libretto/readme_test.go` together, and a check about the README's contents
  belongs with the file it is about. `check-payload` answers whether the payload is
  coherent — a different question with a different failure.
- **It reads the directory, never a list.** Named as the whole point rather than an
  implementation detail, because a list-based version would have passed the exact failure
  this change exists to fix.
- **No version bump** — the user's call, 2026-08-12. `release:patch` at most on the request
  that carries it.
- **Ceiling named:** the guard proves a command's *name* appears somewhere in the Commands
  section. It cannot tell a real description from a placeholder row, and it will not catch
  a command whose row says something false. The replacement, the day that matters, is a
  criterion about what a row must contain — not a longer regex.

## Task breakdown

- [ ] `cmd/libretto/readme_test.go` — the guard, failing first against the current README.
- [ ] `README.md` — the `/libretto-attacca` row, and its line in the first-run walk.
- [ ] the delta applied onto `.agents/specs/readme/spec.md`, and this folder deleted.

## Verification criteria

- **a command file with no mention in the README's Commands section fails the suite**, and
  the failure names it
  Proof: cmd/libretto/readme_test.go TestEveryCommandIsInTheReadme
- the README's sections stay in reading order with Commands where it was
  Proof: cmd/libretto/readme_test.go TestReadmeSectionsAreInReadingOrder
- the first-run walk still names the flow and its stops after gaining a line
  Proof: cmd/libretto/readme_test.go TestReadmeWalksAFirstRun
- every relative link in the README still resolves
  Proof: cmd/libretto/readme_test.go TestReadmeLinksResolve

**The guard is verified the way the flow requires: written failing, watched to fail, then
made to pass.** Its first run happens against the README as it stands today, where
`/libretto-attacca` is genuinely missing — so the red is real rather than manufactured, and
that is the one run worth recording.
