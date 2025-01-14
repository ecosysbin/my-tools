## 安装Kwok

```bash
# 下载最新版本的kwok（kwokctl可以模拟完整得一个集群，这边master通过kind部署，则kwokctl可以不用下载）
wget https://github.com/kwok-io/kwok/releases/latest/download/kwok-linux-amd64.tar.gz

# 解压
tar -zxvf kwok-linux-amd64.tar.gz

# 移动到/usr/local/bin目录
sudo mv kwok /usr/local/bin

```
## 启动Kwok

```bash
kwok   --kubeconfig=~/.kube/config   --manage-all-nodes=false   --manage-nodes-with-annotation-selector=kwok.x-k8s.io/node=fake   --manage-nodes-with-label-selector=   --manage-single-node=   --cidr=10.0.0.1/24   --node-ip=10.0.0.1   --node-lease-duration-seconds=40
```

## 接下来创建对应标签得node, 即可自动加入到kwok管理节点中
```
kubectl apply -f node-1.yaml
```
## node-1.yaml
```
apiVersion: v1
kind: Node
metadata:
  annotations:
    node.alpha.kubernetes.io/ttl: "0"
    kwok.x-k8s.io/node: fake
  labels:
    beta.kubernetes.io/arch: amd64
    beta.kubernetes.io/os: linux
    kubernetes.io/arch: amd64
    kubernetes.io/hostname: kwok-node-0
    kubernetes.io/os: linux
    kubernetes.io/role: agent
    node-role.kubernetes.io/agent: ""
    type: kwok
  name: kwok-node-1
spec:
  taints: # Avoid scheduling actual running pods to fake Node
  - effect: NoSchedule
    key: kwok.x-k8s.io/node
    value: fake
status:
  allocatable:
    cpu: 32
    memory: 256Gi
    pods: 110
  capacity:
    cpu: 32
    memory: 256Gi
    pods: 110
  nodeInfo:
    architecture: amd64
    bootID: ""
    containerRuntimeVersion: ""
    kernelVersion: ""
    kubeProxyVersion: fake
    kubeletVersion: fake
    machineID: ""
    operatingSystem: linux
    osImage: ""
    systemUUID: ""
  phase: Running
```