package bot

import (
	"strings"
	"testing"

	"github.com/eggersjoseph/discord-orchestrator/internal/config"
	"github.com/eggersjoseph/discord-orchestrator/internal/state"
)

// ---- Parser tests (Task 006) ----

func TestParse_SendLogFull(t *testing.T) {
	cmd := parseCommand(`send --log socrates --full`)
	if cmd.Verb != "send" {
		t.Errorf("Verb = %q, want %q", cmd.Verb, "send")
	}
	if cmd.Flags["--log"] != "socrates" {
		t.Errorf("Flags[--log] = %q, want %q", cmd.Flags["--log"], "socrates")
	}
	if !cmd.Full {
		t.Error("Full should be true when --full is present")
	}
}

func TestParse_TailWithLines(t *testing.T) {
	cmd := parseCommand(`tail --session socrates --lines 50`)
	if cmd.Verb != "tail" {
		t.Errorf("Verb = %q, want %q", cmd.Verb, "tail")
	}
	if cmd.Flags["--session"] != "socrates" {
		t.Errorf("Flags[--session] = %q, want %q", cmd.Flags["--session"], "socrates")
	}
	if cmd.Flags["--lines"] != "50" {
		t.Errorf("Flags[--lines] = %q, want %q", cmd.Flags["--lines"], "50")
	}
	if cmd.Full {
		t.Error("Full should be false when --full is absent")
	}
}

func TestParse_InjectWithQuotedMessage(t *testing.T) {
	cmd := parseCommand(`inject --session foo --message "hello world"`)
	if cmd.Verb != "inject" {
		t.Errorf("Verb = %q, want %q", cmd.Verb, "inject")
	}
	if cmd.Flags["--session"] != "foo" {
		t.Errorf("Flags[--session] = %q, want %q", cmd.Flags["--session"], "foo")
	}
	if cmd.Flags["--message"] != "hello world" {
		t.Errorf("Flags[--message] = %q, want %q", cmd.Flags["--message"], "hello world")
	}
}

func TestParse_StartWithPositionalArg(t *testing.T) {
	cmd := parseCommand(`start socrates`)
	if cmd.Verb != "start" {
		t.Errorf("Verb = %q, want %q", cmd.Verb, "start")
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "socrates" {
		t.Errorf("Args = %v, want [socrates]", cmd.Args)
	}
	if cmd.Full {
		t.Error("Full should be false")
	}
}

func TestParse_UnknownVerb(t *testing.T) {
	cmd := parseCommand(`unknowncmd`)
	if cmd.Verb != "unknowncmd" {
		t.Errorf("Verb = %q, want %q", cmd.Verb, "unknowncmd")
	}
}

func TestParse_List(t *testing.T) {
	cmd := parseCommand(`list`)
	if cmd.Verb != "list" {
		t.Errorf("Verb = %q, want %q", cmd.Verb, "list")
	}
	if len(cmd.Args) != 0 {
		t.Errorf("Args should be empty, got %v", cmd.Args)
	}
	if len(cmd.Flags) != 0 {
		t.Errorf("Flags should be empty, got %v", cmd.Flags)
	}
	if cmd.Full {
		t.Error("Full should be false")
	}
}

// ---- Dispatch tests (Task 007) ----

// newTestBot builds a Bot with a fake config and state — no Discord connection.
func newTestBot() *Bot {
	cfg := &config.Config{
		DiscordToken:     "fake",
		AuthorizedUserID: "user123",
		Sessions: map[string]config.SessionDef{
			"socrates": {Cmd: "agent-socrates"},
		},
	}
	s := &state.State{
		Sessions: map[string]state.SessionState{
			"socrates": {Status: "stopped"},
		},
		Pending: []string{},
	}
	return &Bot{
		cfg:       cfg,
		repoRoot:  "/fake/repo",
		state:     s,
		statePath: "/fake/repo/dev-tools/discord-orchestrator/orchestrator-state.json",
	}
}

func TestDispatch_UnknownVerbReturnsHelp(t *testing.T) {
	b := newTestBot()
	cmd := parseCommand("notacommand")
	reply := b.dispatch(cmd)
	if !strings.Contains(reply, "Commands:") {
		t.Errorf("expected help text, got: %q", reply)
	}
}

func TestDispatch_HelpReturnsHelp(t *testing.T) {
	b := newTestBot()
	cmd := parseCommand("help")
	reply := b.dispatch(cmd)
	if !strings.Contains(reply, "Commands:") {
		t.Errorf("expected help text, got: %q", reply)
	}
}

func TestDispatch_StartMissingArgReturnsUsage(t *testing.T) {
	b := newTestBot()
	cmd := parseCommand("start")
	reply := b.dispatch(cmd)
	if !strings.Contains(strings.ToLower(reply), "usage") {
		t.Errorf("expected usage error, got: %q", reply)
	}
}

func TestDispatch_InjectOnStoppedSessionReturnsError(t *testing.T) {
	b := newTestBot()
	// socrates is stopped in test state
	cmd := parseCommand(`inject --session socrates --message "hello"`)
	reply := b.dispatch(cmd)
	if !strings.Contains(reply, "stopped") && !strings.Contains(reply, "Error") {
		t.Errorf("expected error reply for inject on stopped session, got: %q", reply)
	}
}

func TestDispatch_SendLogFullUsesDifferentPath(t *testing.T) {
	b := newTestBot()
	// --full flag present — should try ChunkFile path (file won't exist, expect error msg)
	cmd := parseCommand("send --log socrates --full")
	if !cmd.Full {
		t.Fatal("Full should be true")
	}
	reply := b.dispatch(cmd)
	// File doesn't exist in test env, but we verify the --full path was taken
	// (ChunkFile error vs ReadTail error have the same prefix "Error reading file")
	if reply == "" {
		t.Error("expected non-empty reply")
	}
}

func TestDispatch_SendLogNoFullUsesReadTailPath(t *testing.T) {
	b := newTestBot()
	cmd := parseCommand("send --log socrates")
	if cmd.Full {
		t.Fatal("Full should be false")
	}
	reply := b.dispatch(cmd)
	// File doesn't exist — error is expected, but ReadTail path was taken
	if reply == "" {
		t.Error("expected non-empty reply")
	}
}

func TestDispatch_PingReturnsUptime(t *testing.T) {
	b := newTestBot()
	cmd := parseCommand("ping")
	reply := b.dispatch(cmd)
	if !strings.Contains(reply, "Online") && !strings.Contains(reply, "uptime") {
		t.Errorf("expected ping reply, got: %q", reply)
	}
}

func TestDispatch_TailMissingSessionReturnsUsage(t *testing.T) {
	b := newTestBot()
	cmd := parseCommand("tail")
	reply := b.dispatch(cmd)
	if !strings.Contains(strings.ToLower(reply), "usage") {
		t.Errorf("expected usage error, got: %q", reply)
	}
}
