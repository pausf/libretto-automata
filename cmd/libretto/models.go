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
	if args[0] != "set" {
		return fmt.Errorf("unknown models command %q — try `models` or `models set <model> <agent>…`", args[0])
	}
	return setModels(root, tg, args[1:])
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
		fmt.Printf("  %-*s  %-12s%s\n", width, a.Name, describe(a.Model), note)
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
		fmt.Printf("  %-10s %-12s %s\n", nameOf(m), m.Version, m.Label)
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

	names, all := args[1:], false
	rest := names[:0]
	for _, a := range names {
		if a == "--all" || a == "-a" {
			all = true
			continue
		}
		rest = append(rest, a)
	}
	names = rest

	// Not "all of them". A destructive default that fires on a forgotten argument
	// is how every agent on the machine silently becomes the same model.
	if !all && len(names) == 0 {
		return fmt.Errorf("name the agents, or pass --all — refusing to guess which you meant")
	}
	if all {
		agents, _, err := agentmodel.Agents(agentsDir(tg))
		if err != nil {
			return err
		}
		names = names[:0]
		for _, a := range agents {
			names = append(names, a.Name)
		}
	}

	dir := agentsDir(tg)

	// Which rows are shared has to be read before the write, so the report can say
	// what each one reached.
	agents, _, err := agentmodel.Agents(dir)
	if err != nil {
		return err
	}
	shared := make(map[string]bool, len(agents))
	for _, a := range agents {
		shared[a.Name] = link.Owned(root, a.Path)
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
		fmt.Printf("  %-20s %-12s%s\n", name, describe(model), note)
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

// describe renders a model for a human. The default is the interesting case: an
// absent key is a state, and printing an empty column would look like a bug.
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
