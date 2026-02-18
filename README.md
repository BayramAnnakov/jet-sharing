# Jet Sharing

A fictional scooter-sharing platform used for the **AI-Native Engineering Workshop** series.

## What is this?

This is a sample codebase for practicing AI-assisted development with [Claude Code](https://docs.anthropic.com/en/docs/claude-code). It contains:

- **Go backend** (`main.go`) — REST API with chi router for scooter operations (unlock, lock, status)
- **Flutter mobile app** (`lib/`) — Screens, widgets, and services for the rider app
- **CLAUDE.md** — Team conventions and coding standards (Claude Jet's "onboarding manual")
- **Path-scoped rules** (`.claude/rules/`) — Language-specific rules for Go and Flutter

## Getting Started

### Prerequisites

- [Claude Code](https://docs.anthropic.com/en/docs/claude-code) installed and authenticated
- Go 1.26+ (for backend exercises)
- Flutter 3.29+ / Dart 3.6+ (for mobile exercises)

### Setup

```bash
git clone https://github.com/BayramAnnakov/jet-sharing.git
cd jet-sharing
claude
```

That's it. Claude Code will automatically read `CLAUDE.md` and `.claude/rules/` when you start a session.

## Workshop Exercises

### Exercise 1: Code Explanation
Pick any function and ask Claude Code to explain it:
```
Explain what handleUnlockScooter does, including edge cases and potential issues.
```

### Exercise 2: Technical Questions
Ask context-aware questions about the codebase:
```
Why does the scooter map use a sync.RWMutex instead of a regular Mutex?
```

### Exercise 3: Add a Feature (with Plan Mode)
Use `Shift+Tab` or type `/plan` to review Claude Jet's approach before it writes code:
```
Add a new endpoint POST /api/v1/scooters/{id}/report that lets riders report a damaged scooter. Include the handler, request/response types, and input validation. Add it to the existing router.
```

After Claude Jet generates the code, review what it got **right** (chi router, slog, error wrapping) and what it **missed** (input validation? tests? battery check?). These gaps are what make CLAUDE.md a living document.

**Key takeaway**: Add verification commands to your CLAUDE.md. Instead of just "write tests," add: *"After implementing any handler, run `go test ./...` and ensure all tests pass."* This lets Claude Jet check its own work — the #1 best practice from Anthropic's engineering teams.

### Exercise 4: Improve CLAUDE.md
Run `/init` to see what Claude generates, then compare with the existing `CLAUDE.md`. What's missing? What's wrong?

## Project Structure

```
jet-sharing/
  CLAUDE.md              # Team conventions (shared, checked into git)
  .claude/
    rules/
      go-backend.md      # Go-specific rules (paths: ["**/*.go"])
      flutter.md         # Flutter-specific rules (paths: ["**/*.dart"])
  main.go                # Go backend (chi router, REST API)
  go.mod                 # Go module definition
  lib/                   # Flutter app
    screens/             # App screens
    widgets/             # Reusable widgets
    services/            # API and location services
  pubspec.yaml           # Flutter dependencies
```

## License

This project is for educational purposes as part of the AI-Native Engineering Workshop series.
