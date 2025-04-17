#!/usr/bin/env bash
set -e

# Usage: ./1.3-combined-script.sh <vcluster-id>

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <vcluster-id>"
  exit 1
fi

# vCluster ID
VCLUSTER_ID="$1"
NAMESPACE="vcluster-${VCLUSTER_ID}"
CLUSTERROLE_NAME="vc-${VCLUSTER_ID}-v-vcluster-${VCLUSTER_ID}"
QUOTA_NAME="${VCLUSTER_ID}-quota"
DEPLOYMENT_NAME="${VCLUSTER_ID}"
TARGET_CONTAINER="hooks"
NEW_VERSION="v0.8.1"

#######################################
# STEP 1: 为 ClusterRole 增加 namespaces 和 persistentvolumes 的权限
#######################################
echo "=== Step 1: Patching ClusterRole '${CLUSTERROLE_NAME}' ==="

# 检查 ClusterRole 是否存在
if ! kubectl get clusterrole "$CLUSTERROLE_NAME" > /dev/null 2>&1; then
  echo "ClusterRole '$CLUSTERROLE_NAME' does not exist. Exiting."
  exit 1
fi

# 要添加的权限规则
EXTRA_RULE='{
  "apiGroups": [""],
  "resources": ["namespaces", "persistentvolumes"],
  "verbs": ["create", "delete", "patch", "update", "get", "list", "watch"]
}'

# 使用 JSON Patch 增加权限
kubectl patch clusterrole "$CLUSTERROLE_NAME" --type=json \
  -p "[{\"op\": \"add\", \"path\": \"/rules/-\", \"value\": $EXTRA_RULE}]"

echo "Successfully patched ClusterRole '$CLUSTERROLE_NAME'."

#######################################
# STEP 2: 更新 ResourceQuota，移除 requests.storage
#######################################
echo "=== Step 2: Updating ResourceQuota in namespace '${NAMESPACE}' ==="

# 检查命名空间是否存在
if ! kubectl get namespace "$NAMESPACE" > /dev/null 2>&1; then
  echo "Namespace '$NAMESPACE' does not exist. Exiting."
  exit 1
fi

echo "Namespace '$NAMESPACE' found."

# 检查 ResourceQuota 是否存在
resource_quota=$(kubectl get resourcequota "$QUOTA_NAME" -n "$NAMESPACE" -o jsonpath='{.metadata.name}' || true)
if [[ -z "$resource_quota" ]]; then
    echo "ResourceQuota '$QUOTA_NAME' not found in namespace '$NAMESPACE', skipping ResourceQuota update..."
else
    echo "Found ResourceQuota '$QUOTA_NAME' in namespace '$NAMESPACE'. Checking if 'requests.storage' needs removal..."

    # 检查 requests.storage 是否存在
    has_storage=$(kubectl get resourcequota "$QUOTA_NAME" -n "$NAMESPACE" -o jsonpath='{.spec.hard.requests\.storage}' || true)
    if [[ -z "$has_storage" ]]; then
        echo "'requests.storage' not found in ResourceQuota '$QUOTA_NAME', no action needed."
    else
        echo "'requests.storage' found in ResourceQuota '$QUOTA_NAME', proceeding to remove it."
        # 创建 JSON Patch 数据
        patch='[{"op": "remove", "path": "/spec/hard/requests.storage"}]'
        # 应用 Patch
        kubectl patch resourcequota "$QUOTA_NAME" -n "$NAMESPACE" --type=json -p "$patch"
        echo "Successfully removed 'requests.storage' from ResourceQuota '$QUOTA_NAME'."
    fi
fi

#######################################
# STEP 3: 升级 hooks 容器镜像
#######################################
echo "=== Step 3: Updating hooks container image in Deployment '${DEPLOYMENT_NAME}' ==="

# 检查 Deployment 是否存在
if ! kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" &>/dev/null; then
    echo "Deployment '$DEPLOYMENT_NAME' not found in namespace '$NAMESPACE', skipping hooks upgrade..."
    exit 0
fi

# 获取 hooks 容器信息
container_info=$(kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" -o json | jq -c ".spec.template.spec.containers[] | select(.name == \"$TARGET_CONTAINER\")")

# 如果没有匹配的容器，直接退出
if [[ -z "$container_info" ]]; then
    echo "Container '$TARGET_CONTAINER' not found in Deployment '$DEPLOYMENT_NAME', skipping hooks upgrade..."
    exit 0
fi

current_image=$(echo "$container_info" | jq -r '.image')
image_prefix=${current_image%:*}
new_image="${image_prefix}:${NEW_VERSION}"

echo "Found container '$TARGET_CONTAINER' in Deployment '$DEPLOYMENT_NAME' with image '$current_image'"
echo "Updating image to '$new_image'"

# 更新镜像
kubectl set image deployment/"$DEPLOYMENT_NAME" -n "$NAMESPACE" "$TARGET_CONTAINER"="$new_image" --record
echo "Successfully updated image for container '$TARGET_CONTAINER' in Deployment '$DEPLOYMENT_NAME'."

echo "All steps completed successfully."
