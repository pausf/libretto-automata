# check-clause-coverage-in-anchors

Tracker: none
Targets: payload
Durable decisions: every entry in Prior decisions below

## The ask, verbatim

> Cobertura cláusula→test en `--anchors`. El patrón "un criterio cita un test que cubre
> media cláusula" apareció seis veces en tu ledger de lecciones — es tu bug recurrente
> número uno. Hoy `--anchors` prueba que la cita resuelve, no que el test ejerza la
> condición del `When`/`While`. Aunque sea heurístico, se paga solo.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Extend `spec-drift --anchors` so a cited test is checked — even heuristically — for
exercising the trigger condition named in the criterion's EARS clause. Prior art:
spec-kit `/speckit.analyze` cross-artifact coverage.

## Decision log

### Session 2026-10-22 (attacca — nobody was asked; every answer below is assumed or measured)

- Q: how often does the pattern recur? → A: **eight times in `.agents/lessons.md`, not six**,
  counted by reading every entry, plus a ninth recorded only in auto-memory (2026-08-14, the
  host-neutral marker criterion). The eight: 2026-08-17 add-contributing-guide (twice),
  2026-08-17 split-readme-into-sections, 2026-08-18 add-specs-wiki, 2026-08-18
  wiki-pages-redesign, 2026-08-19 wiki-flow-board, 2026-08-19
  make-specs-and-plans-collaborative (twice). Six cite a Go test, two cite
  `scripts/check-payload`. Measured, not assumed.
- Q: which half of the clause did those proofs miss? → A: **the trigger in one of the eight**
  (wiki-flow-board: "where the ledger is absent or empty of entries"). Six missed an item
  of an enumerated *response* — a link, three article fields, five of six gate lines, a
  tools list, three stops, a guide section — and one asserted on a phrase that also lived in
  another file, so it could not fail. Measured by reading each entry and its criterion.
- Q: does the requested heuristic — the trigger clause's key terms looked up in the cited
  test's body — earn a place in the gate? → A: **no, measured.** Prototyped against every
  criterion in `.agents/specs/*/spec.md` (706 criteria, 719 `Proof:` citations). It flags
  15 of the 52 Go-proved trigger criteria that carry a term (2.1% of all 706); **all ten
  flags hand-checked were false**. Widened to every term in the whole criterion it flags
  106 of 304 (15% of 706); all ten sampled were false again. Run at the pre-fix commit
  of each of the six Go-proved ledger cases, the trigger variant catches none for the right
  reason — its only flag there fires on a placeholder template, and still fires on today's
  fixed test — and the whole-criterion variant catches one. It cannot run before the test
  exists, and it has no body to read for the two `check-payload` cases. Numbers under
  *Approach*. (assumed: that this outweighs the ask's "even heuristically". If wrong, the
  prototype's term extraction ships behind its own flag, warn-only, never inside
  `--anchors`, and its measured false-positive rate goes in the flag's help line.)
- Q: what is built instead? → A: **a structural pointer, warn-only**: a criterion carrying
  two or more `shall` on exactly one `Proof:` line and naming no ceiling is *joined*.
  `--anchors` prints one line counting them; `--trace` lists them (assumed). Measured
  precision 6 of 26 hand-checked (23%), and in this one session it found three gaps live
  today. If wrong — the user wants the gate to stay as it is — nothing ships, and the
  measurement moves into *Prior decisions* as the reason, alone.
- Q: fail or warn? → A: **warn, never fail** (assumed). 36 of the 99 criteria in the ten
  landed EARS-era delta files history recovers are joined, and a hand-checked sample of
  ten held one real gap. A gate failing one criterion in three, wrong three times in
  four, gets muted. If wrong — a delta must not land joined — `joined` returns 1 on
  change documents only, and the escape below is already the fix the author would make.
- Q: what does the `--anchors` line say, per criterion or a count? → A: **one count line**,
  change documents and capability specs separately, naming `--trace` as the list (assumed).
  Per-criterion lines on every gate run is the cry-wolf shape; 72 capability criteria are
  joined today. If wrong, the change-document half prints per criterion, which on today's
  tree is zero lines.
- Q: how does an author answer a joined criterion? → A: **split it, add a second `Proof:`
  line, or write the word *ceiling* in it** — the house phrase is "Ceiling named:"
  (assumed). A declaration in the document, never a flag, the same as `Durable decisions:
  none`. If wrong, a dedicated marker is one more pattern for the awk to read.
- Q: what happens to prose-proved citations (`check-payload`, `spec-drift --self-test`)? →
  A: **they are counted like any other.** The predicate reads the criterion and the number
  of its `Proof:` lines and never opens the proof, so the kind of citation is irrelevant
  (assumed). This is the case the rejected heuristic had to skip, and it is where two of
  the eight ledger cases live.
- Q: does it run inside `--anchors` or as a new gate? → A: inside, as a soft line after
  `--retired` — the recorded *Prior decision* "New checks ride inside `--anchors` rather
  than becoming new gates". No seventh gate and no new flag. Settled by the spec, not
  assumed.
- Q: do the skills that explain `--trace` change? → A: yes — `record-work`'s `--trace` table
  gains the row and `review-spec`'s sentence naming what `--trace` reports gains the
  section, each pinned by one `check_wiring` row (assumed). Both enumerate `--trace`'s
  output today, so leaving them is a list going stale. If wrong, the skills stay and the
  section explains itself in its own heading.
- Q: are the three live gaps the measurement found fixed here? → A: no; they are named
  under *Scope boundaries* as follow-ups (assumed). Two are `cli` and one is a self-test
  case; mixing them in widens a payload change into a cli one. If wrong, each is one test
  case.
- Q: the fork — how is it built? → A: **A, `is_joined` plus `joined_of` in `spec-drift`, a
  count from `anchors_and_ears` and a section in `--trace`** (assumed, the recommended
  option). B — the requested term-in-body heuristic, warn-only — lost on the measurement
  above. C — build nothing and record the measurement as a *Prior decision* — lost because
  the structural check found three live gaps in one run for the cost of one awk function.

### Session 2026-10-22, after review-spec (the run acting as author on the reviewer's findings)

- Q: does `joined_of` list a criterion with no `Proof:`? → A: no. It prints only criteria
  carrying at least one `Proof:`. A zero-proof criterion is already `--trace`'s *unproven
  criterion*, and listing it twice is noise (assumed). The `joined_of` criterion now says
  so and its self-test fixture holds an unproven criterion. If wrong, `joined_of` prints
  `0<TAB>text` too, and `is_joined` still rejects it because its count is not 1.
- Q: what exactly do the two wiring rows pin? → A: for `record-work`, the joined row of
  its `--trace` table, pattern `^\| \*\*joined criterion\*\* \|`, so prose that only
  mentions the word cannot satisfy it. For `review-spec`, the prompt wording, pattern
  `split the criterion or name its ceiling` (assumed). Both criteria now say "Proved as
  wiring only" and name the ceiling, the way `payload` states that limit elsewhere. If
  wrong, each row widens to the word `joined` alone and the ceiling grows to match.
- Q: is the zero case promised? → A: yes. While no criterion is joined, `--anchors` prints
  no joined line, and its output is byte-identical to the three existing sections composed
  as `anchors_and_ears` composes them today. The self-test compares `--anchors` against
  `anchors; ears; retired` called directly in the same fixture (assumed). The old script
  cannot run inside the self-test, so the composition is the stand-in for "today". If
  wrong, the comparison becomes a golden string, which goes stale every time a section's
  wording changes.
- Q: is the `--trace` summary count proved? → A: yes, its own criterion and a self-test case
  over a fixture with two joined criteria, asserting `, 2 joined criterion(a)` exactly. A
  missing count or a wrong one fails it (assumed). If wrong, the count leaves the summary
  line and only the section remains.
- Q: does the existing `payload` criterion "the whole tree is traceable, not only what is
  staged" get a sibling or an amendment? → A: an amendment, at landing. Its list of what
  `--trace` names gains a fourth item, joined criteria. The replacement text is stated
  under *Approach* and the landing has its own box (assumed). A parallel criterion would
  be two sentences describing one map, and the one nobody edited would drift. If wrong,
  the old criterion stays as it is and the `--trace` criterion here lands beside it.

## Outcomes

- `spec-drift --trace` lists every *joined* criterion — two or more `shall`, exactly one
  `Proof:` line, no mention of a ceiling — from capability specs and change documents,
  each with the file that holds it, and counts them in its summary line.
- `spec-drift --anchors` prints one line counting joined criteria in change documents and
  in capability specs, and its exit status is decided exactly as before. With nothing
  joined, its output is byte-identical to today's.
- An author answers a joined criterion by splitting it, citing a second proof, or naming its
  ceiling — and the line goes away.
- `record-work` and `review-spec` describe the new `--trace` section where they already
  describe the others.
- The rejected heuristic and the numbers that rejected it are recorded in `payload`'s
  *Prior decisions*, so the idea is not re-proposed without new evidence.

## Scope boundaries

**In:**

- `skills/record-work/spec-drift`: `is_joined`, `joined_of`, the count line inside
  `--anchors`, the section and summary count in `--trace`, the self-test cases, the
  header comment
- `skills/record-work/SKILL.md`: one row in the `--trace` table, version 1.5 → 1.6
- `skills/review-spec/SKILL.md`: the sentence naming what `--trace` reports, version
  2.0 → 2.1
- `scripts/check-payload`: two `check_wiring` rows
- `.agents/specs/payload/spec.md`, at landing: the existing criterion "the whole tree is
  traceable, not only what is staged" is **amended, not paralleled**. Its list of what
  `--trace` names gains joined criteria, replaced with the text under *Approach*. No
  second criterion about the map is added beside it.

**Out:**

- **the term-in-test-body heuristic**, in any form. Rejected on the measurement; see
  *Prior decisions*.
- **failing on a joined criterion.** Warn-only, in both scopes.
- **opening, parsing or running the cited test.** The predicate reads the criterion and
  counts its citations, nothing else.
- **enumerated responses** — one `shall` promising a list whose proof asserts part of it.
  Six of the eight ledger cases have that shape, and telling an enumeration from a
  sentence with commas is a reading. `work-reviewer` caught all eight; it stays the
  mechanism.
- **the three live gaps the measurement found**, named so they are not lost:
  `cmd/libretto/init_test.go TestInitRefusesOutsideAGitRepository` never exercises the
  "git is unavailable" half of its trigger (`cli`); `TestWikiPageCarriesTheActivitySparkline`'s
  fixture dates sit inside both an eight-month window ending at the latest commit and one
  ending at today's clock, so "never the clock's" cannot fail today (`cli`); and no
  `spec-drift` self-test case stages nothing for `--retired`, so "When no such deletion is
  staged, it shall report nothing" rests on the gate's own daily runs (`payload`). Each is
  a `/libretto-queue` capture, not this change.
- **migrating the 72 joined capability criteria.** A capability migrates when a delta lands
  on it, the same rule `--ears` set.
- **`AGENTS.md`, `docs/`, the CI workflow.** The gate count and the gate list do not move,
  and the comment beside `--anchors` in `AGENTS.md` names what fails it, which this does not.

## Constraints

- bash, `rg`, `awk` and coreutils only — the shell tooling `AGENTS.md` lists. No dependency
  is added.
- BSD awk as macOS ships it: no `gensub`, no `IGNORECASE`, no `\b`. Case folding through
  `tolower`.
- `set -uo pipefail` and `set -f` hold for the whole script. **No pipe ends in `rg -q`**:
  under pipefail an early-exiting reader SIGPIPEs the writer and a 141 reads as no match —
  the trap `has_line` documents. Every test reads its input whole, `rg … >/dev/null`.
- One criterion grammar: a column-0 `- ` under a level-2 `## Verification criteria`
  heading, continuation lines folded, outside fences — `criteria_of`'s and
  `criteria_without_proof`'s.
- `shall` is counted the way `is_ears` matches it — ` shall ` with `*`, `_` and backticks
  stripped — so the two checks never disagree about what a `shall` is.
- `--anchors`' exit status, and the output of its three existing sections, are unchanged.
  Every existing self-test case keeps its expected value, unedited.
- Output truncated by `printf '%.70s'`, as `--ears` does, never by `cut`, which refuses a
  multibyte character split by the cut under a C locale.

## Prior decisions

- **A cited test's body is not searched for the criterion's terms.** Assumed 2026-10-22,
  attacca run, on a measurement. Prototyped over every criterion in `.agents/specs/*/spec.md`
  — 706 criteria, 719 citations — the trigger-clause variant flagged 15 of 52 candidates
  and the whole-criterion variant 106 of 304; twenty flags hand-checked across both were
  twenty false positives. Terms reach a test through helpers, constants and
  `filepath.Join` segments, a negated trigger is *correctly* absent from the body, and a
  `check-payload` citation has no body. On the ledger's six Go-proved cases, each run at
  its pre-fix commit, the trigger variant caught none for the right reason and the
  whole-criterion variant one. **What would reopen it:** a resolver that follows a test
  into its helpers, which is a Go call graph and a dependency this script refuses — or a
  new measurement beating these numbers.
- **A criterion is joined when it carries two or more `shall`, exactly one `Proof:` line,
  and no occurrence of *ceiling*.** Assumed 2026-10-22, attacca run. Structural on
  purpose: it reads the document and never guesses about the test, so it runs before the
  test exists — at `review-spec`, where splitting a criterion costs a line. If wrong,
  the predicate is one function, `is_joined`.
- **Joined criteria warn and never fail, in change documents and capability specs alike.**
  Assumed 2026-10-22, attacca run. Measured: 36 of 99 criteria in the landed EARS-era
  deltas history recovers are joined; 6 of 26 hand-checked flags were real gaps.
  **Ceiling named:** it saw one of the ledger's eight entries — wiki-flow-board, whose two
  joined criteria were both real misses. Three of the eight predate EARS and carry no
  `shall`; four rest one `shall` on an enumerated response. Enumerated responses are
  invisible to it and stay a reading for `work-reviewer`.
- **`--anchors` carries a count, `--trace` carries the list.** Assumed 2026-10-22, attacca
  run. The gate's output is read on every commit; one line there costs nothing, and a list
  of 72 would be the warning that gets muted. `--trace` is the map `review-spec`'s reader
  already runs first.
- **The escape is a declaration in the criterion — split it, a second `Proof:` line, or
  the word *ceiling* — never a flag.** Assumed 2026-10-22, attacca run, on the precedent
  of `Durable decisions: none`: written by whoever knew, at the moment they knew.

## Task breakdown

1. The predicate and the listing: `is_joined`, `joined_of`, the `--trace` section and its
   summary count, with self-test cases driven through the real functions.
2. The count line inside `--anchors`, after `--retired`, outside the exit status, and the
   header comment that lists what `--anchors` and `--trace` ask.
3. The two skills that describe `--trace`, with a `check_wiring` row each and their version
   bumps.

## Verification criteria

- **Where** a verification criterion carries the EARS marker two or more times — counted as
  `is_ears` matches it — sits above exactly one `Proof:` line, and contains no
  case-insensitive `ceiling`, `is_joined` **shall** return true for it, and false for every
  other criterion. The self-test drives the real `is_joined` through six cases: two markers
  on one proof, two markers with "Ceiling named:", two markers on two proofs, one marker, a
  bolded pair, and `marshalling` beside one real marker.
  Proof: skills/record-work/spec-drift --self-test
- **Where** a change document holds `## Verification criteria`, holding one criterion with
  a `Proof:` and one without, then `## Checklist`, `joined_of` **shall** print one line,
  carrying the proven criterion's `Proof:` count. It prints nothing for the unproven
  criterion, which is `--trace`'s *unproven criterion* already, and nothing for the boxes.
  Proof: skills/record-work/spec-drift --self-test
- **When** `spec-drift --trace` runs in a repository whose capability spec and change
  document each hold one joined criterion, it **shall** list both under the joined-criteria
  section, each with its file.
  Proof: skills/record-work/spec-drift --self-test
- **When** `spec-drift --trace` runs in a repository holding exactly two joined criteria,
  its summary line **shall** end in `, 2 joined criterion(a)`. The self-test asserts that
  exact text, so a missing count and a wrong one both fail it.
  Proof: skills/record-work/spec-drift --self-test
- **When** `spec-drift --anchors` runs in a repository holding a joined criterion whose
  citation resolves, it **shall** print exactly one line counting joined criteria in change
  documents and in capability specs.
  Proof: skills/record-work/spec-drift --self-test
- **While** no criterion in the repository is joined, `spec-drift --anchors` **shall**
  print no joined line, and its output is byte-identical to `anchors`, `ears` and `retired`
  run directly in the same repository with the separators `anchors_and_ears` prints between
  them today. That composition stands in for the version before this change, which the
  self-test cannot run.
  Proof: skills/record-work/spec-drift --self-test
- **If** joined criteria are the only finding, **then** `spec-drift --anchors` **shall**
  exit 0.
  Proof: skills/record-work/spec-drift --self-test
- The `record-work` skill **shall** carry joined criteria as a row of its `--trace` table,
  pinned by a `check_wiring` row on the pattern `^\| \*\*joined criterion\*\* \|`. Proved
  as wiring only. **Ceiling named:** the row proves the table names joined criteria. It
  does not prove the row's description matches what `is_joined` does, which is a reading.
  Proof: scripts/check-payload
- The `review-spec` skill **shall** state that `--trace` lists joined criteria as a prompt
  to "split the criterion or name its ceiling", pinned by a `check_wiring` row on the
  pattern `split the criterion or name its ceiling`. Proved as wiring only. **Ceiling
  named:** the row proves the wording is present, not that a reviewer acts on it.
  Proof: scripts/check-payload

## Approach

### Summary

Two small awk-and-bash functions in `spec-drift`: `joined_of` reads each criterion with the
number of `Proof:` lines beneath it, and `is_joined` decides whether it rests several
`shall` on one proof without naming a ceiling. `--trace` lists those criteria; `--anchors`
prints one line counting them and keeps its exit status. The heuristic the ask named — key
terms looked up in the cited test's body — was prototyped first, measured across the whole
specification, and rejected; the numbers are below and land in `payload`'s *Prior
decisions*.

### Technical context

- `skills/record-work/spec-drift`: bash, `set -uo pipefail`, `set -f`, every question asked
  through `rg` and `awk`. `criteria_of` folds a criterion up to its first `Proof:`;
  `criteria_without_proof` reads the same grammar for `--trace`; `is_ears` is the one
  `shall` matcher, called by the gate and the self-test both. `anchors_and_ears` runs
  `anchors`, `ears` and `retired` and is what `--anchors` means. `trace` reads
  `spec_files` only today.
- The self-test builds throwaway repositories under `mktemp -d` with `HOME` and the git
  config pinned (`land`, `retire`, `tick`), and drives the real functions.
- `scripts/check-payload`'s `check_wiring <file> <pattern> <description>` pins a sentence.
- Blast radius: `skills/record-work/spec-drift`, `skills/record-work/SKILL.md`,
  `skills/review-spec/SKILL.md`, `scripts/check-payload`; at landing, the delta onto
  `.agents/specs/payload/spec.md`.
- Gates: the six in `AGENTS.md`.

**What was measured, 2026-10-22.** A throwaway prototype in bash, `rg` and awk, outside the
repository, read every criterion bullet carrying a `Proof:` in the fourteen capability
specs: 706 criteria and 719 citations — 637 Go tests, 61 `scripts/check-payload`, 16
`spec-drift`, 5 other files.

*The requested heuristic.* Trigger clause: from each `When`/`While`/`If`/`Where` to the next
`shall` or `then`. Key terms: backticked spans, double-quoted strings, `--flags`, with
`<placeholders>` and a trailing `()` removed. Go test body: from `func <Name>(` to the first
column-0 `}`. A flag: no term found in the body (or some missing, "partial").

| Variant | Criteria with terms | Go-proved | Flagged (none + partial) | Of all 706 | Hand-checked | False |
|---|---|---|---|---|---|---|
| trigger clause | 82 | 52 | 7 + 8 = 15 | 2.1% | 10 | 10 |
| whole criterion | 373 | 304 | 47 + 59 = 106 | 15.0% | 10 | 10 |

The trigger-clause sample, in full: `checklistPath` (a function name the test reaches
through the wiki), `## <date> · <change> · <phase>` (a template written by a helper),
`Targets:` twice (in `landFixture`), `libretto queue` (called through `runQueue`), `## Reading`
(a negated trigger — correctly absent), `.agents/` (built from `filepath.Join` segments),
`.agents/changes/` (in `writeProposal`), `plan.md`/`tasks.md` (a negated trigger), and a
regex the test exercises as a table of values. Every one exercised its clause.

On the ledger's six Go-proved cases — add-contributing-guide twice,
split-readme-into-sections, add-specs-wiki, wiki-pages-redesign, wiki-flow-board, each run
against its fix commit's parent — the trigger variant caught none for the right reason
(three of the six predate EARS and carry no trigger at all), and the whole-criterion
variant one (`Governs:`, absent from `TestWikiHTMLIsHomeAndPages`).

*The structural predicate built here.*

| Scope | Criteria | Joined | Hand-checked | Real gaps |
|---|---|---|---|---|
| capability specs today | 706 | 72 (10.2%) | 10 | 2 |
| landed EARS-era change deltas, the ten files history recovers | 99 | 36 (36%) | 10 | 1 |
| the three wiki deltas above, pre-fix | — | 6 | 6 | 3 |

Precision 6 of 26 (23%). The real gaps: the sparkline window, `init`'s "git is
unavailable", the `--retired` "nothing staged" clause, wiki-flow-board's two absence
clauses, and the redesign's "remove the superseded … in the same commit", which one test
cannot hold. No change document in flight today holds a criterion with a `Proof:`, so
`--anchors` would print the capability count alone.

### The approach

1. **`joined_of <file>`** — awk over `defenced`, the criterion grammar `criteria_of` uses:
   a column-0 `- ` opens a criterion, indented lines fold into it until the first
   `Proof:`, every `Proof:` line beneath it before the next column-0 bullet or level-2
   heading increments its count. Prints `<count><TAB><text>` only for a criterion with at
   least one `Proof:`. A criterion with none prints nothing here, because it is already
   `--trace`'s *unproven criterion*, and listing it twice is noise.
2. **`is_joined <text> <count>`** — true when `count` is 1, the text with `*`, `_` and
   backticks stripped holds two or more ` shall ` (counted by awk `gsub`), and
   `tolower(text)` holds no `ceiling`. Called by the gate and the self-test both, as
   `is_ears` is.
3. **`joined count`** — over `change_files` and `spec_files`, counts each scope; prints
   nothing when both are zero, else one line:
   `  N joined criterion(a) — X in change documents, Y in capability specs: two or more shall on one Proof:. spec-drift --trace lists them; split each, cite a second proof, or name its ceiling`.
   Returns 0 always. `anchors_and_ears` calls it after `retired`, without `|| rc=1`. The
   separating blank line is printed by `joined` itself, and only with its line. So with
   nothing joined, `--anchors` writes exactly the bytes `anchors`, `ears` and `retired`
   write today, with today's separators, and not one byte more.
4. **`--trace`** gains a fourth section, `── criteria resting several shall on one proof ──`,
   over `spec_files` and `change_files`, one `joined  <file>` line and a `%.70s` head per
   criterion, `none` when empty; the summary line gains `, %d joined criterion(a)` as its
   last field.
5. **The self-test** gains:
   - six `is_joined` cases
   - a `joined_of` case over a criteria section holding one proven criterion and one
     unproven criterion, followed by a checklist. Exactly one line comes back, carrying the
     proven criterion's count
   - one throwaway repository holding a capability spec and a change document with one
     joined criterion each, citing a `t_test.go` holding `func TestX`. It runs once under
     `--trace`, asserting both files are named and the summary line ends in
     `, 2 joined criterion(a)`, and once under `--anchors`, asserting exit 0, the count
     line present and `1 in change documents, 1 in capability specs`
   - one throwaway repository with the same files, both criteria split so nothing is
     joined. It asserts that `--anchors` output holds no `joined` line and equals,
     byte for byte, `anchors; printf '\n'; ears; printf '\n'; retired` run in the same
     repository
6. **The skills.** `record-work`'s `--trace` table gains the row *joined criterion* and its
   "three questions" becomes four; `review-spec`'s sentence naming what `--trace` reports
   gains joined criteria, as a prompt to "split the criterion or name its ceiling", never a
   finding on its own. `check-payload` gains one `check_wiring` row for each:
   `skills/record-work/SKILL.md` on `^\| \*\*joined criterion\*\* \|`, the table-row
   shape, so prose that only mentions the word cannot satisfy it. And
   `skills/review-spec/SKILL.md` on `split the criterion or name its ceiling`.
7. **The header comment** of `spec-drift` lists the count under `--anchors` and the section
   under `--trace`.
8. **At landing**, the existing `payload` criterion "the whole tree is traceable, not only
   what is staged" is amended in place. Its sentence beginning "`--trace` reads the tracked
   tree" is replaced, whole, by:

   > `--trace` reads the tracked tree and names orphan code, a `Governs:` glob matching
   > nothing, a criterion under *Verification criteria* with no `Proof:` beneath it, and a
   > joined criterion — two or more `shall` resting on one `Proof:` with no ceiling named.

   The rest of that criterion and its `--self-test` citation stay as they are. This
   change's own `--trace` criteria land beside it as the proof of the fourth item, not as
   a second description of the map.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| B — the trigger clause's terms looked up in the cited Go test's body, warn-only | ten of ten sampled flags false; caught none of the ledger's Go cases for the right reason; cannot run before the test exists; skips every prose-proved citation |
| the same over every term in the criterion | ten of ten sampled flags false at a 15% flag rate; one ledger case of six |
| following a test into its helpers | a Go call graph, which is a dependency this bash script refuses |
| failing on joined criteria in change documents, like `--ears` | 36% of EARS-era delta criteria joined, about one in four real |
| per-criterion lines inside `--anchors` | 72 lines on every gate run; the list belongs in the map |
| C — build nothing, record the measurement | the structural check found three live gaps in one run for one function's cost |
| a new `--joined` flag | the list already has a home in `--trace`, and a flag nobody runs finds nothing |

### Risks

| Risk | What catches it |
|---|---|
| the count line flips `--anchors` red | the self-test's `--anchors` case asserts exit 0 with a joined criterion present; `joined` is called without `\|\| rc=1` |
| a pipe into `rg -q` turns a SIGPIPE into "no match" | the constraint above; the new code counts in awk and tests with `>/dev/null` |
| `ceiling` mentioned for another reason silences a real joined criterion | accepted: an author who wrote the word was thinking about what the proof misses. Named in *Prior decisions* by the escape it is |
| `joined_of` and `criteria_of` disagree on a criterion's text | both read through `defenced` with the same bullet rule; the self-test's checklist case covers the boundary that has broken before |
| a `%.70s` cut splits a multibyte character | the same cut `--ears` already prints; cosmetic, never a status |

### Validation and rollback

Every new self-test case runs inside `spec-drift --self-test`, and every existing case
keeps its expected value. Two must be forced red before they are believed:

- the `--anchors` case passes if `joined` is never called — exit 0 is the old behaviour.
  Assert the count line is present, then remove the call from `anchors_and_ears` and watch
  that assertion fail.
- `two shall on one proof is joined` passes if `is_joined` answers yes to everything.
  Change its threshold to one and watch `one shall is not joined` fail.
- the zero case passes if `joined` prints nothing at all, so it is forced red the other
  way: make `joined` print its separator unconditionally and watch the byte comparison
  fail.
- the `--trace` count case: hard-code the count to 0 and watch the `, 2 joined
  criterion(a)` assertion fail.

`scripts/check-payload`'s two rows are forced red by deleting the sentence each pins from a
throwaway copy. Rollback: one revert; nothing migrates.

### Complexity deliberately kept

- `joined_of` beside `criteria_of`, rather than teaching `criteria_of` to count proofs:
  `criteria_of` feeds `--ears`, the hard half of the gate, and a second output column there
  is a change to the gate this change promises not to touch.

## Checklist

Phase 6 runs these with `build-and-check`. Each box leaves the tree green on its own.

- [x] **`--trace` lists joined criteria.** `joined_of`, `is_joined`, the fourth `--trace`
      section and its summary count, over capability specs and change documents; the six
      `is_joined` cases, the `joined_of` proven-and-unproven case and the
      throwaway-repository `--trace` case, with its summary count forced red once. Closes:
      the `is_joined` criterion, the `joined_of` criterion, the `--trace` listing criterion,
      the `--trace` summary-count criterion — all on skills/record-work/spec-drift
      --self-test. Waits on: nothing.
- [x] **`--anchors` counts them and stays green.** `joined count`, called from
      `anchors_and_ears` after `retired` and outside `rc`; the header comment; the
      `--anchors` case in the same throwaway repository, forced red once by removing the
      call; the zero-case repository, compared byte for byte, forced red once by an
      unconditional separator. Closes: the count-line criterion, the zero-case criterion
      and the exit-0 criterion, on skills/record-work/spec-drift --self-test. Waits on:
      box 1.
- [x] **The skills say what the section means.** The `record-work` `--trace` table row and
      "four questions", version 1.6; the `review-spec` sentence, version 2.1; one
      `check_wiring` row each in `scripts/check-payload` — `^\| \*\*joined criterion\*\* \|`
      and `split the criterion or name its ceiling` — each forced red on a throwaway copy.
      Closes: the `record-work` criterion and the `review-spec` criterion, on
      scripts/check-payload. Waits on: box 1, so neither skill names a section that does
      not exist. Independent of box 2.
- [ ] **At landing: the map's criterion names its fourth item.** In
      `.agents/specs/payload/spec.md`, the criterion "the whole tree is traceable, not only
      what is staged" gets its `--trace` sentence replaced by the text in *Approach*, step 8.
      It is amended in place, with no parallel criterion, and lands in the same commit as
      this delta and the *Prior decisions*. Closes: that criterion's amended list, on its
      existing skills/record-work/spec-drift --self-test citation. Waits on: boxes 1–3.
