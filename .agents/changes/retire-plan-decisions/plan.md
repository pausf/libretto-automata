# Plan — a plan cannot be deleted without retiring its decisions

Spec: [spec.md](spec.md) · Proposal: [proposal.md](proposal.md)

Durable decisions: the two in *Prior decisions* below — the section-not-file comparison,
and the declaration living in the plan rather than in a flag.

## Summary

One check in `spec-drift`, reading the index. If the staged set deletes a
`.agents/changes/*/plan.md`, then some `.agents/specs/*/spec.md` must have a *different
Prior decisions section* between `HEAD` and the index. The plan being deleted may opt out
by declaring it has nothing to retire.

## Technical context

| | |
|---|---|
| Language | POSIX-ish bash, inside `skills/record-work/spec-drift` (now ~600 lines) |
| Tools | `git`, `awk`, `rg` — all already used there |
| Gates | the same six; this rides inside `--anchors` |
| Reads | the **index** (`git show :path`) against `HEAD` (`git show HEAD:path`) |
| Blast radius | 1 script · 2 skills · `AGENTS.md` · 1 capability spec |

The constraint that shapes it: phase 8 stages the landing commit *before* running the
gates, because `AGENTS.md` requires the gates to prove something about the tree that gets
recorded. So the check reads staged state, exactly as default `spec-drift` already does.

## The approach

Three steps, and the second is where the precision lives.

1. `git diff --cached --diff-filter=D --name-only` for paths matching
   `.agents/changes/*/plan.md`. Empty means nothing to do and the check is silent.
2. For every `.agents/specs/*/spec.md` in the staged set, extract the *Prior decisions*
   section from `git show HEAD:<path>` and from `git show :<path>`, and compare the two
   strings. Any difference is a migration.
3. If none differ, read the deleted plan out of `git show HEAD:<plan>` and look for
   `Durable decisions: none`. Present means pass; absent means fail, naming the change.

Extracting the section is the same awk shape `criteria_of` already uses — heading in,
next heading out.

### Alternatives it beat

| Considered | Why it lost |
|---|---|
| require any edit to a capability spec | the landing commit applies the delta onto that spec **by definition**, so this is green on every landing and measures nothing. It is the version that looks like a gate and is not one |
| a `--no-decisions` flag on the command | typed by whoever is trying to get the commit through, at the moment they want it to stop complaining. The declaration has to be written while the plan is, by the person who knows |
| a commit trailer (`Retires: none`) | same problem one step removed, and it puts the record in a place `git log` carries but the plan's own reader never sees |
| compare the whole spec file | passes on any edit anywhere, including a typo fix in an unrelated section |
| block in `record-work` prose instead of a gate | that is precisely what already exists and what this change was opened to replace |
| a seventh gate | "six gates" is in ten places here; the EARS half already settled this |

## Risks

| Risk | Mitigation |
|---|---|
| the gate fires on a change with genuinely nothing to retire, and becomes a nuisance | the declaration, and it is one line in a document being written anyway |
| the declaration becomes a reflex — every plan carries `none` | not mitigated by a mechanism, and it cannot be: it is a judgment. `review-work` reads the plan and the diff and can see a plan with an alternatives table and a `none` beside it. **Named, not solved** |
| a project using this payload has no `.agents/specs/` at all | the loop over staged spec files is empty, the deletion loop is empty, the check is silent |
| the section heading is worded differently in another project | matched case-insensitively on `prior decisions`, same as the criteria heading already is |

## Validation

`--self-test`, with three cases built on real temporary git repositories rather than on
string fixtures: a landing with a migrated decision, a landing without one, and a landing
that declares `none`. String fixtures would test the awk and not the git plumbing, and
the plumbing is where this can be wrong.

**The case that will pass for the wrong reason is the third one.** A bug that makes the
check silent — a glob that never matches, a `--diff-filter` typo — passes "declares none"
and "migrated a decision" both, and fails nothing. So the middle case is the one that has
to be forced red on purpose before any of the three is believed.

`CLAUDE_HOME` is irrelevant here; nothing touches a target. What matters instead is that
every fixture repo is a `t.TempDir`-equivalent under `mktemp -d`, and that no test runs
`git` in the real repository.

## Rollback

One revert. The check is additive and silent when nothing is being landed, so a tree
without it behaves exactly as today.

## Complexity deliberately kept

**Two `git show` calls per staged spec rather than parsing one diff.** Parsing
`git diff --cached -U0` hunks and deciding which fall inside a section is fiddly and gets
the boundary wrong at a heading; showing both versions and comparing sections is obvious
and slower by an amount nobody will measure. `ponytail:` at the site.
