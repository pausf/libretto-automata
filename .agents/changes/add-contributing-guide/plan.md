# add-contributing-guide — plan

Execution is `build-and-check`, phase 6.

Spec: `.agents/changes/add-contributing-guide/spec.md` (Targets: contributing — new capability)

## Boxes

- [x] **1 · the guard and the guide, together**
      Spec: outcomes 1, 2, 3 and 4.
      Closes: *`CONTRIBUTING.md` exists and names the four contributor-specific things* ·
      *none of the four `AGENTS.md` phrases is restated* · *every relative link resolves*.
      Evidence: `TestContributingIsADoorNotACopy` exit 1 before the file existed —
      `reading CONTRIBUTING.md: no such file or directory`, the honest first red rather than a
      manufactured one — then exit 0. The no-duplication arm proved separately: two forbidden
      phrases planted, exit 1 naming both, exit 0 once removed. It caught `Co-Authored-By`
      **wrapped across a line break**, which is `flat()` doing its job rather than merely being
      present — the thing two earlier guards in this file failed at.
      All six gates green: `gofmt` silent, `vet` silent, `go test ./...` exit 0,
      `check-payload` exit 0, `--self-test` exit 0, `--anchors` 533 citations resolve.
      Waits on: nothing. Can start now.

      **One box.** The guard alone fails; the guide alone can restate `AGENTS.md` freely with
      nothing watching.

- [ ] **2 · the capability and the index**
      Spec: task breakdown 3.
      Closes: nothing new — it lands what box 1 proved.
      Evidence: `.agents/specs/contributing/spec.md` written, `docs/SPEC.md` row added,
      `spec-drift --anchors` green and `--trace` with no dead claim for the new `Governs:`.
      Waits on: box 1. Phase 8, and the commit that deletes this file.

      **`docs/SPEC.md` is the only place the capability list lives**, so the row is not
      optional bookkeeping — a capability absent from the index is a capability the next
      session does not know to read.

2 boxes. Box 1 can start now.
