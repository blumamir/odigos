#!/bin/bash

set -e

K8S_VERSION=$1
KIND_IMAGE=$2
CLUSTER_NAME=simple-demo-${K8S_VERSION}

if [ -z "$K8S_VERSION" ] || [ -z "$KIND_IMAGE" ]; then
  echo "Usage: $0 <k8s-version> <kind-image>"
  echo "Example: $0 v1.20.15 kindest/node:v1.20.15@sha256:a32bf55309..."
  exit 1
fi

# Create cluster
kind create cluster \
  --name "$CLUSTER_NAME" \
  --image "$KIND_IMAGE"

# Images to pull and load into the cluster
IMAGES=(
  "registry.odigos.io/odigos-demo-inventory:v0.1"
  "registry.odigos.io/odigos-demo-membership:v0.1"
  "registry.odigos.io/odigos-demo-coupon:v0.1"
  "registry.odigos.io/odigos-demo-pricing:v0.1"
  "registry.odigos.io/odigos-demo-frontend:v0.2"
)

# Pull and load images
for IMAGE in "${IMAGES[@]}"; do
  docker pull "$IMAGE"
  kind load docker-image --name "$CLUSTER_NAME" "$IMAGE"
done

# Commit control-plane node to a new image
docker commit "${CLUSTER_NAME}-control-plane" keyval/kind-simple-demo:v0.2

# Delete the cluster
kind delete cluster --name "$CLUSTER_NAME"
