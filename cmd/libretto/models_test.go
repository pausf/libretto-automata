package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// agent writes an agent into the repo fixture and returns its path.
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

func TestModelsListsEveryAgentAndChangesNothing(t *testing.T) {
	f := newFixture(t)
	f.agent(t, "review-design", "haiku")
	f.agent(t, "spec-writer", "")
	before := f.agentBody(t, "review-design")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"review-design", "haiku", "spec-writer"} {
		if !strings.Contains(out, want) {
			t.Errorf("output does not mention %q:\n%s", want, out)
		}
	}
	if got := f.agentBody(t, "review-design"); got != before {
		t.Error("models with no arguments modified an agent file")
	}
}

func TestModelsShowsDefaultForAnUndeclaredAgent(t *testing.T) {
	f := newFixture(t)
	f.agent(t, "spec-writer", "")

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
	f.agent(t, "review-design", "")
	f.agent(t, "review-tests", "")
	f.agent(t, "review-security", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku", "review-design", "review-tests"})
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"review-design", "review-tests"} {
		if !strings.Contains(f.agentBody(t, name), "model: haiku") {
			t.Errorf("%s did not get the model", name)
		}
	}
	if strings.Contains(f.agentBody(t, "review-security"), "model:") {
		t.Error("review-security was named by nobody and still changed")
	}
}

func TestModelsSetAllReachesEveryAgent(t *testing.T) {
	f := newFixture(t)
	f.agent(t, "review-design", "")
	f.agent(t, "spec-writer", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "sonnet", "--all"})
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"review-design", "spec-writer"} {
		if !strings.Contains(f.agentBody(t, name), "model: sonnet") {
			t.Errorf("%s did not get the model", name)
		}
	}
}

// A destructive default that fires on a forgotten argument is how every agent on
// the machine silently becomes the same model.
func TestModelsSetWithoutAgentsIsAnError(t *testing.T) {
	f := newFixture(t)
	f.agent(t, "review-design", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku"})
	})
	if err == nil {
		t.Fatal("set with no agents and no --all was accepted")
	}
	if strings.Contains(f.agentBody(t, "review-design"), "model:") {
		t.Error("an agent was written despite the refusal")
	}
}

func TestModelsSetRejectsAnUnknownModel(t *testing.T) {
	f := newFixture(t)
	f.agent(t, "review-design", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "gpt-4", "--all"})
	})
	if err == nil {
		t.Fatal("an unknown model was accepted")
	}
	if strings.Contains(f.agentBody(t, "review-design"), "model:") {
		t.Error("an agent was written despite the unknown model")
	}
}

func TestModelsSetRejectsAnUnknownAgentAndWritesNothing(t *testing.T) {
	f := newFixture(t)
	f.agent(t, "review-design", "")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku", "review-design", "no-such-agent"})
	})
	if err == nil {
		t.Fatal("an unknown agent name was accepted")
	}
	if strings.Contains(f.agentBody(t, "review-design"), "model:") {
		t.Error("the valid agent in the set was written anyway")
	}
}

// The flag promises a project-local effect it cannot deliver: both scopes symlink
// to the same repo file. Saying so at the moment of writing is the difference
// between a documented limit and a surprise three weeks later.
func TestModelsSetUnderProjectScopeSaysTheEffectIsShared(t *testing.T) {
	f := newFixture(t)
	f.agent(t, "review-design", "")

	out, _, err := capture(t, func() error {
		return models(f.Repo, f.project(), []string{"set", "haiku", "review-design"})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(out), "every project") {
		t.Errorf("writing under --project did not warn that the effect is shared:\n%s", out)
	}
}

func TestModelsOutputHasNoEscapeCodes(t *testing.T) {
	f := newFixture(t)
	f.agent(t, "review-design", "haiku")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "\x1b") {
		t.Errorf("output carries escape codes:\n%q", out)
	}
}

// A repo with no agents/ directory is a state, not a crash.
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
