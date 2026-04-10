#!/usr/bin/env bash
set -e
for app in $APPS; do
  echo "📦 $app:" && (cd "apps/$app" && mise run test) && echo ""
done
echo "✅ All tests passed"
