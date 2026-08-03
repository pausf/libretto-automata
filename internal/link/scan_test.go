package link

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
)

// repoWith builds a throwaway repo tree. Keys are paths relative to the repo
// root; a trailing slash means "directory", anything else is a file.
func repoWith(t *testing.T, paths ...string) string {
	t.Helper()
	root := t.TempDir()

	for _, p := range paths {
		full := filepath.Join(root, p)
		if p[len(p)-1] == '/' {
			if err := os.MkdirAll(full, 0o755); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func names(items []Item) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.Name
	}
	return out
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestItems(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		kind  target.Kind
		want  []string
	}{
		{
			name:  "skills are directories, sorted",
			paths: []string{"skills/po-assistant/", "skills/ci-cd-manager/"},
			kind:  target.Skills,
			want:  []string{"ci-cd-manager", "po-assistant"},
		},
		{
			name:  "agents are markdown files",
			paths: []string{"agents/jd-judge-a.md", "agents/sdd-apply.md"},
			kind:  target.Agents,
			want:  []string{"jd-judge-a.md", "sdd-apply.md"},
		},
		{
			name:  "a loose file among skills is not a skill",
			paths: []string{"skills/po-assistant/", "skills/notes.txt"},
			kind:  target.Skills,
			want:  []string{"po-assistant"},
		},
		{
			name:  "a directory among agents is not an agent",
			paths: []string{"agents/real.md", "agents/subdir/"},
			kind:  target.Agents,
			want:  []string{"real.md"},
		},
		{
			name:  "non-markdown files among commands are ignored",
			paths: []string{"commands/review.md", "commands/README.txt"},
			kind:  target.Commands,
			want:  []string{"review.md"},
		},
		{
			name:  "dotfiles are never items",
			paths: []string{"skills/.gitkeep", "skills/.hidden/", "skills/real/"},
			kind:  target.Skills,
			want:  []string{"real"},
		},
		{
			name:  "an empty kind directory yields nothing",
			paths: []string{"skills/"},
			kind:  target.Skills,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := repoWith(t, tt.paths...)

			items, err := Items(root, tt.kind)
			if err != nil {
				t.Fatalf("Items() error = %v", err)
			}
			if got := names(items); !equal(got, tt.want) {
				t.Errorf("Items() = %v, want %v", got, tt.want)
			}
		})
	}
}

// A repo that has never held items of a kind is normal, not broken.
func TestItemsMissingKindDirIsNotAnError(t *testing.T) {
	root := t.TempDir()

	items, err := Items(root, target.Skills)
	if err != nil {
		t.Fatalf("Items() error = %v, want nil", err)
	}
	if len(items) != 0 {
		t.Errorf("Items() = %v, want empty", names(items))
	}
}

func TestItemsCarryAbsoluteRepoPaths(t *testing.T) {
	root := repoWith(t, "skills/po-assistant/")

	items, err := Items(root, target.Skills)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d items, want 1", len(items))
	}

	want := filepath.Join(root, "skills", "po-assistant")
	if items[0].Path != want {
		t.Errorf("Path = %q, want %q", items[0].Path, want)
	}
	if items[0].Kind != target.Skills {
		t.Errorf("Kind = %q, want %q", items[0].Kind, target.Skills)
	}
}

func TestCounts(t *testing.T) {
	root := repoWith(t,
		"skills/a/", "skills/b/",
		"agents/one.md",
		"commands/",
	)
	t.Setenv(target.EnvClaudeHome, t.TempDir())

	counts, err := Counts(root, target.NewClaude())
	if err != nil {
		t.Fatalf("Counts() error = %v", err)
	}

	want := map[target.Kind]int{target.Skills: 2, target.Agents: 1, target.Commands: 0}
	for k, n := range want {
		if counts[k] != n {
			t.Errorf("Counts()[%s] = %d, want %d", k, counts[k], n)
		}
	}
}

// A kind the target does not accept must be absent, not zero — otherwise the UI
// cannot tell "no items" from "not applicable to this target".
func TestCountsOmitsRejectedKinds(t *testing.T) {
	root := repoWith(t, "skills/a/")
	t.Setenv(target.EnvClaudeHome, t.TempDir())

	counts, err := Counts(root, target.NewClaude())
	if err != nil {
		t.Fatal(err)
	}
	if _, present := counts[target.Kind("hooks")]; present {
		t.Error(`Counts() included "hooks", which no target accepts`)
	}
}
