---
name: review-security
description: "Trigger: reviewing a diff, a file or an MR for vulnerabilities — injection, XSS, auth gaps, secrets, unsafe deserialization, SSRF. Reports only findings it is confident are exploitable, with the attacker's path traced. Standalone: works on any diff, with or without review-project."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Find what an attacker can actually use in the change under review, and nothing
else. A security report stuffed with theoretical maybes teaches its reader to skim,
and the one real injection scrolls past with the noise — **the false positive is
the enemy of the true one.**

`evidence` governs: a vulnerability reported is a vulnerability whose path was
traced, not a pattern that matched.

## The one question

Every candidate finding answers this before it is written down:

**Can an attacker reach this, with input they control?**

Trace it — do not assume it. Read where the value comes from, through every hop the
codebase shows. Reviewing the diff does not mean reading only the diff: the diff is
what gets *reported on*, the codebase is what gets *read* to decide whether a
finding is real.

| Attacker-controlled — investigate | Operator-controlled — usually safe |
|---|---|
| request params, body, headers, cookies | settings and config files |
| URL path segments | environment variables |
| file uploads, names and content | hardcoded constants |
| content other users wrote | values only admins can set |
| websocket and queue messages | framework internals |

A pattern fed only by the right-hand column is not a finding.

## What never gets flagged

- a pattern the framework already mitigates — auto-escaped template output, ORM
  parameterization, a sanitizer upstream. Check before flagging, and flag the
  bypasses instead: `mark_safe`, `dangerouslySetInnerHTML`, `v-html`, raw SQL
  string-building
- test files, dead code, comments
- defense-in-depth advice with no exploit behind it — a missing header on an
  internal tool is a suggestion, not a finding
- weak crypto doing a non-crypto job: a checksum is allowed to be md5; a password
  hash is not

## What always gets flagged

- `eval`/`exec`/`system`/`shell=True` fed by anything user-reachable
- deserializing user data with a deserializer that executes (`pickle`,
  `yaml.load`, Java `ObjectInputStream`, PHP `unserialize`)
- SQL, command or template strings built by interpolation from user input
- a secret in the source — key, token, password, private key. This one needs no
  reachability argument; its presence is the finding
- authorization that trusts the client: an id taken from the request and used
  without an ownership check

## Confidence is part of the finding

Three levels, and only two of them reach the report:

- **confirmed** — the path from attacker input to the sink was traced and holds.
  Reported, with severity.
- **needs verification** — the pattern is dangerous but one hop could not be read
  (an external call, generated code). Reported as a question, stating exactly which
  hop is unverified. Never dressed up as confirmed.
- **theoretical** — dropped. Not softened into a "consider…", dropped.

Severity, for the confirmed ones: **critical** (exploitable unauthenticated, severe
impact), **high** (exploitable with conditions), **medium** (narrow conditions,
moderate impact), **low** (real but marginal). One line of impact each — what the
attacker gets, not what category the textbook files it under.

## Output

Per finding: location (`file:line`), what it is, the traced path from attacker
input to sink, impact, severity or the unverified hop, and the fix in one line.

No findings is a statement, not an absence: **"nothing I am confident is
exploitable"** — which is not a certificate that nothing exists.

This lens reports; it never blocks, edits, commits or pushes.
