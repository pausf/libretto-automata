# add-wiki-html-output

Tracker: none

## The ask, verbatim

> ya pero no tiene un visor para ver la wiki? es decir como una pagina web? por
> ejemplo como hace claude con lo artifact
>
> no lo estas entendiendo, lo que quiero es que montes un sistema parecido para que
> se puede ejecutar desde cada projecto con sus specs

## Reading

Not a bug — a new feature, and it is the bring-back condition the add-specs-wiki
contract wrote down ("HTML comes back only if someone actually needs to browse specs
off-forge") arriving one conversation later. The user saw a demo artifact rendering
this repository's specs as a styled web page (sidebar index, per-capability sections,
criteria filter, light/dark) and wants that as a capability of the tool itself:
runnable in any project, against that project's own specs — not a one-off page.

Shape: an HTML output mode on the existing `libretto wiki` subcommand, producing one
self-contained page a browser opens locally. Static file, not a server — the example
given was a page reached by a link, and the markdown wiki's own reasoning (no daemon,
no infrastructure) still holds.
