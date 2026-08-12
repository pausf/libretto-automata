package repo

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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
//
// **A tag that retracts itself is not a release.** `v1.0.2` in this repository has no
// Release and exists only to carry a `retract` block covering the two bad versions before
// it and itself — and it is perfectly valid semver, so it sorted highest and got offered
// as an update. The module proxy honours `retract` and answers `v0.7.0`; a git ref carries
// no retraction at all, so the checkout path had to learn to open the tag and look.
func (s Shell) LatestTag(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", s.Root, "ls-remote", "--tags", "origin")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if text := strings.TrimSpace(string(out)); text != "" {
			return "", fmt.Errorf("%s", text)
		}
		return "", err
	}

	candidates := descendingTags(string(out))
	if len(candidates) > maxRetractProbes {
		candidates = candidates[:maxRetractProbes]
	}
	for _, tag := range candidates {
		if !s.retractsItself(ctx, tag) {
			return tag, nil
		}
	}
	return "", nil
}

// maxRetractProbes bounds the walk down a remote's tags.
//
// ponytail: a constant, and low. Each step past the first is a `git show` on a tag that
// turned out to be retracted, and the real answer is one — this repository has exactly one
// tombstone. Ten is headroom for a project that retracted a bad run of releases; a remote
// with more than that is not one this notice can usefully summarise anyway, and an
// unbounded loop over a remote-controlled list is a subprocess per line.
const maxRetractProbes = 10

// descendingTags reads ls-remote's output — one `<sha>\trefs/tags/<name>` per line — and
// returns every plain release tag, newest first.
//
// `^{}` suffixes mark the commit an annotated tag points at, and this repository's tags
// are annotated — so every release appears twice and the peeled line has to be dropped or
// the name is `v0.2.0^{}`, which parses as nothing.
//
// It returns the whole ordered list rather than the single highest, because the highest is
// no longer necessarily the answer: it has to be opened first, and the runner-up is what
// the answer falls back to.
func descendingTags(out string) []string {
	type tag struct {
		name   string
		parsed [3]int
	}

	var tags []tag
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
		tags = append(tags, tag{name, parsed})
	}

	sort.Slice(tags, func(i, j int) bool { return greater(tags[i].parsed, tags[j].parsed) })

	names := make([]string, 0, len(tags))
	for _, t := range tags {
		// An annotated tag appears twice in ls-remote's output and both lines survive the
		// peel-stripping as the same name. Dropping the repeat here keeps the probe budget
		// counting distinct tags rather than lines.
		if len(names) > 0 && names[len(names)-1] == t.name {
			continue
		}
		names = append(names, t.name)
	}
	return names
}

// retractsItself reports whether tag's own go.mod retracts tag.
//
// The tag's own, not the working tree's: a retraction is published *with* the version it
// withdraws, which is the only arrangement that lets the highest version retract itself.
//
// **An unreadable go.mod is not evidence of a retraction.** `ls-remote` can name a tag the
// local object store has never seen — one pushed since the last fetch — and `git show`
// fails on it. Both answers are wrong in some case and this one is wrong in the rarer,
// cheaper one: a tag that new is a release far more often than a tombstone, and hiding it
// would mean a genuine release going unannounced to everyone whose tags are a day old.
//
// ponytail: no fetch. Reaching the network to settle a speculative background check is the
// hang the cache exists to prevent, arriving from a different direction — and it would
// write to the user's object store for a question nobody asked.
func (s Shell) retractsItself(ctx context.Context, tag string) bool {
	out, err := exec.CommandContext(ctx, "git", "-C", s.Root, "show", tag+":go.mod").Output()
	if err != nil {
		return false
	}
	return retracted(string(out), tag)
}

// retracted reports whether a go.mod's retract directives cover version.
//
// Hand-parsed rather than through `golang.org/x/mod/modfile`, which would be a sixth direct
// dependency for one predicate — AGENTS.md's ladder puts stdlib first and asks before
// adding. All three forms the go.mod grammar allows are read: a single directive, a
// parenthesised block, and a `[low, high]` range in either position.
func retracted(gomod, version string) bool {
	want, ok := parseSemver(version)
	if !ok {
		return false
	}

	inBlock := false
	for _, raw := range strings.Split(gomod, "\n") {
		// Comments carry the reason a version was withdrawn, and `// retract v1.0.2 one
		// day` is a note rather than a directive.
		line := strings.TrimSpace(raw)
		if i := strings.Index(line, "//"); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		if line == "" {
			continue
		}

		switch {
		case inBlock:
			if line == ")" {
				inBlock = false
				continue
			}
		case line == "retract (":
			inBlock = true
			continue
		case strings.HasPrefix(line, "retract "):
			line = strings.TrimSpace(strings.TrimPrefix(line, "retract "))
		default:
			continue
		}

		if covers(line, want) {
			return true
		}
	}
	return false
}

// covers reads one retract entry — `v1.0.2` or `[v1.0.0, v1.0.2]` — against a version.
func covers(entry string, want [3]int) bool {
	if !strings.HasPrefix(entry, "[") {
		got, ok := parseSemver(entry)
		return ok && got == want
	}

	low, high, ok := strings.Cut(strings.Trim(entry, "[]"), ",")
	if !ok {
		return false
	}
	lo, okLo := parseSemver(strings.TrimSpace(low))
	hi, okHi := parseSemver(strings.TrimSpace(high))
	// Inclusive at both ends, which is what the go.mod grammar means by a range.
	return okLo && okHi && !greater(lo, want) && !greater(want, hi)
}

func greater(a, b [3]int) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] > b[i]
		}
	}
	return false
}

// CheckedLatest is LatestTag behind the cache, for a checkout. The cache goes in .git/, so it
// is never committed and goes with the clone.
func CheckedLatest(ctx context.Context, root string, ttl time.Duration) (string, error) {
	return Cached(ctx, filepath.Join(root, ".git", checkFile), ttl, Shell{Root: root}.LatestTag)
}

// Cached puts any "what is the newest release" call behind the same one-a-day cache.
//
// Exported and taking the cache path because there are two such calls and one cache: a
// checkout asks its remote with ls-remote, an installed copy asks the forge for its latest
// release, and each keeps its answer somewhere different. Writing a second cache for the
// second caller would be a second TTL, a second failure policy and a second thing to get
// wrong — and the policy is the interesting part, not the storage.
func Cached(ctx context.Context, path string, ttl time.Duration, ask func(context.Context) (string, error)) (string, error) {
	return cached(ctx, path, ttl, time.Now, ask)
}

// cached takes its clock and its asker, because a test that depends on the wall clock behaves
// differently at midnight and one that depends on the network is not a test.
//
// Failure is cached as an empty answer, deliberately. Caching only successes means a
// machine with no network pays the timeout on every launch — the hang this exists to
// prevent, arriving once per invocation instead of once a day. The caller sees no error
// either way: nobody needs a panel that reports its own inability to check.
//
// A directory that is not there means nowhere to keep the answer. Ask and do not cache, rather
// than fail — that is a machine with nothing installed yet.
func cached(
	ctx context.Context,
	path string,
	ttl time.Duration,
	now func() time.Time,
	ask func(context.Context) (string, error),
) (string, error) {
	cacheable := false
	if info, err := os.Stat(filepath.Dir(path)); err == nil && info.IsDir() {
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
