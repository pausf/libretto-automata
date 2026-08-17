# add-vertical-slicing-to-plan — plan

Execution is `build-and-check`, phase 6. Nothing else reads this file to start work.

Spec: `.agents/changes/add-vertical-slicing-to-plan/spec.md` (Targets: payload)

## Boxes

- [x] **1 · the mandate and its gate, together**
      Spec: outcomes, all four bullets.
      Closes: *`skills/write-plan/SKILL.md` carries the vertical-slicing mandate* ·
      *no bare `Claude` addressee* · *a host-neutral marker is still present*.
      Evidence: `scripts/check-payload` exit 1 with
      `FAIL  a box is cut to stand alone … no longer contains /one badly cut box/` before the
      skill edit, exit 0 with `ok    a box is cut to stand alone …` after. Observed in that
      order. All six gates green: `gofmt -l .` silent, `go vet` silent, `go test ./...`
      exit 0, `check-payload` exit 0, `spec-drift --self-test` all passed,
      `spec-drift --anchors` 531 citations resolve.
      Waits on: nothing. Can start now.

      **This is one box and not two on purpose.** The skill edit alone leaves a mandate
      nothing gates; the `check_wiring` line alone fails the gate it just added. Neither
      half leaves the tree green, so neither half is a box — which is the rule this change
      adds, applied to its own plan. The frontmatter `version:` bump rides in the same box:
      it is the skill's contract moving, and it is one line.

- [ ] **2 · land the delta**
      Spec: task breakdown 3.
      Closes: *every `Proof:` citation in this delta resolves once it lands on `payload`*.
      Evidence: `skills/record-work/spec-drift --anchors` green with the change folder gone.
      Waits on: box 1. This is phase 8's work and lands in the same commit as box 1's code,
      per `payload` — one delta landing, never one per box.

2 boxes. Box 1 can start now.
