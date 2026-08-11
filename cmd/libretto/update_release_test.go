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

// fakeUpdate is an update whose every outside effect is recorded instead of performed.
//
// Three effects, where the previous design had four: the binary and the payload arrive in one
// `go install`, because they travel in the same module. The step that had to be ordered against
// its neighbour no longer exists.
type fakeUpdate struct {
	steps  []string
	latest string
	failAt string
}

func (f *fakeUpdate) deps(running string) updateDeps {
	return updateDeps{
		running: running,
		latest: func(context.Context) (string, error) {
			f.steps = append(f.steps, "latest")
			if f.failAt == "latest" {
				return "", errors.New("no route to host")
			}
			return f.latest, nil
		},
		install: func(_ context.Context, v string) error {
			f.steps = append(f.steps, "install:"+v)
			if f.failAt == "install" {
				return errors.New("checksum mismatch for module")
			}
			return nil
		},
		dirFor: func(v string) string { return "/modcache/libretto-automata@" + v },
		relink: func(root string) error {
			f.steps = append(f.steps, "relink:"+root)
			if f.failAt == "relink" {
				return errors.New("one item was not linked")
			}
			return nil
		},
	}
}

// One install, then a relink of the version that install brought down.
//
// There is no longer an ordering constraint to protect. The old design fetched a payload tarball
// and replaced the binary as two steps, so "payload before binary" had to be enforced and tested
// — a new binary reading an old payload is a state nobody can reason about. Shipping both in one
// module removed the window rather than guarding it.
func TestUpdateInstallsThenRelinksTheVersionItInstalled(t *testing.T) {
	f := &fakeUpdate{latest: "v0.4.0"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("v0.3.0")); err != nil {
		t.Fatal(err)
	}

	want := "latest install:v0.4.0 relink:/modcache/libretto-automata@v0.4.0"
	if got := strings.Join(f.steps, " "); got != want {
		t.Errorf("steps ran as %q, want %q", got, want)
	}
}

// The relink uses the NEW version's directory, not payloadRoot(). This process still reports the
// old version, so resolving again would link the payload it is already on — and the update would
// report success having changed no link at all.
func TestUpdateRelinksTheNewVersionNotTheRunningOne(t *testing.T) {
	f := &fakeUpdate{latest: "v0.4.0"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("v0.3.0")); err != nil {
		t.Fatal(err)
	}
	for _, step := range f.steps {
		if strings.HasPrefix(step, "relink:") && strings.Contains(step, "v0.3.0") {
			t.Errorf("relinked the running version instead of the new one: %q", step)
		}
	}
}

// Relinking after the swap is not redundant. `current` moving means existing links keep
// resolving, but a version that *adds* an item leaves that item unlinked — which is the whole
// complaint `notify-users-of-new-updates` was queued for.
func TestUpdateRelinksSoNewItemsAppear(t *testing.T) {
	f := &fakeUpdate{latest: "v0.4.0"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("v0.3.0")); err != nil {
		t.Fatal(err)
	}
	if len(f.steps) == 0 || !strings.HasPrefix(f.steps[len(f.steps)-1], "relink") {
		t.Errorf("relink is not the last step: %v", f.steps)
	}
}

// The complaint that produced this whole change was the word `git pull` appearing in front of
// somebody who only wanted to use the tool. A promise about output is kept by asserting on
// output.
func TestUpdateFromAReleaseNeverMentionsGit(t *testing.T) {
	for _, failAt := range []string{"", "latest", "install", "relink"} {
		f := &fakeUpdate{latest: "v0.4.0", failAt: failAt}
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
func TestUpdateReportsWhichStepFailed(t *testing.T) {
	for step, want := range map[string]string{
		"latest":  "could not find",
		"install": "nothing was changed",
		"relink":  "link",
	} {
		f := &fakeUpdate{latest: "v0.4.0", failAt: step}
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

// A failed install changes nothing and relinks nothing. The links still point at the version
// that was working, and `go install` either completed or it did not — there is no partially
// downloaded module to reason about, because the go command does not leave one.
func TestAFailedUpdateLeavesThePreviousVersionActive(t *testing.T) {
	f := &fakeUpdate{latest: "v0.4.0", failAt: "install"}
	var out strings.Builder

	err := fromRelease(context.Background(), &out, f.deps("v0.3.0"))
	if err == nil {
		t.Fatal("a failed install was reported as success")
	}
	if !strings.Contains(err.Error(), "nothing was changed") {
		t.Errorf("the error does not say the machine is untouched: %v", err)
	}
	for _, step := range f.steps {
		if strings.HasPrefix(step, "relink") {
			t.Errorf("relinked after the install failed: %v", f.steps)
		}
	}
}

// **Gone with the design: there is no half-succeeded update.** The previous version could
// install a payload and fail to replace the binary, so it reported both and carried on. One
// `go install` either brings down the module or it does not, which deleted that whole branch
// and the report that explained it.

// Already on the newest release is a success that changes nothing and says so. Nothing is
// downloaded and nothing is relinked — a no-op that relinks is a no-op with a report.
func TestUpdateOnTheNewestVersionDoesNothing(t *testing.T) {
	f := &fakeUpdate{latest: "v0.4.0"}
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
func TestUpdateFromNothingInstalled(t *testing.T) {
	f := &fakeUpdate{latest: "v0.4.0"}
	var out strings.Builder

	if err := fromRelease(context.Background(), &out, f.deps("dev")); err != nil {
		t.Fatalf("upgrading with nothing installed: %v", err)
	}
	if got := strings.Join(f.steps, " "); !strings.Contains(got, "install:v0.4.0") {
		t.Errorf("steps = %q, want the version installed", got)
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

// **Gone with the copy: the row no longer names a mechanism.** A test pinned that it said
// `pull` in a checkout and `release` otherwise; the row is a label now and says neither. What
// the command actually did is in its output, which TestUpdateFromAReleaseNeverMentionsGit still
// holds to the same standard.

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
