package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// SessionDef holds the command to run for a named session.
type SessionDef struct {
	Cmd string `json:"cmd"`
}

// Config is the loaded and validated configuration for the orchestrator.
type Config struct {
	DiscordToken     string                `json:"discord_token"`
	AuthorizedUserID string                `json:"authorized_user_id"`
	Sessions         map[string]SessionDef `json:"sessions"`
}

// Load reads and validates the config file at the given path.
// Returns a clear error if any required field is missing.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: cannot open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields() // surface typos, but we catch below
	if err := dec.Decode(&cfg); err != nil {
		// Unknown-field errors are non-fatal per spec — retry without strict mode
		f.Seek(0, 0)
		var cfg2 Config
		if err2 := json.NewDecoder(f).Decode(&cfg2); err2 != nil {
			return nil, fmt.Errorf("config: malformed JSON in %q: %w", path, err2)
		}
		cfg = cfg2
	}

	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("config: missing required field \"discord_token\"")
	}
	if cfg.AuthorizedUserID == "" {
		return nil, fmt.Errorf("config: missing required field \"authorized_user_id\"")
	}

	return &cfg, nil
}
