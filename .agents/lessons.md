# Lessons

Append-only. Entries are spent by `libretto-retro` and gain a `Resolved:` line there;
they are never edited otherwise.

## 2026-08-13 · improve-metrics-report · phase 8

> no pasa las pipeline https://github.com/pausf/libretto-automata/pull/44 algo ha
> pasado , cuidado y apuntalo para libretto retro

Nothing was broken: `gates` passed twice and the red check was `require-release-label`
— the designed refusal that exists because attacca must never label a release. The
correction is about presentation: the attacca run reported the flow complete and the PR
open without saying, loudly and in its own line, **"the label check will now be red on
the PR, that red is the bump question reaching you, here is the one command"**. A
designed refusal that the report does not predict reads as a broken pipeline, costs the
user an alarm and a round trip, and the mode's whole visibility budget is that report.
Resolved: 2026-08-13 · project knowledge — convention recorded in .agents/specs/ci/spec.md
(prior decisions); the payload-side fix stays queued as ask-release-label-at-attacca-end
