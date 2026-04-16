#!/usr/bin/env bash
set -e
for app in $APPS_WITH_GENERATE; do
  echo "📦 $app:" && (cd "apps/$app" && mise run generate) && echo ""
done
echo "📄 docs:" && mise run docs:generate && echo ""
echo "✅ All files generated"
