// Package dist gets a versioned payload from a published release onto the machine, and
// makes it the one the links point at.
//
// It is deliberately separate from internal/repo. That package answers questions about a
// git checkout; nothing here involves git, and folding a downloader into it would give one
// package two unrelated reasons to exist.
package dist

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// File and directory modes come from here, never from the archive.
//
// The payload is markdown. Nothing in it needs to be executable, so an executable bit
// arriving from a download is a bit nobody asked for — and a setuid bit is worse than that.
const (
	fileMode = 0o644
	dirMode  = 0o755
)

// extract unpacks a gzipped tarball into dest, refusing anything it was not built to accept.
//
// **This is a trust boundary.** The archive is remote input, and an extractor that trusts
// its entries writes wherever the archive tells it to. The rules, all of them enforced
// before any byte of an entry is written:
//
//   - regular files and directories only — a whitelist, so a type nobody enumerated is
//     refused rather than falling through to a default
//   - every entry's resolved path must be inside dest
//   - the running total of bytes written must stay under maxBytes
//   - modes are ours
//
// A refused entry aborts the whole extraction. It is not skipped with a warning: an archive
// containing one such entry is not an archive to take the rest of. Cleaning up what was
// already written is the caller's job — install extracts to a temporary directory and
// renames, so a failure leaves nothing that looks installed.
func extract(r io.Reader, dest string, maxBytes int64) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("not a gzipped tarball: %w", err)
	}
	defer gz.Close()

	// Cleaned once, here, because every containment check compares against it. An
	// uncleaned dest with a trailing separator or a `.` in it makes every comparison below
	// subtly wrong.
	root := filepath.Clean(dest)
	if err := os.MkdirAll(root, dirMode); err != nil {
		return err
	}

	var written int64
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading the tarball: %w", err)
		}

		switch h.Typeflag {
		case tar.TypeReg, tar.TypeDir:
		default:
			// Symlinks are refused here along with everything else, and they are the case
			// worth naming: a link to `/` followed by a write through it is the second half
			// of every extraction escape. The payload has no legitimate use for one.
			return fmt.Errorf("refusing %q: only regular files and directories are accepted, not type %q",
				h.Name, string(h.Typeflag))
		}

		path, err := resolve(root, h.Name)
		if err != nil {
			return err
		}

		if h.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(path, dirMode); err != nil {
				return err
			}
			continue
		}

		if written+h.Size > maxBytes {
			return fmt.Errorf("refusing %q: the archive is too large, over the %d byte ceiling",
				h.Name, maxBytes)
		}

		// The parent may have no entry of its own. `tar czf` writes directory entries, but
		// nothing in the format guarantees they arrive before their contents.
		if err := os.MkdirAll(filepath.Dir(path), dirMode); err != nil {
			return err
		}

		n, err := writeFile(path, tr, h.Size)
		written += n
		if err != nil {
			return err
		}
	}
}

// resolve turns an entry name into a path inside root, or refuses it.
//
// The check is on the *cleaned* path with a separator appended, not on the raw string. A
// plain `strings.HasPrefix(path, root)` accepts `/tmp/destevil` as being inside `/tmp/dest`
// — the sibling that shares the prefix, and the reason this is its own function with its
// own test.
func resolve(root, name string) (string, error) {
	if filepath.IsAbs(name) || strings.HasPrefix(name, "/") {
		return "", fmt.Errorf("refusing %q: an absolute path", name)
	}

	path := filepath.Clean(filepath.Join(root, name))
	if path != root && !strings.HasPrefix(path, root+string(os.PathSeparator)) {
		return "", fmt.Errorf("refusing %q: it resolves outside the destination", name)
	}
	return path, nil
}

// writeFile copies exactly size bytes, and returns how many it wrote even on failure so the
// caller's running total stays honest.
//
// io.CopyN rather than io.Copy: the header's size is the contract, and a tar reader that
// yielded more than the header declared would otherwise write past the ceiling the caller
// just checked.
func writeFile(path string, r io.Reader, size int64) (int64, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fileMode)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	n, err := io.CopyN(f, r, size)
	if err != nil && err != io.EOF {
		return n, fmt.Errorf("writing %s: %w", path, err)
	}

	// O_CREATE honours the umask, so a file that already existed keeps its old mode and a
	// new one may be stricter than intended. Setting it explicitly is what makes the mode
	// ours rather than the environment's.
	if err := os.Chmod(path, fileMode); err != nil {
		return n, err
	}
	return n, nil
}
