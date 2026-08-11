# notify-users-of-new-updates Implementation Plan

**Goal:** Every subcommand that reads the payload ends by telling the user a newer
release exists, on stderr, without changing what it exits with.

**Architecture:** One `defer` at the single dispatch point in `run`. The notice text
comes from the existing cached path (`askLatest` → `repo.IsNewer`), so no second
"what is newest" and no second "is this newer". A package-level `var` holds the
producer so the tests never touch the network.

**Tech Stack:** Go 1.26.5, stdlib only. No new dependency.

## Global Constraints

- Never colour the line. `status` and `models` piped carry no escape codes today and
  must still.
- The notice never changes the exit code.
- Silent on: not newer, no cached answer, failed check.
- `doctor`, `update`, `version`, `help` are untouched.
- Every gate passes before the commit: `gofmt -l .`, `go vet ./...`,
  `go test ./... -count=1`, `scripts/check-payload`,
  `skills/record-work/spec-drift --self-test`, `skills/record-work/spec-drift --anchors`.

---

### Task 1: The notice, and the commands that carry it

Spec criteria closed: *goes to stderr*, *silent with nothing to say*, *no escape codes*.

**Files:**
- Create: `cmd/libretto/notice.go`
- Create: `cmd/libretto/notice_test.go`

**Interfaces:**
- Consumes: `askLatest(ctx, root) (string, error)` and `checkTimeout` from
  `cmd/libretto/version.go`; `repo.IsNewer(latest, running string) bool`;
  `invokedAs() string` from `cmd/libretto/main.go`.
- Produces: `carriesNotice(cmd string) bool`, `subcommandNotice(root, running string) string`,
  `noticeAfter(w io.Writer, cmd, root, running string)`, and the package var
  `noticeSource func(root, running string) string`.

- [x] **Step 1: Write the failing tests**

```go
package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSubcommandNoticeGoesToStderr(t *testing.T) {
	stub(t, func(root, running string) string { return "v0.1.0 → v0.2.0 available" })

	var err bytes.Buffer
	noticeAfter(&err, "status", "/root", "v0.1.0")

	if !strings.Contains(err.String(), "v0.2.0 available") {
		t.Fatalf("stderr = %q, want the notice", err.String())
	}
	if !strings.HasSuffix(err.String(), "\n") {
		t.Fatalf("stderr = %q, want a trailing newline", err.String())
	}
}

func TestSubcommandNoticeIsSilentWithNothingToSay(t *testing.T) {
	stub(t, func(root, running string) string { return "" })

	var err bytes.Buffer
	noticeAfter(&err, "status", "/root", "v0.1.0")

	if err.String() != "" {
		t.Fatalf("stderr = %q, want nothing", err.String())
	}
}

func TestDoctorAndUpdateDoNotRepeatTheNotice(t *testing.T) {
	stub(t, func(root, running string) string { return "v0.1.0 → v0.2.0 available" })

	for _, cmd := range []string{"doctor", "update", "version", "help", "unknown"} {
		var err bytes.Buffer
		noticeAfter(&err, cmd, "/root", "v0.1.0")
		if err.String() != "" {
			t.Errorf("%s printed %q, want nothing", cmd, err.String())
		}
	}
}

func TestSubcommandNoticeHasNoEscapeCodes(t *testing.T) {
	stub(t, func(root, running string) string { return "v0.1.0 → v0.2.0 available" })

	for _, cmd := range []string{"status", "preview", "install", "prune", "uninstall", "models"} {
		var err bytes.Buffer
		noticeAfter(&err, cmd, "/root", "v0.1.0")
		if err.Len() == 0 {
			t.Errorf("%s printed nothing, want the notice", cmd)
		}
		if strings.Contains(err.String(), "\x1b") {
			t.Errorf("%s printed an escape code: %q", cmd, err.String())
		}
	}
}

// stub swaps the notice source for the duration of one test. A package var is the
// smallest seam that keeps the tests off the network.
func stub(t *testing.T, f func(root, running string) string) {
	t.Helper()
	prev := noticeSource
	noticeSource = f
	t.Cleanup(func() { noticeSource = prev })
}
```

- [x] **Step 2: Run them and watch them fail**

Run: `go test ./cmd/libretto/ -run 'Notice' -count=1`
Expected: FAIL — `undefined: noticeAfter`, `undefined: noticeSource`.

- [x] **Step 3: Write `cmd/libretto/notice.go`**

```go
package main

import (
	"context"
	"fmt"
	"io"

	"github.com/pausf/libretto-automata/internal/repo"
)

// carriesNotice is which subcommands end by saying a newer release exists.
//
// `doctor` is absent because it already says something on every path, live and
// uncached — a second line would print the same fact twice, once stale. `update` is
// absent because announcing the thing it is doing is noise. `version` and `help`
// answer before the payload is even located.
func carriesNotice(cmd string) bool {
	switch cmd {
	case "status", "preview", "install", "prune", "uninstall", "models":
		return true
	}
	return false
}

// subcommandNotice is the line, or "" when there is nothing to say.
//
// Cached, through the same path the panel uses: a subcommand is not a diagnostic and
// does not get to pay for a network call on every run. Silence on a failed check is
// the panel's rule too — "could not check for a newer release" on every run of a
// command about something else is the noise that gets a notice ignored.
func subcommandNotice(root, running string) string {
	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()

	latest, err := askLatest(ctx, root)
	if err != nil || !repo.IsNewer(latest, running) {
		return ""
	}
	return fmt.Sprintf("%s → %s available · run `%s update`", running, latest, invokedAs())
}

// noticeSource is a var so the tests can answer without reaching a remote.
var noticeSource = subcommandNotice

// noticeAfter writes the notice for cmd, if cmd carries one and there is one to write.
//
// Never coloured. The panel's row is gold because it sits inside a rendered frame; a
// line appended to a command's output is not, and `status` piped carries no escape
// codes by an existing promise.
func noticeAfter(w io.Writer, cmd, root, running string) {
	if !carriesNotice(cmd) {
		return
	}
	if line := noticeSource(root, running); line != "" {
		fmt.Fprintln(w, line)
	}
}
```

- [x] **Step 4: Run them and watch them pass**

Run: `go test ./cmd/libretto/ -run 'Notice' -count=1`
Expected: PASS, all four.

- [x] **Step 5: Gates, then commit**

```bash
gofmt -l . && go vet ./... && go test ./... -count=1
git add cmd/libretto/notice.go cmd/libretto/notice_test.go
git commit -m "feat(cli): say a newer release exists after a subcommand"
```

---

### Task 2: Wire it into the dispatch

Depends on Task 1. Spec criterion closed: *the exit code is whatever the subcommand
returned*.

**Files:**
- Modify: `cmd/libretto/main.go:177` — immediately before `switch args[0]`
- Modify: `cmd/libretto/notice_test.go`

- [x] **Step 1: Write the failing test**

```go
func TestSubcommandNoticeDoesNotChangeTheExitCode(t *testing.T) {
	t.Setenv("CLAUDE_HOME", t.TempDir())
	stub(t, func(root, running string) string { return "v0.1.0 → v0.2.0 available" })

	// A command that succeeds still succeeds.
	if err := run([]string{"status"}); err != nil {
		t.Fatalf("status returned %v, want nil", err)
	}
	// A command that fails still fails, with its own error and not the notice's.
	err := run([]string{"nonesuch"})
	if err == nil || !strings.Contains(err.Error(), "nonesuch") {
		t.Fatalf("unknown command returned %v, want its own error", err)
	}
}
```

- [x] **Step 2: Run it and watch it fail**

Run: `go test ./cmd/libretto/ -run TestSubcommandNoticeDoesNotChangeTheExitCode -count=1`
Expected: FAIL — the notice is not wired in yet, so this passes for the wrong reason.
**If it passes here, add a stderr assertion that proves the notice was printed**, then
watch that fail. A test that is green before the code exists proves nothing.

- [x] **Step 3: Add the defer in `run`**

In `cmd/libretto/main.go`, directly above `switch args[0] {`:

```go
	// After the command, not before: the notice is news about something else, and a
	// user watching `install` scroll should see the result first. `defer` also means
	// it cannot touch what the command returned — being a release behind is not an
	// error, and `install` already uses a non-zero exit for a conflict.
	defer noticeAfter(os.Stderr, args[0], root, version)

	switch args[0] {
```

- [x] **Step 4: Run the whole suite**

Run: `go test ./... -count=1`
Expected: PASS. Watch `TestStatusOutputHasNoEscapeCodes` and `TestRunDispatch`
specifically — they are the ones the new line could break.

- [x] **Step 5: Gates, then commit**

```bash
gofmt -l . && go vet ./... && go test ./... -count=1
git add cmd/libretto/main.go cmd/libretto/notice_test.go
git commit -m "feat(cli): print the release notice after every payload subcommand"
```

---

### Task 3: Land the delta on the capability spec

Depends on Tasks 1 and 2 being green. This is phase 8's work and it happens in the
same commit as nothing else — see `skills/record-work/`.

**Files:**
- Modify: `.agents/specs/cli/spec.md` — under *Saying a newer release exists*, and its
  *Verification criteria*
- Delete: `.agents/changes/notify-users-of-new-updates/`

- [ ] **Step 1:** Fold the delta's outcomes into *Saying a newer release exists* as
  bullets in the section's existing voice, keeping the four decisions (stderr, cached,
  silent, uncoloured) and the two exclusions (`doctor`, `update`).
- [ ] **Step 2:** Move all five `Proof:` lines into the cli spec's *Verification criteria*.
- [ ] **Step 3:** Run the anchor gate.

```bash
skills/record-work/spec-drift --anchors > /tmp/anchors.out; echo $?
```
Expected: exit 0. Every one of the five names a test that now exists.

- [ ] **Step 4:** `rm -r .agents/changes/notify-users-of-new-updates/` and commit the
  spec and the deletion together.

---

## What can start now

Task 1. Tasks 2 and 3 are strictly sequential behind it.
