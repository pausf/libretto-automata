package agentmodel

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// writable reports whether a file can be opened for writing, without changing it.
//
// O_WRONLY with no O_TRUNC and no O_CREATE asks the question and touches nothing.
// Checking the permission bits instead would be a guess: ownership, ACLs and a
// read-only mount all answer differently from what the mode says.
func writable(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	return f.Close()
}

// Agent is one payload agent, the model it currently runs on, and how hard that
// model is told to think.
type Agent struct {
	Name   string
	Model  string
	Effort string
	Path   string
}

// Agents lists every agent file in a directory, sorted by name.
//
// The directory is handed in rather than derived. That is what lets one package serve
// both the repository's own `agents/` and an install target's — without learning what
// a target is, and while staying testable against a bare t.TempDir().
//
// **Every `*.md` here is an agent**, whether this tool created it or not. Ownership is
// the caller's question; this package works on what it is pointed at.
//
// A directory that does not exist reports no agents rather than an error: a target
// that has never had one installed is a state, and making every caller special-case
// os.IsNotExist to render an empty list is how that state becomes a crash.
//
// An entry that cannot be opened at all — a stale symlink whose destination was
// renamed or deleted — is skipped and returned by name in the second value, rather
// than failing the whole listing. Renaming one agent leaves a stale link in every
// target that had the old name, and taking down eleven readable agents because a
// twelfth is dangling is a listing that punishes the ordinary case. `prune` owns
// stale links; this only has to survive them.
//
// **A file that is present and is not an agent is still an error.** Apply's
// all-or-nothing guarantee rests on it — skipping it would let `--all` write around
// something somebody put there deliberately.
//
// Sorted because two surfaces render this list — the CLI and the panel — and a list
// whose order depends on readdir is a list that reorders itself between machines.
func Agents(dir string) (agents []Agent, unreadable []string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}
		path := filepath.Join(dir, name)
		model, err := ReadModel(path)
		if err != nil {
			if os.IsNotExist(err) {
				unreadable = append(unreadable, strings.TrimSuffix(name, ".md"))
				continue
			}
			return nil, nil, err
		}
		// The file has already been read and its frontmatter proven, so the second
		// read cannot fail for a reason the first survived.
		effort, err := ReadEffort(path)
		if err != nil {
			return nil, nil, err
		}
		agents = append(agents, Agent{
			Name:   strings.TrimSuffix(name, ".md"),
			Model:  model,
			Effort: effort,
			Path:   path,
		})
	}

	sort.Slice(agents, func(i, j int) bool { return agents[i].Name < agents[j].Name })
	sort.Strings(unreadable)
	return agents, unreadable, nil
}

// Apply declares model on every named agent, or changes nothing at all.
//
// The whole set is checked before the first file is opened: the model is in the
// catalogue, every name is an agent the repo has, and every one of those files has
// frontmatter to write into. Validating as it went would leave a half-applied set
// and no way to know how far it got — which is worse than a refusal, because the
// user's next move depends on knowing the state.
//
// ponytail: the guarantee covers the failures this tool can foresee, not a disk that
// fills between the third write and the fourth. A staged write into temporary files
// and an atomic rename per file would close that too — worth doing the day an agent
// file is big enough or a target slow enough for it to happen.
func Apply(dir string, names []string, model string) error {
	if len(names) == 0 {
		return fmt.Errorf("no agents named — nothing marked does not mean all of them")
	}
	if model != Default && !Valid(model) {
		return fmt.Errorf("unknown model %q", model)
	}

	targets, err := resolve(dir, names)
	if err != nil {
		return err
	}

	for _, a := range targets {
		if err := SetModel(a.Path, model); err != nil {
			return err
		}
		// A model with no effort levels makes any declared effort inert, and a key
		// that claims a setting the model has no concept of is a lie left in a
		// prompt file. Dropping it is the honest half of the same act; the callers
		// report the row it happened on, because a silent edit is the other
		// failure.
		if !SupportsEffort(model) && a.Effort != Default {
			if err := SetEffort(a.Path, Default); err != nil {
				return err
			}
		}
	}
	return nil
}

// resolve turns a list of agent names into the agents themselves, refusing the whole
// set if any name is unknown or any file cannot be written.
//
// Shared by Apply and ApplyEffort. The all-or-nothing guarantee is the same guarantee
// in both, and two copies of it is one place for the two to drift apart.
func resolve(dir string, names []string) ([]Agent, error) {
	agents, _, err := Agents(dir)
	if err != nil {
		return nil, err
	}
	known := make(map[string]Agent, len(agents))
	for _, a := range agents {
		known[a.Name] = a
	}

	targets := make([]Agent, 0, len(names))
	for _, name := range names {
		a, ok := known[name]
		if !ok {
			return nil, fmt.Errorf("no such agent: %s", name)
		}
		// Agents() already read this file's frontmatter, so a broken one has
		// failed above. Readable is not writable, though, and a read-only file in
		// the middle of the set would strand the ones before it — so ask now,
		// while refusing is still free.
		if err := writable(a.Path); err != nil {
			return nil, err
		}
		targets = append(targets, a)
	}
	return targets, nil
}
