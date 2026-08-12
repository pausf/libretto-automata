package main

import (
	"os/exec"
	"strings"
	"testing"
)

// The label check exists to move release.yml's refusal from after the merge to before
// it. Read as text for the reason gates_test.go gives; the script itself is extracted
// and run under bash, because "refuses zero and two" is behaviour, not shape.

func labelWorkflow(t *testing.T) string {
	t.Helper()
	return repoFile(t, ".github/workflows/require-release-label.yml")
}

// runLabelScript runs the workflow's own script — the same bytes CI runs — with LABELS
// set the way the env: block would set it.
func runLabelScript(t *testing.T, labels string) error {
	t.Helper()
	lines := runScriptLines(labelWorkflow(t))
	if len(lines) == 0 {
		t.Fatal("no run: script found in the label workflow — these tests would pass vacuously")
	}
	cmd := exec.Command("bash", "-c", strings.Join(lines, "\n"))
	cmd.Env = append(cmd.Environ(), "LABELS="+labels)
	return cmd.Run()
}

// The reason this is a separate file: `labeled` and `unlabeled` are not in the default
// pull_request event set, and without them adding the label leaves the check red until
// the next push — a check that can only be satisfied by pushing is a check people
// learn to work around.
func TestLabelWorkflowRunsWhenLabelsChange(t *testing.T) {
	workflow := labelWorkflow(t)

	for _, event := range []string{"labeled", "unlabeled"} {
		if !strings.Contains(workflow, event) {
			t.Errorf("the label workflow does not trigger on %q — the check would not re-evaluate when labels change", event)
		}
	}
}

func TestLabelScriptRefusesZeroLabels(t *testing.T) {
	for _, labels := range []string{"", "bug documentation"} {
		if runLabelScript(t, labels) == nil {
			t.Errorf("the script passed with LABELS=%q — a request with no release: label must fail the check", labels)
		}
	}
}

func TestLabelScriptRefusesTwoLabels(t *testing.T) {
	if runLabelScript(t, "release:patch release:minor") == nil {
		t.Error("the script passed with two release labels — one bump, one label")
	}
}

func TestLabelScriptAcceptsExactlyOneLabel(t *testing.T) {
	for _, bump := range []string{"release:patch", "release:minor", "release:major"} {
		if err := runLabelScript(t, bump+" bug"); err != nil {
			t.Errorf("the script failed with exactly one release label (%s): %v", bump, err)
		}
	}
}

func TestLabelWorkflowIsReadOnly(t *testing.T) {
	workflow := labelWorkflow(t)

	if !strings.Contains(workflow, "permissions: {}") {
		t.Error("the label workflow does not declare an empty permissions block — a job that touches nothing gets no permissions at all")
	}
	if strings.Contains(workflow, "write") {
		t.Error("the label workflow mentions write access — nothing here has any business writing")
	}
}

// The labels arrive in the event payload. A checkout would run contributor code for no
// reason, and a job with no permissions should also have no tree to be tempted by.
func TestLabelWorkflowNeedsNoCheckout(t *testing.T) {
	if strings.Contains(labelWorkflow(t), "actions/checkout") {
		t.Error("the label workflow checks out the tree — the labels are in the event payload and the job needs no code")
	}
}

// Same rule as release.yml: a label is text somebody else wrote, and `${{ }}` expanded
// inside a shell script is that text becoming executable.
func TestLabelWorkflowNeverExpandsUntrustedTextInsideAScript(t *testing.T) {
	lines := runScriptLines(labelWorkflow(t))

	if len(lines) == 0 {
		t.Fatal("no run: script found in the label workflow — this test would pass vacuously")
	}
	for _, line := range lines {
		if strings.Contains(line, "${{") {
			t.Errorf("a run: script expands %q — untrusted text must arrive through env:, not be interpolated into the shell", strings.TrimSpace(line))
		}
	}
}
