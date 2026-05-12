# dev-tools

Agent-switching and skill composition system for the Vigilon monorepo.

---

## Setup

Run once from the repo root:

```bash
./dev-tools/install.sh
source ~/.bashrc
```

This will:
- Symlink steering docs (`CLAUDE.md`) into each package directory
- Register agent aliases in your shell

---

## Running an Agent

After setup, launch an agent by typing its alias in your terminal:

```bash
agent-pair-programming
```

This opens a Claude Code session with that agent's identity and skills pre-loaded.

### Available Agents

| Alias | File | Purpose |
|---|---|---|
| `agent-pair-programming` | `agents/planning-agent.md` | Stress-test plans, resolve architecture decisions, write shared-understanding.md. Evaluates agentic-worker artifacts. |
| `agent-socrates` | `agents/socrates.md` | Index packages, teach through Socratic questioning, scope problems, propose solutions |
| `agent-agentic-worker` | `agents/agentic-worker.md` | Implement designs from shared-understanding.md — plan, task, execute with evaluation gates |

---

## Switching Agents

Agents have clean context boundaries. To switch:

1. Run the transition-doc skill at the end of your session — it writes `.context.md`
2. Quit the session (`/exit` or Ctrl+C)
3. Launch the next agent alias

The next agent picks up `.context.md` automatically at launch.

---

## Adding a New Agent

1. Create a skill file in `dev-tools/skills/your-skill.md`
2. Create an agent file in `dev-tools/agents/your-agent.md` that composes skills via `@`
3. Add a shell function to `install.sh` and re-run it:

```bash
agent-yourname() { claude --system-prompt "$(cat dev-tools/agents/your-agent.md)" "$@"; }
```

4. Run `source ~/.bashrc`
5. Add the session to `dev-tools/discord-orchestrator/config.json` if you want the orchestrator to manage it

---

## File Structure

```
dev-tools/
├── agents/                    # Agent identities (--system-prompt targets)
├── skills/                    # Atomic, reusable skill files
├── steering/                  # Per-package CLAUDE.md templates
├── transitions/               # Session handoff docs
├── scripts/                   # Workspace and setup scripts
├── docs/                      # ADRs and methodology docs
├── discord-orchestrator/      # Discord bot daemon — manages agent tmux sessions
│   ├── main.go
│   ├── go.mod / go.sum
│   ├── config.example.json    # Copy to config.json and fill in token + user ID
│   ├── discord-orchestrator.service  # Systemd unit template
│   ├── .gitignore             # Excludes config.json, logs/, bin/
│   └── internal/
│       ├── bot/               # Discord session, command dispatch, flag parsing
│       ├── config/            # Config loading + validation
│       ├── files/             # Log tailing, file chunking
│       ├── session/           # tmux lifecycle (start/stop/inject/status/list)
│       └── state/             # State persistence (atomic JSON)
├── steering-guide.md          # Agent and skill design principles
└── install.sh                 # Wires everything up + builds/installs orchestrator

agentAI/             # Agent I/O directory — all agent outputs live here
├── socrates/
│   ├── context.json              # State manifest (last run status)
│   ├── high-level-architecture.md  # Output of Index (high-level mode)
│   ├── scoped-architecture.md    # Output of Focus skill
│   ├── proposals.md              # Output of Propose skill
│   └── session-log.md            # Output of Socratic skill
├── pair-programming/
│   ├── context.json              # State manifest + evaluation queue
│   ├── shared-understanding.md   # Authority doc — binding for agentic-worker
│   ├── questions-log.md          # Q&A transcript — pair-programming only
│   └── evaluation.md             # PASS/FAIL verdict for PLAN.md or TASK.md
└── agentic-worker/
    ├── context.json              # Phase + gate state + handoff signal
    ├── PLAN.md                   # Module plan (Gate 1 output)
    └── TASK.md                   # Vertical slice tasks (Gate 2 output, live execution record)
```

### Key Files

| File | Purpose |
|---|---|
| `.context.md` | Active session handoff — read at start, written at end |
| `CLAUDE.md` (root) | Passive rules loaded whenever Claude opens in root |
| `*/CLAUDE.md` (packages) | Symlinked steering, loaded per-package |
| `transitions/current.md` | Latest handoff doc (mirrors `.context.md`) |
| `transitions/template.md` | Blank template for new handoff docs |
| `dev-tools/steering-guide.md` | Design principles for building agents and skills |
| `agentAI/<agent>/context.json` | State manifest — tracks skill run status per agent |

---

## Agent Pipeline

The full system flows across three agents. Each has hard-stop gates.

```
agent-socrates (optional — provides architecture context)
  └── Index   → agentAI/socrates/high-level-architecture.md
        └── Focus   → agentAI/socrates/scoped-architecture.md
              └── Propose → agentAI/socrates/proposals.md
  └── Socratic → agentAI/socrates/session-log.md

agent-pair-programming (design authority)
  └── pair-programming → agentAI/pair-programming/shared-understanding.md
                       → agentAI/pair-programming/questions-log.md
  └── evaluate         → agentAI/pair-programming/evaluation.md (PASS | FAIL)

agent-agentic-worker (implementation — gated by pair-programming evaluation)
  └── Gate 1: PLAN.md  ──► agent-pair-programming evaluates ──► pass/fail loop
  └── Gate 2: TASK.md  ──► agent-pair-programming evaluates ──► pass/fail loop
  └── Phase 3: Execute
        AFK tasks → autonomous implementation + verification
        HITL tasks → pause, surface blocker, wait, resume
```

### Socrates Skill Chain

```
Index → Focus → Propose  (hard stop at each gate)
Socratic runs independently
```

---

## Discord Orchestrator

A persistent Go daemon that lets you manage agent sessions from Discord DMs. Always running as a systemd user service — survives reboots.

### First-time setup

```bash
cp dev-tools/discord-orchestrator/config.example.json dev-tools/discord-orchestrator/config.json
# Edit config.json — fill in discord_token and authorized_user_id
source dev-tools/install.sh
```

`install.sh` builds the binary and installs the systemd service automatically.

### Discord commands

```
start <name>                         Start a named agent session
stop <name>                          Kill a session
restart <name>                       Stop then start
list                                 List all sessions and statuses
status <name>                        Single session status
inject --session <name> --message "…"   Send a message into a session
inject --session <name> --file <path>   Inject file contents into a session
send --file <path>                   Send file tail (last 50KB)
send --file <path> --full            Send entire file in chunks
send --log <name>                    Send session log tail (last 50KB)
send --log <name> --full             Send full session log in chunks
tail --session <name> [--lines N]    Send last N lines of session log (default 20)
ping                                 Uptime + running session count
reload                               Reload state from disk
help                                 Print command list
```

### Known sessions (pre-configured)

| Session name | Agent |
|---|---|
| `socrates` | Socrates teaching agent |
| `pair-programming` | Pair-programming / design agent |
| `agentic-worker` | Agentic worker / implementation agent |

Add new sessions by editing `config.json` — no code changes needed.

---

## Conventions

- **Skills** — single-purpose, reusable, no agent identity
- **Agents** — compose skills via `@filename`, carry a persona and job description
- **Steering** — passive rules in `CLAUDE.md`, always loaded, no workflow logic
- **Context** — `.context.md` is gitignored; never commit session state
- **agentAI/** — I/O contract directory; other scripts expect outputs here; do not move files out
