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

func TestOpencodeAcceptsOnlySkills(t *testing.T) {
	sandbox := t.TempDir()
	t.Setenv(EnvOpencodeHome, sandbox)
	o := NewOpencode()

	if got, want := o.Dir(Skills), filepath.Join(sandbox, "skills"); got != want {
		t.Errorf("Dir(skills) = %q, want %q", got, want)
	}
	if !o.Accepts(Skills) {
		t.Error("Accepts(skills) = false, want true")
	}
	if len(o.Kinds()) != 1 || o.Kinds()[0] != Skills {
		t.Errorf("Kinds() = %v, want [skills]", o.Kinds())
	}
	for _, k := range []Kind{Agents, Commands, Kind("hooks")} {
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
