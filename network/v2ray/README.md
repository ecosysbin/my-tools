# 背景介绍
我们要上网查询一些学习资料时，要下载一些镜像时，常碰到下载不下来，网页打开慢或者打不开的问题，就好像有一张无形的墙拦住了一样。这个文章讲一讲如何跨过这张无形的墙

# 解决方案
## 1. 科学上网
首先，你需要科学上网，科学上网可以让你上网更加顺畅，不受墙的限制。

## 2. 选择合适的镜像源
选择一个合适的镜像源，可以让你下载速度更快，并且可以解决下载不下来的问题。

## 3. 实在还是不能满足可以看看下面的方法
### 3.1 购买一台海外的云服务器（香港的就可以）
购买云服务器的步骤略过。
假设已经买好了服务器，本地通过xshell等ssh工具连接到了这台服务器。

### 3.1 下载安装v2ray
github地址：https://github.com/233boy/v2ray
安装：在云服务上执行如下指令一键安装
bash <(wget -qO- -o- https://git.io/v2ray.sh)
执行完成会有如下输出：
![alt text](image.png)

可以通过命令查询v2ray的运行状态：
systemctl status v2ray

### 3.2 配置客户端（以下以windows系统下安装v2rayN为例）
下载安装包：
wget https://github.com/2dust/v2rayN/releases/download/v3.27/v2rayN-v3.27.zip
解压：
unzip v2rayN-v3.27.zip
启动：
双击v2rayN.exe
配置：
![alt text](image2.png)

### 3.3 配置完成，浏览器设置代理即可正常使用了

### 3.4 要是公司内部使用可能要将v2rayN的端口映射到公司内网，这样就可以在公司内网下访问互联网了。这样可以在内网选择一台有图形界面的服务器，配置好v2rayN，然后通过内网访问互联网。其他服务器就可以将代理设置为这台机器的转发端口正常上网了。
设置代理的指令：
windows下：
set HTTP_PROXY=http://172.20.3.88:1088
set HTTPS_PROXY=http://172.20.3.88:1088
linux下：
export HTTP_PROXY=http://172.20.3.88:1088;export HTTPS_PROXY=http://172.20.3.88:1088;

### linux客户端代理
root@dev02-middle-01:~# ps -ef|grep v2
root     1572911       1  0 Jul03 ?        00:02:54 /usr/bin/containerd-shim-runc-v2 -namespace moby -id f9bd1ceedfffff5dcae9d7581c0da9251b22bbab3d74c2cd96642fd3b771c629 -address /run/containerd/containerd.sock
root     1572932 1572911  0 Jul03 ?        00:33:34 /usr/bin/v2ray -config /etc/v2ray/config.json
root     3901219 3900634  0 05:28 pts/0    00:00:00 grep --color=auto v2
root@dev02-middle-01:~# docker ps
CONTAINER ID   IMAGE                      COMMAND                  CREATED       STATUS       PORTS     NAMES
f9bd1ceedfff   v2fly/v2fly-core:v4.45.2   "/usr/bin/v2ray -con…"   13 days ago   Up 11 days             v2ray
root@dev02-middle-01:~# cat /etc/v2ray/config.json
{
  "log": {
    "loglevel": "warning"
  },
  "policy": {
    "levels": {
      "0": {
        "handshake": 1,
        "connIdle": 10,
        "uplinkOnly": 0,
        "downlinkOnly": 0,
        "bufferSize": 0
      }
    }
  },
  "dns": {
    "hosts": {
      "domain:k8s.io": "172.16.0.5",
      "domain:pkg.dev": "172.16.0.5",
      "geosite:gfw": "172.16.0.5"
    },
    "servers": [
      "https+local://223.5.5.5/dns-query",
      "https+local://223.6.6.6/dns-query",
      "https+local://1.12.12.12/dns-query",
      "https+local://120.53.53.53/dns-query"
    ],
    "queryStrategy": "UseIPv4"
  },
  "routing": {
    "domainStrategy": "AsIs",
    "domainMatcher": "mph",
    "rules": [
      {
        "type": "field",
        "inboundTag": [
          "dns-in"
        ],
        "outboundTag": "dns-out"
      },
      {
        "type": "field",
        "inboundTag": [
          "http-in",
          "https-in"
        ],
        "balancerTag": "balancer"
      }
    ],
    "balancers": [
      {
        "tag": "balancer",
        "selector": [
          "mfcloud"
        ],
        "strategy": {
          "type": "leastPing"
        }
      }
    ]
  },
  "observatory": {
    "subjectSelector": [
      "mfcloud"
    ],
    "probeURL": "https://www.gstatic.com/generate_204",
    "probeInterval": "10s"
  },
  "inbounds": [
    {
      "protocol": "dokodemo-door",
      "tag": "dns-in",
      "port": 53,
      "listen": "172.16.0.5",
      "settings": {
        "address": "114.114.114.114",
        "port": 53,
        "network": "udp"
      }
    },
    {
      "protocol": "dokodemo-door",
      "tag": "http-in",
      "port": 80,
      "listen": "172.16.0.5",
      "settings": {
        "address": "nxdomain",
        "port": 80
      },
      "sniffing": {
        "enabled": true,
        "destOverride": [
          "http"
        ]
      }
    },
    {
      "protocol": "dokodemo-door",
      "tag": "https-in",
      "port": 443,
      "listen": "172.16.0.5",
      "settings": {
        "address": "nxdomain",
        "port": 443
      },
      "sniffing": {
        "enabled": true,
        "destOverride": [
          "tls"
        ]
      }
    }
  ],
  "outbounds": [
    {
      "tag": "default",
      "protocol": "freedom",
      "settings": {
        "domainStrategy": "UseIPv4"
      }
    },
    {
      "tag": "dns-out",
      "protocol": "dns"
    },
    {
      "tag": "mfcloud-0",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn06.tomimser.xyz",
            "port": 9510,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "hk.v.3.az-1.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-1",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn07.tomimser.xyz",
            "port": 9520,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "hk.v.3.az-2.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-2",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn07.tomimser.xyz",
            "port": 9530,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "hk.v.3.az-3.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-3",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.v.3.hk-1.relay-gd-1.tomimser.xyz",
            "port": 32343,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "hk.dia-0.dmit.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-4",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.v.3.hk-1.relay-gd-1.tomimser.xyz",
            "port": 9550,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "hk.v.3.az-5.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-5",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.v.3.hk-1.relay-gd-1.tomimser.xyz",
            "port": 9560,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "hk.v.3.az-4.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-6",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn03.tomimser.xyz",
            "port": 25210,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato/ws",
          "headers": {
            "Host": "hk.v.3.hgc-1.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-7",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35501,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/718d45d67a753.video002.m3u8",
          "headers": {
            "Host": "hk.v.3.hgc-2.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-8",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35502,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/c11b-4f7b-9c58-078d15e269de.video002.m3u8",
          "headers": {
            "Host": "hk.v.3.hgc-3.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-9",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35503,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/098d-42df-94a4-2ccc088d92f9.video002.m3u8",
          "headers": {
            "Host": "hk.v.3.hgc-4.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-10",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35504,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/e00d-431d-bb67-4ceaa30dbef1.video002.m3u8",
          "headers": {
            "Host": "hk.v.3.hgc-5.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-11",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn07.tomimser.xyz",
            "port": 25611,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato/ws",
          "headers": {
            "Host": "hk.v.3.1.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-12",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn07.tomimser.xyz",
            "port": 26211,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato/ws",
          "headers": {
            "Host": "tw.hinet.01.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-13",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn07.tomimser.xyz",
            "port": 26212,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato/ws",
          "headers": {
            "Host": "tw.hinet.02.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-14",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn05.tomimser.xyz",
            "port": 26213,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato/ws",
          "headers": {
            "Host": "tw.hinet.03.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-15",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.v.3.hk-1.relay-gd-1.tomimser.xyz",
            "port": 26214,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato/ws",
          "headers": {
            "Host": "tw.hinet.02.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-16",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35521,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/a704-452c-8ad3-977b187067c6.live208.m3u8",
          "headers": {
            "Host": "tw.hinet.03.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-17",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn05.tomimser.xyz",
            "port": 30440,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "jp.v3-1.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-18",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn05.tomimser.xyz",
            "port": 30441,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "jp.v3-2.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-19",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn05.tomimser.xyz",
            "port": 30442,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "jp.v3-3.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-20",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn07.tomimser.xyz",
            "port": 30443,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "jp.v3-3.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-21",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.v.3.hk-1.relay-gd-1.tomimser.xyz",
            "port": 30444,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "jp.v3-4.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-22",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.v.3.hk-1.relay-gd-1.tomimser.xyz",
            "port": 30445,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "jp.v3-5.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-23",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35511,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/7016-43f0-96ef-46d73adc8cc8.live118.m3u8",
          "headers": {
            "Host": "jp.v3-4.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-24",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn06.tomimser.xyz",
            "port": 30450,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "us.v3-1.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-25",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn03.tomimser.xyz",
            "port": 30451,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "us.v3-2.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-26",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn05.tomimser.xyz",
            "port": 30452,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "us.v3-3.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-27",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn05.tomimser.xyz",
            "port": 30453,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "us.v3-4.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-28",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn07.tomimser.xyz",
            "port": 30455,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "us.v3-5.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-29",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35541,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/07fc-4596-9379-6c4ea64b097e.live108.m3u8",
          "headers": {
            "Host": "us.v3-4.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-30",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35542,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/07fc-4596-9379-6c4ea64b097e.live108.m3u8",
          "headers": {
            "Host": "us.v3-5.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-31",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn06.tomimser.xyz",
            "port": 31443,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "sg.v3-1.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-32",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn06.tomimser.xyz",
            "port": 31444,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "sg.v3-2.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-33",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.v.3.hk-1.relay-gd-1.tomimser.xyz",
            "port": 31445,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "sg.v3-3.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-34",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.v.3.hk-1.relay-gd-1.tomimser.xyz",
            "port": 31446,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "sg.v3-4.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-35",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35531,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/0ff4-432c-90ab-8050cb7260c1.live008.m3u8",
          "headers": {
            "Host": "sg.03.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-36",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35532,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/0ff4-432c-90ab-8050cb7260c1.live008.m3u8",
          "headers": {
            "Host": "sg.01.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-37",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "big-users.c.cdn07.tomimser.xyz",
            "port": 19010,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "kr.v.3.vu-1.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-38",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cnnode.cdn01.tomimser.xyz",
            "port": 35551,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/485c-4bf5-a1df-f64fb78d15e8.live001.m3u8",
          "headers": {
            "Host": "kr.03.tomimser.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-39",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.relay-v2.1.rssvt.xyz",
            "port": 50211,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "hk.v.2-1.rssvt.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-40",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.relay-v2.3.rssvt.xyz",
            "port": 50212,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "hk.v.2-2.rssvt.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-41",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.relay-v2.1.rssvt.xyz",
            "port": 50213,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato2.1",
          "headers": {
            "Host": "hk.v.2-3.rssvt.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-42",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.relay-v2.1.rssvt.xyz",
            "port": 50292,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "jp.v2-1.rssvt.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-43",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.relay-v2.1.rssvt.xyz",
            "port": 50201,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "jp.v2-2.rssvt.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-44",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.relay-v2.3.rssvt.xyz",
            "port": 50202,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "jp.v2-3.rssvt.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-45",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.relay-v2.2.rssvt.xyz",
            "port": 50301,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "us.v2-1.rssvt.xyz"
          }
        }
      }
    },
    {
      "tag": "mfcloud-46",
      "protocol": "vmess",
      "settings": {
        "vnext": [
          {
            "address": "cn.relay-v2.2.rssvt.xyz",
            "port": 50302,
            "users": [
              {
                "id": "bd3cb86c-af49-3484-8b38-d1f973b5050f",
                "alterId": 0,
                "level": 0
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/tomato",
          "headers": {
            "Host": "us.v2-3.rssvt.xyz"
          }
        }
      }
    }
  ]
}
root@dev02-middle-01:~# 