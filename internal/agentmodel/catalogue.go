package agentmodel

// The catalogue of models an agent may declare.
//
// It is static, and that is a decision rather than a shortcut. Finding out which
// models an account may actually use means an authenticated call, and AGENTS.md
// forbids this tool from accepting, storing or transmitting a credential. A list
// that is honest about being a list beats a detection that would need a token.
//
// So a label may name the plan tier a model belongs to. None of them claims the
// user has that tier — the binary does not know, and saying otherwise would be the
// kind of confident wrong answer that gets believed.

// Model is one catalogue entry.
type Model struct {
	// Name is what goes in the frontmatter. Default is the empty string.
	Name string

	// Label is the one line shown beside it in the CLI and the panel.
	Label string
}

// catalogue is ordered cheapest first, with the session default at the top.
//
// Order is contracted, not incidental: the panel and the CLI both render it, and a
// list whose first entry is the most expensive model puts the costly choice under
// the cursor of a feature whose whole purpose is to reduce the bill.
var catalogue = []Model{
	{Default, "the session's model — whatever you are running"},
	{"haiku", "cheapest; fine for pattern-matching over prose"},
	{"sonnet", "the everyday working model"},
	{"opus", "most capable; Max plans, metered on Pro"},
}

// Catalogue returns the legal models, cheapest first.
func Catalogue() []Model {
	out := make([]Model, len(catalogue))
	copy(out, catalogue)
	return out
}

// Valid reports whether name may be written into an agent's frontmatter.
//
// Default is not valid here: removing the key is a different act from declaring a
// model, and the callers that offer both say so separately.
func Valid(name string) bool {
	for _, m := range catalogue {
		if m.Name != Default && m.Name == name {
			return true
		}
	}
	return false
}
