#!/usr/bin/env bash
# 𝄞 Libretto Automata — symlink installer
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEST="${CLAUDE_HOME:-$HOME/.claude}"
DIRS=(skills agents commands)

# ponytail: repo dirs are the source of truth, no manifest file to drift
each_item() {
  for d in "${DIRS[@]}"; do
    for src in "$REPO/$d"/*; do
      [ -e "$src" ] || continue
      "$1" "$src" "$DEST/$d/$(basename "$src")"
    done
  done
}

link() {
  local src="$1" dst="$2"
  if [ -e "$dst" ] && [ ! -L "$dst" ]; then
    echo "  skip  $dst (exists, not a symlink)"
    return
  fi
  mkdir -p "$(dirname "$dst")"
  ln -sfn "$src" "$dst"
  echo "  link  $dst"
}

show() {
  local src="$1" dst="$2"
  if [ "$(readlink "$dst" 2>/dev/null)" = "$src" ]; then
    echo "  ok    $dst"
  else
    echo "  MISS  $dst"
  fi
}

check() {
  local dst="$2"
  [ -L "$dst" ] && [ ! -e "$dst" ] && echo "  BROKEN $dst -> $(readlink "$dst")"
  return 0
}

case "${1:-install}" in
  install) echo "𝄞 installing into $DEST"; each_item link ;;
  status)  echo "𝄞 status"; each_item show ;;
  doctor)  echo "𝄞 doctor"; each_item check; echo "  done" ;;
  *) echo "usage: $0 {install|status|doctor}" >&2; exit 2 ;;
esac
