#!/bin/bash

# Usage: sh 1.3-update-resourcequota.sh <vcluster-id>

set -e

# 检查参数是否传递
if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <vcluster-id>"
  exit 1
fi

# vCluster ID
VCLUSTER_ID=$1
NAMESPACE="vcluster-${VCLUSTER_ID}"
QUOTA_NAME="${VCLUSTER_ID}-quota"

# 检查命名空间是否存在
if ! kubectl get namespace "$NAMESPACE" > /dev/null 2>&1; then
  echo "Namespace '$NAMESPACE' does not exist. Exiting."
  exit 1
fi

echo "Namespace '$NAMESPACE' found."

# 检查 ResourceQuota 是否存在
resource_quota=$(kubectl get resourcequota "$QUOTA_NAME" -n "$NAMESPACE" -o jsonpath='{.metadata.name}' || true)
if [[ -z "$resource_quota" ]]; then
    echo "ResourceQuota '$QUOTA_NAME' not found in namespace '$NAMESPACE', exiting..."
    exit 1
fi

echo "Found ResourceQuota '$QUOTA_NAME' in namespace '$NAMESPACE'"

# 检查 requests.storage 是否存在
has_storage=$(kubectl get resourcequota "$QUOTA_NAME" -n "$NAMESPACE" -o jsonpath='{.spec.hard.requests\.storage}' || true)
if [[ -z "$has_storage" ]]; then
    echo "'requests.storage' not found in ResourceQuota '$QUOTA_NAME', nothing to remove"
    exit 0
fi

echo "'requests.storage' found in ResourceQuota '$QUOTA_NAME', proceeding to remove it"

# 创建 JSON Patch 数据
patch='[{"op": "remove", "path": "/spec/hard/requests.storage"}]'

# 应用 Patch
echo "Patching ResourceQuota '$QUOTA_NAME' to remove 'requests.storage'"
kubectl patch resourcequota "$QUOTA_NAME" -n "$NAMESPACE" --type=json -p "$patch"

echo "Successfully removed 'requests.storage' from ResourceQuota '$QUOTA_NAME' in namespace '$NAMESPACE'"
