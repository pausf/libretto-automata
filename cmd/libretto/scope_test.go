package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
	"github.com/pausf/libretto-automata/internal/ui"
)

// The promise of scopes is isolation: installing into one must leave the other
// exactly as it was. Every test here asserts the *absence* of a change in the
// scope that was not asked for, because that is the failure that would go
// unnoticed — a command that works and also quietly edits the global config.

func TestInstallProjectScopeLeavesGlobalAlone(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")

	if _, _, err := capture(t, func() error { return install(f.Repo, f.project()) }); err != nil {
		t.Fatalf("install --project failed: %v", err)
	}

	if !isSymlinkTo(t, f.projectDest("skills", "alpha"), item) {
		t.Error("alpha is not linked in the project scope")
	}
	if _, err := os.Lstat(f.dest("skills", "alpha")); !os.IsNotExist(err) {
		t.Fatal("a project install wrote into the global config")
	}
}

func TestInstallGlobalScopeLeavesProjectAlone(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")

	if _, _, err := capture(t, func() error { return install(f.Repo, f.global()) }); err != nil {
		t.Fatalf("install --global failed: %v", err)
	}

	if !isSymlinkTo(t, f.dest("skills", "alpha"), item) {
		t.Error("alpha is not linked in the global scope")
	}
	if _, err := os.Lstat(f.projectDest("skills", "alpha")); !os.IsNotExist(err) {
		t.Fatal("a global install wrote into the project")
	}
}

func TestPruneProjectScopeLeavesGlobalAlone(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	// A stale link in each scope: ours, pointing at an item that never existed.
	gone := f.Repo + "/skills/gone"
	f.link(t, gone, f.dest("skills", "gone"))
	f.link(t, gone, f.projectDest("skills", "gone"))

	if _, _, err := capture(t, func() error {
		return prune(f.Repo, f.project(), []string{"--yes"})
	}); err != nil {
		t.Fatalf("prune --project failed: %v", err)
	}

	if _, err := os.Lstat(f.projectDest("skills", "gone")); !os.IsNotExist(err) {
		t.Error("prune did not remove the project's stale link")
	}
	if _, err := os.Lstat(f.dest("skills", "gone")); err != nil {
		t.Fatal("prune --project removed a link from the global config")
	}
}

// ── flag parsing ─────────────────────────────────────────────────────────────

func TestDefaultScopeIsGlobal(t *testing.T) {
	scope, _, rest, err := scopeFlags([]string{"install"})
	if err != nil {
		t.Fatal(err)
	}
	if scope != target.GlobalScope {
		t.Fatalf("default scope is %q; every invocation written before scopes existed meant global", scope)
	}
	if len(rest) != 1 || rest[0] != "install" {
		t.Fatalf("rest is %v, want [install]", rest)
	}
}

func TestScopeFlagsAreRemovedFromTheArguments(t *testing.T) {
	// A flag left in the arguments reaches the subcommand, and `prune --project`
	// would be read as a confirmation that was never given.
	scope, _, rest, err := scopeFlags([]string{"prune", "--project", "--yes"})
	if err != nil {
		t.Fatal(err)
	}
	if scope != target.ProjectScope {
		t.Fatalf("scope is %q, want project", scope)
	}
	if len(rest) != 2 || rest[0] != "prune" || rest[1] != "--yes" {
		t.Fatalf("rest is %v, want [prune --yes]", rest)
	}
}

func TestBothScopeFlagsIsAnError(t *testing.T) {
	for _, args := range [][]string{
		{"install", "--global", "--project"},
		{"install", "--project", "--global"},
		{"-p", "install", "-g"},
	} {
		if _, _, _, err := scopeFlags(args); err == nil {
			t.Errorf("%v was accepted; two answers to one question is a mistake, not a precedence rule", args)
		}
	}
}

func TestRepeatingTheSameScopeFlagIsFine(t *testing.T) {
	// Harmless and unambiguous. Rejecting it would be pedantry.
	scope, _, _, err := scopeFlags([]string{"--project", "install", "--project"})
	if err != nil {
		t.Fatalf("repeating --project was rejected: %v", err)
	}
	if scope != target.ProjectScope {
		t.Fatalf("scope is %q, want project", scope)
	}
}

// ── the output says where it went ────────────────────────────────────────────

func TestOutputNamesTheScopeRoot(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	out, _, err := capture(t, func() error { return install(f.Repo, f.project()) })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, f.project().Root()) {
		t.Errorf("the output never names the root it wrote to.\ngot:\n%s", out)
	}

	out, _, err = capture(t, func() error { return status(f.Repo, f.global()) })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, f.Claude) {
		t.Errorf("status never names the root it read.\ngot:\n%s", out)
	}
}

// ── the panel must not be torn by a long root ────────────────────────────────

// Found by looking at a render, not by reading the code: a temporary directory is
// long enough to push the frame's right border out of alignment, which the panel
// spec forbids at every width.
func TestShortenKeepsRootsInsideTheBudget(t *testing.T) {
	long := "/private/tmp/claude-501/-Users-pau-sanchez-gitrepos-libretto-automata/" +
		"0942a1b3-52f9-4f16-9ad7-9117e82e9a30/scratchpad/proyecto/.claude"

	got := shorten(long)
	if n := len([]rune(got)); n > pathBudget {
		t.Fatalf("shorten returned %d columns, budget is %d: %q", n, pathBudget, got)
	}
	if !strings.HasPrefix(got, "…") {
		t.Errorf("something was removed and the result does not say so: %q", got)
	}
	// The tail is what identifies the directory, so it has to survive.
	if !strings.HasSuffix(got, "/.claude") {
		t.Errorf("the tail was lost: %q", got)
	}
}

func TestShortenLeavesShortPathsAlone(t *testing.T) {
	if got := shorten("/tmp/p/.claude"); got != "/tmp/p/.claude" {
		t.Errorf("a path inside the budget was altered: %q", got)
	}
}

func TestShortenPrefersTildeForTheHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skip("no home directory")
	}
	if got := shorten(home + "/.claude"); got != "~/.claude" {
		t.Errorf("shorten(%q) = %q, want ~/.claude", home+"/.claude", got)
	}
}

func TestShortenSaysWhenThereIsNoRoot(t *testing.T) {
	if got := shorten(""); got != "not configured" {
		t.Errorf("an empty root rendered as %q", got)
	}
}

// ── the strip must be able to tell its destinations apart ─────────────────────

// The bug this pins down: the strip showed `link.Counts`, which counts items in
// the *repo* filtered by the kinds a target accepts. Both scopes accept the same
// three kinds, so both rows showed identical numbers — always, by construction. It
// looked like "what is installed here?" and answered "what does the repo hold?"
// twice.
//
// Reported as "lo mismo y eso me parece raro", which it was.
func TestStripRowsReportTheirOwnState(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	f.skill(t, "beta")

	// One of the two is linked globally. Nothing is linked in the project.
	f.link(t, f.Repo+"/skills/alpha", f.dest("skills", "alpha"))

	_, rows, err := panelData(f.Repo, f.Project, target.GlobalScope)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != len(scopeOrder) {
		t.Fatalf("the strip has %d rows, want %d", len(rows), len(scopeOrder))
	}

	global, project := rows[0].Info, rows[1].Info
	if global == project {
		t.Fatalf("both rows report %q — a strip whose rows cannot differ is a strip that misleads", global)
	}
	if !strings.Contains(global, "linked") {
		t.Errorf("the global row does not mention the item that is linked: %q", global)
	}
	if strings.Contains(project, "linked") {
		t.Errorf("the project row claims something is linked when nothing is: %q", project)
	}
}

// The status row follows the active destination, not a sum of both. A count that
// mixed them would answer a question nobody asked.
func TestStatusRowFollowsTheActiveScope(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	f.link(t, f.Repo+"/skills/alpha", f.dest("skills", "alpha"))

	g, _, err := panelData(f.Repo, f.Project, target.GlobalScope)
	if err != nil {
		t.Fatal(err)
	}
	p, _, err := panelData(f.Repo, f.Project, target.ProjectScope)
	if err != nil {
		t.Fatal(err)
	}

	gs, ps := descOf(g, "status"), descOf(p, "status")
	if gs == ps {
		t.Fatalf("the status row reads %q for both destinations", gs)
	}
	if !strings.Contains(gs, "linked") || !strings.Contains(ps, "missing") {
		t.Errorf("global=%q project=%q — each should describe its own destination", gs, ps)
	}
}

func descOf(menu []ui.MenuItem, label string) string {
	for _, m := range menu {
		if m.Label == label {
			return m.Desc
		}
	}
	return ""
}

// ── the panel's actions actually run ─────────────────────────────────────────

// The bug: selecting `install` set a notice reading "running install…" and ran
// nothing at all. The panel promised work it never did, which is worse than a panel
// with no actions — and there was a test asserting the notice, so the lie looked
// verified.
func TestDispatchRunsTheAction(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")

	if _, _, err := capture(t, func() error {
		return dispatch("install", f.Repo, f.project(), false)
	}); err != nil {
		t.Fatalf("dispatch install failed: %v", err)
	}

	if !isSymlinkTo(t, f.projectDest("skills", "alpha"), item) {
		t.Fatal("dispatch reported no error and installed nothing")
	}
}

// Every menu label must dispatch. A label with no case behind it is a row that does
// nothing, which is the bug this pins down, one action at a time.
func TestEveryMenuLabelDispatches(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	menu, _, err := panelData(f.Repo, f.Project, target.GlobalScope)
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range menu {
		if !item.Enabled {
			continue
		}
		_, _, err := capture(t, func() error { return dispatch(item.Label, f.Repo, f.global(), false) })
		if err != nil && strings.Contains(err.Error(), "unknown action") {
			t.Errorf("menu offers %q and dispatch has no case for it", item.Label)
		}
	}
}

// Chosen from a menu, prune must still be dry. Nothing deletes on one keypress.
func TestDispatchedPruneIsDry(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	f.link(t, f.Repo+"/skills/gone", f.dest("skills", "gone"))

	out, _, err := capture(t, func() error { return dispatch("prune", f.Repo, f.global(), false) })
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(f.dest("skills", "gone")); os.IsNotExist(err) {
		t.Fatal("prune chosen from the menu deleted a link without being asked twice")
	}
	if !strings.Contains(out, "would remove") {
		t.Errorf("prune did not say what it would do:\n%s", out)
	}
}

// One lookup for "where is the project", threaded down. Resolving it separately in
// the strip and in the runner gave two answers that agreed only by accident — and
// they disagreed: the strip read one root while the action wrote to another.
func TestStripAndRunnerAgreeOnTheProjectRoot(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	_, rows, err := panelData(f.Repo, f.Project, target.ProjectScope)
	if err != nil {
		t.Fatal(err)
	}

	want := f.project().Root()
	if !strings.Contains(rows[1].Info, shorten(want)) {
		t.Fatalf("the strip names %q, the runner would write to %q", rows[1].Info, want)
	}

	// Install through the same door the panel uses, then the strip must notice.
	if _, err := runCaptured("install", f.Repo, target.Resolve(scopeOrder[1], f.Project), false); err != nil {
		t.Fatal(err)
	}
	_, rows, err = panelData(f.Repo, f.Project, target.ProjectScope)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rows[1].Info, "linked") {
		t.Fatalf("the strip still reports %q after installing into that very destination", rows[1].Info)
	}
}

// prune from the panel acts on the destination the strip is pointing at, and on no
// other. This is the destructive path, so it gets its own test rather than resting
// on the one that covers install.
func TestPanelPruneActsOnTheActiveDestinationOnly(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	// One of ours, orphaned, in each destination.
	gone := f.Repo + "/skills/gone"
	f.link(t, gone, f.dest("skills", "gone"))
	f.link(t, gone, f.projectDest("skills", "gone"))

	// The panel's prune is dry, so it must remove nothing anywhere.
	if _, err := runCaptured("prune", f.Repo, target.Resolve(scopeOrder[1], f.Project), false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(f.projectDest("skills", "gone")); err != nil {
		t.Fatal("prune from the panel deleted without being asked twice")
	}

	// And with confirmation, only the active destination loses its link.
	if _, _, err := capture(t, func() error {
		return prune(f.Repo, target.Resolve(scopeOrder[1], f.Project), []string{"--yes"})
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(f.projectDest("skills", "gone")); !os.IsNotExist(err) {
		t.Error("the project's stale link survived")
	}
	if _, err := os.Lstat(f.dest("skills", "gone")); err != nil {
		t.Fatal("pruning the project removed a link from the global config")
	}

	// The same isolation holds for the new destinations: pruning codex with its
	// own orphan leaves the global one standing.
	f.link(t, gone, filepath.Join(f.Codex, "skills", "gone"))
	if _, _, err := capture(t, func() error {
		return prune(f.Repo, target.Resolve(target.CodexScope, ""), []string{"--yes"})
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(f.Codex, "skills", "gone")); !os.IsNotExist(err) {
		t.Error("codex's stale link survived")
	}
	if _, err := os.Lstat(f.dest("skills", "gone")); err != nil {
		t.Fatal("pruning codex removed a link from the global config")
	}
}

func TestDestinationFlags(t *testing.T) {
	t.Run("each new flag resolves its destination", func(t *testing.T) {
		for flag, want := range map[string]target.Scope{
			"--codex":    target.CodexScope,
			"--opencode": target.OpencodeScope,
		} {
			scope, chosen, rest, err := scopeFlags([]string{flag, "install"})
			if err != nil {
				t.Fatalf("%s: %v", flag, err)
			}
			if scope != want || chosen != string(want) {
				t.Errorf("%s resolved to %q/%q, want %q", flag, scope, chosen, want)
			}
			if len(rest) != 1 || rest[0] != "install" {
				t.Errorf("%s: rest is %v, want [install]", flag, rest)
			}
		}
	})

	t.Run("any two destination flags is an error", func(t *testing.T) {
		for _, args := range [][]string{
			{"install", "--codex", "--opencode"},
			{"install", "--global", "--codex"},
			{"--opencode", "install", "-p"},
			{"--codex", "install", "-g"},
		} {
			if _, _, _, err := scopeFlags(args); err == nil {
				t.Errorf("%v was accepted; two answers to one question is a mistake, not a precedence rule", args)
			}
		}
	})
}

func TestInstallCodexLeavesOthersAlone(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")

	if _, _, err := capture(t, func() error {
		return install(f.Repo, target.Resolve(target.CodexScope, ""))
	}); err != nil {
		t.Fatalf("install --codex failed: %v", err)
	}

	if !isSymlinkTo(t, filepath.Join(f.Codex, "skills", "alpha"), item) {
		t.Error("alpha is not linked in the codex destination")
	}
	for name, root := range map[string]string{
		"global": f.Claude, "project": f.projectDest(), "opencode": f.Opencode,
	} {
		if _, err := os.Lstat(filepath.Join(root, "skills", "alpha")); !os.IsNotExist(err) {
			t.Errorf("a codex install wrote into the %s destination", name)
		}
	}
	for _, kind := range []string{"agents", "commands"} {
		if _, err := os.Lstat(filepath.Join(f.Codex, kind)); !os.IsNotExist(err) {
			t.Errorf("a codex install created a %s directory the target does not accept", kind)
		}
	}
}

func TestInstallOpencodeLeavesOthersAlone(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")

	if _, _, err := capture(t, func() error {
		return install(f.Repo, target.Resolve(target.OpencodeScope, ""))
	}); err != nil {
		t.Fatalf("install --opencode failed: %v", err)
	}

	if !isSymlinkTo(t, filepath.Join(f.Opencode, "skills", "alpha"), item) {
		t.Error("alpha is not linked in the opencode destination")
	}
	for name, root := range map[string]string{
		"global": f.Claude, "project": f.projectDest(), "codex": f.Codex,
	} {
		if _, err := os.Lstat(filepath.Join(root, "skills", "alpha")); !os.IsNotExist(err) {
			t.Errorf("an opencode install wrote into the %s destination", name)
		}
	}
	for _, kind := range []string{"agents", "commands"} {
		if _, err := os.Lstat(filepath.Join(f.Opencode, kind)); !os.IsNotExist(err) {
			t.Errorf("an opencode install created a %s directory the target does not accept", kind)
		}
	}
}
