#!/usr/bin/env bash
set -euo pipefail

# Upload the latest source .changes file to Launchpad PPA.
# Usage: ./scripts/upload-ppa.sh <launchpad_user> <ppa_name> [changes_file]

if [[ $# -lt 2 ]]; then
  echo "Usage: $0 <launchpad_user> <ppa_name> [changes_file]" >&2
  exit 1
fi

LAUNCHPAD_USER="$1"
PPA_NAME="$2"
CHANGES_FILE="${3:-}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

if [[ -z "${CHANGES_FILE}" ]]; then
  CHANGES_FILE="$(ls -t ../*_source.changes | head -n 1 || true)"
fi

if [[ -z "${CHANGES_FILE}" || ! -f "${CHANGES_FILE}" ]]; then
  echo "No .changes file found. Build source package first." >&2
  exit 1
fi

dput "ppa:${LAUNCHPAD_USER}/${PPA_NAME}" "${CHANGES_FILE}"

echo "Uploaded: ${CHANGES_FILE}"
