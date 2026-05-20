#!/bin/bash
set -euo pipefail

# Build and push the Middleware MCP server image (multi-arch).
#
# Usage:   ./docker-build.sh <tag>
# Example: ./docker-build.sh v1.0.0
#
# Override the target image / platforms with env vars:
#   REGISTRY=ghcr.io/middleware-labs IMAGE=mcp-middleware \
#   PLATFORMS=linux/amd64,linux/arm64 ./docker-build.sh v1.0.0
#
# Set PUSH_LATEST=1 to also tag and push :latest.

if [ -z "${1:-}" ]; then
    echo "Usage: $0 <tag>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

TAG="$1"
REGISTRY="${REGISTRY:-ghcr.io/middleware-labs}"
IMAGE="${IMAGE:-mcp-middleware}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"
REF="${REGISTRY}/${IMAGE}:${TAG}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

TAGS=(-t "${REF}")
if [ "${PUSH_LATEST:-0}" = "1" ]; then
    TAGS+=(-t "${REGISTRY}/${IMAGE}:latest")
fi

echo "Building ${REF} for ${PLATFORMS} from context: $(pwd)"
docker buildx build \
  --platform "${PLATFORMS}" \
  -f ./Dockerfile . \
  "${TAGS[@]}" \
  --build-arg VERSION="${TAG}" \
  --push

echo "Pushed ${REF}"
[ "${PUSH_LATEST:-0}" = "1" ] && echo "Pushed ${REGISTRY}/${IMAGE}:latest"
exit 0
