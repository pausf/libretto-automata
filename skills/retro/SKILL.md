---
name: retro
description: "Trigger: running a retro on the flow; spending the lessons ledger; 'retro', 'lecciones', 'learn from the flow', 'que no vuelva a pasar'. Reads .agents/lessons.md, fixes what belongs to the project, proposes what belongs to the payload. Never applies a payload diff."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Spends the lessons ledger, so the same correction is never paid twice.

`evidence` captures every user correction into `.agents/lessons.md` while the flow
runs. This skill is the other half: read the open entries, put each fix where it
belongs, and mark the entry so it is never re-classified.

It runs where the flow ran, on that project's ledger, and nowhere else.

## Read the ledger

`.agents/lessons.md`, at the project root's `.agents/`. Entries look like:

```
## 2026-08-13 · add-relative-discounts · build-and-check
Said: the tag here is <type>/<scope>, never bare
Did: wrote a bare tag in the commit subject
```

An entry with a `Resolved:` line is spent — skip it. **No ledger, or no entry
without a `Resolved:` line, means there is nothing to do: say so in one line and
stop.** That is a state, not an error, and it costs no wait.

## Classify each open entry

Read `Said:` and `Did:` and ask one question: **where does the fix live?**

| Type | The signal | The fix |
|---|---|---|
| **project knowledge** | the flow lacked a fact of *this* project — a convention, a naming rule, a constraint | record it in this project's contract |
| **flow defect** | the payload skill that was active led the work wrong, and would lead any project wrong | propose the exact diff to that skill |
| **one-off** | a mistake nothing systematic produced, not worth preventing | mark it and move on |

The distinction is load-bearing. A retro that writes one project's manias into the
payload "fixes" the flow for this repository and breaks it for every other one.
When an entry could read both ways, it is project knowledge until the same lesson
shows up somewhere a second time — the cheaper wrong guess.

## Apply, propose, or close

**Project knowledge** — write the convention where this project keeps its contract:
the prior decisions of the capability spec whose `Governs:` owns the corrected
path, or `AGENTS.md` when no capability does. One entry, dated, in the project's
own words. Then mark the ledger entry:

```
Resolved: 2026-08-14 · convention recorded in .agents/specs/cli/spec.md
```

**Flow defect** — name the payload skill, and put the exact diff in the report:
the lines as they are, the lines as they should be. **Never apply it.** The payload
is the product; it does not get edited as a side effect of a retro, and the skill
being proposed against may not even be writable from here — it is installed,
symlinked from its own repository. Mark the entry resolved only when the user says
what they did with the proposal; until then it stays open and the next retro
proposes it again.

**One-off** — mark it resolved with that reading, so the next retro does not spend
context re-deciding it:

```
Resolved: 2026-08-14 · one-off, not worth preventing
```

## The rules

- **Append `Resolved:` lines. Change nothing else.** The ledger is history; an
  edited entry is history that lies.
- **Never push, never commit on the user's behalf.** The marks and any project
  contract edits land like any other work: shown, then committed when the user's
  flow commits.
- **Never touch the payload**, not even a typo, not even invited by the diff being
  small. The proposal is the whole power this skill has over it.
- One pass per invocation. A classification the user disagrees with is corrected by
  them and — like every correction — is itself a lesson.

## Report

- entries found: how many open, how many already resolved
- per entry, one line: its header, the classification, and what was done —
  *recorded in `<file>`*, *proposed (diff below)*, or *one-off*
- every proposed payload diff, in full, after the table
- what the metrics will now show: the count only moves when `evidence` captures,
  so a quiet ledger is a flow that got corrected less — or a capture rule that
  needs proposing against. Say which reading the entries support.
