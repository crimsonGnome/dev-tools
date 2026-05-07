---
name: socrates-propose
description: Propose 3-5 structured solutions to the problem statement scoped by
  the Focus skill. Each proposal includes approach, tradeoffs, complexity,
  files touched, logging gaps, and a recommendation flag. Requires
  scoped-architecture.md to exist first — hard stop if missing.
---

## Prerequisite Check (hard stop)
Before doing anything, verify `agentAI/socrates/scoped-architecture.md` exists.

If it does not exist, stop immediately and say:
> "Focus has not been run. Run the Focus skill first to generate
> scoped-architecture.md before proposing solutions."

Do not proceed until this file exists.

---

## Steps

### Step 1 — Read Context
Read `agentAI/socrates/scoped-architecture.md`. Extract:
- Problem statement
- Scoped components and their roles
- Relevant flows
- Key interfaces
- Missing boundary logging gaps

### Step 2 — Generate Proposals
Generate 3-5 distinct proposals. Proposals should represent meaningfully
different approaches — not variations of the same solution.

Each proposal must contain:

| Field | Description |
|-------|-------------|
| **Title** | One-line name for the approach |
| **Approach** | 2-3 sentences describing the solution |
| **Tradeoffs** | Explicit pros and cons |
| **Complexity** | Low / Medium / High |
| **Files likely touched** | Drawn from scoped-architecture.md |
| **Logging gaps to address** | From the Focus boundary audit |
| **Recommended** | yes / no + one-line reason |

Mark exactly one proposal as Socrates' recommendation.

### Step 3 — Write proposals.md
Write to `agentAI/socrates/proposals.md`.

---

## Quality Gate (internal — do not print, blocks output)
  - `scoped-architecture.md` confirmed to exist
  - Minimum 3 proposals generated
  - Every proposal contains all required fields
  - Exactly one proposal marked as recommended with a reason
  - Logging gaps addressed in at least one proposal
  - `proposals.md` written successfully

---

## Output
`agentAI/socrates/proposals.md`

After writing, update `agentAI/socrates/context.json`:
```json
"propose": {
  "status": "complete",
  "recommendation": "Title of recommended proposal",
  "output": "agentAI/socrates/proposals.md"
}
```
