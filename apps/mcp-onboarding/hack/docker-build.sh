#!/usr/bin/env bash
set -e
VERSION=${VERSION:-$(git rev-parse --short HEAD)}
docker build -t ghcr.io/navikt/mcp-onboarding:$VERSION .
docker tag ghcr.io/navikt/mcp-onboarding:$VERSION ghcr.io/navikt/mcp-onboarding:latest
