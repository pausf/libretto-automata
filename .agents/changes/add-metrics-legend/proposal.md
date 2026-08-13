# add-metrics-legend

Tracker: none

## The ask, verbatim

> yo lo que hecho de menos es una leyenda de lo que significa cada cosa sabes?

## Reading

`libretto metrics` prints a footer explaining what it does NOT measure (`flowCeiling`)
but nothing explaining what each column means. `reopen`, `—` and the landed/in-flight
distinction all need a translator today. Add a legend in the same style as the existing
footer, in English, since it is CLI output. Sketch discussed:

```
  commits  how many commits touched the change folder
  span     first commit to last — calendar clock, not attention
  closed   plan.md boxes that went from open to closed
  reopen   boxes that went back open — tasks called done before they were
  state    landed = folder deleted on landing · in flight = still on disk
  —        the change has no plan.md in its history
```
