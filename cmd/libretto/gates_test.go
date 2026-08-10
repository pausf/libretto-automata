package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The workflow and the Makefile are read as text rather than parsed.
//
// A YAML dependency added to assert on a config file is a dependency added to check a
// checker, and this repository puts a new dependency behind an ask. What these tests
// need is "does this command appear", which is a substring question.

func repoFile(t *testing.T, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", rel))
	if err != nil {
		t.Fatalf("reading %s: %v", rel, err)
	}
	return string(data)
}

// gates are the six AGENTS.md names, as the substrings that identify each one in a
// script. Kept here so both the workflow test and the Makefile test compare against
// one list rather than two that agree today.
var gates = []string{
	"gofmt -l .",
	"go vet ./...",
	"go test ./... -count=1",
	"scripts/check-payload",
	"spec-drift --self-test",
	"spec-drift --anchors",
}

func TestWorkflowRunsEveryGateAgentsNames(t *testing.T) {
	workflow := repoFile(t, ".github/workflows/gates.yml")

	for _, gate := range gates {
		if !strings.Contains(workflow, gate) {
			t.Errorf("the workflow does not run %q — AGENTS.md names six gates and CI must run all of them", gate)
		}
	}
}

// A workflow that can write is a workflow that can be made to write. Nothing here
// pushes, comments or commits, so nothing here needs more than read.
func TestWorkflowIsReadOnly(t *testing.T) {
	workflow := repoFile(t, ".github/workflows/gates.yml")

	if !strings.Contains(workflow, "contents: read") {
		t.Error("the workflow does not declare read-only contents permission")
	}
	for _, forbidden := range []string{"contents: write", "pull-requests: write", "packages: write"} {
		if strings.Contains(workflow, forbidden) {
			t.Errorf("the workflow asks for %q — nothing in it writes", forbidden)
		}
	}
}

// The one that fails open, and therefore the one that gets a real check rather than a
// substring.
//
// `gofmt -l` prints the unformatted files and exits **zero**. A step written as
// `run: gofmt -l .` is green on a repository nobody has formatted — the gate is
// present, named in the summary, and asserting nothing. So the step has to test the
// output, and this test proves the command it uses actually fails on bad input rather
// than trusting that it reads correctly.
func TestGofmtStepFailsOnUnformattedOutput(t *testing.T) {
	workflow := repoFile(t, ".github/workflows/gates.yml")

	// The step must do something with gofmt's output. A bare invocation is the bug.
	var step string
	for _, line := range strings.Split(workflow, "\n") {
		if strings.Contains(line, "gofmt -l .") {
			step = strings.TrimSpace(line)
		}
	}
	if step == "" {
		t.Fatal("no gofmt step in the workflow")
	}
	if strings.HasSuffix(step, "gofmt -l .") {
		t.Fatalf("the gofmt step is %q — gofmt -l exits zero with a list of unformatted files, so this passes on an unformatted repository", step)
	}

	// And the shell the step uses must genuinely fail on a listing. Proving the
	// command rather than the YAML: the same construction is run here against a file
	// that is deliberately misformatted.
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.go")
	if err := os.WriteFile(bad, []byte("package p\nfunc  F( ) {\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("bash", "-c", `test -z "$(gofmt -l .)"`)
	cmd.Dir = dir
	if err := cmd.Run(); err == nil {
		t.Error("the guard construction passed on an unformatted file — it would pass in CI too")
	}

	if err := exec.Command("gofmt", "-w", bad).Run(); err != nil {
		t.Fatal(err)
	}
	cmd = exec.Command("bash", "-c", `test -z "$(gofmt -l .)"`)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Errorf("the guard construction failed on a formatted file: %v", err)
	}
}

// Two lists that agree today is exactly the arrangement that stops agreeing without
// anybody noticing. `make gates` and the workflow run the same six commands, and this
// is what keeps that true after somebody adds a seventh to one of them.
func TestMakeGatesMatchesTheWorkflow(t *testing.T) {
	makefile := repoFile(t, "Makefile")

	// The recipe is the tab-indented run of lines under `gates:`. Splitting on blank
	// lines looked simpler and was wrong — a comment above the target owns the block.
	var recipe strings.Builder
	inTarget := false
	for _, line := range strings.Split(makefile, "\n") {
		switch {
		case strings.HasPrefix(line, "gates:"):
			inTarget = true
		case inTarget && strings.HasPrefix(line, "\t"):
			recipe.WriteString(line + "\n")
		case inTarget:
			inTarget = false
		}
	}
	target := recipe.String()
	if target == "" {
		t.Fatal("no `gates` target in the Makefile, or it has no recipe")
	}

	for _, gate := range gates {
		if !strings.Contains(target, gate) {
			t.Errorf("`make gates` does not run %q — the workflow does, and the two must be one list", gate)
		}
	}

	if !strings.Contains(makefile, ".PHONY") || !strings.Contains(makefilePhony(makefile), "gates") {
		t.Error("`gates` is missing from .PHONY — a file of that name would silently disable the target")
	}
}

func makefilePhony(makefile string) string {
	for _, line := range strings.Split(makefile, "\n") {
		if strings.HasPrefix(line, ".PHONY:") {
			return line
		}
	}
	return ""
}
