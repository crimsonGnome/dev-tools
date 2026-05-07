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
| `agent-worker` | `agents/agentic-worker.md` | Implement designs from shared-understanding.md — plan, task, execute with evaluation gates |

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
3. Add an alias to `install.sh` and re-run it:

```bash
alias agent-yourname="claude --system-prompt \"$(cat dev-tools/agents/your-agent.md)\""
```

4. Run `source ~/.bashrc`

---

## File Structure

```
dev-tools/
├── agents/          # Agent identities (--system-prompt targets)
├── skills/          # Atomic, reusable skill files
├── steering/        # Per-package CLAUDE.md templates
├── transitions/     # Session handoff docs
├── scripts/         # Workspace and setup scripts
├── docs/            # ADRs and methodology docs
├── steering-guide.md  # Agent and skill design principles
└── install.sh       # Wires everything up

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

agent-worker (implementation — gated by pair-programming evaluation)
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

## Conventions

- **Skills** — single-purpose, reusable, no agent identity
- **Agents** — compose skills via `@filename`, carry a persona and job description
- **Steering** — passive rules in `CLAUDE.md`, always loaded, no workflow logic
- **Context** — `.context.md` is gitignored; never commit session state
- **agentAI/** — I/O contract directory; other scripts expect outputs here; do not move files out
