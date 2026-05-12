#!/usr/bin/env bash
set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEV_TOOLS="$REPO_ROOT/dev-tools"

echo "Installing dev-tools from $DEV_TOOLS..."

# ---------------------------------------------------------------------------
# Symlink steering docs into each package
# ---------------------------------------------------------------------------
symlink_steering() {
  local src="$DEV_TOOLS/steering/$1"
  local dest="$REPO_ROOT/$2/CLAUDE.md"
  mkdir -p "$(dirname "$dest")"
  ln -sf "$src" "$dest"
  echo "  Linked $1 → $2/CLAUDE.md"
}

symlink_steering "root.CLAUDE.md"     "."  2>/dev/null || \
  ln -sf "$DEV_TOOLS/steering/root.CLAUDE.md" "$REPO_ROOT/CLAUDE.md" && \
  echo "  Linked root.CLAUDE.md → CLAUDE.md"
symlink_steering "cdk.CLAUDE.md"      "vigilonCDK"
symlink_steering "api.CLAUDE.md"      "vigilonAPIHandler"
symlink_steering "frontend.CLAUDE.md" "front-end"

# ---------------------------------------------------------------------------
# Remove legacy alias block (aliases cause premature backtick expansion)
# ---------------------------------------------------------------------------
if grep -q "dev-tools agent aliases" ~/.bashrc 2>/dev/null; then
  sed -i '/# dev-tools agent aliases/d' ~/.bashrc
  sed -i '/^alias agent-planning=/d'    ~/.bashrc
  sed -i '/^alias agent-pair-programming=/d' ~/.bashrc
  sed -i '/^alias agent-socrates=/d'    ~/.bashrc
  sed -i '/^alias agent-worker=/d'      ~/.bashrc
  echo "  Removed legacy agent aliases from ~/.bashrc"
fi

# ---------------------------------------------------------------------------
# Register agent shell functions
#
# Why functions instead of aliases?
# Aliases expand their stored text at invocation time — if that text contains
# backticks (common in Markdown system-prompt files) bash re-interprets them
# as command substitutions, causing errors like:
#   bash: agentAI/agentic-worker/context.json: No such file or directory
#
# Functions evaluate $(cat ...) only when called; the file content is passed
# as a plain string to claude, so backticks inside .md files are never
# interpreted by the shell.
# ---------------------------------------------------------------------------
FUNCTIONS_BLOCK="
# dev-tools agent functions
agent-pair-programming() { claude --system-prompt \"\$(cat ${DEV_TOOLS}/agents/planning-agent.md)\" \"\$@\"; }
agent-socrates()         { claude --system-prompt \"\$(cat ${DEV_TOOLS}/agents/socrates.md)\" \"\$@\"; }
agent-agentic-coder()    { claude --system-prompt \"\$(cat ${DEV_TOOLS}/agents/agentic-worker.md)\" \"\$@\"; }
agent-agentic-worker()   { claude --system-prompt \"\$(cat ${DEV_TOOLS}/agents/agentic-worker.md)\" \"\$@\"; }
"

if ! grep -q "dev-tools agent functions" ~/.bashrc 2>/dev/null; then
  printf '%s\n' "$FUNCTIONS_BLOCK" >> ~/.bashrc
  echo "  Added agent functions to ~/.bashrc"
else
  echo "  Agent functions already in ~/.bashrc — skipping"
fi

# ---------------------------------------------------------------------------
# Discord Orchestrator — build Go binary and install systemd user service
# ---------------------------------------------------------------------------
ORCH_DIR="$DEV_TOOLS/discord-orchestrator"
ORCH_BIN="$ORCH_DIR/bin/discord-orchestrator"
ORCH_SERVICE_TMPL="$ORCH_DIR/discord-orchestrator.service"
ORCH_SERVICE_DEST="$HOME/.config/systemd/user/discord-orchestrator.service"

echo "Building Discord Orchestrator..."
mkdir -p "$ORCH_DIR/bin"
(cd "$ORCH_DIR" && GOROOT=/usr/local/go GOPATH="$HOME/go" /usr/local/go/bin/go build -o "$ORCH_BIN" .)
echo "  Built $ORCH_BIN"

echo "Installing systemd user service..."
mkdir -p "$(dirname "$ORCH_SERVICE_DEST")"
sed "s|REPO_ROOT_PLACEHOLDER|$REPO_ROOT|g" "$ORCH_SERVICE_TMPL" > "$ORCH_SERVICE_DEST"
echo "  Wrote $ORCH_SERVICE_DEST"

systemctl --user daemon-reload
systemctl --user enable discord-orchestrator
systemctl --user restart discord-orchestrator
echo "  discord-orchestrator service enabled and started"

echo ""
echo "Done. Run 'source ~/.bashrc' to activate the agent functions."
