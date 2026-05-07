---
name: socrates-index
description: Index a package structure, key components, end-to-end flows, and
  dependencies. Dual-mode — high-level indexes the full package, scoped zooms
  into a specific component. High-level writes high-level-architecture.md.
  Scoped output is consumed by the Focus skill.
---

## Mode
- **No argument** → high-level mode → writes `agentAI/socrates/high-level-architecture.md`
- **Scope argument** (e.g. `src/auth`) → scoped mode → output consumed by Focus skill

---

## Indexing Pipeline

### Step 1 — Folder Map
Read the top-level directory structure. List every top-level folder with a
one-line description of its purpose.

### Step 2 — Project Shape
Read: README.md, docs/ directory, and all config files (package.json, Makefile,
go.mod, tsconfig.json, .env.example, etc.). Understand the project's purpose,
dependency graph, and build system before reading any source files.

### Step 3 — Key Source Files
For each important top-level directory, read key source files. Balance depth
across the whole package — do not exhaust context on one module. Prioritize
entry points, interfaces, and files referenced most by other modules.

### Step 4 — Key Interfaces
Identify key types, interfaces, and functions. Record each with:
- File path
- Line number
- One-line description of its role

### Step 5 — End-to-End Flow
Trace at least one complete end-to-end flow through the system from entry point
to output. Map dependencies:
- What this package consumes (external services, libraries, other packages)
- What consumes this package (callers, dependents)

---

## Diagrams
Where applicable, include ASCII or Mermaid diagrams. Diagrams serve two purposes:
1. Help users understand how the architecture is structured
2. Give AI agents a quick visual reference for flows and relationships

Every output file must include at least one diagram.

---

## Scoped Mode — Additional Steps
When running in scoped mode on a specific component:
- Identify all connection points between this component and the rest of the system
- Ensure the scoped flow can be traced back to a flow documented in
  `high-level-architecture.md`
- Flag any API calls or external service calls within scope that lack failure
  logging — list them explicitly

---

## Quality Gate (internal — do not print, blocks output)
Socrates must verify all criteria before writing the output file.

**Structural**
  - Every top-level folder has a one-line description
  - README, config files, and docs/ have been read
  - Key source files read across all important directories
  - No single module consumed disproportionate context

**Content**
  - Key types, interfaces, and functions identified with file path + line number
  - At least one end-to-end flow traced start to finish
  - Dependencies mapped (inbound and outbound)
  - At least one ASCII or Mermaid diagram included

**Summary**
  - 2-3 sentence summary written
  - Output file written successfully

**Scoped mode additionally**
  - Component connection points to the rest of the system identified
  - Scoped flow traceable back to a flow in high-level-architecture.md
  - Missing boundary logging flagged explicitly

---

## Output
- High-level mode: `agentAI/socrates/high-level-architecture.md`
- Scoped mode: returned to Focus skill (not written independently)

After writing, update `agentAI/socrates/context.json`:
```json
"index": {
  "status": "complete",
  "mode": "high-level | scoped"
}
```
