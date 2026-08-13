# add-flow-retro — plan

> **For agentic workers:** this plan is live state, read and executed by
> `build-and-check` (phase 6). One writer: the orchestrator marks boxes; sub-agents
> report. Boxes use `- [ ]` syntax.

**Goal:** the flow captures user corrections into `.agents/lessons.md`, a retro
command spends them, and `libretto metrics` counts them per change.

**Architecture:** capture is one rule in `evidence` (the standing skill every phase
invokes); the ledger is append-only markdown with a fixed countable header; `retro`
is a new payload skill routed by `commands/libretto-retro.md`; the corrections
column in metrics parses the header only. Payload writes, delivery reads.

**Specs:** `spec-payload.md` (Targets: payload) · `spec-cli.md` (Targets: cli), both
in this folder.

## Global constraints

- Ledger path: `.agents/lessons.md`, project-relative. Append-only; entries gain a
  `Resolved:` line, are never edited otherwise.
- Header format, load-bearing for both capabilities:
  `## <date> · <change> · <phase>` — starts with `## `, exactly two ` · `
  separators, three non-empty fields, date not validated. Change field is `-` when
  no change is open.
- A correction is the user saying work already done was wrong. A changed ask is not
  a lesson.
- The retro proposes payload diffs, never applies them.
- No skill references `scripts/` or `docs/`; `name:` frontmatter matches the
  directory; commands route and never implement.
- Gates before every commit: `gofmt -l .` (empty), `go vet ./...`,
  `go test ./... -count=1`, `scripts/check-payload`,
  `skills/record-work/spec-drift --self-test`, `spec-drift --anchors`.
  Never pipe a gate into `head`.

---

### Task 1 — the capture rule in `evidence`

Spec: spec-payload.md. Closes: "the capture rule's decisive words are in
`skills/evidence/`".
Waits on: nothing.

**Files:** Modify: `skills/evidence/SKILL.md`

- [x] **1.1** Add one section to `skills/evidence/SKILL.md`: on a user correction
      (work already done was wrong — a changed ask is not a lesson), append to
      `.agents/lessons.md` one entry — header `## <date> · <change or -> · <phase
      skill>`, a `Said:` line with the user's words, a `Did:` line with what the
      flow had done — then carry on. Never interrupt, never ask, never edit or
      delete existing entries. State the format once, as the contract both the
      retro and `libretto metrics` parse.
- [x] **1.2** Run `scripts/check-payload`; expect pass (no new references yet).
- [x] **1.3** Commit: `feat(payload): evidence captures user corrections into the lessons ledger`

### Task 2 — the `retro` skill

Spec: spec-payload.md. Closes: "the `retro` skill parses, is named correctly, and
references only what installs".
Waits on: Task 1 (the format it reads is stated there).

**Files:** Create: `skills/retro/SKILL.md`

- [x] **2.1** Write `skills/retro/SKILL.md` — frontmatter `name: retro`, a
      description that triggers on "retro", "lecciones", "learn from the flow".
      Body: read `.agents/lessons.md`; take entries without a `Resolved:` line;
      classify each as **project knowledge** (record the convention in the owning
      capability spec's prior decisions or `AGENTS.md`, mark
      `Resolved: <date> · <what was done>`), **flow defect** (name the payload
      skill, propose the exact diff in the report, apply nothing, mark resolved
      only when the user says what they did with it), or **one-off** (mark resolved
      with that reading). Report: entries found, classification, what was written
      where, every proposed diff. No ledger or no open entries: say so in one line
      and stop. Never edit entries beyond appending `Resolved:`; never push; never
      touch the payload.
- [x] **2.2** Run `scripts/check-payload`; expect pass.
- [x] **2.3** Commit: `feat(payload): retro skill — classify lessons, fix the project, propose to the payload`

### Task 3 — the `libretto-retro` command

Spec: spec-payload.md. Closes: "`libretto-retro` routes to `retro` and the
reference resolves".
Waits on: Task 2.

**Files:** Create: `commands/libretto-retro.md`

- [x] **3.1** Write `commands/libretto-retro.md` in the shape of
      `commands/libretto-status.md`: invoke `Skill(skill="retro")`, describe
      nothing the skill already says, carry the one-line why (spend the ledger so
      the same correction is not paid twice).
- [x] **3.2** Run `scripts/check-payload`; expect pass.
- [x] **3.3** Commit: `feat(payload): libretto-retro command routes to the retro skill`

### Task 4 — wiring check in `check-payload`

Spec: spec-payload.md. Closes the same three criteria mechanically — the decisive
words stay where they live.
Waits on: Tasks 1–3 (it checks their artifacts).

**Files:** Modify: `scripts/check-payload`

- [x] **4.1** Extend `scripts/check-payload` following its existing decisive-words
      pattern (the one that guards `review-spec`'s wiring): `skills/evidence/SKILL.md`
      contains `lessons.md` and `Said:`; `skills/retro/SKILL.md` exists and contains
      `Resolved:`; `commands/libretto-retro.md` references the `retro` skill.
- [x] **4.2** Break one of the three by hand in the working tree, run the script,
      expect fail; restore, expect pass. (Observed, not assumed.)
- [x] **4.3** Commit: `test(payload): check-payload guards the retro wiring`

### Task 5 — the corrections column in `libretto metrics`

Spec: spec-cli.md. Closes all three cli criteria.
Waits on: Task 1 only (the header format). Independent of Tasks 2–4.

**Files:** Modify: `cmd/libretto/metrics.go`, `cmd/libretto/metrics_test.go`

- [x] **5.1** Write the three failing tests in `cmd/libretto/metrics_test.go`,
      fixture style as the existing ones (fixed ledger content, no real repos):
      `TestCorrectionsAreCountedPerChange` — a ledger with two entries for change A,
      one for B, expects 2 and 1 in their rows;
      `TestNoLedgerReportsADashNotAZero` — absent file, every row shows `-`;
      `TestMalformedAndChangelessEntriesDoNotCrashTheCount` — a header with one
      separator and a header with change `-`: the first is skipped, the second is
      counted nowhere and named in the report's one-line note.
- [x] **5.2** Run `go test ./cmd/libretto/ -run 'Corrections|NoLedger|Malformed' -count=1`;
      expect FAIL (functions undefined / column absent).
- [x] **5.3** Implement: parse `.agents/lessons.md` (match: `## ` prefix, exactly
      two ` · `, three non-empty fields), count per change, `-` column when the
      file is absent, the note for `-`-change entries, column wired into the
      report next to the existing ones.
- [x] **5.4** Run the same tests; expect PASS. Then the full gates.
- [x] **5.5** Commit: `feat(cli): metrics counts corrections per change from the lessons ledger`

### Task 6 — land the change

Waits on: Tasks 1–5.

- [ ] **6.1** Apply both deltas onto `.agents/specs/payload/spec.md` and
      `.agents/specs/cli/spec.md` (outcomes, scope, prior decisions, criteria with
      their proofs), delete `.agents/changes/add-flow-retro/`, same commit as the
      final reconciliation. All six gates green — `--anchors` now resolves the
      three new test names for real.
