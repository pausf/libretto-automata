package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/agentmodel"
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

// effortOfAgent reads the effort a repo agent declares, through the same reader the
// command uses. Reading the raw body and grepping would pass a `effort:` that landed
// in the body.
func (f fixture) effortOfAgent(t *testing.T, name string) string {
	t.Helper()
	got, err := agentmodel.ReadEffort(filepath.Join(f.Repo, "agents", name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func (f fixture) foreignEffort(t *testing.T, name string) string {
	t.Helper()
	got, err := agentmodel.ReadEffort(filepath.Join(f.Claude, "agents", name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestModelsListsEffortBesideTheModel(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "deep", "opus")
	f.foreign(t, "inherits", "")
	if err := agentmodel.SetEffort(filepath.Join(f.Claude, "agents", "deep.md"), "xhigh"); err != nil {
		t.Fatal(err)
	}

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "xhigh") {
		t.Errorf("the listing never shows a declared effort:\n%s", out)
	}
	// An agent declaring neither key shows the session word twice, and both columns
	// mean the same thing by it.
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "inherits") {
			if strings.Count(line, "(session)") != 2 {
				t.Errorf("an agent declaring neither key should read (session) in both columns, got:\n%s", line)
			}
		}
	}
}

func TestModelsListsTheEffortCatalogue(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "deep", "opus")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}

	for _, level := range []string{"low", "medium", "high", "xhigh", "max"} {
		if !strings.Contains(out, level) {
			t.Errorf("the trailer never names %q:\n%s", level, out)
		}
	}
	if strings.Contains(out, "ultracode") {
		t.Error("ultracode is a session mode, not a level, and no frontmatter accepts it")
	}
	// Haiku carrying the same silent row as the others would read as a level it can
	// take. The listing has to say it cannot.
	if !strings.Contains(out, "no effort levels") {
		t.Errorf("the catalogue never says which model has no effort levels:\n%s", out)
	}
	// The levels table decays exactly as the model versions do — an organisation can
	// cap which are available, and nothing here can ask. The date is what makes that
	// staleness visible instead of silent, which is why one trailer carrying it and
	// the other not was a finding.
	if !strings.Contains(out, agentmodel.Resolved) {
		t.Errorf("the effort trailer never states when it was last checked:\n%s", out)
	}
}

func TestModelsEffortWritesOnlyTheNamedAgents(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "named", "opus")
	f.foreign(t, "spared", "opus")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"effort", "xhigh", "named"})
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := f.foreignEffort(t, "named"); got != "xhigh" {
		t.Errorf("named = %q, want xhigh", got)
	}
	if got := f.foreignEffort(t, "spared"); got != "" {
		t.Errorf("spared was not named and reads %q, want untouched", got)
	}
}

func TestModelsEffortAllReachesEveryAgent(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "one", "opus")
	f.foreign(t, "two", "sonnet")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"effort", "low", "--all"})
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"one", "two"} {
		if got := f.foreignEffort(t, name); got != "low" {
			t.Errorf("%s = %q, want low", name, got)
		}
	}
}

// A destructive default that fires on a forgotten argument is how every agent on the
// machine silently becomes the same thing.
func TestModelsEffortRefusesWithNothingNamed(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "one", "opus")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"effort", "max"})
	})
	if err == nil {
		t.Fatal("models effort with no agents and no --all was accepted")
	}
	if got := f.foreignEffort(t, "one"); got != "" {
		t.Errorf("one reads %q after a refusal, want untouched", got)
	}
}

func TestModelsEffortRefusesAnUnknownLevel(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "one", "opus")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"effort", "ultracode", "one"})
	})
	if err == nil {
		t.Fatal("models effort accepted ultracode")
	}
	for _, level := range []string{"low", "high", "max"} {
		if !strings.Contains(err.Error(), level) {
			t.Errorf("the refusal never names %q: %v", level, err)
		}
	}
	if got := f.foreignEffort(t, "one"); got != "" {
		t.Errorf("one reads %q after a refusal, want untouched", got)
	}
}

func TestModelsEffortRefusesAnAgentThatCannotRunTheLevel(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "capable", "opus")
	f.foreign(t, "cheap", "haiku")

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"effort", "high", "--all"})
	})
	if err == nil {
		t.Fatal("a Haiku agent was given an effort level")
	}
	if !strings.Contains(err.Error(), "cheap") {
		t.Errorf("the refusal does not name the offending agent: %v", err)
	}
	if got := f.foreignEffort(t, "capable"); got != "" {
		t.Errorf("capable was written anyway and reads %q — the refusal was not all-or-nothing", got)
	}
}

func TestModelsEffortDefaultRemovesTheKey(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "one", "opus")
	if err := agentmodel.SetEffort(filepath.Join(f.Claude, "agents", "one.md"), "max"); err != nil {
		t.Fatal(err)
	}

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"effort", "default", "one"})
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := f.foreignEffort(t, "one"); got != "" {
		t.Errorf("one = %q, want the key gone", got)
	}
	if strings.Contains(f.foreignBody(t, "one"), "effort:") {
		t.Error("the key is still in the file")
	}
}

// A dropped effort is an edit to a prompt file. A silent one is the failure the
// clearing exists to avoid.
func TestModelsSetReportsAClearedEffort(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "lens", "opus")
	if err := agentmodel.SetEffort(filepath.Join(f.Claude, "agents", "lens.md"), "xhigh"); err != nil {
		t.Fatal(err)
	}

	out, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"set", "haiku", "lens"})
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := f.foreignEffort(t, "lens"); got != "" {
		t.Errorf("effort after moving to haiku = %q, want cleared", got)
	}
	if !strings.Contains(out, "cleared") {
		t.Errorf("the clearing happened and the report never said so:\n%s", out)
	}
}

// A verb the binary accepts and the help never mentions is a verb nobody finds. That
// is not hypothetical: `models effort` shipped working — the panel key, the footer, the
// command — and was reported as missing, because `libretto help` listed `models set`
// and stopped there.
//
// The list is not maintained here. It is read out of the dispatch's own unknown-command
// error, which already names every verb it accepts, so a third verb added there without
// a line in the help fails this rather than passing quietly.
func TestEveryModelsVerbIsInTheHelp(t *testing.T) {
	f := newFixture(t)

	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"nonesuch"})
	})
	if err == nil {
		t.Fatal("an unknown models verb was accepted")
	}

	// Anchored to the backtick the error quotes each form in. Unanchored, the pattern
	// captured `command` out of the error's own "unknown models command" prefix and
	// then failed looking for it in the help — which is the test finding a bug in
	// itself, and the reason this comment names the anchor rather than assuming it.
	verbs := regexp.MustCompile("`models ([a-z]+)").FindAllStringSubmatch(err.Error(), -1)
	if len(verbs) == 0 {
		t.Fatalf("the dispatch error names no verbs, so this test proves nothing: %v", err)
	}

	_, help, _ := capture(t, func() error { usage(); return nil })
	for _, v := range verbs {
		if !strings.Contains(help, "models "+v[1]) {
			t.Errorf("the dispatch accepts `models %s` and the help never offers it:\n%s", v[1], help)
		}
	}
}

// The version column is only checkable against something. A reader who disagrees with
// `sonnet · Sonnet 4.5` needs to know it was read off the Amazon Bedrock row.
func TestModelsNamesTheProviderItResolvedFor(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "one", "opus")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "resolved for") {
		t.Errorf("the listing never says which provider the versions are for:\n%s", out)
	}
}

// The hole this closed: the catalogue said `sonnet` supports all five levels, which is
// true on the Anthropic API and on no other provider. Under Bedrock the alias is Sonnet
// 4.5, which appears in no row of the host's effort table.
func TestModelsReflectsTheProviderInTheVersionColumn(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "one", "sonnet")
	t.Setenv("CLAUDE_CODE_USE_BEDROCK", "1")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "Amazon Bedrock") {
		t.Errorf("the listing does not name Bedrock:\n%s", out)
	}
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "sonnet ") {
			continue
		}
		if !strings.Contains(line, "Sonnet 4.5") {
			t.Errorf("sonnet should resolve to Sonnet 4.5 on Bedrock: %q", line)
		}
		if !strings.Contains(line, "no effort levels") {
			t.Errorf("Sonnet 4.5 has no effort levels and the row does not say so: %q", line)
		}
		return
	}
	t.Fatalf("no sonnet row in the catalogue:\n%s", out)
}

// And the refusal follows the provider, not the alias. A level on a sonnet agent is legal
// on the Anthropic API and must be refused under Bedrock.
func TestModelsEffortFollowsTheProviderNotTheAlias(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "one", "sonnet")

	// First without a provider: allowed.
	_, _, err := capture(t, func() error {
		return models(f.Repo, f.global(), []string{"effort", "xhigh", "one"})
	})
	if err != nil {
		t.Fatalf("xhigh on sonnet was refused on the Anthropic API: %v", err)
	}

	t.Setenv("CLAUDE_CODE_USE_BEDROCK", "1")
	_, _, err = capture(t, func() error {
		return models(f.Repo, f.global(), []string{"effort", "xhigh", "one"})
	})
	if err == nil {
		t.Fatal("xhigh on sonnet was accepted under Bedrock, where sonnet is Sonnet 4.5")
	}
	if !strings.Contains(err.Error(), "one") {
		t.Errorf("the refusal does not name the agent: %v", err)
	}
}

func TestModelsListsTheRecommendationAndItsReason(t *testing.T) {
	f := newFixture(t)
	// A payload agent this repository has an opinion about, and one it does not.
	f.foreign(t, "review-lens-tests", "haiku")
	f.foreign(t, "somebody-elses-agent", "opus")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}

	// Both on one line, and that is not fussiness: the catalogue trailer already
	// prints "pattern-matching over prose" as haiku's own label, so a bare Contains
	// passes against output this change has not touched. It did, on the first run.
	found := false
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "review-lens-tests") && strings.Contains(line, "pattern-matching over prose") {
			found = true
		}
	}
	if !found {
		t.Errorf("no line carries the agent and its reason together:\n%s", out)
	}
	// Silence for somebody's own agent — absent from the trailer, not blank in it.
	//
	// Asserted against the trailer rather than against a glyph. The first version looked
	// for "→" on that agent's line, and the trailer emits "←"; the assertion could not
	// fail, so deleting the guard that skips unrecommended agents left it green.
	_, trailer, ok := strings.Cut(out, "recommended, and never applied")
	if !ok {
		t.Fatalf("no recommendation trailer at all:\n%s", out)
	}
	if strings.Contains(trailer, "somebody-elses-agent") {
		t.Errorf("an agent we have no opinion about appears in the trailer:\n%s", trailer)
	}
}

func TestModelsMarksAgentsRunningAgainstTheRecommendation(t *testing.T) {
	f := newFixture(t)
	// Declares sonnet, recommended haiku — the real disagreement this ships with.
	f.foreign(t, "review-lens-design", "sonnet")
	// Declares exactly what it is recommended.
	f.foreign(t, "review-lens-tests", "haiku")

	out, _, err := capture(t, func() error { return models(f.Repo, f.global(), nil) })
	if err != nil {
		t.Fatal(err)
	}

	// A recommendation nobody can tell they are ignoring changes nothing.
	if !strings.Contains(out, "runs sonnet today") {
		t.Errorf("an agent running against its recommendation is not marked as doing so:\n%s", out)
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "review-lens-tests") && strings.Contains(line, "today") {
			t.Errorf("an agent already at its recommendation was marked as diverging:\n%s", line)
		}
	}
}

func TestAgentRowsCarryTheRecommendation(t *testing.T) {
	f := newFixture(t)
	f.foreign(t, "review-lens-security", "sonnet")
	f.foreign(t, "somebody-elses-agent", "opus")

	rows, err := agentRows(f.Repo, f.global())
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, r := range rows {
		got[r.Name] = r.Recommended + "/" + r.RecommendedEffort
	}
	if got["review-lens-security"] != "sonnet/xhigh" {
		t.Errorf("review-lens-security row carries %q, want sonnet/xhigh", got["review-lens-security"])
	}
	// The seam is the point: the panel is handed an answer, never the rule.
	if got["somebody-elses-agent"] != "/" {
		t.Errorf("an unrecommended agent's row carries %q, want nothing", got["somebody-elses-agent"])
	}
}
