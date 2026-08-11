package dist

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

// The payload ships inside the module, so `go install` already downloads it. Dir is where it
// lands — nothing is fetched, extracted or verified by this package, because the go command
// did all three and checked the module against the checksum database while it was at it.
func TestDirIsTheModuleCacheEntryForAVersion(t *testing.T) {
	cache := "/home/x/go/pkg/mod"

	got := Dir(cache, "github.com/pausf/libretto-automata", "v0.5.0")
	want := filepath.Join(cache, "github.com", "pausf", "libretto-automata@v0.5.0")
	if got != want {
		t.Errorf("Dir = %q, want %q", got, want)
	}
}

// GOMODCACHE wins, then GOPATH/pkg/mod, then ~/go/pkg/mod — the go command's own order.
// Resolved without shelling out to `go env`: this runs on every invocation that needs the
// payload, and a subprocess for a path that three variables already determine is a subprocess
// on the hot path.
func TestModCacheFollowsTheGoCommandsOrder(t *testing.T) {
	cases := []struct {
		name                     string
		gomodcache, gopath, home string
		want                     string
	}{
		{"GOMODCACHE wins", "/explicit", "/gopath", "/home", "/explicit"},
		{"then GOPATH", "", "/gopath", "/home", filepath.Join("/gopath", "pkg", "mod")},
		{"then home", "", "", "/home", filepath.Join("/home", "go", "pkg", "mod")},
		{"nothing known", "", "", "", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := modCache(c.gomodcache, c.gopath, c.home); got != c.want {
				t.Errorf("modCache = %q, want %q", got, c.want)
			}
		})
	}
}

// The module proxy is asked, not the GitHub API and not the releases redirect. It is the same
// source `go install @latest` resolves against, so the two cannot disagree about what latest
// means — and it answers off **tags**, which is why `go install` works on a repository with no
// published Releases at all.
func TestLatestReadsTheVersionFromTheModuleProxy(t *testing.T) {
	srv := proxyServer(t, http.StatusOK, `{"Version":"v0.5.0","Time":"2026-08-11T10:44:31Z"}`)

	got, err := Latest(context.Background(), srv.Client(), srv.URL, "github.com/pausf/libretto-automata")
	if err != nil {
		t.Fatal(err)
	}
	if got != "v0.5.0" {
		t.Errorf("Latest = %q, want v0.5.0", got)
	}
}

// A module with no tag at all answers with a pseudo-version — `v0.0.0-<date>-<sha>`. That is
// not a release and must not be offered as one: `repo.IsRelease` decides, so there is still one
// notion of what a version is across the whole project.
func TestLatestRefusesAPseudoVersion(t *testing.T) {
	for _, v := range []string{
		"v0.0.0-20260811104431-34c43e91d94f",
		"v1.0.0-rc.1",
		"v0.5.0+incompatible",
	} {
		t.Run(v, func(t *testing.T) {
			srv := proxyServer(t, http.StatusOK, fmt.Sprintf(`{"Version":%q}`, v))

			if got, err := Latest(context.Background(), srv.Client(), srv.URL, "m"); err == nil {
				t.Errorf("Latest accepted %q and returned %q", v, got)
			}
		})
	}
}

// A 404 from the proxy means the module has never been published, which is a state to report
// and not a version to invent.
func TestLatestOnANonOKStatus(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusGone, http.StatusInternalServerError} {
		srv := proxyServer(t, status, "")

		if _, err := Latest(context.Background(), srv.Client(), srv.URL, "m"); err == nil {
			t.Errorf("Latest accepted status %d", status)
		}
	}
}

// Malformed JSON, and an empty body. Both are what a captive portal or a truncated response
// look like, and neither may read as an answer.
func TestLatestRefusesAnUnreadableAnswer(t *testing.T) {
	for _, body := range []string{"", "<html>hello</html>", `{"Version":}`, `{}`} {
		srv := proxyServer(t, http.StatusOK, body)

		if _, err := Latest(context.Background(), srv.Client(), srv.URL, "m"); err == nil {
			t.Errorf("Latest accepted %q", body)
		}
	}
}

// The panel's first paint must not wait on this.
func TestLatestHonoursTheDeadline(t *testing.T) {
	srv := proxyServer(t, http.StatusOK, `{"Version":"v0.5.0"}`)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := Latest(ctx, srv.Client(), srv.URL, "m"); err == nil {
		t.Error("Latest ignored a cancelled context")
	}
}

// A module path with an uppercase letter has to be escaped for the proxy — `!` before the
// lowered rune. This project's path is all lowercase, so the escaping is a correctness
// property nothing here exercises in production; it is written and tested anyway, because the
// day somebody forks to a capitalised account the failure would be a 404 nobody could explain.
func TestLatestEscapesUppercaseInTheModulePath(t *testing.T) {
	var asked string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		asked = r.URL.Path
		fmt.Fprint(w, `{"Version":"v0.5.0"}`)
	}))
	t.Cleanup(srv.Close)

	if _, err := Latest(context.Background(), srv.Client(), srv.URL, "github.com/Gentleman/Tool"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(asked, "!gentleman/!tool") {
		t.Errorf("the proxy was asked for %q, want the escaped path", asked)
	}
}

// Install runs the go command, and takes the runner so no test installs anything.
func TestInstallRunsGoInstallAtTheVersion(t *testing.T) {
	var got []string
	run := func(_ context.Context, name string, args ...string) error {
		got = append([]string{name}, args...)
		return nil
	}

	if err := install(context.Background(), run, "github.com/pausf/libretto-automata", "v0.5.0"); err != nil {
		t.Fatal(err)
	}

	want := "go install github.com/pausf/libretto-automata/cmd/libretto@v0.5.0"
	if strings.Join(got, " ") != want {
		t.Errorf("ran %q, want %q", strings.Join(got, " "), want)
	}
}

// A refused install is reported, not swallowed. No Go on the machine is the common case and it
// has to say so rather than leaving the caller to guess why nothing changed.
func TestInstallReportsAFailure(t *testing.T) {
	run := func(context.Context, string, ...string) error {
		return errors.New("executable file not found in $PATH")
	}

	err := install(context.Background(), run, "m", "v0.5.0")
	if err == nil {
		t.Fatal("a failed go install was reported as success")
	}
	if !strings.Contains(err.Error(), "v0.5.0") {
		t.Errorf("the error does not name the version it was installing: %v", err)
	}
}

func proxyServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/@latest") {
			t.Errorf("unexpected request for %s", r.URL.Path)
		}
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv
}
