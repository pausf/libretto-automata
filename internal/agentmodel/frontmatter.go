// Package agentmodel reads and writes the model each payload agent runs on.
//
// The model lives in one line of an agent file's frontmatter. That file is a
// prompt, so everything here is built around one promise: change the one line and
// leave every other byte exactly as it was.
//
// Frontmatter here is not YAML. It is a fenced block of `key: value` lines, which
// is the same thing scripts/check-payload reads with two `rg` calls. Parsing it as
// general YAML would mean a dependency for a format neither side ever writes.
package agentmodel

import (
	"fmt"
	"os"
	"strings"
)

// Default is the model of an agent that declares none: whatever the session runs
// on. It is the absence of the key, not a word meaning absence — two spellings of
// one state is a difference somebody eventually treats as meaningful.
const Default = ""

const (
	fence = "---"
	key   = "model:"
)

// ReadModel returns the model an agent file declares, or Default when it declares
// none. A file with no frontmatter is an error: it is not an agent file by the same
// rule scripts/check-payload applies, and guessing otherwise invents one.
func ReadModel(path string) (string, error) {
	lines, closing, err := readFrontmatter(path)
	if err != nil {
		return "", err
	}
	for _, line := range lines[1:closing] {
		if value, ok := modelValue(line); ok {
			return value, nil
		}
	}
	return Default, nil
}

// SetModel declares model in an agent file, replacing any value already there.
// Default removes the key instead of writing a word that means its absence.
//
// Every byte outside the one line survives. The file is split, one element is
// changed, added or dropped, and the rest is joined back untouched — no reflow, no
// re-serialisation, nothing that could normalise a blank line into tidiness nobody
// asked for.
//
// Setting a model an agent already declares writes nothing at all. Identical bytes
// still move the mtime, and a tool that dirties a git tree for a no-op is one people
// stop running before a release.
//
// It does not check the model against the catalogue. Apply does, for the whole set
// at once, before any file is opened — validating here too would put the same rule
// in two places and let a caller reach the weaker one.
func SetModel(path, model string) error {
	lines, closing, err := readFrontmatter(path)
	if err != nil {
		return err
	}

	at := -1
	for i := 1; i < closing; i++ {
		if _, ok := modelValue(lines[i]); ok {
			at = i
			break
		}
	}

	current := Default
	if at >= 0 {
		current, _ = modelValue(lines[at])
	}
	if current == model {
		return nil
	}

	switch {
	case model == Default:
		lines = append(lines[:at], lines[at+1:]...)
	case at >= 0:
		lines[at] = key + " " + model
	default:
		// Last line of the block, so an insert never lands between two keys a
		// reader expects to see together.
		lines = append(lines[:closing], append([]string{key + " " + model}, lines[closing:]...)...)
	}

	return writeInPlace(path, strings.Join(lines, "\n"))
}

// writeInPlace rewrites a file, keeping the permissions it already had.
//
// os.WriteFile's mode applies only when it creates the file, so a mode read from
// disk is the only way an existing agent file keeps whatever the user set on it.
func writeInPlace(path, content string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), fi.Mode().Perm())
}

// readFrontmatter reads the file and locates its frontmatter block, returning every
// line and the index of the closing fence.
//
// The search stops at the closing fence rather than scanning the whole file. An
// agent's body is prose *about* models, and a scan that runs past the fence reports
// a model nobody set — from a line that is documentation.
func readFrontmatter(path string) (lines []string, closing int, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}

	lines = splitKeepingShape(string(data))
	if len(lines) == 0 || strings.TrimRight(lines[0], "\r") != fence {
		return nil, 0, fmt.Errorf("%s: no frontmatter — an agent file opens with %s on line 1", path, fence)
	}

	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(lines[i], "\r") == fence {
			return lines, i, nil
		}
	}
	return nil, 0, fmt.Errorf("%s: frontmatter is not closed", path)
}

// modelValue reports whether a frontmatter line declares the model, and its value.
func modelValue(line string) (string, bool) {
	trimmed := strings.TrimRight(line, "\r")
	if !strings.HasPrefix(trimmed, key) {
		return "", false
	}
	return strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, key)), `"'`), true
}

// splitKeepingShape splits on newlines without inventing or losing a trailing one.
//
// strings.Split on a file ending in "\n" yields a final empty element, and joining
// that back reproduces the file exactly. A file *not* ending in a newline yields no
// such element and also round-trips. That property is what makes the byte-for-byte
// promise in SetModel hold for both shapes.
func splitKeepingShape(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
