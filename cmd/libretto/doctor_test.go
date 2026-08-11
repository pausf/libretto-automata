package main

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestDoctorReportsNewerRelease(t *testing.T) {
	line := releaseLine("v0.2.0", func(context.Context) (string, error) { return "v0.3.0", nil })

	for _, want := range []string{"v0.2.0", "v0.3.0", "update"} {
		if !strings.Contains(line, want) {
			t.Errorf("line %q does not mention %q", line, want)
		}
	}
}

func TestDoctorSaysUpToDate(t *testing.T) {
	line := releaseLine("v0.3.0", func(context.Context) (string, error) { return "v0.3.0", nil })

	if !strings.Contains(line, "up to date") {
		t.Errorf("line = %q, want it to say up to date", line)
	}
}

// A failed check is stated, never silent and never an error. `doctor` is the command that
// went looking, so "I could not find out" is one of its legitimate answers — and printing
// nothing would read as "you are up to date", which is a claim nobody verified.
func TestDoctorSaysSoWhenTheCheckFails(t *testing.T) {
	line := releaseLine("v0.2.0", func(context.Context) (string, error) {
		return "", errors.New("no route to host")
	})

	if !strings.Contains(line, "could not check") {
		t.Errorf("line = %q, want it to say the check failed", line)
	}
	if strings.Contains(line, "up to date") {
		t.Errorf("a failed check was reported as up to date: %q", line)
	}
}

// A remote with no releases to offer is neither a failure nor an update. Saying "up to
// date" would be defensible; claiming a version exists would not.
func TestDoctorHandlesARemoteWithNoReleases(t *testing.T) {
	line := releaseLine("v0.2.0", func(context.Context) (string, error) { return "", nil })

	if strings.Contains(line, "available") {
		t.Errorf("a remote with no releases produced an update line: %q", line)
	}
	if line == "" {
		t.Error("doctor said nothing about the release check at all")
	}
}

// A version that cannot be parsed is not out of date. The line has to say what it does
// know — the tag the remote has — without ranking a binary it cannot identify.
func TestDoctorDoesNotRankAnUnidentifiableBinary(t *testing.T) {
	line := releaseLine("dev", func(context.Context) (string, error) { return "v0.3.0", nil })

	if strings.Contains(line, "available") {
		t.Errorf("a `dev` binary was told it is out of date: %q", line)
	}
	if !strings.Contains(line, "v0.3.0") {
		t.Errorf("line %q does not name the tag the remote has", line)
	}
}

// `doctor` names the mode, because "up to date" means something different in each: a checkout
// three commits past a tag is current with its remote and behind the release. The *command* it
// names is always `update` — there is only one.
func TestDoctorNamesTheModeItIsIn(t *testing.T) {
	ask := func(context.Context) (string, error) { return "v0.3.0", nil }

	got := releaseLine("v0.2.0", ask)
	if !strings.Contains(got, "update") {
		t.Errorf("line = %q, want it to name update", got)
	}
	if strings.Contains(got, "upgrade") {
		t.Errorf("line = %q still names a command that does not exist", got)
	}
}
