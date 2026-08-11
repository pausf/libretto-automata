# add-change-queue

Tracker: none

## The ask, verbatim

> tengo otra feature , ya sabes que cuando vamos por flow se crean unos cambios luego
> spec etc.. pero si quiere ir creando ideas en changes y ejecutarla mas tarde como una
> pila no puedo, asi que lo guay seria crea un nuevo camando llamado por ejemplo
> libretto-queue cambiale el nombre si no lo ves, en el cual vamos pidiendo feature y se
> van poniendo en la cola de changes, luego con otra comando podemos ir 1 por 1
> haciendolas

## Reading

Today the flow takes one piece of work and runs it end to end. There is no way to
capture several ideas up front and work through them later. Wanted:

1. A command (working name `/libretto-queue`) that captures feature ideas one after
   another into `.agents/changes/` as queued entries — capture only, no spec, no plan,
   no code.
2. A way to later take the next queued idea and run it through the flow, one at a time.

Open questions for the spec: whether the second half is a new command or `/libretto-flow`
source 1 picking queued entries up; queue order (FIFO — the user said "pila" but
described a queue: "vamos pidiendo… luego 1 por 1"); how a queued idea differs from a
change in flight so `/libretto-status` and find-work report them honestly.
