# close-flow-open-questions — payload delta

Targets: payload

Four amendments to the flow, all closing gaps the flow has already named about
itself: the drift gate's opt-in, the artifact nobody looks at, the reviewer findings
nobody counts, and phase 2's fear of asking.

## Outcomes

1. `spec-drift --block` exists: the same checks as the default mode, exiting 1 when
   anything was warned about, 0 when silent. The default mode still always exits 0.
2. `skills/record-work/SKILL.md` documents a copy-paste pre-commit hook snippet using
   `--block`, presented as the opt-in it is — the warn-never-block default does not
   move, and the prose says so.
3. `skills/build-and-check/SKILL.md` carries a rule with a stated test: **if judging
   the change means looking at it, render it and look before the review seam.** A
   palette, a layout, an image, a panel row all qualify; a flag whose output is read,
   not looked at, does not. What was rendered and what was seen goes in the
   evidence; measured contrast where the change is about colour. Not a new phase,
   not a stop.
4. `docs/FLOW.md`'s Open section closes its "where does the artifact get looked at?"
   question, pointing at the phase-6 rule as the answer.
5. `skills/review-work/SKILL.md` instructs the seam to append one entry per reviewer
   finding to `.agents/lessons.md`, using the ledger's existing header contract
   (`## <date> · <change> · <phase>`) with the phase field exactly `6→7` — the value
   the cli delta keys on to keep findings out of the user-corrections column.
6. `skills/write-spec/SKILL.md` step 4 loses its hard cap: still one
   `AskUserQuestion` call, still zero-is-legitimate, but the count is judgment —
   ask what a wrong guess would make expensive, and when in doubt, ask. The prose
   carries the user's own boundary: better to ask than to stay silent, and never a
   form-length interrogation of things the code already answers.

## Scope boundaries

In: the six outcomes above. Out, named:

- **No hook is installed by anything.** The snippet is documentation; wiring it is
  the user's act in their own repository. `settings.json` and git hooks stay
  unmanaged, per the standing decision.
- **No new phase and no new stop** for the visual rule. It rides inside phase 6.
- **The reviewer stays read-only on the work.** Appending its findings to the ledger
  is the seam's write, made by the orchestrator side of the seam — the reviewer
  subagent itself still writes nothing.
- **No screenshot tooling.** Rendering means the project's own render path
  (`libretto preview`, a test's output, an existing script). Building capture
  infrastructure is out.
- **`--block` gains no options.** No severity levels, no per-path config. One flag,
  one meaning.

## Constraints

- `spec-drift` is POSIX-leaning bash with a `--self-test`; the new flag gets a
  self-test case like the existing flags have.
- The ledger header contract is already parsed by `cmd/libretto/metrics.go
  corrections()` — three fields, two ` · ` separators. The `6→7` phase value must
  survive that parser unchanged (it is a plain string; no code change needed to
  accept it).
- Skills never reference `scripts/` or `docs/` — the snippet lives inside the skill
  and names the script by its installed path.
- `ponytail:` comments and all payload prose in English.

## Prior decisions

- Warn-never-block as the default: settled in FLOW.md and record-work; this change
  builds the opt-in *because* that default is right. (2026-08-14, this change)
- Reviewer findings enter the ledger: chosen by the user over aggregate-only, so all
  three correction sources become measurable. (2026-08-14, AskUserQuestion)
- Phase 2 question cap removed in favour of judgment with a one-call bound: the
  user's words — better to ask out of fear than not to ask, but no thousand
  questions. Supersedes the three-never-four rule of 2026-08-12. (2026-08-14,
  AskUserQuestion)
- The builder reads its own render; no user stop is added. Follows from the stop
  test: a render the user must approve is a fourth stop, and the seam already
  exists to catch what the builder missed. (2026-08-14, this change)

## Task breakdown

1. `spec-drift`: add `--block`, its usage line, and a self-test case.
2. `record-work` SKILL.md: the hook snippet section, opt-in framing.
3. `build-and-check` SKILL.md: the visual-output rule; FLOW.md Open section closed.
4. `review-work` SKILL.md: append findings to the ledger, header contract cited.
5. `write-spec` SKILL.md: step 4 recut — no hard cap, one call, judgment bounds.

## Verification criteria

- `--block` exits 1 when the default mode would have warned, 0 when silent.
  Proof: skills/record-work/spec-drift --self-test
- The payload stays internally consistent — frontmatter, references, reachability.
  Proof: scripts/check-payload
- The record-work skill documents the hook snippet and names it opt-in — a wiring
  row anchors the snippet's presence.
  Proof: scripts/check-payload
- build-and-check carries the visual-output rule — a wiring row anchors it.
  Proof: scripts/check-payload
- FLOW.md's Open section no longer lists the artifact question as open — no
  mechanical anchor reaches `docs/`, and inventing one for a one-time closure is
  tooling for a single edit, so the 6→7 seam reads the section against this
  criterion instead of a `Proof:` citation. A deliberate gap, named.
- review-work instructs the ledger append with the exact header contract — a wiring
  row anchors it.
  Proof: scripts/check-payload
- write-spec's step 4 states the judgment rule and keeps the one-call bound — a
  wiring row anchors it.
  Proof: scripts/check-payload
