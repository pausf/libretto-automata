# make-installed-payload-portable — Implementation Plan

**Goal:** the installed payload works outside this repository — skills find spec-drift
relative to themselves, spec-drift refuses loudly when `rg` is missing, doctor reports
the gate tools.

**Architecture:** no new files, no new plumbing. Wording changes in two SKILL.md files,
a guard block in one bash script, one new check in `scripts/check-payload`, two rows in
`prerequisites()`. The regression check lands first so it is observed failing on the
current tree.

**Tech Stack:** bash, Go stdlib testing.

## Global constraints (from the deltas)

- spec-drift's drift findings keep exiting 0 — warn-never-block. Exit 2 is reserved for
  "cannot run", exit 1 stays `--anchors`' failure.
- The doctor prerequisite report never affects the exit code.
- Skills never reference `scripts/` or `docs/`.
- Both edited SKILL.md files bump frontmatter `version:`.

Tasks 1→2 are ordered (the check must fail before the fix makes it pass). Tasks 3 and 4
are independent of each other and of 1→2.

---

### Task 1 — check-payload forbids absolute `~/.claude/` paths in skills

Spec: spec-payload.md · Criterion: "Skills carry no `~/.claude/` absolute path. Proof: scripts/check-payload"

**Files:** Modify: `scripts/check-payload` (after the scripts/-reference check, ~line 100)

- [ ] **Step 1: add the check**

```bash
# A skill naming its tools by an absolute install path works only in the layout it
# was written in. install --project puts the payload somewhere else, and the skill
# breaks for exactly the users who chose that.
printf '\n── absolute install paths ──\n'
hits=$(rg -n '~/\.claude/' skills/ 2>/dev/null)
if [ -n "$hits" ]; then
  while IFS= read -r line; do fail "absolute install path: $line"; done <<< "$hits"
else
  ok "no ~/.claude/ absolute paths in skills/"
fi
```

- [ ] **Step 2: run it and observe FAIL** — `scripts/check-payload; echo $?` must exit 1
  listing the six known references (write-spec:309,397; record-work:77,87,92,93).
  This red run is the test; do not commit until Task 2 turns it green.

### Task 2 — reword the six references to the skill base directory

Spec: spec-payload.md · Criterion: same as Task 1 (turns it green)

**Files:** Modify: `skills/record-work/SKILL.md` (lines 77, 87, 92–93; frontmatter
version), `skills/write-spec/SKILL.md` (lines 309, 397; frontmatter version 1.1 → 1.2)

- [ ] **Step 1: record-work** — the script is its sibling. Replace the path in the three
  spots; where the path appears in a runnable block, use the announced base directory:

```
<skill-base>/spec-drift --anchors
```

  with one sentence above the first use: "`<skill-base>` is this skill's base
  directory, announced when the skill is invoked." Line 87's prose ("it ships beside
  this file") already says the right thing — keep it, drop the absolute path from it.
  Bump `version:` in frontmatter.

- [ ] **Step 2: write-spec** — same phrasing, sibling hop: `<skill-base>/../record-work/spec-drift`.
  Both install layouts keep skills side by side, so the hop holds. Bump 1.1 → 1.2.

- [ ] **Step 3: run `scripts/check-payload; echo $?`** — expect exit 0, the new check ok.

- [ ] **Step 4: commit** tasks 1+2 together (the check and the fix are one reviewable
  unit; committing the check alone leaves a red gate on the branch).

### Task 3 — spec-drift refuses loudly without rg

Spec: spec-payload.md · Criterion: "spec-drift without `rg` exits 2 naming the missing
tool. Proof: skills/record-work/spec-drift --self-test"

**Files:** Modify: `skills/record-work/spec-drift` (guard near the top, after `set`;
new self-test case at the end of the self-test block)

- [ ] **Step 1: the guard**

```bash
command -v rg >/dev/null 2>&1 || {
  printf 'spec-drift: rg (ripgrep) is required and not on PATH — brew install ripgrep\n' >&2
  exit 2
}
```

- [ ] **Step 2: the self-test case** — follow the existing harness's style on contact
  (read how existing cases invoke and assert). The contract to assert: run the script
  with a `PATH` where `rg` cannot resolve but `bash` can, expect exit 2 and a message
  naming ripgrep. Sketch, to adjust against the real harness:

```bash
out=$(PATH=/dev/null bash "$0" 2>&1); rc=$?
want 'missing rg exits 2' '2' "$rc"
```

- [ ] **Step 3: run `skills/record-work/spec-drift --self-test; echo $?`** — expect exit 0,
  one more check than the current 17. Then `--anchors` and default mode still exit as
  before (0 on this tree).

- [ ] **Step 4: commit.**

### Task 4 — doctor reports rg and jq

Spec: spec-cli.md · Criterion: "report lists rg and jq with attribution.
Proof: cmd/libretto/main_test.go TestPrerequisitesIncludeTheGateTools"

**Files:** Modify: `cmd/libretto/main.go` (`prerequisites()`, the return slice),
`cmd/libretto/main_test.go` (new test beside `TestPrerequisitesDoNotAffectTheExitCode`)

- [ ] **Step 1: failing test**

```go
func TestPrerequisitesIncludeTheGateTools(t *testing.T) {
	names := map[string]bool{}
	for _, p := range prerequisites() {
		names[p.Name] = true
	}
	for _, want := range []string{"rg", "jq"} {
		if !names[want] {
			t.Errorf("prerequisites() missing %q", want)
		}
	}
}
```

  Run: `go test ./cmd/libretto/ -run TestPrerequisitesIncludeTheGateTools -count=1` — expect FAIL.

- [ ] **Step 2: the rows** — in the `return []Prereq{...}` slice, after "git host":

```go
{"rg", onPath("rg"), "record-work, find-work — brew install ripgrep"},
{"jq", onPath("jq"), "find-work — brew install jq"},
```

- [ ] **Step 3: run the test** — expect PASS. Then the neighbour:
  `go test ./cmd/libretto/ -run TestPrerequisites -count=1` (both tests) — PASS.

- [ ] **Step 4: commit.**

---

**Close-out (phase 8, not tasks):** all six gates, apply the two deltas onto
`.agents/specs/payload/spec.md` and `.agents/specs/cli/spec.md`, delete this folder in
the landing commit.
