package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pausf/libretto-automata/internal/repo"
)

// cloneFunc is repo.Clone's shape, taken as a parameter so bootstrap can be tested
// without a network and without a git binary.
type cloneFunc func(ctx context.Context, url, dest string) error

// bootstrapTimeout bounds the clone.
//
// Generous, because this is the payload arriving over somebody's connection and the user
// asked for it — unlike the release check, which is speculative and gets five seconds.
const bootstrapTimeout = 5 * time.Minute

// bootstrap makes sure dest is a clone of this project, and says so before it does.
//
// The announcement comes first, always. The clone lands in the user's home directory, and
// a tool that writes there and then mentions it is a tool they had no chance to decline.
//
// Three outcomes and no fourth:
//
//	already our clone   nothing happens, and nothing is printed
//	nothing there       announce, clone, report the destination
//	something else      refused, named, nothing touched
//
// It does not prompt. Every path here has to work without a TTY, or `libretto install` in
// CI stops working the day it needs to bootstrap.
func bootstrap(ctx context.Context, out io.Writer, dest, url string, clone cloneFunc) error {
	if isRepo(dest) {
		return nil
	}

	if entries, err := os.ReadDir(dest); err == nil && len(entries) > 0 {
		return fmt.Errorf("%s exists and is not a libretto clone (it contains %s)\n"+
			"       nothing was touched — move it aside, or point %s at a clone you made",
			dest, entries[0].Name(), EnvRoot)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot read %s: %w", dest, err)
	}

	fmt.Fprintf(out, "no clone found — cloning the payload into %s\n", dest)
	fmt.Fprintf(out, "         from %s\n", url)

	ctx, cancel := context.WithTimeout(ctx, bootstrapTimeout)
	defer cancel()

	if err := clone(ctx, url, dest); err != nil {
		// A half-finished clone would be a foreign destination forever, and the user
		// would have to work out for themselves that deleting it is the fix. RemoveAll
		// is safe here: nothing but this call has ever written to dest.
		os.RemoveAll(dest)
		return fmt.Errorf("could not clone into %s: %w", dest, err)
	}

	fmt.Fprintf(out, "cloned   %s\n\n", dest)
	return nil
}

// ensureClone resolves the root and bootstraps it when there is nothing there.
//
// This is what every command that needs the payload calls instead of repoRoot(). The
// override is honoured without a clone attempt when it points at something real — someone
// who set LIBRETTO_ROOT has already decided where the clone is.
func ensureClone() (string, error) {
	root, err := repoRoot()
	if err != nil {
		return "", err
	}
	if isRepo(root) {
		return root, nil
	}
	if err := bootstrap(context.Background(), os.Stdout, root, repo.ModuleURL(), repo.Clone); err != nil {
		return "", err
	}
	return root, nil
}
