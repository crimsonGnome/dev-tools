#!/usr/bin/env bash
set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEV_TOOLS="$REPO_ROOT/dev-tools"

echo "Installing dev-tools from $DEV_TOOLS..."

# Symlink steering docs into each package
symlink_steering() {
  local src="$DEV_TOOLS/steering/$1"
  local dest="$REPO_ROOT/$2/CLAUDE.md"
  mkdir -p "$(dirname "$dest")"
  ln -sf "$src" "$dest"
  echo "  Linked $1 → $2/CLAUDE.md"
}

symlink_steering "root.CLAUDE.md"     "."  2>/dev/null || ln -sf "$DEV_TOOLS/steering/root.CLAUDE.md" "$REPO_ROOT/CLAUDE.md" && echo "  Linked root.CLAUDE.md → CLAUDE.md"
symlink_steering "cdk.CLAUDE.md"      "vigilonCDK"
symlink_steering "api.CLAUDE.md"      "vigilonAPIHandler"
symlink_steering "frontend.CLAUDE.md" "front-end"

# Register bash aliases (appends to ~/.bashrc if not already present)
ALIAS_BLOCK='
# dev-tools agent aliases
alias agent-pair-programming="claude --system-prompt \"$(cat '"$DEV_TOOLS"'/agents/planning-agent.md)\""
alias agent-socrates="claude --system-prompt \"$(cat '"$DEV_TOOLS"'/agents/socrates.md)\""
alias agent-worker="claude --system-prompt \"$(cat '"$DEV_TOOLS"'/agents/agentic-worker.md)\""
'

if ! grep -q "dev-tools agent aliases" ~/.bashrc 2>/dev/null; then
  echo "$ALIAS_BLOCK" >> ~/.bashrc
  echo "  Added agent aliases to ~/.bashrc"
else
  echo "  Agent aliases already in ~/.bashrc — skipping"
fi

echo ""
echo "Done. Run 'source ~/.bashrc' to activate aliases."
