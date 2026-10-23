# Pushing, the request and the bump

Read this when the repository has a remote, before the push question in `SKILL.md`.

It is read **after** the push section on purpose. When the exception below sat above
the native-prompt paragraph, what followed it read as reinstating the ask
unconditionally. That was the reviewer's finding on the run that landed the mode.

## The one exception, and it is the whole of `/libretto-attacca`

**Everything in `SKILL.md`'s push section — the question, the native prompt, the
waiting — is the attended flow.** Under `/libretto-attacca` the question is not asked at all, because it was
answered when the command was typed. "Never push unasked" is intact rather than
overridden: the asking happened at the prompt.

The answer covers **this branch and this request and nothing past it** — no merge, no tag,
no release, and no `release:` label the user did not choose. The bump is asked separately
at the very end and is never covered by this consent: typing an answer the user gave is
not the same act as assuming one, and the invocation answered the push, not the version.

And it is paid for in the request's description, which carries three things or the run has
bought silence rather than speed:

- **what the invocation answered** — the one stop and the push — so a reader can tell
  which of the decisions in front of them a person made
- **every question the run assumed past**, each with what changes if it is wrong
- **the bump a person chose**, and that a person chose it — see the bump section below. A
  reviewer has to be able to tell an agreed bump from an assumed one, and that distinction
  is what this whole command is built out of. An unlabeled request says the question went
  unanswered, which is its own honest state and not an omission to be tidied away.

The yes path is the yes path: what follows below about returning to the base branch happens
here too, because the push and the request were both confirmed.

The two used to be asked separately. That bought a second round trip and no safety:
pushing a branch and then declining to open the request for it is a state almost
nobody wants, and the user who wants exactly that says so in the same breath. Asking
twice about one intent is the ceremony this flow is supposed to spend elsewhere.

When the answer is yes, confirm it landed: **the remote tip matches the local tip.** A
push that printed no error is not a push that was accepted. Then read the created
request back from the forge rather than trusting the output of the command that made
it.

## Then return to the base branch, current

Once the push and the request are both confirmed, and only then:

```
git checkout <base>
git pull --ff-only
```

The working tree ends on the base branch, up to date with the remote. **The feature
branch is left alone** — not deleted, not merged, not rebased. The request is open, not
merged, and a branch deleted at that moment takes the only local copy of unmerged work
with it.

The base is what the repository says it is — `git symbolic-ref refs/remotes/origin/HEAD`,
or the branch the request was opened against. Never hardcoded to `main`.

**Two conditions, and both are stops rather than fixes:**

- **a dirty tree.** `git checkout` carries uncommitted work onto the base branch
  silently, which is the "branch before the first write" failure pointed the other way.
  Report what is uncommitted and stay put.
- **a fast-forward that refuses.** A diverged base is a fact worth seeing. Never
  `git pull` without `--ff-only`, and never resolve it here: a merge commit manufactured
  on somebody's base branch by a phase whose job is bookkeeping is the surprise nobody
  wants.

**Only on the yes path.** No push means the work exists nowhere but this branch, and
moving off it buys nothing.

Why this is here at all: a session starts wherever the last one left the working
directory, so a flow that ends on a merged feature branch hands the next phase 1 a stale
base to read. That is not hypothetical — phase 1 of the run that added this reported a
branch as work in flight when it had already been merged and tagged, because local `main`
was seven commits behind. The reading was wrong, it was reported as fact, and the
correction cost a round trip.

## Last of all, under `/libretto-attacca`: the bump

**This is the only question an unattended run asks, and it is the last thing that
happens** — after the report, after the request is confirmed open, after the return to
the base branch. There is nothing downstream of it, and that is what makes it not a
stop: the work is complete and reviewable whether it is answered or not.

It exists because a designed refusal that nobody predicts reads as breakage. The
`release:` label is required before a merge, so an unattended run's request carries a red
check by design — that red *is* the bump question arriving. It arrived as a broken
pipeline once, and cost an alarm and a round trip for a decision the run was standing
next to.

**The run never decides the bump. It asks, and it types the answer.** `AGENTS.md` states
the split: *the reading is yours, the typing is not*. Everything below holds that line.

First, does the repository define the labels at all?

```
gh label list --search 'release:' --limit 100 --json name
```

**Never a bare `gh label list`.** It fetches 30 by default, ascending, so `release:*`
sorts late and falls off the page in any repository past thirty labels — and the answer
comes back "none of the three" for a repository that defines all three. That inverts the
whole check silently, which is the worst way for it to be wrong.

Then match **whole names** against exactly `release:patch`, `release:minor` and
`release:major`. Never a prefix or a substring: `release:patch-hotfix` contains
`release:patch` and is not it. The workflow that reads these labels matches whole names
for the same reason, and a detection looser than the check it feeds is a detection that
finds a label `gh pr edit` will then be refused.

A repository that **defines none of the three** is not asked and is not told why — the
convention does not exist there, so neither does the question. Nothing is created:
inventing `release:minor` in somebody else's repository is deciding that repository's
release convention on their behalf.

Where they do exist, **the bump is asked once** with `AskUserQuestion` — the run's own
reading recommended first, then the others, and room to answer differently.

**`release:major` is present and is never the first option.** Selectable, because a
native question the user answers *is* the asking its standing rule demands. Never
recommended, because recommending it is the announcement that rule forbids — and
announcing a major three times and proceeding is what published `v1.0.0` and `v1.0.1`
from a table read without the paragraph above it. A version number cannot be recalled
once the proxy has cached it.

Then apply exactly one — the workflow refuses two — and **read the request back** off the
forge to confirm the label is on it:

```
gh pr edit <n> --add-label <label>
gh pr view <n> --json labels
```

The same rule the push already carries, for the same reason: a command that printed no
error is not a change the forge accepted.

**Then put it in the description, and this step is not optional.** The description was
written before the question was asked — the request had to exist for there to be anything
to label — so the bump reaches it only by being written back:

```
gh pr edit <n> --body "<the description, plus the bump a person chose>"
```

Without this the third bullet above is a promise the run never keeps: the description
lists what the invocation answered and what it assumed, and a reader has no way to tell
that the bump was neither. The label alone does not say who chose it.

**Unanswered, the run ends exactly as it does today** — unlabeled, and the closing
report's red-check line *is never withdrawn* by the question. That line is written before
the question is asked. A report that promised a red check and then quietly labelled the
request has lied about the state the user will find, which is the same failure as the
silence, wearing the opposite face.

**No default. Ever.** Not patch, not "the safe one". A silently-wrong bump is the failure
that published `v1.0.0` with a politer name on it. Headless runs — `libretto loop` — have
no prompt for `AskUserQuestion` to arrive in, so unanswered is the normal path there and
must be the quiet one.

**Ceiling named:** the question is only as reachable as the terminal it is asked in. A
scheduled or piped run never sees it and lands unlabeled — today's behaviour, which is
why today's behaviour is the fallback rather than an error. The replacement, if that
becomes the common case, is a `gh pr comment` carrying the three commands. Deliberately
not built for a case that has not happened.

**Attended runs do not ask this.** `/libretto-flow` already stops at phase 8 with the user
present and watching; the red check ambushes nobody there. The user's call, 2026-08-13 —
back the day an attended run pays the same round trip.

## Which forge, and whether there is one

Derive it. Do not assume, and do not ask what the repository can answer:

```
git remote get-url origin
```

| The URL contains | The tool |
|---|---|
| `github.com` | `gh pr create` |
| `gitlab` | `glab mr create` |

**When the bump is decided before the request exists, the label rides `gh pr create`** —
`gh pr create --label release:<bump>`, never a second `gh pr edit` after it. The user named
the bump, or authorised one in advance: then there is nothing left to ask, and a request
opened unlabelled fails the label check's `opened` run and mails the user a red check for a
question already answered. That happened three times in one day. The other side is
untouched: under `/libretto-attacca` with no bump decided, the request is opened unlabelled
on purpose and its red check *is* the bump question — see the section on the bump above.

A missing or unauthenticated CLI **stops**, with the install line and nothing else —
the shape `skills/find-work/` already uses for `jira`:

```
brew install gh          # then: gh auth login
brew install glab        # then: glab auth login
```

`gh auth login` and `glab auth login` are interactive. The user runs them
themselves — and **no token ever comes through the conversation.** A credential
pasted to an agent is in the transcript and the logs from that moment on, and
rotating it is the only remedy.

Do not offer a workaround. A hand-built API call with a token found in the
environment is how a stop becomes an exposure, and having the other forge's CLI
installed is not a fallback: `glab` cannot open a request on GitHub.

**Ceiling named:** this is a substring test on one URL. It does not survive a
self-hosted forge on a neutral domain, or Gitea and Forgejo. When that day comes the
answer is an explicit setting read from the repository, not a longer list of guesses.

**The description is written, not omitted.** One or three bullets of what the change
does, the evidence that it works, and what was deliberately left out — the phase 7
report already contains all three. `gitlab-mr-description` or an equivalent, if
installed, may shape it.
