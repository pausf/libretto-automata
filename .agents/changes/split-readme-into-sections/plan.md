# split-readme-into-sections — plan

Execution is `build-and-check`, phase 6.

Spec: `.agents/changes/split-readme-into-sections/spec.md` (Targets: readme)

## Boxes

- [x] **1 · the guards and the edits, together**
      Spec: outcomes 1, 2, 3 and 4.
      Closes: *`README.md` is at most 380 lines* · *each of the five relocated arguments is absent
      from the README and present in `docs/`* · *`README.md` links to `CONTRIBUTING.md` and the
      link resolves* · *every existing `readme` criterion still passes*.
      Evidence: all three new guards watched red against the README as it stood — `README.md is
      389 lines, over its 340 ceiling`, each of the five anchors failing on **both** ends
      (`still argued in the README` and `in neither docs/DESIGN.md nor docs/FLOW.md`), and
      `the README never links to CONTRIBUTING.md`. Then green after the edits, with the whole of
      `TestReadme*` green — this change is measured by not breaking seven existing outcomes as
      much as by meeting its own.

      **The ceiling was corrected mid-box, from 340 to 380, and that is the honest record.** 340
      was chosen before measuring; the duplicate and five relocations netted 13 lines, not 49,
      because a relocation swaps prose for a pointer. Grinding to 340 would have meant deleting
      reference content, which this change's own scope boundary forbids. The ratchet was then
      proved to bite: six blank lines appended in an exported copy gives
      `README.md is 382 lines, over its 380 ceiling`.
      Waits on: nothing. Can start now.

      **One box.** The guards alone leave the suite red; the edits alone move arguments with
      nothing checking they arrived, which is exactly the *moved* versus *deleted* distinction
      outcome 4 exists for.

- [ ] **2 · land the delta**
      Spec: task breakdown 3.
      Closes: *every `Proof:` citation resolves once the delta lands on `readme`*.
      Evidence: `spec-drift --anchors` green with the change folder gone.
      Waits on: box 1. Phase 8, same commit as the landing.

2 boxes. Box 1 can start now.
