# adopt-organic-driven-development

Tracker: none

## The ask, verbatim

> pasa esto que es SDD a ODD, busca informacion sobre el ODD creo que gentleman
> programming habla de esto

> no pero quiero hacer el desarrollo en libretto como si no me copiara de gentle , es
> decir sabes qu enuestro proyecto utiliza sdd le si quiere el proyecto, pues tenemos
> que pasarlo a ODD

## Reading

Not a bug: nothing built behaves wrong. This is a change of intent for the payload itself.

Libretto's flow is spec-driven: every piece of work walks the same eight phases and
produces proposal, spec, plan and tasks before code. The ask is to move it to an
*organic* model — structure proportional to the work, so a small change stays light and
only substantial work earns durable artifacts — designed as Libretto's own, not a copy
of Gentleman Programming's ODD in gentle-ai.

Research that informed the reading (2026-10-21): gentle-ai v4 made Organic Driven
Development its only workflow and retired SDD/OpenSpec. Its core ideas — one feature
document instead of a chain of artifacts, delegation only for a named reason,
verification scaled by risk — are the reference point, not the template.

## Open, for phase 2

- what of the current flow survives: capability specs with EARS + `Proof:`, `--anchors`,
  the trivial lane, landing into `.agents/specs/`
- what decides "small" versus "substantial", and who decides it
- whether the eight phases collapse, or stay as steps a small change may skip
- what happens to the CLI surfaces that read the current layout: `loop`, `metrics`,
  `wiki`, `land`, `find-work`'s scans
