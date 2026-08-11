// Package dist finds the payload an installed copy links from, and moves it to a newer one.
//
// **It downloads nothing and verifies nothing, and that is the design.** The payload —
// `skills/`, `agents/`, `commands/` — ships inside the Go module, so
// `go install <module>/cmd/libretto@latest` already fetches it, already checks it against the
// checksum database, and already puts it somewhere with a version in the path. This package
// says where that is and how to ask for a newer one.
//
// An earlier version of this package downloaded a release tarball, verified a sha256 and
// extracted it through a guarded extractor with five refusals. All of that was reimplementing,
// less well, machinery the go command already has: GOSUMDB is a stronger guarantee than a
// checksum published beside the file it describes, because it is not the publisher who vouches
// for it.
//
// The cost, stated: `go clean -modcache` removes the payload and every link into it breaks at
// once. `libretto install` re-downloads and repairs it, and `doctor` reports it.
package dist

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pausf/libretto-automata/internal/repo"
)

// DefaultProxy is where "what is the newest version" is asked.
//
// The module proxy, and not the GitHub API or the releases redirect: it is the same source
// `go install @latest` resolves against, so the two cannot disagree about what latest means.
// It also answers from **tags**, which is why `go install` works against a repository that has
// published no Releases at all.
const DefaultProxy = "https://proxy.golang.org"

// Dir is the module cache entry for one version of the module — where its payload is.
func Dir(cache, module, version string) string {
	return filepath.Join(cache, filepath.FromSlash(module)+"@"+version)
}

// ModCache is where the go command keeps downloaded modules.
func ModCache() string {
	home, _ := os.UserHomeDir()
	return modCache(os.Getenv("GOMODCACHE"), os.Getenv("GOPATH"), home)
}

// modCache follows the go command's own order, from its inputs rather than by shelling out to
// `go env`. Every command that touches the payload needs this, so a subprocess here is a
// subprocess on the hot path — for a path that three variables already determine.
func modCache(gomodcache, gopath, home string) string {
	switch {
	case gomodcache != "":
		return gomodcache
	case gopath != "":
		return filepath.Join(gopath, "pkg", "mod")
	case home != "":
		return filepath.Join(home, "go", "pkg", "mod")
	default:
		return ""
	}
}

// Latest is the newest published release of the module, from the proxy.
//
// Only a plain release counts, and `repo.IsRelease` is what decides — one notion of what a
// version is across the project. A module with no tags answers with a pseudo-version
// (`v0.0.0-<date>-<sha>`), which is not something to offer anybody as an update.
func Latest(ctx context.Context, client *http.Client, proxy, module string) (string, error) {
	url := strings.TrimSuffix(proxy, "/") + "/" + escapeModule(module) + "/@latest"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("asking the module proxy for the newest version: %s", resp.Status)
	}

	// Bounded: the answer is two fields. Anything longer is not this endpoint's response.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return "", err
	}

	var answer struct{ Version string }
	if err := json.Unmarshal(body, &answer); err != nil {
		return "", fmt.Errorf("the module proxy's answer could not be read: %w", err)
	}
	if !repo.IsRelease(answer.Version) {
		if answer.Version == "" {
			return "", fmt.Errorf("the module proxy named no version")
		}
		return "", fmt.Errorf("the newest version is not a plain release: %q", answer.Version)
	}
	return answer.Version, nil
}

// escapeModule lower-cases the module path the way the proxy protocol requires: an uppercase
// rune becomes `!` followed by its lowercase form.
//
// This project's path has no uppercase in it, so nothing here exercises this in production. It
// is written anyway because the day somebody forks to a capitalised account, the failure would
// be a 404 with no explanation attached to it.
func escapeModule(module string) string {
	var b strings.Builder
	for _, r := range module {
		if r >= 'A' && r <= 'Z' {
			b.WriteByte('!')
			b.WriteRune(r + ('a' - 'A'))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// runner is how the go command is invoked, taken as a parameter so no test installs anything.
type runner func(ctx context.Context, name string, args ...string) error

// Install runs `go install <module>/cmd/libretto@<version>`, which brings down the binary and
// the payload together.
func Install(ctx context.Context, module, version string) error {
	return install(ctx, runCommand, module, version)
}

func install(ctx context.Context, run runner, module, version string) error {
	if err := run(ctx, "go", "install", module+"/cmd/libretto@"+version); err != nil {
		return fmt.Errorf("installing %s: %w", version, err)
	}
	return nil
}

// runCommand keeps the command's own output when it fails. `go install` explains itself well —
// a proxy that refused, a checksum mismatch, a toolchain too old — and replacing that with
// "install failed" would throw away the only useful part.
func runCommand(ctx context.Context, name string, args ...string) error {
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil {
		if text := strings.TrimSpace(string(out)); text != "" {
			return fmt.Errorf("%s", text)
		}
		return err
	}
	return nil
}
