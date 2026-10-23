# name-the-wiki-after-the-repository

Tracker: none
Targets: cli
Durable decisions: the one in Prior decisions below

## The ask, verbatim

> QUE MEJORARIAS? … has todo

(2026-10-23. Item 4 of the improvements proposed after draining the queue.)

## Reading

A bug, observed seven times in the queue run: `libretto wiki` titles the specs after the
directory it runs in (`filepath.Base(projectDir)`), so a run from a git worktree writes
`# Specs — libretto-automata-wt-<change>` into a committed file. Every landing worked
around it with a temporary symlink named after the repository. The title should name the
repository, wherever it is run from.

**Observed:** `go run ./cmd/libretto wiki` inside `../libretto-automata-wt-add-queue-command`
wrote `# Specs — libretto-automata-wt-add-queue-command` on line 3 of
`.agents/specs/README.md`. **Expected:** `# Specs — libretto-automata`.

## Decision log

### Session 2026-10-23
- Q: what names the repository? → A: the directory holding git's common dir — the main checkout's root, shared by every worktree (`git rev-parse --path-format=absolute --git-common-dir`, its parent) (assumed; if wrong, the remote's name, which needs a remote)
- Q: when git cannot answer? → A: the directory the wiki runs in, as today — a project outside git keeps its title (assumed)

### Session 2026-10-23, after review-spec
- Q: a submodule or a --separate-git-dir layout reports a common dir whose parent is not the checkout — which side? → A: the common dir names the repository only when its own base name is `.git`; any other layout falls back to the run directory's name, as today (assumed; if wrong, `--show-toplevel` of the main worktree, one more git call)
- Q: is the main-checkout case proved? → A: yes, its own criterion — no existing test asserts the heading (assumed)
- Q: the HTML criterion — `<title>` alone? → A: both `<title>` and the home `<h1>`, which share one variable today (assumed)

## Outcomes

1. `libretto wiki` run from a linked worktree titles both views — `README.md`'s heading and
   `wiki.html`'s `<title>` and heading — after the main checkout's directory.
2. Run from the main checkout, the title is what it is today.
3. Where git cannot say, or its common dir is not named `.git` (a submodule, a separate git
   dir), the title is the run directory's name, as today.

## Scope boundaries

**In:** the title in `cmd/libretto/wiki.go`'s two renderers, through one helper and a seam.

**Out:** renaming anything else after the repository; reading the remote; a `--title` flag.

## Constraints

- One git call per render, through the `gitRunner` seam the wiki already uses for branches.
- Stdlib only. Every existing wiki test stays unedited.

## Prior decisions

- **The wiki's title names the repository, not the run directory.** Fixed 2026-10-23: a run
  from a worktree wrote the worktree's directory name into the committed view, seven landings
  in a row; the name now comes from the parent of git's common dir, which every worktree
  shares, falling back to the run directory when git cannot say.

## Task breakdown

- the helper and seam, both renderers, the tests.

## Verification criteria

- **When** `libretto wiki` runs in a linked worktree, the `README.md` heading **shall** name
  the main checkout's directory.
  Proof: cmd/libretto/wiki_test.go TestWikiFromAWorktreeNamesTheRepository
- **When** `libretto wiki --html` runs in a linked worktree, `wiki.html`'s `<title>` and its
  home page's `<h1>` **shall** both name the main checkout's directory.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLFromAWorktreeNamesTheRepository
- **If** git cannot report a common dir, **then** the title **shall** be the run directory's
  name.
  Proof: cmd/libretto/wiki_test.go TestWikiOutsideGitKeepsTheDirectoryName
- **If** the common dir git reports is not named `.git`, **then** the title **shall** be the
  run directory's name.
  Proof: cmd/libretto/wiki_test.go TestWikiProjectNameFallsBackOffAPlainGitDir
- **When** `libretto wiki` runs in the main checkout, the `README.md` heading **shall** name
  that checkout's directory.
  Proof: cmd/libretto/wiki_test.go TestWikiFromTheMainCheckoutNamesItsDirectory

## Approach

`wikiProjectName(projectDir)` — a package variable, like `wikiGitBranches` — returns the
base name of the parent of `git rev-parse --path-format=absolute --git-common-dir` when
that dir is named `.git`, and `filepath.Base(projectDir)` on any error or any other layout. Both renderers call it once. The worktree tests
build a real repository plus `git worktree add`, so `make test-short` skips them, as the
other real-git tests are skipped; the outside-git test needs no git.

| Considered | Why it lost |
|---|---|
| the origin remote's name | a repository need not have a remote, and a fork renames it |
| `--show-toplevel` | it is the worktree's own root, which is the bug |

Risk: a submodule and a `--separate-git-dir` layout keep today's title rather than the
repository's — accepted, and proved by the `.git`-name fallback test; a bare repository has no
working tree for the wiki to run in.

## Checklist

- [x] **The wiki names the repository from any worktree.** The helper, both renderers, the
  five tests, each red first. Closes: the five criteria. Waits on: nothing.
