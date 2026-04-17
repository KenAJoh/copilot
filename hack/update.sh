#!/usr/bin/env bash
set -e
for app in $APPS; do
  echo "📦 $app:"
  (cd "apps/$app" && mise run update)
  echo ""
done
echo "📦 scripts/generate-docs:"
(cd scripts/generate-docs && go get -u ./... && go mod tidy)
echo ""
echo "✅ All dependencies updated"
