# 查看镜像列表
ctr -n k8s.io image ls

# 拉远程镜像
ctr -n k8s.io image pull --user "mingfuyun:2gdeH661" registry.cn-beijing.aliyuncs.com/mfimg/nats:2.9.17-alpine

# 镜像打tag
ctr -n k8s.io image tag registry.cn-beijing.aliyuncs.com/mfimg/nats:2.9.17-alpine containercloud-mirror.xaidc.com/nats:2.9.17-alpine

# 导出镜像文件
ctr -n k8s.io image export nat.tar containercloud-mirror.xaidc.com/nats:2.9.17-alpine

# 导入镜像文件
ctr -n k8s.io image import nat.tar containercloud-mirror.xaidc.com/nats:2.9.17-alpine