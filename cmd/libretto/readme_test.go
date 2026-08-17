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

// A payload command that ships without reaching the README is invisible to everyone who
// never opens `commands/`. That happened once, to `libretto-attacca`, under a spec that
// already asked in prose for one line per command.
//
// The whole file, not the `## Commands` section: that table is the *binary's* subcommands,
// and the slash commands live in the first-run door list. Writing this guard against the
// heading that shares their name failed all six on its first run — which is the reason it
// reads the whole README and the reason this comment exists.
//
// The directory is read rather than listed here on purpose: a list in this file is the same
// failure one level down, where somebody adds a command, forgets the list, and the guard
// stays green.
func TestEveryCommandIsInTheReadme(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("..", "..", "commands"))
	if err != nil {
		t.Fatalf("cannot read commands/: %v", err)
	}

	commands := flat(repoFile(t, "README.md"))
	found := 0
	for _, entry := range entries {
		name, ok := strings.CutSuffix(entry.Name(), ".md")
		if entry.IsDir() || !ok {
			continue
		}
		found++
		// Word-bounded, not `strings.Contains`: a substring match lets a new command ride
		// on a longer name already in the file — `libretto-stat` would be satisfied by the
		// `/libretto-status` line and ship unmentioned. `-` is not a word character here,
		// so `libretto-stat\b` refuses `libretto-status` and accepts `/libretto-stat `.
		if !regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\b`).MatchString(commands) {
			t.Errorf("commands/%s.md ships and the README never mentions %s", name, name)
		}
	}
	if found == 0 {
		t.Fatal("no command files matched — the directory moved, not the README")
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

// badgeImage captures the image URL of every `[![alt](image)](link)` badge — the image
// only, so a link target or alt text that happens to mention a build is not caught.
//
// Run it over flat(), never the raw document: `[^)\s]+` cannot cross a newline, so a badge
// whose markdown wraps between `![Build](` and its URL slipped past the scan entirely while
// the six honest badges kept the count non-zero and the guard green. Found by the 6→7
// reviewer with a probe badge, which is the second time this file has been bitten by a
// guard that silently could not fire. `\s*` after the paren is what flat() leaves behind.
var badgeImage = regexp.MustCompile(`!\[[^\]]*\]\(\s*(https://[^)\s]+)\)`)

// A shields.io/badge/… URL is a literal: whatever it says, it says forever. That is fine
// for a fact that does not change — a language version, a tool name, a licence — and a lie
// for a run outcome. The tests badge was `tests-passing-brightgreen` linking to the
// workflow *file*, so it claimed green whether or not anything had run, which is the exact
// failure skills/evidence exists to refuse: nothing is true until it has been observed.
//
// Word-bounded, not substring. The command guard in this same file shipped as a substring
// match and `libretto-stat` would have ridden on `/libretto-status`; here a substring would
// refuse an honest badge for any tool whose name contains `pass` or `build`, and a guard
// that false-positives on honest content is a guard somebody deletes.
func TestNoBadgeAssertsAStatus(t *testing.T) {
	readme := repoFile(t, "README.md")

	claims := regexp.MustCompile(`(?i)\b(passing|failing|pass|fail|build|coverage)\b`)
	badges := badgeImage.FindAllStringSubmatch(flat(readme), -1)
	if len(badges) == 0 {
		t.Fatal("no badges matched — the pattern is broken, not the README")
	}

	for _, match := range badges {
		image := match[1]
		if !strings.Contains(image, "shields.io/badge/") {
			continue // a live endpoint reports a real run; only literals can lie
		}
		if word := claims.FindString(image); word != "" {
			t.Errorf("the badge %s hardcodes %q — a literal cannot report a run.\n"+
				"Use the workflow's own badge endpoint, which turns red on its own.", image, word)
		}
	}

	// The rule above is satisfied by a README with no tests badge at all, so the one badge
	// that must exist is pinned separately. Broad rule plus one anchor; neither covers the
	// other.
	const endpoint = "actions/workflows/gates.yml/badge.svg"
	if !strings.Contains(flat(readme), endpoint) {
		t.Errorf("the README does not carry %s — the tests badge must report the real run", endpoint)
	}
}

var mermaidFence = regexp.MustCompile("(?s)```mermaid(.*?)```")

// Two diagrams live inside What you get — delivery and flow — because that section is
// where a stranger decides whether to install. GitHub renders the fences as images; the
// guard holds the shape (two fences, the stops named), never the drawing.
func TestWhatYouGetCarriesTheDiagrams(t *testing.T) {
	got := section(t, repoFile(t, "README.md"), "## What you get")

	diagrams := mermaidFence.FindAllStringSubmatch(got, -1)
	if len(diagrams) != 2 {
		t.Fatalf("What you get carries %d mermaid diagrams, want 2 — delivery and flow", len(diagrams))
	}

	stopsNamed := false
	for _, d := range diagrams {
		body := flat(d[1])
		if strings.Contains(body, "spec") && strings.Contains(body, "plan") && strings.Contains(body, "push") {
			stopsNamed = true
		}
	}
	if !stopsNamed {
		t.Error("neither diagram names the three stops — spec, plan and push")
	}
}
