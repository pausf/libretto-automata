# Plan — Fable in the model catalogue

Three tasks, strictly sequential. Task 2 is the one the user asked for; task 1 is what
makes its version column true, and task 3 is every place that enumerates the list.

- [x] **1 · `fable` resolves to `Fable 5` on the Anthropic API**
      Spec: task 1 · `internal/agentmodel/provider.go`
      Add `fable` to `anthropicAPI`. Nothing else: the third-party maps stay silent by
      decision, and `effortByVersion` already carries `Fable 5`.
      Closes: *`fable` resolves to Fable 5 on the Anthropic API and runs all five levels*
      Verified by: internal/agentmodel/provider_test.go TestFableResolvesAndRunsAllFiveLevels
      Waits on: nothing — starts now.

- [x] **2 · the `fable` catalogue row**
      Spec: task 2 · `internal/agentmodel/catalogue.go`
      One entry, last, after `opus`: alias, `Fable 5`, and a label that names no plan.
      Closes: *the catalogue lists `fable` last*, and the three existing invariants —
      valid, versioned, claims no plan.
      Verified by: internal/agentmodel/catalogue_test.go TestCatalogueListsTheSubscriptionModels
      Waits on: **task 1.** The row's version column asserts what the alias resolves to;
      adding the row first would state it before anything made it true.

- [x] **3 · every place that enumerates the legal values**
      Spec: task 3
      `catalogue_test.go`'s `want` list, the capability spec's model and provider tables,
      and the README's two `libretto models` listings.
      Closes: nothing new — it is what keeps the criteria above honest rather than
      passing against a stale list.
      Verified by: the six gates. `Proof:` citations live in the spec, not here — the
      anchors gate parses that keyword wherever it finds it, and a plan citing a command
      under it reads as a broken citation to a file named `gofmt`.
      Waits on: **task 2.** The test's `want` list and the code have to move together or
      the suite is red between them.
