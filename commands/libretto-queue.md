---
description: Capture feature ideas into the change queue, one after another. No spec, no plan, no code.
---

Capture only. This command writes proposals and stops there — the queue is drained by
`/libretto-next`, one idea at a time, and that is where the flow actually starts.

Ideas arrive faster than they get built. Without somewhere to put them the choice is
"build it now or lose it", and the thing that loses is whatever was already half
finished.

## The loop

For each idea the user gives:

1. **Name it** from the request — verb-led and readable, `add-relative-discounts`, never
   an invented ticket id. A fake key implies a tracker somebody could consult.
2. Write `.agents/changes/<name>/proposal.md`:

   ```
   # <name>

   Tracker: none
   Queued: 2026-08-11

   ## The ask, verbatim

   > what the user said, in the words they said it in

   ## Reading

   one or two lines on what it is for
   ```

   `Queued:` is today's date, ISO. It is what orders the queue, and it is in the file so
   the order survives a rebase and needs no `git log` walk. Its presence is also what
   marks the change as *not started* — `/libretto-next` removes the line when it picks
   the idea up.

   **The ask goes in verbatim.** Paraphrasing loses the part you did not understand yet,
   and an idea captured now is read weeks later by someone who has forgotten the
   conversation.
3. Say the name back in one line, then **ask for the next idea**.

Keep going until the user says they are done. Then report the whole queue, oldest first,
and point at `/libretto-next`.

## What this never does

- **No spec, no plan, no branch, no code.** An idea is not work started. Commit the
  proposals on the current branch as docs-only commits so the queue is visible from the
  base branch — a branch per captured idea scatters the queue across N branches nobody
  can see.
- **No tracker key as a queued idea.** A key means the tracker is already the source of
  truth: hand it to `/libretto-flow <key>` instead. Copying a ticket into a local
  proposal creates a second one that immediately starts drifting.
- **No priorities, no tags, no reordering.** The queue is FIFO by `Queued:`, and
  `/libretto-next` letting the user pick a different one *is* the reordering mechanism.
- **No editing or deleting through a command.** They are markdown files in a directory.
  Edit them, or `rm -r` the folder.

## Never create the directory in a project that has none

`.agents/changes/` existing means this project already consolidates its changes. If it
does not exist, say so and stop rather than inventing a staging area nobody empties —
the same rule `write-spec` holds.
