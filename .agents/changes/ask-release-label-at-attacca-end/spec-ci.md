# ask-release-label-at-attacca-end — delta

Targets: ci

One amendment, and it is the smallest kind: a prior decision that names this work as
queued stops being true the moment it lands.

## Outcomes

The prior decision *"A designed refusal is predicted where it will be read, or it reads as
breakage"* keeps every word about the prediction — that half is independent and still the
rule for any report ending at this repository's request. Its closing sentence, *"The fuller
fix — asking the bump natively at the end of the run — is queued as
`ask-release-label-at-attacca-end`"*, becomes the record that it was built, and where the
behaviour now lives.

## Scope boundaries

**In:** that sentence, in `.agents/specs/ci/spec.md`.

**Out:** the workflow, the label check, the release target, and every criterion under them.
`require-release-label` is unchanged and stays unchanged — it still refuses an unlabelled
request, still judges nothing about which bump is correct, and still holds no permissions.
**The red check is not what this change removes.** What changes is that somebody is now
asked before they meet it.

## Constraints

`.github/**` and `Makefile` are what this capability governs, and this delta touches
neither. It is a prose amendment to a decision, carried here rather than in the payload
delta because that is the capability the sentence lives in.

## Prior decisions

- **A pointer to queued work is a dated fact, not a permanent one.** Left alone it points
  at a change folder that phase 8 deletes on landing — a citation resolving to nothing, in
  the spec that is supposed to be what is true now.

## Task breakdown

1. `.agents/specs/ci/spec.md`: rewrite the closing sentence of that prior decision.

## Verification criteria

- **no spec cites `ask-release-label-at-attacca-end` as queued** once the change folder is
  gone, and every `Proof:` in the amended spec still resolves
  Proof: skills/record-work/spec-drift --anchors
