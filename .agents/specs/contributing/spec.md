# contributing

Governs: CONTRIBUTING.md

The contributor's door. `README.md` is written for somebody deciding whether to install this;
this file is for somebody about to change it, and they are a different reader with a different
question.

**This capability exists because `readme` says so.** That capability governs the front door and
names `AGENTS.md` out of scope in as many words — *the contributor's door, and a different
reader.* `CONTRIBUTING.md` is that reader, so folding it into `readme` would contradict
`readme`'s own boundary.

## Outcomes

1. **`CONTRIBUTING.md` is at the root**, where GitHub links it from the issue and pull-request
   forms. That link is the whole reason a root file beats a section in `AGENTS.md`: it arrives
   at the moment the contributor needs it, rather than waiting to be found.

2. **It states only what is contributor-specific, and links for everything else.** The five
   things that live nowhere else, or that an outside contributor has no way to guess:

   | | Why it is here and not only in `AGENTS.md` |
   |---|---|
   | the six gates as a runnable block | a contributor needs to paste and run them, not read about them |
   | a merge needs a `release:` label, and the run **refuses** without one | the most surprising thing about this repository from outside, and a designed refusal that has already been read once as a broken pipeline |
   | the bump is a reading of `.agents/specs/`, so a contributor **proposes** it | otherwise a contributor picks the smaller one to avoid the question, which is exactly what `AGENTS.md` forbids |
   | work does not come from a tracker | there is no board; an issue describing the change *is* the contribution, and a ticket id should not be invented |
   | `1.0.0` is undecided, so nothing is stable | a version number cannot be recalled once the proxy has cached it |

   And one about the process rather than the rules: **what review looks like here** — findings
   arrive attributed and unedited, including the ones already fixed, and a reviewer is disagreed
   with by evidence rather than by agreement. Both are unusual enough that meeting them
   unwarned reads as hostility rather than as the house style.

3. **No rule `AGENTS.md` owns is restated.** Four phrases it owns are **absent**:
   `Co-Authored-By`, `type(scope): subject`, `72 chars`, and `stdlib, then a native`.

   **Absent, with no "unless the line links to it" escape clause.** The first draft had one, and
   it could not be implemented: the scan runs over the flattened document, where there are no
   lines to reason about. If the guide needs to raise commit conventions it links without
   quoting the rule, which is the behaviour wanted anyway.

4. **Every relative link resolves**, and the link pattern must match something — a door made of
   no links is not a door, and a scan finding none would pass vacuously.

   **No link points at a path whose existence depends on work being unfinished.**
   `.agents/changes/` is the case: git does not track empty directories, so once every change has
   landed the directory is not in a fresh checkout at all. This shipped as a link and broke CI
   on the one branch where the queue had finally emptied — green on the author's machine, where
   six queued proposals were sitting on `main`, and green on the two branches before it, which
   still had folders left.

   **The guard could not catch it and still cannot.** Link resolution reads the working tree, so
   a link to a directory that exists only while work is in flight passes while work is in flight.
   The directory is named in prose now, never linked. **Ceiling:** nothing gates this — the
   replacement, the day it recurs, is a link check run against `git ls-files` rather than the
   filesystem, which is a real change to how that test works and not worth it for one path.

## Scope boundaries

**In:** `CONTRIBUTING.md`, and the test that holds it.

**Out:**

- **`AGENTS.md`.** Not edited by this capability and **not governed** by it — see prior
  decisions, where the cost is named.
- **Wording.** No criterion pins a sentence. The no-tracker clause is anchored on its
  **heading**, so improving the prose does not break the guard.
- **Issue and pull-request templates.** A separate promise and a separate maintenance burden.
- **A code of conduct.** A real one is a commitment to enforcement by a person.
- **`README.md`.** It already links out to `AGENTS.md` from *Learn more*; a new front-door link
  is `readme`'s change, not this one.
- **Whether the guide is *useful* to a stranger.** Read, never checked — `readme` records why
  pretending otherwise is expensive.

## Constraints

- **A rule stated twice is what this capability exists to prevent**, so outcome 3 is a check and
  not an intention.
- The proof lives at `cmd/libretto/readme_test.go`, which already reads four documents outside
  its own package and owns `repoFile` and `flat`. A new package for one more document is not
  worth its own directory — the decision `readme` already recorded.
- **The scan runs over `flat()`.** That file has shipped two guards unable to fire because
  Markdown wrapped between the words they searched for; the third caught a planted
  `Co-Authored-By` **across a line break**, which is `flat()` earning its place rather than
  merely being present.
- **No version number in prose.** `1.0.0` is discussed as a decision nobody has made, never as
  a version anybody is on.

## Prior decisions

- **A capability of its own rather than a widening of `readme`.** `readme` names `AGENTS.md` out
  of scope as a different reader, and this file is that reader.

- **`AGENTS.md` stays ungoverned, and it is a decision rather than an oversight.** The same batch
  that created this capability made the *opposite* call about `THIRD-PARTY.md`, widening
  `payload` to claim it — so the difference is stated: `THIRD-PARTY.md` is **parsed by a gate**
  and changes rarely, while `AGENTS.md` is the working rulebook and changes constantly.
  Governing it would make every rule edit a spec edit, and a drift check that fires on every
  commit is a check that gets muted — which `payload` already records as the reason drift
  detection warns and never blocks.

  **Ceiling, named because it is the same criticism this repository levelled at its own licence
  files:** `AGENTS.md` is therefore a path where drift is nobody's finding. Accepted because it
  is prose about process rather than a contract any code satisfies, and because six specs
  already reference it, so a contradiction surfaces when one of them is read. The replacement,
  the day `AGENTS.md` contradicts a capability spec unnoticed, is a criterion asserting the
  specific overlapping claims — never a `Governs:` line over the whole file.

- **The guide is short and mostly links.** The alternative, a self-contained guide, is the
  two-sources-of-truth failure with better first-time reading. Outcome 3 keeps the decision
  honest after the fact.

- **The `release:` label rule is stated most prominently**, because it is the designed refusal an
  outside contributor cannot predict, and this repository has already had it read as breakage.

- **No code of conduct.** A commitment to enforcement is a person's to make.

- **Every clause of outcome 2 has its own assertion**, after the 6→7 reviewer measured that two
  of them had none: deleting the whole no-tracker section left the guard green, and so did
  removing five of the six gate lines. A criterion citing one test for four clauses needs four
  assertions or a stated ceiling — never the citation alone, which is the recorded
  *a criterion can cite a gate that tests half of it* failure.

## Task breakdown

Held going forward, not open work:

- when a rule moves into `AGENTS.md` that the guide paraphrases, the guide loses the paraphrase
- when a gate is added or removed, the runnable block and its assertion move together
- when a phrase is added to the forbidden set, it is verified present in `AGENTS.md` first, so
  the set names real rules rather than plausible ones

## Verification criteria

- **`CONTRIBUTING.md` exists and names each contributor-specific thing, one assertion per
  clause** — all six gate commands individually, `release:patch`, `.agents/specs/`, the
  no-tracker section by its heading, and `1.0.0`.
  Proof: cmd/libretto/readme_test.go TestContributingIsADoorNotACopy

- **none of the four phrases `AGENTS.md` owns appears**, matched over the flattened document.
  Proof: cmd/libretto/readme_test.go TestContributingIsADoorNotACopy

  **Ceiling:** it is a named set, so a rule restated in words nobody listed passes — the same
  class as `readme`'s badge word list and its command-name guard. The replacement, the day it
  matters, is a criterion about how long this file may get, not a longer phrase list.

- **every relative link resolves, and the pattern matched something.**
  Proof: cmd/libretto/readme_test.go TestContributingIsADoorNotACopy

  **Every arm watched red, and independently.** The first red was the honest one — the file did
  not exist. Then, in an exported copy so the repository was never touched: a planted
  `Co-Authored-By`, a `type(scope): subject` **wrapped across a newline**, a removed `1.0.0`, a
  removed `release:patch`, a deleted no-tracker heading, five deleted gate lines, a dead
  relative link, and a link-free document. Each fires.
