# Delta — the wiki regenerates when a change lands

Targets: payload

A generated index nobody regenerates is drift wearing a marker comment. The moment
the specification moves is `record-work`'s landing step — the delta is applied onto
the capability spec there, so the wiki refresh belongs in the same breath and the
same commit.

## Verification criteria

- Where `libretto` is on PATH and the project holds a consolidated specs
  directory, the `record-work` skill shall instruct the landing step to run
  `libretto wiki` and include the refreshed index in the landing commit; where the
  binary is absent, it shall say the wiki may be stale and move on rather than
  block the landing.
  Proof: rg -q 'libretto wiki' skills/record-work/SKILL.md

The proof anchors only the instruction's presence — a command citation cannot test
that a session obeys prose. That limit is the same one every skill criterion in the
payload spec already lives with.
