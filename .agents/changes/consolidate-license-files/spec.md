# consolidate-license-files — delta

Targets: payload

Four licence files sit at the root: `LICENSE`, plus `LICENSE-caveman`, `LICENSE-ponytail`
and `LICENSE-superpowers`. The ask read them as *alternatives to choose between*, which
would be confusing. They are not — they are the upstream licences of vendored skills this
repository ships, and `THIRD-PARTY.md` already records which is which.

**So the confusion is real and the fix is presentation, not consolidation.** Deleting one to
tidy the root is a licensing problem wearing a docs problem's clothes: a vendored copy has to
carry its own licence text.

**And the finding this opened with is bigger than the move.** No capability's `Governs:`
claims `LICENSE*` or `THIRD-PARTY.md`, so a path nobody owns is a path where drift is
nobody's finding — the exact sentence the `readme` capability was created to answer about
`README.md`. `scripts/check-payload` *parses* `THIRD-PARTY.md` to derive the vendored skill
list, so that file is already load-bearing for a gate while being governed by nothing.

## Outcomes

1. **`LICENSE` is alone at the root.** The three vendored licence texts live in
   `licenses/`, keeping their names.

2. **`THIRD-PARTY.md` points at them where they now are**, and its links resolve. It also
   **says why they are in a directory** — that none of the three is an alternative to root
   `LICENSE`, and that deleting one to tidy the root would be a licensing failure. The ask
   misread the root listing exactly that way, so a move that leaves no explanation invites the
   same misreading from the next reader; the paragraph is the half of the fix that is
   presentation rather than layout.

   **Added after the 6→7 reviewer found it as unnamed scope.** The prose shipped with no
   outcome behind it, which is how scope arrives without being asked for even when it is
   right — the reviewer was correct to name it, and the answer was an outcome rather than a
   deletion.

3. **`payload` governs the vendoring record.** Its `Governs:` line gains `THIRD-PARTY.md`
   and `licenses/**`. `payload` is the right owner because it already carries the constraint
   that vendored items are copied unmodified *with their licence and version recorded* —
   `THIRD-PARTY.md` is that record.

4. **Every relative link in `THIRD-PARTY.md` resolves, and it is checked.** Nothing tested
   them before, which is what made this move able to break the file silently.

## Scope boundaries

**In:** three `git mv`s, the link lines and the explanatory paragraph in `THIRD-PARTY.md`,
`payload`'s `Governs:` line, and one new test case.

**Out**, named:

- **no licence file is deleted, reworded or merged.** Not one byte of a vendored licence text
  changes. This is what the ask's wording pointed at and the one thing that must not happen
- **`LICENSE` does not move.** GitHub reads the root `LICENSE` to display the repository's
  licence; moving it trades a tidy root for a lost badge and a lost API field
- **the vendored-skill table is not touched.** `scripts/check-payload` parses those rows
- **no `NOTICE` file, no SPDX headers, no licence scanner.** The record exists and is accurate
- **`README.md`'s licence badge and footer are unchanged.** They point at root `LICENSE`
- **no test that a licence text is *correct*.** That is a legal reading, not a check

## Constraints

- **`scripts/check-payload` parses `THIRD-PARTY.md`** at line 116 and fails loudly on a parse
  returning nothing or everything. Run and read after the move, never assumed
- `git mv`, so history follows the files
- the new case joins `cmd/libretto/readme_test.go`, which already reads documents outside its
  own package and owns `repoFile`
- **whitespace is normalised before matching.** That test file has now shipped two guards that
  could not fire because Markdown wrapped; a third would be nobody's accident
- all six gates green

## Prior decisions

- **`licenses/`, lowercase and plural** — *assumed, nobody was asked.* `/libretto-attacca`
  answered the contract stop. It sorts away from `LICENSE`, which is the point of the move.
  What changes if wrong: a rename and the links.
- **The filenames keep their `LICENSE-<origin>` form** — *assumed.* The name says what the file
  is when read alone, which is how a licence file gets read.
- **`payload` takes ownership rather than a new capability** — *assumed.* Four files and a
  record do not justify their own directory. **Ceiling:** `payload`'s `Governs:` is now wider
  than the payload directories, so a licence-only commit reads as payload drift. That is the
  intent; the replacement, if it becomes noise, is a `vendoring` capability, not a narrower
  glob.
- **The link check is added rather than trusted** — *assumed.* The move is exactly the kind of
  change that breaks a link, and nothing was watching.
- **Outcome 3 carries no `Proof:`, deliberately.** The obvious citation was
  `spec-drift --trace`, and it cannot fail: `--trace` is a map and **returns 0 whatever it
  finds**, as its own line 42 says. A criterion citing a gate that can never be red reads as
  proven and is worse than one with no citation. It is an observation in the report instead.
  The alternative was a test asserting one line of one spec file contains two strings, which is
  machinery to check a document the next reader checks for free.

## Task breakdown

1. **The link guard first, red.**
2. **Then the move and the links**, together.
3. Widen `payload`'s `Governs:` and land this delta, deleting the change folder.

**Steps 1 and 2 are one box.** The guard alone is green against the tree as it stands, so it
proves nothing; the move alone leaves dead links with nothing watching.

## Verification criteria

- **every relative link in `THIRD-PARTY.md` resolves**, matched over the flattened document,
  and the link pattern must have matched something — a scan finding no links passes vacuously.
  Proof: cmd/libretto/readme_test.go TestThirdPartyLinksResolve

- **the three vendored licence files exist under `licenses/` and no `LICENSE-*` remains at the
  root**, asserted directly because a link nobody updated still points at a file that is still
  there — so the link half alone passes with one moved and two forgotten.
  Proof: cmd/libretto/readme_test.go TestThirdPartyLinksResolve

- **`scripts/check-payload` still derives the vendored skill list from `THIRD-PARTY.md`.**
  Proof: scripts/check-payload

**What nothing here tests:** that a licence text is the correct one for its upstream, or that
the vendored version recorded beside it is accurate. Both are readings, not checks.
