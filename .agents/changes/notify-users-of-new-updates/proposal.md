# notify-users-of-new-updates

Tracker: none
Queued: 2026-08-11

## The ask, verbatim

> alguna veces me olvido de install desde libreto de nuevos command etc.. que vinieron
> en la nueva update es decir me falta fichero por linkar por que son nuevos, que
> podemos hacer para que lo usuario sepan que tiene una actualizacion nueva ?

## Reading

After pulling a newer payload, items that did not exist before are simply unlinked, and
nothing tells the user until they happen to run `install` again. The idea is to surface
that state — new or missing links, and possibly a newer release — instead of relying on
the user remembering to reinstall.
