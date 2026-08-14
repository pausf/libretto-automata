# consolidate-license-files

Tracker: none
Queued: 2026-08-14

## The ask, verbatim

> Hay múltiples archivos de licencia "alternativos" (LICENSE-caveman, -ponytail,
> -superpowers) que podrían confundir; dejar una única licencia principal y documentar
> alternativas sería más claro.

## Reading

Four files at the root: `LICENSE` plus three named after vendored payloads. They are not
"alternatives" to choose between — they are the upstream licences of third-party skills
this repository ships, and `THIRD-PARTY.md` already documents which is which. The
confusion the ask names is real; the fix is presentation, not consolidation.

**Deleting them is not on the table** — a vendored copy has to carry its own licence text,
and dropping one to tidy the root is a licensing problem rather than a docs problem. The
lazy move is a `licenses/` directory with the three, `LICENSE` alone at the root, and
`THIRD-PARTY.md` pointing at them. Check first that nothing resolves those paths: the
payload gate and `THIRD-PARTY.md` both reference vendored files.
