#!/bin/bash

set -e

# 目标端口
TARGET_PORTS=(6060 8080 7777)

# 遍历所有以 vcluster- 开头的命名空间
for ns in $(kubectl get ns -o jsonpath='{.items[*].metadata.name}' | tr ' ' '\n' | grep '^vcluster-'); do
    # 提取 vcid
    vcid=${ns#vcluster-}
    echo "Checking namespace: $ns with vcid: $vcid"

    # 检查是否存在名为 "vcid-workloads" 的 NetworkPolicy
    policy=$(kubectl get networkpolicy -n "$ns" -o jsonpath="{.items[?(@.metadata.name=='$vcid-workloads')]}")
    if [[ -z "$policy" ]]; then
        echo "No NetworkPolicy named '$vcid-workloads' found in namespace $ns"
        echo
        continue
    fi

    echo "Found NetworkPolicy '$vcid-workloads' in namespace $ns"

    # 检查 NetworkPolicy 中是否已经有目标端口
    current_ports=$(kubectl get networkpolicy "$vcid-workloads" -n "$ns" -o jsonpath='{.spec.egress[*].ports[*].port}' | tr ' ' '\n')
    ports_to_add=()

    for port in "${TARGET_PORTS[@]}"; do
        if ! echo "$current_ports" | grep -q "^$port\$"; then
            ports_to_add+=("{\"port\":$port,\"protocol\":\"TCP\"}")
        fi
    done

    # 如果没有需要添加的端口，跳过
    if [[ ${#ports_to_add[@]} -eq 0 ]]; then
        echo "All target ports already exist in the NetworkPolicy"
        echo
        continue
    fi

    # 创建 JSON Patch 数据
    patch="["
    for ((i = 0; i < ${#ports_to_add[@]}; i++)); do
        if [[ $i -gt 0 ]]; then
            patch+=","
        fi
        patch+="{\"op\": \"add\", \"path\": \"/spec/egress/2/ports/-\", \"value\": ${ports_to_add[i]}}"
    done
    patch+="]"

    # 应用 Patch
    echo "Patching NetworkPolicy '$vcid-workloads' with the following ports: ${ports_to_add[*]}"
    kubectl patch networkpolicy "$vcid-workloads" -n "$ns" --type=json -p "$patch"

    echo "Successfully patched NetworkPolicy '$vcid-workloads' in namespace $ns"
    echo
done
