# Scooter API — Team Conventions

## Architecture
- Go backend uses chi router (v5), not gin or mux
- All handlers follow: validate -> process -> respond pattern
- Use slog structured logger, NEVER fmt.Println or log.Printf
- Keep handlers under 50 LOC, extract business logic to service layer
- Error handling: wrap with fmt.Errorf("operation: %w", err)
- No unnecessary abstractions — simple is better than clever
- HTTP responses: always set Content-Type, use helper functions for JSON responses

## Code Style
- No interface unless 2+ implementations exist
- context.Context always first parameter
- Table-driven tests only
- Use named returns only when it improves readability
- Constants over magic numbers
- Group imports: stdlib, external, internal

## Project Structure
- Handlers in handlers/ package (but for this demo, main.go is fine)
- Business logic in service/ package
- Keep main.go focused on wiring: routes, middleware, server startup
- Config from environment variables, not hardcoded

## Testing
- Table-driven tests with descriptive names
- Use httptest for handler tests
- Test error paths, not just happy paths
- Test names: Test_handleVerbNoun_condition (e.g., Test_handleUnlockScooter_batteryTooLow)

## API Conventions
- All endpoints return JSON with Content-Type: application/json
- Error responses use format: {"error": "human-readable message"}
- Scooter IDs follow pattern: sc-NNNN (e.g., sc-1001)
- Status transitions: available -> in_use -> available (maintenance is manual)
- Always validate scooter status before state transitions

## Common Mistakes to Avoid
- NEVER use global mutable state for request handling without proper locking
- NEVER log sensitive data (payment tokens, user credentials)
- NEVER return raw internal errors to clients
- Always validate scooter ID format before database lookup
- Always check scooter status before unlock (prevent double-unlock)
- Always check battery level before allowing unlock (minimum 10%)
- Use sync.RWMutex — read lock for reads, write lock for mutations
