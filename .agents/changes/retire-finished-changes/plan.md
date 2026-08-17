# Plan — a finished change that never landed is reported, not omitted

Spec: [spec.md](spec.md) · Proposal: [proposal.md](proposal.md)

Durable decisions: one — the signal is the **difference between the two scans phase 1
already runs**, never a third command. See *Prior decisions* below.

## Summary

No new machinery. Phase 1 already runs two `rg -c` scans over every change's checklist,
one for open boxes and one for closed, and `rg -c` prints only files that matched. **A
change that appears in the closed scan and not in the open one is finished and not
landed** — the signal has been sitting in the output the whole time and nothing said what
it meant.

So this is a wording change in two payload files, two `check_wiring` rows to hold it, and
then the landing that was owed: four decisions into `payload`'s *Prior decisions*, both
folders deleted, in one commit.

## Technical context

| | |
|---|---|
| Language | markdown. No Go moves, no test file moves |
| Gates | the same six; `check-payload` is the one that carries this |
| Held up by | `check_wiring <file> <phrase> <description>` in `scripts/check-payload` — a line-scoped literal match |
| Generated | `docs/PAYLOAD.md`, and neither file's frontmatter changes, so it will not drift |
| Blast radius | `skills/find-work/SKILL.md` · `commands/libretto-status.md` · `scripts/check-payload` · `.agents/specs/payload/spec.md` · two folders deleted |

The constraint that shapes the whole approach: **`find-work` may not reference `scripts/`
or `docs/`.** It installs into projects that have neither, so whatever detects this has to
be something the skill can state in prose and a session can run from the skill alone.

## The approach

`rg -c` prints `path:count` for files with at least one match and stays silent for the
rest. That asymmetry is already the answer:

```
rg -c '^\s*- \[ \]' .agents/changes/*/tasks.md .agents/changes/*/plan.md   # open
rg -c '^\s*- \[x\]' .agents/changes/*/tasks.md .agents/changes/*/plan.md   # closed
```

| In the open scan | In the closed scan | What it is |
|---|---|---|
| yes | either | work in flight — the existing report |
| **no** | **yes** | **finished and not landed** — the new one |
| no | no | no checklist: a captured idea, or a change not cut yet. Say nothing |

So the skill gains a paragraph naming the second row and what to do about it, and the
status command gains the same line. Nothing else runs.

The report says: the change's name, that every box is closed, and that its folder is
still present — which together mean `record-work`'s landing did not finish. It does not
say the delta was or was not applied; that is the reading the spec deliberately scoped
out.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| a third `rg` for "changes with no open boxes" | the two scans already carry it. A third command is a third thing to keep consistent with the other two, and the first time they disagree the report is wrong in a way nobody can see |
| `comm`/`diff` over the two outputs, in the skill | machinery a human reader has to execute mentally to know what the skill is claiming. The table above is the same logic and can be read |
| do it in Go, as `libretto status` | that subcommand answers a different question — symlink state — and the flow's own status is a slash command that already exists. Two answers to "what is open" is one too many, which `find-work` already says about the queue |
| make it fail rather than report | there is no commit to fail. `--retired` fires on a deletion; this state is the **absence** of one, and a gate cannot refuse a commit nobody made |
| detect it by checking the capability spec for the delta | requires deciding whether a delta is "present" in a spec, which is a reading. Wrong in one direction it accuses a correct landing; wrong in the other it clears a broken one |
| leave it as prose in `record-work` only | that is what already existed. Step 3 of four was written down, meant, and skipped — this whole change is the evidence that prose alone does not hold |

## Risks

| Risk | Mitigation |
|---|---|
| a session reads only the open scan and reports nothing, exactly as before | the mandate is worded as a **row in a table of three outcomes**, so the empty open scan is a case with a name rather than an absence to notice |
| `check_wiring` passes on a phrase that appears in an example rather than in the mandate | the phrase chosen is one that appears nowhere else in the payload — checked with `rg` across `skills/`, `commands/` and `agents/` before it is wired |
| the report becomes noise in a repository that keeps landed folders on purpose | not mitigated, and named: this payload's contract says a landed change's folder is deleted. A project that disagrees will see a standing warning, and the honest fix then is that project's, not a flag here |
| the two folders are deleted without retiring the decisions | `spec-drift --retired`, which is expected to refuse the first attempt |

## Validation and rollback

Two `check_wiring` rows, plus the six gates.

**The two phrases, named here rather than chosen at build time.** The cutter's finding:
the plan described the *property* the phrase needed — unique across the payload, sitting
in the mandate and not in an example — without naming it, which leaves the gate's entire
contract to be invented by whoever happens to be typing.

```
skills/find-work/SKILL.md       finished and not landed
commands/libretto-status.md     finished and not landed
```

One phrase, both files, because both make the same claim and a row per wording is two
things to keep in step. Confirmed absent from `skills/`, `commands/` and `agents/` before
either row is wired — a row matching a phrase that also appears in an example is a row
that proves nothing.

**The one that will pass for the wrong reason is `check_wiring` itself.** It is a literal
line-scoped match, so it goes green the moment the phrase exists anywhere in the file —
including inside a table row that says the opposite, or a sentence explaining what not to
do. It proves the mandate is *present*, never that it is *correct*, and the spec already
names that ceiling. So both rows are forced red on purpose by deleting the line they
match, before either is believed.

**And the retirement gate gets its first real exercise here.** Staging the two deletions
without touching `payload`'s *Prior decisions* must fail `--anchors`. That refusal is
run and read before the migration is written — a gate whose first real encounter is a
pass has proved nothing about itself.

Rollback: one revert. Nothing migrates, no state changes shape. The deleted folders come
back with the revert like any other file.

## Complexity deliberately kept

None. This change removes an unstated inference rather than adding a mechanism, and the
scan it reads was already there.
