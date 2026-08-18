# Delta — the wiki opens in the browser

Targets: cli

`wiki --open` ends with the viewer on screen: it writes (or refreshes) the HTML
view and hands it to the platform's opener. The panel's wiki row uses it, so a
menu press ends in a browser tab instead of a path the user has to go and open.

## Outcomes

`libretto wiki --open` in a project with specs leaves `wiki.html` current and the
default browser showing it, via `file://` — no server, no port, no daemon.
Pressing the panel's `wiki` row does the same. A failed generation (foreign file,
no specs) opens nothing.

## Scope boundaries

In: the `--open` flag (implies the HTML view), a small opener seam so tests
observe the launch without launching, the platform argv choice as a pure
function, and the panel row dispatching `--open` with its description saying so.

Out, named:

- **an HTTP server / real localhost URL.** `file://` is the same browser tab with
  no lifetime to manage inside a TUI; `--serve` comes back if someone needs a URL
  other people (or a `file://`-blocking browser policy) can reach.
- **choosing a browser.** The platform opener decides, as it does for everything.
- **opening the markdown view.** A browser does not render `README.md`; the open
  path is the HTML view by construction.

## Constraints

- Tests never launch a browser: the opener is a seam (`var`), and the argv
  selection is a pure function tested per platform string.
- `--open` and `--html --open` behave identically — `--open` implies the view it
  opens.
- Failure ordering: generation errors surface exactly as today, and the opener is
  only reached after a successful write.
- If the opener fails after a successful write, then the run shall report the
  written path together with the opener's error and exit non-zero — the file
  stays, and the message hands the user what to open by hand. (Settled
  2026-08-18 from the cutter's finding; the builder does not decide it silently.)

## Prior decisions

- **`file://` over a localhost server.** The ask was "the browser opens with the
  wiki"; the page is proven self-contained from `file://`, and a server inside a
  panel action is a lifetime nobody manages. Bring-back condition for `--serve`:
  a URL someone else must reach, or a browser policy blocking `file://`.
  2026-08-18.
- **The panel row runs `wiki --open`, by instruction** — the user asked for the
  press to end on screen, 2026-08-18. This amends the earlier "the row runs the
  plain command": the row still dispatches the same subcommand, now with the one
  flag that makes a menu press finish in front of the user. The plain run stays
  what landings use.
- **The opener is the platform's own** — `open` on darwin, `xdg-open` elsewhere —
  chosen by a pure function so the mapping is testable without executing it.

## Task breakdown

1. The flag, the seam, the argv function, the panel row switch — with the proofs.

## Verification criteria

- When invoked with `--open`, the `wiki` subcommand shall write the HTML view
  exactly as `--html` does and then invoke the opener seam with the generated
  file's path, after the write succeeds; if the opener errors, the run shall
  surface that error naming the written path, non-zero.
  Proof: cmd/libretto/wiki_test.go TestWikiOpenGeneratesAndOpensTheViewer
- If generation refuses or fails, then the `--open` run shall not invoke the
  opener.
  Proof: cmd/libretto/wiki_test.go TestWikiOpenDoesNotOpenOnFailure
- The opener argv shall name `open` on darwin and `xdg-open` on every other
  platform, with the file's path as the argument.
  Proof: cmd/libretto/wiki_test.go TestOpenerArgvPerPlatform
- When the panel's `wiki` row is dispatched, it shall run the wiki with
  `--open`, and its menu description shall say the viewer opens.
  Proof: cmd/libretto/wiki_test.go TestPanelWikiRowOpensTheViewer

**This replaces a criterion the capability spec landed earlier the same day** —
"the panel shall run the plain `wiki` command … leaving the generated
`README.md`", proven by `TestDispatchRunsWiki`. That sentence is rewritten in the
same commit that renames its test — `--anchors` runs per commit and a capability
citing a renamed test is red immediately, so this one piece of the delta cannot
wait for the landing. `TestDispatchRunsWiki` becomes the proof of the new
behaviour under the new name above. No release ever carried the
old promise — both versions ride the same PR — so this is an amendment, not a
shipped reversal.
