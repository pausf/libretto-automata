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

## 2026-08-14 · adapt-payload-wording-to-three-hosts · 6→7
Said: the named ceiling understates the hole — any occurrence of the marker satisfies the gate, including one in prose unrelated to any mandate
Did: fixed the criterion's ceiling to describe the real behaviour, reproduced with a probe file ("talk about a host's own dog here" passes with ok); the regex is left alone deliberately, because judging relevance is the paragraph parser the delta named out of scope

## 2026-08-14 · adapt-payload-wording-to-three-hosts · 6→7
Said: the delta was applied onto the capability spec but the change folder was not deleted in the same commit, which AGENTS.md requires
Did: fixed properly this time — the folder deletion was amended into the same commit rather than added as a second one. The previous change deferred the same finding to a follow-up commit; the tension is that review-work needs a committed diff to review, and amending an unpushed commit resolves it without rewriting shared history
## 2026-08-14 · add-transformed-agent-targets · 6→7
Said: the widened Owned is not scoped to targets that generate anything, and `prune --claude --yes` offered to delete a hand-written file that merely carried a marker line
Did: fixed — Owned is symlink-only again, OwnedEither is the widened question, and the kind is carried on the Entry so only a generated kind asks it. Watched red by reverting the split: both the unit refusal and the end-to-end prune test fail

## 2026-08-14 · add-transformed-agent-targets · 6→7
Said: a marker naming the repo root or any directory inside it was accepted, and prune --yes deleted both
Did: fixed — strictly inside the root, and a path that exists and is a directory is refused. A missing source still stays ours, or prune could never remove an orphan. Watched red

## 2026-08-14 · add-transformed-agent-targets · 6→7
Said: "moving the repository orphans every generated file and prune is the remedy for both" is false — the file becomes foreign, so prune skips it and can never remove it
Did: fixed the sentence. The remedy is uninstall from the old checkout before moving, or a manual delete after. True of the symlink arm too, and the sentence was new in this change

## 2026-08-14 · add-transformed-agent-targets · 6→7
Said: create on a generated kind loses the "appeared since the scan" refusal, because os.Rename replaces silently where os.Symlink returns file exists
Did: fixed — create uses os.Link (ErrExist when the destination exists), repoint keeps rename because replacing is its intent. The first test written for this passed either way because create's own Lstat caught the case; replaced with one that asserts the write primitive itself, then watched red

## 2026-08-14 · add-transformed-agent-targets · 6→7
Said: the change folder was not deleted, so the same contract text lives in two places
Did: `git rm -r -q -f` silently did nothing because the files were untracked in this session, and -q hid the error. Removed with rm -rf and amended into the one commit. Never trust a quiet git rm on paths that may be untracked

## 2026-08-14 · add-transformed-agent-targets · 6→7
Said: the targets criterion still read "skills and commands for opencode" while the test it cites now asserts three kinds — a criterion contradicting its own proof; and apply.go still claimed it never removes a real file
Did: both fixed. The second is the more interesting one: the comment was true when written and the change is what falsified it

## 2026-08-14 · add-transformed-agent-targets · 6→7
Said: the marker is emitted as an unquoted YAML scalar, so a checkout path containing " #" would be read as a comment and OpenCode throws rather than skipping
Did: the marker is now always emitted double-quoted, and the reader accepts quoted or bare. Six awkward paths covered by test — spaces, #, colon, quote, backslash

## 2026-08-17 · make-test-badge-live · 6→7
Said: the word-check half scans the raw README, and badgeImage's `[^)\s]+` cannot cross a newline — so a literal status badge whose markdown wraps between `![Build](` and its URL slips past entirely, while the six honest badges keep the count non-zero and the guard green. The delta's own Constraints demanded flat() normalisation and only the endpoint half did it
Did: the scan now runs over flat(readme), with `\s*` after the paren for what flat leaves behind. Proved with the reviewer's own probe — a wrapped `build-passing` badge appended to the README, exit 1 naming `build`, exit 0 once removed. **Second time this file has shipped a guard that silently could not fire**: the capability already records `"checksum database"` never matching because the README wrapped between those two words. A new substring guard in this file starts from flat(), not from the raw document
