# consolidate-license-files — plan

Execution is `build-and-check`, phase 6.

Spec: `.agents/changes/consolidate-license-files/spec.md` (Targets: payload)

## Boxes

- [x] **1 · the guard, the move and the links, together**
      Spec: outcomes 1, 2 and 4.
      Closes: *every relative link in `THIRD-PARTY.md` resolves* · *the three vendored licence
      files exist under `licenses/` and no `LICENSE-*` remains at the root* ·
      *`scripts/check-payload` still derives the vendored skill list*.
      Evidence: **both halves watched red independently.** Before the move, exit 1 on layout —
      `licenses/LICENSE-caveman does not exist` and `LICENSE-caveman is still at the root`,
      three times over. After the move with the links untouched, exit 1 on the links —
      `THIRD-PARTY.md links to LICENSE-superpowers, which does not exist`, three times. Green
      once the links were updated. `scripts/check-payload` exit 0 **and its output read**: the
      table parse still returns all seven vendored skills, counted rather than assumed, because
      the gate passes silently on a parse returning everything.
      All six gates green: `gofmt` silent, `vet` silent, `go test ./...` exit 0,
      `check-payload` exit 0, `--self-test` exit 0, `--anchors` all resolve.
      Waits on: nothing.

      **One box.** The guard alone is green against the tree as it stands, so it proves
      nothing; the move alone leaves dead links with nothing watching.

- [ ] **2 · widen ownership and land the delta**
      Spec: outcome 3, and the task breakdown.
      Closes: nothing mechanically — outcome 3 carries no `Proof:` on purpose, because
      `--trace` returns 0 whatever it finds.
      Evidence: the observation. `spec-drift --trace` after the widening: **0 dead claims, and
      orphans 20 → 16** — the two new globs claim four real files. `--anchors` green with the
      change folder gone.
      Waits on: box 1. Phase 8, and it is the commit that deletes this file.

2 boxes.
