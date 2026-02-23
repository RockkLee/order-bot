#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

ARTIFACT_DIR="${ARTIFACT_DIR:-build-artifacts}"
ARCHIVE_NAME="${ARCHIVE_NAME:-frontend-dist.tar.gz}"

if [[ "${1:-}" == "--clean" ]]; then
  rm -rf dist "$ARTIFACT_DIR"
fi

npm ci
npm run build-only

mkdir -p "$ARTIFACT_DIR"
tar -C dist -czf "$ARTIFACT_DIR/$ARCHIVE_NAME" .

cat <<MSG
Frontend build complete.
- Static files: $ROOT_DIR/dist
- Archive for upload: $ROOT_DIR/$ARTIFACT_DIR/$ARCHIVE_NAME

Upload example:
  aws s3 sync "$ROOT_DIR/dist" "s3://<your-bucket>" --delete
MSG
