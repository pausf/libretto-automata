package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Every case runs against a fake gitRunner and a t.TempDir(): the command makes two git
// calls, and both are answered here, so no repository is ever built.

// queueGit answers rev-parse with root and the branch listing with branches, or with
// branchErr when it is set. Every call is logged when log is non-nil.
func queueGit(root, branches string, branchErr error, log *[]string) gitRunner {
	return func(args ...string) (string, error) {
		if log != nil {
			*log = append(*log, strings.Join(args, " "))
		}
		switch {
		case len(args) > 0 && args[0] == "rev-parse":
			return root + "\n", nil
		case len(args) > 0 && args[0] == "branch":
			return branches, branchErr
		}
		return "", errors.New("unexpected git call: " + strings.Join(args, " "))
	}
}

func writeProposal(t *testing.T, root, name, body string) {
	t.Helper()
	dir := filepath.Join(root, ".agents", "changes", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runQueue(t *testing.T, cwd string, git gitRunner, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	var out, errw bytes.Buffer
	err = queue(&out, &errw, args, cwd, git)
	return out.String(), errw.String(), err
}

func lines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}

func TestQueueListsOldestFirst(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "a-later", "# a-later\n\nQueued: 2026-09-02\n\n## Reading\n\nThe later idea. It has more.\n")
	writeProposal(t, root, "z-earlier", "# z-earlier\n\nQueued: 2026-08-01\n\n## Reading\n\nThe earlier\nidea. Second sentence.\n\n## Notes\n\nNot read.\n")

	out, _, err := runQueue(t, root, queueGit(root, "", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"2026-08-01  z-earlier  The earlier idea.",
		"2026-09-02  a-later  The later idea.",
	}
	if got := lines(out); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("queue =\n%s\nwant\n%s", out, strings.Join(want, "\n"))
	}
}

func TestQueueBreaksDateTiesByName(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "zeta", "Queued: 2026-08-01\n")
	writeProposal(t, root, "alpha", "Queued: 2026-08-01\n")
	writeProposal(t, root, "mid", "Queued: 2026-08-01\n")

	out, _, err := runQueue(t, root, queueGit(root, "", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	want := "2026-08-01  alpha\n2026-08-01  mid\n2026-08-01  zeta\n"
	if out != want {
		t.Errorf("queue =\n%s\nwant\n%s", out, want)
	}
}

func TestQueueWithoutAReadingPrintsTheNameAlone(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "bare", "# bare\n\nQueued: 2026-08-01\n\nJust an idea, no reading.\n")

	out, _, err := runQueue(t, root, queueGit(root, "", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	if out != "2026-08-01  bare\n" {
		t.Errorf("queue = %q, want the date and the name only", out)
	}
}

func TestQueueReadsOnlyALineStartingWithQueued(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "quoted", "# quoted\n\n> Queued: 2026-08-01 was said in the ask\n")
	writeProposal(t, root, "inline", "# inline\n\nIt was Queued: 2026-08-01 once.\n")
	writeProposal(t, root, "real", "# real\n\nQueued: 2026-08-02\n")

	out, _, err := runQueue(t, root, queueGit(root, "", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	if out != "2026-08-02  real\n" {
		t.Errorf("queue = %q, want real alone", out)
	}
}

func TestQueueSkipsAProposalWithAChecklistBox(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "started", "# started\n\nQueued: 2026-08-01\n\n## Checklist\n\n- [ ] the work began\n")
	// A box outside the section is part of the ask, not the work.
	writeProposal(t, root, "boxed-ask", "# boxed-ask\n\nQueued: 2026-08-02\n\n- [ ] quoted in the ask\n")

	out, _, err := runQueue(t, root, queueGit(root, "", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	if out != "2026-08-02  boxed-ask\n" {
		t.Errorf("queue = %q, want boxed-ask alone", out)
	}
}

func TestQueueDropsAChangeItsBranchAlreadyTook(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "an-idea", "Queued: 2026-08-01\n")
	writeProposal(t, root, "kept", "Queued: 2026-08-02\n")

	// With no branch, the change is listed: a scan that simply missed it cannot pass below.
	out, _, err := runQueue(t, root, queueGit(root, "", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "2026-08-01  an-idea\n") {
		t.Fatalf("an-idea is not listed with no branches:\n%s", out)
	}

	for _, branch := range []string{"feat/an-idea", "fix/an-idea", "docs/an-idea", "chore/an-idea", "refactor/an-idea", "an-idea"} {
		t.Run(branch, func(t *testing.T) {
			out, _, err := runQueue(t, root, queueGit(root, "main-work\n"+branch+"\n", nil, nil))
			if err != nil {
				t.Fatal(err)
			}
			got := lines(out)
			if len(got) != 2 {
				t.Fatalf("want the kept change and one closing line, got:\n%s", out)
			}
			if got[0] != "2026-08-02  kept" {
				t.Errorf("first line = %q, want the kept change", got[0])
			}
			if strings.Contains(got[0], "an-idea") || !strings.Contains(got[1], "an-idea ("+branch+")") {
				t.Errorf("the closing line does not name an-idea with %s:\n%s", branch, out)
			}
		})
	}
}

func TestQueueKeepsAChangeOnAPartialOrForeignMatch(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "an-idea", "Queued: 2026-08-01\n")

	branches := "feat/an-idea-v2\nan-ide\nfeat/an\nwip/an-idea\nuser/an-idea\n"
	out, _, err := runQueue(t, root, queueGit(root, branches, nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	if out != "2026-08-01  an-idea\n" {
		t.Errorf("queue = %q, want an-idea kept and no closing line", out)
	}
}

func TestQueueAsksGitForUnmergedLocalBranches(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "an-idea", "Queued: 2026-08-01\n")

	var log []string
	if _, _, err := runQueue(t, root, queueGit(root, "", nil, &log)); err != nil {
		t.Fatal(err)
	}
	var asked []string
	for _, call := range log {
		if strings.HasPrefix(call, "branch") {
			asked = append(asked, call)
		}
	}
	if len(asked) != 1 || asked[0] != "branch --format=%(refname:short) --no-merged HEAD" {
		t.Errorf("branch calls = %q, want exactly one with the unmerged-local form", asked)
	}
}

func TestQueueUnreadableBranchesListsUnfiltered(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "an-idea", "Queued: 2026-08-01\n")

	// The listing is there and git still fails: nothing it says may filter anything.
	out, stderr, err := runQueue(t, root, queueGit(root, "feat/an-idea\n", errors.New("exit status 128"), nil))
	if err != nil {
		t.Fatalf("unreadable branches is not an error: %v", err)
	}
	if out != "2026-08-01  an-idea\n" {
		t.Errorf("queue = %q, want an-idea unfiltered", out)
	}
	if got := lines(stderr); stderr == "" || len(got) != 1 {
		t.Errorf("stderr = %q, want exactly one line", stderr)
	}
}

func TestQueueReadsTheRepositoryRoot(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "an-idea", "Queued: 2026-08-01\n")
	sub := filepath.Join(root, "cmd", "deep")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	out, _, err := runQueue(t, sub, queueGit(root, "", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	if out != "2026-08-01  an-idea\n" {
		t.Errorf("from a subdirectory, queue = %q, want the root's an-idea", out)
	}
}

func TestQueueWithoutAChangesDirSaysSo(t *testing.T) {
	root := t.TempDir()

	out, _, err := runQueue(t, root, queueGit(root, "", nil, nil))
	if err != nil {
		t.Fatalf("a missing changes directory is not an error: %v", err)
	}
	if out != "no .agents/changes/ here — nothing is queued\n" {
		t.Errorf("queue = %q", out)
	}
}

func TestQueueEmptyIsAStateNotAnError(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "in-flight", "# in-flight\n\n## Checklist\n\n- [ ] doing it\n")

	out, _, err := runQueue(t, root, queueGit(root, "", nil, nil))
	if err != nil {
		t.Fatalf("an empty queue is not an error: %v", err)
	}
	if out != "nothing is queued\n" {
		t.Errorf("queue = %q", out)
	}
}

func TestQueueAllDispatchedSaysNothingQueuedThenNamesThem(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "one", "Queued: 2026-08-01\n")
	writeProposal(t, root, "two", "Queued: 2026-08-02\n")

	out, stderr, err := runQueue(t, root, queueGit(root, "feat/one\nfix/two\n", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	got := lines(out)
	if len(got) != 2 || got[0] != "nothing is queued" {
		t.Fatalf("stdout = %q, want `nothing is queued` then the closing line", out)
	}
	if !strings.Contains(got[1], "one (feat/one)") || !strings.Contains(got[1], "two (fix/two)") {
		t.Errorf("closing line = %q, want both changes with their branches", got[1])
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want nothing", stderr)
	}
}

func TestQueueWritesNothing(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "an-idea", "# an-idea\n\nQueued: 2026-08-01\n\n## Reading\n\nAn idea.\n")
	writeProposal(t, root, "taken", "Queued: 2026-08-02\n")

	snapshot := func() map[string]string {
		tree := map[string]string{}
		err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			info, err := d.Info()
			if err != nil {
				return err
			}
			v := info.Mode().String()
			if !d.IsDir() {
				b, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				v += " " + string(b)
			}
			tree[p] = v
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		return tree
	}
	before := snapshot()

	out, _, err := runQueue(t, root, queueGit(root, "feat/taken\n", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	// A command that printed nothing also wrote nothing; it has to have run first.
	if !strings.Contains(out, "2026-08-01  an-idea  An idea.") {
		t.Fatalf("the queue was not printed:\n%s", out)
	}
	after := snapshot()
	if len(after) != len(before) {
		t.Errorf("tree went from %d entries to %d", len(before), len(after))
	}
	for p, v := range before {
		if after[p] != v {
			t.Errorf("%s changed", p)
		}
	}
}

func TestQueueRefusesArguments(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "an-idea", "Queued: 2026-08-01\n")

	for _, arg := range []string{"--jsno", "an-idea", "--yes"} {
		t.Run(arg, func(t *testing.T) {
			out, _, err := runQueue(t, root, queueGit(root, "", nil, nil), arg)
			if err == nil || !strings.Contains(err.Error(), arg) {
				t.Errorf("err = %v, want a refusal naming %s", err, arg)
			}
			if out != "" {
				t.Errorf("a refused run printed a queue:\n%s", out)
			}
		})
	}

	// --json is the one argument it reads.
	t.Run("--json", func(t *testing.T) {
		out, _, err := runQueue(t, root, queueGit(root, "", nil, nil), "--json")
		if err != nil {
			t.Fatalf("--json was refused: %v", err)
		}
		if !strings.Contains(out, "an-idea") {
			t.Errorf("--json printed no queue:\n%s", out)
		}
	})
}

// Proven through run, as init proves its own: a payload root holding nothing and a cwd
// that is no repository, so whatever comes back is queue's and never the payload gate's.
func TestQueueNeedsNoPayload(t *testing.T) {
	if needsPayload([]string{"queue"}) {
		t.Fatal("queue reads the project, not the payload tree")
	}
	t.Setenv(EnvRoot, filepath.Join(t.TempDir(), "never-installed"))
	t.Chdir(t.TempDir())

	out, _, err := capture(t, func() error { return run([]string{"queue"}) })
	if err != nil {
		t.Fatalf("queue with no payload: %v", err)
	}
	if !strings.Contains(out, "nothing is queued") {
		t.Errorf("stdout = %q, want queue's own empty state", out)
	}
}

// ── --json ───────────────────────────────────────────────────────────────────

func queueJSON(t *testing.T, cwd string, git gitRunner) map[string]any {
	t.Helper()
	out, _, err := runQueue(t, cwd, git, "--json")
	if err != nil {
		t.Fatalf("queue --json: %v", err)
	}
	return decodeObject(t, out)
}

func TestQueueJSONListsOldestFirst(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "a-later", "# a-later\n\nQueued: 2026-09-02\n\n## Reading\n\nThe later idea. It has more.\n")
	writeProposal(t, root, "z-earlier", "# z-earlier\n\nQueued: 2026-08-01\n\n## Reading\n\nThe earlier\nidea. Second sentence.\n")
	writeProposal(t, root, "b-tie", "Queued: 2026-08-01\n\n## Reading\n\nTied.\n")

	queued := arrayAt(t, queueJSON(t, root, queueGit(root, "", nil, nil)), "queued")
	want := []map[string]any{
		{"name": "b-tie", "date": "2026-08-01", "reading": "Tied."},
		{"name": "z-earlier", "date": "2026-08-01", "reading": "The earlier idea."},
		{"name": "a-later", "date": "2026-09-02", "reading": "The later idea."},
	}
	if len(queued) != len(want) {
		t.Fatalf("queued = %#v, want %d ideas", queued, len(want))
	}
	for i, w := range want {
		got, _ := queued[i].(map[string]any)
		if len(got) != len(w) || got["name"] != w["name"] || got["date"] != w["date"] || got["reading"] != w["reading"] {
			t.Errorf("queued[%d] = %#v, want %#v", i, got, w)
		}
	}
}

func TestQueueJSONOmitsAMissingReading(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "bare", "# bare\n\nQueued: 2026-08-01\n\nJust an idea, no reading.\n")

	queued := arrayAt(t, queueJSON(t, root, queueGit(root, "", nil, nil)), "queued")
	if len(queued) != 1 {
		t.Fatalf("queued = %#v, want bare alone", queued)
	}
	got := queued[0].(map[string]any)
	if _, ok := got["reading"]; ok {
		t.Errorf("queued[0] = %#v, want no reading key", got)
	}
	if got["name"] != "bare" || got["date"] != "2026-08-01" {
		t.Errorf("queued[0] = %#v, want bare on 2026-08-01", got)
	}
}

func TestQueueJSONKeepsReadingTextVerbatim(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "html", "Queued: 2026-08-01\n\n## Reading\n\nKeep <b> & </b> as written.\n")

	out, _, err := runQueue(t, root, queueGit(root, "", nil, nil), "--json")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Keep <b> & </b> as written.") {
		t.Errorf("the reading did not reach stdout as written:\n%s", out)
	}
	queued := arrayAt(t, decodeObject(t, out), "queued")
	if got := queued[0].(map[string]any)["reading"]; got != "Keep <b> & </b> as written." {
		t.Errorf("reading = %#v", got)
	}
}

func TestQueueJSONNamesDispatchedChanges(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "an-idea", "Queued: 2026-08-01\n")
	writeProposal(t, root, "kept", "Queued: 2026-08-02\n")

	doc := queueJSON(t, root, queueGit(root, "main-work\nfeat/an-idea\n", nil, nil))
	dispatched := arrayAt(t, doc, "dispatched")
	if len(dispatched) != 1 {
		t.Fatalf("dispatched = %#v, want an-idea alone", dispatched)
	}
	if got := dispatched[0].(map[string]any); got["name"] != "an-idea" || got["branch"] != "feat/an-idea" || len(got) != 2 {
		t.Errorf("dispatched[0] = %#v, want an-idea on feat/an-idea", got)
	}
	if queued := arrayAt(t, doc, "queued"); len(queued) != 1 || queued[0].(map[string]any)["name"] != "kept" {
		t.Errorf("queued = %#v, want kept alone", queued)
	}

	// Branches read, none of them took anything: a known-empty list.
	if d := arrayAt(t, queueJSON(t, root, queueGit(root, "main-work\n", nil, nil)), "dispatched"); len(d) != 0 {
		t.Errorf("dispatched = %#v, want []", d)
	}
}

func TestQueueJSONUnreadableBranchesOmitsDispatched(t *testing.T) {
	root := t.TempDir()
	writeProposal(t, root, "an-idea", "Queued: 2026-08-01\n")

	out, stderr, err := runQueue(t, root, queueGit(root, "feat/an-idea\n", errors.New("exit status 128"), nil), "--json")
	if err != nil {
		t.Fatalf("unreadable branches is not an error: %v", err)
	}
	doc := decodeObject(t, out)
	if _, ok := doc["dispatched"]; ok {
		t.Errorf("doc = %#v, want no dispatched key when git could not list branches", doc)
	}
	if queued := arrayAt(t, doc, "queued"); len(queued) != 1 {
		t.Errorf("queued = %#v, want an-idea unfiltered", queued)
	}
	if got := lines(stderr); stderr == "" || len(got) != 1 {
		t.Errorf("stderr = %q, want the same one line as without --json", stderr)
	}
}

func TestQueueJSONWithoutAChangesDirIsEmpty(t *testing.T) {
	root := t.TempDir()

	doc := queueJSON(t, root, queueGit(root, "feat/anything\n", nil, nil))
	if queued := arrayAt(t, doc, "queued"); len(queued) != 0 {
		t.Errorf("queued = %#v, want []", queued)
	}
	if _, ok := doc["dispatched"]; ok {
		t.Errorf("doc = %#v, want no dispatched key: git was never asked", doc)
	}
}
