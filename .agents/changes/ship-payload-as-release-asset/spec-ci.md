# Continuous integration — delta

Targets: ci

One `make` target that opens a Release page. **Nothing is attached to it, and no workflow runs
it.**

## Outcomes

```bash
git tag -a v0.5.0 -m "..."
git push origin v0.5.0
make release
```

- **The payload is not an asset, because it does not need to be.** It ships inside the Go module,
  so `go install <module>/cmd/libretto@latest` brings it down with the binary and checks it
  against the checksum database. The module proxy resolves `@latest` from **tags**, which is why
  installing works against a repository with no Releases at all.
- **So this target is for humans:** a Release page carrying the tag's own notes, so somebody can
  read what changed without reading the log. Skipping it costs nothing mechanical.
- **It runs on the author's machine, not in CI.** `AGENTS.md` says a tag is a human act and this
  spec says CI publishes nothing. A workflow releasing on tag push would make both false, and it
  would do it by turning a judgment into an automation nobody re-reads.
- **The push comes before the target, and the order is load-bearing.** `gh release create`
  creates the tag itself when it is not on the remote, at the default branch's HEAD — a second
  tag with your name on it pointing somewhere you did not choose. `--verify-tag` refuses instead.
- **A tag is not a Release.** A tag is a git ref; a Release is a GitHub object. `git push origin
  v0.5.0` creates the first and not the second, and this repository sat at four tags and zero
  Releases — verified, not assumed.
- **It refuses a dirty tree and a HEAD that is not at a tag**, and runs the six gates first. A
  release is the one moment where "the gates were green earlier" is not good enough.
- **Re-running is a no-op** when the Release already exists, rather than an error or a second one.

## Scope boundaries

**In:** one `release` target, its four preconditions, and the gates it runs first.

**Out:**

- **a release workflow.** Named as declined, not forgotten. *Ceiling:* the day releases are
  frequent enough that doing it by hand gets skipped, a workflow gated on a tag a human already
  pushed is the answer — never one that creates the tag.
- **`goreleaser`, cross-compiled binaries, a payload tarball, checksums.** An earlier draft of
  this delta specified all four. They are gone because `go install` needs none of them: the
  module carries the payload and the go command builds for the platform it is on.
- **release notes generation.** `--notes-from-tag` takes what the author wrote. A changelog
  generated from commit subjects is one nobody reads twice.
- **touching `gates.yml`.** The six gates are unchanged.

## Constraints

- **`gh` only**, already required by the flow for pull requests. No new tool, and no `tar` or
  `shasum` any more either.
- **`gh` must be authenticated**, and the target says so and stops rather than half-running. No
  token is read from the environment, echoed, or written anywhere.

## Prior decisions

- **The tag stays a human act.** Reaffirmed rather than reopened.
- **The gates are not duplicated into the target.** It invokes `make gates`, so there is one list.
- **No assets.** The tarball, the checksums and the four cross-compiled binaries were specified,
  and then the module cache made all of them unnecessary. Recorded so the next reader knows the
  absence is a decision.

## Task breakdown

1. `make release`: the four preconditions, `make gates`, `gh release create --verify-tag
   --notes-from-tag`, idempotent against an existing Release.

## Verification criteria

```
Proof: cmd/libretto/gates_test.go TestReleaseTargetRunsTheGatesFirst
Proof: cmd/libretto/gates_test.go TestReleaseVerifiesTheTagRatherThanCreatingIt
```

Two, where the asset-publishing draft had three. The one that went held the Makefile's asset
names against the names Go fetched — a test whose whole job was keeping two languages in
agreement about a filename, and there is no filename now.
