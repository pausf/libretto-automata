package target

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCodexRootResolution(t *testing.T) {
	t.Run("AGENTS_HOME wins", func(t *testing.T) {
		sandbox := t.TempDir()
		t.Setenv(EnvAgentsHome, sandbox)

		c := NewCodex()
		if c.Root() != sandbox {
			t.Fatalf("Root() = %q, want %q", c.Root(), sandbox)
		}
	})

	t.Run("falls back to home", func(t *testing.T) {
		t.Setenv(EnvAgentsHome, "")

		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("no home directory in this environment")
		}
		want := filepath.Join(home, ".agents")

		if got := NewCodex().Root(); got != want {
			t.Fatalf("Root() = %q, want %q", got, want)
		}
	})
}

func TestCodexAcceptsOnlySkills(t *testing.T) {
	sandbox := t.TempDir()
	t.Setenv(EnvAgentsHome, sandbox)
	c := NewCodex()

	if got, want := c.Dir(Skills), filepath.Join(sandbox, "skills"); got != want {
		t.Errorf("Dir(skills) = %q, want %q", got, want)
	}
	if !c.Accepts(Skills) {
		t.Error("Accepts(skills) = false, want true")
	}
	if len(c.Kinds()) != 1 || c.Kinds()[0] != Skills {
		t.Errorf("Kinds() = %v, want [skills]", c.Kinds())
	}
	for _, k := range []Kind{Agents, Commands, Kind("hooks")} {
		if c.Accepts(k) {
			t.Errorf("Accepts(%s) = true, want false", k)
		}
	}
}

func TestCodexExists(t *testing.T) {
	t.Run("present root", func(t *testing.T) {
		t.Setenv(EnvAgentsHome, t.TempDir())
		if !NewCodex().Exists() {
			t.Error("Exists() = false for a real directory")
		}
	})

	t.Run("absent root", func(t *testing.T) {
		t.Setenv(EnvAgentsHome, filepath.Join(t.TempDir(), "nope"))
		if NewCodex().Exists() {
			t.Error("Exists() = true for a missing directory")
		}
	})
}
