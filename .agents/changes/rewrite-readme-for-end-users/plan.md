# rewrite-readme-for-end-users — plan

> **Live state.** The orchestrator is the only writer. A box is marked the moment its task
> is verified, never batched at the end. An agent joining late reads this file and believes
> it.

**Goal:** `README.md` becomes what it is → install → first run → reference, and every
"why" paragraph lands in `docs/` instead of being cut.

**Architecture:** Prose relocation, proven by one Go test that reads `README.md`,
`docs/DESIGN.md` and `docs/FLOW.md`. Same shape as `cmd/libretto/gates_test.go`, which
tests `.github/**` the same way and already provides the `repoFile(t, rel)` helper this
reuses.

**Spec:** [`spec.md`](spec.md) · `Targets: readme` — the capability does not exist yet and
task 4 creates it.

## Global Constraints

- Go 1.26.5 · no new dependency · stdlib only (`os`, `strings`, `regexp`, `testing`)
- **No behaviour changes.** No command, flag, output string or exit code moves.
- `𝄞` stays in `README.md`. `♩♪♫♬` appear nowhere.
- The install line stays `@latest` — the `v1.0.2` tombstone is what keeps that on `0.5.x`.
- No invented version in prose: `@<version>`, never `@v0.5.0`.
- Every relative link resolves from the repository root.
- All six gates pass before every commit. `make gates` runs them.

## Order

```
1 test (red) ──▶ 2 move the reasoning ──▶ 3 rewrite the README ──▶ 4 land the capability
```

Strictly sequential — each task turns part of task 1's test green, and nothing here is
independent. **Task 2 before task 3 is load-bearing:** the reasoning is safe in `docs/`
before the rewrite touches the README, so a badly aimed edit cannot lose it.

**Can start now: task 1.**

---

### Task 1 — the test, red

**Files:** Create `cmd/libretto/readme_test.go`
**Consumes:** `repoFile(t, rel)` from `cmd/libretto/gates_test.go` — reads a path relative
to the repository root and fails the test if it is missing.
**Produces:** `section(t, doc, heading)` — the text under a `## ` heading up to the next
`## `, subsections included. Task 3 relies on nothing else.
**Closes:** every criterion in the spec except outcome 5.
**Spec:** verification criteria, all six.

- [x] **Step 1 — write the file**

**The snippet below is the first version, kept as written.** It shipped with two changes
the corrections at the foot of this file describe: `flat()`, and a second anchor for the
paraphrase the reviewer caught. `cmd/libretto/readme_test.go` is the file; this is the plan.

```go
package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The README is the one document written for somebody who has never seen this project,
// so its shape is a promise and the shape is what gets tested. Wording is not — a test
// that pins prose is a test somebody deletes the first time they improve a sentence.

// section returns the text under heading, up to the next `## `. Subsections stay in:
// reasoning hiding under a `### ` inside Install is still reasoning inside Install.
func section(t *testing.T, doc, heading string) string {
	t.Helper()
	start := strings.Index(doc, heading)
	if start < 0 {
		t.Fatalf("the README has no %q section", heading)
	}
	rest := doc[start+len(heading):]
	if end := strings.Index(rest, "\n## "); end >= 0 {
		return rest[:end]
	}
	return rest
}

func TestReadmeSectionsAreInReadingOrder(t *testing.T) {
	readme := repoFile(t, "README.md")

	previous, at := "", -1
	for _, heading := range []string{
		"## What you get",
		"## Install",
		"## Your first run",
		"## Commands",
		"## Learn more",
	} {
		i := strings.Index(readme, heading)
		if i < 0 {
			t.Fatalf("the README has no %q section", heading)
		}
		if i < at {
			t.Errorf("%q comes before %q — install has to precede what to type", heading, previous)
		}
		previous, at = heading, i
	}
}

func TestInstallSectionIsStepsOnly(t *testing.T) {
	install := section(t, repoFile(t, "README.md"), "## Install")

	for _, want := range []string{
		"go install github.com/pausf/libretto-automata/cmd/libretto@latest",
		"libretto install",
		"1.26",
	} {
		if !strings.Contains(install, want) {
			t.Errorf("the install section never says %q", want)
		}
	}
	// Why the payload rides inside the module is a good answer to a question nobody
	// has yet while they are counting steps. It lives in docs/DESIGN.md now.
	for _, moved := range []string{"GOMODCACHE", "@v0.", "checksum database"} {
		if strings.Contains(install, moved) {
			t.Errorf("the install section still explains %q — that belongs in docs/DESIGN.md", moved)
		}
	}
}

func TestReadmeWalksAFirstRun(t *testing.T) {
	first := section(t, repoFile(t, "README.md"), "## Your first run")

	// Not the wording — the stops. A walk that skips where the flow pauses is a walk
	// that surprises the reader at the first pause.
	for _, want := range []string{"/libretto-flow", "spec", "plan", "push"} {
		if !strings.Contains(first, want) {
			t.Errorf("the first-run walk never mentions %q", want)
		}
	}
}

// Moved and deleted look identical in a diff of one file. Every subject here has to be
// absent from the README and present in docs/ — that is the whole test, and it is the
// reason it reads three files.
func TestMovedReasoningLandedInDocs(t *testing.T) {
	readme := repoFile(t, "README.md")
	docs := repoFile(t, "docs/DESIGN.md") + repoFile(t, "docs/FLOW.md")

	for _, moved := range []struct{ subject, phrase string }{
		{"no --force", "there is no `--force`"},
		{"prune is not uninstall", "Prune deliberately"},
		{"both scope flags is an error", "Two answers to one question"},
		{"prune and uninstall are dry by default", "being asked twice"},
		{"two queue commands and not one", "substitutes different work"},
		{"the payload is not compiled in", "module cache"},
		{"model aliases rather than ids", "two spellings of one state"},
		{"spec-drift warns, never blocks", "a deleted check finds nothing"},
	} {
		if strings.Contains(readme, moved.phrase) {
			t.Errorf("%s is still argued in the README", moved.subject)
		}
		if !strings.Contains(docs, moved.phrase) {
			t.Errorf("%s is in neither docs/DESIGN.md nor docs/FLOW.md — moved means moved, not deleted", moved.subject)
		}
	}
}

// `](#anchor)` never matches: the group demands at least one character that is neither
// `)` nor `#`.
var readmeLink = regexp.MustCompile(`\]\(([^)#]+)(?:#[^)]*)?\)`)

func TestReadmeLinksResolve(t *testing.T) {
	for _, match := range readmeLink.FindAllStringSubmatch(repoFile(t, "README.md"), -1) {
		target := match[1]
		if strings.HasPrefix(target, "http") || strings.HasPrefix(target, "mailto:") {
			continue
		}
		if _, err := os.Stat(filepath.Join("..", "..", target)); err != nil {
			t.Errorf("the README links to %s, which does not exist", target)
		}
	}
}
```

- [x] **Step 2 — run it and read the failure**

```bash
go test ./cmd/libretto/ -run 'TestReadme|TestInstallSection|TestMovedReasoning' -count=1
```

Expected: **FAIL.** `TestReadmeSectionsAreInReadingOrder` fatals on `## What you get`,
which does not exist. `TestMovedReasoningLandedInDocs` reports eight subjects still in
the README. `TestReadmeLinksResolve` should already **pass** — today's links resolve, and
a green there now is what proves the regex works rather than matching nothing.

If `TestReadmeLinksResolve` passes with zero links matched the test is vacuous. Confirm
with `-v` that it examined the real links before trusting it.

- [x] **Step 3 — do not commit it yet**

**The original plan said "gates, then commit" here and that was wrong.** Gate 3 is
`go test ./... -count=1`, so a commit carrying a red test is a commit with a red gate,
and `AGENTS.md` requires the test to ship *in the same commit as the logic it proves*.
Watching it fail is the discipline; committing the failure is not part of it.

`readme_test.go` is therefore staged with task 3, alongside the README it describes.
`gofmt -l .` and `go vet ./...` are clean now — observed — and that is the whole check
this step gets.

---

### Task 2 — move the reasoning into `docs/`

**Files:** Modify `docs/DESIGN.md`, `docs/FLOW.md`, `README.md`
**Consumes:** the eight subjects listed in task 1's `TestMovedReasoningLandedInDocs`.
**Closes:** outcome 4.
**Spec:** outcome 4, and the prior decision that no `docs/WHY.md` is created.

- [x] **Step 1 — copy each paragraph into the file that already owns its subject**

`docs/DESIGN.md` — tool behaviour. It already opens with `Why symlinks, per item`, so
these are new `##` sections beside it:

| Subject | From `README.md` |
|---|---|
| no `--force`, foreign entries reported | the "Two things in one repository" paragraph |
| `prune` is not `uninstall` | the `prune`/`uninstall` contrast |
| dry by default, asked twice | the same paragraph |
| both scope flags is an error | "Two answers to one question" |
| the payload is not compiled in | the checkout section, and the module-cache sentence |
| aliases rather than model ids | "two spellings of one state" |

`docs/FLOW.md` — the flow and the queue:

| Subject | From `README.md` |
|---|---|
| two queue commands and not one | "substitutes different work for what you handed it" |
| `spec-drift` warns and never blocks | "a deleted check finds nothing" |

**Verbatim where it still reads correctly.** A sentence that referred to a nearby README
table has to be repointed, and that is the only edit permitted here — the phrases task 1
asserts on are the anchor, so rewording one breaks the test on purpose.

- [x] **Step 2 — delete those paragraphs from `README.md`, and nothing else yet**

No restructure in this step. Delete, leave the current headings standing. The README will
read badly for one commit; task 3 is what fixes that, and separating the two keeps the
"was anything lost?" question answerable from one diff.

- [x] **Step 3 — the relocation test goes green**

```bash
go test ./cmd/libretto/ -run TestMovedReasoningLandedInDocs -count=1 -v
```

Expected: **PASS**, all eight subjects. The order, install and first-run tests are still
red — that is task 3.

- [x] **Step 4 — gates, then commit**

```bash
make gates
git add docs/DESIGN.md docs/FLOW.md README.md
git commit -m "docs: move the reasoning out of the README into DESIGN and FLOW"
```

---

### Task 3 — rewrite `README.md`

**Files:** Modify `README.md`
**Closes:** outcomes 1, 2, 3, 5.
**Spec:** outcome 1's seven-section table.

- [x] **Step 1 — restructure into the seven sections**

Order, and nothing between them:

1. name · the one-line claim · badges · `docs/panel.svg` · the automaton paragraph
2. `## What you get` — the payload and the CLI in two lines, and what a session looks like
3. `## Install` — Go 1.26+, `go install …@latest`, `libretto install`. Three steps.
4. `## Your first run` — written new, below
5. `## Commands` — today's table, one line per command, no prose between rows
6. `## Where it installs` · `## The five states` · `## Environment` — reference detail
7. `## Learn more` — `docs/FLOW.md`, `docs/DESIGN.md`, `AGENTS.md`, `THIRD-PARTY.md`

Keep, because the README is their only home: badges, the panel image, the automaton
paragraph, the licence, the `THIRD-PARTY.md` attribution and its three-plugin list.

- [x] **Step 2 — write `## Your first run`**

Numbered steps, each saying what the user types and what comes back:

```markdown
## Your first run

You installed the CLI. The flow lives in Claude Code from here — these are slash
commands, not shell commands.

1. **Say what you want.** `/libretto-flow add a --json flag to status`
   It reads the request, names the change, and writes it down before anything else.
2. **It stops at the spec.** You get a contract: what "done" means, what is out of
   scope, and how each promise will be proven. **This is the cheap place to disagree.**
3. **It stops at the plan.** An ordered checklist, with what can start now.
4. **It builds**, marking boxes as it goes, and leaves a test behind for each one.
5. **A fresh reviewer reads it** — one that saw none of the session that wrote the code,
   which is the point — and every finding gets fixed.
6. **It reports**, in the spec's own terms, including what it deliberately did not build.
7. **It asks once about pushing.** That answer is always yours.

Nothing in flight? `/libretto-status`. An idea you do not want to build today?
`/libretto-queue` captures it, and `/libretto-next` picks the oldest one up later.
```

- [x] **Step 3 — the whole file goes green**

```bash
go test ./cmd/libretto/ -count=1
```

Expected: **PASS**, every case. `TestReadmeLinksResolve` matters most here — the rewrite
adds links into `docs/`.

- [x] **Step 4 — read it as a stranger**

No test holds this. Read the file top to bottom and answer out loud: could somebody who
has never seen this project install it and type one useful thing? Say what is still
wrong, or say that nothing is.

- [x] **Step 5 — gates, then commit**

```bash
make gates
git add README.md
git commit -m "docs(readme): what it is, how to install it, what to type"
```

---

### Task 4 — land the capability

**Files:** Create `.agents/specs/readme/spec.md` · Modify `docs/SPEC.md` · Delete
`.agents/changes/rewrite-readme-for-end-users/`
**Closes:** nothing new — it makes the contract true.
**Spec:** the delta's `Targets: readme` preamble.

- [x] **Step 1 — write `.agents/specs/readme/spec.md`** from the delta:
      `Governs: README.md`, the six criteria with their `Proof:` citations intact.
- [x] **Step 2 — add the `readme` row to `docs/SPEC.md`'s capability table.** That table
      is the only place the list lives, and it has drifted before by being counted twice.
- [ ] **Step 3 — delete the change folder**, this plan included.

```bash
skills/record-work/spec-drift --anchors    # the six new citations must resolve
make gates
git add -A && git commit -m "docs(spec): the readme capability"
```

- [ ] **Step 4 — the push question.** Phase 8 asks it. It is not answered here.

---

## State

| Task | Done | Evidence |
|---|---|---|
| 1 · the test, red | ☑ | 4 of 5 cases FAIL for the stated reasons, `TestReadmeLinksResolve` PASS with links matched. `gofmt`, `vet` clean. Uncommitted by design — see task 1 step 3. |
| 2 · move the reasoning | ☑ | `TestMovedReasoningLandedInDocs` PASS, all eight subjects. Six gates exit 0 (283 citations resolve). Commit `83fa8ce`. |
| 3 · rewrite the README | ☑ | All five README cases PASS; the whole `cmd/libretto` suite PASS. Six gates exit 0. Commit `87df489`. |
| 4 · land the capability | ☑ | `.agents/specs/readme/spec.md` + the `docs/SPEC.md` row. `--anchors`: 288 citations, all resolve, the five new ones named. Six gates exit 0. Commit `6069fac`. **The change folder is not deleted yet** — phase 7's reviewer reads the contract, and phase 8 deletes it with the push question. |

## Phase 7 — the fresh reviewer

Five cited proofs re-run by it, all pass, plus `gofmt`, `go vet`, `check-payload` and
`--anchors` (288 citations). Four findings, all fixed in one pass:

| Finding | Fix | Re-run |
|---|---|---|
| a relocated argument survived in the README as a paraphrase — *module cache* → *an installed copy* — with the guard still green | the sentence is gone, replaced by the fact and a link to `DESIGN.md`; a second anchor on `"wins over"` added for the same subject | `TestMovedReasoningLandedInDocs` PASS |
| the Install section still argued ("The relink is not redundant —") | trimmed to the fact: "relinks, so a release that *adds* a skill arrives linked" | `TestInstallSectionIsStepsOnly` PASS |
| `spec.md` and `plan.md` untracked — the contract absent from the diff | both committed, so the merge carries what the code was written against | — |
| 0 of 16 boxes marked while the state table claimed four tasks done | 14 marked; the two open ones belong to phase 8 | — |

Its unresolved question — whether the remaining `prune`/`uninstall` sentences are reference
or residue — is settled in the capability spec's prior decisions, not here: what a command
*does* stays, *why it is that way* goes.

**Two corrections made while running, recorded so they are not rediscovered:**

- Task 1 does not commit. The test ships with task 3. See task 1 step 3.
- `flat()` was added to the test: every substring assertion normalises whitespace first.
  Without it `"checksum database"` never matched, because the README wraps between the two
  words — a guard that silently could not fire. Found by reading the first red run instead
  of just counting the failures.
