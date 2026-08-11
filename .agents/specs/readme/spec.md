# readme

Governs: README.md cmd/libretto/readme_test.go

The front door. `README.md` is the only document in this repository written for somebody
who has never seen the project, and its shape is a promise to that reader: what this is,
how to install it, what to type, and only then the detail.

**This capability exists because nothing claimed the file.** `docs/` is deliberately
uncontracted prose and the README was treated the same way, so it drifted into contributor
notes twice — the second time reaching *what to actually type* on line 209 of 334. A path
no spec governs is a path where drift is nobody's finding.

## Outcomes

1. **The sections appear in reading order**, with no reasoning section between them:

   | Order | Section | Answers |
   |---|---|---|
   | 1 | the name, the one-line claim, the panel image | what is this |
   | 2 | **What you get** | why would I install it |
   | 3 | **Install** | how do I get it |
   | 4 | **Your first run** | what do I type now |
   | 5 | **Commands** | the reference table, one line per command |
   | 6 | **Where it installs** · **The five states** · **Environment** | the reference detail |
   | 7 | **Learn more** | links out to `docs/`, `.agents/specs/`, `AGENTS.md` |

   Attribution, what is not managed here, and the licence follow as a footer.

2. **Install is steps a reader can count.** A Go version floor, one `go install …@latest`,
   one `libretto install`, and one read-only command to confirm it landed. No explanation
   of the module cache, no `$GOMODCACHE` tree, no version-in-the-path argument.

3. **A first-run walk exists**, between Install and Commands, naming `/libretto-flow` and
   every stop the flow makes — including that the push is asked and never assumed.

4. **Reasoning lives in `docs/`, and being there is checked.** For each argument the
   README once carried — no `--force`; `prune` is not `uninstall`; both scope flags is an
   error; dry by default; two queue commands and not one; why the payload is not compiled
   in; aliases rather than model ids; `spec-drift` warns and never blocks — the text is
   **absent from `README.md` and present in `docs/DESIGN.md` or `docs/FLOW.md`**.

   This is the criterion that makes *moved* different from *deleted*. In a single-file diff
   the two are indistinguishable, which is why the proof reads three files.

5. **Every relative link resolves.** A README is mostly links once the reasoning leaves it,
   and a dead link in the front door is worse than a missing paragraph.

## Scope boundaries

**In:** `README.md`'s structure, and the test that holds it.

**Out:**

- **Wording.** No criterion here pins a sentence, except the relocated arguments in
  outcome 4, where the phrase *is* the anchor. A test that pins prose is a test somebody
  deletes the first time they improve a sentence.
- **`docs/**`.** Uncontracted by decision. This capability asserts that certain text is
  *present* there; it does not govern how those files are organised.
- **`AGENTS.md`.** The contributor's door, and a different reader.
- **Whether the result reads well.** No test holds that. See *Prior decisions*.
- **Screenshots and recordings** beyond `docs/panel.svg`, which is captured from the
  binary.
- **Translation.** The repository is English.

## Constraints

- `𝄞` is permitted in `README.md` and nowhere a terminal will render it. `♩♪♫♬` are banned
  outright — ambiguous width tears the layout.
- **The install line stays `@latest`.** The `v1.0.2` tombstone tag is what keeps `@latest`
  resolving on the `0.5.x` line; pinning a version in prose would document around a
  retraction that is already working.
- **No version number in prose.** `@<version>`, never `@v0.5.0`. A hardcoded version in a
  document desynchronises exactly like one in a source file, and silently.
- The proof reads files as text. No Markdown parser — a dependency added to check a
  document is a dependency added to check a checker.

## Prior decisions

- **The reasoning moved out; it was not cut.** Asked and answered by the user, 2026-08-11.
  Keeping everything in one file below a separator, and a minimal edit in place, were both
  offered and declined.
- **It went to `docs/DESIGN.md` and `docs/FLOW.md`, not to a new `docs/WHY.md`.** Those
  files already own those subjects — `DESIGN.md` opens its tail with `Why symlinks, per
  item`. A third file would split one subject three ways.
- **The test lives at `cmd/libretto/readme_test.go`.** It tests no package in that
  directory, which is already true of `gates_test.go` for `.github/**`, and it reuses that
  file's `repoFile` helper. A new package for one file that reads three documents is not
  worth its own directory.
- **Every substring assertion normalises whitespace first**, via `flat()`. Without it
  `"checksum database"` never matched, because the README wraps between those two words —
  a guard that silently could not fire. The bug was in the first version of this test and
  was caught by reading the red run rather than counting it.
- **No criterion for the glyph rule.** `𝄞` is in seven files, because the rule bans it
  where a *terminal* renders it, not where a document discusses it. A test asserting "only
  in README" fails on day one; the part that matters is held by the `panel` capability and
  `internal/ui/logo_test.go`.
- **"Does it read well to a stranger" is verified by reading, and said out loud.** It is
  the one thing here no test can hold, and pretending otherwise is what produced the
  README this replaced — every fact in that version was correct.

## Task breakdown

Held by this capability going forward, not open work:

- keep the seven sections in order when a section is added
- when an argument moves out of the README, land it in `docs/` in the same commit
- when a link is added, `TestReadmeLinksResolve` is the check — it fatals if the pattern
  matches nothing, so it cannot pass vacuously

## Verification criteria

- **Outcome 1** — the five load-bearing headings appear in order.
  Proof: cmd/libretto/readme_test.go TestReadmeSectionsAreInReadingOrder

- **Outcome 2** — the install section contains the install commands and a version floor,
  and none of `GOMODCACHE`, `@v0.`, `checksum database`.
  Proof: cmd/libretto/readme_test.go TestInstallSectionIsStepsOnly

- **Outcome 3** — a `Your first run` section exists and names `/libretto-flow`, the spec
  stop, the plan stop and the push.
  Proof: cmd/libretto/readme_test.go TestReadmeWalksAFirstRun

- **Outcome 4** — all eight relocated arguments absent from `README.md` and present in
  `docs/DESIGN.md` or `docs/FLOW.md`.
  Proof: cmd/libretto/readme_test.go TestMovedReasoningLandedInDocs

- **Outcome 5** — every relative link resolves, and the link pattern matched something.
  Proof: cmd/libretto/readme_test.go TestReadmeLinksResolve
