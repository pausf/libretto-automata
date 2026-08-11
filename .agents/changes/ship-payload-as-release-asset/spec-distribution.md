# Distribution — new capability

Targets: distribution *(new)*

Governs: internal/dist/**

Getting a versioned payload onto the machine from a published release, and making it the
one the links point at.

This is a **new capability**, not an amendment. `repo-sync` answers questions about a git
checkout; nothing here involves git. Folding a downloader into it would give one package
two unrelated reasons to exist and one spec two subjects.

## Outcomes

### The payload has a versioned home

```
~/.local/share/libretto/
├── v0.3.0/          skills/ agents/ commands/
├── v0.4.0/          skills/ agents/ commands/
└── current  →  v0.4.0
```

- **`current` is a symlink, and it is what `~/.claude` links through.** So
  `~/.claude/skills/write-spec` → `…/libretto/current/skills/write-spec` → `…/v0.4.0/skills/write-spec`.
  Real files at the end of the chain, which is what `ownership` and `linking` require and
  what embedding could never give.
- **Activating a version is one atomic symlink swap.** Written to a temporary name and
  renamed over `current`, so there is no instant where `current` is missing. A user whose
  `~/.claude` briefly pointed at nothing would see every skill vanish.
- **The previous version is kept, and only the previous one.** Rolling back is a swap, not
  a download. Two versions on disk is a few hundred kilobytes of markdown; keeping every
  version ever installed is a directory nobody prunes.
  *Ceiling:* if a payload ever ships something large, this becomes a retention setting.
- **This is not `~/.libretto-automata`.** That was the clone bootstrap, and it is removed —
  see the `cli` delta. `~/.local/share` is right here for the reason it was wrong there:
  this *is* opaque application data, not a checkout anybody edits.

### Which version is the latest

- **From the release redirect, not the API.** `https://github.com/<owner>/<repo>/releases/latest`
  answers `302` with `Location: …/releases/tag/vX.Y.Z`. The tag is read from that path.
- **No JSON, no token, no rate limit.** The unauthenticated GitHub API allows 60 requests
  an hour; a redirect costs nothing and needs no parsing. It is also one less thing to
  break when a response shape changes.
- **The redirect is not followed.** Only its `Location` is read. Following it fetches an
  HTML page nobody needs.
- **A non-302, a missing `Location`, or a target that is not a plain `vX.Y.Z` is "could not
  tell", never a guess.** `repo-sync`'s existing semver rules decide what a release is;
  this package does not get a second opinion.
- **It cannot hang**, and it is subject to the same one-a-day cache and the same
  cache-the-failure rule that `repo-sync` already established. An offline machine asks once.

### Fetching a version

- **The URL is deterministic**, so there is nothing to search:
  `…/releases/download/<tag>/payload-<tag>.tar.gz` and `…/payload-<tag>.tar.gz.sha256`.
- **The checksum is fetched and verified before a single byte is extracted.** Not after,
  and not skippable. Downloading from the network and then linking the result into
  `~/.claude` is the step with no later remedy, and "verify afterwards" means the damage is
  already done when the check fails.
- **The download goes to a temporary file, and extraction goes to a temporary directory**
  renamed into place only once it has completed. A half-extracted `v0.4.0/` that a later
  run treats as installed is a payload with holes in it.
- **A version already on disk is not re-downloaded.** `upgrade` twice is one download.

### Extraction is a trust boundary

**This is the one section here that is not about convenience.** A tarball is remote input,
and an extractor that trusts its entries writes wherever the archive says.

- **Every entry's resolved destination must be inside the destination directory.** Checked
  after joining and cleaning, so `../../.ssh/authorized_keys`, an absolute `/etc/passwd`,
  and `a/../../b` are all refused.
- **A refused entry aborts the whole extraction.** Not skipped with a warning: an archive
  containing one such entry is not an archive to take the rest of.
- **Symlinks inside the archive are refused outright.** The payload is files and
  directories; a link in it has no legitimate purpose and is the second half of every
  extraction escape — a link to `/` followed by a write through it.
- **Hard links, devices, FIFOs, setuid bits: refused.** The allowed set is regular files
  and directories, and it is a whitelist rather than a list of things to reject.
- **Sizes are bounded.** A total-bytes ceiling, so a small archive that decompresses to
  fill the disk fails instead of succeeding slowly.
- **Modes are normalised**, not taken from the archive: `0644` for files, `0755` for
  directories. The payload is markdown; nothing in it needs to be executable, and an
  executable bit arriving from a download is a bit nobody asked for.

## Scope boundaries

**In:** the versioned directory layout, the `current` swap, reading the latest tag from the
release redirect, downloading, checksum verification, extraction and its refusals, and
keeping one previous version.

**Out:**

- **replacing the binary.** `cli` owns that, and it already knows how — the mechanism
  landed with the rebuild that writes over `os.Executable()`.
- **git, in any form.** That is `repo-sync`, and the whole point of this capability is a
  path that does not need a checkout.
- **linking.** `linking` and `ownership` own it; this capability only makes real files exist
  at a stable path.
- **signature verification.** *Ceiling, stated rather than hidden:* the checksum proves the
  bytes match what the release publishes, and nothing more — a compromised release
  publishes a matching checksum. Signing (minisign, cosign, or a GPG detached signature) is
  the answer the day the threat model includes a compromised release, and it needs a key
  the project does not have yet.
- **per-platform assets.** The payload is prose. One tarball serves every platform, and a
  matrix would produce identical files under different names.
- **publishing.** `ci` owns `make release`, and by `AGENTS.md` the tag is a human act.
- **mirrors, a proxy, an alternate host.** One source, named in one place.
- **resuming a partial download.** A payload tarball is small; a failed download is
  retried whole. *Ceiling:* if it ever gets large enough for this to matter, ranged
  requests are the answer.

## Constraints

- **`net/http` from the standard library. No new dependency.** The ladder's second rung, and
  it holds: one GET for the redirect, two for the assets.
- **`archive/tar` and `compress/gzip`, also stdlib.** Shelling out to `tar` would move the
  refusals above into a flag soup that differs between GNU and BSD tar — and the BSD one is
  what macOS ships.
- **`GOOS`-independent.** Nothing here reads a platform.
- **Every path is a parameter or derived from one root**, so the tests run entirely inside
  `t.TempDir()`. No test writes to a real `~/.local/share`, the same rule `CLAUDE_HOME` and
  `LIBRETTO_ROOT` already carry.
- **No test reaches the network.** The HTTP client is injectable and the tests serve from
  `net/http/httptest`; the tarballs are built in-process by the test.
- **The owner and repository are named in exactly one place**, derived from the module path
  as `repo.ModuleURL` already does. A constant beside a module path is two declarations of
  one fact.

## Prior decisions

- **A release asset, not `embed.FS`.** Asked and answered. `gentle-ai` embeds and its
  `sync` writes the files out; libretto cannot, because the payload is delivered as
  symlinks and `ownership`, `link-state`, `linking` and `prune` all require a real file at
  the far end. An `embed.FS` has no path on disk.
- **The release redirect, not `api.github.com`.** No rate limit, no JSON, one less shape to
  break. Recorded because the obvious next instinct is to reach for the API.
- **No `goreleaser`.** The payload is platform-independent, so there is no matrix to build,
  and the binary already arrives through `go install`. `gh release create` from a `make`
  target keeps the tag a human act, which `ci`'s spec requires.
- **One tarball, not per-item downloads.** 32 requests to install is 32 ways to end up
  half-installed.

## Task breakdown

1. `internal/dist/paths.go` — the versioned layout, `current`, and the atomic swap.
2. `internal/dist/latest.go` — the release redirect, and the tag read from `Location`.
3. `internal/dist/fetch.go` — download to a temporary file, fetch and verify the checksum.
4. `internal/dist/extract.go` — the guarded extractor: the whitelist, the containment
   check, the size ceiling, mode normalisation.
5. `internal/dist/install.go` — compose: fetch, verify, extract to a temporary directory,
   rename, swap `current`, drop the version before last.

## Verification criteria

```
Proof: internal/dist/paths_test.go TestCurrentPointsAtTheActivatedVersion
Proof: internal/dist/paths_test.go TestActivateIsAtomicAndNeverLeavesCurrentMissing
Proof: internal/dist/paths_test.go TestOnlyThePreviousVersionIsKept
Proof: internal/dist/paths_test.go TestRollbackIsASwapNotADownload
Proof: internal/dist/latest_test.go TestLatestReadsTheTagFromTheRedirectLocation
Proof: internal/dist/latest_test.go TestLatestDoesNotFollowTheRedirect
Proof: internal/dist/latest_test.go TestLatestRejectsANonSemverLocation
Proof: internal/dist/latest_test.go TestLatestOnAnUnexpectedStatusIsCouldNotTell
Proof: internal/dist/latest_test.go TestLatestHonoursTheDeadline
Proof: internal/dist/fetch_test.go TestFetchVerifiesTheChecksumBeforeExtracting
Proof: internal/dist/fetch_test.go TestFetchRefusesAMismatchedChecksumAndKeepsNothing
Proof: internal/dist/fetch_test.go TestFetchDoesNotRedownloadAVersionAlreadyPresent
Proof: internal/dist/extract_test.go TestExtractRefusesAPathEscapingTheDestination
Proof: internal/dist/extract_test.go TestExtractRefusesAnAbsolutePath
Proof: internal/dist/extract_test.go TestExtractRefusesASymlinkEntry
Proof: internal/dist/extract_test.go TestExtractRefusesANonRegularEntry
Proof: internal/dist/extract_test.go TestOneRefusedEntryAbortsTheWholeExtraction
Proof: internal/dist/extract_test.go TestExtractStopsAtTheSizeCeiling
Proof: internal/dist/extract_test.go TestExtractNormalisesModesAndDropsTheExecutableBit
Proof: internal/dist/install_test.go TestInstallLeavesNoPartialVersionOnFailure
Proof: internal/dist/install_test.go TestInstallActivatesTheNewVersion
```

The five extraction refusals are the criteria that matter most here, and they are written
before the extractor exists on purpose: an extractor built first and guarded afterwards is
guarded against the cases its author happened to remember.
