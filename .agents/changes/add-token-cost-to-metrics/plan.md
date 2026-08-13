# add-token-cost-to-metrics — plan

Live state. The orchestrator is the only writer; a box is marked the moment its task is
verified, and it ships in the commit that closed it.

Execution is `build-and-check`, phase 6. Every task names the spec section it derives from
and the criterion that closes it.
Spec: `.agents/changes/add-token-cost-to-metrics/spec.md`.

**This one is real Go, so the discipline is red-first per task, not once at the front.**
The shapes are all known from measured fixtures, so every test below can be written and
watched fail before its production code exists. A test that was green on its first run is
rewritten, never kept.

## Foundation — nothing renders until this parses

- [ ] **1 · The fixture root, and the streaming parse**
      A `t.TempDir()` transcript root with a hand-written `.jsonl` carrying real entry
      shapes: an `assistant` entry with the full `usage` object, a `user` entry, a `mode`
      entry with no `cwd`/`gitBranch`/`timestamp`. Then `cmd/libretto/usage.go`: read line
      by line, decode into a struct where **every field is optional**, sum the four
      numbers from `assistant` entries only, ignore `iterations[]`.
      *From:* outcomes 1 · constraints · *Closes:* criterion 1.
      *Blocks:* everything. Nothing blocks it.
      Proof: `TestTheFourUsageNumbersAreKeptApart`

- [ ] **2 · Everything that is absent, malformed or strange**
      A line that is not JSON, an `assistant` entry with no `usage`, an entry with no
      `gitBranch`, a `<synthetic>` model with an all-zero usage object and null
      `service_tier`. None fatal, each counted where outcome 1 says.
      *From:* constraints · *Closes:* criterion 2. *Waits on:* 1.
      Proof: `TestAMalformedLineDoesNotCostTheRestOfTheFile`

- [ ] **3 · Discovery, including the subagents**
      Encode the repository root forward into the project directory name — every `/` to
      `-` — and **never invert it**. Find the top-level `*.jsonl` and the
      `*/subagents/agent-*.jsonl` beneath it.
      The fixture carries both, with different numbers, so a reader that walked only the
      top level fails rather than under-reporting quietly.
      *From:* outcomes 5 · scope boundaries · *Closes:* criterion 3. *Waits on:* 1.
      Proof: `TestSubagentTranscriptsAreCounted`

- [ ] **4 · The transcript root is never written**
      Snapshot the fixture tree — every path, its size, the SHA-256 of its contents —
      before the read and compare after. Red on a create, a delete, a truncation and an
      in-place rewrite alike.
      *From:* verification criteria · *Closes:* criterion 7. *Waits on:* 3.
      Proof: `TestTheTranscriptRootIsNeverWritten`

## Attribution — where the honesty lives

- [ ] **5 · Per entry, never per file**
      One fixture file whose entries name two different branches. Each entry attributes to
      its own — a per-file reading misattributes everything after a checkout, and that was
      measured on a real session spanning four branches.
      *From:* prior decisions · *Closes:* criterion 4. *Waits on:* 1.
      Proof: `TestBranchIsReadPerEntryNotPerFile`

- [ ] **6 · Prefix stripped, name matched whole**
      `feat/`, `fix/`, `docs/`, `chore/`, `refactor/` stripped; the rest matched **whole**
      against the change names git has seen. The negative case is the point:
      `feat/add-thing-extra` must **not** attribute to `add-thing`.
      *From:* outcomes 2 · A3 · *Closes:* criterion 5. *Waits on:* 5.
      Proof: `TestAPrefixIsStrippedAndTheNameMatchedWhole`

- [ ] **7 · The unattributed bucket, and the invariant**
      `main`, `HEAD`, a branch matching no change, an absent field — all into one bucket.
      The assertion that matters: **attributed + unattributed = corpus**. That invariant is
      independent of rendering and is what stops a quiet drop.
      *From:* outcomes 2 · *Closes:* criterion 6. *Waits on:* 6.
      Proof: `TestUnattributedTokensAreReportedNotDiscarded`

## Rendering — pinned in the spec, so build to the worked example

- [ ] **8 · The corpus block, and the filter invariant**
      The token block after the existing footer, before `flowLegend`. **Corpus-wide under
      both commands** — the totals do not move when a change filter is applied, or the
      invariant from task 7 stops being readable off the output.
      *From:* what each command prints · *Closes:* criterion 10. *Waits on:* 7.
      Proof: `TestTheTokenFooterIsCorpusWideUnderAFilter`

- [ ] **9 · The per-change block, the phases, and the dash**
      Under `metrics <change>`: that change's four numbers, then one row per
      `attributionSkill` value, then an **unattributed row** that is never distributed
      across the phases that did name a skill.
      The dash lives here and nowhere else — three states, per outcome 3.
      *From:* outcomes 3, 4 · *Closes:* criteria 8 and 9. *Waits on:* 8.
      Proof: `TestAChangeWithNoTokensReportsADashNotAZero`,
      `TestPerPhaseCostCarriesAnUnattributedRow`

- [ ] **10 · No root is a state, not an error**
      `CLAUDE_HOME` at an empty temp dir, and at one with no directory for this repository.
      The git-derived report prints in full either way; the token block becomes one line
      saying the measurement was unavailable.
      *From:* constraints · *Closes:* criterion 11. *Waits on:* 8.
      Proof: `TestNoTranscriptRootStillReportsTheGitMetrics`

- [ ] **11 · Wire it into `main.go`**
      `filepath.Join(target.NewClaude().Root(), "projects")`, passed in as a parameter.
      **Injected, never resolved inside the logic** — the same seam `execGit` already uses.
      No file under `internal/target/` is touched; `Root()` is already exported and that
      is what keeps `Targets: cli` correct.
      *From:* scope boundaries · constraints · *Waits on:* 8.

## Closing

- [ ] **12 · The ceiling stops being half-false**
      `flowCeiling` says per-phase measurement "needs a phase to write them down". True for
      duration, false for cost since `attributionSkill` started being recorded. Split the
      sentence.
      The new test asserts **both directions** — the cost sentence present, the old
      undifferentiated claim **absent**. The absence half is what makes it red before the
      work; `TestTheReportNamesWhatItCannotMeasure` checks only that three substrings are
      present and would pass whatever is done here, including nothing.
      *From:* outcomes 6 · A4 · *Closes:* criterion 12. *Waits on:* 9.
      Proof: `TestTheCeilingSeparatesCostFromDuration`

- [ ] **13 · Six gates, then apply the delta**
      `gofmt -l .`, `go vet ./...`, `go test ./... -count=1`, `scripts/check-payload`,
      `spec-drift --self-test`, `spec-drift --anchors`.
      Then fold the twelve criteria onto `.agents/specs/cli/spec.md`, into its `metrics`
      section, and delete the change folder in the same commit as the code.
      `--anchors` must resolve every new `Proof:` — **these cite Go tests by name**, unlike
      the last change, so an invented or renamed test fails the gate rather than passing a
      file-level check.
      *From:* task breakdown 7 · *Waits on:* 12.

## What can start now

**Task 1, alone.** Nothing renders until something parses, and nothing is attributed until
something is read.
