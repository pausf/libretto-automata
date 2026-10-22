package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/link"
	"github.com/pausf/libretto-automata/internal/target"
)

// The shared criteria run one row per command that reads --json. Each row builds two
// states — one the command reports on, one it fails on — and hands back a runner bound to
// that state, so a test can run the same state with and without the flag.
type jsonRow struct {
	name    string
	ok      func(t *testing.T) jsonRunner
	failing func(t *testing.T) jsonRunner
}

// jsonRunner runs the command once, with extra arguments after its name, and returns what
// reached stdout and the error the command returned.
type jsonRunner func(extra ...string) (stdout string, err error)

var jsonRows = []jsonRow{
	{"status", statusReporting, statusFailing},
}

// statusReporting is a payload with one item of each kind and a target to scan.
func statusReporting(t *testing.T) jsonRunner {
	f := newFixture(t)
	f.skill(t, "alpha")
	f.agentItem(t, "beta")
	f.command(t, "gamma")
	return statusRunner(t, f)
}

// statusFailing is a payload whose skills/ is a file: the payload is present, so the
// command gets as far as scanning, and the scan cannot read the directory it expects.
func statusFailing(t *testing.T) jsonRunner {
	f := newFixture(t)
	f.agentItem(t, "beta")
	if err := os.WriteFile(filepath.Join(f.Repo, "skills"), []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	return statusRunner(t, f)
}

func statusRunner(t *testing.T, f fixture) jsonRunner {
	return func(extra ...string) (string, error) {
		out, _, err := capture(t, func() error { return runIn(f.Repo, append([]string{"status"}, extra...)) })
		return out, err
	}
}

// decodeObject parses stdout as one JSON object into a map, so an absent key stays
// observable — a typed decode would read it as zero.
func decodeObject(t *testing.T, out string) map[string]any {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, out)
	}
	return doc
}

// ── shared ───────────────────────────────────────────────────────────────────

func TestJSONIsOneObjectOnStdout(t *testing.T) {
	for _, row := range jsonRows {
		t.Run(row.name, func(t *testing.T) {
			out, err := row.ok(t)("--json")
			if err != nil {
				t.Fatalf("%s --json: %v", row.name, err)
			}
			if !strings.HasSuffix(out, "}\n") {
				t.Errorf("stdout does not end in one object and one newline: %q", out)
			}
			// Unmarshal refuses anything after the first value, a second object included.
			decodeObject(t, out)
		})
	}
}

func TestJSONCarriesSchemaOne(t *testing.T) {
	for _, row := range jsonRows {
		t.Run(row.name, func(t *testing.T) {
			out, err := row.ok(t)("--json")
			if err != nil {
				t.Fatalf("%s --json: %v", row.name, err)
			}
			dec := json.NewDecoder(strings.NewReader(out))
			dec.UseNumber()
			var doc map[string]any
			if err := dec.Decode(&doc); err != nil {
				t.Fatalf("stdout is not JSON: %v", err)
			}
			if got, ok := doc["schema"].(json.Number); !ok || got.String() != "1" {
				t.Errorf("schema = %#v, want the integer 1", doc["schema"])
			}
		})
	}
}

func TestJSONDoesNotChangeTheExitCode(t *testing.T) {
	for _, row := range jsonRows {
		for state, build := range map[string]func(*testing.T) jsonRunner{"reporting": row.ok, "failing": row.failing} {
			t.Run(row.name+"/"+state, func(t *testing.T) {
				r := build(t)
				_, textErr := r()
				_, jsonErr := r("--json")
				if (textErr == nil) != (jsonErr == nil) {
					t.Errorf("without --json err = %v, with it err = %v; the exit code moved", textErr, jsonErr)
				}
			})
		}
	}
}

func TestJSONFailureWritesNothingToStdout(t *testing.T) {
	for _, row := range jsonRows {
		t.Run(row.name, func(t *testing.T) {
			out, err := row.failing(t)("--json")
			if err == nil {
				t.Fatal("the failing state did not fail; the test proves nothing")
			}
			if out != "" {
				t.Errorf("a failed report wrote to stdout: %q", out)
			}
		})
	}
}

func TestJSONErrorsMatchTheTextErrors(t *testing.T) {
	for _, row := range jsonRows {
		t.Run(row.name, func(t *testing.T) {
			r := row.failing(t)
			_, textErr := r()
			_, jsonErr := r("--json")
			if textErr == nil || jsonErr == nil {
				t.Fatalf("want both to fail: text %v, json %v", textErr, jsonErr)
			}
			if textErr.Error() != jsonErr.Error() {
				t.Errorf("errors differ:\n text = %v\n json = %v", textErr, jsonErr)
			}
		})
	}
}

// encoding/json writes a control byte inside a string as \u001b, so this catches a raw
// escape written beside the document — the way a colour path would arrive.
func TestJSONHasNoEscapeCodes(t *testing.T) {
	for _, row := range jsonRows {
		t.Run(row.name, func(t *testing.T) {
			out, err := row.ok(t)("--json")
			if err != nil {
				t.Fatalf("%s --json: %v", row.name, err)
			}
			if out == "" {
				t.Fatal("stdout is empty; nothing to check")
			}
			if strings.Contains(out, "\x1b") {
				t.Errorf("the JSON carries an escape byte: %q", out)
			}
		})
	}
}

// ── status ───────────────────────────────────────────────────────────────────

// statusJSON runs status --json directly against a target and decodes what it wrote.
func statusJSON(t *testing.T, repo string, tg target.Target) map[string]any {
	t.Helper()
	out, _, err := capture(t, func() error { return status(repo, tg, "--json") })
	if err != nil {
		t.Fatalf("status --json: %v", err)
	}
	return decodeObject(t, out)
}

func objectAt(t *testing.T, doc map[string]any, key string) map[string]any {
	t.Helper()
	obj, ok := doc[key].(map[string]any)
	if !ok {
		t.Fatalf("%s = %#v, want an object", key, doc[key])
	}
	return obj
}

func arrayAt(t *testing.T, doc map[string]any, key string) []any {
	t.Helper()
	arr, ok := doc[key].([]any)
	if !ok {
		t.Fatalf("%s = %#v, want an array", key, doc[key])
	}
	return arr
}

func TestStatusJSONNamesTheTarget(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	tgt := objectAt(t, statusJSON(t, f.Repo, f.global()), "target")
	if tgt["name"] != "claude" || tgt["root"] != f.Claude || tgt["configured"] != true {
		t.Errorf("target = %#v, want claude at %s, configured", tgt, f.Claude)
	}

	// A root that does not exist yet is still known; it is only not configured.
	absent := filepath.Join(t.TempDir(), "nowhere")
	t.Setenv(claudeHomeEnv, absent)
	tgt = objectAt(t, statusJSON(t, f.Repo, target.Global()), "target")
	if tgt["root"] != absent || tgt["configured"] != false {
		t.Errorf("target = %#v, want root %s, not configured", tgt, absent)
	}
}

func TestStatusJSONOmitsAnUnknownRoot(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	tgt := objectAt(t, statusJSON(t, f.Repo, target.Project("")), "target")
	if _, ok := tgt["root"]; ok {
		t.Errorf("target = %#v, want no root key for a rootless target", tgt)
	}
	if tgt["name"] != "project" || tgt["configured"] != false {
		t.Errorf("target = %#v, want project, not configured", tgt)
	}
}

func TestStatusJSONListsEveryItem(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "beta")
	f.skill(t, "alpha")
	f.agentItem(t, "delta")
	f.command(t, "gamma")
	staleLink(t, f, "gone")

	want, err := link.Scan(f.Repo, f.global())
	if err != nil {
		t.Fatal(err)
	}
	items := arrayAt(t, statusJSON(t, f.Repo, f.global()), "items")
	if len(items) != len(want) {
		t.Fatalf("items = %d, want %d: %#v", len(items), len(want), items)
	}
	for i, e := range want {
		got, _ := items[i].(map[string]any)
		if got["kind"] != string(e.Kind) || got["name"] != e.Name || got["state"] != string(e.State) {
			t.Errorf("items[%d] = %#v, want %s/%s %s", i, got, e.Kind, e.Name, e.State)
		}
		if len(got) != 3 {
			t.Errorf("items[%d] carries %d keys, want kind, name and state: %#v", i, len(got), got)
		}
	}
}

func TestStatusJSONStatesMatchTheText(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "missing-one")
	linked := f.skill(t, "linked-one")
	f.link(t, linked, f.dest("skills", "linked-one"))
	other := f.skill(t, "other")
	f.skill(t, "wrong-one")
	f.link(t, other, f.dest("skills", "wrong-one"))
	f.skill(t, "conflict-one")
	f.putReal(t, "skills", "conflict-one", "not ours")
	staleLink(t, f, "stale-one")

	text, _, err := capture(t, func() error { return status(f.Repo, f.global()) })
	if err != nil {
		t.Fatal(err)
	}
	// "  <state padded> <kind>/<name>" — the state may hold a space, so it is everything
	// before the last field.
	textState := map[string]string{}
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || !strings.Contains(fields[len(fields)-1], "/") {
			continue
		}
		textState[fields[len(fields)-1]] = strings.Join(fields[:len(fields)-1], " ")
	}

	seen := map[string]bool{}
	for _, it := range arrayAt(t, statusJSON(t, f.Repo, f.global()), "items") {
		item := it.(map[string]any)
		key := item["kind"].(string) + "/" + item["name"].(string)
		state := item["state"].(string)
		seen[state] = true
		if textState[key] != state {
			t.Errorf("%s: JSON says %q, text says %q", key, state, textState[key])
		}
	}
	for _, s := range []link.State{link.Linked, link.Missing, link.WrongTarget, link.Conflict, link.Stale} {
		if !seen[string(s)] {
			t.Errorf("the fixture produced no %q item; the comparison does not cover it", s)
		}
	}
}

func TestStatusJSONEmptyIsAnEmptyArray(t *testing.T) {
	f := newFixture(t)

	out, _, err := capture(t, func() error { return status(f.Repo, f.global(), "--json") })
	if err != nil {
		t.Fatal(err)
	}
	if items := arrayAt(t, decodeObject(t, out), "items"); len(items) != 0 {
		t.Errorf("items = %#v, want []", items)
	}
}
