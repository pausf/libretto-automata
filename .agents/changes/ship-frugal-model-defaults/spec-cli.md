# ship-frugal-model-defaults — the recommendation, and where it is read

Targets: cli

The repository knows which of its agents suit a cheap tier. That knowledge lives in one
paragraph of `skills/review-project/SKILL.md` and on nobody's machine, so every install
configures models by hand or not at all. This ships it as data.

**Recommend only. Nothing is ever written on the strength of a recommendation** — not by a
flag, not by an `--apply`, not by a first run. The reading is the repository's; the typing
is the user's. Same split `AGENTS.md` draws for the release bump, for the same reason: a
tool that acts on its own reading is a tool whose reading nobody audits.

## Why this is `cli` and not `agent-models`

**Because `agent-models` settled, from a real failure, that it does not know this
repository.** Its prior decisions open with *the subject is a directory, not the
repository*, and its constraints say the package *never learns what an install target
is — it is handed a directory and works on every `*.md` in it, whoever created them.*
That decision was written after reading "the agents available right now" as "the agents
this payload ships", which edited seven files installed nowhere while the user's own
twenty-two stayed invisible.

A hardcoded map from `review-lens-security` and `spec-writer` to a tier **is** the
payload's agent list. Putting it in `agent-models` would reintroduce exactly the coupling
that decision exists against, in the one package whose layering claim is that it is
testable against a bare `t.TempDir()`.

`cmd/libretto` already knows what the payload is — it gates commands on it. So the table
goes where that knowledge already lives, and `agent-models` is not touched at all.

## Outcomes

**1 · A per-agent recommendation, in code, in the binary that already knows the payload.**

`recommend(name string) (recommendation, bool)` returns a recommended model alias, a
recommended effort, and **the one-line reason** — never a bare verdict. `false` for an
agent the table does not know, which is every agent a user writes themselves.

**Unknown is silent, never a guess.** A recommendation invented for somebody's own agent
is worse than none.

**2 · `libretto models` prints the recommendation, with its reason.**

This listing has no width budget to defend, so it is where the reasoning goes. Each row
whose agent has a recommendation gains it; rows without one gain nothing.

**The listing is where the reason lives, and that is a decision the panel forces.** The
selector's catalogue is 58 columns at its narrowest and the reasons run to seventy runes,
so a reason rendered there would tear the frame — the same arithmetic that kept a third
column off the row. Mark on the screen, reason in the listing.

**3 · The adapter carries the recommendation to the panel, and nothing else changes.**

`agentRows` in `cmd/libretto/main.go` is already the seam where `internal/ui` and
`internal/agentmodel` are kept apart — it is what filled `Efforts` when the effort
narrowing landed. It gains one field the same way, so `internal/ui` still imports no
`agentmodel`.

**4 · The effort half is silent where effort does not exist.**

`Haiku 4.5` is listed with no effort levels, so an agent recommended onto `haiku` carries
no effort recommendation. Recommending a level that `SupportsEffort` reports unavailable
would be a suggestion the tool would itself refuse to type.

## Scope boundaries

**In:**

- `cmd/libretto/recommend.go` — the table, the lookup, the reason
- `cmd/libretto/recommend_test.go`
- `cmd/libretto/models.go` — the listing
- `cmd/libretto/models_test.go`
- `cmd/libretto/main.go` — one field in `agentRows`
- `.agents/specs/cli/spec.md` — the criteria, applied at phase 8

**Out, and named so it cannot be quietly added:**

- **No `--apply`, no `libretto models recommend`, no first-run write.** The whole point is
  the split. *Brings it back:* nothing; it would be a different change with its own
  argument.
- **`internal/agentmodel/` is not touched.** Reasoned above. *Brings it back:* a second
  consumer of the table outside this binary, which would mean the payload's agent list had
  become somebody else's business too.
- **No editing `agents/*.md`.** The five lenses and two workers keep the `model:` keys they
  declare. A change that shipped a recommendation *and* rewrote the files to match would
  have applied it under another name.
- **No cost or price figures in a reason.** Prices move under the repository, and the
  catalogue already refuses to guess what a user's plan includes.
- **No new dependency, no new package.**

## Constraints

- **Every recommended value passes `agentmodel.Valid` / `ValidEffort`**, and a recommended
  effort only appears where `SupportsEffort` is true for the recommended model. A table
  entry naming `sonnet-5`, or putting `xhigh` on `haiku`, is a recommendation the tool
  would refuse to type — the one thing worse than no recommendation. Proved by walking the
  table, not trusted.
- **Matched whole, by agent name.** No prefix, no substring.
- **`TestModelsListsEffortBesideTheModel` asserts `(session)` appears exactly twice** on the
  inheriting agent's line. It is a counting assertion and it is the tripwire for a
  recommendation rendered into the wrong column. It stays, unedited.
- All six gates pass before any commit.

## Prior decisions

**Settled by reading, this change:**

- **Not in the agent files' frontmatter, and not in a markdown table.** A
  `recommended-model:` key would be a second spelling of a decision the file already
  carries, and `agent-models` records that writing model knowledge into a markdown table
  **was considered and rejected once already** — the selector would write one format and
  parse another that any reformat breaks.
- **A recommendation carries a reason.** The catalogue's model labels answer "when would I
  pick this"; a per-agent verdict with no reason is an instruction, and nothing here gives
  instructions.

**Assumed under `/libretto-attacca`, because nobody was there to answer.** Each names what
changes if it is wrong:

- **A1 · The recommendations themselves**, each derived from what the agent's own
  description says it does:

  | Agent | Recommended | Reason |
  |---|---|---|
  | `review-lens-design` | `haiku` | pattern-matching over prose |
  | `review-lens-tests` | `haiku` | pattern-matching over prose |
  | `review-lens-intent` | `sonnet` · `high` | the only lens that runs commands, and it reads forge payload |
  | `review-lens-reliability` | `sonnet` · `high` | races and error paths are reasoning, not matching |
  | `review-lens-security` | `sonnet` · `xhigh` | the skill says security does not run cheap |
  | `spec-writer` | `sonnet` · `high` | it writes the contract everything later is measured against |
  | `work-reviewer` | `sonnet` · `high` | a cheap miss here is a false green |

  *If wrong:* one table entry each. Nothing structural moves, which is why the values sit
  in one place.

  **The values are reviewed, never tested, and review caught one.** The guards below walk
  the table for what a machine can decide — is this a real model, is that a level the
  model has, is this the priciest tier — and none of them can see a value disagreeing with
  this paragraph. `review-lens-intent` shipped carrying `high` while this table omitted it,
  and both were green. The table was right and the paragraph was short; **a recommendation
  of `high` is not noise**, because an absent key means the *session's* effort and the
  session may have been launched at `low`.

- **A2 · Nothing is recommended onto `opus` or `fable`.** The catalogue orders itself
  cheapest-first precisely because this feature exists to reduce a bill. `sonnet` is the
  ceiling the table uses; an agent that genuinely needs more is a decision a person makes
  at the screen. *If wrong:* a table entry changes.

- **A3 · `review-lens-design` is recommended onto `haiku` while its file declares
  `sonnet`.** A real disagreement on disk, not an oversight: the skill's own reasoning says
  design and tests are pattern-matching over prose and can run cheap, and
  `review-lens-tests.md` declares `haiku` while `review-lens-design.md` declares `sonnet`.
  **The recommendation follows the stated reasoning and the file is left alone**, so the
  disagreement becomes visible rather than settled silently here. *If wrong:* the entry
  becomes `sonnet` — but somebody should say which of the two was the mistake.

## Task breakdown

1. `recommendation` and `recommend`, with the table, test-first.
2. The two integrity guards proved by walking the table.
3. An unknown agent: `false`, no reason.
4. The listing column, and the existing counting assertions still green.
5. One field in `agentRows`.
6. Six gates, then apply this delta onto `.agents/specs/cli/spec.md`.

## Verification criteria

- **every recommendation names a model the catalogue accepts**, and an effort the
  recommended model supports. The table is walked, so a future edit naming `sonnet-5` or
  putting a level on `haiku` fails here rather than at the screen.
  Proof: cmd/libretto/recommend_test.go TestEveryRecommendationIsTypeable

- **an agent the table does not know gets no recommendation and no reason.**
  Proof: cmd/libretto/recommend_test.go TestAnUnknownAgentGetsNoRecommendation

- **every recommendation carries a reason.** A verdict with no reason is an instruction.
  Proof: cmd/libretto/recommend_test.go TestEveryRecommendationCarriesAReason

- **nothing is recommended onto the two priciest tiers.** The catalogue is ordered
  cheapest-first because this feature exists to lower a bill; a table recommending `opus`
  or `fable` inverts the thing it was built for.
  Proof: cmd/libretto/recommend_test.go TestNothingIsRecommendedOntoThePriciestTiers

- **the listing prints the recommendation and its reason, and prints nothing for an agent
  without one.** No blank cell standing in for an opinion that does not exist.
  Proof: cmd/libretto/models_test.go TestModelsListsTheRecommendationAndItsReason

- **the listing says when an agent's declared model differs from what is recommended**,
  because a recommendation nobody can tell they are ignoring is a recommendation that
  changes nothing. `review-lens-design` is that case on the day this ships.
  Proof: cmd/libretto/models_test.go TestModelsMarksAgentsRunningAgainstTheRecommendation

- **the panel gets the recommendation through the existing adapter.** `agentRows` fills it
  exactly as it fills `Efforts`, so `internal/ui` still imports no `agentmodel`.
  Proof: cmd/libretto/models_test.go TestAgentRowsCarryTheRecommendation
