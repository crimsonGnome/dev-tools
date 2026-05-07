# API Handler Package Steering

This package contains Lambda function handlers for the Vigilon platform.

## Scope
- Business logic for API endpoints
- No infrastructure definitions — those live in `vigilonCDK/`

## Rules
- Keep handlers thin — delegate logic to service modules
- All external calls (DB, third-party APIs) must be wrapped for error handling
- Do not commit secrets or credentials; use environment variables from the CDK stack
