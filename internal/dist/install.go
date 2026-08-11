package dist

import (
	"context"
	"net/http"
)

// Install puts a version's payload on disk and makes it the active one.
//
// Four steps, and the order between the last two is what makes an interrupted upgrade a
// no-op:
//
//	fetch     download, verify, extract into <base>/<tag>   — leaves nothing on failure
//	activate  swap current, atomically
//	prune     drop everything but this version and the one it replaced
//
// **Nothing is pruned unless the activation succeeded.** Pruning after a step that did not
// happen is how the version the user was running disappears while the new one is not there
// either — and there is no recovering from that without a network.
//
// Installing the version that is already active is a success and makes no request. `upgrade`
// run twice must not fail the second time.
func Install(ctx context.Context, client *http.Client, host, base, tag string) error {
	previous := Active(base)
	if previous == tag {
		return nil
	}

	if err := fetch(ctx, client, host, base, tag, extractInto); err != nil {
		return err
	}
	if err := Activate(base, tag); err != nil {
		return err
	}

	// Keeping the version just replaced is what makes a rollback a symlink swap rather than
	// a download. A prune failure is not worth failing an upgrade that already worked: the
	// payload is installed and active, and the cost of a leftover directory is disk.
	keep := []string{tag}
	if previous != "" {
		keep = append(keep, previous)
	}
	_ = Prune(base, keep...)
	return nil
}
