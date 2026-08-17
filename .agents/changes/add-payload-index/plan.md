# add-payload-index — plan

Execution is `build-and-check`, phase 6.

Spec: `.agents/changes/add-payload-index/spec.md` (Targets: payload)

## Boxes

- [x] **1 · the staleness gate and the generator, together**
      Spec: outcomes 1, 2 and 4.
      Closes: *`check-payload` fails when `docs/PAYLOAD.md` does not match a fresh generation* ·
      *`--index` regenerates it such that the default run then passes* · *every item under the
      three directories appears in the page*.
      Evidence: gate red with no page — `docs/PAYLOAD.md does not exist — run
      scripts/check-payload --index`, exit 1 — then green after `--index`, with 36 rows against
      22 + 7 + 7 counted from the directories. **Then watched biting three ways**, because a
      staleness gate that has only run against a fresh page has proved nothing: a hand-edited
      description, an added probe skill, and a removed row. All three gave
      `docs/PAYLOAD.md has drifted from the payload`, and the tree was restored to exit 0 after
      each.
      All six gates green: `gofmt` silent, `vet` silent, `go test ./...` exit 0,
      `check-payload` exit 0, `--self-test` exit 0, `--anchors` 537 citations resolve.

      **Scope added deliberately, not silently.** The new page exposed a pre-existing defect in
      `check-payload`'s referenced-paths check: it piped `rg -oN` (which prefixes `file:`) into a
      second `rg` that re-extracted from that prefix, so any file under `docs/` or `scripts/`
      verified *itself* as a referenced path. `docs/PAYLOAD.md` and `docs/PLAN.md` were both
      reported `ok` while nothing referenced them. Fixed with `-I`; noise rather than a false
      pass, since a genuine reference is a text match either way.
      Waits on: nothing. Can start now.

      **One box.** The gate alone fails the whole suite; the generator alone produces a page
      nothing keeps current — which is the typed list this change exists to avoid, wearing a
      script's clothes.

- [ ] **2 · land the delta**
      Spec: task breakdown 3.
      Closes: *every `Proof:` citation resolves once the delta lands on `payload`*.
      Evidence: `spec-drift --anchors` green with the change folder gone.
      Waits on: box 1. Phase 8, same commit as the landing.

2 boxes. Box 1 can start now.
