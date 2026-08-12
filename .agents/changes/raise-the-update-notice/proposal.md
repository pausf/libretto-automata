# raise-the-update-notice

Tracker: none

## The ask, in the words it was asked in

> no se ve dorado ni centra ni nada, yo pensaba que molaria mas como poner esto
> [una imagen de una etiqueta azul con NEW VERSION] encima del logo y dentro poner
> update v0.6.0 → v1.0.2 available

> a,os hacer lo que comente de la imagen este y lo poner sobre el logo asi lo vemos
> mejor lo de actualizar

Invoked under `/libretto-attacca`, so the contract stop, the order stop and the push
were all answered at the prompt.

## What prompted it

The user could not see the notice at all, and three separate reasons stacked up:

- `preview` never wires `WithReleaseCheck`, so the panel's gold row is silent there by
  construction
- the `libretto` on their PATH reports `v0.7.0-35-g13d4766`, which `parseSemver` rejects,
  so `IsNewer` is false and no notice is ever produced
- what they *did* see was `subcommandNotice` — plain and uncoloured on purpose, appended
  after the frame

So the row was working the whole time. The complaint that survives once that is cleared
up is a real one about weight: the notice sits below the menu, in the same visual register
as a target row, and it reads as one more line rather than as news.

## What this change is

Move the update notice from its band between the menu and the destination strip to a
banner above the logo, drawn as a bordered label rather than a bare gold line.

## What it is not

The `v1.0.2` the notice currently offers is the retraction tombstone, and offering it to
a checkout is a real defect — `git ls-remote` cannot see a `retract` block the way the
module proxy can. That is separate work and does not land here.
