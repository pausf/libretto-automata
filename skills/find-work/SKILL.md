---
name: find-work
description: "Trigger: starting the Libretto flow; asking what to work on; resuming a half-finished change; reading a Jira task, issue URL or board URL; a request made in conversation. Phase 1 — finds the work wherever it came from and loads exactly one piece of it."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Phase 1 of the Libretto flow: **find the work**, and load exactly one piece of it.

It stops there. It does not write specs, plan, branch, code or commit.

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

## Source 3 — what the user said

**A request in conversation is a legitimate input, not a fallback.** It is how most work
arrives; every change in this repository so far arrived that way.

There is no key, so two things have to be produced rather than looked up:

- **the change's name**, from the request, verb-led and readable — `add-relative-discounts`,
  never an invented ticket id. A fake key implies a tracker that could be consulted.
- **`proposal.md` recording `Tracker: none`** and what was asked, in the words it was
  asked in. Paraphrasing a request loses the part you did not understand yet.

Confirm the reading before phase 2 starts. A request understood slightly wrong becomes a
spec that is confidently wrong, and the spec is harder to unpick than the sentence.

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

Then stop. Nothing else happens here.
