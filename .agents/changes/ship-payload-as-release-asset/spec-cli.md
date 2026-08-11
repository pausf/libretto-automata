# CLI — delta

Targets: cli

`upgrade` arrives, the clone bootstrap goes, and the two roots the tool has always
conflated get separated.

## Outcomes

### Two roots, not one

The tool has always used one `root` for two unrelated things, and it got away with it
because a clone happened to be both:

| The **payload root** | where `skills/`, `agents/`, `commands/` are — what links point at |
| The **checkout** | a git repository, needed only by `update` and only in development |

`go install` breaks the coincidence: there is a payload root and no checkout.

**The payload root resolves in four rungs:**

| 1 | `LIBRETTO_ROOT` | absolute override, taken as given |
| 2 | the compile-time source directory | when it has `.git` — development |
| 3 | the working directory | when it has `.git` **and** a `go.mod` naming this module |
| 4 | `~/.local/share/libretto/current` | the activated release |

Rungs 1–3 are unchanged and keep their reasoning: **a clone you are standing in still
wins**, which is what keeps editing a skill and seeing it live working. Only rung 4
changes — from `~/.libretto-automata` to the activated release directory.

- **`~/.libretto-automata` is gone, and nothing migrates.** It shipped in no tag, so
  `@latest` could never resolve to it and no machine has one. A migration path for a
  directory that never existed for anybody is code with no caller.
- **Nothing is cloned, ever, by any command.** The bootstrap that announced a destination
  and ran `git clone` is removed outright, along with its timeout.

### `upgrade` replaces `update` for installed users

| `upgrade` | fetch the newest release's payload, activate it, replace the binary, relink |
| `update` | **development only** — `git pull`, rebuild, relink, exactly as it does today |

- **`upgrade` never mentions git**, because it never uses it. That is the whole point of
  this change: a person who typed `go install` should never be told about a pull.
- **`update` outside a checkout says so and points at `upgrade`.** Not a silent alias:
  the two do different things, and a command that quietly becomes another is a command
  whose output nobody can predict. Symmetrically, **`upgrade` inside a checkout refuses
  and points at `update`** — overwriting a developer's working tree with a release tarball
  is the one thing it must never do.
- **The order is fixed: activate the payload, then the binary, then relink.** A binary from
  the new version reading a payload from the old is a state nobody can reason about — the
  same argument the update flow already makes about relinking before compiling.
- **It relinks afterwards, and that is not redundant.** `current` swapping means existing
  links keep resolving, but a version that *adds* an item leaves that item unlinked. This
  is exactly the complaint `notify-users-of-new-updates` was queued for, answered here.
- **Every step is reported with its outcome**, and a failure names which step. "upgrade
  failed" with four possible causes is a message that sends the reader to the source.
- **A failed step leaves the previous version active.** `current` is only swapped once the
  new payload is fully extracted, so an interrupted upgrade is a no-op rather than a
  half-installed tree.
- **The binary is replaced by re-running `go install <module>/cmd/libretto@<tag>`.** Go is
  present by construction — it is how the command got installed. Publishing per-platform
  binaries would remove that assumption and is deliberately out of scope; it is named in
  `distribution`.
- **If the binary cannot be replaced, the payload upgrade still stands** and the report says
  the binary is unchanged — the same split the rebuild already makes for an unwritable
  destination, and for the same reason.

### Knowing there is something to upgrade to

- **Which source is asked depends on the mode, and each mode has exactly one.** A checkout
  asks its remote with `git ls-remote --tags`; an installed copy reads the release redirect.
  Two implementations of "the latest version" would be two answers that can disagree, but
  these are two *questions* — "what has my remote tagged" and "what has the project
  released" — and only one is ever asked on a given machine.
- **`doctor` names the mode it is in**, so the answer is interpretable. "up to date" means
  something different in a checkout three commits past a tag.
- **The panel's notice says `choose upgrade`** where the command is `upgrade`, and
  `choose update` in a checkout. A row naming a command that refuses is worse than no row.

## Scope boundaries

**In:** `upgrade`, the payload-root rungs, separating the payload root from the checkout,
removing bootstrap, the mode-dependent release source, and the usage and README text.

**Out:**

- **downloading, verifying, extracting, activating.** `distribution` owns all four.
- **`install.sh`, brew, a tap, per-platform binaries.** Still out, still named so it is a
  decision. `go install` is the one entry point.
- **migrating `~/.libretto-automata`.** It never shipped.
- **`upgrade --to <version>`.** Rollback is a symlink swap that `distribution` supports, but
  a flag for it is a flag with no reported need. *Ceiling:* the first bad release that has
  to be backed out.
- **auto-upgrade, and prompting to upgrade.** The notice is a row; pressing `update` is the
  user's move.
- **`sync` as a separate command.** `gentle-ai` needs one because its assets live inside the
  binary and have to be written out. Here the payload is already files, and relinking is
  what `install` does — a second command for it would be a synonym.

## Constraints

- `CLAUDE_HOME` still governs targets. **`LIBRETTO_ROOT` still governs the payload root**, and
  a third variable for the release directory is not added — the override already covers
  pointing the tool somewhere else.
- **Every test that could reach `~/.local/share` sets `LIBRETTO_ROOT` or the injected base.**
  Same rule as `CLAUDE_HOME`, extended to the new directory.
- No test performs a real `go install` or a real download.
- The mode check is `isRepo(root)` — the same one-line probe the rungs already use, not a
  second notion of what development means.

## Prior decisions

- **A release asset, not an embedded payload.** See `distribution`. This is the second time
  the question has been settled and the reasoning has not changed: symlinks need real files.
- **`upgrade` and `update` are two commands, not one that guesses.** The alternative — one
  `update` that pulls in a checkout and downloads otherwise — was rejected: a command whose
  mechanism depends on invisible state is a command whose failure modes nobody can predict.
- **The clone bootstrap was the wrong shape, and it is removed rather than kept alongside.**
  Two ways to get a payload is two code paths, and the one nobody runs is the one that
  breaks. It survived one session and no tag.

## Task breakdown

1. Rename the resolution to a payload root; rung 4 becomes the activated release.
2. Remove `bootstrap.go`, `ensureClone`, the bootstrap timeout and their tests.
3. `upgrade`: compose `distribution` and the existing rebuild, in the fixed order, with a
   per-step report.
4. `update` outside a checkout points at `upgrade`; `upgrade` inside one refuses.
5. The release source per mode, and `doctor` naming the mode.
6. The panel notice naming the right command.
7. Usage and `README.md`.

## Verification criteria

```
Proof: cmd/libretto/root_test.go TestPayloadRootFallsBackToTheActivatedRelease
Proof: cmd/libretto/root_test.go TestPayloadRootStillPrefersACheckoutYouAreStandingIn
Proof: cmd/libretto/upgrade_test.go TestUpgradeActivatesThePayloadBeforeTheBinary
Proof: cmd/libretto/upgrade_test.go TestUpgradeRelinksSoNewItemsAppear
Proof: cmd/libretto/upgrade_test.go TestUpgradeRefusesInsideACheckout
Proof: cmd/libretto/upgrade_test.go TestUpgradeReportsWhichStepFailed
Proof: cmd/libretto/upgrade_test.go TestAFailedUpgradeLeavesThePreviousVersionActive
Proof: cmd/libretto/upgrade_test.go TestUpgradeSurvivesAnUnreplaceableBinary
Proof: cmd/libretto/upgrade_test.go TestUpdateOutsideACheckoutPointsAtUpgrade
Proof: cmd/libretto/upgrade_test.go TestUpgradeNeverMentionsGit
Proof: cmd/libretto/doctor_test.go TestDoctorNamesTheModeItIsIn
Proof: cmd/libretto/version_test.go TestReleaseNoticeNamesTheCommandForTheMode
```

`TestUpgradeNeverMentionsGit` is not a joke criterion. The complaint that produced this
change was the word `git pull` appearing in front of somebody who only wanted to use the
tool, and a promise about output is kept by asserting on output.
