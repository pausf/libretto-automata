# keep-the-readme-in-step

Tracker: none

## The ask, verbatim

> has añadido un comando nuevo de attacca pero tre has olvidado de actualizar el readme,
> eso no puede pasar... modifica el readme y asegurte que no vuelva a pasar , no subas de
> version por este cambio

## Reading

Two halves, and the second is the point.

`/libretto-attacca` landed in PR #33 with the README untouched — it still lists five
commands. The first half is the fix: put the command where a reader finds the others.

The second half is *"que no vuelva a pasar"*. A rule written in prose asking people to
remember is what already failed here — `readme/spec.md` exists and did not stop it. So the
answer has to be mechanical: `scripts/check-payload` already walks `commands/`, and a
command file with no mention in the README is exactly the kind of thing it can fail on.

**No version bump for this change** — the user's explicit instruction. `release:patch` at
most, and the request that carries it says so.
