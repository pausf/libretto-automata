package target

import (
	"bytes"
	"strings"
	"testing"
)

const claudeAgent = `---
name: work-reviewer
description: Reviews one finished change against its spec.
tools: Read, Grep, Glob, Bash, Skill
model: sonnet
---

The body. Two paragraphs, and one of them mentions tools: in prose.

Second paragraph.
`

func TestOpencodeAgentTransform(t *testing.T) {
	got, err := Opencode{}.Transform(Agents, "/repo/agents/work-reviewer.md", []byte(claudeAgent))
	if err != nil {
		t.Fatal(err)
	}
	out := string(got)

	// The frontmatter block only — the body deliberately contains "tools:" in prose,
	// so a check over the whole file would pass for the wrong reason.
	fm := out[:strings.Index(out[4:], "\n---\n")+4]

	if strings.Contains(fm, "tools:") {
		t.Error("tools: survived the transform — on OpenCode it is a Record, not a string, and the schema throws")
	}
	if strings.Contains(fm, "model:") {
		t.Error("model: survived — a Claude tier name is not a provider/model-id")
	}
	for _, want := range []string{
		"name: work-reviewer",
		"description: Reviews one finished change against its spec.",
		"mode: subagent",
		MarkerKey + `: "/repo/agents/work-reviewer.md"`,
	} {
		if !strings.Contains(fm, want) {
			t.Errorf("frontmatter is missing %q", want)
		}
	}

	// The body is the prompt. Changing a byte of it changes what the agent does.
	body := out[strings.Index(out[4:], "\n---\n")+len("\n---\n")+4:]
	wantBody := claudeAgent[strings.Index(claudeAgent[4:], "\n---\n")+len("\n---\n")+4:]
	if body != wantBody {
		t.Errorf("body changed.\n got: %q\nwant: %q", body, wantBody)
	}
}

// Two runs must agree byte for byte. Frontmatter emitted from a Go map would not:
// map iteration is randomised, so `status` would report drift on every scan and
// `install` would rewrite files forever.
func TestAgentTransformIsDeterministic(t *testing.T) {
	var first []byte
	for i := 0; i < 20; i++ {
		got, err := Opencode{}.Transform(Agents, "/repo/agents/a.md", []byte(claudeAgent))
		if err != nil {
			t.Fatal(err)
		}
		if first == nil {
			first = got
			continue
		}
		if !bytes.Equal(first, got) {
			t.Fatalf("run %d differs from run 0:\n%q\n%q", i, first, got)
		}
	}
}

// Re-transforming output must not stack a second marker or a second mode.
func TestTransformIsIdempotentOverItsOwnOutput(t *testing.T) {
	once, err := Opencode{}.Transform(Agents, "/repo/agents/a.md", []byte(claudeAgent))
	if err != nil {
		t.Fatal(err)
	}
	twice, err := Opencode{}.Transform(Agents, "/repo/agents/a.md", once)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(once, twice) {
		t.Errorf("transforming twice differs from once:\n%q\n%q", once, twice)
	}
	if n := strings.Count(string(twice), MarkerKey+":"); n != 1 {
		t.Errorf("marker appears %d times, want 1", n)
	}
	if n := strings.Count(string(twice), "mode: subagent"); n != 1 {
		t.Errorf("mode appears %d times, want 1", n)
	}
}

func TestTransformRefusesALinkedKind(t *testing.T) {
	for _, k := range []Kind{Skills, Commands, Kind("hooks")} {
		if _, err := (Opencode{}).Transform(k, "/repo/x", []byte(claudeAgent)); err == nil {
			t.Errorf("Transform(%s) returned no error — a caller asking to transform a linked kind has a bug", k)
		}
	}
}

// A source with no usable frontmatter cannot be transformed. Emitting it anyway would
// produce a file OpenCode throws on, which breaks its whole config load rather than
// skipping one agent.
func TestTransformRefusesAFileWithNoFrontmatter(t *testing.T) {
	for name, content := range map[string]string{
		"no opening":     "name: x\n---\nbody\n",
		"never closed":   "---\nname: x\nbody with no close\n",
		"empty":          "",
		"body only":      "just prose\n",
		"opens mid-file": "\n---\nname: x\n---\n",
	} {
		if _, err := (Opencode{}).Transform(Agents, "/repo/agents/a.md", []byte(content)); err == nil {
			t.Errorf("%s: Transform returned no error", name)
		}
	}
}

// Only opencode transforms. If claude or a project target ever satisfied this
// interface, every destination would change how it installs — silently.
func TestOnlyOpencodeTransforms(t *testing.T) {
	for _, tg := range []Target{NewClaude(), NewCodex(), NewProject("/tmp/x"), NewCodexProject("/tmp/x")} {
		if _, ok := tg.(Transformer); ok {
			t.Errorf("%s implements Transformer — only opencode installs by transform", tg.Name())
		}
	}
	for _, tg := range []Target{NewOpencode(), NewOpencodeProject("/tmp/x")} {
		if _, ok := tg.(Transformer); !ok {
			t.Errorf("%s does not implement Transformer", tg.Name())
		}
	}
}

// The marker is emitted as a quoted YAML scalar. A checkout path may contain a `#`, a
// colon or a trailing space, and unquoted, ` #` starts a YAML comment — so the host
// would read a different path than the one written. OpenCode throws on an agent it
// cannot parse rather than skipping it, so getting this wrong breaks its config load.
func TestMarkerIsQuotedSoAwkwardPathsSurvive(t *testing.T) {
	for _, source := range []string{
		"/repo/agents/a.md",
		"/repo with spaces/agents/a.md",
		"/repo #1/agents/a.md",
		"/repo/a: b/agents/a.md",
		`/repo/with"quote/agents/a.md`,
		`/repo/with\backslash/agents/a.md`,
	} {
		out, err := Opencode{}.Transform(Agents, source, []byte(claudeAgent))
		if err != nil {
			t.Fatalf("%s: %v", source, err)
		}
		want := MarkerKey + ": " + quoteYAML(source)
		if !strings.Contains(string(out), want) {
			t.Errorf("%s: frontmatter does not carry %s", source, want)
		}
		if strings.Contains(string(out), MarkerKey+": "+source+"\n") && strings.ContainsAny(source, ` #:"\`) {
			t.Errorf("%s: emitted unquoted", source)
		}
	}
}
