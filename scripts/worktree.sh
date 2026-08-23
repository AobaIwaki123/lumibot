#!/usr/bin/env bash
# ==============================================================================
# worktree.sh - Git Worktree management helper for lumibot
# ==============================================================================

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

usage() {
  echo "Usage: $0 [create <branch-name>|remove <branch-name>|list]"
  exit 1
}

CMD="${1:-}"

case "${CMD}" in
  create)
    BRANCH="${2:-}"
    if [[ -z "${BRANCH}" ]]; then
      echo "Error: Branch name required."
      usage
    fi
    WORKTREE_PATH=".worktrees/${BRANCH}"
    if [[ -d "${WORKTREE_PATH}" ]]; then
      echo "Worktree '${WORKTREE_PATH}' already exists."
    else
      mkdir -p .worktrees
      git worktree add -b "${BRANCH}" "${WORKTREE_PATH}" main
      echo "Created worktree at ${WORKTREE_PATH} on branch ${BRANCH}"
    fi
    ;;
  remove)
    BRANCH="${2:-}"
    if [[ -z "${BRANCH}" ]]; then
      echo "Error: Branch name required."
      usage
    fi
    WORKTREE_PATH=".worktrees/${BRANCH}"
    if [[ -d "${WORKTREE_PATH}" ]]; then
      git worktree remove "${WORKTREE_PATH}"
      echo "Removed worktree at ${WORKTREE_PATH}"
    else
      echo "Worktree '${WORKTREE_PATH}' does not exist."
    fi
    ;;
  list)
    git worktree list
    ;;
  *)
    usage
    ;;
esac
