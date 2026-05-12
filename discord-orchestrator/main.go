package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eggersjoseph/discord-orchestrator/internal/bot"
	"github.com/eggersjoseph/discord-orchestrator/internal/config"
	"github.com/eggersjoseph/discord-orchestrator/internal/state"
)

func main() {
	// Derive REPO_ROOT from the binary's location:
	//   <repo>/dev-tools/discord-orchestrator/bin/discord-orchestrator
	//   → Dir → bin/  → Dir → discord-orchestrator/  → Dir → dev-tools/  → Dir → <repo>
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("main: resolve executable path: %v", err)
	}
	// Resolve symlinks so the path is absolute and canonical
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		log.Fatalf("main: eval symlinks: %v", err)
	}
	repoRoot, err := filepath.Abs(filepath.Join(filepath.Dir(exePath), "..", "..", ".."))
	if err != nil {
		log.Fatalf("main: resolve repo root: %v", err)
	}

	// Fail fast if tmux is not installed — before connecting to Discord
	if _, err := exec.LookPath("tmux"); err != nil {
		log.Fatalf("main: tmux not found in PATH — install tmux before running the orchestrator")
	}

	// Load configuration
	cfgPath := filepath.Join(repoRoot, "dev-tools", "discord-orchestrator", "config.json")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("main: load config: %v", err)
	}

	// Load state (missing file is a fresh start, not an error)
	statePath := filepath.Join(repoRoot, "dev-tools", "discord-orchestrator", "orchestrator-state.json")
	s, err := state.Load(statePath)
	if err != nil {
		log.Fatalf("main: load state: %v", err)
	}

	// Reconcile: mark sessions whose PIDs are no longer alive as stopped
	state.Reconcile(s)
	if err := state.Save(statePath, s); err != nil {
		log.Printf("main: warning — could not save reconciled state: %v", err)
	}

	// Start the bot — blocks until interrupted
	b := bot.New(cfg, repoRoot, s, statePath)
	if err := b.Start(); err != nil {
		log.Fatalf("main: bot exited with error: %v", err)
	}
}
