package main

import (
	"runtime/debug"
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
