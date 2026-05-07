---
name: socrates-focus
description: Zoom in on a problem statement. Infers relevant components from the
  high-level architecture, runs a scoped index, flags missing boundary logging,
  and writes scoped-architecture.md. Requires high-level-architecture.md to
  exist first — hard stop if missing.
---

## Prerequisite Check (hard stop)
Before doing anything, verify `agentAI/socrates/high-level-architecture.md` exists.

If it does not exist, stop immediately and say:
> "Index has not been run. Run the Index skill first to generate
> high-level-architecture.md before scoping."

Do not proceed until this file exists.

---

## Steps

### Step 1 — Problem Statement
Capture the user's problem statement. If not already provided, ask:
> "What problem are you trying to solve?"

### Step 2 — Component Inference
Read `agentAI/socrates/high-level-architecture.md`. Map the problem statement
against the documented flows and components. Identify which areas of the codebase
are most likely involved. Proceed directly — no confirmation needed.

### Step 3 — Scoped Index
Run the Index skill in scoped mode on the identified components. This populates
the scoped view with flows, key interfaces, and dependency maps for the relevant
area only.

### Step 4 — Boundary Logging Audit
Within the scoped components, identify every API call or external service call.
Flag any that lack failure logging. List them explicitly in the output — do not
skip this step even if logging appears complete.

### Step 5 — Write scoped-architecture.md
Write `agentAI/socrates/scoped-architecture.md`. Structure mirrors
`high-level-architecture.md` but narrowed to the scoped components.

Include:
- Problem statement
- Scoped component map (with one-line descriptions)
- Relevant flows (each traceable back to a flow in high-level-architecture.md)
- Key interfaces within scope (file path + line number)
- Missing boundary logging — list each call site
- Diagrams where applicable (ASCII or Mermaid)
- 2-3 sentence summary

---

## Quality Gate (internal — do not print, blocks output)
  - `high-level-architecture.md` confirmed to exist
  - Problem statement captured
  - Components inferred from high-level architecture
  - Scoped flow traceable back to a flow in high-level-architecture.md
  - Boundary logging gaps explicitly surfaced (even if none found — state that)
  - At least one diagram included
  - 2-3 sentence summary written
  - `scoped-architecture.md` written successfully

---

## Output
`agentAI/socrates/scoped-architecture.md`

After writing, update `agentAI/socrates/context.json`:
```json
"focus": {
  "status": "complete",
  "problem_statement": "...",
  "scoped_components": ["src/auth/session.go", "..."],
  "output": "agentAI/socrates/scoped-architecture.md"
}
```
