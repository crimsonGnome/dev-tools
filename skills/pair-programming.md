---
name: pair-programming
description: Interview the user relentlessly about a plan or design until reaching
  shared understanding, resolving each branch of the decision tree. Writes
  shared-understanding.md (the authority doc) and questions-log.md (the Q&A
  transcript). Use when user wants to stress-test a plan, get grilled on their
  design, or mentions "grill me".
---

Interview me relentlessly about every aspect of this plan until we reach shared
understanding. Walk down each branch of the design tree, resolving dependencies
between decisions one-by-one. For each question, provide your recommended answer.

Ask the questions one at a time.

If a question can be answered by exploring the codebase, explore the codebase
instead of asking.

---

## Output

At the end of the session, write two files:

### 1. agentAI/pair-programming/shared-understanding.md
The authority document. The agentic-worker follows this exactly — no inference.

Structure:
```
## Problem Statement
[One clear sentence — what are we solving]

## Conventions
[File paths, naming, structure decided this session — the agentic-worker
follows these exactly. Listed before everything else so they cannot be missed.]

## Design Decisions
| Decision | Why This Approach | Alternatives Rejected |

## Interface Boundaries
[What each module exposes and consumes. File paths authoritative.]

## Testing Strategy
[How each layer is tested — unit, integration, e2e]
[What "demoable or verifiable" means for this feature]

## Risk Mitigations
[Known risks and how the design addresses them]

## Out of Scope
[Explicit list — if it's not here, the agentic-worker cannot build it]
```

### 2. agentAI/pair-programming/questions-log.md
The full Q&A transcript from this session.

Structure:
```
## Session
Date: [date]
Topic: [problem statement]

## Questions & Answers
Q: [question asked]
Recommendation: [your recommended answer]
A: [user's answer]
Resolution: [decision reached]

[repeat for each exchange]

## Decision Summary
[Bullet list of every decision made, in order]
```

**Important:** questions-log.md is for pair-programming context only.
Do NOT include it in transition docs. Do NOT pass it to the agentic-worker.
It is a session record, not an authority document.

---

## Quality Gate (internal — do not print, blocks output)
  - Problem statement captured and clear
  - All major decision branches resolved
  - Conventions section populated with at least file paths and naming
  - Testing strategy explicitly defined
  - Out of scope explicitly listed (cannot be empty)
  - shared-understanding.md written successfully
  - questions-log.md written successfully
  - context.json updated

---

## Update context.json
After writing both files, update `agentAI/pair-programming/context.json`:
```json
{
  "last_updated": "[date]",
  "topic": "[problem statement]",
  "status": "complete",
  "shared_understanding": "agentAI/pair-programming/shared-understanding.md",
  "questions_log": "agentAI/pair-programming/questions-log.md"
}
```
