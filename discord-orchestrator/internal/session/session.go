package session

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eggersjoseph/discord-orchestrator/internal/config"
)

// LogPath returns the absolute path to the log file for a named session.
func LogPath(repoRoot, name string) string {
	return filepath.Join(repoRoot, "dev-tools", "discord-orchestrator", "logs", name+".log")
}

// Start creates a new tmux session for the given name, sources install.sh,
// and runs the session's command with stdout+stderr redirected to its log file.
func Start(name, repoRoot string, def config.SessionDef) error {
	logPath := LogPath(repoRoot, name)
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return fmt.Errorf("session: create log dir: %w", err)
	}

	installSh := filepath.Join(repoRoot, "dev-tools", "install.sh")

	// Build the shell command that runs inside the tmux session:
	//   source install.sh && <cmd> >> <logfile> 2>&1
	shellCmd := fmt.Sprintf(
		"source %q && %s >> %q 2>&1",
		installSh,
		def.Cmd,
		logPath,
	)

	cmd := exec.Command(
		"tmux", "new-session", "-d",
		"-s", name,
		"bash", "-c", shellCmd,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("session: start %q: %w\n%s", name, err, out)
	}
	return nil
}

// Stop kills the named tmux session.
func Stop(name string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("session: stop %q: %w\n%s", name, err, out)
	}
	return nil
}

// Inject sends a message into the named tmux session via send-keys.
// It refuses to inject into sessions with status "stopped" or "error".
func Inject(name, message, status string) error {
	if status == "stopped" || status == "error" {
		return fmt.Errorf("session %q is %s — cannot inject", name, status)
	}
	// Escape any single quotes in the message to avoid shell injection
	safe := strings.ReplaceAll(message, "'", "'\\''")
	cmd := exec.Command("tmux", "send-keys", "-t", name, safe, "Enter")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("session: inject %q: %w\n%s", name, err, out)
	}
	return nil
}

// Status returns the status string for the named tmux session.
// Returns "running" if the session exists, "stopped" if it does not.
func Status(name string) (string, error) {
	cmd := exec.Command("tmux", "has-session", "-t", name)
	if err := cmd.Run(); err != nil {
		return "stopped", nil
	}
	return "running", nil
}

// List returns all current tmux session names.
func List() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
	out, err := cmd.Output()
	if err != nil {
		// No sessions running is not an error
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return []string{}, nil
		}
		return nil, fmt.Errorf("session: list: %w", err)
	}
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return []string{}, nil
	}
	return strings.Split(raw, "\n"), nil
}
