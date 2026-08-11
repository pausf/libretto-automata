package repo

import "testing"

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
