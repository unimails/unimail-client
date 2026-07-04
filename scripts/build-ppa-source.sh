#!/usr/bin/env bash
set -euo pipefail

# Build a Debian source package suitable for Launchpad PPA upload.
# Usage: ./scripts/build-ppa-source.sh [ubuntu_series] [upstream_version]

SERIES="${1:-noble}"
UPSTREAM_VERSION="${2:-}"
MODIFY_VERSION="${3:-}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

PACKAGE_NAME="$(dpkg-parsechangelog -SSource)"
CURRENT_VERSION="$(dpkg-parsechangelog -SVersion)"
CURRENT_UPSTREAM="${CURRENT_VERSION%%-*}"

if [[ -n "${UPSTREAM_VERSION}" ]]; then
  export DEBFULLNAME="${DEBFULLNAME:-allcloud.top}"
  export DEBEMAIL="${DEBEMAIL:-admin@allcloud.top}"
  DEBIAN_VERSION="${UPSTREAM_VERSION}ppa${MODIFY_VERSION}~${SERIES}"
  dch --force-bad-version --distribution "${SERIES}" --newversion "${DEBIAN_VERSION}" "Automated PPA source build for ${SERIES}."
else
  DEBIAN_VERSION="${CURRENT_VERSION}"
  UPSTREAM_VERSION="${CURRENT_UPSTREAM}"
fi

# Ensure debian/rules is executable in CI environments.
chmod +x debian/rules

ORIG_TARBALL="../${PACKAGE_NAME}_${UPSTREAM_VERSION}.orig.tar.gz"
git archive --format=tar.gz --prefix="${PACKAGE_NAME}-${UPSTREAM_VERSION}/" -o "${ORIG_TARBALL}" HEAD

if [[ "${UNSIGNED:-0}" == "1" ]]; then
  dpkg-buildpackage -S -sa -d -us -uc
else
  dpkg-buildpackage -S -sa -d
fi

CHANGES_FILE="../${PACKAGE_NAME}_${DEBIAN_VERSION}_source.changes"
if [[ ! -f "${CHANGES_FILE}" ]]; then
  echo "Expected changes file not found: ${CHANGES_FILE}" >&2
  exit 1
fi

echo "Built source package successfully: ${CHANGES_FILE}"
