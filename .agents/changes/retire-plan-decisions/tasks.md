# Tasks — a plan cannot be deleted without retiring its decisions

Spec: [spec.md](spec.md) · Plan: [plan.md](plan.md)

Execution: `build-and-check`. One writer — the orchestrator marks the boxes.

- [x] **The check and its escape** — `section_of` extraction, the staged-deletion scan,
      the `HEAD`-versus-index comparison, and the `Durable decisions: none` read out of
      the deleted plan. Wired into `--anchors`, silent when nothing is being landed.
      Spec: *the gate fires on the landing commit and nowhere else*
      Closes on: three self-test cases green, and the middle one forced red on purpose
      by breaking the comparison.
      Waits on: nothing.

- [x] **The two skills that have to say it exists** — `write-plan` gains the
      declaration line in its structure and says when to use it; `record-work` says the
      landing commit carries the retirement.
      Spec: *the escape is a declaration in the plan*
      Closes on: `check-payload` green, both files naming the line verbatim.
      Waits on: the check, so the wording matches what it actually reads.

- [x] **The contract** — `AGENTS.md` gate description, and the `payload` capability
      delta with its three EARS criteria.
      Spec: all three outcomes
      Closes on: the six gates green, `--anchors` resolving the new `Proof:` citations
      and passing every new criterion as EARS.
      Waits on: everything above.
