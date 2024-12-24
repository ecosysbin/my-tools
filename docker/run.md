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