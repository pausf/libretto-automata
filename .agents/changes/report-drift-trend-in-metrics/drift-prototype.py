"""Throwaway prototype: per-capability drift events from git history.

Usage: python3 -I drift.py <repo> [mode]
  mode = strict   : a code commit is covered only when it touches the capability's spec.md
  mode = change   : ... or touches any file in a change folder whose Targets include the capability (default)
  mode = anychange: ... or touches any file under .agents/changes/
Prints per-capability counts, then every event as TSV to events.tsv beside this script.
"""
import fnmatch, os, re, subprocess, sys, collections

repo = sys.argv[1]
mode = sys.argv[2] if len(sys.argv) > 2 else "change"
here = os.path.dirname(os.path.abspath(__file__))

def git(*a):
    return subprocess.run(["git", "-C", repo, *a], capture_output=True, text=True, check=True).stdout

def defenced(text):
    out, f = [], False
    for l in text.split("\n"):
        if re.match(r"^\s*```", l):
            f = not f
            continue
        if not f:
            out.append(l)
    return out

def governs_of(text):
    globs = []
    for l in defenced(text):
        m = re.match(r"^\s*Governs:(.*)", l)
        if m:
            globs += m.group(1).split()
    return globs

def matches_any(p, globs):
    return any(g and fnmatch.fnmatchcase(p, g) for g in globs)

# 1. capabilities at HEAD
caps = {}
for d in sorted(os.listdir(os.path.join(repo, ".agents/specs"))):
    s = os.path.join(repo, ".agents/specs", d, "spec.md")
    if os.path.isfile(s):
        caps[d] = (f".agents/specs/{d}/spec.md", governs_of(open(s).read()))

# 2. change -> targets, from every added Targets: line in history
targets = collections.defaultdict(set)
cur = None
for l in git("log", "--full-history", "-p", "-U0", "--format=", "--", ".agents/changes/").split("\n"):
    if l.startswith("+++ b/.agents/changes/"):
        cur = l[len("+++ b/.agents/changes/"):].split("/")[0]
    elif l.startswith("+++ "):
        cur = None
    elif cur and re.match(r"^\+\s*Targets:", l):
        for t in l.split("Targets:", 1)[1].replace(",", " ").split():
            targets[cur].add(t.strip("`"))

# 3. commits, oldest first, no merges
commits, c = [], None
unit = len(sys.argv) > 3 and sys.argv[3] == "unit"
logargs = ["--first-parent"] if unit else ["--topo-order", "--no-merges"]
for l in git("log", "--reverse", *logargs, "--name-only",
             "--format=C\t%H\t%ct\t%s").split("\n"):
    if l.startswith("C\t"):
        _, h, t, s = l.split("\t", 3)
        c = {"sha": h, "t": int(t), "subj": s, "files": []}
        commits.append(c)
    elif l.strip() and c:
        c["files"].append(l.strip())

spec_paths = {v[0] for v in caps.values()}
def is_doc(p):
    return p in spec_paths or p.startswith(".agents/changes/") or p.startswith(".agents/specs/")

# 4. walk
events = []
warned = []
stats = {}
for cap, (spath, globs) in caps.items():
    pending = []
    code_commits = covered_n = same = 0
    born = not os.environ.get("BORN")
    for c in commits:
        fs = c["files"]
        touches_s = spath in fs
        if not born:
            if not touches_s:
                continue
            born = True
        changes = {f.split("/")[2] for f in fs if f.startswith(".agents/changes/") and f.count("/") >= 3}
        code = [f for f in fs if not is_doc(f) and matches_any(f, globs)]
        if code:
            code_commits += 1
            covered = touches_s
            if mode == "change":
                covered = covered or any(cap in targets.get(ch, ()) for ch in changes)
            elif mode == "anychange":
                covered = covered or bool(changes)
            if not covered and os.environ.get("SHARED"):
                moved = {k for k, v in caps.items() if v[0] in fs}
                covered = all(any(matches_any(f, caps[k][1]) for k in moved) for f in code)
            if covered:
                covered_n += 1
                if touches_s:
                    same += 1
            else:
                pending.append(c)
                warned.append((cap, c["sha"][:7], c["subj"], " ".join(code[:3])))
        if touches_s and pending:
            events.append((cap, c, list(pending)))
            pending = []
    stats[cap] = (code_commits, covered_n, same, len(pending))

per = collections.Counter(e[0] for e in events)
lag = collections.Counter()
for cap, c, p in events:
    lag[cap] += len(p)
print(f"mode={mode} unit={unit} commits={len(commits)} capabilities={len(caps)}")
print(f"{'capability':16} {'code':>5} {'covered':>7} {'same':>5} {'events':>6} {'lagging':>7} {'open':>5}")
for cap in caps:
    cc, cv, sm, op = stats[cap]
    print(f"{cap:16} {cc:5} {cv:7} {sm:5} {per[cap]:6} {lag[cap]:7} {op:5}")
print(f"total events {len(events)}")
with open(os.path.join(here, f"events-{mode}{'-unit' if unit else ''}.tsv"), "w") as f:
    for cap, c, p in events:
        f.write("\t".join([cap, c["sha"][:7], c["subj"], str(len(p)),
                           " | ".join(x["sha"][:7] + " " + x["subj"] for x in p)]) + "\n")

with open(os.path.join(here, "warned.tsv"), "w") as f:
    for w in warned:
        f.write("\t".join(w) + "\n")
