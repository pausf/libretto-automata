package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The README is the one document written for somebody who has never seen this project,
// so its shape is a promise and the shape is what gets tested. Wording is not — a test
// that pins prose is a test somebody deletes the first time they improve a sentence.

// flat collapses every run of whitespace to one space, so a phrase this file looks for
// still matches when Markdown wraps it across a newline. Without it "checksum database"
// silently never matched, and a guard that cannot fire is worse than no guard — every
// substring assertion below runs through here.
func flat(doc string) string {
	return strings.Join(strings.Fields(doc), " ")
}

// section returns the text under heading, up to the next `## `. Subsections stay in:
// reasoning hiding under a `### ` inside Install is still reasoning inside Install.
func section(t *testing.T, doc, heading string) string {
	t.Helper()
	start := strings.Index(doc, heading)
	if start < 0 {
		t.Fatalf("the README has no %q section", heading)
	}
	rest := doc[start+len(heading):]
	if end := strings.Index(rest, "\n## "); end >= 0 {
		return rest[:end]
	}
	return rest
}

func TestReadmeSectionsAreInReadingOrder(t *testing.T) {
	readme := repoFile(t, "README.md")

	previous, at := "", -1
	for _, heading := range []string{
		"## What you get",
		"## Install",
		"## Your first run",
		"## Commands",
		"## Learn more",
	} {
		i := strings.Index(readme, heading)
		if i < 0 {
			t.Fatalf("the README has no %q section", heading)
		}
		if i < at {
			t.Errorf("%q comes before %q — install has to precede what to type", heading, previous)
		}
		previous, at = heading, i
	}
}

func TestInstallSectionIsStepsOnly(t *testing.T) {
	install := flat(section(t, repoFile(t, "README.md"), "## Install"))

	for _, want := range []string{
		"go install github.com/pausf/libretto-automata/cmd/libretto@latest",
		"libretto install",
		"1.26",
	} {
		if !strings.Contains(install, want) {
			t.Errorf("the install section never says %q", want)
		}
	}
	// Why the payload rides inside the module is a good answer to a question nobody has
	// yet while they are counting steps. It lives in docs/DESIGN.md now.
	for _, moved := range []string{"GOMODCACHE", "@v0.", "checksum database"} {
		if strings.Contains(install, moved) {
			t.Errorf("the install section still explains %q — that belongs in docs/DESIGN.md", moved)
		}
	}
}

func TestReadmeWalksAFirstRun(t *testing.T) {
	first := flat(section(t, repoFile(t, "README.md"), "## Your first run"))

	// Not the wording — the stops. A walk that skips where the flow pauses is a walk
	// that surprises the reader at the first pause.
	for _, want := range []string{"/libretto-flow", "spec", "plan", "push"} {
		if !strings.Contains(first, want) {
			t.Errorf("the first-run walk never mentions %q", want)
		}
	}
}

// Moved and deleted look identical in a diff of one file. Every subject here has to be
// absent from the README and present in docs/ — that is the whole test, and it is the
// reason it reads three files.
func TestMovedReasoningLandedInDocs(t *testing.T) {
	readme := flat(repoFile(t, "README.md"))
	docs := flat(repoFile(t, "docs/DESIGN.md")) + " " + flat(repoFile(t, "docs/FLOW.md"))

	for _, moved := range []struct{ subject, phrase string }{
		{"no --force", "there is no `--force`"},
		{"prune is not uninstall", "Prune deliberately"},
		{"both scope flags is an error", "Two answers to one question"},
		{"prune and uninstall are dry by default", "being asked twice"},
		{"two queue commands and not one", "substitutes different work"},
		{"the payload is not compiled in", "module cache"},
		// Two anchors for one subject, because one was defeated by paraphrase: the README
		// kept the argument and swapped "module cache" for "an installed copy", which left
		// the guard green. "wins over" is the load-bearing verb either way.
		{"the payload is not compiled in, paraphrased", "wins over"},
		{"model aliases rather than ids", "two spellings of one state"},
		{"spec-drift warns, never blocks", "a deleted check finds nothing"},
	} {
		if strings.Contains(readme, moved.phrase) {
			t.Errorf("%s is still argued in the README", moved.subject)
		}
		if !strings.Contains(docs, moved.phrase) {
			t.Errorf("%s is in neither docs/DESIGN.md nor docs/FLOW.md — moved means moved, not deleted", moved.subject)
		}
	}
}

// `](#anchor)` never matches: the group demands at least one character that is neither
// `)` nor `#`.
var readmeLink = regexp.MustCompile(`\]\(([^)#]+)(?:#[^)]*)?\)`)

func TestReadmeLinksResolve(t *testing.T) {
	links := readmeLink.FindAllStringSubmatch(repoFile(t, "README.md"), -1)
	if len(links) == 0 {
		t.Fatal("no links matched — the pattern is broken, not the README")
	}

	for _, match := range links {
		target := match[1]
		if strings.HasPrefix(target, "http") || strings.HasPrefix(target, "mailto:") {
			continue
		}
		if _, err := os.Stat(filepath.Join("..", "..", target)); err != nil {
			t.Errorf("the README links to %s, which does not exist", target)
		}
	}
}
