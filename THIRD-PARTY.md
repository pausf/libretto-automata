# Third-party items

This repository ships skills written by other people, so that installing Libretto
Automata gives you a working flow without requiring you to install anything else
first. They are vendored deliberately, not by accident.

## Vendored

| Item | Origin | Version | Licence |
|---|---|---|---|
| `skills/writing-plans/` | [obra/superpowers](https://github.com/obra/superpowers) | 6.2.0 | MIT · Copyright (c) 2025 Jesse Vincent |
| `skills/test-driven-development/` | [obra/superpowers](https://github.com/obra/superpowers) | 6.2.0 | MIT · Copyright (c) 2025 Jesse Vincent |
| `skills/using-git-worktrees/` | [obra/superpowers](https://github.com/obra/superpowers) | 6.2.0 | MIT · Copyright (c) 2025 Jesse Vincent |
| `skills/ponytail/` | [DietrichGebert/ponytail](https://github.com/DietrichGebert/ponytail) | 4.9.0 · `2ed6c52` | MIT · Copyright (c) 2026 DietrichGebert |
| `skills/ponytail-debt/` | [DietrichGebert/ponytail](https://github.com/DietrichGebert/ponytail) | 4.9.0 · `2ed6c52` | MIT · Copyright (c) 2026 DietrichGebert |
| `skills/caveman/` | [JuliusBrussee/caveman](https://github.com/JuliusBrussee/caveman) | 0.1.0 · `11ddc0c` | MIT · Copyright (c) 2026 Julius Brussee |
| `skills/caveman-commit/` | [JuliusBrussee/caveman](https://github.com/JuliusBrussee/caveman) | 0.1.0 · `11ddc0c` | MIT · Copyright (c) 2026 Julius Brussee |

Full licence texts, in [`licenses/`](licenses):
[`LICENSE-superpowers`](licenses/LICENSE-superpowers),
[`LICENSE-ponytail`](licenses/LICENSE-ponytail),
[`LICENSE-caveman`](licenses/LICENSE-caveman).

They live in a directory rather than at the root so that `LICENSE` — this repository's own —
is the only licence file a reader meets there. **None of them is an alternative to it.** Each
is the upstream licence of a vendored skill, and a vendored copy has to carry its own text:
deleting one to tidy the root would be a licensing failure, not a cleanup.

Copied unmodified. If one of them needs changing for this flow, the change goes in
the Libretto skill that calls it — never into the vendored copy, so the copy stays
comparable with upstream.

## Why vendored and not required

The flow's own skills are thin on purpose: `write-plan` delegates plan structure,
`build-and-check` delegates test discipline and worktree setup. A thin skill whose
delegate is missing is not thin, it is broken.

Being installed on the author's machine is not the same as being installed on
yours. These ship with the repository so the flow works on a machine that has
nothing else.

**ponytail and caveman sat in the not-vendored table until 2026-08-10.** The
rationale — vendoring creates a second copy of something the user may already have
chosen a version of — assumed a user who has chosen versions of things. The flow's
target is a machine with nothing on it, where "if `ponytail` is installed" was a
conditional that never came true, about a tool nobody had said exists. The same
brokenness as a missing delegate, at lower volume. Only what the flow calls by name
is vendored: the two cores, the debt ledger, and the commit generator — not the two
plugins' remaining skills.

## Keeping them current

Pinned at the versions in the table. They will drift from upstream, and that is the
accepted cost of working out of the box.

To update: replace the directory from upstream, bump the version in the table above,
and re-read the Libretto skill that calls it in case what it delegates has moved.

## Not vendored

Referenced when present, never required. The flow works without them.

| Item | Origin | Role |
|---|---|---|
| `chained-pr` | gentle-ai | chains dependent branches instead of racing them at the trunk |

## Naming

Vendored items keep their upstream names. Plugin-installed skills are namespaced by
their plugin — `superpowers:test-driven-development`, `ponytail:ponytail` — so a
vendored copy in `~/.claude/skills/` does not collide with a plugin copy. Both can
exist. If both do, they are the same guidance at two versions, and the plugin's is
likely newer — and the plugin is still the path to what is deliberately not
vendored, such as ponytail's always-on hook mode.
