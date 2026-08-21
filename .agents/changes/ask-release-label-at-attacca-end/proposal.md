# ask-release-label-at-attacca-end

Tracker: none
Queued: 2026-08-21

## The ask, verbatim

From `.agents/lessons.md` (2026-08-13 · improve-metrics-report · phase 8), where this
idea was queued and then lost when the changes directory emptied:

> The correction is about presentation: the attacca run reported the flow complete and
> the PR open without saying, loudly and in its own line, **"the label check will now be
> red on the PR, that red is the bump question reaching you, here is the one command"**.
> A designed refusal that the report does not predict reads as a broken pipeline, costs
> the user an alarm and a round trip, and the mode's whole visibility budget is that
> report.

(Re-captured in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Payload-side fix: the attacca closing report must predict the red `require-release-label`
check and hand the user the one labelling command.
