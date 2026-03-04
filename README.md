# Jet Sharing

A fictional scooter-sharing platform used for the **AI-Native Engineering Workshop** series.

## What is this?

This is a sample codebase for practicing AI-assisted development with [Claude Code](https://docs.anthropic.com/en/docs/claude-code). It contains:

- **Go backend** (`main.go`) — REST API with chi router for scooter operations (unlock, lock, status)
- **Order service** (`internal/`) — Order lifecycle, payments, billing, scooter assignment, and task management
- **Flutter mobile app** (`lib/`) — Screens, widgets, and services for the rider app
- **Production logs** (`demos/`) — Synthetic order service logs for incident investigation exercises
- **Database schema** (`DATABASE.md`) — Full table definitions, status codes, and common queries
- **Sub-agents** (`.claude/agents/`) — Specialized Claude Code agents (e.g., incident investigator)
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

## Workshop 3: MCP & Sub-Agents

Workshop 3 focuses on connecting Claude Code to external systems via MCP and building specialized sub-agents.

### MCP Servers

Connect Linear (ticket tracking) and PostgreSQL (order database) as MCP servers:

```bash
# Linear (your tickets) — HTTP transport
claude mcp add --transport http linear-mcp \
  https://mcp.linear.app/mcp

# PostgreSQL (order database) — local process
claude mcp add orders-db -- npx -y @bytebase/dbhub \
  --dsn "postgresql://user:pass@host:port/dbname"

# Check what's connected
claude mcp list
```

### Sub-Agents

The `incident-investigator` sub-agent investigates production incidents using Linear tickets, database queries, log analysis, and source code tracing:

```bash
# List available agents
/agents

# Use the incident investigator
Investigate production incident Linear #JET-5
```

See `.claude/agents/incident-investigator.md` for the full agent definition.

### Order Service Source Code

The `internal/` directory contains Go source code that matches the `caller` field in production logs:

| Package | Files | Purpose |
|---------|-------|---------|
| `order` | `handler.go`, `lifecycle.go`, `status.go` | Order creation, status transitions, HTTP handlers |
| `payment` | `webhook.go` | Stripe webhook processing |
| `billing` | `handler.go` | Per-zone fare calculation |
| `scooter` | `assignment.go` | Nearest-scooter assignment |
| `task` | `manager.go` | Async task lifecycle |

## Previous Workshops

### Workshop 1 & 2 Exercises

**Code Explanation** — Pick any function and ask Claude Code to explain it:
```
Explain what handleUnlockScooter does, including edge cases and potential issues.
```

**Add a Feature (Plan Mode)** — Use `Shift+Tab` to review Claude's approach first:
```
Add a new endpoint POST /api/scooters/{id}/report that lets riders report a damaged scooter.
```

**Improve CLAUDE.md** — Run `/init` to see what Claude generates, then compare with the existing `CLAUDE.md`.

## Project Structure

```
jet-sharing/
  CLAUDE.md                # Team conventions (shared, checked into git)
  DATABASE.md              # Full database schema reference
  .claude/
    agents/
      incident-investigator.md  # Production incident investigation agent
    rules/
      go-backend.md        # Go-specific rules (paths: ["**/*.go"])
      flutter.md           # Flutter-specific rules (paths: ["**/*.dart"])
  main.go                  # Go backend (chi router, REST API)
  go.mod                   # Go module definition
  internal/
    order/                 # Order lifecycle and status management
    payment/               # Stripe webhook processing
    billing/               # Fare calculation
    scooter/               # Scooter assignment
    task/                  # Async task management
  demos/
    order-service-logs.txt # Production logs (~2000 lines)
  lib/                     # Flutter app
    screens/               # App screens
    widgets/               # Reusable widgets
    services/              # API and location services
  pubspec.yaml             # Flutter dependencies
```

## License

This project is for educational purposes as part of the AI-Native Engineering Workshop series.
