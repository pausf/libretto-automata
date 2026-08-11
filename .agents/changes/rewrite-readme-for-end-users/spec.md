# rewrite-readme-for-end-users — delta

Targets: readme

**The capability does not exist yet.** Nothing in `.agents/specs/` claims `README.md`
today, which is why the file drifted into contributor notes without anything noticing.
This delta creates the capability when it lands: `.agents/specs/readme/spec.md`,
`Governs: README.md`.

Giving the front door a contract is the same move `ci` made for `.github/**` — a file no
Go package imports, proven by a Go test that reads it. That precedent is what makes the
criteria below falsifiable instead of decorative.

## Outcomes

What is true when this is done:

1. **`README.md` reads top-to-bottom as a sequence, not as reference.** Its headings
   appear in this order, and no reasoning section sits between them:

   | Order | Section | Answers |
   |---|---|---|
   | 1 | the name, the one-line claim, the panel image | what is this |
   | 2 | **What you get** | why would I install it |
   | 3 | **Install** | how do I get it |
   | 4 | **Your first run** | what do I type now |
   | 5 | **Commands** | the reference table, one line per command |
   | 6 | **Where it installs** · **The five states** · **Environment** | the reference detail |
   | 7 | **Learn more** | links out to `docs/` |

2. **Install is three steps and the reader can count them.** Go 1.26+, one
   `go install …@latest`, one `libretto install`. The paragraph explaining that the
   payload ships inside the module, the `$GOMODCACHE` tree diagram, and the
   version-in-the-path reasoning are not in the install section.

3. **"Your first run" is the section that did not exist.** It walks one real pass of the
   flow as numbered steps: `/libretto-flow "what you want"`, the spec stop, the plan
   stop, build, the review, the report, the push question — each step saying what the
   user does and what comes back. A reader who has installed the CLI and typed nothing
   yet is the reader this section is written for.

4. **Every "why" paragraph is relocated, not deleted.** The reasoning currently in
   `README.md` — no `--force`; `prune` is not `uninstall`; both scope flags is an error;
   dry by default; two queue commands and not one; why the payload is not compiled in;
   why aliases and not model ids; why `spec-drift` warns and never blocks — lands in
   `docs/DESIGN.md` (tool behaviour) or `docs/FLOW.md` (the flow and the queue), and the
   README links there.

5. **Nothing the README is the only home for is lost.** The badges, the panel image, the
   automaton paragraph, the licence, and `THIRD-PARTY.md` attribution stay in
   `README.md`.

## Scope boundaries

**In:** `README.md` rewritten. Reasoning moved into `docs/DESIGN.md` and `docs/FLOW.md`.
One new capability spec at landing. One new test file proving the structure.

**Out**, named so it does not creep in:

- **No new `docs/WHY.md`.** `DESIGN.md` is already "why it is built this way" — it opens
  with `Why symlinks, per item` — and `FLOW.md` is already the flow's reasoning. A third
  file would split the same subject three ways.
- **No change to any command, flag, output string or behaviour.** This change is prose. A
  README that documents behaviour that moved in the same commit cannot be reviewed
  against anything.
- **No screenshots or asciinema beyond the existing `docs/panel.svg`.** A recording goes
  stale silently; the SVG is already captured from the binary.
- **No translation.** The repository is English.
- **`AGENTS.md`, `CLAUDE.md`, `docs/STATE.md`, `docs/PLAN.md` untouched.** `AGENTS.md` is
  the contributor's door and it is doing its job.
- **No rewrite of `docs/FLOW.md` or `docs/DESIGN.md` beyond receiving the moved
  paragraphs.** Received text is appended into the section that already covers the
  subject, not re-argued.

## Constraints

- **`𝄞` stays, and only here.** It is outside the BMP and renders as tofu in a terminal;
  `AGENTS.md` permits it in `README.md` and nowhere else. `♩♪♫♬` are banned outright.
- **The install command must stay `@latest`.** The `v1.0.2` tombstone tag is what keeps
  `@latest` resolving on the `0.5.x` line; a README that pins a version instead would
  document around a retraction that is already working.
- **No version number invented into prose.** The tree diagram today says `@v0.5.0`; any
  path shown after this must read `@<version>` or come out of `git describe`, because a
  hardcoded version in a document desynchronises exactly like one in a source file.
- **Relative links only**, and every one has to resolve from the repository root.
- The six gates pass. `scripts/check-payload` must stay green — no skill or command
  reference may break because a README section moved.

## Prior decisions

- **The "why" moves to `docs/`, it does not stay and it is not cut.** Asked and answered
  by the user this session: the README becomes what it is → install → first steps →
  reference, and each piece of reasoning goes to `docs/` with a link. The alternatives
  offered were keeping everything in one file below a separator, and a minimal
  edit-in-place; both were declined.
- **The reasoning goes into `DESIGN.md` and `FLOW.md`, not into a new file.** Not the
  user's call — the existing files already own those subjects. Recorded here so the next
  session does not rediscover it as an open question.
- **README.md gets a capability spec; `docs/**` still does not.** `docs/` is deliberately
  uncontracted prose. The README is the one document whose shape is a promise to a
  stranger, and this change is the second time it has been rewritten for the wrong
  reader.
- **The structural test lives at `cmd/libretto/readme_test.go`.** It tests no package in
  that directory, which is already true of `gates_test.go` for `.github/**`: the `ci`
  spec states a `Proof:` may cite a test wherever the test can honestly live. A new
  package for one file that reads two documents is not worth its own directory.
- **The test asserts order and relocation, never wording.** A test that pins prose is a
  test that gets deleted the first time somebody improves a sentence, and a deleted test
  proves nothing.

## Task breakdown

1. **Move the reasoning out.** Every paragraph named in outcome 4 into `docs/DESIGN.md`
   or `docs/FLOW.md`, verbatim where it still reads correctly in its new home.
2. **Rewrite `README.md`** to the seven-section order in outcome 1, with **What you get**
   and **Your first run** written new.
3. **Write `cmd/libretto/readme_test.go`** — the four cases cited below.
4. **Create `.agents/specs/readme/spec.md`** from this delta, at landing, and add its row
   to `docs/SPEC.md`. Same commit as the code, and the change folder goes with it.

Order matters for 1 and 2 only: moving first means the rewrite deletes from a README
whose content is already safe somewhere else, so nothing can be lost by a badly aimed
edit.

## Verification criteria

Each criterion names the test that proves it. None of these tests exist yet.

- **Outcome 1** — the seven sections appear in `README.md` in that order, and no heading
  between them belongs to the moved-reasoning set.
  `Proof: cmd/libretto/readme_test.go TestReadmeSectionsAreInReadingOrder`

- **Outcome 2** — the install section contains the `go install …@latest` line and a Go
  version, and does not contain `GOMODCACHE`, a `@v` path, or the module-payload
  explanation.
  `Proof: cmd/libretto/readme_test.go TestInstallSectionIsStepsOnly`

- **Outcome 3** — a `Your first run` section exists, sits between `Install` and
  `Commands`, and names `/libretto-flow`.
  `Proof: cmd/libretto/readme_test.go TestReadmeWalksAFirstRun`

- **Outcome 4** — for every relocated subject, its sentence is absent from `README.md`
  **and** present in `docs/DESIGN.md` or `docs/FLOW.md`. This is the criterion that makes
  "moved" different from "deleted", and it is the reason the test reads two files.
  `Proof: cmd/libretto/readme_test.go TestMovedReasoningLandedInDocs`

- **Constraint: links** — every relative link in `README.md` resolves to a path that
  exists.
  `Proof: cmd/libretto/readme_test.go TestReadmeLinksResolve`

**No criterion for the glyph rule, deliberately.** `𝄞` is in seven files today —
`AGENTS.md`, `install.sh`, `docs/DESIGN.md`, `internal/ui/logo_test.go` among them —
because the rule bans it where a *terminal* renders it, not where a document discusses
it. A test asserting "only in README" fails on day one, and the part that matters is
already held by the `panel` capability and `internal/ui/logo_test.go`. The constraint
stands; it gets no new test.

Observed, not tested: that the result reads well to somebody who has never seen the
project. No test can hold that, and pretending otherwise is what produced the current
README — every fact in it is correct.
