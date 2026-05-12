package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/eggersjoseph/discord-orchestrator/internal/config"
	"github.com/eggersjoseph/discord-orchestrator/internal/files"
	"github.com/eggersjoseph/discord-orchestrator/internal/session"
	"github.com/eggersjoseph/discord-orchestrator/internal/state"
)

// ParsedCommand is the result of tokenizing a DM message.
type ParsedCommand struct {
	Verb  string            // "start", "stop", "send", "tail", "inject", etc.
	Args  []string          // positional args after the verb
	Flags map[string]string // "--session" → "name", "--message" → "text", etc.
	Full  bool              // true when "--full" flag is present
}

// parseCommand tokenizes the raw DM input into a ParsedCommand.
func parseCommand(input string) ParsedCommand {
	// Simple quoted-string-aware tokenizer
	tokens := tokenize(input)
	if len(tokens) == 0 {
		return ParsedCommand{Verb: "unknown", Flags: map[string]string{}, Args: []string{}}
	}

	cmd := ParsedCommand{
		Verb:  strings.ToLower(tokens[0]),
		Flags: map[string]string{},
		Args:  []string{},
	}

	for i := 1; i < len(tokens); i++ {
		tok := tokens[i]
		if tok == "--full" {
			cmd.Full = true
			continue
		}
		if strings.HasPrefix(tok, "--") {
			// Consume next token as value if available
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "--") {
				cmd.Flags[tok] = tokens[i+1]
				i++
			} else {
				cmd.Flags[tok] = ""
			}
			continue
		}
		cmd.Args = append(cmd.Args, tok)
	}

	return cmd
}

// tokenize splits input by whitespace, respecting double-quoted strings.
func tokenize(input string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false

	for _, ch := range input {
		switch {
		case ch == '"' && !inQuote:
			inQuote = true
		case ch == '"' && inQuote:
			inQuote = false
		case (ch == ' ' || ch == '\t') && !inQuote:
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(ch)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

// Bot manages the Discord session and dispatches commands.
type Bot struct {
	cfg       *config.Config
	repoRoot  string
	state     *state.State
	statePath string
	dg        *discordgo.Session
	startedAt time.Time
}

// New creates a new Bot instance. Call Start() to connect.
func New(cfg *config.Config, repoRoot string, s *state.State, statePath string) *Bot {
	return &Bot{
		cfg:       cfg,
		repoRoot:  repoRoot,
		state:     s,
		statePath: statePath,
		startedAt: time.Now(),
	}
}

// Start opens the Discord websocket, registers the message handler, posts a
// startup DM to the authorized user, and blocks until the session is closed.
func (b *Bot) Start() error {
	dg, err := discordgo.New("Bot " + b.cfg.DiscordToken)
	if err != nil {
		return fmt.Errorf("bot: create session: %w", err)
	}
	b.dg = dg

	dg.AddHandler(b.handleMessage)
	dg.Identify.Intents = discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	if err := dg.Open(); err != nil {
		return fmt.Errorf("bot: open websocket: %w", err)
	}
	defer dg.Close()

	// Send startup DM
	restored := len(b.state.Sessions)
	b.sendDM(fmt.Sprintf("✓ Orchestrator online. %d sessions restored.", restored))

	// Block until interrupted
	select {}
}

// Stop closes the Discord session.
func (b *Bot) Stop() {
	if b.dg != nil {
		b.dg.Close()
	}
}

// handleMessage is the Discord message handler. It filters by authorized user
// and dispatches commands. Every code path produces exactly one reply.
func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Silently ignore unauthorized users
	if m.Author.ID != b.cfg.AuthorizedUserID {
		return
	}
	// Only process DMs
	ch, err := s.Channel(m.ChannelID)
	if err != nil || ch.Type != discordgo.ChannelTypeDM {
		return
	}

	cmd := parseCommand(m.Content)
	reply := b.dispatch(cmd)
	s.ChannelMessageSend(m.ChannelID, reply)
}

// sendDM sends a DM to the authorized user.
func (b *Bot) sendDM(msg string) {
	if b.dg == nil {
		return
	}
	ch, err := b.dg.UserChannelCreate(b.cfg.AuthorizedUserID)
	if err != nil {
		return
	}
	b.dg.ChannelMessageSend(ch.ID, msg)
}

// dispatch routes a ParsedCommand to the correct handler and returns the reply string.
// Every case returns exactly one string.
func (b *Bot) dispatch(cmd ParsedCommand) string {
	switch cmd.Verb {
	case "start":
		return b.cmdStart(cmd)
	case "stop":
		return b.cmdStop(cmd)
	case "restart":
		return b.cmdRestart(cmd)
	case "list":
		return b.cmdList()
	case "status":
		return b.cmdStatus(cmd)
	case "inject":
		return b.cmdInject(cmd)
	case "send":
		return b.cmdSend(cmd)
	case "tail":
		return b.cmdTail(cmd)
	case "ping":
		return b.cmdPing()
	case "reload":
		return b.cmdReload()
	case "help", "unknown":
		return helpText()
	default:
		return helpText()
	}
}

func (b *Bot) cmdStart(cmd ParsedCommand) string {
	if len(cmd.Args) == 0 {
		return "Usage: start <name>"
	}
	name := cmd.Args[0]
	def, ok := b.cfg.Sessions[name]
	if !ok {
		return fmt.Sprintf("Unknown session %q. Known sessions: %s", name, b.knownSessions())
	}
	if err := session.Start(name, b.repoRoot, def); err != nil {
		return fmt.Sprintf("Error starting %q: %v", name, err)
	}
	b.state.Sessions[name] = state.SessionState{Status: "running", Started: time.Now().Format(time.RFC3339)}
	state.Save(b.statePath, b.state)
	return fmt.Sprintf("✓ Started session %q", name)
}

func (b *Bot) cmdStop(cmd ParsedCommand) string {
	if len(cmd.Args) == 0 {
		return "Usage: stop <name>"
	}
	name := cmd.Args[0]
	if err := session.Stop(name); err != nil {
		return fmt.Sprintf("Error stopping %q: %v", name, err)
	}
	if ss, ok := b.state.Sessions[name]; ok {
		ss.Status = "stopped"
		ss.PID = 0
		b.state.Sessions[name] = ss
		state.Save(b.statePath, b.state)
	}
	return fmt.Sprintf("✓ Stopped session %q", name)
}

func (b *Bot) cmdRestart(cmd ParsedCommand) string {
	stopReply := b.cmdStop(cmd)
	startReply := b.cmdStart(cmd)
	return stopReply + "\n" + startReply
}

func (b *Bot) cmdList() string {
	names, err := session.List()
	if err != nil {
		return fmt.Sprintf("Error listing sessions: %v", err)
	}
	if len(names) == 0 {
		return "No active sessions."
	}
	var lines []string
	for _, name := range names {
		ss, ok := b.state.Sessions[name]
		if !ok {
			lines = append(lines, fmt.Sprintf("• %s (unknown)", name))
		} else {
			lines = append(lines, fmt.Sprintf("• %s — %s", name, ss.Status))
		}
	}
	return strings.Join(lines, "\n")
}

func (b *Bot) cmdStatus(cmd ParsedCommand) string {
	if len(cmd.Args) == 0 {
		return "Usage: status <name>"
	}
	name := cmd.Args[0]
	status, err := session.Status(name)
	if err != nil {
		return fmt.Sprintf("Error getting status for %q: %v", name, err)
	}
	return fmt.Sprintf("Session %q: %s", name, status)
}

func (b *Bot) cmdInject(cmd ParsedCommand) string {
	name, ok := cmd.Flags["--session"]
	if !ok || name == "" {
		return "Usage: inject --session <name> --message \"...\" | --file <path>"
	}

	ss, exists := b.state.Sessions[name]
	currentStatus := "running"
	if exists {
		currentStatus = ss.Status
	}

	var message string
	if msg, ok := cmd.Flags["--message"]; ok {
		message = msg
	} else if filePath, ok := cmd.Flags["--file"]; ok {
		data, _, err := files.ReadTail(filePath, files.DefaultTailBytes)
		if err != nil {
			return fmt.Sprintf("Error reading file %q: %v", filePath, err)
		}
		message = string(data)
	} else {
		return "Usage: inject --session <name> --message \"...\" | --file <path>"
	}

	if err := session.Inject(name, message, currentStatus); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return fmt.Sprintf("✓ Injected into session %q", name)
}

func (b *Bot) cmdSend(cmd ParsedCommand) string {
	if filePath, ok := cmd.Flags["--file"]; ok {
		return b.sendFileOrTail(filePath, cmd.Full)
	}
	if logName, ok := cmd.Flags["--log"]; ok {
		logPath := session.LogPath(b.repoRoot, logName)
		return b.sendFileOrTail(logPath, cmd.Full)
	}
	return "Usage: send --file <path> [--full] | send --log <name> [--full]"
}

func (b *Bot) sendFileOrTail(path string, full bool) string {
	if full {
		chunks, err := files.ChunkFile(path, files.ChunkSize)
		if err != nil {
			return fmt.Sprintf("Error reading file: %v", err)
		}
		// In production, each chunk would be posted as a Discord attachment.
		// Here we return a summary — the Discord posting happens in handleMessage
		// via the session. For now, return first chunk inline (Discord limit is 2000 chars).
		if len(chunks) == 1 {
			return string(chunks[0])
		}
		return fmt.Sprintf("[Sending %d chunks — first chunk]\n%s", len(chunks), string(chunks[0]))
	}

	data, truncated, err := files.ReadTail(path, files.DefaultTailBytes)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	result := string(data)
	if truncated {
		result = "[Output truncated — use --full for complete file]\n" + result
	}
	return result
}

func (b *Bot) cmdTail(cmd ParsedCommand) string {
	name, ok := cmd.Flags["--session"]
	if !ok || name == "" {
		return "Usage: tail --session <name> [--lines N]"
	}
	lines := files.DefaultTailLines
	if linesStr, ok := cmd.Flags["--lines"]; ok {
		n, err := strconv.Atoi(linesStr)
		if err != nil {
			return fmt.Sprintf("Invalid --lines value %q: must be an integer", linesStr)
		}
		lines = n
	}
	logPath := session.LogPath(b.repoRoot, name)
	result, err := files.Tail(logPath, lines)
	if err != nil {
		return fmt.Sprintf("Error reading log for %q: %v", name, err)
	}
	if result == "" {
		return fmt.Sprintf("Log for %q is empty.", name)
	}
	return result
}

func (b *Bot) cmdPing() string {
	uptime := time.Since(b.startedAt).Round(time.Second)
	running := 0
	for _, ss := range b.state.Sessions {
		if ss.Status == "running" {
			running++
		}
	}
	return fmt.Sprintf("✓ Online — uptime %s — %d session(s) running", uptime, running)
}

func (b *Bot) cmdReload() string {
	s, err := state.Load(b.statePath)
	if err != nil {
		return fmt.Sprintf("Error reloading state: %v", err)
	}
	b.state = s
	return "✓ State reloaded from disk."
}

func (b *Bot) knownSessions() string {
	var names []string
	for name := range b.cfg.Sessions {
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

func helpText() string {
	return strings.TrimSpace(`
Commands:
  start <name>
  stop <name>
  restart <name>
  list
  status <name>
  inject --session <name> --message "..."
  inject --session <name> --file <path>
  send --file <path> [--full]
  send --log <name> [--full]
  tail --session <name> [--lines N]
  ping
  reload
  help`)
}
