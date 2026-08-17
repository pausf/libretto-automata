# make-test-badge-live — delta

Targets: readme

**This is a bug, not a feature.** The README's tests badge is a hardcoded green image:

```
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](.github/workflows/gates.yml)
```

It says *passing* whether or not anything passed. That is the failure this repository
refuses everywhere else — `skills/evidence/` exists because nothing is true until it has
been observed, and a badge asserting a green run nobody observed is a claim the flow would
reject in a report.

So the criterion comes before the fix, and its `Proof:` is red when it is written. The
`readme` capability has seven criteria and **not one of them touches a badge** — that is
the hole the bug came through.

## Outcomes

1. **The tests badge reports a real run.** Its image URL is GitHub's workflow badge
   endpoint for `gates.yml`, which turns red on its own when the workflow fails.

2. **No badge in `README.md` asserts a status from a hardcoded literal.** A badge whose
   image is a `shields.io/badge/…` literal may state a fact that does not change — a
   language version, a tool name, a licence — but never a *run outcome*. The forbidden
   words are the ones that claim a result: `passing`, `failing`, `pass`, `fail`, `build`,
   `coverage` — **matched word-bounded, never as substrings.**

   Word-bounded because the `readme` capability has already paid for the other choice: its
   command guard shipped as a substring match, and `commands/libretto-stat.md` would have been
   satisfied by the existing `/libretto-status` line. Here a substring match would refuse a
   legitimate badge for any tool whose name contains `pass` or `build`, and a guard that
   false-positives on honest content is a guard somebody deletes.

   The other six badges are version and tool labels. They stay, and they are honest as they
   are.

## Scope boundaries

**In:** the one badge line in `README.md`, and one new test case in
`cmd/libretto/readme_test.go` — the file the `readme` capability already governs.

**Out**, named:

- **the other six badges are not touched.** `Go-1.26`, `Bubble Tea-1.3`, `Lip Gloss-1.1`,
  `Claude Code-payload`, `Jira CLI-tracker`, `License-MIT` — none of them claims a run
  outcome, and rewriting a working label to match a convention nothing needs is churn
- **no coverage badge, and no new workflow.** The ask mentioned CI visibility; CI already
  exists and is visible, so the only defect is the false claim. Adding a coverage badge
  would be a new promise arriving on the back of a bug fix
- **no release or version badge.** `AGENTS.md` and the `readme` capability both forbid a
  version number in prose, and a version badge is a version number in prose with an image
  around it
- **the badge is not asserted to be *green*.** A test that required a passing badge would
  fail the moment the workflow legitimately went red, which is the whole point of a live
  badge. The test holds the *endpoint*, never the colour
- **`.github/workflows/gates.yml` is not changed.** The workflow is fine; the badge lied
  about it

## Constraints

- **The proof reads files as text.** No Markdown parser and no network call — the
  `readme` capability already records both, and a test that fetched the badge would fail on
  a machine with no network and turn a documentation guard into a flake
- the new case lives in `cmd/libretto/readme_test.go` and reuses that file's existing
  `repoFile` and `flat()` helpers. A new package for one more case on the same document is
  not worth its own directory — the same decision that capability already recorded
- **whitespace is normalised before matching**, via `flat()`. That capability records a
  guard that silently could not fire because the README wrapped between two words
- the repository slug in the badge URL is `pausf/libretto-automata`
- all six gates stay green

## Prior decisions

- **The endpoint form is `…/actions/workflows/gates.yml/badge.svg`** — *assumed, nobody was
  asked.* `/libretto-attacca` answered the contract stop. It is GitHub's documented
  per-workflow badge URL and it needs no token. What changes if it is wrong: the badge
  renders broken rather than false, which is visible on the front page immediately and is
  strictly better than the current state either way.
- **The guard does two things, and the second is not the thing that was rejected** —
  *assumed.* It forbids status *words* in any literal badge, **and** it pins this one badge's
  endpoint. What was rejected is pinning the URL *as the only check*: that would hold this
  badge and let the next hardcoded `tests-passing` in beside it. The pin is still needed,
  because a word-list check alone is satisfied by deleting the badge — the guard would go
  green on a README with no tests badge at all, which is not the outcome anybody wants.
  Broad rule plus one anchor, and neither covers the other. **Ceiling, and it is a real one:** the word list
  is a list, so a literal badge claiming `[![Tests](…/badge/tests-green-…)]` passes,
  having claimed a run outcome with a word nobody listed. That is the same class of ceiling
  the `readme` capability already names about its command-name guard — it proves a name
  appears, not that a row is honest. The replacement, the day somebody invents a new word
  for green, is a criterion about what a badge may link to, not a longer word list.
- **`build` and `coverage` are on the forbidden list even though no such badge exists
  today** — *assumed.* They are the two most likely next false badges, and the cost of
  listing them now is two words.
- The bug's criterion lands on the `readme` capability rather than on `ci`. `ci` governs
  `.github/**` and `Makefile`; the false claim was in `README.md`, and `readme` is the
  capability that owns it. The hole was `readme` having no badge criterion at all.

## Task breakdown

1. **Write the failing test first.** Add `TestNoBadgeAssertsAStatus` to
   `cmd/libretto/readme_test.go` and watch it fail against the README as it stands, naming
   the `tests-passing` badge. A red run is the only evidence the criterion describes this
   bug and not a neighbouring one.
2. **Then fix the badge line** in `README.md`, and watch the same test go green.

**These are one box.** The test alone leaves the suite red; the fix alone leaves the hole
open for the next false badge. Neither half merges on its own — the rule
`skills/write-plan/` now carries, applied here.

## Verification criteria

- **no badge in `README.md` whose image URL is a `shields.io/badge/` literal contains a
  status word** — `passing`, `failing`, `pass`, `fail`, `build`, `coverage`, matched
  **word-bounded** and case-insensitively, in the image URL only, so a link target or alt text
  mentioning a build is not caught and a tool whose name contains one of the words is not
  refused.
  Proof: cmd/libretto/readme_test.go TestNoBadgeAssertsAStatus

  **Watched red before green.** The run that added it failed naming the `tests-passing`
  badge, then passed once the endpoint replaced it.

- **the tests badge points at the `gates.yml` workflow badge endpoint**, so the guard cannot
  be satisfied by deleting the badge entirely.
  Proof: cmd/libretto/readme_test.go TestNoBadgeAssertsAStatus

- every relative link in the README still resolves after the edit.
  Proof: cmd/libretto/readme_test.go TestReadmeLinksResolve
