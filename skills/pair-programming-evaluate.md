---
name: pair-programming-evaluate
description: Evaluation mode for the pair-programming agent. Reads context.json
  to detect a pending evaluation, then audits PLAN.md or TASK.md against
  shared-understanding.md. Writes evaluation.md with a PASS or FAIL verdict
  and actionable required fixes.
---

## Trigger
Run automatically when `agentAI/agentic-worker/context.json` contains:
```json
"status": "awaiting-evaluation"
```

Read the `"awaiting"` field to know what to evaluate: `"PLAN.md"` or `"TASK.md"`.

---

## Inputs
- `agentAI/pair-programming/shared-understanding.md` — the authority
- `agentAI/socrates/high-level-architecture.md` — architecture reference (if exists)
- `agentAI/socrates/scoped-architecture.md` — scoped reference (if exists)
- `agentAI/agentic-worker/PLAN.md` or `agentAI/agentic-worker/TASK.md`

Do NOT read `questions-log.md` — it is not an authority document.

---

## Evaluation Checks

### When evaluating PLAN.md
  - Every module traces back to a decision in shared-understanding.md
  - File paths and naming match Conventions exactly — no inference
  - No module builds something listed in Out of Scope
  - Interface boundaries match shared-understanding.md
  - Testing strategy is reflected in the plan
  - Deep module pattern applied where applicable

### When evaluating TASK.md
  - Every task is a vertical slice (schema → API → UI → tests) — not horizontal
  - Every task is demoable or verifiable per the testing strategy
  - AFK / HITL tags are correct
  - No task builds something Out of Scope
  - Naming and paths match Conventions exactly
  - Tasks are ordered by dependency (foundational slices first)

---

## Output: agentAI/pair-programming/evaluation.md

```markdown
## Evaluation Target
[PLAN.md | TASK.md]

## Verdict
[PASS | FAIL]

## Checks
✦ Aligns with shared-understanding.md decisions
✦ Follows Conventions exactly (no inference)
✦ Testing strategy honored
✦ Interface boundaries respected
✦ Nothing out of scope has crept in
[Additional checks specific to PLAN or TASK]

## Failures
[If PASS — "None"]
[If FAIL — one item per line:]
- [Specific item] — [what's wrong] — [what shared-understanding.md says]

## Required Fixes
[If PASS — "None. Proceed to next phase."]
[If FAIL — concrete actionable fixes the agentic-worker must make:]
- [Fix 1]
- [Fix 2]
```

---

## After Writing evaluation.md
Update `agentAI/pair-programming/context.json`:
```json
{
  "last_evaluation": "[date]",
  "evaluated": "PLAN.md | TASK.md",
  "verdict": "PASS | FAIL"
}
```

Update `agentAI/agentic-worker/context.json`:
```json
{
  "status": "approved | revision-required",
  "evaluation": "agentAI/pair-programming/evaluation.md"
}
```

---

## Quality Gate (internal — do not print, blocks output)
  - shared-understanding.md read before evaluating anything
  - Every check verified against shared-understanding.md, not inferred
  - Failures cite the specific line or section of shared-understanding.md
  - Required fixes are concrete and actionable — not vague
  - evaluation.md written successfully
  - Both context.json files updated
