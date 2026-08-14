package target

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectRoot(t *testing.T) {
	dir := t.TempDir()
	p := NewProject(dir)

	if got, want := p.Root(), filepath.Join(dir, ".claude"); got != want {
		t.Fatalf("root is %s, want %s", got, want)
	}
	if got := p.Name(); got != "project" {
		t.Fatalf("name is %q", got)
	}

	// The same three kinds as the global target. An item that installs in one
	// scope and vanishes in the other would be indistinguishable from a bug.
	for _, k := range []Kind{Skills, Agents, Commands} {
		if !p.Accepts(k) {
			t.Fatalf("project rejects %s", k)
		}
		if got, want := p.Dir(k), filepath.Join(dir, ".claude", string(k)); got != want {
			t.Fatalf("dir for %s is %s, want %s", k, got, want)
		}
	}
}

func TestProjectExists(t *testing.T) {
	dir := t.TempDir()
	p := NewProject(dir)

	if p.Exists() {
		t.Fatal("a project with no .claude/ reports as configured")
	}

	if err := os.MkdirAll(p.Root(), 0o755); err != nil {
		t.Fatal(err)
	}
	if !p.Exists() {
		t.Fatal("a project with a .claude/ reports as unconfigured")
	}
}

func TestProjectWithNoDirectoryIsInert(t *testing.T) {
	// An empty dir must not resolve to "/.claude" — a target that cannot say
	// where it lives is a target nothing should be written to.
	p := NewProject("")

	if got := p.Root(); got != "" {
		t.Fatalf("root is %q, want empty", got)
	}
	if got := p.Dir(Skills); got != "" {
		t.Fatalf("dir is %q, want empty", got)
	}
	if p.Exists() {
		t.Fatal("an unrooted project reports as configured")
	}
}

func TestResolveScope(t *testing.T) {
	home := t.TempDir()
	t.Setenv(EnvClaudeHome, home)
	project := t.TempDir()

	if got := Resolve(ClaudeTool, GlobalScope, project).Root(); got != home {
		t.Fatalf("global resolved to %s, want %s", got, home)
	}

	want := filepath.Join(project, ".claude")
	if got := Resolve(ClaudeTool, ProjectScope, project).Root(); got != want {
		t.Fatalf("project resolved to %s, want %s", got, want)
	}
}

// Every tool resolves in both scopes, and each pair lands on its own root.
func TestResolveToolScopeMatrix(t *testing.T) {
	agents := t.TempDir()
	t.Setenv(EnvAgentsHome, agents)
	opencode := t.TempDir()
	t.Setenv(EnvOpencodeHome, opencode)
	project := t.TempDir()

	cases := []struct {
		tool  Tool
		scope Scope
		name  string
		root  string
	}{
		{CodexTool, GlobalScope, "codex", agents},
		{CodexTool, ProjectScope, "codex", filepath.Join(project, ".agents")},
		{OpencodeTool, GlobalScope, "opencode", opencode},
		{OpencodeTool, ProjectScope, "opencode", filepath.Join(project, ".opencode")},
	}
	for _, c := range cases {
		got := Resolve(c.tool, c.scope, project)
		if got.Name() != c.name || got.Root() != c.root {
			t.Errorf("%s/%s resolved to %s at %s, want %s at %s",
				c.tool, c.scope, got.Name(), got.Root(), c.name, c.root)
		}
		if got.Accepts(Agents) || got.Accepts(Commands) {
			t.Errorf("%s/%s accepts a kind beyond skills", c.tool, c.scope)
		}
	}
}

// No two (tool, scope) pairs may ever be the same place, or "isolated from the
// global config" is a claim with nothing behind it.
func TestScopesNeverShareARoot(t *testing.T) {
	t.Setenv(EnvClaudeHome, t.TempDir())
	t.Setenv(EnvAgentsHome, t.TempDir())
	t.Setenv(EnvOpencodeHome, t.TempDir())
	project := t.TempDir()

	seen := map[string]string{}
	for _, tool := range Tools {
		for _, s := range []Scope{GlobalScope, ProjectScope} {
			pair := string(tool) + "/" + string(s)
			root := Resolve(tool, s, project).Root()
			if prev, dup := seen[root]; dup {
				t.Fatalf("%s and %s both resolved to %s", prev, pair, root)
			}
			seen[root] = pair
		}
	}
}

// An unknown tool or scope resolves to claude/global rather than to nothing. A
// typo must not silently produce a target with no root that writes nowhere and
// reports success.
func TestUnknownScopeFallsBackToGlobal(t *testing.T) {
	home := t.TempDir()
	t.Setenv(EnvClaudeHome, home)

	if got := Resolve(ClaudeTool, Scope("nonsense"), "").Root(); got != home {
		t.Fatalf("unknown scope resolved to %q, want the global root", got)
	}
	if got := Resolve(Tool("nonsense"), GlobalScope, "").Root(); got != home {
		t.Fatalf("unknown tool resolved to %q, want the global root", got)
	}
}
