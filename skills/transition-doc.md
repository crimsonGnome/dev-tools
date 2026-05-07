---
name: transition-doc
description: Write a session handoff to .context.md at end of session. Captures
  decisions made, open questions, and the next agent's first task.
---

At the end of this session, write a transition doc to `.context.md` at the repo root.

Structure it as follows:

## Session Summary
- Bullet list of decisions made this session

## Open Questions
- Bullet list of unresolved questions or blockers

## Next Agent
- Which agent alias to use next
- First task for that agent, stated clearly in one sentence

## Context
- Any other state the next agent needs to know to avoid repeating work

Write the file using the Write tool at path `.context.md`. Do not summarize — write the full doc.
