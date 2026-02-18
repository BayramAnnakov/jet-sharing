# Jet Sharing — Team Conventions

## Stack
- **Backend**: Go 1.26, chi router (v5), in-memory store, slog logger
- **Frontend**: Flutter 3.41+, Dart 3.11+, Riverpod, go_router (runs as web app via `flutter run -d chrome`)
- **API**: REST JSON, scooter IDs follow pattern sc-NNNN

## Running the Project
Quick start (both backend + frontend):
```bash
./run.sh
```

Or manually:
1. Start the Go backend: `go run main.go` (serves on :8080)
2. Start the Flutter web app: `flutter run -d chrome`
3. The Flutter app connects to `http://localhost:8080/api`

`run.sh` starts the Go server, waits for it to be ready, then launches Flutter in Chrome. When you quit Flutter (`q`), it stops the server automatically.

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
- CORS middleware enabled for local development (allows Flutter web to call the API)

### Testing
- Table-driven tests with descriptive names
- Use httptest for handler tests
- Test names: Test_handleVerbNoun_condition (e.g., Test_handleUnlockScooter_batteryTooLow)

## Flutter App

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
