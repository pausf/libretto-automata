# raise-the-update-notice — plan

Spec: `.agents/changes/raise-the-update-notice/spec.md`

Three tasks, strictly sequential. Task 1 is the only one that can start now: 2 needs the
banner helper to exist, and 3 records what 1 and 2 actually did rather than what they were
meant to do.

## 1 · The banner, and the move

- [ ] `Theme.banner(text, cw)` returns the boxed, captioned rows, and `Render` emits them
      above the logo instead of below the menu — the old band and its separator go
  - From: spec, outcomes + task 1
  - Waits on: nothing
  - Closes: *the banner renders above the logo* · *it is boxed and captioned*
  - Proof: `internal/ui/notice_test.go` `TestUpdateNoticeRendersAboveTheLogo`,
    `TestUpdateNoticeBannerIsBoxedAndCaptioned`
  - Also has to keep passing: `TestPanelOmitsUpdateNoticeWhenEmpty`,
    `internal/ui/fluid_test.go` `TestFrameHoldsWithUpdateNotice`

## 2 · The narrow layout keeps telling the truth

- [ ] `renderNarrow` shows the notice above the mark, still borderless
  - From: spec, task 2
  - Waits on: 1 — the placement rule it mirrors does not exist until then
  - Closes: *the narrow layout keeps it, above the mark and without the box*
  - Proof: `internal/ui/notice_test.go` `TestNarrowLayoutKeepsUpdateNotice`

## 3 · Land the delta

- [ ] Replace the position criterion and its `Proof:` in `.agents/specs/panel/spec.md`,
      fold in the new banner criteria, and delete the change folder
  - From: spec, task 3
  - Waits on: 1 and 2 — a spec describing a placement the code does not have is the drift
    this phase exists to prevent
  - Closes: nothing new. It is what makes the other two true in the one place that lasts
  - Proof: `skills/record-work/spec-drift --anchors` resolves every citation

## Gates

All six, before the commit that carries tasks 1 and 2. Task 3 lands in the same commit as
the final code, per `AGENTS.md`.
