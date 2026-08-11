package repo

import (
	"context"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
)

// A destination with anything in it is refused, never merged into and never forced.
// The tool's whole promise is that it does not write over what it did not create, and
// its own clone directory does not get an exception.
func TestCloneRefusesNonEmptyDestination(t *testing.T) {
	dest := t.TempDir()
	if err := os.WriteFile(filepath.Join(dest, "somebody-elses-file"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Clone(context.Background(), "https://example.invalid/repo.git", dest)
	if err == nil {
		t.Fatal("Clone accepted a destination with a file in it")
	}
	if !strings.Contains(err.Error(), dest) {
		t.Errorf("the error does not name the destination: %v", err)
	}

	// And it left the file alone.
	if _, statErr := os.Stat(filepath.Join(dest, "somebody-elses-file")); statErr != nil {
		t.Errorf("the refused clone touched the existing file: %v", statErr)
	}
}

// An existing empty directory is fine — that is a destination nobody is using, and
// refusing it would mean `mkdir ~/.libretto-automata` permanently breaks bootstrap.
func TestCloneAcceptsAnExistingEmptyDirectory(t *testing.T) {
	src := gitRepoWithACommit(t)
	dest := filepath.Join(t.TempDir(), "empty")
	if err := os.Mkdir(dest, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := Clone(context.Background(), src, dest); err != nil {
		t.Fatalf("Clone into an empty directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, ".git")); err != nil {
		t.Errorf("no .git in the clone: %v", err)
	}
}

func TestCloneCreatesMissingDestination(t *testing.T) {
	src := gitRepoWithACommit(t)
	// Two levels down, because bootstrap can be pointed anywhere by LIBRETTO_ROOT and
	// git only creates the last component itself.
	dest := filepath.Join(t.TempDir(), "nested", "libretto-automata")

	if err := Clone(context.Background(), src, dest); err != nil {
		t.Fatalf("Clone into a missing path: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "marker")); err != nil {
		t.Errorf("the clone has no content: %v", err)
	}
}

func TestCloneHonoursACancelledContext(t *testing.T) {
	src := gitRepoWithACommit(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := Clone(ctx, src, filepath.Join(t.TempDir(), "x")); err == nil {
		t.Error("Clone ignored a cancelled context")
	}
}

// The URL is derived from the module path, so there is one declaration of where this
// project lives rather than a constant sitting beside go.mod waiting to disagree with it.
func TestModuleURLDerivesFromBuildInfo(t *testing.T) {
	info := &debug.BuildInfo{}
	info.Main.Path = "github.com/pausf/libretto-automata"

	if got, want := moduleURL(info), "https://github.com/pausf/libretto-automata.git"; got != want {
		t.Errorf("moduleURL = %q, want %q", got, want)
	}
}

// A test binary has build info but no main module path, and `go install` of a fork has
// a different one. Neither may produce an empty URL that git would reject with something
// unreadable.
func TestModuleURLFallsBackWhenBuildInfoIsUseless(t *testing.T) {
	want := "https://github.com/pausf/libretto-automata.git"

	if got := moduleURL(nil); got != want {
		t.Errorf("moduleURL(nil) = %q, want %q", got, want)
	}
	if got := moduleURL(&debug.BuildInfo{}); got != want {
		t.Errorf("moduleURL(empty) = %q, want %q", got, want)
	}
}

// gitRepoWithACommit is a clonable source: a repository with one commit in it. An empty
// repository clones with a warning and no working tree, which is not what bootstrap does.
func gitRepoWithACommit(t *testing.T) string {
	t.Helper()
	root := gitRepo(t)
	if err := os.WriteFile(filepath.Join(root, "marker"), []byte("payload\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(t, root, "add", "marker")
	run(t, root, "commit", "-q", "-m", "first")
	return root
}
