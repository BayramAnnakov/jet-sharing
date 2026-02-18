---
paths:
  - "**/*.go"
---

# Go Backend Rules

## Handler Pattern
Every HTTP handler MUST follow this structure:
1. Parse and validate input
2. Call service layer
3. Return response

## Error Responses
Use consistent error response format:
- 400: validation errors (bad input from client)
- 404: resource not found
- 409: conflict (e.g., scooter already unlocked)
- 500: internal errors (log full error, return generic message to client)

## Naming
- Handlers: `handleVerbNoun` (e.g., `handleUnlockScooter`)
- Services: `VerbNoun` (e.g., `UnlockScooter`)
- Errors: `ErrNounState` (e.g., `ErrScooterNotFound`)

## Concurrency
- Use sync.RWMutex for in-memory state
- Always defer mu.Unlock() immediately after Lock()
- Prefer read locks (RLock) when not modifying state

## Logging
- Use slog with structured key-value pairs
- Always include resource ID in log context: slog.Info("action", "id", scooterID)
- Log at Warn level for expected client errors (not found, conflict)
- Log at Error level only for unexpected internal failures
- Never log request bodies containing user data
