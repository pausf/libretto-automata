# Tasks — wiki-home-board

Executed by `build-and-check` (phase 6), one fresh session per open box.
Files: `cmd/libretto/wiki.go`, `cmd/libretto/wiki_test.go`. Branch `feat/add-wiki-html-output`.
The capability spec delta lands separately in the final commit; no box below carries it.

- [x] **1. The home board — extraction, checks, project state, render, motion, with all six proofs. One commit.**
  - Traces: all six verification criteria in `spec.md`; plan steps 1–5. One box because
    `spec-drift --anchors` runs per commit and five of the six cited tests do not exist yet.
  - Contains: `wikiCriterion{text, proofFile, proofTest}`; `parseSpec` keeps the first
    `Proof:` per bullet; seams `wikiGitSubject` and `wikiGitTracked`; pure checks
    `isEARS`, `proofResolves`, `capHealthy`, `governedSplit` (`**` handling mirroring
    spec-drift); `readInFlight` (tasks.md, fallback plan.md, skip zero-open, count
    `Queued:` proposals); render rail/strip/bar/dots/footer with absence arms; motion
    inside the reduced-motion block, `/*end-motion*/` kept. Five new tests against temp
    dirs, no clock, no real ~/.claude.
  - Closes when: all six cited proofs pass; every pre-existing wiki test green
    unchanged; force-red on the `**` arm of `governedSplit` observed and restored;
    six gates green; one commit.
  - Waits on: nothing.
  - Evidence: six gates green on commit "feat(cli): the wiki home becomes the
    project board"; the ** arm of governedSplit forced red, observed failing on
    the footer test's own message, restored.

- [x] **2. Render-and-look against this repository.**
  - Closes when: the generated home observed against this repo — rail, strip, bar,
    dots, footer with plausible values. Findings fixed under the same gates.
  - Waits on: 1.
  - Evidence: rendered via throwaway worktree with real seams — health note reads
    "5% EARS · 0 unproven", which is the HONEST number (545 criteria predate the
    syntax by recorded policy; the mockup's 88% was fiction); rail carries real
    landing subjects; the strip shows this very change 0/2; footer 141/36 of 177.
    No findings; the surprising number is the truth doing its job.
