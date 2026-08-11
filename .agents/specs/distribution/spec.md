# Distribution

Governs: internal/dist/**

Finding the payload an installed copy links from, and asking for a newer one.

**It downloads nothing and verifies nothing, and that is the whole design.** The payload —
`skills/`, `agents/`, `commands/` — ships **inside the Go module**, so
`go install <module>/cmd/libretto@latest` already fetches it, already checks it against the
checksum database, and already puts it under a path with the version in it.

This capability exists to say where that is, and to ask which version is newest.

## How this spec got here

Three shapes were tried and two were wrong, so both are recorded — a rejected design with no
reason attached gets proposed again.

**Embedded (`//go:embed`).** Rejected first and still rejected: an `embed.FS` has no path on
disk, so every symlink would become a copy and `ownership`, `link-state`, `linking` and `prune`
would all lose the real file they depend on.

**A release tarball.** Built, tested, then deleted the same session. It downloaded an asset,
verified a sha256 and extracted it through a guarded extractor with five refusals — path
traversal, symlinks, non-regular entries, a size ceiling, mode normalisation. Every line of it
reimplemented, less well, something the go command already does. **`GOSUMDB` is a stronger
guarantee than a checksum published beside the file it describes**, because it is not the
publisher who vouches for it.

**The module.** What is specified below. The question that found it was the ladder's first —
*does this need to exist at all?* — asked against the module cache rather than against the
feature.

## Outcomes

### The payload lives in the module cache

```
$GOMODCACHE/github.com/pausf/libretto-automata@v0.5.0/
├── skills/
├── agents/
└── commands/
```

- **`Dir(cache, module, version)` is that path.** The version is *in* it, which is what makes an
  update a different directory rather than an overwrite of the one currently linked.
- **The entry is read-only**, and that is fine: links point at it and nothing ever writes into
  it. Verified against a real `go install` rather than assumed.
- **`ModCache()` follows the go command's own order** — `GOMODCACHE`, then `GOPATH/pkg/mod`,
  then `~/go/pkg/mod` — **from the environment, without shelling out to `go env`.** Every
  command that touches the payload needs this answer, so a subprocess here is a subprocess on
  the hot path, for a path three variables already determine.

### Which version is newest

- **The module proxy is asked**, not the GitHub API and not the releases redirect. It is the
  same source `go install @latest` resolves against, so the two cannot disagree about what
  *latest* means.
- **And it answers from tags.** That is why `go install` works against a repository with no
  published Releases at all — which was this repository's state at `v0.4.0`: four tags, zero
  Releases.
- **Only a plain release counts, and `repo.IsRelease` decides.** A module with no tags answers
  with a pseudo-version — `v0.0.0-<date>-<sha>` — and that is not something to offer anybody as
  an update. One notion of what a version is, across the project.
- **A non-200, an unreadable body and a missing `Version` are all "could not tell"**, never a
  guess. An empty body and an HTML page are what a captive portal and a truncated response look
  like.
- **It cannot hang**, and it is subject to the same one-a-day cache `repo` already owns. The
  panel's first paint never waits on it.
- **An uppercase module path is escaped** the way the proxy protocol requires — `!` before the
  lowered rune. Nothing here exercises it in production, because this project's path is all
  lowercase; it is written and tested because the day somebody forks to a capitalised account
  the failure would be a 404 with no explanation attached.

### Asking for a newer one

- **`Install(module, version)` runs `go install <module>/cmd/libretto@<version>`.** One step,
  bringing the binary and the payload down together, because they travel in the same module.
- **There is no window where one is new and the other is old.** The tarball design needed a
  fixed step order to avoid exactly that state; one module removed it rather than guarding it.
- **A failure is reported with the go command's own output.** `go install` explains itself well
  — a proxy that refused, a checksum mismatch, a toolchain too old — and replacing that with
  "install failed" throws away the only useful part.

## Scope boundaries

**In:** the module cache path, the go command's resolution order, the newest version from the
proxy, module-path escaping, and running `go install`.

**Out:**

- **downloading, extracting or verifying anything.** The go command does all three. This is the
  boundary that deleted most of this capability.
- **git.** That is `repo-sync`. Nothing here shells to git.
- **linking.** `linking` and `ownership` own it.
- **deciding when to update, or what to say about it.** `cli` decides; this answers.
- **per-platform binaries, an `install.sh`, brew, a tap, npm.** `go install` is the entry point,
  and it is one command that leaves nothing else to do. Each of these was considered and each
  adds a second way in.
- **publishing release assets.** Nothing is attached to a release, so there is nothing to
  publish. `ci` keeps a `release` target for notes a human reads.
- **a payload cache of our own, a `current` symlink, retention.** The module cache is already
  versioned, and the go command already prunes it if asked.

## Constraints

- **No new dependency.** `net/http`, `encoding/json`, `os/exec`, `path/filepath` — stdlib.
- **No test reaches the network**, and none runs `go install`. The proxy is served by
  `httptest`; the runner is a parameter.
- **No test writes outside `t.TempDir()`.** Nothing here has a real-machine path baked in:
  `ModCache` reads the environment and `modCache` takes its three inputs.
- The module path is written down in exactly one place, `cli`'s `moduleLine`, and derived from
  there.

## Prior decisions

- **The module, not an embedded payload and not a release asset.** All three above, with their
  reasons, so none of them gets proposed again without new information.
- **The proxy, not `api.github.com` and not the `/releases/latest` redirect.** Same source as
  `go install`, no rate limit, no token, and it works off tags.
- **`repo.IsRelease` is imported rather than reimplemented.** Marked `ponytail:` at its
  definition with the upgrade path: if a package about downloads importing one about git ever
  costs anything, the semver functions move out and both depend on that.

## The cost, named rather than discovered

**`go clean -modcache` removes the payload, and every link into it breaks at once.**
`libretto install` re-downloads the module and repairs them; `doctor` reports the state.

That is the price of not keeping our own copy, and it is the right price: the alternative was a
download-verify-extract pipeline whose failure modes we owned instead of the go command's.

**Go is required.** It always was — `go install` was already the only entry point — and it is
now load-bearing for updates too. Named because the audience for this tool is not necessarily a
Go developer, and that is a real limit rather than an oversight.

## Task breakdown

1. `Dir` and `ModCache` — the module cache entry, and the go command's resolution order.
2. `Latest` — the proxy, the release check, the escaping.
3. `Install` — `go install`, with the runner as a parameter.

## Verification criteria

```
Proof: internal/dist/dist_test.go TestDirIsTheModuleCacheEntryForAVersion
Proof: internal/dist/dist_test.go TestModCacheFollowsTheGoCommandsOrder
Proof: internal/dist/dist_test.go TestLatestReadsTheVersionFromTheModuleProxy
Proof: internal/dist/dist_test.go TestLatestRefusesAPseudoVersion
Proof: internal/dist/dist_test.go TestLatestOnANonOKStatus
Proof: internal/dist/dist_test.go TestLatestRefusesAnUnreadableAnswer
Proof: internal/dist/dist_test.go TestLatestHonoursTheDeadline
Proof: internal/dist/dist_test.go TestLatestEscapesUppercaseInTheModulePath
Proof: internal/dist/dist_test.go TestInstallRunsGoInstallAtTheVersion
Proof: internal/dist/dist_test.go TestInstallReportsAFailure
```

Ten criteria where the tarball design had twenty-one. The five that mattered most there — the
extractor's refusals — are gone because the thing they guarded is gone, and that is the
difference between deleting a test and deleting a risk.
