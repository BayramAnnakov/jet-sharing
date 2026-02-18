#!/usr/bin/env bash
set -e

# Kill any existing server on :8080
kill $(lsof -ti:8080) 2>/dev/null || true

echo "Starting Go backend on :8080..."
go run main.go &
GO_PID=$!

# Wait for the server to be ready
for i in {1..10}; do
  if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "Backend ready."
    break
  fi
  sleep 1
done

echo "Starting Flutter web app..."
flutter run -d chrome

# When Flutter exits, stop the Go server
kill $GO_PID 2>/dev/null
echo "Server stopped."
