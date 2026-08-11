# ship-payload-as-release-asset — plan

**This file is live state.** A box is marked the moment its task is verified, never
batched. The orchestrator is the only writer.

**Goal:** the payload arrives as a published release asset, so `upgrade` never mentions git.

**Architecture:** `internal/dist` fetches, verifies, extracts and activates a version under
`~/.local/share/libretto/<version>/`, with `current` as the symlink `~/.claude` links
through. `cli` composes it into `upgrade` and stops conflating the payload root with the
checkout. `ci` gains a `make release` that publishes the tarball. `panel` renders one more
menu row.

**Specs:** `spec-distribution.md` · `spec-cli.md` · `spec-ci.md` · `spec-panel.md`

## Global constraints

- Go 1.26.5. **No new dependency** — `net/http`, `archive/tar`, `compress/gzip`,
  `crypto/sha256`, all stdlib. A sixth direct dependency is on `AGENTS.md`'s ask-first list.
- All six gates before every commit: `make gates`.
- **No test reaches the network.** `net/http/httptest` serves; tarballs are built in-process.
- **No test writes to a real `~/.local/share`, `~/.claude` or `~/.libretto-automata`.**
  `CLAUDE_HOME`, `LIBRETTO_ROOT`, and an injected base directory.
- No test runs a real `go install`, and none runs `gh`.
- `ponytail:` comments in English.
- The extractor's refusals are **whitelist, not blacklist**. Regular files and directories;
  everything else refused.

## Files

| File | Responsibility |
|---|---|
| `internal/dist/paths.go` *(new)* | the versioned layout, `current`, the atomic swap, retention |
| `internal/dist/latest.go` *(new)* | the release redirect, and the tag from `Location` |
| `internal/dist/extract.go` *(new)* | the guarded extractor |
| `internal/dist/fetch.go` *(new)* | download to a temp file, fetch and verify the checksum |
| `internal/dist/install.go` *(new)* | compose: fetch → verify → extract → rename → swap |
| `cmd/libretto/root.go` | rung 4 becomes the activated release; payload root named as such |
| `cmd/libretto/bootstrap.go` | **deleted** |
| `cmd/libretto/upgrade.go` *(new)* | the `upgrade` command and its per-step report |
| `cmd/libretto/main.go` | dispatch, `update`'s redirect, `doctor`'s mode line |
| `internal/ui/panel.go` · `model.go` | nothing — the menu comes from `cli` |
| `Makefile` | the `release` target |
| `README.md` | `upgrade`, and the clone section as development-only |

---

## Can start now

**D1, D2, D3, R1** — four independent tasks, no shared file between any two.

**Start with D3.** It is the security-critical one, its refusals are specified, and every
other task is cheaper to reason about once the extractor is known to be honest.

---

### D3 · the guarded extractor

**Spec:** `spec-distribution.md` → *Extraction is a trust boundary*
**Files:** create `internal/dist/extract.go`, `internal/dist/extract_test.go`
**Blocked by:** nothing
**Produces:** `extract(r io.Reader, dest string, maxBytes int64) error`

**Refusals first, then the extractor.** An extractor written first and guarded afterwards is
guarded against the cases its author happened to remember.

- [x] failing tests, each building its own tarball in-process: an entry resolving outside
      `dest` (`../../x`, `a/../../b`), an absolute path (`/etc/passwd`), a symlink entry, a
      non-regular entry (hard link, FIFO, device), one refused entry aborting the whole
      extraction, the total-size ceiling, and modes normalised with the executable bit
      dropped
- [x] run them, watch them fail — `undefined: extract`
- [x] implement: whitelist `tar.TypeReg` and `tar.TypeDir`; containment by
      `filepath.Clean` + a separator-aware prefix check against the cleaned `dest`; a
      running byte total; `0644` / `0755` regardless of the header
- [x] `make gates`, commit — exit 0 *(D3 closed)*

**The containment test was proven to have teeth, not assumed to.** Made the check the naive
`strings.HasPrefix(path, root)` and watched `TestExtractRefusesASiblingSharingThePrefix`
fail on `/tmp/destevil` being accepted as inside `/tmp/dest`, then restored it. Worth the two
minutes: a security test that passes for an unrelated reason is the failure the review seam
caught last time, one level deeper.

Three tests beyond the plan: the happy path (a file whose parent has no directory entry still
lands), a sibling sharing the destination's prefix, and input that is not a gzipped tarball at
all — which is what a truncated download and an HTML error page both look like, and neither
may read as an empty success. `io.CopyN` rather than `io.Copy`, so a reader yielding more than
the header declared cannot write past the ceiling just checked.

**Closes:** `TestExtractRefusesAPathEscapingTheDestination` ·
`TestExtractRefusesAnAbsolutePath` · `TestExtractRefusesASymlinkEntry` ·
`TestExtractRefusesANonRegularEntry` · `TestOneRefusedEntryAbortsTheWholeExtraction` ·
`TestExtractStopsAtTheSizeCeiling` · `TestExtractNormalisesModesAndDropsTheExecutableBit`

**Watch for:** a prefix check on strings alone passes `/tmp/destevil` as inside `/tmp/dest`.
The check is on the cleaned path with a trailing separator, or `filepath.Rel` and a `..`
test.

---

### D1 · the versioned layout

**Spec:** `spec-distribution.md` → *The payload has a versioned home*
**Files:** create `internal/dist/paths.go`, `internal/dist/paths_test.go`
**Blocked by:** nothing
**Produces:** `Base(home string) string` · `VersionDir(base, tag string) string` ·
`Current(base) string` · `Activate(base, tag string) error` · `Prune(base, keep ...string) error`

`Base` takes the home directory rather than reading it, so no test can reach a real
`~/.local/share`.

- [x] failing tests: `current` resolves to the activated version; activating twice moves it;
      **`current` is never absent mid-swap** — asserted by activating over an existing
      `current` and checking it resolves throughout; only the previous version survives a
      prune; rolling back to the previous version is an `Activate` with no download
- [x] run them, watch them fail — `undefined: Activate`
- [x] implement: `os.Symlink` to a temporary name then `os.Rename` over `current`
- [x] `make gates`, commit — exit 0 *(D1 closed)*

Four beyond the plan, each closing a way to lose a payload: `Activate` refuses a version that
is not on disk (pointing `current` at nothing is the vanishing made permanent), `Prune` never
removes the **active** version whatever its keep-list says, `Versions` sorts numerically and
ignores `current` and stray files, and `Base` is under `~/.local/share`.

**`repo.IsRelease` was exported rather than a second semver parser written here** — the spec
required one implementation. The coupling (a package about downloads importing one about git,
compile-time only) is marked `ponytail:` with its upgrade path: if it ever costs anything the
semver functions move to their own package.

**Closes:** `TestCurrentPointsAtTheActivatedVersion` ·
`TestActivateIsAtomicAndNeverLeavesCurrentMissing` · `TestOnlyThePreviousVersionIsKept` ·
`TestRollbackIsASwapNotADownload`

**Watch for:** `os.Symlink` fails if the target name exists, which is why the swap is
symlink-to-temp-then-rename and not remove-then-symlink. Remove-then-symlink is the version
with a window where `~/.claude` points at nothing.

---

### D2 · the latest tag, from the redirect

**Spec:** `spec-distribution.md` → *Which version is the latest*
**Files:** create `internal/dist/latest.go`, `internal/dist/latest_test.go`
**Blocked by:** nothing
**Produces:** `Latest(ctx, client *http.Client, base string) (string, error)`

- [x] failing tests, served by `httptest`: a `302` with
      `Location: /releases/tag/v0.4.0` yields `v0.4.0`; the redirect is **not** followed —
      asserted by a handler that fails the test if the tag page is requested; a `Location`
      that is not a plain `vX.Y.Z` is refused; a `200` or a `404` is "could not tell"; a
      cancelled context stops it
- [x] run them, watch them fail — `undefined: Latest`
- [x] implement with a client whose `CheckRedirect` returns `http.ErrUseLastResponse`, and
      `repo.IsRelease` for what a release is — exported, never reimplemented
- [x] `make gates`, commit — exit 0 *(D2 closed)*

**Also proven to have teeth.** Removed the `CheckRedirect` override and watched *two* tests
fail with `the redirect was followed`. Without that assertion living in the handler, the
no-follow promise is untestable.

Three beyond the plan: an absolute `Location` (which is what GitHub actually sends — depending
on the bare-path form is depending on a detail nobody promised), a `302` with no `Location` at
all, and `Latest` handed a **default** client, because it must not depend on the caller having
configured redirects correctly. The client is shallow-copied so the override cannot mutate one
the caller shares.

**Closes:** `TestLatestReadsTheTagFromTheRedirectLocation` ·
`TestLatestDoesNotFollowTheRedirect` · `TestLatestRejectsANonSemverLocation` ·
`TestLatestOnAnUnexpectedStatusIsCouldNotTell` · `TestLatestHonoursTheDeadline`

**Watch for:** the default `http.Client` follows redirects, so without `CheckRedirect` this
test passes for the wrong reason — it would read the tag out of a followed response's URL
and nobody would notice until GitHub changed the page.

---

### R1 · `make release`

**Spec:** `spec-ci.md` → all outcomes
**Files:** modify `Makefile`
**Blocked by:** nothing

- [x] the target: refuse a dirty tree, refuse a commit that is not an annotated tag, run
      `make gates`, build `payload-<tag>.tar.gz` from `skills/ agents/ commands/` only,
      write the `.sha256`, `gh release create` or replace the assets if the release exists
- [x] verify by reading the produced tarball's entry list — **not by publishing anything**
- [x] `make gates`, commit

**Closes:** the two `gates_test.go` criteria land in R2; this task is the target itself.

---

### D4 · download and verify

**Spec:** `spec-distribution.md` → *Fetching a version*
**Files:** create `internal/dist/fetch.go`, `internal/dist/fetch_test.go`
**Blocked by:** **D1**
**Produces:** `assetNames(tag string) (tarball, checksum string)` ·
`fetch(ctx, client, base, tag, dest string) error`

`assetNames` is the seam R2 asserts against, so the `Makefile` and this package cannot drift.

- [x] failing tests: the checksum is fetched and compared **before** extraction is reached —
      asserted with a mismatched checksum and a spy that fails if extraction is attempted;
      a mismatch leaves nothing on disk; a version directory already present is not
      re-downloaded
- [x] run them, watch them fail
- [x] implement: `sha256` over the downloaded temp file, compared to the fetched digest
- [x] `make gates`, commit

**Closes:** `TestFetchVerifiesTheChecksumBeforeExtracting` ·
`TestFetchRefusesAMismatchedChecksumAndKeepsNothing` ·
`TestFetchDoesNotRedownloadAVersionAlreadyPresent`

---

### D5 · compose the install

**Spec:** `spec-distribution.md` → task 5
**Files:** create `internal/dist/install.go`, `internal/dist/install_test.go`
**Blocked by:** **D1, D3, D4**
**Produces:** `Install(ctx, client, base, tag string) error`

- [x] failing tests: a failure at any step leaves no partial version directory and `current`
      untouched; success activates the new version
- [x] run them, watch them fail
- [x] implement: extract to `<base>/.tmp-<tag>`, rename to `<base>/<tag>`, then `Activate`
- [x] `make gates`, commit

**Closes:** `TestInstallLeavesNoPartialVersionOnFailure` · `TestInstallActivatesTheNewVersion`

---

### C1 · the payload root

**Spec:** `spec-cli.md` → *Two roots, not one*
**Files:** modify `cmd/libretto/root.go`, `cmd/libretto/root_test.go`
**Blocked by:** **D1**

- [x] failing tests: nothing found resolves to the activated release, not
      `~/.libretto-automata`; a checkout you are standing in still wins
- [x] run them, watch them fail
- [x] rung 4 → `dist.Current(base)`; rename to say *payload root*; `BootstrapDir` goes
- [x] `make gates`, commit

**Closes:** `TestPayloadRootFallsBackToTheActivatedRelease` ·
`TestPayloadRootStillPrefersACheckoutYouAreStandingIn`

---

### C2 · delete the bootstrap

**Spec:** `spec-cli.md` → *Two roots, not one*, third bullet
**Files:** delete `cmd/libretto/bootstrap.go`, `cmd/libretto/bootstrap_test.go`; modify
`cmd/libretto/main.go`
**Blocked by:** **C1**

Deleting `bootstrap_test.go` deletes five passing tests and one that was written in the
review seam. **That is correct here and is not weakening a proof:** they prove behaviour
this change removes on purpose. `TestVersionAndHelpDoNotBootstrap` is the one to keep the
*intent* of — `version` and `help` must still not touch the payload root — so it moves to
`root_test.go` rather than dying with the file.

- [x] move `TestVersionAndHelpDoNotBootstrap`'s intent into `root_test.go`, renamed for what
      it now guards
- [x] delete the file and `ensureClone`; `run` resolves the payload root directly
- [x] confirm no reference to `repo.Clone` survives outside its own tests
- [x] `make gates`, commit

**Closes:** no new criterion. A deletion's proof is that the suite still passes with the
behaviour gone, and the moved test.

---

### C3 · `upgrade`

**Spec:** `spec-cli.md` → *`upgrade` replaces `update` for installed users*
**Files:** create `cmd/libretto/upgrade.go`, `cmd/libretto/upgrade_test.go`; modify
`cmd/libretto/main.go` (dispatch, usage)
**Blocked by:** **D2, D5, C1**

The order is fixed: **payload, then binary, then relink.** A new binary reading an old
payload is a state nobody can reason about.

- [x] failing tests: the payload is activated before the binary is touched; it relinks
      afterwards so a *new* item appears; a failure names which step; a failed step leaves
      the previous version active; an unreplaceable binary still leaves the payload upgraded;
      **the output never contains the word `git`**
- [x] run them, watch them fail
- [x] implement, with the binary step as `go install <module>/cmd/libretto@<tag>` behind an
      injected runner so no test installs anything
- [x] `make gates`, commit

**Closes:** `TestUpgradeActivatesThePayloadBeforeTheBinary` ·
`TestUpgradeRelinksSoNewItemsAppear` · `TestUpgradeReportsWhichStepFailed` ·
`TestAFailedUpgradeLeavesThePreviousVersionActive` ·
`TestUpgradeSurvivesAnUnreplaceableBinary` · `TestUpgradeNeverMentionsGit`

---

### C4 · the two commands point at each other

**Spec:** `spec-cli.md` → second bullet of *`upgrade` replaces `update`*
**Files:** modify `cmd/libretto/upgrade.go`, `cmd/libretto/main.go`,
`cmd/libretto/upgrade_test.go`
**Blocked by:** **C3**

- [x] failing tests: `upgrade` inside a checkout refuses and names `update`; `update`
      outside one points at `upgrade`
- [x] run them, watch them fail
- [x] implement, both on `isRepo(root)` — the probe the rungs already use
- [x] `make gates`, commit

**Closes:** `TestUpgradeRefusesInsideACheckout` · `TestUpdateOutsideACheckoutPointsAtUpgrade`

---

### C5 · the release source per mode

**Spec:** `spec-cli.md` → *Knowing there is something to upgrade to*
**Files:** modify `cmd/libretto/version.go`, `cmd/libretto/main.go`,
`cmd/libretto/version_test.go`, `cmd/libretto/doctor_test.go`
**Blocked by:** **D2, C1**

A checkout asks `repo.LatestTag`; an installed copy asks `dist.Latest`. **Two questions,
one asked per machine** — not two answers to one question.

- [x] failing tests: the notice names `upgrade` on an installed copy and `update` in a
      checkout; `doctor` names the mode it is in
- [x] run them, watch them fail
- [x] implement; the cache and TTL are `repo`'s existing ones, unchanged
- [x] `make gates`, commit

**Closes:** `TestReleaseNoticeNamesTheCommandForTheMode` · `TestDoctorNamesTheModeItIsIn`

---

### P1 · the second menu row

**Spec:** `spec-panel.md` → all outcomes
**Files:** modify `cmd/libretto/main.go` (the menu), `internal/ui/menu_test.go`
**Blocked by:** **C5**

Both rows always present; the inapplicable one disabled, which is already `panel`'s rule.

- [x] failing tests: both `upgrade` and `update` are offered; the one that does not apply is
      disabled rather than absent
- [x] run them, watch them fail
- [x] implement
- [x] `make gates`, commit

**Closes:** `TestBothUpgradeAndUpdateAreOffered` ·
`TestTheInapplicableActionIsDisabledNotHidden`

---

### R2 · hold the producer and the consumer together

**Spec:** `spec-ci.md` → verification criteria
**Files:** modify `cmd/libretto/gates_test.go`
**Blocked by:** **R1, D4**

The `Makefile` writes the asset names and Go reads them, in two languages. Nothing but a
test holds them together, and a typo in either is a release that installs on nobody's
machine.

- [x] failing tests: the release target runs the gates first; the tarball carries only the
      three payload directories; **the names in the `Makefile` match `assetNames`**
- [x] run them, watch them fail
- [x] implement by reading the `Makefile` text, as the existing gate tests already do
- [x] `make gates`, commit

**Closes:** `TestReleaseTargetRunsTheGatesFirst` ·
`TestReleaseTarballCarriesOnlyThePayloadDirectories` ·
`TestReleaseAssetNamesMatchWhatDistributionFetches`

---

### Z1 · say so

**Spec:** `spec-cli.md` → scope, *In*
**Files:** modify `cmd/libretto/main.go` (`usage`), `README.md`
**Blocked by:** **C4, C5, P1, R2**

- [x] usage: `upgrade`, and `update` marked as the checkout command
- [x] `README.md`: `upgrade` replaces `git pull` in the update section; the clone section
      stays as the development route; `~/.local/share/libretto` named; **the word `git` does
      not appear in the install-and-update path**
- [x] `AGENTS.md`: the release step — tag, then `make release`
- [x] `make gates` and `scripts/check-payload`, commit

**Closes:** no test. Documentation that agrees with the binary is checked by reading both,
and `TestUpgradeNeverMentionsGit` covers the output half.

---

## What the plan got wrong

**Three tasks merged, each for a reason.** C1+C2: separating them leaves a commit where
`ensureClone` would try to clone *into* the release directory. C3+C4: `upgrade` without the
checkout guard would install a release and relink `~/.claude` away from a developer's tree, in
silence. D4's tests share a file with D5's. In every case the intermediate commit was not just
untidy but actively wrong, which is a different thing.

**The cache had to be generalised, and that fixed a hang the plan never saw.** `repo`'s check
cache was keyed on `root/.git`, which an installed copy does not have — so "no `.git` means ask
without caching" meant an HTTP call on **every panel launch**, precisely the hang the cache
exists to prevent. `repo.Cached` now takes the path; one TTL, one failure policy, two homes.

**Two safety findings that were not in any spec.**

`prune --yes` with no payload would have deleted every link the user has. The root points at
`~/.local/share/libretto/current`, absent on a fresh machine, so every link scans as `stale` —
a destructive command doing exactly what it promises on a false premise. That is why the
missing-payload guard is a stop and not a warning, and `models` is exempt because it reads the
target's agents rather than the payload.

`dist.Prune` never removes the **active** version whatever its keep-list says. A caller with a
wrong list should lose a spare, not the running payload.

**`repo.Clone` and `repo.ModuleURL` were deleted.** Written earlier the same session, and left
with no caller by C2. Dead code whose tests pass proves nothing. `--anchors` then caught
thirteen citations in the landed `cli` spec — seven renamed tests and six deleted ones — which
is the gate doing exactly what it is for.

**Three tests were verified to have teeth rather than assumed to.** The containment check made
naive (`/tmp/destevil` accepted as inside `/tmp/dest`), `CheckRedirect` removed (two tests fail
with "the redirect was followed"), and the release asset name changed to `.tgz`. Each failed,
then was restored.

**`TestEveryMenuLabelDispatches` earned its keep.** It predates this change and caught that the
panel dispatches through its own switch, not through `run()` — so a new menu row without a case
there is a row that does nothing.

**Two bugs of mine caught by the existing suite:** `args[0]` read before checking `len`, and the
payload guard placed ahead of the no-TTY exit that has always been 2.

**Scope that arrived unstated and stayed:** `libretto install` with no payload explains itself
and points at `upgrade`, rather than reporting an empty tree — "nothing to link" is what a fully
linked machine also says.

## Landing

Phase 8: create `.agents/specs/distribution/spec.md` from its delta — a **new capability
directory**, not an amendment — apply the `cli`, `ci` and `panel` deltas, delete this change
folder, and confirm `--anchors` resolves everything. Also update `docs/SPEC.md`, which is an
index of the capabilities and will be one short.

**And the queue:** `notify-users-of-new-updates` is answered by C3's relink and the notice
row that already shipped. It should be dropped rather than left to be picked up as if
untouched — but that is the user's call at phase 8, not this plan's.
