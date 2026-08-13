# ship-frugal-model-defaults — the recommendation itself

Targets: agent-models

The repository knows which agents suit a cheap tier. That knowledge lives in one
paragraph of `skills/review-project/SKILL.md` and in nobody's machine, so every install
configures models by hand or not at all. This ships it as data.

**Recommend only. Nothing is ever written on the strength of a recommendation** — not by a
flag, not by an `--apply`, not by a first run. The reading is the repository's; the typing
is the user's. That is the same split `AGENTS.md` draws for the release bump, and it is
here for the same reason: a tool that acts on its own reading is a tool whose reading
nobody audits.

## Outcomes

**1 · A per-agent recommendation, in code, beside the catalogue it draws from.**

`Recommend(name string) (Recommendation, bool)` returns a recommended model alias, a
recommended effort, and **the one-line reason** — never a bare verdict. The catalogue's
own entries already carry "when would I pick this" labels; this is the same idea keyed by
agent instead of by model, and it is written in the same voice.

`false` for an agent the table does not know, which is every agent a user writes
themselves. **Unknown is silent, never a guess** — the capability already treats what it
cannot verify as capable rather than inventing an answer, and a recommendation invented
for somebody's own agent is worse than none.

**2 · A recommendation is a suggestion about a file, and touches no file.**

Nothing in this capability's writing path consults the table. `Apply`, `SetModel`,
`SetEffort` and `ApplyEffort` behave exactly as they do today, byte for byte.

**3 · The effort half is silent where effort does not exist.**

`Haiku 4.5` is listed with no effort levels, so an agent recommended onto `haiku` carries
no effort recommendation. Recommending a level that `SupportsEffort` reports as
unavailable would be a suggestion the tool itself would refuse to type.

## Scope boundaries

**In:**

- `internal/agentmodel/recommend.go` — the table, the lookup, the reason
- `internal/agentmodel/recommend_test.go`
- `.agents/specs/agent-models/spec.md` — the criteria and a widened `In:` line

**Out, and named so it cannot be quietly added:**

- **No `--apply`, no `libretto models recommend`, no first-run write.** The whole point is
  the split. *Brings it back:* nothing. If it ever arrives it is a different change with
  its own argument.
- **No editing `agents/*.md` in this change.** The five review lenses and two workers keep
  the `model:` keys they declare. A change that shipped a recommendation *and* rewrote the
  files to match it would have applied it under another name.
- **No per-agent reason for an agent the table does not know.** Silence.
- **No new dependency, no new package.** One file beside `catalogue.go`.
- **No cost or price figures in the reason.** Prices move under the repository and the
  catalogue already refuses to guess what a user's plan includes.

## Constraints

- **Every recommended value passes `Valid` / `ValidEffort`.** A table entry naming a model
  the catalogue does not know is a recommendation the tool would refuse to type, which is
  the one thing worse than no recommendation. Proved, not trusted.
- **A recommended effort only appears where `SupportsEffort` is true** for the recommended
  model.
- **The table is keyed by agent name, matched whole.** No prefix, no substring — the same
  discipline the branch attribution in `metrics` learned.
- All six gates pass before any commit.

## Prior decisions

**Settled by reading, this change:**

- **The table lives in `internal/agentmodel/`, not in the agent files' frontmatter and not
  in a markdown table.** A `recommended-model:` key in the frontmatter would be a second
  spelling of a decision the file already carries, and the spec records that writing model
  knowledge into a markdown table **was considered and rejected once already**, because
  the selector would then write one format and parse another that any reformat breaks.
- **A recommendation carries a reason.** The catalogue's labels answer "when would I pick
  this"; a per-agent verdict with no reason is an instruction, and this capability does
  not give instructions.

**Assumed under `/libretto-attacca`, because nobody was there to answer.** Each names what
changes if it is wrong:

- **A1 · The recommendations themselves.** Every one is a judgment call, derived from what
  each agent's own description says it does:

  | Agent | Recommended | Reason, in one line |
  |---|---|---|
  | `review-lens-design` | `haiku` | pattern-matching over prose — over-engineering and YAGNI read off the diff's shape |
  | `review-lens-tests` | `haiku` | pattern-matching over prose — does a test assert behaviour, was one weakened |
  | `review-lens-intent` | `sonnet` | the only lens that runs commands, and it reads the forge's payload |
  | `review-lens-reliability` | `sonnet`, `high` | runtime behaviour, races and error paths are reasoning, not matching |
  | `review-lens-security` | `sonnet`, `xhigh` | the skill says in as many words that security does not run cheap |
  | `spec-writer` | `sonnet`, `high` | it writes a contract the whole flow is later measured against |
  | `work-reviewer` | `sonnet`, `high` | it re-runs proofs and reports independently; a cheap miss here is a false green |

  *If wrong:* one table entry each, and the reason beside it. Nothing structural moves,
  which is the point of keeping the values in one place.

- **A2 · No agent is recommended onto `opus` or `fable`.** The catalogue orders itself
  cheapest-first precisely because this feature exists to reduce a bill, and recommending
  the two priciest tiers would invert that. `sonnet` is the ceiling the table uses, and an
  agent that genuinely needs more is a decision a person makes at the screen. *If wrong:*
  a table entry changes; `Valid` already accepts both.

- **A3 · `review-lens-design` is recommended onto `haiku` while its file declares
  `sonnet`.** This is a real disagreement on disk, not an oversight: the skill's own
  reasoning says "design and tests are pattern-matching over prose and can run on a cheap
  model", and `review-lens-tests.md` declares `haiku` while `review-lens-design.md`
  declares `sonnet`. **The recommendation follows the stated reasoning and the file is
  left alone**, so the disagreement becomes visible on the screen rather than settled
  silently by this change. *If wrong:* the table entry becomes `sonnet` and the divergence
  disappears — but somebody should say which of the two was the mistake.

## Task breakdown

1. `Recommendation` and `Recommend`, with the table, test-first.
2. The two guards proved: every entry valid, and no effort recommended where the model
   supports none.
3. An unknown agent returns `false` and no reason.
4. Six gates, then apply this delta onto `.agents/specs/agent-models/spec.md`.

## Verification criteria

- **every recommendation names a model the catalogue accepts**, and an effort the
  recommended model actually supports. The table is walked and each entry checked, so a
  future edit naming `sonnet-5` or putting `xhigh` on `haiku` fails here rather than at
  the screen.
  Proof: internal/agentmodel/recommend_test.go TestEveryRecommendationIsTypeable

- **an agent the table does not know gets no recommendation and no reason.** Silence, not
  a guess — a user's own agent is not something this repository has an opinion about.
  Proof: internal/agentmodel/recommend_test.go TestAnUnknownAgentGetsNoRecommendation

- **a recommendation onto a model with no effort levels carries no effort.** `haiku` has
  none, so recommending one would be a suggestion the tool would refuse to type.
  Proof: internal/agentmodel/recommend_test.go TestNoEffortIsRecommendedWhereTheModelHasNone

- **every recommendation carries a reason.** A verdict with no reason is an instruction,
  and this capability does not give instructions.
  Proof: internal/agentmodel/recommend_test.go TestEveryRecommendationCarriesAReason

- **the writing path never consults the table.** `Apply`, `SetModel`, `SetEffort` and
  `ApplyEffort` write exactly what they are given, for an agent with a recommendation and
  against it.
  Proof: internal/agentmodel/recommend_test.go TestWritingIgnoresTheRecommendation
