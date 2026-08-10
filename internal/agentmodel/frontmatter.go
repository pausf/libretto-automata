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
