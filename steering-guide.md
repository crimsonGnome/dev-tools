# Steering Guide

Rules and principles for building agents and skills in this system.

---

## Principle 1: Extend, Don't Duplicate
Skills should share a base structure. When building new skills, identify the
common contract (inputs, quality gate, output format) and define it once. Extend
for specifics. Avoid copy-pasting logic across skills — if two skills share
behavior, that behavior belongs in a base definition or shared principle, not
duplicated. Think in base classes and extensions.

## Principle 2: Quality Gates Are Agent-Specific
Every agent defines its own base quality gate appropriate to its domain. Quality
gates are internal completion criteria — the agent checks them before writing
output but does not print them. They are not universal. A teaching agent's gate
(diagram included, summary written) does not apply to a planning agent. Define
the gate for the domain, then extend per skill.

## Principle 3: Log All External Boundaries
Every call to an API, external service, or external data source must log on
failure. This is a refactor target across the codebase. When indexing or scoping,
agents must flag any API or service call that lacks failure logging. Do not assume
it exists — verify it.

## Principle 4: Expect Incomplete Logging in Dev Apps
Development apps will have incomplete logging. Treat missing boundary logs as a
gap to surface explicitly — not an assumption to skip over. Surface it in the
output so developers know what needs to be added.

## Principle 5: Architecture Docs Are Navigation, Source Files Are Evidence
Architecture markdown files (high-level-architecture.md, scoped-architecture.md)
are human-friendly navigation aids. They help agents and developers locate things
quickly. They are never cited as evidence. When an agent needs to back up a claim,
it goes directly to the source file and cites file path + line number.

## Principle 6: Enforce Skill Ordering With Hard Stops
When a skill depends on another skill's output, enforce that dependency with a
hard stop and a clear error message. Never proceed silently when a prerequisite
is missing. Make the dependency chain explicit and visible to the user.

---

## Skill Dependency Chains

### Socrates
```
Index → Focus → Propose
```

---

## agentAI Directory Convention
All agent outputs live under `agentAI/<agent-name>/`. Other scripts and agents
expect inputs and outputs at these paths. Each agent maintains a `context.json`
as a lightweight state manifest — status and metadata only, not content. Content
lives in the markdown output files.
