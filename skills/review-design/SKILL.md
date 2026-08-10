---
name: review-design
description: "Trigger: reviewing a diff, a file or an MR for design — over-engineering, YAGNI and KISS violations, SOLID breaches, code smells, the wrong abstraction. Every finding is a labelled judgment call, and the reviewed project's own conventions override the baseline. Standalone: works on any diff, with or without review-project."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Ask of every piece of the change: **does this earn its complexity?** Most design
debt enters a codebase through review, one plausible abstraction at a time, because
nobody asked the question while the diff was still cheap to change.

Two rules bind everything below:

- **The project overrides the baseline.** A convention the reviewed codebase
  documents or consistently practices wins over anything here. Read how the
  surrounding code does it before flagging how this diff does it — a finding that
  fights the house style is noise.
- **Every finding is a labelled judgment call, never a violation.** "Possible
  speculative generality" is a question the author can answer; a design rule
  enforced as law is how reviews become adversarial and get ignored. Skip anything
  a linter or formatter already enforces.

## The cheap questions first

In order, because the first ones delete the most code:

- **Does this need to exist?** (YAGNI) — an abstraction with one implementation, a
  parameter no caller passes, a config for a value that never changes, a hook for a
  future the MR does not name. The fix is deletion, and the day a second case
  arrives it can be built with the second case in view.
- **Is there a simpler way?** (KISS) — cleverness where boring would do, a regex
  where a comparison works, three layers where one function fits. The test: could
  the person paged at 3am follow this?
- **Is the duplication actually a problem?** — slight duplication beats the wrong
  abstraction. Two similar hunks are a candidate, not a finding; extract when the
  shape has proven shared, roughly the third occurrence, not the second.

## Then the structure

SOLID, as questions matched against the diff — flag the answer, name the principle:

| Question | Smell when the answer is no |
|---|---|
| Does this module change for one reason? | a class collecting unrelated jobs, edits to it scattered through the diff for different reasons |
| Can behaviour be added without editing this again? | the same type-switch recurring across files |
| Could any implementation substitute its interface? | a subtype that stubs or overrides most of what it inherits |
| Does the client use the whole interface it depends on? | empty method implementations, fat interfaces |
| Does the logic depend on an abstraction where one already exists? | business code constructing its own infrastructure |

And the classic smells worth naming when they appear in the change: a name that
does not say what the thing does, a method living off another object's data, the
same few parameters travelling in a pack, a primitive standing in for a domain
concept, a middle man that only delegates, a message chain the caller should not
know.

**One reason to change wins arguments that "SOLID" loses.** Cite what the code
does, not the acronym.

## Proportion

The baseline questions cut both ways: a finding that adds ceremony to a small
change is itself a KISS violation. A script run once does not need dependency
injection, and saying so is a legitimate answer to a finding from this lens.

## Output

Per finding: location (`file:line`), the label (which question or smell), what was
observed in the hunk, and the smaller version — what deletion or simplification
resolves it. Questions, not verdicts.

No findings is a statement: **"nothing here fights its own weight."**

This lens reports; it never blocks, edits, commits or pushes.
