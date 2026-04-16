#!/usr/bin/env bash
set -e
docker run --rm \
  -e GITHUB_APP_ID \
  -e GITHUB_APP_PRIVATE_KEY \
  -e GITHUB_APP_INSTALLATION_ID \
  -e GCP_TEAM_PROJECT_ID \
  copilot-metrics:local --run-once
