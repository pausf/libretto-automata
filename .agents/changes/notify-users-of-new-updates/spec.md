# Say a newer release exists on every subcommand

Targets: cli

This is a delta on `.agents/specs/cli/spec.md`, under *Saying a newer release exists*.

## Outcomes

Today the notice lives in exactly two places, and both need the user to go looking: the
panel's row (needs a TTY) and `libretto doctor` (needs somebody to suspect there is
something to diagnose). After installing, most people run neither for months — the
payload is used inside Claude Code, not through this binary.

So: **whenever the user runs the CLI at all, for any reason, they find out.**

- **Every subcommand that already reads the payload ends by saying a newer release is
  available** — `status`, `preview`, `install`, `prune`, `uninstall`, `models`.
- **The line goes to stderr**, after the command's own output, and reads the same as the
  panel's row: `v0.5.2 → v0.6.0 available · run `libretto update``, naming the command as
  it was invoked.
- **It never changes the exit code.** Being a release behind is news. `install` exiting
  non-zero already means a conflict, and a second meaning for that code would make the
  first one unreadable.
- **It is silent unless there is something to say.** Not newer, no cached answer, check
  failed, remote unreachable — nothing is printed. Exactly the panel's rule.
- **It costs no network call the panel would not already have made.** It goes through the
  same cached path (`askLatest`, `repo.CheckTTL`), so a run inside the TTL prints from the
  cache and a run outside it pays at most `checkTimeout`.

## Scope boundaries

**In:** the six subcommands above, the stderr line, and the tests that prove all three of
its properties.

**Out, and named so it cannot be quietly added:**

- **`doctor` is not touched.** It already says something on every path, live, and it is
  the one command whose answer must not come from the cache — going through the cache
  swallows the error the line exists to report. Adding a second line would print the same
  fact twice, once stale.
- **`update` is not touched.** It is the act of updating; announcing the thing it is doing
  is noise.
- **`version` and `help` are not touched.** They answer before the payload is located and
  read nothing, which is a promise the cli spec already makes.
- **The payload is not reopened.** A skill that announced CLI versions would put this
  repository's tooling inside an installed skill, which `payload` forbids precisely so a
  skill works in a project that never heard of `libretto`. That boundary stays.
- **The limit stays a limit.** If the user never runs the CLI, no CLI notice reaches
  them. This change does not claim to solve that; it guarantees that *every* run tells
  them, where today only one of eight does.
- **No new flag.** No `--no-update-check`, no env var to silence it. Nobody has asked, and
  the line is one line on stderr.

## Constraints

- Piped `status` and `models` carry no escape codes — an existing promise with a test
  behind it (`TestStatusOutputHasNoEscapeCodes`). The line is therefore **never coloured
  at all**, and on stderr rather than stdout so a script parsing `status` sees
  byte-identical output to today. The panel's row is gold because it sits inside a
  rendered frame; a line appended to a command's output is not in a frame, and "colour it
  only when stderr is a TTY" is a branch, a theme dependency and a test to buy something
  nobody asked for.
- One dispatch point exists (`run`'s switch in `cmd/libretto/main.go`); the notice is
  emitted there, once, not in six commands.
- The comparison stays `repo.IsNewer`. Formatting is this side's job; deciding what
  "newer" means is not, and two implementations is two that disagree.

## Prior decisions

- **stderr, not stdout.** The alternative is a line on stdout, and it loses on the fact
  that `status` output is parseable and somebody may be parsing it. Streams already mean
  something here: stdout is the answer, stderr is everything else — `main` already writes
  errors and usage there.
- **Cached, never live.** `doctor` is the live one, by an existing decision recorded in
  the cli spec. A user running `install` moments after a release still waits out the TTL,
  and that is accepted: the alternative is every subcommand paying a network round trip
  for speculative news.
- **Silence on no answer.** `doctor` earns "could not check" because the user typed a
  diagnostic. `status` did not, and "could not check for a newer release" on every run of
  a command about something else is the noise that gets a notice ignored.
- Ceiling, named: **this reaches nobody who does not run the binary.** The upgrade path,
  if that ever matters, is the payload route — and it costs the `payload` boundary above.

## Task breakdown

1. Emit the notice from `run`'s dispatch for the six subcommands, to stderr, after the
   command returns and without touching its error.
2. Tests: stream, silence, exit code, and that `doctor`/`update` do not repeat it.
3. Apply the delta onto `.agents/specs/cli/spec.md`.

## Verification criteria

- **the notice goes to stderr, and stdout is byte-identical to a run without it**
  Proof: cmd/libretto/notice_test.go TestSubcommandNoticeGoesToStderr
- **nothing is printed when there is no newer release, no cached answer, or a failed
  check**
  Proof: cmd/libretto/notice_test.go TestSubcommandNoticeIsSilentWithNothingToSay
- **the exit code is whatever the subcommand returned** — proven on both a zero and a
  non-zero path
  Proof: cmd/libretto/notice_test.go TestSubcommandNoticeDoesNotChangeTheExitCode
- **`doctor` and `update` do not carry it**, so the fact is never printed twice
  Proof: cmd/libretto/notice_test.go TestDoctorAndUpdateDoNotRepeatTheNotice
- **the line carries no escape codes**, on any stream
  Proof: cmd/libretto/notice_test.go TestSubcommandNoticeHasNoEscapeCodes
