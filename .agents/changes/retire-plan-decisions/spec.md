# Delta — a plan cannot be deleted without retiring its decisions

Applies to: `.agents/specs/payload/spec.md`

## Outcomes

A change cannot land quietly discarding the reasoning that produced it.

- The gate fires on the **landing commit**, where the plan is deleted, and nowhere else.
  A commit that touches no `plan.md` deletion sees nothing.
- The escape is a **declaration in the plan**, not a flag on the command line — written
  by the person who knew, at the time they knew it.
- The gate proves the *section moved*, never that what moved is correct. That reading
  stays with `review-work`.

## Scope boundaries

In: `skills/record-work/spec-drift`, `skills/write-plan/SKILL.md`,
`skills/record-work/SKILL.md`, `AGENTS.md`, and the `payload` capability spec.

Out: a seventh gate. Out: retrofitting changes already landed. Out: judging whether the
migrated decision matches the plan it came from.

## Constraints

- No new dependency. `git`, `awk`, `rg` — what `spec-drift` already uses.
- It has to work in a project that is not this one: `spec-drift` ships inside the skill.
- The check reads the **index**, because phase 8 stages the landing commit before it
  runs the gates.

## Prior decisions

- **The section, not the file.** Requiring any edit to a capability spec would pass on
  the delta application alone — which happens in that same commit by definition, so the
  gate would be green on every landing and would measure nothing.
- **The declaration lives in the plan, not in a commit trailer or a flag.** A flag is
  typed by whoever is trying to get a commit through; a line in the plan is written while
  the plan is being written. Same reasoning `write-spec` gives for recording an answer
  next to the decision it produced.
- **Inside `--anchors`**, per the count-in-ten-places argument the EARS half already
  settled.

## Task breakdown

The check and its escape; the self-test; the two skills that must say the line exists;
the capability delta.

## Verification criteria

- **If** a staged commit deletes a change's `plan.md` and no capability spec's *Prior
  decisions* section differs between `HEAD` and the index, **then** `spec-drift`
  **shall** fail and name the change being landed.
  Proof: skills/record-work/spec-drift --self-test
- **Where** the deleted plan carries `Durable decisions: none`, the gate **shall** pass
  with no capability spec change.
  Proof: skills/record-work/spec-drift --self-test
- **When** no `plan.md` deletion is staged, the check **shall** report nothing and
  **shall** not fail.
  Proof: skills/record-work/spec-drift --self-test
