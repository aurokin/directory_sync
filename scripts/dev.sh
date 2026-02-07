#!/usr/bin/env bash
set -euo pipefail

# Run dsync from source via "go run".
#
# Examples:
#   scripts/dev.sh --help
#   scripts/dev.sh init
#   scripts/dev.sh doctor
#   scripts/dev.sh pull --link photos --dry-run

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "${ROOT_DIR}"
exec go run ./cmd/dsync "$@"
