---
name: find-work
description: "Trigger: starting the Libretto flow; asking what to work on; resuming a half-finished change; reading a Jira task, issue URL or board URL; a request made in conversation. Phase 1 — finds the work wherever it came from and loads exactly one piece of it."
license: MIT
metadata:
  author: pausf
  version: "1.2"
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
| 1 · in flight | the open boxes and what can start now have been read off the plan |
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
| **a change already in flight** | there are unchecked boxes in `.agents/changes/*/plan.md` |
| **a tracker key or URL** | one was given |
| **what the user said** | anything else — the request *is* the input |

**Asked in that order, and the order is the point.** Starting something new while a
change sits half-finished is how the half-finished thing gets abandoned, and the cost
is not the wasted work — it is a `.agents/changes/` directory nobody trusts any more.

## Source 1 — a change already in flight

Look at home before looking anywhere else.

```
rg -c '^\s*- \[ \]' .agents/changes/*/plan.md
rg -c '^\s*- \[x\]' .agents/changes/*/plan.md
```

A change with open boxes is work waiting. For each one, report:

- its name, and what its `proposal.md` says it is for in one line
- how many boxes are open out of how many
- **what its plan says can start now** — the plans record dependencies, and a plan whose
  dependencies nobody reads is a list, not a plan

Then ask: continue one, or begin something else. **Never choose.** Picking up somebody's
half-finished work without asking is a decision about their priorities.

No `.agents/changes/` directory, or none with open boxes, means nothing is in flight.
That is a state, not an error — say it in one line and move on.

### A queued idea is not work in flight

`/libretto-queue` writes proposals with a `Queued:` line and nothing else beside them —
no spec, no plan, no branch. They live in the same directory as changes in flight and
they are **not** source 1:

```
rg -l '^Queued:' .agents/changes/*/proposal.md
```

Report them, oldest `Queued:` date first, as their own short list, and point at
`/libretto-next`. **Never ask whether to pick one up, and never let them block.**

Home first exists so *started* work does not get abandoned. An idea costs nothing to
abandon — nothing has been built on it. Making four captured ideas stand between the user
and a Jira task would make capture punitive, and a queue that costs something to add to
is a queue nobody uses.

### A branch is also work in flight

**The scan above cannot see the trivial lane.** A change that needed no spec never creates
a `changes/` directory — by definition it has no spec and no plan — so it lives only as
commits on a branch, and `rg` over `plan.md` reports nothing in flight while the work sits
there unmerged.

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

## Source 2 — a tracker key or URL

Everything below goes through the `jira` CLI (ankitpokhrel/jira-cli). Never reach for an
MCP server or the REST API directly.

## Rules that hold for every command here

**Never pipe a `jira` command into `head`, `tail` or anything else when the exit
code matters.** The pipeline reports the last command's status, so a failed `jira`
call looks like success. Redirect to a file, check `$?`, then read the file.

**`--plain` is not plain.** It emits ANSI colour escapes that pollute the
transcript. Use `--raw` — it returns JSON — and parse with `jq`. Plain mode is for
lists a human will skim; add `--no-truncate` there or columns get silently cut.

**Exit 1 does not mean failure.** `No result found for given query` exits 1. An
empty result and a broken query are the same exit code, so read stderr before
concluding anything.

## Step 1 — Is the CLI installed?

```
which jira
```

Missing → report this and stop:

```
brew install ankitpokhrel/jira-cli/jira-cli
```

## Step 2 — Is it configured?

```
test -f "${JIRA_CONFIG_FILE:-$HOME/.config/.jira/.config.yml}"
```

Present → read `server` and state which site is configured, so a wrong account is
caught before any task is read. `project` and `board` are often empty; that is
fine, nothing below depends on them.

Missing → the user runs this **themselves**, in their own shell:

```
! jira init
```

`jira init` is interactive. It asks for the site, the login email and an API token,
writes the config, and stores the token in the OS keyring — on macOS the Keychain,
under service `jira-cli`.

Where the token comes from:
<https://id.atlassian.com/manage-profile/security/api-tokens>

Then stop and wait. Do not continue in the same turn.

### The token never comes through here

**Never ask the user to paste the token into the conversation. Never accept it if
offered. Never write it to a file, a config, an env var, or an example.**

A token pasted to an agent is in the transcript, the logs and the session history
from that moment on, and rotating it is the only remedy. `jira init` already puts
it in the keyring, which is where it belongs.

If the user pastes one anyway: say plainly that it is now exposed and should be
revoked and reissued, then continue without it.

## Step 3 — Does the token work?

A config file proves nothing — the token may be absent, expired or revoked.

**`jira me` does not test authentication.** It echoes the configured login from the
config file without calling the API, and exits 0 with a dead token.

Make a real authenticated call and keep the exit code:

```
jira project list > /tmp/jira-auth.out 2> /tmp/jira-auth.err; echo $?
```

- exit 0 → authenticated. One line, move on.
- non-zero → read `/tmp/jira-auth.err` before diagnosing:
  - `401` / `403` / `Unauthorized` → the credential is unusable. Point at
    `! jira init` to re-enter the token.
  - a network or DNS error → say so. This is not an auth problem and re-running
    `jira init` will not fix it.
  - anything else → show the error verbatim.

Then stop. Do not guess at the cause. **Do not retry in a loop** — repeated failed
auth against Atlassian gets the account throttled.

## Step 4 — Which task?

Resolve the argument:

| Input | Meaning |
|---|---|
| `EUCAR-1234` | the key, use it |
| `.../browse/EUCAR-1234` | extract the key |
| `.../projects/EUCAR/boards/547` | a board, not a task — list it and ask |
| anything else | a description, confirm before assuming |

**The project key is the issue key's prefix** — `EUCAR-3494` → `EUCAR`. Derive it;
do not depend on `project` being set in the config.

With no argument, or with a board URL, list candidates first so the answer is a
choice and not recall:

```
jira issue list -p EUCAR -q "project = EUCAR AND status != Done" \
  --plain --no-truncate --columns key,type,status,summary
```

Then ask with `AskUserQuestion`: a recommended option, the other candidates, and
room to type a different key. Never choose on the user's behalf.

## Step 5 — Read it

The description is the point of this phase, and **`jira issue view --plain` does
not include it** — it prints only a metadata header. Use `--raw`:

```
jira issue view EUCAR-3494 --raw > /tmp/jira-issue.json; echo $?
```

```
jq -r '{
  key,
  type:    .fields.issuetype.name,
  status:  .fields.status.name,
  summary: .fields.summary,
  parent:  (.fields.parent.key // null),
  desc:    .fields.description
}' /tmp/jira-issue.json
```

`parent` being non-null means the given key is itself a subtask. Say so — the user
may have meant the parent.

Subtasks need the project flag. **`-q "parent = KEY"` alone fails**, because `-q`
runs in a project context and there is none when the config's `project` is empty:

```
jira issue list -p EUCAR -q "parent = EUCAR-3494" \
  --plain --no-truncate --columns key,type,status,summary > /tmp/jira-subs.out 2>&1; echo $?
```

Exit 1 with `No result found` means **no subtasks**, not an error. Report it as
"no subtasks" and continue.

## Output

Report, briefly:

- key, type, status, summary
- whether it is itself a subtask, and of what
- the subtasks with their statuses, or that there are none
- anything in the description ambiguous enough to need phase 4

Then carry straight into phase 2. Nothing else happens *here* — writing the spec is not
this phase's job, but it is the same turn.
