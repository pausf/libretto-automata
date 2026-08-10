package link

// Planning turns a scan into a list of intended actions. It is a pure function of
// the entries handed to it: no filesystem access, no writes, no surprises.
//
// Keeping it pure is what makes the dangerous part testable. Every question about
// what the tool *would* do is answered here, in memory, without a directory in
// sight.

// Act is what to do about one entry.
type Act string

const (
	// Create makes a link that does not exist yet.
	Create Act = "create"

	// Repoint replaces one of our own links that aims at the wrong item.
	Repoint Act = "repoint"

	// Remove deletes one of our own links with no item behind it.
	Remove Act = "remove"

	// Skip touches nothing. It exists so a plan can say out loud what it refuses
	// to do — a tool whose promise is "never clobbers your work" should report
	// the things it declined, not silently omit them.
	Skip Act = "skip"
)

// Action is one intended operation, with the entry that justified it.
type Action struct {
	Act   Act
	Entry Entry
}

// Writes reports whether this action changes the filesystem.
func (a Action) Writes() bool { return a.Act != Skip }

// Plan lists what `install` intends to do.
//
// Linked entries produce nothing: an already-correct tree yields an empty plan,
// which is what makes install safe to run repeatedly (SPEC R9). Conflicts produce
// Skip, never a write. Stale entries are absent entirely — removing things is
// prune's job and install must not do it as a side effect.
func Plan(entries []Entry) []Action {
	var actions []Action
	for _, e := range entries {
		switch e.State {
		case Missing:
			actions = append(actions, Action{Act: Create, Entry: e})
		case WrongTarget:
			actions = append(actions, Action{Act: Repoint, Entry: e})
		case Conflict:
			actions = append(actions, Action{Act: Skip, Entry: e})
		case Linked, Stale:
			// Nothing to do, and deliberately not prune's business either way.
		}
	}
	return actions
}

// PrunePlan lists what `prune` intends to do.
//
// Only Stale entries, which by construction are links this repo owns whose item
// is gone. Nothing else is ever removed: not conflicts, not foreign links, not
// entries that merely look wrong.
func PrunePlan(entries []Entry) []Action {
	var actions []Action
	for _, e := range entries {
		if e.State == Stale {
			actions = append(actions, Action{Act: Remove, Entry: e})
		}
	}
	return actions
}

// UninstallPlan lists what `uninstall` intends to do: undo this repository's work at
// one destination.
//
// It is the pair of Plan, not of PrunePlan. Prune cleans up after the *repo* changed —
// a renamed item leaves a link pointing at nothing — and deliberately ignores links
// that are correct, which is the guarantee that makes it safe to run. Uninstall
// removes links that are **working**, because the user changed their mind rather than
// because the repo did.
//
//	linked        remove   ours, working — the whole point
//	wrong target  remove   ours, misaimed. Undoing does not need it repaired first.
//	stale         remove   ours, orphaned
//	conflict      skip     not ours. Reported, never touched.
//	missing       —        nothing is there
//
// No new applying logic: Apply already re-checks ownership at the moment of writing
// and already removes a link without following it. This is a different question about
// *which* entries, answered by a pure function over machinery that is already proven.
func UninstallPlan(entries []Entry) []Action {
	var actions []Action
	for _, e := range entries {
		switch e.State {
		case Linked, WrongTarget, Stale:
			actions = append(actions, Action{Act: Remove, Entry: e})
		case Conflict:
			// Stated, not omitted. A tool whose promise is "never clobbers your
			// work" reports the refusal rather than going quiet about it.
			actions = append(actions, Action{Act: Skip, Entry: e})
		case Missing:
			// Nothing to undo.
		}
	}
	return actions
}

// Writes returns only the actions that would change the filesystem.
func Writes(actions []Action) []Action {
	var out []Action
	for _, a := range actions {
		if a.Writes() {
			out = append(out, a)
		}
	}
	return out
}
