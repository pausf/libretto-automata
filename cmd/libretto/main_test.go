package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/link"
)

// These tests cover the composition, not the logic it composes. `linking`,
// `link-state`, `targets` and `panel` prove that the pieces behave; this file
// proves the CLI wires them to the right exit codes and refuses what it promised
// to refuse.
//
// Every criterion the cli spec lists as owed a test has one here.

// ── exit codes ───────────────────────────────────────────────────────────────

// A conflict is the case that matters. The item was not linked, so the install is
// incomplete, and a script gating on the exit code has to learn that.
func TestInstallExitsNonZeroOnConflict(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	real := f.putReal(t, "skills", "alpha", "somebody else's file")

	_, _, err := capture(t, func() error { return install(f.Repo, f.global()) })
	if err == nil {
		t.Fatal("install returned nil for a conflicting item; a script would read that as success")
	}

	// And the foreign file is still exactly as it was.
	body, readErr := os.ReadFile(real)
	if readErr != nil {
		t.Fatalf("the conflicting file is gone: %v", readErr)
	}
	if string(body) != "somebody else's file" {
		t.Errorf("the conflicting file was modified: %q", body)
	}
}

func TestInstallExitsZeroWhenEverythingLinks(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")

	_, _, err := capture(t, func() error { return install(f.Repo, f.global()) })
	if err != nil {
		t.Fatalf("install failed on a clean target: %v", err)
	}
	if !isSymlinkTo(t, f.dest("skills", "alpha"), item) {
		t.Error("alpha was reported as linked but the symlink is not there")
	}
}

// Running it twice is not a special case that needed handling — an already
// correct tree produces an empty plan. This is the test that says so.
func TestInstallIsIdempotent(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	if _, _, err := capture(t, func() error { return install(f.Repo, f.global()) }); err != nil {
		t.Fatalf("first install: %v", err)
	}
	out, _, err := capture(t, func() error { return install(f.Repo, f.global()) })
	if err != nil {
		t.Fatalf("second install: %v", err)
	}
	if !strings.Contains(out, "already correct") {
		t.Errorf("second run did not report an empty plan:\n%s", out)
	}
}

func TestDoctorExitsNonZeroWhenSomethingNeedsAttention(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha") // never installed, so: missing

	_, _, err := capture(t, func() error { return doctor(f.Repo, f.global()) })
	if err == nil {
		t.Fatal("doctor returned nil with a missing item")
	}
}

// ── prune is dry by default ──────────────────────────────────────────────────

// staleLink builds the one thing prune is allowed to remove: a link this repo
// owns whose item no longer exists.
func staleLink(t *testing.T, f fixture, name string) string {
	t.Helper()

	dir := f.dest("skills", "")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	dst := f.dest("skills", name)
	if err := os.Symlink(filepath.Join(f.Repo, "skills", name), dst); err != nil {
		t.Fatal(err)
	}
	return dst
}

func TestPruneWithoutYesChangesNothing(t *testing.T) {
	f := newFixture(t)
	stale := staleLink(t, f, "gone")

	out, _, err := capture(t, func() error { return prune(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatalf("a dry run must not fail: %v", err)
	}
	if !exists(stale) {
		t.Fatal("prune removed a link without --yes")
	}
	if !strings.Contains(out, "would remove") {
		t.Errorf("the dry run did not say what it would do:\n%s", out)
	}
}

func TestPruneYesRemovesOnlyWhatThePlanNamed(t *testing.T) {
	f := newFixture(t)
	stale := staleLink(t, f, "gone")

	// Two things that must survive: an item that still exists, and a file that
	// was never ours.
	item := f.skill(t, "alpha")
	if _, _, err := capture(t, func() error { return install(f.Repo, f.global()) }); err != nil {
		t.Fatalf("setup install: %v", err)
	}
	foreign := f.putReal(t, "skills", "theirs", "not ours")

	if _, _, err := capture(t, func() error { return prune(f.Repo, f.global(), []string{"--yes"}) }); err != nil {
		t.Fatalf("prune --yes: %v", err)
	}

	if exists(stale) {
		t.Error("the stale link survived prune --yes")
	}
	if !isSymlinkTo(t, f.dest("skills", "alpha"), item) {
		t.Error("prune removed a link whose item still exists")
	}
	if !exists(foreign) {
		t.Error("prune removed a foreign entry — it must never do that")
	}
}

// ── dispatch ─────────────────────────────────────────────────────────────────

func TestRunDispatch(t *testing.T) {
	f := newFixture(t)

	cases := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"version", []string{"version"}, false},
		{"short version flag", []string{"-v"}, false},
		{"help", []string{"help"}, false},
		{"models", []string{"models"}, false},
		{"models with an unknown verb", []string{"models", "nosuchthing"}, true},
		{"unknown command", []string{"nosuchthing"}, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, _, err := capture(t, func() error { return runIn(f.Repo, c.args) })
			if (err != nil) != c.wantErr {
				t.Errorf("run(%v) error = %v, want error: %v", c.args, err, c.wantErr)
			}
		})
	}
}

// runIn runs a subcommand against a chosen root, the way run does once it has
// resolved the repo. It exists so dispatch can be tested without the process
// discovering the real repository around it.
func runIn(root string, args []string) error {
	t := os.Getenv("LIBRETTO_ROOT")
	os.Setenv("LIBRETTO_ROOT", root)
	defer os.Setenv("LIBRETTO_ROOT", t)
	return run(args)
}

// No TTY and no subcommand exits non-zero rather than opening a panel it cannot
// drive. os.Exit cannot be observed in-process, so this re-executes the test
// binary and reads the real status code.
func TestNoTTYAndNoSubcommandExitsNonZero(t *testing.T) {
	if os.Getenv("LIBRETTO_TEST_NO_ARGS") == "1" {
		// Stdout here is a pipe, never a terminal, which is the condition
		// under test.
		_ = run(nil)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNoTTYAndNoSubcommandExitsNonZero")
	cmd.Env = append(os.Environ(), "LIBRETTO_TEST_NO_ARGS=1", "LIBRETTO_ROOT="+t.TempDir())
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("exited zero with no TTY and no subcommand; it must not")
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() != 2 {
		t.Errorf("exit code = %d, want 2", ee.ExitCode())
	}
	if !strings.Contains(string(out), "show the panel") {
		t.Errorf("usage was not printed:\n%s", out)
	}
}

// ── piped output is plain ────────────────────────────────────────────────────

func TestStatusOutputHasNoEscapeCodes(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	out, _, err := capture(t, func() error { return status(f.Repo, f.global()) })
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if strings.Contains(out, "\x1b") {
		t.Errorf("piped status emitted ANSI escape codes:\n%q", out)
	}
	if !strings.Contains(out, "missing") || !strings.Contains(out, "alpha") {
		t.Errorf("status did not report the item's state:\n%s", out)
	}
}

// ── the command names itself ─────────────────────────────────────────────────

// Linked under another name, every line that names a command must use that name.
// A remedy naming a command the reader does not have is a remedy they cannot run.
func TestOutputNamesTheInvokedCommand(t *testing.T) {
	orig := os.Args[0]
	os.Args[0] = "/somewhere/lb"
	defer func() { os.Args[0] = orig }()

	if got := invokedAs(); got != "lb" {
		t.Fatalf("invokedAs() = %q, want %q", got, "lb")
	}
	if got := remedy(link.Missing); !strings.Contains(got, "lb install") {
		t.Errorf("remedy(missing) = %q, want it to name lb", got)
	}
	if got := remedy(link.Stale); !strings.Contains(got, "lb prune") {
		t.Errorf("remedy(stale) = %q, want it to name lb", got)
	}

	_, stderr, _ := capture(t, func() error { usage(); return nil })
	if !strings.Contains(stderr, "lb install") {
		t.Errorf("usage did not name the invoked command:\n%s", stderr)
	}
	if strings.Contains(stderr, "lib ") {
		t.Errorf("usage still carries a hardcoded name:\n%s", stderr)
	}
}

// ── prerequisites are informational ──────────────────────────────────────────

// rg and jq are what the payload's own gates run on — spec-drift asks every
// question through rg, find-work parses jira through jq. A doctor that omits them
// leaves the silent-false-negative machine looking healthy.
func TestPrerequisitesIncludeTheGateTools(t *testing.T) {
	names := map[string]bool{}
	for _, p := range prerequisites() {
		names[p.Name] = true
	}
	for _, want := range []string{"rg", "jq"} {
		if !names[want] {
			t.Errorf("prerequisites() missing %q", want)
		}
	}
}

// None of the prerequisites is required, so their absence must never reach the
// exit code. Failing on an optional absence trains people to ignore the section.
func TestPrerequisitesDoNotAffectTheExitCode(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")
	if _, _, err := capture(t, func() error { return install(f.Repo, f.global()) }); err != nil {
		t.Fatalf("setup install: %v", err)
	}
	if !isSymlinkTo(t, f.dest("skills", "alpha"), item) {
		t.Fatal("setup did not link the item")
	}

	// Nothing on PATH and an empty home: every prerequisite is absent.
	t.Setenv("PATH", "")
	t.Setenv("HOME", t.TempDir())

	out, _, err := capture(t, func() error { return doctor(f.Repo, f.global()) })
	if err != nil {
		t.Fatalf("doctor failed with every optional tool missing: %v", err)
	}
	if !strings.Contains(out, "prerequisites") {
		t.Errorf("doctor did not print the prerequisite report:\n%s", out)
	}
}

// A companion installed as a plain command is installed. Checking only the plugin
// tree reports it as missing and sends people to install what they already have.
func TestCompanionFoundWhereverItLegitimatelyLives(t *testing.T) {
	for _, where := range []string{
		filepath.Join("plugins", "some-marketplace", "caveman"),
		filepath.Join("plugins", "some-marketplace", "nested", "caveman"),
		filepath.Join("skills", "caveman"),
		filepath.Join("commands", "caveman.md"),
	} {
		t.Run(where, func(t *testing.T) {
			home := t.TempDir()
			path := filepath.Join(home, ".claude", where)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, nil, 0o644); err != nil {
				t.Fatal(err)
			}
			if !companion(home, "caveman") {
				t.Errorf("companion installed at %s reported as missing", where)
			}
		})
	}
}

func TestCompanionAbsentIsAbsent(t *testing.T) {
	if companion(t.TempDir(), "caveman") {
		t.Error("reported a companion that is not installed")
	}
}

// ── reporting helpers ────────────────────────────────────────────────────────

// The tally renders in a fixed order and omits zeros, so the same situation
// always produces the same line.
func TestSummarise(t *testing.T) {
	cases := []struct {
		name   string
		counts map[link.State]int
		want   string
	}{
		{"empty", nil, "nothing to link"},
		{"one state", map[link.State]int{link.Missing: 3}, "3 missing"},
		{
			"fixed order regardless of map iteration",
			map[link.State]int{link.Stale: 1, link.Linked: 2, link.Conflict: 1},
			"2 linked · 1 conflict · 1 stale",
		},
		{"zeros are omitted", map[link.State]int{link.Linked: 1, link.Missing: 0}, "1 linked"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := summarise(c.counts); got != c.want {
				t.Errorf("summarise() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestPlural(t *testing.T) {
	cases := []struct {
		kind string
		n    int
		want string
	}{
		{"skills", 1, "skill"},
		{"skills", 0, "skills"},
		{"skills", 2, "skills"},
		{"commands", 1, "command"},
	}

	for _, c := range cases {
		if got := plural(c.kind, c.n); got != c.want {
			t.Errorf("plural(%q, %d) = %q, want %q", c.kind, c.n, got, c.want)
		}
	}
}

func TestShort(t *testing.T) {
	if got := short("7be3da3aabbccdd"); got != "7be3da3" {
		t.Errorf("short() = %q, want %q", got, "7be3da3")
	}
	if got := short("abc"); got != "abc" {
		t.Errorf("short() truncated something already short: %q", got)
	}
}
