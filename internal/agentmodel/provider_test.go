package agentmodel

import (
	"os"
	"strings"
	"testing"
)

// clearProvider unsets every variable detection reads, so a test starts from a known
// environment rather than from whatever the developer's shell exports.
//
// t.Setenv("", ...) cannot unset, so the values are saved and restored by hand. Without
// this the suite passes on a laptop and fails in a Bedrock CI job, which is the class of
// flake that gets a test deleted rather than fixed.
func clearProvider(t *testing.T) {
	t.Helper()
	for _, n := range []string{
		"CLAUDE_CODE_USE_BEDROCK", "CLAUDE_CODE_USE_MANTLE", "CLAUDE_CODE_USE_VERTEX",
		"ANTHROPIC_FOUNDRY_RESOURCE", "ANTHROPIC_FOUNDRY_BASE_URL",
		"ANTHROPIC_AWS_WORKSPACE_ID", "ANTHROPIC_AWS_BASE_URL",
		"ANTHROPIC_BASE_URL",
		"ANTHROPIC_DEFAULT_OPUS_MODEL", "ANTHROPIC_DEFAULT_SONNET_MODEL",
		"ANTHROPIC_DEFAULT_HAIKU_MODEL", "ANTHROPIC_DEFAULT_FABLE_MODEL",
	} {
		if prev, ok := os.LookupEnv(n); ok {
			os.Unsetenv(n)
			t.Cleanup(func() { os.Setenv(n, prev) })
		}
	}
}

// Nothing set is the Anthropic API, and Detected stays false: "found the default" and
// "found nothing and fell back to the default" are the same answer for different reasons.
func TestNoProviderVariablesMeansTheAnthropicAPI(t *testing.T) {
	clearProvider(t)

	p := DetectProvider()
	if p.Detected {
		t.Error("a provider was reported as detected with nothing set")
	}
	if !strings.Contains(p.Name, "Anthropic API") {
		t.Errorf("provider = %q, want the Anthropic API", p.Name)
	}
	for alias, want := range map[string]string{"opus": "Opus 5", "sonnet": "Sonnet 5", "haiku": "Haiku 4.5"} {
		got, ok := p.Resolve(alias)
		if !ok || got != want {
			t.Errorf("%s resolves to %q (%v), want %q", alias, got, ok, want)
		}
	}
}

// The whole reason this file exists: on Amazon Bedrock `sonnet` is Sonnet 4.5, which
// appears in no row of the effort table. The alias supporting all five was only ever
// right on the Anthropic API.
func TestOnBedrockSonnetHasNoEffortLevels(t *testing.T) {
	clearProvider(t)
	t.Setenv("CLAUDE_CODE_USE_BEDROCK", "1")

	p := DetectProvider()
	if !p.Detected || p.Name != "Amazon Bedrock" {
		t.Fatalf("provider = %q detected=%v, want Amazon Bedrock", p.Name, p.Detected)
	}
	if got := p.effortsFor("sonnet"); len(got) != 0 {
		t.Errorf("sonnet on Bedrock offers %v, want nothing — it resolves to Sonnet 4.5", got)
	}
	if got := p.effortsFor("opus"); len(got) != len(allFive) {
		t.Errorf("opus on Bedrock offers %v, want all five — it resolves to Opus 5", got)
	}
}

// Opus 4.6 has four of the five, which is why EffortsFor returns a slice and not a bool.
func TestOnFoundryOpusLosesXHigh(t *testing.T) {
	clearProvider(t)
	t.Setenv("ANTHROPIC_FOUNDRY_RESOURCE", "some-resource")

	p := DetectProvider()
	if p.Name != "Microsoft Foundry" {
		t.Fatalf("provider = %q, want Microsoft Foundry", p.Name)
	}
	got := p.effortsFor("opus")
	if len(got) != len(noXHigh) {
		t.Fatalf("opus on Foundry offers %v, want %v", got, noXHigh)
	}
	for _, level := range got {
		if level == "xhigh" {
			t.Error("xhigh was offered on Opus 4.6, which does not have it")
		}
	}
}

// Mantle is checked before Bedrock because a session can have both set, and Mantle serves
// Sonnet 5 where the Invoke API serves Sonnet 4.5. Order-dependent, so it is pinned.
func TestMantleWinsOverBedrockWhenBothAreSet(t *testing.T) {
	clearProvider(t)
	t.Setenv("CLAUDE_CODE_USE_BEDROCK", "1")
	t.Setenv("CLAUDE_CODE_USE_MANTLE", "1")

	p := DetectProvider()
	if !strings.Contains(p.Name, "Mantle") {
		t.Fatalf("provider = %q, want Mantle to win", p.Name)
	}
	if got := p.effortsFor("sonnet"); len(got) != len(allFive) {
		t.Errorf("sonnet on Mantle offers %v, want all five — it is Sonnet 5 there", got)
	}
}

// An explicit pin is exactly the thing that overrides a provider default, so it wins.
func TestAPinnedModelWinsOverTheProviderDefault(t *testing.T) {
	clearProvider(t)
	t.Setenv("CLAUDE_CODE_USE_BEDROCK", "1")
	t.Setenv("ANTHROPIC_DEFAULT_SONNET_MODEL", "us.anthropic.claude-sonnet-4-6")

	p := DetectProvider()
	got, ok := p.Resolve("sonnet")
	if !ok || got != "Sonnet 4.6" {
		t.Fatalf("pinned sonnet resolves to %q (%v), want Sonnet 4.6", got, ok)
	}
	// Sonnet 4.6 has four of the five, so the pin changes the offer as well as the name.
	if levels := p.effortsFor("sonnet"); len(levels) != len(noXHigh) {
		t.Errorf("pinned sonnet offers %v, want %v", levels, noXHigh)
	}
}

func TestVersionOfReadsEveryProviderIDShape(t *testing.T) {
	cases := map[string]string{
		"us.anthropic.claude-sonnet-4-6":              "Sonnet 4.6",
		"claude-opus-5":                               "Opus 5",
		"anthropic.claude-haiku-4-5":                  "Haiku 4.5",
		"us.anthropic.claude-haiku-4-5-20251001-v1:0": "Haiku 4.5",
		"eu.anthropic.claude-opus-5":                  "Opus 5",
		"us-gov.anthropic.claude-sonnet-4-5-20250929": "Sonnet 4.5",
		"claude-fable-5":                              "Fable 5",
	}
	for id, want := range cases {
		got, ok := versionOf(id)
		if !ok || got != want {
			t.Errorf("versionOf(%q) = %q (%v), want %q", id, got, ok, want)
		}
	}
}

// An inference profile ARN or a Foundry deployment name is a real answer of "not
// knowable" rather than a fall through to the provider default. Guessing past the pin
// would report the model the user deliberately replaced.
func TestAnUnparseablePinIsUnknownRatherThanTheDefault(t *testing.T) {
	clearProvider(t)
	t.Setenv("CLAUDE_CODE_USE_BEDROCK", "1")
	t.Setenv("ANTHROPIC_DEFAULT_SONNET_MODEL",
		"arn:aws:bedrock:us-east-2:123456789012:application-inference-profile/sonnet-prod")

	p := DetectProvider()
	if got, ok := p.Resolve("sonnet"); ok {
		t.Errorf("an ARN resolved to %q; it is not knowable from here", got)
	}
	// Unknown is capable: refusing a level on a model that may well support it is the
	// more expensive mistake.
	if levels := p.effortsFor("sonnet"); len(levels) != len(allFive) {
		t.Errorf("an unresolvable pin offers %v, want all five", levels)
	}
}

// ANTHROPIC_BASE_URL changes where requests go, not which model answers them, so the
// alias table behind a gateway is unknown — reported, never assumed to be the default.
func TestAGatewayIsReportedAsUnknownRatherThanAssumed(t *testing.T) {
	clearProvider(t)
	t.Setenv("ANTHROPIC_BASE_URL", "https://llm-gateway.internal.example.com")

	p := DetectProvider()
	if !p.Detected {
		t.Error("a gateway was not reported as detected")
	}
	if _, ok := p.Resolve("sonnet"); ok {
		t.Error("an alias resolved behind a gateway, which cannot be known from here")
	}
	if levels := p.effortsFor("sonnet"); len(levels) != len(allFive) {
		t.Errorf("behind a gateway sonnet offers %v, want all five", levels)
	}
}

// Anthropic's own URL is not a gateway. Setting it explicitly is common in CI.
func TestAnthropicsOwnBaseURLIsNotAGateway(t *testing.T) {
	clearProvider(t)
	t.Setenv("ANTHROPIC_BASE_URL", "https://api.anthropic.com")

	if p := DetectProvider(); p.Detected {
		t.Errorf("provider = %q detected=true; Anthropic's own URL is the default", p.Name)
	}
}

// The host documents an empty string as the way to override a stale export from a shell
// profile you do not control. A detection that treated "" as on would strand that user.
func TestAnEmptyOrZeroFlagIsOff(t *testing.T) {
	for _, value := range []string{"", "0", "false"} {
		clearProvider(t)
		t.Setenv("CLAUDE_CODE_USE_BEDROCK", value)
		if p := DetectProvider(); p.Detected {
			t.Errorf("CLAUDE_CODE_USE_BEDROCK=%q was read as on (%s)", value, p.Name)
		}
	}
}

// An agent declaring no model runs on whatever the session runs on, and the session is
// not this process. Unknowable by construction, and capable.
func TestTheSessionDefaultIsNeverResolved(t *testing.T) {
	clearProvider(t)

	p := DetectProvider()
	if got, ok := p.Resolve(Default); ok {
		t.Errorf("the session default resolved to %q", got)
	}
	if !SupportsEffort(Default) {
		t.Error("the session default was refused an effort level")
	}
}

// Detection is os.Getenv and nothing else. The capability forbids this binary touching a
// credential, and reading an API key's value would be touching one — so the names it
// reads are the non-secret ones, and this test is what keeps that true as providers are
// added.
func TestDetectionNeverReadsASecret(t *testing.T) {
	clearProvider(t)

	for _, secret := range []string{
		"ANTHROPIC_FOUNDRY_API_KEY", "ANTHROPIC_FOUNDRY_AUTH_TOKEN",
		"ANTHROPIC_AWS_API_KEY", "AWS_BEARER_TOKEN_BEDROCK", "ANTHROPIC_API_KEY",
	} {
		t.Setenv(secret, "sk-do-not-read-this")
	}

	p := DetectProvider()
	if p.Detected {
		t.Errorf("provider %q was detected from a credential variable alone", p.Name)
	}
}
