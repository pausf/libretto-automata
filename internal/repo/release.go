package repo

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// CheckTTL is how long an answer about the latest release stays good.
//
// A day, because releases are not hourly and the alternative is a subprocess and a
// network round trip on every single panel launch. ponytail: a constant, not a setting —
// a knob nobody turns is a knob to document, test and support.
const CheckTTL = 24 * time.Hour

// checkFile is where the answer is kept: inside .git, so it is never committed, needs no
// .gitignore entry anybody has to remember, and disappears with the clone.
const checkFile = "libretto-update-check"

// LatestTag is the highest plain release tag on the remote.
//
// `git ls-remote --tags` rather than the GitHub API or the Go module proxy: no token, no
// rate limit, no JSON, and it works against whatever remote the user's git can reach —
// including a fork, and including the local path this is tested against.
//
// An empty string with no error means the remote has no release to offer. A remote with
// only prereleases, or no tags at all, is a state and not a failure.
func (s Shell) LatestTag(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", s.Root, "ls-remote", "--tags", "origin")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if text := strings.TrimSpace(string(out)); text != "" {
			return "", fmt.Errorf("%s", text)
		}
		return "", err
	}
	return highestTag(string(out)), nil
}

// highestTag reads ls-remote's output: one `<sha>\trefs/tags/<name>` per line.
//
// `^{}` suffixes mark the commit an annotated tag points at, and this repository's tags
// are annotated — so every release appears twice and the peeled line has to be dropped or
// the name is `v0.2.0^{}`, which parses as nothing.
func highestTag(out string) string {
	best, bestParsed := "", [3]int{}
	for _, line := range strings.Split(out, "\n") {
		_, ref, ok := strings.Cut(strings.TrimSpace(line), "\t")
		if !ok {
			continue
		}
		name := strings.TrimSuffix(strings.TrimPrefix(ref, "refs/tags/"), "^{}")

		parsed, ok := parseSemver(name)
		if !ok {
			continue
		}
		if best == "" || greater(parsed, bestParsed) {
			best, bestParsed = name, parsed
		}
	}
	return best
}

func greater(a, b [3]int) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] > b[i]
		}
	}
	return false
}

// CheckedLatest is LatestTag behind the cache. This is the only entry point the CLI uses.
func CheckedLatest(ctx context.Context, root string, ttl time.Duration) (string, error) {
	return checkedLatest(ctx, root, ttl, time.Now, Shell{Root: root}.LatestTag)
}

// checkedLatest takes its clock and its asker, because a test that depends on the wall
// clock behaves differently at midnight and one that depends on the network is not a test.
//
// Failure is cached as an empty answer, deliberately. Caching only successes means a
// machine with no network pays the timeout on every launch — the hang this exists to
// prevent, arriving once per invocation instead of once a day. The caller sees no error
// either way: nobody needs a panel that reports its own inability to check.
//
// No .git means nowhere to keep the answer, which is the bootstrap case. Ask and do not
// cache, rather than fail.
func checkedLatest(
	ctx context.Context,
	root string,
	ttl time.Duration,
	now func() time.Time,
	ask func(context.Context) (string, error),
) (string, error) {
	path := filepath.Join(root, ".git", checkFile)
	cacheable := false
	if info, err := os.Stat(filepath.Join(root, ".git")); err == nil && info.IsDir() {
		cacheable = true
	}

	if cacheable {
		if tag, at, ok := readCheck(path); ok && now().Sub(at) < ttl {
			return tag, nil
		}
	}

	tag, err := ask(ctx)
	if err != nil {
		tag = ""
	}
	if cacheable {
		// A cache that cannot be written is not worth an error. The answer is already in
		// hand; the only cost is asking again next time.
		_ = os.WriteFile(path,
			[]byte(strconv.FormatInt(now().Unix(), 10)+" "+tag+"\n"), 0o644)
	}
	return tag, nil
}

// readCheck parses `<unix> [tag]`. The tag is absent when the last check failed, which is
// a valid entry and the reason the timestamp comes first.
func readCheck(path string) (string, time.Time, bool) {
	body, err := os.ReadFile(path)
	if err != nil {
		return "", time.Time{}, false
	}
	stamp, tag, _ := strings.Cut(strings.TrimSpace(string(body)), " ")
	secs, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		return "", time.Time{}, false
	}
	return strings.TrimSpace(tag), time.Unix(secs, 0), true
}

// parseSemver reads a plain release tag: `v` then exactly three decimal fields.
//
// Only plain releases count. A prerelease, a `git describe` string and build metadata
// are all rejected rather than ranked — this repository tags plain semver and says so
// in AGENTS.md, so ordering the other forms would be code with no caller.
//
// ponytail: three fields, no prerelease ordering. The first v1.0.0-rc.1 that is meant
// to be offered to users needs real precedence rules here — spec §11 of semver, which
// is thirty lines and currently proves nothing. Until then a prerelease is invisible,
// which is the safe direction: it cannot claim to be newer than the release it precedes.
func parseSemver(tag string) ([3]int, bool) {
	var out [3]int
	if len(tag) < 2 || tag[0] != 'v' {
		return out, false
	}

	field, seen, digits := 0, 0, 0
	for _, r := range tag[1:] {
		switch {
		case r >= '0' && r <= '9':
			// Guard the multiply rather than trusting the input's length: a tag is a
			// remote-controlled string and an overflowing field would compare as
			// negative, which reads as "older" and silently hides an update.
			if digits > 9 {
				return out, false
			}
			field = field*10 + int(r-'0')
			digits++
		case r == '.':
			if digits == 0 || seen == 2 {
				return out, false
			}
			out[seen], seen, field, digits = field, seen+1, 0, 0
		default:
			return out, false
		}
	}
	if seen != 2 || digits == 0 {
		return out, false
	}
	out[2] = field
	return out, true
}

// IsRelease reports whether tag names a plain release.
//
// Exported for internal/dist, which has to decide whether a redirect's target and a
// directory's name are versions. The alternative was a second parser there, and two
// implementations of "what counts as a release" is two that can disagree about a prerelease.
//
// ponytail: dist imports repo for this and for IsNewer, and for nothing else — a package
// about downloads depending on a package about git, at compile time only. If that coupling
// ever costs anything, the four semver functions move to their own package and both depend
// on it. Two exported functions is not worth a package today.
func IsRelease(tag string) bool {
	_, ok := parseSemver(tag)
	return ok
}

// IsNewer reports whether latest is a release strictly ahead of running.
//
// Exported because cmd/libretto formats the notice and must not own a second answer to
// the same question. Two implementations of "is this newer" is two that can disagree,
// and the one that gets read is whichever the reader found first.
//
// Either side being unparseable is false, never true. Telling someone whose binary
// reports `dev` that they are out of date is a guess presented as a fact.
func IsNewer(latest, running string) bool {
	l, ok := parseSemver(latest)
	if !ok {
		return false
	}
	r, ok := parseSemver(running)
	if !ok {
		return false
	}
	for i := range l {
		if l[i] != r[i] {
			return l[i] > r[i]
		}
	}
	return false
}
