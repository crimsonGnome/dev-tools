package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SessionState holds the runtime state for a single tmux session.
type SessionState struct {
	Status  string `json:"status"`  // running | idle | stopped | error
	PID     int    `json:"pid"`
	Started string `json:"started"` // RFC3339 timestamp or empty
}

// State is the full orchestrator state persisted to disk.
type State struct {
	Sessions    map[string]SessionState `json:"sessions"`
	LastCommand string                  `json:"last_command"`
	LastStatus  string                  `json:"last_status"` // success | error
	Pending     []string                `json:"pending"`
}

// Load reads state from path. If the file does not exist, a fresh empty State
// is returned (not an error).
func Load(path string) (*State, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return &State{
			Sessions: make(map[string]SessionState),
			Pending:  []string{},
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("state: open %q: %w", path, err)
	}
	defer f.Close()

	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("state: decode %q: %w", path, err)
	}
	if s.Sessions == nil {
		s.Sessions = make(map[string]SessionState)
	}
	if s.Pending == nil {
		s.Pending = []string{}
	}
	return &s, nil
}

// Save writes state to path atomically via a temp file + rename.
func Save(path string, s *State) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".state-*.tmp")
	if err != nil {
		return fmt.Errorf("state: create temp: %w", err)
	}
	tmpName := tmp.Name()

	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("state: encode: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("state: close temp: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("state: rename to %q: %w", path, err)
	}
	return nil
}

// Reconcile cross-references stored PIDs against /proc/<pid>.
// Any session whose PID no longer exists in /proc is marked stopped.
func Reconcile(s *State) {
	for name, ss := range s.Sessions {
		if ss.PID <= 0 {
			continue
		}
		if _, err := os.Stat(fmt.Sprintf("/proc/%d", ss.PID)); os.IsNotExist(err) {
			ss.Status = "stopped"
			ss.PID = 0
			s.Sessions[name] = ss
		}
	}
}
