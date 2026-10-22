# package-payload-as-claude-plugin

Tracker: none
Targets: distribution
Durable decisions: every entry in Prior decisions below, and the ones in spec-cli.md and spec-readme.md

## The ask, verbatim

> Estratégica: empaquetar el payload como plugin de Claude Code (marketplace,
> versionado) para el target Claude, dejando el CLI para Codex/OpenCode. Los plugins son
> LA unidad de distribución del ecosistema desde principios de 2026.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Ship the payload as a versioned Claude Code plugin for the Claude target, alongside the
symlinking CLI which keeps covering the other tools.

## Decision log

### Session 2026-10-22 (attacca — nobody was asked; every answer below is assumed unless it says measured)

- Q: can the repository itself be the marketplace and the plugin, with no file moving? →
  A: **yes** (measured). A plugin's default component directories are `skills/`,
  `agents/` and `commands/` at the plugin root, and the manifest is optional
  (code.claude.com/docs/en/plugins-reference, *Standard layout*). Two plugins installed on
  this machine are exactly this shape — `obra/superpowers` and `DietrichGebert/ponytail`
  each carry `.claude-plugin/marketplace.json` with one entry whose `source` is `"./"`, and
  their `skills/` sit at the repository root. `claude plugin validate .` on this worktree,
  with no manifest at all, passes: "This folder holds a skills tree … anything else in it
  (agents/, commands/, hooks/, …) is validated and loads with it." The only additions are
  two JSON files under `.claude-plugin/`.
- Q: what are the plugin and the marketplace called? → A: **plugin `libretto`,
  marketplace `libretto-automata`; the install id is `libretto@libretto-automata`**
  (assumed). The plugin name is the namespace every component is listed under, so it is
  the short one the user already types as the binary: `/libretto:libretto-flow`, not
  `/libretto-automata:libretto-flow`. The marketplace takes the repository's name, which
  is not reserved and does not start with `claude-`. If wrong, both names change before the
  first install anywhere — after it, a rename is the marketplace's `renames` map and every
  user's `enabledPlugins` key moves.
- Q: where does the version come from, given "the binary never carries a hardcoded
  version"? → A: **nowhere in a file. Neither `plugin.json` nor the marketplace entry
  carries `version`, and Claude Code derives it from git** (assumed; the rule is measured
  from the vendor docs, the behaviour is not yet observed — see *Risks*). The docs'
  resolution order is manifest `version`, then entry `version`, then the source type; for
  "a relative path inside a Git-hosted marketplace" it is "the commit SHA of the installed
  directory" (code.claude.com/docs/en/plugins/loading, *How Claude Code computes the
  version*). A pinned string would be the second source AGENTS.md forbids, and it would
  pin users: "a manifest that pins `version` keeps every user on the cached copy until its
  author changes the string, however many commits they push." `release.yml` is untouched.
  A user who wants a release rather than `main` adds the marketplace at a tag —
  `/plugin marketplace add pausf/libretto-automata@v0.6.0` — which the vendor accepts as
  `owner/repo@ref`. If wrong — the derived version proves unusable — the fallback is a
  release step that writes `version` from the tag it is about to push, which means the
  workflow commits to `main`; that is its own change with its own `ci` delta.
- Q: if a user has both, every skill exists twice. What is the rule? → A: **the plugin
  wins the Claude tool. While the install record names `libretto@libretto-automata`,
  `libretto install` links nothing for the Claude tool in either scope and says so in one
  line, exit 0; `doctor` reports every owned Claude link as served twice and exits
  non-zero; Codex and OpenCode are untouched** (assumed). The duplicate is measured on this
  machine today: the session's skill listing carries `ponytail` and `ponytail:ponytail`,
  `using-git-worktrees` and `superpowers:using-git-worktrees` — the vendored link and the
  plugin copy, both paid for on every turn ("Every skill in the skill listing adds to your
  context on every turn", code.claude.com/docs/en/skills). Worse than the cost: a bare name
  resolves to the link, not the plugin ("The bare `/fancy` also invokes the skill unless
  another command already uses that name"), and a linked agent outranks a plugin agent of
  the same name (code.claude.com/docs/en/sub-agents, priority 4 over 5) — so a plugin
  update would change nothing the flow actually runs while the user believes it did.
  Existing links are never removed automatically — `libretto uninstall --yes` is the
  user's. If wrong — the user wants both — the skip becomes a `doctor` warning only, and
  install links as today.
- Q: does a bare `Skill(skill="find-work")` resolve inside the plugin? → A: **yes, by
  documentation for skills, unobserved for agents and commands, and the payload is left as
  it is because its names are prose a model reads** (assumed). Every payload skill carries
  frontmatter `name:` (measured, 24 of 24), and a plugin skill with `name` answers to the
  bare name too. For plugin agents the docs say only that they are `<plugin>:<name>` and
  that `--agent` finds a bare name; nothing says `subagent_type` does. The payload already
  relies on prose rather than literal calls on OpenCode (targets' *Prior decisions*,
  2026-08-14), and the listing a model reads names `libretto:spec-writer`. Box 2 observes
  it before any CLI work. If wrong — a bare agent or command name fails in a plugin-only
  session — the payload's agent references and `libretto loop`'s `/libretto-attacca`
  prompt need the namespaced form, which is a payload change and a stop on this one.
- Q: is reading `installed_plugins.json` "settings.json handling", which AGENTS.md puts out
  of scope? → A: **no** (assumed). It is Claude Code's record of what is fetched, read and
  never written; `settings.json`'s `enabledPlugins` is not read at all. The cost: a plugin
  installed but disabled still counts as present, and the remedy is to uninstall it. If
  wrong, the detection is deleted and only the README's "one or the other" remains.
- Q: where is the record? → A: **`$CLAUDE_CODE_PLUGIN_CACHE_DIR/installed_plugins.json`
  when that variable is set, `<claude root>/plugins/installed_plugins.json` otherwise**
  (assumed, from code.claude.com/docs/en/plugins/loading, *Find plugins on disk*: "one
  plugins root, which is `~/.claude/plugins` unless you set `CLAUDE_CODE_PLUGIN_CACHE_DIR`").
  The claude root comes from the target, so `CLAUDE_HOME` relocates it and the suite stays
  off a real machine. If wrong, the variable is dropped and a relocated cache reads as no
  plugin — the pre-change behaviour.
- Q: an unreadable record? → A: **"could not tell" — install links as it does today and
  `doctor` names the file and the error** (assumed). Refusing to install on a corrupt file
  Claude Code itself rebuilds would block the user on another tool's state. If wrong,
  install refuses the Claude tool and exits non-zero.
- Q: which capability owns this? → A: **`distribution` primary; `cli` and `readme` as
  extra capabilities** (assumed). `.claude-plugin/**` is governed by nothing today
  (measured: no `Governs:` line claims it), and distribution is the capability that says
  where the payload an installed copy uses is — now one more answer. The detection lives in
  `internal/dist`; what `install`, `update`, `doctor` and `status` do with it is `cli`'s;
  the front door is `readme`'s. `targets` does not move: the claude root is already a
  target's answer, and the plugins root is derived from it the way the remembered
  destination is. `payload` does not move: no skill, agent or command changes. If wrong,
  the manifests move under `payload`'s `Governs:` and their tests into
  `scripts/check-payload`.
- Q: does distribution's "per-platform binaries, an `install.sh`, brew, a tap, npm … each
  adds a second way in" forbid this? → A: **no, and the scope line is amended to say why**
  (assumed). That list is about the binary. The plugin is the payload's way into one host,
  and it brings no binary: Codex and OpenCode still need `go install`. If wrong, the change
  stops here — the ask is exactly a second way in.
- Q: do the vendored third-party skills ship inside the plugin, and is that licensed? → A:
  **yes, unchanged, and `THIRD-PARTY.md` does not move** (assumed; the licences are
  measured). All three upstreams are MIT (`licenses/LICENSE-superpowers`,
  `LICENSE-ponytail`, `LICENSE-caveman`), which permits redistribution with the notice. The
  plugin root is the repository root, so the copy Claude Code caches carries `licenses/`
  and `THIRD-PARTY.md` beside the skills — the notice travels with the copy by
  construction. A test pins the repository tree under `./`; the cached copy is observed in
  box 2, never tested. The plugin lists them as `libretto:ponytail`, beside a
  user's own `ponytail:ponytail`; `THIRD-PARTY.md`'s *Naming* paragraph already says both
  can exist. If wrong — a vendored skill should not ship in the plugin — it has to leave
  `skills/`, because a manifest `skills` list adds to the default scan and never subtracts
  from it.
- Q: what is published? → A: **nothing to any third-party marketplace and nothing to
  Anthropic's directory; the repository is addable with
  `/plugin marketplace add pausf/libretto-automata`, which is in scope** (assumed, per the
  brief and AGENTS.md's ask-first on any push). Being addable needs no publication: the
  marketplace is the repository on its default branch. No `icon`, no directory listing
  fields. If wrong, a submission is its own change.
- Q: does `libretto models` reach plugin agents? → A: **no, and it is not made to**
  (assumed). The plugin's agents are a copy in Claude Code's cache, a file this tool did
  not create; writing a `model:` line there is the overwrite AGENTS.md forbids, and the next
  update would discard it. Plugin users run the `model:` lines the payload ships. If wrong,
  that needs a plugin-side mechanism, and none exists today: the plugin `settings` key
  takes only `agent` and `subagentStatusLine`.
- Q: does `libretto update` update the plugin? → A: **no** (assumed). `claude plugin
  update` is Claude Code's, and driving it from here acts on another tool's state.
  Third-party marketplaces default to auto-update off, so the README says how to update. If
  wrong, `update` shells to `claude plugin marketplace update` and `claude plugin update`.
- Q: the fork — how is it built? → A: **A: the repository root is the plugin and the
  marketplace, no version field, and `libretto install` defers to the plugin for the
  Claude tool while `doctor` reports a duplicate** (assumed, the recommended option). B —
  a `version` the release workflow writes from the tag — lost: the workflow would have to
  commit to `main` before tagging, and a written field is a second source the moment
  anything else edits it. C — a separate plugin directory (`plugin/claude-code/`, engram's
  shape) — lost: it needs the payload copied into it, since a symlink leading outside a
  plugin is rejected, and a copy is a second source of every skill. D — `doctor` warns, and
  `install` links as today — lost: `install` would keep creating the duplicate `doctor`
  then complains about, and the stale-link failure is silent until somebody reads `doctor`.

### Session 2026-10-22, after review-spec (the run acting as author on the reviewer's findings)

- Q: what does "set" mean for `CLAUDE_CODE_PLUGIN_CACHE_DIR`? → A: **non-empty. Empty or
  unset falls through to `<claude root>/plugins`** (assumed). An exported empty variable is
  how a shell clears one, and treating it as a root would read `installed_plugins.json` from
  the working directory. If wrong, an empty value means "no plugins root" and the CLI never
  defers.
- Q: which Claude root does a `--project` run derive the plugins root from? → A: **the
  global one, for both scopes; a project run never looks under `<cwd>/.claude/plugins`**
  (assumed). Claude Code keeps one plugins root per user (code.claude.com/docs/en/plugins/loading,
  *Find plugins on disk*), and a project's `.claude/` holds none. If wrong, a project run also
  reads `<cwd>/.claude/plugins`, which nothing writes.
- Q: does every install of the plugin make the CLI defer? → A: **only a `user`-scope install.
  A `project` or `local` install suppresses no links; `doctor` names it in one informational
  line and exits 0 on that alone** (assumed). `installed_plugins.json` records each install's
  `scope` (measured on this machine: every install here is `"scope": "user"`). A project
  install serves one repository; deferring the global links for it would strip the payload
  from every other project. If wrong, any scope defers and the informational line goes.
- Q: what does the licence criterion actually prove? → A: **the repository tree under the
  entry's source `./`, never Claude Code's cache copy — its ceiling, named in the criterion**
  (assumed). The cache copy is observed by box 2, not tested. If wrong, a test would have to
  install the plugin, which needs the `claude` binary the gates refuse.
- Q: which standing cli paragraph does the exit-0 skip amend? → A: **"Exit codes carry
  meaning. Non-zero whenever something was not linked or needs attention", amended in place;
  `spec-cli.md` carries the full replacement** (assumed): a Claude target skipped because the
  user-scope plugin serves it is neither "not linked" nor "needs attention", and exits 0. If
  wrong, the skip exits non-zero and the panel reports a failure on every plugin user's run.
- Q: how does `doctor --json` show a link served twice? → A: **in `problems`, with `state:
  "served twice"` — a doctor-only problem state documented beside the `link.State` spellings
  — so `problems` is non-empty whenever `doctor` exits non-zero for it** (assumed). If wrong,
  `link.State` gains a sixth value and `status --json` carries it too.
- Q: does `status --json` print the plugin line? → A: **no; stdout stays exactly one JSON
  object, with no plugin key, as already scoped out** (assumed). If wrong, an additive
  `plugin` key.
- Q: how is `update` proven to defer? → A: **`TestUpdateRelinksNothingForClaudeWhileThePluginIsInstalled`
  drives `update`'s release route with the real relink closure built in `update.go` — the
  production `install` path — never an injected `relink`** (assumed). A spy relink would pass
  whether or not `update` reaches the deferral. If wrong, the test is a spy and the deferral
  in `update` is review's to see.
- Q: how can "the README says they are alternatives" be failed? → A: **it pins the anchor
  phrase "never both", declared as a named exception to the readme's no-pinned-wording
  boundary, the way outcome 4's relocated arguments are** (assumed). If wrong, the criterion
  is dropped.
- Q: are the Updating lines proven? → A: **yes — `TestReadmeUpdatingNamesThePluginRoute`**
  (assumed). If wrong, they are untested prose.
- Q: what text replaces readme Outcome 2 and its criterion? → A: **`spec-readme.md` carries
  both in full, and landing applies them** (assumed). If wrong, landing writes them from
  intent, which is what the reviewer flagged.

## Outcomes

**The repository is a Claude Code marketplace holding one plugin, and the plugin is the
repository.** Two files are added and none move:

- `.claude-plugin/marketplace.json` — `name: libretto-automata`, an `owner`, a
  `description`, and one entry: `name: libretto`, `source: "./"`, a `description`. No
  `version`.
- `.claude-plugin/plugin.json` — `name: libretto`, `description`, `author`, `homepage`,
  `repository`, `license: MIT`. No `version`, and no `skills`, `agents` or `commands` key:
  the default directories are what loads, so a skill added to `skills/` ships with nothing
  else edited.

What a user does:

1. `/plugin marketplace add pausf/libretto-automata` — or `…@v0.6.0` to sit on a release.
2. `/plugin install libretto@libretto-automata`.
3. The payload lists as `libretto:find-work`, `libretto:spec-writer`, `/libretto:libretto-flow`.

**The version is git's.** Claude Code computes it from the installed commit; nothing in
the repository states it, and `release.yml` does not change.

**`internal/dist` answers one more "where is the payload": whether Claude Code has it as
this plugin, and at which scopes.** It reads the install record and reports a `user`-scope
install, any `project` or `local` installs, or could not tell — and it never writes. Only a
`user`-scope install makes the CLI defer. `cli` decides what that means for each command
(`spec-cli.md`).

**The vendored skills and their licences travel inside the plugin**, because the plugin
root is the repository root.

## Scope boundaries

**In:** the two manifests; the detection of `libretto@libretto-automata` in Claude Code's
install record; `install`, `update`, `doctor` and `status` honouring it for the Claude tool
(`spec-cli.md`); the README's plugin route (`spec-readme.md`); amending distribution's
"second way in" boundary.

**Out:**

- **publishing to any third-party marketplace or submitting to Anthropic's directory.**
  Ask first, per AGENTS.md. Also out: `icon`, `documentationUrl` and the other listing
  fields that only the directory reads.
- **a plugin for Codex or OpenCode** (`.codex-plugin/`, `.opencode/`). The ask leaves both
  to the CLI.
- **hooks, MCP servers, `bin/`, `userConfig`, `settings`.** The payload has none, and a
  top-level `bin/` would stop claude.ai and Cowork from installing the plugin at all.
- **reading `settings.json` or `enabledPlugins`.** A disabled-but-installed plugin counts
  as present. Ceiling named: the remedy is `/plugin uninstall`, never a setting.
- **removing links on the user's behalf** when the plugin is found. `libretto uninstall
  --yes` exists; `doctor` names it.
- **`libretto models` for plugin agents**, and **`libretto update` driving `claude plugin
  update`** — both act on files or state this tool did not create.
- **renaming the commands** to drop the `libretto-` prefix that reads twice under the
  plugin namespace. Renaming breaks every linked install's muscle memory for a cosmetic gain.
- **the panel.** With the plugin installed it shows the Claude items as missing, and its
  install reports the skip. Ceiling named: what brings it in is a user finding the panel
  misleading.
- **a `plugin` key in `status --json` or `doctor --json`.** Additive later, no schema bump.
- **the tag, the `release:` label, the push.** `record-work`'s; this is a `feat`, so a
  minor by the table, and the label stays the user's.

## Constraints

- **No new dependency.** `encoding/json`, `os`, `path/filepath`.
- **No test reads a real record.** Every test sets `CLAUDE_HOME` to `t.TempDir()` and sets
  `CLAUDE_CODE_PLUGIN_CACHE_DIR` to empty, so a developer whose shell exports it cannot leak
  their machine into the suite.
- **The six gates never need the `claude` binary.** CI has none. The manifests are pinned by
  Go tests reading the JSON; `claude plugin validate` is evidence observed locally, never a
  gate.
- **No path derived from `~/.claude` outside `internal/target`.** The plugins root is
  derived from the target's `Root()`, as the remembered destination already is.
- **Names the vendor reserves are avoided**: no plugin name starting with `claude-`, no
  marketplace name from the reserved list (code.claude.com/docs/en/plugins/marketplace-reference,
  *Reserved names*).
- **Nothing this tool did not create is written** — the record, the cache, the settings.

## Prior decisions

Carried, unchanged:

- **The module is how the binary and the payload arrive for Codex and OpenCode**, and
  `go install` stays the entry point there. The embedded payload and the release tarball
  stay rejected for the reasons this spec records.
- **The proxy, not `api.github.com`**, for the newest version — unchanged; the plugin has no
  version of its own to ask about.

New, every one assumed under `/libretto-attacca` on 2026-10-22 — each names what changes if
it is wrong in the decision log above:

- **ASSUMED — the repository root is both the marketplace and the plugin**, entry source
  `"./"`, measured against superpowers' and ponytail's installed layouts and a passing
  `claude plugin validate .`. A separate plugin directory lost because it needs a copy.
- **ASSUMED — plugin `libretto`, marketplace `libretto-automata`, install id
  `libretto@libretto-automata`.**
- **ASSUMED — no `version` in either manifest; Claude Code derives it from the installed
  commit.** A written version lost because the release workflow would have to commit, and
  because a pinned string freezes users until somebody edits it. **Ceiling named:** `/plugin`
  shows a 12-character hash rather than `v0.6.0`. What brings a field back: a user who
  needs the semver shown, met by a release step writing it from the tag — its own change.
- **ASSUMED — the plugin wins the Claude tool**; the rule itself is `cli`'s.
- **ASSUMED — the detection reads `installed_plugins.json` under the plugins root and never
  `settings.json`**; `CLAUDE_CODE_PLUGIN_CACHE_DIR` relocates the root.
- **ASSUMED — an unreadable record is "could not tell"**, never a refusal.
- **ASSUMED, after review-spec — `CLAUDE_CODE_PLUGIN_CACHE_DIR` counts only when non-empty.**
- **ASSUMED, after review-spec — the plugins root comes from the global Claude root for both
  scopes.**
- **ASSUMED, after review-spec — `ClaudePlugin` separates a `user`-scope install from
  `project` and `local` ones**, because only the first serves what the links serve.
- **ASSUMED, after review-spec — the licence test checks the repository tree, not the cache
  copy.** **Ceiling named:** that the cached copy carries `licenses/` is observed in box 2,
  never tested.
- **ASSUMED — "a second way in" is a boundary about the binary.** Amended in *Scope
  boundaries* at landing: the per-platform-binary list stays out; the payload's Claude Code
  plugin is in, and it brings no binary.
- **ASSUMED — the vendored skills ship in the plugin, licences beside them, with
  `THIRD-PARTY.md` unchanged.**
- **ASSUMED — bare skill and agent names inside the payload are left as they are**, on the
  strength of the docs for skills and of the payload's names being prose. Observed before
  the CLI work, or the change stops.

## Task breakdown

1. The two manifests, and Go tests pinning their names, their source, the absent version and
   the absent component keys.
2. The licence test: everything `THIRD-PARTY.md` links under `licenses/` exists under the
   entry's source.
3. An observed plugin-only session: bare skill, agent and command names, and a script run by
   path from the loaded plugin.
4. `dist.PluginsRoot` and `dist.ClaudePlugin` — the record, its location, the three answers.
5. `install` and `update` deferring to the plugin for the Claude tool (`spec-cli.md`).
6. `doctor` and `status` reporting it (`spec-cli.md`).
7. The README's plugin route (`spec-readme.md`).

## Verification criteria

- **The plugin manifest at `.claude-plugin/plugin.json` shall name the plugin `libretto`.**
  Proof: internal/dist/plugin_test.go TestPluginManifestNamesLibretto
- **The marketplace manifest at `.claude-plugin/marketplace.json` shall be named
  `libretto-automata`.**
  Proof: internal/dist/plugin_test.go TestMarketplaceIsNamedAfterTheRepository
- **The marketplace manifest shall list exactly one plugin, `libretto`, whose source is the
  repository root `./`.**
  Proof: internal/dist/plugin_test.go TestMarketplaceListsThePluginAtTheRoot
- **Neither the plugin manifest nor the marketplace entry shall carry a `version` key.**
  Proof: internal/dist/plugin_test.go TestPluginCarriesNoVersion
- **The plugin manifest shall declare no `skills`, `agents` or `commands` key.**
  Proof: internal/dist/plugin_test.go TestPluginManifestDeclaresNoComponentPaths
- **Every file `THIRD-PARTY.md` links under `licenses/` shall exist in the repository tree
  under the marketplace entry's source `./`** — ceiling named: this reads the repository,
  never Claude Code's cache copy, which box 2 observes.
  Proof: internal/dist/plugin_test.go TestPluginRootCarriesTheVendoredLicences
- **Where `CLAUDE_CODE_PLUGIN_CACHE_DIR` is non-empty, `PluginsRoot` shall return its value.**
  Proof: internal/dist/plugin_test.go TestPluginsRootHonoursTheCacheDirVariable
- **If `CLAUDE_CODE_PLUGIN_CACHE_DIR` is empty or unset, then `PluginsRoot` shall return
  `<claude root>/plugins`.**
  Proof: internal/dist/plugin_test.go TestPluginsRootTreatsAnEmptyCacheDirAsUnset
- **When the install record lists `libretto@libretto-automata` with a `user`-scope install,
  `ClaudePlugin` shall report a user-scope install.**
  Proof: internal/dist/plugin_test.go TestClaudePluginReadsTheInstallRecord
- **When the install record lists `libretto@libretto-automata` only at `project` or `local`
  scope, `ClaudePlugin` shall report no user-scope install and name each such scope.**
  Proof: internal/dist/plugin_test.go TestClaudePluginSeparatesProjectAndLocalInstalls
- **If no install record exists, then `ClaudePlugin` shall report the plugin not installed,
  with no error.**
  Proof: internal/dist/plugin_test.go TestClaudePluginWithNoRecordIsNotInstalled
- **If the record lists a plugin named `libretto` from another marketplace, then
  `ClaudePlugin` shall report the plugin not installed.**
  Proof: internal/dist/plugin_test.go TestClaudePluginIgnoresAnotherMarketplacesLibretto
- **If the install record cannot be read or is not JSON, then `ClaudePlugin` shall return
  an error naming the file.**
  Proof: internal/dist/plugin_test.go TestClaudePluginReportsAnUnreadableRecord
- **`ClaudePlugin` shall leave the install record byte-identical.**
  Proof: internal/dist/plugin_test.go TestClaudePluginWritesNothing

## Approach

### Summary

Add two JSON manifests under `.claude-plugin/` so the repository is its own marketplace
holding one plugin rooted at `./`, with no version field anywhere — Claude Code derives the
version from the installed commit. Teach `internal/dist` to read Claude Code's install
record and say whether `libretto@libretto-automata` is there, and teach `install`,
`update`, `doctor` and `status` to treat the plugin as the owner of the Claude tool, so the
same payload is never served twice. No skill, agent or command changes.

### Technical context

Go 1.26.5, stdlib only. Tests beside the code; the six gates in AGENTS.md; CI has no
`claude` binary.

Vendor evidence, read on 2026-10-22 against Claude Code 2.1.308:

- code.claude.com/docs/en/plugins-reference — the manifest is optional, `name` is the only
  required key, components default to `skills/`, `agents/`, `commands/` at the plugin
  root, "Every component is namespaced under it", `version` "keeps users on that version
  until you change it", component paths may not escape the plugin root.
- code.claude.com/docs/en/plugins/marketplace-reference — `.claude-plugin/marketplace.json`
  needs `name`, `owner`, `plugins`; an entry needs `name` and `source`; a relative source
  "must start with `./`"; `owner/repo@ref` pins a marketplace to a tag; the reserved names.
- code.claude.com/docs/en/plugins/loading — the version resolution order and the
  commit-SHA rule for a relative path in a Git-hosted marketplace; every non-local plugin
  is copied to `cache/<marketplace>/<plugin>/<version>/`; `installed_plugins.json` and
  `CLAUDE_CODE_PLUGIN_CACHE_DIR`; auto-update is off by default for third-party
  marketplaces.
- code.claude.com/docs/en/plugins/components and code.claude.com/docs/en/skills —
  `/<plugin>:<skill>`, the bare name for a skill with frontmatter `name`, `<plugin>:<agent>`,
  and plugin agents ignoring `permissionMode`, `hooks`, `mcpServers` (the payload's agents
  use none of them — measured).
- code.claude.com/docs/en/sub-agents — plugin agents are priority 5, below
  `~/.claude/agents/` at 4, so a linked agent of the same name wins over the plugin's.

Observed on this machine, read-only: `~/.claude/plugins/installed_plugins.json` (version 2,
keys `name@marketplace`, each an array of installs with `scope`, `installPath`, `version`,
`gitCommitSha`); `known_marketplaces.json`; `marketplaces/superpowers-dev/.claude-plugin/`
and the ponytail marketplace's `.claude-plugin/` — both repository-root plugins with
`"source": "./"`; engram's entry says `0.1.1` while its installed copy is `0.1.0`, which is
the manifest-over-entry precedence the docs state, seen rather than read. Every installed
plugin here pins a version, so **the derived-SHA path has no installed example** — see
*Risks*.

Blast radius: `.claude-plugin/marketplace.json`, `.claude-plugin/plugin.json` (new);
`internal/dist/plugin.go`, `internal/dist/plugin_test.go` (new); `cmd/libretto/main.go`
and a new `cmd/libretto/plugin_test.go`; `README.md`, `cmd/libretto/readme_test.go`;
three capability specs and `docs/SPEC.md` at landing.

### The approach

**Manifests.** Two files, written by hand once. The marketplace entry carries `name`,
`source: "./"` and `description`; `plugin.json` carries metadata only. No component paths,
so the default directories load and a new skill ships by existing. No `version` in either.

**Tests read the manifests as the repository ships them**, from `internal/dist` by relative
path (`../../.claude-plugin/…`), the way `cmd/libretto/readme_test.go` reads `README.md`.
Decoding into `map[string]any` rather than a struct is what lets "carries no `version`
key" be a test at all.

**Detection.** In `internal/dist/plugin.go`:

- `PluginsRoot(claudeRoot string) string` — `CLAUDE_CODE_PLUGIN_CACHE_DIR` when it is
  non-empty, otherwise `filepath.Join(claudeRoot, "plugins")`; an empty value is unset. The
  caller always passes the **global** Claude target's `Root()` — `target.Global().Root()` —
  for both scopes: a `--project` `install` or `doctor` never looks under
  `<cwd>/.claude/plugins`, because Claude Code keeps one plugins root per user.
- `ClaudePlugin(pluginsRoot string) (Installs, error)` — reads `installed_plugins.json`,
  decodes its `plugins` object, and for the key `libretto@libretto-automata` returns
  `Installs{User bool; Other []Install}`: whether any entry has `"scope": "user"`, and every
  `project` or `local` entry with its scope and project path. A missing file is the zero
  value and no error; unreadable or undecodable is the zero value and an error with the path
  in the message. One constant holds the id. **Only `User` makes the CLI defer**; `Other` is
  what `doctor` names in its informational line.

**The rule, applied in `cmd/libretto`** (`spec-cli.md` carries the promises): the Claude
tool's install path asks first and, on a user-scope install, prints one line naming the
plugin and saying nothing was linked, and returns success without planning; on an error, or
on project and local installs only, it links as today. The cli spec's standing exit-code
paragraph is amended in place to say a skip for the user-scope plugin is neither "not
linked" nor "needs attention". `update` reaches the same path because its release route's
`relink` closure in `update.go` calls the production `install`; its test drives that real
closure, never an injected one. `doctor` asks the same question, stops counting a missing
Claude item as attention while the user-scope plugin is installed, reports each linked
Claude item as served twice with `libretto uninstall --yes` as the remedy — in text, and in
`doctor --json` as a `problems` entry with the doctor-only state `served twice` — names
project and local installs in one informational line that does not affect the exit code,
and names an unreadable record. `status` prints one line above its Claude items, in text
only: under `--json` stdout stays one object. Codex and OpenCode never ask.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| B — `version` in `plugin.json`, written by `release.yml` from the tag | the workflow only tags today; writing the field means a commit on `main` before the tag, and the field is then a second source the first time anybody edits it by hand. Comes back if the hash shown in `/plugin` proves to be a real problem for users |
| C — a plugin directory of its own (`plugin/claude-code/`), as engram does | Claude Code rejects a symlink that leads outside the plugin, so the payload would be copied into it — a second copy of every skill, kept in step by hand or by a build step |
| D — `doctor` warns; `install` links as today | `install` keeps manufacturing the duplicate, and the silent failure — the plugin updates, the bare names keep resolving to the older links — lasts until somebody reads `doctor` |
| a pinned `version` matching the last tag, bumped by hand | the exact constant AGENTS.md forbids, and the docs say it freezes every user on the cached copy until somebody remembers |
| reading `enabledPlugins` to skip a disabled plugin | `settings.json` handling is out of scope by decision; the record is enough, and uninstalling a disabled plugin is one command |
| a `skills` list in the manifest naming each skill | a list is a second place every new skill has to be added; the default scan already loads exactly `skills/` |

### Risks

| Risk | What catches it |
|---|---|
| a bare agent name in `subagent_type` does not resolve to `libretto:<agent>` in a plugin-only session | box 2's observation; a failure stops the change, recorded in the decision log, before any CLI work builds on it |
| `/libretto-attacca` in `libretto loop`'s prompt does not resolve as a bare command for a plugin command | box 2 observes `claude -p '/libretto-attacca …'` with only the plugin loaded; same stop |
| the derived version is not the commit SHA, or does not move between commits, and updates never arrive | **unmitigated before merge**: a GitHub-marketplace install needs the manifests on `main`, and `--plugin-dir` loads in place with no version. Accepted, and reported by phase 7 as the first thing to observe after the merge — install from the marketplace in an isolated config and read `version` in its `installed_plugins.json` across two tags. If it fails, alternative B is the next change |
| executable bits lost in the cache copy — `skills/record-work/spec-drift` is run by path | box 2 runs `spec-drift --self-test` from the loaded plugin's path; the cache-copy case shares the risk above |
| a developer's local-directory marketplace loads the checkout in place, and `make build` puts `bin/libretto` on the Bash tool's `PATH` | accepted. The GitHub install has no `bin/` (it is ignored), and the local case runs the developer's own binary |
| `CLAUDE.md` at the plugin root: the docs say `claude plugin validate` warns about it | measured at box 1: with the manifests in place `claude plugin validate .` warns ("CLAUDE.md at the plugin root is not loaded as project context") and still passes. Harmless: the file is repository guidance, not plugin context, and `--strict` is not a gate here |
| a project-scope install of the plugin strips links every other project needs | only a `user`-scope install defers, pinned by `TestProjectOrLocalPluginInstallDoesNotSuppressLinks` |
| a user has another tool's plugin named `libretto` | the id check is the full `libretto@libretto-automata`, pinned by `TestClaudePluginIgnoresAnotherMarketplacesLibretto` |
| Claude Code changes the record's shape | the decode fails, which is "could not tell": install links as before and `doctor` names the file. No mitigation beyond that — the record is the vendor's |

### Validation and rollback

The six gates. New tests in `internal/dist/plugin_test.go` and `cmd/libretto/plugin_test.go`;
`readme_test.go` gains three. **Two have to be forced red before they are believed:**
`TestPluginCarriesNoVersion` — add a `version` to `plugin.json` and watch it fail, since a
test that decodes into the wrong type passes on any file;
`TestInstallLinksNothingForClaudeWhileThePluginIsInstalled` — remove the record from the
fixture and watch links appear, since an install that links nothing for any reason passes it.
And `TestUpdateRelinksNothingForClaudeWhileThePluginIsInstalled` drives `update`'s release
route with the **real** relink closure from `update.go` — the production `install` path —
with `latest`, `install` and `dirFor` stubbed; an injected `relink` would pass whether or not
`update` reaches the deferral.

Box 2 is evidence no gate can produce: a real Claude Code session. Its result goes into the
decision log, measured, before the CLI boxes start.

Rollback is one revert. Nothing migrates; no user state is written. A user who installed the
plugin keeps a working plugin from the cache until they uninstall it.

### Complexity deliberately kept

`CLAUDE_CODE_PLUGIN_CACHE_DIR` is honoured although no user here sets it. Without it, a
relocated cache reads as no plugin and the duplicate comes back silently — the exact failure
this change exists to remove, for the users least likely to notice.

## Checklist

Executed by `build-and-check`. Box 1 first. Box 2 waits on box 1 and gates boxes 3 and 4;
box 4 waits on box 3; box 5 waits on box 1; box 6 waits on all.

- [x] **The manifests, pinned.** `.claude-plugin/marketplace.json` and
      `.claude-plugin/plugin.json` as *The approach* describes; `internal/dist/plugin_test.go`
      decoding both as shipped; `claude plugin validate .` run locally and its output read.
      Closes: TestPluginManifestNamesLibretto, TestMarketplaceIsNamedAfterTheRepository,
      TestMarketplaceListsThePluginAtTheRoot, TestPluginCarriesNoVersion,
      TestPluginManifestDeclaresNoComponentPaths, TestPluginRootCarriesTheVendoredLicences
      (the repository tree under `./`; the cache copy is box 2's).
- [ ] **Observe a plugin-only session — evidence, no code.** In a Claude Code whose config
  **Stopped 2026-10-22:** an isolated config (`CLAUDE_CONFIG_DIR=$(mktemp -d)`) answers "Not logged in · Please run /login", and copying credentials into it is out of bounds. To resume: the user runs `! CLAUDE_CONFIG_DIR=<dir> claude /login` once, then this box runs the observations in that config with `--plugin-dir` pointing at this worktree. Boxes 3–5 wait on it by design.
      carries no libretto links — `CLAUDE_CONFIG_DIR` pointed at an empty directory and
      authenticated, or a person on a second account; **never by unlinking the real
      `~/.claude`** — load the checkout with `claude --plugin-dir <checkout>` and observe:
      `/libretto-flow` invoked bare; a command's `Skill(skill="find-work")` reaching
      `libretto:find-work`; `spec-writer` launched by its bare name; `claude -p
      '/libretto-attacca …'` resolving; `skills/record-work/spec-drift --self-test` run from
      the loaded plugin's path; and, once the manifests are installable from a marketplace in
      that isolated config, that the cached copy carries `licenses/` and `THIRD-PARTY.md`.
      Each result goes into the decision log as measured. **A run
      that cannot get an isolated authenticated session stops on this box rather than
      assuming it.** Closes: no criterion — it is the evidence the approach rests on, and a
      failure stops the change.
- [ ] **`install` and `update` defer to the user-scope plugin for the Claude tool, end to
      end.** `internal/dist/plugin.go` with `PluginsRoot` and `ClaudePlugin`; the Claude
      install path asking it under the global plugins root in both scopes;
      `cmd/libretto/plugin_test.go`, whose `update` test drives the real relink closure.
      Closes: TestPluginsRootHonoursTheCacheDirVariable,
      TestPluginsRootTreatsAnEmptyCacheDirAsUnset, TestClaudePluginReadsTheInstallRecord,
      TestClaudePluginSeparatesProjectAndLocalInstalls,
      TestClaudePluginWithNoRecordIsNotInstalled,
      TestClaudePluginIgnoresAnotherMarketplacesLibretto,
      TestClaudePluginReportsAnUnreadableRecord, TestClaudePluginWritesNothing,
      TestInstallLinksNothingForClaudeWhileThePluginIsInstalled,
      TestInstallNamesThePluginWhenItSkipsClaude, TestInstallSkippingForThePluginExitsZero,
      TestProjectScopeInstallReadsTheGlobalRecord,
      TestProjectOrLocalPluginInstallDoesNotSuppressLinks,
      TestPluginDoesNotChangeCodexOrOpencode, TestUnreadableRecordLinksAsBefore,
      TestUpdateRelinksNothingForClaudeWhileThePluginIsInstalled.
- [ ] **`doctor` and `status` report the plugin.** `served twice` in text and as a
      doctor-only `problems` state in `doctor --json`, documented beside the `link.State`
      spellings; the informational line for project and local installs; `status`'s line in
      text only. Closes: TestDoctorReportsLinksThePluginAlsoServes,
      TestDoctorDoesNotAskForLinksThePluginServes, TestDoctorExitsNonZeroOnADuplicate,
      TestDoctorJSONListsServedTwice, TestProjectScopeDoctorReadsTheGlobalRecord,
      TestDoctorNamesAProjectOrLocalPluginInstall,
      TestDoctorExitsZeroOnAProjectOrLocalPluginInstallAlone,
      TestDoctorNamesAnUnreadableRecord, TestStatusNamesThePlugin,
      TestStatusJSONOmitsThePluginLine.
- [ ] **The README's plugin route.** Install gains the two `/plugin` lines, the "one or the
      other" sentence and how to update; *From a checkout instead* moves to `docs/DESIGN.md`
      with a one-line pointer left, so the README stays at its 380 ceiling;
      `readme_test.go` gains its three tests. Closes: TestInstallNamesThePluginRoute,
      TestReadmeSaysPluginAndLinksAreAlternatives, TestInstallSectionIsStepsOnly,
      TestReadmeUpdatingNamesThePluginRoute, TestReadmeStaysShort.
- [ ] **Land.** Apply this contract onto `.agents/specs/distribution/spec.md` — `Governs:`
      gains `.claude-plugin/**`, the "second way in" boundary is amended, the prior decisions
      move into *Prior decisions* — and `spec-cli.md`, `spec-readme.md` onto theirs;
      in `cli` the standing exit-code paragraph replaced in place with `spec-cli.md`'s text
      and `served twice` added beside the `link.State` spellings; in `readme` Outcome 2, its
      criterion and the wording boundary replaced with `spec-readme.md`'s text;
      `docs/SPEC.md`'s distribution row; delete this folder; the six gates. Closes: every
      criterion above, re-run on the landing tree.
