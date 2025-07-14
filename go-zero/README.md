# 生成服务端代码
goctl api go -api ./greet.api -dir .
# 生成mysql客户端
goctl model mysql ddl --src user.sql --dir .