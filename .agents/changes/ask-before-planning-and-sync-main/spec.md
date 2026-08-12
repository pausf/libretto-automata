# ask-before-planning-and-sync-main — delta

Targets: payload

Two amendments to the flow, unrelated except that both arrived in one sentence.

## Outcomes

### 1 · Phase 8 leaves the base branch current

When the push-and-open-the-request question is answered yes **and the push and the
request are both confirmed landed**, phase 8 does one more thing before it reports:

```
git checkout <base> && git pull --ff-only
```

Afterwards the working tree is on the base branch, up to date with the remote. The
feature branch is left alone — not deleted, not merged, not rebased.

The base is whatever the repository says it is (`git symbolic-ref refs/remotes/origin/HEAD`,
falling back to the branch the request was opened against). Never hardcoded to `main`.

**It runs only on the yes path.** No push means the work exists nowhere else, and moving
off the branch that holds it buys nothing.

### 2 · Phase 2 asks up to three questions before it writes

`write-spec` asks the points it cannot settle from the code — **at most three, and none
when there are none.** The questions go out in one `AskUserQuestion` call, before the
spec file is written, so the answers are in the contract rather than on top of it.

What qualifies: a product tradeoff, two live precedents with nothing choosing between
them, a boundary the request left open, a scope decision that changes what gets built.

What does not: anything the code answers, anything the proposal already states, and
anything whose only available answer is "yes, carry on".

Every answer lands in the spec under **prior decisions**, next to what it settled.

## Scope boundaries

**In:** `skills/record-work/SKILL.md`, `skills/write-spec/SKILL.md`,
`commands/libretto-flow.md`, `docs/FLOW.md`, and this delta landing on
`.agents/specs/payload/spec.md`.

**Out:**

- **questions in phase 5.** The user chose phase 2 alone. The plan already stops for the
  order and what waits on what; a second tranche of questions before the first line of
  code is two round trips where the contract only needed one.
- **a fourth stop.** The three questions ride the stop phase 2 already has. The flow's
  stop count does not move.
- **a minimum question count.** Zero is a legitimate answer, reported in one line. A
  forced quota manufactures questions the code already answers, which `AGENTS.md`
  forbids in as many words.
- **asking in the trivial lane.** No spec means no contract to disagree about; a change
  that skipped phase 2's writing skips its questions with it.
- **deleting the feature branch**, merging it, or waiting on the request. The request is
  open, not merged, and a branch deleted at that moment takes the only local copy of
  unmerged work with it.
- **touching the base branch when the answer was no**, or when the push or the request
  failed to confirm.
- **`git pull` without `--ff-only`.** A merge commit manufactured on the user's base
  branch by a phase whose job is bookkeeping is the surprise nobody wants.

## Constraints

- The tree is clean by the time this runs — phase 8 has just committed everything. If it
  is not, **report and stay put.** `git checkout` carries uncommitted work onto the base
  branch silently, which is the same failure the "branch before the first write" rule
  exists to prevent, pointed the other way.
- A `--ff-only` pull that refuses is reported, never forced. Diverged base is a fact
  worth seeing, not a thing to resolve unasked.
- The three questions are one `AskUserQuestion` call, not three turns. Three round trips
  to build one contract is the ceremony the payload spec spends its length arguing
  against.

## Prior decisions

- **Phase 8 checks out the base rather than updating its ref from the branch.** The
  user's call, this run. `git fetch origin main:main` would leave you on the feature
  branch with a current base, which is cheaper — but the observed failure was a *next
  session* starting from a stale tree, and a session starts wherever the last one left
  the working directory. Landing on the base is what fixes that.
- **Questions in phase 2 only.** The user's call, this run, from three offered options.
- **Up to three, zero allowed.** The user's call, this run. The alternative — always
  three — was offered and declined.
- **This is not hypothetical.** Phase 1 of the very run that produced this change
  reported a merged, tagged branch as work in flight, because local `main` was seven
  commits behind `origin/main`. The reading was wrong, it was reported to the user as
  fact, and the correction cost a round trip.

## Task breakdown

- [ ] `skills/record-work/SKILL.md` — the return to the base branch, on the yes path only
- [ ] `skills/write-spec/SKILL.md` — step 4 becomes up-to-three, one call, zero allowed
- [ ] `commands/libretto-flow.md` — phase 8's line and phase 4's line
- [ ] `docs/FLOW.md` — the reasoning behind both
- [ ] `.agents/specs/payload/spec.md` — this delta applied, change folder deleted

## Verification criteria

- frontmatter still parses and every referenced skill still exists
  Proof: scripts/check-payload
- no skill invokes a path that does not get installed
  Proof: scripts/check-payload
- every `Proof:` citation in every spec resolves, file and test name
  Proof: skills/record-work/spec-drift --anchors

**Neither outcome is verifiable by a test, and no criterion will claim otherwise.** Both
are prose in a skill, and a skill is a prompt — checked by running it. Inventing a Go test
for "phase 8 checks out the base branch" would be a citation to something that cannot
exist, which is the exact fabrication `--anchors` was built to catch.

So both land as **claims until observed**, in the payload spec's own idiom:

- a run answering yes to push ends with the working tree on the base branch, current with
  the remote, and the feature branch intact
- a run answering no ends on the feature branch, unchanged from today
- a phase 2 with real open points asks at most three in one call, and the answers appear
  in the spec under prior decisions
- a phase 2 with nothing to ask says so in one line and writes the spec
