---
name: agentic-worker-tasks
description: Gate 2 of the agentic-worker pipeline. Reads an approved PLAN.md
  and breaks it into tracer-bullet vertical slice tasks. Each task cuts
  end-to-end through every layer. Tagged AFK or HITL. Signals pair-programming
  for evaluation via context.json. Hard stop if PLAN.md is not approved.
---

## Prerequisite Check (hard stop)
Before doing anything, verify `agentAI/agentic-worker/context.json` shows:
```json
"status": "approved",
"awaiting": "PLAN.md"
```

If evaluation has not passed, stop and say:
> "PLAN.md has not been approved by pair-programming. Switch to
> agent-pair-programming to run evaluation before generating tasks."

---

## Steps

### Step 1 — Read Inputs
Read in order:
1. `agentAI/pair-programming/shared-understanding.md` — authority, testing strategy
2. `agentAI/agentic-worker/PLAN.md` — approved module plan

### Step 2 — Break Into Tracer-Bullet Tasks
For each module in PLAN.md, decompose into vertical slice tasks.

**Tracer-bullet rule:** Every task cuts end-to-end through every layer it touches:
schema → API → UI → tests. Never slice horizontally (e.g. "add all schemas" is
a horizontal slice — wrong. "Add user refresh schema + handler + route + test"
is a vertical slice — correct).

**Sizing rule:** Prefer many thin slices over few thick ones. A task should be
completable and verifiable in one focused session.

**Completion rule:** Every task must be demoable or verifiable per the testing
strategy in shared-understanding.md. If a task cannot be verified, it is too
large — split it.

**Tagging rule:**
- **AFK** — can be fully implemented and merged without human interaction
- **HITL** — requires a design decision, review, or human step to proceed
- Prefer AFK. Only tag HITL when genuinely unavoidable.

**Ordering rule:** Tasks must be ordered by dependency. Foundational slices first.

### Step 3 — Write TASK.md
Write `agentAI/agentic-worker/TASK.md`:

```markdown
## Authority
Source: agentAI/pair-programming/shared-understanding.md
Plan: agentAI/agentic-worker/PLAN.md

## Tasks

### Task [NNN] — [Title]
Module: [module from PLAN.md]
Slice: [schema | API | UI | tests — list all layers this task touches]
Tag: AFK | HITL
HITL reason: [if HITL — exactly what human step is required]
Verification: [how this task is verified demoable/complete per testing strategy]
Status: ⬜ Pending
```

---

## Signal for Evaluation
After writing TASK.md, update `agentAI/agentic-worker/context.json`:
```json
{
  "last_updated": "[date]",
  "current_phase": "tasks",
  "status": "awaiting-evaluation",
  "awaiting": "TASK.md",
  "evaluation": null
}
```

Then tell the user:
> "TASK.md is ready for review. Switch to agent-pair-programming to run evaluation."

---

## Quality Gate (internal — do not print, blocks output)
  - PLAN.md approval confirmed before proceeding
  - Every task is a vertical slice — no horizontal slices
  - Every task has a defined verification method
  - AFK / HITL tagged correctly — HITL reason documented
  - Tasks ordered by dependency
  - No task builds anything in Out of Scope
  - Naming and paths match shared-understanding.md Conventions exactly
  - TASK.md written successfully
  - context.json updated with awaiting-evaluation signal
