package target

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpencodeRootResolution(t *testing.T) {
	t.Run("OPENCODE_HOME wins", func(t *testing.T) {
		sandbox := t.TempDir()
		t.Setenv(EnvOpencodeHome, sandbox)

		o := NewOpencode()
		if o.Root() != sandbox {
			t.Fatalf("Root() = %q, want %q", o.Root(), sandbox)
		}
	})

	t.Run("falls back to ~/.config/opencode", func(t *testing.T) {
		t.Setenv(EnvOpencodeHome, "")

		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("no home directory in this environment")
		}
		want := filepath.Join(home, ".config", "opencode")

		if got := NewOpencode().Root(); got != want {
			t.Fatalf("Root() = %q, want %q", got, want)
		}
	})
}

// The commands directory is the plural one dirUnderRoot already produces for every
// kind. OpenCode globs "{command,commands}/**/*.md", so the plural is one of the two
// names it looks for — verified in sst/opencode packages/opencode/src/config/command.ts,
// which is also where `symlink: true` makes a linked file readable.
func TestOpencodeAcceptsSkillsCommandsAndAgents(t *testing.T) {
	sandbox := t.TempDir()
	t.Setenv(EnvOpencodeHome, sandbox)
	o := NewOpencode()

	for _, k := range []Kind{Skills, Commands, Agents} {
		if got, want := o.Dir(k), filepath.Join(sandbox, string(k)); got != want {
			t.Errorf("Dir(%s) = %q, want %q", k, got, want)
		}
		if !o.Accepts(k) {
			t.Errorf("Accepts(%s) = false, want true", k)
		}
	}
	if len(o.Kinds()) != 3 || o.Kinds()[0] != Skills || o.Kinds()[1] != Commands || o.Kinds()[2] != Agents {
		t.Errorf("Kinds() = %v, want [skills commands agents]", o.Kinds())
	}
	for _, k := range []Kind{Kind("hooks"), Kind("modes")} {
		if o.Accepts(k) {
			t.Errorf("Accepts(%s) = true, want false", k)
		}
	}

	t.Run("exists follows the root", func(t *testing.T) {
		if !o.Exists() {
			t.Error("Exists() = false for a real directory")
		}
		t.Setenv(EnvOpencodeHome, filepath.Join(t.TempDir(), "nope"))
		if NewOpencode().Exists() {
			t.Error("Exists() = true for a missing directory")
		}
	})
}
