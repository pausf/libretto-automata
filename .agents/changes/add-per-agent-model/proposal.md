# add-per-agent-model

Tracker: none

## The ask, in the words it was asked in

> me gustaria hacer una nueva feature sea  Modelo más barato para design y tests. No son
> menos tokens, son tokens más baratos — pero si la pregunta real es la factura, acá está
> el 3-5× en esas dos lentes. Son pattern-matching sobre prosa., es decir que el usuario
> por ejemplo desde el CLI del libretto + un comando por ejemplo libretto-save-token
> (busca un nombre mejor que este) pueda cambiar el modelo que esta tirando los agentes ,
> por ejemplo podriamos mostrar en el cli y en el comando somos lo agente que tiene
> disponible en ese momento y que elija por agente que modelo quiere tirar, eso hara que
> se tenga que editar los ficheros de los agentes, importante tambien muestrale la
> opciones de modelo que tiene dispoinible con un subscription de claude

## Reading

Every agent in the payload runs on whatever model the session happens to be on. The four
lenses that are pattern-matching over prose — design and tests above all — pay Opus rates
for work a cheaper model does as well. The lever is not fewer lenses, it is cheaper tokens
per lens.

So: a per-agent model choice, surfaced where the user already is (the panel, and a
command), listing the models a Claude subscription actually gives them, and writing the
answer into each agent file's frontmatter as `model:`.

Today none of `agents/review-intent.md`, `agents/review-lens.md`, `agents/spec-writer.md`
or `agents/work-reviewer.md` declares a `model:`.

## Two doors, one selection model

The choice is offered in both places, and they are not two features:

- **the panel** (`libretto`, the TUI) — the agents listed with their current model, picked
  there
- **the command** — the same list, the same choice, for a session that is not in the panel

And the selection is **multi**: mark one agent, several, or all of them, then apply one
model to everything marked. One-by-one is the degenerate case of bulk, not a separate
mode — a flow that only does one at a time makes the common case (make the four prose
lenses cheap) into four separate acts.

## Open, for phase 2 to settle

- **The command's name.** `libretto-save-token` is wrong twice over — nothing here is a
  token, and nothing is saved. Candidates: `libretto models`, `/libretto-models`.
- **Where the choice is stored.** Editing the payload's own `agents/*.md` in the repo
  edits the source of the symlinks — every project that installed them changes at once,
  and the repo goes dirty.
- **Which model names to offer**, and whether the list is fixed or discovered.
- **`review-lens` is one agent for four lenses.** A per-lens model choice cannot live in
  its frontmatter; the caller would have to pass it.
- Whether `scripts/check-payload` accepts a `model:` key in frontmatter.
