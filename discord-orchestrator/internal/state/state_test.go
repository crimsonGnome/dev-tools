package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	original := &State{
		Sessions: map[string]SessionState{
			"socrates": {Status: "running", PID: 12345, Started: "2026-05-11T10:00:00Z"},
		},
		LastCommand: "start socrates",
		LastStatus:  "success",
		Pending:     []string{},
	}

	if err := Save(path, original); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.LastCommand != original.LastCommand {
		t.Errorf("LastCommand = %q, want %q", loaded.LastCommand, original.LastCommand)
	}
	ss, ok := loaded.Sessions["socrates"]
	if !ok {
		t.Fatal("Sessions[socrates] missing after round-trip")
	}
	if ss.Status != "running" || ss.PID != 12345 {
		t.Errorf("Sessions[socrates] = %+v", ss)
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	s, err := Load(filepath.Join(dir, "nonexistent.json"))
	if err != nil {
		t.Fatalf("Load missing file should not error, got: %v", err)
	}
	if s == nil {
		t.Fatal("Load returned nil State")
	}
	if len(s.Sessions) != 0 {
		t.Errorf("expected empty sessions, got %v", s.Sessions)
	}
}

func TestReconcileMarksStalePIDStopped(t *testing.T) {
	// PID 999999999 almost certainly does not exist
	s := &State{
		Sessions: map[string]SessionState{
			"dead": {Status: "running", PID: 999999999},
		},
	}
	Reconcile(s)
	ss := s.Sessions["dead"]
	if ss.Status != "stopped" {
		t.Errorf("expected stopped, got %q", ss.Status)
	}
	if ss.PID != 0 {
		t.Errorf("expected PID 0 after reconcile, got %d", ss.PID)
	}
}

func TestReconcileLeavesValidPIDAlone(t *testing.T) {
	// Use PID 1 (init/systemd) which always exists on Linux
	s := &State{
		Sessions: map[string]SessionState{
			"init": {Status: "running", PID: 1},
		},
	}
	Reconcile(s)
	ss := s.Sessions["init"]
	if ss.Status != "running" {
		t.Errorf("expected running, got %q", ss.Status)
	}
}

func TestAtomicWriteProducesFinalFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	s := &State{
		Sessions:   map[string]SessionState{},
		LastStatus: "success",
		Pending:    []string{},
	}
	if err := Save(path, s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Final file must exist
	if _, err := os.Stat(path); err != nil {
		t.Errorf("final file missing: %v", err)
	}

	// No temp files should remain
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.Name() != "state.json" {
			t.Errorf("unexpected leftover file: %s", e.Name())
		}
	}
}
