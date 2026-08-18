# add-wiki-panel-entry

Tracker: none

## The ask, verbatim

> me gustari que añadiera la opcion aqui [panel menu screenshot: install /
> uninstall / update / status / doctor / prune, with the scope ▸ project strip]
> ojo solo tiene que salir instalar en projecto no en global hazla en la misma pr

## Reading

Not a bug — an amendment to the panel menu, delivered on the same branch/PR (#67)
as the wiki viewer, as asked. A `wiki` row in the panel, visible **only when the
scope is project** — never under global — because the wiki is a per-project
artifact read from the working directory's specs. Same conditional-row precedent
as `models` (absent when the repo has no agents): a project without a specs
directory gets no row either, since a row that can only error is a promise the
panel cannot keep.
