package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
	"github.com/pausf/libretto-automata/internal/ui"
)

// A full-shaped capability and a legacy one. The legacy fixture is the point:
// specs written before EARS carry prose criteria with no Proof: bullets, and the
// wiki renders what is there instead of erroring on what is not.
const fullSpec = `# Pricing

Governs: src/pricing/** src/cart/total.go

Relative discounts across bundles, and the rounding they promise.
Second sentence nobody quotes.

## Verification criteria

- When a bundle holds two items, the system shall apply the relative discount.
  Proof: src/pricing/bundle_test.go TestRelativeDiscount
- If the total goes negative, then the system shall **clamp** it to ` + "`zero`" + `.
  Proof: src/pricing/bundle_test.go TestClampAtZero
`

const legacySpec = `# Checkout

Governs: src/checkout/**

## Notes

Prose only, written before the syntax existed.
`

func writeSpecs(t *testing.T, root, dir string) string {
	t.Helper()
	specs := filepath.Join(root, dir)
	for name, body := range map[string]string{"pricing": fullSpec, "checkout": legacySpec} {
		if err := os.MkdirAll(filepath.Join(specs, name), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(specs, name, "spec.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return specs
}

func runWiki(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	err := wiki(&out, dir, args)
	return out.String(), err
}

func wikiHTML(t *testing.T, specs string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(specs, "wiki.html"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func readme(t *testing.T, specs string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(specs, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestWikiDiscoversSpecsDirInDriftOrder(t *testing.T) {
	// Populated back to front: each earlier candidate, once added, must win over
	// every later one — the same first-hit-wins walk spec-drift does.
	order := []string{".agents/specs", "specs", "openspec", "docs/specs", "spec"}
	for i := range order {
		dir := t.TempDir()
		for _, later := range order[i:] { // candidate i and everything after it exist
			writeSpecs(t, dir, later)
		}
		if _, err := runWiki(t, dir); err != nil {
			t.Fatalf("with %s present: %v", order[i], err)
		}
		for j, candidate := range order[i:] {
			_, err := os.Stat(filepath.Join(dir, candidate, "README.md"))
			if j == 0 && err != nil {
				t.Fatalf("expected the wiki in %s, not found", candidate)
			}
			if j > 0 && err == nil {
				t.Fatalf("%s present but the wiki also landed in later %s", order[i], candidate)
			}
		}
	}
}

func TestWikiWritesIndexAndSections(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	out, err := runWiki(t, dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "README.md") || !strings.Contains(out, "2 capabilities") {
		t.Fatalf("report line missing or wrong: %q", out)
	}
	got := readme(t, specs)
	for _, want := range []string{
		"| Capability | Governs |",
		"[checkout](checkout/spec.md)",
		"[pricing](pricing/spec.md)",
		"## pricing",
		"`src/pricing/** src/cart/total.go`",
		"Relative discounts across bundles, and the rounding they promise.",
		"shall apply the relative discount",
		"shall **clamp** it to `zero`", // README keeps the raw inline markdown
		"## checkout",
		"`src/checkout/**`",
		"[full spec](pricing/spec.md)",
		"[full spec](checkout/spec.md)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("wiki missing %q", want)
		}
	}
	if strings.Contains(got, "Proof:") {
		t.Error("criteria carried their Proof: lines into the wiki")
	}
	// The legacy capability has no Proof:-backed bullets and no intro before its
	// first heading beyond none — absent renders as absent, never as an error.
	if strings.Contains(got, "Prose only") {
		t.Error("prose under a later heading leaked in as the intro")
	}
}

func TestGeneratedReadmeCarriesTheMarker(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, "specs")
	if _, err := runWiki(t, dir); err != nil {
		t.Fatal(err)
	}
	got := readme(t, specs)
	first, _, _ := strings.Cut(got, "\n")
	if first != wikiMarker {
		t.Fatalf("first line is %q, not the marker", first)
	}
	if !strings.Contains(wikiMarker, "libretto wiki") {
		t.Error("the marker does not name the command that refreshes the file")
	}
}

func TestWikiNeverOverwritesAHandWrittenReadme(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	hand := "# My own index\n\nWritten by a person.\n"
	if err := os.WriteFile(filepath.Join(specs, "README.md"), []byte(hand), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir); err == nil {
		t.Fatal("expected a refusal, got success")
	}
	if got := readme(t, specs); got != hand {
		t.Fatalf("the hand-written README was modified:\n%s", got)
	}
	// A marked one is ours and regenerates.
	if err := os.WriteFile(filepath.Join(specs, "README.md"), []byte(wikiMarker+"\nstale\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir); err != nil {
		t.Fatalf("refused to refresh a file it generated: %v", err)
	}
}

func TestWikiReportsNoSpecsAndExitsNonZero(t *testing.T) {
	dir := t.TempDir()
	_, err := runWiki(t, dir)
	if err == nil || !strings.Contains(err.Error(), "no specs directory") {
		t.Fatalf("missing directory: want a one-line naming it, got %v", err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err = runWiki(t, dir)
	if err == nil || !strings.Contains(err.Error(), "holds no") {
		t.Fatalf("empty directory: want the emptiness named, got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "specs", "README.md")); statErr == nil {
		t.Fatal("an empty specification still got a README written")
	}
}

func TestWikiOutputIsDeterministic(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if _, err := runWiki(t, dir); err != nil {
		t.Fatal(err)
	}
	first := readme(t, specs)
	if _, err := runWiki(t, dir); err != nil {
		t.Fatal(err)
	}
	if second := readme(t, specs); second != first {
		t.Fatal("two runs over unchanged input produced different bytes")
	}
}

func TestWikiWritesNothingButTheReadme(t *testing.T) {
	dir := t.TempDir()
	writeSpecs(t, dir, ".agents/specs")
	before := treeSnapshot(t, dir)
	if _, err := runWiki(t, dir); err != nil {
		t.Fatal(err)
	}
	after := treeSnapshot(t, dir)
	var extra []string
	for p := range after {
		if _, ok := before[p]; !ok {
			extra = append(extra, p)
		}
	}
	if len(extra) != 1 || !strings.HasSuffix(extra[0], filepath.Join(".agents", "specs", "README.md")) {
		t.Fatalf("expected exactly the one README, new files: %v", extra)
	}
	for p, was := range before {
		if after[p] != was {
			t.Errorf("pre-existing file changed: %s", p)
		}
	}
}

func treeSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	snap := map[string]string{}
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		snap[p] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snap
}

func TestWikiHTMLWritesTheViewer(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	out, err := runWiki(t, dir, "--html")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "wiki.html") {
		t.Fatalf("report line missing the file: %q", out)
	}
	got := wikiHTML(t, specs)
	for _, want := range []string{
		`href="#pricing"`, // nav entry
		`href="#checkout"`,
		`id="pricing"`, // section
		"Relative discounts across bundles",
		"src/pricing/** src/cart/total.go",
		"shall apply the relative discount",
		"<strong>clamp</strong>", // bold conversion, after escaping
		"<code>zero</code>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("viewer missing %q", want)
		}
	}
	if _, err := os.Stat(filepath.Join(specs, "README.md")); err == nil {
		t.Error("--html run wrote README.md too")
	}
}

func TestWikiHTMLCarriesTheMarker(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, "specs")
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	first, _, _ := strings.Cut(wikiHTML(t, specs), "\n")
	if first != wikiHTMLMarker {
		t.Fatalf("first line is %q, not the HTML marker", first)
	}
	if !strings.Contains(wikiHTMLMarker, "libretto wiki --html") {
		t.Error("the marker does not name the refresh command")
	}
}

func TestWikiHTMLNeverOverwritesAForeignFile(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	hand := "<html>somebody's page</html>\n"
	if err := os.WriteFile(filepath.Join(specs, "wiki.html"), []byte(hand), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir, "--html"); err == nil {
		t.Fatal("expected a refusal, got success")
	}
	if got := wikiHTML(t, specs); got != hand {
		t.Fatalf("the foreign wiki.html was modified:\n%s", got)
	}
	if err := os.WriteFile(filepath.Join(specs, "wiki.html"), []byte(wikiHTMLMarker+"\nstale\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatalf("refused to refresh a file it generated: %v", err)
	}
}

func TestWikiHTMLEscapesSpecContent(t *testing.T) {
	dir := t.TempDir()
	specs := filepath.Join(dir, ".agents/specs", "sneaky")
	if err := os.MkdirAll(specs, 0o755); err != nil {
		t.Fatal(err)
	}
	spec := "# Sneaky\n\nGoverns: src/**\n\nIntro with <script>alert(1)</script> inside.\n\n" +
		"## Verification criteria\n\n- When poked, the system shall <script>alert(2)</script> not run it.\n  Proof: a_test.go TestX\n"
	if err := os.WriteFile(filepath.Join(specs, "spec.md"), []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, filepath.Join(dir, ".agents/specs"))
	if strings.Contains(got, "<script>alert(") {
		t.Fatal("spec content reached markup position unescaped")
	}
	if !strings.Contains(got, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Error("escaped intro text missing from the page")
	}
}

func TestWikiHTMLIsSelfContained(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	for _, m := range regexp.MustCompile(`(?:src|href)="(https?://[^"]+)"`).FindAllStringSubmatch(got, -1) {
		u := m[1]
		if !strings.HasPrefix(u, "https://fonts.googleapis.com") && !strings.HasPrefix(u, "https://fonts.gstatic.com") {
			t.Errorf("external reference beyond the font hosts: %s", u)
		}
	}
	if strings.Contains(got, "<script src=") {
		t.Error("external script referenced; the filter must be inline")
	}
}

func TestWikiHTMLIsDeterministic(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	first := wikiHTML(t, specs)
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	if second := wikiHTML(t, specs); second != first {
		t.Fatal("two --html runs over unchanged input produced different bytes")
	}
}

func TestPlainWikiRefreshesAMarkedHTMLView(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if err := os.WriteFile(filepath.Join(specs, "wiki.html"), []byte(wikiHTMLMarker+"\nstale\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir); err != nil {
		t.Fatal(err)
	}
	if got := wikiHTML(t, specs); strings.Contains(got, "stale") {
		t.Fatal("plain run left the marked wiki.html stale")
	}
	// A foreign wiki.html is left alone and the plain run still succeeds —
	// erroring here would block every landing regeneration in the project.
	foreign := "<html>not ours</html>\n"
	if err := os.WriteFile(filepath.Join(specs, "wiki.html"), []byte(foreign), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir); err != nil {
		t.Fatalf("plain run failed on a foreign wiki.html: %v", err)
	}
	if got := wikiHTML(t, specs); got != foreign {
		t.Fatal("plain run touched a foreign wiki.html")
	}
}

func TestWikiHTMLWritesNothingButTheOneFile(t *testing.T) {
	dir := t.TempDir()
	writeSpecs(t, dir, ".agents/specs")
	before := treeSnapshot(t, dir)
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	after := treeSnapshot(t, dir)
	var extra []string
	for p := range after {
		if _, ok := before[p]; !ok {
			extra = append(extra, p)
		}
	}
	if len(extra) != 1 || !strings.HasSuffix(extra[0], filepath.Join(".agents", "specs", "wiki.html")) {
		t.Fatalf("expected exactly the one wiki.html, new files: %v", extra)
	}
	for p, was := range before {
		if after[p] != was {
			t.Errorf("pre-existing file changed: %s", p)
		}
	}
}

func TestWikiRejectsAnUnknownFlag(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if _, err := runWiki(t, dir, "--bogus"); err == nil || !strings.Contains(err.Error(), "--bogus") {
		t.Fatalf("want the unknown argument named in an error, got %v", err)
	}
	for _, f := range []string{"README.md", "wiki.html"} {
		if _, err := os.Stat(filepath.Join(specs, f)); err == nil {
			t.Errorf("a rejected invocation still wrote %s", f)
		}
	}
}

func TestPanelOffersWikiOnlyInAProjectWithSpecs(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	writeSpecs(t, f.Project, ".agents/specs")

	wikiRow := func(scope target.Scope) *ui.MenuItem {
		menu, _, err := panelData(f.Repo, f.Project, target.ClaudeTool, scope)
		if err != nil {
			t.Fatal(err)
		}
		for i := range menu {
			if menu[i].Label == "wiki" {
				return &menu[i]
			}
		}
		return nil
	}

	row := wikiRow(target.ProjectScope)
	if row == nil {
		t.Fatal("project scope with specs: the menu does not offer wiki")
	}
	if !row.Enabled {
		t.Error("the wiki row is offered but disabled")
	}
	if row.Destructive {
		t.Error("rendering a wiki is not destructive")
	}
	if wikiRow(target.GlobalScope) != nil {
		t.Fatal("global scope offers wiki — the user drew that line explicitly")
	}

	bare := newFixture(t)
	bare.skill(t, "alpha")
	menu, _, err := panelData(bare.Repo, bare.Project, target.ClaudeTool, target.ProjectScope)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range menu {
		if item.Label == "wiki" {
			t.Fatal("a project with no specs directory still offers wiki")
		}
	}
}

// recordOpener swaps the browser seam for a recorder; no test launches anything.
func recordOpener(t *testing.T) *[]string {
	t.Helper()
	var calls []string
	prev := openViewer
	openViewer = func(path string) error { calls = append(calls, path); return nil }
	t.Cleanup(func() { openViewer = prev })
	return &calls
}

func TestPanelWikiRowOpensTheViewer(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	specs := writeSpecs(t, f.Project, ".agents/specs")
	calls := recordOpener(t)

	if _, err := runCaptured("wiki", f.Repo, f.Project, f.project(), false); err != nil {
		t.Fatalf("dispatching wiki failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(specs, "wiki.html")); err != nil {
		t.Fatal("dispatch reported no error and generated no viewer")
	}
	if len(*calls) != 1 || (*calls)[0] != filepath.Join(specs, "wiki.html") {
		t.Fatalf("the row did not hand the viewer to the opener: %v", *calls)
	}

	menu, _, err := panelData(f.Repo, f.Project, target.ClaudeTool, target.ProjectScope)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range menu {
		if item.Label == "wiki" && !strings.Contains(item.Desc, "open") {
			t.Errorf("the row's description does not say the viewer opens: %q", item.Desc)
		}
	}
}

func TestWikiOpenGeneratesAndOpensTheViewer(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	calls := recordOpener(t)

	out, err := runWiki(t, dir, "--open")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(specs, "wiki.html")
	if _, err := os.Stat(want); err != nil {
		t.Fatal("--open did not write the viewer")
	}
	if len(*calls) != 1 || (*calls)[0] != want {
		t.Fatalf("opener got %v, want the written path", *calls)
	}
	if !strings.Contains(out, "opened") {
		t.Errorf("the report does not say it opened: %q", out)
	}

	// The two spellings are one behaviour.
	if _, err := runWiki(t, dir, "--html", "--open"); err != nil {
		t.Fatalf("--html --open: %v", err)
	}

	// An opener error after a successful write surfaces non-zero, naming the path
	// so the user can open it by hand — the file stays.
	openViewer = func(string) error { return os.ErrPermission }
	if _, err := runWiki(t, dir, "--open"); err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("opener failure not surfaced with the written path: %v", err)
	}
}

func TestWikiOpenDoesNotOpenOnFailure(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	calls := recordOpener(t)
	foreign := "<html>not ours</html>\n"
	if err := os.WriteFile(filepath.Join(specs, "wiki.html"), []byte(foreign), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir, "--open"); err == nil {
		t.Fatal("a refused write still reported success")
	}
	if len(*calls) != 0 {
		t.Fatalf("a refused write still opened a browser: %v", *calls)
	}
}

func TestOpenerArgvPerPlatform(t *testing.T) {
	if got := openerArgv("darwin", "/p/wiki.html"); len(got) != 2 || got[0] != "open" || got[1] != "/p/wiki.html" {
		t.Fatalf("darwin argv: %v", got)
	}
	for _, goos := range []string{"linux", "freebsd", "windows"} {
		if got := openerArgv(goos, "/p/wiki.html"); len(got) != 2 || got[0] != "xdg-open" {
			t.Fatalf("%s argv: %v", goos, got)
		}
	}
}

func TestWikiHTMLMotionRespectsReducedMotion(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	if strings.Contains(got, "addEventListener('scroll'") || strings.Contains(got, `addEventListener("scroll"`) {
		t.Fatal("scroll listeners present; motion must be CSS scroll-driven")
	}
	motionStart := strings.Index(got, "@media (prefers-reduced-motion: no-preference)")
	motionEnd := strings.Index(got, "/*end-motion*/")
	if motionStart < 0 || motionEnd < motionStart {
		t.Fatal("the delimited reduced-motion block is missing")
	}
	inside := got[motionStart:motionEnd]
	if !strings.Contains(inside, "animation-timeline") {
		t.Fatal("no scroll-driven animation inside the motion block")
	}
	outside := got[:motionStart] + got[motionEnd:]
	if strings.Contains(outside, "animation-timeline") || strings.Contains(outside, "@keyframes") {
		t.Fatal("motion declared outside the reduced-motion guard")
	}
}

func TestWikiHTMLThemesAreTokenComplete(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)

	block := func(start string) string {
		i := strings.Index(got, start)
		if i < 0 {
			t.Fatalf("token block %q missing", start)
		}
		j := strings.Index(got[i:], "}")
		return got[i : i+j]
	}
	dark := block(":root{")
	light := block("@media (prefers-color-scheme: light){:root{")

	// The same token set in both — a token defined once renders one theme's
	// text on the other theme's ground.
	tokenRe := regexp.MustCompile(`--[a-z-]+:`)
	darkTokens := tokenRe.FindAllString(dark, -1)
	if len(darkTokens) == 0 {
		t.Fatal("no tokens on the bare :root")
	}
	// Names, not counts: a token renamed in one block would leave a var()
	// resolving in only one theme while the tallies still match.
	lightSet := map[string]bool{}
	for _, tok := range tokenRe.FindAllString(light, -1) {
		lightSet[tok] = true
	}
	for _, tok := range darkTokens {
		if !lightSet[tok] {
			t.Errorf("token %q defined in dark, missing from the light block", tok)
		}
	}

	// No hex literal outside the two token blocks: components take colour
	// only through var().
	styleStart := strings.Index(got, "<style>")
	styleEnd := strings.Index(got, "</style>")
	css := got[styleStart:styleEnd]
	for _, m := range regexp.MustCompile(`#[0-9A-Fa-f]{3,8}\b`).FindAllStringIndex(css, -1) {
		hex := css[m[0]:m[1]]
		if !strings.Contains(dark, hex) && !strings.Contains(light, hex) {
			t.Errorf("hex literal %q outside the token blocks", hex)
		}
	}
}

func TestWikiHTMLIsHomeAndPages(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	for _, want := range []string{
		`<section id="home"`,
		`class="card" href="#pricing"`,
		`class="card" href="#checkout"`, // zero criteria still gets a card
		`style="width:100%"`,            // pricing: 2 of max 2
		`style="width:0%"`,              // checkout: zero, drawn at zero
		`<article class="cap" id="pricing"`,
		`<article class="cap" id="checkout"`, // and a page like every other
		`href="#home"`,                       // the way back
		// The page carries the contract: Governs, intro, criteria — asserted
		// here so this criterion's own citation covers its whole clause.
		`<code>src/pricing/** src/cart/total.go</code>`,
		"Relative discounts across bundles",
		"shall apply the relative discount",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("home/pages missing %q", want)
		}
	}
	if strings.Index(got, `id="home"`) > strings.Index(got, `<article`) {
		t.Error("home is not first in document order — a JS-less render opens mid-wiki")
	}
}

func TestWikiDatesComeFromGitAndDegrade(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	prev := wikiGitDate
	t.Cleanup(func() { wikiGitDate = prev })
	wikiGitDate = func(_, specPath string) string {
		if strings.Contains(specPath, "pricing") {
			return "2026-08-18"
		}
		return "" // checkout: no history — absent, never an error
	}
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	if n := strings.Count(got, "2026-08-18"); n < 2 {
		t.Fatalf("the date should reach the card and the page, found %d occurrence(s)", n)
	}
	first := got
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	if second := wikiHTML(t, specs); second != first {
		t.Fatal("injected dates made two runs differ — determinism lost")
	}
	// Full degrade: the seam yields nothing anywhere and the run still succeeds.
	wikiGitDate = func(_, _ string) string { return "" }
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatalf("a dateless project failed the run: %v", err)
	}
	if strings.Contains(wikiHTML(t, specs), "2026-08-18") {
		t.Fatal("a stale date survived the dateless rerun")
	}
}

func TestWikiHTMLRouterIsAnEnhancement(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	for _, want := range []string{"hashchange", "paged"} {
		if !strings.Contains(got, want) {
			t.Errorf("router missing %q", want)
		}
	}
	if strings.Contains(got, "<script src=") {
		t.Fatal("router must be inline")
	}
	// Visibility is class-toggled only: no display juggling from JS.
	if strings.Contains(got, ".style.display") {
		t.Fatal("router sets styles directly instead of toggling classes")
	}
}

func TestWikiHTMLHomeSearchIsInline(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	if !strings.Contains(got, `id="filter"`) {
		t.Fatal("the home search input is missing")
	}
	if !strings.Contains(got, `data-crit="`) {
		t.Fatal("cards carry no data-crit — the search has nothing to match criteria against")
	}
	if !strings.Contains(got, "shall apply the relative discount") {
		t.Fatal("criteria text missing from the page")
	}
	if !strings.Contains(got, "data-crit") || !strings.Contains(got, "dataset.crit") && !strings.Contains(got, "getAttribute('data-crit')") {
		t.Error("the search script does not read data-crit")
	}
}

// injectWikiGit swaps all three git seams; tests never touch real git state.
func injectWikiGit(t *testing.T, date func(string, string) string, subject func(string, string) string, tracked func(string) []string) {
	t.Helper()
	pd, ps, pt := wikiGitDate, wikiGitSubject, wikiGitTracked
	t.Cleanup(func() { wikiGitDate, wikiGitSubject, wikiGitTracked = pd, ps, pt })
	if date != nil {
		wikiGitDate = date
	}
	if subject != nil {
		wikiGitSubject = subject
	}
	if tracked != nil {
		wikiGitTracked = tracked
	}
}

// healthFixture: "solid" is fully EARS with resolving proofs; "wobbly" carries one
// prose criterion and one EARS criterion citing a file that does not exist.
func healthFixture(t *testing.T, dir string) string {
	t.Helper()
	specs := filepath.Join(dir, ".agents", "specs")
	solid := "# Solid\n\nGoverns: src/**\n\nHolds the line.\n\n## Verification criteria\n\n" +
		"- The system shall hold.\n  Proof: pkg/solid_test.go TestHold\n" +
		"- When poked, the system shall answer.\n  Proof: scripts/check\n"
	wobbly := "# Wobbly\n\nGoverns: tools/one.go\n\n## Verification criteria\n\n" +
		"- the panel opens where it was left\n  Proof: scripts/check\n" +
		"- The system shall wobble.\n  Proof: pkg/gone_test.go TestGone\n"
	for name, body := range map[string]string{"solid": solid, "wobbly": wobbly} {
		if err := os.MkdirAll(filepath.Join(specs, name), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(specs, name, "spec.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(dir, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "pkg", "solid_test.go"), []byte("package pkg\n\nfunc TestHold(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "scripts", "check"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return specs
}

func TestWikiHomeCarriesTheRecentRail(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	injectWikiGit(t,
		func(_, p string) string {
			if strings.Contains(p, "pricing") {
				return "2026-08-18"
			}
			return "2026-08-10"
		},
		func(_, p string) string {
			if strings.Contains(p, "pricing") {
				return "land the discounts <finally>"
			}
			return "first cut"
		}, nil)
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	for _, want := range []string{`class="recent"`, "land the discounts &lt;finally&gt;", "2026-08-18", "first cut"} {
		if !strings.Contains(got, want) {
			t.Errorf("rail missing %q", want)
		}
	}
	if strings.Index(got, "land the discounts") > strings.Index(got, "first cut") {
		t.Error("the rail is not most-recent first")
	}
	// No dates anywhere → no rail at all.
	dir2 := t.TempDir()
	specs2 := writeSpecs(t, dir2, ".agents/specs")
	injectWikiGit(t, func(_, _ string) string { return "" }, func(_, _ string) string { return "" }, nil)
	if _, err := runWiki(t, dir2, "--html"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(wikiHTML(t, specs2), `class="recent"`) {
		t.Fatal("a dateless project still renders the rail")
	}
}

func TestWikiHomeCarriesTheInFlightStrip(t *testing.T) {
	dir := t.TempDir()
	specs := writeSpecs(t, dir, ".agents/specs")
	ch := filepath.Join(dir, ".agents", "changes")
	if err := os.MkdirAll(filepath.Join(ch, "add-discounts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ch, "add-discounts", "tasks.md"),
		[]byte("# Tasks\n\n- [x] one\n- [x] two\n- [ ] three\n- [ ] four\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A landed-shape change (all closed) never shows, and a legacy plan.md counts.
	if err := os.MkdirAll(filepath.Join(ch, "old-style"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ch, "old-style", "plan.md"),
		[]byte("- [ ] only\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(ch, "an-idea"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ch, "an-idea", "proposal.md"),
		[]byte("# an-idea\n\nQueued: 2026-08-14\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	for _, want := range []string{`class="inflight"`, "add-discounts", "2/4", "old-style", "0/1", "1 queued"} {
		if !strings.Contains(got, want) {
			t.Errorf("strip missing %q", want)
		}
	}
	// Queue alone does not summon the strip — the recorded assumption.
	dir2 := t.TempDir()
	specs2 := writeSpecs(t, dir2, ".agents/specs")
	if err := os.MkdirAll(filepath.Join(dir2, ".agents", "changes", "an-idea"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir2, ".agents", "changes", "an-idea", "proposal.md"),
		[]byte("Queued: 2026-08-14\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runWiki(t, dir2, "--html"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(wikiHTML(t, specs2), `class="inflight"`) {
		t.Fatal("a queue alone summoned the strip")
	}
}

func TestWikiHomeMeasuresContractHealth(t *testing.T) {
	dir := t.TempDir()
	specs := healthFixture(t, dir)
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	// 3 of 4 criteria carry shall: green 75, amber the remainder.
	for _, want := range []string{`class="health"`, "width:75%", "width:25%", "1 unproven"} {
		if !strings.Contains(got, want) {
			t.Errorf("health bar missing %q", want)
		}
	}
}

func TestWikiCardsCarryTheHealthDot(t *testing.T) {
	dir := t.TempDir()
	specs := healthFixture(t, dir)
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	if !regexp.MustCompile(`href="#solid"[^>]*data-health="ok"`).MatchString(got) {
		t.Error("solid's card is not marked healthy")
	}
	if !regexp.MustCompile(`href="#wobbly"[^>]*data-health="warn"`).MatchString(got) {
		t.Error("wobbly's card is not marked amber")
	}
}

func TestWikiFooterMeasuresGovernedTree(t *testing.T) {
	dir := t.TempDir()
	specs := healthFixture(t, dir)
	injectWikiGit(t, nil, nil, func(string) []string {
		return []string{"src/a/b.go", "src/c.go", "tools/one.go", "README.md", "x/y.go"}
	})
	if _, err := runWiki(t, dir, "--html"); err != nil {
		t.Fatal(err)
	}
	got := wikiHTML(t, specs)
	// src/** must cross a directory boundary — the ** arm the force-red breaks.
	for _, want := range []string{`class="governed"`, "3 governed", "2 orphan"} {
		if !strings.Contains(got, want) {
			t.Errorf("footer missing %q", want)
		}
	}
	// Git unavailable → no footer.
	dir2 := t.TempDir()
	specs2 := healthFixture(t, dir2)
	injectWikiGit(t, nil, nil, func(string) []string { return nil })
	if _, err := runWiki(t, dir2, "--html"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(wikiHTML(t, specs2), `class="governed"`) {
		t.Fatal("a gitless project still renders the footer")
	}
}
