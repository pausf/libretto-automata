# Lessons

Append-only. Entries are spent by `libretto-retro` and gain a `Resolved:` line there;
they are never edited otherwise.

## 2026-08-13 · improve-metrics-report · phase 8

> no pasa las pipeline https://github.com/pausf/libretto-automata/pull/44 algo ha
> pasado , cuidado y apuntalo para libretto retro

Nothing was broken: `gates` passed twice and the red check was `require-release-label`
— the designed refusal that exists because attacca must never label a release. The
correction is about presentation: the attacca run reported the flow complete and the PR
open without saying, loudly and in its own line, **"the label check will now be red on
the PR, that red is the bump question reaching you, here is the one command"**. A
designed refusal that the report does not predict reads as a broken pipeline, costs the
user an alarm and a round trip, and the mode's whole visibility budget is that report.
Resolved: 2026-08-13 · project knowledge — convention recorded in .agents/specs/ci/spec.md
(prior decisions); the payload-side fix stays queued as ask-release-label-at-attacca-end

## 2026-08-13 · ship-frugal-model-defaults · phase 1

> /libretto-status: un proposal cuya rama ya lo despachó no va en la sección de cola, va
> en vuelo o en ninguna. Un reporte que induce a rehacer trabajo terminado es peor que uno
> que calla.

Said after `/libretto-status` listed `ship-frugal-model-defaults` under **the queue**
while its work was finished and open as PR #48. The scan is literal and reads the working
tree: the `Queued:` line comes out in the pickup commit, which lives on the feature
branch, so from `main` the proposal still carries it. Correct by the letter of the
definition and wrong in every way that matters — the next command typed was
`/libretto-attacca ship-frugal-model-defaults`, which would have branched from `main`,
found the proposal still queued there, and rewritten ~1500 lines against a spec that
already contained them. Two requests, one change, a guaranteed conflict on two capability
specs.

**The `Queued:` line is necessary and not sufficient.** A report that induces rework is
worse than silence, because silence costs a question and this costs the work twice.

## 2026-08-14 · close-flow-open-questions · 6→7
Said: code with no outcome behind it — the delta must name every file it moves
Did: rewrote commands/libretto-flow.md without an outcome naming it; the delta now names all three homes of the rule

## 2026-08-14 · close-flow-open-questions · 6→7
Said: a 6→7 entry is a finding wherever it sits, never a correction
Did: counted an orphan 6→7 entry into the corrections-outside-any-change line; exclusion now outranks the orphan count

## 2026-08-14 · add-system-diagram-to-readme · 6→7
Said: the outcome promised eight phases in order; the diagram named five nodes and the guard cannot count phases
Did: the flow diagram now carries all eight numbers (2-4 share a node, as the flow itself rides asking on the spec) and phase 8 precedes the push stop

## 2026-08-14 · add-system-diagram-to-readme · 6→7
Said: the constraint reads ASCII-safe and the middle dot U+00B7 is not ASCII
Did: label separators switched to plain periods; the contract held as written rather than reworded around the code

## 2026-08-14 · add-system-diagram-to-readme · 6→7
Said: nothing in the artifacts records that the renders were looked at
Did: the looking happened (mermaid-cli, both diagrams, tilde verified at 3x) and the phase 7 report now names it as evidence

## 2026-08-14 · add-multi-tool-skill-targets · 6→7
Said: the env count lives in cli/spec.md, help and README, and all three move or drift
Did: only cli/spec.md carries a count ("Five"); the delta now names the one real home and carries the two replacement rows for landing

## 2026-08-14 · add-multi-tool-skill-targets · 6→7
Said: every occurrence of Claude outside the allowlist fails the gate
Did: the check filtered whole lines, so an addressee sharing a line with "Claude Code" escaped; allowlisted tokens are now deleted per-hit before the re-search, forced red on a mixed line

## 2026-08-14 · add-multi-tool-skill-targets · 7→8
Said: en project solo se puede instalar en claude no tiene sentido; el menú se haría muy largo, hay que comprimirlo
Did: shipped a flat destination list where only claude had a project side — reworked to orthogonal tool × scope axes (one row per tool, s flips the scope) before anything released; the free-text answer "que el usuario pueda elegir" had been read as "global-only phase 1" when it meant the full choice

## 2026-08-14 · add-opencode-command-target · 6→7
Said: docs/STATE.md still says commands for opencode remain queued, and this change is what shipped them
Did: not fixed — the sentence sits under *decisions not to relitigate*, which AGENTS.md marks ask-first, so the edit is the user's call; the exact replacement wording is in the phase 7 report

## 2026-08-14 · add-opencode-command-target · 6→7
Said: the README edit rewrote two prose paragraphs the delta and plan named nowhere, in a file the delta was careful to say it does not own
Did: spec-cli task 3 and plan task 6 now name the paragraphs and why the correction is required rather than optional; the prose itself was accurate and stays

## 2026-08-14 · add-opencode-command-target · 6→7
Said: the change folder carries promises that are also in the capability specs now
Did: pending, not missed — phase 8 applies the delta and deletes the folder in the commit that lands it

## 2026-08-14 · add-opencode-command-target · 6→7
Said: a queue capture (adapt-payload-wording-to-three-hosts) rides this feature branch into its PR
Did: named in the phase 7 report and the request description; the capture is the honest home for the half of the ask this change does not do
