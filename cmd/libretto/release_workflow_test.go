package main

import (
	"strings"
	"testing"
)

// The release workflow is the one thing in this repository that writes, and every
// property below is one that fails silently or dangerously rather than loudly. Read as
// text for the reason `gates_test.go` gives: a YAML dependency to assert on a config file
// is a dependency added to check a checker.

func releaseWorkflow(t *testing.T) string {
	t.Helper()
	return repoFile(t, ".github/workflows/release.yml")
}

// The dangerous one. This job holds `contents: write` and is triggered by a
// pull-request event, so checking out the request's own code would run a contributor's
// tree with the token in the environment. The merged work is already on main.
func TestReleaseWorkflowChecksOutMainAndNotTheRequestHead(t *testing.T) {
	workflow := releaseWorkflow(t)

	if !strings.Contains(workflow, "ref: main") {
		t.Error("the release workflow does not pin the checkout to main")
	}
	for _, forbidden := range []string{"head.sha", "head.ref", "head_ref"} {
		if strings.Contains(workflow, forbidden) {
			t.Errorf("the release workflow references %q — a job with contents: write must never check out the request's own code", forbidden)
		}
	}
}

// The silent one. actions/checkout defaults to depth 1 with no tags, and
// `git describe --tags` against that finds nothing, so the next version computes from
// zero. The symptom is a v0.0.1 tag on a repository with a dozen releases.
func TestReleaseWorkflowFetchesEnoughHistoryToSeeTheTags(t *testing.T) {
	if !strings.Contains(releaseWorkflow(t), "fetch-depth: 0") {
		t.Error("the release workflow does not set fetch-depth: 0 — git describe cannot see tags in a shallow clone")
	}
}

// The one with no error message at all. A tag pushed with GITHUB_TOKEN does not trigger
// workflows, so a second workflow listening on `push: tags` would never fire and the only
// evidence is a tag with no Release. Both halves live in one job.
func TestReleaseWorkflowPublishesInTheSameRunThatPushesTheTag(t *testing.T) {
	workflow := releaseWorkflow(t)

	if !strings.Contains(workflow, "git push origin") {
		t.Error("the release workflow does not push the tag")
	}
	if !strings.Contains(workflow, "gh release create") {
		t.Error("the release workflow does not create the Release")
	}
	if strings.Contains(workflow, "push:\n    tags") || strings.Contains(workflow, "tags:") {
		t.Error("the release workflow triggers on a tag push — GITHUB_TOKEN tag pushes do not trigger workflows, so that job would never run")
	}
}

// Refusing is the feature. The bump is a reading of .agents/specs/ and not of the commit
// log, so a missing label has to stop the run: defaulting to patch is the silently-wrong
// bump this design exists to avoid.
func TestReleaseWorkflowRefusesWithoutABumpLabel(t *testing.T) {
	workflow := releaseWorkflow(t)

	for _, label := range []string{"release:major", "release:minor", "release:patch"} {
		if !strings.Contains(workflow, label) {
			t.Errorf("the release workflow does not recognise the %q label", label)
		}
	}
	if !strings.Contains(workflow, "no release: label") {
		t.Error("the release workflow does not refuse when no bump label is present")
	}
}

// A release is the one moment where "the gates were green earlier" is not good enough:
// the merge commit is a tree no gate has run against.
func TestReleaseWorkflowRunsTheGatesBeforeItTags(t *testing.T) {
	workflow := releaseWorkflow(t)

	gatesAt := strings.Index(workflow, "make gates")
	if gatesAt < 0 {
		t.Fatal("the release workflow does not run the gates")
	}
	pushAt := strings.Index(workflow, "git push origin")
	if pushAt < 0 {
		t.Fatal("the release workflow does not push the tag")
	}
	if gatesAt > pushAt {
		t.Error("the release workflow tags before running the gates")
	}
}

// Two merges landing together would read the same last tag, and the second would try to
// create a tag that already exists.
func TestReleaseWorkflowSerialisesConcurrentMerges(t *testing.T) {
	workflow := releaseWorkflow(t)

	if !strings.Contains(workflow, "group: release") {
		t.Error("the release workflow declares no concurrency group — two merges would race on the last tag")
	}
	if !strings.Contains(workflow, "cancel-in-progress: false") {
		t.Error("the release workflow allows cancellation — a half-finished release is worse than a queued one")
	}
}

// A closed request is not a merged one. Without the guard, closing a request without
// merging would tag main at whatever it happens to be.
func TestReleaseWorkflowOnlyRunsOnAMergedRequest(t *testing.T) {
	if !strings.Contains(releaseWorkflow(t), "github.event.pull_request.merged == true") {
		t.Error("the release workflow does not check that the request was merged rather than merely closed")
	}
}

// The injection one. A label, a title and a body are text somebody else wrote, and `${{ }}`
// expanded inside a shell script is that text becoming executable. Untrusted values reach
// the scripts through the environment.
func TestReleaseWorkflowNeverExpandsUntrustedTextInsideAScript(t *testing.T) {
	lines := runScriptLines(releaseWorkflow(t))

	// Without this the test passes by finding nothing, which is the shape of a green that
	// asserts nothing at all.
	if len(lines) == 0 {
		t.Fatal("no run: script found in the release workflow — this test would pass vacuously")
	}

	for _, line := range lines {
		if strings.Contains(line, "${{") {
			t.Errorf("a run: script expands %q — untrusted text must arrive through env:, not be interpolated into the shell", strings.TrimSpace(line))
		}
	}
}

// runScriptLines returns the lines belonging to `run: |` blocks, which is where an
// expansion becomes executable. Indentation is the block: a line indented further than its
// `run:` key is inside it, and the first line indented no further ends it.
func runScriptLines(workflow string) []string {
	var inside []string
	depth := -1

	for _, line := range strings.Split(workflow, "\n") {
		indent := len(line) - len(strings.TrimLeft(line, " "))

		if depth >= 0 {
			if strings.TrimSpace(line) == "" {
				continue
			}
			if indent > depth {
				inside = append(inside, line)
				continue
			}
			depth = -1
		}
		if strings.HasSuffix(strings.TrimRight(line, " "), "run: |") {
			depth = indent
		}
	}
	return inside
}
