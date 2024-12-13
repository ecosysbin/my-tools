#!/bin/bash

# 关闭扇区
swapoff -a
sed -i 's/\/swap/#\/swap/g' /etc/fstab


# 安装软件包
dpkg -i ./deb-soft/*.deb

# 加载镜像
docker load -i ./images/calicocni.tar
docker load -i ./images/calico-kube-controller.tar
docker load -i ./images/caliconode.tar
docker load -i ./images/calicopoddaemon.tar
docker load -i ./images/coredns.tar
docker load -i ./images/dashboard.tar
docker load -i ./images/etcd.tar
docker load -i ./images/kube-apiserver.tar
docker load -i ./images/kube-contro.tar
docker load -i ./images/kubeproxy.tar
docker load -i ./images/kube-schecduler.tar
docker load -i ./images/metric-scraper.tar
docker load -i ./images/metricserver.tar
docker load -i ./images/pause.tar
docker load -i ./images/pod2daemon-flexvol-v3.21.1.tar
docker load -i ./images/prometheus-v2.36.2.tar
docker load -i ./images/registry-2.8.1.tar
docker load -i ./images/nginx.tar


# 创建集群凭证
mkdir -p /etc/kubernetes/pki/
cp ./pki/basic_auth_file  /etc/kubernetes/pki/basic_auth_file

# 增加default路由
ip route add default via 0.0.0.0 dev docker0 proto dhcp metric 100

# k8s集群初始化
./run.sh