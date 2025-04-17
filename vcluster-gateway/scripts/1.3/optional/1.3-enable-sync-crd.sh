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
ROLE_NAME="${VCLUSTER_ID}"
CLUSTERROLE_NAME="vc-${VCLUSTER_ID}-v-vcluster-${VCLUSTER_ID}"

# 检查命名空间是否存在
if ! kubectl get namespace "$NAMESPACE" > /dev/null 2>&1; then
  echo "Namespace '$NAMESPACE' does not exist. Exiting."
  exit 1
fi

echo "Namespace '$NAMESPACE' found."

# 检查 Deployment 是否存在
DEPLOYMENT_NAME="$VCLUSTER_ID"
if ! kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" > /dev/null 2>&1; then
  echo "Deployment '$DEPLOYMENT_NAME' not found in namespace '$NAMESPACE'. Exiting."
  exit 1
fi

echo "Deployment '$DEPLOYMENT_NAME' found in namespace '$NAMESPACE'."

# 首先获取当前的 CONFIG 环境变量索引
CONFIG_INDEX=$(kubectl get deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" -o json | \
  jq '.spec.template.spec.containers[0].env | map(.name == "CONFIG") | index(true)')

# 构建 Deployment Patch JSON
if [ "$CONFIG_INDEX" != "null" ]; then
  # CONFIG 环境变量存在，更新它
  DEPLOYMENT_PATCH_JSON="[
    {
      \"op\": \"replace\",
      \"path\": \"/spec/template/spec/containers/0/env/$CONFIG_INDEX\",
      \"value\": {
        \"name\": \"CONFIG\",
        \"value\": \"version: v1beta1\nexport:\n  - apiVersion: osm.datacanvas.com/v1alpha1\n    kind: ServiceExporter\n    patches:\n    - op: rewriteName\n      path: spec.serviceName\n    reversePatches:\n    - op: copyFromObject\n      fromPath: status\n      path: status\"
      }
    }
  ]"
else
  # CONFIG 环境变量不存在，添加它
  DEPLOYMENT_PATCH_JSON='[
    {
      "op": "add",
      "path": "/spec/template/spec/containers/0/env/-",
      "value": {
        "name": "CONFIG",
        "value": "version: v1beta1\nexport:\n  - apiVersion: osm.datacanvas.com/v1alpha1\n    kind: ServiceExporter\n    patches:\n    - op: rewriteName\n      path: spec.serviceName\n    reversePatches:\n    - op: copyFromObject\n      fromPath: status\n      path: status"
      }
    }
  ]'
fi

echo "Patching Deployment '$DEPLOYMENT_NAME' in namespace '$NAMESPACE'..."
kubectl patch deployment "$DEPLOYMENT_NAME" -n "$NAMESPACE" --type=json -p "$DEPLOYMENT_PATCH_JSON"
echo "Successfully patched Deployment '$DEPLOYMENT_NAME'."

# 更新 Role
echo "Patching Role '$ROLE_NAME' in namespace '$NAMESPACE'..."
ROLE_PATCH_JSON='{
  "apiGroups": ["osm.datacanvas.com"],
  "resources": ["serviceexporter", "serviceexporters"],
  "verbs": ["create", "delete", "patch", "update", "get", "list", "watch"]
}'

# 直接拼接 JSON Patch 数据
kubectl patch role "$ROLE_NAME" -n "$NAMESPACE" --type=json \
  -p "[{\"op\": \"add\", \"path\": \"/rules/-\", \"value\": $ROLE_PATCH_JSON}]"

echo "Successfully patched Role '$ROLE_NAME' in namespace '$NAMESPACE'."

# 更新 ClusterRole 一次性添加所有规则
echo "Patching ClusterRole '$CLUSTERROLE_NAME'..."

CLUSTERROLE_PATCH_JSON='[
  {
    "apiGroups": ["apiextensions.k8s.io"],
    "resources": ["customresourcedefinitions"],
    "verbs": ["create", "delete", "patch", "update", "get", "list", "watch"]
  },
  {
    "apiGroups": ["osm.datacanvas.com"],
    "resources": ["serviceexporter", "serviceexporters"],
    "verbs": ["create", "delete", "patch", "update", "get", "list", "watch"]
  },
  {
    "apiGroups": [""],
    "resources": ["services", "endpoints"],
    "verbs": ["get", "watch", "list"]
  }
]'

# 使用 JSON Patch 更新 ClusterRole (批量添加所有规则)
kubectl patch clusterrole "$CLUSTERROLE_NAME" --type=json \
  -p "$(echo "$CLUSTERROLE_PATCH_JSON" | jq -c '[.[] | {"op": "add", "path": "/rules/-", "value": .}]')"

echo "Successfully patched ClusterRole '$CLUSTERROLE_NAME'."

echo "All tasks completed successfully."