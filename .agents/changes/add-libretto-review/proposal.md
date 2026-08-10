# add-libretto-review

Tracker: none

## Request (verbatim)

> nueva feature con release 0.5.0 , Hacer un libretto-review donde le pases una URL de
> GitHub/ gitlab y te haga un worktree para revisar el proyecto .
>
> Si el proyecto es muy grande para hacer un worktree que haga un cambio de rama a la
> que tiene que revisar y despues de revisar volver a la rama que estabas

## Reading

A new `libretto-review` capability: given a GitHub/GitLab URL, prepare the project for
review in a worktree; when the project is too large for a worktree, fall back to
switching to the branch under review and switching back to the original branch after
the review. Targeted at release 0.5.0.

Open questions for phase 2:

- Is this payload (a skill/command shipped to `~/.claude`) or CLI (a Go subcommand)?
- "URL de GitHub/GitLab": a repo URL to clone, or a PR/MR URL of an already-cloned repo?
- What defines "too large for a worktree"?
