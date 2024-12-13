目录介绍：
pki：Public Key Infrastructure 公开密钥基础设施


第一步：# 一键部署kubernetes
chmod +x start.sh
./start.sh

等待启动成功浏览器访问https://172.21.28.14:9000, 使用账号:admin 密码:admin 登录dashboard。可以查看到集群部署的资源对象使用以及监控详情

注意：172.21.28.14换成部署kubernetes的linux虚机地址

第二部: # 如下步骤部署一个nginx, 等待创建成功，浏览器登录http://172.21.28.14:32500/ 可以访问到nginx首页
或者登录到linux所在节点执行kubectl apply -f ./example/nginx.yaml
1. 创建deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 80
2. 创建service
apiVersion: v1
kind: Service
metadata:
  name: ngx-service
  labels:
    app: nginx
spec:
  type: NodePort
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
    nodePort: 32535
