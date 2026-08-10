# Delta: the installed payload finds its own tools

Targets: payload

## Outcomes

- `record-work` and `write-spec` reference the `spec-drift` script relative to the
  invoked skill's base directory — the path every skill invocation announces — never
  by the absolute `~/.claude/skills/...` path. The same instructions work under
  `install --global` and `install --project`.
- `spec-drift` run on a machine without `rg` says so on stderr and exits non-zero,
  in every mode. Today it exits 0 having checked nothing — a green gate that lied.
- `scripts/check-payload` fails on any `~/.claude/` absolute path inside `skills/`,
  so the bug class cannot return unnoticed.

## Scope boundaries

In: the six `~/.claude/skills/record-work/spec-drift` references (write-spec:309,397;
record-work:77,87,92,93), a dependency guard at the top of `spec-drift`, one new
check in `check-payload`, one new self-test case.

Out, named:

- No `PATH` installation, wrapper binary, or `libretto`-mediated dispatch for
  spec-drift. The script stays a sibling of its SKILL.md; the reference moves, the
  file does not.
- No guard for `jq` in spec-drift — it does not use `jq` (verified: zero matches).
  `jq` belongs to the cli delta's doctor report, attributed to `find-work`.
- No sweep of other skills for other absolute paths beyond what the new
  check-payload check finds mechanically.

## Constraints

- `spec-drift` default mode keeps exiting 0 **for drift findings** — warn-never-block
  is a prior decision (payload spec). The non-zero exit here is for *inability to
  run*, which is a different statement than "no drift found".
- Skills never reference `scripts/` or `docs/` — the guard lives inside the script,
  the regression check inside `scripts/check-payload` (repo-only, allowed).
- Both SKILL.md files bump their frontmatter version: contract wording changes.

## Prior decisions

- The skill base directory is the anchor because every skill invocation injects
  "Base directory for this skill: <path>" — observed in this very session for both
  skills. `write-spec` reaches the sibling as `<base-dir>/../record-work/spec-drift`;
  `record-work` as `<base-dir>/spec-drift`. Both layouts (global and project) keep
  skills side by side, so the relative hop holds in both.
- Missing `rg` exits `2`, distinct from `--anchors`' failure exit `1`.

## Task breakdown

1. Reword the six references in `skills/record-work/SKILL.md` and
   `skills/write-spec/SKILL.md` to resolve from the skill base directory; bump both
   versions.
2. Add the `rg` guard to `skills/record-work/spec-drift`, plus a self-test case that
   runs the script with an empty `PATH` and expects exit 2.
3. Add the absolute-path check to `scripts/check-payload`.

## Verification criteria

- Skills carry no `~/.claude/` absolute path.
  Proof: scripts/check-payload
- spec-drift without `rg` exits 2 with a message naming the missing tool.
  Proof: skills/record-work/spec-drift --self-test
- The self-test suite still passes end to end.
  Proof: skills/record-work/spec-drift --self-test
