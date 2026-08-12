# keep-the-readme-in-step — the stops are asked natively

Targets: payload

Every stop in the flow, and every question inside one, is asked with `AskUserQuestion`.
None of them is a sentence at the end of a report.

**The rule already exists and covers exactly one of them.** `record-work` says it about
push: *"a question in prose is a paragraph the reader can skim past, and the flow then sits
waiting for an answer to something that read as a summary"*. That argument was never about
push. It is about every place the flow waits, and phases 2 and 5 wait the same way while
asking in prose.

## Outcomes

1. **Phase 2's stop is a native question.** After the spec is written and reported, the
   go-ahead is asked with `AskUserQuestion` — carry on to the plan, or change the contract
   first.

2. **Phase 5's stop is a native question.** Same shape: start the work, or change the order
   first.

3. **Phase 1's in-flight choice is a native question.** It already says *never choose*; what
   it does not say is *how to ask*, and the which-task question three sections below it
   already uses the native prompt. One skill asking two ways is the inconsistency that
   makes the prose version look acceptable.

4. **Every one of them carries the same shape as phase 8's:** the recommended option first,
   saying what will actually run; the real alternatives; and room to answer differently.

5. **`/libretto-attacca` still asks none of them.** It answers stops in advance — that is
   unchanged and untouched by this. What changes is how the attended flow asks, not how
   many times it asks.

## Scope boundaries

**In:** how the three stops and phase 1's in-flight choice are asked. The four skills that
own them, the command that routes them, and the flow document that argues for them.

**Out:**

- **the number of stops.** Still three. This is about the shape of a question, never about
  adding one, and a change that moved the count would be a different spec.
- **the trivial lane.** A change with no spec collapses phases 2 and 5 entirely; there is
  no question to ask natively when there is no question.
- **the phase 7 report.** It stops for nothing and gains no prompt. Its content is a
  report, and a report with a button on it is a stop wearing a disguise.
- **wording the options.** Each skill writes its own; pinning the labels here would put the
  copy in two places.
- **anything a sub-agent asks.** Sub-agents do not ask — unchanged. They return findings
  and the orchestrator asks.

## Constraints

- `AskUserQuestion` is native. **Do not build a prompt system** — `docs/FLOW.md` already
  says this and it holds here.
- A skill states *that* its stop is asked natively and what the options mean. It does not
  restate `record-work`'s argument for why; that argument lives once.
- Every skill stays self-sufficient: each says it in its own file, because a skill invoked
  directly gets no orchestrator to remind it.

## Prior decisions

- **The native prompt for every stop** — the user's call, 2026-08-12: *"todas las preguntas,
  incluso la de pasar de fase, con la pregunta de Claude en el CLI"*. The evidence was in
  front of both of us: the run that specced this passed both of its stops as a sentence at
  the bottom of a report, which is the exact failure `record-work` describes.
- **The reason generalises rather than being re-argued.** `record-work` made the case for
  push in 2026-08-10 and it was never push-specific: a question in prose is skimmable, and
  a flow waiting on an answer nobody realised was a question is a flow that looks hung.
- **Ceiling named:** this is a rule in prose asking a model to reach for a tool, and nothing
  mechanical checks it — `scripts/check-payload` cannot tell a native prompt from a
  paragraph. The replacement, the day it drifts, is a check that every skill owning a stop
  names `AskUserQuestion`, which is a string search and would have caught this one.

## Task breakdown

- [ ] `skills/write-spec/SKILL.md` and `skills/write-plan/SKILL.md` — each stop asked
      natively, with what its options mean.
- [ ] `skills/find-work/SKILL.md` — the in-flight choice, the same way its sibling
      which-task question already is.
- [ ] `commands/libretto-flow.md` and `docs/FLOW.md` — the stops are native, said once
      where the stops are argued for.

## Verification criteria

- every referenced skill and command still exists, and frontmatter still parses
  Proof: scripts/check-payload
- every `Proof:` citation in the spec this lands on resolves, file and test name
  Proof: skills/record-work/spec-drift --anchors

**Neither of those checks the behaviour**, and there is no honest test that would — a skill
is a prompt, and whether a model reaches for the native tool is observed by running the
flow, not by a Go test. The ceiling above names the check that *would* have caught this
one, and it is deliberately not built yet: a string search for `AskUserQuestion` in five
files is a guard against a failure that has happened once.

Claim, not fact, until a run observes it: a `/libretto-flow` that stops at phase 2 presents
a native question rather than a paragraph.
