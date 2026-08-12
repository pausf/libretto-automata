package main

import (
	"fmt"
	"strings"

	"github.com/pausf/libretto-automata/internal/agentmodel"
	"github.com/pausf/libretto-automata/internal/link"
	"github.com/pausf/libretto-automata/internal/target"
)

// models shows which model each agent runs on, and sets it.
//
// The subcommand is named after the field it edits. It was very nearly called
// something about saving a token, which would have been wrong twice: nothing here
// is a token, and nothing is saved in the sense that word implies.
func models(root string, tg target.Target, args []string) error {
	if len(args) == 0 {
		return listModels(root, tg)
	}
	switch args[0] {
	case "set":
		return setModels(root, tg, args[1:])
	case "effort":
		// A verb rather than a flag on `set`. `set opus --effort xhigh` forces the
		// model to be restated to change the depth, and a restated model is a
		// write nobody asked for — the two keys are independent and stay so.
		return setEffort(root, tg, args[1:])
	}
	return fmt.Errorf("unknown models command %q — try `models`, `models set <model> <agent>…` or `models effort <level> <agent>…`", args[0])
}

// sharedMark flags an agent whose file this repository owns.
//
// Writing it reaches the repository's file, and therefore every target linking to it.
// Writing an unmarked agent reaches that target only. The word names the consequence
// rather than the mechanism, because the consequence is what the reader has to act on.
const sharedMark = "shared"

// agentsDir is the directory of the target under discussion. Empty when the target
// does not take agents at all, which lists nothing rather than failing.
func agentsDir(tg target.Target) string {
	if !tg.Accepts(target.Agents) {
		return ""
	}
	return tg.Dir(target.Agents)
}

func listModels(root string, tg target.Target) error {
	agents, unreadable, err := agentmodel.Agents(agentsDir(tg))
	if err != nil {
		return err
	}
	if len(agents) == 0 {
		fmt.Printf("no agents in %s\n", tg.Name())
		return nil
	}

	width := 0
	for _, a := range agents {
		if n := len([]rune(a.Name)); n > width {
			width = n
		}
	}

	fmt.Printf("%s  %s\n", tg.Name(), tg.Root())
	for _, a := range agents {
		note := ""
		if link.Owned(root, a.Path) {
			note = "  " + sharedMark
		}
		fmt.Printf("  %-*s  %-12s%-10s%s\n", width, a.Name, describe(a.Model), describe(a.Effort), note)
	}
	// Skipping in silence would trade a loud failure for a quiet one. `doctor` and
	// `prune` own stale links; this only has to say they are there.
	if len(unreadable) > 0 {
		fmt.Println()
		fmt.Printf("  %d link(s) here point at nothing: %s\n", len(unreadable), strings.Join(unreadable, ", "))
		fmt.Println("  run `libretto doctor` — `prune --yes` removes them")
	}

	fmt.Println()
	fmt.Printf("models available (aliases; versions as of %s):\n", agentmodel.Resolved)
	for _, m := range agentmodel.Catalogue() {
		note := ""
		if !m.Effort {
			note = "  — no effort levels"
		}
		fmt.Printf("  %-10s %-12s %s%s\n", nameOf(m), m.Version, m.Label, note)
	}

	fmt.Println()
	fmt.Println("effort available (weakest first; `models effort <level> <agent>…`):")
	fmt.Printf("  %-10s %s\n", "default", "no key at all: whatever the session runs at")
	for _, e := range agentmodel.Efforts() {
		fmt.Printf("  %-10s %s\n", e.Name, e.Label)
	}
	return nil
}

func setModels(root string, tg target.Target, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("models set needs a model — one of %s", known())
	}
	model := args[0]
	if model == "default" {
		model = agentmodel.Default
	}

	dir := agentsDir(tg)
	names, err := targetNames(dir, args[1:])
	if err != nil {
		return err
	}

	// Which rows are shared, and which carry an effort, both have to be read before
	// the write: the report says what each row reached, and a cleared effort is only
	// visible by comparison with what was there.
	agents, _, err := agentmodel.Agents(dir)
	if err != nil {
		return err
	}
	shared := make(map[string]bool, len(agents))
	effort := make(map[string]string, len(agents))
	for _, a := range agents {
		shared[a.Name] = link.Owned(root, a.Path)
		effort[a.Name] = a.Effort
	}

	if err := agentmodel.Apply(dir, names, model); err != nil {
		return err
	}

	anyShared := false
	for _, name := range names {
		note := ""
		if shared[name] {
			note, anyShared = "  "+sharedMark, true
		}
		// A dropped effort is an edit to a prompt file, and a silent one is the
		// failure the clearing exists to avoid. Say it on the row it happened to.
		cleared := ""
		if !agentmodel.SupportsEffort(model) && effort[name] != agentmodel.Default {
			cleared = fmt.Sprintf("  effort %s cleared — %s has no levels", effort[name], describe(model))
		}
		fmt.Printf("  %-20s %-12s%s%s\n", name, describe(model), note, cleared)
	}

	// Only the shared rows reach beyond this target. Saying it about all of them was
	// true when the subject was always this repository's own file; saying it about a
	// target-local write is the same class of error as the silence it replaced.
	fmt.Println()
	if anyShared {
		fmt.Printf("the rows marked %s live in %s/agents and are linked from more than one\n", sharedMark, root)
		fmt.Println("destination — those take effect for every project on this machine.")
	} else {
		fmt.Printf("written into %s — this destination only.\n", dir)
	}
	return nil
}

// setEffort writes one effort level onto a set of agents.
//
// It is `setModels` with the other key, and it shares that function's argument
// handling rather than restating it: the `--all` refusal is a safety property and two
// copies of one is how a caller reaches the weaker version.
func setEffort(root string, tg target.Target, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("models effort needs a level — one of %s", knownEfforts())
	}
	effort := args[0]
	if effort == "default" {
		effort = agentmodel.Default
	}
	// ApplyEffort refuses an unknown level too, and says which. Naming the five here
	// is the difference between a wrong guess and a next move: `ultracode` is the
	// name a user is most likely to reach for, and it is the one the list corrects.
	if effort != agentmodel.Default && !agentmodel.ValidEffort(effort) {
		return fmt.Errorf("unknown effort %q — one of %s", effort, knownEfforts())
	}

	dir := agentsDir(tg)
	names, err := targetNames(dir, args[1:])
	if err != nil {
		return err
	}

	agents, _, err := agentmodel.Agents(dir)
	if err != nil {
		return err
	}
	shared := make(map[string]bool, len(agents))
	for _, a := range agents {
		shared[a.Name] = link.Owned(root, a.Path)
	}

	if err := agentmodel.ApplyEffort(dir, names, effort); err != nil {
		return err
	}

	anyShared := false
	for _, name := range names {
		note := ""
		if shared[name] {
			note, anyShared = "  "+sharedMark, true
		}
		fmt.Printf("  %-20s %-12s%s\n", name, describe(effort), note)
	}

	fmt.Println()
	if anyShared {
		fmt.Printf("the rows marked %s live in %s/agents and are linked from more than one\n", sharedMark, root)
		fmt.Println("destination — those take effect for every project on this machine.")
	} else {
		fmt.Printf("written into %s — this destination only.\n", dir)
	}
	return nil
}

// targetNames resolves the agent arguments a write acts on: the names given, or every
// agent in the directory when --all is passed.
//
// Shared by both writing verbs. The refusal below is the reason it is shared — it is
// not a convenience, it is the thing standing between a forgotten argument and every
// agent on the machine silently becoming the same thing.
func targetNames(dir string, args []string) ([]string, error) {
	names, all := args, false
	rest := names[:0]
	for _, a := range names {
		if a == "--all" || a == "-a" {
			all = true
			continue
		}
		rest = append(rest, a)
	}
	names = rest

	if !all && len(names) == 0 {
		return nil, fmt.Errorf("name the agents, or pass --all — refusing to guess which you meant")
	}
	if !all {
		return names, nil
	}

	agents, _, err := agentmodel.Agents(dir)
	if err != nil {
		return nil, err
	}
	names = names[:0]
	for _, a := range agents {
		names = append(names, a.Name)
	}
	return names, nil
}

// describe renders a model or an effort for a human. The default is the interesting
// case for both: an absent key is a state, and printing an empty column would look
// like a bug.
func describe(model string) string {
	if model == agentmodel.Default {
		return "(session)"
	}
	return model
}

func nameOf(m agentmodel.Model) string {
	if m.Name == agentmodel.Default {
		return "default"
	}
	return m.Name
}

func known() string {
	out := ""
	for _, m := range agentmodel.Catalogue() {
		if out != "" {
			out += ", "
		}
		out += nameOf(m)
	}
	return out
}

func knownEfforts() string {
	out := "default"
	for _, e := range agentmodel.Efforts() {
		out += ", " + e.Name
	}
	return out
}
