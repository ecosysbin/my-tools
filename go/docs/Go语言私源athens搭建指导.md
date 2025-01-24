# Go私源athens搭建指导

## 1.目的
Go开发环境在有外网的情况下可以直接使用go官方的私有地址https://goproxy.cn拉取第三方包， 也可以通过国内阿里公司提供的开放源https://mirrors.aliyun.com/goproxy拉取第三方包。但是在内网环境，没有外网权限，则需要搭建自己的私源仓库，当前go社区使用最广泛的是通过微软公司提供的免费开源的athens项目搭建私源。

## 2.准备工作
1）. 准备一台可以连接外网的linux服务器（例如: ip为192.168.1.16），并配置go语言开发环境（Go开发环境搭建操作指南）
2）. 下载athens项目代码
cd  $GOPATH/src
git clone https://github.com/gomods/athens.git

## 3.安装athens
1）. 编译athens，编译完成会在$GOPATH/src/athens/目录下生成athens二进制文件
    cd $GOPATH/src/athens/
make
2）. 创建缓存路径
mkdir -p /opt/athens/data
3）. 创建过滤规则
vi ./filterFile
D
### 内网的gitlab不需要通过GlobalEndpoint下载
+ git.example.com
+ github.com/gomods/athens v0.1,v0.2,v0.4.1
4）. 修改配置
vi config.dev.toml
GlobalEndpoint = "https://mirrors.aliyun.com/goproxy"     // 配置代理仓库地址
FilterFile = "./filterFile"   // FilterFile 需要和GlobalEndpoint 同时配置，不然GlobalEndpoint 不会生效
StorageType = "disk"     // 配置存储类型为磁盘存储
Storage.Disk RootPath = "/opt/athens/data"  // 配置缓存存储路径
5）. 启动athens （3000是默认服务端口，可通过配置文件config.dev.toml修改）
./athens -config_file=./config.dev.toml

## 4.客户端配置
  配置客户端GOPROXY的地址为athens服务地址即可
  go env -w GOPROXY=http://192.168.1.16:3000    // ip地址为athens所在服务器ip, 3000为athens服务监听端口