# show-the-panel-in-the-readme

Tracker: none

## What was asked

> en el readme podrias poner una imagen de como se ve el cli

Said in conversation on 2026-08-07, right after verifying the panel live over a
temporary `CLAUDE_HOME` — the ask is to show that panel to people reading the README,
before they install anything.

## Why

The panel is the front door of the CLI and the README describes it only in words. A
reader deciding whether to install judges a TUI by how it looks; a screenshot answers
that in one glance where a paragraph does not.

## Open for the spec to settle

- Screenshot (PNG checked into the repo) versus terminal-to-SVG capture versus asking
  the user to take it — and where the asset lives.
- Which state the panel shows: fresh install (everything missing) or a healthy one
  (everything linked). What it shows is what it promises.
- One caution already known: the README is the one place `𝄞` is allowed, but the
  *panel* renders in a terminal — whatever capture method is used must show the real
  terminal output, not a mock.
