# refine-model-selector

Tracker: none

## The ask, verbatim

> en el cli en el apartado de modelo
>
> ```
>  │  ❯ [ ] review-lens-design       haiku          shared                                            │
>  │    [ ] review-lens-intent       haiku          shared                                            │
>  │    [ ] review-lens-reliability  (session)      shared                                            │
>  │    [ ] review-lens-security     (session)      shared                                            │
>  │    [ ] review-lens-tests        haiku          shared                                            │
>  │    [ ] spec-writer              (session)      shared                                            │
>  │    [ ] work-reviewer            (session)      shared                                            │
> ```
>
> molaria que se separara por modelo del agente es decir una linea que los separe
> para que se vea mas visual y en la leyenda poner que el space es para selecionar
> y faltaria arreglar un apartado que ponga all que si lo seleccionas, se
> selecionand todos de golpe

## Reading

Three things about the models selector screen (`libretto models`, and the selector
screen inside the panel):

1. **Group the rows by the agent's current model**, with a separator line between
   groups, so the screen reads visually instead of as a flat list.
2. **The legend does not say that `space` selects.** Add it.
3. **There is no "all" row.** Selecting it should select every agent at once.
