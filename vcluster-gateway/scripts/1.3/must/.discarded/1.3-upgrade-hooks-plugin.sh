#!/bin/bash

set -e

# 检查参数是否传递
if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <vcluster-id>"
  exit 1
fi

# vCluster ID
VCLUSTER_ID=$1
NAMESPACE="vcluster-${VCLUSTER_ID}"
DEPLOYMENT_NAME="${VCLUSTER_ID}"
TARGET_CONTAINER="hooks"
NEW_VERSION="v0.8.1"

printf "Processing Deployment: %s in Namespace: %s\n" "$DEPLOYMENT_NAME" "$NAMESPACE"

# 检查命名空间是否存在
if ! kubectl get namespace "$NAMESPACE" > /dev/null 2>&1; then
  echo "Namespace '$NAMESPACE' does not exist. Exiting."
  exit 1
fi

# 检查 Deployment 是否存在
if ! kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" &>/dev/null; then
    printf "Deployment '%s' not found in namespace '%s', exiting...\n" "$DEPLOYMENT_NAME" "$NAMESPACE" >&2
    exit 1
fi

# 使用 jq 直接处理目标 Deployment
container_info=$(kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" -o json | jq -c ".spec.template.spec.containers[] | select(.name == \"$TARGET_CONTAINER\")")

# 如果没有匹配的容器，直接退出
if [[ -z "$container_info" ]]; then
    printf "Container '%s' not found in Deployment '%s'. Exiting...\n" "$TARGET_CONTAINER" "$DEPLOYMENT_NAME" >&2
    exit 1
fi

# 提取当前镜像路径并更新
current_image=$(echo "$container_info" | jq -r '.image')
image_prefix=${current_image%:*}
new_image="${image_prefix}:${NEW_VERSION}"

printf "Found container '%s' in Deployment '%s' with image '%s'\n" "$TARGET_CONTAINER" "$DEPLOYMENT_NAME" "$current_image"
printf "Updating image to '%s'\n" "$new_image"

# 更新镜像
kubectl set image deployment/"$DEPLOYMENT_NAME" -n "$NAMESPACE" "$TARGET_CONTAINER"="$new_image" --record
printf "Successfully updated image for container '%s' in Deployment '%s'\n" "$TARGET_CONTAINER" "$DEPLOYMENT_NAME"
