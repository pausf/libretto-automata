# Continuous integration — delta

Targets: ci

One `make` target that publishes a release. **No workflow, and that is the point.**

## Outcomes

```bash
git tag -a v0.4.0 -m "..."
make release                 # builds the payload tarball, attaches it to the release
```

- **`make release` runs on the author's machine, not in CI.** `AGENTS.md` says a tag is a
  human act, and this spec's own scope boundary says CI publishes nothing. A workflow that
  released on tag push would make both sentences false, and it would do it by turning a
  deliberate decision into an automation nobody re-read.
- **It refuses to run on anything but a clean tree at an annotated tag.** A release built
  from a dirty tree ships files that are in no commit, and nobody can ever reconstruct what
  went out.
- **The tarball is `payload-<tag>.tar.gz`, with a `.sha256` beside it.** Both are attached to
  the GitHub release; the checksum is what `distribution` verifies before extracting.
- **It contains `skills/`, `agents/`, `commands/` and nothing else.** Not the Go sources, not
  `docs/`, not `.agents/`, not `scripts/` — `scripts/check-payload` is repo-only tooling by
  `AGENTS.md`, and a tarball carrying the specs would ship the change folder of whatever was
  in flight.
- **One tarball for every platform.** The payload is markdown. A per-platform matrix would
  produce identical files under different names.
- **Paths inside the archive are relative and contain no `..`.** The extractor refuses those
  anyway; producing them would mean shipping an archive our own tool rejects, which is a
  thing to catch here rather than in a bug report.
- **It runs the six gates first and stops if any fails.** A release is the one moment where
  "the gates were green earlier" is not good enough.
- **It is idempotent against an existing release**: re-running replaces the assets rather
  than erroring or creating a second release for the same tag.

## Scope boundaries

**In:** one `release` target in the `Makefile`, its preconditions, the tarball's contents and
its checksum.

**Out:**

- **a release workflow.** Named as the thing being declined, not forgotten. *Ceiling:* the
  day releases are frequent enough that doing it by hand gets skipped, a workflow gated on a
  tag already pushed by a human is the answer — never one that creates the tag.
- **`goreleaser`.** No matrix to build, and the binary comes from `go install`. It would be a
  new tool for a job `tar` and `gh` already do.
- **publishing binaries.** `distribution` names this as out, with its reason.
- **release notes generation.** `gh release create --notes` takes what the author writes. A
  generated changelog from commit subjects is a changelog nobody reads twice.
- **signing the tarball.** `distribution` names this too, with its ceiling.
- **touching `gates.yml`.** The six gates are unchanged.

## Constraints

- **`gh` and `tar` only.** Both already required by the flow — `record-work` uses `gh` for
  requests. No new tool.
- **`gh` must be authenticated**, and the target says so and stops rather than half-running.
  No token is read from the environment, echoed, or written anywhere: the same rule
  `find-work` applies to `jira`.
- The tarball is built with **relative paths from the repository root**, so the archive has no
  leading directory to strip.
- The checksum is `shasum -a 256` output in its standard form, so a human can verify it by
  hand with the tool their machine already has.

## Prior decisions

- **The tag stays a human act.** Already recorded in this spec and reaffirmed rather than
  reopened: the alternative was a tag-push workflow, and the reason it loses is that the
  decision to release would then live in a file instead of in somebody's judgment.
- **The gates are not duplicated into the release target.** It invokes `make gates`, so there
  is still one list.

## Task breakdown

1. `make release`: the preconditions, `make gates`, the tarball, the checksum, `gh release
   create` / asset replace.
2. A test that the target's contents match what `distribution` expects to download.

## Verification criteria

```
Proof: cmd/libretto/gates_test.go TestReleaseTargetRunsTheGatesFirst
Proof: cmd/libretto/gates_test.go TestReleaseTarballCarriesOnlyThePayloadDirectories
Proof: cmd/libretto/gates_test.go TestReleaseAssetNamesMatchWhatDistributionFetches
```

The third is the one that earns its place. The producer and the consumer of that filename are
in different packages and in different languages — a `Makefile` and Go — so nothing but a test
holds them together, and a typo in either is a release that installs on nobody's machine.
