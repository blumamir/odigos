#!/bin/bash

set -e

GENERIC_SCRIPT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/kind-simple-demo.sh"

# Define matrix of K8s versions and corresponding kind images
VERSIONS_AND_IMAGES=(
  "v1.20.15 kindest/node:v1.20.15@sha256:a32bf55309294120616886b5338f95dd98a2f7231519c7dedcec32ba29699394"
  "v1.23.17 kindest/node:v1.23.17@sha256:14d0a9a892b943866d7e6be119a06871291c517d279aedb816a4b4bc0ec0a5b3"
  "v1.32.0 kindest/node:v1.32.0@sha256:2458b423d635d7b01637cac2d6de7e1c1dca1148a2ba2e90975e214ca849e7cb"
)

for entry in "${VERSIONS_AND_IMAGES[@]}"; do
  read -r VERSION IMAGE <<< "$entry"
  echo "=== Building for $VERSION ==="
  $GENERIC_SCRIPT "$VERSION" "$IMAGE"
done
