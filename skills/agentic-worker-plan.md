---
name: agentic-worker-plan
description: Gate 1 of the agentic-worker pipeline. Reads shared-understanding.md
  as the authority, explores the repo, sketches major modules using the deep
  module pattern, and writes PLAN.md. Signals pair-programming for evaluation
  via context.json. Hard stop if shared-understanding.md is missing.
---

## Authority Rule
shared-understanding.md is authoritative. It specifies files, paths, naming
conventions, and structure. Follow exactly. Do not infer.

---

## Prerequisite Check (hard stop)
Before doing anything, verify `agentAI/pair-programming/shared-understanding.md` exists.

If it does not exist, stop immediately and say:
> "shared-understanding.md is missing. Run agent-pair-programming first to
> design the solution before planning implementation."

---

## Steps

### Step 1 — Read Authority
Read `agentAI/pair-programming/shared-understanding.md` in full.
Extract and hold:
- Conventions (file paths, naming, structure — these are binding)
- Design Decisions
- Interface Boundaries
- Testing Strategy
- Out of Scope

### Step 2 — Read Architecture (soft)
If they exist, read:
- `agentAI/socrates/high-level-architecture.md`
- `agentAI/socrates/scoped-architecture.md`

If missing, proceed. Do not block on these.

### Step 3 — Explore the Repo
Explore the codebase enough to ground implementation decisions. Focus on:
- Entry points and existing patterns relevant to the plan
- Existing base classes or interfaces to extend (not duplicate)
- File structure of areas the plan will touch

### Step 4 — Sketch Modules
Identify the major modules to build or modify. For each, apply the deep module
pattern: significant functionality behind a simple, testable interface.

Prefer modules that:
- Hide complexity behind a clean interface
- Can be tested in isolation
- Have a single clear responsibility

### Step 5 — Write PLAN.md
Write `agentAI/agentic-worker/PLAN.md`:

```markdown
## Authority
Source: agentAI/pair-programming/shared-understanding.md
Conventions: [list conventions pulled from shared-understanding]

## Modules

### [Module Name]
Action: Build | Modify
Path: [exact path from shared-understanding — no inference]
Responsibility: [one line]
Interface: [inputs / outputs]
Deep module rationale: [why this has significant functionality behind a simple surface]

[repeat for each module]

## Dependency Order
[Modules listed in build order — which must exist before others]

## Out of Scope
[Pulled verbatim from shared-understanding.md]
```

---

## Signal for Evaluation
After writing PLAN.md, update `agentAI/agentic-worker/context.json`:
```json
{
  "last_updated": "[date]",
  "current_phase": "plan",
  "status": "awaiting-evaluation",
  "awaiting": "PLAN.md",
  "evaluation": null
}
```

Then tell the user:
> "PLAN.md is ready for review. Switch to agent-pair-programming to run evaluation."

---

## Quality Gate (internal — do not print, blocks output)
  - shared-understanding.md confirmed to exist and fully read
  - Conventions extracted and applied — no inference on naming or paths
  - Every module traces to a decision in shared-understanding.md
  - No module builds anything in Out of Scope
  - Deep module pattern applied and documented
  - Dependency order specified
  - PLAN.md written successfully
  - context.json updated with awaiting-evaluation signal
