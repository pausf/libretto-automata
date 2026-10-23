# spec-drift beyond the landing

Read this when running `spec-drift --trace`, or when wiring `--block` into a hook. The
landing itself needs only the three modes in `SKILL.md`.

## `--trace` — the map, not a gate

The first two are both scoped to a commit: one to what is staged, one to citations that
already exist. **Neither can see code that was never staged while a spec existed**, and
in any repository whose specs arrived after its code that is most of the repository. So
the contract can be entirely intact by both checks and still govern a third of the tree.

`--trace` reads the whole tracked tree and answers four questions:

| | |
|---|---|
| **orphan code** | a tracked file no `Governs:` line governs — no plain entry claims it, or a `!` entry gives it back |
| **dead ownership** | a `Governs:` glob matching no tracked file — a package that moved, and a spec still pointing where it was |
| **unproven criterion** | a bullet under *Verification criteria* with no `Proof:` beneath it |
| **joined criterion** | two or more `shall` resting on exactly one `Proof:`, with no ceiling named — in capability specs and change documents. One test rarely holds two promises; split it, cite a second proof, or write the ceiling into it. `--anchors` prints the count and never fails on it |

**It always exits 0, and that is deliberate.** A first run on a real project reports
dozens of orphans, because the honest answer is that most code was never specified. A
report that fails the build on the day it is introduced is a report somebody deletes,
and a deleted report finds nothing. Read it, pick the ones that matter, leave the rest.

**What it will not tell you is whether a `Proof:` names a *good* test.** That is a
reading of the test, and a script that guessed would be wrong exactly where the answer
matters. `skills/review-spec/` asks it.

**It warns; it does not block.** Always exits 0. A check that stops a commit in
someone else's project gets removed the same day, and then there is no check at all.
A deliberate no is a valid answer — say it out loud in the report so the next reader
knows the question was asked and not skipped.

Anyone who wants it to be a gate can wire it into their own `pre-commit` hook or CI.
That is their decision to make, not this flow's to make for them — and once made, it
is one paste. `--block` runs the same checks and turns the warning into exit 1:

```bash
printf '#!/bin/sh\nexec <skill-base>/spec-drift --block\n' > .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

With `<skill-base>` spelled as the real installed path. **This is opt-in and stays
opt-in**: nothing in the flow installs the hook, the default mode still always exits
0, and removing the gate is deleting the hook file.
