#!/usr/bin/env bash
# Run all stages, track failures, report at end.
# Individual scripts (generate.sh, check.sh, build.sh) use set -e internally.
failed=()

mise run generate || failed+=("generate")
mise run check    || failed+=("check")
mise run build    || failed+=("build")

if [[ ${#failed[@]} -gt 0 ]]; then
  echo ""
  echo "❌ Failed stages: ${failed[*]}"
  exit 1
fi
echo "✅ All stages passed"
