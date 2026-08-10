package agentmodel

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

// read is the whole file, for the comparisons that follow. Nothing here checks a
// key in isolation: the promise is about every other byte.
func read(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// An agent file is a prompt. A writer that reflows it, reorders it, drops a blank
// line or normalises a trailing newline has changed the prompt — so the assertion
// is the entire file, not the presence of the key.
func TestSetModelInsertsWithoutDisturbingTheFile(t *testing.T) {
	path := agentFile(t, undeclared)

	if err := SetModel(path, "haiku"); err != nil {
		t.Fatal(err)
	}

	want := `---
name: spec-writer
description: Writes one delta spec.
tools: Read, Grep, Glob, Skill, Write
model: haiku
---

You write **one** spec.
`
	if got := read(t, path); got != want {
		t.Errorf("file after insert:\n%q\nwant:\n%q", got, want)
	}
}

func TestSetModelReplacesInPlace(t *testing.T) {
	path := agentFile(t, declared)

	if err := SetModel(path, "opus"); err != nil {
		t.Fatal(err)
	}

	want := `---
name: review-design
description: The design lens.
tools: Read, Grep, Glob, Skill
model: opus
---

You are one lens of a five-lens review.
`
	if got := read(t, path); got != want {
		t.Errorf("file after replace:\n%q\nwant:\n%q", got, want)
	}
}

// Choosing the default removes the key. An absent key is already the language's way
// of saying "whatever the session runs on"; writing a word that means the same is a
// second spelling of one state.
func TestSetModelDefaultRemovesTheKey(t *testing.T) {
	path := agentFile(t, declared)

	if err := SetModel(path, Default); err != nil {
		t.Fatal(err)
	}

	want := `---
name: review-design
description: The design lens.
tools: Read, Grep, Glob, Skill
---

You are one lens of a five-lens review.
`
	if got := read(t, path); got != want {
		t.Errorf("file after removing the key:\n%q\nwant:\n%q", got, want)
	}
}

// The error return is half the promise. A refusal that already mangled the file is
// the failure worth catching, so both halves are asserted.
func TestSetModelRefusesAFileWithoutFrontmatter(t *testing.T) {
	const body = "Just a document.\n"
	path := agentFile(t, body)

	if err := SetModel(path, "haiku"); err == nil {
		t.Error("SetModel() accepted a file with no frontmatter")
	}
	if got := read(t, path); got != body {
		t.Errorf("refused file was modified:\n%q\nwant:\n%q", got, body)
	}
}

// Setting the model an agent already has must not touch the file at all. A rewrite
// that produces identical bytes still moves the mtime, and a tool that dirties a
// git tree for a no-op is a tool people stop running.
func TestSetModelIsIdempotent(t *testing.T) {
	path := agentFile(t, declared)

	// Stamped into the past on purpose. Comparing against the mtime the file
	// happens to have does not work: a rewrite lands in the same filesystem tick
	// and the timestamps compare equal, so the test passes with the guard removed.
	// That is not a hypothetical — it was written that way first and a mutation
	// proved it decorative.
	stamp := time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(path, stamp, stamp); err != nil {
		t.Fatal(err)
	}

	if err := SetModel(path, "haiku"); err != nil {
		t.Fatal(err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().UTC().Equal(stamp) {
		t.Error("SetModel() rewrote a file that already declared that model")
	}
	if got := read(t, path); got != declared {
		t.Errorf("file changed:\n%q\nwant:\n%q", got, declared)
	}
}

// A file with no trailing newline round-trips too. strings.Split and a naive Join
// disagree about that final element, and the disagreement shows up as a newline
// this tool invented.
func TestSetModelPreservesAMissingTrailingNewline(t *testing.T) {
	path := agentFile(t, "---\nname: x\ndescription: y\n---\nBody, no newline.")

	if err := SetModel(path, "haiku"); err != nil {
		t.Fatal(err)
	}

	want := "---\nname: x\ndescription: y\nmodel: haiku\n---\nBody, no newline."
	if got := read(t, path); got != want {
		t.Errorf("file:\n%q\nwant:\n%q", got, want)
	}
}
