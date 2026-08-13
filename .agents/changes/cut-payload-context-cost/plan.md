# cut-payload-context-cost — plan

Live state. The orchestrator is the only writer; a box is marked the moment its task is
verified, and it ships in the commit that closed it.

Execution is `build-and-check`, phase 6.

Every task names the spec section it derives from and the criterion that closes it.
Spec: `.agents/changes/cut-payload-context-cost/spec.md`.

**Three of the six rules land in one file.** The order below is not convenience: task 1
must be red before tasks 2–4 exist, or the rows prove nothing.

## The red-first task, and why it is first

- [x] **1 · Three `check_wiring` rows, observed RED** — all three failed, exit 1, nothing
      else in the script failed. Recorded:
      `FAIL phase 2 selects the governing spec, not the spec root — … no longer contains /never the corpus/`
      `FAIL a fan-out brief section is named per subtask, vocabulary always — … /always the vocabulary/`
      `FAIL the fan-out states what a mid-phase model switch costs — … /rebills the whole context/`
      Add three rows to `scripts/check-payload`, all over `skills/write-spec/SKILL.md`,
      each matching a phrase the file does **not** yet contain. Run `scripts/check-payload`,
      redirect to a file, check `$?`, read the file. Record all three failures verbatim —
      that recorded output is the evidence, and it can only be collected now.
      *From:* task breakdown 1 · constraints · *Closes:* nothing on its own; it is what
      makes criteria 1–3 capable of failing.
      *Blocks:* 2, 3, 4. Nothing blocks it.

      A row that comes up green here is a row whose pattern already exists in the file.
      Remove it and rewrite the pattern — never keep it and never soften the check.

## The three rules — all in `skills/write-spec/SKILL.md`

These three edit one file, so they are **serial, not parallel**. One writer, and
concurrent edits to one markdown is the lost-update race this plan's own rules forbid.

- [x] **2 · Rule 1 — select, don't load** — row green, rows 3 and 4 still red, nothing
      else in the script disturbed. The duplicate *If the project already has a
      convention* section at the file's tail was collapsed into a pointer, so the layout
      question is answered once rather than twice.
      `write-spec` step 2: the selection instruction (by `Targets:`, `Governs:`, or the
      project's index), naming reading the whole spec root as the thing not to do. Stated
      without naming any repository's layout.
      Same task corrects **A4**: the "If the project already has a convention" paragraph
      claims *"This repository is one."* — false, thirteen capability directories under
      `.agents/specs/`, which `spec-drift` resolves first.
      *From:* outcomes 1 · task breakdown 2 · *Closes:* criterion 1.
      *Waits on:* 1.

- [ ] **3 · Rule 2 — the brief's fixed headings**
      `write-spec` step 2b: the brief's five section headings stated as a **fixed
      enumerated set** (conventions, constraints, settled decisions, vocabulary, the
      six-pillar structure); the sub-agent prompt names the sections its subtask touches
      and **always names the vocabulary**.
      The return contract is **not touched** — `agents/spec-writer.md` already promises it.
      *From:* outcomes 2 · task breakdown 3 · *Closes:* criterion 2.
      *Waits on:* 1. Independent of 2 in content, serial in file.

- [ ] **4 · Rule 3 — cache stability**
      `write-spec` step 2b: a switch of model or effort invalidates the cached prefix and
      rebills the whole context at full input price, stated beside the fan-out where N
      contexts pay it.
      Then `docs/FLOW.md` under *Delegation*, carrying the same reasoning — **uncontracted,
      and no criterion cites it**, because no capability governs `docs/`.
      *From:* outcomes 3 · task breakdown 4 · *Closes:* criterion 3.
      *Waits on:* 1.

## Closing

- [ ] **5 · `scripts/check-payload` green, observed**
      Run it, redirect, check `$?`, read the file. All three rows green, and the rest of
      the script still green. Never piped into `head` — the pipeline reports the last
      command's status.
      *From:* task breakdown 5 · *Waits on:* 2, 3, 4.

- [ ] **6 · All six gates, then apply the delta**
      `gofmt -l .` (nothing), `go vet ./...`, `go test ./... -count=1`,
      `scripts/check-payload`, `spec-drift --self-test`, `spec-drift --anchors`.
      Then fold the three criteria onto `.agents/specs/payload/spec.md`, into the
      flow-wiring group of `## Verification criteria` — **before** the paragraph that
      introduces the bump rows, since that paragraph scopes the bullets after it. Delete
      the change folder in the same commit as the code.
      *From:* task breakdown 6 · constraints · *Waits on:* 5.

      `--anchors` must resolve the three new `Proof:` citations. They cite
      `scripts/check-payload`, which is a file rather than a Go test — the same shape the
      other 18 payload criteria already use.

## What can start now

**Task 1, alone.** Everything else waits on it, by design: the rows have to be capable of
failing before there is prose for them to match.
