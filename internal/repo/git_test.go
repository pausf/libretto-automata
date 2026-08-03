package repo

import "testing"

func TestNeedsRebuild(t *testing.T) {
	cases := []struct {
		name  string
		paths []string
		want  bool
	}{
		{"nothing changed", nil, false},
		{"only payload", []string{"skills/write-spec/SKILL.md", "docs/FLOW.md"}, false},
		{"a go file", []string{"internal/link/plan.go"}, true},
		{"go file among markdown", []string{"README.md", "cmd/libretto/main.go"}, true},
		{"go.mod", []string{"go.mod"}, true},
		{"go.sum", []string{"go.sum"}, true},
		// A path that merely mentions go must not trigger a compile.
		{"markdown about go", []string{"docs/going-further.md"}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NeedsRebuild(c.paths); got != c.want {
				t.Errorf("NeedsRebuild(%v) = %v, want %v", c.paths, got, c.want)
			}
		})
	}
}
