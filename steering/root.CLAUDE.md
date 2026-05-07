# Root Steering

This is the monorepo root for the Vigilon/Ignisight project.

## Packages
- `vigilonCDK/` — CDK infrastructure (AWS deployments)
- `vigilonAPIHandler/` — Lambda function handlers
- `front-end/` — React Native mobile app (Vigilance)

## Agent System
- Agents live in `dev-tools/agents/`
- Skills live in `dev-tools/skills/`
- Switch agents by quitting the session, writing `.context.md`, and running the next alias
- Never carry context across agent boundaries — write a transition doc first

## Steering Guide
- Agent and skill design principles live in `dev-tools/steering-guide.md`
- All agents must follow the principles defined there

## General Rules
- Do not modify `.context.md` directly unless running the transition-doc skill
- Keep commits scoped to a single package where possible
- Check `.context.md` at session start for prior agent handoff
