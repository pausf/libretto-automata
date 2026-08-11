# ship-payload-as-release-asset

Tracker: none

## The ask, verbatim

> no me gusta que diga  # git pull osea lo suyo es que para atualizar haga como este
> https://github.com/Gentleman-Programming/gentle-ai , como se actualiza el codigo a la
> ultima version? como esta versionado para saber que se tiene que actualizar etc ?

## Reading

`libretto update` announcing `git pull` is wrong for somebody who only wants to *use* the
tool. The reference is `gentle-ai`, which was read rather than assumed:

- `internal/assets/assets.go` uses `embed.FS` — its assets ship **inside** the binary
- `gentle-ai upgrade` re-downloads the binary; `gentle-ai sync` then writes the assets out
- versioned with semver tags and `goreleaser`, publishing per-platform binaries
- the latest version comes from `GET api.github.com/repos/.../releases/latest`, `tag_name`
  scraped with `sed`
- `scripts/install.sh` prefers brew, then a release binary, then `go install`

Embedding is the one shape this project cannot take: the payload is delivered as
**symlinks**, and `ownership`, `link-state`, `linking` and `prune` are all built on real
files existing at the other end of one. An `embed.FS` has no path on disk, so embedding
turns every link into a copy and the ownership model goes with it.

So the third shape, which phase 2 of `distribute-via-go-install` never offered and should
have: **the payload is a release asset.** `goreleaser` publishes the binary and a payload
tarball per tag; `libretto upgrade` fetches the tarball for the newer version, verifies it,
extracts it under `~/.local/share/libretto/<version>/`, and points the links there. No git
in the user-facing path, no embed, and symlinks keep pointing at real files.

## What this replaces

The clone bootstrap that `distribute-via-go-install` just landed. It is safe to replace
rather than migrate: **no tag exists**, so `@latest` cannot resolve to it and nobody has
ever installed it.

What that change built and this one keeps: the semver comparison, `LatestTag`, the check
cache, the panel's notice row and its `Init` seam, the version fallback through build info,
`repoRoot`'s four rungs, and the rebuild that replaces the running binary.

## Named at the outset, because neither gets scoped out

- **Extraction is a trust boundary.** A `../../.ssh/authorized_keys` entry inside the
  tarball writes outside the destination. Every path is validated against the destination
  before anything is written, and absolute paths and symlinks escaping the tree are
  refused.
- **The checksum is verified before extraction, not after.** Downloading from the network
  and then linking it into `~/.claude` unverified is the step with no later remedy.
