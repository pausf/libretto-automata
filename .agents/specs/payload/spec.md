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
  it.
- Phase 6 does not fan out. Parallel implementation needs isolation, a serial merge
  queue and a conflict protocol; without those three, concurrency manufactures races.

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
- [ ] **run the flow end to end against a real task.** Nothing here has been executed.
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
