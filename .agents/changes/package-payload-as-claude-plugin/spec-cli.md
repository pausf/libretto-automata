# package-payload-as-claude-plugin — cli

Targets: cli

What `install`, `update`, `doctor` and `status` do when Claude Code already has the payload
as the `libretto` plugin. The vocabulary, the decision log and the approach are in
`proposal.md`; `internal/dist` answers *whether*, this decides *what follows*.

## Outcomes

**Where the record is read.** The install record is always looked up under the plugins root
derived from the **global** Claude root — `dist.PluginsRoot(target.Global().Root())` — for
both scopes. A `--project` run of `install` or `doctor` never looks under
`<cwd>/.claude/plugins`: Claude Code keeps one plugins root per user, and a project's
`.claude/` holds none.

**Only a user-scope install of the plugin makes the CLI defer.** `installed_plugins.json`
records each install with a `scope`; the CLI defers only when `libretto@libretto-automata`
carries a `user`-scope install, because only that one serves every project the links would
serve. A `project` or `local` install serves one repository, so it suppresses nothing.

**While the user-scope plugin is installed, the plugin wins the Claude tool:**

- **`install`, either Claude scope**, plans nothing and links nothing for the Claude tool,
  prints one line — `claude: served by the libretto plugin (libretto@libretto-automata) —
  nothing linked` — and exits 0.
- **`update`** relinks nothing for the Claude tool, because its release route's relink is the
  production `install` path and reaches the same deferral.
- **`doctor`** reports each owned Claude link as **served twice**, naming
  `libretto uninstall --yes` as the remedy, and exits non-zero while any exists. A missing
  Claude item is **not** attention — the plugin serves it, and asking for the link would ask
  for the duplicate.
- **`doctor --json`** lists each such link in `problems` with `state: "served twice"` — a
  doctor-only problem state, documented beside the `link.State` spellings in the `--json`
  shapes table. So `problems` is non-empty whenever `doctor` exits non-zero for a duplicate.
- **`status`** prints one line naming the plugin above the Claude items, which it still lists
  in their real states. **Under `--json` that line is not printed**: stdout stays exactly one
  JSON object, with no plugin key.
- **`--codex` and `--opencode`** never ask the question and behave exactly as before.

**A `project` or `local` install of the plugin** links nothing away: `install` links as
always, and `doctor` prints one informational line naming each such scope and its project
path, exiting 0 on that alone.

**When the record cannot be read**, `install` links as it always has, and `doctor` names the
file and the error. "Could not tell" never blocks anything.

### The standing paragraph this amends, in place

The cli spec's *Outcomes* currently reads:

> **Exit codes carry meaning.** Non-zero whenever something was not linked or needs
> attention, so a script can gate on it. A conflict counts: an item that did not get
> linked is an incomplete install, whatever the reason.

Landing replaces it with:

> **Exit codes carry meaning.** Non-zero whenever something was not linked or needs
> attention, so a script can gate on it. A conflict counts: an item that did not get
> linked is an incomplete install, whatever the reason. **A Claude target skipped because
> the user-scope `libretto` plugin serves it is neither "not linked" nor "needs attention"**:
> the payload is installed, by the route the user chose, so `install` and `update` exit 0 on
> that skip. A link the plugin also serves *is* attention, and `doctor` exits non-zero on it.

## Scope boundaries

**In:** the four commands above, for the Claude tool only, and the one `doctor --json`
problem state.

**Out:**

- **removing links** on finding the plugin. `uninstall --yes` and `prune --yes` are unchanged
  and are the remedy `doctor` names.
- **a plugin key in `status --json` or `doctor --json`.** The only JSON change is the
  `served twice` state inside `problems`; anything more is additive later, with no schema bump.
- **reading `<cwd>/.claude/plugins`** for a project-scope run. Claude Code has no such root.
- **the panel**, which keeps showing Claude items as missing and offering install.
- **`models`.** It keeps listing and writing the linked agents; plugin agents live in Claude
  Code's cache, which this tool did not create.
- **`loop`'s prompt.** It keeps `/libretto-attacca`; `proposal.md`'s box 2 observes that it
  resolves with only the plugin loaded, and a failure stops the change.

## Constraints

- **Tests set `CLAUDE_HOME` to `t.TempDir()` and `CLAUDE_CODE_PLUGIN_CACHE_DIR` to empty**,
  and write their own `installed_plugins.json` fixture. No test reads or writes a real one.
- **The plugins root comes from the global Claude target's `Root()`** through
  `dist.PluginsRoot`, never from a second lookup of the home directory, and never from the
  project target.
- **Exit codes keep their meanings, as the amended paragraph above states them.** `doctor`'s
  non-zero is "something needs attention", which a link served twice is and a project or
  local plugin install is not. `install`'s and `update`'s non-zero is something not linked —
  a conflict or a failure — and a Claude target the user-scope plugin serves is not that.
- **`update`'s proof drives the real relink closure** built in `update.go`, which calls the
  production `install` — never an injected `relink`. The deferral lives in `install`, so a
  test with a spy relink would pass whether or not `update` honours it.

## Prior decisions

- **ASSUMED, `/libretto-attacca`, 2026-10-22 — the plugin wins the Claude tool; `install`
  defers rather than warns.** A warning alone keeps `install` creating the duplicate that a
  linked agent then shadows (`~/.claude/agents/` outranks a plugin's agents), so a plugin
  update would change nothing the flow runs. What changes if wrong: the skip becomes a
  `doctor` line and `install` links as today.
- **ASSUMED — the skip exits 0, and the standing exit-code paragraph is amended in place to
  say so.** Deferring is the outcome the user chose by installing the plugin, not a refusal.
  What changes if wrong: one exit code, and the panel's install would report a failure on
  every run for plugin users.
- **ASSUMED, after review-spec — only a `user`-scope install defers.** A project or local
  install serves one repository; suppressing global links for it would strip the payload from
  every other project. What changes if wrong: any scope defers, and the informational line
  goes.
- **ASSUMED, after review-spec — the record is read under the global Claude root for both
  scopes.** What changes if wrong: a project run also reads `<cwd>/.claude/plugins`, which
  Claude Code does not write.
- **ASSUMED, after review-spec — `served twice` is a doctor-only problem state in
  `doctor --json`**, not a sixth `link.State`. The link's own state is still `linked`; the
  problem is the pairing with the plugin, which only `doctor` judges. What changes if wrong:
  `link.State` gains a value and `status --json` carries it too.
- **ASSUMED — an unreadable record links as before.** What changes if wrong: `install`
  refuses the Claude tool, non-zero, naming the file.
- **ASSUMED — `status` names the plugin in text only, and not under `--json`.** What changes
  if wrong: an additive `plugin` key in `status --json`.

## Task breakdown

1. The Claude install path asks `dist.ClaudePlugin` under the global plugins root and defers
   on a user-scope install; `update`'s real relink inherits it.
2. `doctor` reclassifies Claude items while the user-scope plugin is installed, emits
   `served twice` in text and in `--json`, names project and local installs informationally,
   and names an unreadable record.
3. `status` prints the one line, text only.

## Verification criteria

- **While the user-scope libretto plugin is installed, `install` for the Claude tool shall
  create no link in either scope.**
  Proof: cmd/libretto/plugin_test.go TestInstallLinksNothingForClaudeWhileThePluginIsInstalled
- **While the user-scope libretto plugin is installed, `install` for the Claude tool shall
  print one line naming `libretto@libretto-automata` and saying nothing was linked.**
  Proof: cmd/libretto/plugin_test.go TestInstallNamesThePluginWhenItSkipsClaude
- **When `install` defers the Claude tool to the user-scope plugin, the command shall exit 0.**
  Proof: cmd/libretto/plugin_test.go TestInstallSkippingForThePluginExitsZero
- **When `install` runs in project scope, it shall read the install record under the global
  Claude root's plugins directory, never under `<cwd>/.claude/plugins`.**
  Proof: cmd/libretto/plugin_test.go TestProjectScopeInstallReadsTheGlobalRecord
- **While the record lists `libretto@libretto-automata` only at `project` or `local` scope,
  `install` for the Claude tool shall link as it does with no record present.**
  Proof: cmd/libretto/plugin_test.go TestProjectOrLocalPluginInstallDoesNotSuppressLinks
- **While the user-scope libretto plugin is installed, `install --codex` and
  `install --opencode` shall produce the same links they produce with no record present.**
  Proof: cmd/libretto/plugin_test.go TestPluginDoesNotChangeCodexOrOpencode
- **If the install record cannot be read, then `install` shall link the Claude tool as it
  does with no record present.**
  Proof: cmd/libretto/plugin_test.go TestUnreadableRecordLinksAsBefore
- **While the user-scope libretto plugin is installed, `update`'s release route, relinking
  through its real relink closure, shall link nothing for the Claude tool.**
  Proof: cmd/libretto/plugin_test.go TestUpdateRelinksNothingForClaudeWhileThePluginIsInstalled
- **While the user-scope libretto plugin is installed, `doctor` shall report every owned
  Claude link as served twice, naming `libretto uninstall --yes`.**
  Proof: cmd/libretto/plugin_test.go TestDoctorReportsLinksThePluginAlsoServes
- **While the user-scope libretto plugin is installed, `doctor` shall not count a missing
  Claude item as needing attention.**
  Proof: cmd/libretto/plugin_test.go TestDoctorDoesNotAskForLinksThePluginServes
- **While the user-scope libretto plugin is installed and an owned Claude link exists,
  `doctor` shall exit non-zero.**
  Proof: cmd/libretto/plugin_test.go TestDoctorExitsNonZeroOnADuplicate
- **While the user-scope libretto plugin is installed, `doctor --json` shall list each owned
  Claude link in `problems` with `state` `served twice`.**
  Proof: cmd/libretto/plugin_test.go TestDoctorJSONListsServedTwice
- **When `doctor` runs in project scope, it shall read the install record under the global
  Claude root's plugins directory, never under `<cwd>/.claude/plugins`.**
  Proof: cmd/libretto/plugin_test.go TestProjectScopeDoctorReadsTheGlobalRecord
- **While the record lists `libretto@libretto-automata` at `project` or `local` scope,
  `doctor` shall print one informational line naming each such scope.**
  Proof: cmd/libretto/plugin_test.go TestDoctorNamesAProjectOrLocalPluginInstall
- **If a project or local plugin install is `doctor`'s only finding, then `doctor` shall
  exit 0.**
  Proof: cmd/libretto/plugin_test.go TestDoctorExitsZeroOnAProjectOrLocalPluginInstallAlone
- **If the install record cannot be read, then `doctor` shall name the file and the error.**
  Proof: cmd/libretto/plugin_test.go TestDoctorNamesAnUnreadableRecord
- **While the user-scope libretto plugin is installed, `status` for the Claude tool shall
  print one line naming the plugin before its items.**
  Proof: cmd/libretto/plugin_test.go TestStatusNamesThePlugin
- **While the user-scope libretto plugin is installed, `status --json` shall write exactly
  one JSON object to stdout, carrying no plugin line and no plugin key.**
  Proof: cmd/libretto/plugin_test.go TestStatusJSONOmitsThePluginLine
