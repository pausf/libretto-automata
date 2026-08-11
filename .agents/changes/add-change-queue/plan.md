# add-change-queue — plan

Spec: spec.md (Targets: payload). One writer: the orchestrator.

- [x] 1. `commands/libretto-queue.md` — the capture loop; defines the `Queued:` line
      convention. Waits on: nothing.
      Closes: "both commands parse and every skill they reference exists" —
      scripts/check-payload · passed, exit 0, "all checks passed"
- [x] 2. `commands/libretto-next.md` — pick oldest-first, branch, remove `Queued:`,
      enter phase 2. Waits on: 1 (the convention it reads).
      Closes: scripts/check-payload · passed, exit 0
- [x] 3. `skills/find-work/SKILL.md` — queued proposals reported, never blocking, never
      source 1. Waits on: 1. Version 1.1 → 1.2.
      Closes: scripts/check-payload · passed, exit 0
- [x] 4. `commands/libretto-status.md` — the queue section, delegated to find-work's
      scan. Waits on: 3.
      Closes: scripts/check-payload · passed, exit 0
- [ ] 5. land the delta onto `.agents/specs/payload/spec.md`, delete this folder —
      phase 8, same commit as the last code. Waits on: 1–4.
      Closes: skills/record-work/spec-drift --anchors
