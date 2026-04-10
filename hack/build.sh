#!/usr/bin/env bash
set -e
for app in $APPS; do
  echo "📦 $app:" && (cd "apps/$app" && mise run build) && echo ""
done
echo "✅ All apps built successfully"
