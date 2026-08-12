# fix-vendored-superpowers-refs

Tracker: none
Queued: 2026-08-12

## The ask, verbatim

> lo de arregla el superpowers — skills/writing-plans/SKILL.md references
> superpowers:subagent-driven-development, superpowers:executing-plans and
> superpowers:using-git-worktrees, which the payload does not ship; rewrite the
> refs to point at the payload's own skills or vendor what's missing

## Reading

The vendored `writing-plans` skill (lines 16, 61, 163–167) requires three
`superpowers:`-namespaced skills the payload never installs, so plans generated
by users without superpowers demand a skill that does not exist. Fix is a
decision between rewriting the refs to the payload's own skills
(`write-plan`/`build-and-check`) or vendoring the missing ones —
`scripts/check-payload` should also learn to validate namespaced refs so this
cannot regress.
