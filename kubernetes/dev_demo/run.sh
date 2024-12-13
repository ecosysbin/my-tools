#!/bin/bash

# k8s集群初始化
kubeadm init --config ./bluechart/kubeadm-init.yaml --v=5

# sleep 5s
sleep 5

# 创建kubectl凭证
mkdir -p $HOME/.kube
sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

hostname=`cat /etc/hostname`

# 去除污点
kubectl taint node $hostname node-role.kubernetes.io/master-

# 部署软件
kubectl apply -f ./bluechart/calico.yaml
kubectl apply -f ./bluechart/metricserver.yaml
kubectl apply -f ./bluechart/recommended.yaml