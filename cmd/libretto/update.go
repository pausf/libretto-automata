package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pausf/libretto-automata/internal/dist"
	"github.com/pausf/libretto-automata/internal/target"
)

// updateTimeout bounds the whole thing: a redirect, two downloads, an extraction and a
// compile. Generous, because the user asked for it and is waiting.
const updateTimeout = 10 * time.Minute

// defaultClient is the HTTP client every release call uses.
//
// A timeout on the client as well as the context: a context bounds the whole upgrade, and this
// bounds one stalled connection inside it. ponytail: no retry, no backoff. A failed upgrade is
// re-run by typing the command again, and a retry loop around a download is machinery for a
// failure the user can see.
func defaultClient() *http.Client { return &http.Client{Timeout: 2 * time.Minute} }

// updateDeps are upgrade's four outside effects, taken as parameters.
//
// Not for the sake of abstraction — there is one real implementation and there will not be a
// second. It is so the *order* is provable: a spy that records the steps is the only way to
// assert that the payload is activated before the binary is touched, and asserting on an error
// message would pass for an implementation that did it backwards and cleaned up.
type updateDeps struct {
	running string // the version this binary reports

	latest func(ctx context.Context) (string, error)
	// install brings down the binary and the payload together — one `go install`, because they
	// travel in the same module. What used to be two steps with an ordering constraint between
	// them is one step that cannot be half-done.
	install func(ctx context.Context, version string) error
	// where the payload for a version lands, for the report and for the relink
	dirFor func(version string) string
	relink func(root string) error
}

// fromRelease moves an installed copy to the newest release.
//
// **The order is fixed: payload, then binary, then relink.** A new binary reading an old
// payload is a state nobody can reason about — the same argument the update flow already makes
// about relinking before compiling. And nothing downstream runs if the payload step failed,
// which is what leaves the previous version active rather than half-replaced.
//
// **Relinking last is not redundant.** `current` moving means every existing link keeps
// resolving, but a release that *adds* an item leaves that item unlinked, and nothing would
// tell the user. That was a queued complaint of its own before this command existed.
//
// It says `git` nowhere, in success or in failure. That is the whole reason this command
// exists rather than `update` growing a branch.
func fromRelease(ctx context.Context, out io.Writer, d updateDeps) error {
	version, err := d.latest(ctx)
	if err != nil {
		return fmt.Errorf("could not find the newest version: %w", err)
	}

	if version == d.running {
		fmt.Fprintf(out, "up to date at %s\n", d.running)
		return nil
	}

	from := d.running
	if from == "" || from == "dev" {
		from = "an unidentified build"
	}
	fmt.Fprintf(out, "update   %s → %s\n", from, version)

	// One step, and it either happened or it did not. The binary and the payload are in the same
	// module, so there is no window where one is new and the other is old — which is the state
	// the previous design needed a fixed step order to avoid.
	if err := d.install(ctx, version); err != nil {
		return fmt.Errorf("nothing was changed: %w", err)
	}

	root := d.dirFor(version)
	fmt.Fprintf(out, "payload  %s\n", root)
	fmt.Fprintln(out, "         this process is still the old binary — run it again to use the new one")

	if err := d.relink(root); err != nil {
		return fmt.Errorf("%s is installed but not every item could be linked: %w", version, err)
	}
	return nil
}

// updateFromRelease wires fromRelease to the real world. `update` calls it for an installed
// copy; a checkout takes the other branch and pulls.
//
// There is no second command. `install` links the payload and `update` brings the installation
// up to date — those are the two words the audience already knows, and for them the meaning is
// one thing either way: put my installation on the newest version. Which mechanism runs is
// decided by how the tool got here, which the person who put it there knows perfectly well.
//
// An earlier session split this into `update` and `upgrade` and argued that a command whose
// mechanism depends on invisible state has unpredictable failures. The state is not invisible —
// it is "did you clone this or install it" — and the split bought two commands, two menu rows
// and a pair of mutual refusals to explain.
func updateFromRelease(ctx context.Context) error {
	cache := dist.ModCache()
	if cache == "" {
		return fmt.Errorf("cannot find the module cache — set GOMODCACHE or GOPATH")
	}
	module := moduleImportPath()

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	return fromRelease(ctx, os.Stdout, updateDeps{
		running: version,
		latest: func(ctx context.Context) (string, error) {
			return dist.Latest(ctx, defaultClient(), dist.DefaultProxy, module)
		},
		install: func(ctx context.Context, v string) error {
			return dist.Install(ctx, module, v)
		},
		dirFor: func(v string) string { return dist.Dir(cache, module, v) },
		// The new version's directory, not payloadRoot(): this process reports the old version,
		// so resolving again would link the payload it is already on.
		relink: func(root string) error {
			return install(root, target.Resolve(target.ClaudeTool, target.GlobalScope, ""))
		},
	})
}

// moduleImportPath is this module, taken from the one place it is written down.
func moduleImportPath() string {
	return strings.TrimPrefix(moduleLine, "module ")
}
