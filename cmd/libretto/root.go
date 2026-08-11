package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// EnvRoot overrides where the clone is. It is taken as given, with no validation: an
// escape hatch that checks its answer is a hatch that can refuse the one case you
// needed it for.
const EnvRoot = "LIBRETTO_ROOT"

// BootstrapDir is where the clone lands when there is not one already, under $HOME.
//
// A dotdir in $HOME rather than ~/.local/share, because this is a working git
// repository the user is expected to cd into and edit — not opaque application data.
const BootstrapDir = ".libretto-automata"

// moduleLine identifies a clone of this project by its go.mod. Any git repository would
// otherwise satisfy the working-directory rung, and `libretto install` run inside an
// unrelated project would go looking for a payload that project does not have.
const moduleLine = "module github.com/pausf/libretto-automata"

// repoRoot locates the clone this binary links from, in four rungs:
//
//	1  LIBRETTO_ROOT          absolute override, unvalidated
//	2  the compile-time path  when it is a repo — development, `make build`
//	3  the working directory  when it is a repo of this module
//	4  ~/.libretto-automata   the bootstrap clone, whether or not it exists yet
//
// Rung 2 used to accept any directory with a go.mod beside it. Under `go install` that
// is the read-only module cache, which has a go.mod — so it won, and every link pointed
// into a versioned cache directory that the next install orphaned. The probe is .git
// now, because everything this root is used for needs git: the pull, the rebuild
// decision, the release check.
//
// Rung 4 returns a path that may not exist. That is the caller's cue to bootstrap, and
// it beats the old fallback of returning the working directory — which quietly linked
// ~/.claude against whatever the user happened to be cd'd into.
// The rungs are resolved by resolveRoot, which takes its inputs rather than reading
// them. runtime.Caller is fixed at compile time, so inside this repository rung 2 always
// wins and rungs 3 and 4 are unreachable from a test — the exact rungs that only matter
// once the binary lives somewhere else. A seam here is not ceremony; without it the
// `go install` behaviour is the part with no proof.
func repoRoot() (string, error) {
	var compileTime string
	if _, file, _, ok := runtime.Caller(0); ok {
		// cmd/libretto/root.go -> repo root
		compileTime = filepath.Dir(filepath.Dir(filepath.Dir(file)))
	}
	wd, _ := os.Getwd()

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot locate the clone: %w", err)
	}
	return resolveRoot(os.Getenv(EnvRoot), compileTime, wd, home)
}

func resolveRoot(override, compileTime, wd, home string) (string, error) {
	if override != "" {
		return override, nil
	}
	if compileTime != "" && isRepo(compileTime) {
		return compileTime, nil
	}
	if wd != "" && isRepo(wd) && isThisModule(wd) {
		return wd, nil
	}
	if home == "" {
		return "", fmt.Errorf("no clone found and no home directory to bootstrap into")
	}
	return filepath.Join(home, BootstrapDir), nil
}

// isRepo reports whether dir is the top of a git checkout.
//
// A .git file counts, not just a directory: that is what a worktree has, and this flow
// recommends worktrees. Refusing one would break the tool exactly where its own advice
// put you.
func isRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

func isThisModule(dir string) bool {
	body, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return false
	}
	return strings.HasPrefix(strings.TrimSpace(string(body)), moduleLine)
}
