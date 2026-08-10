---
name: review-tests
description: "Trigger: reviewing a diff, a file or an MR for its tests — does the change carry its proof, do the tests test behaviour or implementation, and did any existing proof get quietly weakened. Test tampering is always a severe finding. Standalone: works on any diff, with or without review-project."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Two questions, asked of every change: **does the new behaviour arrive with its
proof**, and **did any existing proof get quietly weakened to let this diff pass?**

The second is the dangerous one. A missing test is visible debt; a weakened test is
a claim that something is still verified when it no longer is.

`evidence` governs: run the suite when it is runnable and report what was observed;
a coverage opinion formed without running anything says so.

## Question 1 — does the change carry its proof?

Match the diff's logic against the diff's tests:

- **new branch, loop, parser, boundary, money or security path with no test
  touching it** — name the untested path, not a percentage. "The error branch of
  the refund is untested" beats any coverage number.
- **tests that assert the implementation, not the behaviour** — mock-heavy tests
  that re-state the code's internals pass forever and catch nothing; a test that
  would survive a correct rewrite is testing the right thing
- **a bug fix with no regression test** — the bug that happened once is the one
  proven to be reachable; it gets a test that fails on the old code, always
- **flakiness being born** — sleeps standing in for synchronisation, order
  dependence between tests, time and randomness unpinned

Proportion cuts both ways, and the reviewed project's own conventions win: a
one-line change with no logic needs no test, and saying so is a legitimate answer.
A test asserting that a constant equals itself is itself a finding.

## Question 2 — did existing proof get weakened?

Read every hunk that touches an existing test as a claim needing justification.
The red flags, each one **severe until the MR's stated intent explains it**:

| In the diff | What it usually means |
|---|---|
| a skip or ignore added to a passing test | an invariant dropped, silently |
| an assertion loosened — exact to fuzzy, value to type, equals to not-throws | a constraint relaxed to fit the code |
| an expected value edited in the same MR that changes the behaviour | the test was bent to the code — ask which one was wrong, the requirement or the implementation |
| a test deleted "because it tested old code" | proof removed with nothing replacing it |
| a tolerance, timeout or retry widened | a real failure being reclassified as noise |

The question for each: **did the requirement change, or did the test just lose an
argument with the implementation?** Only the first justifies the edit, and the MR's
description should say so. A legitimate behaviour change replaces the old proof
with a new one — it does not just erase it.

## Output

Per finding: location (`file:line`), which question it fails, what was observed in
the hunk, severity, and the missing or restored test in one line — named
concretely, "a test that calls X with an empty list", never "add tests".

If the suite was run, its result in one line, observed. If it was not runnable,
that, stated.

No findings is a statement: **"the change carries its proof, and no existing proof
was weakened."**

This lens reports; it never blocks, edits, commits or pushes.
