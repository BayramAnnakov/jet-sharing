---
paths:
  - "lib/**/*.dart"
  - "test/**/*.dart"
---

# Flutter/Dart Conventions

## Architecture
- Feature-based folder structure: lib/screens/, lib/widgets/, lib/services/, lib/models/
- Services handle all API calls — screens never call http directly
- Use const constructors everywhere possible
- Stateless widgets by default; StatefulWidget only when local state is needed

## State Management
- Riverpod for app-wide state — no Provider, no setState for shared state
- Local UI state (loading, animation) is fine with setState
- Never put business logic in widgets — extract to services or notifiers

## Naming
- Files: snake_case (scooter_map_screen.dart)
- Classes: PascalCase (ScooterMapScreen)
- Screens end with Screen (not Page, not View)
- Widgets are extracted when >50 lines or reused

## Error Handling
- Custom exception classes per service (e.g., ScooterServiceException)
- Always show user-friendly errors — never raw exception messages
- Check `mounted` before calling setState after async gaps

## Testing
- Widget tests in test/widgets/, unit tests in test/services/
- Use `flutter_test` and `mocktail` for mocking
- Test both success and error states for every service method
- Run `flutter analyze && flutter test` before committing

## Do NOT
- NEVER use hardcoded strings — use app localization
- NEVER store tokens in plain text — use flutter_secure_storage
- NEVER call setState after async without checking mounted
- NEVER use print() for logging — use a proper logger (e.g., logger package)
