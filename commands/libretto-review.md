---
description: Review a GitHub/GitLab PR/MR in a worktree, or a branch switch that restores itself
---

Takes one input: the PR/MR URL.

```
Skill(skill="review-project")
```

The skill owns everything — forge derivation, the clone check, worktree versus
branch-switch, the reviewer, the restore. This command routes and never implements;
a second statement of any of those rules here would be the copy nobody edits.

No URL given → ask for one and stop. Nothing else is guessed.
