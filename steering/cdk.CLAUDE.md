# CDK Package Steering

This package contains AWS CDK infrastructure for the Vigilon platform.

## Scope
- Infrastructure definitions only — no business logic
- Changes here affect deployed AWS resources; review carefully before deploying

## Rules
- Always run `cdk diff` before `cdk deploy`
- Keep stacks modular — one stack per concern
- Do not hardcode account IDs or region strings; use environment variables or context values
