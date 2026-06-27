#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scripts/clean_appledouble.sh [--check] [--include-git]

Deletes macOS AppleDouble sidecars (._*) and .DS_Store files from the repo.
Delete mode includes .git so transient sidecars cannot break git commands.
Check mode scans the working tree by default; use --include-git for a full scan.
EOF
}

mode="delete"
include_git="auto"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --check)
      mode="check"
      ;;
    --include-git)
      include_git="yes"
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage >&2
      exit 2
      ;;
  esac
  shift
done

if [[ "$include_git" == "auto" ]]; then
  if [[ "$mode" == "delete" ]]; then
    include_git="yes"
  else
    include_git="no"
  fi
fi

root="$(git rev-parse --show-toplevel)"
tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT
if [[ "$include_git" == "yes" ]]; then
  find "$root" \
    \( -name '._*' -o -name '.DS_Store' \) \
    -type f \
    -print >"$tmp"
else
  find "$root" \
    \( -path "$root/.git" -o -path "$root/.git/*" \) -prune \
    -o \( -name '._*' -o -name '.DS_Store' \) \
    -type f \
    -print >"$tmp"
fi

if [[ ! -s "$tmp" ]]; then
  if [[ "$include_git" == "yes" ]]; then
    echo "No AppleDouble or .DS_Store files found."
  else
    echo "No AppleDouble or .DS_Store files found in the working tree."
  fi
  exit 0
fi

cat "$tmp"
if [[ "$mode" == "check" ]]; then
  echo "AppleDouble or .DS_Store files found." >&2
  exit 1
fi

count=0
while IFS= read -r file; do
  rm -f -- "$file"
  count=$((count + 1))
done <"$tmp"
echo "Deleted $count AppleDouble/.DS_Store file(s)."
