package main

import (
	"testing"

	"github.com/pausf/libretto-automata/internal/agentmodel"
)

func TestAnUnknownAgentGetsNoRecommendation(t *testing.T) {
	for _, name := range []string{"", "my-own-agent", "review-lens", "review-lens-designer"} {
		if r, ok := recommend(name); ok {
			t.Errorf("recommend(%q) = %+v, want no opinion — the match is whole, and a user's own agent is not our business", name, r)
		}
	}
}

func TestEveryRecommendationCarriesAReason(t *testing.T) {
	for name, r := range recommendations {
		if r.reason == "" {
			t.Errorf("%s recommends %q with no reason — a verdict without one is an instruction", name, r.model)
		}
	}
}

func TestEveryRecommendationIsTypeable(t *testing.T) {
	// The guard the table cannot check by review. An entry naming a model the catalogue
	// refuses, or a level the recommended model has none of, is a suggestion this binary
	// would decline to type — which is worse than saying nothing.
	for name, r := range recommendations {
		if !agentmodel.Valid(r.model) {
			t.Errorf("%s recommends %q, which the catalogue does not know", name, r.model)
		}
		if r.effort == "" {
			continue
		}
		if !agentmodel.ValidEffort(r.effort) {
			t.Errorf("%s recommends effort %q, which is not a level", name, r.effort)
		}
		if !agentmodel.SupportsEffort(r.model) {
			t.Errorf("%s recommends effort %q on %s, which has no levels", name, r.effort, r.model)
		}
	}
}

func TestNothingIsRecommendedOntoThePriciestTiers(t *testing.T) {
	// The catalogue is ordered cheapest-first because this feature exists to lower a
	// bill. A table steering onto the two priciest tiers inverts what it was built for.
	for name, r := range recommendations {
		if r.model == "opus" || r.model == "fable" {
			t.Errorf("%s is recommended onto %q — nothing is steered onto the priciest tiers", name, r.model)
		}
	}
}
