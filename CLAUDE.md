# Jet Sharing — Team Conventions

## Stack
- **Backend**: Go 1.26, chi router (v5), PostgreSQL, slog logger
- **Mobile**: Flutter 3.29, Dart 3.6, Riverpod, go_router
- **API**: REST JSON, scooter IDs follow pattern sc-NNNN

## Go Backend

### Architecture
- All handlers follow: validate -> process -> respond pattern
- Use slog structured logger, NEVER fmt.Println or log.Printf
- Keep handlers under 50 LOC, extract business logic to service layer
- Error handling: wrap with fmt.Errorf("operation: %w", err)
- No unnecessary abstractions — simple is better than clever
- No interface unless 2+ implementations exist
- context.Context always first parameter

### Project Structure
- Handlers in handlers/ package (but for this demo, main.go is fine)
- Business logic in service/ package
- Keep main.go focused on wiring: routes, middleware, server startup

### Testing
- Table-driven tests with descriptive names
- Use httptest for handler tests
- Test names: Test_handleVerbNoun_condition (e.g., Test_handleUnlockScooter_batteryTooLow)

## Flutter Mobile App

### Architecture
- Feature-based folders: lib/screens/, lib/widgets/, lib/services/, lib/models/
- Services handle all API calls — screens never call http directly
- Riverpod for app-wide state; setState only for local UI state
- Use const constructors everywhere possible

### Naming
- Files: snake_case (scooter_map_screen.dart)
- Classes: PascalCase (ScooterMapScreen)
- Screens end with Screen (not Page, not View)

### Testing
- Widget tests in test/widgets/, unit tests in test/services/
- Use mocktail for mocking
- Run `flutter analyze && flutter test` before committing

## Common Mistakes to Avoid
- NEVER use global mutable state without proper locking (Go: sync.RWMutex)
- NEVER log sensitive data (payment tokens, user credentials)
- NEVER return raw internal errors to clients
- NEVER use print()/fmt.Println for logging — use slog (Go) or logger (Dart)
- NEVER store tokens in plain text — use flutter_secure_storage (Dart)
- NEVER call setState after async without checking mounted (Flutter)
- Always check battery level before allowing unlock (minimum 10%)
- Always validate scooter status before state transitions
