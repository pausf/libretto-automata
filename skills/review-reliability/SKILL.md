---
name: review-reliability
description: "Trigger: reviewing a diff, a file or an MR for what breaks at runtime — logic errors, edge cases, race conditions, unbounded work, leaked resources, error paths that lose data. Every finding names the input or state that triggers it. Standalone: works on any diff, with or without review-project."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Find what pages someone at 3am. Not style, not architecture — the concrete input or
state under which this change does the wrong thing, and what the wrong thing costs.

`evidence` governs: a bug reported is a bug whose trigger was traced, and what
could not be verified is said out loud, never silently assumed fine.

## The one question

Every candidate finding answers this before it is written down:

**What input or state makes this go wrong, and what happens then?**

"This could be fragile" is not a finding. "An empty list reaches this index on the
first run of the month" is. Trace the trigger through the code — reviewing the diff
means *reporting* on the diff, but *reading* whatever the trace needs.

## Where to look, in cost order

The expensive failures first, because one of these outranks any ten of the rest:

- **error paths that lose data** — a write that partially succeeds, an exception
  swallowed between the mutation and its confirmation, a catch that logs and
  continues as if the operation happened
- **read-then-write races** — two requests interleaving between the check and the
  act; anything that computes on stale state and writes it back
- **unbounded work** — a loop over a collection the user can grow, a query without
  a limit fed to a loop of queries (N+1), memory that accumulates per request,
  missing pagination on something production will make large
- **resources that never close** — handles, connections, locks, temp files on the
  error path; the happy path almost always cleans up, the sad one is the finding
- **the classic borders** — off-by-one, empty input, the zero, the negative, the
  duplicate, null where the author assumed presence, timezone and encoding at any
  boundary that crosses systems, float equality near money

**Check whether it is already handled.** A guard two calls up, a database
constraint, a framework retry — a finding the codebase already answers is noise,
and noise teaches the reader to skim. Look before flagging.

## Severity is impact, not category

- **critical** — data loss or corruption, or a crash any user can trigger
- **high** — wrong behaviour on realistic input, or an outage under realistic load
- **medium** — wrong behaviour on unusual-but-reachable input
- **low** — real, reachable, marginal cost

A race that needs two admins clicking in the same millisecond is low however
exciting it sounds. An off-by-one on every request is not low however boring.

## Output

Per finding: location (`file:line`), the trigger (input or state, concrete), what
happens, severity, and the fix in one line.

Then — this is not optional — **what was not verified**: the paths that could not
be traced, the external calls whose behaviour was assumed, the load that was
guessed at. A review that lists its blind spots is worth more than one that claims
none.

No findings is a statement: **"nothing I could make break, and here is what I
could not check."** Never an invented issue to fill the page.

This lens reports; it never blocks, edits, commits or pushes.
