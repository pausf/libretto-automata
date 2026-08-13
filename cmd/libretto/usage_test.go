package main

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// assistantLine is the shape measured off a real transcript: the four numbers live
// under .message.usage, and iterations[] repeats them. A parser that adds both
// double-counts every entry, so the fixture always carries iterations.
func assistantLine(branch, skill string, in, out, cw, cr int64) string {
	return fmt.Sprintf(`{"type":"assistant","gitBranch":%q,"attributionSkill":%q,`+
		`"message":{"model":"claude-opus-5","usage":{`+
		`"input_tokens":%d,"output_tokens":%d,`+
		`"cache_creation_input_tokens":%d,"cache_read_input_tokens":%d,`+
		`"iterations":[{"input_tokens":%d,"output_tokens":%d,`+
		`"cache_creation_input_tokens":%d,"cache_read_input_tokens":%d}]}}}`,
		branch, skill, in, out, cw, cr, in, out, cw, cr)
}

// transcriptRoot writes a projects/ tree for repoRoot and returns the projects dir.
// Keys are paths relative to the encoded project directory, so a subagent file is
// just "<session>/subagents/agent-1.jsonl".
func transcriptRoot(t *testing.T, repoRoot string, files map[string][]string) string {
	t.Helper()
	projects := filepath.Join(t.TempDir(), "projects")
	dir := filepath.Join(projects, projectDirName(repoRoot))
	for rel, lines := range files {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", p, err)
		}
		body := ""
		for _, l := range lines {
			body += l + "\n"
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", p, err)
		}
	}
	return projects
}

func totals(es []usageEntry) (in, out, cw, cr int64) {
	for _, e := range es {
		in, out, cw, cr = in+e.in, out+e.out, cw+e.cacheW, cr+e.cacheR
	}
	return
}

func TestTheFourUsageNumbersAreKeptApart(t *testing.T) {
	root := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl": {
			assistantLine("main", "write-spec", 2, 104, 36331, 22234),
			assistantLine("main", "write-spec", 3, 200, 1000, 500000),
		},
	})

	es, found := readUsage(root, "/repo")
	if !found {
		t.Fatal("the project directory exists and was not found")
	}
	in, out, cw, cr := totals(es)

	// Each number stands on its own. Summing them, or folding iterations[] back in,
	// changes exactly these four and nothing else — which is what makes this the
	// assertion rather than a total.
	if in != 5 || out != 304 || cw != 37331 || cr != 522234 {
		t.Errorf("got in=%d out=%d cacheW=%d cacheR=%d, want 5 304 37331 522234", in, out, cw, cr)
	}
}

func TestAMalformedLineDoesNotCostTheRestOfTheFile(t *testing.T) {
	root := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl": {
			assistantLine("main", "write-spec", 1, 10, 100, 1000),
			`{"type":"assistant","message":{"usage":`, // truncated mid-object
			`not json at all`,
			``,                                       // a blank line
			`{"type":"mode"}`,                        // an entry type with no envelope at all
			`{"type":"user","message":{"role":"u"}}`, // no usage object
			`{"type":"assistant","message":{"model":"<synthetic>","usage":{` +
				`"input_tokens":0,"output_tokens":0,"cache_creation_input_tokens":0,` +
				`"cache_read_input_tokens":0,"service_tier":null,"iterations":null}}}`,
			assistantLine("main", "write-spec", 2, 20, 200, 2000),
		},
	})

	es, _ := readUsage(root, "/repo")
	in, out, cw, cr := totals(es)
	if in != 3 || out != 30 || cw != 300 || cr != 3000 {
		t.Errorf("got in=%d out=%d cacheW=%d cacheR=%d, want 3 30 300 3000", in, out, cw, cr)
	}
}

func TestSubagentTranscriptsAreCounted(t *testing.T) {
	root := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl": {assistantLine("main", "review-work", 1, 10, 100, 1000)},
		// One review-lens subagent file was measured at 3M cache reads. A reader that
		// walks only the top level undercounts by more than it reports.
		"sess-1/subagents/agent-abc.jsonl": {assistantLine("main", "review-work", 4, 40, 400, 4000)},
	})

	es, _ := readUsage(root, "/repo")
	in, out, cw, cr := totals(es)
	if in != 5 || out != 50 || cw != 500 || cr != 5000 {
		t.Errorf("got in=%d out=%d cacheW=%d cacheR=%d, want 5 50 500 5000 — subagents missing?", in, out, cw, cr)
	}
}

// snapshot records every path under root with its size and content hash. Size alone
// misses an in-place rewrite of the same length; mtime is a coin flip at this
// granularity. The hash is what makes the comparison mean "unchanged".
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	got := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			got[p] = "dir"
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		got[p] = fmt.Sprintf("%d/%x", len(b), sha256.Sum256(b))
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", root, err)
	}
	return got
}

func TestTheTranscriptRootIsNeverWritten(t *testing.T) {
	root := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl":                          {assistantLine("main", "write-spec", 1, 10, 100, 1000)},
		"sess-1/subagents/agent-abc.jsonl": {assistantLine("main", "write-spec", 2, 20, 200, 2000)},
	})

	before := snapshot(t, root)
	readUsage(root, "/repo")
	after := snapshot(t, root)

	for p, want := range before {
		got, ok := after[p]
		if !ok {
			t.Errorf("%s disappeared", p)
			continue
		}
		if got != want {
			t.Errorf("%s changed: %s -> %s", p, want, got)
		}
	}
	for p := range after {
		if _, ok := before[p]; !ok {
			t.Errorf("%s was created", p)
		}
	}
}

func TestAMissingProjectDirectoryIsAStateNotAnError(t *testing.T) {
	root := transcriptRoot(t, "/other-repo", map[string][]string{
		"a.jsonl": {assistantLine("main", "write-spec", 1, 10, 100, 1000)},
	})

	es, found := readUsage(root, "/repo")
	if found {
		t.Error("reported a project directory for a repository that has none")
	}
	if len(es) != 0 {
		t.Errorf("got %d entries for an absent project directory, want 0", len(es))
	}

	if _, found := readUsage(filepath.Join(t.TempDir(), "nope"), "/repo"); found {
		t.Error("reported a transcript root that does not exist")
	}
}

func TestBranchIsReadPerEntryNotPerFile(t *testing.T) {
	// Measured on a real session: one file spanning four branches. A per-file reading
	// misattributes every entry after a checkout.
	root := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl": {
			assistantLine("feat/add-thing", "write-spec", 1, 10, 100, 1000),
			assistantLine("feat/other-thing", "write-spec", 2, 20, 200, 2000),
		},
	})
	es, _ := readUsage(root, "/repo")

	got := map[string]int64{}
	for _, e := range es {
		got[e.branch] += e.in
	}
	if got["feat/add-thing"] != 1 || got["feat/other-thing"] != 2 {
		t.Errorf("got %v, want one entry on each branch", got)
	}
}

func TestAPrefixIsStrippedAndTheNameMatchedWhole(t *testing.T) {
	changes := []string{"add-thing", "add-thing-extra"}

	cases := []struct {
		branch string
		want   string // "" means unattributed
	}{
		{"feat/add-thing", "add-thing"},
		{"fix/add-thing", "add-thing"},
		{"docs/add-thing", "add-thing"},
		{"chore/add-thing", "add-thing"},
		{"refactor/add-thing", "add-thing"},
		{"add-thing", "add-thing"},
		{"add-thing-extra", "add-thing-extra"},
		// The negative is the point. A prefix or substring rule attributes this to
		// add-thing, and one wrong guess becomes a number nobody can audit.
		{"feat/add-thing-more", ""},
		{"feat/unrelated", ""},
		{"main", ""},
		{"HEAD", ""},
		{"", ""},
		// An unrecognised prefix fails safe: unattributed, never the wrong change.
		{"wip/add-thing", ""},
	}
	for _, c := range cases {
		if got := attributeBranch(c.branch, changes); got != c.want {
			t.Errorf("attributeBranch(%q) = %q, want %q", c.branch, got, c.want)
		}
	}
}

func TestUnattributedTokensAreReportedNotDiscarded(t *testing.T) {
	root := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl": {
			assistantLine("feat/add-thing", "write-spec", 1, 10, 100, 1000),
			assistantLine("main", "write-spec", 2, 20, 200, 2000),
			assistantLine("HEAD", "write-spec", 4, 40, 400, 4000),
			assistantLine("feat/never-a-change", "write-spec", 8, 80, 800, 8000),
			assistantLine("", "write-spec", 16, 160, 1600, 16000),
		},
	})
	es, _ := readUsage(root, "/repo")

	byChange, unattributed := attribute(es, []string{"add-thing"})

	if byChange["add-thing"].in != 1 {
		t.Errorf("attributed in=%d, want 1", byChange["add-thing"].in)
	}
	if unattributed.in != 30 {
		t.Errorf("unattributed in=%d, want 30 — main, HEAD, an unmatched branch and an absent one", unattributed.in)
	}

	// The invariant that stops a quiet drop, and the reason this is not just two sums:
	// every token read is in exactly one bucket.
	corpusIn, corpusOut, corpusW, corpusR := totals(es)
	var sum usageTotals
	for _, t := range byChange {
		sum = sum.plus(t)
	}
	sum = sum.plus(unattributed)
	if sum.in != corpusIn || sum.out != corpusOut || sum.cacheW != corpusW || sum.cacheR != corpusR {
		t.Errorf("attributed + unattributed = %+v, want corpus %d %d %d %d",
			sum, corpusIn, corpusOut, corpusW, corpusR)
	}
}
