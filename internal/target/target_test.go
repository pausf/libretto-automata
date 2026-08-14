package target

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKindItemShape(t *testing.T) {
	tests := []struct {
		name  string
		kind  Kind
		isDir bool
		ext   string
	}{
		{"skills are directories", Skills, true, ""},
		{"agents are markdown files", Agents, false, ".md"},
		{"commands are markdown files", Commands, false, ".md"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.kind.ItemsAreDirs(); got != tt.isDir {
				t.Errorf("ItemsAreDirs() = %v, want %v", got, tt.isDir)
			}
			if got := tt.kind.Ext(); got != tt.ext {
				t.Errorf("Ext() = %q, want %q", got, tt.ext)
			}
		})
	}
}

func TestClaudeRootResolution(t *testing.T) {
	t.Run("CLAUDE_HOME wins", func(t *testing.T) {
		sandbox := t.TempDir()
		t.Setenv(EnvClaudeHome, sandbox)

		c := NewClaude()
		if c.Root() != sandbox {
			t.Fatalf("Root() = %q, want %q", c.Root(), sandbox)
		}
	})

	t.Run("falls back to home", func(t *testing.T) {
		t.Setenv(EnvClaudeHome, "")

		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("no home directory in this environment")
		}
		want := filepath.Join(home, ".claude")

		if got := NewClaude().Root(); got != want {
			t.Fatalf("Root() = %q, want %q", got, want)
		}
	})
}

func TestClaudeDirs(t *testing.T) {
	sandbox := t.TempDir()
	t.Setenv(EnvClaudeHome, sandbox)
	c := NewClaude()

	tests := []struct {
		kind Kind
		want string
	}{
		{Skills, filepath.Join(sandbox, "skills")},
		{Agents, filepath.Join(sandbox, "agents")},
		{Commands, filepath.Join(sandbox, "commands")},
	}
	for _, tt := range tests {
		t.Run(string(tt.kind), func(t *testing.T) {
			if got := c.Dir(tt.kind); got != tt.want {
				t.Errorf("Dir(%s) = %q, want %q", tt.kind, got, tt.want)
			}
			if !c.Accepts(tt.kind) {
				t.Errorf("Accepts(%s) = false, want true", tt.kind)
			}
		})
	}
}

func TestClaudeRejectsUnknownKind(t *testing.T) {
	t.Setenv(EnvClaudeHome, t.TempDir())

	if NewClaude().Accepts(Kind("hooks")) {
		t.Error(`Accepts("hooks") = true, want false`)
	}
}

func TestClaudeExists(t *testing.T) {
	t.Run("present root", func(t *testing.T) {
		t.Setenv(EnvClaudeHome, t.TempDir())
		if !NewClaude().Exists() {
			t.Error("Exists() = false for a real directory")
		}
	})

	t.Run("absent root", func(t *testing.T) {
		t.Setenv(EnvClaudeHome, filepath.Join(t.TempDir(), "nope"))
		if NewClaude().Exists() {
			t.Error("Exists() = true for a missing directory")
		}
	})

	t.Run("root is a file, not a directory", func(t *testing.T) {
		f := filepath.Join(t.TempDir(), "claude")
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Setenv(EnvClaudeHome, f)

		if NewClaude().Exists() {
			t.Error("Exists() = true for a regular file")
		}
	})
}

// A target must never leak its root into empty-string joins, or Dir would
// return a path relative to the process working directory and writes would land
// somewhere unexpected.
func TestUnresolvableRootYieldsEmptyDirs(t *testing.T) {
	for _, tg := range []Target{Claude{}, Codex{}, Opencode{}} {
		if e, ok := tg.(interface{ Exists() bool }); ok && e.Exists() {
			t.Errorf("%s: Exists() = true for an empty root", tg.Name())
		}
		for _, k := range tg.Kinds() {
			if got := tg.Dir(k); got != "" {
				t.Errorf("%s: Dir(%s) = %q, want empty", tg.Name(), k, got)
			}
		}
	}
}
