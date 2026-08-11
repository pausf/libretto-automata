package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/pausf/libretto-automata/internal/repo"
)

// checkTimeout bounds the release check.
//
// Five seconds, and short on purpose: this is speculative work nobody asked for. The
// bootstrap clone gets minutes because the user is waiting for the payload; nobody is
// waiting to be told they are up to date.
const checkTimeout = 5 * time.Second

// releaseNotice is the panel's row, or "" when there is nothing to say.
//
// Cached — the panel must not pay for a network call on every launch. `doctor` is the
// command that checks live, because the user typed a diagnostic and can afford to wait.
//
// The comparison is repo.IsNewer, not a second one written here. Formatting is this side's
// job; deciding what "newer" means is not.
func releaseNotice(root, running string) string {
	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()

	latest, err := repo.CheckedLatest(ctx, root, repo.CheckTTL)
	if err != nil || !repo.IsNewer(latest, running) {
		return ""
	}
	return fmt.Sprintf("%s → %s available · choose update", running, latest)
}

// resolveVersion decides what the binary reports: the ldflags stamp, then the module
// version the toolchain recorded, then `dev`.
//
// ldflags first, and not the other way round. `make build` stamps `git describe --tags
// --always --dirty`, and build info cannot know a tree is dirty or how many commits past
// the tag it sits — so preferring build info would quietly turn v0.2.0-3-gabc123-dirty
// back into a clean v0.2.0.
//
// The fallback is what makes `go install pkg@version` honest: that build gets no ldflags
// but does record the module version, so it can say v0.3.0 instead of reporting `dev` to
// everyone who did not build from a clone. `(devel)` is not a version — it is what the
// toolchain records for a build from a working tree — so it falls through.
//
// It takes the build info rather than reading it, so both branches have a test. A
// function calling debug.ReadBuildInfo() inside itself can only be tested for whichever
// answer the test binary happens to give.
func resolveVersion(stamped string, info *debug.BuildInfo) string {
	if stamped != "" && stamped != "dev" {
		return stamped
	}
	if info != nil && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}

// buildVersion is resolveVersion against this process's own build info.
func buildVersion(stamped string) string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return resolveVersion(stamped, nil)
	}
	return resolveVersion(stamped, info)
}
