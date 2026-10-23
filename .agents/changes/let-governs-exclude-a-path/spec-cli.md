Targets: cli

# let-governs-exclude-a-path — the cli half

## Outcomes

1. The wiki's home footer counts a tracked file as governed when some capability governs it
   under spec-drift's rule: a plain entry matches and none of that capability's `!` entries
   does. Each capability's entries are matched as one list, never pooled with another's.
2. cli's `Governs:` gives up `readme_test.go`, `release_workflow_test.go`,
   `label_workflow_test.go` and `gates_test.go`.

## Scope boundaries

**In:** `governedSplit` and its caller in `cmd/libretto/wiki.go`; cli's `Governs:` line.

**Out:** how the line is rendered; any other wiki count.

## Constraints

- The rule is spec-drift's `governed_by`, mirrored; cli's promise that the footer's matching
  mirrors spec-drift's stands, now with the exclusion in it.
- Every existing wiki test stays unedited.

## Prior decisions

- **The wiki footer honours `!` per capability.** Decided 2026-10-23 with payload's
  exclusion: pooling every capability's globs would let one capability's `!` hide a path
  another claims, and two matching authorities is what cli already forbids.

## Task breakdown

- `governedSplit` per capability, with the exclusion; cli's `Governs:` line.

## Verification criteria

- **Where** a capability's `Governs:` holds a `!` entry matching a tracked file and no other
  capability claims it, the wiki footer **shall** count that file as unclaimed.
  Proof: cmd/libretto/wiki_test.go TestWikiFooterHonoursAnExclusion
- **Where** one capability excludes a tracked file another capability claims, the wiki footer
  **shall** count that file as claimed.
  Proof: cmd/libretto/wiki_test.go TestWikiFooterExclusionIsPerCapability
