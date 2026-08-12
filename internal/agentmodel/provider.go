package agentmodel

import (
	"os"
	"regexp"
	"strings"
)

// What an alias means on this machine.
//
// `opus` and `sonnet` do not name the same model everywhere. The host resolves them per
// provider, and effort levels are a property of the concrete model rather than of the
// alias — so a catalogue that answers "sonnet supports all five levels" is answering for
// the Anthropic API and quietly guessing everywhere else. On Amazon Bedrock `sonnet` is
// Sonnet 4.5, which supports no effort at all.
//
// This reads the environment to answer honestly. **It is not the network call the
// capability refuses.** No request, no credential, no token: `os.Getenv` on variables the
// host already documents, and — deliberately — never on a variable that holds a secret.
// Presence of `ANTHROPIC_FOUNDRY_RESOURCE` says as much as reading an API key would, and
// a binary that never touches the key cannot leak it by accident.
//
// Read off code.claude.com/docs/en/model-config (*Model aliases*, *Adjust effort level*),
// /docs/en/amazon-bedrock (*Pin model versions*) and /docs/en/env-vars, as of Resolved.

// Provider is where the environment says requests go, and what the aliases mean there.
type Provider struct {
	// Name is what the listing prints. Phrased as the docs phrase it.
	Name string

	// Detected is false when nothing in the environment named a provider, which means
	// the Anthropic API. Stated rather than implied: "detected the default" and "found
	// nothing and assumed the default" are the same answer for different reasons, and
	// only one of them is worth a second look.
	Detected bool

	// versions maps a catalogue alias to the model it resolves to here. An alias absent
	// from the map is one this table cannot answer for.
	versions map[string]string
}

// The five levels, and the four. Named so the tables below read as data.
var (
	allFive = []string{"low", "medium", "high", "xhigh", "max"}
	noXHigh = []string{"low", "medium", "high", "max"}
)

// effortByVersion is the host's own table, keyed by the concrete model.
//
// **Presence in this map means "known".** A nil value is a model known to have no effort
// levels; an absent key is a model this build has never heard of, and those are treated
// as capable — the same posture the rest of this package takes towards what it cannot
// verify. Collapsing the two into one bool would make an unrecognised pin
// indistinguishable from Haiku, and refusing on a model that may well support effort is
// the more expensive mistake.
var effortByVersion = map[string][]string{
	"Fable 5":    allFive,
	"Opus 5":     allFive,
	"Opus 4.8":   allFive,
	"Opus 4.7":   allFive,
	"Sonnet 5":   allFive,
	"Opus 4.6":   noXHigh,
	"Sonnet 4.6": noXHigh,

	// Listed in no row of the effort table, so known to have none.
	"Sonnet 4.5": nil,
	"Haiku 4.5":  nil,
}

// providers is what each alias resolves to per provider, from the docs' own table.
var providers = []struct {
	name     string
	env      func() bool
	versions map[string]string
}{
	{
		// Checked before Bedrock: a Mantle session can have both variables set, and it
		// serves Sonnet 5 where the Invoke API serves Sonnet 4.5.
		name: "Amazon Bedrock (Mantle)",
		env:  func() bool { return truthy("CLAUDE_CODE_USE_MANTLE") },
		// Mantle's lineup is granted per account and its IDs carry no version suffix,
		// so only the alias defaults are answerable here.
		versions: map[string]string{"opus": "Opus 5", "sonnet": "Sonnet 5", "haiku": "Haiku 4.5"},
	},
	{
		name:     "Amazon Bedrock",
		env:      func() bool { return truthy("CLAUDE_CODE_USE_BEDROCK") },
		versions: map[string]string{"opus": "Opus 5", "sonnet": "Sonnet 4.5", "haiku": "Haiku 4.5"},
	},
	{
		name:     "Google Cloud's Agent Platform",
		env:      func() bool { return truthy("CLAUDE_CODE_USE_VERTEX") },
		versions: map[string]string{"opus": "Opus 5", "sonnet": "Sonnet 4.5", "haiku": "Haiku 4.5"},
	},
	{
		// The resource name and the base URL, never the API key or the auth token.
		name:     "Microsoft Foundry",
		env:      func() bool { return set("ANTHROPIC_FOUNDRY_RESOURCE", "ANTHROPIC_FOUNDRY_BASE_URL") },
		versions: map[string]string{"opus": "Opus 4.6", "sonnet": "Sonnet 4.5", "haiku": "Haiku 4.5"},
	},
	{
		// The workspace id, never the workspace API key.
		name:     "Claude Platform on AWS",
		env:      func() bool { return set("ANTHROPIC_AWS_WORKSPACE_ID", "ANTHROPIC_AWS_BASE_URL") },
		versions: map[string]string{"opus": "Opus 5", "sonnet": "Sonnet 4.6", "haiku": "Haiku 4.5"},
	},
}

// anthropicAPI is what the aliases mean with nothing set.
//
// `fable` appears here and in none of the provider tables above, and that is an answer
// rather than an omission: those tables transcribe the host's own per-provider alias table,
// which does not name Fable. An alias absent from a provider's map resolves to *not
// knowable*, which is this package's standing posture towards what it cannot verify — and
// unknown is treated as capable, so `fable` is offered all five levels there, which is what
// Fable 5 runs anyway.
var anthropicAPI = map[string]string{
	"fable":  "Fable 5",
	"opus":   "Opus 5",
	"sonnet": "Sonnet 5",
	"haiku":  "Haiku 4.5",
}

// DetectProvider reads the environment and reports which provider it names.
//
// **A gateway is its own answer rather than a guess, and the docs say two things about
// that which read as contradictory.** Model configuration notes that `ANTHROPIC_BASE_URL`
// "changes where requests are sent, not which model answers them" — which argues the
// alias table is unaffected. But the model-name check section says that behind an LLM
// gateway or a custom `ANTHROPIC_BASE_URL`, "your provider or gateway defines the model
// names, so Claude Code passes any string through without checking it."
//
// The second is the one that settles it: if the gateway defines the names, what `sonnet`
// means there is the gateway's business and not this binary's. So a non-Anthropic host is
// reported as unresolvable. The weaker line is quoted here because reading it alone is
// enough to talk yourself into deleting this branch, which nearly happened.
func DetectProvider() Provider {
	for _, p := range providers {
		if p.env() {
			return Provider{Name: p.name, Detected: true, versions: p.versions}
		}
	}

	if url := os.Getenv("ANTHROPIC_BASE_URL"); url != "" && !strings.Contains(url, "api.anthropic.com") {
		// No versions map at all: every alias is unanswerable, which is the honest
		// answer for a gateway that may route anywhere.
		return Provider{Name: "a gateway (ANTHROPIC_BASE_URL)", Detected: true}
	}

	return Provider{Name: "the Anthropic API", versions: anthropicAPI}
}

// Resolve reports the concrete model a catalogue alias means here, and whether that is
// knowable at all.
//
// An explicit pin wins over the provider's default, because that is exactly what a pin
// is for. Default — an agent declaring no model — is unknowable by construction: it runs
// on whatever the session runs on, and the session is not this process.
func (p Provider) Resolve(alias string) (string, bool) {
	if alias == Default {
		return "", false
	}
	if pinned := os.Getenv("ANTHROPIC_DEFAULT_" + strings.ToUpper(alias) + "_MODEL"); pinned != "" {
		if v, ok := versionOf(pinned); ok {
			return v, true
		}
		// A pin this build cannot parse — an application inference profile ARN, a
		// deployment name — is a real answer of "not knowable", not a fall through to
		// the provider default. The pin is what is in force; guessing past it would
		// report the model the user deliberately replaced.
		return "", false
	}
	v, ok := p.versions[alias]
	return v, ok
}

// pinPattern picks the tier and version out of a model id, whatever the provider wraps
// it in: `us.anthropic.claude-sonnet-4-6`, `claude-opus-5`,
// `anthropic.claude-haiku-4-5-20251001-v1:0`.
var pinPattern = regexp.MustCompile(`claude-(fable|opus|sonnet|haiku)-(\d+)(?:-(\d+))?`)

// versionOf turns a model id into the catalogue's own version wording.
func versionOf(id string) (string, bool) {
	m := pinPattern.FindStringSubmatch(strings.ToLower(id))
	if m == nil {
		return "", false
	}
	tier := strings.ToUpper(m[1][:1]) + m[1][1:]
	version := m[2]
	if m[3] != "" {
		version += "." + m[3]
	}
	return tier + " " + version, true
}

// truthy reports whether an on/off variable is on. The host documents `1` or `true`, and
// an empty string is how a stale export is overridden.
func truthy(name string) bool {
	switch strings.ToLower(os.Getenv(name)) {
	case "", "0", "false":
		return false
	}
	return true
}

// set reports whether any of the named variables carries a value. Only ever called with
// non-secret names — see the note at the top of this file.
func set(names ...string) bool {
	for _, n := range names {
		if os.Getenv(n) != "" {
			return true
		}
	}
	return false
}
