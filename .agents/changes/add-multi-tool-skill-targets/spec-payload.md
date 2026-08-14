# Delta: the payload reads correctly in all three tools

Targets: payload

## Outcomes

- No skill or command instructs by naming Claude where it means *the agent
  running this skill*. "Ask Claude to…" becomes an instruction any of the
  three tools can follow about itself.
- Where a skill depends on a Claude-Code-only mechanism (the
  `AskUserQuestion` tool, the `Skill` tool, `CLAUDE_HOME`), the dependency
  stays — but phrased with its generic fallback: ask the user in
  conversation when the native prompt does not exist.
- Literal `Claude` survives only where it is a fact, not an addressee: real
  paths (`~/.claude`, `CLAUDE.md`, `CLAUDE_HOME`), the product name when the
  sentence is about Claude Code specifically, and quotes.
- `scripts/check-payload` gains a prose check: every occurrence of `Claude`
  in skill and command bodies that does not match the allowlist of factual
  patterns (paths, `CLAUDE.md`, `CLAUDE_HOME`, product statements) fails the
  gate. The classification lives entirely in the allowlist — the check
  itself exercises no judgment — so the property survives the next skill
  written and two implementers of the check cannot disagree.

## Scope boundaries

**In:** the wording pass over `skills/**` and `commands/**` prose, the
check-payload rule, its allowlist.

**Out, named:**
- Per-tool conditional content ("if you are Codex, do…"). One text that
  reads correctly everywhere, never three texts.
- Rewriting `agents/*.md`. Agent definitions are Claude-Code-only until the
  `add-transformed-agent-targets` change ships them elsewhere; renaming
  their addressee now serves nobody.
- Feature parity across tools. A skill that uses `AskUserQuestion` degrades
  to asking in conversation; making every tool equally capable is not
  wording work.

## Constraints

- The check lives in `scripts/check-payload` — repo-only tooling, never
  referenced from a skill, exactly as every other payload gate.
- `ponytail:` comments stay English and stay named `ponytail:` — the marker
  is a contract with `ponytail-debt`, not prose.
- Skill frontmatter `version:` bumps where a skill's instructions change
  meaningfully; a pure pronoun swap does not force one, a changed fallback
  instruction does.

## Prior decisions

- **The wording pass ships inside this change** — user decision, 2026-08-14,
  overriding the recommendation to queue it separately. Consequence accepted:
  one review covers CLI mechanics and payload prose.
- The allowlist is a list in the check, not a lint pragma in the files —
  payload prose stays clean of tooling markers.

## Task breakdown

1. Audit: list every `Claude` occurrence in `skills/**` and `commands/**`
   prose, classify addressee vs fact.
2. Rewrite the addressee cases; name generic fallbacks where a Claude-only
   mechanism is invoked.
3. The check-payload prose rule with its allowlist.

## Verification criteria

- no skill or command addresses Claude where it means the running agent, and
  the check fails on a reintroduction
  Proof: scripts/check-payload
- factual uses (paths, env vars, product statements) pass the check
  unflagged
  Proof: scripts/check-payload
