# add-wiki-open

Tracker: none

## The ask, verbatim

> vale pero cuando le doy a wiki (render this project's specs into
> …/libretto-automata/.agents/specs) que me abra un navegador localhost con el html

## Reading

Not a bug — the panel's wiki row should end with the viewer on screen, not with a
file on disk. Pressing the row generates/refreshes the HTML view and opens it in
the default browser.

"localhost" read as *the browser opens with the wiki*, not as a demand for an HTTP
daemon: the page is self-contained by contract and proven to work from `file://`,
which is the same browser tab with no server lifetime to manage inside a TUI. If a
real http URL is ever needed (sharing it, or a browser policy blocking `file://`),
that is the recorded bring-back condition for `--serve`.

Delivered on the same branch/PR (#67), following the panel-row pattern.
