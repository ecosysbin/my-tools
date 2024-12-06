#!/bin/bash
 
# 第一阶段: 加载BIOS/UEFI
echo "加载BIOS/UEFI"
 
# 第二阶段: 加载GRUB
echo "加载GRUB"
 
# 第三阶段: 加载内核
echo "加载内核"
 
# 内核开始执行自己的初始化过程
echo "内核初始化硬件"
echo "加载驱动程序"
echo "挂载根文件系统"
 
# 第四阶段: 用户空间初始化
echo "系统环境设置"
echo "网络设置"
echo "挂载其他文件系统"
echo "启动系统服务"
 
# 第五阶段: 登录
echo "等待用户登录"