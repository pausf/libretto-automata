# fix-metrics-span-overlap

Tracker: none
Queued: 2026-08-13

## The ask, verbatim

> 1. El total de span suma solapamientos (metrics.go:290). totalSpan += m.span() suma el
> span de cada change como si fueran secuenciales. Si dos changes estuvieron abiertos la
> misma semana, el footer reporta dos semanas de "wall clock" que en el calendario fueron
> una. El total dice "X de reloj" y no es reloj de nadie — es la suma de rangos que se
> pisan. Arreglo barato: o fusionar los rangos antes de sumar, o renombrar el footer a
> algo honesto tipo "span acumulado por change".

## Reading

The only one of the metrics findings that is a wrong number rather than missing
ergonomics. Two fixes on the table: merge overlapping [first, last] ranges before
summing, or relabel the footer so it stops claiming to be wall clock. Per-change spans
in the rows are correct either way; only the total lies.
