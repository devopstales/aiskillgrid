#!/usr/bin/env bash
set -u

# Cursor wrapper that delegates to the shared hook implementation.
exec ./.agents/hooks/refresh-indexes.sh
