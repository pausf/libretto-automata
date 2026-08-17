# add-payload-index — delta

Targets: payload

`skills/`, `agents/` and `commands/` are the payload — the reason the project exists — and
there is **no one page listing what ships.** 22 skills, 7 agents, 7 commands, and mapping a
phase to the directory that implements it is left to the reader of `docs/FLOW.md`'s prose.

## The trap, which the proposal named before the work started

**A hand-written list of directories drifts.** This repository has been bitten twice: a
paragraph in `AGENTS.md` said *ten* over eleven spec directories, which is why `docs/SPEC.md`
is now the only place the capability list lives. A typed `docs/PAYLOAD.md` would reproduce
that failure exactly, one directory deeper — and the copy nobody edited is the one that reads
as authoritative.

So the index is **generated and gated**, never typed. `scripts/check-payload` already walks
all three directories and parses every frontmatter `name:` and `description:`, so it has the
data already; what it lacks is a way to emit it and a check that the committed page still
matches.

## Outcomes

1. **`scripts/check-payload --index` writes `docs/PAYLOAD.md`** — every skill, agent and
   command that ships, each with the `description:` from its own frontmatter, grouped by kind
   and sorted within each group.

2. **The default run fails when `docs/PAYLOAD.md` is stale.** It regenerates the index in
   memory and compares; a mismatch names what changed and says to re-run with `--index`. This
   is the whole reason the page can be trusted: an item added without regenerating breaks the
   payload gate, which is one of the six.

3. **No skill, agent or command references the index** — as a **constraint**, not a gated
   outcome, because the first draft of this said the existing gate covered it and that was
   wrong.

   `docs/` is not installed, so a reference would be broken for everyone who installed the
   payload rather than cloning it. But the uninstalled-path gate is **scoped to executables**,
   `scripts/` and `bin/`, and its own comment gives the reason: prose mentioning `docs/` is
   describing rather than instructing, and failing those *"turns this check into noise nobody
   reads."*

   So it stays a constraint, in the same position `payload` already puts `docs/FLOW.md` — a
   path this repository has that the user's project does not. **Not widened to `docs/`**,
   because that would fail every legitimate mention in the payload's own prose, and a gate
   nobody reads is worse than a stated rule.

4. **`docs/PAYLOAD.md` says it is generated**, in its first lines, so nobody edits it by hand
   and loses the edit on the next regeneration.

## Scope boundaries

**In:** `scripts/check-payload`, the generated `docs/PAYLOAD.md`, and this delta.

**And one thing added to scope during the work, stated rather than slipped in:** the new page
exposed a pre-existing defect in the same script's referenced-paths check. It piped `rg -oN`,
which prefixes each hit with its filename, into a second `rg` that re-extracted from that
prefix — so `docs/PAYLOAD.md:scripts/check-payload` yielded *both* strings and any file under
`docs/` or `scripts/` verified itself as a referenced path. `docs/PAYLOAD.md` and `docs/PLAN.md`
were both reported `ok` with nothing referencing them.

Fixed here, with `-I`, because it is one flag and this change is what made it visible. It is
**noise rather than a false pass** — a genuine reference is a text match either way — but a check
that reports paths nobody cites is a check whose output stops being read, which is the failure
mode this repository already records about drift warnings.

**Out**, named:

- **the description text is not written here.** It is whatever the item's own frontmatter says,
  which is what makes the page generated rather than authored. A page with hand-tuned blurbs is
  a page that drifts, and it would need a second gate to hold it
- **no phase-to-directory mapping.** The proposal's reading mentions it and it is refused: a
  phase is a concept in `docs/FLOW.md`, several items serve one phase, and a mapping asserted in
  a generated file would need a source of truth that does not exist. `FLOW.md` names its skills
  in prose already
- **`docs/FLOW.md`, `docs/SPEC.md` and `README.md` are unchanged.** A link to the index from
  `README.md` is a change to the `readme` capability and its section test — not this one
- **no new script.** `check-payload` already has the parse; a second tool reading the same
  frontmatter is a second parser to keep in step
- **nothing installed changes.** No skill, agent or command file is touched
- **no `--fix`-style flag on the default run.** A gate that repairs what it measures is a gate
  that never fails, which is the whole failure mode of a generated file with no check

## Constraints

- **`scripts/` is repo-only tooling and is never referenced from a skill.** `payload` already
  states this and gates it
- the generator reuses the frontmatter parse already in the file — the same `rg -N -m1 '^name:'`
  and `'^description:'` extraction, so a page and a gate can never disagree about what an item
  is called
- **a `description:` is a single frontmatter line here**, which is how every payload item is
  already written and what the existing gate assumes. A multi-line YAML description would come
  out truncated — named under prior decisions with what replaces it
- the comparison must not depend on trailing-newline handling differing between the writer and
  the reader, or the gate fails on a file it just wrote
- `set -uo pipefail` is in force; the script does not use `-e`, so a failure is reported through
  `fail` rather than by aborting
- all six gates green

## Prior decisions

- **Generated and gated, not typed** — the proposal decided this before the work began and the
  repository has the scar tissue to justify it. Not revisited.
- **`--index` writes the file; the default run only compares** — *assumed, nobody was asked.*
  `/libretto-attacca` answered the contract stop. A gate that rewrites the tree it is checking
  makes `git status` part of its output and can never fail. What changes if this is wrong:
  nothing about the page, only who types one command.
- **The page lives at `docs/PAYLOAD.md`** — *assumed.* `docs/` is where uninstalled,
  human-facing documents live, and the index is exactly that.
- **Grouped by kind, sorted within the group** — *assumed.* Sorted rather than in flow order,
  because flow order is a judgment that would have to be maintained by hand — which is the
  drift this whole change exists to avoid. `docs/FLOW.md` is where order is explained.
- **A single-line `description:` is assumed** — *assumed, and this is the sharpest ceiling.*
  Every payload item today writes one line, and the existing gate already depends on that. A
  multi-line YAML description would be truncated to its first line in the page **and the gate
  would still pass**, because generator and comparison share the same parse. The replacement,
  the day an item needs a folded description, is a `fail` when a `description:` value ends in
  `>` or `|` — never a YAML parser, because `payload` already records that a dependency added
  to check a document is a dependency added to check a checker.
- **No entry in `docs/SPEC.md`** — *assumed.* That file indexes *capabilities*, and the payload
  index is not one. Adding it there would blur the one thing that file is for.

## Task breakdown

1. **The staleness gate first, red.** Add the comparison to the default run and prove it fails
   with no `docs/PAYLOAD.md` present at all — the honest first red.
2. **Then `--index` and the generated page**, and watch the gate go green. Then prove the gate
   bites: hand-edit a description in the page and see it fail; add a probe skill and see it
   fail.
3. Land this delta onto `payload` and delete the change folder.

**Steps 1 and 2 are one box.** The gate alone fails the whole suite; the generator alone
produces a page nothing keeps current, which is the typed list this change exists to avoid
wearing a script's clothes.

## Verification criteria

- **`scripts/check-payload` fails when `docs/PAYLOAD.md` does not match a fresh generation**,
  and the failure says to re-run with `--index`.
  Proof: scripts/check-payload

  **Watched red before green, and then watched biting**: absent file first, then a hand-edited
  description, then an added item. A staleness gate that has only ever run against a fresh page
  has proved nothing.

- **`scripts/check-payload --index` regenerates `docs/PAYLOAD.md` such that the default run then
  passes** — the two halves agree by construction, because they share one parse.
  Proof: scripts/check-payload

- **every item under `skills/`, `agents/` and `commands/` appears in `docs/PAYLOAD.md`.** The
  directories are read, never a list in the script — a list is the same drift one level down.

  **This is the same check as the first criterion, not a second one**, and it is written
  separately only because it is the promise a reader of the page cares about while staleness is
  the mechanism that keeps it. Adding an item without regenerating *is* a mismatch. Said out
  loud because two criteria citing one gate reads as two checks, and then nobody notices when
  one of them was never implemented.
  Proof: scripts/check-payload

**No criterion for outcome 3.** The uninstalled-path gate is scoped to executables by design, so
nothing mechanically stops a skill from citing `docs/PAYLOAD.md` — exactly as nothing stops one
citing `docs/FLOW.md` today. It is a constraint `payload` already carries in prose, and widening
the gate to `docs/` would fail every legitimate mention in the payload's own writing.

**What nothing here tests:** whether a `description:` is *accurate*, or whether the grouping
reads well to somebody navigating the repository for the first time. Both are readings.
