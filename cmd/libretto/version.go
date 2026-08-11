package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pausf/libretto-automata/internal/dist"
	"github.com/pausf/libretto-automata/internal/repo"
)

// checkTimeout bounds the release check.
//
// Five seconds, and short on purpose: this is speculative work nobody asked for. The
// bootstrap clone gets minutes because the user is waiting for the payload; nobody is
// waiting to be told they are up to date.
const checkTimeout = 5 * time.Second

// askLatest is the "what is the newest release" call for the mode this machine is in, behind
// the shared cache.
//
// **Two questions, not two answers to one.** A checkout asks *its remote* what it has tagged;
// an installed copy asks *the project* what it has released. Only one is ever asked on a given
// machine, so these are not two implementations that can disagree — and the cache policy they
// share is repo's, written once.
func askLatest(ctx context.Context, root string) (string, error) {
	if isRepo(root) {
		return repo.CheckedLatest(ctx, root, repo.CheckTTL)
	}

	// The installed copy keeps its answer beside the payloads. Not in .git/ — there is no
	// .git — and without somewhere to keep it the panel would make an HTTP call on every
	// single launch, which is the hang the cache exists to prevent.
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base := dist.Base(home)
	return repo.Cached(ctx, filepath.Join(base, checkFile), repo.CheckTTL,
		func(ctx context.Context) (string, error) {
			return dist.Latest(ctx, defaultClient(), forgeHost)
		})
}

// checkFile is where an installed copy keeps the cached answer. A dotfile, so it is not
// mistaken for a version directory — `dist.Versions` ignores it either way, and belt and
// braces on a directory the tool prunes is cheap.
const checkFile = ".update-check"

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

	latest, err := askLatest(ctx, root)
	if err != nil || !repo.IsNewer(latest, running) {
		return ""
	}
	// Always `update`, in either mode. The row does not have to name a command per machine
	// any more, which is one of the things collapsing two commands into one bought.
	return fmt.Sprintf("%s → %s available · choose update", running, latest)
}

// releaseLine is `doctor`'s version: one line, always something, never silence.
//
// The panel stays quiet when there is nothing to say; a diagnostic does not get to. The
// user typed a command that went looking, so "I could not find out" is one of its
// legitimate answers — and printing nothing would read as "you are up to date", which is a
// claim nobody verified.
//
// ask is a parameter so this has a test that does not reach a remote. `doctor` passes
// Shell.LatestTag directly rather than the cached path: it checks live, because the user is
// waiting for a diagnosis, and going through the cache would both write a file and swallow
// the very error this line exists to report.
func releaseLine(running string, ask func(context.Context) (string, error)) string {
	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()

	latest, err := ask(ctx)
	switch {
	case err != nil:
		return "could not check for a newer release — " + firstLine(err.Error())
	case latest == "":
		return "the remote has no releases to offer"
	case repo.IsNewer(latest, running):
		return fmt.Sprintf("%s → %s available · run `%s update`", running, latest, invokedAs())
	case running == latest:
		return "up to date at " + running
	default:
		// Either ahead of the remote, or a version that will not parse. Both get the
		// facts and no ranking: telling somebody running `dev` that they are out of date
		// is a guess presented as a fact.
		return fmt.Sprintf("running %s · the latest release is %s", running, latest)
	}
}

// firstLine keeps a git error to one line. Its multi-line advice is for somebody at a
// prompt, and doctor's report is a table.
func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
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
