package agentmodel

import (
	"strings"
	"testing"
)

func TestCatalogueListsTheSubscriptionModels(t *testing.T) {
	got := Catalogue()

	var names []string
	for _, m := range got {
		names = append(names, m.Name)
	}
	want := []string{Default, "haiku", "sonnet", "opus"}

	if len(names) != len(want) {
		t.Fatalf("Catalogue() = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("Catalogue()[%d] = %q, want %q", i, names[i], want[i])
		}
	}
}

// Every entry has to say something the user can act on. A catalogue row with an
// empty label renders as a blank line in the panel, which is worse than no
// catalogue at all: it looks like a model with no name.
func TestEveryCatalogueEntryIsLabelled(t *testing.T) {
	for _, m := range Catalogue() {
		if strings.TrimSpace(m.Label) == "" {
			t.Errorf("model %q has no label", m.Name)
		}
	}
}

// The binary cannot ask Anthropic what this account is paying for, and AGENTS.md
// forbids it from touching a credential to find out. So a label may name a plan
// tier, and must never claim the user has it.
func TestNoLabelClaimsTheUserHasAPlan(t *testing.T) {
	forbidden := []string{"your plan", "you have", "your subscription"}
	for _, m := range Catalogue() {
		low := strings.ToLower(m.Label)
		for _, phrase := range forbidden {
			if strings.Contains(low, phrase) {
				t.Errorf("model %q label claims to know the user's plan: %q", m.Name, m.Label)
			}
		}
	}
}

func TestUnknownModelIsRefused(t *testing.T) {
	for _, name := range []string{"gpt-4", "claude", "Haiku", "haiku ", "opus-5"} {
		if Valid(name) {
			t.Errorf("Valid(%q) = true, want false", name)
		}
	}
}

func TestEveryCatalogueEntryIsValid(t *testing.T) {
	for _, m := range Catalogue() {
		if m.Name == Default {
			continue
		}
		if !Valid(m.Name) {
			t.Errorf("Valid(%q) = false for a model the catalogue lists", m.Name)
		}
	}
}
