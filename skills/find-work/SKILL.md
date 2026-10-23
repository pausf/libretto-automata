---
name: find-work
description: "Trigger: starting the Libretto flow; asking what to work on; resuming a half-finished change; reading a Jira task, issue URL or board URL; a request made in conversation. Phase 1 — finds the work wherever it came from and loads exactly one piece of it."
license: MIT
metadata:
  author: pausf
  version: "1.4"
---

## What this does

Phase 1 of the Libretto flow: **find the work**, and load exactly one piece of it.

It stops there. It does not write a spec, a plan, or any code.

## Phase 1 does not end until its artifact exists

**Reporting what was found is not finishing.** A source-3 request leaves
`.agents/changes/<name>/proposal.md` **on disk** before this phase reports, and a phase
that only described the work has not run.

The rule was written below, in prose, under source 3. It read as a description of what a
proposal contains rather than as a step that must happen, and it was skipped — the ask
was reported, confirmation was requested, and no file existed. The user had to ask where
the change was. **A step stated as prose is a step that gets read and not done**, so it
is stated here as an output instead:

| Source | Phase 1 has not finished until |
|---|---|
| 1 · in flight | the open boxes and what can start now have been read off the checklist |
| 2 · tracker | the issue was fetched and its description read — not just its header |
| 3 · what the user said | **`proposal.md` exists**, with `Tracker: none` and the ask verbatim |

Writing that file means committing it, and committing means a branch. **So phase 1
creates the branch when it writes the proposal** — the same rule as
`skills/build-and-check/` step 0, arriving earlier because the first write does. Phase 6
then *ensures* a branch rather than creating one, and finds the work already on it.

Then state the reading and **carry on into phase 2 in the same turn**. Stating it comes
after the file, not instead of it: a reading recorded nowhere is re-derived next session
and does not come out the same. But stating is not asking — the place to disagree with
the reading is the spec, which is the next thing written and the first thing this flow
stops on.

**Phase 1 does not wait.** It used to, and the wait bought a round trip: the user read a
paraphrase of their own sentence back and said yes. What they actually want to correct is
the contract, and the contract does not exist yet.

**One exception, and it is not ceremony:** when source 1 finds work already in flight,
choosing between it and something new is the user's, and there is no work to do until
they answer. That is the input arriving, not a phase transition. Source 2's missing or
unconfigured tracker stops for the same reason — nothing downstream exists.

The flow does not begin at a tracker. A tracker is one of three places work comes
from, and in practice the least common:

| Source | When |
|---|---|
| **a change already in flight** | a change's checklist under `.agents/changes/` has an open box |
| **a tracker key or URL** | one was given |
| **what the user said** | anything else — the request *is* the input |

**Asked in that order, and the order is the point.** Starting something new while a
change sits half-finished is how the half-finished thing gets abandoned, and the cost
is not the wasted work — it is a `.agents/changes/` directory nobody trusts any more.

## Source 1 — a change already in flight

Look at home before looking anywhere else.

```
bash <<'SCAN'
scan() {  # ${1}: ' ' for the open scan, 'xX' for the closed one
  for d in .agents/changes/*/; do
    for f in proposal.md tasks.md plan.md; do
      [ -f "$d$f" ] || continue
      awk -v doc="$([ $f = proposal.md ] && echo 1)" -v want="${1}" '
        /^[[:space:]]*```/ { fenced = !fenced; next }
        fenced { next }
        doc && /^##[[:space:]]/ { listed = ($(0) ~ /^##[[:space:]]+Checklist[[:space:]]*$/); next }
        doc && !listed { next }
        /^[[:space:]]*[-*][[:space:]]+\[( |x|X)\]/ { all++; if (index(want, substr($(0), index($(0), "[") + 1, 1))) n++ }
        END { if (n) print FILENAME ":" n; exit !all }' "$d$f" && break
      [ $f = proposal.md ] || break
    done
  done
}
scan ' '; scan xX
SCAN
```

**Spelled `${1}` and `$(0)`, run through `bash`, on purpose.** Claude Code substitutes a
bare dollar-digit in a skill body with the words of the invocation's arguments, so the
plain spelling rendered as `want="a"` under `/libretto-attacca Add a …` and the scan read
nothing. And zsh aborts an unmatched glob, so an empty `.agents/changes/` errored instead
of reporting nothing in flight. Both measured in a trial run of the flow.

**Three names, because the checklist moved twice.** It was `plan.md` until 2026-08-17,
then `tasks.md`, and now it is the `## Checklist` section of `proposal.md`. A change
created before each move still carries the old layout, and a search that saw only the
newest name would report the work in flight as nothing in flight — which is precisely the
failure this source exists to prevent, arriving through the fix for a different one.

**A box is a line matching `^\s*[-*]\s+\[( |x|X)\]`, and a box inside a fence does not
count.** That is the CLI's definition, and the scan above is it. In `proposal.md` only the
boxes in a section opened by a `## Checklist` heading and closed by the next `^##\s`
heading count; in `tasks.md` and `plan.md` every box outside a fence does. A bare count of
`- [ ]` lines over the whole file is not this: it counts the boxes in a fenced example,
and in `proposal.md` a breakdown or a quoted ask that happens to look like a task.

**`proposal.md` is the checklist when its `## Checklist` holds a box** — otherwise
`tasks.md`, otherwise `plan.md`. A file existing never decides it alone, because every
change in every layout has a `proposal.md`, and its existence would hide every old
change's boxes. A queued capture has no `## Checklist`, so it falls through to the third
row below.

### Read both scans, not just the first one

Each scan prints a line only for a change whose count is not zero, and stays silent for
the rest. That asymmetry is not a quirk of the output — **it is the second question this
source answers**, and until 2026-08-18 nothing said so, which is how two changes shipped
with their landing half-done and phase 1 walked past them twice.

| In the open scan | In the closed scan | What it is |
|---|---|---|
| yes | either | **work in flight** — the report below |
| **no** | **yes** | **finished and not landed** — every box closed and the folder still here |
| no | no | no checklist yet: a queued idea, a trivial-lane capture, or a change not cut. Say nothing |

**The middle row is reported by name, and it is not work in flight.** There is nothing to
continue — the boxes are all closed. What it means is that `record-work`'s landing did not
finish: the delta may have been applied, and the change folder was never deleted. Say the
change is **finished and not landed**, and leave what to do about it to the user.

Do not go and check whether the delta was actually applied. That means reading a capability
spec and deciding whether a delta is *present* in it, which is a reading — and a wrong one
either accuses a correct landing or clears a broken one. Zero open boxes is the signal; the
user interprets it.

A change with open boxes is work waiting. For each one, report:

- its name, and what its `proposal.md` says it is for in one line
- how many boxes are open out of how many
- **what its checklist says can start now** — it records dependencies, and dependencies
  nobody reads make it a list rather than a plan

Then ask: continue one, or begin something else. **Never choose.** Picking up somebody's
half-finished work without asking is a decision about their priorities.

**Ask it with `AskUserQuestion`** — or in conversation where the native prompt does not
exist; the stop is the rule, not the widget — each open change as an option, the one with the fewest
boxes left recommended, and room to start something else. The same native prompt the
which-task question in `tracker.md` uses: one skill asking two ways is what makes the prose
version look acceptable everywhere else.

**Under `/libretto-attacca` this one is already answered:** with a task given, that task;
with none, the oldest change with open boxes. Nothing else changes here — a tracker CLI
that is missing or unauthorised still stops, because that is the input failing to arrive
and no invocation can supply it. The command carries why.

No `.agents/changes/` directory, or none with open boxes, means nothing is in flight.
That is a state, not an error — say it in one line and move on.

### A branch is also work in flight

**The scan above cannot see the trivial lane.** Source 3 writes `proposal.md` before phase
2 decides the lane, so a change that needed no spec does get a folder — a capture with no
`## Checklist`, which the scan reads as the third row, and which the trivial lane deletes in
its own commit. After that it lives only as commits on a branch, and the scan reports
nothing in flight while the work sits there unmerged.

That is not hypothetical: the lane's own first run produced a commit on a local branch,
and the next phase 1 reported an empty house.

So look at the branches too, and read the result rather than assuming it:

```
git branch --format='%(refname:short)' --no-merged main
git log --oneline main..<branch>
gh pr list --json number,headRefName,state      # or glab mr list
```

Report each unmerged branch with its commits and whether a request is open for it.
**Unpushed and un-requested is the state worth naming** — that is work nobody but this
machine has.

`main` here means whatever the base branch actually is. Do not hardcode it if the
repository says otherwise.

## The queue — captured, never started

**When the caller asks for the queue, read `queue.md` beside this file.** That is
`/libretto-status`, `/libretto-next`, or any run in reporting mode. The queue is not a
source and never blocks one: a captured idea costs nothing to abandon.

## Source 3 — what the user said

**A request in conversation is a legitimate input, not a fallback.** It is how most work
arrives; every change in this repository so far arrived that way.

There is no key, so two things have to be produced rather than looked up:

- **the change's name**, from the request, verb-led and readable — `add-relative-discounts`,
  never an invented ticket id. A fake key implies a tracker that could be consulted.
- **`proposal.md` recording `Tracker: none`** and what was asked, in the words it was
  asked in. Paraphrasing a request loses the part you did not understand yet.

State the reading before phase 2 starts, and keep going. A request understood slightly
wrong becomes a spec that is confidently wrong — but the spec is where that gets caught,
and phase 2 stops for exactly this. Saying the reading out loud is what makes the wrong
one visible; waiting for a yes is what makes it expensive.

### Say whether it is a bug, and say it here

**A bug arriving as generic work is how the contract stops growing.** Whatever the source
— a tracker key, a branch, a sentence — the reading has to name it: *this is a report that
something already built behaved wrong*, as against something not built yet.

It costs one line and it changes phase 2's whole shape. A bug is a hole in the contract,
so the criterion that would have caught it comes before the fix, with a `Proof:` that is
red when it is written. Phase 2 knows how; it cannot apply it to work that never
announced itself as a failure.

**And record it in `proposal.md` in the reporter's own words** — the observed behaviour,
the expected behaviour, and what produced it. Paraphrasing costs the reproduction, and
without a reproduction there is no failing test, so there is no proof and the criterion
is a sentence somebody hopes is true.

`tracker.md` makes the same demand of a tracker issue: read its type, and do not infer
one from the summary. "Discount wrong on bundles" is a bug or a change of intent, and
those are opposite pieces of work.

## Source 2 — a tracker key or URL

**When a tracker key or URL was given, read `tracker.md` beside this file.** It covers
the `jira` CLI, its install and auth checks, which task to read, and how to read its
description and subtasks. A missing or unauthenticated CLI stops there, and the token
never comes through the conversation.

## Output

Report the work found, in one short paragraph: what it is, where it came from, and for a
tracker task what `tracker.md` lists. Then carry straight into phase 2. Nothing else
happens *here*: writing the spec is not this phase's job, but it is the same turn.
