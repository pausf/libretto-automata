package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/link"
	"github.com/pausf/libretto-automata/internal/target"
)

// The shared criteria run one row per command that reads --json. Each row builds two
// states — one the command reports on, one it fails on — and hands back a runner bound to
// that state, so a test can run the same state with and without the flag. more holds any
// other state whose exit code is worth comparing.
type jsonRow struct {
	name    string
	ok      func(t *testing.T) jsonRunner
	failing func(t *testing.T) jsonRunner
	more    map[string]func(t *testing.T) jsonRunner
}

// jsonRunner runs the command once, with extra arguments after its name, and returns what
// reached stdout and the error the command returned.
type jsonRunner func(extra ...string) (stdout string, err error)

var jsonRows = []jsonRow{
	{name: "status", ok: unlinkedPayload("status"), failing: unreadablePayload("status")},
	{
		name: "doctor", ok: linkedPayload("doctor"), failing: unreadablePayload("doctor"),
		more: map[string]func(*testing.T) jsonRunner{"needs attention": unlinkedPayload("doctor")},
	},
	{name: "queue", ok: queueRunner(), failing: queueRunner("bogus")},
	{name: "metrics", ok: metricsRunner(), failing: metricsRunner("nope")},
}

// metricsRunner is a history of two changes, run with lead before any extra argument. A
// change git never saw is refused before anything is written.
func metricsRunner(lead ...string) func(*testing.T) jsonRunner {
	return func(t *testing.T) jsonRunner {
		return func(extra ...string) (string, error) {
			var out strings.Builder
			err := metricsWith(&out, append(append([]string{}, lead...), extra...), twoChangeGit(1700000000), "", "")
			return out.String(), err
		}
	}
}

// queueRunner is a project with one queued idea, run with lead before any extra argument.
// A lead argument is one queue refuses, which is its only failure.
func queueRunner(lead ...string) func(*testing.T) jsonRunner {
	return func(t *testing.T) jsonRunner {
		root := t.TempDir()
		writeProposal(t, root, "an-idea", "Queued: 2026-08-01\n\n## Reading\n\nAn idea.\n")
		return func(extra ...string) (string, error) {
			out, _, err := runQueue(t, root, queueGit(root, "", nil, nil), append(append([]string{}, lead...), extra...)...)
			return out, err
		}
	}
}

// unlinkedPayload is a payload with one item of each kind, none of them linked.
func unlinkedPayload(cmd string) func(*testing.T) jsonRunner {
	return func(t *testing.T) jsonRunner {
		f := newFixture(t)
		f.skill(t, "alpha")
		f.agentItem(t, "beta")
		f.command(t, "gamma")
		return payloadRunner(t, f, cmd)
	}
}

// linkedPayload is a payload whose one item is linked: nothing needs attention.
func linkedPayload(cmd string) func(*testing.T) jsonRunner {
	return func(t *testing.T) jsonRunner {
		f := newFixture(t)
		f.link(t, f.skill(t, "alpha"), f.dest("skills", "alpha"))
		return payloadRunner(t, f, cmd)
	}
}

// unreadablePayload is a payload whose skills/ is a file: the payload is present, so the
// command gets as far as scanning, and the scan cannot read the directory it expects.
func unreadablePayload(cmd string) func(*testing.T) jsonRunner {
	return func(t *testing.T) jsonRunner {
		f := newFixture(t)
		f.agentItem(t, "beta")
		if err := os.WriteFile(filepath.Join(f.Repo, "skills"), []byte("not a directory"), 0o644); err != nil {
			t.Fatal(err)
		}
		return payloadRunner(t, f, cmd)
	}
}

func payloadRunner(t *testing.T, f fixture, cmd string) jsonRunner {
	return func(extra ...string) (string, error) {
		out, _, err := capture(t, func() error { return runIn(f.Repo, append([]string{cmd}, extra...)) })
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
		states := map[string]func(*testing.T) jsonRunner{"reporting": row.ok, "failing": row.failing}
		for name, build := range row.more {
			states[name] = build
		}
		for state, build := range states {
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

// ── doctor ───────────────────────────────────────────────────────────────────

// doctorJSON runs doctor --json directly and decodes what it wrote. The exit code is
// doctor's verdict, not this helper's, so a non-nil error is returned rather than fatal.
func doctorJSON(t *testing.T, repo string, tg target.Target) (map[string]any, error) {
	t.Helper()
	out, _, err := capture(t, func() error { return doctor(repo, tg, "--json") })
	return decodeObject(t, out), err
}

func TestDoctorJSONListsOnlyProblems(t *testing.T) {
	f := newFixture(t)
	f.link(t, f.skill(t, "linked-one"), f.dest("skills", "linked-one"))
	f.skill(t, "missing-one")
	f.skill(t, "conflict-one")
	f.putReal(t, "skills", "conflict-one", "not ours")
	staleLink(t, f, "stale-one")

	entries, err := link.Scan(f.Repo, f.global())
	if err != nil {
		t.Fatal(err)
	}
	var want []link.Entry
	for _, e := range entries {
		if e.State.NeedsAttention() {
			want = append(want, e)
		}
	}

	doc, _ := doctorJSON(t, f.Repo, f.global())
	problems := arrayAt(t, doc, "problems")
	if len(problems) != len(want) {
		t.Fatalf("problems = %#v, want %d entries", problems, len(want))
	}
	for i, e := range want {
		got, _ := problems[i].(map[string]any)
		if got["kind"] != string(e.Kind) || got["name"] != e.Name || got["state"] != string(e.State) || len(got) != 3 {
			t.Errorf("problems[%d] = %#v, want %s/%s %s and nothing else", i, got, e.Kind, e.Name, e.State)
		}
		if got["name"] == "linked-one" {
			t.Errorf("a linked item was listed as a problem: %#v", got)
		}
	}

	// Nothing wrong is a known-empty list, not an absent one.
	clean := newFixture(t)
	clean.link(t, clean.skill(t, "alpha"), clean.dest("skills", "alpha"))
	doc, err = doctorJSON(t, clean.Repo, clean.global())
	if err != nil {
		t.Fatalf("a clean target failed: %v", err)
	}
	if p := arrayAt(t, doc, "problems"); len(p) != 0 {
		t.Errorf("problems = %#v, want []", p)
	}
}

func TestDoctorJSONNamesTheMode(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	doc, _ := doctorJSON(t, f.Repo, f.global())
	rel := objectAt(t, doc, "release")
	if rel["mode"] != "installed" || rel["payload_root"] != f.Repo {
		t.Errorf("release = %#v, want installed at %s", rel, f.Repo)
	}

	if err := os.Mkdir(filepath.Join(f.Repo, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	doc, _ = doctorJSON(t, f.Repo, f.global())
	if rel := objectAt(t, doc, "release"); rel["mode"] != "checkout" || rel["payload_root"] != f.Repo {
		t.Errorf("release = %#v, want checkout at %s", rel, f.Repo)
	}
}

func TestDoctorJSONCarriesTheRunningVersion(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	withVersion(t, "v0.4.2")

	doc, _ := doctorJSON(t, f.Repo, f.global())
	if got := objectAt(t, doc, "release")["running"]; got != "v0.4.2" {
		t.Errorf("running = %#v, want v0.4.2", got)
	}

	// Present even when the check failed — the running version never depended on it.
	stubRelease(t, "", errors.New("no route to host"))
	doc, _ = doctorJSON(t, f.Repo, f.global())
	if got := objectAt(t, doc, "release")["running"]; got != "v0.4.2" {
		t.Errorf("running after a failed check = %#v, want v0.4.2", got)
	}
}

func TestDoctorJSONCarriesTheLatestRelease(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	stubRelease(t, "v9.9.9", nil)

	doc, _ := doctorJSON(t, f.Repo, f.global())
	if got := objectAt(t, doc, "release")["latest"]; got != "v9.9.9" {
		t.Errorf("latest = %#v, want v9.9.9", got)
	}
}

func TestDoctorJSONRanksWhatRanks(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	withVersion(t, "v0.1.0")

	for latest, want := range map[string]bool{"v9.9.9": true, "v0.1.0": false} {
		stubRelease(t, latest, nil)
		doc, _ := doctorJSON(t, f.Repo, f.global())
		if got, ok := objectAt(t, doc, "release")["newer"].(bool); !ok || got != want {
			t.Errorf("latest %s against v0.1.0: newer = %#v, want %v", latest, objectAt(t, doc, "release")["newer"], want)
		}
	}
}

func TestDoctorJSONOmitsAnUnrankableComparison(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	cases := []struct{ name, running, latest string }{
		{"a dev build", "dev", "v0.3.0"},
		{"ahead of the remote", "v1.0.0", "v0.3.0"},
		{"no releases", "v0.1.0", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withVersion(t, c.running)
			stubRelease(t, c.latest, nil)
			doc, _ := doctorJSON(t, f.Repo, f.global())
			rel := objectAt(t, doc, "release")
			if _, ok := rel["newer"]; ok {
				t.Errorf("release = %#v, want no newer key", rel)
			}
			if _, ok := rel["latest"]; ok != (c.latest != "") {
				t.Errorf("release = %#v, latest present = %v, want %v", rel, ok, c.latest != "")
			}
		})
	}
}

func TestDoctorJSONReportsAFailedCheck(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	stubRelease(t, "", errors.New("no route to host\nhint: check the network"))

	doc, _ := doctorJSON(t, f.Repo, f.global())
	rel := objectAt(t, doc, "release")
	if rel["error"] != "no route to host" {
		t.Errorf("error = %#v, want the failure's first line", rel["error"])
	}
	for _, key := range []string{"latest", "newer"} {
		if _, ok := rel[key]; ok {
			t.Errorf("a failed check carries %s: %#v", key, rel)
		}
	}

	stubRelease(t, "v0.3.0", nil)
	doc, _ = doctorJSON(t, f.Repo, f.global())
	if _, ok := objectAt(t, doc, "release")["error"]; ok {
		t.Errorf("a check that answered carries error: %#v", objectAt(t, doc, "release"))
	}
}

func TestDoctorJSONListsPrerequisites(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	want := prerequisites()
	doc, _ := doctorJSON(t, f.Repo, f.global())
	got := arrayAt(t, doc, "prerequisites")
	if len(got) != len(want) {
		t.Fatalf("prerequisites = %#v, want %d", got, len(want))
	}
	for i, p := range want {
		obj, _ := got[i].(map[string]any)
		if obj["name"] != p.Name || obj["found"] != p.Found || len(obj) != 2 {
			t.Errorf("prerequisites[%d] = %#v, want name %q found %v and nothing else", i, obj, p.Name, p.Found)
		}
	}
}
