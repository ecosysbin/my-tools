# 启动一个grpc服务同时是一个http服务

# 测试验证
## 启动grpc服务
```
root@demo:/home/goproject/src/grpchttp# go build -o test ./cmd/server/main.go 
root@demo:/home/goproject/src/grpchttp# ls
buf.gen.yaml  buf.yaml  cmd  gen  go.mod  go.sum  greet  grpct.exe  test
root@demo:/home/goproject/src/grpchttp# ./test 
hello grpc
```

## 访问http接口
```
(env) zetyun@demo:~$ curl --header "Content-Type: application/json" --data '{"name": "Jane"}'  http://localhost:8080/greet.v1.GreetService/Greet
{"greeting":"Hello, Jane!"}(env) zetyun@sd-k8s-master-1:~$ 
```

## 访问grpc接口(需要提前安装grpcurl工具，并在工程根目录下执行)
```
root@demo:/home/goproject/src/grpchttp# go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
root@demo:/home/goproject/src/grpchttp# grpcurl -protoset <(buf build -o -) -plaintext -d '{"name": "Jane"}' localhost:8080 greet.v1.GreetService/Greet
{
  "greeting": "Hello, Jane!"
}

```
# 参考文档：
https://connectrpc.com/docs/go/getting-started/