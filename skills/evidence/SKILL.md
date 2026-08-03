---
name: evidence
description: "Trigger: running tests, builds, linters or any command whose result decides what happens next; reporting work as done; a test that fails; deciding whether to keep going. Applies at every phase of the flow. Nothing is true here until it has been observed."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## The one rule

**Nothing is true until it has been observed.**

Everything below is that sentence applied somewhere. These are not a phase. They
hold during every phase of the Libretto flow, and an agent that treats them as a
checklist to complete once has already broken them.

## A failing test is information, not an obstacle

When a test fails, exactly two moves are legal: **fix the cause**, or **stop and
say why**.

Never edit the assertion, loosen the tolerance, add a skip, delete the case,
comment it out, or wrap it in a conditional that happens to be false today. A test
you changed in order to pass it proves nothing — it now measures whatever you just
decided, which is the thing under suspicion.

The pressure to do this arrives exactly when the work is nearly finished, which is
why it is written down instead of assumed.

If the failure is not understood after a genuine attempt, stop. A stopped item with
an honest reason is worth more than a green suite that lies.

## Read the result, do not infer it

Any command whose outcome decides the next step runs **in the foreground**, with
its output kept, and the output is read before anything is concluded.

Never end a turn with a build, a test run or a gate still running in the
background. A finished process is not a passed process.

```
go test ./... > /tmp/gate.out 2>&1; echo "exit=$?"
```

Then read `/tmp/gate.out`. Do not summarise from memory of what you expected.

**A pipe destroys the exit code.** `cmd | head` reports the status of `head`, so a
failed command reads as success. This is not hypothetical — it happened while
building `find-work`: a `jira` call printed an error and the shell reported
`exit=0`, because the output went through `head`. Redirect to a file instead. If a
pipe is genuinely needed, `set -o pipefail` or check `${PIPESTATUS[0]}`.

**And a non-zero exit is not automatically a failure.** The same CLI exits 1 for
"no results found", which is a perfectly good answer to "does this task have
subtasks?". Read stderr and decide what the code actually means before reporting
anything as broken.

## Say only what you saw

Before the final message of a turn, every claim in it must trace to something
observed in this session:

- **"tests pass"** → the run finished in the foreground and its output was read
- **"it builds"** → same
- **"committed"** → the commit exists; `git log` was consulted, not assumed
- **"pushed"** → the remote tip matches the local tip
- **"done"** → the criteria in the spec were checked one at a time

Anything that cannot be traced that way gets stated as what it really is: not done,
not verified, or not attempted. A hedge is honest. A confident claim about
unobserved work is not.

If a step was skipped, the report says which one and why. Reporting partial work
accurately is a success; reporting complete work inaccurately is the only real
failure mode here.

## Verification runs on what was recorded

A gate proves something about the tree it ran against. If the tree was dirty, the
gate proved nothing about what got committed.

So the state under test and the state recorded must be the same. Work-in-progress
commits are cheap and can be squashed later; verifying uncommitted work and then
committing it is a claim about code nobody tested.

After a rebase or a conflict resolution, the tree is new. Everything green before
the rebase is unproven again.

## Criteria name their proof

In the spec, a verification criterion that does not name the test that proves it is
not a criterion — it is a hope.

Each criterion cites the file and the case. Written before the code, which means
citing tests that do not exist yet. That is the point: the citation is what makes
the test get written, and it is what makes the spec falsifiable instead of
decorative.

At the end, the criterion and its test are checked against each other. A criterion
whose test was never written is an unmet criterion, whatever the suite says.

## Know when to stop

Repetition is not progress. Two concrete limits:

**Two failed gates on the same item → stop that item.** Not a third attempt. Report
what failed, what was tried, and what is still unknown, then move to something
else or hand back.

**Two items stopped in a row → stop everything.** That pattern is not bad luck, it
is a wrong assumption upstream, and the next item will hit it too.

When stopping for that second reason, **write the handoff note** — the project's own,
whatever it is called, and `STATE.md` beside the specs if it has none. Four things:
what stands, what is blocked, what was learned, and the one thing the next session
should read first.

A session that ends without it leaves the next one to rediscover everything, and the
thing most worth passing on is exactly the thing that made this session stop.

Loops around a flaky external dependency are the same failure with better
manners — two attempts, then treat it as down and say so.

## Where words go

The chat gets the verdict. The artifacts get the reasoning.

One line per finished item is enough: what it was, what happened, where the
evidence is. The spec holds the decisions, the plan holds the state, `STATE.md`
holds the handoff.

Raw tool output is never the final message. Prose repeated into the chat that
already exists in a file is duplication that will drift from its source.
