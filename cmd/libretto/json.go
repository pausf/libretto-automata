package main

import (
	"encoding/json"
	"io"
	"time"

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

type queuedJSON struct {
	Name    string `json:"name"`
	Date    string `json:"date"`
	Reading string `json:"reading,omitempty"`
}

type dispatchedJSON struct {
	Name   string `json:"name"`
	Branch string `json:"branch"`
}

// queueDoc's Dispatched is a pointer so the two states stay apart: nil is unknown and
// omitted, a pointer to an empty slice is known-empty and written [].
type queueDoc struct {
	Schema     int               `json:"schema"`
	Queued     []queuedJSON      `json:"queued"`
	Dispatched *[]dispatchedJSON `json:"dispatched,omitempty"`
}

func queueDocOf(r queueReport) queueDoc {
	doc := queueDoc{Schema: jsonSchema, Queued: make([]queuedJSON, 0, len(r.kept))}
	for _, i := range r.kept {
		doc.Queued = append(doc.Queued, queuedJSON{Name: i.name, Date: i.date, Reading: i.reading})
	}
	if r.branchesRead {
		d := make([]dispatchedJSON, 0, len(r.dropped))
		for _, i := range r.dropped {
			d = append(d, dispatchedJSON{Name: i.name, Branch: i.branch})
		}
		doc.Dispatched = &d
	}
	return doc
}

// changeJSON is one row of the metrics table. A pointer is a value that may be unknown:
// nil is omitted, a set pointer to zero is a known zero and stays present.
type changeJSON struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	Commits     *int   `json:"commits,omitempty"`
	First       string `json:"first,omitempty"`
	Last        string `json:"last,omitempty"`
	Closed      *int   `json:"closed,omitempty"`
	Boxes       *int   `json:"boxes,omitempty"`
	Reopened    *int   `json:"reopened,omitempty"`
	Corrections *int   `json:"corrections,omitempty"`
}

type correctionsJSON struct {
	ByPhase        map[string]int `json:"by_phase"`
	OutsideChanges int            `json:"outside_changes"`
}

type usageJSON struct {
	Input      int64 `json:"input"`
	Output     int64 `json:"output"`
	CacheWrite int64 `json:"cache_write"`
	CacheRead  int64 `json:"cache_read"`
}

// tokenChangeJSON's ByPhase is a pointer for the same reason dispatched is: a change no
// session reached has no phases to report, which is not the same as an empty set.
type tokenChangeJSON struct {
	Name    string                `json:"name"`
	Totals  *usageJSON            `json:"totals,omitempty"`
	ByPhase *map[string]usageJSON `json:"by_phase,omitempty"`
	NoPhase *usageJSON            `json:"no_phase,omitempty"`
}

type tokensJSON struct {
	Attributed          usageJSON        `json:"attributed"`
	Unattributed        usageJSON        `json:"unattributed"`
	Entries             int              `json:"entries"`
	UnattributedEntries int              `json:"unattributed_entries"`
	Change              *tokenChangeJSON `json:"change,omitempty"`
}

type loopRowJSON struct {
	Name      string   `json:"name"`
	Rounds    int      `json:"rounds"`
	NetClosed int      `json:"net_closed"`
	CostUSD   *float64 `json:"cost_usd,omitempty"`
	Unpriced  int      `json:"unpriced"`
}

// loopJSON carries either the rows and the skipped count, or the error that kept the log
// from being read — never both.
type loopJSON struct {
	Changes *[]loopRowJSON `json:"changes,omitempty"`
	Skipped *int           `json:"skipped,omitempty"`
	Error   string         `json:"error,omitempty"`
}

type metricsDoc struct {
	Schema      int              `json:"schema"`
	Changes     []changeJSON     `json:"changes"`
	Corrections *correctionsJSON `json:"corrections,omitempty"`
	Tokens      *tokensJSON      `json:"tokens,omitempty"`
	Loop        *loopJSON        `json:"loop,omitempty"`
}

func usageJSONOf(t usageTotals) usageJSON {
	return usageJSON{Input: t.in, Output: t.out, CacheWrite: t.cacheW, CacheRead: t.cacheR}
}

func intPtr(n int) *int { return &n }

func metricsDocOf(r metricsReport) metricsDoc {
	doc := metricsDoc{Schema: jsonSchema, Changes: make([]changeJSON, 0, len(r.rows))}

	for _, row := range r.rows {
		c := changeJSON{Name: row.name, State: "unreadable"}
		if r.ledger.seen {
			c.Corrections = intPtr(r.ledger.counts[row.name])
		}
		if row.err == nil {
			m := row.m
			c.State = "in flight"
			if m.landed {
				c.State = "landed"
			}
			c.Commits = intPtr(m.commits)
			c.First = m.first.UTC().Format(time.RFC3339)
			c.Last = m.last.UTC().Format(time.RFC3339)
			if m.planSeen {
				c.Closed, c.Boxes, c.Reopened = intPtr(m.checked-m.uncheck), intPtr(m.boxes), intPtr(m.uncheck)
			}
		}
		doc.Changes = append(doc.Changes, c)
	}

	if r.ledger.seen {
		doc.Corrections = &correctionsJSON{ByPhase: r.ledger.byPhase, OutsideChanges: r.ledger.orphans}
	}

	if t := r.tokens; t.found {
		tok := &tokensJSON{
			Attributed: usageJSONOf(t.attributed), Unattributed: usageJSONOf(t.unattributed),
			Entries: t.entries, UnattributedEntries: t.offEntries,
		}
		if t.only != "" {
			tok.Change = &tokenChangeJSON{Name: t.only}
			if t.attributedHere {
				total := usageJSONOf(t.total)
				byPhase := make(map[string]usageJSON, len(t.byPhase))
				for p, u := range t.byPhase {
					byPhase[p] = usageJSONOf(u)
				}
				tok.Change.Totals, tok.Change.ByPhase = &total, &byPhase
				if t.noPhaseSeen {
					noPhase := usageJSONOf(t.noPhase)
					tok.Change.NoPhase = &noPhase
				}
			}
		}
		doc.Tokens = tok
	}

	if l := r.loop; l.exists {
		if l.err != nil {
			doc.Loop = &loopJSON{Error: l.err.Error()}
		} else {
			rows := make([]loopRowJSON, 0, len(l.names))
			for _, n := range l.names {
				t := l.rows[n]
				row := loopRowJSON{Name: n, Rounds: t.rounds, NetClosed: t.net, Unpriced: t.unpriced}
				if t.priced {
					cost := t.cost
					row.CostUSD = &cost
				}
				rows = append(rows, row)
			}
			doc.Loop = &loopJSON{Changes: &rows, Skipped: intPtr(l.skipped)}
		}
	}
	return doc
}
