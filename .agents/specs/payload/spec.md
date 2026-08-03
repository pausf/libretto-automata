# Payload

Governs: skills/** commands/** agents/** scripts/**

The reason the project exists. Everything else is delivery.

The flow itself — its eight phases and its reasoning — is [docs/FLOW.md](../../../docs/FLOW.md).
This spec is the contract the payload's *artifacts* have to satisfy.

## Outcomes

Installing this repository gives a working flow on a machine that has nothing else.

- one command, `libretto-flow`, that routes and never implements
- `libretto-status`, read-only, reporting what is in flight
- one skill per phase, each of which stops where its phase stops
- three standing rules that hold at every phase: **ask**, **commit**, **evidence**
- **ceremony proportional to the change** — four stops for a change with a contract, one
  for a change too small to have one
- the vendored delegates the thin skills depend on, so a thin skill is never a broken
  one
- drift detection that ships with the skill that uses it

**Every skill is self-sufficient once installed.** A skill that only works inside this
repository is a skill that works for nobody.

## Scope boundaries

**In:** skills, agents, commands, the drift-detection tooling, the vendored
third-party items and their attribution.

**Out:**

- **requiring ponytail or caveman.** Called when present, never required. `doctor`
  reports them; nothing fails without them.
- creating a `changes/` directory in a project that has none. A staging area nobody
  empties is a second source of truth.
- accepting a secret. No skill asks for a token, and any skill that is handed one says
  it is now exposed and continues without it.
- pushing, or committing from a sub-agent.
- **a flag, profile or setting for how big a change is.** It is knowable from the change.
  Asking the user to declare it invents a second source of truth about a fact already on
  disk, and the two will disagree exactly when it matters.
- **a second flow for small work.** One flow with a proportionate gear, not two to keep
  in sync.

### Ceremony is proportional to the change

**Four waits are the price of a contract, not the price of the flow.** Phase 1 reports and
waits, phase 2 waits for the go-ahead, phase 5 waits, phase 7 waits — and every one of
those exists so the user can say "no, not that" about something they might disagree with.

When phase 2 answers **no spec needed**, there is nothing to disagree about, so the waits
go with the spec: phases 6, 7 and 8 run in one turn and exactly one question is asked at
the end. **That question is phase 8's, and it survives every collapse** — pushing is the
user's decision, not ceremony.

**What collapses is the wait, never the saying.** Phase 7 still reports what was done, its
evidence, and what was left out with the condition that brings it back. A phase that skips
the report because the change looked small is how the one omission that mattered goes
unmentioned.

The cost of getting this wrong is measured, not theoretical: a session spent four round
trips updating two documentation files, and every one of the four was mandated by this
payload. A flow that charges a typo the price of a feature gets routed around — for typos
first, then for small features, until what is left is a ritual reserved for work important
enough to deserve it.

### A phase is invoked even when it has nothing to do

**A decision belongs to the phase that owns it.** `write-spec` decides whether a spec is
needed; `build-and-check` decides where the work lands and how much to check. An
orchestrator that makes those calls itself and proceeds has not saved a step — it has
moved a decision somewhere nobody can audit.

So every phase the work reaches is invoked, **including the ones whose answer is "nothing
here"**, and the declining is reported in one line. Invoking is not gating: a phase that
declines adds no wait.

The cost is one line. What it buys is the difference between a skip and an omission, which
from outside are identical — a session made exactly the right call about a small change,
never invoked phases 2 or 6 to make it, and the user asked why the flow had not run. It
had, in substance. Nothing said so.

### Work is found, not fetched

**The flow does not begin at a tracker.** Phase 1 asks three sources, in this order:

| Source | When |
|---|---|
| a change already in flight | unchecked boxes in `.agents/changes/*/plan.md` |
| a tracker key or URL | one was given |
| what the user said | anything else — the request *is* the input |

Home first, and the order carries the reason: starting something new while a change sits
half-finished is how the half-finished thing gets abandoned, and the cost is not the lost
work but a `.agents/changes/` directory nobody trusts.

**A request in conversation is a legitimate input.** Every change in this repository so
far arrived that way, and the phase used to treat it as a fallback in a table. With no key
the change is named from the request, verb-led, and `proposal.md` records `Tracker: none`
plus what was asked in the words it was asked in — a fake key would imply a tracker
somebody could consult.

Two other implementations agree that the tracker is not the door. gentle-ai mentions none
at all in its propose or explore phases; CodelyTV/agent-harness takes the task as an
argument and offers GitHub Issues as an optional variant.

**Reporting and choosing are one skill.** `/libretto-status` invokes phase 1 in reporting
mode rather than scanning the same directories its own way. Two answers to "what work
exists" is one too many, and the one that disagrees is always the one nobody is reading.

**Landed changes are deleted, not archived.** gentle-ai moves them to `changes/archive/`;
git history is the archive here, and a directory nobody reopens is growth. A decision,
not an oversight.

## Constraints

**A skill may only invoke what gets installed.** `libretto install` links `skills/`,
`agents/` and `commands/` — nothing else. A skill telling the user to run `scripts/foo`
is broken for everyone who installed it, because `scripts/` never reaches their
machine. **A tool a skill needs ships inside the skill's own directory.**

For the same reason no skill names a path from this repository as though it existed in
the user's project. `docs/FLOW.md` is not somewhere they have.

**Frontmatter `name:` must equal the directory or filename.** A mismatch means the
caller reads one name off disk and invokes another, and the failure is a sub-agent that
silently never runs.

**Every file has exactly one author.** The brief is written once and read many times;
each spec delta has one writer; the plan has one writer, the orchestrator. This is what
makes concurrency safe without locks.

**Sub-agents start with no rules loaded.** Every launch carries them explicitly, or the
default takes over — and the default is to guess plausibly.

**Vendored items are copied unmodified**, with their licence and version recorded. A
change needed for this flow goes in the skill that calls them, never into the copy, so
the copy stays comparable with upstream.

## Prior decisions

- The tracker is read through its CLI, never MCP, never the REST API.
- The API token lives in the OS keyring, put there by `jira init`, run by the user in
  their own shell. It never enters a conversation.
- Specs are per **capability**, never per ticket: a capability spec accumulates and
  stays true, a ticket spec is dead the day the ticket closes.
- Deltas live in `.agents/changes/<change>/` and are applied onto the capability spec in
  the commit that lands the change, which then deletes the change folder.
- Drift detection **warns and never blocks**. A check that stops a commit in someone
  else's project is a check that gets deleted.
- Phase 2 may decide no spec is needed. Skipping the phase is a legitimate outcome of
  it, and the "no" collapses phase 7's gate with it.
- Phase 6 does not fan out. Parallel implementation needs isolation, a serial merge
  queue and a conflict protocol; without those three, concurrency manufactures races.
- **Phase 6 owns the branch, at step 0, before the first file is written** — not before
  the first commit. `git checkout -b` carries uncommitted work, so editing on the base
  branch and branching at commit time succeeds until the base has moved or touches one
  of your files. Phase 8 keeps the same check as a backstop, names phase 6 as its owner,
  and reports rather than silently fixes: a backstop that covers for the rule it backs up
  is how the rule stops being followed.
- **Push and the pull request are one question.** Asked separately they bought a second
  round trip and no safety — a pushed branch with no request opened is a state almost
  nobody wants, and whoever wants it says so in the same breath.
- **The forge is derived, never assumed:** `git remote get-url origin`, `github.com` →
  `gh`, `gitlab` → `glab`. No remote means no question. **Ceiling named:** a substring
  test on one URL, which does not survive a self-hosted forge on a neutral domain, or
  Gitea and Forgejo. The replacement that day is an explicit setting read from the
  repository, not a longer list of guesses in a prompt.
- **A missing forge CLI stops**, with its install line and nothing else — the shape phase
  1 already uses for `jira`. Never a hand-built API call using a token found in the
  environment: that turns a stop into an exposure. The other forge's CLI is not a
  fallback.
- **Whether a change needs a spec stays a judgment, with no heuristic behind it** — not a
  diff-size threshold, not a file count. The real test is whether two people could
  reasonably disagree about what "done" means. A number would be wrong in both directions
  and trusted anyway.

## Task breakdown

- [x] `find-work` — phase 1, three sources
- [x] `libretto-status` — read-only, what is in flight
- [x] `write-spec` — phases 0, 2, 3, including the fan-out and its brief
- [x] `write-plan` — phase 5
- [x] `build-and-check` — phase 6
- [x] `present-work` — phase 7
- [x] `record-work` — phase 8, including consolidation
- [x] `evidence` — the standing rules
- [x] `libretto-flow` — the routing command
- [x] vendored delegates with attribution
- [x] `spec-drift` and `check-payload`
- [x] **the flow run end to end against a real task.** `right-size-the-flow` — phase 1
      found it in flight, 2 wrote the delta, 5 the plan, 6 closed six tasks with a commit
      each, 7 reported, 8 consolidated. The forge stop was exercised for real, not
      simulated: `gh` was absent, phase 8 stopped with the install line.
- [ ] **the failure paths.** Still written and never run: an unconfigured tracker, a board
      URL where a key was expected, the trivial lane actually collapsing on a one-line
      change, and a sub-agent hitting a question it must not answer itself.
- [ ] an independent verifier: check the implementation against the spec's criteria,
      never run by whoever wrote the code

## Verification criteria

- frontmatter parses, and `name:` matches the directory or filename
  Proof: scripts/check-payload
- no stray file sits where the linker would install it as an item
  Proof: scripts/check-payload
- every referenced skill exists
  Proof: scripts/check-payload
- **no skill invokes a path that does not get installed**
  Proof: scripts/check-payload
- glob matching, capability derivation and citation extraction behave
  Proof: skills/record-work/spec-drift --self-test
- **an anchor inside a fenced code block is an illustration, not a declaration.** Any
  document explaining the convention shows it in a fence; reading that literally made
  the index page claim to govern the linking package. Knowing the difference is the
  extractor's job, not the author's job to avoid illustrating.
  Proof: skills/record-work/spec-drift --self-test
- **a glob that matches real directories is still a pattern.** Unquoted, the shell
  path-expands it before it is ever compared, so every glob matching real files
  silently stopped working and the check reported drift nowhere. `set -f` off, and a
  false negative in a checker is worse than no checker.
  Proof: skills/record-work/spec-drift --self-test
- **an invented test name is rejected rather than accepted**
  Proof: skills/record-work/spec-drift --self-test
- every `Proof:` citation in every spec resolves, file and test name
  Proof: skills/record-work/spec-drift --anchors
- **the status command delegates rather than restating the scan**, and every referenced
  skill survived the rename
  Proof: scripts/check-payload

**What none of this verifies is behaviour.** A skill is a prompt, and a prompt is
checked by running it. The static checks above catch what silently degrades one — a
broken reference, an unreachable tool, frontmatter a host cannot parse — and nothing
more.

The failure paths are the ones worth exercising first, because they are the ones written
and never run: an unconfigured tracker, a board URL where a key was expected, a trivial
task that should skip phase 2, and a sub-agent that hits a question it must not answer
itself.

**What the first real run observed**, stated as observations rather than citations, because
a criterion citing a test that cannot exist is the fabrication the anchor exists to
prevent:

- phase 6 created the branch before the first edit; nothing landed on `main`
- the missing-CLI stop fired for real — `gh` absent, install line given, no workaround
  offered and no token asked for
- a phase declining was reported in one line, and cost no wait
- the four stops still happened, because this change had a spec — the gear is
  proportionate, not removed
- `prune` removed exactly the one stale link it planned to and left the other thirteen,
  the first exercise of that path outside a temporary directory

Still unobserved, and therefore still claims rather than facts: the collapsed lane on a
change that needs no spec, and every remaining failure path above.
