#!/usr/bin/env bash
set -e
echo "Validating allowlist.json..."

if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
  echo "Using existing server on port 8080..."
  if curl -f http://localhost:8080/v0.1/servers > /dev/null 2>&1; then
    echo "✅ Validation successful"
    exit 0
  else
    echo "❌ Validation failed"
    exit 1
  fi
fi

go run . &
SERVER_PID=$!
sleep 3
if curl -f http://localhost:8080/health > /dev/null 2>&1 && \
   curl -f http://localhost:8080/v0.1/servers > /dev/null 2>&1; then
  echo "✅ Validation successful"
  EXIT_CODE=0
else
  echo "❌ Validation failed"
  EXIT_CODE=1
fi
kill $SERVER_PID 2>/dev/null || true
exit $EXIT_CODE
