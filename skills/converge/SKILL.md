---
name: converge
description: "Trigger: a checklist stopped describing the code; open boxes that may already be built, or built otherwise; 'converge', 're-derive the boxes', 'the plan drifted'. Launches one fresh converge-reader over the tree as it is, then reports every open box as open, closable or drifted. A side door, not a phase."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Re-derives one change's open boxes against the code as it stands, for a checklist that
drifted from what got built. `loop` drives a checklist; this re-derives one that stopped
describing the code.

It is a door the user opens, like `/libretto-status` — not a phase. Neither the flow nor
`attacca` routes here.

**The reader judges; the session decides.** Whether a box is done is a reading, and the
reading runs in a **fresh `converge-reader`** that never had this conversation — the session
running converge is often the one whose plan drifted. The reader has no shell, so nothing it
says is evidence.

`skills/evidence/` governs every step here.

## Step 0 — Is the tree clean?

```
git status --porcelain --untracked-files=no
```

Anything listed means **a dirty tree stops converge before the reader**: name the files, in
one line each, and stop. A reading of uncommitted code says nothing about what was recorded,
and a modified checklist file would mix converge's view with somebody else's edit.

Tracked files only. An untracked file — a generated `wiki.html`, a scratch note — never
stops converge: no proof runs against it unless something imports it, and stopping on one
would stop converge on every tree that carries one.

## Step 1 — Which change?

**Named** in the arguments: that one. Refused, in one line, if its folder holds no checklist
box.

**Empty**: invoke `Skill(skill="find-work")` — on another host, the host's own skill tool —
and run **source 1 only, never the queue**: its open-box scan and its unmerged branches. A
queued capture has no checklist, so there is nothing in it to converge.

- one change in flight: that one
- several: **ask which one** — one `AskUserQuestion` call, the host's native prompt, or a
  plain question in conversation where there is none — and start nothing until it is
  answered
- none: say so in one line and stop

**`find-work` defines the box** and which file holds the checklist; this skill restates
neither. Read the checklist from the file that skill's rule picks, counted by its rule.

## Step 2 — Launch exactly one reader

Launch **exactly one `converge-reader`**, foreground, and read its return in full. Its
world, all paths and no narration:

- the tree as it is now — the repository root, and **never a diff range**: drift is a
  property of the code that exists, wherever it was written
- the change folder, and the checklist file Step 1 found in it
- the capability specs the contract's `Targets:` lines name

It gets **none of the conversation** — not the session's account of what was built, not
which boxes the session thinks are done. Telling a reader what to find is priming the
witness.

One, not a panel: a second reader over the same tree is correlation dressed as
corroboration.

**Silence is not a verdict.** A reader that died or returned nothing did not read the
change. Say so and launch it again — once. A reader that cannot finish twice is a stopped
item, per `skills/evidence/`.

## Step 3 — Drifted boxes

A drifted box is not touched — not its text, not its state, not a note beneath it:
**drift is reported, never rewritten.** The report carries the box's own words
**beside the reader's file and line** for what the code does instead.

A box traces to the contract. Rewriting it to match the code ratifies the drift without the
contract moving: an edit to what the change promises, made by the back door. Resolving drift
is the user's, at the contract. Say so in the report, and say that `loop` would build a
drifted box as written, against the code.

## Step 4 — The report

This version reads and reports; it writes nothing and commits nothing.

One line per open box, in checklist order, in one of the states the reader returned —
**open**, **closable**, **drifted** — each with the reader's file or symbol. A closable box
is reported as closable and stays open: **the reader's verdict closes nothing.**

Then, as the reader returned them:

- each `uncovered` criterion, with its `Proof:` line
- each `gap` — an outcome or scope item no criterion covers — reported, never boxed:
  writing the missing criterion is editing the contract, and that is phase 2's
- anything the reader marked **unverified**, as a question, never as a verdict

Or the reader's explicit **"nothing found"**, relayed as a statement, not an absence.
