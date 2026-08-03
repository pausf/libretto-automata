// Package link enumerates repo items and reconciles them with a target.
//
// Everything here is read-only. Mutation lives in apply.go (see docs/PLAN.md
// phase 3) and does not exist yet.
package link

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pausf/libretto-automata/internal/target"
)

// Item is one thing the repo owns and can link into a target.
type Item struct {
	Kind target.Kind
	Name string // basename as it appears in both repo and target
	Path string // absolute path inside the repo
}

// Items lists the repo's items of one kind, sorted by name.
//
// A missing kind directory is not an error — the repo may simply hold no items
// of that kind yet. Entries of the wrong shape are skipped: a stray file among
// skills, or a non-.md file among agents, is somebody's note, not an item.
func Items(repoRoot string, k target.Kind) ([]Item, error) {
	dir := filepath.Join(repoRoot, string(k))

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var items []Item
	for _, e := range entries {
		if !isItem(e, k) {
			continue
		}
		items = append(items, Item{
			Kind: k,
			Name: e.Name(),
			Path: filepath.Join(dir, e.Name()),
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

// isItem reports whether a directory entry is an item of this kind.
func isItem(e os.DirEntry, k target.Kind) bool {
	name := e.Name()

	// Dotfiles are never items. This is what keeps .gitkeep out of the count.
	if name == "" || name[0] == '.' {
		return false
	}
	if e.IsDir() != k.ItemsAreDirs() {
		return false
	}
	if ext := k.Ext(); ext != "" && filepath.Ext(name) != ext {
		return false
	}
	return true
}

// Counts returns the number of items per kind, for every kind the target
// accepts. Kinds the target rejects are absent from the map rather than zero, so
// a caller can tell "none" from "not applicable".
func Counts(repoRoot string, t target.Target) (map[target.Kind]int, error) {
	counts := make(map[target.Kind]int, len(t.Kinds()))
	for _, k := range t.Kinds() {
		items, err := Items(repoRoot, k)
		if err != nil {
			return nil, err
		}
		counts[k] = len(items)
	}
	return counts, nil
}
