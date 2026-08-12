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

// clonedRemote is taggedRemote's sibling for the retraction tests: a real clone, so the
// tag objects exist locally and `git show <tag>:go.mod` has something to read.
//
// The distinction is the whole reason this helper exists rather than a flag on the other
// one. `ls-remote` answers about the remote; reading a retraction reads a blob, and a blob
// has to be in the local object store. A repository that merely has `origin` set — which
// is what taggedRemote builds — can name a tag it cannot open.
func clonedRemote(t *testing.T, goMod map[string]string, tags ...string) Shell {
	t.Helper()
	upstream := gitRepoWithACommit(t)
	for _, tag := range tags {
		write(t, upstream, "go.mod", goMod[tag])
		commit(t, upstream, "go.mod for "+tag)
		run(t, upstream, "tag", "-a", tag, "-m", tag)
	}

	local := t.TempDir()
	clone := filepath.Join(local, "clone")
	run(t, local, "clone", "-q", upstream, clone)
	return Shell{Root: clone}
}

const plainGoMod = "module example.invalid\n\ngo 1.26.5\n"

// The tombstone case, and the reason this exists: `v1.0.2` in this repository has no
// Release and exists only to retract itself and the two versions before it. It is a
// perfectly valid semver tag, so it sorted highest and got offered as an update — to a
// checkout only, because the module proxy honours `retract` and `git ls-remote` cannot.
func TestLatestTagSkipsATagThatRetractsItself(t *testing.T) {
	s := clonedRemote(t, map[string]string{
		"v0.7.0": plainGoMod,
		"v1.0.2": plainGoMod + `
retract (
	v1.0.0 // not a major; the tool's contract never changed.
	v1.0.1 // not a major; the tool's contract never changed.
	v1.0.2 // exists only to carry these retractions. Never a release.
)
`,
	}, "v0.7.0", "v1.0.2")

	got, err := s.LatestTag(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got != "v0.7.0" {
		t.Errorf("LatestTag = %q, want v0.7.0 — the tombstone was offered as an update", got)
	}
}

// A tag whose go.mod cannot be read is offered, not hidden.
//
// **This test asserted the opposite until the tombstone reached a user, and the reversal is
// deliberate rather than a weakened test.** It used to require an unreadable `go.mod` to be
// read as permission to offer the tag, on the reasoning that a tag too new to have been
// fetched is a release far more often than a tombstone.
//
// That reasoning holds in general and fails here permanently: this repository's tombstone is
// the highest tag by design and stays there, so "rare and transient" is certain and forever
// for any clone lacking that one blob. The old rule therefore guaranteed the one outcome it
// was written to prevent, and the spec criterion moved with the code.
func TestLatestTagDoesNotOfferATagWhoseGoModCannotBeRead(t *testing.T) {
	// taggedRemote, deliberately: `origin` is set and the tags are not local, so every
	// `git show` fails. That is a `--no-tags` or shallow clone, built out of the ordinary
	// harness.
	s := taggedRemote(t, "v0.7.0", "v1.0.2")

	got, err := s.LatestTag(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Errorf("LatestTag = %q, want \"\" — nothing here is verifiable, so there is "+
			"nothing to offer", got)
	}
}

func TestRetractedReadsEveryDirectiveForm(t *testing.T) {
	tests := []struct {
		name  string
		gomod string
		tag   string
		want  bool
	}{
		{"no retract block at all", plainGoMod, "v1.0.2", false},
		{"single directive naming it", "retract v1.0.2\n", "v1.0.2", true},
		{"single directive naming another", "retract v1.0.0\n", "v1.0.2", false},
		{"block naming it", "retract (\n\tv1.0.0\n\tv1.0.2\n)\n", "v1.0.2", true},
		{"block naming only others", "retract (\n\tv1.0.0\n\tv1.0.1\n)\n", "v1.0.2", false},
		{"block with a reason comment", "retract (\n\tv1.0.2 // tombstone\n)\n", "v1.0.2", true},
		{"a range covering it", "retract [v1.0.0, v1.0.2]\n", "v1.0.1", true},
		{"a range at its upper bound", "retract [v1.0.0, v1.0.2]\n", "v1.0.2", true},
		{"a range past it", "retract [v1.0.0, v1.0.2]\n", "v1.1.0", false},
		{"a range inside a block", "retract (\n\t[v1.0.0, v1.0.2]\n)\n", "v1.0.1", true},
		// The word appearing in a comment or a module path is not a directive.
		{"the word in a comment", "// retract v1.0.2 one day\n", "v1.0.2", false},
		{"a require line that merely contains it", "require retracted.dev/x v1.0.2\n", "v1.0.2", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := retracted(tt.gomod, tt.tag); got != tt.want {
				t.Errorf("retracted(%q, %q) = %v, want %v", tt.gomod, tt.tag, got, tt.want)
			}
		})
	}
}

func TestLatestTagPicksHighestPlainSemver(t *testing.T) {
	// Deliberately not in order, and v0.9.0 after v0.10.0: the whole point is that this
	// is not a string sort.
	//
	// clonedRemote rather than taggedRemote: since an unverifiable tag is skipped, a
	// fixture whose tag objects are absent answers "" for every ordering question and the
	// test would pass or fail for a reason unrelated to sorting.
	// The go.mod bodies differ only by a comment: clonedRemote commits between tags, and
	// an identical file is nothing to commit.
	s := clonedRemote(t, map[string]string{
		"v0.1.0":  plainGoMod + "// v0.1.0\n",
		"v0.10.0": plainGoMod + "// v0.10.0\n",
		"v0.9.0":  plainGoMod + "// v0.9.0\n",
		"v0.2.0":  plainGoMod + "// v0.2.0\n",
	}, "v0.1.0", "v0.10.0", "v0.9.0", "v0.2.0")

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
		got, err := cached(context.Background(), cachePath(root), time.Hour, nowStub, ask)
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

	if _, err := cached(context.Background(), cachePath(root), time.Hour, nowStub, ask); err != nil {
		t.Fatal(err)
	}
	// Same clock, zero TTL: the entry is already stale.
	if _, err := cached(context.Background(), cachePath(root), 0, nowStub, ask); err != nil {
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
		got, err := cached(context.Background(), cachePath(root), time.Hour, nowStub, ask)
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

// With nowhere to keep the cache there is still an answer. That is a machine with nothing
// installed yet, and it means "ask, do not cache" rather than "fail".
func TestCheckCacheWithoutAGitDirectoryStillAnswers(t *testing.T) {
	root := t.TempDir()

	calls := 0
	ask := func(context.Context) (string, error) {
		calls++
		return "v0.3.0", nil
	}

	for i := 0; i < 2; i++ {
		got, err := cached(context.Background(), cachePath(root), time.Hour, nowStub, ask)
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

// cachePath is where the checkout's answer is kept.
func cachePath(root string) string { return filepath.Join(root, ".git", checkFile) }

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

// gitRepoWithACommit is a clonable, taggable source: a repository with one commit in it. It
// moved here from clone_test.go when Clone was removed — the helper outlived the feature it
// was written for, because ls-remote needs a remote with history too.
func gitRepoWithACommit(t *testing.T) string {
	t.Helper()
	root := gitRepo(t)
	if err := os.WriteFile(filepath.Join(root, "marker"), []byte("payload\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(t, root, "add", "marker")
	run(t, root, "commit", "-q", "-m", "first")
	return root
}

// The hole the tombstone fell through, reversing a prior decision in the repo-sync spec.
//
// `taggedRemote` builds what a `--no-tags` or shallow clone looks like: `ls-remote` names
// every tag on the remote and the local object store holds none of their objects, so
// `git show <tag>:go.mod` fails on all of them.
//
// The old rule offered such a tag, reasoning that a tag too new to be fetched is a release
// far more often than a tombstone. True in general, and wrong here in a way that does not
// self-correct: **this repository's tombstone is the highest tag permanently, by design**,
// so the "rare and transient" case is certain and forever for anyone whose store lacks that
// object. Reproduced on a real clone before this test existed — `LatestTag` returned
// `v1.0.2`, which is the one version that must never be offered.
//
// The two errors are not symmetric. A missed notice is repaired by the next fetch, and
// `update` fetches. A notice naming a retracted version is wrong every time it prints, and
// it points at a number that can never be un-published.
func TestAnUnverifiableTagIsNotOffered(t *testing.T) {
	s := taggedRemote(t, "v0.9.0", "v1.0.2")

	got, err := s.LatestTag(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "v1.0.2" {
		t.Fatalf("offered the tombstone: a tag whose go.mod cannot be read must not be " +
			"offered as an update")
	}
	if got != "" {
		t.Fatalf("no tag here is verifiable, so there is nothing to offer; got %q", got)
	}
}
