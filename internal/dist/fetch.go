package dist

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// maxPayloadBytes is the extraction ceiling.
//
// ponytail: 64 MiB, a constant rather than a setting. The payload is markdown and the current
// one is well under a megabyte, so this is two orders of magnitude of room. If a payload ever
// legitimately approaches it, the number moves — a knob nobody turns is a knob to document,
// test and support.
const maxPayloadBytes = 64 << 20

// assetNames are the two files `make release` attaches to a release.
//
// This is one half of a contract whose other half is written in the Makefile, in a different
// language. TestReleaseAssetNamesMatchWhatDistributionFetches holds them together, because a
// typo in either is a release that installs on nobody's machine and nothing else would catch
// it.
func assetNames(tag string) (tarball, checksum string) {
	tarball = "payload-" + tag + ".tar.gz"
	return tarball, tarball + ".sha256"
}

// extractor is what fetch hands the verified bytes to. A parameter so the ordering — verify,
// then extract — can be proven by a spy rather than inferred from an error message.
type extractor func(r io.Reader, dest string) error

// extractInto is the real extractor, with the ceiling applied.
func extractInto(r io.Reader, dest string) error {
	return extract(r, dest, maxPayloadBytes)
}

// fetch downloads a release's payload, verifies it, and extracts it into its version
// directory.
//
// **The checksum is fetched and compared before a single byte is extracted.** Not afterwards:
// by then the payload is on disk and one symlink away from ~/.claude, and there is no later
// remedy for that. The order here is the security property, and it is the one the tests pin.
//
// A version already on disk is not re-downloaded. That makes `upgrade` twice one download,
// and — more importantly — it means a payload that is currently linked is never overwritten
// underneath the links pointing into it.
//
// Nothing is left behind on failure: the download goes to a temporary file that is removed on
// every path, and extraction goes to a temporary directory renamed into place only once it
// has completed.
func fetch(ctx context.Context, client *http.Client, host, base, tag string, unpack extractor) error {
	dest := VersionDir(base, tag)
	if info, err := os.Stat(dest); err == nil && info.IsDir() {
		return nil
	}

	tarball, checksum := assetNames(tag)

	// The checksum first. It is a few dozen bytes, so fetching it before the archive costs
	// nothing and means a release published without one fails before the download.
	want, err := fetchDigest(ctx, client, host, tag, checksum)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(base, dirMode); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(base, ".download-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	got, err := download(ctx, client, assetURL(host, tag, tarball), tmp)
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("refusing %s: checksum mismatch\n  published %s\n  got       %s",
			tarball, want, got)
	}

	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Extract beside the destination, then rename. A half-extracted v0.4.0/ that a later run
	// treats as installed is a payload with holes in it, and `fetch` above would then skip
	// the download that would have fixed it.
	staging, err := os.MkdirTemp(base, ".unpack-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)

	if err := unpack(tmp, staging); err != nil {
		return err
	}
	return os.Rename(staging, dest)
}

// download copies the response body into w and returns its sha256, so the bytes are hashed
// as they arrive rather than read a second time.
func download(ctx context.Context, client *http.Client, url string, w io.Writer) (string, error) {
	resp, err := get(ctx, client, url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	sum := sha256.New()
	if _, err := io.Copy(io.MultiWriter(w, sum), io.LimitReader(resp.Body, maxPayloadBytes+1)); err != nil {
		return "", fmt.Errorf("downloading %s: %w", url, err)
	}
	return hex.EncodeToString(sum.Sum(nil)), nil
}

// fetchDigest reads the published digest.
//
// `shasum -a 256` writes `<hex>  <filename>`, two fields. Reading the whole line as the digest
// would compare a hash against a hash-plus-a-name and never match — which would look like
// every release being corrupt.
func fetchDigest(ctx context.Context, client *http.Client, host, tag, name string) (string, error) {
	resp, err := get(ctx, client, assetURL(host, tag, name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// A digest line is under a hundred bytes; anything longer is not one.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		return "", err
	}

	digest, _, _ := strings.Cut(strings.TrimSpace(string(body)), " ")
	if len(digest) != sha256.Size*2 {
		return "", fmt.Errorf("%s does not contain a sha256 digest", name)
	}
	if _, err := hex.DecodeString(digest); err != nil {
		return "", fmt.Errorf("%s does not contain a sha256 digest: %w", name, err)
	}
	return digest, nil
}

func assetURL(host, tag, name string) string {
	return strings.TrimSuffix(host, "/") + "/releases/download/" + tag + "/" + name
}

// get is a GET whose non-200 is an error. A 404 on an asset means a release published without
// it — which is what a hand-run `make release` that half-failed looks like — and a body of
// HTML is not a payload.
func get(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("fetching %s: %s", url, resp.Status)
	}
	return resp, nil
}
