#!/bin/bash

# This script is used to install common tools on Linux
apt install dos2unix
dos2unix ./hack/local-up-volcano.sh

# 查看磁盘占用进行排行
du -sh /home/* 2>/dev/null | sort -rh | head -n 10