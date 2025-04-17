#!/usr/bin/env bash
set -e

# 检查参数是否传递
if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <vcluster-id>"
  exit 1
fi

# vCluster ID
VCLUSTER_ID=$1
CLUSTERROLE_NAME="vc-${VCLUSTER_ID}-v-vcluster-${VCLUSTER_ID}"

# 检查 ClusterRole 是否存在
if ! kubectl get clusterrole "$CLUSTERROLE_NAME" > /dev/null 2>&1; then
  echo "ClusterRole '$CLUSTERROLE_NAME' does not exist. Exiting."
  exit 1
fi

echo "Found ClusterRole '$CLUSTERROLE_NAME'. Patching..."

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
