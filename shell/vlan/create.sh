#!/bin/bash

#!/bin/bash

# 创建虚拟网桥并启用
sudo brctl addbr br0
sudo ip link set br0 up

# 创建VLAN 10设备并配置相关参数
sudo vconfig add br0 10
sudo ip addr add 192.168.10.1/24 dev br0.10
sudo ip link set br0.10 up

# 创建VLAN 20设备并配置相关参数
sudo vconfig add br0 20
sudo ip addr add 192.168.20.1/24 dev br0.20
sudo ip link set br0.20 up

echo "VLAN设备创建及配置完成，可进行通信测试以验证隔离效果"