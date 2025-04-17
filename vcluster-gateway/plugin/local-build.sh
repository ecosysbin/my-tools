#!/bin/bash

# Exit on any error
set -euo pipefail

# Variables
CONTEXT="alaya.dev"
DEPLOYMENT_NAME="vc5r3tpuzosa"
CONTAINER_NAME="hooks"
REPO="harbor.zetyun.cn/aidc/vcluster/vcluster-plugin"
PLATFORM="linux/amd64"
VERSION=${VERSION:-$(date +"%Y%m%d%H%M%S")} # Default to current timestamp

# Print environment
echo "Building for PLATFORM: ${PLATFORM}"
echo "Using VERSION: ${VERSION}"
echo "Target Deployment: ${DEPLOYMENT_NAME} (Container: ${CONTAINER_NAME}) in Context: ${CONTEXT}"

# Step 1: Build the Go binary
echo "Building Go binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o plugin main.go
echo "Go binary built successfully."

# Step 2: Build the Docker image
echo "Building Docker image..."
docker build \
    --platform "${PLATFORM}" \
    -t "${REPO}:${VERSION}" \
    -f Dockerfile-local .
echo "Docker image built successfully: ${REPO}:${VERSION}"

# Step 3: Push the Docker image
echo "Pushing Docker image..."
docker push "${REPO}:${VERSION}"
echo "Docker image pushed successfully: ${REPO}:${VERSION}"

# Step 4: Update the Deployment
echo "Updating Deployment ${DEPLOYMENT_NAME} in context ${CONTEXT}..."
kubectl set image -n vcluster-${DEPLOYMENT_NAME} deployment/${DEPLOYMENT_NAME} \
  ${CONTAINER_NAME}=${REPO}:${VERSION} \
  --context=${CONTEXT}
echo "Deployment ${DEPLOYMENT_NAME} updated successfully."


echo "Build, push, and deployment update process completed!"
