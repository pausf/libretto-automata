package repo

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
