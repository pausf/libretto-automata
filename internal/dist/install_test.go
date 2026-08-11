package dist

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallActivatesTheNewVersion(t *testing.T) {
	body := goodTarball(t)
	srv := assetServer(t, "v0.4.0", body, digestOf(body))
	base := t.TempDir()

	if err := Install(context.Background(), srv.Client(), srv.URL, base, "v0.4.0"); err != nil {
		t.Fatalf("Install: %v", err)
	}

	if got := Active(base); got != "v0.4.0" {
		t.Errorf("active version is %q, want v0.4.0", got)
	}
	// The whole chain resolves: current → v0.4.0 → a real file. That is what ~/.claude
	// links through, and a link needs something real at the end.
	if _, err := os.Stat(filepath.Join(Current(base), "skills", "marker.md")); err != nil {
		t.Errorf("the payload does not resolve through current: %v", err)
	}
}

// A failure at any step leaves the previous version active and no partial directory behind.
// An interrupted upgrade has to be a no-op, not a half-installed tree.
func TestInstallLeavesNoPartialVersionOnFailure(t *testing.T) {
	base := t.TempDir()
	mkVersion(t, base, "v0.3.0")
	if err := Activate(base, "v0.3.0"); err != nil {
		t.Fatal(err)
	}

	body := goodTarball(t)
	srv := assetServer(t, "v0.4.0", body, digestOf([]byte("wrong")))

	if err := Install(context.Background(), srv.Client(), srv.URL, base, "v0.4.0"); err == nil {
		t.Fatal("Install accepted a mismatched checksum")
	}

	if got := Active(base); got != "v0.3.0" {
		t.Errorf("active version is %q after a failed install, want v0.3.0", got)
	}
	if _, err := os.Stat(VersionDir(base, "v0.4.0")); !os.IsNotExist(err) {
		t.Error("a failed install left a v0.4.0 directory behind")
	}
	// And the version that was running still resolves.
	if _, err := os.Stat(filepath.Join(Current(base), "skills", "marker.md")); err != nil {
		t.Errorf("the previously active payload no longer resolves: %v", err)
	}
}

// Install keeps the version it replaced, and drops the one before that. Rolling back has to
// be a swap, and a directory nobody prunes is the other failure.
func TestInstallKeepsTheReplacedVersionAndDropsTheRest(t *testing.T) {
	base := t.TempDir()
	for _, v := range []string{"v0.1.0", "v0.2.0", "v0.3.0"} {
		mkVersion(t, base, v)
	}
	if err := Activate(base, "v0.3.0"); err != nil {
		t.Fatal(err)
	}

	body := goodTarball(t)
	srv := assetServer(t, "v0.4.0", body, digestOf(body))

	if err := Install(context.Background(), srv.Client(), srv.URL, base, "v0.4.0"); err != nil {
		t.Fatal(err)
	}

	for _, kept := range []string{"v0.4.0", "v0.3.0"} {
		if _, err := os.Stat(VersionDir(base, kept)); err != nil {
			t.Errorf("%s should have been kept: %v", kept, err)
		}
	}
	for _, gone := range []string{"v0.1.0", "v0.2.0"} {
		if _, err := os.Stat(VersionDir(base, gone)); !os.IsNotExist(err) {
			t.Errorf("%s should have been pruned", gone)
		}
	}
}

// Installing what is already active is a no-op that succeeds. `upgrade` run twice must not
// fail the second time, and must not make a network request it does not need.
func TestInstallTheActiveVersionIsANoOp(t *testing.T) {
	base := t.TempDir()
	mkVersion(t, base, "v0.4.0")
	if err := Activate(base, "v0.4.0"); err != nil {
		t.Fatal(err)
	}

	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
	}))
	t.Cleanup(srv.Close)

	if err := Install(context.Background(), srv.Client(), srv.URL, base, "v0.4.0"); err != nil {
		t.Fatalf("re-installing the active version: %v", err)
	}
	if requests != 0 {
		t.Errorf("made %d requests re-installing the active version", requests)
	}
	if got := Active(base); got != "v0.4.0" {
		t.Errorf("active version is %q", got)
	}
}

// A failed activation must not prune. Pruning after a step that did not happen is how the
// version you were running disappears while the new one is not there either.
func TestInstallDoesNotPruneWhenActivationFails(t *testing.T) {
	base := t.TempDir()
	mkVersion(t, base, "v0.1.0")
	mkVersion(t, base, "v0.2.0")
	mkVersion(t, base, "v0.3.0")
	if err := Activate(base, "v0.3.0"); err != nil {
		t.Fatal(err)
	}

	body := goodTarball(t)
	srv := assetServer(t, "v0.4.0", body, digestOf([]byte("wrong")))
	_ = Install(context.Background(), srv.Client(), srv.URL, base, "v0.4.0")

	for _, v := range []string{"v0.1.0", "v0.2.0", "v0.3.0"} {
		if _, err := os.Stat(VersionDir(base, v)); err != nil {
			t.Errorf("%s was pruned after a failed install: %v", v, err)
		}
	}
}
