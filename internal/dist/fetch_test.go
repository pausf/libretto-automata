package dist

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssetNamesFollowTheTag(t *testing.T) {
	tarball, checksum := assetNames("v0.4.0")

	if tarball != "payload-v0.4.0.tar.gz" {
		t.Errorf("tarball = %q", tarball)
	}
	if checksum != "payload-v0.4.0.tar.gz.sha256" {
		t.Errorf("checksum = %q", checksum)
	}
}

// The checksum is fetched and compared before a single byte is extracted. Not after: by then
// the payload is already on disk and one symlink away from ~/.claude, and there is no later
// remedy for that.
//
// The spy is what makes the ordering provable. Asserting "it failed" would pass for an
// implementation that extracted first and cleaned up afterwards.
func TestFetchVerifiesTheChecksumBeforeExtracting(t *testing.T) {
	body := goodTarball(t)
	srv := assetServer(t, "v0.4.0", body, "0000000000000000000000000000000000000000000000000000000000000000")

	base := t.TempDir()
	extracted := false
	err := fetch(context.Background(), srv.Client(), srv.URL, base, "v0.4.0",
		func(io.Reader, string) error {
			extracted = true
			return nil
		})

	if err == nil {
		t.Fatal("fetch accepted a mismatched checksum")
	}
	if extracted {
		t.Error("extraction was attempted before the checksum was verified")
	}
	if !strings.Contains(err.Error(), "checksum") {
		t.Errorf("the error does not say what failed: %v", err)
	}
}

func TestFetchExtractsWhenTheChecksumMatches(t *testing.T) {
	body := goodTarball(t)
	srv := assetServer(t, "v0.4.0", body, digestOf(body))

	base := t.TempDir()
	var got []byte
	err := fetch(context.Background(), srv.Client(), srv.URL, base, "v0.4.0",
		func(r io.Reader, dest string) error {
			var rerr error
			got, rerr = io.ReadAll(r)
			return rerr
		})
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if len(got) != len(body) {
		t.Errorf("the extractor got %d bytes, want %d", len(got), len(body))
	}
}

// A mismatch keeps nothing: no version directory, and no temporary file left in the base for
// a later run to trip over.
func TestFetchRefusesAMismatchedChecksumAndKeepsNothing(t *testing.T) {
	body := goodTarball(t)
	srv := assetServer(t, "v0.4.0", body, digestOf([]byte("something else")))

	base := t.TempDir()
	if err := fetch(context.Background(), srv.Client(), srv.URL, base, "v0.4.0", extractInto); err == nil {
		t.Fatal("fetch accepted a mismatched checksum")
	}

	if _, err := os.Stat(VersionDir(base, "v0.4.0")); !os.IsNotExist(err) {
		t.Error("a refused fetch left a version directory behind")
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("a refused fetch left %d entries in the base: %v", len(entries), entries)
	}
}

// `upgrade` twice is one download. Re-fetching a version already on disk wastes the network
// and, worse, would overwrite a payload that is currently linked.
func TestFetchDoesNotRedownloadAVersionAlreadyPresent(t *testing.T) {
	base := t.TempDir()
	mkVersion(t, base, "v0.4.0")

	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	if err := fetch(context.Background(), srv.Client(), srv.URL, base, "v0.4.0", extractInto); err != nil {
		t.Fatalf("fetch on a version already present: %v", err)
	}
	if requests != 0 {
		t.Errorf("fetch made %d requests for a version already on disk", requests)
	}
	// And the existing payload is untouched.
	if _, err := os.Stat(filepath.Join(VersionDir(base, "v0.4.0"), "skills", "marker.md")); err != nil {
		t.Errorf("the present version was disturbed: %v", err)
	}
}

// A 404 on either asset is a release published without them — which is what a hand-run
// `make release` that half-failed looks like.
func TestFetchReportsAMissingAsset(t *testing.T) {
	base := t.TempDir()

	for _, missing := range []string{"tarball", "checksum"} {
		t.Run(missing, func(t *testing.T) {
			body := goodTarball(t)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.HasSuffix(r.URL.Path, ".sha256"):
					if missing == "checksum" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					fmt.Fprintf(w, "%s  payload-v0.4.0.tar.gz\n", digestOf(body))
				default:
					if missing == "tarball" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					w.Write(body)
				}
			}))
			t.Cleanup(srv.Close)

			if err := fetch(context.Background(), srv.Client(), srv.URL, base, "v0.4.0", extractInto); err == nil {
				t.Errorf("fetch accepted a release with no %s", missing)
			}
		})
	}
}

// shasum writes "<hex>  <filename>". The filename half is not the digest, and reading the
// whole line as one would compare a hash against a hash-plus-a-name and never match.
func TestFetchReadsShasumsTwoFieldFormat(t *testing.T) {
	body := goodTarball(t)
	srv := assetServer(t, "v0.4.0", body, digestOf(body))

	base := t.TempDir()
	if err := fetch(context.Background(), srv.Client(), srv.URL, base, "v0.4.0",
		func(io.Reader, string) error { return nil }); err != nil {
		t.Fatalf("fetch could not read shasum's output format: %v", err)
	}
}

func TestFetchHonoursTheDeadline(t *testing.T) {
	body := goodTarball(t)
	srv := assetServer(t, "v0.4.0", body, digestOf(body))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := fetch(ctx, srv.Client(), srv.URL, t.TempDir(), "v0.4.0", extractInto); err == nil {
		t.Error("fetch ignored a cancelled context")
	}
}

// assetServer serves both release assets for one tag, with whatever digest it is told —
// so a test can hand it a wrong one without building a corrupt archive.
func assetServer(t *testing.T, tag string, body []byte, digest string) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tarball, checksum := assetNames(tag)
		switch {
		case strings.HasSuffix(r.URL.Path, checksum):
			fmt.Fprintf(w, "%s  %s\n", digest, tarball)
		case strings.HasSuffix(r.URL.Path, tarball):
			w.Write(body)
		default:
			t.Errorf("unexpected request for %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func goodTarball(t *testing.T) []byte {
	t.Helper()
	r := tarball(t, entry{name: "skills/marker.md", body: "# marker\n"})
	body, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func digestOf(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
