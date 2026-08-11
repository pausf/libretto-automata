package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeUpgrade is an upgrade whose every outside effect is recorded instead of performed. The
// order of `steps` is the property most of these tests are about.
type fakeUpgrade struct {
	steps  []string
	latest string
	failAt string
}

func (f *fakeUpgrade) deps(running string) upgradeDeps {
	return upgradeDeps{
		running: running,
		base:    "/nowhere",
		latest: func(context.Context) (string, error) {
			f.steps = append(f.steps, "latest")
			if f.failAt == "latest" {
				return "", errors.New("no route to host")
			}
			return f.latest, nil
		},
		installPayload: func(_ context.Context, tag string) error {
			f.steps = append(f.steps, "payload:"+tag)
			if f.failAt == "payload" {
				return errors.New("checksum mismatch")
			}
			return nil
		},
		replaceBinary: func(_ context.Context, tag string) (string, error) {
			f.steps = append(f.steps, "binary:"+tag)
			if f.failAt == "binary" {
				return "", errors.New("permission denied")
			}
			return "", nil
		},
		relink: func() error {
			f.steps = append(f.steps, "relink")
			if f.failAt == "relink" {
				return errors.New("one item was not linked")
			}
			return nil
		},
	}
}

// The order is fixed and it is not cosmetic: a new binary reading an old payload is a state
// nobody can reason about, which is the same argument the update flow already makes about
// relinking before compiling.
func TestUpgradeActivatesThePayloadBeforeTheBinary(t *testing.T) {
	f := &fakeUpgrade{latest: "v0.4.0"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("v0.3.0")); err != nil {
		t.Fatal(err)
	}

	want := []string{"latest", "payload:v0.4.0", "binary:v0.4.0", "relink"}
	if got := strings.Join(f.steps, " "); got != strings.Join(want, " ") {
		t.Errorf("steps ran as %q, want %q", got, strings.Join(want, " "))
	}
}

// Relinking after the swap is not redundant. `current` moving means existing links keep
// resolving, but a version that *adds* an item leaves that item unlinked — which is the whole
// complaint `notify-users-of-new-updates` was queued for.
func TestUpgradeRelinksSoNewItemsAppear(t *testing.T) {
	f := &fakeUpgrade{latest: "v0.4.0"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("v0.3.0")); err != nil {
		t.Fatal(err)
	}
	if len(f.steps) == 0 || f.steps[len(f.steps)-1] != "relink" {
		t.Errorf("relink is not the last step: %v", f.steps)
	}
}

// The complaint that produced this whole change was the word `git pull` appearing in front of
// somebody who only wanted to use the tool. A promise about output is kept by asserting on
// output.
func TestUpdateFromAReleaseNeverMentionsGit(t *testing.T) {
	for _, failAt := range []string{"", "latest", "payload", "binary", "relink"} {
		f := &fakeUpgrade{latest: "v0.4.0", failAt: failAt}
		var out strings.Builder

		err := fromRelease(context.Background(), &out, f.deps("v0.3.0"))

		said := out.String()
		if err != nil {
			said += " " + err.Error()
		}
		if strings.Contains(strings.ToLower(said), "git") {
			t.Errorf("failAt=%q said %q, which mentions git", failAt, said)
		}
	}
}

// A failure names which step. "upgrade failed" with four possible causes is a message that
// sends the reader to the source.
func TestUpgradeReportsWhichStepFailed(t *testing.T) {
	for step, want := range map[string]string{
		"latest":  "could not find",
		"payload": "payload",
		"relink":  "link",
	} {
		f := &fakeUpgrade{latest: "v0.4.0", failAt: step}
		var out strings.Builder

		err := fromRelease(context.Background(), &out, f.deps("v0.3.0"))
		if err == nil {
			t.Errorf("failAt=%q succeeded", step)
			continue
		}
		if !strings.Contains(err.Error(), want) {
			t.Errorf("failAt=%q said %q, want it to mention %q", step, err, want)
		}
	}
}

// A failed payload step leaves the binary alone, which is what keeps the previous version
// active: `dist.Install` only swaps `current` once the extraction completed, so nothing
// downstream may run and make the machine inconsistent.
func TestAFailedUpgradeLeavesThePreviousVersionActive(t *testing.T) {
	f := &fakeUpgrade{latest: "v0.4.0", failAt: "payload"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("v0.3.0")); err == nil {
		t.Fatal("a failed payload step was reported as success")
	}
	for _, step := range f.steps {
		if strings.HasPrefix(step, "binary") || step == "relink" {
			t.Errorf("%q ran after the payload step failed: %v", step, f.steps)
		}
	}
}

// The payload upgrade still stands when the binary cannot be replaced. The same split the
// rebuild already makes for an unwritable destination, and for the same reason: rolling back
// a payload that installed correctly loses more than it saves.
func TestUpgradeSurvivesAnUnreplaceableBinary(t *testing.T) {
	f := &fakeUpgrade{latest: "v0.4.0", failAt: "binary"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("v0.3.0")); err != nil {
		t.Fatalf("an unreplaceable binary failed the whole upgrade: %v", err)
	}
	if !strings.Contains(out.String(), "unchanged") {
		t.Errorf("the report does not say the binary is unchanged: %q", out.String())
	}
	// And it still relinked: the payload is new, so the links have to be.
	if f.steps[len(f.steps)-1] != "relink" {
		t.Errorf("relink was skipped after the binary step failed: %v", f.steps)
	}
}

// Already on the newest release is a success that changes nothing and says so. Nothing is
// downloaded and nothing is relinked — a no-op that relinks is a no-op with a report.
func TestUpgradeOnTheNewestReleaseDoesNothing(t *testing.T) {
	f := &fakeUpgrade{latest: "v0.4.0"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("v0.4.0")); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(f.steps, " "); got != "latest" {
		t.Errorf("steps = %q, want only the check", got)
	}
	if !strings.Contains(out.String(), "up to date") {
		t.Errorf("the report does not say it is up to date: %q", out.String())
	}
}

// A machine with no payload yet — a fresh `go install` — upgrades rather than refusing.
// `libretto upgrade` is how the payload arrives the first time.
func TestUpgradeFromNothingInstalled(t *testing.T) {
	f := &fakeUpgrade{latest: "v0.4.0"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("dev")); err != nil {
		t.Fatalf("upgrading with nothing installed: %v", err)
	}
	if got := strings.Join(f.steps, " "); !strings.Contains(got, "payload:v0.4.0") {
		t.Errorf("steps = %q, want the payload installed", got)
	}
}

// **One command, two routes.** `update` in a checkout pulls the tree; anywhere else it takes
// the release route. There is no second command and no refusal to explain: for the person
// typing it the meaning is one thing, and how the tool got here decides the mechanism.
func TestUpdateTakesTheRouteThisInstallationCameBy(t *testing.T) {
	f := newFixture(t)

	// A checkout with no remote: `update` gets far enough to prove it took the git route —
	// the release route would never mention a remote at all.
	checkout := gitDir(t)
	writeFile(t, filepath.Join(checkout, "go.mod"), moduleLine+"\n")
	if err := os.MkdirAll(filepath.Join(checkout, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	err := update(checkout, f.global())
	if err != nil && strings.Contains(err.Error(), "release") {
		t.Errorf("a checkout took the release route: %v", err)
	}
}

// The menu row names the mechanism this machine will use, so the one row is honest about what
// pressing it does — without becoming two rows.
func TestTheUpdateRowNamesTheMechanism(t *testing.T) {
	if got := updateDesc(true); !strings.Contains(got, "pull") {
		t.Errorf("in a checkout the row says %q, want it to mention the pull", got)
	}
	if got := updateDesc(false); !strings.Contains(got, "release") {
		t.Errorf("installed, the row says %q, want it to mention the release", got)
	}
	if strings.Contains(updateDesc(false), "git") {
		t.Errorf("the installed row mentions git: %q", updateDesc(false))
	}
}

// `install` with no payload on disk explains itself instead of reporting an empty tree. A
// fresh `go install` has a binary and nothing else, and "nothing to link" would read as
// success.
func TestInstallWithNoPayloadPointsAtUpdate(t *testing.T) {
	missing := fmt.Sprintf("%s/never-installed", t.TempDir())
	t.Setenv(EnvRoot, missing)
	t.Setenv("CLAUDE_HOME", t.TempDir())

	err := run([]string{"install"})
	if err == nil {
		t.Fatal("install succeeded with no payload")
	}
	if !strings.Contains(err.Error(), "update") {
		t.Errorf("the error does not point at `update`: %v", err)
	}
}

// **The dangerous one.** With no payload every link in the target resolves to nothing, so a
// scan reports all of them `stale` — and `prune --yes` would remove every item the user has,
// doing exactly what it promises on a premise that is false. It has to refuse before scanning.
func TestPruneWithNoPayloadRefusesInsteadOfDeletingEverything(t *testing.T) {
	f := newFixture(t)
	t.Setenv(EnvRoot, filepath.Join(t.TempDir(), "never-installed"))

	// A link the user has, of the kind a false `stale` verdict would sweep away.
	planted := filepath.Join(f.Claude, "skills", "write-spec")
	f.link(t, filepath.Join(t.TempDir(), "write-spec"), planted)

	err := run([]string{"prune", "--yes"})
	if err == nil {
		t.Fatal("prune --yes ran with no payload on disk")
	}
	if !strings.Contains(err.Error(), "update") {
		t.Errorf("the refusal does not point at `update`: %v", err)
	}
	if _, lerr := os.Lstat(planted); lerr != nil {
		t.Errorf("prune removed a link while the payload was missing: %v", lerr)
	}
}

// `models` reads the target's agents directory, not the payload, so it still works on a
// machine with nothing installed — which is exactly when somebody might want to see what
// model an agent is on.
func TestModelsWorksWithNoPayload(t *testing.T) {
	newFixture(t)
	t.Setenv(EnvRoot, filepath.Join(t.TempDir(), "never-installed"))

	if err := run([]string{"models"}); err != nil {
		t.Errorf("models with no payload: %v", err)
	}
}
