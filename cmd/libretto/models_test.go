package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
)

// agent writes an agent into the repo and returns its path. Used as the source a
// target's symlink points at.
func (f fixture) agent(t *testing.T, name, model string) string {
	t.Helper()

	dir := filepath.Join(f.Repo, "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: " + name + "\ndescription: fixture\n"
	if model != "" {
		body += "model: " + model + "\n"
	}
	body += "---\n\nBody.\n"

	path := filepath.Join(dir, name+".md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func (f fixture) agentBody(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.Repo, "agents", name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// installed links a repo agent into the target, the way `install` would. Writing it
// reaches the repository's file — this is the "shared" case.
func (f fixture) installed(t *testing.T, name, model string) {
	t.Helper()
	f.link(t, f.agent(t, name, model), f.dest("agents", name+".md"))
}

// foreign plants a real agent file in the target that libretto did not create — the
// 22-agents-in-~/.claude case this whole change exists for.
func (f fixture) foreign(t *testing.T, name, model string) string {
	t.Helper()

	dir := filepath.Join(f.Claude, "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: " + name + "\ndescription: not ours\n"
	if model != "" {
		body += "model: " + model + "\n"
	}
	body += "---\n\nBody.\n"

	path := filepath.Join(dir, name+".md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func (f fixture) foreignBody(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.Claude, "agents", name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestModelsListsTheTargetsAgents(t *testing.T) {
	f := newFixture(t)
	f.installed(t, "review-design", "haiku")
	f.foreign(t, "sdd-apply", "sonnet")
	f.agent(t, "never-installed", "") // in the repo, absent from the target

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"review-design", "sdd-apply"} {
		if !strings.Contains(out, want) {
			t.Errorf("the target holds %q and it was not listed:\n%s", want, out)
		}
	}
	if strings.Contains(out, "never-installed") {
		t.Errorf("an agent absent from the target was listed:\n%s", out)
	}
}

// The whole reason for this change. Before it, an agent the user had and libretto did
// not ship was invisible — and on the machine that reported the bug that was 22 of
// them.
func TestModelsEditsAnAgentTheRepositoryDoesNotShip(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku", "sdd-apply"})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(f.foreignBody(t, "sdd-apply"), "model: haiku") {
		t.Error("an agent libretto did not create was not written")
	}
}

func TestModelsMarksSharedAgents(t *testing.T) {
	f := newFixture(t)
	f.installed(t, "review-design", "haiku")
	f.foreign(t, "sdd-apply", "sonnet")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}

	for _, line := range strings.Split(out, "\n") {
		switch {
		case strings.Contains(line, "review-design"):
			if !strings.Contains(line, "shared") {
				t.Errorf("a symlink into the repo was not marked shared: %q", line)
			}
		case strings.Contains(line, "sdd-apply"):
			if strings.Contains(line, "shared") {
				t.Errorf("a target-local file was marked shared: %q", line)
			}
		}
	}
}

func TestModelsSetSaysWhenAWriteIsShared(t *testing.T) {
	f := newFixture(t)
	f.installed(t, "review-design", "")

	out, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku", "review-design"})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(out), "every project") {
		t.Errorf("writing a shared agent did not say the effect is machine-wide:\n%s", out)
	}
}

// The shipped message told everyone their write was machine-wide. That is now true of
// the symlinked rows only, and saying it about a local file is the same class of error
// as the silence it replaced.
func TestModelsSetDoesNotOverclaimALocalWrite(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")

	out, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku", "sdd-apply"})
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(out), "every project") {
		t.Errorf("a target-local write claimed to be machine-wide:\n%s", out)
	}
}

func TestModelsListingDiffersBetweenScopes(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")

	globalOut, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}
	projectOut, _, err := capture(t, func() error { return models(f.Repo, f.project(), nil) })
	if err != nil {
		t.Fatal(err)
	}
	if globalOut == projectOut {
		t.Error("the two scopes produced identical listings — the flag changes nothing")
	}
}

func TestModelsShowsDefaultForAnUndeclaredAgent(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "session") {
		t.Errorf("an undeclared agent should read as the session's model:\n%s", out)
	}
}

func TestModelsSetAppliesToEveryNamedAgent(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")
	f.foreign(t, "jd-judge-a", "")
	f.foreign(t, "review-risk", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku", "sdd-apply", "jd-judge-a"})
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"sdd-apply", "jd-judge-a"} {
		if !strings.Contains(f.foreignBody(t, name), "model: haiku") {
			t.Errorf("%s did not get the model", name)
		}
	}
	if strings.Contains(f.foreignBody(t, "review-risk"), "model:") {
		t.Error("review-risk was named by nobody and still changed")
	}
}

func TestModelsSetAllReachesEveryAgent(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")
	f.installed(t, "review-design", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "sonnet", "--all"})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(f.foreignBody(t, "sdd-apply"), "model: sonnet") {
		t.Error("the foreign agent did not get the model")
	}
	if !strings.Contains(f.agentBody(t, "review-design"), "model: sonnet") {
		t.Error("the shared agent's destination did not get the model")
	}
}

// A destructive default that fires on a forgotten argument is how every agent on the
// machine silently becomes the same model — and the set is now 22, not 7.
func TestModelsSetWithoutAgentsIsAnError(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku"})
	})
	if err == nil {
		t.Fatal("set with no agents and no --all was accepted")
	}
	if strings.Contains(f.foreignBody(t, "sdd-apply"), "model:") {
		t.Error("an agent was written despite the refusal")
	}
}

func TestModelsSetRejectsAnUnknownModel(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "gpt-4", "--all"})
	})
	if err == nil {
		t.Fatal("an unknown model was accepted")
	}
	if strings.Contains(f.foreignBody(t, "sdd-apply"), "model:") {
		t.Error("an agent was written despite the unknown model")
	}
}

func TestModelsSetRejectsAnUnknownAgentAndWritesNothing(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku", "sdd-apply", "no-such-agent"})
	})
	if err == nil {
		t.Fatal("an unknown agent name was accepted")
	}
	if strings.Contains(f.foreignBody(t, "sdd-apply"), "model:") {
		t.Error("the valid agent in the set was written anyway")
	}
}

// The all-or-nothing guarantee was written for this repository's own tidy files. It
// now meets whatever is actually in somebody's ~/.claude/agents.
func TestModelsSetRefusesAStrayFileAndWritesNothing(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "")

	stray := filepath.Join(f.Claude, "agents", "NOTES.md")
	if err := os.WriteFile(stray, []byte("Just a document.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku", "--all"})
	})
	if err == nil {
		t.Fatal("a directory containing a non-agent file was accepted")
	}
	if strings.Contains(f.foreignBody(t, "sdd-apply"), "model:") {
		t.Error("an agent was written despite the stray file")
	}
	if got := string(mustRead(t, stray)); got != "Just a document.\n" {
		t.Errorf("the stray file was modified: %q", got)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestModelsOutputHasNoEscapeCodes(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "haiku")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "\x1b") {
		t.Errorf("output carries escape codes:\n%q", out)
	}
}

// A target that has never had an agent installed is a state, not a crash.
func TestModelsWithNoAgentsSaysSo(t *testing.T) {
	f := newFixture(t)

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no agents") {
		t.Errorf("expected a plain statement that there are none:\n%s", out)
	}
}

// The menu row answers "how much of this destination is still expensive". Counting a
// directory the user is not looking at would make the one row that reports live state
// report somebody else's.
func TestMenuTallyCountsTheActiveTargetsAgents(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "sdd-apply", "haiku")
	f.foreign(t, "jd-judge-a", "opus")

	menu, _, err := panelData(f.Repo, f.Project, target.GlobalScope)
	if err != nil {
		t.Fatal(err)
	}

	var models string
	for _, item := range menu {
		if item.Label == "models" {
			models = item.Desc
		}
	}
	if models == "" {
		t.Fatal("no models row in the menu")
	}
	for _, want := range []string{"1 on haiku", "1 on opus"} {
		if !strings.Contains(models, want) {
			t.Errorf("tally = %q, want it to contain %q", models, want)
		}
	}
}
