#!/usr/bin/env bash
set -e

MANIFEST="internal/discovery/copilot-manifest.json"
BEFORE=$(shasum -a 256 "$MANIFEST")

echo "Generating manifest..."
go run ./cmd/generate-manifest

AFTER=$(shasum -a 256 "$MANIFEST")

if [ "$BEFORE" = "$AFTER" ]; then
  echo "✅ Manifest is up to date"
else
  echo "❌ Manifest is out of date. Run 'mise generate' to update it."
  git --no-pager diff "$MANIFEST" || true
  git checkout "$MANIFEST"
  exit 1
fi
