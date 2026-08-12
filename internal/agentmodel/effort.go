package agentmodel

import "fmt"

// The catalogue of effort levels an agent may declare.
//
// Effort controls adaptive reasoning: how often and how deeply the model thinks on
// each step. It is the other half of what `model:` buys — that key chooses the tier,
// this one chooses the depth inside it, and the two are independent because keeping a
// model while spending fewer of its tokens is the case the feature exists for.
//
// Static for the same reason the model catalogue is: an organisation can cap which
// levels are available, and finding out means an authenticated call this tool is
// forbidden to make. The Resolved date above covers this table too — it decays for
// the same reason and a second date is a second thing to forget.
//
// `ultracode` is deliberately absent. It sends xhigh to the model and additionally
// has Claude orchestrate dynamic workflows, which makes it a Claude Code session
// setting rather than a level; no frontmatter accepts it.

// Effort is one catalogue entry.
type Effort struct {
	// Name is what goes in the frontmatter. Default is the empty string.
	Name string

	// Label is the one line shown beside it in the CLI and the panel.
	Label string
}

// efforts is ordered weakest first, which is the same contract the model catalogue
// carries: two surfaces render this list, and the cheap end goes under the cursor of
// a feature whose purpose is to reduce the bill.
var efforts = []Effort{
	{"low", "short, scoped work that is not intelligence-sensitive"},
	{"medium", "cost-sensitive work that can trade off some intelligence"},
	{"high", "the balance point — and the host's own default"},
	{"xhigh", "deeper reasoning at higher spend"},
	{"max", "the deepest; prone to overthinking. Measure before adopting"},
}

// Efforts returns the legal levels, weakest first.
func Efforts() []Effort {
	out := make([]Effort, len(efforts))
	copy(out, efforts)
	return out
}

// ValidEffort reports whether name may be written into an agent's frontmatter.
//
// Default is not valid here, for the reason Valid gives: removing the key is a
// different act from declaring a level, and the callers that offer both say so
// separately.
func ValidEffort(name string) bool {
	for _, e := range efforts {
		if e.Name == name {
			return true
		}
	}
	return false
}

// SupportsEffort reports whether a catalogue model runs effort levels at all.
//
// Haiku does not, and an unknown model is treated as though it does — the same
// posture the rest of this package takes towards things it cannot verify. Writing a
// level into an agent whose model has no concept of one would leave a line in a
// prompt file claiming a setting that does nothing, which is the confident wrong
// answer the model catalogue's own comment warns about.
//
// Default — an agent declaring no model — supports it. The session model is
// unknowable from here, and refusing would be a guess in the other direction.
func SupportsEffort(model string) bool {
	for _, m := range catalogue {
		if m.Name == model {
			return m.Effort
		}
	}
	return true
}

// ReadEffort returns the effort an agent file declares, or Default when it declares
// none — meaning the agent runs at whatever the session runs at, not at the host's
// default of `high`. That distinction is the host's business; this reports the file.
func ReadEffort(path string) (string, error) { return readKey(path, effortKey) }

// SetEffort declares effort in an agent file, replacing any value already there.
// Default removes the key. Every byte outside the one line survives — see SetModel,
// which shares the implementation and states the promise in full.
func SetEffort(path, effort string) error { return setKey(path, effortKey, effort) }

// ApplyEffort declares effort on every named agent, or changes nothing at all.
//
// The whole set is validated before the first file is opened: the level is in the
// catalogue, every name is an agent, every one of those files is writable, and — the
// check this function exists for — every one of them runs a model that has effort
// levels. Validating as it went would leave half a set written and no way to know how
// far it got.
func ApplyEffort(dir string, names []string, effort string) error {
	if len(names) == 0 {
		return fmt.Errorf("no agents named — nothing marked does not mean all of them")
	}
	if effort != Default && !ValidEffort(effort) {
		return fmt.Errorf("unknown effort %q", effort)
	}

	targets, err := resolve(dir, names)
	if err != nil {
		return err
	}

	// Only when a level is being written. Removing one is legal on any agent — a
	// Haiku agent that somehow carries the key should be able to shed it.
	if effort != Default {
		for _, a := range targets {
			if !SupportsEffort(a.Model) {
				return fmt.Errorf("%s runs on %s, which has no effort levels", a.Name, a.Model)
			}
		}
	}

	for _, a := range targets {
		if err := SetEffort(a.Path, effort); err != nil {
			return err
		}
	}
	return nil
}
