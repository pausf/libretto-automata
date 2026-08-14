package target

import (
	"bytes"
	"fmt"
	"strings"
)

// MarkerKey is the frontmatter key a generated file carries to say where it came
// from.
//
// It is what makes a generated file provably ours: `ownership` reads it and refuses
// anything without it. The `x-` prefix is not decoration — OpenCode sweeps
// frontmatter keys it does not know into an `options` map rather than rejecting the
// file (packages/core/src/v1/config/agent.ts, normalize()), so the marker is inert
// in the host that reads it.
const MarkerKey = "x-libretto-source"

// Transformer is implemented by a target that installs some kind by writing a
// derived file rather than by symlinking the source.
//
// Discovered by interface assertion, never required of every target — the same way
// Exists() is, and for the same reason: a target that has nothing to say about
// transforming should not have to answer.
//
// The caller reads the source and passes its bytes. source is used only to write the
// marker. Keeping it a pure function of the bytes is what lets `link-state` compute
// what *should* be at a destination without touching the filesystem, and what keeps
// the comparison free of a read that could race.
// Transforms is asked first and answers per kind, because a target may transform one
// kind and symlink another — OpenCode does exactly that. Without it a caller would
// have to call Transform to find out, and an error would be indistinguishable from
// "this kind is linked": the first version of this interface had only Transform, and
// every skill and command in a transforming target was classified as a conflict
// because the error came back for a kind that was perfectly fine.
type Transformer interface {
	Transforms(k Kind) bool
	Transform(k Kind, source string, content []byte) ([]byte, error)
}

// quoteYAML renders a path as a YAML double-quoted scalar.
//
// Always quoted, never conditionally. A checkout path may contain a `#`, a leading or
// trailing space, or a colon — and unquoted, ` #` starts a YAML comment, so the value
// the host reads would be a different path than the one written. OpenCode *throws* on
// an agent it cannot parse rather than skipping it, so a path that needs quoting and
// does not get it breaks that host's whole config load.
//
// Quoting unconditionally also keeps the transform deterministic in one branch fewer.
func quoteYAML(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"', '\\':
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	b.WriteByte('"')
	return b.String()
}

// frontmatter splits a markdown file into its frontmatter lines and the rest.
//
// The block is the region between an opening `---` on line 1 and the next `---`. A
// file that does not open with `---`, or never closes it, has no frontmatter — and
// that is an error for a transform rather than an empty result, because emitting an
// agent with no frontmatter would produce a file OpenCode throws on.
func frontmatter(content []byte) (keys []string, body []byte, err error) {
	text := string(content)
	if !strings.HasPrefix(text, "---\n") {
		return nil, nil, fmt.Errorf("no frontmatter: file does not open with ---")
	}
	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		// A file whose frontmatter is the whole content, closing on the last line.
		if strings.HasSuffix(rest, "\n---") {
			return strings.Split(strings.TrimSuffix(rest, "\n---"), "\n"), nil, nil
		}
		return nil, nil, fmt.Errorf("no frontmatter: opening --- is never closed")
	}
	return strings.Split(rest[:end], "\n"), []byte(rest[end+len("\n---\n"):]), nil
}

// Transforms reports which kinds OpenCode installs by writing a derived file.
// Agents only: skills and commands are read from symlinks unchanged.
func (o Opencode) Transforms(k Kind) bool { return k == Agents }

// Transform rewrites one Claude-format agent into the file OpenCode reads.
//
// Only Agents is transformed. Skills and commands are symlinked, and asking for a
// transform of one is a caller bug rather than something to answer with the bytes
// unchanged — an error here is how that bug surfaces instead of shipping.
func (o Opencode) Transform(k Kind, source string, content []byte) ([]byte, error) {
	if k != Agents {
		return nil, fmt.Errorf("opencode installs %s by symlink, not by transform", k)
	}

	keys, body, err := frontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", source, err)
	}

	// Kept in the order encountered, so output is a function of input alone. Emitting
	// from a map would reorder between runs and every scan would report drift.
	var kept []string
	for _, line := range keys {
		switch {
		case line == "":
			continue
		// tools: is dropped rather than mapped. On Claude it is a comma-separated
		// string and a safety property — what keeps a review lens read-only. OpenCode
		// types it Record<string, boolean> over *its* tool names, which do not
		// correspond one-to-one and have no entry for Skill, so a mapping would be a
		// guess whose failure mode is a reviewer that can write. Leaving the key out
		// leaves OpenCode's own default in charge, which its user chose.
		case strings.HasPrefix(line, "tools:"):
			continue
		// model: is dropped for the narrower reason that a Claude tier name is not a
		// provider/model-id and means nothing there.
		case strings.HasPrefix(line, "model:"):
			continue
		// A source that already carries a marker is being re-transformed. Drop the old
		// one rather than emitting two.
		case strings.HasPrefix(line, MarkerKey+":"):
			continue
		case strings.HasPrefix(line, "mode:"):
			continue
		default:
			kept = append(kept, line)
		}
	}

	var out bytes.Buffer
	out.WriteString("---\n")
	for _, line := range kept {
		out.WriteString(line)
		out.WriteByte('\n')
	}
	out.WriteString("mode: subagent\n")
	out.WriteString(MarkerKey + ": " + quoteYAML(source) + "\n")
	out.WriteString("---\n")
	out.Write(body)
	return out.Bytes(), nil
}
