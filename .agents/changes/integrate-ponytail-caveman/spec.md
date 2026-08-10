# Delta: ponytail and caveman ship in the payload

Targets: payload

The flow's prose defers to ponytail and caveman at almost every phase — "if `ponytail`
is installed…" — but neither ships, so a user starting from zero installs the flow and
gets conditionals that are never true, about tools nobody told them exist. The same
gap the superpowers delegates closed in this same file's history: a thin skill whose
delegate is missing is not thin, it is broken. The conditional prose is the same
brokenness at lower volume.

Both upstream cores are standalone `SKILL.md` files — verified 2026-08-10: no hooks,
no plugin machinery, MIT both. They can be vendored exactly like the three
superpowers skills already are.

## Outcomes

- **Four new vendored skills in `skills/`, copied unmodified, upstream names kept:**

  | Skill | Origin | Why this one |
  |---|---|---|
  | `ponytail` | DietrichGebert/ponytail | the ladder — `write-spec`, `build-and-check` and `libretto-flow` already call it |
  | `ponytail-debt` | DietrichGebert/ponytail | `write-spec`'s prior-decisions pillar reads its ledger |
  | `caveman` | JuliusBrussee/caveman | `present-work` already defers to it |
  | `caveman-commit` | JuliusBrussee/caveman | `record-work` already offers it for messages |

  The selection rule is ponytail's own first rung applied to vendoring: **only what
  the flow already calls by name.**

- `libretto install` links them like any other skill — no linker change, they are
  ordinary `skills/` directories.
- **THIRD-PARTY.md moves ponytail and caveman from *Not vendored* to *Vendored***,
  with pinned version/commit, licence line, and full licence texts beside
  `LICENSE-superpowers` (`LICENSE-ponytail`, `LICENSE-caveman`).
- **The calling skills' prose stops being conditional where the condition is now
  always true post-install**, and plugin-namespaced references (`ponytail:ponytail`)
  become the vendored plain names. Per the vendoring constraint, every adaptation
  lands in the calling skill, never in the copy.
- Docs that state the old decision say the new one: `docs/FLOW.md`'s "present rather
  than vendored" paragraph, README's companions section.

## Scope boundaries

**In:** the four skills above, their attribution, the calling-skill prose, the two
docs, THIRD-PARTY.md.

**Out, named:**

- **the rest of both plugins** — `ponytail-review`, `ponytail-audit`, `ponytail-gain`,
  `ponytail-help`, `caveman-compress`, `caveman-review`, `caveman-stats`,
  `caveman-help`, `cavecrew`. The flow calls none of them, and `ponytail-review`
  overlaps the review seam this payload already owns. Back in when the flow grows a
  call site, not before.
- **always-on mode.** As plugins, both can inject themselves into every session via
  hooks; libretto does not manage `settings.json` or hooks (STATE.md, out of scope).
  The vendored skills activate when the flow invokes them and when their descriptions
  trigger — which is precisely where this payload wants them. The upstream plugin
  remains the path to always-on, and both copies coexist by namespace.
- **requiring them.** Nothing fails without them; `prune`/`uninstall` removes them
  like any other item. Shipped is not required.
- **any CLI change.** `companion()` already counts `~/.claude/skills/` as a
  legitimate home, so `doctor` reports the vendored copies as present the moment
  `install` runs. Zero Go moves.

## Constraints

- **Copied unmodified** — the payload spec's standing vendoring rule. Adaptations go
  in the calling skill so the copy stays comparable with upstream.
- Frontmatter `name:` must equal the directory — upstream already satisfies this;
  verify rather than assume, `check-payload` is the gate.
- Vendored skills must be self-sufficient once installed: no references back to their
  upstream repo's `scripts/` or docs. Verified for the two cores; verify for
  `ponytail-debt` and `caveman-commit` at copy time, and if one reaches outside
  itself, its tool ships inside its directory or the skill is dropped from the set —
  never patched.
- Pinned at the upstream version/commit current at vendor time, recorded in
  THIRD-PARTY.md. Drift from upstream is the accepted cost, same as superpowers at
  6.2.0.

## Prior decisions

- **Vendor, reversing THIRD-PARTY.md's *Not vendored* entry — the user's explicit
  call, 2026-08-10.** The old rationale (a second copy of something the user may
  already have chosen a version of) assumed a user who has versions of things; the
  target is now a machine with nothing on it. The collision half of that rationale
  was already answered by the *Naming* section: plugins namespace, vendored copies
  do not, both coexist, the plugin's is likely newer.
- "Installer prints the install commands" was specced first and discarded in review
  of this same change — printed commands still leave the fresh user one manual step
  from a flow that works as written.
- STATE.md's *never required* stands untouched. This change makes them present by
  default, not mandatory.

## Task breakdown

- [ ] 1. vendor `ponytail` and `ponytail-debt`, pinned, with `LICENSE-ponytail`
- [ ] 2. vendor `caveman` and `caveman-commit`, pinned, with `LICENSE-caveman`
- [ ] 3. THIRD-PARTY.md: move both to *Vendored*, record versions, update *Not
       vendored* and *Naming*
- [ ] 4. calling skills: plain names for `ponytail:ponytail`, conditional prose
       reconciled in `write-spec`, `build-and-check`, `present-work`, `record-work`,
       `libretto-flow`
- [ ] 5. docs: FLOW.md's vendoring paragraph and README's companions section state
       the new decision
- [ ] 6. all six gates green over the enlarged payload

## Verification criteria

- every vendored skill's frontmatter parses and `name:` equals its directory
  Proof: scripts/check-payload
- every skill the flow references — including the four new names — exists in the
  payload
  Proof: scripts/check-payload
- no vendored skill invokes a path that does not get installed
  Proof: scripts/check-payload
- the payload spec's own anchors still resolve after the prose edits
  Proof: skills/record-work/spec-drift --anchors

Unprovable by a gate, so stated as a recorded fact instead: the copies match upstream
at the pinned commit — THIRD-PARTY.md records which commit, and `diff` against a fresh
clone is the check when doubt arises.
