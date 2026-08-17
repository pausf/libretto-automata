# add-vertical-slicing-to-plan — delta

Targets: payload

Phase 5 requires a plan to be **ordered** and **derived**. It has never required a box to
be worth closing on its own. This delta adds that requirement to `write-plan`'s contract
and gates the mandate the same way every other prose mandate in the payload is gated.

## Outcomes

`skills/write-plan/SKILL.md` states, as contract rather than advice, that a box is cut so
that closing it leaves the tree green and mergeable on its own:

- one box is one end-to-end change — the code and its proof in the same box. Not a layer
  of one. **The capability spec is not part of a box**: a delta still lands on it once, in
  the final commit, exactly as `payload` already requires. What a box owes is the proof
  that closes it, never a slice of the delta
- two boxes where the first only makes sense once the second lands are **one badly cut
  box**, and the fix is to merge them rather than to order them
- a box that cannot be cut that way is a signal that phase 2's task breakdown was cut
  along components, and `write-plan`'s existing "derived, not invented" rule already says
  what to do with a task the spec got wrong: go back and find out which end is wrong
- the mandate names what it buys `libretto loop`: the loop runs one fresh session per open
  box, so a box that does not stand alone leaves the tree broken between sessions

`scripts/check-payload` fails when that mandate is absent from `skills/write-plan/SKILL.md`.

## Scope boundaries

**In:** the mandate in `skills/write-plan/SKILL.md`, one `check_wiring` line in
`scripts/check-payload`, and this delta landing on `.agents/specs/payload/spec.md`.

**Out**, named so it cannot arrive quietly:

- **`skills/write-spec/` is not changed.** The horizontal cut originates in phase 2's task
  breakdown, and amending both skills would put one rule in two files — the failure mode
  `CLAUDE.md` opens by naming. `write-plan` already owns the authority to reject a task the
  spec cut wrong, so the rule lands where it is enforced rather than where it originates
- **no check on the content of a `plan.md`.** Whether a box is genuinely mergeable alone is
  judgment, and the payload spec already records that a check exercising judgment drifts.
  The gate checks that the skill carries the mandate, nothing more
- **`skills/build-and-check/` is not changed.** Phase 6 executes a box; it does not cut one
- **the vendored `skills/writing-plans/` is not edited.** Vendored copies stay
  byte-comparable with upstream; a divergence lives in `write-plan`'s override table
- **`libretto loop` is not changed.** The mandate makes its per-box sessions safer; the
  loop's own behaviour is untouched
- **RPI is not adopted as a named methodology.** Two of its three phases are already
  covered better here, and importing the vocabulary would add a second name for phase 3
  and for `evidence`

## Constraints

- **A skill may only invoke what gets installed.** The mandate cites no path outside
  `skills/write-plan/`; `scripts/check-payload` is repo-only tooling and the skill never
  references it
- **Prose addresses the agent, never Claude by name**, and the file's existing
  `native prompt` marker must survive the edit
- the mandate must be searchable as a fixed phrase on **one line**. `check_wiring` runs
  `rg -qN`, which is line-scoped; the newline-squashing the payload spec describes belongs to
  a different check in the same script and does not apply here
- `gofmt`, `go vet` and `go test ./...` are untouched by this change and must stay green

## Prior decisions

- **The rule lands in `write-plan` only** — *assumed, nobody was asked.* `/libretto-attacca`
  answered the contract stop, so this was decided rather than raised. What changes if it is
  wrong: a spec whose task breakdown is cut along components forces phase 5 to reject tasks
  it could instead have received correctly cut, which shows up as phase 5 sending work back
  to phase 2 repeatedly. The fix then is a sentence in `write-spec`'s task-breakdown pillar,
  and the mandate here stays as it is.
- **The gate is a marker-phrase `check_wiring`, not a new kind of check** — *assumed.* It is
  the mechanism the payload already uses 30-odd times for exactly this. **Ceiling, and it is
  the one the payload spec already names about its other marker checks:** the check asks
  whether the phrase appears anywhere in the file, so it proves the mandate is present and
  proves nothing about whether a plan obeyed it. What replaces it, the day a plan is cut
  horizontally with the mandate sitting in the file above it, is a phase-6 or `review-work`
  finding — never a longer regex.
- The attribution in the original ask is removed everywhere, at the user's explicit
  request. Recorded in `proposal.md` as an asked-for edit rather than an accident.
- Vertical slicing is added as a property of a box, not as a methodology with a name to
  learn. The payload's house style is to state the rule and what it costs, which is why
  there is no new vocabulary here.

## Task breakdown

1. Add the mandate to `skills/write-plan/SKILL.md`, and bump its frontmatter `version:`
   — the skill's own contract moved, which is what that field tracks.
2. Add the `check_wiring` line to `scripts/check-payload` and watch it fail against a
   `write-plan` that does not yet carry the mandate, then pass against one that does.
3. Land this delta onto `.agents/specs/payload/spec.md` and delete the change folder, in
   the same commit as the code.

**Tasks 1 and 2 are one box, not two.** Task 1 alone leaves the payload with a mandate
nothing gates; task 2 alone fails the gate it just added. Cut apart, neither leaves the tree
green — which is the rule this change is about, applied to itself. Task 3 is phase 8 and is
the same commit.

## Verification criteria

- **`skills/write-plan/SKILL.md` carries the vertical-slicing mandate.** The gate searches
  the file for the literal phrase `one badly cut box` — the wording that states two dependent
  boxes are one box wrongly cut. Absent, the payload gate fails.

  **`check_wiring` is `rg -qN`, so the match is line-scoped and not newline-squashed.** The
  phrase has to sit on one line in the skill, and the criterion says so because the first
  draft of it claimed squashing that this gate does not do — a criterion describing a
  stronger check than the one it cites is how a promise goes green for the wrong reason.
  Proof: scripts/check-payload
- **the mandate introduces no bare `Claude` addressee.** The existing allowlist check reads
  the whole file, so the new prose is covered by it.
  Proof: scripts/check-payload
- **a host-neutral marker is still present in `skills/write-plan/SKILL.md`** after the edit.
  Separate from the addressee criterion above because they are separate checks with separate
  failure modes, and joined with an `and` one of them can be half-met and read as done.
  Proof: scripts/check-payload
- frontmatter still parses and `name:` still matches the directory after the version bump.
  Proof: scripts/check-payload
- every `Proof:` citation in this delta resolves once it lands on `payload`.
  Proof: skills/record-work/spec-drift --anchors

**What no criterion here tests, stated rather than left to be discovered:** that a plan
written under the mandate is actually sliced vertically. That is judgment, it is scoped out
above, and the ceiling is recorded under prior decisions.
