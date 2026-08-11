package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pausf/libretto-automata/internal/dist"
)

// EnvRoot overrides where the payload root is. It is taken as given, with no validation:
// an escape hatch that checks its answer is a hatch that can refuse the one case you
// needed it for.
const EnvRoot = "LIBRETTO_ROOT"

// moduleLine identifies a clone of this project by its go.mod. Any git repository would
// otherwise satisfy the working-directory rung, and `libretto install` run inside an
// unrelated project would go looking for a payload that project does not have.
const moduleLine = "module github.com/pausf/libretto-automata"

// payloadRoot locates the tree this binary links from, in four rungs:
//
//	1  LIBRETTO_ROOT                        absolute override, unvalidated
//	2  the compile-time path                when it is a checkout — development
//	3  the working directory                when it is a checkout of this module
//	4  ~/.local/share/libretto/current      the activated release
//
// **The payload root and the checkout are two different things**, and this tool conflated
// them for as long as a clone happened to be both. `go install` breaks the coincidence:
// there is a payload and no checkout. Rungs 2 and 3 answer both questions at once and
// rung 4 answers only the first — which is why `update` asks isRepo() before pulling and
// `upgrade` asks it before refusing.
//
// Rung 2 probes for `.git`, not for a `go.mod`. Under `go install` a `go.mod` matches the
// read-only module cache, which would win, and every link would point into a versioned
// cache directory that the next install orphans.
//
// Rung 4 returns a path that may not exist: nothing has been installed yet. That is not
// something to fix by cloning — `upgrade` downloads a release — and it still beats the
// fallback it replaced, which returned the working directory and quietly linked
// ~/.claude against whatever the user happened to be cd'd into.
//
// The rungs are resolved by resolveRoot, which takes its inputs rather than reading
// them. runtime.Caller is fixed at compile time, so inside this repository rung 2 always
// wins and rungs 3 and 4 are unreachable from a test — the exact rungs that only matter
// once the binary lives somewhere else. A seam here is not ceremony; without it the
// `go install` behaviour is the part with no proof.
func payloadRoot() (string, error) {
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
		return "", fmt.Errorf("no payload found and no home directory to look for one in")
	}
	return dist.Current(dist.Base(home)), nil
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
