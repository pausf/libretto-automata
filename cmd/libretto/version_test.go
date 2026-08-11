package main

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
)

// ldflags wins whenever it says anything but `dev`. It has to: `make build` stamps
// `git describe --tags --always --dirty`, and build info cannot know a tree is dirty or
// how many commits past the tag it is. Preferring build info would turn
// v0.2.0-3-gabc123-dirty back into a clean v0.2.0 the moment one existed.
func TestVersionPrefersLdflagsOverBuildInfo(t *testing.T) {
	info := &debug.BuildInfo{}
	info.Main.Version = "v0.9.9"

	for _, stamped := range []string{"v0.2.0", "v0.2.0-3-gabc123", "v0.2.0-dirty"} {
		if got := resolveVersion(stamped, info); got != stamped {
			t.Errorf("resolveVersion(%q, v0.9.9) = %q, want the stamped value", stamped, got)
		}
	}
}

// `go install pkg@version` produces no ldflags but does record the module version, so
// the binary can still say what it is instead of reporting `dev` to everyone who did not
// build from a clone.
func TestVersionFallsBackToBuildInfo(t *testing.T) {
	info := &debug.BuildInfo{}
	info.Main.Version = "v0.3.0"

	if got := resolveVersion("dev", info); got != "v0.3.0" {
		t.Errorf("resolveVersion(dev, v0.3.0) = %q, want v0.3.0", got)
	}
	if got := resolveVersion("", info); got != "v0.3.0" {
		t.Errorf("resolveVersion(empty, v0.3.0) = %q, want v0.3.0", got)
	}
}

// The notice names both versions and the action. A row saying "an update is available"
// cannot be checked against `libretto version`, and one with no action is a notification
// with nowhere to go.
func TestReleaseNoticeNamesBothVersionsAndTheAction(t *testing.T) {
	root := rootWithCachedTag(t, "v0.3.0")

	got := releaseNotice(root, "v0.2.0")
	for _, want := range []string{"v0.2.0", "v0.3.0", "update"} {
		if !strings.Contains(got, want) {
			t.Errorf("notice %q does not mention %q", got, want)
		}
	}
}

// Silence is the default. Up to date, ahead of the remote, and a version that cannot be
// parsed all produce no row — the last one because telling somebody running `dev` that they
// are out of date is a guess presented as a fact.
func TestReleaseNoticeIsSilentWhenThereIsNothingToSay(t *testing.T) {
	cases := map[string]struct{ cached, running string }{
		"up to date":          {"v0.3.0", "v0.3.0"},
		"ahead of the remote": {"v0.2.0", "v0.3.0"},
		"unidentifiable":      {"v0.3.0", "dev"},
		"a dirty build":       {"v0.3.0", "v0.2.0-3-gabc123-dirty"},
		"nothing cached":      {"", "v0.2.0"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := releaseNotice(rootWithCachedTag(t, c.cached), c.running); got != "" {
				t.Errorf("notice = %q, want silence", got)
			}
		})
	}
}

// rootWithCachedTag is a clone whose release check has already been answered, so the notice
// is built without a remote, a network or a subprocess.
func rootWithCachedTag(t *testing.T, tag string) string {
	t.Helper()
	root := t.TempDir()
	git := filepath.Join(root, ".git")
	if err := os.MkdirAll(git, 0o755); err != nil {
		t.Fatal(err)
	}
	// The cache format: a unix timestamp, a space, then the tag. Far enough in the future
	// that no TTL can expire it, so the check never reaches for a remote that is not there.
	writeFile(t, filepath.Join(git, "libretto-update-check"), "99999999999 "+tag+"\n")
	return root
}

// The notice names the command that will actually work, and the two modes need different ones.
// A row naming a command that refuses is worse than no row.
func TestReleaseNoticeNamesTheCommandForTheMode(t *testing.T) {
	checkout := rootWithCachedTag(t, "v0.3.0")
	if got := releaseNotice(checkout, "v0.2.0"); !strings.Contains(got, "update") {
		t.Errorf("in a checkout the notice says %q, want it to name `update`", got)
	}
	if got := releaseNotice(checkout, "v0.2.0"); strings.Contains(got, "upgrade") {
		t.Errorf("in a checkout the notice names `upgrade`, which would refuse: %q", got)
	}

	// An installed copy has no .git, so `upgradeCommand` is what decides — asserted directly,
	// because reaching releaseNotice there would need a live forge.
	installed := t.TempDir()
	if got := upgradeCommand(installed); got != "upgrade" {
		t.Errorf("outside a checkout the command is %q, want upgrade", got)
	}
	if got := upgradeCommand(checkout); got != "update" {
		t.Errorf("inside a checkout the command is %q, want update", got)
	}
}

// Neither source means `dev`, and `dev` is the honest answer: a binary that cannot prove
// its version says so instead of claiming one. `(devel)` is what the toolchain records
// for a build from a working tree, which is not a version either.
func TestVersionIsDevWhenNothingKnows(t *testing.T) {
	cases := map[string]*debug.BuildInfo{
		"no build info": nil,
		"empty":         {},
	}
	devel := &debug.BuildInfo{}
	devel.Main.Version = "(devel)"
	cases["(devel)"] = devel

	for name, info := range cases {
		t.Run(name, func(t *testing.T) {
			if got := resolveVersion("dev", info); got != "dev" {
				t.Errorf("resolveVersion(dev, %s) = %q, want dev", name, got)
			}
		})
	}
}
