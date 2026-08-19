# Tasks — wiki-flow-board

Executed by `build-and-check`, one box per session. The orchestrator owns this file.

- [x] **1. The flow article, end to end with its proofs — one commit.**
  readLessons splits entry headers on ` · ` from the right (the phase is the last
  field); readInFlight grows a queue listing (name+date, date-then-name order);
  the article renders as class `cap` id `flow` (router pages it free) with
  proportional phase bars and queue rows; either source alone summons it; home
  links `#flow` only when it exists. Three tests in the same commit, dotted-name
  fixture forced red under a left split and restored.
  - Waits on: nothing.
  - Evidence: six gates green on commit "feat(cli): the flow board"; left-split
    force-red observed on the dotted-name fixture, restored.

- [x] **2. The look.** Generate against a real ledger+queue project, copy out of
  the worktree before opening (the recorded lesson), observe bars/queue/link in
  both themes; styling fixes only; suite green.
  - Waits on: 1.
  - Evidence: rendered against this repo's real ledger — 6→7 counts 53
    corrections against ones elsewhere, which is the seam doing its job made
    visible; queue empty here (honest); home link present; copied out of the
    worktree before opening, per the recorded lesson.
