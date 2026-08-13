package main

// Which model each payload agent suits, and why.
//
// This lives here rather than in internal/agentmodel deliberately. That package settled,
// from a real failure, that **the subject is a directory, not the repository** — it is
// handed a directory and works on every *.md in it, whoever created them. A map from
// `review-lens-security` to a tier is this payload's agent list, and putting it there
// would reintroduce the exact coupling that decision exists against. cmd/libretto already
// knows what the payload is; it gates commands on it.
//
// **Recommend only.** Nothing in this binary writes a model or an effort because of
// anything below. The reading is the repository's; the typing is the user's.

// recommendation is a suggestion about one agent, and never a bare verdict: the reason is
// the half that lets somebody disagree with it.
type recommendation struct {
	model  string
	effort string // empty when the recommended model has no effort levels
	reason string
}

// recommendations is keyed by agent name, matched whole. Every value is checked against
// the catalogue by a test rather than by review — an entry naming a model the tool would
// refuse to type is the one thing worse than no recommendation at all.
//
// Nothing is recommended onto opus or fable. The catalogue orders itself cheapest-first
// because this exists to lower a bill, and recommending the two priciest tiers inverts
// the thing it was built for. An agent that genuinely needs more is a decision a person
// makes at the screen.
var recommendations = map[string]recommendation{
	"review-lens-design": {"haiku", "", "pattern-matching over prose"},
	"review-lens-tests":  {"haiku", "", "pattern-matching over prose"},
	"review-lens-intent": {"sonnet", "high", "the only lens that runs commands, and it reads forge payload"},

	"review-lens-reliability": {"sonnet", "high", "races and error paths are reasoning, not matching"},
	"review-lens-security":    {"sonnet", "xhigh", "the skill says security does not run cheap"},

	"spec-writer":   {"sonnet", "high", "it writes the contract everything later is measured against"},
	"work-reviewer": {"sonnet", "high", "a cheap miss here is a false green"},
}

// recommend returns what this repository suggests for an agent, and whether it has an
// opinion at all.
//
// false for every agent a user writes themselves, and that is the point: silence, never a
// guess. A recommendation invented for somebody else's agent is worse than none.
func recommend(name string) (recommendation, bool) {
	r, ok := recommendations[name]
	return r, ok
}
