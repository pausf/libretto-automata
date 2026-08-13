# ask-release-label-at-attacca-end — delta

Targets: payload

An unattended run ends at a request whose `require-release-label` check is red by design.
That red *is* the bump question reaching the user, and until now it reached them as a
broken pipeline: one alarm, one round trip, paid for a decision the run was standing
right next to. The retro of 2026-08-13 bought the cheap half — the closing report predicts
the red check. This is the other half: ask the bump natively, once, after the run is over.

**The run never decides the bump. It asks, and it types the answer.** That sentence is
what keeps this from reversing the non-goal it sits next to, and every criterion below
exists to hold it.

## Outcomes

After `/libretto-attacca` has pushed, opened the request and given its closing report —
and only then — the flow asks the bump once with `AskUserQuestion`: `release:patch`,
`release:minor`, `release:major`. The answer is applied with `gh pr edit --add-label`, and
the run confirms the label landed by reading the request back.

- **The question is the last thing that happens.** After the report, after the request is
  confirmed open, after the return to the base branch. There is nothing downstream of it,
  which is what makes it not a stop: the work is already complete and reviewable whether
  it is answered or not.
- **Unanswered ends exactly as today** — unlabeled, with the closing report's red-check
  line intact. That line is written before the question is asked and is never withdrawn by
  it: a report that promised a red check and then silently labelled the request has lied
  about the state the user will find.
- **A run with no `release:*` labels defined in the repository does not ask.** The
  question is offered only where the convention exists.
- **`release:major` is selectable and never recommended.**
- The request's description gains the bump the user chose and the fact that they chose it,
  alongside what the invocation answered and what the run assumed past.

## Scope boundaries

**In:**

- `/libretto-attacca`, at the end of phase 8, in `skills/record-work/`.
- Detecting whether the repository defines the three labels, and staying silent when it
  does not.
- Applying exactly one label, and reading the request back to confirm.

**Out, named so it cannot arrive quietly:**

- **The attended `/libretto-flow`.** Its phase 8 already stops with the user present and
  watching; the red check does not ambush anyone there. The user's call, 2026-08-13. Back
  the day an attended run pays the same round trip.
- **Merging, tagging, releasing, or creating a Release.** Untouched. This change reverses
  exactly one word of the standing non-goal — *labelling* — and leaves the other three
  where they are.
- **Choosing the bump when the question is unanswered.** No default, no "patch is safe".
  A silently-wrong bump is the failure that published `v1.0.0`; defaulting is that failure
  with a politer name.
- **Creating the labels when they are absent.** A run that invents `release:minor` in
  somebody else's repository has decided that repository's release convention.
- **Re-asking.** One question per run. A request that already carries a `release:*` label
  is not asked about again.
- **A flag, setting or env var to suppress the question.** Consent is per run, and the
  suppression already exists: do not answer it.

## Constraints

- **`libretto-attacca` describes no phase.** The payload spec forbids it restating what a
  skill owns, so the behaviour lives in `skills/record-work/` and the command references
  it at most. A second description of phase 8 inside the command is the drift that rule
  exists to prevent.
- **`gh` only** — never the REST API, never a hand-built call. The same rule that governs
  `jira`, and for the same reason: a missing or unauthenticated CLI is a stop, not
  something to route around.
- **`gh pr edit --add-label` fails on a label that does not exist**, which is why the
  detection above is correctness rather than politeness.
- **Headless runs (`libretto loop`) have no interactive prompt.** `AskUserQuestion` cannot
  arrive, so the unanswered path is the normal path there and must be the quiet one.
- Matching is by whole label name. The workflow already refuses two `release:*` labels, so
  applying a second one turns a red check into a differently red check.

## Prior decisions

- **The non-goal is reversed in exactly one word, and the reversal is conditional.** The
  payload spec's *"an unattended run merging, tagging, releasing, or labelling the
  request"* becomes *"…merging, tagging, or releasing the request, or labelling it with a
  bump the user did not choose."* The bump stays a reading of `.agents/specs/`, made by a
  person. What the run gains is the typing, which `AGENTS.md` already assigns to it: *"The
  reading is yours. The typing is not."*
- **`release:major` keeps its ask-and-wait rule with no exception while this is `0.x`, and
  a native question satisfies it.** The rule forbids *announcing* a major and proceeding;
  a prompt the user answers is the asking the rule demands. But the run must never put
  major first: recommending it is the announcement, and announcing it three times is what
  published `v1.0.0` and `v1.0.1` off a table read without the paragraph above it.
- **Attacca only** — the user's call, 2026-08-13. What changes if it is wrong: an attended
  run keeps paying the round trip the unattended one stopped paying.
- **Detect the labels rather than hardcoding this repository** — the user's call,
  2026-08-13. What changes if it is wrong: a project whose `release:*` labels mean
  something else gets asked a question that does not apply to it, and answers by ignoring
  it.
- **The question comes after the report, never mid-run.** From the lesson of 2026-08-13:
  the mode's whole visibility budget is that closing report, and a question that interrupts
  it spends the budget on the interruption.
- **Ceiling named:** the question is only as reachable as the terminal it is asked in. A
  scheduled or piped run never sees it and lands unlabeled — the same state as today,
  which is why today's behaviour is the fallback rather than an error. The replacement, if
  that becomes the common case, is `gh pr comment` carrying the three commands; deliberately
  not built for a case that has not happened.

## Task breakdown

1. `skills/record-work/`: the bump question at the end of the attacca path — detection,
   the native prompt with major never first, `gh pr edit --add-label`, and the read-back.
2. `skills/record-work/`: the unanswered and no-labels paths, and the rule that the
   report's red-check line is never withdrawn.
3. `skills/present-work/` or `record-work`: the chosen bump and the fact that a person
   chose it, into the request's description beside the invocation's answers.
4. `commands/libretto-attacca.md`: amend the `Never` list so *labelling* reads as its
   conditional self rather than an absolute the flow now breaks.
5. `scripts/check-payload`: the wiring check — the decisive words are in the file that
   owns them, and the command does not restate the phase.
6. `.agents/specs/ci/spec.md`: the prior decision naming this as queued becomes the
   decision naming it as built.

## Verification criteria

- **the bump question exists in `skills/record-work/`, on the attacca path, after the
  push and the request are confirmed, and the label is read back off the request** — and
  not in `commands/libretto-attacca.md`, which describes no phase. The read-back is named
  here rather than left to the task list because it is the same rule the push already
  carries: a command that printed no error is not a change the forge accepted
  Proof: scripts/check-payload
- **`release:major` is never the recommended option.** A static read of the file that owns
  the prompt: major is present and is not first
  Proof: scripts/check-payload
- **the no-labels and unanswered paths are stated**, and the closing report's red-check
  line is required to survive the question
  Proof: scripts/check-payload
- **the command's `Never` list no longer contradicts the skill.** Labelling reads as
  conditional on the user having chosen the bump; merging, tagging and releasing stay
  absolute
  Proof: scripts/check-payload

**Proved as wiring, and the difference is stated rather than assumed.** A prompt's conduct
is checked by running it; what a script can check is that the phase exists where it should,
that the words carrying the constraint are still in the file that owns them, and that the
command has not grown a second copy. `review-spec` deleted from the flow once left
`check-payload` green, which is why every payload criterion says this out loud.
