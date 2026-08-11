package dist

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// entry is one thing to put in a test tarball. Zero Typeflag means a regular file.
type entry struct {
	name     string
	body     string
	typeflag byte
	linkname string
	mode     int64
}

// tarball builds a .tar.gz in memory. Building the archives here rather than committing
// fixture files is what makes the refusals below readable: the malicious entry is visible
// in the test that expects it to be refused.
func tarball(t *testing.T, entries ...entry) *bytes.Reader {
	t.Helper()

	var raw bytes.Buffer
	gz := gzip.NewWriter(&raw)
	tw := tar.NewWriter(gz)

	for _, e := range entries {
		mode := e.mode
		if mode == 0 {
			mode = 0o644
		}
		h := &tar.Header{
			Name:     e.name,
			Mode:     mode,
			Size:     int64(len(e.body)),
			Typeflag: e.typeflag,
			Linkname: e.linkname,
		}
		if e.typeflag == tar.TypeDir || e.typeflag == tar.TypeSymlink || e.typeflag == tar.TypeLink {
			h.Size = 0
		}
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if h.Size > 0 {
			if _, err := tw.Write([]byte(e.body)); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(raw.Bytes())
}

const noCeiling = 1 << 30

func TestExtractWritesFilesAndDirectories(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "v0.4.0")

	err := extract(tarball(t,
		entry{name: "skills/", typeflag: tar.TypeDir},
		entry{name: "skills/write-spec/SKILL.md", body: "# write-spec\n"},
		entry{name: "agents/spec-writer.md", body: "# spec-writer\n"},
	), dest, noCeiling)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	body, err := os.ReadFile(filepath.Join(dest, "skills", "write-spec", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "# write-spec\n" {
		t.Errorf("content = %q", body)
	}
	// A file whose parent directory has no entry of its own must still land: `tar czf`
	// includes directory entries, but nothing guarantees they come first.
	if _, err := os.Stat(filepath.Join(dest, "agents", "spec-writer.md")); err != nil {
		t.Errorf("a file whose parent had no entry was not written: %v", err)
	}
}

// The classic extraction escape. `../../.ssh/authorized_keys` inside an archive writes
// outside the destination, and this is the one thing that turns a download into a
// compromise.
func TestExtractRefusesAPathEscapingTheDestination(t *testing.T) {
	cases := map[string]string{
		"parent traversal":  "../escaped",
		"deep traversal":    "../../../escaped",
		"traversal midway":  "skills/../../escaped",
		"cleaned traversal": "a/b/../../../escaped",
	}

	for name, path := range cases {
		t.Run(name, func(t *testing.T) {
			base := t.TempDir()
			dest := filepath.Join(base, "v0.4.0")

			err := extract(tarball(t, entry{name: path, body: "owned"}), dest, noCeiling)
			if err == nil {
				t.Fatalf("extract accepted %q", path)
			}
			if !strings.Contains(err.Error(), path) {
				t.Errorf("the refusal does not name the entry: %v", err)
			}

			// And nothing was written anywhere above the destination.
			if _, statErr := os.Stat(filepath.Join(base, "escaped")); !os.IsNotExist(statErr) {
				t.Errorf("%q wrote outside the destination", path)
			}
		})
	}
}

// A destination whose name is a prefix of a sibling is the check people get wrong: a string
// comparison passes `/tmp/destevil` as being inside `/tmp/dest`.
func TestExtractRefusesASiblingSharingThePrefix(t *testing.T) {
	base := t.TempDir()
	dest := filepath.Join(base, "dest")

	err := extract(tarball(t, entry{name: "../destevil/x", body: "owned"}), dest, noCeiling)
	if err == nil {
		t.Fatal("extract accepted a sibling directory sharing the destination's prefix")
	}
	if _, statErr := os.Stat(filepath.Join(base, "destevil")); !os.IsNotExist(statErr) {
		t.Error("the entry wrote into a sibling of the destination")
	}
}

func TestExtractRefusesAnAbsolutePath(t *testing.T) {
	for _, path := range []string{"/etc/passwd", "/tmp/owned"} {
		dest := filepath.Join(t.TempDir(), "v0.4.0")

		if err := extract(tarball(t, entry{name: path, body: "owned"}), dest, noCeiling); err == nil {
			t.Errorf("extract accepted the absolute path %q", path)
		}
	}
}

// A symlink in the payload has no legitimate purpose, and it is the second half of every
// extraction escape: a link to `/` followed by a write through it.
func TestExtractRefusesASymlinkEntry(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "v0.4.0")

	err := extract(tarball(t,
		entry{name: "skills/evil", typeflag: tar.TypeSymlink, linkname: "/"},
	), dest, noCeiling)
	if err == nil {
		t.Fatal("extract accepted a symlink entry")
	}
	if !strings.Contains(err.Error(), "skills/evil") {
		t.Errorf("the refusal does not name the entry: %v", err)
	}
}

// A whitelist, not a list of things to reject: anything that is not a regular file or a
// directory is refused, including types nobody thought to enumerate.
func TestExtractRefusesANonRegularEntry(t *testing.T) {
	cases := map[string]byte{
		"hard link":    tar.TypeLink,
		"fifo":         tar.TypeFifo,
		"char device":  tar.TypeChar,
		"block device": tar.TypeBlock,
	}

	for name, flag := range cases {
		t.Run(name, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "v0.4.0")
			e := entry{name: "skills/odd", typeflag: flag}
			if flag == tar.TypeLink {
				e.linkname = "skills/other"
			}

			if err := extract(tarball(t, e), dest, noCeiling); err == nil {
				t.Errorf("extract accepted a %s", name)
			}
		})
	}
}

// Not skipped with a warning. An archive containing one such entry is not an archive to
// take the rest of, so a legitimate entry after the refused one must not land either.
func TestOneRefusedEntryAbortsTheWholeExtraction(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "v0.4.0")

	err := extract(tarball(t,
		entry{name: "skills/first.md", body: "fine"},
		entry{name: "../escaped", body: "owned"},
		entry{name: "skills/third.md", body: "also fine"},
	), dest, noCeiling)
	if err == nil {
		t.Fatal("extract accepted an archive with a refused entry in the middle")
	}

	if _, statErr := os.Stat(filepath.Join(dest, "skills", "third.md")); !os.IsNotExist(statErr) {
		t.Error("extraction continued past the refused entry")
	}
}

// A small archive that decompresses to fill the disk fails rather than succeeding slowly.
func TestExtractStopsAtTheSizeCeiling(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "v0.4.0")
	big := strings.Repeat("x", 4096)

	err := extract(tarball(t,
		entry{name: "skills/a.md", body: big},
		entry{name: "skills/b.md", body: big},
	), dest, 5000)
	if err == nil {
		t.Fatal("extract ignored the size ceiling")
	}
	if !strings.Contains(err.Error(), "too large") {
		t.Errorf("the error does not say what happened: %v", err)
	}
}

// Modes come from us, not from the archive. The payload is markdown; nothing in it needs to
// be executable, and an executable bit arriving from a download is a bit nobody asked for.
func TestExtractNormalisesModesAndDropsTheExecutableBit(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "v0.4.0")

	err := extract(tarball(t,
		entry{name: "skills/", typeflag: tar.TypeDir, mode: 0o777},
		entry{name: "skills/runme.md", body: "# not a script\n", mode: 0o755},
		entry{name: "skills/setuid.md", body: "# nor this\n", mode: 0o4755},
	), dest, noCeiling)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	for _, name := range []string{"runme.md", "setuid.md"} {
		info, err := os.Stat(filepath.Join(dest, "skills", name))
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != 0o644 {
			t.Errorf("%s has mode %v, want 0644", name, got)
		}
		if info.Mode()&fs.ModeSetuid != 0 {
			t.Errorf("%s kept the setuid bit", name)
		}
	}

	info, err := os.Stat(filepath.Join(dest, "skills"))
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o755 {
		t.Errorf("the directory has mode %v, want 0755", got)
	}
}

// Not a tarball at all, and not gzip. Both are what a truncated download or an HTML error
// page look like, and neither may read as an empty-but-successful extraction.
func TestExtractRefusesSomethingThatIsNotAGzippedTarball(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "v0.4.0")

	if err := extract(bytes.NewReader([]byte("<html>404</html>")), dest, noCeiling); err == nil {
		t.Error("extract accepted something that is not gzip")
	}

	var gzipped bytes.Buffer
	gz := gzip.NewWriter(&gzipped)
	gz.Write([]byte("not a tar"))
	gz.Close()

	if err := extract(bytes.NewReader(gzipped.Bytes()), dest, noCeiling); err == nil {
		t.Error("extract accepted gzip that is not a tarball")
	}
}
