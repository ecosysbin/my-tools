#!/bin/bash

# This script is used to install common tools on Linux
apt install dos2unix
dos2unix ./hack/local-up-volcano.sh

# 查看磁盘占用进行排行
du -sh /home/* 2>/dev/null | sort -rh | head -n 10

# 遍历查找包含关键字的文件，输出文件名称
find . -type f -exec grep -q "docker.io" {} \; -print

# 遍历查找包含关键字的文件，输出文件名称，并输出关键字内容行
grep -r --include="*" "docker.io" .

# 查找到文件替换为其他字符串
grep -rl "docker.io" . | xargs sed -i 's/docker.io/containercloud-mirror.xaidc.com/g'

# 替换字符串存在"/"字符冲突时，这样写
grep -rl "gcr.io" . | xargs sed -i 's#gcr.io#m.daocloud.io/gcr.io#g'

# 磁盘占用过高，查看根目录下的大文件（结果是vc使用openEBS pvc没有删除，限制是5GB,结果用了71G）
find . -type f -size +100M 2>/dev/null | xargs du -sh

# 查看Containerd镜像和缓存占用
du -sh /var/lib/containerd/*
# 查看容器日志占用（假设日志在/var/log/containers）
du -sh /var/log/containers/*
# 查看容器运行时数据占用
du -sh /var/lib/kubelet/*