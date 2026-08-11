package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pausf/libretto-automata/internal/dist"
	"github.com/pausf/libretto-automata/internal/target"
)

// upgradeTimeout bounds the whole thing: a redirect, two downloads, an extraction and a
// compile. Generous, because the user asked for it and is waiting.
const upgradeTimeout = 10 * time.Minute

// defaultClient is the HTTP client every release call uses.
//
// A timeout on the client as well as the context: a context bounds the whole upgrade, and this
// bounds one stalled connection inside it. ponytail: no retry, no backoff. A failed upgrade is
// re-run by typing the command again, and a retry loop around a download is machinery for a
// failure the user can see.
func defaultClient() *http.Client { return &http.Client{Timeout: 2 * time.Minute} }

// forgeHost is where releases are published. Derived from the module path, so there is one
// declaration of where this project lives.
const forgeHost = "https://github.com/pausf/libretto-automata"

// upgradeDeps are upgrade's four outside effects, taken as parameters.
//
// Not for the sake of abstraction — there is one real implementation and there will not be a
// second. It is so the *order* is provable: a spy that records the steps is the only way to
// assert that the payload is activated before the binary is touched, and asserting on an error
// message would pass for an implementation that did it backwards and cleaned up.
type upgradeDeps struct {
	running string // the version this binary reports
	base    string // where payloads live, for the report

	latest         func(ctx context.Context) (string, error)
	installPayload func(ctx context.Context, tag string) error
	replaceBinary  func(ctx context.Context, tag string) (note string, err error)
	relink         func() error
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
func fromRelease(ctx context.Context, out io.Writer, d upgradeDeps) error {
	tag, err := d.latest(ctx)
	if err != nil {
		return fmt.Errorf("could not find the newest release: %w", err)
	}
	if tag == "" {
		return fmt.Errorf("could not find the newest release: none is published")
	}

	if tag == d.running {
		fmt.Fprintf(out, "up to date at %s\n", d.running)
		return nil
	}

	from := d.running
	if from == "" || from == "dev" {
		from = "nothing installed"
	}
	fmt.Fprintf(out, "upgrade  %s → %s\n", from, tag)

	if err := d.installPayload(ctx, tag); err != nil {
		return fmt.Errorf("the payload was not installed: %w", err)
	}
	fmt.Fprintf(out, "payload  %s\n", dist.VersionDir(d.base, tag))

	// A binary that cannot be replaced does not undo a payload that installed correctly. The
	// same split `rebuild` already makes for an unwritable destination, and the same reason:
	// rolling back work that succeeded loses more than it saves.
	if note, berr := d.replaceBinary(ctx, tag); berr != nil {
		fmt.Fprintf(out, "binary   unchanged — %v\n", berr)
		fmt.Fprintf(out, "         the payload is on %s; the command you run is still the old build\n", tag)
	} else if note != "" {
		fmt.Fprintf(out, "binary   %s\n", note)
	}

	if err := d.relink(); err != nil {
		return fmt.Errorf("the payload is on %s but not every item could be linked: %w", tag, err)
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
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	base := dist.Base(home)

	ctx, cancel := context.WithTimeout(ctx, upgradeTimeout)
	defer cancel()

	return fromRelease(ctx, os.Stdout, upgradeDeps{
		running: version,
		base:    base,
		latest: func(ctx context.Context) (string, error) {
			return dist.Latest(ctx, defaultClient(), forgeHost)
		},
		installPayload: func(ctx context.Context, tag string) error {
			return dist.Install(ctx, defaultClient(), forgeHost, base, tag)
		},
		replaceBinary: goInstall,
		relink: func() error {
			// Resolved again, deliberately: `current` has just moved, and the target is
			// whatever scope the flags said.
			newRoot, rerr := payloadRoot()
			if rerr != nil {
				return rerr
			}
			return install(newRoot, target.Resolve(target.GlobalScope, ""))
		},
	})
}

// goInstall replaces the running command with the release's build.
//
// `go install` rather than a published binary, because Go is present by construction — it is
// how this command got here. Publishing per-platform binaries would remove that assumption and
// is deliberately out of scope.
//
// The word `go` appears in the command and not in the report. What the user is told is the path
// that changed.
func goInstall(ctx context.Context, tag string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	cmd := exec.CommandContext(ctx, "go", "install", moduleImportPath()+"/cmd/libretto@"+tag)
	if out, err := cmd.CombinedOutput(); err != nil {
		if text := strings.TrimSpace(string(out)); text != "" {
			return "", fmt.Errorf("%s", firstLine(text))
		}
		return "", err
	}
	return exe + " replaced", nil
}

// moduleImportPath is this module, taken from the one place it is written down.
func moduleImportPath() string {
	return strings.TrimPrefix(moduleLine, "module ")
}
