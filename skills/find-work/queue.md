# The queue — captured, never started

Read this when the caller asks for the queue: `/libretto-status`, `/libretto-next`, or any
run in reporting mode. It assumes the branch scan in `SKILL.md` has already run.


**This is not source 1, and it is not a fourth source.** It is a scan this skill owns so
that nothing else has to walk the same directory: `/libretto-status` reports it,
`/libretto-next` picks from it, and both delegate here rather than writing their own
version. Two answers to "what is queued" is one too many, and the one that disagrees is
always the one nobody is reading.

`/libretto-queue` writes proposals carrying a `Queued:` line. **That line is where the
definition starts** — absent means the change is underway; present means captured and not
started, unless one of the two conditions below says otherwise:

```
rg -l '^Queued:' .agents/changes/*/proposal.md
```

Run it whenever the caller asks for the queue, including in reporting mode. For each one,
oldest `Queued:` date first: its name and what its proposal says it is for, in one line.
Two proposals carrying the same `Queued:` date go by change name, ascending — so
`/libretto-next`'s pick and the first line `libretto queue` prints are the same change.

A proposal holding a box under its `## Checklist` is not queued, whatever its `Queued:`
line says: the box means the work started.

`libretto queue` implements this section; where the two disagree, the section is the rule.
The binary is a reader, never the source — on a machine with the skills and no binary this
section is still the whole rule.

### The line is necessary and not sufficient — drop the ones already **dispatched by a branch**

**That scan reads the working tree, and the working tree is one branch's opinion.** The
`Queued:` line comes out in the pickup commit, and the pickup commit lives on the feature
branch — so from the base branch a change that has been built, reviewed and pushed still
carries the line and still reads as captured-and-not-started.

So filter the result against the branch scan this phase already ran:

**A queued name that matches an unmerged branch is not queued. It is in flight**, and it
belongs in the in-flight report with its commits and its request — never in the queue.

Match the change name against each branch with a leading `feat/`, `fix/`, `docs/`,
`chore/` or `refactor/` removed, and match it **whole**. An unrecognised prefix fails
safe: the branch does not match, so the change stays in the queue and the worst case is
the report this rule exists to prevent, rather than a captured idea silently vanishing.

**This is a bug fix, and the failure is worth keeping.** `/libretto-status` listed a
change under *the queue* while its work was finished and open as a pull request. The next
command typed was `/libretto-attacca <that change>` — which would have branched from the
base, found the proposal still queued there, and rebuilt work that already existed, ending
in two requests for one change and a conflict on every spec it had already amended.

**A report that induces rework is worse than one that says nothing**, because silence
costs a question and this costs the work twice. Reported as the user put it, and it is the
rule now.

**Never ask whether to pick one up, and never let one block.** Home first exists so
*started* work does not get abandoned; an idea costs nothing to abandon, because nothing
has been built on it. Making four captured ideas stand between the user and a Jira task
would make capture punitive, and a queue that is expensive to add to is a queue nobody
uses. Picking one up is `/libretto-next`'s job, and it asks.
