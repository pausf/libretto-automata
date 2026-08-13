# ask-release-label-at-attacca-end — plan

> **For agentic workers:** phase 6 is `build-and-check`. This file is **live state**, not a
> handoff — mark a box the moment its task is genuinely finished, and commit the mark in the
> same commit as the work. A mark left in the working tree is a mark that never happened, and
> phase 1 reads these boxes to decide what is in flight.

**Goal:** after an unattended run has pushed, opened the request and reported, ask the bump
once and type the user's answer onto the request.

**Architecture:** the behaviour lives in `skills/record-work/SKILL.md`, on the attacca path,
below the push confirmation. `commands/libretto-attacca.md` gains no description of it —
that command describes no phase. Each criterion is proved as *wiring* by a `check_wiring`
row in `scripts/check-payload`: not that the prompt behaves, but that the words carrying the
constraint are still in the file that owns them.

**Tech stack:** markdown prompts, `bash` + `rg` in `scripts/check-payload`, `gh` at runtime.

## Global constraints

- `gh` only. Never the REST API, never a hand-built call, never a token from the environment.
- `gh pr edit --add-label` fails on a label the repository does not define — detection is
  correctness, not politeness.
- Exactly one `release:*` label. The workflow refuses two.
- `release:major` is selectable and **never first**.
- No default when unanswered. No label invented where the repository defines none.
- No `~/.claude/` absolute path and no `scripts/` reference inside a `SKILL.md` —
  `check-payload` fails both, and both would break for everyone who installed the payload.
- Six gates before any commit. `gofmt -l .` · `go vet ./...` · `go test ./... -count=1` ·
  `scripts/check-payload` · `spec-drift --self-test` · `spec-drift --anchors`.

## The red step is real here

A `check_wiring` row whose pattern is absent prints `FAIL` and exits 1, so "write the row,
watch it fail" is an observable red — not a formality. Write the row **before** the prose it
describes, every task. Reversed, the row gets written to match prose that already exists and
nothing ever proved it could fail.

---

### Task 1 — the question, and the label typed onto the request

**Files:** modify `skills/record-work/SKILL.md` (the attacca section, below *"Then return to
the base branch"*) · modify `scripts/check-payload` (the `check_wiring` block, ~line 180)

**Closes:** criteria 1, 2 and 4. **Waits on:** nothing.

**Produces:** the section heading and the decisive strings tasks 2 and 3 extend —
`the bump is asked once`, `gh pr edit --add-label`, `read the request back`.

- [ ] **1.1 Write the three rows red**

```bash
check_wiring skills/record-work/SKILL.md 'the bump is asked once' 'attacca asks the bump at the end'
check_wiring skills/record-work/SKILL.md 'read the request back'  'the applied label is read back off the request'
check_wiring skills/record-work/SKILL.md 'never the first option' 'release:major is never recommended'
```

- [ ] **1.2 Run it and watch it fail**

Run: `scripts/check-payload > /tmp/cp.out 2>&1; echo $?`
Expected: `1`, and three `FAIL … no longer contains /…/` lines in `/tmp/cp.out`.

- [ ] **1.3 Write the section in `skills/record-work/SKILL.md`**

Below the base-branch return, inside the attacca path. It states, in prose:

- the bump is asked once, after the report and after the request is confirmed open —
  there is nothing downstream of it, which is why it is not a stop
- `gh label list --json name` first: with none of `release:patch|minor|major` defined, say
  nothing at all
- `AskUserQuestion` with the run's own reading recommended first, then the others.
  **`release:major` is present and never the first option** — recommending it is the
  announcement `AGENTS.md` forbids, and announcing it is what published `v1.0.0`
- apply exactly one with `gh pr edit <n> --add-label <label>`
- **read the request back** off the forge and confirm the label is on it. A command that
  printed no error is not a change the forge accepted — the same rule the push already carries

- [ ] **1.4 Run it and watch it pass**

Run: `scripts/check-payload > /tmp/cp.out 2>&1; echo $?` → `0`, three new `ok` lines.

- [ ] **1.5 Six gates, then commit**

```bash
git add skills/record-work/SKILL.md scripts/check-payload .agents/changes/ask-release-label-at-attacca-end/plan.md
git commit -m "feat(record-work): ask the bump once at attacca's end"
```

---

### Task 2 — the two quiet paths

**Files:** modify `skills/record-work/SKILL.md` (the section task 1 wrote) · modify
`scripts/check-payload`

**Closes:** criteria 5 and 6. **Waits on:** task 1 — it extends that section.

**Consumes:** the section heading and `gh label list` sentence from task 1.

- [ ] **2.1 Write the two rows red**

```bash
check_wiring skills/record-work/SKILL.md 'defines none of the three' 'a repository without the labels is not asked'
check_wiring skills/record-work/SKILL.md 'is never withdrawn'        'the red-check line survives an unanswered question'
```

- [ ] **2.2 Run it and watch it fail** — `scripts/check-payload; echo $?` → `1`, two `FAIL`.

- [ ] **2.3 Write the two paths**

- **no labels:** a repository that *defines none of the three* is not asked and is not told
  why. Nothing is created — inventing `release:minor` in somebody else's repository decides
  that repository's release convention.
- **unanswered:** the run ends exactly as it does today — unlabeled, and the closing
  report's red-check line *is never withdrawn* by the question. The line is written before
  the question is asked; a report that promised a red check and then quietly labelled the
  request has lied about the state the user will find. Headless (`libretto loop`) has no
  prompt, so this is the normal path there and must be the quiet one.
- **no default.** Not patch, not "the safe one". A silently-wrong bump is the failure that
  published `v1.0.0` wearing a politer name.

- [ ] **2.4 Run it and watch it pass** → `0`, two new `ok`.

- [ ] **2.5 Six gates, then commit**

```bash
git add skills/record-work/SKILL.md scripts/check-payload .agents/changes/ask-release-label-at-attacca-end/plan.md
git commit -m "feat(record-work): the no-labels and unanswered paths stay quiet"
```

---

### Task 3 — the description says a person chose it

**Files:** modify `skills/record-work/SKILL.md` (the description paragraph at *"paid for in
the request's description"*) · modify `scripts/check-payload`

**Closes:** criterion 8 — the one phase 5's coverage pass found missing.
**Waits on:** task 1 — the bump has to exist before the description can carry it.

- [ ] **3.1 Write the row red**

```bash
check_wiring skills/record-work/SKILL.md 'the bump a person chose' 'the request records who chose the bump'
```

- [ ] **3.2 Run it and watch it fail** → `1`, one `FAIL`.

- [ ] **3.3 Extend the description paragraph**

It already carries two things — what the invocation answered, and every question the run
assumed past. It gains a third: **the bump a person chose**, and that a person chose it. A
reviewer has to be able to tell an agreed bump from an assumed one, and that distinction is
what this whole command is built out of.

- [ ] **3.4 Run it and watch it pass** → `0`.

- [ ] **3.5 Six gates, then commit**

```bash
git add skills/record-work/SKILL.md scripts/check-payload .agents/changes/ask-release-label-at-attacca-end/plan.md
git commit -m "feat(record-work): the request says who chose the bump"
```

---

### Task 4 — the command stops contradicting the skill, and stays thin

**Files:** modify `commands/libretto-attacca.md` (the `Never` list, ~line 105) · modify
`scripts/check-payload` (a new `check_absent` helper beside `check_wiring`)

**Closes:** criteria 3 and 7. **Waits on:** nothing — start it in parallel with task 1.

**Produces:** `check_absent <file> <pattern> <description>`, the inverse of `check_wiring`.

- [ ] **4.1 Write the helper and the three rows red**

```bash
check_absent() {  # <file> <pattern> <description>
  if rg -qN -- "$2" "$1" 2>/dev/null; then fail "$3 — $1 still contains /$2/"; else ok "$3"; fi
}
check_absent  commands/libretto-attacca.md 'or labelling the request'    'the command no longer forbids labelling outright'
check_absent  commands/libretto-attacca.md 'gh pr edit'                  'the command does not restate phase 8'
check_wiring  commands/libretto-attacca.md 'a bump the user did not choose' 'labelling is forbidden only when unchosen'
```

- [ ] **4.2 Run it and watch it fail**

Run: `scripts/check-payload > /tmp/cp.out 2>&1; echo $?`
Expected: `1`. Two `FAIL` — the `or labelling the request` absence and the conditional
wording. The `gh pr edit` row passes already, and that is the point: it is a **regression
guard**, red only the day somebody copies the phase into the command.

- [ ] **4.3 Rewrite the one bullet in the `Never` list**

From *"merge, tag, release, or put a `release:` label on the request"* to: merge, tag or
release it, **or label it with a bump the user did not choose**. Merging, tagging and
releasing stay absolute. Say in the same bullet why the exception holds — the run asks and
types, it never reads; the reading is the user's, per `AGENTS.md`.

- [ ] **4.4 Run it and watch it pass** → `0`, three `ok`.

- [ ] **4.5 Six gates, then commit**

```bash
git add commands/libretto-attacca.md scripts/check-payload .agents/changes/ask-release-label-at-attacca-end/plan.md
git commit -m "docs(attacca): forbid labelling only when the user did not choose the bump"
```

---

### Task 5 — the ci spec stops pointing at a queue

**Files:** modify `.agents/specs/ci/spec.md:222-229`

**Closes:** the ci delta's two criteria. **Waits on:** nothing — but land it **last**, in
the same commit as the delta fold, so the sentence and the folder it names die together.

- [ ] **5.1 Rewrite the decision's closing sentence**

*"The fuller fix — asking the bump natively at the end of the run — is queued as
`ask-release-label-at-attacca-end`"* becomes the record that it was built and where it
lives: `skills/record-work/`, on the attacca path. Every other word of that decision stays
— the prediction rule is independent and still governs any report ending at this
repository's request.

- [ ] **5.2 Confirm no spec still calls it queued**

Run: `rg -n 'ask-release-label-at-attacca-end' .agents/specs/`
Expected: the amended sentence, and no occurrence of the word *queued* beside it.
**This is the criterion with no `Proof:`** — `--anchors` knows nothing about change-folder
names, so this is a human read, done here where the folder is deleted.

- [ ] **5.3 Six gates, then commit** — with the delta fold, per phase 8.

---

## What can start now

**Tasks 1 and 4.** They touch different files and neither waits on anything.

Task 2 and task 3 both extend the section task 1 writes, so they follow it — and they follow
it **one at a time**, because they edit the same region and the second would be rewriting
ground the first just moved. Task 5 lands last, with the fold.

## What is deliberately not in this plan

- **A test that the prompt behaves.** A prompt is checked by running it. Every row above
  proves the wiring is still there to behave, and the spec says so rather than letting a
  green `check-payload` read as proof of conduct.
- **A guard against re-asking a labelled request.** Attacca opens the request in the same
  run; it cannot arrive labelled. It comes back the day the scope widens to the attended
  flow, and not before.
