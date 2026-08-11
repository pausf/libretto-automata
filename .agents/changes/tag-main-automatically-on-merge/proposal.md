# tag-main-automatically-on-merge

Tracker: none
Queued: 2026-08-11

## The ask, verbatim

> podemos automatizar esto git tag -a v0.5.1 -m "..." / git push origin v0.5.1 / make release ?

> podemos automatizar esta parte con github ci ? cuando se hace un merge a main?

> como debemos hacer esto para que sea automatico al hacer merge a main para que se
> tagge la version

## Why this is queued and not started

It **reverses a promise**, and that makes it a bigger piece of work than it looks.

`.agents/specs/ci/spec.md` puts two things out of scope explicitly: *"no tags, no
publishing from a workflow"*, and *"anything that writes — a workflow with write access is
a workflow that can be made to write"*. It also names the ceiling in advance: a workflow
gated on a tag **a human already pushed** is the allowed shape, one that **creates** the
tag never is.

So by `AGENTS.md#Versioning`, a promise reversed is a **major** bump: the merge that lands
this is the one that takes the repository to `v1.0.0`. That is not a reason not to do it —
it is a reason it does not ride along in someone else's pull request.

## The one question that has to be answered first

**Where does the bump come from?** It is not derivable from the commit log, and that is
the whole difficulty rather than a detail.

`AGENTS.md` makes *"a new promise in an existing spec"* a minor and *"a promise removed or
reversed"* a major. Both are readings of `.agents/specs/`, not of `type:` prefixes. So a
workflow that reads only commit types picks the floor, and is **silently wrong exactly
when a contract moves** — the one case where getting it wrong matters.

| | Bump source | Cost |
|---|---|---|
| **recommended** | a `release:patch\|minor\|major` label on the pull request; fail the run when it is absent | the human decides at review time, CI only executes. Never silently wrong; refuses instead |
| | commit types alone | fully automatic, and wrong the day a spec promise moves |

Deliberately not defaulting to patch when the label is missing: assuming patch *is* the
silent-wrong-bump failure, wearing a different hat.

## The four traps, worked out already

Recorded so they are not rediscovered, because three of them fail in ways that do not
look like the cause:

1. **`fetch-depth: 0` on checkout.** The default is depth 1 with no tags, so
   `git describe --tags` finds nothing and the bump computes from zero.
2. **A tag pushed with `GITHUB_TOKEN` does not trigger workflows.** GitHub suppresses it to
   prevent recursion. So `gh release create` must run in the **same job**; a second
   workflow on `push: tags` would never fire, and the symptom is a tag with no Release and
   no clue why.
3. **`ref: main`, never the pull request head.** This is a job holding `contents: write`
   fired by a pull-request event. Checking out contributor code there hands over the token.
4. **`concurrency: group: release`.** Two merges landing close together read the same last
   tag, and the second tries to create one that exists.

And the gates run in the job before it tags: `.agents/specs/ci/spec.md` already says a
release is the one moment where *"the gates were green earlier"* is not good enough.

## What it will have to settle

- **Whether `make release` stays.** If CI publishes, the hand-run target is either dead or
  the documented fallback, and picking one is part of this.
- **Where the reversal is recorded.** `ci`'s scope boundaries and the *Never* list in
  `AGENTS.md` both say the opposite of this today; both have to move in the same commit.
- **Whether a human can still stop it.** A merge that tags immediately removes the last
  point where somebody notices the bump is wrong.
