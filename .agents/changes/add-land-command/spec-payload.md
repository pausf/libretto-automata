Targets: payload

# Delta — record-work gains a guarded `libretto land` clause

One clause in one file. The `libretto land` verifier exists (sibling delta,
`Targets: cli`); a verifier nothing invokes verifies nothing, so the landing
step of `skills/record-work/SKILL.md` gains the invocation — guarded the same
way its existing `libretto wiki` clause already is (lines 104–108 today).

## Outcomes

- The "Landing a change consolidates it" section of `skills/record-work/SKILL.md`
  carries one new clause, beside the existing `libretto wiki` one and in its
  shape: **where `libretto` is on PATH, run `libretto land` before the landing
  commit.** On a non-zero exit, fix the named missing part and re-run —
  never commit past it, never bypass. Where the binary is absent, say the
  landing is unverified and continue: a missing convenience never blocks a
  landing.
- The payload capability spec gains the matching criterion at consolidation
  (this delta applied onto `.agents/specs/payload/spec.md`, per the landing
  contract — not a task of this change).

## Scope boundaries

**In:** the one edit to `skills/record-work/SKILL.md`, its frontmatter
`version:` bump, and this delta.

**Out — non-goals:**

- **anything about what `libretto land` does.** Flags, which parts it checks,
  exit codes, discovery, output — all owned by the sibling delta
  (`spec.md`, `Targets: cli`). This clause names the command and reads its
  exit status; it describes nothing.
- **making the check a gate.** The seventh-gate question does not arise: the
  clause is skill prose, the same standing as the wiki clause. Part 3 of the
  landing stays owned by `spec-drift --retired`, inside `--anchors`.
- **installing the binary, or telling the user to.** Absent means unverified
  and continue, full stop. (The wiki clause sets the precedent: "may be stale
  and continue".)
- **the Go binary learning about this clause.** Payload references delivery
  here, never the reverse.

## Constraints

- Skills never reference `scripts/` or `docs/`; a skill may reference the
  `libretto` binary **only** guarded by "where `libretto` is on PATH", with an
  explicit absent-binary path. The wiki clause is the precedent and the new
  clause matches its grammar.
- `scripts/check-payload` must stay green after the edit — no uninstalled
  path, no bare `Claude` addressee introduced.
- Register: bold the guarded instruction, plain prose for the absent path,
  exactly as the wiki clause does. One paragraph; the clause states the rule
  and does not restate the command's contract.
- `libretto land` is verify-only over the staged index (sibling's contract),
  so the clause runs it when the landing commit is staged — after the delta
  application, the retirement and the folder deletion are in the index, and
  after `libretto wiki` so the refreshed index is staged too.
- Skill frontmatter `version:` tracks the skill's own contract; a meaningfully
  changed instruction bumps it (`1.3` → `1.4`).

## Prior decisions

From `decisions.md`, session 2026-08-21, carried as it stands:

- Q: does the payload learn about the command? → A: yes, minimally — a second
  delta on the payload capability: record-work gains "where `libretto` is on
  PATH, run `libretto land` before the landing commit", same shape as the
  existing wiki clause, absent-binary path unchanged. A verifier nothing
  invokes verifies nothing. If wrong: the clause is one sentence to remove.
  (assumed)

And two already standing in the payload capability spec, which this clause
must not disturb:

- **Every skill is self-sufficient once installed.** The guard is what keeps
  the skill working on a machine that installed only the payload: `libretto`
  is delivery, not a dependency.
- Drift and landing checks in this flow **warn or stop the author, never
  block someone else's commit** unless they opted in. The clause instructs
  the agent running the flow; it installs no hook.

## Task breakdown

- [ ] **The one edit.** Add the `libretto land` clause to the "Landing a
      change consolidates it" section of `skills/record-work/SKILL.md`,
      beside the `libretto wiki` clause; bump frontmatter `version:` to
      `1.4`. Proof in the same commit: `scripts/check-payload` run and read
      green.

Small by design. No second task exists; the sibling delta owns the command.

## Verification criteria

- **Where** `libretto` is on PATH, the `record-work` skill **shall** instruct
  the landing step to run `libretto land` before the landing commit; **if**
  the command exits non-zero, **then** the skill **shall** instruct fixing the
  missing part it names and re-running, never committing past it; **where**
  the binary is absent, it **shall** say the landing is unverified and
  continue rather than block. **Ceiling named:** the anchor keeps the
  instruction findable in its file, never proves a session obeyed prose — the
  same limit the wiki-clause criterion beside it lives with.
  Proof: skills/record-work/SKILL.md
- **When** the clause lands, `scripts/check-payload` **shall** pass: the
  reference is to a binary on PATH, not an uninstalled repository path, and
  the guard is what makes that true. This proves the reference is legal, not
  that the clause is followed.
  Proof: scripts/check-payload
