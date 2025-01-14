# 创建一个kind集群

## 安装kind

``` bash
# 下载kind
curl -Lo kind https://kind.sigs.k8s.io/dl/v0.11.1/kind-linux-amd64

# 移动kind到/usr/local/bin目录
sudo mv kind /usr/local/bin/kind
```

## 创建集群（--name 自定义集群名称）

``` bash
# 创建集群
kind create cluster
```

## 删除集群（--name 自定义集群名称）

``` bash
# 删除集群
kind delete cluster
```

## 创建完成，集群的context自动塞入到~/.kube/config文件中，可以使用kubectl命令进行操作