# Delta — Fable in the model catalogue

Targets: agent-models
Governs: internal/agentmodel/catalogue.go internal/agentmodel/provider.go

Claude's `fable` alias is absent from the catalogue, so `libretto models` cannot offer it
and no agent can be put on Fable without editing a file by hand — which is the one thing
this capability exists to remove.

Half the plumbing is already here and was built without the entry: `pinPattern` already
reads `fable` out of a model id, and `effortByVersion` already knows `Fable 5` runs all
five levels. The alias is the missing row, not the mechanism.

## Outcomes

`fable` becomes a legal value of `model:`, everywhere the other three already are:
`libretto models` lists it, `libretto models set fable <agents>` writes it, `Valid("fable")`
is true, and the panel's selector offers it. `libretto models effort` offers it all five
levels, because that is what Fable 5 runs.

The capability spec's model table gains one row:

| Value | Resolves to | Means |
|---|---|---|
| `fable` | Fable 5 | the most capable, and the most expensive per token |

- **It sorts last, after `opus`.** The catalogue's order is contracted, not incidental —
  cheapest first, so the cheap choice sits under the cursor of a feature whose purpose is
  to reduce the bill. Fable's tokens cost more than Opus's, so it goes at the far end.
  This is the existing contract applied, not a new one.
- **`fable` resolves to `Fable 5` on the Anthropic API**, and the version column carries
  that with the same `Resolved` date as the rest. No second date.
- **Its label names no plan.** Same standing refusal as every other row: the binary cannot
  ask Anthropic what this account is paying for, so it does not say.

## Scope boundaries

**In:** one `catalogue` row, `fable` in the Anthropic-API alias map, and the tests and spec
tables that enumerate the legal values.

**Out, named so it cannot arrive quietly:**

- **A `fable` row in the third-party provider tables.** See prior decisions — absent is a
  deliberate answer, not an omission.
- **A cost or price column.** "More expensive than Opus" is in the label as prose. A number
  in a table is a second thing that decays, and unlike the version column nothing here
  makes the staleness visible.
- **`mythos`.** A sibling of Fable that is invitation-only, so an alias offered to everyone
  would mostly be an alias that does not resolve. Not asked for either.
- **Any change to the effort catalogue.** Fable 5 runs the same five levels; the table is
  already right and `effortByVersion` already says so.

## Constraints

- **The catalogue's order is a contract**, read by both the CLI listing and the panel
  selector. An entry inserted in the middle moves what sits under the cursor.
- **`TestCatalogueListsTheSubscriptionModels` asserts the exact list, in order.** It fails
  the moment a row is added, which is the point — the list is not something to extend
  silently.
- **Every catalogue invariant already under test applies to the new row**: it carries a
  label, it names its version, and its label claims no plan. Three existing tests cover it
  the moment it exists.

## Prior decisions

- **The third-party provider tables get no `fable` entry.** *(Assumed under
  `/libretto-attacca`, 2026-08-12 — nobody was asked.)* `providers` in `provider.go`
  transcribes the host's own per-provider alias table, and that table is not readable from
  here for Fable. An alias absent from a provider's map is already this package's honest
  answer — *not knowable* — and unknown is treated as capable, so `fable` on Bedrock offers
  all five levels, which is what Fable 5 runs anyway. **What changes if this is wrong:** one
  map entry per provider that serves it. Nothing user-visible moves unless a provider serves
  a Fable with fewer than five levels, and no such model exists today.

  **Ceiling named, because verifying this found it:** the listing falls back to the
  catalogue's own `Version` column when a provider cannot resolve an alias, and prints the
  trailer `resolved for Amazon Bedrock` underneath. So the `fable` row there reads
  `Fable 5` on the catalogue's authority while the trailer implies the provider's. That is
  pre-existing behaviour for every unresolvable alias — it is how the gateway case already
  reads — and `fable` is the first row to hit it on a *named* provider rather than on a
  gateway. What would lift it is distinguishing "resolved here" from "the catalogue's
  claim, unresolved here" in the version column, which is a change to how every row
  renders and is out of scope for one catalogue entry.
- **The label names no plan tier.** *(Assumed, same run.)* `opus` says "Max plans, metered
  on Pro" because that is documented; which plans include Fable is not verifiable from here,
  and the catalogue's whole posture forbids guessing at it. **What changes if this is wrong:**
  one label string.
- **`fable` sorts last rather than beside `opus`.** *(Assumed, same run.)* Read off the
  existing cheapest-first contract and Fable's per-token cost. **What changes if this is
  wrong:** one line's position, and which row the selector opens on.

## Task breakdown

1. `internal/agentmodel/provider.go`: `fable` resolves to `Fable 5` on the Anthropic API.
2. `internal/agentmodel/catalogue.go`: the `fable` row, last, with its version and label.
3. The tests and docs that enumerate the legal values: the catalogue list test, the spec
   tables, the README listing.

## Verification criteria

- the catalogue lists `fable` last, after `opus`
  Proof: internal/agentmodel/catalogue_test.go TestCatalogueListsTheSubscriptionModels
- `fable` is a legal frontmatter value and an unknown name still is not
  Proof: internal/agentmodel/catalogue_test.go TestEveryCatalogueEntryIsValid
  Proof: internal/agentmodel/catalogue_test.go TestUnknownModelIsRefused
- it names the version its alias resolves to, and its label claims no plan
  Proof: internal/agentmodel/catalogue_test.go TestEveryRealModelNamesItsVersion
  Proof: internal/agentmodel/catalogue_test.go TestNoLabelClaimsTheUserHasAPlan
- **`fable` resolves to Fable 5 on the Anthropic API and runs all five effort levels**
  Proof: internal/agentmodel/provider_test.go TestFableResolvesAndRunsAllFiveLevels
- an agent can be moved onto `fable` and keep a declared effort, because Fable runs it
  Proof: internal/agentmodel/effort_test.go TestApplyModelKeepsEffortWhenTheModelSupportsIt
