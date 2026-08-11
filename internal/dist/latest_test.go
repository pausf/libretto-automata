package dist

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// releaseServer serves the /releases/latest redirect and fails the test if the tag page it
// points at is ever requested. That second half is the whole point: the default http.Client
// follows redirects, so without CheckRedirect a "reads the tag" test passes by reading it
// out of a followed response — and nobody finds out until GitHub changes the page.
func releaseServer(t *testing.T, status int, location string) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/releases/latest"):
			if location != "" {
				w.Header().Set("Location", location)
			}
			w.WriteHeader(status)
		default:
			t.Errorf("the redirect was followed: something requested %s", r.URL.Path)
			w.WriteHeader(http.StatusOK)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestLatestReadsTheTagFromTheRedirectLocation(t *testing.T) {
	srv := releaseServer(t, http.StatusFound, "/pausf/libretto-automata/releases/tag/v0.4.0")

	got, err := Latest(context.Background(), srv.Client(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if got != "v0.4.0" {
		t.Errorf("Latest = %q, want v0.4.0", got)
	}
}

// An absolute Location is what GitHub actually sends. Both forms have to work, or this
// breaks the first time the header changes shape.
func TestLatestAcceptsAnAbsoluteLocation(t *testing.T) {
	srv := releaseServer(t, http.StatusFound,
		"https://github.com/pausf/libretto-automata/releases/tag/v1.2.3")

	got, err := Latest(context.Background(), srv.Client(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if got != "v1.2.3" {
		t.Errorf("Latest = %q, want v1.2.3", got)
	}
}

// The redirect is read, never followed. releaseServer fails the test if the tag page is
// fetched, so this asserts on the request that must not happen.
func TestLatestDoesNotFollowTheRedirect(t *testing.T) {
	srv := releaseServer(t, http.StatusFound, "/pausf/libretto-automata/releases/tag/v0.4.0")

	if _, err := Latest(context.Background(), srv.Client(), srv.URL); err != nil {
		t.Fatal(err)
	}
	// The assertion lives in the handler: any request for a path other than
	// /releases/latest fails the test.
}

// repo's semver rules decide what a release is. A Location naming a prerelease, a branch or
// nothing recognisable is "could not tell" — not a tag to offer somebody.
func TestLatestRejectsANonSemverLocation(t *testing.T) {
	for _, loc := range []string{
		"/pausf/libretto-automata/releases/tag/v1.0.0-rc.1",
		"/pausf/libretto-automata/releases/tag/nightly",
		"/pausf/libretto-automata/releases/tag/",
		"/pausf/libretto-automata/releases",
		"/",
	} {
		t.Run(loc, func(t *testing.T) {
			srv := releaseServer(t, http.StatusFound, loc)

			if got, err := Latest(context.Background(), srv.Client(), srv.URL); err == nil {
				t.Errorf("Latest accepted %q and returned %q", loc, got)
			}
		})
	}
}

// A repository with no releases answers 404; a proxy or an outage answers 200 with a page.
// Neither is a version, and neither may be guessed at.
func TestLatestOnAnUnexpectedStatusIsCouldNotTell(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNotFound, http.StatusInternalServerError} {
		srv := releaseServer(t, status, "")

		if _, err := Latest(context.Background(), srv.Client(), srv.URL); err == nil {
			t.Errorf("Latest accepted status %d", status)
		}
	}
}

// A 302 with no Location at all — malformed, and it must not read as an empty success.
func TestLatestWithNoLocationHeader(t *testing.T) {
	srv := releaseServer(t, http.StatusFound, "")

	if _, err := Latest(context.Background(), srv.Client(), srv.URL); err == nil {
		t.Error("Latest accepted a redirect with no Location")
	}
}

// A network that accepts a connection and never answers must not hold the panel's first
// paint. This is the check that keeps `libretto` from looking hung on bad DNS.
func TestLatestHonoursTheDeadline(t *testing.T) {
	srv := releaseServer(t, http.StatusFound, "/pausf/libretto-automata/releases/tag/v0.4.0")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := Latest(ctx, srv.Client(), srv.URL); err == nil {
		t.Error("Latest ignored a cancelled context")
	}
}

// The default client follows redirects, and Latest must work when handed one — it cannot
// depend on the caller having configured CheckRedirect correctly.
func TestLatestWorksWithADefaultClient(t *testing.T) {
	srv := releaseServer(t, http.StatusFound, "/pausf/libretto-automata/releases/tag/v0.4.0")

	got, err := Latest(context.Background(), &http.Client{}, srv.URL)
	if err != nil {
		t.Fatalf("Latest with a default client: %v", err)
	}
	if got != "v0.4.0" {
		t.Errorf("Latest = %q, want v0.4.0", got)
	}
}
