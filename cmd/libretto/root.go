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
//	4  $GOMODCACHE/<module>@<version>       the payload `go install` brought down
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
// **Rung 4 is the module cache, and that is the whole distribution story.** The payload ships
// inside the Go module, so `go install <module>/cmd/libretto@latest` downloads it along with the
// binary, checks it against the checksum database, and puts it under a path with the version in
// it. Nothing is extracted, nothing is unpacked, and no release asset is involved.
//
// It resolves to *this binary's own* version, from build info — so the payload linked is the one
// that shipped with the command doing the linking, and the two can never be a version apart.
//
// The path can be missing: `go clean -modcache` removes it. `libretto install` re-downloads and
// repairs, and `doctor` reports it.
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

	return resolveRoot(os.Getenv(EnvRoot), compileTime, wd)
}

func resolveRoot(override, compileTime, wd string) (string, error) {
	if override != "" {
		return override, nil
	}
	if compileTime != "" && isRepo(compileTime) {
		return compileTime, nil
	}
	if wd != "" && isRepo(wd) && isThisModule(wd) {
		return wd, nil
	}
	if cache := dist.ModCache(); cache != "" {
		return dist.Dir(cache, moduleImportPath(), buildVersion(version)), nil
	}
	return "", fmt.Errorf("no payload found and no module cache to look in")
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

// payloadPresent reports whether root holds a payload to link from.
//
// One of the three kind directories is enough: a tree with skills/ and no commands/ is a
// legitimate payload, and requiring all three would refuse a future release that dropped one.
func payloadPresent(root string) bool {
	for _, kind := range []string{"skills", "agents", "commands"} {
		if info, err := os.Stat(filepath.Join(root, kind)); err == nil && info.IsDir() {
			return true
		}
	}
	return false
}

// needsPayload reports whether the command about to run reads the payload tree.
//
// The panel and every command that scans the tree do. `models` does not — it reads the
// *target's* agents directory, so it works on a machine with nothing installed, which is
// exactly when somebody might want to see what model an agent is on.
//
// `upgrade` is what fixes a missing payload, so it cannot be gated on having one.
func needsPayload(args []string) bool {
	if len(args) == 0 {
		return true // the panel shows the tree's state
	}
	switch args[0] {
	case "models", "update", "loop", "metrics", "land":
		// `update` is what fixes a missing payload on an installed copy, so it cannot be
		// gated on having one. In a checkout the payload is the tree and is always there.
		//
		// `loop` reads the *project's* .agents/changes/<change>/tasks.md and relaunches a
		// session that resolves its own skills. Gating it on this repository's payload tree
		// would refuse the loop on every machine that installed the binary and nothing else
		// — which is the machine the loop is for.
		//
		// `land` reads the project being landed, never ~/.claude or the payload tree —
		// the machine record-work invokes it on is exactly one that installed the binary
		// and nothing else.
		return false
	default:
		return true
	}
}
