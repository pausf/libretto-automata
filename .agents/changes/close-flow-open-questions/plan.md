# close-flow-open-questions — plan

Header: executed by `build-and-check` (phase 6). One writer: the orchestrator.

Derived from `spec-payload.md` and `spec-cli.md`. Each task names its delta and the
criterion that closes it.

## Tasks

- [x] 1. `spec-drift --block`: same checks as default, exit 1 when it warned, plus
      usage line and a self-test case. (payload · Proof: spec-drift --self-test)
- [x] 2. record-work SKILL.md: the copy-paste pre-commit snippet using `--block`,
      framed opt-in, default unchanged; wiring row in check-payload. Waits on 1.
      (payload · Proof: scripts/check-payload)
- [ ] 3. build-and-check SKILL.md: the visual-output rule — if judging it means
      looking at it, render and look before the seam — wiring row; and FLOW.md's
      Open section closes the artifact question. Independent.
      (payload · Proof: scripts/check-payload + seam readback for FLOW.md)
- [ ] 4. review-work SKILL.md: append one ledger entry per finding, header contract
      `## <date> · <change> · 6→7`, wiring row. Independent.
      (payload · Proof: scripts/check-payload)
- [ ] 5. write-spec SKILL.md step 4 recut: no hard cap, one call, zero legitimate,
      better-ask-than-silent, never a form; wiring row. Independent.
      (payload · Proof: scripts/check-payload)
- [ ] 6. `corrections()` returns per-phase counts; `6→7` excluded from the
      per-change corrections column. Independent — the `6→7` value is fixed by the
      spec, not by task 4 landing.
      (cli · Proof: TestCorrectionsCountsByPhase, TestReviewerFindingsStayOutOfCorrections)
- [ ] 7. Corpus report renders the phase breakdown with its explanation line;
      absent ledger → not-in-use line; present-but-empty → one-line statement.
      Waits on 6.
      (cli · Proof: TestMetricsReportsCorrectionsByPhase, TestMetricsPhaseBreakdownAbsentLedger)

## Can start now

1, 3, 4, 5, 6. Task 2 waits on 1; task 7 waits on 6.
