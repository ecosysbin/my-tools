# 使用root用户启动容器
docker run -u root -itd -p 8080:8080 -p 50000:50000 --restart=on-failure jenkins:v5

# 将docker挂载进容器内，以便可以在容器内打包镜像
docker run -u root -itd -v /var/run/docker.sock:/var/run/docker.sock -v /usr/bin/docker:/usr/bin/docker -p 8080:8080 -p 50000:50000 --restart=on-failure jenkins:v5

# 以后台进程方式启动容器
docker run -u root -d -p 8080:8080 -p 50000:50000 --restart=on-failure jenkins:v5

# 容器的运行态修改，提交为新的镜像
docker commit <container_id> <image_name>
# 如：
```
docker commit b7941f7b4f94 jenkins:v5
```    

# 通过bash进入容器（通过sh进入容器时，会提示找不到命令，操作会受限）
docker exec -it <container_id> /bin/bash
或
docker exec -it <container_id> bash
# 如：
``` 
docker exec -it b7941f7b4f94 /bin/bash
```

# docker build
## 使用docker buildx build 构建多平台镜像, 为了解决容器内拉取go mod包网络错误问题，可以配置HTTPS_PROXY和HTTP_PROXY环境变量
```
docker buildx build -t "volcanosh/vc-controller-manager:79255a8dec66b3bb1a15d3e709240031a2ebffd8" . -f ./installer/dockerfile/controller-manager/Dockerfile --output=type="docker" --platform "linux/amd64" --build-arg APK_MIRROR= --build-arg OPEN_EULER_IMAGE_TAG=22.03-lts-sp2 --build-arg HTTP_PROXY=http://172.20.3.88:1088 --build-arg HTTPS_PROXY=http://172.20.3.88:1088;
```

# 清理所有未使用的镜像
docker image prune -a