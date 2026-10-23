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
  (fixed in v0.37.0). `check-payload`'s marker search has the same shape but is not
  vulnerable, measured: it squashes newlines first, so `rg` reads the one line whole before it
  can match. It is rewritten anyway, so the new row holds with no exception;
- a `release:` label added after `gh pr create` failed the label check's `opened` run and
  emailed the user a false alarm, three times.

## Decision log

### Session 2026-10-23
- Q: fail or warn on an unclosed fence? → A: fail — a fence never closed hides every citation after it, so the gate it lives in is green for the wrong reason; there is no false positive to fear, an odd count is always a defect (assumed; if wrong, a warning line like joined)
- Q: where does the fence check live? → A: inside `--anchors`, over the same files it reads — never a seventh gate (assumed; the count six is a standing decision)
- Q: what counts as a pipe into `rg -q`? → A: a `|` that is not part of `||`, followed by `rg` with `q` among its short flags, in any file under `skills/` or `scripts/` (assumed; if wrong, the pattern widens or narrows, nothing else moves)
- Q: the label rule — attended runs too? → A: wherever the bump is decided before the request exists, it rides `gh pr create --label`; the bump question itself stays where it is (assumed)

### Session 2026-10-23, after review-spec
- Q: how does `fences` get the fence pattern without copying it? → A: one `FENCE` variable holds the awk pattern; `defenced` and `fences` both read it (assumed)
- Q: what is a pipe into `rg -q`, exactly? → A: on one line, a `|` that is not part of `||`, then `rg`, then — anywhere among its arguments before the next `|` — a short-flag group containing `q`, or `--quiet` (assumed)
- Q: which side of the label rule are the three alarms on? → A: all three were requests this session opened with the bump already decided (the user's standing authorisation for the queue run), then labelled in a second command. attacca with no decided bump is the other side: its request is opened unlabelled and the red check *is* the bump question, as payload states — unchanged (assumed)
- Q: where does `fences` sit in `anchors_and_ears`, and what does it print? → A: after `retired`, before `joined`; like `joined` it prints its own leading blank line only together with its lines, so with every fence closed the output is byte-identical to today's (assumed)
- Q: file, or file and line? → A: the `rg -q` row names file and line (`rg -n` has both); the fence check names the file and its fence-line count, since an unclosed fence has no single line to blame (assumed)

## Outcomes

1. `spec-drift --anchors` fails, naming the file, when a capability spec or change document
   holds an odd number of lines that open or close a fence — the same lines `defenced`
   toggles on.
2. With every fence closed, `--anchors` prints exactly what it printed before.
3. `scripts/check-payload` fails, naming the file and line, when a file under `skills/` or
   `scripts/` pipes into `rg` with a short-flag group containing `q`, or `--quiet`.
4. `check-payload`'s marker search no longer pipes into `rg -q`. Measured: it was not
   vulnerable (one squashed line is read whole), so this is the row holding without an
   exception, not a fix.
5. `record-work` says to pass the `release:` label to `gh pr create` whenever the bump is
   already decided before the request exists — the user named it, or authorised it in
   advance — and why. attacca with no decided bump is untouched: its red check is still the
   question.

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
- The fence definition lives once, in a `FENCE` variable `defenced` and `fences` both read.
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
  change document (fixed v0.37.0). Measured 2026-10-23: 40,000 lines through `| rg -q` exit
  141, the same lines through `rg … >/dev/null` exit 0, and a newline-squashed single line
  exits 0 either way — the reader must finish the line before it can match. Read the
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
- **If** a file under `skills/` or `scripts/` pipes into `rg` with a short-flag group
  containing `q`, or with `--quiet`, **then** `scripts/check-payload` **shall** fail naming
  that file and line. This also closes Outcome 4: the marker search's old pipe is such a line.
  Proof: scripts/check-payload
- **Where** the bump is decided before the request exists, `record-work` **shall** say to
  pass the label to `gh pr create`. Proved as wiring only, pinned on the phrase `the label
  rides gh pr create`. Ceiling named: the row pins the phrase, not that a session follows it.
  Proof: scripts/check-payload

## Approach

1. `spec-drift`: the awk fence pattern moves into one `FENCE` variable, read by `defenced`
   and by a new `fences` function. `fences` runs in `anchors_and_ears` after `retired` and
   before `joined`, with `|| rc=1`. Per file in `spec_files` and `change_files` it counts the
   fence lines; for each odd count it prints, after one leading blank line printed only with
   them, `  UNCLOSED  <file> — <n> fence line(s), one never closed`. With every count even it
   prints nothing. Two self-test cases on a temporary repo: one odd (exit non-zero, the file
   named); one even, compared byte for byte against `anchors; printf '\n'; ears; printf
   '\n'; retired` — the composition the existing zero-joined case already uses.
2. `check-payload`: line 576 becomes `tr '\n' ' ' < "$f" | rg "…" >/dev/null`. A new row
   runs `rg -n` over `skills/` and `scripts/` for a non-`||` pipe into `rg` whose arguments,
   before any further `|`, hold a short-flag group with `q` or `--quiet`; the pattern is built
   from parts so the row's own source line does not match it, and it fails naming each
   `file:line`. Forced red on a scratch copy holding `x | rg -N -q y` and `x | rg --quiet y`.
3. `record-work`: one paragraph beside the forge table, and `version` 1.6 → 1.7.

| Considered | Why it lost |
|---|---|
| a warning for the fence | an unclosed fence makes the gate it sits in lie; warning about a lie is still a green gate |
| a shellcheck dependency | a new dependency for two patterns; the ladder says no |
| fixing only the instances | the class came back twice in two days; a row is cheaper than a third time |

Risks: the `rg -q` row's pattern could miss a multi-line pipe — accepted, none exists, and
the row names the shape it pins.

## Checklist

- [x] **`--anchors` fails on an unclosed fence.** `fences` in spec-drift, two self-test
  cases, forced red by removing the call. Closes: the two fence criteria. Waits on: nothing.
- [x] **No pipe into `rg -q` in the payload's scripts.** The row, line 576 fixed, the row
  observed red on today's tree first. Closes: the `rg -q` criterion. Waits on: nothing.
- [ ] **The label rides `gh pr create`.** The paragraph, the row, record-work 1.7. Closes:
  the label criterion. Waits on: nothing.
