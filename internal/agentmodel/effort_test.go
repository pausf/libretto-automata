package agentmodel

import (
	"os"
	"path/filepath"
	"testing"
)

// dirWithModels builds a throwaway directory of agent files, each declaring a model.
//
// The map is name → model, with the empty string meaning "declares none". dirWith
// above cannot express that, and the pairing rule is the whole subject here: a level
// is legal or not depending on the model of the agent it is being written to.
func dirWithModels(t *testing.T, models map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, model := range models {
		body := "---\nname: " + name + "\ndescription: A lens.\n"
		if model != Default {
			body += "model: " + model + "\n"
		}
		body += "---\n\nBody.\n"
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func effortOf(t *testing.T, dir, name string) string {
	t.Helper()
	got, err := ReadEffort(filepath.Join(dir, name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func writeAgent(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "a.md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadEffortReturnsTheDeclaredEffort(t *testing.T) {
	path := writeAgent(t, "---\nname: a\nmodel: opus\neffort: xhigh\n---\n\nBody.\n")

	got, err := ReadEffort(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "xhigh" {
		t.Errorf("ReadEffort() = %q, want %q", got, "xhigh")
	}
}

func TestReadEffortReportsDefaultWhenTheKeyIsAbsent(t *testing.T) {
	path := writeAgent(t, "---\nname: a\nmodel: opus\n---\n\nBody.\n")

	got, err := ReadEffort(path)
	if err != nil {
		t.Fatalf("an absent effort is a state, not an error: %v", err)
	}
	if got != Default {
		t.Errorf("ReadEffort() = %q, want the default (empty)", got)
	}
}

// The body of an agent is prose about effort levels. A reader that scans past the
// closing fence reports a level nobody set, from a line that is documentation.
func TestReadEffortIgnoresTheBody(t *testing.T) {
	path := writeAgent(t, "---\nname: a\n---\n\nSet\neffort: max\nto think harder.\n")

	got, err := ReadEffort(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != Default {
		t.Errorf("ReadEffort() = %q from the body, want the default", got)
	}
}

func TestSetEffortInsertsWithoutDisturbingTheFile(t *testing.T) {
	before := "---\nname: a\ndescription: A lens.\nmodel: opus\n---\n\nBody.\n"
	path := writeAgent(t, before)

	if err := SetEffort(path, "high"); err != nil {
		t.Fatal(err)
	}

	want := "---\nname: a\ndescription: A lens.\nmodel: opus\neffort: high\n---\n\nBody.\n"
	if got := read(t, path); got != want {
		t.Errorf("SetEffort() produced\n%q\nwant\n%q", got, want)
	}
}

func TestSetEffortReplacesInPlace(t *testing.T) {
	path := writeAgent(t, "---\nname: a\neffort: low\nmodel: opus\n---\n\nBody.\n")

	if err := SetEffort(path, "max"); err != nil {
		t.Fatal(err)
	}

	want := "---\nname: a\neffort: max\nmodel: opus\n---\n\nBody.\n"
	if got := read(t, path); got != want {
		t.Errorf("SetEffort() produced\n%q\nwant\n%q", got, want)
	}
}

func TestSetEffortDefaultRemovesTheKey(t *testing.T) {
	path := writeAgent(t, "---\nname: a\neffort: xhigh\nmodel: opus\n---\n\nBody.\n")

	if err := SetEffort(path, Default); err != nil {
		t.Fatal(err)
	}

	want := "---\nname: a\nmodel: opus\n---\n\nBody.\n"
	if got := read(t, path); got != want {
		t.Errorf("SetEffort(Default) produced\n%q\nwant the key gone:\n%q", got, want)
	}
}

// Identical bytes still move the mtime, and a tool that dirties a git tree for a
// no-op is one people stop running before a release.
func TestSetEffortIsIdempotent(t *testing.T) {
	path := writeAgent(t, "---\nname: a\neffort: high\n---\n\nBody.\n")
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := SetEffort(path, "high"); err != nil {
		t.Fatal(err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Error("SetEffort() rewrote a file that already declared the level")
	}
}

// The feature exists so an agent can stay on Opus while its effort drops. Two keys
// that disturb each other would make that impossible to express.
func TestTheTwoKeysDoNotDisturbEachOther(t *testing.T) {
	path := writeAgent(t, "---\nname: a\n---\n\nBody.\n")

	if err := SetModel(path, "opus"); err != nil {
		t.Fatal(err)
	}
	if err := SetEffort(path, "xhigh"); err != nil {
		t.Fatal(err)
	}
	if got, _ := ReadModel(path); got != "opus" {
		t.Errorf("model after an effort write = %q, want opus", got)
	}

	if err := SetModel(path, "sonnet"); err != nil {
		t.Fatal(err)
	}
	if got, _ := ReadEffort(path); got != "xhigh" {
		t.Errorf("effort after a model write = %q, want xhigh", got)
	}

	if err := SetEffort(path, Default); err != nil {
		t.Fatal(err)
	}
	if got, _ := ReadModel(path); got != "sonnet" {
		t.Errorf("model after removing the effort = %q, want sonnet", got)
	}
}

func TestEffortCatalogueListsTheFiveLevels(t *testing.T) {
	want := []string{"low", "medium", "high", "xhigh", "max"}

	got := Efforts()
	if len(got) != len(want) {
		t.Fatalf("Efforts() returned %d levels, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Name != want[i] {
			t.Errorf("Efforts()[%d].Name = %q, want %q — weakest first is contracted", i, got[i].Name, want[i])
		}
		if got[i].Label == "" {
			t.Errorf("Efforts()[%d] (%s) carries no label", i, got[i].Name)
		}
	}
}

// ultracode sends xhigh and turns on workflow orchestration. It is a Claude Code
// session mode, not a level, and no frontmatter accepts it.
func TestUnknownEffortIsRefused(t *testing.T) {
	for _, name := range []string{"ultracode", "highest", "HIGH", "", "xxhigh"} {
		if ValidEffort(name) {
			t.Errorf("ValidEffort(%q) = true, want false", name)
		}
	}
	for _, name := range []string{"low", "medium", "high", "xhigh", "max"} {
		if !ValidEffort(name) {
			t.Errorf("ValidEffort(%q) = false, want true", name)
		}
	}
}

func TestWhichModelsSupportEffort(t *testing.T) {
	if SupportsEffort("haiku") {
		t.Error("SupportsEffort(haiku) = true — Haiku appears in no row of the docs' effort table")
	}
	for _, model := range []string{"opus", "sonnet", Default} {
		if !SupportsEffort(model) {
			t.Errorf("SupportsEffort(%q) = false, want true", model)
		}
	}
}

func TestApplyEffortReachesEveryAgentInTheSet(t *testing.T) {
	dir := dirWithModels(t, map[string]string{"one": "opus", "two": "sonnet", "three": "opus"})

	if err := ApplyEffort(dir, []string{"one", "two"}, "xhigh"); err != nil {
		t.Fatal(err)
	}

	if got := effortOf(t, dir, "one"); got != "xhigh" {
		t.Errorf("one = %q, want xhigh", got)
	}
	if got := effortOf(t, dir, "two"); got != "xhigh" {
		t.Errorf("two = %q, want xhigh", got)
	}
	if got := effortOf(t, dir, "three"); got != Default {
		t.Errorf("three was not named and reads %q, want the default", got)
	}
}

// All-or-nothing, and the ones before the bad agent in the list are the point: a
// partial apply is a state the user cannot read back off the file.
func TestApplyEffortWritesNothingWhenAnyAgentCannotRunIt(t *testing.T) {
	dir := dirWithModels(t, map[string]string{"first": "opus", "cheap": "haiku"})

	err := ApplyEffort(dir, []string{"first", "cheap"}, "high")
	if err == nil {
		t.Fatal("ApplyEffort() accepted a level on a Haiku agent")
	}
	if got := effortOf(t, dir, "first"); got != Default {
		t.Errorf("first was written anyway and reads %q — the refusal was not all-or-nothing", got)
	}
	if got := effortOf(t, dir, "cheap"); got != Default {
		t.Errorf("cheap reads %q, want untouched", got)
	}
}

// The session model is unknowable from here, so refusing would be a guess in the
// other direction.
func TestApplyEffortAllowsAnAgentOnTheSessionModel(t *testing.T) {
	dir := dirWithModels(t, map[string]string{"inherits": Default})

	if err := ApplyEffort(dir, []string{"inherits"}, "low"); err != nil {
		t.Fatalf("ApplyEffort() refused an agent that declares no model: %v", err)
	}
	if got := effortOf(t, dir, "inherits"); got != "low" {
		t.Errorf("inherits = %q, want low", got)
	}
}

func TestApplyEffortRefusesAnUnknownLevel(t *testing.T) {
	dir := dirWithModels(t, map[string]string{"one": "opus"})

	if err := ApplyEffort(dir, []string{"one"}, "ultracode"); err == nil {
		t.Fatal("ApplyEffort() accepted ultracode, which is not a level")
	}
	if got := effortOf(t, dir, "one"); got != Default {
		t.Errorf("one reads %q after a refused level, want untouched", got)
	}
}

func TestAgentsReportsEachCurrentEffort(t *testing.T) {
	dir := dirWithModels(t, map[string]string{"declared": "opus", "inherits": Default})
	if err := SetEffort(filepath.Join(dir, "declared.md"), "max"); err != nil {
		t.Fatal(err)
	}

	agents, _, err := Agents(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := make(map[string]string, len(agents))
	for _, a := range agents {
		got[a.Name] = a.Effort
	}
	if got["declared"] != "max" {
		t.Errorf("declared.Effort = %q, want max", got["declared"])
	}
	if got["inherits"] != Default {
		t.Errorf("inherits.Effort = %q, want the default", got["inherits"])
	}
}

// A key claiming a setting the model has no concept of is the confident wrong answer
// the catalogue's own comment warns about.
func TestApplyModelClearsEffortWhenTheModelSupportsNone(t *testing.T) {
	dir := dirWithModels(t, map[string]string{"lens": "opus"})
	if err := SetEffort(filepath.Join(dir, "lens.md"), "xhigh"); err != nil {
		t.Fatal(err)
	}

	if err := Apply(dir, []string{"lens"}, "haiku"); err != nil {
		t.Fatal(err)
	}

	if got := effortOf(t, dir, "lens"); got != Default {
		t.Errorf("effort after moving to haiku = %q, want cleared", got)
	}
	if got := modelOf(t, dir, "lens"); got != "haiku" {
		t.Errorf("model = %q, want haiku", got)
	}
}

// The clearing is narrow on purpose: a model that does support effort leaves the
// level alone, because changing the tier is not a request to change the depth.
func TestApplyModelKeepsEffortWhenTheModelSupportsIt(t *testing.T) {
	dir := dirWithModels(t, map[string]string{"lens": "opus"})
	if err := SetEffort(filepath.Join(dir, "lens.md"), "xhigh"); err != nil {
		t.Fatal(err)
	}

	if err := Apply(dir, []string{"lens"}, "sonnet"); err != nil {
		t.Fatal(err)
	}

	if got := effortOf(t, dir, "lens"); got != "xhigh" {
		t.Errorf("effort after moving to sonnet = %q, want xhigh", got)
	}
}
