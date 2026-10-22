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
2. One sentence: **the plugin or `libretto install` for Claude Code, never both** — with the
   plugin installed, `libretto install` links nothing for Claude and `doctor` reports any link
   left over.
3. Under *Updating*, the plugin's own: `claude plugin marketplace update libretto-automata`
   then `claude plugin update libretto@libretto-automata`, since a third-party marketplace
   does not auto-update by default.

**Criterion 2 of the readme spec is amended**, because "a Go version floor, one `go install`,
one `libretto install`" stops being the only count: the plugin route is a second list of
steps, still countable, still free of explanation.

**The README stays at or under its 380-line ceiling.** It is at 380 today (measured), so the
route arrives by moving lines out, never by raising the ratchet.

## Scope boundaries

**In:** the Install section, its Updating subsection, and the amended criterion 2.

**Out:**

- **explaining the plugin's versioning or the duplicate rule** in the README. Both are
  reasoning; they live in the distribution and cli specs, and in `docs/` if anywhere.
- **renaming commands in the first-run walk.** `/libretto-flow` stays the name shown; under the
  plugin it is `/libretto:libretto-flow`, and whether the bare name resolves is observed by
  `proposal.md`'s box 2.

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

## Task breakdown

1. The two `/plugin` lines and the one-or-the-other sentence in Install.
2. The plugin's update line under Updating.
3. The checkout subsection moved out to keep the ceiling.
4. Criterion 2 and `TestInstallSectionIsStepsOnly` amended to admit the plugin route.

## Verification criteria

- **The Install section shall name `/plugin marketplace add pausf/libretto-automata` and
  `/plugin install libretto@libretto-automata`.**
  Proof: cmd/libretto/readme_test.go TestInstallNamesThePluginRoute
- **The README shall state that the plugin and `libretto install` for Claude Code are
  alternatives.**
  Proof: cmd/libretto/readme_test.go TestReadmeSaysPluginAndLinksAreAlternatives
- **With the plugin route added, the Install section shall still carry none of
  `GOMODCACHE`, `@v0.` or "checksum database".**
  Proof: cmd/libretto/readme_test.go TestInstallSectionIsStepsOnly
- **`README.md` shall stay at most 380 lines with the plugin route in it.**
  Proof: cmd/libretto/readme_test.go TestReadmeStaysShort
