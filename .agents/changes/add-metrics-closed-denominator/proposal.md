# add-metrics-closed-denominator

Tracker: none

## The ask, verbatim

> 2. closed sin denominador. Ver "5 cerradas" no te dice si el plan tenía 5 cajas o 20.
> Un 5/5 versus 5/18 cambia completamente la lectura de un change in flight — uno está
> por aterrizar, el otro está empantanado. El dato ya está en el diff de plan.md que el
> código ya parsea; es una columna más, no una fuente nueva.

## Reading

Show `closed` as `n/total` so an in-flight change reads as nearly-landed or bogged down
at a glance. The total box count is derivable from the same plan.md history `measure`
already walks — no new data source.
