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

## 2026-08-17 · consolidate-license-files · 6→7
Said: THIRD-PARTY.md gained a four-line explanatory paragraph that no outcome asked for — the In list named only the link lines
Did: added it to outcome 2 rather than deleting it. The prose is the half of this fix that is presentation rather than layout: the ask misread the root listing as offering alternatives, so a move with no explanation invites the same misreading. Scope arriving without an outcome behind it is still unnamed scope even when it is right

## 2026-08-17 · consolidate-license-files · 6→7
Said: outcome 3's Governs: widening exists only as an uncommitted working-tree change, so --anchors green for the new payload citation was measured against a dirty tree
Did: not a defect — box 2 was open and phase 8 in progress. Re-ran every gate after the landing commit rather than carrying the dirty-tree measurement into the report. Worth keeping because "green" measured mid-phase-8 is a real way to report a gate that was never run on what shipped

## 2026-08-17 · consolidate-license-files · 6→7
Said: "LICENSE does not move" and "the vendored table is not touched" were confirmed by inspection with nothing to catch a regression — is the root LICENSE worth one line in the same test?
Did: yes, added. The existing loop says which files must leave the root and is satisfied by a root with no licence at all, which is the wrong tidy — GitHub reads root LICENSE for the displayed licence and the API field. Proved it can fire by exporting the branch with git archive and moving LICENSE aside, so the repo was never touched. The vendored table stays unguarded: check-payload already fails loudly on a parse returning nothing or everything

## 2026-08-17 · add-contributing-guide · 6→7
Said: criterion 1's "work does not come from a tracker" clause was tested by nothing — deleting the whole section from a copy left the guard at zero failures
Did: one assertion per clause, anchored on the heading rather than a sentence so improving the prose does not break the guard. This is the recorded "a criterion can cite a gate that tests half of it" class, and the delta declared no ceiling for it — a four-clause criterion citing one test needs four assertions or a stated ceiling, never the citation alone

## 2026-08-17 · add-contributing-guide · 6→7
Said: "the six gates as a runnable block" was held by one of the six — five gate lines removed from a copy and the guard stayed green
Did: all six asserted individually. The value of that section is that a contributor can paste and run it, so five commands vanishing undetected destroys the outcome while leaving the criterion green. A block asserted by one of its lines is a block nothing holds

## 2026-08-17 · add-contributing-guide · 6→7
Said: `## What to expect on review` is content no outcome names — scope arrived without asking
Did: made it an outcome rather than deleting it. Second time in this batch the same finding landed, after THIRD-PARTY.md's paragraph — writing prose that is clearly wanted still needs an outcome behind it, or the contract does not cover what shipped

## 2026-08-17 · add-contributing-guide · 6→7
Said: unverified question — the guide attached `--force` to the gates, while AGENTS.md uses it about never overwriting what the tool did not create
Did: reworded. Gates have no `--force`; the sentence conflated two rules and would have taught a contributor something untrue about the gates. Raised as a question rather than a finding, and it was the most wrong line in the file

## 2026-08-17 · split-readme-into-sections · 6→7
Said: the `"costs a line"` anchor could not fail on the docs side — the same diff that moved the argument into FLOW.md introduced the identical phrase into DESIGN.md, and the assertion reads both files concatenated, so deleting the guarded paragraph left the guard green
Did: re-anchored on "more expensive than this one". **Then made the same mistake again in the fix**: the replacement anchor for the vendoring move, "not thin, it is broken", was already in FLOW.md at base, so it passed on a pre-existing sentence in the wrong file. An anchor must be counted against the BASE of every destination document, not just against the branch — and counted the way flat() counts, because rg is line-scoped and these phrases wrap. Third and fourth instances of the guard-that-reads-green-while-unable-to-fire class in this one file

## 2026-08-17 · split-readme-into-sections · 6→7
Said: outcome 1 claims five relocated arguments and only four anchors exist — the vendoring move had nothing holding it, so it is the one relocation that could be silently reverted
Did: fifth anchor added. The count in the outcome and the count in the test are two places holding one number, which is the drift AGENTS.md opens by naming — and here the spec was the copy that was right

## 2026-08-17 · split-readme-into-sections · 6→7
Said: the plan's box-1 evidence describes five anchors failing on both ends, and four anchors existed — the plan is what the next reader trusts
Did: corrected. Writing evidence from what the run intended rather than from what it observed is the failure `evidence` exists to prevent, committed in the artifact that exists to prevent it

## 2026-08-17 · split-readme-into-sections · 6→7
Said: unverified question — three behaviour statements left the README beyond the phrases the table names: the unresolvable-provider fallback, "never a variable that holds a secret", and the merges/tags half of attacca's refusal
Did: all three restored compactly. The capability's own prior decision settles it — what a command *does* is a reference fact and stays, why it is that way is an argument and goes. The second one is a security promise, which is on the never-scoped-out list. And anchoring on "never merges, tags or releases" had made a behaviour fact unmentionable in the README: an anchor chosen from the *fact* rather than the *reasoning* drags the fact out with the argument

## 2026-08-17 · split-readme-into-sections · 6→7
Said: the CONTRIBUTING assertion scans the whole README while the outcome names Learn more — a link arriving in the footer would satisfy the guard without satisfying the outcome
Did: scoped to the Learn more section. The footer already mentions the licence and third-party files, so it is exactly where a stray link would plausibly land

## 2026-08-17 · add-payload-index · 6→7
Said: outcome 1 unmet for 4 of 36 rows with the gate green — caveman, caveman-commit, ponytail and ponytail-debt write `description: >` as a folded YAML scalar, so the page carried a literal `>` and no text. The delta's constraint claimed "a description: is a single frontmatter line here, which is how every payload item is already written", and that was false at HEAD
Did: the parse now reads the two block forms that exist here, `>` and `|` with or without a chomping indicator, and nothing else — not a YAML parser, because payload records that a dependency added to check a document is a dependency added to check a checker. **The lesson is not the parser.** The delta named this as a ceiling and deferred the fix to "the day an item needs a folded description", and that day was already four items in the past. A ceiling is a statement about the future; before writing one, check whether the present already violates it

## 2026-08-17 · add-payload-index · 6→7
Said: with four rows having no text after the name, `sort` decides their order on punctuation alone — byte order puts `caveman-commit` first, glibc's punctuation-ignoring collation puts `caveman` first, and CI runs on ubuntu-latest. So the gate could report drift on a page nobody edited
Did: pinned `LC_ALL=C` inside the pipeline, so the caller's locale cannot move it and the page a contributor generates is the page CI compares. Verified byte-identical under C, en_US.UTF-8 and en_GB.UTF-8. **The reviewer could not observe GNU sort on this machine and said so rather than concluding** — an unverified question reported as one is worth more than a confident guess, and pinning the collation settles it without needing the observation

## 2026-08-17 · add-payload-index · 6→7
Said: the page claims "what libretto install links into a target, and nothing else" — false for OpenCode, whose agents are installed by writing a derived file rather than linking, and Codex takes skills only
Did: reworded after verifying `func (o Opencode) Transforms(k Kind) bool { return k == Agents }`. A generated page's boilerplate is prose nobody reviews on regeneration, so a false claim in it survives every future run of the generator — the one part of a generated file that needs reviewing like hand-written text

## 2026-08-17 · add-payload-index · 6→7
Said: the new comment hardcodes "22 skills, 7 agents and 7 commands" — a typed count sitting next to a generator whose stated reason for existing is that a typed count drifted
Did: removed the count. Writing the failure and then committing it two lines below the explanation

## 2026-08-17 · add-payload-index · CI
Said: gates failed on CI with `CONTRIBUTING.md links to .agents/changes/, which does not exist` — red on the branch, green on this machine
Did: de-linked the path; it is named in prose now. **The mechanism is the lesson.** Git does not track empty directories, so once the last change folder landed, `.agents/changes/` was absent from a fresh checkout. It passed locally because `main` still had six queued proposals on disk, and it passed on the two branches before it because each still had folders left — so the guard was green three times for a reason that had nothing to do with correctness. Link resolution reads the working tree, so a link to a directory that exists only while work is unfinished passes while work is unfinished. Reproduced with `git archive`, which drops empty directories exactly as a checkout does — that is the way to see what CI sees without waiting for CI

## 2026-08-18 · retire-finished-changes · 6→7
Said: write-plan is modified but not in the delta's In: list, and it is the one added mandate here with no criterion and no check_wiring row — the other two each got one
Did: added it to scope, gave it a criterion and a row matching `never the list`, and forced that row red before believing it. The unwired mandate was the one added late, as a fix to something the cutter found; the two planned ones were wired without being asked

## 2026-08-18 · retire-finished-changes · 6→7
Said: Targets: payload omits cli while the delta lands on .agents/specs/cli/spec.md, and no gate catches it — targets_of feeds only the non-blocking drift warning
Did: Targets: payload cli. Left unfixed, a later staged edit under cmd/libretto/** would not have been credited to this delta

## 2026-08-18 · retire-finished-changes · 6→7
Said: Task breakdown still says four decisions, contradicted by Scope boundaries and by the landing commit, which retires six
Did: corrected it. Same shape as the "ten over eleven directories" failure AGENTS.md already records — a number stated in two places, fixed in one

## 2026-08-18 · add-specs-wiki · 6→7
Said: TestWikiWritesIndexAndSections asserts the index-row links but not the per-section [full spec] link, so half of criterion 2 could regress with the proof still green
Did: added the two per-section link assertions, re-ran the test, green. The known half-a-clause pattern, caught by a fresh reviewer
