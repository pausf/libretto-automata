package main

import (
	"os"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
	"github.com/pausf/libretto-automata/internal/ui"
)

// uninstall is the pair of install: it removes links that are *working*, which prune
// deliberately spares. Every test here checks either that it removed ours, or that it
// left alone what is not.

func TestUninstallWithoutYesChangesNothing(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")
	if _, _, err := capture(t, func() error { return install(f.Repo, f.global()) }); err != nil {
		t.Fatal(err)
	}

	out, _, err := capture(t, func() error { return uninstall(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}
	if !isSymlinkTo(t, f.dest("skills", "alpha"), item) {
		t.Fatal("uninstall removed a link without being asked twice")
	}
	if !strings.Contains(out, "would remove") || !strings.Contains(out, "Nothing has been changed") {
		t.Errorf("the dry run does not say what it would do:\n%s", out)
	}
}

func TestUninstallYesRemovesOurLinks(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	f.skill(t, "beta")
	if _, _, err := capture(t, func() error { return install(f.Repo, f.global()) }); err != nil {
		t.Fatal(err)
	}

	out, _, err := capture(t, func() error { return uninstall(f.Repo, f.global(), []string{"--yes"}) })
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"alpha", "beta"} {
		if _, err := os.Lstat(f.dest("skills", n)); !os.IsNotExist(err) {
			t.Errorf("%s survived", n)
		}
	}
	if !strings.Contains(out, "2 removed") {
		t.Errorf("the report does not count what it removed:\n%s", out)
	}
}

// One destination per run, like everything else.
func TestUninstallProjectScopeLeavesGlobalAlone(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	for _, tg := range []target.Target{f.global(), f.project()} {
		if _, _, err := capture(t, func() error { return install(f.Repo, tg) }); err != nil {
			t.Fatal(err)
		}
	}

	if _, _, err := capture(t, func() error {
		return uninstall(f.Repo, f.project(), []string{"--yes"})
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Lstat(f.projectDest("skills", "alpha")); !os.IsNotExist(err) {
		t.Error("the project's link survived")
	}
	if _, err := os.Lstat(f.dest("skills", "alpha")); err != nil {
		t.Fatal("uninstalling the project removed a link from the global config")
	}
}

// A conflict is kept and reported. It is not ours, and the exit code has to say the
// destination is not in the state a bare "removed" would imply.
func TestUninstallReportsConflicts(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	f.skill(t, "theirs")
	if _, _, err := capture(t, func() error { return install(f.Repo, f.global()) }); err != nil {
		t.Fatal(err)
	}
	// Replace one of ours with somebody else's real file.
	_ = os.Remove(f.dest("skills", "theirs"))
	real := f.putReal(t, "skills", "theirs", "hours of work")

	out, _, err := capture(t, func() error { return uninstall(f.Repo, f.global(), []string{"--yes"}) })
	if err != nil {
		t.Fatalf("a kept conflict should not be an error on its own: %v", err)
	}
	if !strings.Contains(out, "keep") || !strings.Contains(out, "not ours") {
		t.Errorf("the conflict is not reported as kept:\n%s", out)
	}
	if body, rerr := os.ReadFile(real); rerr != nil || string(body) != "hours of work" {
		t.Fatal("uninstall destroyed a foreign entry")
	}
	if _, err := os.Lstat(f.dest("skills", "alpha")); !os.IsNotExist(err) {
		t.Error("our own link survived alongside the conflict")
	}
}

// Nothing of ours installed is a state, not an error.
func TestUninstallOnACleanDestinationSaysSo(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	out, _, err := capture(t, func() error { return uninstall(f.Repo, f.global(), []string{"--yes"}) })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "nothing of ours") {
		t.Errorf("output = %q, want it to say there is nothing installed", out)
	}
}

// The property that matters most: install, uninstall, and the destination is as it was.
func TestInstallThenUninstallIsARoundTrip(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	f.skill(t, "beta")

	before, err := os.ReadDir(f.Claude)
	if err != nil {
		t.Fatal(err)
	}

	for _, run := range []func() error{
		func() error { return install(f.Repo, f.global()) },
		func() error { return uninstall(f.Repo, f.global(), []string{"--yes"}) },
	} {
		if _, _, err := capture(t, run); err != nil {
			t.Fatal(err)
		}
	}

	// The skills directory is created by install and deliberately left behind, so the
	// round trip is about entries, not about the tree being byte-identical.
	rest, err := os.ReadDir(f.dest("skills", ""))
	if err != nil {
		t.Fatalf("the destination directory was removed: %v", err)
	}
	if len(rest) != 0 {
		names := make([]string, 0, len(rest))
		for _, e := range rest {
			names = append(names, e.Name())
		}
		t.Fatalf("uninstall left %v behind", names)
	}
	_ = before
}

// ── the panel ────────────────────────────────────────────────────────────────

func TestPanelOffersUninstallAsDestructive(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	menu, _, err := panelData(f.Repo, f.Project, target.ClaudeTool, target.GlobalScope)
	if err != nil {
		t.Fatal(err)
	}

	var found *ui.MenuItem
	for i := range menu {
		if menu[i].Label == "uninstall" {
			found = &menu[i]
		}
	}
	if found == nil {
		t.Fatal("the menu does not offer uninstall")
	}
	if !found.Enabled {
		t.Error("uninstall is offered but disabled")
	}
	if !found.Destructive {
		t.Error("uninstall is not marked destructive, so one press would remove everything")
	}
	if !strings.Contains(found.Desc, shorten(f.global().Root())) {
		t.Errorf("the row does not name the destination: %q", found.Desc)
	}
}

func TestPanelUninstallNeedsTwoPresses(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")
	if _, _, err := capture(t, func() error { return install(f.Repo, f.global()) }); err != nil {
		t.Fatal(err)
	}

	// One press: the panel runs it dry.
	if _, err := runCaptured("uninstall", f.Repo, f.global(), false); err != nil {
		t.Fatal(err)
	}
	if !isSymlinkTo(t, f.dest("skills", "alpha"), item) {
		t.Fatal("one press from the panel removed the link")
	}

	// Confirmed: gone.
	if _, err := runCaptured("uninstall", f.Repo, f.global(), true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(f.dest("skills", "alpha")); !os.IsNotExist(err) {
		t.Error("the confirmed press did not remove the link")
	}
}
