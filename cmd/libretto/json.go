package main

import (
	"encoding/json"
	"io"

	"github.com/pausf/libretto-automata/internal/link"
	"github.com/pausf/libretto-automata/internal/target"
)

// jsonSchema is the shape version every --json document carries. Adding a key never
// bumps it; removing, renaming or retyping one does, and that commit carries a `!`.
const jsonSchema = 1

// jsonFlag reports whether --json is among a command's arguments and returns the rest.
//
// A per-command flag, read after the command name: only the read-only reports honour
// it, and a global flag would need every other command to refuse it.
func jsonFlag(args []string) (bool, []string) {
	found := false
	rest := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--json" {
			found = true
			continue
		}
		rest = append(rest, a)
	}
	return found, rest
}

// writeJSON writes v as one indented document and one trailing newline.
//
// HTML escaping is off so a reading carrying `<`, `>` or `&` reaches the consumer as
// written. The documents are their own types, never the internal ones with tags on —
// every key is a promise, and a tag on an internal type makes each refactor of it a
// contract change.
func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

// targetJSON is the destination a report describes. Root is empty, and so omitted, only
// for a target that has none.
type targetJSON struct {
	Name       string `json:"name"`
	Root       string `json:"root,omitempty"`
	Configured bool   `json:"configured"`
}

// itemJSON is one scanned entry, its state spelled as the text report spells it.
type itemJSON struct {
	Kind  string `json:"kind"`
	Name  string `json:"name"`
	State string `json:"state"`
}

type statusDoc struct {
	Schema int        `json:"schema"`
	Target targetJSON `json:"target"`
	Items  []itemJSON `json:"items"`
}

func targetOf(tg target.Target) targetJSON {
	return targetJSON{Name: tg.Name(), Root: tg.Root(), Configured: configured(tg)}
}

// itemsOf is never nil: an empty scan is a known-empty list, written [] and not null.
func itemsOf(entries []link.Entry) []itemJSON {
	items := make([]itemJSON, 0, len(entries))
	for _, e := range entries {
		items = append(items, itemJSON{Kind: string(e.Kind), Name: e.Name, State: string(e.State)})
	}
	return items
}

// releaseJSON is what doctor's release check found. Latest is empty, and so omitted, when
// the check failed or the remote has nothing; Newer is nil when the two do not rank.
type releaseJSON struct {
	Mode        string `json:"mode"`
	PayloadRoot string `json:"payload_root"`
	Running     string `json:"running"`
	Latest      string `json:"latest,omitempty"`
	Newer       *bool  `json:"newer,omitempty"`
	Error       string `json:"error,omitempty"`
}

type prereqJSON struct {
	Name  string `json:"name"`
	Found bool   `json:"found"`
}

type doctorDoc struct {
	Schema        int          `json:"schema"`
	Target        targetJSON   `json:"target"`
	Problems      []itemJSON   `json:"problems"`
	Release       releaseJSON  `json:"release"`
	Prerequisites []prereqJSON `json:"prerequisites"`
}

func doctorDocOf(tg target.Target, problems []link.Entry, root string, c releaseCheck, prereqs []Prereq) doctorDoc {
	rel := releaseJSON{Mode: "installed", PayloadRoot: root, Running: c.running}
	if isRepo(root) {
		rel.Mode = "checkout"
	}
	if c.err != nil {
		rel.Error = firstLine(c.err.Error())
	} else {
		rel.Latest = c.latest
	}
	if newer, ranked := c.newer(); ranked {
		rel.Newer = &newer
	}

	ps := make([]prereqJSON, 0, len(prereqs))
	for _, p := range prereqs {
		ps = append(ps, prereqJSON{Name: p.Name, Found: p.Found})
	}
	return doctorDoc{Schema: jsonSchema, Target: targetOf(tg), Problems: itemsOf(problems), Release: rel, Prerequisites: ps}
}
