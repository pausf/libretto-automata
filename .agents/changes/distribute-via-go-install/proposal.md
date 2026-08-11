# distribute-via-go-install

Tracker: none
Queued: 2026-08-11

## The ask, verbatim

> ahora para instalar o actualizar tengo que hacerlo desde github, pero veo que mucha
> gente instala sus cosas desde npx etc.. o por ejemplo go install y el repo de github ,
> me gustaria hacer un sistema de actualiza y instalar desde go y asi sería mas comodo y
> aparte cada vez que haya una actualizacion que en el libretto por ejemplo en un carte
> aqui

```
          ╭──────────────────────────────────────────────────────────────────────────────────────────────────╮
          │  ░▒▓█ ════════════════════════════════════════════════════════════════════════════════════ █▓▒░  │
          │                                                                                                  │
          │                        ▄▀▀▄                                                                      │
          │                       ▐▌  ▐▌    █    ▀█▀  █▀▄  █▀▄  █▀▀  ▀█▀  ▀█▀  ▄▀▄                           │
          │                     ──█▄▄▄▀──   █     █   █▀▄  █▀▄  █▀    █    █   █ █  ──                       │
          │                     ──█▀▀▀▄──   █▄▄  ▄█▄  █▄▀  █ ▀  █▄▄   █    █   ▀▄▀  ──                       │
          │                       ▐▌  ▐▌                                                                     │
          │                       ▐▙▄▄▟▘  ▏ A U T O M A T A                                                  │
          │                        ▐▌     ▏ the libretto is written first ·                                  │
          │                       ▄▀      ▏ the automaton performs it                                        │
          │                               ▏ b y   p a u s f                                                  │
          │                                                                                                  │
          │  ░▒▓█ ════════════════════════════════════════════════════════════════════════════════════ █▓▒░  │
          ├──────────────────────────────────────────────────────────────────────────────────────────────────┤
          │                                                                                                  │
          │  ❯ ▸ install     link the score into …/libretto-automata/.claude                                 │
          │    ▸ uninstall   take it back out of …/libretto-automata/.claude                                 │
          │    ▸ update      git pull · relink · report                                                      │
          │    ▸ status      32 linked                                                                       │
          │    ▸ models      3 on haiku · 4 on session                                                       │
          │    ▸ doctor      diagnose the orchestra                                                          │
          │    ▸ prune       drop links whose source is gone                                                 │
          │                                                                                                  │
          ├──────────────────────────────────────────────────────────────────────────────────────────────────┤
          │    ● global   32 missing              ~/.claude                                                  │
          │  ❯ ● project  32 linked               …/libretto-automata/.claude                                │
          ╰────────────────────────────────────────────────────────────────────────────
```

> que ponga actualiza a la siguiente version y te lo ponga para que puedas actualizar y
> saber que hay una nueva

## Reading

Installing and updating today means cloning the repo and pulling by hand. The idea is a
`go install`-style distribution so the binary comes from the module path, plus a line in
the panel that says a newer version exists and offers to move to it. Overlaps with
[[notify-users-of-new-updates]] on the "there is a new version" half.
