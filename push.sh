#!/usr/bin/env bash
set -euo pipefail

# Push local RouteFast repo to GitHub
#
# Usage:
#   chmod +x push-routefast-to-github.sh
#   GITHUB_USERNAME=yourname ./push-routefast-to-github.sh
#
# Assumes:
#   - you are inside the repo folder
#   - commits already exist
#   - tags may or may not exist

GITHUB_USERNAME="${GITHUB_USERNAME:-YOUR_USERNAME}"
REPO_NAME="routefast-ce"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || { echo "Missing required command: $1" >&2; exit 1; }
}

need_cmd git

if [ ! -d .git ]; then
  echo "Not inside a git repository"
  exit 1
fi

REMOTE_URL="https://github.com/${GITHUB_USERNAME}/${REPO_NAME}.git"

echo "[1/5] Checking repo status..."
git status

echo "[2/5] Adding all files..."
git add .

if ! git diff --cached --quiet; then
  echo "[3/5] Creating commit..."
  git commit -m "RouteFast CD initial import"
else
  echo "[3/5] No new changes to commit"
fi

echo "[4/5] Setting remote..."
git remote remove origin >/dev/null 2>&1 || true
git remote add origin "$REMOTE_URL"

echo "[5/5] Pushing to GitHub..."

# push main branch
git push -u origin main

# push tags if they exist
if [ -n "$(git tag)" ]; then
  git push origin --tags
fi

echo
echo "✅ Done."
echo "Repo URL:"
echo "  https://github.com/${GITHUB_USERNAME}/${REPO_NAME}"
