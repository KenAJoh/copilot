#!/usr/bin/env bash
failed=()

for app in $APPS; do
  echo "📦 $app:"
  if (cd "apps/$app" && mise run test); then
    echo ""
  else
    failed+=("$app")
    echo ""
  fi
done

echo "🧭 nav-pilot:"
if mise run nav-pilot:test; then
  echo ""
else
  failed+=("nav-pilot")
  echo ""
fi

if [[ ${#failed[@]} -gt 0 ]]; then
  echo "❌ Tests failed for: ${failed[*]}"
  exit 1
fi
echo "✅ All tests passed"
