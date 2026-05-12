package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config-*.json")
	if err != nil {
		t.Fatalf("createTemp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `{
		"discord_token": "tok123",
		"authorized_user_id": "uid456",
		"sessions": {
			"socrates": {"cmd": "agent-socrates"}
		}
	}`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DiscordToken != "tok123" {
		t.Errorf("DiscordToken = %q, want %q", cfg.DiscordToken, "tok123")
	}
	if cfg.AuthorizedUserID != "uid456" {
		t.Errorf("AuthorizedUserID = %q, want %q", cfg.AuthorizedUserID, "uid456")
	}
	if def, ok := cfg.Sessions["socrates"]; !ok || def.Cmd != "agent-socrates" {
		t.Errorf("Sessions[socrates] = %+v", cfg.Sessions["socrates"])
	}
}

func TestLoad_MissingDiscordToken(t *testing.T) {
	path := writeTemp(t, `{
		"authorized_user_id": "uid456"
	}`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing discord_token, got nil")
	}
}

func TestLoad_MissingAuthorizedUserID(t *testing.T) {
	path := writeTemp(t, `{
		"discord_token": "tok123"
	}`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing authorized_user_id, got nil")
	}
}

func TestLoad_UnknownFieldsIgnored(t *testing.T) {
	path := writeTemp(t, `{
		"discord_token": "tok123",
		"authorized_user_id": "uid456",
		"unknown_field": "should_be_ignored"
	}`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error for unknown field: %v", err)
	}
	if cfg.DiscordToken != "tok123" {
		t.Errorf("DiscordToken = %q", cfg.DiscordToken)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
