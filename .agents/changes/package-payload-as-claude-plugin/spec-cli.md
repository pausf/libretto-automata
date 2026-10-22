# package-payload-as-claude-plugin — cli

Targets: cli

What `install`, `update`, `doctor` and `status` do when Claude Code already has the payload
as the `libretto` plugin. The vocabulary, the decision log and the approach are in
`proposal.md`; `internal/dist` answers *whether*, this decides *what follows*.

## Outcomes

**The plugin wins the Claude tool, and nothing is served twice by this tool's hand.** While
`dist.ClaudePlugin` reports `libretto@libretto-automata` installed:

- **`install`, either Claude scope**, plans nothing and links nothing for the Claude tool,
  prints one line — `claude: served by the libretto plugin (libretto@libretto-automata) —
  nothing linked` — and exits 0. A skip it chose is not a skip it suffered, so the "non-zero
  if anything was skipped" rule does not fire.
- **`update`** relinks nothing for the Claude tool, by reaching the same path.
- **`doctor`** reports each owned Claude link as **served twice**, naming
  `libretto uninstall --yes` as the remedy, and exits non-zero while any exists. A missing
  Claude item is **not** attention — the plugin serves it, and asking for the link would
  ask for the duplicate.
- **`status`** prints one line naming the plugin above the Claude items, which it still
  lists in their real states.
- **`--codex` and `--opencode`** never ask the question and behave exactly as before.

**When the record cannot be read**, `install` links as it always has, and `doctor` names the
file and the error. "Could not tell" never blocks anything.

## Scope boundaries

**In:** the four commands above, for the Claude tool only.

**Out:**

- **removing links** on finding the plugin. `uninstall --yes` and `prune --yes` are unchanged
  and are the remedy `doctor` names.
- **`--json` shapes.** `status --json` and `doctor --json` carry no plugin key in this change;
  adding one later is additive and needs no schema bump.
- **the panel**, which keeps showing Claude items as missing and offering install.
- **`models`.** It keeps listing and writing the linked agents; plugin agents live in Claude
  Code's cache, which this tool did not create.
- **`loop`'s prompt.** It keeps `/libretto-attacca`; `proposal.md`'s box 2 observes that it
  resolves with only the plugin loaded, and a failure stops the change.

## Constraints

- **Tests set `CLAUDE_HOME` to `t.TempDir()` and `CLAUDE_CODE_PLUGIN_CACHE_DIR` to empty**,
  and write their own `installed_plugins.json` fixture. No test reads or writes a real one.
- **The plugins root comes from the Claude target's `Root()`** through `dist.PluginsRoot`,
  never from a second lookup of the home directory.
- **Exit codes keep their meanings.** `doctor`'s non-zero is "something needs attention",
  which a duplicate is; `install`'s non-zero stays a conflict or a failure.

## Prior decisions

- **ASSUMED, `/libretto-attacca`, 2026-10-22 — the plugin wins the Claude tool; `install`
  defers rather than warns.** A warning alone keeps `install` creating the duplicate that a
  linked agent then shadows (`~/.claude/agents/` outranks a plugin's agents), so a plugin
  update would change nothing the flow runs. What changes if wrong: the skip becomes a
  `doctor` line and `install` links as today.
- **ASSUMED — the skip exits 0.** Deferring is the outcome the user chose by installing the
  plugin, not a refusal. What changes if wrong: one exit code, and the panel's install would
  report a failure on every run for plugin users — which is why it is not.
- **ASSUMED — an unreadable record links as before.** What changes if wrong: `install`
  refuses the Claude tool, non-zero, naming the file.
- **ASSUMED — `status` names the plugin in text only.** What changes if wrong: an additive
  `plugin` key in `status --json`.

## Task breakdown

1. The Claude install path asks `dist.ClaudePlugin` and defers on `true`; `update` inherits it.
2. `doctor` reclassifies Claude items while the plugin is installed and names an unreadable
   record.
3. `status` prints the one line.

## Verification criteria

- **While the libretto plugin is installed, `install` for the Claude tool shall create no
  link in either scope.**
  Proof: cmd/libretto/plugin_test.go TestInstallLinksNothingForClaudeWhileThePluginIsInstalled
- **While the libretto plugin is installed, `install` for the Claude tool shall print one
  line naming `libretto@libretto-automata` and saying nothing was linked.**
  Proof: cmd/libretto/plugin_test.go TestInstallNamesThePluginWhenItSkipsClaude
- **When `install` defers the Claude tool to the plugin, the command shall exit 0.**
  Proof: cmd/libretto/plugin_test.go TestInstallSkippingForThePluginExitsZero
- **While the libretto plugin is installed, `install --codex` and `install --opencode` shall
  produce the same links they produce with no record present.**
  Proof: cmd/libretto/plugin_test.go TestPluginDoesNotChangeCodexOrOpencode
- **If the install record cannot be read, then `install` shall link the Claude tool as it
  does with no record present.**
  Proof: cmd/libretto/plugin_test.go TestUnreadableRecordLinksAsBefore
- **While the libretto plugin is installed, `update` shall relink nothing for the Claude
  tool.**
  Proof: cmd/libretto/plugin_test.go TestUpdateRelinksNothingForClaudeWhileThePluginIsInstalled
- **While the libretto plugin is installed, `doctor` shall report every owned Claude link as
  served twice, naming `libretto uninstall --yes`.**
  Proof: cmd/libretto/plugin_test.go TestDoctorReportsLinksThePluginAlsoServes
- **While the libretto plugin is installed, `doctor` shall not count a missing Claude item as
  needing attention.**
  Proof: cmd/libretto/plugin_test.go TestDoctorDoesNotAskForLinksThePluginServes
- **While the libretto plugin is installed and an owned Claude link exists, `doctor` shall
  exit non-zero.**
  Proof: cmd/libretto/plugin_test.go TestDoctorExitsNonZeroOnADuplicate
- **If the install record cannot be read, then `doctor` shall name the file and the error.**
  Proof: cmd/libretto/plugin_test.go TestDoctorNamesAnUnreadableRecord
- **While the libretto plugin is installed, `status` for the Claude tool shall print one line
  naming the plugin before its items.**
  Proof: cmd/libretto/plugin_test.go TestStatusNamesThePlugin
