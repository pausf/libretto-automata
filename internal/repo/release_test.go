package repo

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// taggedRemote is a repository with `origin` pointing at another repository carrying tags.
// A local path is a perfectly good remote for ls-remote, so this needs no network.
func taggedRemote(t *testing.T, tags ...string) Shell {
	t.Helper()
	upstream := gitRepoWithACommit(t)
	for _, tag := range tags {
		// Annotated, because AGENTS.md says releases are `git tag -a`, and an annotated
		// tag is exactly what makes ls-remote emit a second `^{}` line per tag. A
		// lightweight tag here would leave the peel-stripping untested.
		run(t, upstream, "tag", "-a", tag, "-m", tag)
	}

	local := gitRepo(t)
	run(t, local, "remote", "add", "origin", upstream)
	return Shell{Root: local}
}

func TestLatestTagPicksHighestPlainSemver(t *testing.T) {
	// Deliberately not in order, and v0.9.0 after v0.10.0: the whole point is that this
	// is not a string sort.
	s := taggedRemote(t, "v0.1.0", "v0.10.0", "v0.9.0", "v0.2.0")

	got, err := s.LatestTag(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got != "v0.10.0" {
		t.Errorf("LatestTag = %q, want v0.10.0", got)
	}
}

// A remote whose only tags are prereleases or not versions at all has no release to
// offer, and that is an empty answer rather than an error. Nothing is wrong.
func TestLatestTagIgnoresPrereleaseAndNonSemverTags(t *testing.T) {
	s := taggedRemote(t, "v1.0.0-rc.1", "nightly", "v2.0.0-beta", "release-candidate")

	got, err := s.LatestTag(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Errorf("LatestTag = %q, want no release", got)
	}
}

// A remote with no tags at all is the same state: nothing to offer, nothing wrong.
func TestLatestTagOnARemoteWithNoTags(t *testing.T) {
	s := taggedRemote(t)

	got, err := s.LatestTag(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Errorf("LatestTag = %q, want empty", got)
	}
}

// The call cannot hang. A network that accepts a connection and never answers would
// otherwise block the panel's first paint, and the user's only recourse is ⌃C on a tool
// that looks broken.
func TestLatestTagHonoursDeadline(t *testing.T) {
	s := taggedRemote(t, "v0.1.0")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := s.LatestTag(ctx); err == nil {
		t.Error("LatestTag ignored a cancelled context")
	}
}

// A repository with no remote cannot be asked, and saying so is not an error the user
// needs to see — `update` already treats a missing remote as a state.
func TestLatestTagWithNoRemote(t *testing.T) {
	s := Shell{Root: gitRepo(t)}

	if _, err := s.LatestTag(context.Background()); err == nil {
		t.Error("LatestTag succeeded with no remote configured")
	}
}

// The cache is what keeps a network call off every launch, and it caches failure too: a
// machine with no network would otherwise pay the timeout every single time.
func TestCheckCacheSuppressesCallsInsideTTL(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	calls := 0
	ask := func(context.Context) (string, error) {
		calls++
		return "v0.3.0", nil
	}

	for i := 0; i < 3; i++ {
		got, err := checkedLatest(context.Background(), root, time.Hour, nowStub, ask)
		if err != nil {
			t.Fatal(err)
		}
		if got != "v0.3.0" {
			t.Errorf("call %d returned %q, want v0.3.0", i, got)
		}
	}
	if calls != 1 {
		t.Errorf("asked the remote %d times inside the TTL, want 1", calls)
	}
}

func TestCheckCacheAsksAgainOnceTheTTLExpires(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	calls := 0
	ask := func(context.Context) (string, error) {
		calls++
		return "v0.3.0", nil
	}

	if _, err := checkedLatest(context.Background(), root, time.Hour, nowStub, ask); err != nil {
		t.Fatal(err)
	}
	// Same clock, zero TTL: the entry is already stale.
	if _, err := checkedLatest(context.Background(), root, 0, nowStub, ask); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("asked %d times across an expired TTL, want 2", calls)
	}
}

// An offline machine must not pay the timeout on every launch. The failure is recorded as
// an answer — an empty one — so the next launch inside the TTL is silent and instant.
func TestCheckCacheRecordsFailureSoOfflineDoesNotRetry(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	calls := 0
	ask := func(context.Context) (string, error) {
		calls++
		return "", errors.New("no route to host")
	}

	for i := 0; i < 3; i++ {
		got, err := checkedLatest(context.Background(), root, time.Hour, nowStub, ask)
		if err != nil {
			t.Errorf("call %d surfaced the failure to the caller: %v", i, err)
		}
		if got != "" {
			t.Errorf("call %d returned %q, want empty", i, got)
		}
	}
	if calls != 1 {
		t.Errorf("retried %d times while offline, want 1", calls)
	}
}

// Without a .git directory there is nowhere to keep the cache. That is the bootstrap
// case, and it means "ask, do not cache" rather than "fail".
func TestCheckCacheWithoutAGitDirectoryStillAnswers(t *testing.T) {
	root := t.TempDir()

	calls := 0
	ask := func(context.Context) (string, error) {
		calls++
		return "v0.3.0", nil
	}

	for i := 0; i < 2; i++ {
		got, err := checkedLatest(context.Background(), root, time.Hour, nowStub, ask)
		if err != nil {
			t.Fatal(err)
		}
		if got != "v0.3.0" {
			t.Errorf("call %d returned %q", i, got)
		}
	}
	if calls != 2 {
		t.Errorf("asked %d times with nowhere to cache, want 2", calls)
	}
}

// nowStub is a fixed clock. Date.now in a test is a test that behaves differently at
// midnight.
func nowStub() time.Time { return time.Unix(1_770_000_000, 0) }

func TestNewerComparesFieldsNumerically(t *testing.T) {
	cases := []struct {
		name           string
		latest, runing string
		want           bool
	}{
		{"a patch ahead", "v0.2.1", "v0.2.0", true},
		{"the same version", "v0.2.0", "v0.2.0", false},
		{"behind", "v0.2.0", "v0.2.1", false},
		// String comparison gets this one backwards, and it is two releases away.
		{"ten beats nine", "v0.10.0", "v0.9.0", true},
		{"nine does not beat ten", "v0.9.0", "v0.10.0", false},
		{"major wins over minor", "v1.0.0", "v0.99.0", true},
		{"minor wins over patch", "v0.3.0", "v0.2.99", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsNewer(c.latest, c.runing); got != c.want {
				t.Errorf("IsNewer(%q, %q) = %v, want %v", c.latest, c.runing, got, c.want)
			}
		})
	}
}

func TestParseSemverAcceptsOnlyPlainReleases(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"v0.2.0", true},
		{"v10.20.30", true},
		{"0.2.0", false},            // the v is part of how this repo tags
		{"v0.2", false},             // three fields or nothing
		{"v0.2.0.1", false},         // four is not semver either
		{"v1.0.0-rc.1", false},      // a prerelease is not offered
		{"v0.2.0-3-gabc123", false}, // git describe output, not a release
		{"v0.2.0-dirty", false},     // nor is a dirty build
		{"dev", false},              // the honest no-version
		{"", false},                 // and nothing at all
		{"v0.2.0+meta", false},      // build metadata is not ranked
		{"vx.y.z", false},           // digits, or it is not a version
		{"v-1.2.3", false},          // no negatives: the minus is not a digit
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			if _, ok := parseSemver(c.in); ok != c.want {
				t.Errorf("parseSemver(%q) ok = %v, want %v", c.in, ok, c.want)
			}
		})
	}
}

// A binary that cannot say what version it is must never be told it is out of date.
// `dev` and a git-describe string are both common — `dev` is what a plain `go build`
// produces, and every commit past a tag produces the other.
func TestNewerIsFalseForUnparseableRunningVersion(t *testing.T) {
	for _, running := range []string{"dev", "", "v0.2.0-3-gabc123", "v0.2.0-dirty"} {
		if IsNewer("v9.9.9", running) {
			t.Errorf("IsNewer(v9.9.9, %q) = true, want false", running)
		}
	}
}

func TestNewerIsFalseForUnparseableLatest(t *testing.T) {
	for _, latest := range []string{"", "dev", "v1.0.0-rc.1"} {
		if IsNewer(latest, "v0.2.0") {
			t.Errorf("IsNewer(%q, v0.2.0) = true, want false", latest)
		}
	}
}
