# Decisions

### Session 2026-08-21

- Q: approach — self-contained init.go, an internal/bootstrap package, or init with
  guided onboarding content? → A: self-contained `cmd/libretto/init.go`, one
  `rev-parse` through the existing gitRunner seam, stdlib only — the mould every
  command follows. internal/ package lost (one consumer); onboarding content lost
  (ponytail — the list grows one file at a time when someone asks). If wrong: the
  logic lifts the day a second consumer exists. (assumed)

- Q: what exactly does `init` create? → A: the `.agents/` layout only, minimally:
  `.agents/specs/.gitkeep`, `.agents/changes/.gitkeep`, and `.agents/lessons.md` with
  the standard append-only header. No sample spec, no `docs/SPEC.md`, no CLAUDE.md
  (out of scope by standing decision). If wrong: the list grows one file at a time.
  (assumed)
- Q: what if the project already consolidates elsewhere (`specs/`, `openspec/`,
  `docs/specs/`, `spec/`)? → A: refuse and name what was found — the project already
  has a convention and `init` bootstraps only the `.agents/` layout; inventing a
  second layout is two sources of truth. A partial `.agents/` gets its missing pieces
  filled. If wrong: a flag to adopt another layout is a later change. (assumed)
- Q: idempotent? → A: yes — a second run creates nothing, overwrites nothing, reports
  everything as already present, exits 0. Never overwrite anything, ever — the tool's
  cardinal rule. If wrong: nothing; this one has no other defensible reading. (assumed)
- Q: outside a git repository? → A: refuse, same shape as `metrics`/`land` — the
  whole point of the keep-files is surviving git, and a non-repo has nothing to
  survive. If wrong: a warn-and-continue mode is one branch. (assumed)
- Q: does `init` commit? → A: never — it creates files and reports; committing the
  user's repository is the user's act. If wrong: unthinkable here; committing would
  break "never commit unasked". (assumed)
- Q: does this repository dogfood the fix? → A: yes — `.agents/changes/.gitkeep` is
  created and committed as part of this change, closing the vanishing-queue hole that
  lost ask-release-label-at-attacca-end. `.agents/specs/` already survives (it has
  content). If wrong: one file to delete. (assumed)
- Q: (writer finding) foreign root AND a partial `.agents/` both present — fill or
  refuse? → A: refusal wins, derived from the foreign-layout decision's own
  two-sources-of-truth rationale; the report names both findings. If wrong: one
  branch reorders. (assumed)
- Q: what happens to the payload non-goal ("no creating a changes/ directory in a
  project that has none")? → A: it is rewritten, not deleted: implicit creation by
  skills stays forbidden; `libretto init` becomes the named, explicit, user-invoked
  owner of bootstrap. The two refusal sites (write-spec, libretto-queue) gain a
  guarded one-line pointer ("where `libretto` is on PATH: `libretto init`"), absent
  path unchanged. If wrong: the pointers are one sentence each to remove. (assumed)
