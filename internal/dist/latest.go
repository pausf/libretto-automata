package dist

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pausf/libretto-automata/internal/repo"
)

// tagSegment is what a release page's path ends with, and what Location is read for.
const tagSegment = "/releases/tag/"

// Latest is the newest published release tag, read from the /releases/latest redirect.
//
// Not api.github.com. The unauthenticated API allows sixty requests an hour — gentle-ai
// lives with that — while a redirect costs no quota, needs no token and parses no JSON. It
// is also one less response shape to break when the provider changes something.
//
// host is the forge's base URL, so a test can point it at httptest and a fork can be reached
// without a second code path.
//
// The redirect is read, never followed: following it fetches an HTML page nobody needs. The
// client's own CheckRedirect is overridden rather than trusted — Latest cannot depend on
// every caller having configured one, and the default client follows redirects, which would
// make it read the tag out of a followed response instead.
//
// Anything other than a 302 carrying a plain release tag is an error, never a guess. A
// repository with no releases answers 404; an outage or a captive portal answers 200 with a
// page. Neither is a version.
func Latest(ctx context.Context, client *http.Client, host string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimSuffix(host, "/")+"/releases/latest", nil)
	if err != nil {
		return "", err
	}

	// A shallow copy, so overriding CheckRedirect cannot mutate a client the caller shares
	// with something else. Everything the caller configured — timeouts, transport, proxy —
	// is carried over.
	noFollow := *client
	noFollow.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := noFollow.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
		return "", fmt.Errorf("asking for the latest release: unexpected status %s", resp.Status)
	}

	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", fmt.Errorf("asking for the latest release: the redirect carried no Location")
	}
	return tagFromLocation(loc)
}

// tagFromLocation pulls the release tag out of a redirect target.
//
// Both an absolute URL and a bare path have to work: GitHub sends the absolute form today,
// and depending on that is depending on a detail nobody promised.
//
// repo.IsRelease decides whether what came back is a release. A prerelease, a branch name or
// anything else is refused here for the same reason it is refused there — invisible is the
// safe direction, and this package does not get a second opinion about semver.
func tagFromLocation(loc string) (string, error) {
	p := loc
	if u, err := url.Parse(loc); err == nil && u.Path != "" {
		p = u.Path
	}

	i := strings.LastIndex(p, tagSegment)
	if i < 0 {
		return "", fmt.Errorf("the redirect does not point at a release: %s", loc)
	}

	tag := path.Base(p[i+len(tagSegment):])
	if !repo.IsRelease(tag) {
		return "", fmt.Errorf("the latest release is not a plain version: %q", tag)
	}
	return tag, nil
}
