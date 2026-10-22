package dist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// The repository is its own Claude Code marketplace holding one plugin rooted at "./". These
// read the manifests as the repository ships them. Decoding into map[string]any rather than a
// struct is what makes "carries no version key" testable: a struct without the field would
// decode any file and pass.
func manifest(t *testing.T, name string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".claude-plugin", name))
	if err != nil {
		t.Fatalf("reading .claude-plugin/%s: %v", name, err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decoding .claude-plugin/%s: %v", name, err)
	}
	return m
}

func marketplaceEntries(t *testing.T) []map[string]any {
	t.Helper()
	raw, ok := manifest(t, "marketplace.json")["plugins"].([]any)
	if !ok {
		t.Fatal("marketplace.json has no plugins array")
	}
	var entries []map[string]any
	for _, e := range raw {
		entry, ok := e.(map[string]any)
		if !ok {
			t.Fatalf("marketplace.json plugin entry is not an object: %v", e)
		}
		entries = append(entries, entry)
	}
	return entries
}

// The plugin name is the namespace every component is listed under, so it is the short one
// the user already types as the binary: /libretto:libretto-flow.
func TestPluginManifestNamesLibretto(t *testing.T) {
	if got := manifest(t, "plugin.json")["name"]; got != "libretto" {
		t.Errorf("plugin.json name = %v, want libretto", got)
	}
}

func TestMarketplaceIsNamedAfterTheRepository(t *testing.T) {
	if got := manifest(t, "marketplace.json")["name"]; got != "libretto-automata" {
		t.Errorf("marketplace.json name = %v, want libretto-automata", got)
	}
}

func TestMarketplaceListsThePluginAtTheRoot(t *testing.T) {
	entries := marketplaceEntries(t)
	if len(entries) != 1 {
		t.Fatalf("marketplace.json lists %d plugins, want exactly 1", len(entries))
	}
	if got := entries[0]["name"]; got != "libretto" {
		t.Errorf("entry name = %v, want libretto", got)
	}
	if got := entries[0]["source"]; got != "./" {
		t.Errorf("entry source = %v, want ./", got)
	}
}

// Claude Code derives the version from the installed commit when neither manifest pins one.
// A pinned string is a second source of the version, and it freezes every user on the cached
// copy until somebody edits it.
func TestPluginCarriesNoVersion(t *testing.T) {
	if _, ok := manifest(t, "plugin.json")["version"]; ok {
		t.Error("plugin.json carries a version key")
	}
	for _, entry := range marketplaceEntries(t) {
		if _, ok := entry["version"]; ok {
			t.Errorf("marketplace entry %v carries a version key", entry["name"])
		}
	}
}

// The default directories are what loads, so a skill added to skills/ ships with nothing
// else edited. A component list would be a second place every new item must be added.
func TestPluginManifestDeclaresNoComponentPaths(t *testing.T) {
	m := manifest(t, "plugin.json")
	for _, key := range []string{"skills", "agents", "commands"} {
		if _, ok := m[key]; ok {
			t.Errorf("plugin.json declares %q", key)
		}
	}
}

// The vendored skills ship inside the plugin, so their licences must too. This reads the
// repository tree under the entry's source — never Claude Code's cache copy.
func TestPluginRootCarriesTheVendoredLicences(t *testing.T) {
	entries := marketplaceEntries(t)
	if len(entries) != 1 {
		t.Fatalf("marketplace.json lists %d plugins, want exactly 1", len(entries))
	}
	source, _ := entries[0]["source"].(string)
	root := filepath.Join("..", "..", source)

	notice, err := os.ReadFile(filepath.Join(root, "THIRD-PARTY.md"))
	if err != nil {
		t.Fatalf("the plugin root carries no THIRD-PARTY.md: %v", err)
	}
	links := regexp.MustCompile(`\]\((licenses/[^)]+)\)`).FindAllStringSubmatch(string(notice), -1)
	if len(links) == 0 {
		t.Fatal("THIRD-PARTY.md links nothing under licenses/")
	}
	for _, l := range links {
		if _, err := os.Stat(filepath.Join(root, l[1])); err != nil {
			t.Errorf("THIRD-PARTY.md links %s, missing under the plugin root: %v", l[1], err)
		}
	}
}
