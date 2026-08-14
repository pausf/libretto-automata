package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
)

// The panel's remembered destination.
//
// Every case runs through newFixture, so CLAUDE_HOME points at a temporary directory
// and the preference is written there rather than into a real ~/.claude.

// Nothing remembered is not an error. It is the state of every machine before this
// feature existed, and it must keep opening where it always did.
func TestNothingRememberedOpensGlobal(t *testing.T) {
	newFixture(t)

	if tool, scope := rememberedScope(); tool != target.ClaudeTool || scope != target.GlobalScope {
		t.Errorf("with nothing remembered got %q/%q, want claude/global", tool, scope)
	}
}

// A word nobody recognises must not produce a destination nobody chose. This is
// target.Resolve's rule, followed rather than reinvented.
func TestUnrecognisedPreferenceIsGlobal(t *testing.T) {
	for _, content := range []string{"", "\n", "  ", "globl", "PROJECT", "0"} {
		t.Run("content="+content, func(t *testing.T) {
			newFixture(t)
			writePreference(t, content)

			if tool, scope := rememberedScope(); tool != target.ClaudeTool || scope != target.GlobalScope {
				t.Errorf("%q gave %q/%q, want claude/global", content, tool, scope)
			}
		})
	}
}

// The round trip, and the surrounding whitespace a human editing the file by hand
// would leave behind.
func TestRememberThenReadIsARoundTrip(t *testing.T) {
	newFixture(t)

	remember(target.CodexTool, target.ProjectScope)
	if tool, scope := rememberedScope(); tool != target.CodexTool || scope != target.ProjectScope {
		t.Fatalf("got %q/%q, want codex/project", tool, scope)
	}

	// The legacy one-word file a pre-tools session left behind still reads.
	writePreference(t, "  project\n")
	if tool, scope := rememberedScope(); tool != target.ClaudeTool || scope != target.ProjectScope {
		t.Errorf("legacy value gave %q/%q, want claude/project", tool, scope)
	}

	remember(target.ClaudeTool, target.GlobalScope)
	if tool, scope := rememberedScope(); tool != target.ClaudeTool || scope != target.GlobalScope {
		t.Errorf("after remembering claude/global got %q/%q", tool, scope)
	}
}

// It sits at the root, not under skills/, agents/ or commands/, and it is a regular
// file rather than a link this tool created — so the destructive paths have no reason
// to touch it. Checked rather than assumed, because the failure mode is `uninstall`
// removing a file it does not own.
func TestUninstallLeavesThePreferenceAlone(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	if _, _, err := capture(t, func() error { return install(f.Repo, f.global()) }); err != nil {
		t.Fatal(err)
	}
	remember(target.ClaudeTool, target.ProjectScope)

	if _, _, err := capture(t, func() error {
		return uninstall(f.Repo, f.global(), []string{"--yes"})
	}); err != nil {
		t.Fatal(err)
	}

	if _, scope := rememberedScope(); scope != target.ProjectScope {
		t.Errorf("uninstall took the preference with it: now %q", scope)
	}
}

// The feature, at the seam a test can actually reach: with no flag, the panel opens
// where it was left.
func TestTheRememberedDestinationOpensThePanel(t *testing.T) {
	newFixture(t)
	remember(target.CodexTool, target.ProjectScope)

	if tool, scope := openingScope(target.ClaudeTool, target.GlobalScope, ""); tool != target.CodexTool || scope != target.ProjectScope {
		t.Errorf("opened on %q/%q, want codex/project", tool, scope)
	}
}

// A flag says *this run*. The file says *where I work*. A flag typed once must not
// silently re-answer the question for every future session, so it wins and writes
// nothing.
func TestAnExplicitFlagBeatsTheRememberedDestination(t *testing.T) {
	newFixture(t)
	remember(target.ClaudeTool, target.ProjectScope)

	if _, scope := openingScope(target.ClaudeTool, target.GlobalScope, "global"); scope != target.GlobalScope {
		t.Errorf("--global gave %q, want %q", scope, target.GlobalScope)
	}
	if _, scope := rememberedScope(); scope != target.ProjectScope {
		t.Errorf("the flag overwrote the preference: now %q", scope)
	}
}

// The promise that makes one machine-wide value tolerable: a command typed into a
// terminal does not change meaning because of state left by an earlier session.
//
// Asserted at the scopeFlags seam rather than by running `install`: with a remembered
// `project` a regression would install into <cwd>/.claude, and <cwd> under `go test`
// is this package's own directory. A test that proves the bug by writing into the
// repository is not a test worth having.
func TestSubcommandsIgnoreTheRememberedDestination(t *testing.T) {
	f := newFixture(t)
	remember(target.ClaudeTool, target.ProjectScope)

	_, scope, chosen, rest, err := scopeFlags([]string{"install"})
	if err != nil {
		t.Fatal(err)
	}
	if chosen != "" {
		t.Errorf("chosen = %q, want empty for an argument list with no flag", chosen)
	}
	if scope != target.GlobalScope {
		t.Errorf("scope = %q, want %q — subcommands do not read the preference", scope, target.GlobalScope)
	}
	if got := target.Resolve(target.ClaudeTool, scope, f.Project).Root(); got != f.global().Root() {
		t.Errorf("resolved to %q, want the global root %q", got, f.global().Root())
	}
	if len(rest) != 1 || rest[0] != "install" {
		t.Errorf("rest = %v, want [install]", rest)
	}
}

// `tab` is what remembers. The panel asks for a fresh view of the destination it is
// moving to, and that is the one moment a switch has really happened.
func TestSwitchingDestinationRemembersIt(t *testing.T) {
	f := newFixture(t)

	// toolOrder is [claude, codex, opencode], so index 1 is codex.
	scope := target.GlobalScope
	menu, rows, err := panelRefresh(f.Repo, f.Project, &scope)(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(menu) == 0 || len(rows) == 0 {
		t.Fatal("the refresh returned an empty view")
	}
	if tool, _ := rememberedScope(); tool != target.CodexTool {
		t.Errorf("after switching to codex, remembered %q", tool)
	}
}

// The file has to agree with the screen. A failed refresh leaves the panel where it
// was — `panel`'s spec is explicit that showing one destination's counts under
// another's name is worse than showing an error — so nothing may be persisted.
func TestAFailedSwitchRemembersNothing(t *testing.T) {
	newFixture(t)

	if _, _, err := panelRefresh(t.TempDir(), t.TempDir(), nil)(len(toolOrder)); err == nil {
		t.Fatal("an out-of-range destination refreshed successfully")
	}
	if _, err := os.Stat(preferencePath()); !os.IsNotExist(err) {
		t.Errorf("a failed switch wrote the preference anyway (stat err: %v)", err)
	}
}

// Remembering is a convenience. A convenience that can fail the switch it decorates
// is worse than no convenience, so an unwritable root changes nothing but the
// remembering.
func TestAFailedWriteDoesNotFailTheSwitch(t *testing.T) {
	f := newFixture(t)

	// A read-only root, not a root that is a regular file. Pointing CLAUDE_HOME at a
	// file breaks panelData as well — it reads the target's `skills/` — and a test
	// where the whole refresh fails proves the opposite of this criterion. The root
	// stays readable so the strip can still be assembled; only the write is blocked.
	if err := os.Chmod(f.Claude, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(f.Claude, 0o755) })

	menu, rows, err := panelRefresh(f.Repo, f.Project, nil)(1)
	if err != nil {
		t.Fatalf("an unwritable preference failed the switch: %v", err)
	}
	if len(menu) == 0 || len(rows) == 0 {
		t.Error("the refresh returned an empty view")
	}
	// And it really was blocked, so this is not passing for the wrong reason.
	if tool, _ := rememberedScope(); tool != target.ClaudeTool {
		t.Errorf("the write was not actually blocked: remembered %q", tool)
	}
}

// writePreference plants the file directly, so a test can set up states remember()
// would never produce.
func writePreference(t *testing.T, content string) {
	t.Helper()

	path := preferencePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRememberedDestinationRecognisesNewTargets(t *testing.T) {
	for _, want := range []target.Tool{target.CodexTool, target.OpencodeTool} {
		t.Run(string(want), func(t *testing.T) {
			newFixture(t)
			remember(want, target.ProjectScope)
			if tool, scope := rememberedScope(); tool != want || scope != target.ProjectScope {
				t.Fatalf("remembered %q/project, got back %q/%q", want, tool, scope)
			}
		})
	}

	t.Run("garbage still falls back to claude/global", func(t *testing.T) {
		newFixture(t)
		remember(target.Tool("garbage"), target.Scope("nonsense"))
		if tool, scope := rememberedScope(); tool != target.ClaudeTool || scope != target.GlobalScope {
			t.Fatalf("garbage resolved to %q/%q, want claude/global", tool, scope)
		}
	})
}
