#!/bin/bash

set -e

# === 参数和环境变量 ===

# 检查必需的环境变量是否已设置
REQUIRED_VARS=("VCLUSTER_ID" "IMAGE_NAME" "SIZE" "VOLUMENAME" "TENANT_ID" "APIKEY")

for var in "${REQUIRED_VARS[@]}"; do
  if [[ -z "${!var}" ]]; then
    echo "Error: Environment variable '$var' is not set."
    echo "Please set all required environment variables and retry."
    exit 1
  fi
done

# 设置变量
VCLUSTER_ID="$VCLUSTER_ID"
IMAGE_NAME="$IMAGE_NAME"
SIZE="$SIZE"
VOLUMENAME="$VOLUMENAME"
TENANT_ID="$TENANT_ID"
APIKEY="$APIKEY"

NAMESPACE="vcluster-${VCLUSTER_ID}"
DEPLOYMENT_NAME="${VCLUSTER_ID}"
TARGET_CONTAINER="syncer"
NEW_CONTAINER_NAME="multi-storage-requester"

# === 功能实现 ===

echo "=== 开始配置 Deployment '${DEPLOYMENT_NAME}' 在命名空间 '${NAMESPACE}' ==="

# 1. 检查命名空间是否存在
if ! kubectl get namespace "$NAMESPACE" > /dev/null 2>&1; then
  echo "Error: Namespace '$NAMESPACE' does not exist. Exiting."
  exit 1
fi
echo "Namespace '$NAMESPACE' found."

# 2. 检查 Deployment 是否存在
if ! kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" > /dev/null 2>&1; then
  echo "Error: Deployment '$DEPLOYMENT_NAME' not found in namespace '$NAMESPACE'. Exiting."
  exit 1
fi
echo "Deployment '$DEPLOYMENT_NAME' found in namespace '$NAMESPACE'."

# 3. 添加 Volume 'kubeconfig-shared'（如果尚未存在）
echo "Checking if volume 'kubeconfig-shared' exists in Deployment '$DEPLOYMENT_NAME'..."
volume_exists=$(kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" -o json | jq -r '.spec.template.spec.volumes[]? | select(.name=="kubeconfig-shared")')

if [[ -z "$volume_exists" ]]; then
  echo "Adding volume 'kubeconfig-shared' to Deployment '$DEPLOYMENT_NAME'..."
  kubectl patch deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" --type='json' -p='[
    {
      "op": "add",
      "path": "/spec/template/spec/volumes/-",
      "value": {
        "emptyDir": {},
        "name": "kubeconfig-shared"
      }
    }
  ]'
  echo "Volume 'kubeconfig-shared' added."
else
  echo "Volume 'kubeconfig-shared' already exists. Skipping addition."
fi

# 4. 在容器 'syncer' 中添加 VolumeMount（如果尚未存在）
echo "Checking if volumeMount for 'kubeconfig-shared' exists in container '$TARGET_CONTAINER'..."
volume_mount_exists=$(kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" -o json | \
  jq -r --arg TARGET_CONTAINER "$TARGET_CONTAINER" \
       '.spec.template.spec.containers[] | select(.name==$TARGET_CONTAINER) | .volumeMounts[]? | select(.name=="kubeconfig-shared")')

if [[ -z "$volume_mount_exists" ]]; then
  echo "Adding volumeMount to container '$TARGET_CONTAINER'..."
  kubectl patch deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" --type='json' -p='[
    {
      "op": "add",
      "path": "/spec/template/spec/containers/0/volumeMounts/-",
      "value": {
        "mountPath": "/root/.kube",
        "name": "kubeconfig-shared"
      }
    }
  ]'
  echo "VolumeMount added to container '$TARGET_CONTAINER'."
else
  echo "VolumeMount for 'kubeconfig-shared' already exists in container '$TARGET_CONTAINER'. Skipping addition."
fi

# 5. 添加新容器 'multi-storage-requester'（如果尚未存在）
echo "Checking if container '$NEW_CONTAINER_NAME' exists in Deployment '$DEPLOYMENT_NAME'..."
new_container_exists=$(kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" -o json | \
  jq -r '.spec.template.spec.containers[] | select(.name=="'"$NEW_CONTAINER_NAME"'")')

if [[ -z "$new_container_exists" ]]; then
  echo "Adding container '$NEW_CONTAINER_NAME' to Deployment '$DEPLOYMENT_NAME'..."

  # 使用 jq 构建 JSON Patch
  PATCH_JSON=$(jq -n \
    --arg name "$NEW_CONTAINER_NAME" \
    --arg image "$IMAGE_NAME" \
    --arg imagePullPolicy "Always" \
    --arg VCLUSTER_PLUGIN_ADDRESS "localhost:14001" \
    --arg VCLUSTER_PLUGIN_NAME "multi-storage-requester" \
    --arg NVIDIA_VISIBLE_DEVICES "none" \
    --arg SIZE "$SIZE" \
    --arg VOLUMENAME "$VOLUMENAME" \
    --arg TENANT_ID "$TENANT_ID" \
    --arg VCLUSTER_ID "$VCLUSTER_ID" \
    --arg APIKEY "$APIKEY" \
    '[
      {
        "op": "add",
        "path": "/spec/template/spec/containers/-",
        "value": {
          "name": $name,
          "image": $image,
          "imagePullPolicy": $imagePullPolicy,
          "env": [
            {
              "name": "VCLUSTER_PLUGIN_ADDRESS",
              "value": $VCLUSTER_PLUGIN_ADDRESS
            },
            {
              "name": "VCLUSTER_PLUGIN_NAME",
              "value": $VCLUSTER_PLUGIN_NAME
            },
            {
              "name": "NVIDIA_VISIBLE_DEVICES",
              "value": $NVIDIA_VISIBLE_DEVICES
            }
          ],
          "volumeMounts": [
            {
              "mountPath": "/root/.kube",
              "name": "kubeconfig-shared"
            }
          ],
          "args": [
            "--storage-types=capacity:\($SIZE):\($VOLUMENAME)",
            "--organization-id=\($TENANT_ID)",
            "--vcluster-id=\($VCLUSTER_ID)",
            "--apikey=\($APIKEY)"
          ]
        }
      }
    ]')

  # 输出生成的 PATCH_JSON，帮助调试
  echo "Generated PATCH_JSON: $PATCH_JSON"

  # 应用 Patch
  kubectl patch deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" --type='json' -p="$PATCH_JSON"

  echo "Container '$NEW_CONTAINER_NAME' added successfully."
else
  echo "Container '$NEW_CONTAINER_NAME' already exists in Deployment '$DEPLOYMENT_NAME'. Skipping addition."
fi

echo "=== Deployment '$DEPLOYMENT_NAME' 配置完成 ==="