#!/bin/bash

# 查看模块列表
lsmod
# xor                    24576  2 async_xor,btrfs
# raid6_pq              122880  4 async_pq,btrfs,raid456,async_raid6_recov
# libcrc32c              16384  7 nf_conntrack,nf_nat,openvswitch,btrfs,nf_tables,raid456,ip_vs

# history： lsmod |grep vfio   lsmod|grep nouveau  lsmod|grep yr
# 会显示模块的名称、大小（以字节为单位）以及使用该模块的其他模块数量等信息

# 加载模块（必须自己写一个）
insmod my_driver.ko
# .ko是 Linux 内核模块文件的扩展名。insmod命令相对比较简单，它不会自动处理模块的依赖关系。这意味着如果被加载的模块依赖于其他模块，需要先手动加载那些依赖模块

modprobe module_name
# modprobe命令会自动处理模块的依赖关系，并加载所有依赖的模块。如果模块不存在，则会自动下载并安装。其中module_name是要加载的模块名称，不需要带.ko扩展名。


# 查看内核日志
dmesg
# 显示内核的日志信息，包括系统启动信息、硬件错误、模块加载信息等。