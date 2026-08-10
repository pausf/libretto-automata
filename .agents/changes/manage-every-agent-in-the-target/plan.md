# manage-every-agent-in-the-target — plan

**Goal:** the model selector operates on the agents of the **selected target** — all 22
of them on this machine, not the 7 this repository ships — and says which writes reach
every project.

**Architecture:** one parameter change carries most of it. `agentmodel.Agents` and
`Apply` take the agents *directory* instead of a repo root they join `agents/` onto, so
the package gains the reach without learning what a target is. `cmd` resolves the
directory from the active target and asks `link.Owned` which rows are shared;
`internal/ui` gains one bool on a row and reloads on `tab`.

Specs: `spec.md` (agent-models) · `spec-cli.md` (cli) · `spec-panel.md` (panel), all in
this folder. One writer: the orchestrator marks these boxes; sub-agents report.

## Global constraints

- six gates before every commit: `gofmt -l .` silent, `go vet ./...`,
  `go test ./... -count=1`, `scripts/check-payload`, `spec-drift --self-test`,
  `spec-drift --anchors`
- **`--anchors` is red for this change until task 4.** The deltas cite 14 tests the
  plan creates. The per-task check is **zero BROKEN outside this change**:

  ```
  skills/record-work/spec-drift --anchors > /tmp/anchors.out 2>&1
  rg -N '^  BROKEN' /tmp/anchors.out | rg -v 'manage-every-agent-in-the-target' | wc -l
  ```

  It must be fully green before phase 8.
- **Task 2 retires three cli criteria, so three citations in `.agents/specs/cli/spec.md`
  go broken and stay broken until phase 8.** They are named in `spec-cli.md` under
  *Criteria this delta retires*. From task 2 onward the check is **zero BROKEN outside
  this change and outside those three names** — a rule stated because otherwise the
  next reader cannot tell a retired criterion from rot, and would either panic or stop
  reading the gate.
- **`rg` needs `--hidden` to see `.agents`.** The obvious spelling of any check over the
  specs reports a clean repo while the specs say otherwise. It did exactly that last
  change.
- **This is a signature change, so existing tests break by design.** `apply_test.go`'s
  `repoWith` builds a repo shape; the cmd tests build a fixture repo. They move to bare
  directories rather than being deleted — a test deleted to make a signature change
  compile is the failing test weakened to get a green gate.
- `CLAUDE_HOME` points at a temp dir in anything that touches a target; never a real
  `~/.claude`. **This change makes that rule load-bearing rather than tidy** — the code
  under test now writes into target directories.
- **No test writes into the payload's own `agents/`.** Fixtures are `t.TempDir()`.
- Go 1.26.5, standard library only. `internal/ui` still imports neither
  `internal/target` nor `internal/agentmodel`.
- no push from any step; the release tag is phase 8's call.

---

### Task 1: `agentmodel` takes a directory

Depends on: nothing. **Start here** — tasks 2 and 4 both wait on the signature.
From spec: `spec.md` breakdown 1. Closes: all six agent-models criteria.

**Files:** modify `internal/agentmodel/apply.go`, `internal/agentmodel/apply_test.go`

- `Agents(dir string)` and `Apply(dir string, names []string, model string)`; the `dir`
  constant and the `filepath.Join(repoRoot, "agents")` go with them
- **a missing directory reports no agents, not an error** — a target with no `agents/`
  yet is a state, and the CLI must not have to special-case `os.IsNotExist`
- the existing `repoWith` helper becomes a helper that builds a bare agents directory;
  every current test keeps its assertion and changes only how its fixture is laid out
- **a symlinked agent file is written through to its destination.** That is ordinary
  file behaviour, and it is the mechanism behind the whole `shared` marker — so it gets
  a test rather than a comment
- [x] **1.1** `TestAgentsOnAMissingDirectoryIsEmptyNotAnError`,
  `TestApplyThroughASymlinkWritesTheDestination` written and failing first; the existing
  tests re-pointed at the new fixture shape
- [x] **1.2** implement until they pass
- [x] **1.3** six gates
- [x] **1.4** committed: 12ae329

### Task 2: `models` acts on the target

Depends on: Task 1.
From spec: `spec-cli.md` breakdown 2. Closes: all ten cli criteria.

**Files:** modify `cmd/libretto/models.go`, `cmd/libretto/models_test.go`

- the subject is `tg.Dir(target.Agents)`, not `<root>/agents`; a target that does not
  accept agents lists nothing rather than failing
- `link.Owned(root, path)` decides the `shared` marker. **One definition of ownership** —
  do not write a second predicate here
- **the post-write line names which kind was written.** The shipped message told
  everyone their write was machine-wide; that is now true of the symlinked rows only,
  and saying it about a local file is the same class of error as the silence it
  replaced
- `agentsReaching` and the `· not linked here` marker go: every listed agent is in the
  target by construction now, so the note has nothing left to say
- [x] **2.1** the ten tests named in `spec-cli.md`, failing first
- [x] **2.2** implement until they pass
- [x] **2.3** six gates
- [x] **2.4** committed: 633b206

### Task 3: the row learns it is shared, and `tab` reloads

Depends on: nothing — `internal/ui` never touches the filesystem, so this is built
against callbacks and fake rows. **Independent of tasks 1 and 2.**
From spec: `spec-panel.md` breakdown 3. Closes: five of the seven panel criteria.

**Files:** modify `internal/ui/models.go`, `internal/ui/model.go`,
`internal/ui/models_test.go`

- `AgentRow.Shared bool`, rendered as the word `shared` — **not a colour.** The strip
  already shipped that mistake once and had correct behaviour reported as a bug
- `tab` inside the selector switches destination *and* reloads the rows. A screen
  showing one target's agents under another's name is exactly the lie the strip exists
  to prevent
- a failed reload keeps the rows it had and says so, the way `nextScope` already
  handles a failed refresh
- an empty set renders a plain statement, not an empty box
- [x] **3.1** `TestSharedAgentsAreMarked`, `TestSharedMarkerIsLegibleWithoutColour`,
  `TestTabReloadsTheSelectorForTheNewDestination`,
  `TestAFailedReloadKeepsTheRowsAndSaysSo`, `TestAnEmptyAgentSetSaysSo` — failing first
- [x] **3.2** implement until they pass
- [x] **3.3** six gates
- [x] **3.4** committed: 3f4041e

### Task 4: wire the panel to the active target

Depends on: Tasks 1, 2 and 3.
From spec: `spec-panel.md` breakdown 4. Closes: the tally criterion, and turns
`--anchors` green.

**Files:** modify `cmd/libretto/main.go`, `cmd/libretto/models_test.go`

- `agentRows` builds from the active target's directory and sets `Shared` from
  `link.Owned`; the selector's callbacks close over the destination the panel is on,
  **passed in, never captured at startup** — the same rule the runner already follows,
  and for the same reason
- the menu row's tally counts that target's agents
- [x] **4.1** `TestMenuTallyCountsTheActiveTargetsAgents` — failing first
- [x] **4.2** implement until it passes
- [x] **4.3** six gates, `--anchors` fully green, plus `make preview` looked at
- [x] **4.4** committed: 935fbe7

### Task 5: namespace the lens agents

Depends on: nothing in Go. Payload only. Added mid-change, at the user's request, so
the collision is fixed before this work is reviewed rather than after.
From spec: `spec-review-project.md` breakdown 5. Closes: its check-payload criteria.

**Files:** rename `agents/review-{security,design,reliability,tests,intent}.md` to
`agents/review-lens-*.md`; modify `skills/review-project/SKILL.md`

- the agents are renamed; **the skills they apply are not** — they did not collide,
  and `install --global` reports exactly one conflict on the reporting machine
- **the reference hunt is the work, not the rename.** A dangling agent name fails as a
  lens that silently never ran:

  ```
  rg -n --hidden --glob '!.git/*' 'agents/review-(security|design|reliability|tests|intent)\.md'
  ```

  `--hidden` is not optional — `.agents` starts with a dot, and the obvious spelling
  reported a clean repo while the specs still named a deleted agent, last time
- [x] **5.1** rename the five files and their `name:` frontmatter
- [x] **5.2** update the launch table and the `Governs:` line
- [x] **5.3** the reference hunt returns nothing
- [x] **5.4** six gates
- [ ] **5.5** committed

---

## What can start now

**Tasks 1 and 3.** Task 3 is the only one that touches no Go outside `internal/ui`, so
it is the only one that can run truly beside another without two writers in one package.

Then: 2 (after 1) → 4 (after 1, 2, 3).

## Owed before this is done, not by a test

- **`libretto models --global` run against the real `~/.claude`**, listing 22 agents
  rather than 7, and one write confirmed on a file libretto did not create. Until that
  run the fix is a claim about the exact thing the last one got wrong.
- `make preview`, and the selector driven by hand with `tab`.

## Fixed here after all

- **The `review-reliability` collision is closed by task 5.** It was written up as a
  separate change; the user asked for it fixed before review rather than after, which
  is the better call — a reviewer reading this diff would otherwise have found a
  payload that cannot install on the machine it was built on.
