package link

import (
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
)

// entry builds a minimal Entry for planning tests. Planning is pure, so nothing
// here needs to exist on disk.
func entry(state State, name string) Entry {
	return Entry{
		Target:   "claude code",
		Kind:     target.Skills,
		Name:     name,
		State:    state,
		RepoPath: "/repo/skills/" + name,
		DestPath: "/home/.claude/skills/" + name,
	}
}

func acts(actions []Action) []Act {
	out := make([]Act, 0, len(actions))
	for _, a := range actions {
		out = append(out, a.Act)
	}
	return out
}

func sameActs(a, b []Act) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestPlan(t *testing.T) {
	cases := []struct {
		name    string
		entries []Entry
		want    []Act
	}{
		{
			name:    "missing becomes create",
			entries: []Entry{entry(Missing, "a")},
			want:    []Act{Create},
		},
		{
			name:    "wrong target becomes repoint",
			entries: []Entry{entry(WrongTarget, "a")},
			want:    []Act{Repoint},
		},
		{
			name:    "conflict becomes skip, never a write",
			entries: []Entry{entry(Conflict, "a")},
			want:    []Act{Skip},
		},
		{
			name:    "linked produces nothing",
			entries: []Entry{entry(Linked, "a")},
			want:    nil,
		},
		{
			name:    "install never removes stale entries",
			entries: []Entry{entry(Stale, "a")},
			want:    nil,
		},
		{
			name: "order follows the scan",
			entries: []Entry{
				entry(Missing, "a"), entry(Linked, "b"),
				entry(Conflict, "c"), entry(WrongTarget, "d"),
			},
			want: []Act{Create, Skip, Repoint},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := acts(Plan(c.entries))
			if !sameActs(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

// R9: running install against a correct tree must do nothing at all.
func TestPlanIsEmptyForACorrectTree(t *testing.T) {
	entries := []Entry{entry(Linked, "a"), entry(Linked, "b"), entry(Linked, "c")}

	if plan := Plan(entries); len(plan) != 0 {
		t.Fatalf("an already-correct tree produced %d action(s): %v", len(plan), plan)
	}
}

func TestPrunePlanTakesOnlyStale(t *testing.T) {
	entries := []Entry{
		entry(Linked, "a"), entry(Missing, "b"), entry(Conflict, "c"),
		entry(WrongTarget, "d"), entry(Stale, "e"), entry(Stale, "f"),
	}

	plan := PrunePlan(entries)
	if want := []Act{Remove, Remove}; !sameActs(acts(plan), want) {
		t.Fatalf("got %v, want %v", acts(plan), want)
	}
	for _, a := range plan {
		if a.Entry.State != Stale {
			t.Fatalf("prune planned %s for a %s entry", a.Act, a.Entry.State)
		}
	}
}

func TestWritesExcludesSkip(t *testing.T) {
	plan := Plan([]Entry{entry(Missing, "a"), entry(Conflict, "b"), entry(WrongTarget, "c")})

	if got := len(plan); got != 3 {
		t.Fatalf("plan has %d actions, want 3", got)
	}
	if got := acts(Writes(plan)); !sameActs(got, []Act{Create, Repoint}) {
		t.Fatalf("writes are %v, want create and repoint only", got)
	}
}

func TestPlanIsPure(t *testing.T) {
	// A plan for paths that do not exist must still be produced. Anything that
	// touched the filesystem here would fail.
	entries := []Entry{{
		State:    Missing,
		Name:     "ghost",
		RepoPath: "/nowhere/skills/ghost",
		DestPath: "/nowhere/.claude/skills/ghost",
	}}

	if plan := Plan(entries); len(plan) != 1 || plan[0].Act != Create {
		t.Fatalf("planning consulted the filesystem: %v", plan)
	}
}
