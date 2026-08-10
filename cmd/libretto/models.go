package main

import (
	"fmt"
	"os"
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

// notLinkedHere marks an agent the repository has but this target does not.
//
// The model itself is one file and cannot differ between targets. Which agents
// actually reach a target can, and that is the one honest thing the scope flag
// changes about this command — without it the two listings are identical and the
// flag is decoration.
const notLinkedHere = "· not linked here"

func listModels(root string, tg target.Target) error {
	agents, err := agentmodel.Agents(root)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no agents in this repository")
			return nil
		}
		return err
	}
	if len(agents) == 0 {
		fmt.Println("no agents in this repository")
		return nil
	}

	reaches, err := agentsReaching(root, tg)
	if err != nil {
		return err
	}

	fmt.Printf("%s  %s\n", tg.Name(), tg.Root())
	for _, a := range agents {
		note := ""
		if !reaches[a.Name] {
			note = "  " + notLinkedHere
		}
		fmt.Printf("  %-20s %s%s\n", a.Name, describe(a.Model), note)
	}
	fmt.Println()
	fmt.Printf("models available (aliases; versions as of %s):\n", agentmodel.Resolved)
	for _, m := range agentmodel.Catalogue() {
		fmt.Printf("  %-10s %-12s %s\n", nameOf(m), m.Version, m.Label)
	}
	return nil
}

// agentsReaching reports which agents are actually installed in a target, keyed by
// the name agentmodel uses — without the .md that link items carry.
//
// Only `Linked` counts. A conflict is somebody else's file sitting where ours would
// go, and an agent whose slot is occupied does not reach the target however much the
// repository wishes it did.
func agentsReaching(root string, tg target.Target) (map[string]bool, error) {
	if !tg.Accepts(target.Agents) {
		return nil, nil
	}

	entries, err := link.Scan(root, tg)
	if err != nil {
		return nil, err
	}

	reaches := make(map[string]bool, len(entries))
	for _, e := range entries {
		if e.Kind == target.Agents && e.State == link.Linked {
			reaches[strings.TrimSuffix(e.Name, ".md")] = true
		}
	}
	return reaches, nil
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
		agents, err := agentmodel.Agents(root)
		if err != nil {
			return err
		}
		names = names[:0]
		for _, a := range agents {
			names = append(names, a.Name)
		}
	}

	if err := agentmodel.Apply(root, names, model); err != nil {
		return err
	}

	for _, name := range names {
		fmt.Printf("  %-20s %s\n", name, describe(model))
	}

	// The scope flag promises something it cannot deliver. Both targets symlink to
	// the same file in this repository, so the change is shared — and a user who
	// typed --project has every reason to expect otherwise.
	fmt.Println()
	fmt.Printf("written into %s/agents — in effect for every project on this machine,\n", root)
	fmt.Printf("because %s links to those files rather than holding copies.\n", tg.Root())
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
