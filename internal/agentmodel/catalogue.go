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
	//
	// It is an alias, not a version. `opus` keeps meaning "the Opus tier" after
	// the model behind it is replaced, which is why the frontmatter carries the
	// alias and not a pinned id — an agent file written today should not need
	// editing the day a new Opus ships.
	Name string

	// Version is the model the alias resolves to as of Resolved. Empty for the
	// session default, which resolves to whatever the session is running.
	//
	// This is the field that decays. It is here anyway because `opus` alone does
	// not answer the question the user is actually asking, and an answer that
	// goes stale in public beats no answer at all — Resolved is what makes the
	// staleness visible instead of silent.
	Version string

	// Label is the one line shown beside it in the CLI and the panel.
	Label string

	// Effort reports whether this model runs effort levels at all.
	//
	// It lives on the model rather than in a list beside the effort catalogue
	// because it is a fact about the model, and a second list is a second thing to
	// forget when an entry is added. SupportsEffort reads it.
	Effort bool
}

// Resolved is when the Version column was last checked against the model
// catalogue.
//
// ponytail: a hand-maintained date, and it will rot exactly the way this
// repository's own "117 tests" badge rotted. The upgrade path, the day it
// matters: read the aliases from the host rather than restating them. Nothing
// here can ask, so until then the date is the honesty mechanism — a reader can
// see how old the claim is instead of trusting it blind.
const Resolved = "2026-08"

// catalogue is ordered cheapest first, with the session default at the top.
//
// Order is contracted, not incidental: the panel and the CLI both render it, and a
// list whose first entry is the most expensive model puts the costly choice under
// the cursor of a feature whose whole purpose is to reduce the bill.
// The Effort column is read off the host's own table, which lists Fable 5, Opus 5,
// Sonnet 5, Opus 4.8, 4.7 and 4.6, and Sonnet 4.6 — and no Haiku. The session default
// claims support because the session's model is unknowable from here, and refusing
// would be the guess.
var catalogue = []Model{
	{Default, "", "the session's model — whatever you are running", true},
	{"haiku", "Haiku 4.5", "cheapest; fine for pattern-matching over prose", false},
	{"sonnet", "Sonnet 5", "the everyday working model", true},
	{"opus", "Opus 5", "most capable; Max plans, metered on Pro", true},
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
