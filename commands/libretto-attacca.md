---
description: Run the flow without stopping — attacca, the score's instruction to go on to the next movement without pausing
---

The flow in `docs/FLOW.md`, run end to end without waiting at a question, ending at a
pushed branch with a request open on it.

Input: `$ARGUMENTS` — a Jira issue key, an issue URL, or empty.

**This command knows exactly one thing `/libretto-flow` does not: what the invocation
already answered.** Everything else it delegates, and it describes no phase — the phases
are `/libretto-flow`'s to route and each skill's to explain. Two descriptions of one phase
is one description too many, and the maintained one is never the copy.

## What the invocation answers, and what it cannot

**A stop is a place where the user changes something. A gate is a place where the code is
measured.** This command is defined over the first list only. That line is the whole of
it, and getting it wrong is the only way this becomes dangerous: a mode that answers a
gate is not unattended, it is unverified.

| At | Under this command |
|---|---|
| 2–3 · the contract stop | **answered** — the spec is written, committed and reported, and the flow carries on |
| 5 · the order stop | **answered** — same |
| 8 · push and open the request | **answered yes**, both halves |
| a failing gate, twice on one task | **stops** — and two stopped tasks stops the session, unchanged |
| a missing or unauthenticated `jira`, `gh`, `glab` | **stops** — the input never arrived; there is nothing to assume |
| a question no reading of the code settles | **assumed, recorded, carried** — see below |

The user typed this command and not the other one. That is the answer, given in advance,
and it is why "never push unasked" is intact rather than overridden — the asking happened
at the prompt. **The consent covers this branch and this request and nothing past it.**

**The request's description says which of its decisions a person made.** What the
invocation answered — the two stops and the push — and every question the run assumed past.
Without the first half a reviewer cannot tell an agreed contract from an assumed one, and
that is the difference this command is entirely built out of.

## A question it cannot derive becomes an assumption, never a prompt

When a phase hits something no reading of the code settles — a product tradeoff, two live
precedents, anything phase 4 would have asked — it takes the option it would have
recommended and writes it down in three places:

1. the spec, under prior decisions, **marked as assumed** and naming what changes if it is
   wrong
2. the phase 7 report
3. the request's description

Never a default left silent, and never a guess dressed as a decision. The flow already
does this after the plan, where an unsettled question becomes a finding; this moves the
same rule to the front rather than inventing a new one.

**Ceiling named:** an assumption is only as visible as the request carrying it. A run
whose request nobody reads has bought silence, not speed. The replacement, the day that
bites, is refusing to open the request when an assumption was made — never a prompt
mid-run.

## The phases

Standing rules unchanged: `evidence` at every phase, `ponytail` throughout, every phase
invoked even when its answer is "nothing here".

`Skill(skill="…")` below is Claude Code's spelling; on another host, load that skill with
the host's own skill tool. This command answers the stops rather than the widgets, so
where a phase would have used a native prompt there is nothing to substitute — but the one
question that survives, the bump at the very end, asks with the host's own native prompt
or in conversation where there is none.

```
Skill(skill="find-work")          1 · with a task given, that task; with none, home first —
                                      the oldest change with open boxes
Skill(skill="write-spec")         2–3 · then carry on
Skill(skill="review-spec")        3→5 · reads the contract, only when one was written.
                                      **Here it matters more than in the attended flow**:
                                      nobody is reading the spec at the stop, so this is
                                      the only thing between an unfailable criterion and
                                      a plan built on it. It still only reports — see
                                      below for who acts on it
Skill(skill="write-plan")         5 · then carry on
Skill(skill="build-and-check")    6
Skill(skill="review-work")        6→7 · it already fixes rather than asks
Skill(skill="present-work")       7 · naming what this command answered, and every
                                      assumption the run made
Skill(skill="record-work")        8 · push and open the request, both already answered
```

A change phase 2 decides needs no spec collapses the same way it always does. There was
nothing to wait for here either; what collapses is the writing, and phase 7 still says
everything.

### Who acts on `review-spec`'s findings

**The reviewer reports and never edits the spec — that rule does not bend here.** What is
different is who is left holding the report. In the attended flow it is the phase-3 author
at the stop; unattended there is no stop, so **this run is the phase-3 author** and amends
its own spec before phase 5 reads it.

That is not the same as the reviewer fixing its own findings, and the distinction is the
one that keeps the review worth running: a reviewer that edits what it reviews has no
independent reading left to report. Reviewer reports, author amends, and unattended those
are two roles held by one run rather than two.

**A finding acted on is named in the phase 7 report and in the request**, exactly like an
assumption — for the same reason. A contract quietly rewritten between review and plan is
a contract nobody agreed to.

## Never

- **merge, tag or release the request, or label it with a bump the user did not choose.**
  It ends at a request open for review. The bump is a reading of the specs rather than of
  the commits, and a version number cannot be recalled once the proxy has cached it — so
  the reading stays the user's, always. What the run may do is **type** an answer it was
  given: phase 8 asks once at the very end, and `skills/record-work/` owns every rule
  about it. Unanswered, headless, or a repository that defines no `release:` labels all
  end here unlabeled, exactly as before.
- **skip, reorder or soften a gate**, or reach for `--force`, `--no-verify` or anything
  else that buys a green result. This removes waits, never checks.
- **weaken, skip or delete a failing test.** Fix the cause or stop and say why.
- **drain the queue.** One invocation is one piece of work, exactly as `/libretto-flow
  <task>` means that task. `/libretto-next` picks one and asks.
- **assume past a missing credential**, with a hand-built API call or a token found in the
  environment. That turns a stop into an exposure.

## When it stops

It stops on its branch, with everything committed that was finished, and **no request
opened**. Then it reports what it observed, what it was doing when it stopped, and what
would let it continue. A stopped run that opens a request anyway is asking someone to
review work that was never proven.
