# Plan — add-opencode-command-target

Execution is phase 6, `build-and-check`. This file is live state: a box is marked the
moment its task is verified, and it ships in the commit that closed it.

Six tasks. Two carry the contract; four are the prose and the proof that make it true.

## Task 1 — OpenCode accepts the commands kind

- [x] `Opencode.Kinds()` returns `{Skills, Commands}`, and the type comment stops
      claiming it accepts only skills.
- From: spec-targets, task 1
- Closes: *opencode serves skills at `<root>/skills` and commands at `<root>/commands`,
  and rejects agents*
- Waits on: nothing. **This is the only production change, and every other task
  depends on it.**
- Evidence: `internal/target/opencode_test.go` TestOpencodeAcceptsSkillsAndCommands

## Task 2 — the unit criterion, rewritten rather than extended

- [x] `TestOpencodeAcceptsOnlySkills` becomes `TestOpencodeAcceptsSkillsAndCommands`:
      `Dir(Commands)` is `<root>/commands`, `Accepts(Commands)` is true, `Kinds()` is
      exactly `{Skills, Commands}`, and `Agents` plus an invented kind are still
      rejected.
- From: spec-targets, tasks 1–2
- Closes: the same criterion as task 1
- Waits on: task 1
- **This renames a test the current spec cites.** The old name has to leave the spec in
  the same change, or `spec-drift --anchors` breaks on a citation to a test that no
  longer exists.

## Task 3 — the scope matrix

- [x] `TestResolveToolScopeMatrix` expects skills **and commands** for opencode in both
      scopes, skills alone for codex, three kinds for claude.
- From: spec-targets, task 2
- Closes: *every tool resolves in both scopes onto its own root — skills only for
  codex, skills and commands for opencode*
- Waits on: task 1

## Task 4 — the install-isolation test, whose contract moved

- [x] `TestInstallOpencodeLeavesOthersAlone` asserts a command is linked **as a
      symlink** into `<root>/commands`, the skill still is, no `agents` directory
      appears, and no other destination was written to.
- From: spec-targets, task 3
- Closes: *an opencode install links a command as a symlink into `<root>/commands`* …
- Waits on: task 1
- **The existing test asserts the opposite** — that no `commands` directory exists. That
  assertion is the old contract, not a test to protect: it is edited because the
  promise changed, and the change is recorded in spec-targets. Nothing here weakens a
  test to get a green gate; the `agents` half of the same assertion stays exactly as
  it is.
- The fixture needs a command item. `newFixture` builds skills; check whether it can
  make a `commands/*.md` and add the smallest thing that works if it cannot.

## Task 5 — the two spec prose corrections

- [x] `.agents/specs/targets/spec.md`: the table's opencode row, the skills-only
      sentence and the two criteria the deltas replace.
- [x] `.agents/specs/cli/spec.md`: the `--opencode` help row and the skills-only
      paragraph, which now names Codex alone.
- From: spec-targets task 2, spec-cli task 2
- Closes: nothing on its own — this is the delta landing, phase 8's job, listed so it
  is not forgotten
- Waits on: tasks 1–4 green

## Task 6 — help and README

- [x] `cmd/libretto/main.go`: `--opencode` reads *skills and commands*.
- [x] `README.md`: the `--opencode` destination row loses "— skills only". Owned by the
      **readme** capability; no `readme` criterion moves.
- From: spec-cli tasks 1 and 3
- Closes: *`help` still names every destination flag and every destination environment
  variable* — which this cannot fail, so the row's accuracy is held by reading
- Waits on: task 1

## What can start now

**Task 1, alone.** Everything else waits on it, and it is two lines.

Tasks 2, 3, 4 and 6 are independent of each other once task 1 is in — they touch four
different files. Task 5 is last by definition: it is the contract landing on the
capability specs.

## The commit shape

**One commit, not two.** `spec-drift --anchors` is a gate that must pass before any
commit, and the deltas cite `TestOpencodeAcceptsSkillsAndCommands` — a test that does
not exist until task 2. A spec-only commit would fail the gate it is supposed to pass,
and `AGENTS.md` already requires the spec to ship in the same commit as the code that
taught it and the test in the same commit as the logic it proves. This change is one
finished unit of work, so it is one commit.
