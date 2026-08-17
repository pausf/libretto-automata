# make-test-badge-live

Tracker: none

## The ask, verbatim

> Si el objetivo es adopción externa, añadir CI (workflows visibles) o badges de
> build/test facilitaría confianza inmediata (aunque hay tests locales).

## Reading

The premise is wrong on both halves, and what it uncovers is worse than what it asked for.
CI exists and is visible: `.github/workflows/` carries `gates.yml`, `release.yml` and
`require-release-label.yml`. Badges exist too — seven of them in the README's first
thirteen lines.

The defect is that the tests badge is **static**:

```
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](.github/workflows/gates.yml)
```

That is a hardcoded green image linking to a workflow file. It says *passing* whether or
not anything passed, which is the exact failure this repository refuses everywhere else —
`skills/evidence/` exists because nothing is true until observed, and a badge asserting a
green run nobody observed is a claim the flow would reject in a report.

The fix is one line: the real endpoint,
`https://github.com/pausf/libretto-automata/actions/workflows/gates.yml/badge.svg`, which
turns red on its own. The other six are version and tool labels, not run status, and are
honest as they are.
