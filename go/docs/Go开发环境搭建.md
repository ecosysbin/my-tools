# Go开发环境搭建操作指南
下载并配置Go 开发语言环境 (go1.19.3)
## 下载地址
Go官网下载地址：https://golang.org/dl/
Go官方镜像站：https://golang.google.cn/dl/
如：wget https://golang.google.cn/dl/go1.23.3.linux-amd64.tar.gz
/home/zetyun/wb/go/bin
go version
开发者社区下载地址：https://studygolang.com/dl



## 环境配置
GOROOT是go 语言环境的安装路径，go的一些内置工具链通过GOROOT环境变量运行
GOPATH是一个环境变量，用来表明你写的go项目的存放路径（工作目录），GOPATH路径最好只设置一个，我们写的所有Go项目代码都放到GOPATH的src目录下。
GOPATH配置好以后，就可以通过go get 来安装一些基于Go lang的三方库，但国内用户无法直接访问部分库，可以通过配置国内源来解决这个问题。这里我们使用阿里云的源，在命令行输入下面两个命令：
go env -w GO111MODULE=on   // go 1.11 后默认打开，可通过go env查看
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy
Linux下GOROOT、GOPATH配置参考：
root@wangbin:/home/goproject/src/athens/cmd/proxy# cat /etc/profile.d/go_env
export GOROOT=/opt/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=/home/goproject
export PATH=$PATH:$GOPATH
root@wangbin:/home/goproject/src/athens/cmd/proxy#
环境变量生效： 
source /etc/profile.d/go_env

## 2.集成开发环境IDE, 推荐使用Visual Studio Code
   IDE下载地址：https://code.visualstudio.com/Download
   IDE安装插件：go   // Rich Go language support for Visual Studio Code
   安装方式：1. 如下在插件市场选择go install          
2. 在没有外网情况下可以进入插件vs code插件市场下载csvx格式的插件文件，通过导入方式安装
   插件市场：https://marketplace.visualstudio.com/items?itemName=golang.Go

## 3.安装开发工具链插件
1）.在有外网的情况下也可选择如下安装方式：
打开VSCode  ->  Ctrol + Shift + p  ->  Go: Install/Update Tools
选择所有的go插件，点击ok，等待安装完成
  2）.在没外网的情况下将如下开发插件放入$GOPATH目录下解压，重启vscode