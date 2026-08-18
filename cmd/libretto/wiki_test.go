package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
- If the total goes negative, then the system shall clamp it to zero.
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

func runWiki(t *testing.T, dir string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	err := wiki(&out, dir)
	return out.String(), err
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
		"shall clamp it to zero",
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
