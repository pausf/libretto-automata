---
name: converge
description: "Trigger: a checklist stopped describing the code; open boxes that may already be built, or built otherwise; 'converge', 're-derive the boxes', 'the plan drifted'. Launches one fresh converge-reader over the tree as it is, runs the cited proofs itself, closes only what it watched pass, reports every other open box as open, drifted or unmapped, appends a box for each of the change's own uncovered criteria whose proof fails, and commits the checklist alone. A side door, not a phase."
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

Launch **exactly one `converge-reader`** and read its return in full — the host may run it
in the background, and nothing proceeds until its return is read. Its
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

## Step 3 — The anchors, before any proof

Run **`spec-drift --anchors` before any proof** — the script beside the sibling
`record-work` skill, at `<skill-base>/../record-work/spec-drift` from this skill's announced
base:

```
<skill-base>/../record-work/spec-drift --anchors > /tmp/converge-anchors.out 2>&1; echo "exit=$?"
```

Read the file. Any citation of this change it names as unresolved — a file that is gone, a
test name that does not exist — is not passed, and is not run.

This comes first because the run that follows cannot see it: a test name that does not exist
exits 0 under `go test -run`. The anchors catch the name; the pass rule below catches the run.

## Step 4 — Closable boxes

A box the reader returned closable is a box worth running, nothing more. For each:

1. **Map it.** Match the box's `Closes:` text against the contract's criteria — the
   `## Verification criteria` of `proposal.md` and of each `spec-<capability>.md` beside it,
   or an old-layout change's delta. Each criterion it names must match exactly one. No
   match, more than one candidate, or no `Closes:` at all — as an old-layout box may have —
   and the box stays open, reported as `unmapped`. A mapping is never guessed, and the
   reader's opinion of one is not a mapping: the user sharpens the `Closes:`.
2. **Run it.** Take each mapped criterion's `Proof:` line from the contract itself — never
   from the reader's quote of it — and run it in the foreground the way `work-reviewer`
   runs a citation: output to a file, the exit read, the file read. Two kinds, two rules:
   - **test-named**, `<file> TestX` — run with `-v`. It passes only when the output
     carries `--- PASS: <TestName>` for that test, whatever its exit status. A name that
     does not exist prints `PASS` and no such line.
   - **file-only**, a command with no test name — it passes when it runs to the end,
     exits 0 and prints no FAIL line.
3. **Mark it.** Every citation of every mapped criterion passed: turn its `[ ]` into `[x]`
   and append, on the box's last line, `— converged <date>: <citation>, <citation> passed`,
   naming each one that passed. A box wraps over several lines, and its last one is where a
   reader finishes it. Anything short of that and the box stays open; the report names the
   citation that did not pass and what it printed.

The reader cannot run anything, and a reader's verdict never closes a box. The checklist has
one writer, and the writer is the one that must have watched the pass.

## Step 5 — Drifted boxes

A drifted box is not touched — not its text, not its state, not a note beneath it:
**drift is reported, never rewritten.** The report carries the box's own words
**beside the reader's file and line** for what the code does instead.

A box traces to the contract. Rewriting it to match the code ratifies the drift without the
contract moving: an edit to what the change promises, made by the back door. Resolving drift
is the user's, at the contract. Say so in the report, and say that `loop` would build a
drifted box as written, against the code.

## Step 6 — Uncovered criteria

Here **uncovered means the change's own criteria**: those under `## Verification criteria`
in `proposal.md` and in each `spec-<capability>.md` beside it — an old-layout change's delta
in their place — that no box's `Closes:` names. Never a criterion of a capability spec the
change targets: that spec is what is true already, and its criteria belong to changes that
landed. Counting them would turn a checklist tool into a regression sweep.

For each `uncovered` line the reader returned, check it against that definition, then run
the criterion's `Proof:` line from the contract by the rules of Step 3 and Step 4.

- **Not passed**: one open box, appended after the checklist's last box, in this shape —

  ```
  - [ ] **<criterion, in a phrase>** (converge <date>). Closes: <criterion> on <citation>. Waits on: nothing.
  ```

  One box per criterion, marked `(converge <date>)` so the next reader knows where it came
  from.
- **Passed**: the criterion passes already, so no box is appended — a task whose criterion
  is met is not a task — and the report says so.

A `gap` the reader names — an outcome or a scope item no criterion covers — is reported,
never boxed. A box names the criterion that closes it, and a box with none is a box nobody
can close; writing the missing criterion is editing the contract, and that is phase 2's.

## Step 7 — The writes

converge writes to the checklist and nothing else: the marks of Step 4 and the boxes Step 6
appends. Not the contract, not `## Approach`, not the decision log, `Targets:` or
`Durable decisions:` — every pillar stays as it was.

It never unchecks a closed box, and never re-runs one. The ask is about open boxes; a
closed box whose proof fails today is a regression, and `libretto metrics` already counts a
reopened box as the sign of one.

## Step 8 — The commit

A mark never committed was never marked. Anything marked or appended: stage the checklist
file alone and commit it, with that file as the only pathspec.

```
git add -- <checklist>
git commit -m 'docs(<change>): converge the checklist' -- <checklist>
```

The commit is `docs(<change>): converge the checklist` and it stages nothing else — a
converge commit is reviewed alone. Step 0 already stopped on anything staged, because a
staged file is a tracked modification.

If nothing was marked or appended, no commit.

converge never pushes and never creates a branch: the push is the user's, asked at
`record-work`.

## Step 9 — The report

One line per open box, in checklist order, in exactly one state:

- **open** — with what the reader found missing, or the citation that did not pass
- **closed** — with every citation that passed
- **drifted** — the box's words beside the reader's file and line
- **unmapped** — with what its `Closes:` failed to match

Then the proofs run, one line each, pass or fail. Then:

- each uncovered criterion, with its proof's result and the box appended for it, or
  "passes already"
- each `gap`, reported as a contract gap for phase 2
- anything the reader marked **unverified**, as a question, never as a verdict

Or the reader's explicit **"nothing found"**, relayed as a statement, not an absence.

Last, the commit's hash, or "no commit".
