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

- [x] **Step 1: add the check**

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

- [x] **Step 2: run it and observe FAIL** — observed: exit 1, exactly the six known
  references (write-spec:309,397; record-work:77,87,92,93).

### Task 2 — reword the six references to the skill base directory

Spec: spec-payload.md · Criterion: same as Task 1 (turns it green)

**Files:** Modify: `skills/record-work/SKILL.md` (lines 77, 87, 92–93; frontmatter
version), `skills/write-spec/SKILL.md` (lines 309, 397; frontmatter version 1.1 → 1.2)

- [x] **Step 1: record-work** — the script is its sibling. Replace the path in the three
  spots; where the path appears in a runnable block, use the announced base directory:

```
<skill-base>/spec-drift --anchors
```

  with one sentence above the first use: "`<skill-base>` is this skill's base
  directory, announced when the skill is invoked." Line 87's prose ("it ships beside
  this file") already says the right thing — keep it, drop the absolute path from it.
  Bump `version:` in frontmatter.

- [x] **Step 2: write-spec** — same phrasing, sibling hop: `<skill-base>/../record-work/spec-drift`.
  Both install layouts keep skills side by side, so the hop holds. Bump 1.1 → 1.2.

- [x] **Step 3: run `scripts/check-payload; echo $?`** — observed: exit 0, new check ok.

- [x] **Step 4: commit** tasks 1+2 together — d800acf.

### Task 3 — spec-drift refuses loudly without rg

Spec: spec-payload.md · Criterion: "spec-drift without `rg` exits 2 naming the missing
tool. Proof: skills/record-work/spec-drift --self-test"

**Files:** Modify: `skills/record-work/spec-drift` (guard near the top, after `set`;
new self-test case at the end of the self-test block)

- [x] **Step 1: the guard**

```bash
command -v rg >/dev/null 2>&1 || {
  printf 'spec-drift: rg (ripgrep) is required and not on PATH — brew install ripgrep\n' >&2
  exit 2
}
```

- [x] **Step 2: the self-test case** — done in the harness's `want` style: bash invoked
  by absolute path with `PATH=/dev/null`, so only rg is missing; asserts exit 2 and
  a message naming ripgrep. Two new checks (code and message), 20 total.

- [x] **Step 3: run `skills/record-work/spec-drift --self-test; echo $?`** — observed:
  exit 0, 20 checks ok. Default mode still 0. `--anchors` exited 1 mid-change — the
  cli delta's Proof named a test task 4 had not written yet; green after task 4.

- [x] **Step 4: commit** — b6b0e90.

### Task 4 — doctor reports rg and jq

Spec: spec-cli.md · Criterion: report lists rg and jq with attribution.
Proof: cmd/libretto/main_test.go TestPrerequisitesIncludeTheGateTools

**Files:** Modify: `cmd/libretto/main.go` (`prerequisites()`, the return slice),
`cmd/libretto/main_test.go` (new test beside `TestPrerequisitesDoNotAffectTheExitCode`)

- [x] **Step 1: failing test**

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

- [x] **Step 2: the rows** — in the `return []Prereq{...}` slice, after "git host":

```go
{"rg", onPath("rg"), "record-work, find-work — brew install ripgrep"},
{"jq", onPath("jq"), "find-work — brew install jq"},
```

- [x] **Step 3: run the test** — observed: red first (`missing "rg"`, `missing "jq"`),
  then both prerequisite tests PASS after the rows.

- [x] **Step 4: commit** — 1aff3e2.

---

**Close-out (phase 8, not tasks):** all six gates, apply the two deltas onto
`.agents/specs/payload/spec.md` and `.agents/specs/cli/spec.md`, delete this folder in
the landing commit.
