# add-specs-wiki

Tracker: none

## The ask, verbatim

> nueva feature, se me ha ocurrido crear una wiki por proyecto de las specs que se van
> creando es, decir sdaber que ahora tenemos specs por dominicio de la aplicacion, pueda
> que depende de donde este en que proyecto etc.. se cree un wiki de eso, no se muy bien
> como hacerlo desde el libretto si a la hora de instalar hacer un mini script de node
> etc,... busca informacion de esto y dime opciones que podemos hacer, como todos son
> symlink no se me ocurre

## Reading

Not a bug — a new feature. The flow writes specs per capability/domain into
`.agents/specs/<capability>/spec.md` in each project that uses the payload. The user
wants a per-project **wiki view** of those specs: something browsable that shows what
capabilities exist and what they promise, generated from the spec files themselves.

Open questions the user explicitly raised: how to deliver it given that libretto
installs everything as symlinks — a node mini-script at install time? something else?
The user asked for research and options first, not an implementation.
