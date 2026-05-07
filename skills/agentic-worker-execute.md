---
name: agentic-worker-execute
description: Phase 3 of the agentic-worker pipeline. Executes approved tasks
  from TASK.md one at a time. AFK tasks run to completion autonomously. HITL
  tasks pause execution, surface the required human step, and resume after
  input. Updates TASK.md in place as the live execution record.
---

## Prerequisite Check (hard stop)
Before doing anything, verify `agentAI/agentic-worker/context.json` shows:
```json
"status": "approved",
"awaiting": "TASK.md"
```

If evaluation has not passed, stop and say:
> "TASK.md has not been approved by pair-programming. Switch to
> agent-pair-programming to run evaluation before executing."

---

## Authority Rule
Read `agentAI/pair-programming/shared-understanding.md` before executing any task.
Conventions, interface boundaries, and testing strategy are binding throughout execution.
Do not infer. If a decision is not in shared-understanding.md, stop and surface it
as a HITL blocker — do not guess.

---

## Execution Loop

### For each task in TASK.md (in order):

**Update status to In Progress:**
```
Status: 🔄 In Progress
```

**If AFK:**
1. Implement the vertical slice (all layers the task touches)
2. Run verification:
   - Tests pass for this slice
   - No regressions in adjacent passing tests
   - Behavior is observable/demoable
3. If verification passes → mark complete:
   ```
   Status: ✅ Complete
   Verified: [what was verified — test output, endpoint response, rendered component, etc.]
   ```
4. If no tests exist for a layer this task touches → do not mark complete:
   ```
   Status: ❌ Failed
   Reason: No tests exist for [layer]. Gap must be addressed before this task can complete.
   ```
5. Proceed to next task.

**If HITL:**
1. Stop execution immediately.
2. Update TASK.md:
   ```
   Status: ⏸ Blocked (HITL)
   Blocker: [exact human step required — be specific]
   ```
3. Update `agentAI/agentic-worker/context.json`:
   ```json
   {
     "status": "blocked",
     "blocked_task": "Task [NNN]",
     "blocker": "[exact human step required]"
   }
   ```
4. Tell the user clearly:
   > "Execution paused at Task [NNN]. Required: [exact human step].
   > Respond with your decision and I will resume."
5. After human responds → implement → run verification → mark complete → resume.

---

## Resuming After Block
On session start, read `agentAI/agentic-worker/context.json`. If status is
`"blocked"`, resume from the blocked task after the human step is resolved.

---

## Completion
When all tasks are ✅ Complete or ❌ Failed (with documented gaps):

Update `agentAI/agentic-worker/context.json`:
```json
{
  "status": "complete",
  "current_phase": "execute",
  "blocked_tasks": ["Task NNN — reason"],
  "failed_tasks": ["Task NNN — reason"]
}
```

Tell the user:
> "Execution complete. [N] tasks completed. [N] blocked. [N] failed.
> See TASK.md for full record."

---

## Quality Gate (internal — do not print, blocks output)
  - shared-understanding.md read before first task
  - Every task verified before marking complete — no self-certification
  - HITL tasks never skipped or guessed through
  - Missing test coverage flagged as failure — not skipped
  - TASK.md updated in place throughout (single live record)
  - context.json updated after every status change
