# add-unattended-flow-mode

Targets: payload

An unattended mode for the flow, entered through its own command: run the phases end to
end without stopping at a question, and take the branch to a pushed request.

**The change is a classification, not a feature.** The flow's stops are already
enumerated — three in `payload/spec.md`'s table, two exceptions in phase 1, one ask rule
bounded to phases 1, 2 and 5. What this adds is a second reading of that same list: which
of them are *questions the invocation can answer in advance*, and which are not questions
at all.

Getting that line wrong is the whole risk. A mode that answers a gate is not unattended,
it is unverified.

## Outcomes

`/libretto-attacca <task>` runs phases 1 through 8 in sequence and ends at a pushed branch
with a request open on it, having asked nothing.

**It routes; it does not restate.** Like `/libretto-next`, it knows one thing the flow does
not — here, what the invocation already answered — and delegates every phase to the same
skill `/libretto-flow` delegates to. Nothing about the eight phases is written twice.

- **the two in-work stops are answered by the invocation.** The spec stop and the plan
  stop carry on; the artifacts are still written, still committed, still reported.
- **phase 8's question is answered yes**, both halves: push, and open the request. This
  is the only place the mode overrides an "ask first" rule in `AGENTS.md`, and it does it
  because the invocation *is* the answer — the user typed this command and not the other
  one, knowing what it does.
- **a question the flow cannot derive is assumed, recorded and passed on.** It goes into
  the spec under prior decisions marked as assumed, into the phase 7 report, and into the
  request's description. It never becomes a prompt and it is never left invisible.
- **every gate still runs, and a failing one still stops.** `build-and-check`'s two-failed-
  gates rule and `evidence`'s never-weaken-a-test rule are untouched. A stopped run ends
  where it stopped, on its branch, with what it observed — and with no request opened.
- **a stop that is the input failing to arrive still stops**: a missing or unauthenticated
  `jira`, `gh` or `glab`. There is nothing downstream to run and no answer to assume.
- **the run says, once, what the mode answered on the user's behalf.** The report and the
  request description both carry it, so the command's cost is legible after the fact.

`/libretto-flow` is unchanged — the same three stops, in the same places, and one line
naming the other door.

## Scope boundaries

**In:** `commands/libretto-attacca.md`, and the one line each stop-owning skill needs so its
own "then wait" does not contradict the mode that is running it.

**Out:**

- **merging, tagging, releasing, or putting a `release:` label on the request.** The mode
  ends at a request open for review. The bump is a reading of the specs rather than of the
  commits, `release:major` is asked-and-waited-for by standing rule, and a version number
  cannot be recalled once the proxy has it. Nothing here goes near that.
- **skipping, reordering or softening a gate**, and `--force`, `--no-verify` or any other
  flag that buys a green result. Unattended removes waits, never checks.
- **draining the queue.** `payload/spec.md` already scopes that out and this does not
  reopen it: one invocation is one piece of work, exactly as `/libretto-flow <task>` means
  that task.
- **restating the routing.** `/libretto-attacca` names the phases it delegates to and nothing
  about how they work. The moment it explains a phase, there are two descriptions of that
  phase and only one of them gets maintained.
- **an unattended `/libretto-next` or `/libretto-review`.** Both are their own commands
  with their own first question, and reaching them is a second decision. It returns the
  day somebody wants it.
- **a setting, profile, env var or config file for the mode.** It is consent for one run,
  and consent that persists is consent nobody remembers giving. It lives in which command
  was typed and nowhere else.
- **a flag on `/libretto-flow` that does the same thing.** Two doors into one mode is the
  drift the single command exists to avoid, and the argument-shaped one is the one that
  gets typed by accident.
- **assuming past a missing credential.** A stop for a CLI that is absent or unauthorised
  stays a stop; no hand-built API call, no token read from the environment.
- **a sub-agent pushing or committing.** Unchanged, and the mode does not reach it.

## Constraints

- The mode is carried by which command was invoked and by nothing else. There is no state
  on disk saying a run is unattended, so nothing can be left switched on.
- **A skill is self-sufficient once installed.** Each stop-owning skill states what
  unattended does to its own stop, in its own file — a skill that only behaves correctly
  when the orchestrator remembers to tell it is a skill that behaves incorrectly when it
  is invoked directly.
- `evidence` holds unchanged at every phase: nothing reported that was not observed, no
  gate in the background, no test edited into silence.
- The mode changes what the flow *waits for*, never what it *writes*. Proposal, spec,
  plan, commits and report all still exist, because they are what makes an unattended run
  reviewable at all.

## Prior decisions

- **Its own command, not a flag on `/libretto-flow`** — the user's call, 2026-08-12,
  reversing the reading this spec was first written with. The argument against was
  `payload/spec.md` refusing "a second flow for small work", and it does not reach: that
  refuses a second *flow*, and this is a second *door* onto the same one. The argument for
  is the precedent already in the payload — `/libretto-next` is its own command for
  exactly this reason, because an invocation that behaves differently has to be
  unmistakable in the history rather than an argument that can be typed by accident. What
  the refused rule does still bind is the shape: `/libretto-attacca` delegates to the same
  skills and describes none of them, so there is one description of each phase and it is
  the phase's own.
- **`attacca`, because it is the instruction and not a metaphor for it.** It is what a
  score writes to say *go on to the next movement without pausing*, which is the mode
  exactly. The cost is that it means nothing to a reader who does not read music, and one
  line of `description:` frontmatter pays it. The rejected alternatives were `auto`, which
  says something is automatic without saying what still stops, and `solo`, whose musical
  sense points at a prominent voice rather than at an absent pause.
- **A question the flow cannot derive is assumed, recorded and carried on from** — the
  user's call, 2026-08-12. The two rejected alternatives were stopping to ask, which makes
  the mode unreliable in exactly the case it exists for, and stopping without asking,
  which trades a wait for a dead run. The mechanism is not new: the flow already turns an
  unsettled post-plan question into a finding with what was assumed and what changes if
  the assumption is wrong. Unattended moves that rule earlier rather than inventing one.
  **Ceiling named:** an assumption is only as visible as the report carrying it, so a run
  whose request nobody reads has bought silence rather than speed. The replacement, the
  day that bites, is refusing to open the request when an assumption was made — not a
  prompt mid-run.
- **Push and the request are answered by the invocation, not assumed.** "Never push
  unasked" is intact: the asking happened when the command was typed. What the mode may
  not do is widen that consent — it covers this branch, this request, and nothing past it.
- **A gate is not a stop and cannot be answered.** A stop is a place where the user
  changes something; a gate is a place where the code is measured. The mode is defined
  over the first list only, which is why the classification is the spec and the command is
  the detail.
- **With a task argument, that task; with none, home first.** Both are already the flow's
  rules — `/libretto-flow <task>` does that task, and phase 1 asks its three sources in
  order. Phase 1's in-flight choice is the one exception the mode resolves rather than
  inherits: it continues the oldest change with open boxes, because "home first" is the
  order the flow already argues for and picking the newest would need a reason nobody has
  given.

## Task breakdown

- [ ] `commands/libretto-attacca.md` — the new command: the classification it carries —
      what the invocation answers, what stays a stop, what an underivable question becomes
      — and delegation to the phases, describing none of them.
- [ ] `commands/libretto-flow.md` — one line naming the other door, the way it already
      names `/libretto-status`.
- [ ] `skills/find-work/SKILL.md` — the in-flight choice and the which-task choice under
      unattended.
- [ ] `skills/write-spec/SKILL.md` — the contract stop, and step 4's question becoming a
      recorded assumption under prior decisions.
- [ ] `skills/write-plan/SKILL.md` — the order stop.
- [ ] `skills/record-work/SKILL.md` — push and the request answered by the invocation, and
      what the request description must carry.
- [ ] `skills/present-work/SKILL.md` — the report names what the mode answered and every
      assumption it made.
- [ ] `docs/FLOW.md` — the mode beside the stops it removes.
- [ ] the delta applied onto `.agents/specs/payload/spec.md`, and this folder deleted.

## Verification criteria

- every skill and command referenced by the new prose exists, and nothing invokes a path
  that does not get installed
  Proof: scripts/check-payload
- frontmatter still parses and `name:` still matches, for every file touched
  Proof: scripts/check-payload
- every `Proof:` citation in this spec and in the capability spec it lands on resolves,
  file and test name
  Proof: skills/record-work/spec-drift --anchors

**These check the payload, not the mode.** A skill is a prompt and a prompt is verified by
running it — `payload/spec.md` says so about itself, and it is more true here, because
what this change ships is entirely a set of instructions about when not to wait.

So the following are **claims until a run observes them**, and they are the ones worth
running first:

- an unattended run passes the spec stop and the plan stop without prompting, and both
  artifacts exist on disk afterwards
- it ends at a pushed branch with a request open, and the request's description names what
  the mode answered
- a failing gate stops it, on its branch, with no request opened — the case where the mode
  must look exactly like the attended flow
- an absent `gh` stops it with the install line and nothing else, the same stop phase 8
  already exercised for real
- a question it could not derive appears as an assumption in the spec, the report and the
  request, and nowhere as a prompt
