#!/usr/bin/env bash
# Pulls each package from GitHub at a specified path.
# Usage: ./create-workspace.sh
set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "Creating workspace at $REPO_ROOT..."

# Add package clone commands here, e.g.:
# git clone git@github.com:your-org/vigilonCDK.git "$REPO_ROOT/vigilonCDK"
# git clone git@github.com:your-org/vigilonAPIHandler.git "$REPO_ROOT/vigilonAPIHandler"
# git clone git@github.com:your-org/front-end.git "$REPO_ROOT/front-end"

echo "Workspace ready."
