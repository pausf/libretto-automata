# manage-every-agent-in-the-target

Tracker: none

## The ask, in the words it was asked in

From the original request that produced `agent-models`:

> podriamos mostrar en el cli y en el comando somos lo agente que tiene disponible en
> ese momento y que elija por agente que modelo quiere tirar

And, on finding the shipped feature did not do that:

> le estoy dnado al enter […] y no me esta cambiando el modelo a los agentes...

## What actually happened

The feature works. `agents/review-intent.md` was written from the panel and is modified
on disk — the write, the refresh and the redraw all ran.

**It operates on the wrong set.** Measured, not guessed:

| | |
|---|---|
| agents in the user's `~/.claude/agents` | **22** |
| of those, managed by libretto | **0** |
| agents the selector lists | **7** — only what this repo ships |

*"Los agentes que tiene disponible en ese momento"* means the 22 — `sdd-apply`,
`jd-judge-a`, `review-risk`, `agent-planner`, `workflow-orchestrator` and the rest. It
was read as "the payload's agents" and built over `agents/` in this repository. That
reading was wrong, and nothing in the contract caught it because the contract was
written from the same wrong reading.

Every row in the global listing says `· not linked here`, which is the feature
correctly reporting that it is editing files that reach nothing.

## The decision

**Operate on every agent in the selected target**, including real files libretto did
not create.

Asked explicitly, with the cost stated: `AGENTS.md` forbids *overwriting* anything the
tool did not create, and this crosses a related line — writing a frontmatter key inside
somebody else's file. The answer was to do it. Editing one key on request is a
different act from clobbering an item with a symlink, but it is a line to cross
deliberately and record, not to slide across.

## Second problem this exposed

```
conflict     agents/review-reliability.md
```

Splitting `review-lens` into four created `agents/review-reliability.md`, and the user
**already had an agent of that name**. `install` reports the conflict and never touches
it — correct — but the lens can therefore never link into the global target under that
name. The same collision is possible for `review-security`, `review-design` and
`review-tests`, which are also generic names.

## Settled in phase 2

- **Write foreign files without asking.** The scope answer chose the option that
  writes over the one that confirms each write. Recorded in `spec-cli.md`, with the
  distinction that makes it legal: `AGENTS.md` forbids *overwriting* what the tool did
  not create; this replaces one frontmatter line in a file the user named.
- **The package takes a directory, not a target.** `agentmodel` never learns what a
  target is, so the layering survives the wider reach. `spec.md`.
- **Ours versus theirs is shown as `shared`**, and it marks the blast radius rather
  than the ownership: a symlink into the repo reaches every project, a real file
  reaches one target. `spec-cli.md`, `spec-panel.md`.
- **The name collision is deferred to its own change.** It belongs to
  `review-project`, not to the model selector, and folding an install-time defect into
  a scope change would make both harder to review. **It is not fixed by this work** —
  `review-reliability` still cannot install into this machine's global target.

## Open, for phase 2 to settle

- **Whether a foreign file is written without asking.** The scope decision says write;
  it did not say silently.
- **What "ours" versus "theirs" means on screen**, if anything. An owned symlink and a
  real file behave identically to a frontmatter edit but not to `prune`.
- **The name collision.** Rename the lenses to something this payload owns, accept that
  they cannot be installed globally on this machine, or something else.
- **Whether `agent-models` still governs `internal/agentmodel/**` alone.** Reading a
  target's directory means the package learns what a target is, which the current spec
  forbids in as many words.
- **What happens to a target with no `agents/` directory**, and to a file in it that is
  not an agent at all.
