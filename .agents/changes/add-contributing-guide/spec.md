# add-contributing-guide — delta

Targets: contributing — **a new capability.** No existing one fits, and that is the discovery
rather than a problem to route around.

`readme` governs `README.md` and names `AGENTS.md` out of scope in as many words: *the
contributor's door, and a different reader.* `CONTRIBUTING.md` is that same reader. Folding it
into `readme` would contradict that capability's own boundary; folding it into `payload` would
put a document about how to contribute inside the spec for what gets installed.

`docs/SPEC.md` gains its row, because that file is the index **and the only place the list
lives** — the lesson `AGENTS.md` records about a count kept in two places.

## The risk this change is mostly about

**Not the missing file. The duplicate.** `AGENTS.md` already carries the gates, the commit
convention, the `release:` label rule and the boundaries, and it carries them well. A
`CONTRIBUTING.md` restating any of that is the failure `CLAUDE.md` opens by naming: two files
kept in sync is two sources of truth, and the one that wins is the one nobody edited.

So the guide's value is **routing and the contributor-specific**, never restatement.

## Outcomes

1. **`CONTRIBUTING.md` exists at the root**, so GitHub links it from the issue and PR forms.
   That link is the whole reason a root file beats a section in `AGENTS.md` — it arrives at the
   moment the contributor needs it.

2. **It states only what is contributor-specific, and points at `AGENTS.md` for the rest.**
   What is contributor-specific and lives nowhere else today:

   - **the six gates as a runnable block**, so a contributor can paste and run them before
     opening anything
   - **that a PR merges with a `release:patch`, `release:minor` or `release:major` label and
     the run refuses without one** — the single most surprising thing about this repository
     from outside, and a designed refusal that reads as a broken pipeline when unpredicted
   - **that the bump is a reading of `.agents/specs/`, so a contributor proposes it and does
     not decide it**
   - **where work comes from**, because it is not a tracker and an outside contributor has no
     way to guess that
   - **that `1.0.0` has not been decided**, so nothing here is stable

3. **No rule is stated twice.** Four phrases that `AGENTS.md` owns are **absent** from
   `CONTRIBUTING.md` — the exact four are named once, in the criterion below, and nowhere else
   in this delta. This is what makes the guide a door rather than a copy.

   **Absent, with no "or links to it" escape clause.** The first draft allowed the phrase if the
   line carrying it linked to `AGENTS.md`, and the scan is over the flattened document — where
   there are no lines, so the clause was unimplementable. Rather than un-flattening for it: if
   the guide needs to raise commit conventions it links to `AGENTS.md` without quoting the rule,
   which is the behaviour wanted anyway.

4. **Every relative link in `CONTRIBUTING.md` resolves.** A door made of links whose links are
   dead is worse than no door.

## Scope boundaries

**In:** `CONTRIBUTING.md`, the new capability spec, its row in `docs/SPEC.md`, and one test
case.

**Out**, named:

- **`AGENTS.md` is not edited, and not governed.** See prior decisions — this is the decision
  most likely to be argued with, so it is recorded rather than assumed
- **no issue or PR templates**, no `.github/ISSUE_TEMPLATE/`. The ask was a contributing guide;
  templates are a separate promise and a separate maintenance burden
- **no code of conduct.** A real one is a commitment to enforcement, and inventing that on the
  user's behalf is not this run's call
- **no `README.md` change.** `readme` outcome 7 already links out to `AGENTS.md` from *Learn
  more*, and the `readme` capability's own test holds the section order — a new front-door link
  is a change to that capability, not this one
- **no test that the guide reads well.** `readme` already records that this is the one thing no
  test holds, and pretending otherwise is what produced the README it replaced
- **no duplication of the six gates' *reasoning*.** The commands, yes — a contributor needs to
  run them. Why each exists is `AGENTS.md`'s and `docs/`'s

## Constraints

- **A rule stated twice is the thing this change must not produce**, so outcome 3 is a check
  and not an intention
- the new case joins `cmd/libretto/readme_test.go`, which already reads three documents outside
  its own package and owns `repoFile` and `flat`. A new package for one more document is not
  worth its own directory
- **the scan runs over `flat()`.** That file has now shipped two guards unable to fire because
  Markdown wrapped between the words they searched for
- `docs/SPEC.md` carries no requirements and no `Governs:` line — it is a signpost, and adding
  a row must not turn it into a second spec
- **no version number in prose.** The same rule `readme` records; `1.0.0` is discussed as a
  decision nobody has made, never as a version anybody is on
- all six gates green

## Prior decisions

- **A new `contributing` capability rather than widening `readme`** — *assumed, nobody was
  asked.* `/libretto-attacca` answered the contract stop. `readme` names `AGENTS.md` out of
  scope as a different reader, and `CONTRIBUTING.md` is that reader. What changes if this is
  wrong: the spec moves into `readme` and its `Governs:` grows one path.

- **`AGENTS.md` stays ungoverned, and this is a decision rather than an oversight** —
  *assumed.* Change 3 of this batch made exactly the opposite call about `THIRD-PARTY.md`, so
  the difference has to be stated: `THIRD-PARTY.md` is *parsed by a gate* and changes rarely,
  while `AGENTS.md` is the working rulebook and changes constantly. Governing it would make
  every rule edit a spec edit, and a drift check that fires on every commit is a check that
  gets muted — which `payload` already records as the reason drift detection warns and never
  blocks.

  **Ceiling, named because it is the same criticism this batch levelled at licence files:**
  `AGENTS.md` is therefore a path where drift is nobody's finding. Accepted because it is
  *prose about process* rather than a contract any code satisfies, and because six specs
  already reference it, so a contradiction surfaces when one of them is read. The replacement,
  the day `AGENTS.md` contradicts a capability spec unnoticed, is a criterion asserting the
  specific overlapping claims — never a `Governs:` line over the whole file.

- **The guide is short and mostly links** — *assumed.* The alternative, a self-contained guide,
  is the two-sources-of-truth failure with better first-time reading. Outcome 3 exists to keep
  the decision honest after the fact.

- **The `release:` label rule is the one thing stated most prominently**, because it is the
  designed refusal an outside contributor cannot predict. This repository has already had that
  refusal read as a broken pipeline once.

- **No code of conduct** — *assumed, and it is the assumption most worth overriding.* A code of
  conduct is a commitment to enforcement by a person, and a run cannot make that commitment on
  the user's behalf.

## Task breakdown

1. **The no-duplication guard first, red.** Add `TestContributingIsADoorNotACopy` and watch it
   fail — before `CONTRIBUTING.md` exists it fails on the missing file, which is the honest
   first red rather than a manufactured one.
2. **Then `CONTRIBUTING.md`**, and watch it go green.
3. Write `.agents/specs/contributing/spec.md`, add the `docs/SPEC.md` row, delete the change
   folder.

**Steps 1 and 2 are one box.** The guard alone fails; the file alone can restate `AGENTS.md`
freely with nothing watching. Neither half merges.

## Verification criteria

- **`CONTRIBUTING.md` exists and names the four contributor-specific things** — the gate
  commands, the `release:` label requirement, that work does not come from a tracker, and that
  `1.0.0` is undecided.
  Proof: cmd/libretto/readme_test.go TestContributingIsADoorNotACopy

- **no rule `AGENTS.md` owns is restated in `CONTRIBUTING.md`.** These four phrases are absent
  — this is the only place the set is written, so it cannot drift against a second copy:

  | Phrase | Owned by `AGENTS.md` under |
  |---|---|
  | `Co-Authored-By` | Commits — no AI attribution |
  | `type(scope): subject` | Commits — the format |
  | `72 chars` | Commits — the subject length |
  | `stdlib, then a native` | Ask first — the dependency ladder |

  All four verified present in `AGENTS.md` before this criterion was written, so the set names
  real rules rather than plausible ones. Matched over the flattened document, and **absent means
  absent** — no escape clause, because a clause about "the line containing it" cannot be
  implemented against a document with its lines flattened away.
  Proof: cmd/libretto/readme_test.go TestContributingIsADoorNotACopy

  **Ceiling:** it is a named set, so a rule restated in words nobody listed passes. Same class
  as the badge word list and the command-name guard — the replacement, the day it matters, is a
  criterion about how long the file may get, not a longer phrase list.

- **every relative link in `CONTRIBUTING.md` resolves**, and the pattern must have matched
  something, so a file with no links cannot pass vacuously.
  Proof: cmd/libretto/readme_test.go TestContributingIsADoorNotACopy

**What nothing here tests:** whether the guide is *useful* to somebody who has never seen the
project. That is read, not checked, and `readme` already records why pretending otherwise is
expensive.
