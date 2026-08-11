package repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
)

// fallbackModule is the one place this project's own location is written down as a
// literal. Everything else derives it from the module path.
const fallbackModule = "github.com/pausf/libretto-automata"

// ModuleURL is the https clone URL for the module this binary was built from.
func ModuleURL() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return moduleURL(nil)
	}
	return moduleURL(info)
}

// moduleURL takes the build info rather than reading it, so the fallback branch has a
// test. A function calling debug.ReadBuildInfo() inside itself can only ever be tested
// for whichever answer the test binary happens to give.
//
// Deriving the URL means a fork installed with `go install <fork>/cmd/libretto@latest`
// bootstraps from the fork, which is the behaviour somebody working on a fork wants and
// would otherwise have to know to set LIBRETTO_ROOT for.
func moduleURL(info *debug.BuildInfo) string {
	path := fallbackModule
	if info != nil && info.Main.Path != "" {
		path = info.Main.Path
	}
	return "https://" + path + ".git"
}

// Clone creates the working clone the tool needs at dest.
//
// Real `git clone` for the same reason Pull shells out: the fetch has to work with the
// user's credential helper, ssh agent, proxy and .gitconfig, and git is the only
// implementation guaranteed to honour all of them.
//
// A destination holding anything at all is refused. Not merged into, not forced — the
// same promise the linker makes, and the tool's own directory does not get an exception.
// A missing destination is created, including its parents, because LIBRETTO_ROOT can
// point anywhere and git only creates the last component itself.
func Clone(ctx context.Context, url, dest string) error {
	entries, err := os.ReadDir(dest)
	switch {
	case err == nil && len(entries) > 0:
		// Naming what was found is what makes this actionable. "refusing to clone" on
		// its own sends the user to look at a directory they cannot see the problem in,
		// and on macOS the problem is often a single .DS_Store.
		return fmt.Errorf("refusing to clone into %s — it already contains %s", dest, entries[0].Name())
	case err != nil && !errors.Is(err, os.ErrNotExist):
		return fmt.Errorf("cannot read %s: %w", dest, err)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "git", "clone", "--quiet", url, dest)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if text := strings.TrimSpace(string(out)); text != "" {
			return fmt.Errorf("%s", text)
		}
		return err
	}
	return nil
}
