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
		// Four more, moved out when the README was cut from 389 lines to fit under its own
		// ceiling. Each anchor is the phrase a paraphrase would have to lose — the lesson
		// above about "module cache" is why they are chosen that way and not by convenience.
		{"the spec stop is the cheap place to disagree", "costs a line"},
		{"phase 1's source order is the point", "gets abandoned"},
		{"attacca will not answer a gate", "never merges, tags or releases"},
		{"an alias resolves per provider from the environment", "no request, no"},
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

	// Resolving the links it has says nothing about the one it must have: a README with no
	// CONTRIBUTING.md link at all passes the loop below. The contributor's door is reached
	// from the front door or it is reached by nobody browsing.
	if !strings.Contains(flat(repoFile(t, "README.md")), "CONTRIBUTING.md") {
		t.Error("the README never links to CONTRIBUTING.md — the contributor's door is unreachable from the front")
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

// A front door that grows never shrinks back, and this one had reached 389 lines with
// arguments in it that docs/ already owned. The capability's own ceilings proposed this
// mechanism for themselves before it was needed: a criterion about how long the file may get.
//
// It is a ratchet against growth, not a measure of readability — a README rewritten as 340
// very long lines passes. Readability is the one thing no test here holds, which the
// capability records and which is why that is said out loud rather than implied.
func TestReadmeStaysShort(t *testing.T) {
	// 380, and the first draft said 340 — chosen before anything was measured. Removing one
	// duplicate and relocating five arguments netted 13 lines, not 49: the relocations swap
	// prose for a pointer, so each one saves less than it reads like it should. What remains is
	// reference — the commands table, the model and effort listings, the five states, the
	// environment table — and cutting those to hit a number would leave a worse README that
	// passes its own guard, which is the failure this whole capability exists to prevent.
	//
	// So this is a ratchet against growth, deliberately tight. Raising it is the wrong fix; the
	// right one is moving an argument out, which is what the message says.
	const ceiling = 380
	if n := strings.Count(repoFile(t, "README.md"), "\n"); n > ceiling {
		t.Errorf("README.md is %d lines, over its %d ceiling — move an argument to docs/ rather than raising this", n, ceiling)
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

// THIRD-PARTY.md is the record that every vendored skill carries its own licence and
// version, and check-payload already parses its table to derive the vendored list — so the
// file was load-bearing for a gate while no capability's Governs: claimed it. Moving the
// three licence texts into licenses/ is exactly the change that breaks its links, and
// nothing was watching.
//
// The scan runs over flat(), like every substring assertion in this file: two guards here
// have already shipped unable to fire because Markdown wrapped between the words they
// looked for. A third would not be an accident.
func TestThirdPartyLinksResolve(t *testing.T) {
	links := readmeLink.FindAllStringSubmatch(flat(repoFile(t, "THIRD-PARTY.md")), -1)
	if len(links) == 0 {
		t.Fatal("no links matched — the pattern is broken, not THIRD-PARTY.md")
	}

	for _, match := range links {
		target := strings.TrimSpace(match[1])
		if strings.HasPrefix(target, "http") || strings.HasPrefix(target, "mailto:") {
			continue
		}
		if _, err := os.Stat(filepath.Join("..", "..", target)); err != nil {
			t.Errorf("THIRD-PARTY.md links to %s, which does not exist", target)
		}
	}

	// Checking the links alone would pass with one licence moved and two forgotten, since a
	// link nobody updated points at a file that is still there. So the layout is asserted
	// directly: three under licenses/, none left at the root.
	for _, name := range []string{"LICENSE-caveman", "LICENSE-ponytail", "LICENSE-superpowers"} {
		if _, err := os.Stat(filepath.Join("..", "..", "licenses", name)); err != nil {
			t.Errorf("licenses/%s does not exist — a vendored copy must carry its licence text", name)
		}
		if _, err := os.Stat(filepath.Join("..", "..", name)); err == nil {
			t.Errorf("%s is still at the root — LICENSE is the only licence file that belongs there", name)
		}
	}

	// The loop above says which licence files must leave the root and would be satisfied by a
	// root with no licence at all. GitHub reads root LICENSE to display the repository's
	// licence and to fill the API field, so the tidying this test enforces must not take it:
	// the 6→7 reviewer confirmed it by inspection and asked what would catch a regression.
	if _, err := os.Stat(filepath.Join("..", "..", "LICENSE")); err != nil {
		t.Error("LICENSE is not at the root — GitHub reads it there to display the licence")
	}
}

// CONTRIBUTING.md is a door, not a copy. GitHub links it from the issue and PR forms, which
// is the whole reason it beats a section in AGENTS.md — it arrives when the contributor needs
// it. The risk is not the missing file, it is the duplicate: AGENTS.md already carries the
// gates, the commit convention and the label rule, and two files kept in sync is two sources
// of truth where the one that wins is the one nobody edited.
//
// So the guard holds both halves — what only the guide says, and what it must not repeat.
// Absent means absent, with no "unless the line links to AGENTS.md" escape: the scan runs
// over flat(), where there are no lines to reason about.
func TestContributingIsADoorNotACopy(t *testing.T) {
	guide := flat(repoFile(t, "CONTRIBUTING.md"))

	// What lives nowhere else, and what an outside contributor cannot guess.
	for _, must := range []string{
		"release:patch", // the designed refusal that reads as a broken pipeline unpredicted
		".agents/specs/",
		"1.0.0",
		// The 6→7 reviewer measured that deleting this whole section left the guard at zero
		// failures: the clause was in the criterion and in nothing that could fail. The
		// heading is the anchor rather than a sentence, so improving the prose does not
		// break the guard.
		"Work does not come from a tracker",
	} {
		if !strings.Contains(guide, must) {
			t.Errorf("CONTRIBUTING.md never mentions %q — it is one of the few things only this file says", must)
		}
	}

	// All six, not one of six. The reviewer removed five gate lines from a copy and the guard
	// stayed green — so the contributor's paste-and-run block could lose five of the six
	// commands undetected, which is the whole value of that section.
	for _, gate := range []string{
		"gofmt -l .",
		"go vet ./...",
		"go test ./... -count=1",
		"scripts/check-payload",
		"spec-drift --self-test",
		"spec-drift --anchors",
	} {
		if !strings.Contains(guide, gate) {
			t.Errorf("CONTRIBUTING.md omits the gate %q — all six pass or the change is not ready", gate)
		}
	}

	// Rules AGENTS.md owns. Restating one here is the failure this file exists to prevent.
	for _, owned := range []string{
		"Co-Authored-By",
		"type(scope): subject",
		"72 chars",
		"stdlib, then a native",
	} {
		if strings.Contains(guide, owned) {
			t.Errorf("CONTRIBUTING.md restates %q, which AGENTS.md owns — link to it instead", owned)
		}
	}

	links := readmeLink.FindAllStringSubmatch(guide, -1)
	if len(links) == 0 {
		t.Fatal("no links matched — a door made of no links is not a door")
	}
	for _, match := range links {
		target := strings.TrimSpace(match[1])
		if strings.HasPrefix(target, "http") || strings.HasPrefix(target, "mailto:") {
			continue
		}
		if _, err := os.Stat(filepath.Join("..", "..", target)); err != nil {
			t.Errorf("CONTRIBUTING.md links to %s, which does not exist", target)
		}
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
