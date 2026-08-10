package agentmodel

import (
	"os"
	"path/filepath"
	"testing"
)

// agentFile writes body to a temp file and returns its path.
func agentFile(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "some-agent.md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

const declared = `---
name: review-design
description: The design lens.
tools: Read, Grep, Glob, Skill
model: haiku
---

You are one lens of a five-lens review.
`

const undeclared = `---
name: spec-writer
description: Writes one delta spec.
tools: Read, Grep, Glob, Skill, Write
---

You write **one** spec.
`

func TestReadModelReturnsTheDeclaredModel(t *testing.T) {
	got, err := ReadModel(agentFile(t, declared))
	if err != nil {
		t.Fatal(err)
	}
	if got != "haiku" {
		t.Errorf("ReadModel() = %q, want %q", got, "haiku")
	}
}

// An agent that declares nothing runs on whatever the session runs on. That is a
// state, not a failure, and reporting it as an error would make the ordinary case
// the loud one.
func TestReadModelReportsDefaultWhenTheKeyIsAbsent(t *testing.T) {
	got, err := ReadModel(agentFile(t, undeclared))
	if err != nil {
		t.Fatal(err)
	}
	if got != Default {
		t.Errorf("ReadModel() = %q, want the default %q", got, Default)
	}
}

// The body of an agent file is prose about models. A reader that scans the whole
// file finds the first line that looks like a key and reports a model nobody set.
func TestReadModelIgnoresTheBody(t *testing.T) {
	body := `---
name: review-design
description: The design lens.
---

Some agents declare a model, like this:

model: opus

That line is documentation, not configuration.
`
	got, err := ReadModel(agentFile(t, body))
	if err != nil {
		t.Fatal(err)
	}
	if got != Default {
		t.Errorf("ReadModel() = %q from the body, want the default %q", got, Default)
	}
}

func TestReadModelRefusesAFileWithoutFrontmatter(t *testing.T) {
	if _, err := ReadModel(agentFile(t, "Just a document.\n")); err == nil {
		t.Error("ReadModel() accepted a file with no frontmatter")
	}
}
