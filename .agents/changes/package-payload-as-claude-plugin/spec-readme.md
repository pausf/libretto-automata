# package-payload-as-claude-plugin — readme

Targets: readme

The front door learns the second route in. The vocabulary, the decision log and the approach
are in `proposal.md`.

## Outcomes

**Install names two routes for Claude Code and says they are alternatives.** Under
`## Install`, after the `go install` steps:

1. For Claude Code only: `/plugin marketplace add pausf/libretto-automata`, then
   `/plugin install libretto@libretto-automata`. Codex and OpenCode still take
   `libretto install --codex` / `--opencode`.
2. One sentence carrying the anchor phrase **"never both"**: the plugin or `libretto install`
   for Claude Code, never both — with the plugin installed for the user, `libretto install`
   links nothing for Claude and `doctor` reports any link left over.
3. Under *Updating*, the plugin's own: `claude plugin marketplace update libretto-automata`
   then `claude plugin update libretto@libretto-automata`, since a third-party marketplace
   does not auto-update by default.

**The README stays at or under its 380-line ceiling.** It is at 380 today (measured), so the
route arrives by moving lines out, never by raising the ratchet.

### Outcome 2 and its criterion, replaced at landing

The readme spec's Outcome 2 currently reads:

> 2. **Install is steps a reader can count.** A Go version floor, one `go install …@latest`,
>    one `libretto install`, and one read-only command to confirm it landed. No explanation
>    of the module cache, no `$GOMODCACHE` tree, no version-in-the-path argument.

Landing replaces it with:

> 2. **Install is steps a reader can count, on one of two routes.** A Go version floor, one
>    `go install …@latest`, one `libretto install`, and one read-only command to confirm it
>    landed. For Claude Code, alternatively, the two `/plugin` lines — marketplace add, then
>    install — and one sentence saying the two routes are **never both** for Claude Code. No
>    explanation of the module cache, no `$GOMODCACHE` tree, no version-in-the-path argument,
>    and no explanation of how the plugin is versioned.

Its criterion currently reads:

> - **Outcome 2** — the install section contains the install commands and a version floor,
>   and none of `GOMODCACHE`, `@v0.`, `checksum database`.
>   (cited proof: cmd/libretto/readme_test.go TestInstallSectionIsStepsOnly)

Landing replaces it with the three Install criteria below — the existing guard and the two
plugin-route criteria — and adds the Updating criterion beside them.

## Scope boundaries

**In:** the Install section, its Updating subsection, the amended Outcome 2 and its criteria.

**Out:**

- **explaining the plugin's versioning or the duplicate rule** in the README. Both are
  reasoning; they live in the distribution and cli specs, and in `docs/` if anywhere.
- **renaming commands in the first-run walk.** `/libretto-flow` stays the name shown; under the
  plugin it is `/libretto:libretto-flow`, and whether the bare name resolves is observed by
  `proposal.md`'s box 2.
- **pinning any other sentence.** The readme's *Wording* boundary — "No criterion here pins a
  sentence, except the relocated arguments in outcome 4, where the phrase *is* the anchor" —
  gains exactly one more named exception, **"never both"**, for the same reason: it is the
  anchor a paraphrase would lose. Landing adds it to that boundary by name.

## Constraints

- **`TestReadmeStaysShort`'s 380 is not raised.**
- **`TestInstallSectionIsStepsOnly` keeps refusing `GOMODCACHE`, `@v0.` and "checksum
  database"** in the Install section. The pinned-tag example therefore stays out of the
  README: `@v0.` is one of its refused strings, and the tag form is reference for the docs.

## Prior decisions

- **ASSUMED, `/libretto-attacca`, 2026-10-22 — the lines come from moving *From a checkout
  instead* into `docs/DESIGN.md`**, beside the "why the payload is not compiled into the
  binary" section it already links, leaving a one-line pointer. It is the developer's path,
  not the user's. What changes if wrong: it moves to `CONTRIBUTING.md` instead, which makes
  `contributing` a fourth capability on this change.
- **ASSUMED — the README carries no `@v0.x` pin example**, because the existing guard refuses
  `@v0.` in Install. What changes if wrong: the guard is narrowed, which is its own decision.
- **ASSUMED, after review-spec — "never both" is pinned, as a named exception to the wording
  boundary.** A criterion that "the README states they are alternatives" cannot be failed
  without an anchor. What changes if wrong: the criterion is dropped and the sentence is
  review's to see.
- **ASSUMED, after review-spec — the Updating lines are a criterion of their own.** What
  changes if wrong: they are prose under Updating with no test.

## Task breakdown

1. The two `/plugin` lines and the "never both" sentence in Install.
2. The plugin's update lines under Updating.
3. The checkout subsection moved out to keep the ceiling.
4. Outcome 2, its criterion and the wording boundary amended as written above.

## Verification criteria

- **The Install section shall name `/plugin marketplace add pausf/libretto-automata` and
  `/plugin install libretto@libretto-automata`.**
  Proof: cmd/libretto/readme_test.go TestInstallNamesThePluginRoute
- **The Install section shall carry the phrase "never both" in the sentence naming the plugin
  and `libretto install`.**
  Proof: cmd/libretto/readme_test.go TestReadmeSaysPluginAndLinksAreAlternatives
- **With the plugin route added, the Install section shall still carry none of
  `GOMODCACHE`, `@v0.` or "checksum database".**
  Proof: cmd/libretto/readme_test.go TestInstallSectionIsStepsOnly
- **The Updating subsection shall name `claude plugin marketplace update libretto-automata`
  and `claude plugin update libretto@libretto-automata`.**
  Proof: cmd/libretto/readme_test.go TestReadmeUpdatingNamesThePluginRoute
- **`README.md` shall stay at most 380 lines with the plugin route in it.**
  Proof: cmd/libretto/readme_test.go TestReadmeStaysShort
