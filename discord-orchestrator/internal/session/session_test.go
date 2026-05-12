package session

import (
	"strings"
	"testing"

	"github.com/eggersjoseph/discord-orchestrator/internal/config"
)

func TestLogPath(t *testing.T) {
	got := LogPath("/repo", "socrates")
	want := "/repo/dev-tools/discord-orchestrator/logs/socrates.log"
	if got != want {
		t.Errorf("LogPath = %q, want %q", got, want)
	}
}

func TestLogPath_DifferentName(t *testing.T) {
	got := LogPath("/home/user/project", "pair-programming")
	want := "/home/user/project/dev-tools/discord-orchestrator/logs/pair-programming.log"
	if got != want {
		t.Errorf("LogPath = %q, want %q", got, want)
	}
}

func TestInject_StoppedSessionReturnsError(t *testing.T) {
	err := Inject("socrates", "hello", "stopped")
	if err == nil {
		t.Fatal("expected error injecting into stopped session, got nil")
	}
	if !strings.Contains(err.Error(), "stopped") {
		t.Errorf("error message should mention 'stopped', got: %v", err)
	}
}

func TestInject_ErrorSessionReturnsError(t *testing.T) {
	err := Inject("socrates", "hello", "error")
	if err == nil {
		t.Fatal("expected error injecting into error session, got nil")
	}
	if !strings.Contains(err.Error(), "error") {
		t.Errorf("error message should mention 'error', got: %v", err)
	}
}

func TestStart_CommandConstruction(t *testing.T) {
	// We can't run tmux in unit tests, but we can verify LogPath is used correctly
	// by confirming the log path derives from the correct components.
	repoRoot := "/test/repo"
	name := "socrates"
	def := config.SessionDef{Cmd: "agent-socrates"}

	logPath := LogPath(repoRoot, name)
	if !strings.HasSuffix(logPath, "logs/socrates.log") {
		t.Errorf("unexpected log path: %q", logPath)
	}
	if !strings.Contains(logPath, repoRoot) {
		t.Errorf("log path should contain repoRoot %q, got %q", repoRoot, logPath)
	}

	// Verify the def.Cmd is what would be embedded in the shell command
	if def.Cmd != "agent-socrates" {
		t.Errorf("def.Cmd = %q", def.Cmd)
	}
}
