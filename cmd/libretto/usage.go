package main

// Token cost, read off the Claude Code session transcripts.
//
// Nothing here is instrumented: the transcripts are written by the host for its own
// reasons and this is a free rider on them, the same bar the corrections count clears
// against the lessons ledger. Read-only, always.

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// usageEntry is one assistant turn's cost, with what little the transcript says about
// where it happened. Both labels are frequently absent — measured at 4,216 of 8,944
// entries carrying no attributionSkill — and that is the subject of the unattributed
// bucket, not a parsing failure.
type usageEntry struct {
	branch string
	skill  string

	in     int64
	out    int64
	cacheW int64
	cacheR int64
}

// transcriptEntry is deliberately almost entirely optional. Whole entry types — mode,
// last-prompt, ai-title, file-history-delta — carry no envelope at all, so a struct
// that assumes any field crashes on a real file.
type transcriptEntry struct {
	Type      string `json:"type"`
	GitBranch string `json:"gitBranch"`
	Skill     string `json:"attributionSkill"`
	Message   *struct {
		Usage *struct {
			In     int64 `json:"input_tokens"`
			Out    int64 `json:"output_tokens"`
			CacheW int64 `json:"cache_creation_input_tokens"`
			CacheR int64 `json:"cache_read_input_tokens"`
			// iterations repeats the four numbers above. It is read into nothing on
			// purpose: adding both double-counts every entry that has one.
		} `json:"usage"`
	} `json:"message"`
}

// projectDirName encodes a repository root the way the host does: every separator
// becomes a dash.
//
// Forward only. The encoding is lossy — "-Users-x-gitrepos-promofarma-v3" could decode
// to two different paths — so this is never inverted to learn what a directory holds.
func projectDirName(repoRoot string) string {
	return strings.ReplaceAll(filepath.Clean(repoRoot), string(filepath.Separator), "-")
}

// readUsage returns every assistant turn recorded for repoRoot under projectsDir, and
// whether a project directory for that repository was found at all.
//
// found=false is a state, not an error: a checkout that never hosted a session is the
// normal case, and the caller reports the measurement as unavailable rather than as
// zero.
func readUsage(projectsDir, repoRoot string) (entries []usageEntry, found bool) {
	dir := filepath.Join(projectsDir, projectDirName(repoRoot))
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		return nil, false
	}

	for _, p := range transcriptFiles(dir) {
		entries = append(entries, readTranscript(p)...)
	}
	return entries, true
}

// transcriptFiles lists the two places a session's entries live: one file per session
// at the top level, and one per subagent beneath it.
//
// Enumerated rather than walked for every *.jsonl under the directory, because the
// host also keeps tool-results/ and memory/ there. Counting a file whose entries mirror
// ones already read would inflate the total silently, which is the one failure mode a
// cost report cannot afford.
func transcriptFiles(dir string) []string {
	var out []string
	top, _ := filepath.Glob(filepath.Join(dir, "*.jsonl"))
	sub, _ := filepath.Glob(filepath.Join(dir, "*", "subagents", "*.jsonl"))
	out = append(out, top...)
	return append(out, sub...)
}

// readTranscript streams one file. Files reach 6.8 MB and the directory is 81 MB, so
// nothing here holds a whole file.
//
// bufio.Reader rather than Scanner: a Scanner caps line length and reports the overrun
// as an error, which would silently drop the rest of a file whose entry carried a large
// tool result. A line is read whole however long it is.
func readTranscript(path string) []usageEntry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var out []usageEntry
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if e := parseEntry(line); e != nil {
			out = append(out, *e)
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				// A read error mid-file keeps what was already parsed. Reporting
				// nothing because the last line was unreadable is the reverse of the
				// skip-and-continue rule.
				break
			}
			break
		}
	}
	return out
}

// parseEntry returns the cost of one line, or nil when the line carries none — which
// covers a malformed line, a blank one, and every entry type that is not `assistant`.
// Measured: usage appears on 8,944 assistant entries and on no other type, ever.
func parseEntry(line string) *usageEntry {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	var te transcriptEntry
	if err := json.Unmarshal([]byte(line), &te); err != nil {
		return nil
	}
	if te.Type != "assistant" || te.Message == nil || te.Message.Usage == nil {
		return nil
	}
	u := te.Message.Usage
	return &usageEntry{
		branch: te.GitBranch,
		skill:  te.Skill,
		in:     u.In,
		out:    u.Out,
		cacheW: u.CacheW,
		cacheR: u.CacheR,
	}
}
