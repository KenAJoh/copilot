#!/usr/bin/env bash
set -e
mise run generate
mise run check
mise run test
mise run build
