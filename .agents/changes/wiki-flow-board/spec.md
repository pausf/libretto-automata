# Delta — the flow board

Targets: cli

A page of the flow itself: where corrections surface, and which ideas wait.
Read from the ledgers at generation time, absent when they are.

## Outcomes

When the project holds a lessons ledger or queued proposals, the wiki carries a
flow article at `#flow` — reachable from the home like any capability page —
showing corrections counted by the phase field of the ledger's entry headers,
as labeled proportional bars, and the queued proposals with their `Queued:`
dates, oldest first. Neither source present, no article and no home link.

## Scope boundaries

In: the ledger parse, the queue listing, the article render, the home link.

Out, named: the per-change cost table — `libretto metrics` already answers it
at the CLI and the wiki version means reusing that git walker; it returns as
its own change when somebody wants the numbers on the page (assumed under
attacca). Ponytail-debt counts — same shape, same condition.

## Constraints

- The ledger's entry header is `## <date> · <change> · <phase>` — the same
  contract metrics reads; phases counted as written, no normalization beyond
  trimming.
- Determinism, absence arms, no clock — the standing wiki rules.
- Existing wiki proofs stay green unchanged.

## Prior decisions

- **The cost table is out with its condition named.** Assumed 2026-08-19: the
  CLI answers it today; the wiki version returns as its own change reusing the
  metrics walker.
- **The flow board is an article in the router's page system** — `#flow`, class
  `cap`, so navigation costs nothing new.

## Task breakdown

1. Parse + render + home link with the proofs, one commit; then the look.

## Verification criteria

- Where the project holds a lessons ledger with `## <date> · <change> · <phase>`
  entry headers, the page shall carry a flow article at `#flow` with one labeled
  bar per phase, widths proportional to the counts; where the ledger is absent
  or empty of entries, the article shall be absent.
  Proof: cmd/libretto/wiki_test.go TestWikiFlowBoardCountsCorrections
- Where queued proposals exist, the flow article shall list each with its
  `Queued:` date, oldest first; with none, the queue block shall be absent —
  and a queue alone shall still summon the article.
  Proof: cmd/libretto/wiki_test.go TestWikiFlowBoardListsTheQueue
- Where the flow article exists, the home shall link to `#flow`; where it does
  not, the home shall carry no such link.
  Proof: cmd/libretto/wiki_test.go TestWikiHomeLinksTheFlowBoard
- The additions shall keep every existing wiki proof green unchanged.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLIsDeterministic
