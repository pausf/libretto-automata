# gate-recurring-defect-classes

Tracker: none
Targets: payload
Durable decisions: the two in Prior decisions below

## The ask, verbatim

> QUE MEJORARIAS? … has todo

(2026-10-23. The improvements proposed after draining the queue; items 1 and 3 of that list.)

## Reading

Three defect classes reached this repository once each in the queue run and were caught by
hand. Each becomes a check, so the class stops arriving instead of being caught again:

- a fence opened in prose and never closed hid 157 `Proof:` citations from `--anchors`;
- a pipe into `rg -q` under `pipefail` read a SIGPIPE as "no match" — live in `--retired`
  (fixed in v0.37.0), latent in `check-payload`'s marker search (fixed here);
- a `release:` label added after `gh pr create` failed the label check's `opened` run and
  emailed the user a false alarm, three times.

## Decision log

### Session 2026-10-23
- Q: fail or warn on an unclosed fence? → A: fail — a fence never closed hides every citation after it, so the gate it lives in is green for the wrong reason; there is no false positive to fear, an odd count is always a defect (assumed; if wrong, a warning line like joined)
- Q: where does the fence check live? → A: inside `--anchors`, over the same files it reads — never a seventh gate (assumed; the count six is a standing decision)
- Q: what counts as a pipe into `rg -q`? → A: a `|` that is not part of `||`, followed by `rg` with `q` among its short flags, in any file under `skills/` or `scripts/` (assumed; if wrong, the pattern widens or narrows, nothing else moves)
- Q: the label rule — attended runs too? → A: wherever the bump is decided before the request exists, it rides `gh pr create --label`; the bump question itself stays where it is (assumed)

## Outcomes

1. `spec-drift --anchors` fails, naming the file, when a capability spec or change document
   holds an odd number of lines that open or close a fence — the same lines `defenced`
   toggles on.
2. With every fence closed, `--anchors` prints exactly what it printed before.
3. `scripts/check-payload` fails, naming the file and line, when a file under `skills/` or
   `scripts/` pipes into `rg` with `-q`.
4. `check-payload`'s marker search reads each file whole, so a skill past the pipe buffer
   cannot read as non-compliant.
5. `record-work` says to pass the `release:` label to `gh pr create` whenever the bump is
   already decided, and why.

## Scope boundaries

**In:** `skills/record-work/spec-drift` (the fence check, its self-test cases),
`scripts/check-payload` (two rows and the line-576 fix), `skills/record-work/SKILL.md` (one
paragraph, version bump).

**Out:**
- checking fences in skills or commands — `--anchors` reads specs and change documents, and
  a skill's fence is prose a reader sees rendered.
- `~~~` fences: `defenced` does not toggle on them, and one definition of a fence stays one.
- a generic shell linter. Two patterns that bit, not a style guide.
- relabelling past requests.

## Constraints

- No new dependency; bash, `rg`, `awk`, as the scripts already use.
- The fence definition is `defenced`'s regex, reused, not copied.
- Every existing self-test case and `check-payload` row stays unedited.
- The new check itself never pipes into `rg -q`.

## Prior decisions

- **An unclosed fence fails `--anchors`.** Decided 2026-10-23 on the queue run's evidence: a
  criterion that wrote three backticks in prose opened a fence `defenced` never closed, and
  157 `Proof:` citations after it went unchecked while the gate stayed green. An odd count is
  always a defect, so the check has no false positive to warn about.
- **Nothing in the payload's scripts pipes into `rg -q`.** Decided 2026-10-23: under
  `set -o pipefail`, `rg -q` stops reading at its first match, the writer before it dies of
  SIGPIPE, and the pipeline's 141 reads as no match. It made `--retired` fail open on a long
  change document (fixed v0.37.0) and sat latent in `check-payload`'s marker search. Read the
  input whole — `rg … >/dev/null` — and a `check-payload` row keeps it that way.

## Task breakdown

- the fence check in `--anchors`, with self-test cases;
- the `rg -q` row and the line-576 fix in `check-payload`;
- the label paragraph in `record-work`, with its row.

## Verification criteria

- **If** a capability spec or change document holds an odd number of lines matching
  `defenced`'s fence pattern, **then** `spec-drift --anchors` **shall** exit non-zero naming
  that file.
  Proof: skills/record-work/spec-drift --self-test
- **While** every capability spec and change document holds an even number of fence lines,
  `spec-drift --anchors` **shall** print no fence line.
  Proof: skills/record-work/spec-drift --self-test
- **If** a file under `skills/` or `scripts/` pipes into `rg` with `-q`, **then**
  `scripts/check-payload` **shall** fail naming that file.
  Proof: scripts/check-payload
- **Where** the bump is decided before the request exists, `record-work` **shall** say to
  pass the label to `gh pr create`. Proved as wiring only. Ceiling named: the row pins the
  phrase, not that a session follows it.
  Proof: scripts/check-payload

## Approach

1. `spec-drift`: a `fences` function next to `retired`, called from `anchors_and_ears` with
   `|| rc=1`. It counts, per file in `spec_files` and `change_files`, the lines matching
   `^[[:space:]]*` followed by three backticks — the regex `defenced` toggles on — with
   `awk`, and prints `  UNCLOSED  <file> — <n> fence line(s), one never closed` per odd file
   and nothing at all when every count is even. Two self-test cases on a temporary repo: one
   odd, one even (byte-identical output to the three sections alone).
2. `check-payload`: line 576 becomes `tr '\n' ' ' < "$f" | rg "…" >/dev/null`. A new row
   scans `skills/` and `scripts/` for `[^|]\|[[:space:]]*rg[[:space:]]+-[A-Za-z]*q`, written so
   the row's own source line does not match it; fails naming each hit.
3. `record-work`: one paragraph beside the forge table, and `version` 1.6 → 1.7.

| Considered | Why it lost |
|---|---|
| a warning for the fence | an unclosed fence makes the gate it sits in lie; warning about a lie is still a green gate |
| a shellcheck dependency | a new dependency for two patterns; the ladder says no |
| fixing only the instances | the class came back twice in two days; a row is cheaper than a third time |

Risks: the `rg -q` row's pattern could miss a multi-line pipe — accepted, none exists, and
the row names the shape it pins.

## Checklist

- [ ] **`--anchors` fails on an unclosed fence.** `fences` in spec-drift, two self-test
  cases, forced red by removing the call. Closes: the two fence criteria. Waits on: nothing.
- [ ] **No pipe into `rg -q` in the payload's scripts.** The row, line 576 fixed, the row
  observed red on today's tree first. Closes: the `rg -q` criterion. Waits on: nothing.
- [ ] **The label rides `gh pr create`.** The paragraph, the row, record-work 1.7. Closes:
  the label criterion. Waits on: nothing.
