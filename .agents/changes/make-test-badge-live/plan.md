# make-test-badge-live — plan

Execution is `build-and-check`, phase 6.

Spec: `.agents/changes/make-test-badge-live/spec.md` (Targets: readme)

## Boxes

- [x] **1 · the guard and the honest badge, together**
      Spec: outcomes 1 and 2.
      Closes: *no badge whose image URL is a `shields.io/badge/` literal contains a status
      word, word-bounded* · *the tests badge points at the `gates.yml` badge endpoint* ·
      *every relative link still resolves*.
      Evidence: `TestNoBadgeAssertsAStatus` **red first**, exit 1 on both halves —
      `the badge https://img.shields.io/badge/tests-passing-brightgreen.svg hardcodes
      "passing"` and `the README does not carry actions/workflows/gates.yml/badge.svg` —
      then green after the badge line changed. Red before green is the only run that proves
      the criterion describes this bug and not a neighbouring one, and this is a bug, so the
      order is fixed rather than preferred. All six gates green: `gofmt` silent, `vet`
      silent, `go test ./...` exit 0, `check-payload` exit 0, `--self-test` exit 0,
      `--anchors` 531 citations resolve.
      Waits on: nothing. Can start now.

      **One box, not two.** The test alone leaves the suite red; the badge fix alone leaves
      the hole open for the next false badge. Neither half merges on its own.

- [ ] **2 · land the delta**
      Spec: task breakdown.
      Closes: *every `Proof:` citation resolves once the delta lands on `readme`*.
      Evidence: `spec-drift --anchors` green with the change folder gone.
      Waits on: box 1. Phase 8, same commit as the landing.

2 boxes. Box 1 can start now.
